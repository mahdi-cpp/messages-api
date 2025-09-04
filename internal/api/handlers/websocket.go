package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mahdi-cpp/messages-api/internal/application"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	appManager *application.Manager
}

// ServeWs is the legacy function for backward compatibility
func ServeWs(appManager *application.Manager, w http.ResponseWriter, r *http.Request) {
	handler1 := NewWebSocketHandler(appManager)
	handler1.ServeHTTP(w, r)
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(appManager *application.Manager) *WebSocketHandler {
	return &WebSocketHandler{
		appManager: appManager,
	}
}

// ServeHTTP handles HTTP requests and upgrades them to WebSocket
func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	fmt.Println("user_id:", userID)

	username := r.URL.Query().Get("username")
	if username == "" {
		username = userID
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	log.Printf("WebSocket connection established for user: %s (%s)", username, userID)

	h.appManager.CreateClient(conn, userID, username)
}
