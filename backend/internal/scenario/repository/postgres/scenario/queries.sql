-- name: CreateScenario :one
INSERT INTO onboarding.scenarios (
    project_id,
    name,
    description,
    page_pattern
)
VALUES (
    sqlc.arg(project_id),
    sqlc.arg(name),
    sqlc.arg(description),
    sqlc.arg(page_pattern)
)
RETURNING id,
          project_id,
          name,
          description,
          page_pattern,
          status,
          published_at,
          created_at,
          updated_at;

-- name: ListScenariosByProjectID :many
SELECT s.id,
       s.project_id,
       s.name,
       s.description,
       s.page_pattern,
       s.status,
       s.published_at,
       s.created_at,
       s.updated_at,
       (
           SELECT COUNT(*)::bigint
           FROM onboarding.steps AS st
           WHERE st.scenario_id = s.id
             AND st.deleted_at IS NULL
       ) AS steps_count,
       COUNT(*) OVER ()::bigint AS total
FROM onboarding.scenarios AS s
WHERE s.project_id = sqlc.arg(project_id)
  AND s.deleted_at IS NULL
  AND (
      sqlc.narg(status)::text IS NULL
      OR s.status::text = sqlc.narg(status)::text
  )
ORDER BY s.created_at DESC, s.id DESC
LIMIT sqlc.arg(page_limit)::integer
OFFSET sqlc.arg(page_offset)::integer;

-- name: GetScenarioByID :one
SELECT id,
       project_id,
       name,
       description,
       page_pattern,
       status,
       published_at,
       created_at,
       updated_at
FROM onboarding.scenarios
WHERE id = sqlc.arg(scenario_id)
  AND deleted_at IS NULL;

-- name: UpdateScenario :one
UPDATE onboarding.scenarios
SET name         = COALESCE(sqlc.narg(name), name),
    description  = COALESCE(sqlc.narg(description), description),
    page_pattern = COALESCE(sqlc.narg(page_pattern), page_pattern)
WHERE id = sqlc.arg(scenario_id)
  AND deleted_at IS NULL
RETURNING id,
          project_id,
          name,
          description,
          page_pattern,
          status,
          published_at,
          created_at,
          updated_at;

-- name: DeleteScenario :one
UPDATE onboarding.scenarios
SET deleted_at = NOW()
WHERE id = sqlc.arg(scenario_id)
  AND deleted_at IS NULL
RETURNING id;

-- name: LockActiveScenario :one
SELECT id,
       project_id,
       name,
       description,
       page_pattern,
       status,
       published_at,
       created_at,
       updated_at
FROM onboarding.scenarios
WHERE id = sqlc.arg(scenario_id)
  AND deleted_at IS NULL
FOR UPDATE;

-- name: PublishScenario :one
UPDATE onboarding.scenarios
SET status       = 'enabled',
    published_at = COALESCE(published_at, NOW())
WHERE id = sqlc.arg(scenario_id)
  AND status = 'in_development'
  AND deleted_at IS NULL
RETURNING id,
          project_id,
          name,
          description,
          page_pattern,
          status,
          published_at,
          created_at,
          updated_at;

-- name: EnableScenario :one
UPDATE onboarding.scenarios
SET status = 'enabled'
WHERE id = sqlc.arg(scenario_id)
  AND status = 'disabled'
  AND deleted_at IS NULL
RETURNING id,
          project_id,
          name,
          description,
          page_pattern,
          status,
          published_at,
          created_at,
          updated_at;

-- name: DisableScenario :one
UPDATE onboarding.scenarios
SET status = 'disabled'
WHERE id = sqlc.arg(scenario_id)
  AND status = 'enabled'
  AND deleted_at IS NULL
RETURNING id,
          project_id,
          name,
          description,
          page_pattern,
          status,
          published_at,
          created_at,
          updated_at;
