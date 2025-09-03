package models

import "time"

type Room struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"createdBy"`
	IsPrivate   bool      `json:"isPrivate"`
	Members     []string  `json:"members"` // List of user IDs
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
