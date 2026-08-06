package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
		LockActive(ctx context.Context, projectID uuid.UUID, elementID uuid.UUID) error
	}

	projectRepository interface {
		LockActive(ctx context.Context, projectID uuid.UUID) error
	}

	elementUsageChecker interface {
		IsElementUsedBySteps(ctx context.Context, elementID uuid.UUID) (bool, error)
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

	resultElements := make([]entity.Element, 0)

	err := e.transactor.WithTx(ctx, func(ctx context.Context) error {
		if err := e.projectRepository.LockActive(ctx, projectID); err != nil {
			if errors.Is(err, errs.ErrProjectNotFound) {
				return errs.ErrProjectNotFound
			}
			return fmt.Errorf("element usecase - list: %w", err)
		}

		elements, err := e.elementRepository.ListByProjectID(ctx, projectID)
		if err != nil {
			return e.wrapListError(err, projectID)
		}

		resultElements = append(resultElements, elements...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resultElements, nil
}

func (e *elementService) Create(
	ctx context.Context,
	params port.CreateElementParams,
) (entity.Element, error) {
	element := entity.Element{
		ProjectID:   params.ProjectID,
		Key:         strings.TrimSpace(params.Key),
		Label:       strings.TrimSpace(params.Label),
		Description: strings.TrimSpace(params.Description),
	}

	if err := element.Validate(); err != nil {
		return entity.Element{}, fmt.Errorf("element usecase - create: validation error: %w", err)
	}

	var resultElement entity.Element

	err := e.transactor.WithTx(ctx, func(ctx context.Context) error {
		if err := e.projectRepository.LockActive(ctx, element.ProjectID); err != nil {
			if errors.Is(err, errs.ErrProjectNotFound) {
				return errs.ErrProjectNotFound
			}

			return fmt.Errorf("element usecase - create: %w", err)
		}

		createdElement, err := e.elementRepository.Create(ctx, element)
		if err != nil {
			if errors.Is(err, errs.ErrElementKeyAlreadyExists) ||
				errors.Is(err, errs.ErrProjectNotFound) {
				return err
			}
			return e.wrapCreateError(err, params)
		}

		resultElement = createdElement

		return nil
	})
	if err != nil {
		return entity.Element{}, err
	}

	return resultElement, nil
}

func (e *elementService) Update(
	ctx context.Context,
	params port.UpdateElementParams,
) (entity.Element, error) {
	if err := params.Validate(); err != nil {
		return entity.Element{}, fmt.Errorf("element usecase - update: validation error: %w", err)
	}

	var resultElement entity.Element

	err := e.transactor.WithTx(ctx, func(ctx context.Context) error {
		if err := e.projectRepository.LockActive(ctx, params.ProjectID); err != nil {
			if errors.Is(err, errs.ErrProjectNotFound) {
				return errs.ErrProjectNotFound
			}

			return fmt.Errorf("element usecase - update: %w", err)
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
				return err
			}
			if errors.Is(err, errs.ErrElementNotFound) {
				return err
			}
			return e.wrapUpdateError(err, params)
		}

		resultElement = updatedElement
		return nil
	})
	if err != nil {
		return entity.Element{}, err
	}

	return resultElement, nil
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
		if err := e.projectRepository.LockActive(ctx, projectID); err != nil {
			if errors.Is(err, errs.ErrProjectNotFound) {
				return errs.ErrProjectNotFound
			}

			return fmt.Errorf("element usecase - delete: %w", err)
		}

		if err := e.elementRepository.LockActive(ctx, projectID, elementID); err != nil {
			if errors.Is(err, errs.ErrElementNotFound) {
				return errs.ErrElementNotFound
			}

			return fmt.Errorf("element usecase - delete: %w", err)
		}

		used, err := e.elementUsageChecker.IsElementUsedBySteps(ctx, elementID)
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
		zap.String("project_id", projectID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("element usecase - list: project_id=%v: %w", projectID, err)
}

func (e *elementService) wrapCreateError(err error, params port.CreateElementParams) error {
	e.logger.Error("element usecase - create failed",
		zap.String("project_id", params.ProjectID.String()),
		zap.String("key", params.Key),
		zap.Error(err),
	)

	return fmt.Errorf("element usecase - create: project_id=%v key=%s: %w", params.ProjectID, params.Key, err)
}

func getStringFromPtr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func (e *elementService) wrapUpdateError(err error, params port.UpdateElementParams) error {
	e.logger.Error("element usecase - update failed",
		zap.String("project_id", params.ProjectID.String()),
		zap.String("element_id", params.ElementID.String()),
		zap.String("key", getStringFromPtr(params.Key)),
		zap.Error(err),
	)

	return fmt.Errorf("element usecase - update: project_id=%v element_id=%v: %w", params.ProjectID, params.ElementID, err)
}

func (e *elementService) wrapDeleteError(err error, projectID uuid.UUID, elementID uuid.UUID) error {
	e.logger.Error("element usecase - delete failed",
		zap.String("project_id", projectID.String()),
		zap.String("element_id", elementID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("element usecase - delete: project_id=%v element_id=%v: %w", projectID, elementID, err)
}
