package create

import (
	"context"
	"encoding/json"
	"github.com/malinatrash/tabhub/pkg/xml"
	"log/slog"
	"net/http"
	"time"
)

type request struct {
	Name    string `json:"name" validate:"required"`
	OwnerID int    `json:"owner_id" validate:"required"`
	State   []byte `json:"state,omitempty"`
}

type response struct {
	Message string `json:"message"`
	ID      *int   `json:"project_id,omitempty"`
}

type ProjectManager interface {
	CreateProject(ctx context.Context, name string, ownerID int, state []byte) error
	GetProjectId(ctx context.Context, name string, ownerId int) (*int, error)
}

func Handler(log *slog.Logger, manager ProjectManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var project request
		if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		project.State = xml.GenerateEmptyMusicXML()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := manager.CreateProject(ctx, project.Name, project.OwnerID, project.State); err != nil {
			http.Error(w, "Failed to create project", http.StatusInternalServerError)
			log.Error("Error creating project: %v", err)
			return
		}

		id, err := manager.GetProjectId(ctx, project.Name, project.OwnerID)
		if err != nil {
			log.Error("Project does not exist: %v", err)
			http.Error(w, "Project not found", http.StatusInternalServerError)
			return
		}

		resp := response{
			Message: "Project created successfully",
			ID:      id,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
