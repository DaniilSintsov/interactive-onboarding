CREATE SCHEMA IF NOT EXISTS onboarding

CREATE TABLE onboarding.elements
(
    ID          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    ProjectID   UUID        NOT NULL REFERENCES onboarding.project (id) ON DELETE CASCADE,
    Key         TEXT        NOT NULL,
    Label       TEXT        NOT NULL,
    Description TEXT,
    CreatedAt   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UpdatedAt   TIMESTAMPTZ NOT NULL DEFAULT NOW()
)