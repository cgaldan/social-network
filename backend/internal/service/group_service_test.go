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
	userID3 := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "test3@example.com",
		Password:    "password123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "testuser3",
		Gender:      "male",
		IsPublic:    true,
	})

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:   userID,
		Title:       "Test Group",
		Description: "This is a test group.",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	if group.ID == 0 {
		t.Error("Expected non-zero group ID")
	}

	if group.ConversationID == 0 {
		t.Error("Expected non-zero conversation ID")
	}

	members, err := services.Group.GetMembersByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get members: %v", err)
	}
	if len(members) != 1 {
		t.Errorf("Expected 1 member, got %d", len(members))
	}
	if members[0].UserID != group.CreatorID {
		t.Errorf("Expected member 0 to be user %d, got %d", group.CreatorID, members[0].UserID)
	}

	err = services.Group.AddMember(group.ConversationID, group.ID, userID3, "member")
	if err != nil {
		t.Fatalf("Failed to add member: %v", err)
	}
	err = services.Group.AddMember(group.ConversationID, group.ID, userID2, "member")
	if err != nil {
		t.Fatalf("Failed to add member: %v", err)
	}

	members, err = services.Group.GetMembersByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get members: %v", err)
	}
	if len(members) != 3 {
		t.Errorf("Expected 3 members, got %d", len(members))
	}

	if members[0].UserID != userID {
		t.Errorf("Expected member 0 to be user %d, got %d", userID, members[0].UserID)
	}
	if members[1].UserID != userID3 {
		t.Errorf("Expected member 1 to be user %d, got %d", userID3, members[1].UserID)
	}

	err = services.Group.RemoveMember(group.ConversationID, group.ID, userID3)
	if err != nil {
		t.Fatalf("Failed to remove member: %v", err)
	}

	members, err = services.Group.GetMembersByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get members: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members))
	}
	if members[0].UserID != userID {
		t.Errorf("Expected member 0 to be user %d, got %d", userID, members[0].UserID)
	}
}
