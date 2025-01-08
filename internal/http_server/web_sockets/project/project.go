package project

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/malinatrash/tabhub/internal/storage/models"
	"github.com/malinatrash/tabhub/internal/storage/redis"
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

type cacheManager interface {
	PushProject(ctx context.Context, projectID int, state string) error
	DeleteProject(ctx context.Context, projectID int) error
}

type ProjectManager interface {
	Project(ctx context.Context, projectID int) (*models.Project, error)
	UpdateProjectState(ctx context.Context, project *models.Project) error
}

var projects sync.Map // Используем sync.Map для потокобезопасного хранения проектов

func Handler(log *slog.Logger, cManager *redis.Client, pManager ProjectManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := chi.URLParam(r, "id")

		iProjectID, err := strconv.Atoi(projectID)
		if err != nil {
			log.Error("Invalid project id")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		project, exists := projects.Load(projectID)
		if !exists {
			project = &projectWS{
				ID:        projectID,
				Clients:   make(map[*websocket.Conn]bool),
				Broadcast: make(chan string),
			}

			projectInStorage, err := pManager.Project(context.Background(), iProjectID)
			if err != nil {
				log.Error("Failed to get project")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			err = cManager.PushProject(context.Background(), iProjectID, projectInStorage.State)
			if err != nil {
				log.Error("Failed to push project in cache", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			projects.Store(projectID, project)
			go handleMessages(log, project.(*projectWS), cManager, pManager)
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Connection error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		project.(*projectWS).mu.Lock()
		project.(*projectWS).Clients[ws] = true
		project.(*projectWS).mu.Unlock()

		(project.(*projectWS)).Handler(w, r, ws, log, cManager, pManager, iProjectID)
	}
}

func (p *projectWS) logClients(log *slog.Logger, action string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	clientAddrs := make([]string, 0, len(p.Clients))
	for client := range p.Clients {
		remoteAddr := client.RemoteAddr().String()
		clientAddrs = append(clientAddrs, remoteAddr)
	}

	log.Info("WebSocket Clients Update",
		"action", action,
		"total_clients", len(p.Clients),
		"client_addresses", clientAddrs,
	)
}

func (p *projectWS) broadcastMessage(log *slog.Logger, msg string, pManager ProjectManager) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Info("Broadcasting project message",
		"message", msg,
		"current_clients", len(p.Clients),
		"project_id", p.ID,
	)

	_, err := strconv.Atoi(p.ID)
	if err != nil {
		log.Error("Invalid project ID",
			"error", err,
			"project_id", p.ID,
		)
		return err
	}

	payload := map[string]string{
		"state": msg,
	}

	for client := range p.Clients {
		err := client.WriteJSON(payload)
		if err != nil {
			log.Error("Failed to send message to client",
				"error", err,
				"client_address", client.RemoteAddr().String(),
				"project_id", p.ID,
			)

			delete(p.Clients, client)
			p.logClients(log, "problematic_client_removed")

			err := client.Close()
			if err != nil {
				log.Error("Error closing problematic client",
					"error", err,
					"project_id", p.ID,
				)
			}
		}
	}

	return nil
}

func (p *projectWS) saveProjectState(log *slog.Logger, cManager *redis.Client, pManager ProjectManager, lastMsg string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Clients) == 0 {
		return nil
	}

	log.Info("Evaluating project state save",
		"project_id", p.ID,
		"current_clients", len(p.Clients),
		"last_message_length", len(lastMsg),
	)

	// Always attempt to save, even if there are clients
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := strconv.Atoi(p.ID)
	if err != nil {
		log.Error("Invalid project id for state saving",
			"error", err,
			"project_id", p.ID,
		)
		return err
	}

	projectToUpdate := &models.Project{
		ID:    id,
		State: lastMsg,
	}

	err = pManager.UpdateProjectState(ctx, projectToUpdate)
	if err != nil {
		log.Error("Failed to update project state",
			"error", err,
			"project_id", id,
			"last_message", lastMsg,
		)
		return err
	}

	log.Info("Project state successfully updated",
		"project_id", id,
		"state_length", len(lastMsg),
	)

	err = cManager.DeleteProject(ctx, id)
	if err != nil {
		log.Error("Failed to delete project in cache",
			"error", err,
			"project_id", id,
		)
		return err
	}

	// Only delete from projects map if no clients
	if len(p.Clients) == 0 {
		projects.Delete(p.ID)
		log.Info("Project removed due to no active clients",
			"project_id", p.ID,
		)
	}

	return nil
}

func handleMessages(log *slog.Logger, project *projectWS, cManager *redis.Client, pManager ProjectManager) {
	lastMsg := ""
	for {
		select {
		case msg := <-project.Broadcast:
			lastMsg = msg

			log.Info("Processing incoming message",
				"project_id", project.ID,
				"message_length", len(msg),
			)

			// Broadcast with structured message
			if err := project.broadcastMessage(log, msg, pManager); err != nil {
				log.Error("Failed to broadcast message",
					"error", err,
					"project_id", project.ID,
				)
			}

			// Log and ignore any save errors
			if err := project.saveProjectState(log, cManager, pManager, lastMsg); err != nil {
				log.Error("Project state save failed",
					"error", err,
					"project_id", project.ID,
				)
			}
		}
	}
}

func (p *projectWS) Handler(w http.ResponseWriter, r *http.Request, ws *websocket.Conn, log *slog.Logger, cManager *redis.Client, pManager ProjectManager, iProjectID int) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	log.Info("WebSocket Handler started",
		"project_id", p.ID,
		"client_address", ws.RemoteAddr().String(),
	)

	p.logClients(log, "client_connected")

	done := make(chan struct{})

	go func() {
		defer close(done)
		defer func() {
			log.Info("WebSocket connection cleanup started",
				"project_id", p.ID,
				"client_address", ws.RemoteAddr().String(),
			)

			p.mu.Lock()
			delete(p.Clients, ws)
			p.mu.Unlock()

			p.logClients(log, "client_disconnected")

			if err := ws.Close(); err != nil {
				log.Error("Error closing WebSocket",
					"error", err,
					"project_id", p.ID,
				)
			}

			log.Info("WebSocket connection cleanup completed",
				"project_id", p.ID,
				"client_address", ws.RemoteAddr().String(),
			)
		}()

		for {
			select {
			case <-ctx.Done():
				log.Info("WebSocket read loop context cancelled", "project_id", p.ID)
				return
			default:
				_, msg, err := ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Error("WebSocket read error",
							"error", err,
							"project_id", p.ID,
							"error_type", fmt.Sprintf("%T", err),
						)
					}
					return
				}

				log.Info("Received WebSocket message",
					"msg_length", len(msg),
					"project_id", p.ID,
					"client_address", ws.RemoteAddr().String(),
				)

				err = cManager.PushProject(ctx, iProjectID, string(msg))
				if err != nil {
					log.Error("Failed to push project in cache",
						"error", err,
						"project_id", p.ID,
					)
				}
				p.Broadcast <- string(msg)
			}
		}
	}()

	<-done
}
