package projecthttp

import (
	"context"
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	elementService interface {
		List(ctx context.Context, projectID uuid.UUID, page *string) ([]entity.Element, error)
		ListPages(ctx context.Context, projectID uuid.UUID) ([]string, error)
		Create(ctx context.Context, params port.CreateElementParams) (entity.Element, error)
		Update(ctx context.Context, params port.UpdateElementParams) (entity.Element, error)
		Delete(ctx context.Context, projectID uuid.UUID, elementID uuid.UUID) error
	}

	projectService interface {
		Create(ctx context.Context, params port.CreateProjectParams) (port.ProjectWithElements, error)
		List(ctx context.Context, limit, offset int) (port.ListProjectsResult, error)
		GetByID(ctx context.Context, projectID uuid.UUID) (port.ProjectWithElements, error)
		Update(ctx context.Context, projectID uuid.UUID, name string) (port.ProjectWithElements, error)
		Delete(ctx context.Context, projectID uuid.UUID) error
	}
)

type Handler struct {
	elementService elementService
	projectService projectService
	logger         *zap.Logger
}

func NewHandler(
	elementService elementService,
	projectService projectService,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		elementService: elementService,
		projectService: projectService,
		logger:         logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/projects", h.createProject)
	mux.HandleFunc("GET /api/v1/projects", h.listProjects)

	mux.HandleFunc("GET /api/v1/projects/{projectId}", h.getProject)
	mux.HandleFunc("PATCH /api/v1/projects/{projectId}", h.updateProject)
	mux.HandleFunc("DELETE /api/v1/projects/{projectId}", h.deleteProject)

	mux.HandleFunc("GET /api/v1/projects/{projectId}/elements", h.listElements)
	mux.HandleFunc("GET /api/v1/projects/{projectId}/pages", h.listPages)
	mux.HandleFunc("POST /api/v1/projects/{projectId}/elements", h.createElement)
	mux.HandleFunc("PATCH /api/v1/projects/{projectId}/elements/{elementId}", h.updateElement)
	mux.HandleFunc("DELETE /api/v1/projects/{projectId}/elements/{elementId}", h.deleteElement)
}
