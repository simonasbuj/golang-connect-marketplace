-- +goose Up
-- +goose StatementBegin
CREATE TYPE auth.oauth_providers AS ENUM (
    'github',
    'google'
);

CREATE TABLE IF NOT EXISTS auth.oauth_users (
    id VARCHAR(30) PRIMARY KEY,
    user_id VARCHAR(30) NOT NULL 
        REFERENCES auth.users(id),
    provider_user_id VARCHAR(50) NOT NULL,
    provider auth.oauth_providers NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth.oauth_users;
DROP TYPE IF EXISTS auth.oauth_providers;
-- +goose StatementEnd
