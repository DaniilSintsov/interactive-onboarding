package model

import "encoding/json"

type ResolveScenarioRequest struct {
	Page   string `json:"page" validate:"required,min=1,max=2048"`
	UserID string `json:"user_id" validate:"required,min=1,max=255"`
}

type RuntimeScenario struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	PagePattern string        `json:"page_pattern"`
	Steps       []RuntimeStep `json:"steps"`
}

type RuntimeStep struct {
	ID           string          `json:"id"`
	StepNum      int             `json:"step_num"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	FrontendData json.RawMessage `json:"frontend_data"`
	Element      RuntimeElement  `json:"element"`
}

type RuntimeElement struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type User struct {
	UserId    string `json:"user_id"`
	Onboarded bool   `json:"onboarded"`
}
