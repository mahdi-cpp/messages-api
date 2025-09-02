package main

import (
	"fmt"
	"time"

	"github.com/mahdi-cpp/iris-tools/collection_manager_v3"
	"github.com/mahdi-cpp/iris-tools/search"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
)

// Assuming you have your Chat and Member types defined here or in another package

func main() {

	collection, err := collection_manager_v3.NewCollectionManager[*chat.Chat]("/app/github.com/mahdi.cpp/messages-api/internal/documents/chats.json", true)
	if err != nil {
		panic(err)
	}

	chats, err := collection.GetAll()

	// Find for chats that have a specific user as a member
	userID := "018f3a8b-1b32-7290-b1d5-92716a445330"
	results := search.Find(chats, chat.HasMemberWith(chat.MemberWithUserID(userID)))

	// Process results
	for _, result := range results {
		fmt.Printf("Chat ID: %s, Title: %s\n", result.Value.ID, result.Value.Title)
	}

	// Find for chats with active admins who joined in the last month
	t := time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC)
	activeAdminResults := search.Find(chats, chat.HasMemberWith(chat.ActiveAdminsJoinedAfter(t)))
	fmt.Println("\nActiveAdminsJoinedAfter----------------------------------")
	for _, result := range activeAdminResults {
		fmt.Printf("Chat ID: %s, Title: %s\n", result.Value.ID, result.Value.Title)
	}

	// Find for chats with members having "manager" in their title
	managerResults := search.Find(chats, chat.HasMemberWith(chat.MemberWithCustomTitle("Mahdi Abdolmaleki")))
	fmt.Println("\nMemberWithCustomTitle----------------------------------")
	for _, result := range managerResults {
		fmt.Printf("Chat ID: %s, Title: %s\n", result.Value.ID, result.Value.Title)
	}
}
