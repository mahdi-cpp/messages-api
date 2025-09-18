package application

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mahdi-cpp/iris-tools/image_loader"
	"github.com/mahdi-cpp/messages-api/internal/chat_manager"
	"github.com/mahdi-cpp/messages-api/internal/collection_manager"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
	"github.com/mahdi-cpp/messages-api/internal/config"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
	"github.com/mahdi-cpp/messages-api/internal/hub"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024 * 10,
	WriteBufferSize: 1024 * 10,
}

type AppManager struct {
	mu                    sync.RWMutex
	usersStatus           map[string]*UserStatusData //key is userID
	ChatCollectionManager *collection_manager.Manager[*chat.Chat]
	chatManagers          map[uuid.UUID]*chat_manager.Manager // Maps chatIDs to their Manager
	hub                   *hub.Hub
	iconLoader            *image_loader.ImageLoader
	// Added a channel to receive messages from the Hub for saving to a file.
	// یک کانال برای دریافت پیام‌ها از Hub جهت ذخیره در فایل اضافه شده است.
	messagesToSave chan *hub.Message

	createChat chan *chat.Chat
}

func (m *AppManager) GetHub() *hub.Hub {
	return m.hub
}

func NewApplicationManager() (*AppManager, error) {

	manager := &AppManager{
		messagesToSave: make(chan *hub.Message, 1000), // Initialize the channel
		chatManagers:   make(map[uuid.UUID]*chat_manager.Manager),
		createChat:     make(chan *chat.Chat, 100),
	}

	// Pass the new channel to the Hub
	// کانال جدید را به Hub پاس می‌دهیم.
	//manager.hub = hub.NewHub(manager.messagesToSave)
	//go manager.hub.Run()

	// Start goroutine to listen for messages and save them to chats file.
	// یک goroutine برای گوش دادن به پیام‌ها و ذخیره آن‌ها در فایل راه‌اندازی می‌کنیم.
	//go manager.saveMessagesToFile()

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	fmt.Printf("Initial Alloc: %d B\n", m1.Alloc)

	var err error
	var chatsDirectory = config.GetPath("test/chats")
	manager.ChatCollectionManager, err = collection_manager.New[*chat.Chat](chatsDirectory)
	if err != nil {
		panic(err)
	}

	// Get final memory stats
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	fmt.Printf("Final Alloc: %d B\n", m2.Alloc)
	fmt.Printf("Map increased memory usage by: %d B\n", m2.Alloc-m1.Alloc)

	_, err = manager.GetChatManager(config.ChatID1)
	if err != nil {
		fmt.Println(err)
	}

	return manager, nil
}

func (m *AppManager) GetChatManager(chatID uuid.UUID) (*chat_manager.Manager, error) {

	chatManager, ok := m.chatManagers[chatID]
	if ok {
		return chatManager, nil
	}

	chat1, err := m.ChatCollectionManager.Read(chatID)
	if err != nil {
		fmt.Println("chat not found in cash")
		return nil, err
	}

	chatManager, err = chat_manager.New(chat1)
	if err != nil {
		fmt.Println(err)
	}

	m.chatManagers[chatID] = chatManager // add to cash

	return chatManager, nil
}

func (m *AppManager) createChatTimout() {
	//for createChat := range m.createChat {
	//
	//}
}

// ---

func (m *AppManager) MessageCreate(newMessage *message.Message) (*message.Message, error) {

	chatManager, ok := m.chatManagers[newMessage.ChatID]
	if !ok {
		fmt.Println("chat not found.")
		return nil, errors.New("chat not found")
	}

	id, err := helpers.GenerateUUID()
	if err != nil {
		fmt.Printf("Failed to generate uuid: %v", err)
		return nil, err
	}

	newMessage.ID = id

	_, err = chatManager.Messages.Create(newMessage)
	if err != nil {
		fmt.Println("Failed to create message to file.")
		return nil, err
	}

	return newMessage, nil
}

func (m *AppManager) ReadAllMessages(with *message.SearchOptions) ([]*message.Message, error) {

	chatManager, ok := m.chatManagers[with.ChatID]
	if !ok {
		return nil, fmt.Errorf("chatId not found")
	}

	all, err := chatManager.ReadAllMessages()
	if err != nil {
		return nil, err
	}

	return message.Search(all, with), nil
}
