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

// NewClient creates a new chat_client instance
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
		chats:  make(map[string]bool),
	}
}

// handleMessage processes different types of incoming message
func (c *Client) handleMessage(rawMessage []byte) {
	// Forward all message to the hub for processing
	c.hub.HandleClientMessage(c, rawMessage)
}

// IsInChat checks if the chat_client is in a specific chat
func (c *Client) IsInChat(chatID string) bool {
	_, ok := c.chats[chatID]
	return ok
}

// JoinChat adds the chat_client to a chat
func (c *Client) JoinChat(chatID string) {
	c.chats[chatID] = true
}

// LeaveChat removes the chat_client from a chat
func (c *Client) LeaveChat(chatID string) {
	delete(c.chats, chatID)
}

// GetChats returns all chats the chat_client is in
func (c *Client) GetChats() []string {
	chats := make([]string, 0, len(c.chats))
	for chatID := range c.chats {
		chats = append(chats, chatID)
	}
	return chats
}

// SendMessage sends a message directly to this chat_client
func (c *Client) SendMessage(message interface{}) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	select {
	case c.send <- messageBytes:
		return nil
	default:
		// Channel is full, chat_client might be disconnected
		return ErrClientSendBufferFull
	}
}

// Close gracefully closes the chat_client connection
func (c *Client) Close() {
	close(c.send)
	c.conn.Close()
}

// UserID returns the chat_client's user ID
func (c *Client) UserID() string {
	return c.userID
}

// ReadPump handles message from the WebSocket connection
func (c *Client) ReadPump() {

	defer func() {
		// Clean up when chat_client disconnects
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
				log.Printf("WebSocket error from chat_client %s: %v", c.userID, err)
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
				log.Printf("Error creating writer for chat_client %s: %v", c.userID, err)
				return
			}
			writer.Write(message)

			// Add queued message to the current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				writer.Write(<-c.send)
			}

			if err := writer.Close(); err != nil {
				log.Printf("Error closing writer for chat_client %s: %v", c.userID, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping to chat_client %s: %v", c.userID, err)
				return
			}
		}
	}
}
