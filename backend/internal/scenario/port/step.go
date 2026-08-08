package port

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateStepParams struct {
	ScenarioID   uuid.UUID
	ElementID    uuid.UUID
	StepNum      *int
	Title        string
	Description  string
	FrontendData json.RawMessage
}

type UpdateStepParams struct {
	ScenarioID   uuid.UUID
	StepID       uuid.UUID
	ElementID    *uuid.UUID
	Title        *string
	Description  *string
	FrontendData *json.RawMessage
}

type ReorderStepsParams struct {
	ScenarioID     uuid.UUID
	OrderedStepIDs []uuid.UUID
}
