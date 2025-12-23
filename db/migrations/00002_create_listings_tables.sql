-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS listings;

CREATE TYPE listings.listing_status AS ENUM (
    'open',
    'sold',
    'canceled',
    'refunded'
);

CREATE TABLE IF NOT EXISTS listings.categories (
    id VARCHAR(30) PRIMARY KEY,
    title VARCHAR(30) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS listings.listings (
    id VARCHAR(30) PRIMARY KEY,
    user_id VARCHAR(30) NOT NULL
        REFERENCES auth.users(id),
    category_id VARCHAR(30)
        REFERENCES listings.categories(id)
        ON DELETE SET NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    price_in_cents INTEGER NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status listings.listing_status NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS listings.listings_images (
    id VARCHAR(30) PRIMARY KEY,
    listing_id VARCHAR(30)
        REFERENCES listings.listings(id)
        ON DELETE SET NULL,
    path VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS listings.listings_images;
DROP TABLE IF EXISTS listings.listings;
DROP TABLE IF EXISTS listings.categories;
DROP TYPE IF EXISTS listings.listing_status;
DROP SCHEMA IF EXISTS listings;
-- +goose StatementEnd
