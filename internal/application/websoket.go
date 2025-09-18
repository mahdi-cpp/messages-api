package application

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
	"github.com/mahdi-cpp/messages-api/internal/hub"
)

func (m *AppManager) CreateWebsocketClient(w http.ResponseWriter, r *http.Request, userID uuid.UUID, username string) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	log.Printf("WebSocket connection established for user: %s (%s)", username, userID)

	client := hub.NewClient(m.hub, conn, userID)
	m.hub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump()

	// Send welcome message only to this chat_client
	welcomeMessage := map[string]interface{}{
		"type":    "system",
		"message": "Welcome to the chat!",
		"userId":  userID,
		"success": true,
	}

	if err := client.SendMessage(welcomeMessage); err != nil {
		log.Printf("Failed to send welcome message to user %s: %v", userID, err)
	}

	chatId, err := uuid.NewV7()
	if err != nil {
		return
	}

	// Notify all users in default chat about new user
	m.notifyUserJoined(userID, chatId, username)
}

// notifyUserJoined sends a notification when a user joins a chat
func (m *AppManager) notifyUserJoined(userID, chatID uuid.UUID, username string) {

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

func (m *AppManager) saveMessagesToFile() {
	for msg := range m.messagesToSave {

		chatManager, err := m.GetChatManager(msg.ChatID)
		if err != nil {
			return
		}

		id, err := helpers.GenerateUUID()
		if err != nil {
			fmt.Printf("Failed to generate uuid: %v", err)
			return
		}

		newMessage := &message.Message{
			//MessageType: "message",
			ID:        id,
			UserID:    msg.UserID,
			ChatID:    msg.ChatID,
			Caption:   msg.Content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   "1",
		}

		err = chatManager.CreateMessage(newMessage)
		if err != nil {
			fmt.Println("Failed to create message to file.")
			return
		}

		m.hub.BroadcastToChat(msg.ChatID, newMessage)
	}
}
