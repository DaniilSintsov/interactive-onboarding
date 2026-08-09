package analyticshttp

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/analytics/service"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: svc}
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	httpserver.WriteJSON(w, status, httpserver.ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func (h *AnalyticsHandler) GetScenarioAnalyticsTotal(w http.ResponseWriter, r *http.Request) {
	id, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid_scenario_id", err.Error())
		return
	}

	from, to, err := parseTimeRange(r)
	if err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid_time_range", err.Error())
		return
	}

	analytics, err := h.service.GetScenarioAnalytics(r.Context(), id.String(), from, to)
	if err != nil {
		if errors.Is(err, service.ErrScenarioNotFound) {
			writeJSONError(w, http.StatusNotFound, "scenario_not_found", "Scenario not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetDetailedScenarioAnalytics(w http.ResponseWriter, r *http.Request) {
	id, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid_scenario_id", err.Error())
		return
	}

	from, to, err := parseTimeRange(r)
	if err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid_time_range", err.Error())
		return
	}

	analytics, err := h.service.GetDetailedScenarioAnalytics(r.Context(), id.String(), from, to)
	if err != nil {
		if errors.Is(err, service.ErrScenarioNotFound) {
			writeJSONError(w, http.StatusNotFound, "scenario_not_found", "Scenario not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetProjectAnalyticsTotal(w http.ResponseWriter, r *http.Request) {
	id, err := httpserver.ParseUUIDPath(r, "projectId", "invalid_project_id")
	if err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid_project_id", err.Error())
		return
	}

	from, to, err := parseTimeRange(r)
	if err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid_time_range", err.Error())
		return
	}

	analytics, err := h.service.GetProjectAnalytics(r.Context(), id.String(), from, to)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			writeJSONError(w, http.StatusNotFound, "project_not_found", "Project not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, analytics)
}

func parseTimeRange(r *http.Request) (*time.Time, *time.Time, error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" && toStr == "" {
		return nil, nil, nil
	}

	if fromStr == "" || toStr == "" {
		return nil, nil, errors.New("both 'from' and 'to' must be provided together")
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid 'from' format: %w", err)
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid 'to' format: %w", err)
	}

	if !from.Before(to) {
		return nil, nil, errors.New("'from' must be before 'to'")
	}

	return &from, &to, nil
}
