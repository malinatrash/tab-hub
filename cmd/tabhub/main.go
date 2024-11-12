package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/malinatrash/tabhub/internal/config"
	projectsCreate "github.com/malinatrash/tabhub/internal/http-server/handlers/projects/create"
	projectsGet "github.com/malinatrash/tabhub/internal/http-server/handlers/projects/get"
	usersCreate "github.com/malinatrash/tabhub/internal/http-server/handlers/users/create"
	usersGet "github.com/malinatrash/tabhub/internal/http-server/handlers/users/get"
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

	router.Route("/projects", func(r chi.Router) {
		r.Post("/", projectsCreate.Handler(log, storage))
		r.Get("/", projectsGet.Handler(log, storage))
	})
	router.Route("/users", func(r chi.Router) {
		r.Post("/create", usersCreate.Handler(log, storage))
		r.Get("/get", usersGet.Handler(log, storage))
	})

	log.Info("Server starting!")
	log.Info("ENV:", cfg)
	if err = http.ListenAndServe(cfg.Server.Address, router); err != nil {
		log.Error(err.Error())
	}
}
