package users

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/tabhub")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close(ctx)

	_, err = dbpool.Exec(ctx, `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)`, user.Username, user.PasswordHash)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
