-- name: NewSession :one
INSERT INTO sessions
(user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSessionByTokenHash :one
SELECT s.*, u.username as session_username
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE token_hash = $1
    AND revoked_at IS NULL
    AND expires_at > NOW();

-- name: RevokeSession :exec
UPDATE sessions
SET revoked_at = NOW()
WHERE id = $1;
