package domain

import "time"

type RegisterRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Nickname    string    `json:"nickname"`
	Gender      string    `json:"gender"`
	AvatarPath  string    `json:"avatar_path"`
	AboutMe     string    `json:"about_me"`
	IsPublic    bool      `json:"is_public"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type CreatePostRequest struct {
	Title        string `json:"title"`
	Content      string `json:"content"`
	Category     string `json:"category"`
	PrivacyLevel string `json:"privacy_level"`
	MediaURL     string `json:"media_url,omitempty"`
}

type CreateCommentRequest struct {
	Content  string `json:"content"`
	MediaURL string `json:"media_url,omitempty"`
}

type FollowRequest struct {
	FollowerID int    `json:"follower_id"`
	FolloweeID int    `json:"followee_id"`
	Status     string `json:"status"`
}

type DirectConversationRequest struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
}

type SendMessageRequest struct {
	ConversationID int    `json:"conversation_id"`
	Content        string `json:"content"`
}

type CreateGroupRequest struct {
	CreatorID      int    `json:"creator_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	ConversationID int    `json:"conversation_id"`
}

type InviteToGroupRequest struct {
	GroupID   int `json:"group_id"`
	InviterID int `json:"inviter_id"`
	InviteeID int `json:"invitee_id"`
}

type JoinGroupRequest struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}

type CreateGroupEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
}

type GroupEventRSVPRequest struct {
	Response string `json:"response"`
}
