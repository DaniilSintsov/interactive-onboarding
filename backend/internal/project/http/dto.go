package projecthttp

import (
	"time"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name     string                 `json:"name"`
	Elements []CreateElementRequest `json:"elements"`
}

type CreateElementRequest struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Page        string `json:"page"`
}

type UpdateProjectRequest struct {
	Name *string `json:"name"`
}

type UpdateElementRequest struct {
	Key         *string `json:"key"`
	Label       *string `json:"label"`
	Description *string `json:"description"`
}

type ProjectResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	ProjectKey string    `json:"project_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ElementResponse struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"project_id"`
	Key         string    `json:"key"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Page        string    `json:"page"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectWithElementsResponse struct {
	ProjectResponse
	Elements []ElementResponse `json:"elements"`
}

type ProjectListResponse struct {
	Items  []ProjectResponse `json:"items"`
	Total  int64             `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

type ProjectPageResponse struct {
	Page string `json:"page"`
}

type ProjectPagesResponse struct {
	Items []ProjectPageResponse `json:"items"`
}
