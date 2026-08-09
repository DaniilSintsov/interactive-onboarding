package step

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	projecterrs "github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/errs"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	stepRepository interface {
		Create(ctx context.Context, step entity.Step) (entity.Step, error)
		Update(ctx context.Context, params port.UpdateStepParams) (entity.Step, error)
		Delete(ctx context.Context, scenarioID uuid.UUID, stepID uuid.UUID) error
		ListByScenarioID(ctx context.Context, scenarioID uuid.UUID) ([]entity.Step, error)
		GetNextNumber(ctx context.Context, scenarioID uuid.UUID) (int, error)
		LockActive(ctx context.Context, scenarioID uuid.UUID, stepID uuid.UUID) (entity.Step, error)
		Reorder(ctx context.Context, scenarioID uuid.UUID, orderedStepIDs []uuid.UUID) error
	}

	scenarioRepository interface {
		GetByID(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
		LockActive(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
	}

	elementLocker interface {
		LockActive(ctx context.Context, projectID uuid.UUID, elementID uuid.UUID) error
	}

	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) error
	}
)

type stepService struct {
	stepRepository     stepRepository
	scenarioRepository scenarioRepository
	elementLocker      elementLocker
	transactor         transactor
	logger             *zap.Logger
}

func NewStepService(
	stepRepository stepRepository,
	scenarioRepository scenarioRepository,
	elementLocker elementLocker,
	transactor transactor,
	logger *zap.Logger,
) *stepService {
	return &stepService{
		stepRepository:     stepRepository,
		scenarioRepository: scenarioRepository,
		elementLocker:      elementLocker,
		transactor:         transactor,
		logger:             logger,
	}
}

func (service *stepService) List(
	ctx context.Context,
	scenarioID uuid.UUID,
) ([]entity.Step, error) {
	if scenarioID == uuid.Nil {
		return nil, fmt.Errorf("step usecase - list: validation error: %w", errs.ErrStepScenarioIDRequired)
	}

	if _, err := service.scenarioRepository.GetByID(ctx, scenarioID); err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) {
			return nil, err
		}
		return nil, service.wrapListError(err, scenarioID)
	}

	steps, err := service.stepRepository.ListByScenarioID(ctx, scenarioID)
	if err != nil {
		return nil, service.wrapListError(err, scenarioID)
	}

	return steps, nil
}

func (service *stepService) Create(
	ctx context.Context,
	params port.CreateStepParams,
) (entity.Step, error) {
	step, err := stepFromCreateParams(params)
	if err != nil {
		return entity.Step{}, fmt.Errorf("step usecase - create: validation error: %w", err)
	}

	var result entity.Step

	err = service.transactor.WithTx(ctx, func(ctx context.Context) error {
		scenario, err := service.lockMutableScenario(ctx, step.ScenarioID)
		if err != nil {
			return err
		}

		if err := service.elementLocker.LockActive(ctx, scenario.ProjectID, step.ElementID); err != nil {
			return err
		}

		if params.StepNum == nil {
			step.StepNum, err = service.stepRepository.GetNextNumber(ctx, step.ScenarioID)
			if err != nil {
				return err
			}
		}

		createdStep, err := service.stepRepository.Create(ctx, step)
		if err != nil {
			return err
		}

		result = createdStep
		return nil
	})
	if err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) ||
			errors.Is(err, errs.ErrScenarioImmutable) ||
			errors.Is(err, errs.ErrInvalidStepNumber) ||
			errors.Is(err, errs.ErrStepNumberAlreadyExists) ||
			errors.Is(err, projecterrs.ErrElementNotFound) {
			return entity.Step{}, err
		}
		return entity.Step{}, service.wrapCreateError(err, params.ScenarioID)
	}

	return result, nil
}

