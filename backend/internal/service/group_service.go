package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type GroupService struct {
	groupRepo   repository.GroupRepositoryInterface
	convService ConversationServiceInterface
	logger      *logger.Logger
}

func NewGroupService(groupRepo repository.GroupRepositoryInterface, convService ConversationServiceInterface, logger *logger.Logger) *GroupService {
	return &GroupService{
		groupRepo:   groupRepo,
		convService: convService,
		logger:      logger,
	}
}

func (s *GroupService) CreateGroup(group *domain.Group) (*domain.Group, error) {
	conversation, err := s.convService.CreateGroupConversation(group.Title, group.CreatorID)
	if err != nil {
		s.logger.Error("Failed to create group conversation", "error", err)
		return nil, fmt.Errorf("failed to create group conversation")
	}
	group.ConversationID = conversation.ID

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

	err = s.groupRepo.AddMember(group.ID, group.CreatorID, "admin")
	if err != nil {
		s.logger.Error("Failed to add member", "error", err)
		return nil, fmt.Errorf("failed to add member")
	}

	return group, nil
}

func (s *GroupService) GetMembersByGroupID(groupID int) ([]domain.GroupMember, error) {
	return s.groupRepo.GetMembersByGroupID(groupID)
}

func (s *GroupService) AddMember(convID, groupID, userID int, role string) error {
	err := s.groupRepo.AddMember(groupID, userID, role)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	err = s.convService.AddConversationParticipant(convID, userID)
	if err != nil {
		return fmt.Errorf("failed to add conversation participant: %w", err)
	}

	return nil
}

func (s *GroupService) RemoveMember(convID, groupID, userID int) error {
	err := s.groupRepo.RemoveMember(groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	err = s.convService.RemoveConversationParticipant(convID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove conversation participant: %w", err)
	}

	return nil
}
