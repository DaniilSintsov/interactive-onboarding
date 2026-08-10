package scenariohttp

import (
	"encoding/json"
	"time"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/google/uuid"
)

type CreateScenarioRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PagePattern string `json:"page_pattern"`
}

type UpdateScenarioRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	PagePattern *string `json:"page_pattern"`
}

type CreateStepRequest struct {
	ElementID    uuid.UUID       `json:"element_id"`
	StepNum      *int            `json:"step_num"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	FrontendData json.RawMessage `json:"frontend_data"`
}

type UpdateStepRequest struct {
	ElementID    *uuid.UUID       `json:"element_id"`
	Title        *string          `json:"title"`
	Description  *string          `json:"description"`
	FrontendData *json.RawMessage `json:"frontend_data"`
}

type ReorderStepsRequest struct {
	OrderedStepIDs []uuid.UUID `json:"ordered_step_ids"`
}

type ScenarioResponse struct {
	ID          uuid.UUID             `json:"id"`
	ProjectID   uuid.UUID             `json:"project_id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	PagePattern string                `json:"page_pattern"`
	Status      entity.ScenarioStatus `json:"status"`
	PublishedAt *time.Time            `json:"published_at"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type ScenarioSummaryResponse struct {
	ID          uuid.UUID             `json:"id"`
	ProjectID   uuid.UUID             `json:"project_id"`
	Name        string                `json:"name"`
	PagePattern string                `json:"page_pattern"`
	Status      entity.ScenarioStatus `json:"status"`
	StepsCount  int64                 `json:"steps_count"`
	PublishedAt *time.Time            `json:"published_at"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type ScenarioListResponse struct {
	Items  []ScenarioSummaryResponse `json:"items"`
	Total  int64                     `json:"total"`
	Limit  int                       `json:"limit"`
	Offset int                       `json:"offset"`
}

type StepResponse struct {
	ID           uuid.UUID       `json:"id"`
	ScenarioID   uuid.UUID       `json:"scenario_id"`
	ElementID    uuid.UUID       `json:"element_id"`
	StepNum      int             `json:"step_num"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	FrontendData json.RawMessage `json:"frontend_data"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type ScenarioWithStepsResponse struct {
	ScenarioResponse
	Steps []StepResponse `json:"steps"`
}

type ScenarioTestTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
