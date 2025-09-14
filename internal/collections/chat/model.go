package chat

import (
	"time"

	"github.com/google/uuid"
)

func (c *Chat) SetID(id uuid.UUID) { c.ID = id }
func (c *Chat) GetID() uuid.UUID   { return c.ID }

type Chat struct {
	ID                    uuid.UUID       `json:"id"`
	Type                  string          `json:"type"`
	Title                 string          `json:"title"`
	Username              string          `json:"username,omitempty"`
	Description           string          `json:"description,omitempty"`
	Avatar                string          `json:"avatar,omitempty"`
	PinnedMessageID       string          `json:"pinnedMessageId,omitempty"` // Changed type to string for consistency
	MessageAutoDeleteTime int             `json:"messageAutoDeleteTime,omitempty"`
	Permissions           Permissions     `json:"permissions"`
	SlowModeDelay         int             `json:"slowModeDelay,omitempty"`
	StickerSetName        string          `json:"stickerSetName,omitempty"`
	CanSetStickerSet      bool            `json:"canSetStickerSet,omitempty"`
	IsVerified            bool            `json:"isVerified"`
	IsRestricted          bool            `json:"isRestricted"`
	IsCreator             bool            `json:"isCreator"`
	IsScam                bool            `json:"isScam"`
	IsFake                bool            `json:"isFake"`
	InviteLink            string          `json:"inviteLink,omitempty"`
	LinkedChatID          string          `json:"linkedChatId,omitempty"` // Changed type to string for consistency
	Location              *Location       `json:"location,omitempty"`
	Members               []Member        `json:"members,omitempty"`
	ParticipantsCount     int             `json:"participantsCount"`
	ActiveUsernames       []string        `json:"activeUsernames,omitempty"`
	AvailableReactions    []string        `json:"availableReactions,omitempty"`
	Theme                 string          `json:"theme,omitempty"`
	UnreadCount           int             `json:"unreadCount,omitempty"`
	LastMessage           *MessagePreview `json:"lastMessage,omitempty"`
	IsPinned              bool            `json:"isPinned,omitempty"`
	PinOrder              int             `json:"pinOrder,omitempty"`
	MuteUntil             *time.Time      `json:"muteUntil,omitempty"`
	CreatedAt             time.Time       `json:"createdAt"`
	UpdatedAt             time.Time       `json:"updatedAt"`
	DeletedAt             *time.Time      `json:"deletedAt,omitempty"`
	Version               string          `json:"version"`
}

type Member struct {
	UserID      uuid.UUID `json:"userId"`
	Role        string    `json:"role"`        // "member", "admin", "creator"
	CustomTitle string    `json:"customTitle"` // For custom admin titles
	IsActive    bool      `json:"isActive"`
	LastActive  time.Time `json:"lastActive"` // Added for sorting by activity
	JoinedAt    time.Time `json:"joinedAt"`
}

type Permissions struct {
	CanSendMessages       bool `json:"canSendMessages"`
	CanSendMedia          bool `json:"canSendMedia"`
	CanSendPolls          bool `json:"canSendPolls"`
	CanSendOtherMessages  bool `json:"canSendOtherMessages"`
	CanAddWebPagePreviews bool `json:"canAddWebPagePreviews"`
	CanChangeInfo         bool `json:"canChangeInfo"`
	CanInviteUsers        bool `json:"canInviteUsers"`
	CanPinMessages        bool `json:"canPinMessages"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
}

type MessagePreview struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Type      string    `json:"type"`
	AuthorID  string    `json:"authorId"` // Changed name and type for consistency
	Timestamp time.Time `json:"timestamp"`
}
