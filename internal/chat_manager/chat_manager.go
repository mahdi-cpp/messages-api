package chat_manager

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/collection_manager"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
)

const (
	root        = "/app/iris/com.iris.messages/chats"
	chatMessage = "/metadata/v1/messages"
)

type Manager struct {
	chat     *chat.Chat
	Messages *collection_manager.Manager[*message.Message]
}

func New(chat *chat.Chat) (*Manager, error) {

	manager := &Manager{
		chat: chat,
	}

	var err error
	// In this layer, we still initialize the collection manager without a context.
	// The timeout for initial loading can be handled internally by the New function itself.
	var messagesDir = filepath.Join(root, chat.ID.String(), chatMessage)
	manager.Messages, err = collection_manager.New[*message.Message](messagesDir)
	if err != nil {
		return nil, fmt.Errorf("error initializing chat message manager: %w", err)
	}

	return manager, nil
}

// CreateMessage adds a new message to the chat. No context is passed here.
func (m *Manager) CreateMessage(addMessage *message.Message) error {
	_, err := m.Messages.Create(addMessage)
	if err != nil {
		return err
	}
	return nil
}

// ReadMessage retrieves a message by its ID.
func (m *Manager) ReadMessage(messageId uuid.UUID) (*message.Message, error) {
	selectMessage, err := m.Messages.Read(messageId)
	if err != nil {
		return nil, fmt.Errorf("error reading message %s: %w", messageId, err)
	}
	return selectMessage, nil
}

// ReadAllMessages retrieves all messages in the chat.
func (m *Manager) ReadAllMessages() ([]*message.Message, error) {
	all, err := m.Messages.ReadAll()
	if err != nil {
		return nil, err
	}
	return all, nil
}

// UpdateMessage updates a message.
func (m *Manager) UpdateMessage(updateOptions message.UpdateOptions) (*message.Message, error) {
	msg, err := m.Messages.Read(updateOptions.MessageID)
	if err != nil {
		return nil, err
	}
	message.Update(msg, updateOptions)
	msg, err = m.Messages.Update(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// DeleteMessage deletes a message.
func (m *Manager) DeleteMessage(messageID uuid.UUID) error {
	err := m.Messages.Delete(messageID)
	if err != nil {
		return err
	}
	return nil
}
