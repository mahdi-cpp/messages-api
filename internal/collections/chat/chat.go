package chat

import (
	"time"

	"github.com/mahdi-cpp/iris-tools/update"
)

func (a *Chat) SetID(id string)          { a.ID = id }
func (a *Chat) SetCreatedAt(t time.Time) { a.CreatedAt = t }
func (a *Chat) SetUpdatedAt(t time.Time) { a.UpdatedAt = t }
func (a *Chat) GetID() string            { return a.ID }
func (a *Chat) GetCreatedAt() time.Time  { return a.CreatedAt }
func (a *Chat) GetUpdatedAt() time.Time  { return a.UpdatedAt }

type Chat struct {
	ID                    string          `json:"id"`
	Type                  string          `json:"type"` // "private", "group", "channel", "supergroup"
	Title                 string          `json:"title"`
	Username              string          `json:"username"` // Unique identifier for public channels/groups
	Description           string          `json:"description"`
	Avatar                string          `json:"avatar"`                // Chat profile photo
	PinnedMessageId       int             `json:"pinnedMessageId"`       // ID of pinned message
	MessageAutoDeleteTime int             `json:"messageAutoDeleteTime"` // Auto-delete timer
	Permissions           Permissions     `json:"permissions"`           // Default chat permissions
	SlowModeDelay         int             `json:"slowModeDelay"`         // Slow mode delay in seconds
	StickerSetName        string          `json:"stickerSetName"`        // Name of group sticker set
	CanSetStickerSet      bool            `json:"canSetStickerSet"`      // Can set sticker set
	IsVerified            bool            `json:"isVerified"`
	IsRestricted          bool            `json:"isRestricted"`
	IsCreator             bool            `json:"isCreator"`
	IsScam                bool            `json:"isScam"`
	IsFake                bool            `json:"isFake"`
	InviteLink            string          `json:"inviteLink"`         // Generated invite link
	LinkedChatId          int             `json:"linkedChatId"`       // Linked discussion chat for channels
	Location              *Location       `json:"location"`           // For location-based chats
	Members               []Member        `json:"members"`            // Detailed member list
	ParticipantsCount     int             `json:"participantsCount"`  // Cache member count
	ActiveUsernames       []string        `json:"activeUsernames"`    // For multiple usernames
	AvailableReactions    []string        `json:"availableReactions"` // Available emoji reactions
	Theme                 string          `json:"theme"`              // Chat theme
	UnreadCount           int             `json:"unreadCount"`        // Unread message count
	LastMessage           *MessagePreview `json:"lastMessage"`        // Last message preview
	IsPinned              bool            `json:"isPinned"`           // Pinned in user's list
	PinOrder              int             `json:"pinOrder"`           // Position in pinned list
	MuteUntil             time.Time       `json:"muteUntil"`          // Mute notification until
	CreatedAt             time.Time       `json:"createdAt"`
	UpdatedAt             time.Time       `json:"updatedAt"`
	DeletedAt             *time.Time      `json:"deletedAt"` // Use pointer for optional field
	Version               int             `json:"version"`   // Use int for versioning
}

type Member struct {
	UserID      string    `json:"userID"`
	Role        string    `json:"role"` // "member", "admin", "creator"
	JoinedAt    time.Time `json:"joinedAt"`
	CustomTitle string    `json:"customTitle"` // For custom admin titles
	IsActive    bool      `json:"IsActive"`
}

