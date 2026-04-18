package service

import (
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type Services struct {
	Auth    AuthServiceInterface
	Content ContentServiceInterface
}

func NewServices(repos *repository.Repositories, logger *logger.Logger) *Services {
	return &Services{
		Auth:    NewAuthService(repos.User, repos.Session, logger),
		Content: NewContentService(repos.Post, logger),
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
