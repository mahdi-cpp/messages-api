package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a connected user
type Client struct {
	hub            *Hub
	conn           *websocket.Conn
	userID         string
	username       string
	send           chan []byte
	rooms          map[string]bool       // Track which rooms the user is in
	messageHandler func(*Client, []byte) // Add this field
	mutex          sync.RWMutex
}

// NewClient creates a new client instance
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
		rooms:  make(map[string]bool),
	}
}

// SetMessageHandler sets the external message handler
func (c *Client) SetMessageHandler(handler func(*Client, []byte)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.messageHandler = handler
}

// IsInRoom checks if the client is in a specific room
func (c *Client) IsInRoom(roomID string) bool {
	_, ok := c.rooms[roomID]
	return ok
}

// JoinRoom adds the client to a room
func (c *Client) JoinRoom(roomID string) {
	c.rooms[roomID] = true
}

// LeaveRoom removes the client from a room
func (c *Client) LeaveRoom(roomID string) {
	delete(c.rooms, roomID)
}

// GetRooms returns all rooms the client is in
func (c *Client) GetRooms() []string {
	rooms := make([]string, 0, len(c.rooms))
	for roomID := range c.rooms {
		rooms = append(rooms, roomID)
	}
	return rooms
}

// ReadPump handles messages from the WebSocket connection
func (c *Client) ReadPump() {

	defer func() {
		// Clean up when client disconnects
		c.hub.UnregisterClient(c)
		c.conn.Close()
		log.Printf("Client %s disconnected", c.userID)
	}()

	// Configure connection settings
	c.conn.SetReadLimit(10 * 1024) // Max message size 10KB in bytes
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error from client %s: %v", c.userID, err)
			}
			break
		}

		// Process incoming message
		c.handleMessage(message)
	}
}

// WritePump sends messages to the WebSocket connection
func (c *Client) WritePump() {

	ticker := time.NewTicker(54 * time.Second) // Ping interval
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:

			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error creating writer for client %s: %v", c.userID, err)
				return
			}
			writer.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				writer.Write(<-c.send)
			}

			if err := writer.Close(); err != nil {
				log.Printf("Error closing writer for client %s: %v", c.userID, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping to client %s: %v", c.userID, err)
				return
			}
		}
	}
}

// handleMessage processes different types of incoming messages
func (c *Client) handleMessage(rawMessage []byte) {

	// Instead of handling messages here, forward them to the hub/handler
	// This ensures consistent message processing

	c.mutex.RLock()
	handler := c.messageHandler
	c.mutex.RUnlock()

	if handler != nil {
		// Use external message handler
		handler(c, rawMessage)
	} else {
		// Fallback to local handling
		c.handleMessageLocally(rawMessage)
	}
}

// handleMessageLocally handles messages when no external handler is set
func (c *Client) handleMessageLocally(rawMessage []byte) {
	var baseMessage struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(rawMessage, &baseMessage); err != nil {
		log.Printf("Error parsing message from client %s: %v", c.userID, err)
		return
	}

	// Basic local handling for critical messages
	switch baseMessage.Type {
	case "ping":
		// Respond to ping
		err := c.SendMessage(map[string]interface{}{
			"type": "pong",
		})
		if err != nil {
			return
		}
	default:
		log.Printf("No message handler for type: %s from client %s", baseMessage.Type, c.userID)
	}
}

// handleTypingMessage processes typing indicators
func (c *Client) handleTypingMessage(rawMessage []byte) {
	var typingStatus struct {
		ChatID string `json:"chatId"`
		Typing bool   `json:"typing"`
	}

	if err := json.Unmarshal(rawMessage, &typingStatus); err != nil {
		log.Printf("Error parsing typing status from client %s: %v", c.userID, err)
		return
	}

	// Broadcast typing status to all users in the room
	c.hub.BroadcastToRoom(typingStatus.ChatID, map[string]interface{}{
		"type":   "typing",
		"userId": c.userID,
		"chatId": typingStatus.ChatID,
		"typing": typingStatus.Typing,
	})
}

// handleChatMessage processes chat messages
func (c *Client) handleChatMessage(rawMessage []byte) {

	var message struct {
		Content string `json:"content"`
		ChatID  string `json:"chatId"`
	}

	if err := json.Unmarshal(rawMessage, &message); err != nil {
		log.Printf("Error parsing chat message from client %s: %v", c.userID, err)
		return
	}

	chatID, err := generateUUID()
	if err != nil {
		return
	}

	// Create message with metadata
	chatMessage := map[string]interface{}{
		"type":      "message",
		"id":        chatID,
		"userId":    c.userID,
		"content":   message.Content,
		"chatId":    message.ChatID,
		"timestamp": time.Now(),
	}

	// Broadcast message to all users in the room
	c.hub.BroadcastToRoom(message.ChatID, chatMessage)

	log.Printf("Message from client %s in room %s: %s", c.userID, message.ChatID, message.Content)
}

