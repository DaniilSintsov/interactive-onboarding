-- name: CreateScenarioTestToken :one
INSERT INTO onboarding.scenario_test_tokens (
    scenario_id,
    hash,
    expires_at
)
VALUES (
    sqlc.arg(scenario_id),
    sqlc.arg(hash),
    sqlc.arg(expires_at)
)
RETURNING id, scenario_id, hash, created_at, expires_at;
