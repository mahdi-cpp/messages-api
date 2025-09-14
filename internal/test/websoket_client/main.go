package main

import (
	"github.com/goccy/go-json"
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

// Define your user information
var userID = "userid_12"
var username = "mahdi.cpp"

func main() {

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

	// Use a 'done' channel to signal when the program should stop.
	done := make(chan struct{})

	// Start a goroutine to continuously read message.
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("Received: %s", message)
		}
	}()
	sendJoinChat(c, "friends_chat")
	sendJoinChat(c, "family_chat")

	select {}

	//// create a channel to listen for interrupt signals (like Ctrl+C).
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)
	//
	//// Main loop to keep the program running until an interrupt signal is received.
	//for {
	//	select {
	//	case <-done:
	//		return
	//	case <-interrupt:
	//		log.Println("interrupt")
	//		// Cleanly close the connection by sending a close message to the server.
	//		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	//		if err != nil {
	//			log.Println("write close:", err)
	//			return
	//		}
	//		select {
	//		case <-done:
	//		case <-time.After(time.Second):
	//		}
	//		return
	//	}
	//}

}

type JoinChat struct {
	Type   string `json:"type"`
	ChatID string `json:"chatId"`
}

func sendJoinChat(c *websocket.Conn, chatID string) {

	// 1. create the Go struct with the message data.
	// You'll need to install the UUID package: `go get github.com/google/uuid`
	msg := JoinChat{
		Type:   "join_chat",
		ChatID: chatID,
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

func sendMessage(c *websocket.Conn, currentUserID, currentUsername, content, currentChat string) {

	// 1. create the Go struct with the message data.
	// You'll need to install the UUID package: `go get github.com/google/uuid`
	msg := Message{
		Type:      "message",
		ID:        uuid.New().String(),
		UserID:    currentUserID,
		Username:  currentUsername,
		Content:   content,
		ChatID:    currentChat,
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
