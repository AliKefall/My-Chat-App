-- name: CreateUser :one
insert into users(username, password_hash)
values ($1, $2)
returning *;

-- name: GetUserByUsername :one
select * from users where username = $1;

-- name: CreateMessage :one
insert into messages (user_id, content, iv, salt)
values ($1, $2, $3, $4)
returning *;

-- name: ListMessages :many
select * from messages
order by created_at ASC
LIMIT $1;


