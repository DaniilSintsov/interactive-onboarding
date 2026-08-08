package scenariohttp

import (
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
)

func (h *Handler) createStep(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request CreateStepRequest
	if err := httpserver.ParseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.stepService.Create(
		r.Context(),
		createStepRequestToParams(request, scenarioID),
	)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusCreated, stepToResponse(result))
}

func (h *Handler) updateStep(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	stepID, err := httpserver.ParseUUIDPath(r, "stepId", "invalid_step_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request UpdateStepRequest
	if err := httpserver.ParseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.stepService.Update(
		r.Context(),
		updateStepRequestToParams(request, scenarioID, stepID),
	)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, stepToResponse(result))
}

func (h *Handler) deleteStep(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	stepID, err := httpserver.ParseUUIDPath(r, "stepId", "invalid_step_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if err := h.stepService.Delete(r.Context(), scenarioID, stepID); err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) reorderSteps(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	var request ReorderStepsRequest
	if err := httpserver.ParseJSON(w, r, &request); err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.stepService.Reorder(
		r.Context(),
		reorderStepsRequestToParams(request, scenarioID),
	)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, stepsToResponse(result))
}
