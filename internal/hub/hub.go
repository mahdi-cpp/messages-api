package hub

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Room represents a chat room
type Room struct {
	ID        string
	Name      string
	Clients   map[string]*Client // userID -> Client
	CreatedAt time.Time
}

// Hub manages all connected clients and rooms
type Hub struct {
	rooms     map[string]*Room
	clients   map[string]*Client // userID -> Client
	mutex     sync.RWMutex
	startTime time.Time
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		rooms:     make(map[string]*Room),
		clients:   make(map[string]*Client),
		startTime: time.Now(),
	}
}

// Run starts the hub (maintain for compatibility)
func (h *Hub) Run() {
	log.Println("Hub started and running")
	// This method can be used for background tasks if needed
}

// RegisterClient adds a client to the hub
func (h *Hub) RegisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client.userID] = client
	log.Printf("Client registered: %s", client.userID)
}

// UnregisterClient removes a client from the hub and all rooms
func (h *Hub) UnregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Remove client from all rooms
	for roomID := range client.rooms {
		if room, exists := h.rooms[roomID]; exists {
			delete(room.Clients, client.userID)
			log.Printf("Client %s removed from room %s", client.userID, roomID)
		}
	}

	// Remove client from main clients map
	delete(h.clients, client.userID)
	log.Printf("Client unregistered: %s", client.userID)
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// CreateRoom creates a new chat room
func (h *Hub) CreateRoom(roomID, roomName string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.rooms[roomID]; !exists {
		h.rooms[roomID] = &Room{
			ID:        roomID,
			Name:      roomName,
			Clients:   make(map[string]*Client),
			CreatedAt: time.Now(),
		}
		log.Printf("Room created: %s (%s)", roomName, roomID)
	}
}

// JoinRoom adds a client to a room
func (h *Hub) JoinRoom(roomID, userID string, client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if room, exists := h.rooms[roomID]; exists {
		room.Clients[userID] = client
		client.JoinRoom(roomID)
		log.Printf("User %s joined room %s", userID, roomID)
	}
}

// LeaveRoom removes a client from a room
func (h *Hub) LeaveRoom(roomID, userID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if room, exists := h.rooms[roomID]; exists {
		delete(room.Clients, userID)
		log.Printf("User %s left room %s", userID, roomID)
	}
}

// BroadcastToRoom sends a message to all clients in a room
func (h *Hub) BroadcastToRoom(roomID string, message interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Convert message to JSON bytes
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	if room, exists := h.rooms[roomID]; exists {
		for userID, client := range room.Clients {
			select {
			case client.send <- messageBytes:
				// Message sent successfully
			default:
				// Channel is full, client might be disconnected
				log.Printf("Client %s send buffer full, disconnecting", userID)
				close(client.send)
				delete(room.Clients, userID)
				delete(h.clients, userID)
			}
		}
	}
}

// GetRoomList returns a map of room IDs to room names
func (h *Hub) GetRoomList() map[string]string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	roomList := make(map[string]string)
	for id, room := range h.rooms {
		roomList[id] = room.Name
	}
	return roomList
}

// GetRoomUsers returns the list of user IDs in a specific room
func (h *Hub) GetRoomUsers(roomID string) []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if room, exists := h.rooms[roomID]; exists {
		users := make([]string, 0, len(room.Clients))
		for userID := range room.Clients {
			users = append(users, userID)
		}
		return users
	}
	return []string{}
}

// GetRoomStats returns statistics for all rooms
func (h *Hub) GetRoomStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	stats := make(map[string]interface{})
	for roomID, room := range h.rooms {
		stats[roomID] = map[string]interface{}{
			"name":       room.Name,
			"user_count": len(room.Clients),
			"created_at": room.CreatedAt,
			"age":        time.Since(room.CreatedAt).String(),
		}
	}
	return stats
}

// GetStartTime returns the hub start time
func (h *Hub) GetStartTime() time.Time {
	return h.startTime
}

// GetRoom returns a room by ID (thread-safe)
func (h *Hub) GetRoom(roomID string) (*Room, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	room, exists := h.rooms[roomID]
	return room, exists
}

// RoomExists checks if a room exists
func (h *Hub) RoomExists(roomID string) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	_, exists := h.rooms[roomID]
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
			// Channel is full, client might be disconnected
			log.Printf("Client %s send buffer full, disconnecting", userID)
			close(client.send)
			delete(h.clients, userID)
		}
	}
}

// CleanupInactiveRooms removes rooms with no clients (optional)
func (h *Hub) CleanupInactiveRooms() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for roomID, room := range h.rooms {
		if len(room.Clients) == 0 && roomID != "general" {
			delete(h.rooms, roomID)
			log.Printf("Removed inactive room: %s", roomID)
		}
	}
}

// GetHubStats returns comprehensive hub statistics
func (h *Hub) GetHubStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return map[string]interface{}{
		"total_clients": len(h.clients),
		"total_rooms":   len(h.rooms),
		"uptime":        time.Since(h.startTime).String(),
		"start_time":    h.startTime,
		"room_stats":    h.GetRoomStats(),
	}
}
