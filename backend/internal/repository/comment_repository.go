package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/domain"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(userID, postID int, content, mediaURL string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO comments (
			user_id, 
			post_id, 
			content, 
			media_url
		)
		VALUES (?, ?, ?, ?)`,
		userID,
		postID,
		content,
		mediaURL,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to create comment: %w", err)
	}

	return result.LastInsertId()
}

func (r *CommentRepository) GetCommentsByPostID(postID int) ([]domain.Comment, error) {
	rows, err := r.db.Query(`
		SELECT 
			c.id, 
			c.post_id, 
			c.user_id, 
			c.content, 
			c.media_url, 
			c.created_at, 
			c.updated_at, 
			u.nickname
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC`, postID)

	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.MediaURL,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author); err != nil {
			continue
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) GetCommentByID(commentID int) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.QueryRow(`
		SELECT 
			c.id, 
			c.post_id, 
			c.user_id, 
			c.content, 
			c.media_url, 
			c.created_at, 
			c.updated_at, 
			u.nickname
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?`, commentID,
	).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
		&comment.MediaURL,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Author,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("comment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	return &comment, nil
}

func (r *CommentRepository) GetCommentsByUserID(userID int, limit, offset int) ([]domain.Comment, error) {
	rows, err := r.db.Query(`
		SELECT 
			c.id, 
			c.post_id, 
			c.user_id, 
			c.content, 
			c.media_url, 
			c.created_at, 
			c.updated_at, 
			u.nickname
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.user_id = ?
		ORDER BY c.created_at DESC LIMIT ? OFFSET ?`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get user comments: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.MediaURL,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author); err != nil {
			continue
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) UpdateComment(userID, commentID int, content, mediaURL string) error {
	result, err := r.db.Exec(`
		UPDATE comments
		SET 
		content = ?, 
		media_url = ?, 
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?`,
		content,
		mediaURL,
		commentID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("comment not found")
	}
	return nil
}

func (r *CommentRepository) DeleteComment(userID, commentID int) error {
	result, err := r.db.Exec(`DELETE FROM comments WHERE id = ? AND user_id = ?`, commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("comment not found")
	}
	return nil
}
