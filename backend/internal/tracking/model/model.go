package model

import "time"

type ResolveScenarioRequest struct {
	Page   string `json:"page"`
	UserID string `json:"user_id"`
}

type StartSessionRequest struct {
	ScenarioID string `json:"scenario_id"`
	UserID     string `json:"user_id"`
}

type CreateEventRequest struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	StepID     *string   `json:"step_id"`
	Type       EventType `json:"type"`
	Data       []byte    `json:"data"`
	OccurredAt string    `json:"occurred_at"`
}

type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusSkipped   SessionStatus = "skipped"
)

type OnboardingSession struct {
	ID         string        `json:"id"`
	ScenarioID string        `json:"scenario_id"`
	UserID     string        `json:"user_id"`
	Status     SessionStatus `json:"status"`
	StartedAt  time.Time     `json:"started_at"`
	FinishedAt *time.Time    `json:"finished_at"`
}

type EventType string

const (
	EventTypeStepShown           EventType = "step_shown"
	EventTypeStepCompleted       EventType = "step_completed"
	EventTypeStepSkipped         EventType = "step_skipped"
	EventTypeOnboardingCompleted EventType = "onboarding_completed"
	EventTypeOnboardingSkipped   EventType = "onboarding_skipped"
)

type OnboardingEvent struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	StepID     *string   `json:"step_id"`
	Type       EventType `json:"type"`
	Data       []byte    `json:"data"`
	OccurredAt time.Time `json:"occurred_at"`
	ReceivedAt time.Time `json:"received_at"`
}

type EventAcceptedResponse struct {
	Event     OnboardingEvent `json:"event"`
	Duplicate bool            `json:"duplicate"`
}

type SessionCompletion struct {
	Status     SessionStatus
	FinishedAt string
}
