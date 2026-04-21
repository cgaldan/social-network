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

func (r *MessageRepository) SendMessage(message *domain.Message) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO messages (conversation_id, sender_id, content)
		VALUES (?, ?, ?)`, message.ConversationID, message.SenderID, message.Content)

	if err != nil {
		return 0, fmt.Errorf("failed to create message: %w", err)
	}

	return result.LastInsertId()
}
