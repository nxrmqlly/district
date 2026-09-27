-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE lower(username) = lower(sqlc.arg(username)::text);

-- name: NewUser :one
INSERT INTO users
(id, username, email, passwd_hash, created_at)
VALUES (uuidv7(), $1, $2, $3, NOW())
RETURNING *;
