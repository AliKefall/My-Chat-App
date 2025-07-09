
-- name: CreateUser :exec
INSERT INTO users(id, username,password ,email, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

