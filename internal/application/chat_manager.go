package application

import (
	"context"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/collection_manager_gemini"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
)

const (
	root        = "/app/iris/com.iris.messages/chats"
	chatMessage = "/metadata/v1/messages"
)

type ChatManager struct {
	mu                sync.RWMutex
	chat              *chat.Chat
	messages          *collection_manager_gemini.Manager[*message.Message]
	maintenanceCtx    context.Context
	cancelMaintenance context.CancelFunc
	statsMu           sync.Mutex
}

func NewChatManager(chat *chat.Chat) (*ChatManager, error) {

	manager := &ChatManager{
		chat: chat,
	}

	var err error
	manager.messages, err = collection_manager_gemini.New[*message.Message](filepath.Join(root, chat.ID.String(), chatMessage))
	if err != nil {
		fmt.Println("Error opening chat manager:", err)
	}

	return manager, nil
}

func (m *ChatManager) CreateMessage(addMessage *message.Message) error {
	_, err := m.messages.Create(addMessage)
	if err != nil {
		return err
	}
	return nil
}

func (m *ChatManager) ReadMessage(messageId uuid.UUID) (*message.Message, error) {

	selectMessage, err := m.messages.Read(messageId)
	if err != nil {
		return nil, fmt.Errorf("error get messages")
	}

	return selectMessage, nil
}

func (m *ChatManager) ReadAllMessages() ([]*message.Message, error) {

	all, err := m.messages.ReadAll()
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (m *ChatManager) UpdateMessage(updateOptions message.UpdateOptions) (*message.Message, error) {

	msg, err := m.messages.Read(updateOptions.MessageID)
	if err != nil {
		return nil, err
	}

	message.Update(msg, updateOptions)
	msg, err = m.messages.Update(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *ChatManager) DeleteMessage(messageID uuid.UUID) error {
	err := m.messages.Delete(messageID)
	if err != nil {
		return err
	}
	return nil
}

//func (m *ChatManager) DeleteAllMessages() error {
//
//}
