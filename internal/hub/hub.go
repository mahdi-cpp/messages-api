package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Chat represents a chat chat
type Chat struct {
	ID        string
	Name      string
	Clients   map[string]*Client // userID -> Client
	CreatedAt time.Time
}

// Hub manages all connected clients and chats
type Hub struct {
	chats     map[string]*Chat
	clients   map[string]*Client // userID -> Client
	mutex     sync.RWMutex
	startTime time.Time
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	hub := &Hub{
		chats:     make(map[string]*Chat),
		clients:   make(map[string]*Client),
		startTime: time.Now(),
	}

	// Create default chat
	hub.CreateChat("general", "General Chat")

	return hub
}

// Run starts the hub (maintain for compatibility)
func (h *Hub) Run() {
	log.Println("Hub started and running")
	// This method can be used for background tasks if needed
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		h.CleanupInactiveChats()
	}
}

// RegisterClient adds a client to the hub
func (h *Hub) RegisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client.userID] = client
	log.Printf("Client registered: %s. Total clients: %d", client.userID, len(h.clients))
}

// UnregisterClient removes a client from the hub and all chats
func (h *Hub) UnregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Remove client from all chats
	for chatID := range client.chats {
		if chat, exists := h.chats[chatID]; exists {
			delete(chat.Clients, client.userID)
			log.Printf("Client %s removed from chat %s", client.userID, chatID)
		}
	}

	// Remove client from main clients map
	delete(h.clients, client.userID)
	log.Printf("Client unregistered: %s. Remaining clients: %d", client.userID, len(h.clients))
}

// LeaveChat removes a client from a specific chat
func (h *Hub) LeaveChat(chatID, userID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if chat, exists := h.chats[chatID]; exists {
		if _, clientExists := chat.Clients[userID]; clientExists {
			delete(chat.Clients, userID)
			log.Printf("User %s left chat %s. Users remaining: %d", userID, chatID, len(chat.Clients))

			// Also remove from client's chat list
			if client, clientExists := h.clients[userID]; clientExists {
				client.LeaveChat(chatID)
			}
		} else {
			log.Printf("User %s not found in chat %s", userID, chatID)
		}
	} else {
		log.Printf("Chat %s not found", chatID)
	}
}

// JoinChat adds a client to a chat
func (h *Hub) JoinChat(chatID, userID string, client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Ensure chat exists
	if _, exists := h.chats[chatID]; !exists {
		log.Printf("Chat %s does not exist, creating it", chatID)
		h.chats[chatID] = &Chat{
			ID:        chatID,
			Name:      chatID, // Use ID as name for auto-created chats
			Clients:   make(map[string]*Client),
			CreatedAt: time.Now(),
		}
	}

	chat := h.chats[chatID]

	// Add client to the chat
	chat.Clients[userID] = client
	client.JoinChat(chatID)

	log.Printf("User %s joined chat %s (Total in chat: %d)", userID, chatID, len(chat.Clients))
}

// CreateChat creates a new chat chat
func (h *Hub) CreateChat(chatID, chatName string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.chats[chatID]; !exists {
		h.chats[chatID] = &Chat{
			ID:        chatID,
			Name:      chatName,
			Clients:   make(map[string]*Client),
			CreatedAt: time.Now(),
		}
		log.Printf("Chat created: %s (%s). Total chats: %d", chatName, chatID, len(h.chats))
	}
}

// BroadcastToChat sends a message to all clients in a chat
func (h *Hub) BroadcastToChat(chatID string, message interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Convert message to JSON bytes
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	if chat, exists := h.chats[chatID]; exists {
		for userID, client := range chat.Clients {
			// Check if client is still connected and channel is not full
			select {
			case client.send <- messageBytes:
				// Message sent successfully
				log.Printf("Message sent to user %s in chat %s", userID, chatID)
			default:
				// Channel is full, client might be disconnected
				log.Printf("Client %s send buffer full, potentially disconnected", userID)
			}
		}
	} else {
		log.Printf("Chat %s not found for broadcasting", chatID)
	}
}

