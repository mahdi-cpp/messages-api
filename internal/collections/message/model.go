package message

import (
	"time"

	"github.com/google/uuid"
)

func (a *Message) SetID(id uuid.UUID) { a.ID = id }
func (a *Message) GetID() uuid.UUID   { return a.ID }

func (i *Index) SetID(id uuid.UUID) { i.ID = id }
func (i *Index) GetID() uuid.UUID   { return i.ID }

// Index ساختاری برای نگهداری اطلاعات کلیدی و ایندکس‌گذاری است.
type Index struct {
	ID          uuid.UUID `json:"id" `
	ChatID      uuid.UUID `json:"chatId" `
	UserID      uuid.UUID `json:"userId" `
	IsEdited    bool      `json:"isEdited" `
	IsPinned    bool      `json:"isPinned" `
	IsDeleted   bool      `json:"isDeleted" `
	MediaUnread bool      `json:"mediaUnread" `
	CreatedAt   time.Time `json:"createdAt" `
}

type Message struct {
	ID        uuid.UUID `json:"id" index:"true"`
	ChatID    uuid.UUID `json:"chatId" index:"true"`
	UserID    uuid.UUID `json:"userId" index:"true"`
	Caption   string    `json:"caption"`
	Directory string    `json:"directory"`

	// Data types
	AssetType string    `json:"assetType"`
	Medias    []*Media  `json:"medias,omitempty"`
	Voice     *Voice    `json:"voice,omitempty"`
	Music     *Music    `json:"music,omitempty"`
	Document  *Document `json:"document,omitempty"`
	Contact   *Contact  `json:"contact,omitempty"`
	Location  *Location `json:"location,omitempty"`
	Poll      *Poll     `json:"poll,omitempty"`

	// Message attributes
	ReplyToMessageID *uuid.UUID   `json:"replyToMessageId,omitempty"`
	ForwardedFrom    *ForwardInfo `json:"forwardedFrom,omitempty"`
	Entities         []Entity     `json:"entities,omitempty"`
	Views            int          `json:"views,omitempty"`
	Reactions        []Reaction   `json:"reactions"`
	IsEdited         bool         `json:"isEdited" index:"true" `
	IsPinned         bool         `json:"isPinned" index:"true"`
	IsDeleted        bool         `json:"isDeleted" index:"true" `
	MediaUnread      bool         `json:"mediaUnread" `
	Silent           bool         `json:"silent"`

	// Timestamps and metadata
	CreatedAt     time.Time `json:"createdAt" index:"true" `
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     time.Time `json:"deletedAt,omitempty"`
	EncryptionKey string    `json:"encryptionKey,omitempty"`
	Version       string    `json:"version"`
}

//--- Data Types

type Media struct {
	ID          uuid.UUID `json:"id"`
	FileSize    int       `json:"fileSize"`
	MimeType    string    `json:"mimeType"`
	Duration    int       `json:"duration,omitempty"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Orientation string    `json:"orientation"`
	Tags        []Tag     `json:"tags"`
}

type Tag struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
}

type Music struct {
	ID       uuid.UUID `json:"id"`
	Artist   string    `json:"artist"`
	Album    string    `json:"album"`
	FileSize int64     `json:"fileSize"`
	MimeType string    `json:"mimeType"`
	Duration int       `json:"duration"`
}

type Voice struct {
	ID       uuid.UUID `json:"id"`
	FileSize int64     `json:"fileSize"`
	MimeType string    `json:"mimeType"`
	Duration int       `json:"duration"`
}

type Document struct {
	PhoneNumber string    `json:"phoneNumber"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	UserID      uuid.UUID `json:"userId"` // If the contact is a registered user
}

type Contact struct {
	PhoneNumber string    `json:"phoneNumber"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	UserID      uuid.UUID `json:"userId"` // If the contact is a registered user
}
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Accuracy  float64 `json:"accuracy"` // Accuracy radius in meters
}

type Poll struct {
	Question              string       `json:"question"`
	Options               []PollOption `json:"options"`
	TotalVotes            int          `json:"totalVotes"`
	IsAnonymous           bool         `json:"isAnonymous"`
	Type                  string       `json:"type"`
	AllowsMultipleAnswers bool         `json:"allowsMultipleAnswers"`
	CloseDate             time.Time    `json:"closeDate,omitempty"`
}

type PollOption struct {
	Text     string      `json:"text"`
	Votes    int         `json:"votes"`
	VoterIDs []uuid.UUID `json:"voterIds"`
}

//---

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

type TypingStatus struct {
	ChatID uuid.UUID `json:"chatId"`
	UserID uuid.UUID `json:"userId"`
	Typing bool      `json:"typing"`
}
