-- +goose Up
CREATE TABLE onboarding.scenarios
(
    id           UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID       NOT NULL,
    name         TEXT       NOT NULL,
    description  TEXT       NOT NULL DEFAULT '',
    page_pattern TEXT       NOT NULL,
    status       TEXT       NOT NULL DEFAULT 'in_development',
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    CONSTRAINT scenarios_project_fk
    FOREIGN KEY (project_id) 
    REFERENCES onboarding.projects(id)
    ON DELETE RESTRICT,

    CONSTRAINT scenarios_name_length
    CHECK (char_length(name) BETWEEN 1 AND 255),

    CONSTRAINT scenarios_page_pattern_length
    CHECK (char_length(page_pattern) BETWEEN 1 AND 2000),

    CONSTRAINT scenarios_status_check
    CHECK (status IN ('in_development', 'enabled', 'disabled'))
);

CREATE INDEX scenarios_project_id_idx ON onboarding.scenarios (project_id);
CREATE INDEX scenarios_status_idx ON onboarding.scenarios (status);

CREATE TRIGGER scenarios_set_updated_at
BEFORE UPDATE
ON onboarding.scenarios
FOR EACH ROW
EXECUTE FUNCTION onboarding.set_updated_at();

-- +goose Down
DROP TABLE onboarding.scenarios;