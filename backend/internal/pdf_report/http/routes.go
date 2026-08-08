package pdfhttp

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *PDFHandler) {
	mux.HandleFunc("GET /api/v1/scenarios/{scenarioId}/report/pdf", h.GenerateScenarioPDFReport)
}
