package telegram

import (
	"time"
)

// Channel represents a channel/supergroup/geogroup
type Channel struct {
	ID                  string            `json:"id"`
	Title               string            `json:"title"`
	Username            string            `json:"username,omitempty"`
	Description         string            `json:"description,omitempty"`
	Type                ChannelType       `json:"type"`
	ParticipantsCount   int               `json:"participantsCount"`
	IsBroadcast         bool              `json:"isBroadcast"` // Channel (broadcast)
	IsMegagroup         bool              `json:"isMegagroup"` // Supergroup
	IsGigagroup         bool              `json:"isGigagroup"` // Gigagroup (very large)
	IsGeo               bool              `json:"isGeo"`       // Geogroup
	Location            *GeoPoint         `json:"location,omitempty"`
	Photo               *ChatPhoto        `json:"photo,omitempty"`
	SlowModeSeconds     int               `json:"slowModeSeconds"`
	SignaturesEnabled   bool              `json:"signaturesEnabled"`
	PreHistoryHidden    bool              `json:"preHistoryHidden"`
	ParticipantsHidden  bool              `json:"participantsHidden"`
	DiscussionGroupID   string            `json:"discussionGroupId,omitempty"`
	StickerSet          *StickerSet       `json:"stickerSet,omitempty"`
	CreatedAt           time.Time         `json:"createdAt"`
	UpdatedAt           time.Time         `json:"updatedAt"`
	IsInactive          bool              `json:"isInactive"`
	AccessHash          int64             `json:"accessHash"`
	RestrictionReason   string            `json:"restrictionReason,omitempty"`
	AdminRights         *ChatAdminRights  `json:"adminRights,omitempty"`
	BannedRights        *ChatBannedRights `json:"bannedRights,omitempty"`
	DefaultBannedRights *ChatBannedRights `json:"defaultBannedRights,omitempty"`
}

type ChannelType string

const (
	ChannelTypeBroadcast ChannelType = "broadcast"
	ChannelTypeMegagroup ChannelType = "megagroup"
	ChannelTypeGigagroup ChannelType = "gigagroup"
	ChannelTypeGeo       ChannelType = "geo"
)

// GeoPoint represents a geographical location
type GeoPoint struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Accuracy  float64 `json:"accuracy,omitempty"`
}

// ChatPhoto represents a chat/channel photo
type ChatPhoto struct {
	SmallFileID string `json:"smallFileId"`
	BigFileID   string `json:"bigFileId"`
}

// StickerSet represents a sticker set
type StickerSet struct {
	ID         string `json:"id"`
	AccessHash int64  `json:"accessHash"`
	Title      string `json:"title"`
	ShortName  string `json:"shortName"`
	Count      int    `json:"count"`
	Hash       int32  `json:"hash"`
}

// ChatAdminRights represents admin rights in a channel/supergroup
type ChatAdminRights struct {
	ChangeInfo     bool `json:"changeInfo"`
	PostMessages   bool `json:"postMessages"`
	EditMessages   bool `json:"editMessages"`
	DeleteMessages bool `json:"deleteMessages"`
	BanUsers       bool `json:"banUsers"`
	InviteUsers    bool `json:"inviteUsers"`
	PinMessages    bool `json:"pinMessages"`
	AddAdmins      bool `json:"addAdmins"`
	Anonymous      bool `json:"anonymous"`
	ManageCall     bool `json:"manageCall"`
	Other          bool `json:"other"`
}

// ChatBannedRights represents banned rights in a channel/supergroup
type ChatBannedRights struct {
	ViewMessages bool  `json:"viewMessages"`
	SendMessages bool  `json:"sendMessages"`
	SendMedia    bool  `json:"sendMedia"`
	SendStickers bool  `json:"sendStickers"`
	SendGifs     bool  `json:"sendGifs"`
	SendGames    bool  `json:"sendGames"`
	SendInline   bool  `json:"sendInline"`
	EmbedLinks   bool  `json:"embedLinks"`
	SendPolls    bool  `json:"sendPolls"`
	ChangeInfo   bool  `json:"changeInfo"`
	InviteUsers  bool  `json:"inviteUsers"`
	PinMessages  bool  `json:"pinMessages"`
	UntilDate    int64 `json:"untilDate"`
}

