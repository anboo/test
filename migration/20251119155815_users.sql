-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR,
    password VARCHAR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE UNIQUE INDEX udx_users_username ON users (username);

-- +goose Down
DROP INDEX udx_users_username;
DROP TABLE users;
