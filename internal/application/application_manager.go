package application

import (
	"context"
	"sync"

	"github.com/mahdi-cpp/iris-tools/image_loader"
)

type Manager struct {
	mu         sync.RWMutex
	chats      map[string]*ChatManager // Maps chatIDs to their ChatManager
	iconLoader *image_loader.ImageLoader
	ctx        context.Context
}

func NewApplicationManager() (*Manager, error) {

	// Handler the manager
	manager := &Manager{
		chats: make(map[string]*ChatManager),
		ctx:   context.Background(),
	}

	return manager, nil
}
