-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS payments;

CREATE TYPE payments.provider AS ENUM ('stripe', 'klix', 'polar.sh');

CREATE TABLE IF NOT EXISTS payments.seller_accounts (
    id varchar(50) PRIMARY KEY,
    user_id VARCHAR(30) NOT NULL 
        REFERENCES auth.users(id),
    provider payments.provider NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS payments.payments (
    id varchar(30) PRIMARY KEY,
    listing_id VARCHAR(30) NOT NULL
        REFERENCES listings.listings(id),
    buyer_id VARCHAR(30) NOT NULL
        REFERENCES auth.users(id),
    provider_payment_id VARCHAR(50) NOT NULL,
    provider payments.provider NOT NULL,
    amount_in_cents INTEGER NOT NULL,
    fee_amount_in_cents INTEGER NOT NULL,
    currency varchar(3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    refunded_at TIMESTAMPTZ
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments.payments;
DROP TABLE IF EXISTS payments.seller_accounts;
DROP TYPE IF EXISTS payments.provider;
-- +goose StatementEnd
