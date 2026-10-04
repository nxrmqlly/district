-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    username    TEXT        NOT NULL UNIQUE,
    email       TEXT        NOT NULL UNIQUE,
    passwd_hash TEXT        NOT NULL,              -- Argon2id is text
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT username_length CHECK (
        char_length(username) >= 3 AND char_length(username) <= 24
    )
);

-- This makes josh, JOSH, JoSh, JOsh resolve to the same thing
CREATE UNIQUE INDEX users_username_lower_idx
    ON users (lower(username));

CREATE TABLE IF NOT EXISTS posts (
    id          BIGSERIAL   PRIMARY KEY,
    author_id   UUID        NOT NULL REFERENCES users(id),
    title       TEXT        NOT NULL,
    embed_url   TEXT        NOT NULL DEFAULT '',
    body        TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- +goose Down
DROP TABLE posts;
DROP INDEX users_username_lower_idx;
DROP TABLE users;
