package message

import (
	"time"

	"github.com/google/uuid"
)

func (a *Message) SetID(id uuid.UUID) { a.ID = id }
func (a *Message) GetID() uuid.UUID   { return a.ID }

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypePhoto    MessageType = "photo"
	MessageTypeVideo    MessageType = "video"
	MessageTypeDocument MessageType = "document"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeVoice    MessageType = "voice"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeLocation MessageType = "location"
	MessageTypeContact  MessageType = "contact"
	MessageTypePoll     MessageType = "poll"
)

type Message struct {
	ID          uuid.UUID   `json:"id"`
	ChatID      uuid.UUID   `json:"chatId"`
	UserID      uuid.UUID   `json:"userId"`
	Content     string      `json:"content"`
	MessageType MessageType `json:"type"` // Changed name to avoid conflict

	// Media fields
	MediaURL     string `json:"mediaUrl,omitempty"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	FileSize     int64  `json:"fileSize,omitempty"`
	Duration     int    `json:"duration,omitempty"`
	Width        int    `json:"width,omitempty"`
	Height       int    `json:"height,omitempty"`
	MimeType     string `json:"mimeType,omitempty"`

	// Message attributes
	ReplyToMessageID uuid.UUID    `json:"replyToMessageId,omitempty"`
	ForwardedFrom    *ForwardInfo `json:"forwardedFrom,omitempty"`
	Entities         []Entity     `json:"entities,omitempty"`
	Views            int          `json:"views,omitempty"`
	Reactions        []Reaction   `json:"reactions,omitempty"`
	IsEdited         bool         `json:"isEdited,omitempty"`
	IsPinned         bool         `json:"isPinned,omitempty"`
	IsDeleted        bool         `json:"isDeleted,omitempty"`

	// Additional data types
	Poll     *Poll     `json:"poll,omitempty"`
	Location *Location `json:"location,omitempty"`
	Contact  *Contact  `json:"contact,omitempty"`

	// Timestamps and metadata
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     time.Time `json:"deletedAt,omitempty"`
	ExpiryAt      time.Time `json:"expiryAt,omitempty"`
	EncryptionKey string    `json:"encryptionKey,omitempty"`
	Version       string    `json:"version,omitempty"`
}

// ForwardInfo Supporting structs
type ForwardInfo struct {
	FromChatID    uuid.UUID `json:"fromChatId"`
	FromMessageID uuid.UUID `json:"fromMessageId"`
	FromUserID    uuid.UUID `json:"fromUserId"`
	OriginalDate  time.Time `json:"originalDate"`
}

type Entity struct {
	Type   string    `json:"type"`             // mention, hashtag, bot_command, url, etc.
	Offset int       `json:"offset"`           // Offset in UTF-16 code units
	Length int       `json:"length"`           // Length in UTF-16 code units
	URL    string    `json:"url,omitempty"`    // For "text_link" only
	UserID uuid.UUID `json:"userId,omitempty"` // For "mention" only
}

type Reaction struct {
	Emoji   string      `json:"emoji"`
	Count   int         `json:"count"`
	UserIDs []uuid.UUID `json:"userIds,omitempty"` // Users who used this reaction
}

type Poll struct {
	Question              string       `json:"question"`
	Options               []PollOption `json:"options"`
	TotalVotes            int          `json:"totalVotes"`
	IsAnonymous           bool         `json:"isAnonymous"`
	Type                  string       `json:"type"` // regular or quiz
	AllowsMultipleAnswers bool         `json:"allowsMultipleAnswers"`
	CloseDate             time.Time    `json:"closeDate,omitempty"`
}

type PollOption struct {
	Text     string      `json:"text"`
	Votes    int         `json:"votes"`
	VoterIDs []uuid.UUID `json:"voterIds,omitempty"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Accuracy  float64 `json:"accuracy,omitempty"` // Accuracy radius in meters
}

type Contact struct {
	PhoneNumber string    `json:"phoneNumber"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName,omitempty"`
	UserID      uuid.UUID `json:"userId,omitempty"` // If the contact is a registered user
}

type TypingStatus struct {
	ChatID uuid.UUID `json:"chatId"`
	UserID uuid.UUID `json:"userId"`
	Typing bool      `json:"typing"`
}
