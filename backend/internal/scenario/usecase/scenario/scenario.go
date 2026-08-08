package scenario

import (
	"context"
	"errors"
	"fmt"
	"math"
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
	scenarioRepository interface {
		Create(ctx context.Context, scenario entity.Scenario) (entity.Scenario, error)
		Update(ctx context.Context, params port.UpdateScenarioParams) (entity.Scenario, error)
		Delete(ctx context.Context, scenarioID uuid.UUID) error
		ListByProjectID(ctx context.Context, projectID uuid.UUID, status *entity.ScenarioStatus, limit, offset int32) ([]port.ScenarioSummary, int64, error)
		GetByID(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
		LockActive(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
		Publish(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
		Enable(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
		Disable(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
	}

	stepReader interface {
		ListByScenarioID(ctx context.Context, scenarioID uuid.UUID) ([]entity.Step, error)
	}

	projectLocker interface {
		LockActive(ctx context.Context, projectID uuid.UUID) error
	}

	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) error
	}
)

type scenarioService struct {
	scenarioRepository scenarioRepository
	stepReader         stepReader
	projectLocker      projectLocker
	transactor         transactor
	logger             *zap.Logger
}

func NewScenarioService(
	scenarioRepository scenarioRepository,
	stepReader stepReader,
	projectLocker projectLocker,
	transactor transactor,
	logger *zap.Logger,
) *scenarioService {
	return &scenarioService{
		scenarioRepository: scenarioRepository,
		stepReader:         stepReader,
		projectLocker:      projectLocker,
		transactor:         transactor,
		logger:             logger,
	}
}

func (service *scenarioService) Create(
	ctx context.Context,
	params port.CreateScenarioParams,
) (entity.Scenario, error) {
	scenario := entity.Scenario{
		ProjectID:   params.ProjectID,
		Name:        params.Name,
		Description: params.Description,
		PagePattern: params.PagePattern,
		Status:      entity.ScenarioStatusInDevelopment,
	}
	if err := scenario.Validate(); err != nil {
		return entity.Scenario{}, fmt.Errorf("scenario usecase - create: validation error: %w", err)
	}

	var result entity.Scenario

	err := service.transactor.WithTx(ctx, func(ctx context.Context) error {
		if err := service.projectLocker.LockActive(ctx, scenario.ProjectID); err != nil {
			if errors.Is(err, projecterrs.ErrProjectNotFound) {
				return errs.ErrProjectNotFound
			}
			return err
		}

		createdScenario, err := service.scenarioRepository.Create(ctx, scenario)
		if err != nil {
			return err
		}

		result = createdScenario
		return nil
	})
	if err != nil {
		if errors.Is(err, errs.ErrProjectNotFound) {
			return entity.Scenario{}, err
		}
		return entity.Scenario{}, service.wrapProjectOperationError("create", err, scenario.ProjectID)
	}

	return result, nil
}

const maxListLimit = 100

func (service *scenarioService) List(
	ctx context.Context,
	params port.ListScenariosParams,
) (port.ListScenariosResult, error) {
	if params.ProjectID == uuid.Nil {
		return port.ListScenariosResult{}, fmt.Errorf(
			"scenario usecase - list: validation error: %w",
			errs.ErrScenarioProjectIDRequired,
		)
	}
	if params.Status != nil && !params.Status.IsValid() {
		return port.ListScenariosResult{}, fmt.Errorf(
			"scenario usecase - list: validation error: %w",
			errs.ErrScenarioStatusUnknown,
		)
	}
	if params.Limit < 1 || params.Limit > maxListLimit {
		return port.ListScenariosResult{}, fmt.Errorf(
			"scenario usecase - list: validation error: %w",
			errs.ErrScenarioLimitInvalid,
		)
	}
	if params.Offset < 0 || params.Offset > math.MaxInt32 {
		return port.ListScenariosResult{}, fmt.Errorf(
			"scenario usecase - list: validation error: %w",
			errs.ErrScenarioOffsetInvalid,
		)
	}

	var result port.ListScenariosResult

	err := service.transactor.WithTx(ctx, func(ctx context.Context) error {
		if err := service.projectLocker.LockActive(ctx, params.ProjectID); err != nil {
			if errors.Is(err, projecterrs.ErrProjectNotFound) {
				return errs.ErrProjectNotFound
			}
			return err
		}

		scenarios, total, err := service.scenarioRepository.ListByProjectID(
			ctx,
			params.ProjectID,
			params.Status,
			int32(params.Limit),
			int32(params.Offset),
		)
		if err != nil {
			return err
		}

		result = port.ListScenariosResult{
			Scenarios: scenarios,
			Total:     total,
			Limit:     params.Limit,
			Offset:    params.Offset,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errs.ErrProjectNotFound) {
			return port.ListScenariosResult{}, err
		}
		return port.ListScenariosResult{}, service.wrapProjectOperationError("list", err, params.ProjectID)
	}

	return result, nil
}

func (service *scenarioService) GetByID(
	ctx context.Context,
	scenarioID uuid.UUID,
) (port.ScenarioWithSteps, error) {
	if scenarioID == uuid.Nil {
		return port.ScenarioWithSteps{}, fmt.Errorf(
			"scenario usecase - get by id: validation error: %w",
			errs.ErrScenarioIDRequired,
		)
	}

	var result port.ScenarioWithSteps

	err := service.transactor.WithTx(ctx, func(ctx context.Context) error {
		scenario, err := service.scenarioRepository.GetByID(ctx, scenarioID)
		if err != nil {
			return err
		}

		steps, err := service.stepReader.ListByScenarioID(ctx, scenarioID)
		if err != nil {
			return err
		}

		result = port.ScenarioWithSteps{
			Scenario: scenario,
			Steps:    steps,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) {
			return port.ScenarioWithSteps{}, err
		}
		return port.ScenarioWithSteps{}, service.wrapScenarioOperationError("get by id", err, scenarioID)
	}

	return result, nil
}

func (service *scenarioService) Update(
	ctx context.Context,
	params port.UpdateScenarioParams,
) (entity.Scenario, error) {
	normalizedParams, err := normalizeAndValidateUpdateParams(params)
	if err != nil {
		return entity.Scenario{}, fmt.Errorf("scenario usecase - update: validation error: %w", err)
	}

	var result entity.Scenario

	err = service.transactor.WithTx(ctx, func(ctx context.Context) error {
		currentScenario, err := service.scenarioRepository.LockActive(ctx, normalizedParams.ScenarioID)
		if err != nil {
			return err
		}
		if currentScenario.Status == entity.ScenarioStatusEnabled {
			return errs.ErrScenarioImmutable
		}

		updatedScenario, err := service.scenarioRepository.Update(ctx, normalizedParams)
		if err != nil {
			return err
		}

		result = updatedScenario
		return nil
	})
	if err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) || errors.Is(err, errs.ErrScenarioImmutable) {
			return entity.Scenario{}, err
		}
		return entity.Scenario{}, service.wrapScenarioOperationError("update", err, normalizedParams.ScenarioID)
	}

	return result, nil
}

func (service *scenarioService) Delete(
	ctx context.Context,
	scenarioID uuid.UUID,
) error {
	if scenarioID == uuid.Nil {
		return fmt.Errorf("scenario usecase - delete: validation error: %w", errs.ErrScenarioIDRequired)
	}

	if err := service.scenarioRepository.Delete(ctx, scenarioID); err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) {
			return err
		}
		return service.wrapScenarioOperationError("delete", err, scenarioID)
	}

	return nil
}

func (service *scenarioService) Publish(
	ctx context.Context,
	scenarioID uuid.UUID,
) (entity.Scenario, error) {
	return service.transitionStatus(
		ctx,
		scenarioID,
		"publish",
		entity.ScenarioStatusInDevelopment,
		service.scenarioRepository.Publish,
	)
}

func (service *scenarioService) Enable(
	ctx context.Context,
	scenarioID uuid.UUID,
) (entity.Scenario, error) {
	return service.transitionStatus(
		ctx,
		scenarioID,
		"enable",
		entity.ScenarioStatusDisabled,
		service.scenarioRepository.Enable,
	)
}

func (service *scenarioService) Disable(
	ctx context.Context,
	scenarioID uuid.UUID,
) (entity.Scenario, error) {
	return service.transitionStatus(
		ctx,
		scenarioID,
		"disable",
		entity.ScenarioStatusEnabled,
		service.scenarioRepository.Disable,
	)
}

func normalizeAndValidateUpdateParams(params port.UpdateScenarioParams) (port.UpdateScenarioParams, error) {
	if params.ScenarioID == uuid.Nil {
		return port.UpdateScenarioParams{}, errs.ErrScenarioIDRequired
	}
	if params.Name == nil && params.Description == nil && params.PagePattern == nil {
		return port.UpdateScenarioParams{}, errs.ErrScenarioUpdateFieldsRequired
	}

	if params.Name != nil {
		name := strings.TrimSpace(*params.Name)
		if name == "" {
			return port.UpdateScenarioParams{}, errs.ErrScenarioNameInvalid
		}
		if utf8.RuneCountInString(name) > entity.MaxScenarioNameLength {
			return port.UpdateScenarioParams{}, errs.ErrScenarioNameTooLong
		}
		params.Name = &name
	}

	if params.Description != nil {
		description := strings.TrimSpace(*params.Description)
		if utf8.RuneCountInString(description) > entity.MaxScenarioDescriptionLength {
			return port.UpdateScenarioParams{}, errs.ErrScenarioDescriptionTooLong
		}
		params.Description = &description
	}

	if params.PagePattern != nil {
		pagePattern := strings.TrimSpace(*params.PagePattern)
		if pagePattern == "" {
			return port.UpdateScenarioParams{}, errs.ErrScenarioPagePatternInvalid
		}
		if utf8.RuneCountInString(pagePattern) > entity.MaxScenarioPagePatternLength {
			return port.UpdateScenarioParams{}, errs.ErrScenarioPagePatternTooLong
		}
		params.PagePattern = &pagePattern
	}

	return params, nil
}

func (service *scenarioService) transitionStatus(
	ctx context.Context,
	scenarioID uuid.UUID,
	operation string,
	expectedStatus entity.ScenarioStatus,
	transition func(context.Context, uuid.UUID) (entity.Scenario, error),
) (entity.Scenario, error) {
	if scenarioID == uuid.Nil {
		return entity.Scenario{}, fmt.Errorf(
			"scenario usecase - %s: validation error: %w",
			operation,
			errs.ErrScenarioIDRequired,
		)
	}

	var result entity.Scenario

	err := service.transactor.WithTx(ctx, func(ctx context.Context) error {
		currentScenario, err := service.scenarioRepository.LockActive(ctx, scenarioID)
		if err != nil {
			return err
		}
		if currentScenario.Status != expectedStatus {
			return errs.ErrInvalidScenarioStatusTransition
		}

		updatedScenario, err := transition(ctx, scenarioID)
		if err != nil {
			return err
		}

		result = updatedScenario
		return nil
	})
	if err != nil {
		if errors.Is(err, errs.ErrScenarioNotFound) ||
			errors.Is(err, errs.ErrInvalidScenarioStatusTransition) {
			return entity.Scenario{}, err
		}
		return entity.Scenario{}, service.wrapScenarioOperationError(operation, err, scenarioID)
	}

	return result, nil
}

func (service *scenarioService) wrapProjectOperationError(
	operation string,
	err error,
	projectID uuid.UUID,
) error {
	service.logger.Error(
		"scenario usecase operation failed",
		zap.String("operation", operation),
		zap.String("project_id", projectID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("scenario usecase - %s: project_id=%v: %w", operation, projectID, err)
}

func (service *scenarioService) wrapScenarioOperationError(
	operation string,
	err error,
	scenarioID uuid.UUID,
) error {
	service.logger.Error(
		"scenario usecase operation failed",
		zap.String("operation", operation),
		zap.String("scenario_id", scenarioID.String()),
		zap.Error(err),
	)

	return fmt.Errorf("scenario usecase - %s: scenario_id=%v: %w", operation, scenarioID, err)
}
