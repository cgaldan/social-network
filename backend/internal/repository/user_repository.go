package repository

import (
	"database/sql"
	"fmt"
	"social-network/internal/domain"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(email, passwordHash, firstName, lastName string, dateOfBirth time.Time, nickname, gender, avatar_path, aboutMe string, isPublic bool) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO users (
			email,
			password_hash,
			first_name, 
			last_name,
			date_of_birth, 
			nickname, 
			gender,
			avatar_path,
			about_me,
			is_public
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		email,
		passwordHash,
		firstName,
		lastName,
		dateOfBirth,
		nickname,
		gender,
		avatar_path,
		aboutMe,
		isPublic,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return result.LastInsertId()
}

func (r *UserRepository) GetUserByID(userID int) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`
		SELECT 
			id, 
			email, 
			first_name,
			last_name, 
			date_of_birth, 
			nickname, 
			gender, 
			avatar_path, 
			about_me,
			following_count, 
			followers_count,
			is_online, 
			is_public, 
			created_at, 
			last_seen
		FROM users 
		WHERE id = ?`, userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Nickname,
		&user.Gender,
		&user.AvatarPath,
		&user.AboutMe,
		&user.FollowingCount,
		&user.FollowersCount,
		&user.IsOnline,
		&user.IsPublic,
		&user.CreatedAt,
		&user.LastSeen,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByIdentifier(identifier string) (*domain.User, string, error) {
	var user domain.User
	var passwordHash string

	err := r.db.QueryRow(`
		SELECT
			id, 
			email, 
			password_hash,
			first_name,
			last_name, 
			date_of_birth, 
			nickname, 
			gender, 
			avatar_path,
			about_me,
			following_count, 
			followers_count,
			is_online, 
			is_public, 
			created_at, 
			last_seen
		FROM users
		WHERE nickname = ? OR email = ?`,
		identifier, identifier,
	).Scan(
		&user.ID,
		&user.Email,
		&passwordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Nickname,
		&user.Gender,
		&user.AvatarPath,
		&user.AboutMe,
		&user.FollowingCount,
		&user.FollowersCount,
		&user.IsOnline,
		&user.IsPublic,
		&user.CreatedAt,
		&user.LastSeen,
	)

	if err == sql.ErrNoRows {
		return nil, "", fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	return &user, passwordHash, nil
}

func (r *UserRepository) UpdateLastSeen(userID int) error {
	_, err := r.db.Exec("UPDATE users SET last_seen = CURRENT_TIMESTAMP WHERE id = ?", userID)
	return err
}
