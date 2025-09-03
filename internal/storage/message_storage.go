package storage

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mahdi-cpp/messages-api/internal/collections/message"
)

func SaveMessage(message *message.Message) error {
	messagePath := getMessagePath(message.ChatID, message.ID)
	data, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(messagePath, data, 0644)
}

func LoadMessage(chatID, messageID string) (*message.Message, error) {
	messagePath := getMessagePath(chatID, messageID)
	data, err := ioutil.ReadFile(messagePath)
	if err != nil {
		return nil, err
	}

	var message1 message.Message
	if err := json.Unmarshal(data, &message1); err != nil {
		return nil, err
	}

	return &message1, nil
}

func LoadChatMessages(chatID string) ([]*message.Message, error) {
	messagesPath := getChatMessagesPath(chatID)
	entries, err := ioutil.ReadDir(messagesPath)
	if err != nil {
		return nil, err
	}

	var messages []*message.Message
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > 5 && entry.Name()[len(entry.Name())-5:] == ".json" {
			messageID := entry.Name()[:len(entry.Name())-5]
			message1, err := LoadMessage(chatID, messageID)
			if err != nil {
				continue // Skip messages that can't be loaded
			}

			messages = append(messages, message1)
		}
	}

	return messages, nil
}
