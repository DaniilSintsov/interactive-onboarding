package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/internal/platform/httpserver"
	platformModel "github.com/DaniilSintsov/interactive-onboarding/internal/platform/model"
	trackingModel "github.com/DaniilSintsov/interactive-onboarding/internal/tracking/model"
	trackingService "github.com/DaniilSintsov/interactive-onboarding/internal/tracking/service"
)

type TrackingService interface {
	StartSession(context.Context, *trackingModel.StartSessionRequest) (*trackingModel.OnboardingSession, error)
	CreateEvent(context.Context, *trackingModel.CreateEventRequest) (*trackingModel.EventAcceptedResponse, error)
}

type TrackingHandler struct {
	service TrackingService
}

func NewTrackingHandler(srvc TrackingService) *TrackingHandler {
	return &TrackingHandler{
		service: srvc,
	}
}

func (h *TrackingHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	startSessionReq := new(trackingModel.StartSessionRequest)
	if err := httpserver.ParseJson(r, startSessionReq); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
		return
	}

	onboardingSession, err := h.service.StartSession(r.Context(), startSessionReq)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	_ = httpserver.WriteJson(w, http.StatusCreated, onboardingSession)
}

func (h *TrackingHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	eventReq := new(trackingModel.CreateEventRequest)
	if err := httpserver.ParseJson(r, eventReq); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
		return
	}

	response, err := h.service.CreateEvent(r.Context(), eventReq)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	status := http.StatusAccepted
	if response.Duplicate {
		status = http.StatusOK
	}
	_ = httpserver.WriteJson(w, status, response)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, trackingService.ErrInvalidRequest):
		writeError(w, http.StatusUnprocessableEntity, "validation_error", err.Error())
	case errors.Is(err, trackingService.ErrScenarioNotFound),
		errors.Is(err, trackingService.ErrSessionNotFound),
		errors.Is(err, trackingService.ErrStepNotFound):
		writeError(w, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, trackingService.ErrSessionNotActive),
		errors.Is(err, trackingService.ErrStepScenarioMismatch):
		writeError(w, http.StatusConflict, "conflict", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	_ = httpserver.WriteJson(w, status, platformModel.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
