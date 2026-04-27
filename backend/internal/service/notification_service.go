package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type NotificationService struct {
	notificationRepo repository.NotificationRepositoryInterface
	logger           *logger.Logger
}

func NewNotificationService(notificationRepo repository.NotificationRepositoryInterface, logger *logger.Logger) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		logger:           logger,
	}
}

func (s *NotificationService) CreateNotification(input domain.CreateNotificationRequest) (*domain.Notification, error) {
	notification, err := s.notificationRepo.CreateNotification(&domain.Notification{
		RecipientID: input.RecipientID,
		ActorID:     input.ActorID,
		Type:        input.Type,
		Title:       input.Title,
		Body:        input.Body,
		EntityType:  input.EntityType,
		EntityID:    input.EntityID,
		ActionURL:   input.ActionURL,
		Metadata:    input.Metadata,
	})
	if err != nil {
		s.logger.Error("Failed to create notification", "error", err, "recipientID", input.RecipientID)
		return nil, fmt.Errorf("failed to create notification")
	}

	return notification, nil
}

func (s *NotificationService) ListNotifications(userID, limit, offset int) ([]domain.Notification, error) {
	limit, offset = s.validateLimitAndOffset(limit, offset)

	notifications, err := s.notificationRepo.ListNotifications(userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list notifications", "error", err, "userID", userID, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to list notifications")
	}

	return notifications, nil
}

func (s *NotificationService) CountUnread(userID int) (int, error) {
	count, err := s.notificationRepo.CountUnreadNotifications(userID)
	if err != nil {
		s.logger.Error("Failed to count unread notifications", "error", err, "userID", userID)
		return 0, fmt.Errorf("failed to count unread notifications")
	}

	return count, nil
}

func (s *NotificationService) MarkRead(userID, notificationID int) error {
	if err := s.notificationRepo.MarkNotificationRead(notificationID, userID); err != nil {
		s.logger.Error("Failed to mark notification read", "error", err, "userID", userID, "notificationID", notificationID)
		return fmt.Errorf("failed to mark notification read")
	}

	return nil
}

func (s *NotificationService) MarkAllRead(userID int) error {
	if err := s.notificationRepo.MarkAllNotificationsRead(userID); err != nil {
		s.logger.Error("Failed to mark all notifications read", "error", err, "userID", userID)
		return fmt.Errorf("failed to mark all notifications read")
	}

	return nil
}

func (s *NotificationService) validateLimitAndOffset(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
