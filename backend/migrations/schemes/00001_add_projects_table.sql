CREATE SCHEMA IF NOT EXISTS onboarding

CREATE TABLE onboarding.projects
(
    ID         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    Name       TEXT        NOT NULL,
    ProjectKey TEXT        NOT NULL UNIQUE,
    CreatedAt  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UpdatedAt  TIMESTAMPTZ NOT NULL DEFAULT NOW()
)