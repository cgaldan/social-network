package service

import (
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type Services struct {
	Auth         AuthServiceInterface
	Content      ContentServiceInterface
	Post         PostServiceInterface
	Comment      CommentServiceInterface
	Follow       FollowServiceInterface
	Message      MessageServiceInterface
	Conversation ConversationServiceInterface
}

func NewServices(repos *repository.Repositories, logger *logger.Logger) *Services {
	return &Services{
		Auth:         NewAuthService(repos.User, repos.Session, logger),
		Content:      NewContentService(repos.Post, logger),
		Post:         NewPostService(repos.Post, logger),
		Comment:      NewCommentService(repos.Comment, repos.Post, logger),
		Follow:       NewFollowService(repos.Follow, repos.User, logger),
		Message:      NewMessageService(repos.Message, repos.User, repos.Conversation, logger),
		Conversation: NewConversationService(repos.Conversation, repos.Follow, logger),
	}
}

type AuthServiceInterface interface {
	Register(registrationData domain.RegisterRequest) (*domain.User, string, error)
	Login(loginData domain.LoginRequest) (*domain.User, string, error)
	Logout(sessionID string) error
	ValidateSession(sessionID string) (*domain.User, error)
}

type ContentServiceInterface interface {
	CreatePost(userID int, postData domain.CreatePostRequest) (*domain.Post, error)
}

type PostServiceInterface interface {
	GetPostByID(postID int) (*domain.Post, error)
	ListPosts(category string, limit, offset int) ([]domain.Post, error)
	GetPostsByUserID(userID int, limit, offset int) ([]domain.Post, error)
}

type CommentServiceInterface interface {
	CreateComment(userID int, postID int, commentData domain.CreateCommentRequest) (*domain.Comment, error)
	GetCommentsByPostID(postID int) ([]domain.Comment, error)
	GetCommentsByUserID(userID, limit, offset int) ([]domain.Comment, error)
}

type FollowServiceInterface interface {
	FollowUser(followData domain.FollowRequest) (status string, err error)
}

type MessageServiceInterface interface {
	SendMessage(convID, senderID int, content string) (*domain.Message, error)
}

type ConversationServiceInterface interface {
	CreateDirectConversation(convData domain.DirectConversationRequest) (*domain.Conversation, error)
}
