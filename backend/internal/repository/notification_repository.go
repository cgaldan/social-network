package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/domain"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(notification *domain.Notification) (*domain.Notification, error) {
	result, err := r.db.Exec(`
		INSERT INTO notifications (
			recipient_id,
			actor_id,
			type,
			title,
			body,
			entity_type,
			entity_id,
			action_url,
			metadata,
			read_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		notification.RecipientID,
		notification.ActorID,
		notification.Type,
		notification.Title,
		notification.Body,
		notification.EntityType,
		notification.EntityID,
		notification.ActionURL,
		notification.Metadata,
		notification.ReadAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	notificationID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get notification id: %w", err)
	}

	return r.GetNotificationByID(int(notificationID), notification.RecipientID)
}

func (r *NotificationRepository) ListNotifications(recipientID, limit, offset int) ([]domain.Notification, error) {
	rows, err := r.db.Query(`
		SELECT id, recipient_id, actor_id, type, title, body, entity_type, entity_id, action_url, metadata, read_at, created_at
		FROM notifications
		WHERE recipient_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?`,
		recipientID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	notifications := []domain.Notification{}
	for rows.Next() {
		notification, err := scanNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, *notification)
	}

	return notifications, nil
}

func (r *NotificationRepository) CountUnreadNotifications(recipientID int) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM notifications
		WHERE recipient_id = ? AND read_at IS NULL`,
		recipientID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}

	return count, nil
}

func (r *NotificationRepository) MarkNotificationRead(notificationID, recipientID int) error {
	result, err := r.db.Exec(`
		UPDATE notifications
		SET read_at = COALESCE(read_at, CURRENT_TIMESTAMP)
		WHERE id = ? AND recipient_id = ?`,
		notificationID,
		recipientID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark notification read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check notification update: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

func (r *NotificationRepository) MarkAllNotificationsRead(recipientID int) error {
	_, err := r.db.Exec(`
		UPDATE notifications
		SET read_at = COALESCE(read_at, CURRENT_TIMESTAMP)
		WHERE recipient_id = ? AND read_at IS NULL`,
		recipientID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications read: %w", err)
	}

	return nil
}

func (r *NotificationRepository) GetNotificationByID(notificationID, recipientID int) (*domain.Notification, error) {
	row := r.db.QueryRow(`
		SELECT id, recipient_id, actor_id, type, title, body, entity_type, entity_id, action_url, metadata, read_at, created_at
		FROM notifications
		WHERE id = ? AND recipient_id = ?`,
		notificationID,
		recipientID,
	)

	notification, err := scanNotification(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("notification not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	return notification, nil
}

type notificationScanner interface {
	Scan(dest ...interface{}) error
}

func scanNotification(scanner notificationScanner) (*domain.Notification, error) {
	var notification domain.Notification
	var actorID sql.NullInt64
	var entityType sql.NullString
	var entityID sql.NullInt64
	var actionURL sql.NullString
	var metadata sql.NullString
	var readAt sql.NullTime

	err := scanner.Scan(
		&notification.ID,
		&notification.RecipientID,
		&actorID,
		&notification.Type,
		&notification.Title,
		&notification.Body,
		&entityType,
		&entityID,
		&actionURL,
		&metadata,
		&readAt,
		&notification.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if actorID.Valid {
		value := int(actorID.Int64)
		notification.ActorID = &value
	}
	if entityType.Valid {
		notification.EntityType = &entityType.String
	}
	if entityID.Valid {
		value := int(entityID.Int64)
		notification.EntityID = &value
	}
	if actionURL.Valid {
		notification.ActionURL = &actionURL.String
	}
	if metadata.Valid {
		notification.Metadata = &metadata.String
	}
	if readAt.Valid {
		notification.ReadAt = &readAt.Time
	}

	return &notification, nil
}
