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

func TestGroupService_ListGroups(t *testing.T) {
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

	titles := []string{"Group 1", "Group 2", "Group 3", "Group 4", "Group 5"}
	for _, title := range titles {
		_, err := services.Group.CreateGroup(&domain.Group{
			CreatorID:   userID,
			Title:       title,
			Description: "Test description",
		})
		if err != nil {
			t.Fatalf("Failed to create group: %v", err)
		}
	}

	t.Run("list all groups", func(t *testing.T) {
		groups, err := services.Group.ListGroups(10, 0)
		if err != nil {
			t.Fatalf("Failed to list groups: %v", err)
		}
		if len(groups) != 5 {
			t.Errorf("Expected 5 groups, got %d", len(groups))
		}
	})

	t.Run("list groups with limit", func(t *testing.T) {
		groups, err := services.Group.ListGroups(2, 0)
		if err != nil {
			t.Fatalf("Failed to list groups with limit: %v", err)
		}
		if len(groups) != 2 {
			t.Errorf("Expected 2 groups with limit 2, got %d", len(groups))
		}
	})

	t.Run("list groups with offset", func(t *testing.T) {
		groups, err := services.Group.ListGroups(2, 2)
		if err != nil {
			t.Fatalf("Failed to list groups with offset: %v", err)
		}
		if len(groups) != 2 {
			t.Errorf("Expected 2 groups with limit 2 offset 2, got %d", len(groups))
		}
	})

	t.Run("list groups with invalid limit defaults", func(t *testing.T) {
		groups, err := services.Group.ListGroups(0, -1)
		if err != nil {
			t.Fatalf("Failed to list groups with invalid pagination: %v", err)
		}
		if len(groups) != 5 {
			t.Errorf("Expected 5 groups with defaulted pagination, got %d", len(groups))
		}
	})
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
	if len(requests) != 2 {
		t.Errorf("Expected 2 join request, got %d", len(requests))
	}
}

func TestGroupService_CreateGroupEvent(t *testing.T) {
	services := SetupTestServices(t)

	creatorID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "creator@example.com",
		Password:    "password123",
		FirstName:   "Creator",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Nickname:    "creator",
		Gender:      "male",
		IsPublic:    true,
	})
	outsiderID := CreateTestUser(t, services, domain.RegisterRequest{
		Email:       "outsider@example.com",
		Password:    "password123",
		FirstName:   "Outsider",
		LastName:    "User",
		DateOfBirth: time.Now().AddDate(-30, 0, 0),
		Nickname:    "outsider",
		Gender:      "female",
		IsPublic:    true,
	})

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:   creatorID,
		Title:       "Event Group",
		Description: "This is an event group.",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	event, err := services.Group.CreateGroupEvent(creatorID, group.ID, domain.CreateGroupEventRequest{
		Title:       "Group Meetup",
		Description: "Planning the next group meetup.",
		StartsAt:    time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("Failed to create group event: %v", err)
	}
	if event.GroupID != group.ID {
		t.Errorf("Expected group ID %d, got %d", group.ID, event.GroupID)
	}
	if event.CreatorID != creatorID {
		t.Errorf("Expected creator ID %d, got %d", creatorID, event.CreatorID)
	}

	t.Run("member lists group events", func(t *testing.T) {
		events, err := services.Group.ListGroupEvents(creatorID, group.ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to list group events: %v", err)
		}
		if len(events) != 1 {
			t.Errorf("Expected 1 group event, got %d", len(events))
		}
	})

	t.Run("non-member cannot create group event", func(t *testing.T) {
		_, err := services.Group.CreateGroupEvent(outsiderID, group.ID, domain.CreateGroupEventRequest{
			Title:       "Outsider Event",
			Description: "This should not be created.",
			StartsAt:    time.Now().Add(24 * time.Hour),
		})
		if err == nil {
			t.Error("Expected error when non-member creates group event")
		}
	})
}
