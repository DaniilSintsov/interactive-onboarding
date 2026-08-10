package port

import (
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/google/uuid"
)

type CreateScenarioParams struct {
	ProjectID   uuid.UUID
	Name        string
	Description string
	PagePattern string
}

type UpdateScenarioParams struct {
	ScenarioID  uuid.UUID
	Name        *string
	Description *string
	PagePattern *string
}

type ListScenariosParams struct {
	ProjectID uuid.UUID
	Status    *entity.ScenarioStatus
	Limit     int
	Offset    int
}

type ScenarioSummary struct {
	Scenario   entity.Scenario
	StepsCount int64
}

type ListScenariosResult struct {
	Scenarios []ScenarioSummary
	Total     int64
	Limit     int
	Offset    int
}

type ScenarioWithSteps struct {
	Scenario entity.Scenario
	Steps    []entity.Step
}
