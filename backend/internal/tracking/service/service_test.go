package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	trackingModel "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestCreateEventRecordsEventInTransaction(t *testing.T) {
	ctx := testContext()
	sessions := &sessionFake{}
	events := &eventFake{getErr: sql.ErrNoRows, txSessions: sessions}
	service := newEventService(sessions, events)
	request := validEvent(trackingModel.EventTypeStepShown)

	response, err := service.CreateEvent(ctx, request)
	if err != nil {
		t.Fatalf("CreateEvent() error = %v", err)
	}
	if !events.transactionCalled {
		t.Fatal("CreateEvent() did not use a transaction")
	}
	if events.recorded == nil || events.recorded.ID != request.ID {
		t.Fatalf("RecordEvent() received %#v, want event %q", events.recorded, request.ID)
	}
	if !events.recorded.OccurredAt.Equal(mustParseTime(t, request.OccurredAt)) {
		t.Fatalf("OccurredAt = %s, want %s", events.recorded.OccurredAt, request.OccurredAt)
	}
	if sessions.updateCalls != 0 {
		t.Fatalf("UpdateSessionStatus() calls = %d, want 0", sessions.updateCalls)
	}
	if response.Duplicate {
		t.Fatal("response marked a newly recorded event as duplicate")
	}
}

func TestCreateEventCompletionUpdatesOnlyItsSession(t *testing.T) {
	sessions := &sessionFake{}
	events := &eventFake{getErr: sql.ErrNoRows, txSessions: sessions}
	service := newEventService(sessions, events)
	request := validEvent(trackingModel.EventTypeOnboardingCompleted)

	_, err := service.CreateEvent(testContext(), request)
	if err != nil {
		t.Fatalf("CreateEvent() error = %v", err)
	}
	if sessions.updateCalls != 1 {
		t.Fatalf("UpdateSessionStatus() calls = %d, want 1", sessions.updateCalls)
	}
	if sessions.updatedSessionID != request.SessionID {
		t.Fatalf("updated session = %q, want %q", sessions.updatedSessionID, request.SessionID)
	}
	if sessions.updatedStatus != trackingModel.SessionStatusCompleted {
		t.Fatalf("updated status = %q, want completed", sessions.updatedStatus)
	}
	if !sessions.finishedAt.Equal(mustParseTime(t, request.OccurredAt)) {
		t.Fatalf("finished_at = %s, want %s", sessions.finishedAt, request.OccurredAt)
	}
}

func TestCreateEventReturnsExistingDuplicateWithoutChangingSession(t *testing.T) {
	existing := &trackingModel.EventAcceptedResponse{Duplicate: true}
	sessions := &sessionFake{}
	events := &eventFake{getResponse: existing, txSessions: sessions}
	service := newEventService(sessions, events)

	response, err := service.CreateEvent(testContext(), validEvent(trackingModel.EventTypeOnboardingCompleted))
	if err != nil {
		t.Fatalf("CreateEvent() error = %v", err)
	}
	if response != existing {
		t.Fatal("CreateEvent() did not return the existing event response")
	}
	if events.recorded != nil || sessions.updateCalls != 0 {
		t.Fatal("duplicate event was recorded or changed its session")
	}
}

func TestCreateEventRejectsInvalidRequestBeforeTransaction(t *testing.T) {
	events := &eventFake{}
	service := newEventService(&sessionFake{}, events)
	request := validEvent(trackingModel.EventTypeStepShown)
	request.OccurredAt = "not-a-timestamp"

	_, err := service.CreateEvent(testContext(), request)
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("CreateEvent() error = %v, want ErrInvalidRequest", err)
	}
	if events.transactionCalled {
		t.Fatal("invalid request started a transaction")
	}
}

func TestCreateEventHandlesConcurrentDuplicate(t *testing.T) {
	existing := &trackingModel.EventAcceptedResponse{Duplicate: true}
	events := &eventFake{
		getErr:       sql.ErrNoRows,
		recordErr:    &pgconn.PgError{Code: "23505"},
		getAfterFail: existing,
		txSessions:   &sessionFake{},
	}
	service := newEventService(&sessionFake{}, events)

	response, err := service.CreateEvent(testContext(), validEvent(trackingModel.EventTypeStepShown))
	if err != nil {
		t.Fatalf("CreateEvent() error = %v", err)
	}
	if response != existing {
		t.Fatal("CreateEvent() did not return the concurrently created event")
	}
}

