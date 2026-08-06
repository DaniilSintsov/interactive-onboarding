-- name: CreateSession :one
INSERT INTO "OnboardingSession"
  ("session_id", "scenario_id", "user_id", "status", "started_at", "finished_at")
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSessionByScenarioAndUser :one
SELECT *
FROM "OnboardingSession"
WHERE "scenario_id" = $1 AND "user_id" = $2;

-- name: ChangeSessionStatus :one
UPDATE "OnboardingSession"
SET "status" = $1, "finished_at" = $2
WHERE "session_id" = $3
RETURNING *;
