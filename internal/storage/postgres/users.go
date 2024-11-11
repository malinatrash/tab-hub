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
	err := s.db.QueryRowContext(ctx, `SELECT id FROM users WHERE username = $1 and password_hash = $2`, username, passwordHash).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *Storage) CreateUser(ctx context.Context, username string, passwordHash string) error {
	var existingID int
	err := s.db.QueryRowContext(ctx, `SELECT id FROM users WHERE username = $1`, username).Scan(&existingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", myErrors.ErrDBInsert, err)
	}
	if err == nil {
		return fmt.Errorf("%w: %s", myErrors.ErrUserAlreadyExists, username)
	}

	_, err = s.db.ExecContext(ctx, `INSERT INTO users (username, password_hash) VALUES ($1, $2)`, username, passwordHash)
	if err != nil {
		return fmt.Errorf("%w: %v", myErrors.ErrDBInsert, err)
	}

	return nil
}
