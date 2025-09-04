package telegram

import (
	"time"
)

// Chat represents a basic group chat
type Chat struct {
	ID                  string            `json:"id"`
	Title               string            `json:"title"`
	Type                ChatType          `json:"type"`
	ParticipantsCount   int               `json:"participantsCount"`
	Photo               *ChatPhoto        `json:"photo,omitempty"`
	DefaultBannedRights *ChatBannedRights `json:"defaultBannedRights,omitempty"`
	About               string            `json:"about,omitempty"`
	MigratedTo          *MigratedInfo     `json:"migratedTo,omitempty"`
	CreatedAt           time.Time         `json:"createdAt"`
	UpdatedAt           time.Time         `json:"updatedAt"`
	IsDeactivated       bool              `json:"isDeactivated"`
	IsCreator           bool              `json:"isCreator"`
	IsLeft              bool              `json:"isLeft"`
	IsKicked            bool              `json:"isKicked"`
}

type ChatType string

const (
	ChatTypeBasic      ChatType = "basic"
	ChatTypeSupergroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
)

type MigratedInfo struct {
	ChannelID  string `json:"channelId"`
	AccessHash string `json:"accessHash,omitempty"`
}

// ChatFull represents full info about a basic group
type ChatFull struct {
	Chat                   *Chat              `json:"chat"`
	Participants           []*ChatParticipant `json:"participants"`
	ChatPhoto              *ChatPhoto         `json:"chatPhoto,omitempty"`
	NotifySettings         *NotifySettings    `json:"notifySettings,omitempty"`
	ExportedInvite         *ExportedInvite    `json:"exportedInvite,omitempty"`
	BotInfo                []*BotInfo         `json:"botInfo,omitempty"`
	PinnedMessageID        string             `json:"pinnedMessageId,omitempty"`
	FolderID               int                `json:"folderId,omitempty"`
	CanSetUsername         bool               `json:"canSetUsername"`
	HasScheduled           bool               `json:"hasScheduled"`
	Call                   *Call              `json:"call,omitempty"`
	TTLPeriod              int                `json:"ttlPeriod,omitempty"`
	GroupcallDefaultJoinAs *Peer              `json:"groupcallDefaultJoinAs,omitempty"`
}

type ChatParticipant struct {
	UserID    string          `json:"userId"`
	InviterID string          `json:"inviterId,omitempty"`
	Date      time.Time       `json:"date"`
	Role      ParticipantRole `json:"role"`
}

type NotifySettings struct {
	ShowPreviews bool   `json:"showPreviews"`
	Silent       bool   `json:"silent"`
	MuteUntil    int64  `json:"muteUntil"`
	Sound        string `json:"sound,omitempty"`
}

type ExportedInvite struct {
	URL         string `json:"url"`
	ExpireDate  int64  `json:"expireDate,omitempty"`
	UsageLimit  int    `json:"usageLimit,omitempty"`
	Usage       int    `json:"usage,omitempty"`
	IsPermanent bool   `json:"isPermanent"`
	IsRevoked   bool   `json:"isRevoked"`
}

