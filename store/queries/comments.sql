-- name: GetCommentsByPost :many
SELECT c.*, u.username AS author_username
FROM comments c
JOIN users u ON u.id = c.author_id
WHERE c.post_id = $1
ORDER BY c.created_at DESC, c.id ASC;

-- name: CreateComment :one
INSERT INTO comments
(post_id, parent_id, author_id, body)
SELECT $1, $2, $3, $4     -- creates a custom row here
WHERE $2::BIGINT IS NULL  -- stop here if parent comment
OR EXISTS (
    SELECT 1              -- discard, returns true if succeeds
    FROM comments
    WHERE id = $2::BIGINT --  check if the parent comment exists on the post
        AND post_id = $1
)
RETURNING *;