func (service *stepService) Update(
	ctx context.Context,
	params port.UpdateStepParams,
) (entity.Step, error) {
	normalizedParams, err := normalizeAndValidateUpdateParams(params)
	if err != nil {
		return entity.Step{}, fmt.Errorf("step usecase - update: validation error: %w", err)
	}

	var result entity.Step

	err = service.transactor.WithTx(ctx, func(ctx context.Context) error {
		scenario, err := service.lockMutableScenario(ctx, normalizedParams.ScenarioID)
		if err != nil {
			return err
		}

		if _, err := service.stepRepository.LockActive(
			ctx,
			normalizedParams.ScenarioID,
			normalizedParams.StepID,
		); err != nil {
			return err
		}

		if normalizedParams.ElementID != nil {
			if err := service.elementLocker.LockActive(
				ctx,
				scenario.ProjectID,
				*normalizedParams.ElementID,
			); err != nil {
				return err
			}
		}

		updatedStep, err := service.stepRepository.Update(ctx, normalizedParams)
		if err != nil {
			return err
		}

		result = updatedStep
		return nil
	})
	if err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) ||
			errors.Is(err, errs.ErrScenarioImmutable) ||
			errors.Is(err, errs.ErrStepNotFound) ||
			errors.Is(err, projecterrs.ErrElementNotFound) {
			return entity.Step{}, err
		}
		return entity.Step{}, service.wrapUpdateError(err, params.ScenarioID, params.StepID)
	}

	return result, nil
}

func (service *stepService) Delete(
	ctx context.Context,
	scenarioID uuid.UUID,
	stepID uuid.UUID,
) error {
	if scenarioID == uuid.Nil {
		return fmt.Errorf("step usecase - delete: validation error: %w", errs.ErrStepScenarioIDRequired)
	}
	if stepID == uuid.Nil {
		return fmt.Errorf("step usecase - delete: validation error: %w", errs.ErrStepIDRequired)
	}

	err := service.transactor.WithTx(ctx, func(ctx context.Context) error {
		if _, err := service.lockMutableScenario(ctx, scenarioID); err != nil {
			return err
		}

		if err := service.stepRepository.Delete(ctx, scenarioID, stepID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) ||
			errors.Is(err, errs.ErrScenarioImmutable) ||
			errors.Is(err, errs.ErrStepNotFound) {
			return err
		}
		return service.wrapDeleteError(err, scenarioID, stepID)
	}

	return nil
}

func (service *stepService) Reorder(
	ctx context.Context,
	params port.ReorderStepsParams,
) ([]entity.Step, error) {
	if err := validateReorderParams(params); err != nil {
		return nil, fmt.Errorf("step usecase - reorder: validation error: %w", err)
	}

	var result []entity.Step

	err := service.transactor.WithTx(ctx, func(ctx context.Context) error {
		if _, err := service.lockMutableScenario(ctx, params.ScenarioID); err != nil {
			return err
		}

		if err := service.stepRepository.Reorder(ctx, params.ScenarioID, params.OrderedStepIDs); err != nil {
			return err
		}

		steps, err := service.stepRepository.ListByScenarioID(ctx, params.ScenarioID)
		if err != nil {
			return err
		}

		result = steps
		return nil
	})
	if err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) ||
			errors.Is(err, errs.ErrScenarioImmutable) ||
			errors.Is(err, errs.ErrInvalidStepNumber) ||
			errors.Is(err, errs.ErrStepDoesNotBelongToScenario) {
			return nil, err
		}
		return nil, service.wrapReorderError(err, params.ScenarioID)
	}

	return result, nil
}

func stepFromCreateParams(params port.CreateStepParams) (entity.Step, error) {
	stepNumber := 1
	if params.StepNum != nil {
		stepNumber = *params.StepNum
	}

	frontendData := params.FrontendData
	if len(frontendData) == 0 {
		frontendData = []byte("{}")
	}

	step := entity.Step{
		ScenarioID:   params.ScenarioID,
		ElementID:    params.ElementID,
		StepNum:      stepNumber,
		Title:        params.Title,
		Description:  params.Description,
		FrontendData: frontendData,
	}
	if err := step.Validate(); err != nil {
		return entity.Step{}, err
	}

	return step, nil
}

