-- name: CreateProject :one
INSERT INTO onboarding.projects (name, project_key)
VALUES (sqlc.arg(name), sqlc.arg(project_key))
RETURNING id, name, project_key, created_at, updated_at;

-- name: ListProjects :many
SELECT id, name, project_key, created_at, updated_at, COUNT(*) OVER ()::bigint AS total
FROM onboarding.projects
WHERE deleted_at IS NULL
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(page_limit)::integer
OFFSET sqlc.arg(page_offset)::integer;

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