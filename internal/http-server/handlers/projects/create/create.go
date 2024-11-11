package create

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/malinatrash/tabhub/pkg/xml"
	"log/slog"
	"net/http"
	"time"
)

type request struct {
	Name    string `json:"name" validate:"required"`
	OwnerID int    `json:"owner_id" validate:"required"`
	State   []byte `json:"state,omitempty"`
	Private bool   `json:"private"`
}

type response struct {
	Message string `json:"message"`
	ID      *int   `json:"project_id,omitempty"`
}

type ProjectManager interface {
	CreateProject(ctx context.Context, name string, ownerID int, state []byte, private bool) (*int, error)
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

		id, err := manager.CreateProject(ctx, project.Name, project.OwnerID, project.State, project.Private)
		if err != nil {
			errorDescription := fmt.Sprintf("Failed to create project: %s", err.Error())
			http.Error(w, errorDescription, http.StatusInternalServerError)
			log.Error(errorDescription)
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
