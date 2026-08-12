-- +goose Up
INSERT INTO onboarding.projects (name, project_key)
VALUES ('example', 'pk_demo_avito')
ON CONFLICT (project_key) DO NOTHING;

-- +goose Down
DELETE FROM onboarding.projects
WHERE project_key = 'pk_demo_avito';