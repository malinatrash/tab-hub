package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/malinatrash/tabhub/internal/config"
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

	rows, err := db.Query(`
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public';`)
	if err != nil {
		log.Fatalf("Ошибка получения списка таблиц: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Ошибка сканирования имени таблицы: %v", err)
		}

		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", tableName))
		if err != nil {
			log.Printf("Ошибка удаления таблицы %s: %v", tableName, err)
		} else {
			log.Printf("Таблица %s успешно удалена", tableName)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Ошибка перебора таблиц: %v", err)
	}
}
