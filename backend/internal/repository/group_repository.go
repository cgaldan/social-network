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
