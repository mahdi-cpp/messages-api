package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mahdi-cpp/messages-api/internal/chat_client"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

func main() {

	clientConfig := chat_client.ClientChatConfig{
		ServerURL: "ws://localhost:50151/ws",
		UserID:    config.Mahdi,
		Timeout:   30 * time.Second,
	}

	chatClient, err := chat_client.NewChatClient(clientConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Connect to server
	if err := chatClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func(chatClient *chat_client.ChatClient) {
		err := chatClient.Close()
		if err != nil {
		}
	}(chatClient)

	fmt.Println("=== Go Chat Client ===")
	fmt.Println("Commands:")
	fmt.Println("  /join <chat>    - Join a chat")
	fmt.Println("  /leave <chat>   - Leave a chat")
	fmt.Println("  /create <chat>  - create a chat")
	fmt.Println("  /open <chat>    - Open a chat")
	fmt.Println("  /list           - List available chats")
	fmt.Println("  /chats          - Show joined chats")
	fmt.Println("  /current        - Show current chat")
	fmt.Println("  /typing         - Show is typing to other users")
	fmt.Println("  /seen           - Show is seen message")
	fmt.Println("  /exit           - Exit the chat_client")
	fmt.Println("  /help           - Show this help")
	fmt.Println()
	fmt.Printf("userID: %s connected\n", clientConfig.UserID)
	fmt.Printf("Current chat: %s\n", chatClient.GetCurrentChat())
	fmt.Println("Type your message and press Enter to send:")
	fmt.Println()

	// Start message input handler
	go handleUserInput(chatClient)

	// Wait for interrupt or connection close
	chatClient.WaitForInterrupt()
}

func handleUserInput(chatClient *chat_client.ChatClient) {

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

func handleCommand(chatClient *chat_client.ChatClient, command string) {

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

	case "/open":
		if len(parts) < 2 {
			fmt.Println("Usage: /open <chatID>")
			return
		}
		chatID := strings.Join(parts[1:], " ")
		if err := chatClient.CreateChat(chatID); err != nil {
			log.Printf("Failed to open chat: %v", err)
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

	case "/seen":
		err := chatClient.SendSeenIndicator()
		if err != nil {
			log.Printf("Failed to seen chat: %v", err)
		} else {
			fmt.Printf("sent seen command:\n")
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
		fmt.Println("  /create <chat>  - create a chat")
		fmt.Println("  /list           - List available chats")
		fmt.Println("  /chats          - Show joined chats")
		fmt.Println("  /current        - Show current chat")
		fmt.Println("  /exit           - Exit the chat_client")
		fmt.Println("  /help           - Show this help")

	default:
		fmt.Printf("Unknown command: %s\n", parts[0])
		fmt.Println("Type /help for available commands")
	}
}
