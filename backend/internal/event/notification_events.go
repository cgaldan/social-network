package event

import (
	"strconv"
	"time"
)

const (
	FollowRequestedEvent        EventType = "follow.requested"
	GroupInvitationCreatedEvent EventType = "group.invitation.created"
	GroupJoinRequestedEvent     EventType = "group.join_requested"
	GroupEventCreatedEvent      EventType = "group.event.created"
)

type FollowRequestedEventData struct {
	BaseEvent
	FollowerID int
	FolloweeID int
	FollowID   int
	ActorName  string
}

func NewFollowRequestedEvent(followerID, followeeID, followID int, actorName string) *FollowRequestedEventData {
	return &FollowRequestedEventData{
		BaseEvent: BaseEvent{
			EventType: FollowRequestedEvent,
			CreatedAt: time.Now(),
			ID:        strconv.Itoa(followID),
		},
		FollowerID: followerID,
		FolloweeID: followeeID,
		FollowID:   followID,
		ActorName:  actorName,
	}
}

type GroupInvitationCreatedEventData struct {
	BaseEvent
	GroupID      int
	InvitationID int
	InviterID    int
	InviteeID    int
	GroupTitle   string
	ActorName    string
}

func NewGroupInvitationCreatedEvent(groupID, invitationID, inviterID, inviteeID int, groupTitle, actorName string) *GroupInvitationCreatedEventData {
	return &GroupInvitationCreatedEventData{
		BaseEvent: BaseEvent{
			EventType: GroupInvitationCreatedEvent,
			CreatedAt: time.Now(),
			ID:        strconv.Itoa(invitationID),
		},
		GroupID:      groupID,
		InvitationID: invitationID,
		InviterID:    inviterID,
		InviteeID:    inviteeID,
		GroupTitle:   groupTitle,
		ActorName:    actorName,
	}
}

type GroupJoinRequestedEventData struct {
	BaseEvent
	GroupID     int
	RequestID   int
	RequesterID int
	RecipientID int
	GroupTitle  string
	ActorName   string
}

func NewGroupJoinRequestedEvent(groupID, requestID, requesterID, recipientID int, groupTitle, actorName string) *GroupJoinRequestedEventData {
	return &GroupJoinRequestedEventData{
		BaseEvent: BaseEvent{
			EventType: GroupJoinRequestedEvent,
			CreatedAt: time.Now(),
			ID:        strconv.Itoa(requestID),
		},
		GroupID:     groupID,
		RequestID:   requestID,
		RequesterID: requesterID,
		RecipientID: recipientID,
		GroupTitle:  groupTitle,
		ActorName:   actorName,
	}
}

type GroupEventCreatedEventData struct {
	BaseEvent
	GroupID     int
	EventID     int
	CreatorID   int
	RecipientID int
	GroupTitle  string
	EventTitle  string
	ActorName   string
}

func NewGroupEventCreatedEvent(groupID, eventID, creatorID, recipientID int, groupTitle, eventTitle, actorName string) *GroupEventCreatedEventData {
	return &GroupEventCreatedEventData{
		BaseEvent: BaseEvent{
			EventType: GroupEventCreatedEvent,
			CreatedAt: time.Now(),
			ID:        strconv.Itoa(eventID),
		},
		GroupID:     groupID,
		EventID:     eventID,
		CreatorID:   creatorID,
		RecipientID: recipientID,
		GroupTitle:  groupTitle,
		EventTitle:  eventTitle,
		ActorName:   actorName,
	}
}
