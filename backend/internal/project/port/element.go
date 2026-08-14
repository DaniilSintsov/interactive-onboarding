package port

import (
	"strings"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	"github.com/google/uuid"
)

type CreateElementParams struct {
	ProjectID   uuid.UUID
	Key         string
	Label       string
	Description string
	Page        string
}

type UpdateElementParams struct {
	ProjectID   uuid.UUID
	ElementID   uuid.UUID
	Key         *string
	Label       *string
	Description *string
}

func getStringFromPtr(s *string) string {
	if s == nil {
		return "empty"
	}
	return *s
}

func (params *UpdateElementParams) Validate() error {
	params.Normalize()
	if params.ProjectID == uuid.Nil {
		return errs.ErrElementProjectIDRequired
	}
	if params.ElementID == uuid.Nil {
		return errs.ErrElementIDRequired
	}
	if params.Key == nil && params.Label == nil && params.Description == nil {
		return errs.ErrEmptyElementUpdateParams
	}
	return entity.CheckElementContent(
		getStringFromPtr(params.Key),
		getStringFromPtr(params.Label),
		getStringFromPtr(params.Description),
	)
}

func (params *UpdateElementParams) Normalize() {
	if params.Key != nil {
		*params.Key = strings.TrimSpace(*params.Key)
	}
	if params.Label != nil {
		*params.Label = strings.TrimSpace(*params.Label)
	}
	if params.Description != nil {
		*params.Description = strings.TrimSpace(*params.Description)
	}
}
