-- +goose Up
CREATE TABLE IF NOT EXISTS sessions (
    id         UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL, -- sha256
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    last_used  TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ
);

CREATE INDEX sessions_user_id_idx
ON sessions(user_id);

-- +goose Down
DROP INDEX sessions_user_id_idx;
DROP TABLE sessions;