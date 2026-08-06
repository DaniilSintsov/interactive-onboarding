package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/port"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/repository/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	elementRepository interface {
		ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]entity.Element, error)
		Create(ctx context.Context, element entity.Element) (entity.Element, error)
		Update(ctx context.Context, params model.UpdateElementParams) (entity.Element, error)
		Delete(ctx context.Context, projectID uuid.UUID, elementID uuid.UUID) error
	}

	projectRepository interface {
		GetByID(ctx context.Context, projectID uuid.UUID) (entity.Project, error)
	}

	elementUsageChecker interface {
		IsElementUsed(ctx context.Context, elementID uuid.UUID) (bool, error)
	}

	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) (err error)
	}
)

type elementService struct {
	elementRepository   elementRepository
	projectRepository   projectRepository
	elementUsageChecker elementUsageChecker
	transactor          transactor
	logger              *zap.Logger
}

func New(
	elementRepository elementRepository,
	projectRepository projectRepository,
	elementUsageChecker elementUsageChecker,
	transactor transactor,
	logger *zap.Logger,
) *elementService {
	return &elementService{
		elementRepository:   elementRepository,
		projectRepository:   projectRepository,
		elementUsageChecker: elementUsageChecker,
		transactor:          transactor,
		logger:              logger,
	}
}

func (e *elementService) List(
	ctx context.Context,
	projectID uuid.UUID,
) ([]entity.Element, error) {
	if projectID == uuid.Nil {
		return nil, fmt.Errorf("element usecase - list: %w", errs.ErrElementProjectIDRequired)
	}

	if _, err := e.projectRepository.GetByID(ctx, projectID); err != nil {
		return nil, fmt.Errorf("element usecase - list: %w", errs.ErrProjectNotFound)
	}

	elements, err := e.elementRepository.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, e.wrapListError(err, projectID)
	}

	return elements, nil
}

func (e *elementService) Create(
	ctx context.Context,
	params port.CreateElementParams,
) (entity.Element, error) {
	element := entity.Element{
		ProjectID:   params.ProjectID,
		Key:         params.Key,
		Label:       params.Label,
		Description: params.Description,
	}

	if err := element.Validate(); err != nil {
		return entity.Element{}, fmt.Errorf("element usecase - create: validation error: %w", err)
	}

	var resultElement *entity.Element

	err := e.transactor.WithTx(ctx, func(ctx context.Context) error {
		if _, err := e.projectRepository.GetByID(ctx, element.ProjectID); err != nil {
			return fmt.Errorf("element usecase - create: %w", errs.ErrProjectNotFound)
		}

		createdElement, err := e.elementRepository.Create(ctx, element)
		if err != nil {
			if errors.Is(err, errs.ErrElementKeyAlreadyExists) {
				return err
			}
			if errors.Is(err, errs.ErrProjectNotFound) {
				return err
			}
			return e.wrapCreateError(err, params)
		}

		*resultElement = createdElement

		return nil
	})
	if err != nil {
		return entity.Element{}, err
	}

	return *resultElement, nil
}

func (e *elementService) Update(
	ctx context.Context,
	params port.UpdateElementParams,
) (entity.Element, error) {
	if err := params.Validate(); err != nil {
		return entity.Element{}, fmt.Errorf("element usecase - update: validation error: %w", err)
	}

	if _, err := e.projectRepository.GetByID(ctx, params.ProjectID); err != nil {
		return entity.Element{}, fmt.Errorf("element usecase - update: %w", errs.ErrProjectNotFound)
	}

	updatedElement, err := e.elementRepository.Update(ctx, model.UpdateElementParams{
		ProjectID:   params.ProjectID,
		ElementID:   params.ElementID,
		Key:         params.Key,
		Label:       params.Label,
		Description: params.Description,
	})
	if err != nil {
		if errors.Is(err, errs.ErrElementKeyAlreadyExists) {
			return entity.Element{}, err
		}
		if errors.Is(err, errs.ErrElementNotFound) {
			return entity.Element{}, err
		}
		return entity.Element{}, e.wrapUpdateError(err, params)
	}

	return updatedElement, nil
}

func (e *elementService) Delete(
	ctx context.Context,
	projectID uuid.UUID,
	elementID uuid.UUID,
) error {
	if projectID == uuid.Nil {
		return fmt.Errorf("element usecase - delete: %w", errs.ErrElementProjectIDRequired)
	}

	if elementID == uuid.Nil {
		return fmt.Errorf("element usecase - delete: %w", errs.ErrElementIDRequired)
	}

	err := e.transactor.WithTx(ctx, func(ctx context.Context) error {
		if _, err := e.projectRepository.GetByID(ctx, projectID); err != nil {
			return fmt.Errorf("element usecase - delete: %w", errs.ErrProjectNotFound)
		}

		used, err := e.elementUsageChecker.IsElementUsed(ctx, elementID)
		if err != nil {
			return e.wrapDeleteError(err, projectID, elementID)
		}

		if used {
			return fmt.Errorf("element usecase - delete: %w", errs.ErrElementInUse)
		}

		err = e.elementRepository.Delete(ctx, projectID, elementID)
		if err != nil {
			if errors.Is(err, errs.ErrElementNotFound) {
				return err
			}
			return e.wrapDeleteError(err, projectID, elementID)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *elementService) wrapListError(err error, projectID uuid.UUID) error {
	e.logger.Error("element usecase - list failed",
		zap.String("projectID", projectID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("element usecase - list: projectID=%v: %w ", projectID, err)
}

func (e *elementService) wrapCreateError(err error, params port.CreateElementParams) error {
	e.logger.Error("element usecase - create failed",
		zap.String("projectID", params.ProjectID.String()),
		zap.String("key", params.Key),
		zap.String("label", params.Label),
		zap.String("description", params.Description),
		zap.Error(err),
	)

	return fmt.Errorf("element usecase - create: params=%v: %w", params, err)
}

func getStingFromPtr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func (e *elementService) wrapUpdateError(err error, params port.UpdateElementParams) error {
	e.logger.Error("element usecase - update failed",
		zap.String("projectID", params.ProjectID.String()),
		zap.String("elementID", params.ElementID.String()),
		zap.String("key", getStingFromPtr(params.Key)),
		zap.String("label", getStingFromPtr(params.Label)),
		zap.String("description", getStingFromPtr(params.Description)),
		zap.Error(err),
	)

	return fmt.Errorf("element usecase - update: project_id=%v element_id=%v: %w", params.ProjectID, params.ElementID, err)
}

func (e *elementService) wrapDeleteError(err error, projectID uuid.UUID, elementID uuid.UUID) error {
	e.logger.Error("element usecase - delete failed",
		zap.String("projectID", projectID.String()),
		zap.String("elementID", elementID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("element usecase - delete: projectID=%v element_id=%v: %w ", projectID, elementID, err)
}
