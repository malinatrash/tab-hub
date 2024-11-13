package delete

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type request struct {
	OwnerID   int `json:"owner_id"`
	UserID    int `json:"user_id"`
	ProjectID int `json:"project_id"`
}

type response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type PermissionsManager interface {
	DeletePermission(ctx context.Context, projectID int, userID int, ownerID int) error
}

func Handler(log *slog.Logger, manager PermissionsManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "delete.Handler"
		var permission request
		if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
			log.Error("error in op: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := manager.DeletePermission(ctx, permission.ProjectID, permission.UserID, permission.OwnerID)
		if err != nil {
			log.Error("error in op: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := response{
			Message: "Permission deleted successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("error in op: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
