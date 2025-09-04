package hub

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected user
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	userID string
	send   chan []byte
	chats  map[string]bool // Track which chats the user is in
	mutex  sync.RWMutex
}

// NewClient creates a new client instance
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
		chats:  make(map[string]bool),
	}
}

// IsInChat checks if the client is in a specific chat
func (c *Client) IsInChat(chatID string) bool {
	_, ok := c.chats[chatID]
	return ok
}

// JoinChat adds the client to a chat
func (c *Client) JoinChat(chatID string) {
	c.chats[chatID] = true
}

// LeaveChat removes the client from a chat
func (c *Client) LeaveChat(chatID string) {
	delete(c.chats, chatID)
}

// GetChats returns all chats the client is in
func (c *Client) GetChats() []string {
	chats := make([]string, 0, len(c.chats))
	for chatID := range c.chats {
		chats = append(chats, chatID)
	}
	return chats
}

// ReadPump handles message from the WebSocket connection
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

// WritePump sends message to the WebSocket connection
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

			// Add queued message to the current WebSocket message
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

// handleMessage processes different types of incoming message
func (c *Client) handleMessage(rawMessage []byte) {

	// Forward all message to the hub for processing
	c.hub.HandleClientMessage(c, rawMessage)

	// Instead of handling message here, forward them to the hub/handler
	// This ensures consistent message processing

	//c.mutex.RLock()
	//handler := c.messageHandler
	//c.mutex.RUnlock()
	//
	//if handler != nil {
	//	// Use external message handler
	//	handler(c, rawMessage)
	//} else {
	//	// Fallback to local handling
	//	c.handleMessageLocally(rawMessage)
	//}
}

// handleMessageLocally handles message when no external handler is set
//func (c *Client) handleMessageLocally(rawMessage []byte) {
//	var baseMessage struct {
//		Type string `json:"type"`
//	}
//
//	if err := json.Unmarshal(rawMessage, &baseMessage); err != nil {
//		log.Printf("Error parsing message from client %s: %v", c.userID, err)
//		return
//	}
//
//	// Basic local handling for critical message
//	switch baseMessage.Type {
//	case "ping":
//		// Respond to ping
//		err := c.SendMessage(map[string]interface{}{
//			"type": "pong",
//		})
//		if err != nil {
//			return
//		}
//	default:
//		log.Printf("No message handler for type: %s from client %s", baseMessage.Type, c.userID)
//	}
//}

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

// UserID returns the client's user ID
func (c *Client) UserID() string {
	return c.userID
}
