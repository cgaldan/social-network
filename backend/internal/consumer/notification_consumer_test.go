package consumer

import (
	"social-network/internal/event"
	"testing"
)

func TestNotificationConsumer_PublishedEventsCreateNotifications(t *testing.T) {
	t.Run("follow requested", func(t *testing.T) {
		services, eventBus, pusher := setupNotificationConsumerTest(t)
		actorID := createNotificationConsumerUser(t, services, "follow-actor@example.com", "followactor")
		recipientID := createNotificationConsumerUser(t, services, "follow-recipient@example.com", "followrecipient")

		err := eventBus.Publish(event.NewFollowRequestedEvent(actorID, recipientID, 101, "Follow Actor"))
		if err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}

		assertNotificationCreated(t, services, pusher, recipientID, actorID, "follow.requested", "follow_request", 101)
	})

	t.Run("group invitation created", func(t *testing.T) {
		services, eventBus, pusher := setupNotificationConsumerTest(t)
		actorID := createNotificationConsumerUser(t, services, "invite-actor@example.com", "inviteactor")
		recipientID := createNotificationConsumerUser(t, services, "invite-recipient@example.com", "inviterecipient")

		err := eventBus.Publish(event.NewGroupInvitationCreatedEvent(201, 202, actorID, recipientID, "Book Club", "Invite Actor"))
		if err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}

		assertNotificationCreated(t, services, pusher, recipientID, actorID, "group.invitation.created", "group_invitation", 202)
	})

	t.Run("group join requested", func(t *testing.T) {
		services, eventBus, pusher := setupNotificationConsumerTest(t)
		actorID := createNotificationConsumerUser(t, services, "join-actor@example.com", "joinactor")
		recipientID := createNotificationConsumerUser(t, services, "join-recipient@example.com", "joinrecipient")

		err := eventBus.Publish(event.NewGroupJoinRequestedEvent(301, 302, actorID, recipientID, "Hiking Group", "Join Actor"))
		if err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}

		assertNotificationCreated(t, services, pusher, recipientID, actorID, "group.join_requested", "group_join_request", 302)
	})

	t.Run("group event created", func(t *testing.T) {
		services, eventBus, pusher := setupNotificationConsumerTest(t)
		actorID := createNotificationConsumerUser(t, services, "event-actor@example.com", "eventactor")
		recipientID := createNotificationConsumerUser(t, services, "event-recipient@example.com", "eventrecipient")

		err := eventBus.Publish(event.NewGroupEventCreatedEvent(401, 402, actorID, recipientID, "Chess Club", "Weekly Match", "Event Actor"))
		if err != nil {
			t.Fatalf("Failed to publish event: %v", err)
		}

		assertNotificationCreated(t, services, pusher, recipientID, actorID, "group.event.created", "group_event", 402)
	})
}
