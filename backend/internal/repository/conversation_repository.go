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

func (r *ConversationRepository) CreateConversation(conversation *domain.Conversation) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO conversations (name, type)
		VALUES (?, ?)`, conversation.Name, conversation.Type)

	if err != nil {
		return 0, fmt.Errorf("failed to create conversation: %w", err)
	}

	conversationID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve conversation ID: %w", err)
	}

	return conversationID, nil
}

func (r *ConversationRepository) GetConversationByID(conversationID int) (*domain.Conversation, error) {
	var conversation domain.Conversation
	err := r.db.QueryRow(`
		SELECT 
			con.id, 
			con.name, 
			con.type, 
			con.created_at
		FROM conversations con
		WHERE con.id = ?`, conversationID,
	).Scan(
		&conversation.ID,
		&conversation.Name,
		&conversation.Type,
		&conversation.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("conversation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conversation, nil
}
