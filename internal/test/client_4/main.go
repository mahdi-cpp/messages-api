// test_broadcast.go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mahdi-cpp/messages-api/internal/client"
)

func main() {
	// Test with two clients in the same chat
	client1 := client.NewChatClient(client.Config{
		ServerURL: "ws://localhost:8089/ws",
		UserID:    "test_user_1",
		Username:  "TestUser1",
	})

	client2 := client.NewChatClient(client.Config{
		ServerURL: "ws://localhost:8089/ws",
		UserID:    "test_user_2",
		Username:  "TestUser2",
	})

	// Connect both clients
	if err := client1.Connect(); err != nil {
		log.Fatal("Client1 connection failed:", err)
	}
	defer client1.Close()

	if err := client2.Connect(); err != nil {
		log.Fatal("Client2 connection failed:", err)
	}
	defer client2.Close()

	// Both join the same chat
	if err := client1.JoinChat("family_group"); err != nil {
		log.Fatal("Client1 join failed:", err)
	}

	if err := client2.JoinChat("family_group"); err != nil {
		log.Fatal("Client2 join failed:", err)
	}

	// Wait a bit for connections to stabilize
	time.Sleep(1 * time.Second)

	// Client1 sends a message
	fmt.Println("Client1 sending message...")
	if err := client1.SendMessage("Hello from Client1!"); err != nil {
		log.Fatal("Client1 send failed:", err)
	}

	// Wait for message to be delivered
	time.Sleep(2 * time.Second)

	fmt.Printf("Client1 received %d message\n", client1.GetMessageCount())
	fmt.Printf("Client2 received %d message\n", client2.GetMessageCount())

	// Client2 should have received the message from Client1
	if client2.GetMessageCount() > 0 {
		fmt.Println("SUCCESS: Broadcasting is working!")
	} else {
		fmt.Println("FAIL: Client2 did not receive the message")
	}
}
