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
