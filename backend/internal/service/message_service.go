package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
	"strings"
	"time"
)

type MessageService struct {
	messageRepo repository.MessageRepositoryInterface
	userRepo    repository.UserRepositoryInterface
	convRepo    repository.ConversationRepositoryInterface
	logger      *logger.Logger
}

func NewMessageService(messageRepo repository.MessageRepositoryInterface, userRepo repository.UserRepositoryInterface, convRepo repository.ConversationRepositoryInterface, logger *logger.Logger) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		convRepo:    convRepo,
		logger:      logger,
	}
}

func (s *MessageService) SendMessage(convID, senderID int, content string) (*domain.Message, error) {
	IsMember, err := s.convRepo.IsUserInConversation(convID, senderID)
	if err != nil {
		s.logger.Error("Failed to check conversation membership", "error", err, "convID", convID, "senderID", senderID)
		return nil, fmt.Errorf("failed to send message")
	}
	if !IsMember {
		return nil, fmt.Errorf("sender is not part of the conversation")
	}

	if err := s.validateMessage(strings.TrimSpace(content)); err != nil {
		return nil, err
	}

	messageID, err := s.messageRepo.CreateMessage(&domain.Message{
		ConversationID: convID,
		SenderID:       senderID,
		Content:        content,
	})
	if err != nil {
		s.logger.Error("Failed to create message", "error", err)
		return nil, fmt.Errorf("failed to send message")
	}

	message := &domain.Message{
		ID:             int(messageID),
		ConversationID: convID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	// s.hub.BroadcastMessage(message, receiverID)

	s.logger.Info("Message sent successfully", "messageID", messageID, "conversationID", convID, "senderID", senderID)
	return message, nil
}

func (s *MessageService) validateMessage(content string) error {
	if len(content) == 0 {
		return fmt.Errorf("message content cannot be empty")
	}
	if len(content) > 1000 {
		return fmt.Errorf("message content cannot exceed 1000 characters")
	}
	return nil
}
