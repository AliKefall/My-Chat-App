-- name: GetUsers :many
select * from users
order by created_at ASC;
