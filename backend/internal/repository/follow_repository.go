package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/domain"
)

type FollowRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) CreateFollow(followerID, followingID int, status string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO follows (
			follower_id,
			following_id,
			status
		)
		VALUES (?, ?, ?)`,
		followerID,
		followingID,
		status,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to create follow: %w", err)
	}

	return result.LastInsertId()
}

func (r *FollowRepository) GetFollowByID(followID int) (*domain.Follow, error) {
	var follow domain.Follow
	err := r.db.QueryRow(`
		SELECT 
			id,
			follower_id,
			following_id,
			status,
			created_at
		FROM follows 
		WHERE id = ?`, followID,
	).Scan(
		&follow.ID,
		&follow.FollowerID,
		&follow.FollowingID,
		&follow.Status,
		&follow.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("follow not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get follow: %w", err)
	}

	return &follow, nil
}

func (r *FollowRepository) GetFollowByUsers(followerID, followingID int) (*domain.Follow, error) {
	var follow domain.Follow
	err := r.db.QueryRow(`
		SELECT
			id,
			follower_id,
			following_id,
			status,
			created_at
		FROM follows
		WHERE follower_id = ? AND following_id = ?`,
		followerID,
		followingID,
	).Scan(
		&follow.ID,
		&follow.FollowerID,
		&follow.FollowingID,
		&follow.Status,
		&follow.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get follow by users: %w", err)
	}

	return &follow, nil
}

func (r *FollowRepository) GetFollowRequestsByFollowingID(followingID int, limit, offset int) ([]domain.Follow, error) {
	rows, err := r.db.Query(`
		SELECT 
			id,
			follower_id,
			following_id,
			status,
			created_at
		FROM follows 
		WHERE following_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, followingID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get follow requests: %w", err)
	}
	defer rows.Close()

	var followRequests []domain.Follow
	for rows.Next() {
		var followRequest domain.Follow
		err := rows.Scan(
			&followRequest.ID,
			&followRequest.FollowerID,
			&followRequest.FollowingID,
			&followRequest.Status,
			&followRequest.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follow request: %w", err)
		}
		followRequests = append(followRequests, followRequest)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating follow requests: %w", err)
	}

	return followRequests, nil
}

func (r *FollowRepository) GetFollowRequestsByFollowerID(followerID int, limit, offset int) ([]domain.Follow, error) {
	rows, err := r.db.Query(`
		SELECT 
			id,
			follower_id,
			following_id,
			status,
			created_at
		FROM follows 
		WHERE follower_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, followerID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get follow requests by follower ID: %w", err)
	}
	defer rows.Close()

	var followRequests []domain.Follow
	for rows.Next() {
		var followRequest domain.Follow
		err := rows.Scan(
			&followRequest.ID,
			&followRequest.FollowerID,
			&followRequest.FollowingID,
			&followRequest.Status,
			&followRequest.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follow request: %w", err)
		}
		followRequests = append(followRequests, followRequest)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating follow requests: %w", err)
	}

	return followRequests, nil
}

func (r *FollowRepository) UpdateFollowStatus(followID int, status string) error {
	result, err := r.db.Exec(`
		UPDATE follows
		SET status = ?
		WHERE id = ?`, status, followID)

	if err != nil {
		return fmt.Errorf("failed to update follow status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("follow not found")
	}

	return nil
}

func (r *FollowRepository) DeleteFollow(followID int) error {
	result, err := r.db.Exec(`
		DELETE FROM follows
		WHERE id = ?`, followID)

	if err != nil {
		return fmt.Errorf("failed to delete follow: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("follow not found")
	}

	return nil
}

func (r *FollowRepository) GetFollowStatusByFollowID(followID int) (string, error) {
	var status string
	err := r.db.QueryRow(`
		SELECT status FROM follows 
		WHERE id = ?`, followID).Scan(&status)

	if err == sql.ErrNoRows {

	}

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("follow relationship not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get follow status: %w", err)
	}

	return status, nil
}

func (r *FollowRepository) EitherUserFollows(userID1, userID2 int) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM follows 
		WHERE (follower_id = ? AND following_id = ?)
		OR (follower_id = ? AND following_id = ?)
		AND status = 'accepted'`,
		userID1, userID2, userID2, userID1,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check follow relationship: %w", err)
	}

	return count > 0, nil
}
