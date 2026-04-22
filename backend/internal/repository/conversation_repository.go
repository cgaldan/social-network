package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/domain"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}
func (r *ConversationRepository) CreateDirectConversation(userID1, userID2 int) (*domain.Conversation, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO conversations (type) 
		VALUES ('private')`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	convID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO conversation_participants (conversation_id, user_id) 
		VALUES (?, ?), (?, ?)`,
		convID, userID1, convID, userID2,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add participants: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &domain.Conversation{
		ID:   int(convID),
		Type: "private",
	}, nil
}

func (r *ConversationRepository) GetDirectConversation(userID1, userID2 int) (*domain.Conversation, error) {
	var conversation domain.Conversation
	err := r.db.QueryRow(`
		SELECT 
			con.id, 
			con.name, 
			con.type, 
			con.created_at
		FROM conversations con
		JOIN conversation_participants cp1 ON con.id = cp1.conversation_id AND cp1.user_id = ?
		JOIN conversation_participants cp2 ON con.id = cp2.conversation_id AND cp2.user_id = ?
		WHERE con.type = 'private'`,
		userID1, userID2,
	).Scan(
		&conversation.ID,
		&conversation.Name,
		&conversation.Type,
		&conversation.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get direct conversation: %w", err)
	}

	return &conversation, nil
}

func (r *ConversationRepository) IsUserInConversation(conversationID, userID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM conversation_participants 
		WHERE conversation_id = ? AND user_id = ?`, conversationID, userID,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
