package repository

import (
	"testing"
	"time"
)

func TestUserRepository_CreateUser(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, err := repo.CreateUser(
		"test@example.com",
		"hashedpass",
		"John",
		"Doe",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"testuser",
		"male",
		"",
		"",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if userID == 0 {
		t.Error("Expected non-zero user ID")
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, err := repo.CreateUser(
		"test@example.com",
		"hashedpass",
		"John",
		"Doe",
		time.Date(1999, time.June, 1, 0, 0, 0, 0, time.UTC),
		"testuser",
		"male",
		"",
		"",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err := repo.GetUserByID(int(userID))
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user.Nickname != "testuser" {
		t.Errorf("Expected nickname 'testuser', got '%s'", user.Nickname)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestUserRepository_GetUserByIdentifier(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, err := repo.CreateUser(
		"test@example.com",
		"hashedpass",
		"John",
		"Doe",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"testuser",
		"male",
		"https://example.com/avatar.jpg",
		"I am a test user.",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, pass, err := repo.GetUserByIdentifier("testuser")
	if err != nil {
		t.Fatalf("Failed to get user by identifier: %v", err)
	}

	if user.ID != int(userID) {
		t.Errorf("Expected user ID %d, got %d", userID, user.ID)
	}
	if pass != "hashedpass" {
		t.Errorf("Expected password hash 'hashedpass', got '%s'", pass)
	}
}

func TestUserRepository_UpdateLastSeen(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, err := repo.CreateUser(
		"test@example.com",
		"hashedpass",
		"John",
		"Doe",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"testuser",
		"male",
		"",
		"",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	err = repo.UpdateLastSeen(int(userID))
	if err != nil {
		t.Fatalf("Failed to update last seen: %v", err)
	}

	user, _ := repo.GetUserByID(int(userID))
	if user.LastSeen.Before(time.Now().Add(-1 * time.Minute)) {
		t.Error("Expected last seen to be updated to recent time")
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, err := repo.CreateUser(
		"before@example.com",
		"hashedpass",
		"John",
		"Doe",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"before_nickname",
		"male",
		"",
		"",
		true,
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	err = repo.UpdateUser(
		int(userID),
		"after@example.com",
		"Jane",
		"Smith",
		time.Date(1995, time.June, 1, 0, 0, 0, 0, time.UTC),
		"after_nickname",
		"female",
		"/avatar.png",
		"updated about me",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	user, err := repo.GetUserByID(int(userID))
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if user.Email != "after@example.com" {
		t.Errorf("Expected email to be updated, got %s", user.Email)
	}
	if user.Nickname != "after_nickname" {
		t.Errorf("Expected nickname to be updated, got %s", user.Nickname)
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	repos := SetupTestDB(t)
	repo := repos.User

	userID, err := repo.CreateUser(
		"delete@example.com",
		"hashedpass",
		"John",
		"Doe",
		time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		"delete_nickname",
		"male",
		"",
		"",
		true,
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if err := repo.DeleteUser(int(userID)); err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	_, err = repo.GetUserByID(int(userID))
	if err == nil {
		t.Error("Expected deleted user to not exist")
	}
}
