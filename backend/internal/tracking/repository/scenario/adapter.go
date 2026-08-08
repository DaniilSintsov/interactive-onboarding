package scenario

import (
	"context"
	"time"

	trackingModel "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/model"
	scenarioSQLC "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/repository/scenario/sqlc"
	trackingService "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ScenarioRepository struct {
	queries *scenarioSQLC.Queries
}

var _ trackingService.SceanrioRepository = (*ScenarioRepository)(nil)

func NewScenarioRepository(db scenarioSQLC.DBTX) *ScenarioRepository {
	return &ScenarioRepository{queries: scenarioSQLC.New(db)}
}

func (r *ScenarioRepository) GetScenarioByIdAndProjectKey(
	ctx context.Context, scenarioID, projectKey string,
) (*trackingModel.Scenario, error) {
	id, err := uuid.Parse(scenarioID)
	if err != nil {
		return nil, err
	}

	found, err := r.queries.GetScenarioByIdAndProjectKey(ctx, scenarioSQLC.GetScenarioByIdAndProjectKeyParams{
		ID:         id,
		ProjectKey: projectKey,
	})
	if err != nil {
		return nil, err
	}

	return adaptScenario(found), nil
}

func adaptScenario(source scenarioSQLC.OnboardingScenario) *trackingModel.Scenario {
	return &trackingModel.Scenario{
		ID:          source.ID.String(),
		ProjectID:   source.ProjectID.String(),
		Name:        source.Name,
		Description: source.Description,
		PagePattern: source.PagePattern,
		Status:      trackingModel.ScenarioStatus(source.Status),
		PublishedAt: nullableTime(source.PublishedAt),
		CreatedAt:   source.CreatedAt.Time,
		UpdatedAt:   source.UpdatedAt.Time,
		DeletedAt:   nullableTime(source.DeletedAt),
	}
}

func nullableTime(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}
