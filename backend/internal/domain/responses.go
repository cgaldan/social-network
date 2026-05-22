package domain

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
	Token   string `json:"token,omitempty"`
}

type PostsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Posts   []Post `json:"posts"`
	HasMore bool   `json:"has_more"`
}

type PostResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Post    *Post  `json:"post,omitempty"`
}

type PostDetailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Post    *Post  `json:"post,omitempty"`
}

type CommentResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Comment *Comment `json:"comment,omitempty"`
}

type CommentsResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Comments []Comment `json:"comments"`
	HasMore  bool      `json:"has_more"`
}

type FollowResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

type FollowRequestsResponse struct {
	Success  bool                   `json:"success"`
	Message  string                 `json:"message,omitempty"`
	Requests []FollowRequestSummary `json:"requests"`
	HasMore  bool                   `json:"has_more"`
}

type MessageResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Msg     *Message `json:"msg,omitempty"`
}

type MessagesResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Messages []Message `json:"messages"`
	HasMore  bool      `json:"has_more"`
}

type ConversationResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message,omitempty"`
	Conversation *Conversation `json:"conversation,omitempty"`
}

type ConversationsResponse struct {
	Success       bool                  `json:"success"`
	Message       string                `json:"message,omitempty"`
	Conversations []ConversationSummary `json:"conversations"`
	HasMore       bool                  `json:"has_more"`
}

type UsersResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message,omitempty"`
	Users   []UserSummary `json:"users"`
	HasMore bool          `json:"has_more"`
}

type UserResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	User    *UserProfile `json:"user,omitempty"`
}

type GroupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Group   *Group `json:"group,omitempty"`
}

type GroupsResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
	Groups  []Group `json:"groups"`
	HasMore bool    `json:"has_more"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

type UploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	URL     string `json:"url,omitempty"`
}

type GroupInvitationResponse struct {
	Success    bool             `json:"success"`
	Message    string           `json:"message,omitempty"`
	Invitation *GroupInvitation `json:"invitation,omitempty"`
}

type GroupJoinRequestResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Request *GroupJoinRequest `json:"request,omitempty"`
}

type GroupEventResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Event   *GroupEvent `json:"event,omitempty"`
}

type GroupEventsResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Events  []GroupEvent `json:"events"`
	HasMore bool         `json:"has_more"`
}

type GroupEventRSVPResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	RSVP    *GroupEventRSVP `json:"rsvp,omitempty"`
}

type NotificationsResponse struct {
	Success       bool           `json:"success"`
	Message       string         `json:"message,omitempty"`
	Notifications []Notification `json:"notifications"`
	HasMore       bool           `json:"has_more"`
}

type NotificationResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message,omitempty"`
	Notification *Notification `json:"notification,omitempty"`
}

type NotificationUnreadCountResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message,omitempty"`
	UnreadCount int    `json:"unread_count"`
}
