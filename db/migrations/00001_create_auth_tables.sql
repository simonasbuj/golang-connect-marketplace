-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA auth;

CREATE TYPE auth.user_role AS ENUM (
    'admin',
    'customer'
);

CREATE TABLE IF NOT EXISTS auth.users (
    id VARCHAR(30) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(60) NOT NULL,
    name VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    username VARCHAR(40) NOT NULL,
    role auth.user_role DEFAULT 'customer' NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS auth.refresh_tokens (
    token TEXT PRIMARY KEY,
    user_id VARCHAR(30) REFERENCES auth.users(id),
    expires_at TIMESTAMPTZ NOT NULL
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth.refresh_tokens;
DROP TABLE IF EXISTS auth.users;
DROP TYPE IF EXISTS auth.user_role;
DROP SCHEMA IF EXISTS auth;
-- +goose StatementEnd
