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
