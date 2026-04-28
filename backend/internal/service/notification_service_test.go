package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestNotificationService_CreateNotificationPushesPersistedNotification(t *testing.T) {
	pusher := &fakeNotificationPusher{}
	services, _ := SetupTestServicesWithEventBus(t, pusher)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "push@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "push",
		Gender:      "female",
		IsPublic:    true,
	})

	notification, err := services.Notification.CreateNotification(domain.CreateNotificationRequest{
		RecipientID: userID,
		Type:        "follow.requested",
		Title:       "Follow request",
		Body:        "A user requested to follow you.",
	})
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}

	if notification.ID == 0 {
		t.Fatal("Expected persisted notification ID")
	}
	if !pusher.called {
		t.Fatal("Expected notification pusher to be called")
	}
	if pusher.userID != userID {
		t.Errorf("Expected push user ID %d, got %d", userID, pusher.userID)
	}
	if pusher.notification == nil {
		t.Fatal("Expected pushed notification")
	}
	if pusher.notification.ID != notification.ID {
		t.Errorf("Expected pushed notification ID %d, got %d", notification.ID, pusher.notification.ID)
	}
	if pusher.notification.CreatedAt.IsZero() {
		t.Error("Expected pushed notification to include created at timestamp")
	}
}

func TestNotificationService_ListNotificationsPaginationDefaults(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "notifications@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "notifications",
		Gender:      "female",
		IsPublic:    true,
	})

	for i := 0; i < 25; i++ {
		_, err := services.Notification.CreateNotification(domain.CreateNotificationRequest{
			RecipientID: userID,
			Type:        "follow.requested",
			Title:       "Follow request",
			Body:        "A user requested to follow you.",
		})
		if err != nil {
			t.Fatalf("Failed to create notification: %v", err)
		}
	}

	notifications, err := services.Notification.ListNotifications(userID, 0, -1)
	if err != nil {
		t.Fatalf("Failed to list notifications: %v", err)
	}
	if len(notifications) != 20 {
		t.Errorf("Expected default limit of 20, got %d", len(notifications))
	}
}

func TestNotificationService_ReadStateBehavior(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "reader@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "reader",
		Gender:      "female",
		IsPublic:    true,
	})

	first, err := services.Notification.CreateNotification(domain.CreateNotificationRequest{
		RecipientID: userID,
		Type:        "follow.requested",
		Title:       "First",
		Body:        "First body",
	})
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}
	_, err = services.Notification.CreateNotification(domain.CreateNotificationRequest{
		RecipientID: userID,
		Type:        "group.invitation.created",
		Title:       "Second",
		Body:        "Second body",
	})
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}

	count, err := services.Notification.CountUnread(userID)
	if err != nil {
		t.Fatalf("Failed to count unread notifications: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected 2 unread notifications, got %d", count)
	}

	err = services.Notification.MarkRead(userID, first.ID)
	if err != nil {
		t.Fatalf("Failed to mark notification read: %v", err)
	}

	count, err = services.Notification.CountUnread(userID)
	if err != nil {
		t.Fatalf("Failed to count unread notifications: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 unread notification, got %d", count)
	}

	err = services.Notification.MarkAllRead(userID)
	if err != nil {
		t.Fatalf("Failed to mark all notifications read: %v", err)
	}

	count, err = services.Notification.CountUnread(userID)
	if err != nil {
		t.Fatalf("Failed to count unread notifications: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 unread notifications, got %d", count)
	}
}

func TestNotificationService_MarkReadRequiresOwner(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "owner@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "owner",
		Gender:      "female",
		IsPublic:    true,
	})
	otherUserID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "other@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "other",
		Gender:      "male",
		IsPublic:    true,
	})

	notification, err := services.Notification.CreateNotification(domain.CreateNotificationRequest{
		RecipientID: userID,
		Type:        "follow.requested",
		Title:       "Follow request",
		Body:        "A user requested to follow you.",
	})
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}

	err = services.Notification.MarkRead(otherUserID, notification.ID)
	if err == nil {
		t.Fatal("Expected error when another user marks notification read")
	}

	count, err := services.Notification.CountUnread(userID)
	if err != nil {
		t.Fatalf("Failed to count unread notifications: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected owner notification to remain unread, got %d", count)
	}
}
