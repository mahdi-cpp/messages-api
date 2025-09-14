package application

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/iris-tools/search"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/config"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
)

func (m *AppManager) ChatCreate(requestChat *chat.Chat) (*chat.Chat, error) {

	//err := requestChat.Validate()
	//if err != nil {
	//	return nil, err
	//}

	// Step 2: Generate a unique ID for the new chat
	chatID, err := helpers.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate chat ID: %w", err)
	}
	requestChat.ID = chatID

	// Step 3: create the chat in the database
	_, err = m.ChatCollectionManager.Create(requestChat)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat in database: %w", err)
	}

	return requestChat, nil
}

func (m *AppManager) ReadChat(chatID uuid.UUID) (*chat.Chat, error) {

	chat1, err := m.ChatCollectionManager.Read(chatID)
	if err != nil {
		return nil, err
	}

	return chat1, nil
}

func (m *AppManager) ReadAllChats(chatOptions *chat.SearchOptions) ([]*chat.Chat, error) {

	chats, err := m.ChatCollectionManager.ReadAll()
	if err != nil {
		return nil, err
	}

	var userChats []*chat.Chat
	results := search.Find(chats, chat.HasMemberWith(chat.MemberWithUserID(config.Mahdi)))

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

func (m *AppManager) ReadUserChats(userID uuid.UUID) ([]*chat.Chat, error) {

	chats, err := m.ChatCollectionManager.ReadAll()
	if err != nil {
		return nil, err
	}

	//searchOptions := &chat.SearchOptions{
	//	Page: 0,
	//	Size:  10,
	//}
	//filterChats := chat.Search(ChatCollectionManager, searchOptions)

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

func (m *AppManager) UpdateChats(updateOptions chat.UpdateOptions) error {

	for _, chatID := range updateOptions.ChatIDs {
		chat1, err := m.ChatCollectionManager.Read(chatID)
		if err != nil {
			return fmt.Errorf("failed to read chat %s: %w", chatID, err)
		}
		
		chat.Update(chat1, updateOptions)

		_, err = m.ChatCollectionManager.Update(chat1)
		if err != nil {
			return fmt.Errorf("failed to update chat %s: %w", chatID, err)
		}
	}
	return nil
}

func (m *AppManager) ChatDelete(chatID uuid.UUID) error {

	err := m.ChatCollectionManager.Delete(chatID)
	if err != nil {
		fmt.Println("error deleting chat")
		return err
	}

	delete(m.chatManagers, chatID)
	return nil
}
