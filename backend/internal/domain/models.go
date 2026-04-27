package domain

import "time"

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	Nickname       string    `json:"nickname"`
	Gender         string    `json:"gender"`
	AvatarPath     string    `json:"avatar_path"`
	AboutMe        string    `json:"about_me"`
	FollowingCount int       `json:"following_count"`
	FollowersCount int       `json:"followers_count"`
	IsOnline       bool      `json:"is_online"`
	IsPublic       bool      `json:"is_public"`
	CreatedAt      time.Time `json:"created_at"`
	LastSeen       time.Time `json:"last_seen"`
}

type Post struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	GroupID      int       `json:"group_id,omitempty"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Category     string    `json:"category"`
	PrivacyLevel string    `json:"privacy_level"` // public, almost_private, private
	MediaURL     string    `json:"media_url,omitempty"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Author       string    `json:"author"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	MediaURL  string    `json:"media_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    string    `json:"author"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Conversation struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type ConversationUser struct {
	ConversationID int       `json:"conversation_id"`
	UserID         int       `json:"user_id"`
	JoinedAt       time.Time `json:"joined_at"`
	LastReadAt     time.Time `json:"last_read_at"`
}

type ConversationView struct {
	ConversationID   int       `json:"conversation_id"`
	ConversationName string    `json:"conversation_name"`
	ConversationType string    `json:"conversation_type"`
	LastMessage      string    `json:"last_message"`
	LastTime         time.Time `json:"last_time"`
	UnreadCount      int       `json:"unread_count"`
}

type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	SenderID       int       `json:"sender_id"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

type Follow struct {
	ID          int       `json:"id"`
	FollowerID  int       `json:"follower_id"`
	FollowingID int       `json:"following_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type PostDetail struct {
	Post     *Post     `json:"post"`
	Comments []Comment `json:"comments"`
}

type UserStatus struct {
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
	IsOnline bool   `json:"is_online"`
}

type Group struct {
	ID             int       `json:"id"`
	CreatorID      int       `json:"creator_id"`
	Title          string    `json:"name"`
	Description    string    `json:"description"`
	ConversationID int       `json:"conversation_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type GroupMember struct {
	ID       int       `json:"id"`
	GroupID  int       `json:"group_id"`
	UserID   int       `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type GroupInvitation struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	InviterID int       `json:"inviter_id"`
	InviteeID int       `json:"invitee_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupJoinRequest struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupEvent struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartsAt    time.Time `json:"starts_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupEventRSVP struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserID    int       `json:"user_id"`
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Notification struct {
	ID          int        `json:"id"`
	RecipientID int        `json:"recipient_id"`
	ActorID     *int       `json:"actor_id,omitempty"`
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	EntityType  *string    `json:"entity_type,omitempty"`
	EntityID    *int       `json:"entity_id,omitempty"`
	ActionURL   *string    `json:"action_url,omitempty"`
	Metadata    *string    `json:"metadata,omitempty"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
