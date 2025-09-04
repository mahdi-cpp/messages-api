package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// HandleClientMessage processes a message from a client
func (h *Hub) HandleClientMessage(client *Client, rawMessage []byte) {

	var message struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		ChatID  string `json:"chatId"`
	}

	if err := json.Unmarshal(rawMessage, &message); err != nil {
		log.Printf("Error parsing message from client %s: %v", client.userID, err)
		return
	}

	switch message.Type {
	case "message":
		h.HandleChatMessage(client, message.Content, message.ChatID)
	case "typing":
		h.HandleTypingIndicator(client, message.Content, message.ChatID)
	case "seen":
		h.HandleSeenIndicator(client, message.Content, message.ChatID)
	case "join_chat":
		h.HandleJoinChat(client, message.ChatID)
	case "leave_chat":
		h.HandleLeaveChat(client, message.ChatID)
	case "create_chat":
		h.HandleCreateChat(client, message.Content)
	case "get_chats":
		h.HandleGetChats(client)
	default:
		log.Printf("Unknown message type from client %s: %s", client.userID, message.Type)
	}
}

// HandleChatMessage processes and broadcasts chat message
func (h *Hub) HandleChatMessage(client *Client, content, chatID string) {

	if content == "" {
		return
	}

	messageID, err := generateUUID()
	if err != nil {
		fmt.Printf("Error generating message ID: %v", err)
		return
	}

	// Create the message to broadcast
	chatMessage := map[string]interface{}{
		"type":      "message",
		"id":        messageID,
		"userId":    client.userID,
		"content":   content,
		"chatID":    chatID,
		"timestamp": time.Now(),
	}

	log.Printf("Broadcasting message from %s in chat %s:", chatID, content)

	// Broadcast to all clients in the chat
	h.BroadcastToChat(chatID, chatMessage)
}

// HandleTypingIndicator broadcasts typing status
func (h *Hub) HandleTypingIndicator(client *Client, typing, chatID string) {
	typingMessage := map[string]interface{}{
		"type":      "typing",
		"userId":    client.userID,
		"chatId":    chatID,
		"typing":    typing == "true",
		"timestamp": time.Now(),
	}

	h.BroadcastToChat(chatID, typingMessage)
}

func (h *Hub) HandleSeenIndicator(client *Client, typing, chatID string) {
	typingMessage := map[string]interface{}{
		"type":      "seen",
		"userId":    client.userID,
		"chatId":    chatID,
		"timestamp": time.Now(),
	}

	h.BroadcastToChat(chatID, typingMessage)
}

// HandleJoinChat handles chat joining
func (h *Hub) HandleJoinChat(client *Client, chatID string) {

	h.JoinChat(chatID, client.userID, client)

	// Notify chat about new user
	joinMessage := map[string]interface{}{
		"type":      "user_joined",
		"userId":    client.userID,
		"message":   client.userID + " joined the chat",
		"chatId":    chatID,
		"timestamp": time.Now(),
	}

	h.BroadcastToChat(chatID, joinMessage)
}

// HandleLeaveChat handles chat leaving
func (h *Hub) HandleLeaveChat(client *Client, chatID string) {
	h.LeaveChat(chatID, client.userID)

	leaveMessage := map[string]interface{}{
		"type":      "user_left",
		"userId":    client.userID,
		"message":   client.userID + " left the chat",
		"chatId":    chatID,
		"timestamp": time.Now(),
	}

	h.BroadcastToChat(chatID, leaveMessage)
}

// HandleCreateChat handles chat creation
func (h *Hub) HandleCreateChat(client *Client, chatName string) {

	chatID, err := generateUUID()
	if err != nil {
		fmt.Printf("Error generating chat id: %v", err)
		return
	}

	h.CreateChat(chatID, chatName)
	h.JoinChat(chatID, client.userID, client)

	// Notify about chat creation
	chatMessage := map[string]interface{}{
		"type":      "chat_created",
		"chatId":    chatID,
		"chatName":  chatName,
		"userId":    client.userID,
		"timestamp": time.Now(),
	}

	h.BroadcastToAll(chatMessage)
}

// HandleGetChats handles chat list requests
func (h *Hub) HandleGetChats(client *Client) {
	chatList := h.GetChatList()

	response := map[string]interface{}{
		"type":     "chat_list",
		"chatList": chatList,
	}

	if err := client.SendMessage(response); err != nil {
		log.Printf("Failed to send chat list to client %s: %v", client.userID, err)
	}
}
