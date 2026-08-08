-- name: CreateEvent :one
WITH "active_session" AS (
  SELECT "id"
  FROM onboarding.sessions
  WHERE "id" = $2
    AND "status" = 'active'
  FOR UPDATE
)
INSERT INTO onboarding.events
  ("id", "session_id", "step_id", "type", "data", "occurred_at", "received_at")
SELECT
  $1, $2, $3, $4, $5, $6, $7
FROM "active_session"
RETURNING *;

-- name: GetEventById :one
SELECT *
FROM onboarding.events
WHERE "id" = $1;
