-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA auth;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS auth;
-- +goose StatementEnd
