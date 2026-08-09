package step

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/postgres/transactor"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/errs"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/port"
	sqlc "github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/repository/postgres/step/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	uniqueViolationCode     = "23505"
	foreignKeyViolationCode = "23503"

	stepNumberUniqueConstraint = "steps_active_scenario_step_num_unique"
	stepScenarioFKConstraint   = "steps_scenario_fk"
)

type stepRepository struct {
	queries    *sqlc.Queries
	transactor transactor.Transactor
}

func NewRepository(
	db sqlc.DBTX,
	transactor transactor.Transactor,
) *stepRepository {
	return &stepRepository{
		queries:    sqlc.New(db),
		transactor: transactor,
	}
}

func (repo *stepRepository) getQueries(ctx context.Context) *sqlc.Queries {
	if tx, err := transactor.ExtractTx(ctx); err == nil {
		return repo.queries.WithTx(tx)
	}

	return repo.queries
}

func (repo *stepRepository) Create(
	ctx context.Context,
	step entity.Step,
) (entity.Step, error) {
	if step.StepNum < 1 || step.StepNum > entity.MaxStepNumber {
		return entity.Step{}, errs.ErrInvalidStepNumber
	}

	row, err := repo.getQueries(ctx).CreateStep(ctx, sqlc.CreateStepParams{
		ScenarioID:   step.ScenarioID,
		ElementID:    step.ElementID,
		StepNum:      int32(step.StepNum),
		Title:        step.Title,
		Description:  step.Description,
		FrontendData: []byte(step.FrontendData),
	})
	if err != nil {
		switch {
		case isConstraintError(err, uniqueViolationCode, stepNumberUniqueConstraint):
			return entity.Step{}, errs.ErrStepNumberAlreadyExists
		case isConstraintError(err, foreignKeyViolationCode, stepScenarioFKConstraint):
			return entity.Step{}, errs.ErrScenarioNotFound
		default:
			return entity.Step{}, fmt.Errorf("step repository - create: %w", err)
		}
	}

	return entity.Step{
		ID:           row.ID,
		ScenarioID:   row.ScenarioID,
		ElementID:    row.ElementID,
		StepNum:      int(row.StepNum),
		Title:        row.Title,
		Description:  row.Description,
		FrontendData: json.RawMessage(row.FrontendData),
		CreatedAt:    row.CreatedAt.Time.UTC(),
		UpdatedAt:    row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *stepRepository) Update(
	ctx context.Context,
	params port.UpdateStepParams,
) (entity.Step, error) {
	row, err := repo.getQueries(ctx).UpdateStep(ctx, sqlc.UpdateStepParams{
		ElementID:    uuidFromPtr(params.ElementID),
		Title:        textFromPtr(params.Title),
		Description:  textFromPtr(params.Description),
		FrontendData: jsonFromPtr(params.FrontendData),
		ScenarioID:   params.ScenarioID,
		StepID:       params.StepID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Step{}, errs.ErrStepNotFound
		}

		if isConstraintError(err, foreignKeyViolationCode, stepScenarioFKConstraint) {
			return entity.Step{}, errs.ErrScenarioNotFound
		}

		return entity.Step{}, fmt.Errorf("step repository - update: %w", err)
	}

	return entity.Step{
		ID:           row.ID,
		ScenarioID:   row.ScenarioID,
		ElementID:    row.ElementID,
		StepNum:      int(row.StepNum),
		Title:        row.Title,
		Description:  row.Description,
		FrontendData: json.RawMessage(row.FrontendData),
		CreatedAt:    row.CreatedAt.Time.UTC(),
		UpdatedAt:    row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *stepRepository) Delete(
	ctx context.Context,
	scenarioID uuid.UUID,
	stepID uuid.UUID,
) error {
	_, err := repo.getQueries(ctx).DeleteStep(ctx, sqlc.DeleteStepParams{
		ScenarioID: scenarioID,
		StepID:     stepID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrStepNotFound
		}

		return fmt.Errorf("step repository - delete: %w", err)
	}

	return nil
}

func (repo *stepRepository) ListByScenarioID(
	ctx context.Context,
	scenarioID uuid.UUID,
) ([]entity.Step, error) {
	rows, err := repo.getQueries(ctx).ListStepsByScenarioID(ctx, scenarioID)
	if err != nil {
		return nil, fmt.Errorf("step repository - list by scenario id: %w", err)
	}

	steps := make([]entity.Step, 0, len(rows))
	for _, row := range rows {
		steps = append(steps, entity.Step{
			ID:           row.ID,
			ScenarioID:   row.ScenarioID,
			ElementID:    row.ElementID,
			StepNum:      int(row.StepNum),
			Title:        row.Title,
			Description:  row.Description,
			FrontendData: json.RawMessage(row.FrontendData),
			CreatedAt:    row.CreatedAt.Time.UTC(),
			UpdatedAt:    row.UpdatedAt.Time.UTC(),
		})
	}

	return steps, nil
}

func (repo *stepRepository) IsElementUsedBySteps(
	ctx context.Context,
	elementID uuid.UUID,
) (bool, error) {
	used, err := repo.getQueries(ctx).IsElementUsedBySteps(ctx, elementID)
	if err != nil {
		return false, fmt.Errorf("step repository - check element usage: %w", err)
	}

	return used, nil
}

func (repo *stepRepository) GetByID(
	ctx context.Context,
	scenarioID uuid.UUID,
	stepID uuid.UUID,
) (entity.Step, error) {
	row, err := repo.getQueries(ctx).GetStepByID(ctx, sqlc.GetStepByIDParams{
		ScenarioID: scenarioID,
		StepID:     stepID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Step{}, errs.ErrStepNotFound
		}

		return entity.Step{}, fmt.Errorf("step repository - get by id: %w", err)
	}

	return entity.Step{
		ID:           row.ID,
		ScenarioID:   row.ScenarioID,
		ElementID:    row.ElementID,
		StepNum:      int(row.StepNum),
		Title:        row.Title,
		Description:  row.Description,
		FrontendData: json.RawMessage(row.FrontendData),
		CreatedAt:    row.CreatedAt.Time.UTC(),
		UpdatedAt:    row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *stepRepository) GetNextNumber(
	ctx context.Context,
	scenarioID uuid.UUID,
) (int, error) {
	maxStepNumber, err := repo.getQueries(ctx).GetMaxStepNumber(ctx, scenarioID)
	if err != nil {
		return 0, fmt.Errorf("step repository - get next number: %w", err)
	}
	if maxStepNumber >= entity.MaxStepNumber {
		return 0, errs.ErrInvalidStepNumber
	}

	return int(maxStepNumber) + 1, nil
}

func (repo *stepRepository) LockActive(
	ctx context.Context,
	scenarioID uuid.UUID,
	stepID uuid.UUID,
) (entity.Step, error) {
	row, err := repo.getQueries(ctx).LockActiveStep(ctx, sqlc.LockActiveStepParams{
		ScenarioID: scenarioID,
		StepID:     stepID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Step{}, errs.ErrStepNotFound
		}

		return entity.Step{}, fmt.Errorf("step repository - lock active: %w", err)
	}

	return entity.Step{
		ID:           row.ID,
		ScenarioID:   row.ScenarioID,
		ElementID:    row.ElementID,
		StepNum:      int(row.StepNum),
		Title:        row.Title,
		Description:  row.Description,
		FrontendData: json.RawMessage(row.FrontendData),
		CreatedAt:    row.CreatedAt.Time.UTC(),
		UpdatedAt:    row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *stepRepository) LockActiveIDsByScenarioID(
	ctx context.Context,
	scenarioID uuid.UUID,
) ([]uuid.UUID, int32, error) {
	rows, err := repo.getQueries(ctx).LockActiveStepsByScenarioID(ctx, scenarioID)
	if err != nil {
		return nil, 0, fmt.Errorf("step repository - lock active by scenario id: %w", err)
	}

	stepIDs := make([]uuid.UUID, 0, len(rows))
	var maxStepNumber int32
	for _, row := range rows {
		stepIDs = append(stepIDs, row.ID)
		if row.StepNum > maxStepNumber {
			maxStepNumber = row.StepNum
		}
	}

	return stepIDs, maxStepNumber, nil
}

func (repo *stepRepository) Reorder(
	ctx context.Context,
	scenarioID uuid.UUID,
	orderedStepIDs []uuid.UUID,
) error {
	if len(orderedStepIDs) > entity.MaxStepNumber {
		return errs.ErrInvalidStepNumber
	}

	err := repo.transactor.WithTx(ctx, func(ctx context.Context) error {
		activeStepIDs, maxStepNumber, err := repo.LockActiveIDsByScenarioID(ctx, scenarioID)
		if err != nil {
			return fmt.Errorf("lock active steps: %w", err)
		}
		if !sameStepIDs(activeStepIDs, orderedStepIDs) {
			return errs.ErrStepDoesNotBelongToScenario
		}
		if int64(maxStepNumber)*2 > int64(entity.MaxStepNumber) {
			return errs.ErrInvalidStepNumber
		}

		queries := repo.getQueries(ctx)

		if err = queries.MoveStepsOutOfOrderRange(ctx, scenarioID); err != nil {
			return fmt.Errorf("move steps out of order range: %w", err)
		}

		affected, err := queries.ApplyStepOrder(ctx, sqlc.ApplyStepOrderParams{
			ScenarioID:     scenarioID,
			OrderedStepIds: orderedStepIDs,
		})
		if err != nil {
			return fmt.Errorf("apply step order: %w", err)
		}
		if affected != int64(len(orderedStepIDs)) {
			return errs.ErrStepDoesNotBelongToScenario
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("step repository - reorder: %w", err)
	}

	return nil
}

func sameStepIDs(activeStepIDs []uuid.UUID, orderedStepIDs []uuid.UUID) bool {
	if len(activeStepIDs) != len(orderedStepIDs) {
		return false
	}

	active := make(map[uuid.UUID]struct{}, len(activeStepIDs))
	for _, stepID := range activeStepIDs {
		active[stepID] = struct{}{}
	}

	seen := make(map[uuid.UUID]struct{}, len(orderedStepIDs))
	for _, stepID := range orderedStepIDs {
		if _, exists := active[stepID]; !exists {
			return false
		}
		if _, duplicate := seen[stepID]; duplicate {
			return false
		}
		seen[stepID] = struct{}{}
	}

	return true
}

func textFromPtr(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{String: *value, Valid: true}
}

func uuidFromPtr(value *uuid.UUID) pgtype.UUID {
	if value == nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{Bytes: [16]byte(*value), Valid: true}
}

func jsonFromPtr(value *json.RawMessage) []byte {
	if value == nil {
		return nil
	}

	return []byte(*value)
}

func isConstraintError(err error, code string, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) &&
		pgErr.Code == code &&
		pgErr.ConstraintName == constraint
}
