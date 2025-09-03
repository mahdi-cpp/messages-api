package main

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type      string    `json:"type"`
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	ChatID    string    `json:"chatId"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {

	//u := url.URL{Scheme: "ws", Host: "localhost:8089", Path: "/ws"}

	// Define your user information
	userID := "userid_12"
	username := "mahdi.cpp"

	// Construct the URL with query parameters
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8089", // Adjust if your server runs on a different host/port
		Path:   "/ws",
	}

	log.Printf("Connecting to %s", u.String())

	// Add query parameters for user_id and username
	q := u.Query()
	q.Set("user_id", userID)
	q.Set("username", username)
	u.RawQuery = q.Encode()

	log.Printf("Connecting to %s", u.String())

	// Connect to the WebSocket server.
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	sendMessage(c, "user_12", "username1", "Hello", "room1")

	// Wait for and read the response from the server.
	_, response, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("Received from server: %s", response)

	// Keep the client running for a moment to see the output.
	time.Sleep(10 * time.Second)
}

// Assume `c` is a *websocket.Conn connected to the server.
// Assume `currentUserID`, `currentUsername`, `content`, and `currentRoom` are defined.

func sendMessage(c *websocket.Conn, currentUserID, currentUsername, content, currentRoom string) {

	// 1. Create the Go struct with the message data.
	// You'll need to install the UUID package: `go get github.com/google/uuid`
	msg := Message{
		Type:      "message",
		ID:        uuid.New().String(),
		UserID:    currentUserID,
		Username:  currentUsername,
		Content:   content,
		ChatID:    currentRoom,
		Timestamp: time.Now(),
	}

	// 2. Marshal the struct into a JSON byte slice.
	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	// 3. Send the JSON message to the WebSocket server.
	err = c.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		log.Printf("Error writing message to WebSocket: %v", err)
		return
	}

	log.Printf("Message sent: %s", jsonMessage)
}
