package analyticshttp

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, h *AnalyticsHandler) {
	r.Route("/api/v1/scenarios", func(r chi.Router) {
		r.Get("/{scenarioId}/analytics/total", h.GetScenarioAnalyticsTotal)
		r.Get("/{scenarioId}/analytics/detailed", h.GetDetailedScenarioAnalytics)
	})
	r.Route("/api/v1/projects", func(r chi.Router) {
		r.Get("/{projectId}/analytics/total", h.GetProjectAnalyticsTotal)
	})
}
