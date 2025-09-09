package message

import (
	"time"

	"github.com/mahdi-cpp/iris-tools/update"
)

type UpdateOptions struct {
	MessageIDs []string `json:"messageIds,omitempty"`
	UserID     string   `json:"userId"`
	ChatID     string   `json:"chatId"`
	MessageID  string   `json:"messageId"`

	Content     string      `json:"content"`
	MessageType MessageType `json:"type"` // Changed name to avoid conflict

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
			a.Content = u.Content
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
