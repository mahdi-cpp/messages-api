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

	collection, err := collection_manager_v3.NewCollectionManager[*chat.Chat]("/app/github.com/mahdi.cpp/message-api/internal/documents/chats.json", true)
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

	// create a sample chat with members
	chat1 := chat.Chat{
		Members: []chat.Member{
			{
				UserID:     "user1",
				Role:       "admin",
				JoinedAt:   time.Now().AddDate(0, -2, 0), // 2 months ago
				IsActive:   true,
				LastActive: time.Now().AddDate(0, 0, -1), // 1 day ago
			},
			{
				UserID:     "user2",
				Role:       "member",
				JoinedAt:   time.Now().AddDate(0, -1, 0), // 1 month ago
				IsActive:   false,
				LastActive: time.Now().AddDate(0, 0, -10), // 10 days ago
			},
			{
				UserID:     "user3",
				Role:       "creator",
				JoinedAt:   time.Now().AddDate(0, -3, 0), // 3 months ago
				IsActive:   true,
				LastActive: time.Now(), // now
			},
			{
				UserID:     "user4",
				Role:       "admin",
				JoinedAt:   time.Now().AddDate(0, -1, 5), // 1 month 5 days ago
				IsActive:   true,
				LastActive: time.Now().AddDate(0, 0, -2), // 2 days ago
			},
			{
				UserID:     "user5",
				Role:       "admin",
				JoinedAt:   time.Now().AddDate(0, -1, 5), // 1 month 5 days ago
				IsActive:   true,
				LastActive: time.Now().AddDate(0, 0, -2), // 2 days ago
			},
			{
				UserID:     "user6",
				Role:       "admin",
				JoinedAt:   time.Now().AddDate(0, -1, 5), // 1 month 5 days ago
				IsActive:   true,
				LastActive: time.Now().AddDate(0, 0, -2), // 2 days ago
			},
			{
				UserID:     "user7",
				Role:       "admin",
				JoinedAt:   time.Now().AddDate(0, -1, 5), // 1 month 5 days ago
				IsActive:   true,
				LastActive: time.Now().AddDate(0, 0, -2), // 2 days ago
			},
		},
	}

	// Example 1: Count all active members
	activeCount := chat.CountMembers(&chat1, chat.MemberIsActive())
	fmt.Printf("\nActive members: %d\n", activeCount) // Output: Active members: 2

	// Example 2: Count all admins
	adminCount := chat.CountMembers(&chat1, chat.MemberWithRole("admin"))
	fmt.Printf("Admin members: %d\n", adminCount) // Output: Admin members: 2

	// Example 3: Count members with "Manager" in their title
	managerCount := chat.CountMembers(&chat1, chat.MemberWithCustomTitle("Manager"))
	fmt.Printf("Manager members: %d\n", managerCount) // Output: Manager members: 2

	//// Example 4: Count members who joined in the last 2 months
	//recentMembers := chat.CountMembers(&chat1, func(member chat.Member) bool {
	//	return member.JoinedAt.After(time.Now().AddDate(0, -2, 0))
	//})
	//fmt.Printf("Recent members (joined in last 2 months): %d\n", recentMembers) // Output: Recent members: 2
	//
	//// Example 5: Count active admins
	//activeAdmins := chat.CountMembers(&chat1, func(member chat.Member) bool {
	//	return member.Role == "admin" && member.IsActive
	//})
	//fmt.Printf("Active admins: %d\n", activeAdmins) // Output: Active admins: 1

	//-----------------------------------------------------------------------

	fmt.Println("Sorting----------------------------")

	fmt.Println("Original order:")
	for i, m := range chat1.Members {
		fmt.Printf("%d: %s (%s) - Joined: %s, Active: %t, Last Active: %s\n", i, m.UserID, m.Role, m.JoinedAt.Format("2006-01-02"), m.IsActive, m.LastActive.Format("2006-01-02"))
	}

	// Sort by role
	chat.SortByRole(chat1.Members)
	fmt.Println("\nSorted by role:")
	for i, m := range chat1.Members {
		fmt.Printf("%d: %s (%s)\n", i, m.UserID, m.Role)
	}

	// Sort by join date (newest first)
	chat.SortByJoinedAt(chat1.Members, false)
	fmt.Println("\nSorted by join date (newest first):")
	for i, m := range chat1.Members {
		fmt.Printf("%d: %s - Joined: %s\n", i, m.UserID, m.JoinedAt.Format("2006-01-02"))
	}

	// Sort by last activity (most recent first)
	chat.SortByLastActive(chat1.Members, false)
	fmt.Println("\nSorted by last activity (most recent first):")
	for i, m := range chat1.Members {
		fmt.Printf("%d: %s - Last Active: %s\n", i, m.UserID, m.LastActive.Format("2006-01-02"))
	}

	// Sort by activity status (active first)
	chat.SortByActivityStatus(chat1.Members)
	fmt.Println("\nSorted by activity status (active first):")
	for i, m := range chat1.Members {
		fmt.Printf("%d: %s - Active: %t, Last Active: %s\n",
			i, m.UserID, m.IsActive, m.LastActive.Format("2006-01-02"))
	}

	// Multi-level sort: first by role, then by join date
	chat.SortByRoleThenJoinDate(chat1.Members)
	fmt.Println("\nSorted by role then join date:")
	for i, m := range chat1.Members {
		fmt.Printf("%d: %s (%s) - Joined: %s\n",
			i, m.UserID, m.Role, m.JoinedAt.Format("2006-01-02"))
	}
}
