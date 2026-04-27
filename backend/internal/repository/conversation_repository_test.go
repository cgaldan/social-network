package repository

import (
	"testing"
	"time"
)

func TestConversationRepository_CreateDirectConversation(t *testing.T) {
	repos := SetupTestDB(t)
	conversationRepo := repos.Conversation

	userID1, err := repos.User.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user 1: %v", err)
	}

	userID2, err := repos.User.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user 2: %v", err)
	}

	conversation, err := conversationRepo.CreateDirectConversation(int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("Failed to create direct conversation: %v", err)
	}

	if conversation.ID == 0 {
		t.Error("Expected non-zero conversation ID")
	}

	conversation2, err := conversationRepo.GetDirectConversation(int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("Failed to get direct conversation: %v", err)
	}

	if conversation2 == nil {
		t.Fatal("Expected conversation but got nil")
	}

	if conversation2.ID != conversation.ID {
		t.Errorf("Expected conversation ID %d, got %d", conversation.ID, conversation2.ID)
	}

	isUserInConversation, err := conversationRepo.IsUserInConversation(conversation.ID, int(userID1))
	if err != nil {
		t.Fatalf("Failed to check if user is in conversation: %v", err)
	}
	if !isUserInConversation {
		t.Error("Expected user to be in conversation")
	}
}

func TestConversationRepository_CreateGroupConversation(t *testing.T) {
	repos := SetupTestDB(t)
	conversationRepo := repos.Conversation

	userID1, err := repos.User.CreateUser("user1@example.com", "hashedpass1", "User", "One", time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC), "user1", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user 1: %v", err)
	}

	userID2, err := repos.User.CreateUser("user2@example.com", "hashedpass2", "User", "Two", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user2", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user 2: %v", err)
	}

	userID3, err := repos.User.CreateUser("user3@example.com", "hashedpass3", "User", "Three", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "user3", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user 3: %v", err)
	}

	conversation, err := conversationRepo.CreateGroupConversation("Test Group", int(userID1), int(userID2), int(userID3))
	if err != nil {
		t.Fatalf("Failed to create group conversation: %v", err)
	}

	if conversation.ID == 0 {
		t.Error("Expected non-zero conversation ID")
	}

	if conversation.Name != "Test Group" {
		t.Errorf("Expected name 'Test Group', got '%s'", conversation.Name)
	}
	if conversation.Type != "group" {
		t.Errorf("Expected type 'group', got '%s'", conversation.Type)
	}

	isUserInConversation, err := conversationRepo.IsUserInConversation(conversation.ID, int(userID1))
	if err != nil {
		t.Fatalf("Failed to check if user is in conversation: %v", err)
	}

	if !isUserInConversation {
		t.Error("Expected user to be in conversation")
	}

	isUserInConversation, err = conversationRepo.IsUserInConversation(conversation.ID, int(userID2))
	if err != nil {
		t.Fatalf("Failed to check if user is in conversation: %v", err)
	}
	if !isUserInConversation {
		t.Error("Expected user to be in conversation")
	}

	isUserInConversation, err = conversationRepo.IsUserInConversation(conversation.ID, int(userID3))
	if err != nil {
		t.Fatalf("Failed to check if user is in conversation: %v", err)
	}
	if !isUserInConversation {
		t.Error("Expected user to be in conversation")
	}

	err = conversationRepo.RemoveConversationParticipant(conversation.ID, int(userID3))
	if err != nil {
		t.Fatalf("Failed to remove user 3 from conversation: %v", err)
	}

	isUserInConversation, err = conversationRepo.IsUserInConversation(conversation.ID, int(userID3))
	if err != nil {
		t.Fatalf("Failed to check if user is in conversation: %v", err)
	}
	if isUserInConversation {
		t.Error("Expected user to be removed from conversation")
	}
}
