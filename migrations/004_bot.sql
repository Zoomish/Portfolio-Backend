-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    admin BOOLEAN NOT NULL DEFAULT FALSE,
    name TEXT NOT NULL,
    username TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE telegram_session (
    user_id BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    waiting_for TEXT NOT NULL DEFAULT '',
    last_bot_message_id BIGINT,
    error_message_id BIGINT
);

CREATE TABLE user_event (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    kind TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX user_event_kind_created_at_idx ON user_event (kind, created_at);

-- +goose Down
DROP TABLE IF EXISTS user_event;
DROP TABLE IF EXISTS telegram_session;
DROP TABLE IF EXISTS users;
