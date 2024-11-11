package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/malinatrash/tabhub/internal/storage/models"
)

func (s *Storage) CreateProject(ctx context.Context, name string, ownerID int, state []byte, private bool) (*int, error) {
	var id int
	err := s.db.QueryRowContext(ctx, `INSERT INTO projects (name, owner_id, state, private,  created_at, updated_at)  VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`, name, ownerID, state, private).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *Storage) Project(ctx context.Context, id int) (*models.Project, error) {
	var project models.Project
	err := s.db.QueryRowContext(ctx, `SELECT id FROM projects WHERE id = $1`, id).Scan(&project.Name, &project.OwnerID, &project.State)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("project with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}
