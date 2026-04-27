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
	FollowStatusRejected = "declined"
)

func (s *FollowService) FollowUser(followData domain.FollowRequest) (status string, err error) {
	if followData.FollowerID == followData.FolloweeID || followData.FollowerID == 0 || followData.FolloweeID == 0 {
		return "", errors.New("invalid follower or followee ID")
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

	existingFollow, err := s.followRepo.GetFollowByUsers(followData.FollowerID, followData.FolloweeID)
	if err != nil {
		s.logger.Error("Failed to get existing follow relationship", "error", err, "followerID", followData.FollowerID, "followingID", followData.FolloweeID)
		return "", err
	}

	if existingFollow != nil {
		switch existingFollow.Status {
		case FollowStatusAccepted:
			return "", errors.New("you are already following this user")
		case FollowStatusPending:
			return "", errors.New("there is already a pending follow request for this user")
		case FollowStatusRejected:
			err = s.followRepo.UpdateFollowStatus(existingFollow.ID, followData.Status)
			if err != nil {
				s.logger.Error("Failed to re-open declined follow relationship", "error", err, "followID", existingFollow.ID, "followerID", followData.FollowerID, "followingID", followData.FolloweeID)
				return "", err
			}
			return followData.Status, nil
		default:
			return "", errors.New("follow relationship is in an unknown state")
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

type endFollowAction int

const (
	endFollowUnfollow endFollowAction = iota
	endFollowRemoveFollower
)

func (s *FollowService) endFollowRelationship(followerID, followingID int, action endFollowAction) error {
	existingFollow, err := s.followRepo.GetFollowByUsers(followerID, followingID)
	if err != nil {
		s.logger.Error("Failed to get existing follow relationship", "error", err, "followerID", followerID, "followingID", followingID)
		return err
	}

	if existingFollow == nil {
		if action == endFollowUnfollow {
			return errors.New("there is no follow relationship between you and this user")
		}
		return errors.New("this user is not following you")
	}

	switch existingFollow.Status {
	case FollowStatusAccepted, FollowStatusPending:
		if err = s.followRepo.DeleteFollow(existingFollow.ID); err != nil {
			s.logger.Error("Failed to delete follow relationship", "error", err, "followerID", followerID, "followingID", followingID)
			return err
		}
		return nil
	case FollowStatusRejected:
		if action == endFollowUnfollow {
			return errors.New("you are not following this user")
		}
		return errors.New("this user has rejected your follow request")
	default:
		return errors.New("follow relationship is in an unknown state")
	}
}

func (s *FollowService) UnfollowUser(unfollowData domain.UnfollowRequest) (err error) {
	if unfollowData.FollowerID == unfollowData.FolloweeID || unfollowData.FollowerID == 0 || unfollowData.FolloweeID == 0 {
		s.logger.Error("Invalid follower or followee ID", "followerID", unfollowData.FollowerID, "followeeID", unfollowData.FolloweeID)
		return errors.New("invalid follower or followee ID")
	}

	return s.endFollowRelationship(unfollowData.FollowerID, unfollowData.FolloweeID, endFollowUnfollow)
}

func (s *FollowService) RemoveFollower(removeData domain.RemoveFollowerRequest) (err error) {
	if removeData.FolloweeID == removeData.FollowerID || removeData.FolloweeID == 0 || removeData.FollowerID == 0 {
		s.logger.Error("Invalid follower or followee ID", "followerID", removeData.FollowerID, "followeeID", removeData.FolloweeID)
		return errors.New("invalid follower or followee ID")
	}

	return s.endFollowRelationship(removeData.FollowerID, removeData.FolloweeID, endFollowRemoveFollower)
}

func (s *FollowService) GetFollowByID(followID int) (*domain.Follow, error) {
	return s.followRepo.GetFollowByID(followID)
}
