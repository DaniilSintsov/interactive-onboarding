package pdfhttp

import (
	"bytes"
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"code":    code,
		"message": message,
	})
}

func (h *PDFHandler) GenerateScenarioPDFReport(w http.ResponseWriter, r *http.Request) {
	scenarioID, err := httpserver.ParseUUIDPath(r, "scenarioId", "invalid_scenario_id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_scenario_id", err.Error())
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
		// После отправки заголовков JSON-ошибку уже не вернуть, только логируем
		// TODO: добавить логирование ошибки
	}
}

func generatePDF(data service.DetailedScenarioAnalytics) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
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
	pdf.Cell(40, 10, fmt.Sprintf("Average Time: %.2f sec", data.AverageCompletionTimeSeconds))
	pdf.Ln(12)

	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(30, 10, "Step")
	pdf.Cell(30, 10, "Position")
	pdf.Cell(40, 10, "Shown")
	pdf.Cell(40, 10, "Completed")
	pdf.Cell(40, 10, "Completion Rate")
	pdf.Ln(8)

	pdf.SetFont("DejaVu", "", 11)
	for _, step := range data.Steps {
		// Обрезаем через []rune для корректной работы с UTF-8
		title := step.Title
		runes := []rune(title)
		if len(runes) > 20 {
			title = string(runes[:20]) + "..."
		}
		pdf.Cell(30, 8, title)
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
