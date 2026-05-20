package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

const visibilityClause = `(
	p.privacy_level = 'public'
	OR p.user_id = ?
	OR (p.privacy_level = 'almost_private' AND EXISTS (
		SELECT 1 FROM follows f
		WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
	))
	OR (p.privacy_level = 'private' AND EXISTS (
		SELECT 1 FROM post_audiences pa
		WHERE pa.post_id = p.id AND pa.user_id = ?
	))
)`

func (r *PostRepository) CreatePost(userID int, title, content, category, privacyLevel, mediaURL string, groupID int) (int64, error) {
	var groupIDArg interface{}
	if groupID > 0 {
		groupIDArg = groupID
	}

	result, err := r.db.Exec(`
		INSERT INTO posts (
			user_id, 
			title, 
			content, 
			category, 
			privacy_level, 
			media_url,
			group_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID,
		title,
		content,
		category,
		privacyLevel,
		mediaURL,
		groupIDArg,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to create post: %w", err)
	}

	return result.LastInsertId()
}

func (r *PostRepository) GetPostByID(postID int) (*domain.Post, error) {
	var post domain.Post
	var groupID sql.NullInt64
	err := r.db.QueryRow(`
		SELECT 
			p.id, 
			p.user_id, 
			p.group_id,
			p.title, 
			p.content, 
			p.category, 
			p.privacy_level, 
			p.media_url,
			p.like_count,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count,
			p.created_at, 
			p.updated_at, 
			u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`, postID,
	).Scan(
		&post.ID,
		&post.UserID,
		&groupID,
		&post.Title,
		&post.Content,
		&post.Category,
		&post.PrivacyLevel,
		&post.MediaURL,
		&post.LikeCount,
		&post.CommentCount,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Author,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if groupID.Valid {
		post.GroupID = int(groupID.Int64)
	}

	return &post, nil
}

func (r *PostRepository) ListPosts(category string, viewerID, limit, offset int) ([]domain.Post, error) {
	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = r.db.Query(`
			SELECT
				p.id,
				p.user_id,
				p.title,
				p.content,
				p.category,
				p.privacy_level,
				p.media_url,
				p.like_count,
				(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count,
				p.created_at,
				p.updated_at,
				u.nickname
			FROM posts p
			JOIN users u ON p.user_id = u.id
			WHERE p.category = ? AND p.group_id IS NULL AND `+visibilityClause+`
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?`, category, viewerID, viewerID, viewerID, limit, offset)
	} else {
		rows, err = r.db.Query(`
			SELECT
				p.id,
				p.user_id,
				p.title,
				p.content,
				p.category,
				p.privacy_level,
				p.media_url,
				p.like_count,
				(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count,
				p.created_at,
				p.updated_at,
				u.nickname
			FROM posts p
			JOIN users u ON p.user_id = u.id
			WHERE p.group_id IS NULL AND `+visibilityClause+`
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?`, viewerID, viewerID, viewerID, limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.Category,
			&post.PrivacyLevel,
			&post.MediaURL,
			&post.LikeCount,
			&post.CommentCount,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Author)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) ListPostsByGroupID(groupID, limit, offset int) ([]domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT 
			p.id, 
			p.user_id, 
			p.group_id,
			p.title, 
			p.content, 
			p.category, 
			p.privacy_level, 
			p.media_url,
			p.like_count,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count,
			p.created_at, 
			p.updated_at, 
			u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.group_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`, groupID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to list group posts: %w", err)
	}
	defer rows.Close()

	posts := []domain.Post{}
	for rows.Next() {
		var post domain.Post
		var postGroupID sql.NullInt64
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&postGroupID,
			&post.Title,
			&post.Content,
			&post.Category,
			&post.PrivacyLevel,
			&post.MediaURL,
			&post.LikeCount,
			&post.CommentCount,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Author)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		if postGroupID.Valid {
			post.GroupID = int(postGroupID.Int64)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) GetPostsByUserID(targetUserID, viewerID, limit, offset int) ([]domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT
			p.id,
			p.user_id,
			p.title,
			p.content,
			p.category,
			p.privacy_level,
			p.media_url,
			p.like_count,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count,
			p.created_at,
			p.updated_at,
			u.nickname
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ? AND p.group_id IS NULL AND `+visibilityClause+`
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`, targetUserID, viewerID, viewerID, viewerID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.Category,
			&post.PrivacyLevel,
			&post.MediaURL,
			&post.LikeCount,
			&post.CommentCount,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Author)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) PostExists(postID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)`, postID).Scan(&exists)
	return exists, err
}

func (r *PostRepository) UpdatePost(userID, postID int, title, content, category, privacyLevel, mediaURL string) error {
	result, err := r.db.Exec(`
		UPDATE posts
			SET 
			title = ?,
			content = ?,
			category = ?,
			privacy_level = ?,
			media_url = ?,
			updated_at = CURRENT_TIMESTAMP
			WHERE id = ? AND user_id = ?`,
		title,
		content,
		category,
		privacyLevel,
		mediaURL,
		postID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("post not found")
	}
	return nil
}

func (r *PostRepository) DeletePost(userID, postID int) error {
	result, err := r.db.Exec(`DELETE FROM posts WHERE id = ? AND user_id = ?`, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("post not found")
	}
	return nil
}

func (r *PostRepository) SetPostAudience(postID int, userIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM post_audiences WHERE post_id = ?`, postID); err != nil {
		return fmt.Errorf("failed to clear audience: %w", err)
	}

	if len(userIDs) > 0 {
		stmt, err := tx.Prepare(`INSERT OR IGNORE INTO post_audiences (post_id, user_id) VALUES (?, ?)`)
		if err != nil {
			return fmt.Errorf("failed to prepare audience insert: %w", err)
		}
		defer stmt.Close()
		for _, uid := range userIDs {
			if _, err := stmt.Exec(postID, uid); err != nil {
				return fmt.Errorf("failed to insert audience row: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *PostRepository) IsInPostAudience(postID, userID int) (bool, error) {
	var n int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM post_audiences WHERE post_id = ? AND user_id = ?`, postID, userID).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("failed to check audience: %w", err)
	}
	return n > 0, nil
}
