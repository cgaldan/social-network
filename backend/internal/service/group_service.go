package service

import (
	"fmt"
	"social-network/internal/domain"
	"social-network/internal/event"
	"social-network/internal/repository"
	"social-network/packages/logger"
	"time"
)

type GroupService struct {
	groupRepo   repository.GroupRepositoryInterface
	userRepo    repository.UserRepositoryInterface
	convService ConversationServiceInterface
	eventBus    event.EventBus
	logger      *logger.Logger
}

func NewGroupService(groupRepo repository.GroupRepositoryInterface, userRepo repository.UserRepositoryInterface, convService ConversationServiceInterface, eventBus event.EventBus, logger *logger.Logger) *GroupService {
	return &GroupService{
		groupRepo:   groupRepo,
		userRepo:    userRepo,
		convService: convService,
		eventBus:    eventBus,
		logger:      logger,
	}
}

func (s *GroupService) CreateGroup(group *domain.Group) (*domain.Group, error) {
	conversation, err := s.convService.CreateGroupConversation(group.Title, group.CreatorID)
	if err != nil {
		s.logger.Error("Failed to create group conversation", "error", err)
		return nil, fmt.Errorf("failed to create group conversation")
	}
	group.ConversationID = conversation.ID

	groupID, err := s.groupRepo.CreateGroup(group)
	if err != nil {
		s.logger.Error("Failed to create group", "error", err)
		return nil, fmt.Errorf("failed to create group")
	}

	group, err = s.groupRepo.GetGroupByID(int(groupID))
	if err != nil {
		s.logger.Error("Failed to get group by ID", "error", err)
		return nil, fmt.Errorf("failed to get group by ID")
	}

	err = s.groupRepo.AddMember(group.ID, group.CreatorID, "admin")
	if err != nil {
		s.logger.Error("Failed to add member", "error", err)
		return nil, fmt.Errorf("failed to add member")
	}

	return group, nil
}

func (s *GroupService) ListGroups(limit, offset int) ([]domain.Group, error) {
	limit, offset = s.validateLimitAndOffset(limit, offset)

	groups, err := s.groupRepo.ListGroups(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list groups", "error", err, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to list groups")
	}

	return groups, nil
}