type BotInfo struct {
	UserID      string        `json:"userId"`
	Description string        `json:"description"`
	Commands    []*BotCommand `json:"commands,omitempty"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type Call struct {
	ID         string `json:"id"`
	AccessHash string `json:"accessHash"`
}

type Peer struct {
	Type PeerType `json:"type"`
	ID   string   `json:"id"`
}

type PeerType string

const (
	PeerTypeUser    PeerType = "user"
	PeerTypeChat    PeerType = "chat"
	PeerTypeChannel PeerType = "channel"
)

// MessageReadParticipant represents a user who read a message
type MessageReadParticipant struct {
	UserID    string    `json:"userId"`
	Date      time.Time `json:"date"`
	MessageID string    `json:"messageId"`
	ChatID    string    `json:"chatId"`
}

// EmojiStickerSet represents a custom emoji stickerset for supergroups
type EmojiStickerSet struct {
	StickerSetID string `json:"stickerSetId"`
	AccessHash   string `json:"accessHash"`
	Title        string `json:"title"`
	ShortName    string `json:"shortName"`
	Count        int    `json:"count"`
	Hash         string `json:"hash"`
	IsEmoji      bool   `json:"isEmoji"`
	BoostLevel   int    `json:"boostLevel"`
}

// Request and Response structs for message methods

type GetMessageReadParticipantsRequest struct {
	ChatID    string `json:"chatId"`
	MessageID string `json:"messageId"`
}

type GetMessageReadParticipantsResponse struct {
	Participants []*MessageReadParticipant `json:"participants"`
	Count        int                       `json:"count"`
	ExpiresAt    time.Time                 `json:"expiresAt"`
}

type AddChatUserRequest struct {
	ChatID   string `json:"chatId"`
	UserID   string `json:"userId"`
	FwdLimit int    `json:"fwdLimit,omitempty"`
}

type AddChatUserResponse struct {
	Chat     *Chat  `json:"chat"`
	InviteID string `json:"inviteId,omitempty"`
}

type CreateChatRequest struct {
	Title       string   `json:"title"`
	UserIDs     []string `json:"userIds"`
	IsBroadcast bool     `json:"isBroadcast,omitempty"`
}

type CreateChatResponse struct {
	Chat *Chat `json:"chat"`
}

type DeleteChatUserRequest struct {
	ChatID string `json:"chatId"`
	UserID string `json:"userId"`
}

type DeleteChatUserResponse struct {
	Chat *Chat `json:"chat"`
}

type EditChatAboutRequest struct {
	ChatID string `json:"chatId"`
	About  string `json:"about"`
}

type EditChatAdminRequest struct {
	ChatID  string `json:"chatId"`
	UserID  string `json:"userId"`
	IsAdmin bool   `json:"isAdmin"`
}

type EditChatDefaultBannedRightsRequest struct {
	ChatID       string            `json:"chatId"`
	BannedRights *ChatBannedRights `json:"bannedRights"`
}

type EditChatPhotoRequest struct {
	ChatID string    `json:"chatId"`
	Photo  ChatPhoto `json:"photo"`
}

type EditChatTitleRequest struct {
	ChatID string `json:"chatId"`
	Title  string `json:"title"`
}

type GetChatsRequest struct {
	ChatIDs []string `json:"chatIds"`
}

type GetChatsResponse struct {
	Chats []*Chat `json:"chats"`
}

type DeleteChatRequest struct {
	ChatID string `json:"chatId"`
}

type GetCommonChatsRequest struct {
	UserID string `json:"userId"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

type GetCommonChatsResponse struct {
	Chats []*Chat `json:"chats"`
	Count int     `json:"count"`
}

type GetFullChatRequest struct {
	ChatID string `json:"chatId"`
}

type GetFullChatResponse struct {
	FullChat *ChatFull `json:"fullChat"`
}

type MigrateChatRequest struct {
	ChatID string `json:"chatId"`
}

type MigrateChatResponse struct {
	Channel *Channel `json:"channel"`
}

type ConvertToGigagroupRequest struct {
	ChannelID string `json:"channelId"`
}

type ConvertToGigagroupResponse struct {
	Channel *Channel `json:"channel"`
}

type SetEmojiStickersRequest struct {
	ChannelID    string `json:"channelId"`
	StickerSetID string `json:"stickerSetId"`
	AccessHash   string `json:"accessHash,omitempty"`
}

type SetEmojiStickersResponse struct {
	Success bool `json:"success"`
}

// Additional helper structs for message methods

type ServiceMessage struct {
	Type      ServiceMessageType `json:"type"`
	ChatID    string             `json:"chatId"`
	UserID    string             `json:"userId,omitempty"`
	InviterID string             `json:"inviterId,omitempty"`
	Date      time.Time          `json:"date"`
	OldTitle  string             `json:"oldTitle,omitempty"`
	NewTitle  string             `json:"newTitle,omitempty"`
}

type ServiceMessageType string

const (
	ServiceMessageUserJoined       ServiceMessageType = "user_joined"
	ServiceMessageUserLeft         ServiceMessageType = "user_left"
	ServiceMessageChatCreated      ServiceMessageType = "chat_created"
	ServiceMessageChatTitleChanged ServiceMessageType = "chat_title_changed"
	ServiceMessageChatPhotoChanged ServiceMessageType = "chat_photo_changed"
	ServiceMessageChatMigrated     ServiceMessageType = "chat_migrated"
)

// ChatMember represents a member in a basic group
type ChatMember struct {
	UserID      string          `json:"userId"`
	InviterID   string          `json:"inviterId,omitempty"`
	Date        time.Time       `json:"date"`
	Role        ParticipantRole `json:"role"`
	CanBeEdited bool            `json:"canBeEdited,omitempty"`
	IsMember    bool            `json:"isMember"`
	IsLeft      bool            `json:"isLeft"`
	IsKicked    bool            `json:"isKicked"`
}

// ChatInvite represents a chat invite
type ChatInvite struct {
	URL               string     `json:"url"`
	Title             string     `json:"title"`
	Photo             *ChatPhoto `json:"photo,omitempty"`
	ParticipantsCount int        `json:"participantsCount"`
	Participants      []*User    `json:"participants,omitempty"`
	IsChannel         bool       `json:"isChannel"`
	IsBroadcast       bool       `json:"isBroadcast"`
	IsPublic          bool       `json:"isPublic"`
	IsMegagroup       bool       `json:"isMegagroup"`
	IsVerified        bool       `json:"isVerified"`
	IsRestricted      bool       `json:"isRestricted"`
	IsRevoked         bool       `json:"isRevoked"`
	ExpireDate        int64      `json:"expireDate,omitempty"`
	UsageLimit        int        `json:"usageLimit,omitempty"`
	Usage             int        `json:"usage,omitempty"`
}

// User represents a Telegram user
type User struct {
	ID           string     `json:"id"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName,omitempty"`
	Username     string     `json:"username,omitempty"`
	Phone        string     `json:"phone,omitempty"`
	Photo        *UserPhoto `json:"photo,omitempty"`
	Status       UserStatus `json:"status,omitempty"`
	IsBot        bool       `json:"isBot"`
	IsVerified   bool       `json:"isVerified"`
	IsRestricted bool       `json:"isRestricted"`
	IsDeleted    bool       `json:"isDeleted"`
}

type UserPhoto struct {
	SmallFileID string `json:"smallFileId"`
	BigFileID   string `json:"bigFileId"`
}

type UserStatus struct {
	Type      UserStatusType `json:"type"`
	WasOnline time.Time      `json:"wasOnline,omitempty"`
	Expires   int64          `json:"expires,omitempty"`
}

type UserStatusType string

const (
	UserStatusOnline    UserStatusType = "online"
	UserStatusOffline   UserStatusType = "offline"
	UserStatusRecently  UserStatusType = "recently"
	UserStatusLastWeek  UserStatusType = "last_week"
	UserStatusLastMonth UserStatusType = "last_month"
	UserStatusEmpty     UserStatusType = "empty"
)