func normalizeAndValidateUpdateParams(params port.UpdateStepParams) (port.UpdateStepParams, error) {
	if params.ScenarioID == uuid.Nil {
		return port.UpdateStepParams{}, errs.ErrStepScenarioIDRequired
	}
	if params.StepID == uuid.Nil {
		return port.UpdateStepParams{}, errs.ErrStepIDRequired
	}
	if params.ElementID == nil && params.Title == nil && params.Description == nil && params.FrontendData == nil {
		return port.UpdateStepParams{}, errs.ErrStepUpdateFieldsRequired
	}

	if params.ElementID != nil && *params.ElementID == uuid.Nil {
		return port.UpdateStepParams{}, errs.ErrStepElementIDRequired
	}

	if params.Title != nil {
		title := strings.TrimSpace(*params.Title)
		if title == "" {
			return port.UpdateStepParams{}, errs.ErrStepTitleRequired
		}
		if utf8.RuneCountInString(title) > entity.MaxStepTitleLength {
			return port.UpdateStepParams{}, errs.ErrStepTitleTooLong
		}
		params.Title = &title
	}

	if params.Description != nil {
		description := strings.TrimSpace(*params.Description)
		if utf8.RuneCountInString(description) > entity.MaxStepDescriptionLength {
			return port.UpdateStepParams{}, errs.ErrStepDescriptionTooLong
		}
		params.Description = &description
	}

	return params, nil
}

func validateReorderParams(params port.ReorderStepsParams) error {
	if params.ScenarioID == uuid.Nil {
		return errs.ErrStepScenarioIDRequired
	}
	if len(params.OrderedStepIDs) == 0 {
		return errs.ErrStepOrderRequired
	}

	seen := make(map[uuid.UUID]struct{}, len(params.OrderedStepIDs))
	for _, stepID := range params.OrderedStepIDs {
		if stepID == uuid.Nil {
			return errs.ErrStepDoesNotBelongToScenario
		}
		if _, exists := seen[stepID]; exists {
			return errs.ErrStepDoesNotBelongToScenario
		}
		seen[stepID] = struct{}{}
	}

	return nil
}

func (service *stepService) lockMutableScenario(
	ctx context.Context,
	scenarioID uuid.UUID,
) (entity.Scenario, error) {
	scenario, err := service.scenarioRepository.LockActive(ctx, scenarioID)
	if err != nil {
		return entity.Scenario{}, err
	}
	if scenario.Status == entity.ScenarioStatusEnabled {
		return entity.Scenario{}, errs.ErrScenarioImmutable
	}

	return scenario, nil
}

func (service *stepService) wrapListError(err error, scenarioID uuid.UUID) error {
	service.logger.Error("step usecase - list failed",
		zap.String("scenario_id", scenarioID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("step usecase - list: scenario_id=%v: %w", scenarioID, err)
}

func (service *stepService) wrapCreateError(err error, scenarioID uuid.UUID) error {
	service.logger.Error("step usecase - create failed",
		zap.String("scenario_id", scenarioID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("step usecase - create: scenario_id=%v: %w", scenarioID, err)
}

func (service *stepService) wrapUpdateError(err error, scenarioID uuid.UUID, stepID uuid.UUID) error {
	service.logger.Error("step usecase - update failed",
		zap.String("scenario_id", scenarioID.String()),
		zap.String("step_id", stepID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("step usecase - update: scenario_id=%v step_id=%v: %w", scenarioID, stepID, err)
}

func (service *stepService) wrapDeleteError(err error, scenarioID uuid.UUID, stepID uuid.UUID) error {
	service.logger.Error("step usecase - delete failed",
		zap.String("scenario_id", scenarioID.String()),
		zap.String("step_id", stepID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("step usecase - delete: scenario_id=%v step_id=%v: %w", scenarioID, stepID, err)
}

func (service *stepService) wrapReorderError(err error, scenarioID uuid.UUID) error {
	service.logger.Error("step usecase - reorder failed",
		zap.String("scenario_id", scenarioID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("step usecase - reorder: scenario_id=%v: %w", scenarioID, err)
}
