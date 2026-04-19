package repository

import (
	"database/sql"
	"social-network/internal/domain"
	"time"
)

type Repositories struct {
	User     UserRepositoryInterface
	Session  SessionRepositoryInterface
	Post     PostRepositoryInterface
	Comment  CommentRepositoryInterface
	Follower FollowerRepositoryInterface
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:     NewUserRepository(db),
		Session:  NewSessionRepository(db),
		Post:     NewPostRepository(db),
		Comment:  NewCommentRepository(db),
		Follower: NewFollowerRepository(db),
	}
}

type UserRepositoryInterface interface {
	CreateUser(email, passwordHash, firstName, lastName string, dateOfBirth time.Time, nickname, gender, avatar_path, aboutMe string, isPublic bool) (int64, error)
	GetUserByID(userID int) (*domain.User, error)
	GetUserByIdentifier(identifier string) (*domain.User, string, error)
	UpdateLastSeen(userID int) error
}

type SessionRepositoryInterface interface {
	CreateSession(sessionID string, userID int, expiresAt time.Time) error
	GetSessionBySessionID(sessionID string) (*domain.Session, error)
	DeleteSession(sessionID string) error
}

type PostRepositoryInterface interface {
	CreatePost(userID int, title, content, category, privacyLevel, mediaURL string) (int64, error)
	GetPostByID(postID int) (*domain.Post, error)
	ListPosts(category string, limit, offset int) ([]domain.Post, error)
	GetPostsByUserID(userID int, limit, offset int) ([]domain.Post, error)
	PostExists(postID int) (bool, error)
}

type CommentRepositoryInterface interface {
	CreateComment(userID, postID int, content, mediaURL string) (int64, error)
	GetCommentsByPostID(postID int) ([]domain.Comment, error)
	GetCommentByID(commentID int) (*domain.Comment, error)
	GetCommentsByUserID(userID int, limit, offset int) ([]domain.Comment, error)
}

type FollowerRepositoryInterface interface {
	CreateFollower(followerID, followingID int, status string) (int64, error)
	GetFollowerByID(followerID int) (*domain.Follower, error)
	GetFollowersByUserID(userID int, limit, offset int) ([]domain.Follower, error)
	GetFollowingByUserID(userID int, limit, offset int) ([]domain.Follower, error)
	UpdateFollowerStatus(followerID int, status string) error
	DeleteFollower(followerID int) error
	FollowExists(followerID, followingID int) (bool, error)
	GetFollowStatus(followerID, followingID int) (string, error)
}
