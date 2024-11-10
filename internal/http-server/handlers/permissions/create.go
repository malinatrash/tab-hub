package permissions

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
)

type projectPermissionsRequest struct {
	OwnerID   int `json:"owner_id"`
	UserID    int `json:"user_id"`
	ProjectID int `json:"project_id"`
}

func SetProjectPermissionsHandler(writer http.ResponseWriter, request *http.Request) {
	var permission projectPermissionsRequest
	if err := json.NewDecoder(request.Body).Decode(&permission); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/tabhub")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close(ctx)

	var projectOwnerID int
	err = dbpool.QueryRow(ctx, `
		SELECT owner_id FROM projects WHERE id = $1`, permission.ProjectID).Scan(&projectOwnerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(writer, "Project not found", http.StatusNotFound)
		} else {
			http.Error(writer, "Failed to check project owner", http.StatusInternalServerError)
		}
		log.Printf("Error checking project owner: %v", err)
		return
	}

	if projectOwnerID != permission.OwnerID {
		http.Error(writer, "User is not the project owner", http.StatusForbidden)
		return
	}

	_, err = dbpool.Exec(ctx, `
		INSERT INTO project_permissions (user_id, project_id, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())`,
		permission.UserID, permission.ProjectID,
	)

	if err != nil {
		http.Error(writer, "Failed to set permissions", http.StatusInternalServerError)
		log.Printf("Error setting permissions: %v", err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Permissions set successfully"))
}
