package scenariohttp

import (
	"context"
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	scenarioService interface {
		Create(ctx context.Context, params port.CreateScenarioParams) (entity.Scenario, error)
		List(ctx context.Context, params port.ListScenariosParams) (port.ListScenariosResult, error)
		GetByID(ctx context.Context, scenarioID uuid.UUID) (port.ScenarioWithSteps, error)
		Update(ctx context.Context, params port.UpdateScenarioParams) (entity.Scenario, error)
		Delete(ctx context.Context, scenarioID uuid.UUID) error
		Publish(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
		Enable(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
		Disable(ctx context.Context, scenarioID uuid.UUID) (entity.Scenario, error)
	}

	stepService interface {
		Create(ctx context.Context, params port.CreateStepParams) (entity.Step, error)
		Update(ctx context.Context, params port.UpdateStepParams) (entity.Step, error)
		Delete(ctx context.Context, scenarioID uuid.UUID, stepID uuid.UUID) error
		Reorder(ctx context.Context, params port.ReorderStepsParams) ([]entity.Step, error)
	}

	testTokenService interface {
		Create(ctx context.Context, scenarioID uuid.UUID) (port.CreatedScenarioTestToken, error)
	}
)

type Handler struct {
	scenarioService  scenarioService
	stepService      stepService
	testTokenService testTokenService
	logger           *zap.Logger
}

func NewHandler(
	scenarioService scenarioService,
	stepService stepService,
	testTokenService testTokenService,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		scenarioService:  scenarioService,
		stepService:      stepService,
		testTokenService: testTokenService,
		logger:           logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/projects/{projectId}/scenarios", h.createScenario)
	mux.HandleFunc("GET /api/v1/projects/{projectId}/scenarios", h.listScenarios)

	mux.HandleFunc("GET /api/v1/scenarios/{scenarioId}", h.getScenario)
	mux.HandleFunc("PATCH /api/v1/scenarios/{scenarioId}", h.updateScenario)
	mux.HandleFunc("DELETE /api/v1/scenarios/{scenarioId}", h.deleteScenario)

	mux.HandleFunc("POST /api/v1/scenarios/{scenarioId}/publish", h.publishScenario)
	mux.HandleFunc("POST /api/v1/scenarios/{scenarioId}/enable", h.enableScenario)
	mux.HandleFunc("POST /api/v1/scenarios/{scenarioId}/disable", h.disableScenario)
	mux.HandleFunc("POST /api/v1/scenarios/{scenarioId}/test-tokens", h.createTestToken)

	mux.HandleFunc("POST /api/v1/scenarios/{scenarioId}/steps", h.createStep)
	mux.HandleFunc("PATCH /api/v1/scenarios/{scenarioId}/steps/{stepId}", h.updateStep)
	mux.HandleFunc("DELETE /api/v1/scenarios/{scenarioId}/steps/{stepId}", h.deleteStep)
	mux.HandleFunc("PUT /api/v1/scenarios/{scenarioId}/steps/order", h.reorderSteps)
}
