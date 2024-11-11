package create

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/malinatrash/tabhub/internal/storage/myErrors"
	"github.com/malinatrash/tabhub/pkg/hash"
	"log/slog"
	"net/http"
	"time"
)

type request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type response struct {
	ID *int `json:"user_id,omitempty"`
}

type UserManager interface {
	CreateUser(ctx context.Context, username string, password string) error
	User(ctx context.Context, username string, password string) (*int, error)
}

func Handler(log *slog.Logger, manager UserManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user request
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			log.Error("Invalid user data: username or password missing")
			return
		}

		hashPassword, err := hash.Password(user.Password)
		if err != nil {
			http.Error(w, "Invalid password", http.StatusBadRequest)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = manager.CreateUser(ctx, user.Username, hashPassword)
		if err != nil {
			switch {
			case errors.Is(err, myErrors.ErrDBInsert):
				http.Error(w, "User creation failed", http.StatusInternalServerError)
				log.Error("Failed to create user", err)
			case errors.Is(err, myErrors.ErrUserAlreadyExists):
				http.Error(w, "User already exists", http.StatusBadRequest)
				log.Error("Failed to create user", err)
			default:
				http.Error(w, "Unknown error occurred", http.StatusInternalServerError)
				log.Error("Unexpected error", err)
			}
			return
		}

		id, err := manager.User(ctx, user.Username, hashPassword)
		if err != nil {
			log.Error("error while creating user: %v", err)
			http.Error(w, "error while creating user", http.StatusInternalServerError)
			return
		}

		resp := response{
			ID: id,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
