package hub

import (
	"log"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

// Chat represents a chat
type Chat struct {
	ID        uuid.UUID
	Name      string
	Clients   map[uuid.UUID]*Client // userID -> Client
	CreatedAt time.Time
}

// Message represents the data structure of a chat message.
// پیام، ساختار داده‌ای یک پیام چت را نشان می‌دهد.
type Message struct {
	ChatID  uuid.UUID `json:"chatID"`
	UserID  uuid.UUID `json:"userID"`
	Content string    `json:"content"`
}

// Hub manages all connected clients and chats
type Hub struct {
	chats     map[uuid.UUID]*Chat
	clients   map[uuid.UUID]*Client // userID -> Client
	mutex     sync.RWMutex
	startTime time.Time
	// Added a channel to send messages to the Manager.
	// یک کانال برای ارسال پیام‌ها به Manager اضافه شده است.
	messagesToManager chan *Message
}

// NewHub creates a new Hub instance
func NewHub(messages chan *Message) *Hub {

	hub := &Hub{
		chats:             make(map[uuid.UUID]*Chat),
		clients:           make(map[uuid.UUID]*Client),
		startTime:         time.Now(),
		messagesToManager: messages,
	}

	// create default chat
	hub.CreateChat(config.ChatID1, "Admin Chat")

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

// RegisterClient adds a chat_client to the hub
func (h *Hub) RegisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client.userID] = client
	log.Printf("Client registered: %s. Total clients: %d", client.userID, len(h.clients))
}

// UnregisterClient removes a chat_client from the hub and all chats
func (h *Hub) UnregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Remove chat_client from all chats
	for chatID := range client.chats {
		if chat, exists := h.chats[chatID]; exists {
			delete(chat.Clients, client.userID)
			log.Printf("Client %s removed from chat %s", client.userID, chatID)
		}
	}

	// Remove chat_client from main clients map
	delete(h.clients, client.userID)
	log.Printf("Client unregistered: %s. Remaining clients: %d", client.userID, len(h.clients))
}

// LeaveChat removes a chat_client from a specific chat
func (h *Hub) LeaveChat(chatID, userID uuid.UUID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if chat, exists := h.chats[chatID]; exists {
		if _, clientExists := chat.Clients[userID]; clientExists {
			delete(chat.Clients, userID)
			log.Printf("User %s left chat %s. Users remaining: %d", userID, chatID, len(chat.Clients))

			// Also remove from chat_client's chat list
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

// JoinChat adds a chat_client to a chat
func (h *Hub) JoinChat(chatID, userID uuid.UUID, client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Ensure chat exists
	if _, exists := h.chats[chatID]; !exists {
		log.Printf("Chat %s does not exist, creating it", chatID)
		h.chats[chatID] = &Chat{
			ID:        chatID,
			Name:      chatID.String(), // Use ID as name for auto-created chats
			Clients:   make(map[uuid.UUID]*Client),
			CreatedAt: time.Now(),
		}
	}

	chat := h.chats[chatID]

	// Add chat_client to the chat
	chat.Clients[userID] = client
	client.JoinChat(chatID)

	log.Printf("User %s joined chat %s (Total in chat: %d)", userID, chatID, len(chat.Clients))
}

// CreateChat creates a new chat
func (h *Hub) CreateChat(chatID uuid.UUID, chatName string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.chats[chatID]; !exists {
		h.chats[chatID] = &Chat{
			ID:        chatID,
			Name:      chatName,
			Clients:   make(map[uuid.UUID]*Client),
			CreatedAt: time.Now(),
		}
		log.Printf("Chat created: %s (%s). Total chats: %d", chatName, chatID, len(h.chats))
	}
}

// BroadcastToChat sends a message to all clients in a chat
func (h *Hub) BroadcastToChat(chatID uuid.UUID, message interface{}) {

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
			// Check if chat_client is still connected and channel is not full
			select {
			case client.send <- messageBytes:
				// Message sent successfully
				log.Printf("Message sent to user %s in chat %s", userID, chatID)
			default:
				// Channel is full, chat_client might be disconnected
				log.Printf("Client %s send buffer full, potentially disconnected", userID)
			}
		}
	} else {
		log.Printf("Chat %s not found for broadcasting", chatID)
	}
}

// GetChatList returns a map of chat IDs to chat names
func (h *Hub) GetChatList() map[uuid.UUID]string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	chatList := make(map[uuid.UUID]string)
	for id, chat := range h.chats {
		chatList[id] = chat.Name
	}
	return chatList
}

// GetChatUsers returns the list of user IDs in a specific chat
func (h *Hub) GetChatUsers(chatID uuid.UUID) []uuid.UUID {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if chat, exists := h.chats[chatID]; exists {
		users := make([]uuid.UUID, 0, len(chat.Clients))
		for userID := range chat.Clients {
			users = append(users, userID)
		}
		return users
	}
	return []uuid.UUID{}
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
		stats[chatID.String()] = map[string]interface{}{
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
func (h *Hub) GetChat(chatID uuid.UUID) (*Chat, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	chat, exists := h.chats[chatID]
	return chat, exists
}

// ChatExists checks if a chat exists
func (h *Hub) ChatExists(chatID uuid.UUID) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	_, exists := h.chats[chatID]
	return exists
}

// GetClient returns a chat_client by user ID
func (h *Hub) GetClient(userID uuid.UUID) (*Client, bool) {
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
		if len(chat.Clients) == 0 && chatID.String() != "general" {
			delete(h.chats, chatID)
			log.Printf("Removed inactive chat: %s", chatID)
		}
	}
}

// RemoveUserFromAllChats removes a user from all chats
func (h *Hub) RemoveUserFromAllChats(userID uuid.UUID) {
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
func (h *Hub) GetUserChats(userID uuid.UUID) []uuid.UUID {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var userChats []uuid.UUID
	for chatID, chat := range h.chats {
		if _, exists := chat.Clients[userID]; exists {
			userChats = append(userChats, chatID)
		}
	}
	return userChats
}

// IsUserInChat checks if a user is in a specific chat
func (h *Hub) IsUserInChat(userID, chatID uuid.UUID) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if chat, exists := h.chats[chatID]; exists {
		_, userExists := chat.Clients[userID]
		return userExists
	}
	return false
}
