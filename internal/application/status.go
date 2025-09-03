package application

import "time"

type UserStatus string

const (
	StatusOnline       UserStatus = "online"
	StatusOffline      UserStatus = "offline"
	StatusIdle         UserStatus = "idle"
	StatusDoNotDisturb UserStatus = "do_not_disturb"
	StatusInvisible    UserStatus = "invisible" // Special privacy status
)

// PrivacySetting controls who can see the user's status
type PrivacySetting string

const (
	PrivacyEveryone PrivacySetting = "everyone"
	PrivacyContacts PrivacySetting = "contacts"
	PrivacyNobody   PrivacySetting = "nobody"
)

// UserStatusData represents the complete status information stored in Redis
type UserStatusData struct {
	UserID          string         `json:"userID" redis:"userID"`
	Status          UserStatus     `json:"status" redis:"status"`
	LastOnline      time.Time      `json:"last_online" redis:"last_online"`
	CustomStatus    string         `json:"custom_status,omitempty" redis:"custom_status"`
	DevicePlatform  string         `json:"device_platform,omitempty" redis:"device_platform"`
	LastSeenPrivacy PrivacySetting `json:"last_seen_privacy" redis:"last_seen_privacy"`
	OnlinePrivacy   PrivacySetting `json:"online_privacy" redis:"online_privacy"`
}

// HeartbeatRequest is sent periodically by clients to maintain connection
type HeartbeatRequest struct {
	UserID    string `json:"userID"`
	DeviceID  string `json:"device_id"`
	Platform  string `json:"platform"`
	Timestamp int64  `json:"timestamp"` // Unix millis
}

// StatusUpdateRequest is sent when user manually changes status
type StatusUpdateRequest struct {
	UserID       string     `json:"userID"`
	NewStatus    UserStatus `json:"new_status"`
	CustomStatus string     `json:"custom_status,omitempty"`
}

// StatusResponse is sent to clients requesting status information
type StatusResponse struct {
	UserID       string     `json:"userID"`
	Status       UserStatus `json:"status"`
	LastOnline   time.Time  `json:"last_online,omitempty"`
	CustomStatus string     `json:"custom_status,omitempty"`
	Visible      bool       `json:"visible"` // Whether this user's status is visible to requester
}

//Internal Service Communication---------------------------------------

// StatusEvent is published to the message bus when status changes
type StatusEvent struct {
	UserID          string     `json:"userID"`
	PreviousStatus  UserStatus `json:"previous_status"`
	NewStatus       UserStatus `json:"new_status"`
	Timestamp       time.Time  `json:"timestamp"`
	ChangedBySystem bool       `json:"changed_by_system"` // True if triggered by timeout/heartbeat
}

// StatusQuery is used to query user status from other services
type StatusQuery struct {
	RequestingUserID string   `json:"requesting_userID"` // Who is asking
	TargetUserIDs    []string `json:"target_userIDs"`    // Whose status they want to see
}

// BatchStatusResponse contains status for multiple users
type BatchStatusResponse struct {
	Statuses map[string]StatusResponse `json:"statuses"` // userID -> StatusResponse
}

//5. Privacy Configuration------------------------------------------------

// UserPrivacySettings stores all privacy preferences for a user
type UserPrivacySettings struct {
	UserID            string         `json:"userID" redis:"userID"`
	LastSeenPrivacy   PrivacySetting `json:"LastSeenPrivacy" redis:"last_seen_privacy"`
	OnlinePrivacy     PrivacySetting `json:"online_privacy" redis:"online_privacy"`
	ProfilePhoto      PrivacySetting `json:"profile_photo_privacy" redis:"profile_photo_privacy"`
	ForwardedMessages PrivacySetting `json:"forwarded_messages_privacy" redis:"forwarded_messages_privacy"`
	GroupAdd          PrivacySetting `json:"group_add_privacy" redis:"group_add_privacy"`
}

func (u *UserPrivacySettings) RedisKey() string {
	return "user_privacy:" + u.UserID
}
