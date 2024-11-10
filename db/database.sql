CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id INT REFERENCES users(id),
	state BYTEA NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_owner_project UNIQUE (owner_id, name)
);

CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    password_hash VARCHAR(255)
);

CREATE TABLE project_permissions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, project_id)
);


