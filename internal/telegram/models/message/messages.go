package message

import (
	"time"

	"github.com/mahdi-cpp/messages-api/internal/telegram/models"
)

// Message represents a Telegram message
type Message struct {
	ID                string       `json:"id"`
	FromID            string       `json:"fromId,omitempty"`
	PeerID            string       `json:"peerId"`
	Date              time.Time    `json:"date"`
	Message           string       `json:"message,omitempty"`
	Media             *Media       `json:"media,omitempty"`
	ReplyTo           *ReplyInfo   `json:"replyTo,omitempty"`
	Views             int          `json:"views,omitempty"`
	Forwards          int          `json:"forwards,omitempty"`
	Replies           *RepliesInfo `json:"replies,omitempty"`
	EditDate          time.Time    `json:"editDate,omitempty"`
	PostAuthor        string       `json:"postAuthor,omitempty"`
	GroupedID         string       `json:"groupedId,omitempty"`
	IsOut             bool         `json:"isOut"`
	Mentioned         bool         `json:"mentioned"`
	MediaUnread       bool         `json:"mediaUnread"`
	Silent            bool         `json:"silent"`
	Post              bool         `json:"post"`
	FromScheduled     bool         `json:"fromScheduled"`
	Legacy            bool         `json:"legacy"`
	EditHide          bool         `json:"editHide"`
	Pinned            bool         `json:"pinned"`
	NoForwards        bool         `json:"noForwards"`
	ViaBotID          string       `json:"viaBotId,omitempty"`
	Entities          []*Entity    `json:"entities,omitempty"`
	Reactions         []*Reaction  `json:"reactions,omitempty"`
	RestrictionReason string       `json:"restrictionReason,omitempty"`
}

type Media struct {
	Type     MediaType          `json:"type"`
	Photo    *Photo             `json:"photo,omitempty"`
	Document *Document          `json:"document,omitempty"`
	Geo      *telegram.GeoPoint `json:"geo,omitempty"`
	Contact  *Contact           `json:"contact,omitempty"`
	Poll     *Poll              `json:"poll,omitempty"`
	Venue    *Venue             `json:"venue,omitempty"`
	Game     *Game              `json:"game,omitempty"`
	Invoice  *Invoice           `json:"invoice,omitempty"`
	WebPage  *WebPage           `json:"webPage,omitempty"`
}

// FileLocation represents the location of a file in Telegram's distributed storage
type FileLocation struct {
	DCID          int    `json:"dcId"`
	VolumeID      string `json:"volumeId"`
	LocalID       int    `json:"localId"`
	Secret        string `json:"secret"`
	FileReference []byte `json:"fileReference,omitempty"`
}

// DocumentAttribute represents attributes of a document
type DocumentAttribute struct {
	Type                DocumentAttributeType `json:"type"`
	Duration            int                   `json:"duration,omitempty"`            // For audio/video
	Width               int                   `json:"width,omitempty"`               // For image/video
	Height              int                   `json:"height,omitempty"`              // For image/video
	FileName            string                `json:"fileName,omitempty"`            // For document
	Alt                 string                `json:"alt,omitempty"`                 // For image
	Sticker             *StickerAttribute     `json:"sticker,omitempty"`             // For sticker
	Video               *VideoAttribute       `json:"video,omitempty"`               // For video
	Audio               *AudioAttribute       `json:"audio,omitempty"`               // For audio
	Animated            bool                  `json:"animated,omitempty"`            // For animated image
	Voice               bool                  `json:"voice,omitempty"`               // For voice message
	RoundMessage        bool                  `json:"roundMessage,omitempty"`        // For round video/voice
	SupportsStreaming   bool                  `json:"supportsStreaming,omitempty"`   // For video
	VideoStartTimestamp float64               `json:"videoStartTimestamp,omitempty"` // For video
	Performer           string                `json:"performer,omitempty"`           // For audio
	Title               string                `json:"title,omitempty"`               // For audio
	Waveform            []byte                `json:"waveform,omitempty"`            // For voice
}

type DocumentAttributeType string

const (
	DocumentAttributeTypeImageSize   DocumentAttributeType = "imageSize"
	DocumentAttributeTypeAnimated    DocumentAttributeType = "animated"
	DocumentAttributeTypeSticker     DocumentAttributeType = "sticker"
	DocumentAttributeTypeVideo       DocumentAttributeType = "video"
	DocumentAttributeTypeAudio       DocumentAttributeType = "audio"
	DocumentAttributeTypeFilename    DocumentAttributeType = "filename"
	DocumentAttributeTypeHasStickers DocumentAttributeType = "hasStickers"
)

