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
		INSERT INTO conversations (type, pair_key) 
		VALUES ('private', ?)`,
		makePairKey(userID1, userID2),
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
		SELECT id, type, created_at 
		FROM conversations 
		WHERE pair_key = ?`, makePairKey(userID1, userID2),
	).Scan(
		&conversation.ID,
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

func (r *ConversationRepository) CreateGroupConversation(name string, initialUserIDs ...int) (*domain.Conversation, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO conversations (name, type) VALUES (?, 'group')`, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	convID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	for _, userID := range initialUserIDs {
		_, err = tx.Exec(`
			INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?)`, convID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to add participants: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &domain.Conversation{
		ID:   int(convID),
		Name: name,
		Type: "group",
	}, nil
}

func (r *ConversationRepository) GetGroupConversationByID(conversationID int) (*domain.Conversation, error) {
	var conversation domain.Conversation
	err := r.db.QueryRow(`
		SELECT id, name, type, created_at 
		FROM conversations 
		WHERE id = ?`, conversationID,
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
		return nil, fmt.Errorf("failed to get group conversation: %w", err)
	}

	return &conversation, nil
}

func (r *ConversationRepository) AddConversationParticipant(conversationID, userID int) error {
	_, err := r.db.Exec(`
		INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?)`, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}
	return nil
}

func (r *ConversationRepository) RemoveConversationParticipant(conversationID, userID int) error {
	_, err := r.db.Exec(`
		DELETE FROM conversation_participants WHERE conversation_id = ? AND user_id = ?`, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}
	return nil
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

func makePairKey(userID1, userID2 int) string {
	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}
	return fmt.Sprintf("%d-%d", userID2, userID1)
}
