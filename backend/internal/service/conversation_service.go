package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
	"strings"
)

type ConversationService struct {
	convRepo   repository.ConversationRepositoryInterface
	followRepo repository.FollowRepositoryInterface
	logger     *logger.Logger
}

func NewConversationService(convRepo repository.ConversationRepositoryInterface, followRepo repository.FollowRepositoryInterface, logger *logger.Logger) *ConversationService {
	return &ConversationService{
		convRepo:   convRepo,
		followRepo: followRepo,
		logger:     logger,
	}
}

func (s *ConversationService) CreateDirectConversation(convData domain.DirectConversationRequest) (*domain.Conversation, error) {
	userID1 := convData.SenderID
	userID2 := convData.ReceiverID

	if userID1 == userID2 {
		return nil, fmt.Errorf("cannot create conversation with oneself")
	}

	canChat, err := s.followRepo.EitherUserFollows(userID1, userID2)
	if err != nil {
		s.logger.Error("Failed to check chat permissions", "error", err, "userID1", userID1, "userID2", userID2)
		return nil, fmt.Errorf("failed to create conversation")
	}
	if !canChat {
		return nil, fmt.Errorf("users must follow each other to start a conversation")
	}

	existingConv, err := s.convRepo.GetDirectConversation(userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	if existingConv != nil {
		return existingConv, nil
	}

	newConv, err := s.convRepo.CreateDirectConversation(userID1, userID2)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			conv, err := s.convRepo.GetDirectConversation(userID1, userID2)
			if conv == nil && err == nil {
				return nil, fmt.Errorf("conversation creation failed with UNIQUE constraint but conversation could not be retrieved")
			}
			return conv, err
		}
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return newConv, nil
}

func (s *ConversationService) CreateGroupConversation(name string, initialUserIDs ...int) (*domain.Conversation, error) {
	conversation, err := s.convRepo.CreateGroupConversation(name, initialUserIDs...)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conversation, nil
}

func (s *ConversationService) AddConversationParticipant(convID, userID int) error {
	err := s.convRepo.AddConversationParticipant(convID, userID)
	if err != nil {
		return fmt.Errorf("failed to add conversation participant: %w", err)
	}
	return nil
}

func (s *ConversationService) RemoveConversationParticipant(convID, userID int) error {
	err := s.convRepo.RemoveConversationParticipant(convID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove conversation participant: %w", err)
	}
	return nil
}
