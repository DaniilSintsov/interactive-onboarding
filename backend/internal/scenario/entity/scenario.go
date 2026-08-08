package entity

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/errs"
	"github.com/google/uuid"
)

type ScenarioStatus string

const (
	ScenarioStatusEnabled       ScenarioStatus = "enabled"
	ScenarioStatusDisabled      ScenarioStatus = "disabled"
	ScenarioStatusInDevelopment ScenarioStatus = "in_development"
)

const (
	MaxScenarioNameLength        = 255
	MaxScenarioDescriptionLength = 2000
	MaxScenarioPagePatternLength = 2048
)

type Scenario struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Name        string
	Description string
	PagePattern string
	Status      ScenarioStatus
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (status ScenarioStatus) IsValid() bool {
	switch status {
	case ScenarioStatusInDevelopment,
		ScenarioStatusEnabled,
		ScenarioStatusDisabled:
		return true
	default:
		return false
	}
}

func (scenario *Scenario) Normalize() {
	scenario.Name = strings.TrimSpace(scenario.Name)
	scenario.Description = strings.TrimSpace(scenario.Description)
	scenario.PagePattern = strings.TrimSpace(scenario.PagePattern)
}

func (scenario *Scenario) Validate() error {
	scenario.Normalize()

	if scenario.ProjectID == uuid.Nil {
		return errs.ErrScenarioProjectIDRequired
	}
	if !scenario.Status.IsValid() {
		return errs.ErrScenarioStatusUnknown
	}
	if strings.TrimSpace(scenario.Name) == "" {
		return errs.ErrScenarioNameRequired
	}
	if strings.TrimSpace(scenario.PagePattern) == "" {
		return errs.ErrScenarioPagePatternRequired
	}
	if utf8.RuneCountInString(scenario.Name) > MaxScenarioNameLength {
		return errs.ErrScenarioNameTooLong
	}
	if utf8.RuneCountInString(scenario.Description) > MaxScenarioDescriptionLength {
		return errs.ErrScenarioDescriptionTooLong
	}
	if utf8.RuneCountInString(scenario.PagePattern) > MaxScenarioPagePatternLength {
		return errs.ErrScenarioPagePatternTooLong
	}
	return nil
}
