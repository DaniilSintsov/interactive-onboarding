package handlers

import (
	"net/http"

	"github.com/DaniilSintsov/interactive-onboarding/internal/platform/httpserver"
	trackingModel "github.com/DaniilSintsov/interactive-onboarding/internal/tracking/model"
)

type TrackingService interface {
	StartSession(session *trackingModel.StartSessionRequest) (*trackingModel.OnboardingSession, error)
	CreateEvent(event *trackingModel.CreateEventRequest) (*trackingModel.EventAcceptedResponse, error)
}

type TrackingHandler struct {
	service TrackingService
}

func NewTrackingService(srvc TrackingService) *TrackingHandler {
	return &TrackingHandler{
		service: srvc,
	}
}

func (h *TrackingHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	startSessionReq := new(trackingModel.StartSessionRequest)
	httpserver.ParseJson(r, startSessionReq)
	onboardingSession, err := h.service.StartSession(startSessionReq)
	if err != nil {
		httpserver.WriteJson(w, http.StatusBadRequest, nil)
	}
	httpserver.WriteJson(w, http.StatusOK, onboardingSession)
}

func (h *TrackingHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	startSessionReq := new(trackingModel.StartSessionRequest)
	httpserver.ParseJson(r, startSessionReq)
	onboardingSession, err := h.service.StartSession(startSessionReq)
	if err != nil {
		httpserver.WriteJson(w, http.StatusBadRequest, nil)
	}
	httpserver.WriteJson(w, http.StatusOK, onboardingSession)
}
