package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
)

const basePath = "/app/iris/com.iris.message/data"

func getChatPath(chatID string) string {
	return filepath.Join(basePath, fmt.Sprintf("chat_%s", chatID))
}

func getChatMetadataPath(chatID string) string {
	return filepath.Join(getChatPath(chatID), "metadata.json")
}

func getChatMessagesPath(chatID string) string {
	return filepath.Join(getChatPath(chatID), "message")
}

func getMessagePath(chatID, messageID string) string {
	return filepath.Join(getChatMessagesPath(chatID), fmt.Sprintf("%s.json", messageID))
}

func SaveChat(chat *chat.Chat) error {
	// Create chat directory if it doesn't exist
	chatPath := getChatPath(chat.ID)
	if err := os.MkdirAll(chatPath, 0755); err != nil {
		return err
	}

	// Create message directory if it doesn't exist
	messagesPath := getChatMessagesPath(chat.ID)
	if err := os.MkdirAll(messagesPath, 0755); err != nil {
		return err
	}

	// Save metadata
	metadataPath := getChatMetadataPath(chat.ID)
	data, err := json.MarshalIndent(chat, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(metadataPath, data, 0644)
}

func LoadChat(chatID string) (*chat.Chat, error) {
	metadataPath := getChatMetadataPath(chatID)
	data, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var chat1 chat.Chat
	if err := json.Unmarshal(data, &chat1); err != nil {
		return nil, err
	}

	return &chat1, nil
}

func LoadAllChats() ([]*chat.Chat, error) {
	var chats []*chat.Chat

	// Iterate through all chat directories
	entries, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) > 5 && entry.Name()[:5] == "chat_" {
			chatID := entry.Name()[5:]
			chat1, err := LoadChat(chatID)
			if err != nil {
				continue // Skip chats that can't be loaded
			}

			chats = append(chats, chat1)
		}
	}

	return chats, nil
}
