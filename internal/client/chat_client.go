package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ChatClient represents a client connection to the chat server
type ChatClient struct {
	conn         *websocket.Conn
	serverURL    string
	userID       string
	username     string
	messageChan  chan Message
	errorChan    chan error
	closeChan    chan struct{}
	rooms        map[string]bool
	currentRoom  string
	isConnected  bool
	mutex        sync.RWMutex
	messageCount int
}

// Message represents a chat message
type Message struct {
	Type      string    `json:"type"`
	ID        string    `json:"id,omitempty"`
	UserID    string    `json:"userId,omitempty"`
	Username  string    `json:"username,omitempty"`
	Content   string    `json:"content,omitempty"`
	ChatID    string    `json:"chatId,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Success   bool      `json:"success,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// Config holds client configuration
type Config struct {
	ServerURL string
	UserID    string
	Username  string
	Timeout   time.Duration
}

// NewChatClient creates a new chat client instance
func NewChatClient(config Config) *ChatClient {
	if config.UserID == "" {
		config.UserID = uuid.New().String()
	}
	if config.Username == "" {
		config.Username = "user_" + config.UserID[:8]
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &ChatClient{
		serverURL:   config.ServerURL,
		userID:      config.UserID,
		username:    config.Username,
		messageChan: make(chan Message, 100),
		errorChan:   make(chan error, 10),
		closeChan:   make(chan struct{}),
		rooms:       make(map[string]bool),
		currentRoom: "general",
	}
}

// Connect establishes a connection to the chat server
func (c *ChatClient) Connect() error {
	// Parse server URL
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %v", err)
	}

	// Add query parameters
	q := u.Query()
	q.Add("user_id", c.userID)
	q.Add("username", c.username)
	u.RawQuery = q.Encode()

	// Establish WebSocket connection
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	c.conn = conn
	c.isConnected = true

	log.Printf("Connected to chat server: %s", c.serverURL)
	log.Printf("User: %s (%s)", c.username, c.userID)

	// Start message handlers
	go c.readPump()
	go c.handleMessages()

	// Join default room
	if err := c.JoinRoom("general"); err != nil {
		return err
	}

	return nil
}

// readPump handles incoming messages from the server
func (c *ChatClient) readPump() {
	defer func() {
		c.conn.Close()
		c.isConnected = false
	}()

	for {
		select {
		case <-c.closeChan:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				c.errorChan <- fmt.Errorf("read error: %v", err)
				return
			}

			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				c.errorChan <- fmt.Errorf("invalid message format: %v", err)
				continue
			}

			c.messageChan <- msg
			c.messageCount++
		}
	}
}

// handleMessages processes incoming messages
func (c *ChatClient) handleMessages() {
	for {
		select {
		case msg := <-c.messageChan:
			c.handleMessage(msg)
		case <-c.closeChan:
			return
		}
	}
}

// handleMessage processes different message types
func (c *ChatClient) handleMessage(msg Message) {
	switch msg.Type {
	case "message":
		log.Printf("[%s] %s: %s", msg.ChatID, msg.Username, msg.Content)
	case "system":
		log.Printf("SYSTEM: %s", msg.Content)
	case "user_joined":
		log.Printf("-> %s joined the chat", msg.Username)
	case "user_left":
		log.Printf("<- %s left the chat", msg.Username)
	case "typing":
		if msg.Content == "true" {
			log.Printf("%s is typing...", msg.Username)
		}
	case "room_joined":
		log.Printf("Joined room: %s", msg.ChatID)
	case "room_list":
		//if rooms, ok := msg.Content.(map[string]interface{}); ok {
		//	log.Printf("Available rooms: %v", rooms)
		//}
	case "error":
		log.Printf("ERROR: %s", msg.Error)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// SendMessage sends a chat message to the current room
func (c *ChatClient) SendMessage(content string) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type:      "message",
		Content:   content,
		ChatID:    c.currentRoom,
		UserID:    c.userID,
		Username:  c.username,
		Timestamp: time.Now(),
	}

	return c.sendJSON(message)
}

// JoinRoom joins a specific chat room
func (c *ChatClient) JoinRoom(roomID string) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type:   "join_room",
		ChatID: roomID,
		UserID: c.userID,
	}

	if err := c.sendJSON(message); err != nil {
		return err
	}

	c.mutex.Lock()
	c.rooms[roomID] = true
	c.currentRoom = roomID
	c.mutex.Unlock()

	log.Printf("Joining room: %s", roomID)
	return nil
}

// LeaveRoom leaves a specific chat room
func (c *ChatClient) LeaveRoom(roomID string) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type:   "leave_room",
		ChatID: roomID,
		UserID: c.userID,
	}

	if err := c.sendJSON(message); err != nil {
		return err
	}

	c.mutex.Lock()
	delete(c.rooms, roomID)
	if c.currentRoom == roomID {
		c.currentRoom = "general"
	}
	c.mutex.Unlock()

	log.Printf("Left room: %s", roomID)
	return nil
}

// CreateRoom creates a new chat room
func (c *ChatClient) CreateRoom(roomName string) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	roomID := generateRoomID(roomName)

	message := Message{
		Type:     "create_room",
		ChatID:   roomID,
		Content:  roomName,
		UserID:   c.userID,
		Username: c.username,
	}

	return c.sendJSON(message)
}

// ListRooms requests the list of available rooms
func (c *ChatClient) ListRooms() error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type: "get_rooms",
	}

	return c.sendJSON(message)
}

// SendTypingIndicator sends a typing indicator
func (c *ChatClient) SendTypingIndicator(typing bool) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type:    "typing",
		ChatID:  c.currentRoom,
		UserID:  c.userID,
		Content: fmt.Sprintf("%t", typing),
	}

	return c.sendJSON(message)
}

// sendJSON sends a JSON message to the server
func (c *ChatClient) sendJSON(message interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return errors.New("not connected to server")
	}

	return c.conn.WriteJSON(message)
}

// Close gracefully closes the connection
func (c *ChatClient) Close() error {
	close(c.closeChan)
	c.isConnected = false

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsConnected returns the connection status
func (c *ChatClient) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}

// GetCurrentRoom returns the current room ID
func (c *ChatClient) GetCurrentRoom() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.currentRoom
}

// GetRooms returns the list of rooms the client is in
func (c *ChatClient) GetRooms() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	rooms := make([]string, 0, len(c.rooms))
	for room := range c.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

// GetMessageCount returns the number of messages received
func (c *ChatClient) GetMessageCount() int {
	return c.messageCount
}

// WaitForInterrupt waits for OS interrupt signals
func (c *ChatClient) WaitForInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	log.Println("Received interrupt signal, closing connection...")
	c.Close()
}

// generateRoomID creates a URL-friendly room ID from a name
func generateRoomID(name string) string {
	// Simple implementation - in production, use a proper slug generator
	return url.QueryEscape(name)
}
