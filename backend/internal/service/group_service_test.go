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

func TestGroupService_CreateGroupInvitation(t *testing.T) {
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

	err = services.Group.CreateGroupInvitation(group.ID, userID, userID2)
	if err != nil {
		t.Fatalf("Failed to create group invitation: %v", err)
	}

	invitations, err := services.Group.GetGroupInvitationsByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get group invitations: %v", err)
	}
	if len(invitations) != 1 {
		t.Errorf("Expected 1 invitation, got %d", len(invitations))
	}

	err = services.Group.AcceptGroupInvitation(userID2, &invitations[0])
	if err != nil {
		t.Fatalf("Failed to accept group invitation: %v", err)
	}

	members, err := services.Group.GetMembersByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get members: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members))
	}
	if members[0].UserID != userID {
		t.Errorf("Expected member 0 to be user %d, got %d", userID, members[0].UserID)
	}
	if members[1].UserID != userID2 {
		t.Errorf("Expected member 1 to be user %d, got %d", userID2, members[1].UserID)
	}

	invitations, err = services.Group.GetGroupInvitationsByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get group invitations: %v", err)
	}
	if invitations[0].Status != "accepted" {
		t.Errorf("Expected invitation status accepted, got %q", invitations[0].Status)
	}

	err = services.Group.DeclineGroupInvitation(userID2, &invitations[0])
	if err != nil {
		if err.Error() != "invitation is not pending" {
			t.Fatalf("Expected invitation is not pending error, got %v", err)
		}
	}

	err = services.Group.CreateGroupInvitation(group.ID, userID, userID3)
	if err != nil {
		t.Fatalf("Failed to create group invitation: %v", err)
	}

	invitations, err = services.Group.GetGroupInvitationsByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get group invitations: %v", err)
	}
	if len(invitations) != 2 {
		t.Errorf("Expected 2 invitations, got %d", len(invitations))
	}

	err = services.Group.DeclineGroupInvitation(userID2, &invitations[0])
	if err != nil {
		if err.Error() != "invitation is not pending" {
			t.Fatalf("Expected invitation is not pending error, got %v", err)
		}
	}

	invitations, err = services.Group.GetGroupInvitationsByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get group invitations: %v", err)
	}

	if len(invitations) != 2 {
		t.Errorf("Expected 2 invitation, got %d", len(invitations))
	} else if invitations[0].Status != "accepted" {
		t.Errorf("Expected remaining invitation to be accepted, got %q", invitations[0].Status)
	}
}

func TestGroupService_CreateGroupJoinRequest(t *testing.T) {
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

	err = services.Group.CreateGroupJoinRequest(group.ID, userID2)
	if err != nil {
		t.Fatalf("Failed to create group join request: %v", err)
	}

	requests, err := services.Group.GetGroupJoinRequestsByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get group join requests: %v", err)
	}
	if len(requests) != 1 {
		t.Errorf("Expected 1 join request, got %d", len(requests))
	}

	err = services.Group.AcceptGroupJoinRequest(userID, &requests[0])
	if err != nil {
		t.Fatalf("Failed to accept group join request: %v", err)
	}

	members, err := services.Group.GetMembersByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get members: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members))
	}
	if members[0].UserID != userID {
		t.Errorf("Expected member 0 to be user %d, got %d", userID, members[0].UserID)
	}
	if members[1].UserID != userID2 {
		t.Errorf("Expected member 1 to be user %d, got %d", userID2, members[1].UserID)
	}

	requests, err = services.Group.GetGroupJoinRequestsByGroupID(group.ID)
	if requests[0].Status != "accepted" {
		t.Errorf("Expected join request status accepted, got %q", requests[0].Status)
	}

	err = services.Group.DeclineGroupJoinRequest(userID, &requests[0])
	if err != nil {
		if err.Error() != "join request is not pending" {
			t.Fatalf("Expected join request is not pending error, got %v", err)
		}
	}

	err = services.Group.CreateGroupJoinRequest(group.ID, userID3)
	if err != nil {
		t.Fatalf("Failed to create group join request: %v", err)
	}

	requests, err = services.Group.GetGroupJoinRequestsByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get group join requests: %v", err)
	}
	if len(requests) != 2 {
		t.Errorf("Expected 1 join request, got %d", len(requests))
	}

	err = services.Group.DeclineGroupJoinRequest(userID, &requests[1])
	if err != nil {
		t.Fatalf("Failed to decline group join request: %v", err)
	}

	requests, err = services.Group.GetGroupJoinRequestsByGroupID(group.ID)
	if err != nil {
		t.Fatalf("Failed to get group join requests: %v", err)
	}
	if len(requests) != 1 {
		t.Errorf("Expected 1 join request, got %d", len(requests))
	}
}
