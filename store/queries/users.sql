-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE lower(username) = lower(sqlc.arg(username)::text);

-- name: NewUser :one
INSERT INTO users (username, email, passwd_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByLogin :one
SELECT * FROM users
WHERE lower(email) = lower(sqlc.arg(login)::text)
    OR lower(username) = lower(sqlc.arg(login)::text);
