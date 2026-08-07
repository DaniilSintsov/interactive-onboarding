package analyticshttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"interactive-onboarding/internal/analytics/service"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: svc}
}

func (h *AnalyticsHandler) GetScenarioAnalyticsTotal(c *gin.Context) {
	id, err := parseScenarioID(c)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid_scenario_id", err.Error())
		return
	}

	from, to := parseTimeRange(c)

	analytics, err := h.service.GetScenarioAnalytics(c.Request.Context(), id.String(), from, to)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	writeJSON(c, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetDetailedScenarioAnalytics(c *gin.Context) {
	id, err := parseScenarioID(c)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid_scenario_id", err.Error())
		return
	}

	from, to := parseTimeRange(c)

	analytics, err := h.service.GetDetailedScenarioAnalytics(c.Request.Context(), id.String(), from, to)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	writeJSON(c, http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetProjectAnalyticsTotal(c *gin.Context) {
	id, err := parseProjectID(c)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid_project_id", err.Error())
		return
	}

	from, to := parseTimeRange(c)

	analytics, err := h.service.GetProjectAnalytics(c.Request.Context(), id.String(), from, to)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	writeJSON(c, http.StatusOK, analytics)
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

func parseScenarioID(c *gin.Context) (uuid.UUID, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return uuid.Nil, errors.New("missing scenario id")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid scenario id: %w", err)
	}
	return id, nil
}

func parseProjectID(c *gin.Context) (uuid.UUID, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return uuid.Nil, errors.New("missing project id")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid project id: %w", err)
	}
	return id, nil
}

func writeJSON(c *gin.Context, status int, data any) {
	c.Header("Content-Type", "application/json")
	c.JSON(status, data)
}

func writeJSONError(c *gin.Context, status int, code, message string) {
	c.Header("Content-Type", "application/json")
	c.JSON(status, gin.H{
		"code":    code,
		"message": message,
	})
}

func parseJSON(c *gin.Context, dst any) error {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("request body must contain exactly one JSON value")
	}

	return nil
}
