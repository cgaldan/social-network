package consumer

import (
	"social-network/internal/domain"
	"testing"
)

func TestSourceFlows_FollowRequestCreatesNotification(t *testing.T) {
	services, _, pusher := setupNotificationConsumerTest(t)

	followerID := createNotificationConsumerUser(t, services, "follower-source@example.com", "followersource")
	privateUserID := createPrivateNotificationConsumerUser(t, services, "private-source@example.com", "privatesource")

	status, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: followerID,
		FolloweeID: privateUserID,
	})
	if err != nil {
		t.Fatalf("Failed to follow private user: %v", err)
	}
	if status != "pending" {
		t.Fatalf("Expected status pending, got %s", status)
	}

	notifications, err := services.Notification.ListNotifications(privateUserID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list notifications: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}
	notification := notifications[0]
	if notification.Type != "follow.requested" {
		t.Errorf("Expected notification type follow.requested, got %s", notification.Type)
	}
	if notification.RecipientID != privateUserID {
		t.Errorf("Expected recipient ID %d, got %d", privateUserID, notification.RecipientID)
	}
	if notification.ActorID == nil || *notification.ActorID != followerID {
		t.Errorf("Expected actor ID %d, got %v", followerID, notification.ActorID)
	}
	if notification.EntityType == nil || *notification.EntityType != "follow_request" {
		t.Errorf("Expected entity type follow_request, got %v", notification.EntityType)
	}

	if len(pusher.calls) != 1 {
		t.Fatalf("Expected 1 pusher call, got %d", len(pusher.calls))
	}
	if pusher.calls[0].userID != privateUserID {
		t.Errorf("Expected push user ID %d, got %d", privateUserID, pusher.calls[0].userID)
	}
}

func TestSourceFlows_FollowPublicUserDoesNotCreateNotification(t *testing.T) {
	services, _, pusher := setupNotificationConsumerTest(t)

	followerID := createNotificationConsumerUser(t, services, "follower-public@example.com", "followerpublic")
	publicUserID := createNotificationConsumerUser(t, services, "public-source@example.com", "publicsource")

	status, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: followerID,
		FolloweeID: publicUserID,
	})
	if err != nil {
		t.Fatalf("Failed to follow public user: %v", err)
	}
	if status != "accepted" {
		t.Fatalf("Expected status accepted, got %s", status)
	}

	notifications, err := services.Notification.ListNotifications(publicUserID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list notifications: %v", err)
	}
	if len(notifications) != 0 {
		t.Errorf("Expected 0 notifications for accepted public follow, got %d", len(notifications))
	}
	if len(pusher.calls) != 0 {
		t.Errorf("Expected 0 pusher calls for accepted public follow, got %d", len(pusher.calls))
	}
}
