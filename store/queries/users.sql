-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE lower(username) = lower(sqlc.arg(username)::text);

-- name: NewUser :one
INSERT INTO users (username, email, passwd_hash)
VALUES ($1, $2, $3)
RETURNING *;
