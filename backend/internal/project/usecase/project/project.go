package project

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	projectRepository interface {
		Create(ctx context.Context, project entity.Project) (entity.Project, error)
		Update(ctx context.Context, projectID uuid.UUID, name string) (entity.Project, error)
		Delete(ctx context.Context, projectID uuid.UUID) error
		List(ctx context.Context, limit int32, offset int32) ([]entity.Project, error)
		Count(ctx context.Context) (int64, error)
		GetByID(ctx context.Context, projectID uuid.UUID) (entity.Project, error)
	}
)

type projectService struct {
	projectRepository projectRepository
	logger            *zap.Logger
}

func NewProjectService(
	projectRepository projectRepository,
	logger *zap.Logger,
) *projectService {
	return &projectService{
		projectRepository: projectRepository,
		logger:            logger,
	}
}

func (service *projectService) Create(
	ctx context.Context,
	name string,
	projectKey string,
) (entity.Project, error) {
	project := entity.Project{
		Name:       strings.TrimSpace(name),
		ProjectKey: strings.TrimSpace(projectKey),
	}
	if err := project.Validate(); err != nil {
		return entity.Project{}, fmt.Errorf("project usecase - create: validation error: %w", err)
	}
	createdProject, err := service.projectRepository.Create(ctx, project)
	if err != nil {
		if errors.Is(err, errs.ErrProjectKeyAlreadyExists) {
			return entity.Project{}, err
		}
		return entity.Project{}, service.wrapCreateError(err, name)
	}

	return createdProject, nil
}

const (
	MaxPageSize = 100
)

func (service *projectService) List(
	ctx context.Context,
	page int,
	pageSize int,
) (port.ListProjectsResult, error) {
	if page < 1 {
		return port.ListProjectsResult{}, fmt.Errorf("project usecase - list: validation error: %w", errs.ErrPageInvalid)
	}
	if pageSize < 1 || pageSize > MaxPageSize {
		return port.ListProjectsResult{}, fmt.Errorf("project usecase - list: validation error: %w", errs.ErrPageSizeInvalid)
	}

	limit := pageSize
	offset := (page - 1) * pageSize

	list, err := service.projectRepository.List(ctx, int32(limit), int32(offset))
	if err != nil {
		return port.ListProjectsResult{}, service.wrapListError(err, limit, offset)
	}

	count, err := service.projectRepository.Count(ctx)
	if err != nil {
		return port.ListProjectsResult{}, service.wrapListError(err, limit, offset)
	}

	return port.ListProjectsResult{
		Projects: list,
		Total:    count,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (service *projectService) GetByID(
	ctx context.Context,
	projectID uuid.UUID,
) (entity.Project, error) {
	if projectID == uuid.Nil {
		return entity.Project{}, fmt.Errorf("project usecase - get by id: validation error: %w", errs.ErrProjectIDRequired)
	}

	project, err := service.projectRepository.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, errs.ErrProjectNotFound) {
			return entity.Project{}, err
		}
		return entity.Project{}, service.wrapGetByIDError(err, projectID)
	}

	return project, nil
}

func (service *projectService) Update(
	ctx context.Context,
	projectID uuid.UUID,
	name string,
) (entity.Project, error) {
	if projectID == uuid.Nil {
		return entity.Project{}, fmt.Errorf("project usecase - update: validation error: %w", errs.ErrProjectIDRequired)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return entity.Project{},
			fmt.Errorf("project usecase - update: validation error: %w", errs.ErrProjectNameRequired)
	}
	if utf8.RuneCountInString(name) > 255 {
		return entity.Project{},
			fmt.Errorf("project usecase - update: validation error: %w", errs.ErrProjectNameTooLong)
	}

	updatedProject, err := service.projectRepository.Update(ctx, projectID, name)
	if err != nil {
		if errors.Is(err, errs.ErrProjectNotFound) {
			return entity.Project{}, err
		}
		return entity.Project{}, service.wrapUpdateError(err, projectID, name)
	}

	return updatedProject, nil
}

func (service *projectService) Delete(
	ctx context.Context,
	projectID uuid.UUID,
) error {
	if projectID == uuid.Nil {
		return fmt.Errorf("project usecase - delete: validation error: %w", errs.ErrProjectIDRequired)
	}

	err := service.projectRepository.Delete(ctx, projectID)
	if err != nil {
		if errors.Is(err, errs.ErrProjectNotFound) {
			return err
		}
		return service.wrapDeleteError(err, projectID)
	}

	return nil
}

func (service *projectService) wrapCreateError(err error, name string) error {
	service.logger.Error("project usecase - create failed",
		zap.String("name", name),
		zap.Error(err),
	)

	return fmt.Errorf("project usecase - create: name=%s: %w", name, err)
}

func (service *projectService) wrapListError(err error, limit int, offset int) error {
	service.logger.Error("project usecase - list failed",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.Error(err),
	)

	return fmt.Errorf("project usecase - list: limit=%d offset=%d: %w", limit, offset, err)
}

func (service *projectService) wrapGetByIDError(err error, projectID uuid.UUID) error {
	service.logger.Error("project usecase - get by id failed",
		zap.String("project_id", projectID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("project usecase - get by id: project_id=%v: %w", projectID, err)
}

func (service *projectService) wrapUpdateError(err error, projectID uuid.UUID, name string) error {
	service.logger.Error("project usecase - update failed",
		zap.String("project_id", projectID.String()),
		zap.String("name", name),
		zap.Error(err),
	)

	return fmt.Errorf("project usecase - update: project_id=%v name=%s: %w", projectID, name, err)
}

func (service *projectService) wrapDeleteError(err error, projectID uuid.UUID) error {
	service.logger.Error("project usecase - delete failed",
		zap.String("project_id", projectID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("project usecase - delete: project_id=%v: %w", projectID, err)
}
