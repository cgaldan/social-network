package consumer

import (
	"social-network/internal/domain"
	"testing"
	"time"
)

func TestSourceFlows_FollowRequestCreatesNotification(t *testing.T) {
	services, _, pusher := setupNotificationConsumerTest(t)

	followerID := createNotificationConsumerUser(t, services, "follower-source@example.com", "followersource")
	privateUserID := createPrivateNotificationConsumerUser(t, services, "private-source@example.com", "privatesource")

	status, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: followerID,
		FolloweeID: privateUserID,
	})
	if err != nil {
		t.Fatalf("Failed to follow private user: %v", err)
	}
	if status != "pending" {
		t.Fatalf("Expected status pending, got %s", status)
	}

	notifications, err := services.Notification.ListNotifications(privateUserID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list notifications: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}
	notification := notifications[0]
	if notification.Type != "follow.requested" {
		t.Errorf("Expected notification type follow.requested, got %s", notification.Type)
	}
	if notification.RecipientID != privateUserID {
		t.Errorf("Expected recipient ID %d, got %d", privateUserID, notification.RecipientID)
	}
	if notification.ActorID == nil || *notification.ActorID != followerID {
		t.Errorf("Expected actor ID %d, got %v", followerID, notification.ActorID)
	}
	if notification.EntityType == nil || *notification.EntityType != "follow_request" {
		t.Errorf("Expected entity type follow_request, got %v", notification.EntityType)
	}

	if len(pusher.calls) != 1 {
		t.Fatalf("Expected 1 pusher call, got %d", len(pusher.calls))
	}
	if pusher.calls[0].userID != privateUserID {
		t.Errorf("Expected push user ID %d, got %d", privateUserID, pusher.calls[0].userID)
	}
}

func TestSourceFlows_FollowPublicUserDoesNotCreateNotification(t *testing.T) {
	services, _, pusher := setupNotificationConsumerTest(t)

	followerID := createNotificationConsumerUser(t, services, "follower-public@example.com", "followerpublic")
	publicUserID := createNotificationConsumerUser(t, services, "public-source@example.com", "publicsource")

	status, err := services.Follow.FollowUser(domain.FollowRequest{
		FollowerID: followerID,
		FolloweeID: publicUserID,
	})
	if err != nil {
		t.Fatalf("Failed to follow public user: %v", err)
	}
	if status != "accepted" {
		t.Fatalf("Expected status accepted, got %s", status)
	}

	notifications, err := services.Notification.ListNotifications(publicUserID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list notifications: %v", err)
	}
	if len(notifications) != 0 {
		t.Errorf("Expected 0 notifications for accepted public follow, got %d", len(notifications))
	}
	if len(pusher.calls) != 0 {
		t.Errorf("Expected 0 pusher calls for accepted public follow, got %d", len(pusher.calls))
	}
}

func TestSourceFlows_GroupInvitationCreatesNotification(t *testing.T) {
	services, _, pusher := setupNotificationConsumerTest(t)

	creatorID := createNotificationConsumerUser(t, services, "group-creator@example.com", "groupcreator")
	inviteeID := createNotificationConsumerUser(t, services, "group-invitee@example.com", "groupinvitee")

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:   creatorID,
		Title:       "Hiking Crew",
		Description: "We hike trails together.",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	if err := services.Group.CreateGroupInvitation(group.ID, creatorID, inviteeID); err != nil {
		t.Fatalf("Failed to create group invitation: %v", err)
	}

	notifications, err := services.Notification.ListNotifications(inviteeID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list notifications: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}
	notification := notifications[0]
	if notification.Type != "group.invitation.created" {
		t.Errorf("Expected notification type group.invitation.created, got %s", notification.Type)
	}
	if notification.RecipientID != inviteeID {
		t.Errorf("Expected recipient ID %d, got %d", inviteeID, notification.RecipientID)
	}
	if notification.ActorID == nil || *notification.ActorID != creatorID {
		t.Errorf("Expected actor ID %d, got %v", creatorID, notification.ActorID)
	}
	if notification.EntityType == nil || *notification.EntityType != "group_invitation" {
		t.Errorf("Expected entity type group_invitation, got %v", notification.EntityType)
	}

	if len(pusher.calls) != 1 {
		t.Fatalf("Expected 1 pusher call, got %d", len(pusher.calls))
	}
	if pusher.calls[0].userID != inviteeID {
		t.Errorf("Expected push user ID %d, got %d", inviteeID, pusher.calls[0].userID)
	}
}

