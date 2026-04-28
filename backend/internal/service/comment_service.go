package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type CommentService struct {
	commentRepo repository.CommentRepositoryInterface
	postRepo    repository.PostRepositoryInterface
	groupRepo   repository.GroupRepositoryInterface
	logger      *logger.Logger
}

func NewCommentService(commentRepo repository.CommentRepositoryInterface, postRepo repository.PostRepositoryInterface, groupRepo repository.GroupRepositoryInterface, logger *logger.Logger) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		groupRepo:   groupRepo,
		logger:      logger,
	}
}

func (s *CommentService) CreateComment(userID int, postID int, commentData domain.CreateCommentRequest) (*domain.Comment, error) {
	if err := s.validateComment(commentData); err != nil {
		return nil, err
	}

	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		s.logger.Error("Failed to get post for comment", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to create comment")
	}

	if post.GroupID != 0 {
		isMember, err := s.groupRepo.IsUserInGroup(post.GroupID, userID)
		if err != nil {
			s.logger.Error("Failed to check group membership", "error", err, "groupID", post.GroupID, "userID", userID)
			return nil, fmt.Errorf("failed to create comment")
		}
		if !isMember {
			return nil, fmt.Errorf("user is not a member of this group")
		}
	}

	commentID, err := s.commentRepo.CreateComment(userID, postID, commentData.Content, commentData.MediaURL)
	if err != nil {
		s.logger.Error("Failed to create comment", "error", err)
		return nil, fmt.Errorf("failed to create comment")
	}

	comment, err := s.commentRepo.GetCommentByID(int(commentID))
	if err != nil {
		s.logger.Error("Failed to retrieve created comment", "error", err, "commentID", commentID)
		return nil, fmt.Errorf("failed to retrieve created comment")
	}

	s.logger.Info("Comment created successfully", "commentID", commentID, "userID", userID, "postID", postID)
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(userID, postID int) ([]domain.Comment, error) {
	exists, err := s.postRepo.PostExists(postID)
	if err != nil {
		s.logger.Error("Failed to check if post exists", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to get comments for post")
	}

	if exists {
		post, err := s.postRepo.GetPostByID(postID)
		if err != nil {
			s.logger.Error("Failed to get post for comments", "error", err, "postID", postID)
			return nil, fmt.Errorf("failed to get comments for post")
		}

		if post.GroupID != 0 {
			isMember, err := s.groupRepo.IsUserInGroup(post.GroupID, userID)
			if err != nil {
				s.logger.Error("Failed to check group membership", "error", err, "groupID", post.GroupID, "userID", userID)
				return nil, fmt.Errorf("failed to get comments for post")
			}
			if !isMember {
				return nil, fmt.Errorf("user is not a member of this group")
			}
		}
	}

	comments, err := s.commentRepo.GetCommentsByPostID(postID)
	if err != nil {
		s.logger.Error("Failed to get comments by post ID", "error", err, "postID", postID)
		return nil, fmt.Errorf("failed to get comments for post")
	}

	return comments, nil
}

func (s *CommentService) GetCommentsByUserID(userID, limit, offset int) ([]domain.Comment, error) {
	comments, err := s.commentRepo.GetCommentsByUserID(userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get comments by user ID", "error", err, "userID", userID)
		return nil, fmt.Errorf("failed to get comments for user")
	}

	return comments, nil
}

func (s *CommentService) UpdateComment(userID, postID, commentID int, data domain.UpdateCommentRequest) (*domain.Comment, error) {
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		s.logger.Error("Failed to get comment for update", "error", err, "commentID", commentID)
		return nil, err
	}
	if comment.UserID != userID {
		return nil, fmt.Errorf("user is not the owner of this comment")
	}

	validateData := domain.CreateCommentRequest{Content: data.Content, MediaURL: data.MediaURL}
	if err := s.validateComment(validateData); err != nil {
		return nil, err
	}

	if comment.PostID != postID {
		return nil, fmt.Errorf("comment is not associated with this post")
	}

	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		s.logger.Error("Failed to get post for comment update", "error", err, "postID", postID)
		return nil, err
	}

	if post.GroupID != 0 {
		isMember, err := s.groupRepo.IsUserInGroup(post.GroupID, userID)
		if err != nil {
			s.logger.Error("Failed to check group membership", "error", err, "groupID", post.GroupID, "userID", userID)
			return nil, err
		}
		if !isMember {
			return nil, fmt.Errorf("user is not a member of this group")
		}
	}

	if err := s.commentRepo.UpdateComment(userID, commentID, data.Content, data.MediaURL); err != nil {
		s.logger.Error("Failed to update comment", "error", err, "userID", userID, "commentID", commentID)
		return nil, err
	}

	updated, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		s.logger.Error("Failed to retrieve updated comment", "error", err, "commentID", commentID)
		return nil, err
	}

	return updated, nil
}

func (s *CommentService) DeleteComment(userID, postID, commentID int) error {
	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		s.logger.Error("Failed to get comment for delete", "error", err, "commentID", commentID)
		return err
	}

	if comment.UserID != userID {
		return fmt.Errorf("user is not the owner of this comment")
	}

	if comment.PostID != postID {
		return fmt.Errorf("comment is not associated with this post")
	}

	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		s.logger.Error("Failed to get post for comment delete", "error", err, "postID", postID)
		return err
	}

	if post.GroupID != 0 {
		isMember, err := s.groupRepo.IsUserInGroup(post.GroupID, userID)
		if err != nil {
			s.logger.Error("Failed to check group membership", "error", err, "groupID", post.GroupID, "userID", userID)
			return err
		}
		if !isMember {
			return fmt.Errorf("user is not a member of this group")
		}
	}

	if err := s.commentRepo.DeleteComment(userID, commentID); err != nil {
		s.logger.Error("Failed to delete comment", "error", err, "userID", userID, "commentID", commentID)
		return err
	}

	return nil
}

func (s *CommentService) validateComment(commentData domain.CreateCommentRequest) error {
	if commentData.Content == "" || len(commentData.Content) < 1 {
		return fmt.Errorf("comment cannot be empty")
	}
	return nil
}