// StickerAttribute represents sticker-specific attributes
type StickerAttribute struct {
	Alt                  string      `json:"alt"`
	StickerSetID         string      `json:"stickerSetId,omitempty"`
	StickerSetAccessHash string      `json:"stickerSetAccessHash,omitempty"`
	Mask                 bool        `json:"mask,omitempty"`
	MaskCoords           *MaskCoords `json:"maskCoords,omitempty"`
}

// VideoAttribute represents video-specific attributes
type VideoAttribute struct {
	RoundMessage      bool `json:"roundMessage,omitempty"`
	SupportsStreaming bool `json:"supportsStreaming,omitempty"`
	Duration          int  `json:"duration,omitempty"`
	Width             int  `json:"width,omitempty"`
	Height            int  `json:"height,omitempty"`
}

// AudioAttribute represents audio-specific attributes
type AudioAttribute struct {
	Voice     bool   `json:"voice,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Title     string `json:"title,omitempty"`
	Performer string `json:"performer,omitempty"`
	Waveform  []byte `json:"waveform,omitempty"`
}

// MaskCoords represents coordinates for mask stickers
type MaskCoords struct {
	N    int     `json:"n"`    // Mask coordinate type
	X    float64 `json:"x"`    // X coordinate
	Y    float64 `json:"y"`    // Y coordinate
	Zoom float64 `json:"zoom"` // Zoom level
}

// PollAnswer represents an answer option in a poll
type PollAnswer struct {
	Text     string `json:"text"`
	Option   []byte `json:"option"` // Binary data representing the option
	Voters   int    `json:"voters,omitempty"`
	Chosen   bool   `json:"chosen,omitempty"`
	Correct  bool   `json:"correct,omitempty"`  // For quiz polls
	OptionID string `json:"optionId,omitempty"` // Unique identifier for the option
}

// LabeledPrice represents a price with label (used for invoices)
type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int64  `json:"amount"` // Amount in the smallest currency unit (e.g., cents)
}

// Additional helper structs for document attributes

// ImageSize represents dimensions of an image
type ImageSize struct {
	Type   string `json:"type"` // Size type (e.g., "s", "m", "l", "x")
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int    `json:"size,omitempty"` // File size in bytes
}

// VideoSize represents dimensions and duration of a video
type VideoSize struct {
	Type         string  `json:"type"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	Size         int     `json:"size,omitempty"`
	Duration     int     `json:"duration,omitempty"` // Duration in seconds
	VideoStartTs float64 `json:"videoStartTs,omitempty"`
}

// AudioSize represents duration and metadata of an audio file
type AudioSize struct {
	Duration  int    `json:"duration"` // Duration in seconds
	Title     string `json:"title,omitempty"`
	Performer string `json:"performer,omitempty"`
}

// StickerSize represents sticker metadata
type StickerSize struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Emoji  string `json:"emoji,omitempty"`
	SetID  string `json:"setId,omitempty"`
}

// Additional constants for document attribute types
const (
	StickerMaskCoordsPoint  = 0
	StickerMaskCoordsVector = 1
	StickerMaskCoordsMask   = 2
)

// File represents a generic file in Telegram
type File struct {
	ID             string `json:"id"`
	AccessHash     string `json:"accessHash,omitempty"`
	Size           int64  `json:"size"`
	DCID           int    `json:"dcId"`
	KeyFingerprint int    `json:"keyFingerprint,omitempty"`
	Key            []byte `json:"key,omitempty"`
	IV             []byte `json:"iv,omitempty"`
	Parts          int    `json:"parts,omitempty"`
	Name           string `json:"name,omitempty"`
	MD5Checksum    string `json:"md5Checksum,omitempty"`
}

// FileRef represents a file reference
type FileRef struct {
	DCID          int    `json:"dcId"`
	FileID        string `json:"fileId"`
	AccessHash    string `json:"accessHash"`
	FileReference []byte `json:"fileReference,omitempty"`
}

type MediaType string

const (
	MediaTypePhoto    MediaType = "photo"
	MediaTypeDocument MediaType = "document"
	MediaTypeGeo      MediaType = "geo"
	MediaTypeContact  MediaType = "contact"
	MediaTypePoll     MediaType = "poll"
	MediaTypeVenue    MediaType = "venue"
	MediaTypeGame     MediaType = "game"
	MediaTypeInvoice  MediaType = "invoice"
	MediaTypeWebPage  MediaType = "webPage"
)

