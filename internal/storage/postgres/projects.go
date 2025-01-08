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

func (s *Storage) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	query := "SELECT id, name, owner_id, private, created_at, updated_at FROM projects WHERE 1=1"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.OwnerID, &p.Private, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.State = ""
		projects = append(projects, p)
	}

	return projects, nil
}