// GetChatList returns a map of chat IDs to chat names
func (h *Hub) GetChatList() map[string]string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	chatList := make(map[string]string)
	for id, chat := range h.chats {
		chatList[id] = chat.Name
	}
	return chatList
}

// GetChatUsers returns the list of user IDs in a specific chat
func (h *Hub) GetChatUsers(chatID string) []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if chat, exists := h.chats[chatID]; exists {
		users := make([]string, 0, len(chat.Clients))
		for userID := range chat.Clients {
			users = append(users, userID)
		}
		return users
	}
	return []string{}
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// GetChatStats returns statistics for all chats
func (h *Hub) GetChatStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	stats := make(map[string]interface{})
	for chatID, chat := range h.chats {
		stats[chatID] = map[string]interface{}{
			"name":       chat.Name,
			"user_count": len(chat.Clients),
			"created_at": chat.CreatedAt,
		}
	}
	return stats
}

// GetStartTime returns the hub start time
func (h *Hub) GetStartTime() time.Time {
	return h.startTime
}

// GetHubStats returns comprehensive hub statistics
func (h *Hub) GetHubStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return map[string]interface{}{
		"total_clients": len(h.clients),
		"total_chats":   len(h.chats),
		"uptime":        time.Since(h.startTime).String(),
		"start_time":    h.startTime,
		"chat_stats":    h.GetChatStats(),
	}
}

// GetChat returns a chat by ID
func (h *Hub) GetChat(chatID string) (*Chat, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	chat, exists := h.chats[chatID]
	return chat, exists
}

// ChatExists checks if a chat exists
func (h *Hub) ChatExists(chatID string) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	_, exists := h.chats[chatID]
	return exists
}

// GetClient returns a client by user ID
func (h *Hub) GetClient(userID string) (*Client, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	client, exists := h.clients[userID]
	return client, exists
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(message interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for userID, client := range h.clients {
		select {
		case client.send <- messageBytes:
			// Message sent successfully
		default:
			log.Printf("Client %s send buffer full", userID)
		}
	}
}

// CleanupInactiveChats removes chats with no clients (optional)
func (h *Hub) CleanupInactiveChats() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for chatID, chat := range h.chats {
		if len(chat.Clients) == 0 && chatID != "general" {
			delete(h.chats, chatID)
			log.Printf("Removed inactive chat: %s", chatID)
		}
	}
}

// RemoveUserFromAllChats removes a user from all chats
func (h *Hub) RemoveUserFromAllChats(userID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for chatID, chat := range h.chats {
		if _, exists := chat.Clients[userID]; exists {
			delete(chat.Clients, userID)
			log.Printf("User %s removed from chat %s", userID, chatID)
		}
	}
}

// GetUserChats returns all chats a user is in
func (h *Hub) GetUserChats(userID string) []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var userChats []string
	for chatID, chat := range h.chats {
		if _, exists := chat.Clients[userID]; exists {
			userChats = append(userChats, chatID)
		}
	}
	return userChats
}

// IsUserInChat checks if a user is in a specific chat
func (h *Hub) IsUserInChat(userID, chatID string) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if chat, exists := h.chats[chatID]; exists {
		_, userExists := chat.Clients[userID]
		return userExists
	}
	return false
}

//-----------------------------------------------
// hub/hub.go - Add a method to handle client message through the hub

