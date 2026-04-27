package service

import (
	"errors"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type FollowService struct {
	followRepo repository.FollowRepositoryInterface
	userRepo   repository.UserRepositoryInterface
	logger     *logger.Logger
}

func NewFollowService(followRepo repository.FollowRepositoryInterface, userRepo repository.UserRepositoryInterface, logger *logger.Logger) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
		logger:     logger,
	}
}

const (
	FollowStatusPending  = "pending"
	FollowStatusAccepted = "accepted"
	FollowStatusRejected = "rejected"
)

func (s *FollowService) FollowUser(followData domain.FollowRequest) (status string, err error) {
	if followData.FollowerID == followData.FolloweeID || followData.FollowerID == 0 || followData.FolloweeID == 0 {
		return "", errors.New("invalid follower or followee ID")
	}

	exists, err := s.followRepo.FollowExists(followData.FollowerID, followData.FolloweeID)
	if err != nil {
		s.logger.Error("Failed to check if follow relationship exists", "error", err, "followerID", followData.FollowerID, "followingID", followData.FolloweeID)
		return "", err
	}
	if exists {
		return "", errors.New("follow relationship already exists")
	}

	isPublic, err := s.userRepo.GetUserPrivacyByUserID(followData.FolloweeID)
	if err != nil {
		s.logger.Error("Failed to get user privacy settings", "error", err, "userID", followData.FolloweeID)
		return "", err
	}

	if isPublic {
		followData.Status = FollowStatusAccepted
	} else {
		followData.Status = FollowStatusPending
	}

	_, err = s.followRepo.CreateFollow(followData.FollowerID, followData.FolloweeID, followData.Status)
	if err != nil {
		s.logger.Error("Failed to create follow relationship", "error", err, "followerID", followData.FollowerID, "followingID", followData.FolloweeID)
		return "", err
	}

	return followData.Status, nil
}
