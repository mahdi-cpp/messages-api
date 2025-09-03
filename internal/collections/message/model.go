package message

import (
	"time"
)

func (a *Message) SetID(id string)          { a.ID = id }
func (a *Message) SetCreatedAt(t time.Time) { a.CreatedAt = t }
func (a *Message) SetUpdatedAt(t time.Time) { a.UpdatedAt = t }
func (a *Message) GetID() string            { return a.ID }
func (a *Message) GetCreatedAt() time.Time  { return a.CreatedAt }
func (a *Message) GetUpdatedAt() time.Time  { return a.UpdatedAt }

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
	ID               string       `json:"id"`
	ChatID           string       `json:"chatID"`           // Identifier of the chat
	UserID           string       `json:"userID"`           // Sender identifier
	Content          string       `json:"content"`          // Text content or caption
	Type             string       `json:"type"`             // text, photo, video, audio, document, voice, location, contact, poll, etc.
	MediaURL         string       `json:"mediaUrl"`         // URL to media file
	ThumbnailURL     string       `json:"thumbnailUrl"`     // URL to thumbnail
	FileSize         int64        `json:"fileSize"`         // Size of media in bytes
	Duration         int          `json:"duration"`         // For audio/video in seconds
	Width            int          `json:"width"`            // For images/videos
	Height           int          `json:"height"`           // For images/videos
	MimeType         string       `json:"mimeType"`         // MIME type of media
	ReplyToMessageID string       `json:"replyToMessageId"` // ID of message being replied to
	ForwardedFrom    *ForwardInfo `json:"forwardedFrom"`    // Info about forwarded message
	Entities         []Entity     `json:"entities"`         // Text formatting entities
	Views            int          `json:"views"`            // View count (for channels)
	Reactions        []Reaction   `json:"reactions"`        // Emoji reactions
	IsEdited         bool         `json:"isEdited"`         // Has the message been edited?
	IsPinned         bool         `json:"isPinned"`         // Is this message pinned?
	IsDeleted        bool         `json:"isDeleted"`        // Soft delete flag
	Poll             *Poll        `json:"poll"`             // Poll data if message is a poll
	Location         *Location    `json:"location"`         // Location data if message contains location
	Contact          *Contact     `json:"contact"`          // Contact data if message contains contact
	EncryptionKey    string       `json:"encryptionKey"`    // For end-to-end encryption
	EditAt           time.Time    `json:"editAt"`           // When was it last edited
	CreatedAt        time.Time    `json:"createdAt"`        // When was it deleted
	DeletedAt        time.Time    `json:"deletedAt"`        // When was it created
	UpdatedAt        time.Time    `json:"updatedAt"`        // When was it last modified
	ExpiryAt         time.Time    `json:"expiryAt"`         // For self-destructing messages
	Version          string       `json:"version"`          // For optimistic concurrency control
}

// ForwardInfo Supporting structs
type ForwardInfo struct {
	FromChatID    int       `json:"fromChatId"`
	FromMessageID int       `json:"fromMessageId"`
	FromUserID    int       `json:"fromUserId"`
	OriginalDate  time.Time `json:"originalDate"`
}

type Entity struct {
	Type   string `json:"type"`   // mention, hashtag, bot_command, url, etc.
	Offset int    `json:"offset"` // Offset in UTF-16 code units
	Length int    `json:"length"` // Length in UTF-16 code units
	URL    string `json:"url"`    // For "text_link" only
	UserID string `json:"userId"` // For "mention" only
}

type Reaction struct {
	Emoji   string   `json:"emoji"`   // The emoji used
	Count   int      `json:"count"`   // Number of users who used this reaction
	UserIDs []string `json:"userIds"` // Users who used this reaction
}

type Poll struct {
	Question              string       `json:"question"`
	Options               []PollOption `json:"options"`
	TotalVotes            int          `json:"totalVotes"`
	IsAnonymous           bool         `json:"isAnonymous"`
	Type                  string       `json:"type"` // regular or quiz
	AllowsMultipleAnswers bool         `json:"allowsMultipleAnswers"`
	CloseDate             time.Time    `json:"closeDate"`
}

type PollOption struct {
	Text     string   `json:"text"`
	Votes    int      `json:"votes"`
	VoterIDs []string `json:"voterIds"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Accuracy  float64 `json:"accuracy"` // Accuracy radius in meters
}

type Contact struct {
	PhoneNumber string `json:"phoneNumber"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	UserID      string `json:"userId"` // If the contact is a registered user
}

type TypingStatus struct {
	ChatID string `json:"chatID"`
	UserID string `json:"userID"`
	Typing bool   `json:"typing"`
}
