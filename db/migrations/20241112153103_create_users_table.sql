-- +goose Up
CREATE TABLE Users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(255) UNIQUE,
	password_hash VARCHAR(255)
);
-- +goose Down
DROP TABLE Users;
