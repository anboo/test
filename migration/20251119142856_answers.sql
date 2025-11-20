-- +goose Up
CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    question_id INT NOT NULL,
    user_id TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_answers_question_id ON answers (question_id);
CREATE INDEX idx_answers_created_at ON answers (created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_answers_question_id;
DROP INDEX IF EXISTS idx_answers_created_at;
DROP TABLE IF EXISTS answers;
