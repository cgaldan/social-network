package service

import (
	"fmt"

	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type ContentService struct {
	postRepo repository.PostRepositoryInterface
	logger   *logger.Logger
}

func NewContentService(
	postRepo repository.PostRepositoryInterface,
	logger *logger.Logger,
) *ContentService {
	return &ContentService{
		postRepo: postRepo,
		logger:   logger,
	}
}

func (s *ContentService) CreatePost(userID int, postData domain.CreatePostRequest) (*domain.Post, error) {
	if postData.PrivacyLevel == "" {
		postData.PrivacyLevel = "public"
	}

	if err := s.validatePost(postData); err != nil {
		return nil, err
	}

	postID, err := s.postRepo.CreatePost(userID, postData.Title, postData.Content, postData.Category, postData.PrivacyLevel, postData.MediaURL)
	if err != nil {
		s.logger.Error("Failed to create post", "error", err)
		return nil, fmt.Errorf("failed to create post")
	}

	post, err := s.postRepo.GetPostByID(int(postID))
	if err != nil {
		s.logger.Error("Failed to retrieve created post", "error", err)
		return nil, fmt.Errorf("failed to retrieve created post")
	}

	s.logger.Info("Post created successfully", "postID", postID, "userID", userID)
	return post, nil
}

func (s *ContentService) validatePost(data domain.CreatePostRequest) error {
	if data.Title == "" || len(data.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters")
	}
	if data.Content == "" || len(data.Content) < 10 {
		return fmt.Errorf("content must be at least 10 characters")
	}
	if data.Category == "" {
		return fmt.Errorf("category is required")
	}

	switch data.PrivacyLevel {
	case "public", "almost_private", "private":
		return nil
	default:
		return fmt.Errorf("invalid privacy level")
	}
}
