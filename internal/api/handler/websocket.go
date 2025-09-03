package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mahdi-cpp/messages-api/internal/hub"
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
	hub *hub.Hub
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *hub.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// ServeHTTP handles HTTP requests and upgrades them to WebSocket
func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		fmt.Println("User ID required")
		http.Error(w, "User ID required", http.StatusUnauthorized)
		return
	}

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

	client := hub.NewClient(h.hub, conn, userID)
	h.hub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump()

	welcomeMessage := map[string]interface{}{
		"type":    "system",
		"message": "Welcome to the chat!",
		"userId":  userID,
		"success": true,
	}

	if err := client.SendMessage(welcomeMessage); err != nil {
		log.Printf("Failed to send welcome message to user %s: %v", userID, err)
	}

	h.notifyUserJoined(userID, username)
}

// notifyUserJoined sends a notification when a user joins
func (h *WebSocketHandler) notifyUserJoined(userID, username string) {
	joinMessage := map[string]interface{}{
		"type":     "user_joined",
		"userId":   userID,
		"username": username,
		"message":  username + " joined the chat",
	}

	h.hub.BroadcastToRoom("general", joinMessage)
}

// ServeWs is the legacy function for backward compatibility
func ServeWs(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	handler := NewWebSocketHandler(hub)
	handler.ServeHTTP(w, r)
}

// HealthCheckHandler handles WebSocket connection health checks
func HealthCheckHandler(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"clients":   hub.GetClientCount(),
		"rooms":     len(hub.GetRoomList()),
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetConnectionsHandler returns current connection statistics
func GetConnectionsHandler(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	stats := hub.GetHubStats() // Use the new method

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
