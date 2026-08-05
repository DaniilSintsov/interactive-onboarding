package entity

import (
	"strings"
	"time"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
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

func (e Element) Validate() error {
	if e.ProjectID == uuid.Nil {
		return errs.ErrElementProjectIDRequired
	}
	if strings.TrimSpace(e.Key) == "" {
		return errs.ErrElementKeyRequired
	}
	if strings.TrimSpace(e.Label) == "" {
		return errs.ErrElementLabelRequired
	}
	return nil
}

type Project struct {
	ID         uuid.UUID
	Name       string
	ProjectKey string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

func (p Project) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errs.ErrProjectNameRequired
	}
	if strings.TrimSpace(p.ProjectKey) == "" {
		return errs.ErrProjectKeyRequired
	}
	return nil
}
