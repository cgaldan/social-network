package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
	"strings"
	"time"
)

type MessageBroadcaster interface {
	BroadcastMessage(message *domain.Message, receiverIDs []int)
}

type MessageService struct {
	messageRepo repository.MessageRepositoryInterface
	userRepo    repository.UserRepositoryInterface
	convRepo    repository.ConversationRepositoryInterface
	followRepo  repository.FollowRepositoryInterface
	broadcaster MessageBroadcaster
	logger      *logger.Logger
}

func NewMessageService(messageRepo repository.MessageRepositoryInterface, userRepo repository.UserRepositoryInterface, convRepo repository.ConversationRepositoryInterface, followRepo repository.FollowRepositoryInterface, broadcaster MessageBroadcaster, logger *logger.Logger) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		convRepo:    convRepo,
		followRepo:  followRepo,
		broadcaster: broadcaster,
		logger:      logger,
	}
}

func (s *MessageService) SendMessage(convID, senderID int, content string) (*domain.Message, error) {
	isMember, err := s.convRepo.IsUserInConversation(convID, senderID)
	if err != nil {
		s.logger.Error("Failed to check conversation membership", "error", err, "convID", convID, "senderID", senderID)
		return nil, fmt.Errorf("failed to send message")
	}
	if !isMember {
		return nil, fmt.Errorf("sender is not part of the conversation")
	}

	conv, err := s.convRepo.GetConversationByID(convID)
	if err != nil {
		s.logger.Error("Failed to load conversation", "error", err, "convID", convID)
		return nil, fmt.Errorf("failed to send message")
	}
	if conv != nil && conv.Type == "private" && s.followRepo != nil {
		participants, err := s.convRepo.GetParticipantIDs(convID)
		if err != nil {
			s.logger.Error("Failed to load participants", "error", err, "convID", convID)
			return nil, fmt.Errorf("failed to send message")
		}
		var other int
		for _, p := range participants {
			if p != senderID {
				other = p
				break
			}
		}
		if other != 0 {
			eligible, err := s.followRepo.EitherUserFollows(senderID, other)
			if err != nil {
				s.logger.Error("Failed to check follow eligibility", "error", err)
				return nil, fmt.Errorf("failed to send message")
			}
			if !eligible {
				return nil, fmt.Errorf("you must follow each other to send messages")
			}
		}
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

	if s.broadcaster != nil {
		participants, err := s.convRepo.GetParticipantIDs(convID)
		if err != nil {
			s.logger.Error("Failed to fetch participants for broadcast", "error", err, "convID", convID)
		} else if len(participants) > 0 {
			s.broadcaster.BroadcastMessage(message, participants)
		}
	}

	s.logger.Info("Message sent successfully", "messageID", messageID, "conversationID", convID, "senderID", senderID)
	return message, nil
}

func (s *MessageService) ListMessages(convID, userID, limit, beforeID int) ([]domain.Message, error) {
	isMember, err := s.convRepo.IsUserInConversation(convID, userID)
	if err != nil {
		s.logger.Error("Failed to check conversation membership", "error", err, "convID", convID, "userID", userID)
		return nil, fmt.Errorf("failed to list messages")
	}
	if !isMember {
		return nil, fmt.Errorf("user is not part of the conversation")
	}

	return s.messageRepo.ListMessagesByConversationID(convID, limit, beforeID)
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
