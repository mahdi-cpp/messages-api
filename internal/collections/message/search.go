package message

import (
	"time"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/iris-tools/search"
)

const MaxLimit = 1000

type SearchOptions struct {
	UserID    uuid.UUID `form:"userId"`
	ChatID    uuid.UUID `form:"chatId"`
	MessageID uuid.UUID `form:"messageId"`
	Content   string    `form:"content"`
	Type      string    `form:"type"`
	IsEdited  *bool     `form:"isEdited"`
	IsPinned  *bool     `form:"isPinned"`
	IsDeleted *bool     `form:"isDeleted"`

	// Date filters
	CreatedAfter  *time.Time `form:"createdAfter"`
	CreatedBefore *time.Time `form:"createdBefore"`
	ActiveAfter   *time.Time `form:"activeAfter"`

	// Sorting
	Sort      string `form:"sort,omitempty"`
	SortOrder string `form:"sortOrder,omitempty"`

	// Pagination
	Page int `form:"page,omitempty"`
	Size int `form:"size,omitempty"`
}

var LessFunks = map[string]search.LessFunction[*Message]{
	"id":        func(a, b *Message) bool { return a.ID.String() < b.ID.String() },
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
		if with.MessageID != uuid.Nil && c.ID != with.MessageID {
			return false
		}
		if with.Content != "" && c.Caption != with.Content {
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
	if with.Sort != "" {
		lessFn := GetLessFunc(with.Sort, with.SortOrder)
		if lessFn != nil {
			search.SortIndexedItems(results, lessFn)
		}
	}

	// Extract final assets
	final := make([]*Message, len(results))
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
		return []*Message{}
	}

	end := start + with.Size
	if end > len(final) {
		end = len(final)
	}
	return final[start:end]
}
