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
	Posts   []Post `json:"posts,omitempty"`
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

type FollowResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

type MessageResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Msg     *Message `json:"msg,omitempty"`
}

type ConversationResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message,omitempty"`
	Conversation *Conversation `json:"conversation,omitempty"`
}

type GroupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Group   *Group `json:"group,omitempty"`
}

type GroupsResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
	Groups  []Group `json:"groups,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
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
	Events  []GroupEvent `json:"events,omitempty"`
}

type GroupEventRSVPResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	RSVP    *GroupEventRSVP `json:"rsvp,omitempty"`
}
