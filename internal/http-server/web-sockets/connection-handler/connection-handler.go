package connection_handler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/malinatrash/tabhub/internal/storage/models"
	"github.com/malinatrash/tabhub/internal/storage/redis"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type projectWS struct {
	ID        string
	Clients   map[*websocket.Conn]bool
	Broadcast chan string
	mu        sync.Mutex
}

type CacheManager interface {
	PushProject(ctx context.Context, projectID int, state string) error
	DeleteProject(ctx context.Context, projectID int) error
}

type ProjectManager interface {
	Project(ctx context.Context, projectID int) (*models.Project, error)
	UpdateProjectState(ctx context.Context, project *models.Project) error
}

var projects = make(map[string]*projectWS)

func Handler(log *slog.Logger, cManager *redis.Client, pManager ProjectManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := chi.URLParam(r, "id")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Hour)
		defer cancel()

		iProjectID, err := strconv.Atoi(projectID)
		if err != nil {
			log.Error("Invalid project id")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		project, exists := projects[projectID]
		if !exists {
			project = &projectWS{
				ID:        projectID,
				Clients:   make(map[*websocket.Conn]bool),
				Broadcast: make(chan string),
			}

			projectInStorage, err := pManager.Project(ctx, iProjectID)
			if err != nil {
				log.Error("Failed to get project")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			err = cManager.PushProject(ctx, iProjectID, projectInStorage.State)
			if err != nil {
				log.Error("Failed to push project in cache", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			projects[projectID] = project
			go handleMessages(log, project, cManager, pManager)
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Connection error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer func(ws *websocket.Conn) {
			err := ws.Close()
			if err != nil {
				log.Error("Connection error", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}(ws)

		project.mu.Lock()
		project.Clients[ws] = true
		project.mu.Unlock()

		for {
			var data map[string]interface{}
			err := ws.ReadJSON(&data)
			if err != nil {
				log.Error("Read error", err)
				break
			}

			log.Info("Received message", data)

			msg, ok := data["message"].(string)
			if !ok {
				log.Error("Invalid message format")
				continue
			}

			err = cManager.PushProject(ctx, iProjectID, msg)
			if err != nil {
				log.Error("Failed to push project in cache", err)
			}
			project.Broadcast <- msg
		}
	}
}

func handleMessages(log *slog.Logger, project *projectWS, cManager *redis.Client, pManager ProjectManager) {
	for {
		msg := <-project.Broadcast

		project.mu.Lock()
		if len(project.Clients) == 0 {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			id, err := strconv.Atoi(project.ID)
			if err != nil {
				log.Error("Invalid project id")
			}
			err = cManager.DeleteProject(ctx, id)
			if err != nil {
				log.Error("Failed to delete project in cache", err)
			}
			delete(projects, project.ID)
			log.Info("Проект %s был удалён, так как все клиенты отключились", project.ID)

			go func() {
				projectInStorage, err := pManager.Project(ctx, id)
				if err != nil {
					log.Error("Failed to fetch project for update")
				} else {
					projectInStorage.State = msg
					err = pManager.UpdateProjectState(ctx, projectInStorage)
					if err != nil {
						log.Error("Failed to update project state", err)
					}
				}
			}()
			project.mu.Unlock()
			return
		}
		project.mu.Unlock()

		for client := range project.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Error("error", err)
				err := client.Close()
				if err != nil {
					log.Error("error while client close", err)
				}
				project.mu.Lock()
				delete(project.Clients, client)
				project.mu.Unlock()
			}
		}
	}
}
