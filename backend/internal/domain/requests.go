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
	SenderID   int `json:"user_id_1"`
	ReceiverID int `json:"user_id_2"`
}

type SendMessageRequest struct {
	ConversationID int    `json:"conversation_id"`
	Content        string `json:"content"`
}
