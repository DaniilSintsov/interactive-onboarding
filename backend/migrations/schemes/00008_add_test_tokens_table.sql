-- +goose Up
CREATE TABLE onboarding.scenario_test_tokens
(
    id          UUID PRIMARY KEY     DEFAULT DEFAULT gen_random_uuid(),
    scenario_id UUID        NOT NULL,
    hash        BYTEA       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL,

    CONSTRAINT scenario_test_tokens_scenario_fk
        FOREIGN KEY (scenario_id)
            REFERENCES onboarding.scenarios (id)
            ON DELETE RESTRICT,

    CONSTRAINT scenario_test_tokens_hash_unique
        UNIQUE (token_hash),

    CONSTRAINT scenario_test_tokens_hash_length
        CHECK (octet_length(token_hash) = 32),

    CONSTRAINT scenario_test_tokens_expiration_valid
        CHECK (expires_at > created_at)
);


CREATE INDEX scenario_test_tokens_scenario_id_idx
    ON onboarding.scenario_test_tokens (scenario_id);

CREATE INDEX scenario_test_tokens_expires_at_idx
    ON onboarding.scenario_test_tokens (expires_at);

-- +goose Down

DROP TABLE oonboarding.scenario_test_tokens;