// ChannelParticipant represents a participant in a channel/supergroup
type ChannelParticipant struct {
	UserID       string            `json:"userId"`
	ChannelID    string            `json:"channelId"`
	Date         time.Time         `json:"date"`
	AdminRights  *ChatAdminRights  `json:"adminRights,omitempty"`
	BannedRights *ChatBannedRights `json:"bannedRights,omitempty"`
	Role         ParticipantRole   `json:"role"`
	InviterID    string            `json:"inviterId,omitempty"`
	PromotedBy   string            `json:"promotedBy,omitempty"`
}

type ParticipantRole string

const (
	RoleMember     ParticipantRole = "member"
	RoleAdmin      ParticipantRole = "admin"
	RoleCreator    ParticipantRole = "creator"
	RoleBanned     ParticipantRole = "banned"
	RoleLeft       ParticipantRole = "left"
	RoleRestricted ParticipantRole = "restricted"
)

// AdminLogEvent represents an event in the admin log
type AdminLogEvent struct {
	ID        string         `json:"id"`
	Date      time.Time      `json:"date"`
	UserID    string         `json:"userId"`
	ChannelID string         `json:"channelId"`
	Action    AdminLogAction `json:"action"`
	PrevValue interface{}    `json:"prevValue,omitempty"`
	NewValue  interface{}    `json:"newValue,omitempty"`
}

type AdminLogAction struct {
	Type    string      `json:"type"`
	Details interface{} `json:"details,omitempty"`
}

// MessageLink represents a message link with embed info
type MessageLink struct {
	URL         string `json:"url"`
	HTML        string `json:"html,omitempty"`
	IsEmbedded  bool   `json:"isEmbedded"`
	MessageID   string `json:"messageId"`
	ChannelID   string `json:"channelId"`
	IsPermanent bool   `json:"isPermanent"`
}

// Request and Response structs for API methods

type CreateChannelRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	IsBroadcast bool      `json:"isBroadcast,omitempty"`
	IsMegagroup bool      `json:"isMegagroup,omitempty"`
	IsGeo       bool      `json:"isGeo,omitempty"`
	Location    *GeoPoint `json:"location,omitempty"`
}

type CreateChannelResponse struct {
	Channel *Channel `json:"channel"`
}

type GetInactiveChannelsRequest struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type GetInactiveChannelsResponse struct {
	Channels []*Channel `json:"channels"`
	Count    int        `json:"count"`
}

type DeleteChannelRequest struct {
	ChannelID string `json:"channelId"`
}

type DeleteHistoryRequest struct {
	ChannelID string `json:"channelId"`
	MaxID     string `json:"maxId,omitempty"`
}

type DeleteMessagesRequest struct {
	ChannelID  string   `json:"channelId"`
	MessageIDs []string `json:"messageIds"`
}

type DeleteParticipantHistoryRequest struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
}

type EditAdminRequest struct {
	ChannelID   string           `json:"channelId"`
	UserID      string           `json:"userId"`
	AdminRights *ChatAdminRights `json:"adminRights"`
}

type EditBannedRequest struct {
	ChannelID    string            `json:"channelId"`
	UserID       string            `json:"userId"`
	BannedRights *ChatBannedRights `json:"bannedRights"`
}

type EditCreatorRequest struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
	Password  string `json:"password"`
}

type EditLocationRequest struct {
	ChannelID string   `json:"channelId"`
	Location  GeoPoint `json:"location"`
}

type EditPhotoRequest struct {
	ChannelID string    `json:"channelId"`
	Photo     ChatPhoto `json:"photo"`
}

type EditTitleRequest struct {
	ChannelID string `json:"channelId"`
	Title     string `json:"title"`
}

type ExportMessageLinkRequest struct {
	ChannelID   string `json:"channelId"`
	MessageID   string `json:"messageId"`
	IsPermanent bool   `json:"isPermanent,omitempty"`
}

type ExportMessageLinkResponse struct {
	Link *MessageLink `json:"link"`
}

