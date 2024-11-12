package postgres

import (
	"context"
	"fmt"
	"github.com/malinatrash/tabhub/internal/storage/models"
)

func (s *Storage) CreateProject(ctx context.Context, name string, ownerID int, state []byte, private bool) (*int, error) {
	var id int
	query := `INSERT INTO projects (name, owner_id, state, private, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	err := s.db.GetContext(ctx, &id, query, name, ownerID, state, private)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return &id, nil
}

func (s *Storage) Project(ctx context.Context, id int) (*models.Project, error) {
	var project models.Project
	query := `SELECT id, name, owner_id, state, private FROM projects WHERE id = $1`
	err := s.db.GetContext(ctx, &project, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project with id %d: %w", id, err)
	}

	return &project, nil
}
