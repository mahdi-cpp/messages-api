package handlers

import (
	"fmt"
	"net/http"

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

	h.appManager.CreateWebsocketClient(w, r, userID, username)
}
