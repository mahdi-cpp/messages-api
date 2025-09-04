package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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

	// Set the message handler for this client
	//client.SetMessageHandler(h.handleClientMessage)

	// Set user info in client (you'll need to add this field to Client struct)
	client.SetUserInfo(username)

	go client.WritePump()
	go client.ReadPump()

	// Send welcome message only to this client
	welcomeMessage := map[string]interface{}{
		"type":    "system",
		"message": "Welcome to the chat!",
		"userId":  userID,
		"success": true,
	}

	if err := client.SendMessage(welcomeMessage); err != nil {
		log.Printf("Failed to send welcome message to user %s: %v", userID, err)
	}

	// Notify all users in default room about new user
	h.notifyUserJoined(userID, username, "general")
}

// notifyUserJoined sends a notification when a user joins a room
func (h *WebSocketHandler) notifyUserJoined(userID, username, roomID string) {

	joinMessage := map[string]interface{}{
		"type":      "user_joined",
		"userId":    userID,
		"username":  username,
		"message":   username + " joined the room",
		"roomId":    roomID,
		"timestamp": time.Now(),
	}

	// Broadcast to the specific room
	h.hub.BroadcastToRoom(roomID, joinMessage)
}

//// Handle incoming messages from clients
//func (h *WebSocketHandler) handleClientMessage(client *hub.Client, rawMessage []byte) {
//
//	var message struct {
//		Type    string `json:"type"`
//		Content string `json:"content"`
//		ChatID  string `json:"chatId"`
//	}
//
//	if err := json.Unmarshal(rawMessage, &message); err != nil {
//		log.Printf("Error parsing message from client %s: %v", client.UserID(), err)
//		return
//	}
//
//	switch message.Type {
//	case "message":
//		h.handleChatMessage(client, message.Content, message.ChatID)
//	case "typing":
//		h.handleTypingIndicator(client, message.Content, message.ChatID)
//	case "join_room":
//		h.handleJoinRoom(client, message.ChatID)
//	case "leave_room":
//		h.handleLeaveRoom(client, message.ChatID)
//	case "create_room":
//		h.handleCreateRoom(client, message.Content)
//	}
//}
//
//// handleChatMessage processes and broadcasts chat messages
//func (h *WebSocketHandler) handleChatMessage(client *hub.Client, content, chatID string) {
//	if content == "" {
//		return
//	}
//
//	// Create the message to broadcast
//	chatMessage := map[string]interface{}{
//		"type":      "message",
//		"id":        generateUUID(),
//		"userId":    client.UserID(),
//		"username":  client.Username(),
//		"content":   content,
//		"chatID":    chatID,
//		"timestamp": time.Now(),
//	}
//
//	log.Printf("Broadcasting message from %s in room: %s: %s",
//		client.Username(), chatID, content)
//
//	// Broadcast to all clients in the room
//	h.hub.BroadcastToRoom(chatID, chatMessage)
//}
//
//// handleTypingIndicator broadcasts typing status
//func (h *WebSocketHandler) handleTypingIndicator(client *hub.Client, typing, roomID string) {
//
//	typingMessage := map[string]interface{}{
//		"type":      "typing",
//		"userId":    client.UserID(),
//		"username":  client.Username(),
//		"roomId":    roomID,
//		"typing":    typing == "true",
//		"timestamp": time.Now(),
//	}
//
//	h.hub.BroadcastToRoom(roomID, typingMessage)
//}
//
//// handleJoinRoom handles room joining
//func (h *WebSocketHandler) handleJoinRoom(client *hub.Client, roomID string) {
//	h.hub.JoinRoom(roomID, client.UserID(), client)
//
//	// Notify room about new user
//	h.notifyUserJoined(client.UserID(), client.Username(), roomID)
//}
//
//// handleLeaveRoom handles room leaving
//func (h *WebSocketHandler) handleLeaveRoom(client *hub.Client, roomID string) {
//
//	h.hub.LeaveRoom(roomID, client.UserID())
//
//	leaveMessage := map[string]interface{}{
//		"type":      "user_left",
//		"userId":    client.UserID(),
//		"username":  client.Username(),
//		"message":   client.Username() + " left the room",
//		"roomId":    roomID,
//		"timestamp": time.Now(),
//	}
//
//	h.hub.BroadcastToRoom(roomID, leaveMessage)
//}
//
//// handleCreateRoom handles room creation
//func (h *WebSocketHandler) handleCreateRoom(client *hub.Client, roomName string) {
//
//	roomID, err := generateUUID()
//	if err != nil {
//		fmt.Printf("Error generating room id: %v", err)
//		return
//	}
//
//	h.hub.CreateRoom(roomID, roomName)
//	h.hub.JoinRoom(roomID, client.UserID(), client)
//
//	// Notify about room creation
//	roomMessage := map[string]interface{}{
//		"type":      "room_created",
//		"roomId":    roomID,
//		"roomName":  roomName,
//		"userId":    client.UserID(),
//		"username":  client.Username(),
//		"timestamp": time.Now(),
//	}
//
//	h.hub.BroadcastToAll(roomMessage)
//}

// ServeWs is the legacy function for backward compatibility
func ServeWs(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	handler := NewWebSocketHandler(hub)
	handler.ServeHTTP(w, r)
}

func generateUUID() (string, error) {
	u7, err2 := uuid.NewV7()
	if err2 != nil {
		return "", fmt.Errorf("error generating UUIDv7: %w", err2)
	}
	return u7.String(), nil
}
