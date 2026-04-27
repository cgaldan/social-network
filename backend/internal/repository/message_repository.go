package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/domain"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateMessage(message *domain.Message) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO messages (conversation_id, sender_id, content)
		VALUES (?, ?, ?)`, message.ConversationID, message.SenderID, message.Content)

	if err != nil {
		return 0, fmt.Errorf("failed to create message: %w", err)
	}

	return result.LastInsertId()
}

func (r *MessageRepository) GetMessageByID(messageID int) (*domain.Message, error) {
	var message domain.Message
	err := r.db.QueryRow(`
		SELECT 
			m.id, 
			m.conversation_id, 
			m.sender_id, 
			m.content, 
			m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.id = ?`, messageID,
	).Scan(
		&message.ID,
		&message.ConversationID,
		&message.SenderID,
		&message.Content,
		&message.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("message not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil
}
