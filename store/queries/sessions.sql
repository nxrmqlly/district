-- name: NewSession :one
INSERT INTO sessions
(user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSessionByTokenHash :one
SELECT * from sessions
WHERE token_hash = $1
    AND revoked_at IS NULL
    AND expires_at > NOW();