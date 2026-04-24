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

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