type Photo struct {
	ID         string       `json:"id"`
	AccessHash string       `json:"accessHash"`
	Width      int          `json:"width"`
	Height     int          `json:"height"`
	FileSize   int          `json:"fileSize,omitempty"`
	Date       time.Time    `json:"date"`
	Sizes      []*PhotoSize `json:"sizes"`
	DCID       int          `json:"dcId"`
}

type PhotoSize struct {
	Type     string        `json:"type"`
	Width    int           `json:"width"`
	Height   int           `json:"height"`
	Size     int           `json:"size"`
	Location *FileLocation `json:"location,omitempty"`
}

type Document struct {
	ID         string               `json:"id"`
	AccessHash string               `json:"accessHash"`
	FileSize   int                  `json:"fileSize"`
	MimeType   string               `json:"mimeType"`
	Date       time.Time            `json:"date"`
	DCID       int                  `json:"dcId"`
	Attributes []*DocumentAttribute `json:"attributes,omitempty"`
	Thumb      *PhotoSize           `json:"thumb,omitempty"`
}

type Contact struct {
	PhoneNumber string `json:"phoneNumber"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName,omitempty"`
	UserID      string `json:"userId,omitempty"`
}

type Poll struct {
	ID             string        `json:"id"`
	Question       string        `json:"question"`
	Answers        []*PollAnswer `json:"answers"`
	TotalVoters    int           `json:"totalVoters"`
	Closed         bool          `json:"closed"`
	PublicVoters   bool          `json:"publicVoters"`
	MultipleChoice bool          `json:"multipleChoice"`
	Quiz           bool          `json:"quiz"`
}

type Venue struct {
	Geo       *telegram.GeoPoint `json:"geo"`
	Title     string             `json:"title"`
	Address   string             `json:"address"`
	Provider  string             `json:"provider,omitempty"`
	VenueID   string             `json:"venueId,omitempty"`
	VenueType string             `json:"venueType,omitempty"`
}

