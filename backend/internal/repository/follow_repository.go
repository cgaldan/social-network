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
		return 0, fmt.Errorf("failed to create follower: %w", err)
	}

	return result.LastInsertId()
}

func (r *FollowRepository) GetFollowByID(followerID int) (*domain.Follow, error) {
	var follow domain.Follow
	err := r.db.QueryRow(`
		SELECT 
			id,
			follower_id,
			following_id,
			status,
			created_at
		FROM follows 
		WHERE id = ?`, followerID,
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

func (r *FollowRepository) GetFollowersByUserID(userID int, limit, offset int) ([]domain.Follow, error) {
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
		LIMIT ? OFFSET ?`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	defer rows.Close()

	var followers []domain.Follow
	for rows.Next() {
		var follower domain.Follow
		err := rows.Scan(
			&follower.ID,
			&follower.FollowerID,
			&follower.FollowingID,
			&follower.Status,
			&follower.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follower: %w", err)
		}
		followers = append(followers, follower)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating followers: %w", err)
	}

	return followers, nil
}

func (r *FollowRepository) GetFollowingByUserID(userID int, limit, offset int) ([]domain.Follow, error) {
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
		LIMIT ? OFFSET ?`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	defer rows.Close()

	var following []domain.Follow
	for rows.Next() {
		var follower domain.Follow
		err := rows.Scan(
			&follower.ID,
			&follower.FollowerID,
			&follower.FollowingID,
			&follower.Status,
			&follower.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan following: %w", err)
		}
		following = append(following, follower)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating following: %w", err)
	}

	return following, nil
}

func (r *FollowRepository) UpdateFollowStatus(followerID int, status string) error {
	result, err := r.db.Exec(`
		UPDATE follows
		SET status = ?
		WHERE id = ?`, status, followerID)

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

func (r *FollowRepository) DeleteFollow(followerID int) error {
	result, err := r.db.Exec(`
		DELETE FROM follows
		WHERE id = ?`, followerID)

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

func (r *FollowRepository) FollowExists(followerID, followingID int) (bool, error) {
	var id int
	err := r.db.QueryRow(`
		SELECT id FROM follows 
		WHERE follower_id = ? AND following_id = ?`,
		followerID, followingID,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if follow exists: %w", err)
	}

	return true, nil
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
