package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/application"
)

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

	idString := r.URL.Query().Get("user_id")
	if idString == "" {
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(idString)
	if err != nil {
		return
	}

	fmt.Println("user_id:", idString)

	username := r.URL.Query().Get("username")
	if username == "" {
		username = idString
	}

	h.appManager.CreateWebsocketClient(w, r, userID, username)
}
