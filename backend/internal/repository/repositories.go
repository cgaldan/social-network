package repository

import (
	"database/sql"
	"social-network/internal/domain"
)

type Repositories struct {
	User UserRepositoryInterface
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
	}
}

type UserRepositoryInterface interface {
	CreateUser(email, passwordHash, firstName, lastName string, dateOfBirth int, nickname, gender, avatar_path, aboutMe string, isPublic bool) (int64, error)
	GetUserByID(userID int) (*domain.User, error)
	GetUserByIdentifier(identifier string) (*domain.User, string, error)
	UpdateLastSeen(userID int) error
}
