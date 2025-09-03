package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mahdi-cpp/messages-api/internal/client"
)

func main() {

	// Configuration
	config := client.Config{
		ServerURL: "ws://localhost:8089/ws",
		UserID:    "user_130",
		Username:  "GoClientUser_5",
		Timeout:   30 * time.Second,
	}

	// Create client
	chatClient := client.NewChatClient(config)

	// Connect to server
	if err := chatClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer chatClient.Close()

	fmt.Println("=== Go Chat Client ===")
	fmt.Println("Commands:")
	fmt.Println("  /join <room>    - Join a room")
	fmt.Println("  /leave <room>   - Leave a room")
	fmt.Println("  /create <room>  - Create a room")
	fmt.Println("  /list           - List available rooms")
	fmt.Println("  /rooms          - Show joined rooms")
	fmt.Println("  /current        - Show current room")
	fmt.Println("  /exit           - Exit the client")
	fmt.Println("  /help           - Show this help")
	fmt.Println()
	fmt.Printf("Connected as: %s\n", config.Username)
	fmt.Printf("Current room: %s\n", chatClient.GetCurrentRoom())
	fmt.Println("Type your message and press Enter to send:")
	fmt.Println()

	// Start message input handler
	go handleUserInput(chatClient)

	// Wait for interrupt or connection close
	chatClient.WaitForInterrupt()
}

func handleUserInput(chatClient *client.ChatClient) {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		// Handle commands
		if strings.HasPrefix(text, "/") {
			handleCommand(chatClient, text)
			continue
		}

		// Send regular message
		if err := chatClient.SendMessage(text); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
}

func handleCommand(chatClient *client.ChatClient, command string) {

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "/join":
		if len(parts) < 2 {
			fmt.Println("Usage: /join <room>")
			return
		}
		if err := chatClient.JoinRoom(parts[1]); err != nil {
			log.Printf("Failed to join room: %v", err)
		}

	case "/leave":
		if len(parts) < 2 {
			fmt.Println("Usage: /leave <room>")
			return
		}
		if err := chatClient.LeaveRoom(parts[1]); err != nil {
			log.Printf("Failed to leave room: %v", err)
		}

	case "/create":
		if len(parts) < 2 {
			fmt.Println("Usage: /create <room_name>")
			return
		}
		roomName := strings.Join(parts[1:], " ")
		if err := chatClient.CreateRoom(roomName); err != nil {
			log.Printf("Failed to create room: %v", err)
		}

	case "/list":
		if err := chatClient.ListRooms(); err != nil {
			log.Printf("Failed to list rooms: %v", err)
		}

	case "/rooms":
		rooms := chatClient.GetRooms()
		fmt.Printf("Joined rooms: %v\n", rooms)

	case "/current":
		fmt.Printf("Current room: %s\n", chatClient.GetCurrentRoom())

	case "/exit":
		fmt.Println("Goodbye!")
		chatClient.Close()
		os.Exit(0)

	case "/help":
		fmt.Println("Commands:")
		fmt.Println("  /join <room>    - Join a room")
		fmt.Println("  /leave <room>   - Leave a room")
		fmt.Println("  /create <room>  - Create a room")
		fmt.Println("  /list           - List available rooms")
		fmt.Println("  /rooms          - Show joined rooms")
		fmt.Println("  /current        - Show current room")
		fmt.Println("  /exit           - Exit the client")
		fmt.Println("  /help           - Show this help")

	default:
		fmt.Printf("Unknown command: %s\n", parts[0])
		fmt.Println("Type /help for available commands")
	}
}
