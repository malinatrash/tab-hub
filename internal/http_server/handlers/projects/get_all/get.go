package get_all

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/malinatrash/tabhub/internal/storage/models"
)

type ProjectManager interface {
	GetAllProjects(ctx context.Context) ([]models.Project, error)
}

func Handler(log *slog.Logger, manager ProjectManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		projects, err := manager.GetAllProjects(ctx)
		if err != nil {
			log.Error("Failed to retrieve projects", "error", err)
			http.Error(w, "Failed to retrieve projects", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(projects); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
