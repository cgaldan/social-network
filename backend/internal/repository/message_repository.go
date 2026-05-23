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

func (r *MessageRepository) ListMessagesByConversationID(conversationID, limit, beforeID int) ([]domain.Message, error) {
	const query = `
		SELECT id, conversation_id, sender_id, content, created_at, updated_at
		FROM messages
		WHERE conversation_id = ?
		  AND (? = 0 OR id < ?)
		ORDER BY id DESC
		LIMIT ?`
	rows, err := r.db.Query(query, conversationID, beforeID, beforeID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Message, 0)
	for rows.Next() {
		var m domain.Message
		var updatedAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.CreatedAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			m.UpdatedAt = &t
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *MessageRepository) GetMessageByID(messageID int) (*domain.Message, error) {
	var message domain.Message
	var updatedAt sql.NullTime
	err := r.db.QueryRow(`
		SELECT
			m.id,
			m.conversation_id,
			m.sender_id,
			m.content,
			m.created_at,
			m.updated_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.id = ?`, messageID,
	).Scan(
		&message.ID,
		&message.ConversationID,
		&message.SenderID,
		&message.Content,
		&message.CreatedAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("message not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	if updatedAt.Valid {
		t := updatedAt.Time
		message.UpdatedAt = &t
	}

	return &message, nil
}

func (r *MessageRepository) UpdateMessage(messageID, senderID int, content string) error {
	result, err := r.db.Exec(`
		UPDATE messages
		SET content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND sender_id = ?`,
		content, messageID, senderID,
	)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("message not found")
	}
	return nil
}

func (r *MessageRepository) DeleteMessage(messageID, senderID int) error {
	result, err := r.db.Exec(`
		DELETE FROM messages WHERE id = ? AND sender_id = ?`,
		messageID, senderID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("message not found")
	}
	return nil
}