// HandleClientMessage processes a message from a client
func (h *Hub) HandleClientMessage(client *Client, rawMessage []byte) {
	var message struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		ChatID  string `json:"chatId"`
	}

	if err := json.Unmarshal(rawMessage, &message); err != nil {
		log.Printf("Error parsing message from client %s: %v", client.userID, err)
		return
	}

	switch message.Type {
	case "message":
		h.HandleChatMessage(client, message.Content, message.ChatID)
	case "typing":
		h.HandleTypingIndicator(client, message.Content, message.ChatID)
	case "join_chat":
		h.HandleJoinChat(client, message.ChatID)
	case "leave_chat":
		h.HandleLeaveChat(client, message.ChatID)
	case "create_chat":
		h.HandleCreateChat(client, message.Content)
	case "get_chats":
		h.HandleGetChats(client)
	default:
		log.Printf("Unknown message type from client %s: %s", client.userID, message.Type)
	}
}

// HandleChatMessage processes and broadcasts chat message
func (h *Hub) HandleChatMessage(client *Client, content, chatID string) {

	if content == "" {
		return
	}

	messageID, err := generateUUID()
	if err != nil {
		fmt.Printf("Error generating message ID: %v", err)
		return
	}

	// Create the message to broadcast
	chatMessage := map[string]interface{}{
		"type":      "message",
		"id":        messageID,
		"userId":    client.userID,
		"username":  client.username,
		"content":   content,
		"chatID":    chatID,
		"timestamp": time.Now(),
	}

	log.Printf("Broadcasting message from %s in chat %s: %s",
		client.username, chatID, content)

	// Broadcast to all clients in the chat
	h.BroadcastToChat(chatID, chatMessage)
}

// HandleTypingIndicator broadcasts typing status
func (h *Hub) HandleTypingIndicator(client *Client, typing, chatID string) {
	typingMessage := map[string]interface{}{
		"type":      "typing",
		"userId":    client.userID,
		"username":  client.username,
		"chatId":    chatID,
		"typing":    typing == "true",
		"timestamp": time.Now(),
	}

	h.BroadcastToChat(chatID, typingMessage)
}

// HandleJoinChat handles chat joining
func (h *Hub) HandleJoinChat(client *Client, chatID string) {
	h.JoinChat(chatID, client.userID, client)

	// Notify chat about new user
	joinMessage := map[string]interface{}{
		"type":      "user_joined",
		"userId":    client.userID,
		"username":  client.username,
		"message":   client.username + " joined the chat",
		"chatId":    chatID,
		"timestamp": time.Now(),
	}

	h.BroadcastToChat(chatID, joinMessage)
}

// HandleLeaveChat handles chat leaving
func (h *Hub) HandleLeaveChat(client *Client, chatID string) {
	h.LeaveChat(chatID, client.userID)

	leaveMessage := map[string]interface{}{
		"type":      "user_left",
		"userId":    client.userID,
		"username":  client.username,
		"message":   client.username + " left the chat",
		"chatId":    chatID,
		"timestamp": time.Now(),
	}

	h.BroadcastToChat(chatID, leaveMessage)
}

// HandleCreateChat handles chat creation
func (h *Hub) HandleCreateChat(client *Client, chatName string) {

	chatID, err := generateUUID()
	if err != nil {
		fmt.Printf("Error generating chat id: %v", err)
		return
	}

	h.CreateChat(chatID, chatName)
	h.JoinChat(chatID, client.userID, client)

	// Notify about chat creation
	chatMessage := map[string]interface{}{
		"type":      "chat_created",
		"chatId":    chatID,
		"chatName":  chatName,
		"userId":    client.userID,
		"username":  client.username,
		"timestamp": time.Now(),
	}

	h.BroadcastToAll(chatMessage)
}

// HandleGetChats handles chat list requests
func (h *Hub) HandleGetChats(client *Client) {
	chatList := h.GetChatList()

	response := map[string]interface{}{
		"type":     "chat_list",
		"chatList": chatList,
	}

	if err := client.SendMessage(response); err != nil {
		log.Printf("Failed to send chat list to client %s: %v", client.userID, err)
	}
}

func generateUUID() (string, error) {
	u7, err2 := uuid.NewV7()
	if err2 != nil {
		return "", fmt.Errorf("error generating UUIDv7: %w", err2)
	}
	return u7.String(), nil
}
