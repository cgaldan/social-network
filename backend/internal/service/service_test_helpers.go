package service

import (
	"io"
	"social-network/internal/database"
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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

func CreateTestUser(t *testing.T, services *Services, nickname, email, password, firstName, lastName string, dateOfBirth time.Time, gender string) int {
	t.Helper()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	authSvc := services.Auth.(*AuthService)

	userID, err := authSvc.userRepo.CreateUser(email, string(hashedPassword), firstName, lastName, dateOfBirth, nickname, gender, "", "", true)
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
