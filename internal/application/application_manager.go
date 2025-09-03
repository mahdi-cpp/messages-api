package application

import (
	"context"
	"sync"

	"github.com/mahdi-cpp/iris-tools/collection_manager_v3"
	"github.com/mahdi-cpp/iris-tools/image_loader"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
)

type Manager struct {
	mu           sync.RWMutex
	usersStatus  map[string]*UserStatusData //key is userID
	chats        *collection_manager_v3.Manager[*chat.Chat]
	chatManagers map[string]*ChatManager // Maps chatIDs to their ChatManager
	iconLoader   *image_loader.ImageLoader
	ctx          context.Context
}

func NewApplicationManager() (*Manager, error) {

	manager := &Manager{
		ctx: context.Background(),
	}

	var err error
	manager.chats, err = collection_manager_v3.NewCollectionManager[*chat.Chat]("/app/iris/com.iris.messages/chats/metadata", false)
	if err != nil {
		panic(err)
	}

	return manager, nil
}

func (m *Manager) loadChatContent(chatID string) {

	chatManager, err := NewChatManager(chatID)
	if err != nil {
		panic(err)
	}

	m.chatManagers[chatID] = chatManager
}
