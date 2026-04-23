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
