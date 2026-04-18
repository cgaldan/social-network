package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type PostService struct {
	postRepo    repository.PostRepositoryInterface
	commentRepo repository.CommentRepositoryInterface
	logger      *logger.Logger
}

func NewPostService(postRepo repository.PostRepositoryInterface, logger *logger.Logger) *PostService {
	return &PostService{
		postRepo: postRepo,
		logger:   logger,
	}
}

func (s *PostService) GetPostByID(postID int) (*domain.Post, error) {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		s.logger.Error("Failed to get post by ID", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to get post")
	}

	return post, nil
}

func (s *PostService) ListPosts(category string, limit, offset int) ([]domain.Post, error) {
	limit, offset = s.validateLimitAndOffset(limit, offset)

	posts, err := s.postRepo.ListPosts(category, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list posts", "error", err, "category", category, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to list posts")
	}

	return posts, nil
}

func (s *PostService) GetPostsByUserID(userID int, limit, offset int) ([]domain.Post, error) {
	limit, offset = s.validateLimitAndOffset(limit, offset)

	posts, err := s.postRepo.GetPostsByUserID(userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get posts by user ID", "error", err, "userID", userID, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to get posts by user ID")
	}

	return posts, nil
}

func (s *PostService) validateLimitAndOffset(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
