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

func (e *Element) Normalize() {
	e.Key = strings.TrimSpace(e.Key)
	e.Label = strings.TrimSpace(e.Label)
	e.Description = strings.TrimSpace(e.Description)
}

func CheckElementFields(key string, label string, description string) error {
	if key == "" {
		return errs.ErrElementKeyRequired
	}
	if label == "" {
		return errs.ErrElementLabelRequired
	}
	if utf8.RuneCountInString(key) > MaxElementKeyLength {
		return errs.ErrElementKeyTooLong
	}
	if utf8.RuneCountInString(label) > MaxElementLabelLength {
		return errs.ErrElementLabelTooLong
	}
	if utf8.RuneCountInString(description) > MaxElementDescriptionLength {
		return errs.ErrElementDescriptionTooLong
	}
	return nil
}

func (e *Element) Validate() error {
	e.Normalize()
	if e.ProjectID == uuid.Nil {
		return errs.ErrElementProjectIDRequired
	}
	return CheckElementFields(e.Key, e.Label, e.Description)
}

type Project struct {
	ID         uuid.UUID
	Name       string
	ProjectKey string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

func (p *Project) Normalize() {
	p.Name = strings.TrimSpace(p.Name)
	p.ProjectKey = strings.TrimSpace(p.ProjectKey)
}

func (p *Project) Validate() error {
	p.Normalize()
	if p.Name == "" {
		return errs.ErrProjectNameRequired
	}
	if utf8.RuneCountInString(p.Name) > MaxProjectNameLength {
		return errs.ErrProjectNameTooLong
	}
	return nil
}