type GetAdminLogRequest struct {
	ChannelID string   `json:"channelId"`
	Query     string   `json:"query,omitempty"`
	Events    []string `json:"events,omitempty"`
	Admins    []string `json:"admins,omitempty"`
	MaxID     string   `json:"maxId,omitempty"`
	MinID     string   `json:"minId,omitempty"`
	Limit     int      `json:"limit,omitempty"`
}

type GetAdminLogResponse struct {
	Events []*AdminLogEvent `json:"events"`
	Count  int              `json:"count"`
}

type GetAdminPublicChannelsRequest struct {
	ByLocation bool `json:"byLocation,omitempty"`
	CheckLimit bool `json:"checkLimit,omitempty"`
}

type GetAdminPublicChannelsResponse struct {
	Channels []*Channel `json:"channels"`
}

type GetChannelsRequest struct {
	ChannelIDs []string `json:"channelIds"`
}

type GetChannelsResponse struct {
	Channels []*Channel `json:"channels"`
}

type GetFullChannelRequest struct {
	ChannelID string `json:"channelId"`
}

type GetFullChannelResponse struct {
	Channel         *Channel `json:"channel"`
	Participants    int      `json:"participantsCount"`
	Admins          int      `json:"adminsCount"`
	Kicked          int      `json:"kickedCount"`
	Banned          int      `json:"bannedCount"`
	Online          int      `json:"onlineCount"`
	ReadInboxMaxID  string   `json:"readInboxMaxId"`
	ReadOutboxMaxID string   `json:"readOutboxMaxId"`
	UnreadCount     int      `json:"unreadCount"`
}

type GetGroupsForDiscussionResponse struct {
	Groups []*Channel `json:"groups"`
}

type GetMessagesRequest struct {
	ChannelID  string   `json:"channelId"`
	MessageIDs []string `json:"messageIds"`
}

type GetParticipantRequest struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
}

type GetParticipantResponse struct {
	Participant *ChannelParticipant `json:"participant"`
}

type GetParticipantsRequest struct {
	ChannelID string             `json:"channelId"`
	Filter    ParticipantsFilter `json:"filter"`
	Offset    int                `json:"offset,omitempty"`
	Limit     int                `json:"limit,omitempty"`
}

type ParticipantsFilter string

const (
	FilterRecent     ParticipantsFilter = "recent"
	FilterAdmins     ParticipantsFilter = "admins"
	FilterKicked     ParticipantsFilter = "kicked"
	FilterBanned     ParticipantsFilter = "banned"
	FilterBot        ParticipantsFilter = "bots"
	FilterRestricted ParticipantsFilter = "restricted"
)

type GetParticipantsResponse struct {
	Participants []*ChannelParticipant `json:"participants"`
	Count        int                   `json:"count"`
}

type InviteToChannelRequest struct {
	ChannelID string   `json:"channelId"`
	UserIDs   []string `json:"userIds"`
}

type JoinChannelRequest struct {
	ChannelID string `json:"channelId"`
}

type LeaveChannelRequest struct {
	ChannelID string `json:"channelId"`
}

type ReadHistoryRequest struct {
	ChannelID string `json:"channelId"`
	MaxID     string `json:"maxId,omitempty"`
}

type ReadMessageContentsRequest struct {
	ChannelID  string   `json:"channelId"`
	MessageIDs []string `json:"messageIds"`
}

type SetDiscussionGroupRequest struct {
	BroadcastChannelID string `json:"broadcastChannelId"`
	DiscussionGroupID  string `json:"discussionGroupId"`
}

type SetStickersRequest struct {
	ChannelID    string `json:"channelId"`
	StickerSetID string `json:"stickerSetId"`
}

type TogglePreHistoryHiddenRequest struct {
	ChannelID string `json:"channelId"`
	Hidden    bool   `json:"hidden"`
}

type ToggleSignaturesRequest struct {
	ChannelID string `json:"channelId"`
	Enabled   bool   `json:"enabled"`
}

type ToggleSlowModeRequest struct {
	ChannelID       string `json:"channelId"`
	SlowModeSeconds int    `json:"slowModeSeconds"`
}

type ToggleParticipantsHiddenRequest struct {
	ChannelID string `json:"channelId"`
	Hidden    bool   `json:"hidden"`
}
