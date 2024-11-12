package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/malinatrash/tabhub/internal/config"
	"github.com/pressly/goose/v3"
	"log"
)

func main() {
	cfg := config.MustLoad()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	if err := goose.Down(db, "db/migrations"); err != nil {
		log.Fatalf("Ошибка отката миграции: %v", err)
	}
}
