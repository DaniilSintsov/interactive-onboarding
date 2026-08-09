package pdfhttp

import (
	"net/http"
)

func (h *PDFHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/scenarios/{scenarioId}/report/pdf", h.GenerateScenarioPDFReport)
}
