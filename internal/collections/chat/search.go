package chat

import (
	"strings"
	"time"

	"github.com/mahdi-cpp/iris-tools/search"
)

const MaxLimit = 1000

var LessFunks = map[string]search.LessFunction[*Chat]{
	"id":        func(a, b *Chat) bool { return a.ID < b.ID },
	"createdAt": func(a, b *Chat) bool { return a.CreatedAt.Before(b.CreatedAt) },
	"updatedAt": func(a, b *Chat) bool { return a.UpdatedAt.Before(b.UpdatedAt) },
}

func GetLessFunc(sortBy, sortOrder string) search.LessFunction[*Chat] {

	fn, exists := LessFunks[sortBy]
	if !exists {
		return nil
	}

	if sortOrder == "end" {
		return func(a, b *Chat) bool { return !fn(a, b) }
	}
	return fn
}

func BuildChatCriteria(with *SearchOptions) search.Criteria[*Chat] {

	return func(c *Chat) bool {

		// ID filter
		if with.ID != "" && c.ID != with.ID {
			return false
		}

		// Boolean flags
		if with.CanSetStickerSet != nil && c.CanSetStickerSet != *with.CanSetStickerSet {
			return false
		}
		if with.IsVerified != nil && c.IsVerified != *with.IsVerified {
			return false
		}
		if with.IsRestricted != nil && c.IsRestricted != *with.IsRestricted {
			return false
		}
		if with.IsCreator != nil && c.IsCreator != *with.IsCreator {
			return false
		}
		if with.IsScam != nil && c.IsScam != *with.IsScam {
			return false
		}
		if with.IsFake != nil && c.IsFake != *with.IsFake {
			return false
		}

		// Date filters
		if with.CreatedAfter != nil && c.CreatedAt.Before(*with.CreatedAfter) {
			return false
		}
		if with.CreatedBefore != nil && c.CreatedAt.After(*with.CreatedBefore) {
			return false
		}

		return true
	}
}

func Search(chats []*Chat, with *SearchOptions) []*Chat {

	// Build criteria
	criteria := BuildChatCriteria(with)

	// Execute search_manager
	results := search.Find(chats, criteria)

	// Sort results if needed
	if with.SortBy != "" {
		lessFn := GetLessFunc(with.SortBy, with.SortOrder)
		if lessFn != nil {
			search.SortIndexedItems(results, lessFn)
		}
	}

	// Extract final assets
	final := make([]*Chat, len(results))
	for i, item := range results {
		final[i] = item.Value
	}

	if with.Limit == 0 { // if not set default is MAX_LIMIT
		with.Limit = MaxLimit
	}

	// Apply pagination
	start := with.Offset
	end := start + with.Limit
	if end > len(final) {
		end = len(final)
	}
	return final[start:end]
}

// Chat-specific search functions
// ---------------------------------------------------------------------

// HasMemberWith creates a criteria that checks if a chat has at least one member
// matching the provided member criteria
func HasMemberWith(memberCriteria search.Criteria[*Member]) search.Criteria[*Chat] {
	return func(chat *Chat) bool {
		for _, member := range chat.Members {
			if memberCriteria(&member) {
				return true
			}
		}
		return false
	}
}

// Member-specific criteria functions
// ---------------------------------------------------------------------

// MemberWithUserID checks if a member has a specific user ID
func MemberWithUserID(userID string) search.Criteria[*Member] {
	return func(member *Member) bool {
		return member.UserID == userID
	}
}

// MemberWithRole checks if a member has a specific role
func MemberWithRole(role string) search.Criteria[*Member] {
	return func(member *Member) bool {
		return member.Role == role
	}
}

// MemberWithCustomTitle checks if a member's custom title contains the query
func MemberWithCustomTitle(query string) search.Criteria[*Member] {
	return func(member *Member) bool {
		return search.StringContains(member.CustomTitle, query)
	}
}

// MemberIsActive checks if a member is active
func MemberIsActive() search.Criteria[*Member] {
	return func(member *Member) bool {
		return member.IsActive
	}
}

// MemberJoinedAfter checks if a member joined after a specific time
func MemberJoinedAfter(time time.Time) search.Criteria[*Member] {
	return func(member *Member) bool {
		return member.JoinedAt.After(time)
	}
}

// MemberJoinedBefore checks if a member joined before a specific time
func MemberJoinedBefore(time time.Time) search.Criteria[*Member] {
	return func(member *Member) bool {
		return member.JoinedAt.Before(time)
	}
}

// ActiveAdminsJoinedAfter finds active admins who joined after a specific time
func ActiveAdminsJoinedAfter(time time.Time) search.Criteria[*Member] {
	return func(member *Member) bool {
		return member.Role == "admin" && member.IsActive && member.JoinedAt.After(time)
	}
}

// MembersWithTitlePattern finds members with a specific title pattern
func MembersWithTitlePattern(pattern string) search.Criteria[*Member] {
	return func(member *Member) bool {
		return strings.Contains(strings.ToLower(member.CustomTitle),
			strings.ToLower(pattern))
	}
}

// CountMembers returns the number of members matching the criteria in a chat
func CountMembers(chat Chat, criteria search.Criteria[Member]) int {
	count := 0
	for _, member := range chat.Members {
		if criteria(member) {
			count++
		}
	}
	return count
}

// GetMatchingMembers returns all members in a chat that match the criteria
func GetMatchingMembers(chat Chat, criteria search.Criteria[Member]) []Member {
	var matches []Member
	for _, member := range chat.Members {
		if criteria(member) {
			matches = append(matches, member)
		}
	}
	return matches
}
