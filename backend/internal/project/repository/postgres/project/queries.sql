-- name: CreateProject :one
INSERT INTO onboarding.projects (name, project_key)
VALUES (sqlc.arg(name), sqlc.arg(project_key))
RETURNING id, name, project_key, created_at, updated_at;

-- name: ListProjects :many
WITH filtered AS (
    SELECT id, name, project_key, created_at, updated_at
    FROM onboarding.projects
    WHERE deleted_at IS NULL
),
total AS (
    SELECT COUNT(*)::bigint AS value
    FROM filtered
),
page AS (
    SELECT id, name, project_key, created_at, updated_at
    FROM filtered
    ORDER BY created_at DESC, id DESC
    LIMIT sqlc.arg(page_limit)::integer
    OFFSET sqlc.arg(page_offset)::integer
)
SELECT page.id,
       page.name,
       page.project_key,
       page.created_at,
       page.updated_at,
       total.value AS total
FROM total
LEFT JOIN page ON TRUE
ORDER BY page.created_at DESC NULLS LAST, page.id DESC NULLS LAST;

-- name: GetProjectByID :one
SELECT id, name, project_key, created_at, updated_at
FROM onboarding.projects
WHERE id = sqlc.arg(project_id)
  AND deleted_at IS NULL;

-- name: GetProjectIDByKey :one
SELECT id
FROM onboarding.projects
WHERE project_key = sqlc.arg(project_key)
  AND deleted_at IS NULL;

-- name: UpdateProject :one
UPDATE onboarding.projects
SET name = sqlc.arg(name)
WHERE id = sqlc.arg(project_id)
  AND deleted_at IS NULL
RETURNING id, name, project_key, created_at, updated_at;

-- name: DeleteProject :one
UPDATE onboarding.projects
SET deleted_at = NOW()
WHERE id = sqlc.arg(project_id)
  AND deleted_at IS NULL
RETURNING id;

-- name: LockActiveProject :one
SELECT id
FROM onboarding.projects
WHERE id = sqlc.arg(project_id)
  AND deleted_at IS NULL
    FOR SHARE;
