package analyticshttp

import (
	"net/http"
)

func (h *AnalyticsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/scenarios/{scenarioId}/analytics/total", h.GetScenarioAnalyticsTotal)
	mux.HandleFunc("GET /api/v1/scenarios/{scenarioId}/analytics/detailed", h.GetDetailedScenarioAnalytics)
	mux.HandleFunc("GET /api/v1/projects/{projectId}/analytics/total", h.GetProjectAnalyticsTotal)
}
