-- name: CreateEvent :one
WITH "active_session" AS (
  SELECT "session_id"
  FROM "OnboardingSession"
  WHERE "session_id" = $2
    AND "status" = 'active'
  FOR UPDATE
)
INSERT INTO "OnboardingEvent"
  ("event_id", "session_id", "step_id", "type", "data", "occurred_at", "received_at")
SELECT
  $1, $2, $3, $4, $5, $6, $7
FROM "active_session"
RETURNING *;

-- name: GetEventById :one
SELECT *
FROM "OnboardingEvent"
WHERE "event_id" = $1;
