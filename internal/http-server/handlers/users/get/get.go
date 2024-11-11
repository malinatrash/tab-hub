package get

import (
	"context"
	"encoding/json"
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
	User(ctx context.Context, username string, passwordHash string) (*int, error)
}

func Handler(log *slog.Logger, manager UserManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user request
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		hashPassword, err := hash.Password(user.Password)
		if err != nil {
			http.Error(w, "Invalid password", http.StatusBadRequest)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Info(hashPassword)

		id, err := manager.User(ctx, user.Username, hashPassword)
		if err != nil {
			http.Error(w, "error while getting user", http.StatusInternalServerError)
			return
		}

		resp := response{
			ID: id,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
