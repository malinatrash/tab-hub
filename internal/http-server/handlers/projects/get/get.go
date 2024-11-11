package get

import (
	"context"
	"encoding/json"
	"github.com/malinatrash/tabhub/internal/storage/models"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type response struct {
	Name    string `json:"name" validate:"required"`
	OwnerID int    `json:"owner_id" validate:"required"`
	State   []byte `json:"state,omitempty"`
}

type ProjectManager interface {
	Project(ctx context.Context, id int) (*models.Project, error)
}

func Handler(log *slog.Logger, manager ProjectManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Missing project ID", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid project ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		project, err := manager.Project(ctx, id)
		if err != nil {
			log.Error("Project does not exist: %v", err)
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}

		resp := response{
			Name:    project.Name,
			OwnerID: project.OwnerID,
			State:   project.State,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
