package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/malinatrash/tabhub/internal/config"
)

type Storage struct {
	db *sqlx.DB
}

func New(cfg config.Database) (*Storage, error) {
	const op = "storage.postgres.New"
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
