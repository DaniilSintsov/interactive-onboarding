package entity

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	"github.com/google/uuid"
)

const (
	MaxProjectNameLength        = 255
	MaxElementKeyLength         = 255
	MaxElementLabelLength       = 255
	MaxElementDescriptionLength = 2000
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
	if utf8.RuneCountInString(e.Key) > MaxElementKeyLength {
		return errs.ErrElementKeyTooLong
	}
	if utf8.RuneCountInString(e.Label) > MaxElementLabelLength {
		return errs.ErrElementLabelTooLong
	}
	if utf8.RuneCountInString(e.Description) > MaxElementDescriptionLength {
		return errs.ErrElementDescriptionTooLong
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
	if utf8.RuneCountInString(p.Name) > MaxProjectNameLength {
		return errs.ErrProjectNameTooLong
	}
	return nil
}
