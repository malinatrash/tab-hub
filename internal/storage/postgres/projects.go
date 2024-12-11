package postgres

import (
	"context"
	"fmt"
	"github.com/malinatrash/tabhub/internal/storage/models"
)

func (s *Storage) CreateProject(ctx context.Context, name string, ownerID int, state string, private bool) (*int, error) {
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

func (s *Storage) UpdateProjectState(ctx context.Context, project *models.Project) error {
	query := `UPDATE projects SET state = $1, updated_at = NOW() WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, project.State, project.ID)
	if err != nil {
		return fmt.Errorf("failed to update project with id %d: %w", project.ID, err)
	}

	return nil
}
