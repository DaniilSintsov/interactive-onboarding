package analyticshttp

import (
	"net/http"
	"time"

	"interactive-onboarding/internal/analytics/service"
	"interactive-onboarding/internal/platform/httpserver"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: svc}
}

func (h *AnalyticsHandler) GetScenarioAnalyticsTotal(w http.ResponseWriter, r *http.Request) {
	id, err := httpserver.ParseUUIDPath(r, "id")
	if err != nil {
		httpserver.WriteJSONError(w, http.StatusBadRequest, "invalid_scenario_id", err.Error())
		return
	}

	from, to := parseTimeRange(r)

	analytics, err := h.service.GetScenarioAnalytics(r.Context(), id.String(), from, to)
	if err != nil {
		httpserver.WriteJSONError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetDetailedScenarioAnalytics(w http.ResponseWriter, r *http.Request) {
	id, err := httpserver.ParseUUIDPath(r, "id")
	if err != nil {
		httpserver.WriteJSONError(w, http.StatusBadRequest, "invalid_scenario_id", err.Error())
		return
	}

	from, to := parseTimeRange(r)

	analytics, err := h.service.GetDetailedScenarioAnalytics(r.Context(), id.String(), from, to)
	if err != nil {
		httpserver.WriteJSONError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetProjectAnalyticsTotal(w http.ResponseWriter, r *http.Request) {
	id, err := httpserver.ParseUUIDPath(r, "id")
	if err != nil {
		httpserver.WriteJSONError(w, http.StatusBadRequest, "invalid_project_id", err.Error())
		return
	}

	from, to := parseTimeRange(r)

	analytics, err := h.service.GetProjectAnalytics(r.Context(), id.String(), from, to)
	if err != nil {
		httpserver.WriteJSONError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, analytics)
}

func parseTimeRange(r *http.Request) (*time.Time, *time.Time) {
	var from, to *time.Time

	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &t
		}
	}
	if toStr := r.URL.Query().Get("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}

	return from, to
}
