package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/malinatrash/tabhub/internal/config"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Database) (*Storage, error) {
	const op = "storage.postgres.New"
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateProject(ctx context.Context, name string, ownerID int, state []byte) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO projects (name, owner_id, state, created_at, updated_at) 
			    VALUES ($1, $2, $3, NOW(), NOW())`, name, ownerID, state,
	)
	return err
}

func (s *Storage) GetProjectId(ctx context.Context, name string, ownerId int) (*int, error) {
	var id int
	err := s.db.QueryRowContext(ctx, `
        SELECT id FROM projects WHERE name = $1 and owner_id = $2
    `, name, ownerId).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("project with name %s not found", name)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &id, nil
}
