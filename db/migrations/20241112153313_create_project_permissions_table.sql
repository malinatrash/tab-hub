-- +goose Up
CREATE TABLE project_permissions (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id) ON DELETE CASCADE,
	project_id INT REFERENCES projects(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW(),
	UNIQUE(user_id, project_id)
);

-- +goose Down
DROP TABLE project_permissions;