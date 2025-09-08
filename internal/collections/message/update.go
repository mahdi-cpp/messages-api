package message

type UpdateOptions struct {
	MessageIDs  []string    `json:"messageIds,omitempty"`
	ChatID      string      `json:"chatId"`
	UserID      string      `json:"userId"`
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
