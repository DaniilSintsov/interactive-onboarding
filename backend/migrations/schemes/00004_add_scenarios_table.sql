-- +goose Up

CREATE TYPE onboarding.scenario_status AS ENUM ('enabled', 'disabled', 'in_development');

CREATE TABLE onboarding.scenarios
(
    id           UUID PRIMARY KEY                    DEFAULT gen_random_uuid(),
    project_id   UUID                       NOT NULL,
    name         TEXT                       NOT NULL,
    description  TEXT                       NOT NULL DEFAULT '',
    page_pattern TEXT                       NOT NULL,
    status       onboarding.scenario_status NOT NULL DEFAULT 'in_development',
    published_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ                NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ                NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    CONSTRAINT scenario_project_fk
        FOREIGN KEY (project_id)
            REFERENCES onboarding.projects (id)
            ON DELETE RESTRICT,

    CONSTRAINT scenario_name_length
        CHECK (char_length(name) BETWEEN 1 AND 255),

    CONSTRAINT scenario_page_pattern_length
        CHECK (char_length(page_pattern) BETWEEN 1 AND 2048),

    CONSTRAINT scenario_description_length
        CHECK (char_length(description) <= 2000)
);

CREATE INDEX scenarios_active_project_status_idx
    ON onboarding.scenarios (project_id, status) WHERE deleted_at IS NULL;

CREATE INDEX scenarios_runtime_lookup_idx
    ON onboarding.scenarios (
                             project_id,
                             page_pattern,
                             published_at,
                             id
        ) WHERE status = 'enabled'
      AND deleted_at IS NULL;

CREATE TRIGGER scenarios_set_updated_at
    BEFORE UPDATE
    ON onboarding.scenarios
    FOR EACH ROW
    EXECUTE FUNCTION onboarding.set_updated_at();

-- +goose Down

DROP TABLE onboarding.scenarios;
DROP TYPE onboarding.scenario_status;
