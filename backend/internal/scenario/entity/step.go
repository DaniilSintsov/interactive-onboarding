package entity

import (
	"encoding/json"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/errs"
	"github.com/google/uuid"
)

const (
	MaxStepTitleLength       = 255
	MaxStepDescriptionLength = 2000
)

type Step struct {
	ID           uuid.UUID
	ScenarioID   uuid.UUID
	StepNum      int
	Title        string
	Description  string
	ElementID    uuid.UUID
	FrontendData json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func (step *Step) Normalize() {
	step.Title = strings.TrimSpace(step.Title)
	step.Description = strings.TrimSpace(step.Description)
}

func (step *Step) Validate() error {
	step.Normalize()
	
	if step.ElementID == uuid.Nil {
		return errs.ErrStepElementIDRequired
	}
	if step.ScenarioID == uuid.Nil {
		return errs.ErrStepScenarioIDRequired
	}
	if strings.TrimSpace(step.Title) == "" {
		return errs.ErrStepTitleRequired
	}
	if step.StepNum < 1 {
		return errs.ErrInvalidStepNumber
	}
	if utf8.RuneCountInString(step.Title) > MaxStepTitleLength {
		return errs.ErrStepTitleTooLong
	}
	if utf8.RuneCountInString(step.Description) > MaxStepDescriptionLength {
		return errs.ErrStepDescriptionTooLong
	}
	return nil
}
