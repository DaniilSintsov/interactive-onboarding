package scenario

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/postgres/transactor"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/errs"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/port"
	sqlc "github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/repository/postgres/scenario/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	foreignKeyViolationCode = "23503"

	scenarioProjectFKConstraint = "scenario_project_fk"
)

type scenarioRepository struct {
	queries *sqlc.Queries
}

func NewRepository(db sqlc.DBTX) *scenarioRepository {
	return &scenarioRepository{
		queries: sqlc.New(db),
	}
}

func (repo *scenarioRepository) getQueries(ctx context.Context) *sqlc.Queries {
	if tx, err := transactor.ExtractTx(ctx); err == nil {
		return repo.queries.WithTx(tx)
	}

	return repo.queries
}

func (repo *scenarioRepository) Create(
	ctx context.Context,
	scenario entity.Scenario,
) (entity.Scenario, error) {
	row, err := repo.getQueries(ctx).CreateScenario(ctx, sqlc.CreateScenarioParams{
		ProjectID:   scenario.ProjectID,
		Name:        scenario.Name,
		Description: scenario.Description,
		PagePattern: scenario.PagePattern,
	})
	if err != nil {
		if isConstraintError(err, foreignKeyViolationCode, scenarioProjectFKConstraint) {
			return entity.Scenario{}, errs.ErrProjectNotFound
		}

		return entity.Scenario{}, fmt.Errorf("scenario repository - create: %w", err)
	}

	return entity.Scenario{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Name:        row.Name,
		Description: row.Description,
		PagePattern: row.PagePattern,
		Status:      entity.ScenarioStatus(row.Status),
		PublishedAt: timePtr(row.PublishedAt),
		CreatedAt:   row.CreatedAt.Time.UTC(),
		UpdatedAt:   row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *scenarioRepository) Update(
	ctx context.Context,
	params port.UpdateScenarioParams,
) (entity.Scenario, error) {
	row, err := repo.getQueries(ctx).UpdateScenario(ctx, sqlc.UpdateScenarioParams{
		Name:        textFromPtr(params.Name),
		Description: textFromPtr(params.Description),
		PagePattern: textFromPtr(params.PagePattern),
		ScenarioID:  params.ScenarioID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Scenario{}, errs.ErrScenarioNotFound
		}

		return entity.Scenario{}, fmt.Errorf("scenario repository - update: %w", err)
	}

	return entity.Scenario{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Name:        row.Name,
		Description: row.Description,
		PagePattern: row.PagePattern,
		Status:      entity.ScenarioStatus(row.Status),
		PublishedAt: timePtr(row.PublishedAt),
		CreatedAt:   row.CreatedAt.Time.UTC(),
		UpdatedAt:   row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *scenarioRepository) Delete(
	ctx context.Context,
	scenarioID uuid.UUID,
) error {
	_, err := repo.getQueries(ctx).DeleteScenario(ctx, scenarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrScenarioNotFound
		}

		return fmt.Errorf("scenario repository - delete: %w", err)
	}

	return nil
}

func (repo *scenarioRepository) ListByProjectID(
	ctx context.Context,
	projectID uuid.UUID,
	status *entity.ScenarioStatus,
	limit int32,
	offset int32,
) ([]port.ScenarioSummary, int64, error) {
	rows, err := repo.getQueries(ctx).ListScenariosByProjectID(
		ctx,
		sqlc.ListScenariosByProjectIDParams{
			ProjectID:  projectID,
			Status:     statusFromPtr(status),
			PageOffset: offset,
			PageLimit:  limit,
		},
	)
	if err != nil {
		return nil, 0, fmt.Errorf("scenario repository - list by project id: %w", err)
	}

	scenarios := make([]port.ScenarioSummary, 0, len(rows))
	var total int64

	for i, row := range rows {
		if i == 0 {
			total = row.Total
		}

		scenarios = append(scenarios, port.ScenarioSummary{
			Scenario: entity.Scenario{
				ID:          row.ID,
				ProjectID:   row.ProjectID,
				Name:        row.Name,
				Description: row.Description,
				PagePattern: row.PagePattern,
				Status:      entity.ScenarioStatus(row.Status),
				PublishedAt: timePtr(row.PublishedAt),
				CreatedAt:   row.CreatedAt.Time.UTC(),
				UpdatedAt:   row.UpdatedAt.Time.UTC(),
			},
			StepsCount: row.StepsCount,
		})
	}

	return scenarios, total, nil
}

func (repo *scenarioRepository) GetByID(
	ctx context.Context,
	scenarioID uuid.UUID,
) (entity.Scenario, error) {
	row, err := repo.getQueries(ctx).GetScenarioByID(ctx, scenarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Scenario{}, errs.ErrScenarioNotFound
		}

		return entity.Scenario{}, fmt.Errorf("scenario repository - get by id: %w", err)
	}

	return entity.Scenario{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Name:        row.Name,
		Description: row.Description,
		PagePattern: row.PagePattern,
		Status:      entity.ScenarioStatus(row.Status),
		PublishedAt: timePtr(row.PublishedAt),
		CreatedAt:   row.CreatedAt.Time.UTC(),
		UpdatedAt:   row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *scenarioRepository) LockActive(
	ctx context.Context,
	scenarioID uuid.UUID,
) (entity.Scenario, error) {
	row, err := repo.getQueries(ctx).LockActiveScenario(ctx, scenarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Scenario{}, errs.ErrScenarioNotFound
		}

		return entity.Scenario{}, fmt.Errorf("scenario repository - lock active: %w", err)
	}

	return entity.Scenario{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Name:        row.Name,
		Description: row.Description,
		PagePattern: row.PagePattern,
		Status:      entity.ScenarioStatus(row.Status),
		PublishedAt: timePtr(row.PublishedAt),
		CreatedAt:   row.CreatedAt.Time.UTC(),
		UpdatedAt:   row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *scenarioRepository) Publish(
	ctx context.Context,
	scenarioID uuid.UUID,
) (entity.Scenario, error) {
	row, err := repo.getQueries(ctx).PublishScenario(ctx, scenarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Scenario{}, errs.ErrInvalidScenarioStatusTransition
		}

		return entity.Scenario{}, fmt.Errorf("scenario repository - publish: %w", err)
	}

	return entity.Scenario{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Name:        row.Name,
		Description: row.Description,
		PagePattern: row.PagePattern,
		Status:      entity.ScenarioStatus(row.Status),
		PublishedAt: timePtr(row.PublishedAt),
		CreatedAt:   row.CreatedAt.Time.UTC(),
		UpdatedAt:   row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *scenarioRepository) Enable(
	ctx context.Context,
	scenarioID uuid.UUID,
) (entity.Scenario, error) {
	row, err := repo.getQueries(ctx).EnableScenario(ctx, scenarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Scenario{}, errs.ErrInvalidScenarioStatusTransition
		}

		return entity.Scenario{}, fmt.Errorf("scenario repository - enable: %w", err)
	}

	return entity.Scenario{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Name:        row.Name,
		Description: row.Description,
		PagePattern: row.PagePattern,
		Status:      entity.ScenarioStatus(row.Status),
		PublishedAt: timePtr(row.PublishedAt),
		CreatedAt:   row.CreatedAt.Time.UTC(),
		UpdatedAt:   row.UpdatedAt.Time.UTC(),
	}, nil
}

func (repo *scenarioRepository) Disable(
	ctx context.Context,
	scenarioID uuid.UUID,
) (entity.Scenario, error) {
	row, err := repo.getQueries(ctx).DisableScenario(ctx, scenarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Scenario{}, errs.ErrInvalidScenarioStatusTransition
		}

		return entity.Scenario{}, fmt.Errorf("scenario repository - disable: %w", err)
	}

	return entity.Scenario{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Name:        row.Name,
		Description: row.Description,
		PagePattern: row.PagePattern,
		Status:      entity.ScenarioStatus(row.Status),
		PublishedAt: timePtr(row.PublishedAt),
		CreatedAt:   row.CreatedAt.Time.UTC(),
		UpdatedAt:   row.UpdatedAt.Time.UTC(),
	}, nil
}

func timePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}

	t := value.Time.UTC()
	return &t
}

func textFromPtr(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{String: *value, Valid: true}
}

func statusFromPtr(value *entity.ScenarioStatus) pgtype.Text {
	if value == nil {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{String: string(*value), Valid: true}
}

func isConstraintError(err error, code string, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) &&
		pgErr.Code == code &&
		pgErr.ConstraintName == constraint
}
