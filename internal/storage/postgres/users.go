package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/malinatrash/tabhub/internal/storage/myErrors"
)

func (s *Storage) User(ctx context.Context, username string, passwordHash string) (*int, error) {
	var id int
	query := `SELECT id FROM users WHERE username = $1 and password_hash = $2`
	err := s.db.GetContext(ctx, &id, query, username, passwordHash)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *Storage) CreateUser(ctx context.Context, username string, passwordHash string) error {
	var existingID int
	query := `SELECT id FROM users WHERE username = $1`
	err := s.db.GetContext(ctx, &existingID, query, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", myErrors.ErrDBInsert, err)
	}
	if existingID != 0 {
		return fmt.Errorf("%w: %s", myErrors.ErrUserAlreadyExists, username)
	}

	_, err = s.db.ExecContext(ctx, `INSERT INTO users (username, password_hash) VALUES ($1, $2)`, username, passwordHash)
	if err != nil {
		return fmt.Errorf("%w: %v", myErrors.ErrDBInsert, err)
	}

	return nil
}
