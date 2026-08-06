package projecthttp

import "net/http"

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var request CreateProjectRequest
	if err := parseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.projectService.Create(
		r.Context(),
		createProjectRequestToParams(request),
	)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		projectWithElementsToResponse(result),
	)
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDPath(r, "projectId")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request UpdateProjectRequest
	if err := parseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	if request.Name == nil {
		h.handleError(
			w,
			r,
			newRequestError(
				"invalid_request_body",
				"name is required",
				map[string]any{"field": "name"},
				nil,
			),
		)
		return
	}

	result, err := h.projectService.Update(r.Context(), projectID, *request.Name)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		projectWithElementsToResponse(result),
	)
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDPath(r, "projectId")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := h.projectService.Delete(r.Context(), projectID); err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := parsePagination(r)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.projectService.List(r.Context(), limit, offset)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, projectListToResponse(result))
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUIDPath(r, "projectId")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.projectService.GetByID(r.Context(), projectID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, projectWithElementsToResponse(result))
}
