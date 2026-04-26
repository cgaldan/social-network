package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/domain"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) CreateGroup(group *domain.Group) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO groups (
			creator_id, 
			title, 
			description, 
			conversation_id
		)
		VALUES (?, ?, ?, ?)`,
		group.CreatorID,
		group.Title,
		group.Description,
		group.ConversationID,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to create group: %w", err)
	}

	return result.LastInsertId()
}

func (r *GroupRepository) GetGroupByID(groupID int) (*domain.Group, error) {
	group := &domain.Group{}
	err := r.db.QueryRow(`
		SELECT id, creator_id, title, description, conversation_id, created_at
		FROM groups
		WHERE id = ?`,
		groupID,
	).Scan(
		&group.ID,
		&group.CreatorID,
		&group.Title,
		&group.Description,
		&group.ConversationID,
		&group.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("group not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	return group, nil
}

func (r *GroupRepository) AddMember(groupID, userID int, role string) error {
	_, err := r.db.Exec(`
		INSERT INTO group_members (
			group_id,
			user_id,
			role
		)
		VALUES (?, ?, ?)`,
		groupID,
		userID,
		role,
	)

	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}

func (r *GroupRepository) RemoveMember(groupID, userID int) error {
	_, err := r.db.Exec(`
		DELETE FROM group_members
		WHERE group_id = ? AND user_id = ?`,
		groupID,
		userID,
	)
	return err
}

func (r *GroupRepository) GetMembersByGroupID(groupID int) ([]domain.GroupMember, error) {
	rows, err := r.db.Query(`
		SELECT id, group_id, user_id, role, joined_at
		FROM group_members
		WHERE group_id = ?`,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	defer rows.Close()

	members := []domain.GroupMember{}
	for rows.Next() {
		var member domain.GroupMember
		err := rows.Scan(&member.ID, &member.GroupID, &member.UserID, &member.Role, &member.JoinedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, member)
	}
	return members, nil
}

func (r *GroupRepository) CreateGroupInvitation(groupID, inviterID, inviteeID int) error {
	_, err := r.db.Exec(`
		INSERT INTO group_invitations (
			group_id,
			inviter_id,
			invitee_id
		)
		VALUES (?, ?, ?)`,
		groupID,
		inviterID,
		inviteeID,
	)

	if err != nil {
		return fmt.Errorf("failed to create group invitation: %w", err)
	}
	return nil
}

func (r *GroupRepository) CreateGroupJoinRequest(groupID, userID int) error {
	_, err := r.db.Exec(`
		INSERT INTO group_join_requests (
			group_id,
			user_id
		)
		VALUES (?, ?)`,
		groupID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to create group join request: %w", err)
	}
	return nil
}

func (r *GroupRepository) GetGroupInvitationByID(invitationID int) (*domain.GroupInvitation, error) {
	invitation := &domain.GroupInvitation{}
	err := r.db.QueryRow(`
		SELECT id, group_id, inviter_id, invitee_id, status, created_at
		FROM group_invitations
		WHERE id = ?`,
		invitationID,
	).Scan(&invitation.ID, &invitation.GroupID, &invitation.InviterID, &invitation.InviteeID, &invitation.Status, &invitation.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("group invitation not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get group invitation: %w", err)
	}

	return invitation, nil
}

func (r *GroupRepository) GetGroupJoinRequestByID(requestID int) (*domain.GroupJoinRequest, error) {
	request := &domain.GroupJoinRequest{}
	err := r.db.QueryRow(`
		SELECT id, group_id, user_id, status, created_at
		FROM group_join_requests
		WHERE id = ?`,
		requestID,
	).Scan(&request.ID, &request.GroupID, &request.UserID, &request.Status, &request.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("group join request not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get group join request: %w", err)
	}

	return request, nil
}

func (r *GroupRepository) GetGroupInvitationsByGroupID(groupID int) ([]domain.GroupInvitation, error) {
	rows, err := r.db.Query(`
		SELECT id, group_id, inviter_id, invitee_id, status, created_at
		FROM group_invitations
		WHERE group_id = ?`,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get group invitations: %w", err)
	}
	defer rows.Close()

	invitations := []domain.GroupInvitation{}
	for rows.Next() {
		var invitation domain.GroupInvitation
		err := rows.Scan(&invitation.ID, &invitation.GroupID, &invitation.InviterID, &invitation.InviteeID, &invitation.Status, &invitation.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, invitation)
	}
	return invitations, nil
}

func (r *GroupRepository) GetGroupJoinRequestsByGroupID(groupID int) ([]domain.GroupJoinRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, group_id, user_id, status, created_at
		FROM group_join_requests
		WHERE group_id = ?`,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get group join requests: %w", err)
	}
	defer rows.Close()

	requests := []domain.GroupJoinRequest{}
	for rows.Next() {
		var request domain.GroupJoinRequest
		err := rows.Scan(&request.ID, &request.GroupID, &request.UserID, &request.Status, &request.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func (r *GroupRepository) UpdateGroupInvitationStatus(invitationID int, status string) error {
	_, err := r.db.Exec(`
		UPDATE group_invitations
		SET status = ?
		WHERE id = ?`,
		status,
		invitationID,
	)
	if err != nil {
		return fmt.Errorf("failed to update group invitation status: %w", err)
	}
	return nil
}

func (r *GroupRepository) UpdateGroupJoinRequestStatus(requestID int, status string) error {
	_, err := r.db.Exec(`
		UPDATE group_join_requests
		SET status = ?
		WHERE id = ?`,
		status,
		requestID,
	)
	if err != nil {
		return fmt.Errorf("failed to update group join request status: %w", err)
	}
	return nil
}

func (r *GroupRepository) DeleteGroupInvitation(invitationID int) error {
	_, err := r.db.Exec(`
		DELETE FROM group_invitations
		WHERE id = ?`,
		invitationID,
	)
	return err
}

func (r *GroupRepository) DeleteGroupJoinRequest(requestID int) error {
	_, err := r.db.Exec(`
		DELETE FROM group_join_requests
		WHERE id = ?`,
		requestID,
	)
	return err
}

func (r *GroupRepository) IsUserInGroup(groupID, userID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM group_members
		WHERE group_id = ? AND user_id = ?`,
		groupID,
		userID,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check if user is in group: %w", err)
	}

	return count > 0, nil
}

func (r *GroupRepository) IsUserAdmin(groupID, userID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM group_members
		WHERE group_id = ? AND user_id = ? AND role = 'admin'`,
		groupID,
		userID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is admin: %w", err)
	}

	return count > 0, nil
}
