package session

import (
	"context"
	"fmt"
	"time"

	trackingModel "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/model"
	session "github.com/DaniilSintsov/interactive-onboarding/backend/internal/tracking/repository/session/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SessionRepository struct {
	queries *session.Queries
}

func NewSessionRepository(db session.DBTX) *SessionRepository {
	q := session.New(db)
	return &SessionRepository{
		queries: q,
	}
}

func (r *SessionRepository) CreateSession(
	ctx context.Context, onboarding *trackingModel.OnboardingSession,
) (*trackingModel.OnboardingSession, error) {
	if onboarding == nil {
		return nil, fmt.Errorf("onboarding session is required")
	}

	sessionID, err := uuid.Parse(onboarding.ID)
	if err != nil {
		return nil, err
	}
	scenarioID, err := uuid.Parse(onboarding.ScenarioID)
	if err != nil {
		return nil, err
	}
	created, err := r.queries.CreateSession(ctx, session.CreateSessionParams{
		ID:         sessionID,
		ScenarioID: scenarioID,
		UserID:     onboarding.UserID,
		Status:     string(onboarding.Status),
		StartedAt:  pgtype.Timestamptz{Time: onboarding.StartedAt, Valid: true},
		FinishedAt: timestamp(onboarding.FinishedAt),
	})
	if err != nil {
		return nil, err
	}

	return adaptSession(created), nil
}

func (r *SessionRepository) UpdateSessionStatus(
	ctx context.Context, sessionID string, status trackingModel.SessionStatus, finishedAt time.Time,
) (*trackingModel.OnboardingSession, error) {
	parsedSessionID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, err
	}
	updated, err := r.queries.ChangeSessionStatus(ctx, session.ChangeSessionStatusParams{
		ID:         parsedSessionID,
		Status:     string(status),
		FinishedAt: timestamp(&finishedAt),
	})
	if err != nil {
		return nil, err
	}

	return adaptSession(updated), nil
}

func (r *SessionRepository) GetSessionByScenarioAndUser(
	ctx context.Context, scenarioID, userID string,
) (*trackingModel.OnboardingSession, error) {
	parsedScenarioID, err := uuid.Parse(scenarioID)
	if err != nil {
		return nil, err
	}
	found, err := r.queries.GetSessionByScenarioAndUser(ctx, session.GetSessionByScenarioAndUserParams{
		ScenarioID: parsedScenarioID,
		UserID:     userID,
	})
	if err != nil {
		return nil, err
	}

	return adaptSession(found), nil
}

func (r *SessionRepository) GetSessionById(
	ctx context.Context, sessionID string,
) (*trackingModel.OnboardingSession, error) {
	parsedSessionID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, err
	}

	found, err := r.queries.SelectSessionById(ctx, parsedSessionID)
	if err != nil {
		return nil, err
	}

	return adaptSession(found), nil
}

func adaptSession(source session.OnboardingSession) *trackingModel.OnboardingSession {
	return &trackingModel.OnboardingSession{
		ID:         source.ID.String(),
		ScenarioID: source.ScenarioID.String(),
		UserID:     source.UserID,
		Status:     trackingModel.SessionStatus(source.Status),
		StartedAt:  source.StartedAt.Time,
		FinishedAt: nullableTime(source.FinishedAt),
	}
}

func timestamp(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func nullableTime(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}
