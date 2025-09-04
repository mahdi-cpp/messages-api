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

	"github.com/gorilla/websocket"
)

// ChatClient represents a client connection to the chat server
type ChatClient struct {
	conn         *websocket.Conn
	serverURL    string
	userID       string
	UserID       string
	messageChan  chan Message
	errorChan    chan error
	closeChan    chan struct{}
	chats        map[string]bool
	currentChat  string
	isConnected  bool
	mutex        sync.RWMutex
	messageCount int
}

// Message represents a chat message
type Message struct {
	Type      string    `json:"type"`
	ID        string    `json:"id,omitempty"`
	UserID    string    `json:"userId,omitempty"`
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
	Timeout   time.Duration
}

// NewChatClient creates a new chat client instance
func NewChatClient(config Config) (*ChatClient, error) {

	if config.UserID == "" {
		return nil, errors.New("user id is required")
	}
	//if config.UserID == "" {
	//	return nil, errors.New("UserID is required")
	//}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &ChatClient{
		serverURL: config.ServerURL,
		userID:    config.UserID,
		//UserID:    config.UserID,
		messageChan: make(chan Message, 100),
		errorChan:   make(chan error, 10),
		closeChan:   make(chan struct{}),
		chats:       make(map[string]bool),
		currentChat: "general",
	}, nil
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
	q.Add("UserID", c.UserID)
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
	log.Printf("User: %s (%s)", c.UserID, c.userID)

	// Start message handlers
	go c.readPump()
	go c.handleMessages()

	// Join default chat
	if err := c.JoinChat("general"); err != nil {
		return err
	}

	return nil
}

// readPump handles incoming message from the server
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

// handleMessages processes incoming message
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
		log.Printf("{ChatID:%s} %s: %s", msg.ChatID, msg.UserID, msg.Content)
	case "system":
		log.Printf("SYSTEM: %s", msg.Content)
	case "user_joined":
		log.Printf("-> %s joined the chat", msg.UserID)
	case "user_left":
		log.Printf("<- %s left the chat", msg.UserID)
	case "typing":
		log.Printf("%s is typing...", msg.UserID)
	case "seen":
		log.Printf("%s is seen message", msg.UserID)
	case "chat_joined":
		log.Printf("Joined chat: %s", msg.ChatID)
	case "chat_list":
		//if chats, ok := msg.Content.(map[string]interface{}); ok {
		//	log.Printf("Available chats: %v", chats)
		//}
	case "error":
		log.Printf("ERROR: %s", msg.Error)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// SendMessage sends a chat message to the current chat
func (c *ChatClient) SendMessage(content string) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type:      "message",
		Content:   content,
		ChatID:    c.currentChat,
		UserID:    c.userID,
		Timestamp: time.Now(),
	}

	return c.sendJSON(message)
}

// JoinChat joins a specific chat
func (c *ChatClient) JoinChat(chatID string) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type:   "join_chat",
		ChatID: chatID,
		UserID: c.userID,
	}

	if err := c.sendJSON(message); err != nil {
		return err
	}

	c.mutex.Lock()
	c.chats[chatID] = true
	c.currentChat = chatID
	c.mutex.Unlock()

	log.Printf("Joining chat: %s", chatID)
	return nil
}

// LeaveChat leaves a specific chat
func (c *ChatClient) LeaveChat(chatID string) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type:   "leave_chat",
		ChatID: chatID,
		UserID: c.userID,
	}

	if err := c.sendJSON(message); err != nil {
		return err
	}

	c.mutex.Lock()
	delete(c.chats, chatID)
	if c.currentChat == chatID {
		c.currentChat = "general"
	}
	c.mutex.Unlock()

	log.Printf("Left chat: %s", chatID)
	return nil
}

// CreateChat creates a new chat
func (c *ChatClient) CreateChat(chatName string) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	chatID := generateChatID(chatName)

	message := Message{
		Type:    "create_chat",
		ChatID:  chatID,
		Content: chatName,
		UserID:  c.userID,
	}

	return c.sendJSON(message)
}

// ListChats requests the list of available chats
func (c *ChatClient) ListChats() error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type: "get_chats",
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
		ChatID:  c.currentChat,
		UserID:  c.userID,
		Content: fmt.Sprintf("%t", typing),
	}

	return c.sendJSON(message)
}

func (c *ChatClient) SendSeenIndicator() error {

	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message := Message{
		Type:   "seen",
		ChatID: c.currentChat,
		UserID: c.userID,
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

// GetCurrentChat returns the current chat ID
func (c *ChatClient) GetCurrentChat() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.currentChat
}

// GetChats returns the list of chats the client is in
func (c *ChatClient) GetChats() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	chats := make([]string, 0, len(c.chats))
	for chat := range c.chats {
		chats = append(chats, chat)
	}
	return chats
}

// GetMessageCount returns the number of message received
func (c *ChatClient) GetMessageCount() int {
	return c.messageCount
}

// WaitForInterrupt waits for OS interrupt signals
func (c *ChatClient) WaitForInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	log.Println("Received interrupt signal, closing connection...")
	err := c.Close()
	if err != nil {
		return
	}
}

// generateChatID creates a URL-friendly chat ID from a name
func generateChatID(name string) string {
	// Simple implementation - in production, use a proper slug generator
	return url.QueryEscape(name)
}
