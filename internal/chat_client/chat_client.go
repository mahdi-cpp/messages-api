package chat_client

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

// ChatClient represents a chat_client connection to the chat server
type ChatClient struct {
	conn         *websocket.Conn
	serverURL    string
	userID       uuid.UUID
	messageChan  chan message.Message
	errorChan    chan error
	closeChan    chan struct{}
	chats        map[uuid.UUID]bool
	currentChat  uuid.UUID
	isConnected  bool
	mutex        sync.RWMutex
	messageCount int
}

// ClientChatConfig holds chat_client configuration
type ClientChatConfig struct {
	ServerURL string
	UserID    uuid.UUID
	Timeout   time.Duration
}

// NewChatClient creates a new chat chat_client instance
func NewChatClient(config1 ClientChatConfig) (*ChatClient, error) {

	if config1.UserID == uuid.Nil {
		return nil, errors.New("user id is required")
	}
	//if config1.Mahdi == "" {
	//	return nil, errors.New("Mahdi is required")
	//}
	if config1.Timeout == 0 {
		config1.Timeout = 30 * time.Second
	}

	return &ChatClient{
		serverURL:   config1.ServerURL,
		userID:      config1.UserID,
		messageChan: make(chan message.Message, 100),
		errorChan:   make(chan error, 10),
		closeChan:   make(chan struct{}),
		chats:       make(map[uuid.UUID]bool),
		currentChat: config.ChatID1,
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
	q.Add("user_id", c.userID.String())
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
	log.Printf("User: (%s)", c.userID)

	// Start message handlers
	go c.readPump()
	go c.handleMessages()

	// Join default chat
	if err := c.JoinChat(config.ChatID1); err != nil {
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
			_, message1, err := c.conn.ReadMessage()
			if err != nil {
				c.errorChan <- fmt.Errorf("read error: %v", err)
				return
			}

			var msg message.Message
			if err := json.Unmarshal(message1, &msg); err != nil {
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
func (c *ChatClient) handleMessage(msg message.Message) {

	switch msg.MessageType {
	case "message":
		log.Printf("{ChatID1:%s} %s: %s", msg.ChatID, msg.UserID, msg.Caption)
	case "system":
		log.Printf("SYSTEM: %s", msg.Caption)
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
		//if chats, ok := msg.Caption.(map[string]interface{}); ok {
		//	log.Printf("Available chats: %v", chats)
		//}
	case "error":
		log.Printf("ERROR: ")
	default:
		log.Printf("Unknown message type: %s", msg.MessageType)
	}
}

// SendMessage sends a chat message to the current chat
func (c *ChatClient) SendMessage(content string) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message1 := message.Message{
		MessageType: "message",
		Caption:     content,
		ChatID:      c.currentChat,
		UserID:      c.userID,
		CreatedAt:   time.Now(),
	}

	return c.sendJSON(message1)
}

// JoinChat joins a specific chat
func (c *ChatClient) JoinChat(chatID uuid.UUID) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message1 := message.Message{
		MessageType: "join_chat",
		ChatID:      chatID,
		UserID:      c.userID,
	}

	if err := c.sendJSON(message1); err != nil {
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
func (c *ChatClient) LeaveChat(chatID uuid.UUID) error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message1 := message.Message{
		MessageType: "leave_chat",
		ChatID:      chatID,
		UserID:      c.userID,
	}

	if err := c.sendJSON(message1); err != nil {
		return err
	}

	c.mutex.Lock()
	delete(c.chats, chatID)
	if c.currentChat == chatID {
		c.currentChat = config.ChatID1
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

	//chatID := generateChatID(chatName)

	chatID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	message1 := message.Message{
		MessageType: "create_chat",
		ChatID:      chatID,
		Caption:     chatName,
		UserID:      c.userID,
	}

	return c.sendJSON(message1)
}

func (c *ChatClient) OpenChat(chatID uuid.UUID) error {

	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message1 := message.Message{
		MessageType: "chat_open",
		ChatID:      chatID,
		UserID:      c.userID,
	}

	return c.sendJSON(message1)
}

// ListChats requests the list of available chats
func (c *ChatClient) ListChats() error {
	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message1 := message.Message{
		MessageType: "get_chats",
	}

	return c.sendJSON(message1)
}

// SendTypingIndicator sends a typing indicator
func (c *ChatClient) SendTypingIndicator(typing bool) error {

	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message1 := message.Message{
		MessageType: "typing",
		ChatID:      c.currentChat,
		UserID:      c.userID,
		Caption:     fmt.Sprintf("%t", typing),
	}

	return c.sendJSON(message1)
}

func (c *ChatClient) SendSeenIndicator() error {

	if !c.isConnected {
		return errors.New("not connected to server")
	}

	message1 := message.Message{
		MessageType: "seen",
		ChatID:      c.currentChat,
		UserID:      c.userID,
	}

	return c.sendJSON(message1)
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
func (c *ChatClient) GetCurrentChat() uuid.UUID {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.currentChat
}

// GetChats returns the list of chats the chat_client is in
func (c *ChatClient) GetChats() []uuid.UUID {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	chats := make([]uuid.UUID, 0, len(c.chats))
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
//func generateChatID(name string) string {
//	// Simple implementation - in production, use a proper slug generator
//	return url.QueryEscape(name)
//}
