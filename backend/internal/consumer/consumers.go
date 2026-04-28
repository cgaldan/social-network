package consumer

import (
	"social-network/internal/event"
	"social-network/internal/service"
	"social-network/packages/logger"
)

type Consumers struct {
	NotificationConsumer *NotificationConsumer
}

func NewConsumers(notificationService service.NotificationServiceInterface, eventBus event.EventBus, logger *logger.Logger) *Consumers {
	notificationConsumer := NewNotificationConsumer(notificationService, eventBus, logger)

	return &Consumers{
		NotificationConsumer: notificationConsumer,
	}
}

type NotificationConsumerInterface interface {
	RegisterHandlers() error
}

func (c *Consumers) RegisterHandlers() error {
	if err := c.NotificationConsumer.RegisterHandlers(); err != nil {
		return err
	}
	return nil
}
