package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/domain"
)

type FollowerRepository struct {
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) *FollowerRepository {
	return &FollowerRepository{db: db}
}

func (r *FollowerRepository) CreateFollower(followerID, followingID int, status string) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO followers (
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

func (r *FollowerRepository) GetFollowerByID(followerID int) (*domain.Follower, error) {
	var follower domain.Follower
	err := r.db.QueryRow(`
		SELECT 
			id,
			follower_id,
			following_id,
			status,
			created_at
		FROM followers 
		WHERE id = ?`, followerID,
	).Scan(
		&follower.ID,
		&follower.FollowerID,
		&follower.FollowingID,
		&follower.Status,
		&follower.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("follower not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get follower: %w", err)
	}

	return &follower, nil
}

func (r *FollowerRepository) GetFollowersByUserID(userID int, limit, offset int) ([]domain.Follower, error) {
	rows, err := r.db.Query(`
		SELECT 
			id,
			follower_id,
			following_id,
			status,
			created_at
		FROM followers 
		WHERE following_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}
	defer rows.Close()

	var followers []domain.Follower
	for rows.Next() {
		var follower domain.Follower
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

func (r *FollowerRepository) GetFollowingByUserID(userID int, limit, offset int) ([]domain.Follower, error) {
	rows, err := r.db.Query(`
		SELECT 
			id,
			follower_id,
			following_id,
			status,
			created_at
		FROM followers 
		WHERE follower_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	defer rows.Close()

	var following []domain.Follower
	for rows.Next() {
		var follower domain.Follower
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

func (r *FollowerRepository) UpdateFollowerStatus(followerID int, status string) error {
	result, err := r.db.Exec(`
		UPDATE followers
		SET status = ?
		WHERE id = ?`, status, followerID)

	if err != nil {
		return fmt.Errorf("failed to update follower status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("follower not found")
	}

	return nil
}

func (r *FollowerRepository) DeleteFollower(followerID int) error {
	result, err := r.db.Exec(`
		DELETE FROM followers
		WHERE id = ?`, followerID)

	if err != nil {
		return fmt.Errorf("failed to delete follower: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("follower not found")
	}

	return nil
}

func (r *FollowerRepository) FollowExists(followerID, followingID int) (bool, error) {
	var id int
	err := r.db.QueryRow(`
		SELECT id FROM followers 
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

func (r *FollowerRepository) GetFollowStatus(followerID, followingID int) (string, error) {
	var status string
	err := r.db.QueryRow(`
		SELECT status FROM followers 
		WHERE follower_id = ? AND following_id = ?`,
		followerID, followingID,
	).Scan(&status)

	if err == sql.ErrNoRows {
		return "", fmt.Errorf("follow relationship not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get follow status: %w", err)
	}

	return status, nil
}
