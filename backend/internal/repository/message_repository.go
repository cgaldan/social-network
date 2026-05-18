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

func (r *MessageRepository) ListMessagesByConversationID(conversationID, limit, offset int) ([]domain.Message, error) {
	rows, err := r.db.Query(`
		SELECT id, conversation_id, sender_id, content, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?`,
		conversationID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Message, 0)
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
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
