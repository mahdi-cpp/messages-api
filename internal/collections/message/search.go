package message

import (
	"time"

	"github.com/mahdi-cpp/iris-tools/search"
)

const MaxLimit = 1000

type SearchOptions struct {
	ID        string `json:"id"`
	ChatID    string `json:"chatID"`    // Identifier of the chat
	UserID    string `json:"userID"`    // Sender identifier
	Content   string `json:"content"`   // Text content or caption
	Type      string `json:"type"`      // text, photo, video, audio, document, voice, location, contact, poll, etc.
	IsEdited  *bool  `json:"isEdited"`  // Has the message been edited?
	IsPinned  *bool  `json:"isPinned"`  // Is this message pinned?
	IsDeleted *bool  `json:"isDeleted"` // Soft delete flag

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

var LessFunks = map[string]search.LessFunction[*Message]{
	"id":        func(a, b *Message) bool { return a.ID < b.ID },
	"createdAt": func(a, b *Message) bool { return a.CreatedAt.Before(b.CreatedAt) },
	"updatedAt": func(a, b *Message) bool { return a.UpdatedAt.Before(b.UpdatedAt) },
}

func GetLessFunc(sortBy, sortOrder string) search.LessFunction[*Message] {

	fn, exists := LessFunks[sortBy]
	if !exists {
		return nil
	}

	if sortOrder == "end" {
		return func(a, b *Message) bool { return !fn(a, b) }
	}
	return fn
}

func BuildMessageCriteria(with *SearchOptions) search.Criteria[*Message] {

	return func(c *Message) bool {

		// ID filter
		if with.ID != "" && c.ID != with.ID {
			return false
		}

		// Boolean flags
		if with.IsEdited != nil && c.IsEdited != *with.IsEdited {
			return false
		}
		if with.IsPinned != nil && c.IsPinned != *with.IsPinned {
			return false
		}
		if with.IsDeleted != nil && c.IsDeleted != *with.IsDeleted {
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

func Search(chats []*Message, with *SearchOptions) []*Message {

	// Build criteria
	criteria := BuildMessageCriteria(with)

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
	final := make([]*Message, len(results))
	for i, item := range results {
		final[i] = item.Value
	}

	if with.Limit == 0 { // if not set default is MAX_LIMIT
		with.Limit = MaxLimit
	}

	// Apply pagination
	start := with.Offset

	// Check if the start index is out of bounds. If so, return an empty slice.
	if start >= len(final) {
		return []*Message{}
	}

	end := start + with.Limit
	if end > len(final) {
		end = len(final)
	}
	return final[start:end]
}
