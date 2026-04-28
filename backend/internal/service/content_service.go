package service

import (
	"fmt"

	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type ContentService struct {
	postRepo  repository.PostRepositoryInterface
	groupRepo repository.GroupRepositoryInterface
	logger    *logger.Logger
}

func NewContentService(
	postRepo repository.PostRepositoryInterface,
	groupRepo repository.GroupRepositoryInterface,
	logger *logger.Logger,
) *ContentService {
	return &ContentService{
		postRepo:  postRepo,
		groupRepo: groupRepo,
		logger:    logger,
	}
}

func (s *ContentService) CreatePost(userID int, postData domain.CreatePostRequest) (*domain.Post, error) {
	return s.createPost(userID, 0, postData)
}

func (s *ContentService) CreateGroupPost(userID, groupID int, postData domain.CreatePostRequest) (*domain.Post, error) {
	isMember, err := s.groupRepo.IsUserInGroup(groupID, userID)
	if err != nil {
		s.logger.Error("Failed to check group membership", "error", err, "groupID", groupID, "userID", userID)
		return nil, fmt.Errorf("failed to create post")
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of this group")
	}

	return s.createPost(userID, groupID, postData)
}

func (s *ContentService) createPost(userID, groupID int, postData domain.CreatePostRequest) (*domain.Post, error) {
	if postData.PrivacyLevel == "" {
		postData.PrivacyLevel = "public"
	}

	if err := s.validatePost(postData); err != nil {
		return nil, err
	}

	postID, err := s.postRepo.CreatePost(userID, postData.Title, postData.Content, postData.Category, postData.PrivacyLevel, postData.MediaURL, groupID)
	if err != nil {
		s.logger.Error("Failed to create post", "error", err)
		return nil, fmt.Errorf("failed to create post")
	}

	post, err := s.postRepo.GetPostByID(int(postID))
	if err != nil {
		s.logger.Error("Failed to retrieve created post", "error", err)
		return nil, fmt.Errorf("failed to retrieve created post")
	}

	s.logger.Info("Post created successfully", "postID", postID, "userID", userID, "groupID", groupID)
	return post, nil
}

func (s *ContentService) UpdatePost(userID, postID int, data domain.UpdatePostRequest) (*domain.Post, error) {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		s.logger.Error("Failed to retrieve post", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to retrieve post")
	}
	if post.UserID != userID {
		return nil, fmt.Errorf("user is not the owner of this post")
	}

	if data.PrivacyLevel == "" {
		data.PrivacyLevel = "public"
	}

	validateData := domain.CreatePostRequest{
		Title:        data.Title,
		Content:      data.Content,
		Category:     data.Category,
		PrivacyLevel: data.PrivacyLevel,
		MediaURL:     data.MediaURL,
	}
	if err := s.validatePost(validateData); err != nil {
		return nil, err
	}

	if err := s.postRepo.UpdatePost(userID, postID, data.Title, data.Content, data.Category, data.PrivacyLevel, data.MediaURL); err != nil {
		s.logger.Error("Failed to update post", "error", err, "userID", userID, "postID", postID)
		return nil, fmt.Errorf("failed to update post")
	}

	post, err = s.postRepo.GetPostByID(postID)
	if err != nil {
		s.logger.Error("Failed to retrieve updated post", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to retrieve updated post")
	}

	return post, nil
}

func (s *ContentService) DeletePost(userID, postID int) error {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		s.logger.Error("Failed to retrieve post", "error", err, "postID", postID)
		return fmt.Errorf("failed to retrieve post")
	}

	if post.UserID != userID {
		return fmt.Errorf("user is not the owner of this post")
	}

	if err := s.postRepo.DeletePost(userID, postID); err != nil {
		s.logger.Error("Failed to delete post", "error", err, "userID", userID, "postID", postID)
		return fmt.Errorf("failed to delete post")
	}
	return nil
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
