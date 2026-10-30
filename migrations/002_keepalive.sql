-- +goose Up
CREATE TABLE IF NOT EXISTS keepalive (
    id SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    touched_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    seq BIGINT NOT NULL DEFAULT 0
);

INSERT INTO keepalive (id)
VALUES (1)
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS keepalive;
