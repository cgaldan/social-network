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

	GroupName := "Test Group"

	conversation, err := convRepo.CreateGroupConversation(GroupName, int(userID), int(userID2))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	groupID, err := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID),
		Title:          GroupName,
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

	GroupName := "Test Group"

	conversation, err := convRepo.CreateGroupConversation(GroupName, int(userID), int(userID2))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	groupID, err := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID),
		Title:          GroupName,
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

	GroupName := "Test Group"

	conversation, err := conversationRepo.CreateGroupConversation(GroupName, int(userID1), int(userID2))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	groupID, err := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID1),
		Title:          GroupName,
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

func TestGroupRepository_CreateGroupInvitation(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	groupRepo := repos.Group
	convRepo := repos.Conversation

	userID1, err := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID2, err := userRepo.CreateUser("test2@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser2", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	GroupName := "Test Group"

	conversation, err := convRepo.CreateGroupConversation(GroupName, int(userID2))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	groupID, err := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID2),
		Title:          GroupName,
		Description:    "This is a test group.",
		ConversationID: conversation.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	err = groupRepo.CreateGroupInvitation(int(groupID), int(userID2), int(userID1))
	if err != nil {
		t.Fatalf("Failed to create group invitation: %v", err)
	}

	invitations, err := groupRepo.GetGroupInvitationsByGroupID(int(groupID))
	if err != nil {
		t.Fatalf("Failed to get group invitations: %v", err)
	}
	if len(invitations) != 1 {
		t.Errorf("Expected 1 invitation, got %d", len(invitations))
	}
	if invitations[0].GroupID != int(groupID) {
		t.Errorf("Expected group ID %d, got %d", groupID, invitations[0].GroupID)
	}
	if invitations[0].InviterID != int(userID2) {
		t.Errorf("Expected inviter ID %d, got %d", userID2, invitations[0].InviterID)
	}
	if invitations[0].InviteeID != int(userID1) {
		t.Errorf("Expected invitee ID %d, got %d", userID1, invitations[0].InviteeID)
	}

	err = groupRepo.UpdateGroupInvitationStatus(invitations[0].ID, "accepted")
	if err != nil {
		t.Fatalf("Failed to update group invitation status: %v", err)
	}

	invitations, err = groupRepo.GetGroupInvitationsByGroupID(int(groupID))
	if err != nil {
		t.Fatalf("Failed to get group invitations: %v", err)
	}
	if len(invitations) != 1 {
		t.Errorf("Expected 1 invitation, got %d", len(invitations))
	}
	if invitations[0].Status != "accepted" {
		t.Errorf("Expected status 'accepted', got '%s'", invitations[0].Status)
	}

	err = groupRepo.DeleteGroupInvitation(invitations[0].ID)
	if err != nil {
		t.Fatalf("Failed to delete group invitation: %v", err)
	}

	invitations, err = groupRepo.GetGroupInvitationsByGroupID(int(groupID))
	if err != nil {
		t.Fatalf("Failed to get group invitations: %v", err)
	}
	if len(invitations) != 0 {
		t.Errorf("Expected 0 invitations, got %d", len(invitations))
	}
}

func TestGroupRepository_ListGroups(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	groupRepo := repos.Group
	convRepo := repos.Conversation

	userID, err := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	titles := []string{"Group 1", "Group 2", "Group 3", "Group 4", "Group 5"}
	for _, title := range titles {
		conversation, err := convRepo.CreateGroupConversation(title, int(userID))
		if err != nil {
			t.Fatalf("Failed to create conversation: %v", err)
		}
		_, err = groupRepo.CreateGroup(&domain.Group{
			CreatorID:      int(userID),
			Title:          title,
			Description:    "Test description",
			ConversationID: conversation.ID,
		})
		if err != nil {
			t.Fatalf("Failed to create group: %v", err)
		}
	}

	groups, err := groupRepo.ListGroups(10, 0)
	if err != nil {
		t.Fatalf("Failed to list groups: %v", err)
	}
	if len(groups) != 5 {
		t.Errorf("Expected 5 groups, got %d", len(groups))
	}

	groups, err = groupRepo.ListGroups(2, 0)
	if err != nil {
		t.Fatalf("Failed to list groups with limit: %v", err)
	}
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups with limit 2, got %d", len(groups))
	}

	groups, err = groupRepo.ListGroups(2, 2)
	if err != nil {
		t.Fatalf("Failed to list groups with offset: %v", err)
	}
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups with limit 2 offset 2, got %d", len(groups))
	}

	groups, err = groupRepo.ListGroups(10, 5)
	if err != nil {
		t.Fatalf("Failed to list groups with offset past end: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("Expected 0 groups with offset 5, got %d", len(groups))
	}
}

func TestGroupRepository_CreateGroupJoinRequest(t *testing.T) {
	repos := SetupTestDB(t)
	userRepo := repos.User
	groupRepo := repos.Group
	convRepo := repos.Conversation

	userID1, err := userRepo.CreateUser("test@example.com", "hashedpass", "John", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser", "male", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID2, err := userRepo.CreateUser("test2@example.com", "hashedpass", "Jane", "Doe", time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC), "testuser2", "female", "", "", false)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	GroupName := "Test Group"

	conversation, err := convRepo.CreateGroupConversation(GroupName, int(userID1))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	groupID, err := groupRepo.CreateGroup(&domain.Group{
		CreatorID:      int(userID1),
		Title:          GroupName,
		Description:    "This is a test group.",
		ConversationID: conversation.ID,
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	err = groupRepo.CreateGroupJoinRequest(int(groupID), int(userID2))
	if err != nil {
		t.Fatalf("Failed to create group join request: %v", err)
	}

	requests, err := groupRepo.GetGroupJoinRequestsByGroupID(int(groupID))
	if err != nil {
		t.Fatalf("Failed to get group join requests: %v", err)
	}
	if len(requests) != 1 {
		t.Errorf("Expected 1 join request, got %d", len(requests))
	}
	if requests[0].GroupID != int(groupID) {
		t.Errorf("Expected group ID %d, got %d", groupID, requests[0].GroupID)
	}
	if requests[0].UserID != int(userID2) {
		t.Errorf("Expected user ID %d, got %d", userID2, requests[0].UserID)
	}

	err = groupRepo.UpdateGroupJoinRequestStatus(requests[0].ID, "accepted")
	if err != nil {
		t.Fatalf("Failed to update group join request status: %v", err)
	}

	requests, err = groupRepo.GetGroupJoinRequestsByGroupID(int(groupID))
	if err != nil {
		t.Fatalf("Failed to get group join requests: %v", err)
	}
	if len(requests) != 1 {
		t.Errorf("Expected 1 join request, got %d", len(requests))
	}
	if requests[0].Status != "accepted" {
		t.Errorf("Expected status 'accepted', got '%s'", requests[0].Status)
	}

	err = groupRepo.DeleteGroupJoinRequest(requests[0].ID)
	if err != nil {
		t.Fatalf("Failed to delete group join request: %v", err)
	}

	requests, err = groupRepo.GetGroupJoinRequestsByGroupID(int(groupID))
	if err != nil {
		t.Fatalf("Failed to get group join requests: %v", err)
	}
	if len(requests) != 0 {
		t.Errorf("Expected 0 join requests, got %d", len(requests))
	}
}
