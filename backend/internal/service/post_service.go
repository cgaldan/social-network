package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type PostService struct {
	postRepo  repository.PostRepositoryInterface
	groupRepo repository.GroupRepositoryInterface
	logger    *logger.Logger
}

func NewPostService(postRepo repository.PostRepositoryInterface, groupRepo repository.GroupRepositoryInterface, logger *logger.Logger) *PostService {
	return &PostService{
		postRepo:  postRepo,
		groupRepo: groupRepo,
		logger:    logger,
	}
}

func (s *PostService) GetPostByID(userID, postID int) (*domain.Post, error) {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		s.logger.Error("Failed to get post by ID", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to get post")
	}

	if post.GroupID != 0 {
		isMember, err := s.groupRepo.IsUserInGroup(post.GroupID, userID)
		if err != nil {
			s.logger.Error("Failed to check group membership", "error", err, "groupID", post.GroupID, "userID", userID)
			return nil, fmt.Errorf("failed to get post")
		}
		if !isMember {
			return nil, fmt.Errorf("user is not a member of this group")
		}
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

func (s *PostService) ListPostsByGroupID(userID, groupID, limit, offset int) ([]domain.Post, error) {
	isMember, err := s.groupRepo.IsUserInGroup(groupID, userID)
	if err != nil {
		s.logger.Error("Failed to check group membership", "error", err, "groupID", groupID, "userID", userID)
		return nil, fmt.Errorf("failed to list group posts")
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of this group")
	}

	limit, offset = s.validateLimitAndOffset(limit, offset)

	posts, err := s.postRepo.ListPostsByGroupID(groupID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list group posts", "error", err, "groupID", groupID, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to list group posts")
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
