package chat

import (
	"time"

	"github.com/mahdi-cpp/iris-tools/update"
)

type UpdateOptions struct {
	ChatIDs []string `json:"chatIDs,omitempty"` // Asset Ids

	Type        string `json:"type"` // "private", "group", "channel", "supergroup"
	Title       string `json:"title"`
	Username    string `json:"username"` // Unique identifier for public channels/groups
	Description string `json:"description"`
	Avatar      string `json:"avatar"` // Chat profile photo

	CanSetStickerSet *bool `json:"canSetStickerSet"` // Can set sticker set
	IsVerified       *bool `json:"isVerified"`
	IsRestricted     *bool `json:"isRestricted"`
	IsCreator        *bool `json:"isCreator"`
	IsScam           *bool `json:"isScam"`
	IsFake           *bool `json:"isFake"`

	ActiveUsernames       *[]string `json:"users,omitempty"`                 // Full users replacement
	AddActiveUsernames    []string  `json:"AddActiveUsernames,omitempty"`    // Users to add
	RemoveActiveUsernames []string  `json:"removeActiveUsernames,omitempty"` // Users to remove

	Members        *[]Member
	AddMembers     []Member
	RemoveMembers  []Member
	MembersUpdates []update.NestedFieldUpdate[Member]
}

// Key extractors for nested structs
func memberKeyExtractor(m Member) string { return m.UserID }

// Initialize updater
var metadataUpdater = update.NewUpdater[Chat, UpdateOptions]()

func init() {

	// Configure scalar field updates
	metadataUpdater.AddScalarUpdater(func(a *Chat, u UpdateOptions) {
		if u.Type != "" {
			a.Type = u.Type
		}
	})

	metadataUpdater.AddScalarUpdater(func(a *Chat, u UpdateOptions) {
		if u.Username != "" {
			a.Username = u.Username
		}
	})

	metadataUpdater.AddScalarUpdater(func(a *Chat, u UpdateOptions) {
		if u.Description != "" {
			a.Description = u.Description
		}
	})

	metadataUpdater.AddScalarUpdater(func(a *Chat, u UpdateOptions) {
		if u.Avatar != "" {
			a.Avatar = u.Avatar
		}
	})

	// Configure collection operations
	metadataUpdater.AddCollectionUpdater(func(a *Chat, u UpdateOptions) {
		op := update.CollectionUpdateOp[string]{
			FullReplace: u.ActiveUsernames,
			Add:         u.AddActiveUsernames,
			Remove:      u.RemoveActiveUsernames,
		}
		a.ActiveUsernames = update.ApplyCollectionUpdate(a.ActiveUsernames, op)
	})

	// Members (ID-based updates)
	metadataUpdater.AddNestedUpdater(func(p *Chat, u UpdateOptions) {

		op := update.CollectionUpdateOp[Member]{
			FullReplace: u.Members,
			Add:         u.AddMembers,
			Remove:      u.RemoveMembers,
		}
		p.Members = update.ApplyCollectionUpdateByID(
			p.Members,
			op,
			memberKeyExtractor,
		)

		// Apply field-level updates to existing comments
		p.Members = update.ApplyNestedUpdate(
			p.Members,
			u.MembersUpdates,
			memberKeyExtractor,
		)
	})

	// Set modification timestamp
	metadataUpdater.AddPostUpdateHook(func(a *Chat) {
		a.UpdatedAt = time.Now()
	})
}

func Update(p *Chat, update UpdateOptions) *Chat {
	metadataUpdater.Apply(p, update)
	return p
}
