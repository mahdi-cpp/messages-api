package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mahdi-cpp/iris-tools/collection_manager_v3"
	"github.com/mahdi-cpp/iris-tools/image_loader"
	"github.com/mahdi-cpp/iris-tools/search"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
	"github.com/mahdi-cpp/messages-api/internal/config"
	"github.com/mahdi-cpp/messages-api/internal/hub"
	"github.com/mahdi-cpp/messages-api/internal/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024 * 10,
	WriteBufferSize: 1024 * 10,
}

type Manager struct {
	mu              sync.RWMutex
	usersStatus     map[string]*UserStatusData //key is userID
	chatsCollection *collection_manager_v3.Manager[*chat.Chat]
	chatManagers    map[string]*ChatManager // Maps chatIDs to their ChatManager
	hub             *hub.Hub
	iconLoader      *image_loader.ImageLoader
	ctx             context.Context
	// Added a channel to receive messages from the Hub for saving to a file.
	// یک کانال برای دریافت پیام‌ها از Hub جهت ذخیره در فایل اضافه شده است.
	messagesToSave chan *hub.Message
}

func (m *Manager) GetHub() *hub.Hub {
	return m.hub
}

func NewApplicationManager() (*Manager, error) {

	manager := &Manager{
		ctx:            context.Background(),
		messagesToSave: make(chan *hub.Message, 1000), // Initialize the channel
		chatManagers:   make(map[string]*ChatManager),
	}

	// Pass the new channel to the Hub
	// کانال جدید را به Hub پاس می‌دهیم.
	manager.hub = hub.NewHub(manager.messagesToSave)
	go manager.hub.Run()

	// Start a goroutine to listen for messages and save them to a file.
	// یک goroutine برای گوش دادن به پیام‌ها و ذخیره آن‌ها در فایل راه‌اندازی می‌کنیم.
	go manager.saveMessagesToFile()

	var err error
	manager.chatsCollection, err = collection_manager_v3.NewCollectionManager[*chat.Chat]("/app/iris/com.iris.messages/metadata", true)
	if err != nil {
		panic(err)
	}

	//all, err := manager.chatsCollection.GetAll()
	//if err != nil {
	//	return nil, err
	//}

	//for _, chat1 := range all {
	//	fmt.Println("chatId", chat1.ID)
	//}

	_, err = manager.GetChatManager(config.ChatID)
	if err != nil {
		fmt.Println(err)
	}

	return manager, nil
}

func (m *Manager) GetChatManager(chatID string) (*ChatManager, error) {

	chatManager, ok := m.chatManagers[chatID]
	if ok {
		return chatManager, nil
	}

	chat1, err := m.chatsCollection.Get(chatID)
	if err != nil {
		fmt.Println("chat not found in cash")
		return nil, err
	}

	chatManager, err = NewChatManager(chat1)
	if err != nil {
		fmt.Println(err)
	}

	m.chatManagers[chatID] = chatManager // add to cash

	return chatManager, nil
}

func (m *Manager) CreateWebsocketClient(w http.ResponseWriter, r *http.Request, userID string, username string) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	log.Printf("WebSocket connection established for user: %s (%s)", username, userID)

	client := hub.NewClient(m.hub, conn, userID)
	m.hub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump()

	// Send welcome message only to this chat_client
	welcomeMessage := map[string]interface{}{
		"type":    "system",
		"message": "Welcome to the chat!",
		"userId":  userID,
		"success": true,
	}

	if err := client.SendMessage(welcomeMessage); err != nil {
		log.Printf("Failed to send welcome message to user %s: %v", userID, err)
	}

	// Notify all users in default chat about new user
	m.notifyUserJoined(userID, username, "general")
}

// notifyUserJoined sends a notification when a user joins a chat
func (m *Manager) notifyUserJoined(userID, username, chatID string) {

	joinMessage := map[string]interface{}{
		"type":      "user_joined",
		"userId":    userID,
		"username":  username,
		"message":   username + " joined the chat",
		"chatID":    chatID,
		"timestamp": time.Now(),
	}

	// Broadcast to the specific chat
	m.hub.BroadcastToChat(chatID, joinMessage)
}

