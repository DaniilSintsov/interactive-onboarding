package projecthttp

import "net/http"

func (h *Handler) createElement(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDPath(r, "projectId")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request CreateElementRequest
	if err := parseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.elementService.Create(r.Context(), createElementRequestToParams(request, projectID))
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		elementToResponse(result),
	)
}

func (h *Handler) updateElement(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDPath(r, "projectId")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	elementID, err := parseUUIDPath(r, "elementId")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request UpdateElementRequest
	if err := parseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.elementService.Update(r.Context(), updateElementRequestToParams(request, projectID, elementID))
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		elementToResponse(result),
	)
}

func (h *Handler) deleteElement(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDPath(r, "projectId")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	elementID, err := parseUUIDPath(r, "elementId")
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
	projectID, err := parseUUIDPath(r, "projectId")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.elementService.List(r.Context(), projectID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		elementsToResponse(result),
	)
}
