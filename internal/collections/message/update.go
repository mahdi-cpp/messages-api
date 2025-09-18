package message

import (
	"time"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/iris-tools/update"
)

type UpdateOptions struct {
	MessageIDs []uuid.UUID `json:"messageIds,omitempty"`
	UserID     uuid.UUID   `json:"userId"`
	ChatID     uuid.UUID   `json:"chatId"`
	MessageID  uuid.UUID   `json:"messageId"`

	Content string `json:"content"`

	// Message attributes
	ReplyToMessageID string       `json:"replyToMessageId,omitempty"`
	ForwardedFrom    *ForwardInfo `json:"forwardedFrom,omitempty"`
	Entities         []Entity     `json:"entities,omitempty"`
	Reactions        []Reaction   `json:"reactions,omitempty"`
	IsEdited         *bool        `json:"isEdited,omitempty"`
	IsPinned         *bool        `json:"isPinned,omitempty"`
	IsDeleted        *bool        `json:"isDeleted,omitempty"`

	// Additional data types
	Poll     *Poll     `json:"poll,omitempty"`
	Location *Location `json:"location,omitempty"`
	Contact  *Contact  `json:"contact,omitempty"`
}

// Initialize updater
var metadataUpdater = update.NewUpdater[Message, UpdateOptions]()

func init() {

	// Configure scalar field updates
	metadataUpdater.AddScalarUpdater(func(a *Message, u UpdateOptions) {
		if u.Content != "" {
			a.Caption = u.Content
		}
	})

	// Set modification timestamp
	metadataUpdater.AddPostUpdateHook(func(a *Message) {
		a.UpdatedAt = time.Now()
	})
}

func Update(p *Message, update UpdateOptions) *Message {
	metadataUpdater.Apply(p, update)
	return p
}
