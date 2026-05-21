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

func (r *ConversationRepository) GetParticipantIDs(conversationID int) ([]int, error) {
	rows, err := r.db.Query(`
		SELECT user_id
		FROM conversation_participants
		WHERE conversation_id = ?`, conversationID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get participant ids: %w", err)
	}
	defer rows.Close()

	ids := make([]int, 0)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan participant id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *ConversationRepository) ListConversationsByUserID(userID, limit, offset int) ([]domain.ConversationSummary, error) {
	rows, err := r.db.Query(`
		SELECT
			c.id,
			c.name,
			c.type,
			c.created_at,
			(SELECT m.id FROM messages m WHERE m.conversation_id = c.id ORDER BY m.created_at DESC, m.id DESC LIMIT 1),
			(SELECT m.sender_id FROM messages m WHERE m.conversation_id = c.id ORDER BY m.created_at DESC, m.id DESC LIMIT 1),
			(SELECT m.content FROM messages m WHERE m.conversation_id = c.id ORDER BY m.created_at DESC, m.id DESC LIMIT 1),
			(SELECT m.created_at FROM messages m WHERE m.conversation_id = c.id ORDER BY m.created_at DESC, m.id DESC LIMIT 1),
			(SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id AND m.sender_id != ? AND (cp.last_read_at IS NULL OR m.created_at > cp.last_read_at))
		FROM conversations c
		JOIN conversation_participants cp ON cp.conversation_id = c.id
		WHERE cp.user_id = ?
		ORDER BY COALESCE(
			(SELECT m.created_at FROM messages m WHERE m.conversation_id = c.id ORDER BY m.created_at DESC LIMIT 1),
			c.created_at
		) DESC
		LIMIT ? OFFSET ?`,
		userID, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	out := make([]domain.ConversationSummary, 0)
	for rows.Next() {
		var c domain.ConversationSummary
		var (
			msgID        sql.NullInt64
			msgSenderID  sql.NullInt64
			msgContent   sql.NullString
			msgCreatedAt sql.NullTime
		)
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Type, &c.CreatedAt,
			&msgID, &msgSenderID, &msgContent, &msgCreatedAt,
			&c.UnreadCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		if msgID.Valid {
			c.LastMessage = &domain.Message{
				ID:             int(msgID.Int64),
				ConversationID: c.ID,
				SenderID:       int(msgSenderID.Int64),
				Content:        msgContent.String,
				CreatedAt:      msgCreatedAt.Time,
			}
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range out {
		participants, err := r.GetParticipants(out[i].ID)
		if err != nil {
			return nil, err
		}
		out[i].Participants = participants
	}
	return out, nil
}

func (r *ConversationRepository) GetParticipants(conversationID int) ([]domain.ConversationParticipant, error) {
	rows, err := r.db.Query(`
		SELECT
			u.id,
			u.nickname,
			u.first_name,
			u.last_name,
			u.avatar_path,
			EXISTS(SELECT 1 FROM sessions s WHERE s.user_id = u.id AND s.expires_at > NOW())
		FROM conversation_participants cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.conversation_id = ?`, conversationID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	defer rows.Close()

	participants := make([]domain.ConversationParticipant, 0)
	for rows.Next() {
		var p domain.ConversationParticipant
		if err := rows.Scan(&p.UserID, &p.Nickname, &p.FirstName, &p.LastName, &p.AvatarPath, &p.IsOnline); err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		participants = append(participants, p)
	}
	return participants, rows.Err()
}

func makePairKey(userID1, userID2 int) string {
	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}
	return fmt.Sprintf("%d-%d", userID2, userID1)
}

func (r *ConversationRepository) GetConversationByID(conversationID int) (*domain.Conversation, error) {
	var c domain.Conversation
	err := r.db.QueryRow(`
		SELECT id, name, type, created_at
		FROM conversations
		WHERE id = ?`, conversationID,
	).Scan(&c.ID, &c.Name, &c.Type, &c.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("conversation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	return &c, nil
}

func (r *ConversationRepository) MarkConversationRead(conversationID, userID int) error {
	res, err := r.db.Exec(`
		UPDATE conversation_participants
		SET last_read_at = CURRENT_TIMESTAMP
		WHERE conversation_id = ? AND user_id = ?`,
		conversationID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark conversation read: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to mark conversation read: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("user is not part of the conversation")
	}
	return nil
}
