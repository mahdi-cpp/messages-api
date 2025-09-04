package application

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mahdi-cpp/iris-tools/collection_manager_v3"
	"github.com/mahdi-cpp/iris-tools/image_loader"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/hub"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	mu           sync.RWMutex
	hub          *hub.Hub
	usersStatus  map[string]*UserStatusData //key is userID
	chats        *collection_manager_v3.Manager[*chat.Chat]
	chatManagers map[string]*ChatManager // Maps chatIDs to their ChatManager
	iconLoader   *image_loader.ImageLoader
	ctx          context.Context
}

func (m *Manager) GetHub() *hub.Hub {
	return m.hub
}

func NewApplicationManager() (*Manager, error) {

	manager := &Manager{
		ctx: context.Background(),
	}

	manager.hub = hub.NewHub()
	go manager.hub.Run()

	var err error
	manager.chats, err = collection_manager_v3.NewCollectionManager[*chat.Chat]("/app/iris/com.iris.message/chats/metadata", false)
	if err != nil {
		panic(err)
	}

	return manager, nil
}

func (m *Manager) loadChatContent(chatID string) {

	chatManager, err := NewChatManager(chatID)
	if err != nil {
		panic(err)
	}

	m.chatManagers[chatID] = chatManager
}

func (m *Manager) CreateWebsocketClient(w http.ResponseWriter, r *http.Request, userID string, username string) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	log.Printf("WebSocket connection established for user: %s (%s)", username, userID)

	client := hub.NewClient(m.hub, conn, userID)
	m.hub.RegisterClient(client)

	// Set user info in client (you'll need to add this field to Client struct)
	client.SetUserInfo(username)

	go client.WritePump()
	go client.ReadPump()

	// Send welcome message only to this client
	welcomeMessage := map[string]interface{}{
		"type":    "system",
		"message": "Welcome to the chat!",
		"userId":  userID,
		"success": true,
	}

	if err := client.SendMessage(welcomeMessage); err != nil {
		log.Printf("Failed to send welcome message to user %s: %v", userID, err)
	}

	// Notify all users in default chat about new user
	m.notifyUserJoined(userID, username, "general")
}

// notifyUserJoined sends a notification when a user joins a chat
func (m *Manager) notifyUserJoined(userID, username, chatID string) {

	joinMessage := map[string]interface{}{
		"type":      "user_joined",
		"userId":    userID,
		"username":  username,
		"message":   username + " joined the chat",
		"chatID":    chatID,
		"timestamp": time.Now(),
	}

	// Broadcast to the specific chat
	m.hub.BroadcastToChat(chatID, joinMessage)
}