func TestSourceFlows_GroupJoinRequestCreatesNotificationsForAllAdmins(t *testing.T) {
	services, _, pusher := setupNotificationConsumerTest(t)

	creatorID := createNotificationConsumerUser(t, services, "join-creator@example.com", "joincreator")
	coAdminID := createNotificationConsumerUser(t, services, "join-coadmin@example.com", "joincoadmin")
	memberID := createNotificationConsumerUser(t, services, "join-member@example.com", "joinmember")
	requesterID := createNotificationConsumerUser(t, services, "join-requester@example.com", "joinrequester")

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:   creatorID,
		Title:       "Open Group",
		Description: "Anyone can request to join.",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	if err := services.Group.AddMember(group.ConversationID, group.ID, coAdminID, "admin"); err != nil {
		t.Fatalf("Failed to add co-admin: %v", err)
	}
	if err := services.Group.AddMember(group.ConversationID, group.ID, memberID, "member"); err != nil {
		t.Fatalf("Failed to add member: %v", err)
	}

	if err := services.Group.CreateGroupJoinRequest(group.ID, requesterID); err != nil {
		t.Fatalf("Failed to create group join request: %v", err)
	}

	creatorNotifications, err := services.Notification.ListNotifications(creatorID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list creator notifications: %v", err)
	}
	if len(creatorNotifications) != 1 {
		t.Errorf("Expected 1 notification for creator, got %d", len(creatorNotifications))
	}

	coAdminNotifications, err := services.Notification.ListNotifications(coAdminID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list co-admin notifications: %v", err)
	}
	if len(coAdminNotifications) != 1 {
		t.Errorf("Expected 1 notification for co-admin, got %d", len(coAdminNotifications))
	}

	memberNotifications, err := services.Notification.ListNotifications(memberID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list member notifications: %v", err)
	}
	if len(memberNotifications) != 0 {
		t.Errorf("Expected 0 notifications for non-admin member, got %d", len(memberNotifications))
	}

	if len(pusher.calls) != 2 {
		t.Fatalf("Expected 2 pusher calls (one per admin), got %d", len(pusher.calls))
	}
	pushedRecipients := map[int]bool{}
	for _, call := range pusher.calls {
		pushedRecipients[call.userID] = true
	}
	if !pushedRecipients[creatorID] || !pushedRecipients[coAdminID] {
		t.Errorf("Expected pushes for creator (%d) and co-admin (%d), got %v", creatorID, coAdminID, pushedRecipients)
	}
}

func TestSourceFlows_GroupEventCreatedCreatesNotificationForAllMembers(t *testing.T) {
	services, _, pusher := setupNotificationConsumerTest(t)

	creatorID := createNotificationConsumerUser(t, services, "event-creator@example.com", "eventcreator")
	memberID := createNotificationConsumerUser(t, services, "event-member@example.com", "eventmember")
	memberID2 := createNotificationConsumerUser(t, services, "event-member2@example.com", "eventmember2")
	memberID3 := createNotificationConsumerUser(t, services, "event-member3@example.com", "eventmember3")

	group, err := services.Group.CreateGroup(&domain.Group{
		CreatorID:   creatorID,
		Title:       "Open Group",
		Description: "Anyone can request to join.",
	})
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}

	if err := services.Group.AddMember(group.ConversationID, group.ID, memberID, "member"); err != nil {
		t.Fatalf("Failed to add member: %v", err)
	}
	if err := services.Group.AddMember(group.ConversationID, group.ID, memberID2, "member"); err != nil {
		t.Fatalf("Failed to add admin: %v", err)
	}
	if err := services.Group.AddMember(group.ConversationID, group.ID, memberID3, "member"); err != nil {
		t.Fatalf("Failed to add member: %v", err)
	}

	event, err := services.Group.CreateGroupEvent(group.ID, creatorID, domain.CreateGroupEventRequest{
		Title:       "Open Group Event",
		Description: "Anyone can join this event.",
		StartsAt:    time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("Failed to create group event: %v", err)
	}

	notifications, err := services.Notification.ListNotifications(memberID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list notifications: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("Expected 1 notification, got %d", len(notifications))
	}
	notification := notifications[0]
	if notification.Type != "group.event.created" {
		t.Errorf("Expected notification type group.event.created, got %s", notification.Type)
	}
	if notification.RecipientID != memberID {
		t.Errorf("Expected recipient ID %d, got %d", memberID, notification.RecipientID)
	}
	if notification.ActorID == nil || *notification.ActorID != creatorID {
		t.Errorf("Expected actor ID %d, got %v", creatorID, notification.ActorID)
	}
	if notification.EntityType == nil || *notification.EntityType != "group_event" {
		t.Errorf("Expected entity type group_event, got %v", notification.EntityType)
	}
	if notification.EntityID == nil || *notification.EntityID != event.ID {
		t.Errorf("Expected entity ID %d, got %v", event.ID, notification.EntityID)
	}
	if len(pusher.calls) != 3 {
		t.Fatalf("Expected 3 pusher calls, got %d", len(pusher.calls))
	}
	if pusher.calls[0].userID != memberID {
		t.Errorf("Expected push user ID %d, got %d", memberID, pusher.calls[0].userID)
	}
	if pusher.calls[1].userID != memberID2 {
		t.Errorf("Expected push user ID %d, got %d", memberID2, pusher.calls[1].userID)
	}
	if pusher.calls[2].userID != memberID3 {
		t.Errorf("Expected push user ID %d, got %d", memberID3, pusher.calls[2].userID)
	}
}
