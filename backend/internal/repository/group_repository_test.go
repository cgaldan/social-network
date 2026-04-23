package repository

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestGroupRepository_CreateGroup(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	groupRepo := repos.Group
	convRepo := repos.Conversation

	userID, err := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID2, err := userRepo.CreateUser("test2@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser2", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	conversation, err := convRepo.CreateDirectConversation(int(userID), int(userID2))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	groupID, err := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID),
		Title:          "Test Group",
		Description:    "This is a test group.",
		ConversationID: conversation.ID,
	})

	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	if groupID == 0 {
		t.Error("Expected non-zero group ID")
	}
}

func TestGroupRepository_AddMember(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	groupRepo := repos.Group
	convRepo := repos.Conversation

	userID, err := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID2, err := userRepo.CreateUser("test2@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser2", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	conversation, err := convRepo.CreateDirectConversation(int(userID), int(userID2))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	groupID, err := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID),
		Title:          "Test Group",
		Description:    "This is a test group.",
		ConversationID: conversation.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	err = groupRepo.AddMember(int(groupID), int(userID2), "member")
	if err != nil {
		t.Fatalf("Failed to add member: %v", err)
	}

	members, err := groupRepo.GetMembersByGroupID(int(groupID))
	if err != nil {
		t.Fatalf("Failed to get members: %v", err)
	}

	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(members))
	}

	if len(members) != 1 {
	}
}

func TestGroupRepository_RemoveMember(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	groupRepo := repos.Group
	conversationRepo := repos.Conversation

	userID1, err := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID2, err := userRepo.CreateUser("test2@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser2", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	conversation, err := conversationRepo.CreateDirectConversation(int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	groupID, err := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID1),
		Title:          "Test Group",
		Description:    "This is a test group.",
		ConversationID: conversation.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	err = groupRepo.AddMember(int(groupID), int(userID2), "member")
	if err != nil {
		t.Fatalf("Failed to remove member: %v", err)
	}

	err = groupRepo.RemoveMember(int(groupID), int(userID2))
	if err != nil {
		t.Fatalf("Failed to remove member: %v", err)
	}

	members, err := groupRepo.GetMembersByGroupID(int(groupID))
	if err != nil {
		t.Fatalf("Failed to get members: %v", err)
	}
	if len(members) != 0 {
		t.Errorf("Expected 0 members, got %d", len(members))
	}
}
