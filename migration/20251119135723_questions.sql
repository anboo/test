-- +goose Up
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
CREATE INDEX idx_questions_created_at ON questions (created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_questions_created_at;
DROP TABLE IF EXISTS questions;
