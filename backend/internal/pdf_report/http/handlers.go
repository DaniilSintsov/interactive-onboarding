package pdfhttp

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/analytics/service"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/platform/httpserver"
	"github.com/jung-kurt/gofpdf/v2"
)

type PDFHandler struct {
	analyticsService *service.AnalyticsService
}

func NewPDFHandler(svc *service.AnalyticsService) *PDFHandler {
	return &PDFHandler{analyticsService: svc}
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	httpserver.WriteJSON(w, status, httpserver.ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func (h *PDFHandler) GenerateScenarioPDFReport(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid_scenario_id", err.Error())
		return
	}

	from, to, err := parseTimeRange(r)
	if err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, "invalid_time_range", err.Error())
		return
	}

	analytics, err := h.analyticsService.GetDetailedScenarioAnalytics(r.Context(), scenarioID.String(), from, to)
	if err != nil {
		if errors.Is(err, service.ErrScenarioNotFound) {
			writeJSONError(w, http.StatusNotFound, "scenario_not_found", "Scenario not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
		return
	}

	pdfData, err := generatePDF(analytics)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "pdf_generation_error", "Failed to generate PDF")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=report_%s.pdf", scenarioID.String()))
	if _, err := w.Write(pdfData); err != nil {
		return
	}
}

func generatePDF(data service.DetailedScenarioAnalytics) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	pdf.AddUTF8Font("DejaVu", "", "fonts/DejaVuSans.ttf")
	pdf.AddUTF8Font("DejaVu", "B", "fonts/DejaVuSans.ttf")
	pdf.SetFont("DejaVu", "", 14)

	pdf.Cell(40, 10, "Scenario Analytics Report")
	pdf.Ln(12)

	pdf.SetFont("DejaVu", "", 12)
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
	pdf.Cell(40, 10, fmt.Sprintf("Skip Rate: %.2f%%", data.SkipRate*100))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Average Time: %.2f sec", data.AverageCompletionTimeSeconds))
	pdf.Ln(12)

	var printHeader = func() {
		pdf.SetFont("DejaVu", "B", 11)
		pdf.Cell(45, 10, "Step")
		pdf.Cell(20, 10, "Pos")
		pdf.Cell(25, 10, "Shown")
		pdf.Cell(30, 10, "Completed")
		pdf.Cell(25, 10, "Skipped")
		pdf.Cell(30, 10, "Completion %")
		pdf.Cell(25, 10, "Skip %")
		pdf.Cell(30, 10, "Drop-off %")
		pdf.Ln(8)
	}

	printHeader()

	pdf.SetFont("DejaVu", "", 10)
	for _, step := range data.Steps {
		if pdf.GetY() > 180 {
			pdf.AddPage()
			printHeader()
		}

		title := truncateStringByWidth(pdf, step.Title, 45)
		pdf.Cell(45, 7, title)
		pdf.Cell(20, 7, fmt.Sprintf("%d", step.Position))
		pdf.Cell(25, 7, fmt.Sprintf("%d", step.Shown))
		pdf.Cell(30, 7, fmt.Sprintf("%d", step.Completed))
		pdf.Cell(25, 7, fmt.Sprintf("%d", step.Skipped))
		pdf.Cell(30, 7, fmt.Sprintf("%.1f%%", step.CompletionRate*100))
		pdf.Cell(25, 7, fmt.Sprintf("%.1f%%", step.SkipRate*100))
		pdf.Cell(30, 7, fmt.Sprintf("%.1f%%", step.DropOffRate*100))
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

func truncateStringByWidth(pdf *gofpdf.Fpdf, s string, maxWidth float64) string {
	if pdf.GetStringWidth(s) <= maxWidth {
		return s
	}
	runes := []rune(s)
	for i := len(runes); i > 0; i-- {
		truncated := string(runes[:i]) + "..."
		if pdf.GetStringWidth(truncated) <= maxWidth {
			return truncated
		}
	}
	return "..."
}
