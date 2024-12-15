package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/malinatrash/tabhub/internal/config"
	permissionsCreate "github.com/malinatrash/tabhub/internal/http_server/handlers/permissions/create"
	permissionsDelete "github.com/malinatrash/tabhub/internal/http_server/handlers/permissions/delete"
	projectsCreate "github.com/malinatrash/tabhub/internal/http_server/handlers/projects/create"
	projectsGet "github.com/malinatrash/tabhub/internal/http_server/handlers/projects/get"
	usersCreate "github.com/malinatrash/tabhub/internal/http_server/handlers/users/create"
	usersGet "github.com/malinatrash/tabhub/internal/http_server/handlers/users/get"
	wsManager "github.com/malinatrash/tabhub/internal/http_server/web_sockets/project"
	"github.com/malinatrash/tabhub/internal/lib/logger"
	"github.com/malinatrash/tabhub/internal/storage/postgres"
	"github.com/malinatrash/tabhub/internal/storage/redis"

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

	redisClient, err := redis.New(cfg.Cache)
	if err != nil {
		log.Error("failed to initialize redisClient", err.Error())
		os.Exit(1)
	}
	_ = redisClient

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/projects", func(r chi.Router) {

		r.Post("/", projectsCreate.Handler(log, storage))
		r.Get("/{id}", projectsGet.Handler(log, storage))

		r.Get("/{id}/ws", wsManager.Handler(log, redisClient, storage))

		r.Route("/permissions", func(r chi.Router) {
			r.Post("/", permissionsCreate.Handler(log, storage))
			r.Delete("/", permissionsDelete.Handler(log, storage))
		})
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