// handleJoinRoom processes room join requests
func (c *Client) handleJoinRoom(rawMessage []byte) {

	var joinRequest struct {
		RoomID string `json:"roomId"`
	}

	if err := json.Unmarshal(rawMessage, &joinRequest); err != nil {
		log.Printf("Error parsing join_room request from client %s: %v", c.userID, err)
		return
	}

	// Join the room through the hub
	c.hub.JoinRoom(joinRequest.RoomID, c.userID, c)

	// Send confirmation to client
	confirmation := map[string]interface{}{
		"type":    "room_joined",
		"roomId":  joinRequest.RoomID,
		"success": true,
	}

	if message, err := json.Marshal(confirmation); err == nil {
		c.send <- message
	}

	log.Printf("Client %s joined room %s", c.userID, joinRequest.RoomID)
}

// handleLeaveRoom processes room leave requests
func (c *Client) handleLeaveRoom(rawMessage []byte) {
	var leaveRequest struct {
		RoomID string `json:"roomId"`
	}

	if err := json.Unmarshal(rawMessage, &leaveRequest); err != nil {
		log.Printf("Error parsing leave_room request from client %s: %v", c.userID, err)
		return
	}

	// Leave the room through the hub
	//c.hub.LeaveRoom(leaveRequest.RoomID, c.userID)

	// Send confirmation to client
	confirmation := map[string]interface{}{
		"type":    "room_left",
		"roomId":  leaveRequest.RoomID,
		"success": true,
	}

	if message, err := json.Marshal(confirmation); err == nil {
		c.send <- message
	}

	log.Printf("Client %s left room %s", c.userID, leaveRequest.RoomID)
}

// handleCreateRoom processes room creation requests
func (c *Client) handleCreateRoom(rawMessage []byte) {

	var createRequest struct {
		RoomID   string `json:"roomId"`
		RoomName string `json:"roomName"`
	}

	if err := json.Unmarshal(rawMessage, &createRequest); err != nil {
		log.Printf("Error parsing create_room request from client %s: %v", c.userID, err)
		return
	}

	// Create the room through the hub
	c.hub.CreateRoom(createRequest.RoomID, createRequest.RoomName)

	// Auto-join the room after creation
	c.hub.JoinRoom(createRequest.RoomID, c.userID, c)

	// Send confirmation to client
	confirmation := map[string]interface{}{
		"type":     "room_created",
		"roomId":   createRequest.RoomID,
		"roomName": createRequest.RoomName,
		"success":  true,
	}

	if message, err := json.Marshal(confirmation); err == nil {
		c.send <- message
	}

	log.Printf("Client %s created room %s (%s)", c.userID, createRequest.RoomName, createRequest.RoomID)
}

// handleGetRooms sends the list of available rooms to the client
func (c *Client) handleGetRooms() {
	roomList := c.hub.GetRoomList()

	response := map[string]interface{}{
		"type":     "room_list",
		"roomList": roomList,
	}

	if message, err := json.Marshal(response); err == nil {
		c.send <- message
	}

	log.Printf("Sent room list to client %s", c.userID)
}

// handleGetRoomUsers sends the list of users in a specific room
func (c *Client) handleGetRoomUsers(rawMessage []byte) {

	var request struct {
		RoomID string `json:"roomId"`
	}

	if err := json.Unmarshal(rawMessage, &request); err != nil {
		log.Printf("Error parsing get_room_users request from client %s: %v", c.userID, err)
		return
	}

	// Get users from the hub (you'll need to implement GetRoomUsers in hub.go)
	users := c.hub.GetRoomUsers(request.RoomID)

	response := map[string]interface{}{
		"type":   "room_users",
		"roomId": request.RoomID,
		"users":  users,
	}

	if message, err := json.Marshal(response); err == nil {
		c.send <- message
	}

	log.Printf("Sent user list for room %s to client %s", request.RoomID, c.userID)
}

// SendMessage sends a message directly to this client
func (c *Client) SendMessage(message interface{}) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	select {
	case c.send <- messageBytes:
		return nil
	default:
		// Channel is full, client might be disconnected
		return ErrClientSendBufferFull
	}
}

// Close gracefully closes the client connection
func (c *Client) Close() {
	close(c.send)
	c.conn.Close()
}

// ErrClientSendBufferFull Custom errors
var (
	ErrClientSendBufferFull = errors.New("client send buffer is full")
)

// Username returns the client's username
func (c *Client) Username() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.username
}

// UserID returns the client's user ID
func (c *Client) UserID() string {
	return c.userID
}

// SetUserInfo sets the username for the client
func (c *Client) SetUserInfo(username string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.username = username
}

// generateUUID generates a unique identifier for messages
func generateUUID() (string, error) {

	u7, err2 := uuid.NewV7()
	if err2 != nil {
		return "", fmt.Errorf("error generating UUIDv7: %w", err2)
	}

	return u7.String(), nil
}
