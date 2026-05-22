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

// ListMessagesByConversationID returns messages newest-first (DESC).
// Pagination is cursor-based: pass beforeID=0 for the newest page, or pass the
// smallest id you currently hold to fetch the next older page. Offset-based
// pagination is wrong for chat — new messages get inserted at the head between
// fetches, shifting offsets and producing duplicates or gaps.
func (r *MessageRepository) ListMessagesByConversationID(conversationID, limit, beforeID int) ([]domain.Message, error) {
	const query = `
		SELECT id, conversation_id, sender_id, content, created_at
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
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
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
