package model

import "github.com/google/uuid"

type UpdateElementParams struct {
	ProjectID   uuid.UUID
	ElementID   uuid.UUID
	Key         *string
	Label       *string
	Description *string
}
