-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;

-- name: NewPost :one
INSERT INTO posts
(author_id, title, embed_url, body, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPostsByUsername :many
SELECT p.*, u.username as author_username
FROM posts p 
JOIN users u ON u.id = p.author_id
WHERE lower(u.username) = lower(sqlc.arg(username)::text)
ORDER BY p.created_at DESC;