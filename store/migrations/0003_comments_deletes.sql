-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
    id         BIGSERIAL   PRIMARY KEY,
    post_id    BIGINT      NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id  BIGINT      REFERENCES comments(id) ON DELETE SET NULL,
    author_id  UUID        NOT NULL REFERENCES users(id),
    body       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX comments_post_id_idx
    ON comments(post_id);

CREATE INDEX comments_parent_id_idx
    ON comments(parent_id);

-- +goose Down
DROP INDEX comments_parent_id_idx;
DROP INDEX comments_post_id_idx;
DROP TABLE comments;
