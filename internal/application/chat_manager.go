package application

import (
	"context"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"sync"

	"github.com/mahdi-cpp/iris-tools/collection_manager_v3"
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
	messages          *collection_manager_v3.Manager[*message.Message]
	maintenanceCtx    context.Context
	cancelMaintenance context.CancelFunc
	statsMu           sync.Mutex
}

func NewChatManager(chat *chat.Chat) (*ChatManager, error) {

	manager := &ChatManager{
		chat: chat,
	}

	var err error
	manager.messages, err = collection_manager_v3.NewCollectionManager[*message.Message](filepath.Join(root, chat.ID, chatMessage), true)
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

func (m *ChatManager) ReadMessage(messageId string) ([]*message.Message, error) {

	all, err := m.messages.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error get messages")
	}

	return all, nil
}

func (m *ChatManager) ReadAllMessages() ([]*message.Message, error) {

	all, err := m.messages.GetAll()
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (m *ChatManager) UpdateMessage(updateMessage *message.Message) error {
	_, err := m.messages.Update(updateMessage)
	if err != nil {
		return err
	}
	return nil
}

func (m *ChatManager) DeleteMessage(messageID string) error {
	err := m.messages.Delete(messageID)
	if err != nil {
		return err
	}
	return nil
}
