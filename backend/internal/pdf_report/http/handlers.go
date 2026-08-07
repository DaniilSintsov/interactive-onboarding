package pdfhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf/v2"

	"interactive-onboarding/internal/analytics/service"
)

type PDFHandler struct {
	analyticsService *service.AnalyticsService
}

func NewPDFHandler(svc *service.AnalyticsService) *PDFHandler {
	return &PDFHandler{analyticsService: svc}
}

func (h *PDFHandler) GenerateScenarioPDFReport(c *gin.Context) {
	scenarioID, err := parseScenarioID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	from, to := parseTimeRange(c)

	analytics, err := h.analyticsService.GetDetailedScenarioAnalytics(
		c.Request.Context(),
		scenarioID.String(),
		from,
		to,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfData, err := generatePDF(analytics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate PDF"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=report_%s.pdf", scenarioID.String()))
	c.Data(http.StatusOK, "application/pdf", pdfData)
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

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(30, 10, "Step")
	pdf.Cell(30, 10, "Position")
	pdf.Cell(40, 10, "Shown")
	pdf.Cell(40, 10, "Completed")
	pdf.Cell(40, 10, "Completion Rate")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 11)
	for _, step := range data.Steps {
		pdf.Cell(30, 8, step.Title)
		pdf.Cell(30, 8, fmt.Sprintf("%d", step.Position))
		pdf.Cell(40, 8, fmt.Sprintf("%d", step.Shown))
		pdf.Cell(40, 8, fmt.Sprintf("%d", step.Completed))
		pdf.Cell(40, 8, fmt.Sprintf("%.2f%%", step.CompletionRate*100))
		pdf.Ln(7)
	}

	return pdf.Output(nil)
}

func parseScenarioID(c *gin.Context) (uuid.UUID, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("missing scenario id")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid scenario id: %w", err)
	}
	return id, nil
}

func parseTimeRange(c *gin.Context) (*time.Time, *time.Time) {
	var from, to *time.Time

	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}
	return from, to
}
