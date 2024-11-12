-- +goose Up
CREATE TABLE projects (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	owner_id INT REFERENCES users(id),
	state TEXT NOT NULL,
	private BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	CONSTRAINT unique_owner_project UNIQUE (owner_id, name)
);
-- +goose Down
DROP TABLE projects;