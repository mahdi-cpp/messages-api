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
	hub := &Hub{
		rooms:     make(map[string]*Room),
		clients:   make(map[string]*Client),
		startTime: time.Now(),
	}

	// Create default room
	hub.CreateRoom("general", "General Chat")

	return hub
}

// Run starts the hub (maintain for compatibility)
func (h *Hub) Run() {
	log.Println("Hub started and running")
	// This method can be used for background tasks if needed
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		h.CleanupInactiveRooms()
	}
}

// RegisterClient adds a client to the hub
func (h *Hub) RegisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client.userID] = client
	log.Printf("Client registered: %s. Total clients: %d", client.userID, len(h.clients))
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
	log.Printf("Client unregistered: %s. Remaining clients: %d", client.userID, len(h.clients))
}

// LeaveRoom removes a client from a specific room
func (h *Hub) LeaveRoom(roomID, userID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if room, exists := h.rooms[roomID]; exists {
		if _, clientExists := room.Clients[userID]; clientExists {
			delete(room.Clients, userID)
			log.Printf("User %s left room %s. Users remaining: %d", userID, roomID, len(room.Clients))

			// Also remove from client's room list
			if client, clientExists := h.clients[userID]; clientExists {
				client.LeaveRoom(roomID)
			}
		} else {
			log.Printf("User %s not found in room %s", userID, roomID)
		}
	} else {
		log.Printf("Room %s not found", roomID)
	}
}

// JoinRoom adds a client to a room
func (h *Hub) JoinRoom(roomID, userID string, client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Ensure room exists
	if _, exists := h.rooms[roomID]; !exists {
		log.Printf("Room %s does not exist, creating it", roomID)
		h.rooms[roomID] = &Room{
			ID:        roomID,
			Name:      roomID, // Use ID as name for auto-created rooms
			Clients:   make(map[string]*Client),
			CreatedAt: time.Now(),
		}
	}

	room := h.rooms[roomID]

	// Add client to the room
	room.Clients[userID] = client
	client.JoinRoom(roomID)

	log.Printf("User %s joined room %s (Total in room: %d)", userID, roomID, len(room.Clients))
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
		log.Printf("Room created: %s (%s). Total rooms: %d", roomName, roomID, len(h.rooms))
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
			// Check if client is still connected and channel is not full
			select {
			case client.send <- messageBytes:
				// Message sent successfully
				log.Printf("Message sent to user %s in room %s", userID, roomID)
			default:
				// Channel is full, client might be disconnected
				log.Printf("Client %s send buffer full, potentially disconnected", userID)
			}
		}
	} else {
		log.Printf("Room %s not found for broadcasting", roomID)
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

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
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
		"total_rooms":   len(h.rooms),
		"uptime":        time.Since(h.startTime).String(),
		"start_time":    h.startTime,
		"room_stats":    h.GetRoomStats(),
	}
}

// GetRoom returns a room by ID
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
			log.Printf("Client %s send buffer full", userID)
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

// RemoveUserFromAllRooms removes a user from all rooms
func (h *Hub) RemoveUserFromAllRooms(userID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for roomID, room := range h.rooms {
		if _, exists := room.Clients[userID]; exists {
			delete(room.Clients, userID)
			log.Printf("User %s removed from room %s", userID, roomID)
		}
	}
}

// GetUserRooms returns all rooms a user is in
func (h *Hub) GetUserRooms(userID string) []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var userRooms []string
	for roomID, room := range h.rooms {
		if _, exists := room.Clients[userID]; exists {
			userRooms = append(userRooms, roomID)
		}
	}
	return userRooms
}

// IsUserInRoom checks if a user is in a specific room
func (h *Hub) IsUserInRoom(userID, roomID string) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if room, exists := h.rooms[roomID]; exists {
		_, userExists := room.Clients[userID]
		return userExists
	}
	return false
}