type Game struct {
	ID          string    `json:"id"`
	AccessHash  string    `json:"accessHash"`
	ShortName   string    `json:"shortName"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Photo       *Photo    `json:"photo,omitempty"`
	Document    *Document `json:"document,omitempty"`
}

type Invoice struct {
	Currency                 string          `json:"currency"`
	Prices                   []*LabeledPrice `json:"prices"`
	Test                     bool            `json:"test"`
	NameRequested            bool            `json:"nameRequested"`
	PhoneRequested           bool            `json:"phoneRequested"`
	EmailRequested           bool            `json:"emailRequested"`
	ShippingAddressRequested bool            `json:"shippingAddressRequested"`
	Flexible                 bool            `json:"flexible"`
}

type WebPage struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	DisplayURL  string `json:"displayUrl"`
	Type        string `json:"type,omitempty"`
	SiteName    string `json:"siteName,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Photo       *Photo `json:"photo,omitempty"`
	EmbedURL    string `json:"embedUrl,omitempty"`
	EmbedWidth  int    `json:"embedWidth,omitempty"`
	EmbedHeight int    `json:"embedHeight,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	Author      string `json:"author,omitempty"`
}

type ReplyInfo struct {
	MessageID string `json:"messageId"`
	PeerID    string `json:"peerId,omitempty"`
}

type RepliesInfo struct {
	Count          int              `json:"count"`
	Replies        int              `json:"replies"`
	RecentRepliers []*telegram.Peer `json:"recentRepliers,omitempty"`
	ChannelID      string           `json:"channelId,omitempty"`
	MaxID          string           `json:"maxId,omitempty"`
	ReadMaxID      string           `json:"readMaxId,omitempty"`
}

type Entity struct {
	Type     EntityType `json:"type"`
	Offset   int        `json:"offset"`
	Length   int        `json:"length"`
	URL      string     `json:"url,omitempty"`
	UserID   string     `json:"userId,omitempty"`
	Language string     `json:"language,omitempty"`
}

type EntityType string

const (
	EntityTypeMention     EntityType = "mention"
	EntityTypeHashtag     EntityType = "hashtag"
	EntityTypeBotCommand  EntityType = "botCommand"
	EntityTypeURL         EntityType = "url"
	EntityTypeEmail       EntityType = "email"
	EntityTypeBold        EntityType = "bold"
	EntityTypeItalic      EntityType = "italic"
	EntityTypeCode        EntityType = "code"
	EntityTypePre         EntityType = "pre"
	EntityTypeTextURL     EntityType = "textUrl"
	EntityTypeMentionName EntityType = "mentionName"
	EntityTypePhone       EntityType = "phone"
	EntityTypeCashtag     EntityType = "cashtag"
	EntityTypeUnderline   EntityType = "underline"
	EntityTypeStrike      EntityType = "strike"
	EntityTypeBlockquote  EntityType = "blockquote"
	EntityTypeBankCard    EntityType = "bankCard"
)

type Reaction struct {
	Emoticon      string `json:"emoticon,omitempty"`
	CustomEmojiID string `json:"customEmojiId,omitempty"`
	Count         int    `json:"count"`
	Chosen        bool   `json:"chosen,omitempty"`
}

// Request and Response structs for message methods

type DeleteHistoryRequest struct {
	PeerID    string `json:"peerId"`
	MaxID     string `json:"maxId,omitempty"`
	JustClear bool   `json:"justClear,omitempty"`
	Revoke    bool   `json:"revoke,omitempty"`
}

type DeleteHistoryResponse struct {
	DeletedCount int `json:"deletedCount"`
}

type DeleteMessagesRequest struct {
	PeerID     string   `json:"peerId,omitempty"`
	MessageIDs []string `json:"messageIds"`
	Revoke     bool     `json:"revoke,omitempty"`
}

type DeleteMessagesResponse struct {
	DeletedCount int `json:"deletedCount"`
}

type EditMessageRequest struct {
	PeerID    string    `json:"peerId"`
	MessageID string    `json:"messageId"`
	Message   string    `json:"message,omitempty"`
	Media     *Media    `json:"media,omitempty"`
	Entities  []*Entity `json:"entities,omitempty"`
	NoWebPage bool      `json:"noWebPage,omitempty"`
}

type EditMessageResponse struct {
	Message *Message `json:"message"`
}

type ForwardMessagesRequest struct {
	FromPeerID        string   `json:"fromPeerId"`
	ToPeerID          string   `json:"toPeerId"`
	MessageIDs        []string `json:"messageIds"`
	Silent            bool     `json:"silent,omitempty"`
	Background        bool     `json:"background,omitempty"`
	WithMyScore       bool     `json:"withMyScore,omitempty"`
	DropAuthor        bool     `json:"dropAuthor,omitempty"`
	DropMediaCaptions bool     `json:"dropMediaCaptions,omitempty"`
}

type ForwardMessagesResponse struct {
	Messages []*Message `json:"message"`
}

type GetHistoryRequest struct {
	PeerID     string `json:"peerId"`
	OffsetID   string `json:"offsetId,omitempty"`
	OffsetDate int64  `json:"offsetDate,omitempty"`
	AddOffset  int    `json:"addOffset,omitempty"`
	Limit      int    `json:"limit"`
	MaxID      string `json:"maxId,omitempty"`
	MinID      string `json:"minId,omitempty"`
	Hash       string `json:"hash,omitempty"`
}

type GetHistoryResponse struct {
	Messages []*Message `json:"message"`
	Count    int        `json:"count"`
}

type GetSearchResultsPositionsRequest struct {
	PeerID   string          `json:"peerId"`
	Filter   *MessagesFilter `json:"filter"`
	OffsetID string          `json:"offsetId"`
	Limit    int             `json:"limit"`
}

type SearchResultPosition struct {
	MessageID string    `json:"messageId"`
	Date      time.Time `json:"date"`
	Offset    int       `json:"offset"`
}

type GetSearchResultsPositionsResponse struct {
	Positions []*SearchResultPosition `json:"positions"`
	Count     int                     `json:"count"`
}

type GetMessageEditDataRequest struct {
	PeerID    string `json:"peerId"`
	MessageID string `json:"messageId"`
}

type GetMessageEditDataResponse struct {
	CanEdit bool `json:"canEdit"`
	Caption bool `json:"caption,omitempty"`
}

type GetOutboxReadDateRequest struct {
	PeerID    string `json:"peerId"`
	MessageID string `json:"messageId"`
}

type GetOutboxReadDateResponse struct {
	ReadDate time.Time `json:"readDate"`
}

type GetMessagesRequest struct {
	MessageIDs []string `json:"messageIds"`
}

type GetMessagesResponse struct {
	Messages []*Message `json:"message"`
}

type GetMessagesViewsRequest struct {
	PeerID     string   `json:"peerId"`
	MessageIDs []string `json:"messageIds"`
	Increment  bool     `json:"increment,omitempty"`
}

type GetMessagesViewsResponse struct {
	Views []*Views `json:"views"`
}

type Views struct {
	MessageID string `json:"messageId"`
	Views     int    `json:"views"`
}

type GetRecentLocationsRequest struct {
	PeerID string `json:"peerId"`
	Limit  int    `json:"limit"`
	Hash   string `json:"hash,omitempty"`
}

type GetRecentLocationsResponse struct {
	Messages []*Message `json:"message"`
	Count    int        `json:"count"`
}

type GetSearchCountersRequest struct {
	PeerID string          `json:"peerId"`
	Filter *MessagesFilter `json:"filter"`
}

type SearchCounter struct {
	Filter *MessagesFilter `json:"filter"`
	Count  int             `json:"count"`
}

type GetSearchCountersResponse struct {
	Counters []*SearchCounter `json:"counters"`
}

type GetUnreadMentionsRequest struct {
	PeerID    string `json:"peerId"`
	OffsetID  string `json:"offsetId,omitempty"`
	AddOffset int    `json:"addOffset,omitempty"`
	Limit     int    `json:"limit"`
	MaxID     string `json:"maxId,omitempty"`
	MinID     string `json:"minId,omitempty"`
}

type GetUnreadMentionsResponse struct {
	Messages []*Message `json:"message"`
	Count    int        `json:"count"`
}

type ReadHistoryRequest struct {
	PeerID string `json:"peerId"`
	MaxID  string `json:"maxId,omitempty"`
}

type ReadHistoryResponse struct {
	ReadCount int `json:"readCount"`
}

type ReadMentionsRequest struct {
	PeerID string `json:"peerId"`
}

type ReadMentionsResponse struct {
	ReadCount int `json:"readCount"`
}

type ReadMessageContentsRequest struct {
	PeerID     string   `json:"peerId"`
	MessageIDs []string `json:"messageIds"`
}

type ReadMessageContentsResponse struct {
	ReadCount int `json:"readCount"`
}

type ReceivedMessagesRequest struct {
	MaxID string `json:"maxId"`
}

type ReceivedMessagesResponse struct {
	ReceivedCount int `json:"receivedCount"`
}

type SearchRequest struct {
	PeerID    string          `json:"peerId"`
	Query     string          `json:"query,omitempty"`
	Filter    *MessagesFilter `json:"filter,omitempty"`
	MinDate   int64           `json:"minDate,omitempty"`
	MaxDate   int64           `json:"maxDate,omitempty"`
	OffsetID  string          `json:"offsetId,omitempty"`
	AddOffset int             `json:"addOffset,omitempty"`
	Limit     int             `json:"limit"`
	MaxID     string          `json:"maxId,omitempty"`
	MinID     string          `json:"minId,omitempty"`
	Hash      string          `json:"hash,omitempty"`
}

type SearchResponse struct {
	Messages []*Message `json:"message"`
	Count    int        `json:"count"`
}

type GetSearchResultsCalendarRequest struct {
	PeerID     string          `json:"peerId"`
	Filter     *MessagesFilter `json:"filter"`
	OffsetID   string          `json:"offsetId,omitempty"`
	OffsetDate int64           `json:"offsetDate,omitempty"`
}

type SearchResultsCalendar struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
	MinID string    `json:"minId"`
	MaxID string    `json:"maxId"`
}

type GetSearchResultsCalendarResponse struct {
	Results []*SearchResultsCalendar `json:"results"`
	Count   int                      `json:"count"`
}

type SearchGlobalRequest struct {
	Query      string          `json:"query,omitempty"`
	Filter     *MessagesFilter `json:"filter,omitempty"`
	MinDate    int64           `json:"minDate,omitempty"`
	MaxDate    int64           `json:"maxDate,omitempty"`
	OffsetRate int             `json:"offsetRate,omitempty"`
	OffsetPeer *telegram.Peer  `json:"offsetPeer,omitempty"`
	OffsetID   string          `json:"offsetId,omitempty"`
	Limit      int             `json:"limit"`
}

type SearchGlobalResponse struct {
	Messages []*Message       `json:"message"`
	Chats    []*telegram.Chat `json:"chats"`
	Users    []*telegram.User `json:"users"`
	Count    int              `json:"count"`
}

type SearchSentMediaRequest struct {
	Query  string          `json:"query,omitempty"`
	Filter *MessagesFilter `json:"filter,omitempty"`
	Limit  int             `json:"limit"`
}

type SearchSentMediaResponse struct {
	Messages []*Message `json:"message"`
	Count    int        `json:"count"`
}

type SearchPostsRequest struct {
	ChannelID string `json:"channelId"`
	Hashtag   string `json:"hashtag"`
	OffsetID  string `json:"offsetId,omitempty"`
	Limit     int    `json:"limit"`
}

type SearchPostsResponse struct {
	Messages []*Message `json:"message"`
	Count    int        `json:"count"`
}

type SendMediaRequest struct {
	PeerID     string `json:"peerId"`
	Media      *Media `json:"media"`
	Message    string `json:"message,omitempty"`
	ReplyTo    string `json:"replyTo,omitempty"`
	Silent     bool   `json:"silent,omitempty"`
	Background bool   `json:"background,omitempty"`
	ClearDraft bool   `json:"clearDraft,omitempty"`
}

type SendMediaResponse struct {
	Message *Message `json:"message"`
}

type SendMessageRequest struct {
	PeerID     string    `json:"peerId"`
	Message    string    `json:"message"`
	ReplyTo    string    `json:"replyTo,omitempty"`
	Entities   []*Entity `json:"entities,omitempty"`
	Silent     bool      `json:"silent,omitempty"`
	Background bool      `json:"background,omitempty"`
	ClearDraft bool      `json:"clearDraft,omitempty"`
	Noforwards bool      `json:"noforwards,omitempty"`
}

type SendMessageResponse struct {
	Message *Message `json:"message"`
}

type SendMultiMediaRequest struct {
	PeerID     string        `json:"peerId"`
	MultiMedia []*InputMedia `json:"multiMedia"`
	ReplyTo    string        `json:"replyTo,omitempty"`
	Silent     bool          `json:"silent,omitempty"`
	Background bool          `json:"background,omitempty"`
	ClearDraft bool          `json:"clearDraft,omitempty"`
}

type InputMedia struct {
	Type     MediaType `json:"type"`
	Media    *Media    `json:"media"`
	Message  string    `json:"message,omitempty"`
	Entities []*Entity `json:"entities,omitempty"`
}

type SendMultiMediaResponse struct {
	Messages []*Message `json:"message"`
}

type UpdatePinnedMessageRequest struct {
	PeerID    string `json:"peerId"`
	MessageID string `json:"messageId"`
	Silent    bool   `json:"silent,omitempty"`
	Unpin     bool   `json:"unpin,omitempty"`
	PmOneside bool   `json:"pmOneside,omitempty"`
}

type UpdatePinnedMessageResponse struct {
	Message *Message `json:"message"`
}

type UnpinAllMessagesRequest struct {
	PeerID string `json:"peerId"`
}

type UnpinAllMessagesResponse struct {
	UnpinnedCount int `json:"unpinnedCount"`
}

type ToggleNoForwardsRequest struct {
	PeerID  string `json:"peerId"`
	Enabled bool   `json:"enabled"`
}

type ToggleNoForwardsResponse struct {
	Success bool `json:"success"`
}

type SaveDefaultSendAsRequest struct {
	PeerID string         `json:"peerId"`
	SendAs *telegram.Peer `json:"sendAs"`
}

type SaveDefaultSendAsResponse struct {
	Success bool `json:"success"`
}

type GetSendAsRequest struct {
	PeerID string `json:"peerId"`
}

type GetSendAsResponse struct {
	Peers []*telegram.Peer `json:"peers"`
}

// MessagesFilter for search operations
type MessagesFilter struct {
	Type FilterType `json:"type"`
}

type FilterType string

const (
	FilterTypeEmpty      FilterType = "empty"
	FilterTypePhotos     FilterType = "photos"
	FilterTypeVideo      FilterType = "video"
	FilterTypePhotoVideo FilterType = "photoVideo"
	FilterTypeDocument   FilterType = "document"
	FilterTypeURL        FilterType = "url"
	FilterTypeGIF        FilterType = "gif"
	FilterTypeVoice      FilterType = "voice"
	FilterTypeMusic      FilterType = "music"
	FilterTypeChatPhotos FilterType = "chatPhotos"
	FilterTypePhoneCalls FilterType = "phoneCalls"
	FilterTypeRoundVoice FilterType = "roundVoice"
	FilterTypeRoundVideo FilterType = "roundVideo"
	FilterTypeMyMentions FilterType = "myMentions"
	FilterTypeGeo        FilterType = "geo"
	FilterTypeContacts   FilterType = "contacts"
	FilterTypePinned     FilterType = "pinned"
)
