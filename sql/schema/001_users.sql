-- +goose Up
CREATE TABLE users (
	id UUID primary key,
	created_at timestamp,
	updated_at timestamp,
	username string unique not null,
	password string unique not null, 
	email text not null
);

-- +goose Down
drop table users;