func (m *Manager) ChatCreate(newChat *chat.Chat) (*chat.Chat, error) {

	chatID, err := utils.GenerateUUID()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	newChat.ID = chatID

	_, err = m.chatsCollection.Create(newChat)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return newChat, nil
	//chatManager, err := NewChatManager(create)
	//if err != nil {
	//	panic(err)
	//}
}

func (m *Manager) UpdateChat(chatID string, updateOptions chat.UpdateOptions) error {

	fmt.Println(chatID)

	chat1, err := m.chatsCollection.Get(chatID)
	if err != nil {
		return err
	}

	chat.Update(chat1, updateOptions)

	_, err = m.chatsCollection.Update(chat1)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) MessageCreate(newMessage *message.Message) (*message.Message, error) {

	chatManager, ok := m.chatManagers[newMessage.ChatID]
	if !ok {
		fmt.Println("chat not found.")
		return nil, errors.New("chat not found")
	}

	id, err := utils.GenerateUUID()
	if err != nil {
		fmt.Printf("Failed to generate uuid: %v", err)
		return nil, err
	}

	newMessage.ID = id

	_, err = chatManager.messages.Create(newMessage)
	if err != nil {
		fmt.Println("Failed to create message to file.")
		return nil, err
	}

	return newMessage, nil
}

func (m *Manager) saveMessagesToFile() {

	for msg := range m.messagesToSave {

		chatManager, err := m.GetChatManager(msg.ChatID)
		if err != nil {
			return
		}

		id, err := utils.GenerateUUID()
		if err != nil {
			fmt.Printf("Failed to generate uuid: %v", err)
			return
		}

		newMessage := &message.Message{
			MessageType: "message",
			ID:          id,
			Width:       450,
			UserID:      msg.UserID,
			ChatID:      msg.ChatID,
			Content:     msg.Content,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Version:     "1",
		}

		err = chatManager.CreateMessage(newMessage)
		if err != nil {
			fmt.Println("Failed to create message to file.")
			return
		}

		m.hub.BroadcastToChat(msg.ChatID, newMessage)
	}
}

func (m *Manager) ReadChat(chatID string) (*chat.Chat, error) {

	chat1, err := m.chatsCollection.Get(chatID)
	if err != nil {
		return nil, err
	}

	return chat1, nil
}

func (m *Manager) ReadAllChats(chatOptions *chat.SearchOptions) ([]*chat.Chat, error) {

	chats, err := m.chatsCollection.GetAll()
	if err != nil {
		return nil, err
	}

	var userChats []*chat.Chat
	results := search.Find(chats, chat.HasMemberWith(chat.MemberWithUserID(config.UserId)))

	lessFn := chat.GetLessFunc("updatedAt", "start")
	if lessFn != nil {
		search.SortIndexedItems(results, lessFn)
	}

	//fmt.Println("ReadAllChats: ", len(results))

	for _, result := range results {
		userChats = append(userChats, result.Value)
	}

	filterChats := chat.Search(userChats, chatOptions)

	return filterChats, nil
}

func (m *Manager) ReadUserChats(userID string) ([]*chat.Chat, error) {

	chats, err := m.chatsCollection.GetAll()
	if err != nil {
		return nil, err
	}

	//searchOptions := &chat.SearchOptions{
	//	Page: 0,
	//	Size:  10,
	//}
	//filterChats := chat.Search(chatsCollection, searchOptions)

	var filterChats []*chat.Chat
	results := search.Find(chats, chat.HasMemberWith(chat.MemberWithUserID(userID)))

	lessFn := chat.GetLessFunc("updatedAt", "start")
	if lessFn != nil {
		search.SortIndexedItems(results, lessFn)
	}

	for _, result := range results {
		filterChats = append(filterChats, result.Value)
	}

	return filterChats, nil
}

// -------------------------------------------------------------------------------

func (m *Manager) ReadAllMessages(with *message.SearchOptions) ([]*message.Message, error) {

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
