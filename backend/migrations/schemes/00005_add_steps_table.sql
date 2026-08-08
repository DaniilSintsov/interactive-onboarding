-- +goose Up

CREATE TABLE onboarding.steps
(
    id            UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    scenario_id   UUID        NOT NULL,
    element_id    UUID        NOT NULL,
    step_num      INTEGER     NOT NULL,
    title         TEXT        NOT NULL,
    description   TEXT        NOT NULL DEFAULT '',
    frontend_data JSONB       NOT NULL DEFAULT '{}'::JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,

    CONSTRAINT steps_scenario_fk
        FOREIGN KEY (scenario_id)
            REFERENCES onboarding.scenarios (id)
            ON DELETE RESTRICT,

    CONSTRAINT steps_element_fk
        FOREIGN KEY (element_id)
            REFERENCES onboarding.elements (id)
            ON DELETE RESTRICT,

    CONSTRAINT steps_positive_step_num
        CHECK (step_num > 0),

    CONSTRAINT step_title_length
        CHECK (char_length(title) BETWEEN 1 AND 255),

    CONSTRAINT step_description_length
        CHECK (char_length(description) <= 2000)
);

CREATE UNIQUE INDEX steps_active_scenario_step_num_unique
    ON onboarding.steps (scenario_id, step_num) WHERE deleted_at IS NULL;

CREATE INDEX steps_active_element_id_idx
    ON onboarding.steps (element_id) WHERE deleted_at IS NULL;

CREATE TRIGGER steps_set_updated_at
    BEFORE UPDATE
    ON onboarding.steps
    FOR EACH ROW
    EXECUTE FUNCTION onboarding.set_updated_at();

-- +goose Down

DROP TABLE onboarding.steps;
