package consumer

import (
	"encoding/json"
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/event"
	"social-network/internal/service"
	"social-network/packages/logger"
)

type NotificationConsumer struct {
	notificationService service.NotificationServiceInterface
	eventBus            event.EventBus
	logger              *logger.Logger
}

func NewNotificationConsumer(notificationService service.NotificationServiceInterface, eventBus event.EventBus, logger *logger.Logger) *NotificationConsumer {
	return &NotificationConsumer{
		notificationService: notificationService,
		eventBus:            eventBus,
		logger:              logger,
	}
}

func (c *NotificationConsumer) RegisterHandlers() error {
	if err := c.eventBus.Subscribe(event.FollowRequestedEvent, c.handleFollowRequested); err != nil {
		return fmt.Errorf("failed to subscribe to follow requested events: %w", err)
	}
	if err := c.eventBus.Subscribe(event.GroupInvitationCreatedEvent, c.handleGroupInvitationCreated); err != nil {
		return fmt.Errorf("failed to subscribe to group invitation created events: %w", err)
	}
	if err := c.eventBus.Subscribe(event.GroupJoinRequestedEvent, c.handleGroupJoinRequested); err != nil {
		return fmt.Errorf("failed to subscribe to group join requested events: %w", err)
	}
	if err := c.eventBus.Subscribe(event.GroupEventCreatedEvent, c.handleGroupEventCreated); err != nil {
		return fmt.Errorf("failed to subscribe to group event created events: %w", err)
	}

	c.logger.Info("Notification consumer registered successfully")
	return nil
}

func (c *NotificationConsumer) handleFollowRequested(e event.Event) error {
	evt, ok := e.(*event.FollowRequestedEventData)
	if !ok {
		return fmt.Errorf("invalid event type for follow requested handler")
	}

	actorID := evt.FollowerID
	entityType := "follow_request"
	entityID := evt.FollowID
	actionURL := "/followers/requests"

	_, err := c.notificationService.CreateNotification(domain.CreateNotificationRequest{
		RecipientID: evt.FolloweeID,
		ActorID:     &actorID,
		Type:        string(event.FollowRequestedEvent),
		Title:       "New follow request",
		Body:        "Someone requested to follow you.",
		EntityType:  &entityType,
		EntityID:    &entityID,
		ActionURL:   &actionURL,
		Metadata:    metadata(map[string]string{"actor_name": evt.ActorName}),
	})
	return err
}

func (c *NotificationConsumer) handleGroupInvitationCreated(e event.Event) error {
	evt, ok := e.(*event.GroupInvitationCreatedEventData)
	if !ok {
		return fmt.Errorf("invalid event type for group invitation created handler")
	}

	actorID := evt.InviterID
	entityType := "group_invitation"
	entityID := evt.InvitationID
	actionURL := "/groups/invitations"

	_, err := c.notificationService.CreateNotification(domain.CreateNotificationRequest{
		RecipientID: evt.InviteeID,
		ActorID:     &actorID,
		Type:        string(event.GroupInvitationCreatedEvent),
		Title:       "New group invitation",
		Body:        fmt.Sprintf("You were invited to join %s.", displayName(evt.GroupTitle, "a group")),
		EntityType:  &entityType,
		EntityID:    &entityID,
		ActionURL:   &actionURL,
		Metadata:    metadata(map[string]string{"group_title": evt.GroupTitle, "actor_name": evt.ActorName}),
	})
	return err
}

func (c *NotificationConsumer) handleGroupJoinRequested(e event.Event) error {
	evt, ok := e.(*event.GroupJoinRequestedEventData)
	if !ok {
		return fmt.Errorf("invalid event type for group join requested handler")
	}

	actorID := evt.RequesterID
	entityType := "group_join_request"
	entityID := evt.RequestID
	actionURL := "/groups/join"

	_, err := c.notificationService.CreateNotification(domain.CreateNotificationRequest{
		RecipientID: evt.RecipientID,
		ActorID:     &actorID,
		Type:        string(event.GroupJoinRequestedEvent),
		Title:       "New group join request",
		Body:        fmt.Sprintf("Someone requested to join %s.", displayName(evt.GroupTitle, "your group")),
		EntityType:  &entityType,
		EntityID:    &entityID,
		ActionURL:   &actionURL,
		Metadata:    metadata(map[string]string{"group_title": evt.GroupTitle, "actor_name": evt.ActorName}),
	})
	return err
}

func (c *NotificationConsumer) handleGroupEventCreated(e event.Event) error {
	evt, ok := e.(*event.GroupEventCreatedEventData)
	if !ok {
		return fmt.Errorf("invalid event type for group event created handler")
	}

	actorID := evt.CreatorID
	entityType := "group_event"
	entityID := evt.EventID
	actionURL := fmt.Sprintf("/groups/%d/events", evt.GroupID)

	_, err := c.notificationService.CreateNotification(domain.CreateNotificationRequest{
		RecipientID: evt.RecipientID,
		ActorID:     &actorID,
		Type:        string(event.GroupEventCreatedEvent),
		Title:       "New group event",
		Body:        fmt.Sprintf("%s was created in %s.", displayName(evt.EventTitle, "A new event"), displayName(evt.GroupTitle, "your group")),
		EntityType:  &entityType,
		EntityID:    &entityID,
		ActionURL:   &actionURL,
		Metadata:    metadata(map[string]string{"group_title": evt.GroupTitle, "event_title": evt.EventTitle, "actor_name": evt.ActorName}),
	})
	return err
}

func metadata(values map[string]string) *string {
	filtered := make(map[string]string)
	for key, value := range values {
		if value != "" {
			filtered[key] = value
		}
	}
	if len(filtered) == 0 {
		return nil
	}

	data, err := json.Marshal(filtered)
	if err != nil {
		return nil
	}

	value := string(data)
	return &value
}

func displayName(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
