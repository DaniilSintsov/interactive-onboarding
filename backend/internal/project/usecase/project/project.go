package project

import (
	"context"
	"errors"
	"fmt"
	"math"
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
		List(ctx context.Context, limit int32, offset int32) ([]entity.Project, int64, error)
		GetByID(ctx context.Context, projectID uuid.UUID) (entity.Project, error)
	}

	elementRepository interface {
		ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]entity.Element, error)
	}

	projectKeyGenerator interface {
		Generate() (string, error)
	}

	elementCreator interface {
		Create(ctx context.Context, params port.CreateElementParams) (entity.Element, error)
	}

	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) (err error)
	}
)

type projectService struct {
	projectRepository   projectRepository
	elementRepository   elementRepository
	projectKeyGenerator projectKeyGenerator
	elementCreator      elementCreator
	transactor          transactor
	logger              *zap.Logger
}

func NewProjectService(
	projectRepository projectRepository,
	elementRepository elementRepository,
	projectKeyGenerator projectKeyGenerator,
	elementCreator elementCreator,
	transactor transactor,
	logger *zap.Logger,
) *projectService {
	return &projectService{
		projectRepository:   projectRepository,
		elementRepository:   elementRepository,
		projectKeyGenerator: projectKeyGenerator,
		elementCreator:      elementCreator,
		transactor:          transactor,
		logger:              logger,
	}
}

const maxProjectKeyAttempts = 3

func (service *projectService) Create(
	ctx context.Context,
	params port.CreateProjectParams,
) (port.ProjectWithElements, error) {
	if err := params.Validate(); err != nil {
		return port.ProjectWithElements{}, fmt.Errorf("project usecase - create: validation error: %w", err)
	}

	project := entity.Project{
		Name: strings.TrimSpace(params.Name),
	}
	if err := project.Validate(); err != nil {
		return port.ProjectWithElements{}, fmt.Errorf("project usecase - create: validation error: %w", err)
	}

	var result port.ProjectWithElements
	var cycleErr error

	for i := 0; i < maxProjectKeyAttempts; i++ {
		projectKey, err := service.projectKeyGenerator.Generate()
		if err != nil {
			cycleErr = service.wrapCreateError(err, params.Name)
			continue
		}
		project.ProjectKey = projectKey

		err = service.transactor.WithTx(ctx, func(ctx context.Context) error {
			createdProject, err := service.projectRepository.Create(ctx, project)
			if err != nil {
				if errors.Is(err, errs.ErrProjectKeyAlreadyExists) {
					return err
				}
				return service.wrapCreateError(err, params.Name)
			}

			createdElements := make([]entity.Element, 0, len(params.Elements))
			for _, element := range params.Elements {
				createdElement, err := service.elementCreator.Create(ctx, port.CreateElementParams{
					ProjectID:   createdProject.ID,
					Key:         element.Key,
					Label:       element.Label,
					Description: element.Description,
				})
				if err != nil {
					return service.wrapCreateError(err, params.Name)
				}
				createdElements = append(createdElements, createdElement)
			}

			result = port.ProjectWithElements{
				Project:  createdProject,
				Elements: createdElements,
			}

			return nil
		})
		cycleErr = err
		if err != nil && errors.Is(err, errs.ErrProjectKeyAlreadyExists) {
			cycleErr = errs.ErrFailedGenerateUniqueKey
			continue
		}
		break
	}
	if cycleErr != nil {
		return port.ProjectWithElements{}, cycleErr
	}

	return result, nil
}

const (
	MaxLimit = 100
)

func (service *projectService) List(
	ctx context.Context,
	limit int,
	offset int,
) (port.ListProjectsResult, error) {
	if offset < 0 || offset > math.MaxInt32 {
		return port.ListProjectsResult{}, fmt.Errorf("project usecase - list: validation error: %w", errs.ErrOffsetInvalid)
	}
	if limit < 1 || limit > MaxLimit {
		return port.ListProjectsResult{}, fmt.Errorf("project usecase - list: validation error: %w", errs.ErrLimitInvalid)
	}

	list, total, err := service.projectRepository.List(ctx, int32(limit), int32(offset))
	if err != nil {
		return port.ListProjectsResult{}, service.wrapListError(err, limit, offset)
	}

	return port.ListProjectsResult{
		Projects: list,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}, nil
}

func (service *projectService) GetByID(
	ctx context.Context,
	projectID uuid.UUID,
) (port.ProjectWithElements, error) {
	if projectID == uuid.Nil {
		return port.ProjectWithElements{}, fmt.Errorf("project usecase - get by id: validation error: %w", errs.ErrProjectIDRequired)
	}

	var result port.ProjectWithElements

	err := service.transactor.WithTx(ctx, func(ctx context.Context) error {
		project, err := service.projectRepository.GetByID(ctx, projectID)
		if err != nil {
			if errors.Is(err, errs.ErrProjectNotFound) {
				return err
			}
			return service.wrapGetByIDError(err, projectID)
		}

		elements, err := service.elementRepository.ListByProjectID(ctx, project.ID)
		if err != nil {
			return service.wrapGetByIDError(err, projectID)
		}

		result = port.ProjectWithElements{
			Project:  project,
			Elements: elements,
		}

		return nil
	})
	if err != nil {
		return port.ProjectWithElements{}, err
	}

	return result, nil
}

func (service *projectService) Update(
	ctx context.Context,
	projectID uuid.UUID,
	name string,
) (port.ProjectWithElements, error) {
	if projectID == uuid.Nil {
		return port.ProjectWithElements{}, fmt.Errorf("project usecase - update: validation error: %w", errs.ErrProjectIDRequired)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return port.ProjectWithElements{},
			fmt.Errorf("project usecase - update: validation error: %w", errs.ErrProjectNameRequired)
	}
	if utf8.RuneCountInString(name) > entity.MaxProjectNameLength {
		return port.ProjectWithElements{},
			fmt.Errorf("project usecase - update: validation error: %w", errs.ErrProjectNameTooLong)
	}

	var result port.ProjectWithElements

	err := service.transactor.WithTx(ctx, func(ctx context.Context) error {
		updatedProject, err := service.projectRepository.Update(ctx, projectID, name)
		if err != nil {
			if errors.Is(err, errs.ErrProjectNotFound) {
				return err
			}
			return service.wrapUpdateError(err, projectID, name)
		}

		elements, err := service.elementRepository.ListByProjectID(ctx, updatedProject.ID)
		if err != nil {
			return service.wrapUpdateError(err, projectID, name)
		}

		result = port.ProjectWithElements{
			Project:  updatedProject,
			Elements: elements,
		}
		return nil
	})
	if err != nil {
		return port.ProjectWithElements{}, err
	}

	return result, nil
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
