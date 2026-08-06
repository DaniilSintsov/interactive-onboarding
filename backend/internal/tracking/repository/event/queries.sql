-- name: CreateEvent :one
INSERT INTO "OnboardingEvent"
  ("event_id", "session_id", "step_id", "type", "data", "occurred_at", "received_at")
VALUES
  ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetEventById :one
SELECT *
FROM "OnboardingEvent"
WHERE "event_id" = $1;
