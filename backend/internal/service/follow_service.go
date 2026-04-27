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

	exists, err := s.followRepo.AcceptedFollowRequestExists(followData.FollowerID, followData.FolloweeID)
	if err != nil {
		s.logger.Error("Failed to check if accepted follow request exists", "error", err, "followerID", followData.FollowerID, "followingID", followData.FolloweeID)
		return "", err
	}
	if exists {
		return "", errors.New("you are already following this user")
	}

	isPublic, err := s.userRepo.GetUserPrivacyByUserID(followData.FolloweeID)
	if err != nil {
		s.logger.Error("Failed to get user privacy settings", "error", err, "userID", followData.FolloweeID)
		return "", err
	}

	if isPublic {
		followData.Status = FollowStatusAccepted
	} else {
		exists, err = s.followRepo.PendingFollowRequestExists(followData.FollowerID, followData.FolloweeID)
		if err != nil {
			s.logger.Error("Failed to check if follow request exists", "error", err, "followerID", followData.FollowerID, "followingID", followData.FolloweeID)
			return "", err
		}
		if !exists {
			followData.Status = FollowStatusPending
		} else {
			return "", errors.New("there is already a pending follow request for this user")
		}
	}

	_, err = s.followRepo.CreateFollow(followData.FollowerID, followData.FolloweeID, followData.Status)
	if err != nil {
		s.logger.Error("Failed to create follow relationship", "error", err, "followerID", followData.FollowerID, "followingID", followData.FolloweeID)
		return "", err
	}

	return followData.Status, nil
}

func (s *FollowService) AcceptFollowRequest(userID int, followRequest *domain.Follow) (err error) {
	if followRequest.FollowingID != userID {
		return errors.New("user is not the following")
	}

	if followRequest.Status != FollowStatusPending {
		return errors.New("follow request is not pending")
	}

	err = s.followRepo.UpdateFollowStatus(followRequest.ID, FollowStatusAccepted)
	if err != nil {
		s.logger.Error("Failed to update follow status", "error", err, "followID", followRequest.ID)
		return err
	}

	return nil
}

func (s *FollowService) DeclineFollowRequest(userID int, followRequest *domain.Follow) (err error) {
	if followRequest.FollowingID != userID {
		return errors.New("user is not the following")
	}

	if followRequest.Status != FollowStatusPending {
		return errors.New("follow request is not pending")
	}

	err = s.followRepo.UpdateFollowStatus(followRequest.ID, FollowStatusRejected)
	if err != nil {
		s.logger.Error("Failed to update follow status", "error", err, "followID", followRequest.ID)
		return err
	}
	return nil
}
