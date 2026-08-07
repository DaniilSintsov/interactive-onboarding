-- name: CreateSession :one
INSERT INTO onboarding.sessions
  ("id", "scenario_id", "user_id", "status", "started_at", "finished_at")
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSessionByScenarioAndUser :one
SELECT *
FROM onboarding.sessions
WHERE "id" = $1 AND "user_id" = $2;

-- name: SelectSessionById :one
SELECT *
FROM onboarding.sessions
WHERE "id" = $1;

-- name: ChangeSessionStatus :one
UPDATE onboarding.sessions
SET "status" = $1, "finished_at" = $2
WHERE "id" = $3 AND "status" = 'active'
RETURNING *;
