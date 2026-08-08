package scenariohttp

import (
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/scenario/port"
)

func (h *Handler) createScenario(w http.ResponseWriter, r *http.Request) {
	projectID, err := httpserver.ParseUUIDPath(r, "projectId", "invalid_project_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request CreateScenarioRequest
	if err := httpserver.ParseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.scenarioService.Create(
		r.Context(),
		createScenarioRequestToParams(request, projectID),
	)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusCreated, scenarioToResponse(result))
}

func (h *Handler) listScenarios(w http.ResponseWriter, r *http.Request) {
	projectID, err := httpserver.ParseUUIDPath(r, "projectId", "invalid_project_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	limit, offset, err := httpserver.ParsePagination(r)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	status, err := parseScenarioStatus(r)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.scenarioService.List(r.Context(), port.ListScenariosParams{
		ProjectID: projectID,
		Status:    status,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, scenarioListToResponse(result))
}

func (h *Handler) getScenario(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.scenarioService.GetByID(r.Context(), scenarioID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, scenarioWithStepsToResponse(result))
}

func (h *Handler) updateScenario(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request UpdateScenarioRequest
	if err := httpserver.ParseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.scenarioService.Update(
		r.Context(),
		updateScenarioRequestToParams(request, scenarioID),
	)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, scenarioToResponse(result))
}

func (h *Handler) deleteScenario(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := h.scenarioService.Delete(r.Context(), scenarioID); err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) publishScenario(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.scenarioService.Publish(r.Context(), scenarioID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, scenarioToResponse(result))
}

func (h *Handler) enableScenario(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.scenarioService.Enable(r.Context(), scenarioID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, scenarioToResponse(result))
}

func (h *Handler) disableScenario(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.scenarioService.Disable(r.Context(), scenarioID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, scenarioToResponse(result))
}

func parseScenarioStatus(r *http.Request) (*entity.ScenarioStatus, error) {
	rawStatus := r.URL.Query().Get("status")
	if rawStatus == "" {
		return nil, nil
	}

	status := entity.ScenarioStatus(rawStatus)
	if !status.IsValid() {
		return nil, httpserver.NewRequestError(
			"invalid_status",
			"status must be one of: in_development, enabled, disabled",
			map[string]any{"parameter": "status"},
			nil,
		)
	}

	return &status, nil
}
