package application

import (
	"context"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"sync"

	"github.com/mahdi-cpp/iris-tools/collection_manager_v3"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
)

type ChatManager struct {
	mu                sync.RWMutex // Protects all indexes and maps
	messages          *collection_manager_v3.Manager[*message.Message]
	maintenanceCtx    context.Context
	cancelMaintenance context.CancelFunc
	statsMu           sync.Mutex
}

func NewChatManager(chatID string) (*ChatManager, error) {
	manager := &ChatManager{}

	var root = "/app/iris/com.iris.message/chats"
	var chatMessage = "/metadata/v1/message"

	var err error
	manager.messages, err = collection_manager_v3.NewCollectionManager[*message.Message](filepath.Join(root, chatID, chatMessage), false)
	if err != nil {
		panic(err)
	}

	return manager, nil
}

func (m *ChatManager) fetchMessages(userID string) ([]*message.Message, error) {

	return nil, nil
}
