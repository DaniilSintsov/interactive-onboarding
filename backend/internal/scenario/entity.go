package scellar

import (
	"time"

	"github.com/google/uuid"
)

type Element struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Key         string
	Label       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Project struct {
	ID         uuid.UUID
	Name       string
	ProjectKey string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ScenarioStatus string

const (
	ScenarioStatusEnable        ScenarioStatus = "enable"
	ScenarioStatusDisable       ScenarioStatus = "disable"
	ScenarioStatusInDevelopment ScenarioStatus = "in_development"
)

type Scenario struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Name        string
	Description string
	PagePattern string
	Status      ScenarioStatus
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type Step struct {
	ID           uuid.UUID
	ScenarioID   uuid.UUID
	StepNum      int
	Title        string
	Description  string
	ElementID    uuid.UUID
	FrontendInfo any
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
