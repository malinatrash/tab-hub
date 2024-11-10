package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/malinatrash/tabhub/internal/config"
	"github.com/malinatrash/tabhub/internal/http-server/handlers/projects/create"
	"github.com/malinatrash/tabhub/internal/lib/logger"
	"github.com/malinatrash/tabhub/internal/storage/postgres"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	storage, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to initialize storage", err.Error())
		os.Exit(1)
	}
	_ = storage

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/project", create.Handler(log, storage))

	if err = http.ListenAndServe(cfg.Server.Address, nil); err != nil {
		log.Error(err.Error())
	}
}
