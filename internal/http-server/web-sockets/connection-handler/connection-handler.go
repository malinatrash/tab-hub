package connection_handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
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
}

var projects = make(map[string]*projectWS)

func Handler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := chi.URLParam(r, "id")
		project, exists := projects[projectID]
		if !exists {
			project = &projectWS{
				ID:        projectID,
				Clients:   make(map[*websocket.Conn]bool),
				Broadcast: make(chan string),
			}
			projects[projectID] = project
			go handleMessages(log, project)
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Connection error:", err)
			return
		}
		defer ws.Close()

		project.Clients[ws] = true
		for {
			var msg string
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Error("Read error:", err)
				break
			}
			project.Broadcast <- msg
		}
	}
}

func handleMessages(log *slog.Logger, project *projectWS) {
	for {
		if len(project.Clients) == 0 {
			delete(projects, project.ID)
			log.Info("Проект %s был удалён, так как все клиенты отключились", project.ID)
			return
		}

		msg := <-project.Broadcast

		for client := range project.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Error("error: %v", err)
				err := client.Close()
				if err != nil {
					log.Error("error while client close:", err)
					return
				}
				delete(project.Clients, client)
			}
		}
	}
}
