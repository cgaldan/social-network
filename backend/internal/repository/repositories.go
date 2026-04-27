package repository

import (
	"database/sql"
	"social-network/internal/domain"
	"time"
)

type Repositories struct {
	User         UserRepositoryInterface
	Session      SessionRepositoryInterface
	Post         PostRepositoryInterface
	Comment      CommentRepositoryInterface
	Follow       FollowRepositoryInterface
	Message      MessageRepositoryInterface
	Conversation ConversationRepositoryInterface
	Group        GroupRepositoryInterface
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		Session:      NewSessionRepository(db),
		Post:         NewPostRepository(db),
		Comment:      NewCommentRepository(db),
		Follow:       NewFollowRepository(db),
		Message:      NewMessageRepository(db),
		Conversation: NewConversationRepository(db),
		Group:        NewGroupRepository(db),
	}
}

type UserRepositoryInterface interface {
	CreateUser(email, passwordHash, firstName, lastName string, dateOfBirth time.Time, nickname, gender, avatar_path, aboutMe string, isPublic bool) (int64, error)
	GetUserByID(userID int) (*domain.User, error)
	GetUserByIdentifier(identifier string) (*domain.User, string, error)
	UpdateLastSeen(userID int) error
	GetUserPrivacyByUserID(userID int) (bool, error)
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

type FollowRepositoryInterface interface {
	CreateFollow(followerID, followingID int, status string) (int64, error)
	GetFollowByID(followerID int) (*domain.Follow, error)
	GetFollowersByUserID(userID int, limit, offset int) ([]domain.Follow, error)
	GetFollowingByUserID(userID int, limit, offset int) ([]domain.Follow, error)
	UpdateFollowStatus(followerID int, status string) error
	DeleteFollow(followerID int) error
	FollowExists(followerID, followingID int) (bool, error)
	GetFollowStatusByFollowID(followID int) (string, error)
	EitherUserFollows(userID1, userID2 int) (bool, error)
}

type ConversationRepositoryInterface interface {
	IsUserInConversation(conversationID, userID int) (bool, error)
	CreateDirectConversation(userID1, userID2 int) (*domain.Conversation, error)
	GetDirectConversation(userID1, userID2 int) (*domain.Conversation, error)

	CreateGroupConversation(name string, initialUserIDs ...int) (*domain.Conversation, error)
	GetGroupConversationByID(conversationID int) (*domain.Conversation, error)
	AddConversationParticipant(conversationID, userID int) error
	RemoveConversationParticipant(conversationID, userID int) error
}

type MessageRepositoryInterface interface {
	CreateMessage(message *domain.Message) (int64, error)
	GetMessageByID(messageID int) (*domain.Message, error)
}

type GroupRepositoryInterface interface {
	CreateGroup(group *domain.Group) (int64, error)
	GetGroupByID(groupID int) (*domain.Group, error)
	ListGroups(limit, offset int) ([]domain.Group, error)

	AddMember(groupID, userID int, role string) error
	RemoveMember(groupID, userID int) error
	GetMembersByGroupID(groupID int) ([]domain.GroupMember, error)

	CreateGroupInvitation(groupID, inviterID, inviteeID int) error
	CreateGroupJoinRequest(groupID, userID int) error

	GetGroupInvitationByID(invitationID int) (*domain.GroupInvitation, error)
	GetGroupJoinRequestByID(requestID int) (*domain.GroupJoinRequest, error)

	GetGroupInvitationsByGroupID(groupID int) ([]domain.GroupInvitation, error)
	GetGroupJoinRequestsByGroupID(groupID int) ([]domain.GroupJoinRequest, error)

	UpdateGroupInvitationStatus(invitationID int, status string) error
	UpdateGroupJoinRequestStatus(requestID int, status string) error

	DeleteGroupInvitation(invitationID int) error
	DeleteGroupJoinRequest(requestID int) error

	IsUserInGroup(groupID, userID int) (bool, error)
	IsUserAdmin(groupID, userID int) (bool, error)
}
