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
	err := manager.Open()
	if err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *ChatManager) Open() error {

	var root = "/app/iris/com.iris.messages/chats"
	var chatMessage = "/metadata/v1/messages"

	var err error
	m.messages, err = collection_manager_v3.NewCollectionManager[*message.Message](filepath.Join(root, m.chat.ID, chatMessage), true)
	if err != nil {
		panic(err)
	}
	return nil
}

func (m *ChatManager) Read() error {

	var root = "/com.iris.ali/chats"
	var chatMessage = "/metadata/v1/messages"

	var err error
	m.messages, err = collection_manager_v3.NewCollectionManager[*message.Message](filepath.Join(root, m.chat.ID, chatMessage), false)
	if err != nil {
		panic(err)
	}
	return nil
}

func (m *ChatManager) Update(updateOptions chat.UpdateOptions) error {
	return nil
}

func (m *ChatManager) Delete() error {
	return nil
}

//Message handlers----------------------------

func (m *ChatManager) CreateMessage(addMessage *message.Message) error {
	_, err := m.messages.Create(addMessage)
	if err != nil {
		return err
	}
	return nil
}

func (m *ChatManager) ReadMessages() ([]*message.Message, error) {
	if m.messages == nil {
		return nil, fmt.Errorf("no messages")
	}

	all, err := m.messages.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error get messages")
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