type SearchOptions struct {
	ID                    string          `json:"id"`
	Type                  string          `json:"type"` // "private", "group", "channel", "supergroup"
	Title                 string          `json:"title"`
	Username              string          `json:"username"` // Unique identifier for public channels/groups
	Description           string          `json:"description"`
	Avatar                string          `json:"avatar"`                // Chat profile photo
	PinnedMessageId       int             `json:"pinnedMessageId"`       // ID of pinned message
	MessageAutoDeleteTime int             `json:"messageAutoDeleteTime"` // Auto-delete timer
	Permissions           Permissions     `json:"permissions"`           // Default chat permissions
	SlowModeDelay         int             `json:"slowModeDelay"`         // Slow mode delay in seconds
	StickerSetName        string          `json:"stickerSetName"`        // Name of group sticker set
	CanSetStickerSet      *bool           `json:"canSetStickerSet"`      // Can set sticker set
	IsVerified            *bool           `json:"isVerified"`
	IsRestricted          *bool           `json:"isRestricted"`
	IsCreator             *bool           `json:"isCreator"`
	IsScam                *bool           `json:"isScam"`
	IsFake                *bool           `json:"isFake"`
	InviteLink            string          `json:"inviteLink"`         // Generated invite link
	LinkedChatId          int             `json:"linkedChatId"`       // Linked discussion chat for channels
	Location              *Location       `json:"location"`           // For location-based chats
	Members               *[]Member       `json:"members"`            // Detailed member list
	ParticipantsCount     int             `json:"participantsCount"`  // Cache member count
	ActiveUsernames       []string        `json:"activeUsernames"`    // For multiple usernames
	AvailableReactions    []string        `json:"availableReactions"` // Available emoji reactions
	Theme                 string          `json:"theme"`              // Chat theme
	UnreadCount           int             `json:"unreadCount"`        // Unread message count
	LastMessage           *MessagePreview `json:"lastMessage"`        // Last message preview
	IsPinned              *bool           `json:"isPinned"`           // Pinned in user's list
	PinOrder              int             `json:"pinOrder"`           // Position in pinned list
	MuteUntil             time.Time       `json:"muteUntil"`          // Mute notification until

	// Date filters
	CreatedAfter  *time.Time `json:"createdAfter,omitempty"`
	CreatedBefore *time.Time `json:"createdBefore,omitempty"`
	ActiveAfter   *time.Time `json:"activeAfter,omitempty"`

	// Pagination
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`

	// Sorting
	SortBy    string `json:"sortBy,omitempty"`    // "title", "created", "members", "lastActivity"
	SortOrder string `json:"sortOrder,omitempty"` // "asc" or "desc"
}

type Permissions struct {
	CanSendMessages      bool `json:"canSendMessages"`
	CanSendMediaMessages bool `json:"canSendMediaMessages"`
	CanSendPolls         bool `json:"canSendPolls"`
	CanChangeInfo        bool `json:"canChangeInfo"`
	CanInviteUsers       bool `json:"canInviteUsers"`
	CanPinMessages       bool `json:"canPinMessages"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
}

type MessagePreview struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Type      string    `json:"type"` // "text", "photo", "video", etc.
	AuthorId  int       `json:"authorId"`
	Timestamp time.Time `json:"timestamp"`
}

type UpdateOptions struct {
	ChatIDs []string `json:"chatIDs,omitempty"` // Asset Ids

	Type        string `json:"type"` // "private", "group", "channel", "supergroup"
	Title       string `json:"title"`
	Username    string `json:"username"` // Unique identifier for public channels/groups
	Description string `json:"description"`
	Avatar      string `json:"avatar"` // Chat profile photo

	CanSetStickerSet *bool `json:"canSetStickerSet"` // Can set sticker set
	IsVerified       *bool `json:"isVerified"`
	IsRestricted     *bool `json:"isRestricted"`
	IsCreator        *bool `json:"isCreator"`
	IsScam           *bool `json:"isScam"`
	IsFake           *bool `json:"isFake"`

	ActiveUsernames       *[]string `json:"users,omitempty"`                 // Full users replacement
	AddActiveUsernames    []string  `json:"AddActiveUsernames,omitempty"`    // Users to add
	RemoveActiveUsernames []string  `json:"removeActiveUsernames,omitempty"` // Users to remove

	Members        *[]Member
	AddMembers     []Member
	RemoveMembers  []Member
	MembersUpdates []update.NestedFieldUpdate[Member]
}

// Key extractors for nested structs
func memberKeyExtractor(m Member) string { return m.UserID }
