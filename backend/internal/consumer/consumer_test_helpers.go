package consumer

import (
	"io"
	"social-network/internal/database"
	"social-network/internal/domain"
	"social-network/internal/event"
	"social-network/internal/repository"
	"social-network/internal/service"
	"social-network/packages/logger"
	"testing"
	"time"
)

type pushedNotification struct {
	userID       int
	notification *domain.Notification
}

type fakeNotificationPusher struct {
	calls []pushedNotification
}

func (p *fakeNotificationPusher) PushNotification(userID int, notification *domain.Notification) {
	p.calls = append(p.calls, pushedNotification{userID: userID, notification: notification})
}

func setupNotificationConsumerTest(t *testing.T) (*service.Services, event.EventBus, *fakeNotificationPusher) {
	t.Helper()

	db, err := database.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	testLogger := logger.NewLogger(io.Discard, logger.InfoLevel)
	eventBus := event.NewInMemoryBus(testLogger)
	pusher := &fakeNotificationPusher{}
	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, testLogger, pusher)
	consumers := NewConsumers(services.Notification, eventBus, testLogger)
	if err := consumers.RegisterHandlers(); err != nil {
		t.Fatalf("Failed to register notification handlers: %v", err)
	}

	return services, eventBus, pusher
}

func createNotificationConsumerUser(t *testing.T, services *service.Services, email, nickname string) int {
	t.Helper()

	user, _, err := services.Auth.Register(domain.RegisterRequest{
		Email:       email,
		Password:    "password123",
		FirstName:   "Test",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    nickname,
		Gender:      "other",
		IsPublic:    true,
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	return user.ID
}

func assertNotificationCreated(t *testing.T, services *service.Services, pusher *fakeNotificationPusher, recipientID, actorID int, notificationType, entityType string, entityID int) {
	t.Helper()

	notifications, err := services.Notification.ListNotifications(recipientID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list notifications: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}

	notification := notifications[0]
	if notification.Type != notificationType {
		t.Errorf("Expected notification type %s, got %s", notificationType, notification.Type)
	}
	if notification.ActorID == nil || *notification.ActorID != actorID {
		t.Errorf("Expected actor ID %d, got %v", actorID, notification.ActorID)
	}
	if notification.EntityType == nil || *notification.EntityType != entityType {
		t.Errorf("Expected entity type %s, got %v", entityType, notification.EntityType)
	}
	if notification.EntityID == nil || *notification.EntityID != entityID {
		t.Errorf("Expected entity ID %d, got %v", entityID, notification.EntityID)
	}
	if len(pusher.calls) != 1 {
		t.Fatalf("Expected 1 pusher call, got %d", len(pusher.calls))
	}
	if pusher.calls[0].userID != recipientID {
		t.Errorf("Expected pusher user ID %d, got %d", recipientID, pusher.calls[0].userID)
	}
	if pusher.calls[0].notification == nil || pusher.calls[0].notification.ID != notification.ID {
		t.Errorf("Expected pusher notification ID %d, got %v", notification.ID, pusher.calls[0].notification)
	}
}
