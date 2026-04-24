package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type GroupService struct {
	groupRepo repository.GroupRepositoryInterface
	logger    *logger.Logger
}

func NewGroupService(groupRepo repository.GroupRepositoryInterface, logger *logger.Logger) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
		logger:    logger,
	}
}

func (s *GroupService) CreateGroup(group *domain.Group) (*domain.Group, error) {
	groupID, err := s.groupRepo.CreateGroup(group)
	if err != nil {
		s.logger.Error("Failed to create group", "error", err)
		return nil, fmt.Errorf("failed to create group")
	}

	group, err = s.groupRepo.GetGroupByID(int(groupID))
	if err != nil {
		s.logger.Error("Failed to get group by ID", "error", err)
		return nil, fmt.Errorf("failed to get group by ID")
	}

	err = s.groupRepo.AddMember(int(groupID), group.CreatorID, "admin")
	if err != nil {
		s.logger.Error("Failed to add member", "error", err)
		return nil, fmt.Errorf("failed to add member")
	}

	return group, nil
}
