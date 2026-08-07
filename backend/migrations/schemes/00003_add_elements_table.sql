-- +goose Up

CREATE TABLE onboarding.elements
(
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    project_id  UUID        NOT NULL,
    key         TEXT        NOT NULL,
    label       TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT elements_project_fk
        FOREIGN KEY (project_id)
            REFERENCES onboarding.projects (id)
            ON DELETE RESTRICT,

    CONSTRAINT elements_key_length
        CHECK (char_length(key) BETWEEN 1 AND 255),

    CONSTRAINT elements_label_length
        CHECK (char_length(label) BETWEEN 1 AND 255),

    CONSTRAINT elements_description_length
        CHECK (char_length(description) <= 2000)
);

CREATE UNIQUE INDEX elements_project_id_key_unique
    ON onboarding.elements (project_id, key)
    WHERE deleted_at IS NULL;

CREATE TRIGGER elements_set_updated_at
    BEFORE UPDATE
    ON onboarding.elements
    FOR EACH ROW
    EXECUTE FUNCTION onboarding.set_updated_at();

-- +goose Down

DROP TABLE onboarding.elements;
