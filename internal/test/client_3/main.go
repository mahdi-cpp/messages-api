package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/client"
)

func generateUUID() (string, error) {
	u7, err2 := uuid.NewV7()
	if err2 != nil {
		return "", fmt.Errorf("error generating UUIDv7: %w", err2)
	}
	return u7.String(), nil
}

func main() {

	userID, err := generateUUID()
	if err != nil {
		return
	}

	config := client.Config{
		ServerURL: "ws://localhost:8089/ws",
		UserID:    "userID_" + userID,
		Timeout:   30 * time.Second,
	}

	chatClient, err := client.NewChatClient(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Connect to server
	if err := chatClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func(chatClient *client.ChatClient) {
		err := chatClient.Close()
		if err != nil {
		}
	}(chatClient)

	fmt.Println("=== Go Chat Client ===")
	fmt.Println("Commands:")
	fmt.Println("  /join <chat>    - Join a chat")
	fmt.Println("  /leave <chat>   - Leave a chat")
	fmt.Println("  /create <chat>  - Create a chat")
	fmt.Println("  /list           - List available chats")
	fmt.Println("  /chats          - Show joined chats")
	fmt.Println("  /current        - Show current chat")
	fmt.Println("  /typing         - Show is typing to other users")
	fmt.Println("  /exit           - Exit the client")
	fmt.Println("  /help           - Show this help")
	fmt.Println()
	fmt.Printf("Connected as: %s\n", config.UserID)
	fmt.Printf("Current chat: %s\n", chatClient.GetCurrentChat())
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
			fmt.Println("Usage: /join <chat>")
			return
		}
		if err := chatClient.JoinChat(parts[1]); err != nil {
			log.Printf("Failed to join chat: %v", err)
		}

	case "/leave":
		if len(parts) < 2 {
			fmt.Println("Usage: /leave <chat>")
			return
		}
		if err := chatClient.LeaveChat(parts[1]); err != nil {
			log.Printf("Failed to leave chat: %v", err)
		}

	case "/create":
		if len(parts) < 2 {
			fmt.Println("Usage: /create <chat_name>")
			return
		}
		chatName := strings.Join(parts[1:], " ")
		if err := chatClient.CreateChat(chatName); err != nil {
			log.Printf("Failed to create chat: %v", err)
		}

	case "/list":
		if err := chatClient.ListChats(); err != nil {
			log.Printf("Failed to list chats: %v", err)
		}

	case "/chats":
		chats := chatClient.GetChats()
		fmt.Printf("Joined chats: %v\n", chats)

	case "/typing":
		err := chatClient.SendTypingIndicator(true)
		if err != nil {
			log.Printf("Failed to typing chat: %v", err)
		} else {
			fmt.Printf("sent typing command:\n")
		}

	case "/current":
		fmt.Printf("Current chat: %s\n", chatClient.GetCurrentChat())

	case "/exit":
		fmt.Println("Goodbye!")
		err := chatClient.Close()
		if err != nil {
			return
		}
		os.Exit(0)

	case "/help":
		fmt.Println("Commands:")
		fmt.Println("  /join <chat>    - Join a chat")
		fmt.Println("  /leave <chat>   - Leave a chat")
		fmt.Println("  /create <chat>  - Create a chat")
		fmt.Println("  /list           - List available chats")
		fmt.Println("  /chats          - Show joined chats")
		fmt.Println("  /current        - Show current chat")
		fmt.Println("  /exit           - Exit the client")
		fmt.Println("  /help           - Show this help")

	default:
		fmt.Printf("Unknown command: %s\n", parts[0])
		fmt.Println("Type /help for available commands")
	}
}
