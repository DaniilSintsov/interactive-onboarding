package project

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
	DeletedAt   *time.Time
}

type Project struct {
	ID         uuid.UUID
	Name       string
	ProjectKey string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
