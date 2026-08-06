package port

import (
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
	"github.com/google/uuid"
)

type CreateElementParams struct {
	ProjectID   uuid.UUID
	Key         string
	Label       string
	Description string
}

type UpdateElementParams struct {
	ProjectID   uuid.UUID
	ElementID   uuid.UUID
	Key         *string
	Label       *string
	Description *string
}

func (params UpdateElementParams) Validate() error {
	if params.ProjectID == uuid.Nil {
		return errs.ErrElementProjectIDRequired
	}
	if params.ElementID == uuid.Nil {
		return errs.ErrElementIDRequired
	}
	if params.Key == nil && params.Label == nil && params.Description == nil {
		return errs.ErrEmptyElementUpdateParams
	}
	return nil
}

type ListProjectsResult struct {
	Projects []entity.Project
	Total    int64
	Page     int
	PageSize int
}
