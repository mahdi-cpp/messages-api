package chat

func test() {

	//var chats []*Chat

	// Example 1: Find all chats that have a specific user as a member
	//userID := "user123"
	//chatsWithUser := search.Search(chats, HasMemberWith(MemberWithUserID(userID)))
	//
	//fmt.Println(chatsWithUser)
	//
	//// Example 2: Find all chats that have an admin member
	//chatsWithAdmins := search.Search(chats, HasMemberWith(MemberWithRole("admin")))
	//fmt.Println(chatsWithAdmins)
	//
	//// Example 3: Find all chats that have active members who joined in the last week
	//recentActiveMembers := Search(chats, HasMemberWith(
	//	func(member Member) bool {
	//		return member.IsActive &&
	//			member.JoinedAt.After(time.Now().AddDate(0, 0, -7))
	//	},
	//))
	//fmt.Println(recentActiveMembers)

	//// Example 4: Find all chats that have members with "manager" in their custom title
	//chatsWithManagers := Search(chats, search.HasMemberWith(
	//	MemberWithCustomTitle("manager"),
	//))
	//
	//// Example 5: Combine multiple member criteria
	//activeAdmins := Search(chats, search.HasMemberWith(
	//	func(member Member) bool {
	//		return member.Role == "admin" &&
	//			member.IsActive &&
	//			member.JoinedAt.After(time.Now().AddDate(0, -1, 0)) // joined in last month
	//	},
	//))
}
