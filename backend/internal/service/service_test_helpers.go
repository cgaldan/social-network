package service

import (
	"io"
	"social-network/internal/database"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type fakeNotificationPusher struct {
	called       bool
	userID       int
	notification *domain.Notification
}

func (p *fakeNotificationPusher) PushNotification(userID int, notification *domain.Notification) {
	p.called = true
	p.userID = userID
	p.notification = notification
}

func SetupTestServicesWithNotificationPusher(t *testing.T, pusher NotificationPusher) *Services {
	t.Helper()

	db, err := database.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repos := repository.NewRepositories(db)
	testLogger := logger.NewLogger(io.Discard, logger.InfoLevel)

	return NewServices(repos, testLogger, pusher)
}

func SetupTestServices(t *testing.T) *Services {
	t.Helper()

	db, err := database.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repos := repository.NewRepositories(db)

	testLogger := logger.NewLogger(io.Discard, logger.InfoLevel)

	// hub := websocket.NewHub(testLogger, repos.User)

	services := NewServices(repos, testLogger)

	return services
}

func CreateTestUser(t *testing.T, services *Services, req domain.RegisterRequest) int {
	t.Helper()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	authSvc := services.Auth.(*AuthService)

	userID, err := authSvc.userRepo.CreateUser(req.Email, string(hashedPassword), req.FirstName, req.LastName, req.DateOfBirth, req.Nickname, req.Gender, req.AvatarPath, req.AboutMe, req.IsPublic)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return int(userID)
}

func CreateTestPost(t *testing.T, services *Services, userID int, postData domain.CreatePostRequest) *domain.Post {
	t.Helper()

	post, err := services.Content.CreatePost(userID, postData)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	return post
}

func CreateTestComment(t *testing.T, services *Services, userID int, postID int, commentData domain.CreateCommentRequest) *domain.Comment {
	t.Helper()

	comment, err := services.Comment.CreateComment(userID, postID, commentData)
	if err != nil {
		t.Fatalf("Failed to create test comment: %v", err)
	}

	return comment
}

func CreateTestDirectChat(t *testing.T, services *Services, userID1, userID2 int) *domain.Conversation {
	t.Helper()

	conversation, err := services.Conversation.CreateDirectConversation(domain.DirectConversationRequest{
		SenderID:   userID1,
		ReceiverID: userID2,
	})
	if err != nil {
		t.Fatalf("Failed to create test direct chat: %v", err)
	}

	if conversation == nil {
		t.Fatalf("CreateDirectConversation returned nil conversation with no error")
	}

	return conversation
}

func CreateTestGroupConversation(t *testing.T, services *Services, name string, initialUserIDs ...int) *domain.Conversation {
	t.Helper()

	conversation, err := services.Conversation.CreateGroupConversation(name, initialUserIDs...)
	if err != nil {
		t.Fatalf("Failed to create test group conversation: %v", err)
	}

	return conversation
}
