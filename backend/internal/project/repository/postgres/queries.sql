-- name: CreateProject :one
INSERT INTO onboarding.projects (name, project_key)
VALUES (sqlc.arg(name), sqlc.arg(project_key))
RETURNING id, name, project_key, created_at, updated_at;

-- name: ListProjects :many
SELECT id, name, project_key, created_at, updated_at
FROM onboarding.projects
WHERE deleted_at IS NULL
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(limit) OFFSET sqlc.arg(offset);

-- name: CountProjects :one
SELECT COUNT(*)
FROM onboarding.projects
WHERE deleted_at IS NULL;

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

-- name: ListElementsByProjectID :many
SELECT id, project_id, key, label, description, created_at, updated_at
FROM onboarding.elements
WHERE project_id = sqlc.arg(project_id)
  AND deleted_at IS NULL
ORDER BY created_at, id;

-- name: CreateElement :one
INSERT INTO onboarding.elements (project_id, key, label, description)
VALUES (sqlc.arg(project_id), sqlc.arg(key), sqlc.arg(label), sqlc.arg(description))
RETURNING id, project_id, key, label, description, created_at, updated_at;

-- name: UpdateElement :one
UPDATE onboarding.elements
SET key         = COALESCE(sqlc.narg(key), key),
    label       = COALESCE(sqlc.narg(label), label),
    description = COALESCE(sqlc.narg(description), description)
WHERE project_id = sqlc.arg(project_id)
  AND id = sqlc.arg(element_id)
  AND deleted_at IS NULL
RETURNING id, project_id, key, label, description, created_at, updated_at;

-- name: DeleteElement :one
UPDATE onboarding.elements
SET deleted_at = NOW()
WHERE project_id = sqlc.arg(project_id)
  AND id = sqlc.arg(element_id)
  AND deleted_at IS NULL
RETURNING id;
