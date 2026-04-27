package repository

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestNotificationRepository_CreateNotification(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	notificationRepo := repos.Notification

	recipientID, err := userRepo.CreateUser("recipient@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "recipient", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create recipient: %v", err)
	}
	actorID, err := userRepo.CreateUser("actor@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "actor", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create actor: %v", err)
	}

	actorIDInt := int(actorID)
	entityType := "follow_request"
	entityID := 123
	actionURL := "/followers/requests"
	metadata := `{"actor_nickname":"actor"}`

	notification, err := notificationRepo.CreateNotification(&domain.Notification{
		RecipientID: int(recipientID),
		ActorID:     &actorIDInt,
		Type:        "follow.requested",
		Title:       "New follow request",
		Body:        "actor requested to follow you.",
		EntityType:  &entityType,
		EntityID:    &entityID,
		ActionURL:   &actionURL,
		Metadata:    &metadata,
	})
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}

	if notification.ID == 0 {
		t.Error("Expected non-zero notification ID")
	}
	if notification.RecipientID != int(recipientID) {
		t.Errorf("Expected recipient ID %d, got %d", recipientID, notification.RecipientID)
	}
	if notification.ActorID == nil || *notification.ActorID != int(actorID) {
		t.Errorf("Expected actor ID %d, got %v", actorID, notification.ActorID)
	}
	if notification.Type != "follow.requested" {
		t.Errorf("Expected type 'follow.requested', got '%s'", notification.Type)
	}
	if notification.EntityType == nil || *notification.EntityType != entityType {
		t.Errorf("Expected entity type '%s', got %v", entityType, notification.EntityType)
	}
	if notification.EntityID == nil || *notification.EntityID != entityID {
		t.Errorf("Expected entity ID %d, got %v", entityID, notification.EntityID)
	}
	if notification.ActionURL == nil || *notification.ActionURL != actionURL {
		t.Errorf("Expected action URL '%s', got %v", actionURL, notification.ActionURL)
	}
	if notification.Metadata == nil || *notification.Metadata != metadata {
		t.Errorf("Expected metadata '%s', got %v", metadata, notification.Metadata)
	}
	if notification.ReadAt != nil {
		t.Error("Expected unread notification")
	}
	if notification.CreatedAt.IsZero() {
		t.Error("Expected created at timestamp")
	}
}

func TestNotificationRepository_ListNotifications(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	notificationRepo := repos.Notification

	recipientID, _ := userRepo.CreateUser("recipient@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "recipient", "female", "", "", false)
	otherUserID, _ := userRepo.CreateUser("other@example.com", "hashedpass", "Other", "User", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "other", "male", "", "", false)

	notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(recipientID), Type: "follow.requested", Title: "First", Body: "First body"})
	notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(recipientID), Type: "group.invitation.created", Title: "Second", Body: "Second body"})
	notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(recipientID), Type: "group.join_requested", Title: "Third", Body: "Third body"})
	notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(otherUserID), Type: "follow.requested", Title: "Other", Body: "Other body"})

	notifications, err := notificationRepo.ListNotifications(int(recipientID), 2, 0)
	if err != nil {
		t.Fatalf("Failed to list notifications: %v", err)
	}
	if len(notifications) != 2 {
		t.Fatalf("Expected 2 notifications, got %d", len(notifications))
	}
	if notifications[0].Title != "Third" {
		t.Errorf("Expected newest notification first, got '%s'", notifications[0].Title)
	}
	if notifications[1].Title != "Second" {
		t.Errorf("Expected second notification, got '%s'", notifications[1].Title)
	}

	notifications, err = notificationRepo.ListNotifications(int(recipientID), 2, 2)
	if err != nil {
		t.Fatalf("Failed to list notifications with offset: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}
	if notifications[0].Title != "First" {
		t.Errorf("Expected first notification after offset, got '%s'", notifications[0].Title)
	}
}

func TestNotificationRepository_CountUnreadNotifications(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	notificationRepo := repos.Notification

	recipientID, _ := userRepo.CreateUser("recipient@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "recipient", "female", "", "", false)
	otherUserID, _ := userRepo.CreateUser("other@example.com", "hashedpass", "Other", "User", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "other", "male", "", "", false)

	first, _ := notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(recipientID), Type: "follow.requested", Title: "First", Body: "First body"})
	notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(recipientID), Type: "group.invitation.created", Title: "Second", Body: "Second body"})
	notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(otherUserID), Type: "follow.requested", Title: "Other", Body: "Other body"})

	err := notificationRepo.MarkNotificationRead(first.ID, int(recipientID))
	if err != nil {
		t.Fatalf("Failed to mark notification read: %v", err)
	}

	count, err := notificationRepo.CountUnreadNotifications(int(recipientID))
	if err != nil {
		t.Fatalf("Failed to count unread notifications: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 unread notification, got %d", count)
	}
}

func TestNotificationRepository_MarkNotificationRead(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	notificationRepo := repos.Notification

	recipientID, _ := userRepo.CreateUser("recipient@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "recipient", "female", "", "", false)
	otherUserID, _ := userRepo.CreateUser("other@example.com", "hashedpass", "Other", "User", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "other", "male", "", "", false)

	notification, _ := notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(recipientID), Type: "follow.requested", Title: "First", Body: "First body"})

	err := notificationRepo.MarkNotificationRead(notification.ID, int(otherUserID))
	if err == nil {
		t.Fatal("Expected error when marking another user's notification read")
	}

	err = notificationRepo.MarkNotificationRead(notification.ID, int(recipientID))
	if err != nil {
		t.Fatalf("Failed to mark notification read: %v", err)
	}

	notification, err = notificationRepo.GetNotificationByID(notification.ID, int(recipientID))
	if err != nil {
		t.Fatalf("Failed to get notification: %v", err)
	}
	if notification.ReadAt == nil {
		t.Error("Expected notification read timestamp")
	}
}

func TestNotificationRepository_MarkAllNotificationsRead(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	notificationRepo := repos.Notification

	recipientID, _ := userRepo.CreateUser("recipient@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "recipient", "female", "", "", false)
	otherUserID, _ := userRepo.CreateUser("other@example.com", "hashedpass", "Other", "User", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "other", "male", "", "", false)

	notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(recipientID), Type: "follow.requested", Title: "First", Body: "First body"})
	notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(recipientID), Type: "group.invitation.created", Title: "Second", Body: "Second body"})
	notificationRepo.CreateNotification(&domain.Notification{RecipientID: int(otherUserID), Type: "follow.requested", Title: "Other", Body: "Other body"})

	err := notificationRepo.MarkAllNotificationsRead(int(recipientID))
	if err != nil {
		t.Fatalf("Failed to mark all notifications read: %v", err)
	}

	count, err := notificationRepo.CountUnreadNotifications(int(recipientID))
	if err != nil {
		t.Fatalf("Failed to count unread notifications: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 unread notifications, got %d", count)
	}

	count, err = notificationRepo.CountUnreadNotifications(int(otherUserID))
	if err != nil {
		t.Fatalf("Failed to count unread notifications for other user: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected other user's notification to remain unread, got %d", count)
	}
}