func validEvent(eventType trackingModel.EventType) *trackingModel.CreateEventRequest {
	return &trackingModel.CreateEventRequest{
		ID:         uuid.NewString(),
		SessionID:  uuid.NewString(),
		Type:       eventType,
		OccurredAt: "2026-08-06T10:11:12.123456789Z",
	}
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		t.Fatalf("parse time %q: %v", value, err)
	}
	return parsed
}

func testContext() context.Context {
	return context.WithValue(context.Background(), "projectKey", "project-key")
}

func newEventService(sessions *sessionFake, events *eventFake) *TrackingService {
	return NewTrackingService(sessions, events, &scenarioFake{}, &stepFake{})
}

type sessionFake struct {
	updateCalls      int
	updatedSessionID string
	updatedStatus    trackingModel.SessionStatus
	finishedAt       time.Time
	getSession       *trackingModel.OnboardingSession
	getSessionErr    error
}

func (*sessionFake) CreateSession(context.Context, *trackingModel.OnboardingSession) (*trackingModel.OnboardingSession, error) {
	return nil, errors.New("not implemented")
}

func (s *sessionFake) UpdateSessionStatus(_ context.Context, sessionID string, status trackingModel.SessionStatus, finishedAt time.Time) (*trackingModel.OnboardingSession, error) {
	s.updateCalls++
	s.updatedSessionID = sessionID
	s.updatedStatus = status
	s.finishedAt = finishedAt
	return &trackingModel.OnboardingSession{}, nil
}

func (*sessionFake) GetSessionByScenarioAndUser(context.Context, string, string) (*trackingModel.OnboardingSession, error) {
	return nil, errors.New("not implemented")
}

func (s *sessionFake) GetSessionById(context.Context, string) (*trackingModel.OnboardingSession, error) {
	if s.getSession != nil || s.getSessionErr != nil {
		return s.getSession, s.getSessionErr
	}
	return &trackingModel.OnboardingSession{
		ScenarioID: "scenario-1",
		Status:     trackingModel.SessionStatusActive,
	}, nil
}

type scenarioFake struct {
	scenario *trackingModel.Scenario
	err      error
}

func (s *scenarioFake) GetScenarioByIdAndProjectKey(context.Context, string, string) (*trackingModel.Scenario, error) {
	if s.scenario != nil || s.err != nil {
		return s.scenario, s.err
	}
	return &trackingModel.Scenario{Status: trackingModel.ScenarioStatusEnabled}, nil
}

type stepFake struct {
	step *trackingModel.Step
	err  error
}

func (s *stepFake) GetStepById(context.Context, string) (*trackingModel.Step, error) {
	return s.step, s.err
}

type eventFake struct {
	getResponse       *trackingModel.EventAcceptedResponse
	getAfterFail      *trackingModel.EventAcceptedResponse
	getErr            error
	recordErr         error
	recorded          *trackingModel.OnboardingEvent
	transactionCalled bool
	txSessions        *sessionFake
}

func (e *eventFake) RecordEvent(_ context.Context, event *trackingModel.OnboardingEvent) (*trackingModel.EventAcceptedResponse, error) {
	e.recorded = event
	if e.recordErr != nil {
		return nil, e.recordErr
	}
	return &trackingModel.EventAcceptedResponse{Event: *event}, nil
}

func (e *eventFake) GetEventById(context.Context, string) (*trackingModel.EventAcceptedResponse, error) {
	if e.getAfterFail != nil && e.recorded != nil {
		return e.getAfterFail, nil
	}
	if e.getResponse != nil {
		return e.getResponse, nil
	}
	return nil, e.getErr
}

func (e *eventFake) WithinTransaction(_ context.Context, fn func(SessionRepository, EventRepository) error) error {
	e.transactionCalled = true
	if e.txSessions == nil {
		e.txSessions = &sessionFake{}
	}
	return fn(e.txSessions, e)
}
