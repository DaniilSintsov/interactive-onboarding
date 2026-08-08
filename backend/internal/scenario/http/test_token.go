package scenariohttp

import (
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
)

func (h *Handler) createTestToken(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	result, err := h.testTokenService.Create(r.Context(), scenarioID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusCreated, testTokenToResponse(result))
}
