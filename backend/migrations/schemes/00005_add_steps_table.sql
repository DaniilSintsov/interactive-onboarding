-- +goose Up
CREATE TABLE onboarding.steps
(
    id            UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    scenario_id   UUID        NOT NULL,
    element_id    UUID        NOT NULL,
    step_num      INTEGER     NOT NULL,
    title         TEXT        NOT NULL,
    description   TEXT        NOT NULL DEFAULT '',
    frontend_data JSONB        NOT NULL DEFAULT '{}',
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

    CONSTRAINT steps_title_length
        CHECK (char_length(title) BETWEEN 1 AND 255),

    CONSTRAINT steps_description_length
        CHECK (char_length(description) <= 2000),

    CONSTRAINT steps_step_num_positive
        CHECK (step_num >= 1)
);

CREATE UNIQUE INDEX steps_scenario_id_step_num_unique
    ON onboarding.steps (scenario_id, step_num)
    WHERE deleted_at IS NULL;

CREATE INDEX steps_scenario_id_idx ON onboarding.steps (scenario_id);
CREATE INDEX steps_element_id_idx ON onboarding.steps (element_id);

CREATE TRIGGER steps_set_updated_at
    BEFORE UPDATE
    ON onboarding.steps
    FOR EACH ROW
    EXECUTE FUNCTION onboarding.set_updated_at();

-- +goose Down
DROP TABLE onboarding.steps;