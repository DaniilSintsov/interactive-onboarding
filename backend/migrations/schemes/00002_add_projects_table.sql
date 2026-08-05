CREATE TABLE onboarding.projects
(
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    project_key TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT projects_project_key_unique
        UNIQUE (project_key),

    CONSTRAINT projects_name_length
        CHECK (char_length(name) BETWEEN 1 AND 255)
);

CREATE TRIGGER projects_set_updated_at
    BEFORE UPDATE
    ON onboarding.projects
    FOR EACH ROW
EXECUTE FUNCTION onboarding.set_updated_at();