func (s *GroupService) validateLimitAndOffset(limit, offset int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func (s *GroupService) GetMembersByGroupID(groupID int) ([]domain.GroupMember, error) {
	return s.groupRepo.GetMembersByGroupID(groupID)
}

func (s *GroupService) AddMember(convID, groupID, userID int, role string) error {
	err := s.groupRepo.AddMember(groupID, userID, role)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	err = s.convService.AddConversationParticipant(convID, userID)
	if err != nil {
		if remErr := s.groupRepo.RemoveMember(groupID, userID); remErr != nil {
			s.logger.Error("failed to roll back group membership after conversation add failed", "error", remErr, "group_id", groupID, "user_id", userID)
		}
		return fmt.Errorf("failed to add conversation participant: %w", err)
	}

	return nil
}

func (s *GroupService) RemoveMember(convID, groupID, userID int) error {
	err := s.groupRepo.RemoveMember(groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	err = s.convService.RemoveConversationParticipant(convID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove conversation participant: %w", err)
	}

	return nil
}

func (s *GroupService) CreateGroupInvitation(groupID, inviterID, inviteeID int) error {
	isUserInGroup, err := s.groupRepo.IsUserInGroup(groupID, inviterID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if !isUserInGroup {
		return fmt.Errorf("inviter is not in group")
	}

	invitations, err := s.groupRepo.GetGroupInvitationsByGroupID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group invitation: %w", err)
	}
	for _, invitation := range invitations {
		if invitation.InviteeID == inviteeID && invitation.Status == "pending" {
			return fmt.Errorf("there is already a pending invitation for this user")
		}
	}

	joinRequests, err := s.groupRepo.GetGroupJoinRequestsByGroupID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group join requests: %w", err)
	}
	for _, request := range joinRequests {
		if request.UserID == inviteeID && request.Status == "pending" {
			return fmt.Errorf("there is already a pending join request for this user")
		}
	}

	isUserInGroup, err = s.groupRepo.IsUserInGroup(groupID, inviteeID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if isUserInGroup {
		return fmt.Errorf("user is already in group")
	}

	invitationID, err := s.groupRepo.CreateGroupInvitation(groupID, inviterID, inviteeID)
	if err != nil {
		return fmt.Errorf("failed to create group invitation: %w", err)
	}

	s.publishGroupInvitationCreated(groupID, int(invitationID), inviterID, inviteeID)

	return nil
}

func (s *GroupService) publishGroupInvitationCreated(groupID, invitationID, inviterID, inviteeID int) {
	if s.eventBus == nil {
		return
	}

	groupTitle := s.lookupGroupTitle(groupID)
	actorName := s.lookupActorName(inviterID)

	if err := s.eventBus.Publish(event.NewGroupInvitationCreatedEvent(groupID, invitationID, inviterID, inviteeID, groupTitle, actorName)); err != nil {
		s.logger.Error("Failed to publish group invitation created event", "error", err, "invitationID", invitationID)
	}
}

func (s *GroupService) CreateGroupJoinRequest(groupID, userID int) error {
	joinRequests, err := s.groupRepo.GetGroupJoinRequestsByGroupID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group join requests: %w", err)
	}
	for _, request := range joinRequests {
		if request.UserID == userID && request.Status == "pending" {
			return fmt.Errorf("there is already a pending join request for this user")
		}
	}

	invitations, err := s.groupRepo.GetGroupInvitationsByGroupID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group invitations: %w", err)
	}
	for _, invitation := range invitations {
		if invitation.InviteeID == userID && invitation.Status == "pending" {
			return fmt.Errorf("there is already a pending invitation for this user")
		}
	}

	isUserInGroup, err := s.groupRepo.IsUserInGroup(groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if isUserInGroup {
		return fmt.Errorf("user is already in group")
	}

	requestID, err := s.groupRepo.CreateGroupJoinRequest(groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to create group join request: %w", err)
	}

	s.publishGroupJoinRequested(groupID, int(requestID), userID)

	return nil
}

func (s *GroupService) publishGroupJoinRequested(groupID, requestID, requesterID int) {
	if s.eventBus == nil {
		return
	}

	members, err := s.groupRepo.GetMembersByGroupID(groupID)
	if err != nil {
		s.logger.Error("Failed to get group members for join requested event", "error", err, "groupID", groupID)
		return
	}

	groupTitle := s.lookupGroupTitle(groupID)
	actorName := s.lookupActorName(requesterID)

	for _, member := range members {
		if member.Role != "admin" {
			continue
		}
		if err := s.eventBus.Publish(event.NewGroupJoinRequestedEvent(groupID, requestID, requesterID, member.UserID, groupTitle, actorName)); err != nil {
			s.logger.Error("Failed to publish group join requested event", "error", err, "requestID", requestID, "recipientID", member.UserID)
		}
	}
}

func (s *GroupService) lookupGroupTitle(groupID int) string {
	group, err := s.groupRepo.GetGroupByID(groupID)
	if err != nil || group == nil {
		return ""
	}
	return group.Title
}

func (s *GroupService) lookupActorName(userID int) string {
	if s.userRepo == nil {
		return ""
	}
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil || user == nil {
		return ""
	}
	return displayUserName(user)
}

func (s *GroupService) AcceptGroupInvitation(userID int, invitation *domain.GroupInvitation) error {
	if userID != invitation.InviteeID {
		return fmt.Errorf("user is not the invitee")
	}

	if invitation.Status != "pending" {
		return fmt.Errorf("invitation is not pending")
	}

	isUserInGroup, err := s.groupRepo.IsUserInGroup(invitation.GroupID, invitation.InviteeID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if isUserInGroup {
		return fmt.Errorf("user is already in group")
	}

	isUserInGroup, err = s.groupRepo.IsUserInGroup(invitation.GroupID, invitation.InviterID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if !isUserInGroup {
		return fmt.Errorf("inviter is not in group")
	}

	group, err := s.groupRepo.GetGroupByID(invitation.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	err = s.AddMember(group.ConversationID, group.ID, invitation.InviteeID, "member")
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	err = s.groupRepo.UpdateGroupInvitationStatus(invitation.ID, "accepted")
	if err != nil {
		if remErr := s.RemoveMember(group.ConversationID, group.ID, invitation.InviteeID); remErr != nil {
			s.logger.Error("failed to roll back after invitation status update failed", "error", remErr, "group_id", group.ID, "invitation_id", invitation.ID)
		}
		return fmt.Errorf("failed to accept group invitation: %w", err)
	}

	return nil
}

func (s *GroupService) AcceptGroupJoinRequest(answererID int, request *domain.GroupJoinRequest) error {
	if request.Status != "pending" {
		return fmt.Errorf("join request is not pending")
	}

	isUserInGroup, err := s.groupRepo.IsUserInGroup(request.GroupID, answererID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if !isUserInGroup {
		return fmt.Errorf("answerer is not in group")
	}

	isUserInGroup, err = s.groupRepo.IsUserInGroup(request.GroupID, request.UserID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if isUserInGroup {
		return fmt.Errorf("user is already in group")
	}

	isAdmin, err := s.groupRepo.IsUserAdmin(request.GroupID, answererID)
	if err != nil {
		return fmt.Errorf("failed to check if user is admin: %w", err)
	}
	if !isAdmin {
		return fmt.Errorf("answerer is not authorized to accept join requests")
	}

	group, err := s.groupRepo.GetGroupByID(request.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	err = s.AddMember(group.ConversationID, group.ID, request.UserID, "member")
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	err = s.groupRepo.UpdateGroupJoinRequestStatus(request.ID, "accepted")
	if err != nil {
		if remErr := s.RemoveMember(group.ConversationID, group.ID, request.UserID); remErr != nil {
			s.logger.Error("failed to roll back after join request status update failed", "error", remErr, "group_id", group.ID, "request_id", request.ID)
		}
		return fmt.Errorf("failed to accept group join request: %w", err)
	}

	return nil
}

func (s *GroupService) DeclineGroupInvitation(userID int, invitation *domain.GroupInvitation) error {
	if invitation.Status != "pending" {
		return fmt.Errorf("invitation is not pending")
	}

	isUserInGroup, err := s.groupRepo.IsUserInGroup(invitation.GroupID, invitation.InviteeID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if isUserInGroup {
		return fmt.Errorf("user is already in group")
	}

	isUserInGroup, err = s.groupRepo.IsUserInGroup(invitation.GroupID, invitation.InviterID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if !isUserInGroup {
		return fmt.Errorf("inviter is not in group")
	}

	err = s.groupRepo.UpdateGroupInvitationStatus(invitation.ID, "declined")
	if err != nil {
		return fmt.Errorf("failed to decline group invitation: %w", err)
	}

	// err = s.groupRepo.DeleteGroupInvitation(invitation.ID)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete group invitation: %w", err)
	// }

	return nil
}

func (s *GroupService) DeclineGroupJoinRequest(answererID int, request *domain.GroupJoinRequest) error {
	if request.Status != "pending" {
		return fmt.Errorf("join request is not pending")
	}

	isUserInGroup, err := s.groupRepo.IsUserInGroup(request.GroupID, answererID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if !isUserInGroup {
		return fmt.Errorf("answerer is not in group")
	}

	isUserInGroup, err = s.groupRepo.IsUserInGroup(request.GroupID, request.UserID)
	if err != nil {
		return fmt.Errorf("failed to check if user is in group: %w", err)
	}
	if isUserInGroup {
		return fmt.Errorf("user is already in group")
	}

	isAdmin, err := s.groupRepo.IsUserAdmin(request.GroupID, answererID)
	if err != nil {
		return fmt.Errorf("failed to check if user is admin: %w", err)
	}
	if !isAdmin {
		return fmt.Errorf("answerer is not authorized to decline join requests")
	}

	err = s.groupRepo.UpdateGroupJoinRequestStatus(request.ID, "declined")
	if err != nil {
		return fmt.Errorf("failed to decline group join request: %w", err)
	}

	// err = s.groupRepo.DeleteGroupJoinRequest(request.ID)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete group join request: %w", err)
	// }

	return nil
}

func (s *GroupService) GetGroupInvitationByID(invitationID int) (*domain.GroupInvitation, error) {
	return s.groupRepo.GetGroupInvitationByID(invitationID)
}

func (s *GroupService) GetGroupJoinRequestByID(requestID int) (*domain.GroupJoinRequest, error) {
	return s.groupRepo.GetGroupJoinRequestByID(requestID)
}

func (s *GroupService) GetGroupInvitationsByGroupID(groupID int) ([]domain.GroupInvitation, error) {
	return s.groupRepo.GetGroupInvitationsByGroupID(groupID)
}

func (s *GroupService) GetGroupJoinRequestsByGroupID(groupID int) ([]domain.GroupJoinRequest, error) {
	return s.groupRepo.GetGroupJoinRequestsByGroupID(groupID)
}

func (s *GroupService) CreateGroupEvent(userID, groupID int, eventData domain.CreateGroupEventRequest) (*domain.GroupEvent, error) {
	isMember, err := s.groupRepo.IsUserInGroup(groupID, userID)
	if err != nil {
		s.logger.Error("Failed to check group membership", "error", err, "groupID", groupID, "userID", userID)
		return nil, fmt.Errorf("failed to create group event")
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of this group")
	}

	if err := s.validateGroupEvent(eventData); err != nil {
		return nil, err
	}

	eventID, err := s.groupRepo.CreateGroupEvent(&domain.GroupEvent{
		GroupID:     groupID,
		CreatorID:   userID,
		Title:       eventData.Title,
		Description: eventData.Description,
		StartsAt:    eventData.StartsAt,
	})
	if err != nil {
		s.logger.Error("Failed to create group event", "error", err, "groupID", groupID, "userID", userID)
		return nil, fmt.Errorf("failed to create group event")
	}

	event, err := s.groupRepo.GetGroupEventByID(int(eventID))
	if err != nil {
		s.logger.Error("Failed to get group event by ID", "error", err, "eventID", eventID)
		return nil, fmt.Errorf("failed to get group event by ID")
	}

	s.publishGroupEventCreated(groupID, int(eventID), userID, eventData.Title)

	return event, nil
}

func (s *GroupService) publishGroupEventCreated(groupID, eventID, creatorID int, eventTitle string) {
	if s.eventBus == nil {
		return
	}

	members, err := s.groupRepo.GetMembersByGroupID(groupID)
	if err != nil {
		s.logger.Error("Failed to get group members for group event created event", "error", err, "groupID", groupID)
		return
	}

	groupTitle := s.lookupGroupTitle(groupID)
	actorName := s.lookupActorName(creatorID)

	for _, member := range members {
		if member.UserID == creatorID {
			continue
		}
		if err := s.eventBus.Publish(event.NewGroupEventCreatedEvent(groupID, eventID, creatorID, member.UserID, groupTitle, eventTitle, actorName)); err != nil {
			s.logger.Error("Failed to publish group event created event", "error", err, "eventID", eventID)
		}
	}
}

func (s *GroupService) ListGroupEvents(userID, groupID, limit, offset int) ([]domain.GroupEvent, error) {
	isMember, err := s.groupRepo.IsUserInGroup(groupID, userID)
	if err != nil {
		s.logger.Error("Failed to check group membership", "error", err, "groupID", groupID, "userID", userID)
		return nil, fmt.Errorf("failed to list group events")
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of this group")
	}

	limit, offset = s.validateLimitAndOffset(limit, offset)

	events, err := s.groupRepo.ListGroupEvents(groupID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list group events", "error", err, "groupID", groupID, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to list group events")
	}

	return events, nil
}

func (s *GroupService) SetGroupEventRSVP(userID, groupID, eventID int, response string) (*domain.GroupEventRSVP, error) {
	isMember, err := s.groupRepo.IsUserInGroup(groupID, userID)
	if err != nil {
		s.logger.Error("Failed to check group membership", "error", err, "groupID", groupID, "userID", userID)
		return nil, fmt.Errorf("failed to set group event rsvp")
	}
	if !isMember {
		return nil, fmt.Errorf("user is not a member of this group")
	}

	event, err := s.groupRepo.GetGroupEventByID(eventID)
	if err != nil {
		s.logger.Error("Failed to get group event by ID", "error", err, "eventID", eventID)
		return nil, fmt.Errorf("failed to get group event")
	}
	if event.GroupID != groupID {
		return nil, fmt.Errorf("event does not belong to this group")
	}

	if response != "going" && response != "not_going" {
		return nil, fmt.Errorf("invalid rsvp response")
	}

	err = s.groupRepo.SetGroupEventRSVP(eventID, userID, response)
	if err != nil {
		s.logger.Error("Failed to set group event rsvp", "error", err, "eventID", eventID, "userID", userID)
		return nil, fmt.Errorf("failed to set group event rsvp")
	}

	rsvp, err := s.groupRepo.GetGroupEventRSVP(eventID, userID)
	if err != nil {
		s.logger.Error("Failed to get group event rsvp", "error", err, "eventID", eventID, "userID", userID)
		return nil, fmt.Errorf("failed to get group event rsvp")
	}

	return rsvp, nil
}

func (s *GroupService) validateGroupEvent(eventData domain.CreateGroupEventRequest) error {
	if eventData.Title == "" || len(eventData.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters")
	}
	if eventData.Description == "" {
		return fmt.Errorf("description is required")
	}
	if eventData.StartsAt.IsZero() {
		return fmt.Errorf("starts_at is required")
	}
	if eventData.StartsAt.Before(time.Now()) {
		return fmt.Errorf("starts_at must be in the future")
	}
	return nil
}
