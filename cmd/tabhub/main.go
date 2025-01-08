package main

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/malinatrash/tabhub/internal/config"
	permissionsCreate "github.com/malinatrash/tabhub/internal/http_server/handlers/permissions/create"
	permissionsDelete "github.com/malinatrash/tabhub/internal/http_server/handlers/permissions/delete"
	projectsCreate "github.com/malinatrash/tabhub/internal/http_server/handlers/projects/create"
	projectsGet "github.com/malinatrash/tabhub/internal/http_server/handlers/projects/get"
	projectsGetAll "github.com/malinatrash/tabhub/internal/http_server/handlers/projects/get_all"
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
	defer storage.Close()

	redisClient, err := redis.New(cfg.Cache)
	if err != nil {
		log.Error("failed to initialize redisClient", err.Error())
		os.Exit(1)
	}
	defer redisClient.Close()

	router := chi.NewRouter()

	// Add CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/projects", func(r chi.Router) {

		r.Get("/", projectsGetAll.Handler(log, storage))
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

	addr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)

	if err = http.ListenAndServe(addr, router); err != nil {
		log.Error(err.Error())
	}
}
