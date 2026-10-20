-- +goose Up

ALTER TABLE posts
    ALTER COLUMN title DROP NOT NULL,
    ALTER COLUMN body DROP NOT NULL,
    ALTER COLUMN embed_url DROP NOT NULL;

ALTER TABLE posts
    ADD CONSTRAINT posts_content_integrity CHECK (
        (title IS NOT NULL AND body IS NOT NULL AND embed_url IS NOT NULL)
        OR
        (deleted_at IS NOT NULL AND title IS NULL AND body IS NULL AND embed_url IS NULL)
    );

---

ALTER TABLE comments
    ALTER COLUMN body DROP NOT NULL;

ALTER TABLE comments
    ADD CONSTRAINT posts_content_integrity CHECK (
        (body IS NOT NULL)
        OR
        (deleted_at IS NOT NULL AND body IS NULL)
    );

-- ensure acyclic comment graph
ALTER TABLE comments
    ADD CONSTRAINT comments_parents_acyclic
    CHECK (parent_id IS NULL OR parent_id <> id);

-- +goose Down

ALTER TABLE comments
    DROP CONSTRAINT comments_parents_acyclic;

ALTER TABLE comments
    DROP CONSTRAINT comments_content_integrity;

ALTER TABLE comments
    ALTER COLUMN body SET NOT NULL;

---

ALTER TABLE posts
    DROP CONSTRAINT posts_content_integrity;

ALTER TABLE posts
    ALTER COLUMN title SET NOT NULL,
    ALTER COLUMN body SET NOT NULL,
    ALTER COLUMN embed_url SET NOT NULL;

