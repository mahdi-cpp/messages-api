package chat

import (
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/iris-tools/search"
)

type SearchOptions struct {
	ChatID                uuid.UUID       `form:"chatId,omitempty"`
	MessageID             uuid.UUID       `form:"messageId,omitempty"`
	Type                  string          `form:"type,omitempty"` // "private", "group", "channel", "supergroup"
	Title                 string          `form:"title,omitempty"`
	Username              string          `form:"username,omitempty"` // Unique identifier for public channels/groups
	Description           string          `form:"description,omitempty"`
	Avatar                string          `form:"avatar,omitempty"`                // Chat profile photo
	PinnedMessageId       int             `form:"pinnedMessageId,omitempty"`       // ID of pinned message
	MessageAutoDeleteTime int             `form:"messageAutoDeleteTime,omitempty"` // Auto-delete timer
	Permissions           Permissions     `form:"permissions,omitempty"`           // Default chat permissions
	SlowModeDelay         int             `form:"slowModeDelay,omitempty"`         // Slow mode delay in seconds
	StickerSetName        string          `form:"stickerSetName,omitempty"`        // Name of group sticker set
	CanSetStickerSet      *bool           `form:"canSetStickerSet,omitempty"`      // Can set sticker set
	IsVerified            *bool           `form:"isVerified,omitempty"`
	IsRestricted          *bool           `form:"isRestricted,omitempty"`
	IsCreator             *bool           `form:"isCreator,omitempty"`
	IsScam                *bool           `form:"isScam,omitempty"`
	IsFake                *bool           `form:"isFake,omitempty"`
	InviteLink            string          `form:"inviteLink,omitempty"`         // Generated invite link
	LinkedChatID          int             `form:"linkedChatId,omitempty"`       // Linked discussion chat for channels
	Location              *Location       `form:"location,omitempty"`           // For location-based chats
	Members               *[]Member       `form:"members,omitempty"`            // Detailed member list
	ParticipantsCount     int             `form:"participantsCount,omitempty"`  // Cache member count
	ActiveUsernames       []string        `form:"activeUsernames,omitempty"`    // For multiple usernames
	AvailableReactions    []string        `form:"availableReactions,omitempty"` // Available emoji reactions
	Theme                 string          `form:"theme,omitempty"`              // Chat theme
	UnreadCount           int             `form:"unreadCount,omitempty"`        // Unread message count
	LastMessage           *MessagePreview `form:"lastMessage,omitempty"`        // Last message preview
	IsPinned              *bool           `form:"isPinned,omitempty"`           // Pinned in user's list
	PinOrder              int             `form:"pinOrder,omitempty"`           // Position in pinned list
	MuteUntil             time.Time       `form:"muteUntil,omitempty"`          // Mute notification until

	// Date filters
	CreatedAfter  *time.Time `form:"createdAfter,omitempty"`
	CreatedBefore *time.Time `form:"createdBefore,omitempty"`
	ActiveAfter   *time.Time `form:"activeAfter,omitempty"`

	// Sorting
	Sort      string `form:"sort,omitempty"`      // "title", "created", "members", "lastActivity"
	SortOrder string `form:"sortOrder,omitempty"` // "asc" or "desc"

	// Pagination
	Page int `form:"page,omitempty"`
	Size int `form:"size,omitempty"`
}

const MaxLimit = 1000

var LessFunks = map[string]search.LessFunction[*Chat]{
	"id":        func(a, b *Chat) bool { return a.ID.String() < b.ID.String() },
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
		if with.ChatID != uuid.Nil && c.ID != with.ChatID {
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
	if with.Sort != "" {
		lessFn := GetLessFunc(with.Sort, with.SortOrder)
		if lessFn != nil {
			search.SortIndexedItems(results, lessFn)
		}
	}

	// Extract final assets
	final := make([]*Chat, len(results))
	for i, item := range results {
		final[i] = item.Value
	}

	if with.Size == 0 { // if not set default is MAX_LIMIT
		with.Size = MaxLimit
	}

	// Apply pagination
	start := with.Page

	// Check if the start index is out of bounds. If so, return an empty slice.
	if start >= len(final) {
		return []*Chat{}
	}

	end := start + with.Size
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
func MemberWithUserID(userID uuid.UUID) search.Criteria[*Member] {
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
func CountMembers(chat *Chat, criteria search.Criteria[*Member]) int {
	count := 0
	for _, member := range chat.Members {
		if criteria(&member) {
			count++
		}
	}
	return count
}

// GetMatchingMembers returns all members in a chat that match the criteria
func GetMatchingMembers(chat Chat, criteria search.Criteria[*Member]) []Member {
	var matches []Member
	for _, member := range chat.Members {
		if criteria(&member) {
			matches = append(matches, member)
		}
	}
	return matches
}

// Sorting functions for members
// ---------------------------------------------------------------------

// SortByRole sorts members by role (creator > admin > member)
func SortByRole(members []Member) {
	rolePriority := map[string]int{
		"creator": 0,
		"admin":   1,
		"member":  2,
	}

	sort.Slice(members, func(i, j int) bool {
		return rolePriority[members[i].Role] < rolePriority[members[j].Role]
	})
}

// SortByJoinedAt sorts members by join date (newest first by default)
func SortByJoinedAt(members []Member, ascending bool) {
	sort.Slice(members, func(i, j int) bool {
		if ascending {
			return members[i].JoinedAt.Before(members[j].JoinedAt)
		}
		return members[i].JoinedAt.After(members[j].JoinedAt)
	})
}

// SortByLastActive sorts members by last activity time (most recent first by default)
func SortByLastActive(members []Member, ascending bool) {
	sort.Slice(members, func(i, j int) bool {
		if ascending {
			return members[i].LastActive.Before(members[j].LastActive)
		}
		return members[i].LastActive.After(members[j].LastActive)
	})
}

// SortByActivityStatus sorts members (active first, then inactive)
func SortByActivityStatus(members []Member) {
	sort.Slice(members, func(i, j int) bool {
		// Active members first
		if members[i].IsActive && !members[j].IsActive {
			return true
		}
		if !members[i].IsActive && members[j].IsActive {
			return false
		}
		// If both have same status, sort by last active
		return members[i].LastActive.After(members[j].LastActive)
	})
}

// SortByRoleThenJoinDate Multi-level sorting: First by role, then by join date
func SortByRoleThenJoinDate(members []Member) {
	sort.Slice(members, func(i, j int) bool {
		// First, sort by role priority
		rolePriority := map[string]int{
			"creator": 0,
			"admin":   1,
			"member":  2,
		}

		if rolePriority[members[i].Role] != rolePriority[members[j].Role] {
			return rolePriority[members[i].Role] < rolePriority[members[j].Role]
		}

		// If same role, sort by join date (newest first)
		return members[i].JoinedAt.After(members[j].JoinedAt)
	})
}
