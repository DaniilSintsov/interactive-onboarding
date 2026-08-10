package projecthttp

import (
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
)

func (h *Handler) createElement(w http.ResponseWriter, r *http.Request) {
	projectID, err := httpserver.ParseUUIDPath(r, "projectId", "invalid_project_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request CreateElementRequest
	if err = httpserver.ParseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.elementService.Create(r.Context(), createElementRequestToParams(request, projectID))
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(
		w,
		http.StatusCreated,
		elementToResponse(result),
	)
}

func (h *Handler) updateElement(w http.ResponseWriter, r *http.Request) {
	projectID, err := httpserver.ParseUUIDPath(r, "projectId", "invalid_project_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	elementID, err := httpserver.ParseUUIDPath(r, "elementId", "invalid_element_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request UpdateElementRequest
	if err = httpserver.ParseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.elementService.Update(r.Context(), updateElementRequestToParams(request, projectID, elementID))
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(
		w,
		http.StatusOK,
		elementToResponse(result),
	)
}

func (h *Handler) deleteElement(w http.ResponseWriter, r *http.Request) {
	projectID, err := httpserver.ParseUUIDPath(r, "projectId", "invalid_project_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	elementID, err := httpserver.ParseUUIDPath(r, "elementId", "invalid_element_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	err = h.elementService.Delete(r.Context(), projectID, elementID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listElements(w http.ResponseWriter, r *http.Request) {
	projectID, err := httpserver.ParseUUIDPath(r, "projectId", "invalid_project_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.elementService.List(r.Context(), projectID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(
		w,
		http.StatusOK,
		elementsToResponse(result),
	)
}
