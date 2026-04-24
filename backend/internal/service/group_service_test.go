package service

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestGroupService_CreateGroup(t *testing.T) {
	services := SetupTestServices(t)

	userID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "test@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "testuser",
		Gender:      "male",
		IsPublic:    true,
	})

	userID2 := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "test2@example.com",
		Password:    "password123",
		FirstName:   "Jane",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "testuser2",
		Gender:      "female",
		IsPublic:    true,
	})

	_, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: userID,
		FolloweeID: userID2,
	})
	if err != nil {
		t.Fatalf("Failed to create follow relationship: %v", err)
	}

	conversation := CreateTestDirectChat(t, services, userID, userID2)

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:      userID,
		Title:          "Test Group",
		Description:    "This is a test group.",
		ConversationID: conversation.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	if group.ID == 0 {
		t.Error("Expected non-zero group ID")
	}
}
