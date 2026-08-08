package pdfhttp

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf/v2"

	"interactive-onboarding/internal/analytics/service"
	"interactive-onboarding/internal/platform/httpserver"
)

type PDFHandler struct {
	analyticsService *service.AnalyticsService
}

func NewPDFHandler(svc *service.AnalyticsService) *PDFHandler {
	return &PDFHandler{analyticsService: svc}
}

func (h *PDFHandler) GenerateScenarioPDFReport(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "id")
	if err != nil {
		httpserver.WriteJSONError(w, http.StatusBadRequest, "invalid_scenario_id", err.Error())
		return
	}

	from, to := parseTimeRange(r)

	analytics, err := h.analyticsService.GetDetailedScenarioAnalytics(r.Context(), scenarioID.String(), from, to)
	if err != nil {
		httpserver.WriteJSONError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	pdfData, err := generatePDF(analytics)
	if err != nil {
		httpserver.WriteJSONError(w, http.StatusInternalServerError, "pdf_generation_error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=report_%s.pdf", scenarioID.String()))
	w.Write(pdfData)
}

func generatePDF(data service.DetailedScenarioAnalytics) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Scenario Analytics Report")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Scenario ID: %s", data.ScenarioID))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Started: %d", data.Started))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Completed: %d", data.Completed))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Skipped: %d", data.Skipped))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Completion Rate: %.2f%%", data.CompletionRate*100))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Average Time: %.2f sec", data.AverageCompletionTimeSeconds))
	pdf.Ln(12)

	// Заголовок таблицы
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(30, 10, "Step")
	pdf.Cell(30, 10, "Position")
	pdf.Cell(40, 10, "Shown")
	pdf.Cell(40, 10, "Completed")
	pdf.Cell(40, 10, "Completion Rate")
	pdf.Ln(8)

	// Данные таблицы
	pdf.SetFont("Arial", "", 11)
	for _, step := range data.Steps {
		pdf.Cell(30, 8, step.Title)
		pdf.Cell(30, 8, fmt.Sprintf("%d", step.Position))
		pdf.Cell(40, 8, fmt.Sprintf("%d", step.Shown))
		pdf.Cell(40, 8, fmt.Sprintf("%d", step.Completed))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f%%", step.CompletionRate*100))
		pdf.Ln(7)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
