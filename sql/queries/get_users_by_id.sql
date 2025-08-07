-- name: GetUsersByID :one
select * from users
where id = $1;
