-- +goose Up

CREATE SCHEMA onboarding;

-- +goose StatementBegin
CREATE FUNCTION onboarding.set_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- +goose Down

DROP FUNCTION onboarding.set_updated_at();
DROP SCHEMA onboarding;
