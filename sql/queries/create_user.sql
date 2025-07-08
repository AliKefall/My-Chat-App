-- name: CreateUser :exec
insert into users(id, username, email, created_at, updated_at)
values(
	$1,
	$2,
	$3,
	$4,
	$5
);

-- name: GetUser :one
select * from users where id = $1;
