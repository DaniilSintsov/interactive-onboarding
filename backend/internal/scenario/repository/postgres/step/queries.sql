-- name: CreateStep :one
INSERT INTO onboarding.steps (
    scenario_id,
    element_id,
    step_num,
    title,
    description,
    frontend_data
)
VALUES (
    sqlc.arg(scenario_id),
    sqlc.arg(element_id),
    sqlc.arg(step_num),
    sqlc.arg(title),
    sqlc.arg(description),
    sqlc.arg(frontend_data)
)
RETURNING id,
          scenario_id,
          element_id,
          step_num,
          title,
          description,
          frontend_data,
          created_at,
          updated_at;

-- name: UpdateStep :one
UPDATE onboarding.steps
SET element_id    = COALESCE(sqlc.narg(element_id), element_id),
    title         = COALESCE(sqlc.narg(title), title),
    description   = COALESCE(sqlc.narg(description), description),
    frontend_data = COALESCE(sqlc.narg(frontend_data), frontend_data)
WHERE scenario_id = sqlc.arg(scenario_id)
  AND id = sqlc.arg(step_id)
  AND deleted_at IS NULL
    RETURNING id,
          scenario_id,
          element_id,
          step_num,
          title,
          description,
          frontend_data,
          created_at,
          updated_at;

-- name: DeleteStep :one
UPDATE onboarding.steps
SET deleted_at = NOW()
WHERE scenario_id = sqlc.arg(scenario_id)
  AND id = sqlc.arg(step_id)
  AND deleted_at IS NULL
    RETURNING id, step_num;

-- name: ListStepsByScenarioID :many
SELECT id,
       scenario_id,
       element_id,
       step_num,
       title,
       description,
       frontend_data,
       created_at,
       updated_at
FROM onboarding.steps
WHERE scenario_id = sqlc.arg(scenario_id)
  AND deleted_at IS NULL
ORDER BY step_num, id;

-- name: IsElementUsedBySteps :one
SELECT EXISTS (
    SELECT 1
    FROM onboarding.steps AS st
    JOIN onboarding.scenarios AS s ON s.id = st.scenario_id
    WHERE st.element_id = sqlc.arg(element_id)
      AND st.deleted_at IS NULL
      AND s.deleted_at IS NULL
) AS is_used;

-- name: GetNextStepNumber :one
SELECT (COALESCE(MAX(step_num), 0) + 1)::integer
FROM onboarding.steps
WHERE scenario_id = sqlc.arg(scenario_id)
  AND deleted_at IS NULL;

-- name: GetStepByID :one
SELECT id,
       scenario_id,
       element_id,
       step_num,
       title,
       description,
       frontend_data,
       created_at,
       updated_at
FROM onboarding.steps
WHERE scenario_id = sqlc.arg(scenario_id)
  AND id = sqlc.arg(step_id)
  AND deleted_at IS NULL;

-- name: LockActiveStep :one
SELECT id,
       scenario_id,
       element_id,
       step_num,
       title,
       description,
       frontend_data,
       created_at,
       updated_at
FROM onboarding.steps
WHERE scenario_id = sqlc.arg(scenario_id)
  AND id = sqlc.arg(step_id)
  AND deleted_at IS NULL
FOR UPDATE;

-- name: LockActiveStepsByScenarioID :many
SELECT id, step_num
FROM onboarding.steps
WHERE scenario_id = sqlc.arg(scenario_id)
  AND deleted_at IS NULL
ORDER BY step_num, id
FOR UPDATE;

-- name: MoveStepsOutOfOrderRange :exec
WITH current_range AS (
    SELECT COALESCE(MAX(step_num), 0)::integer AS offset
    FROM onboarding.steps
    WHERE scenario_id = sqlc.arg(scenario_id)
      AND deleted_at IS NULL
)
UPDATE onboarding.steps AS s
SET step_num = s.step_num + current_range.offset
FROM current_range
WHERE s.scenario_id = sqlc.arg(scenario_id)
  AND s.deleted_at IS NULL
  AND current_range.offset > 0;

-- name: ApplyStepOrder :execrows
WITH requested_order AS (
    SELECT input.step_id,
           input.ordinality::integer AS step_num
    FROM unnest(sqlc.arg(ordered_step_ids)::uuid[])
             WITH ORDINALITY AS input(step_id, ordinality)
)
UPDATE onboarding.steps AS s
SET step_num = requested_order.step_num
FROM requested_order
WHERE s.scenario_id = sqlc.arg(scenario_id)
  AND s.id = requested_order.step_id
  AND s.deleted_at IS NULL;
