package event

import (
	"context"

	trackingModel "github.com/DaniilSintsov/interactive-onboarding/internal/tracking/model"
	event "github.com/DaniilSintsov/interactive-onboarding/internal/tracking/repository/event/sqlc"
	sessionRepository "github.com/DaniilSintsov/interactive-onboarding/internal/tracking/repository/session"
	trackingService "github.com/DaniilSintsov/interactive-onboarding/internal/tracking/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type EventRepository struct {
	conn    *pgx.Conn
	queries *event.Queries
}

func NewEventRepository(db *pgx.Conn) *EventRepository {
	q := event.New(db)
	return &EventRepository{
		conn:    db,
		queries: q,
	}
}

func (e *EventRepository) WithinTransaction(
	ctx context.Context,
	fn func(trackingService.SessionRepository, trackingService.EventRepository) error,
) (err error) {
	tx, err := e.conn.Begin(ctx)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	txEvents := &EventRepository{queries: e.queries.WithTx(tx)}
	txSessions := sessionRepository.NewSessionRepository(tx)
	if err = fn(txSessions, txEvents); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	committed = true
	return nil
}

func (e *EventRepository) RecordEvent(
	ctx context.Context, onboarding *trackingModel.OnboardingEvent,
) (*trackingModel.EventAcceptedResponse, error) {
	eventId, err := uuid.Parse(onboarding.ID)
	if err != nil {
		return nil, err
	}
	sessionId, err := uuid.Parse(onboarding.SessionID)
	if err != nil {
		return nil, err
	}
	var stepId pgtype.UUID
	if onboarding.StepID != nil {
		step, err := uuid.Parse(*onboarding.StepID)
		if err != nil {
			return nil, err
		}
		stepId = pgtype.UUID{
			Bytes: step,
			Valid: true,
		}
	}
	createEvent := event.CreateEventParams{
		EventID:    eventId,
		SessionID:  sessionId,
		StepID:     stepId,
		Type:       string(onboarding.Type),
		Data:       onboarding.Data,
		OccurredAt: pgtype.Timestamp{Time: onboarding.OccurredAt, Valid: true},
		ReceivedAt: pgtype.Timestamp{Time: onboarding.ReceivedAt, Valid: true},
	}

	createdEvent, err := e.queries.CreateEvent(ctx, createEvent)
	if err != nil {
		return nil, err
	}

	return adaptEvent(createdEvent, false), nil
}

func (e *EventRepository) GetEventById(
	ctx context.Context, eventID string,
) (*trackingModel.EventAcceptedResponse, error) {
	parsedEventID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, err
	}

	found, err := e.queries.GetEventById(ctx, parsedEventID)
	if err != nil {
		return nil, err
	}

	return adaptEvent(found, true), nil
}

func adaptEvent(source event.OnboardingEvent, duplicate bool) *trackingModel.EventAcceptedResponse {
	var stepID *string
	if source.StepID.Valid {
		value := uuid.UUID(source.StepID.Bytes).String()
		stepID = &value
	}

	return &trackingModel.EventAcceptedResponse{
		Event: trackingModel.OnboardingEvent{
			ID:         source.EventID.String(),
			SessionID:  source.SessionID.String(),
			StepID:     stepID,
			Type:       trackingModel.EventType(source.Type),
			Data:       source.Data,
			OccurredAt: source.OccurredAt.Time,
			ReceivedAt: source.ReceivedAt.Time,
		},
		Duplicate: duplicate,
	}
}
