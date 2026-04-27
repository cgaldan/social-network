package service

import (
	"social-network/internal/domain"
	"social-network/internal/repository"
	"social-network/packages/logger"
)

type Services struct {
	Auth         AuthServiceInterface
	Content      ContentServiceInterface
	Post         PostServiceInterface
	Comment      CommentServiceInterface
	Follow       FollowServiceInterface
	Message      MessageServiceInterface
	Conversation ConversationServiceInterface
	Group        GroupServiceInterface
	Notification NotificationServiceInterface
}

func NewServices(repos *repository.Repositories, logger *logger.Logger, notificationPushers ...NotificationPusher) *Services {
	var notificationPusher NotificationPusher
	if len(notificationPushers) > 0 {
		notificationPusher = notificationPushers[0]
	}

	authService := NewAuthService(repos.User, repos.Session, logger)
	contentService := NewContentService(repos.Post, repos.Group, logger)
	postService := NewPostService(repos.Post, repos.Group, logger)
	commentService := NewCommentService(repos.Comment, repos.Post, repos.Group, logger)
	followService := NewFollowService(repos.Follow, repos.User, logger)
	messageService := NewMessageService(repos.Message, repos.User, repos.Conversation, logger)
	conversationService := NewConversationService(repos.Conversation, repos.Follow, logger)
	groupService := NewGroupService(repos.Group, conversationService, logger)
	notificationService := NewNotificationService(repos.Notification, logger, notificationPusher)

	return &Services{
		Auth:         authService,
		Content:      contentService,
		Post:         postService,
		Comment:      commentService,
		Follow:       followService,
		Message:      messageService,
		Conversation: conversationService,
		Group:        groupService,
		Notification: notificationService,
	}
}

type AuthServiceInterface interface {
	Register(registrationData domain.RegisterRequest) (*domain.User, string, error)
	Login(loginData domain.LoginRequest) (*domain.User, string, error)
	Logout(sessionID string) error
	ValidateSession(sessionID string) (*domain.User, error)
}

type ContentServiceInterface interface {
	CreatePost(userID int, postData domain.CreatePostRequest) (*domain.Post, error)
	CreateGroupPost(userID, groupID int, postData domain.CreatePostRequest) (*domain.Post, error)
}

type PostServiceInterface interface {
	GetPostByID(userID, postID int) (*domain.Post, error)
	ListPosts(category string, limit, offset int) ([]domain.Post, error)
	ListPostsByGroupID(userID, groupID, limit, offset int) ([]domain.Post, error)
	GetPostsByUserID(userID int, limit, offset int) ([]domain.Post, error)
}

type CommentServiceInterface interface {
	CreateComment(userID int, postID int, commentData domain.CreateCommentRequest) (*domain.Comment, error)
	GetCommentsByPostID(userID, postID int) ([]domain.Comment, error)
	GetCommentsByUserID(userID, limit, offset int) ([]domain.Comment, error)
}

type FollowServiceInterface interface {
	FollowUser(followData domain.FollowRequest) (status string, err error)
}

type MessageServiceInterface interface {
	SendMessage(convID, senderID int, content string) (*domain.Message, error)
}

type ConversationServiceInterface interface {
	CreateDirectConversation(convData domain.DirectConversationRequest) (*domain.Conversation, error)
	CreateGroupConversation(name string, initialUserIDs ...int) (*domain.Conversation, error)
	AddConversationParticipant(convID, userID int) error
	RemoveConversationParticipant(convID, userID int) error
}

type GroupServiceInterface interface {
	CreateGroup(group *domain.Group) (*domain.Group, error)
	ListGroups(limit, offset int) ([]domain.Group, error)
	GetMembersByGroupID(groupID int) ([]domain.GroupMember, error)
	AddMember(convID, groupID, userID int, role string) error
	RemoveMember(convID, groupID, userID int) error

	CreateGroupInvitation(groupID, inviterID, inviteeID int) error
	CreateGroupJoinRequest(groupID, userID int) error
	AcceptGroupInvitation(userID int, invitation *domain.GroupInvitation) error
	AcceptGroupJoinRequest(answererID int, request *domain.GroupJoinRequest) error
	DeclineGroupInvitation(userID int, invitation *domain.GroupInvitation) error
	DeclineGroupJoinRequest(answererID int, request *domain.GroupJoinRequest) error

	GetGroupInvitationByID(invitationID int) (*domain.GroupInvitation, error)
	GetGroupJoinRequestByID(requestID int) (*domain.GroupJoinRequest, error)
	GetGroupInvitationsByGroupID(groupID int) ([]domain.GroupInvitation, error)
	GetGroupJoinRequestsByGroupID(groupID int) ([]domain.GroupJoinRequest, error)

	CreateGroupEvent(userID, groupID int, eventData domain.CreateGroupEventRequest) (*domain.GroupEvent, error)
	ListGroupEvents(userID, groupID, limit, offset int) ([]domain.GroupEvent, error)
	SetGroupEventRSVP(userID, groupID, eventID int, response string) (*domain.GroupEventRSVP, error)
}

type NotificationServiceInterface interface {
	CreateNotification(input domain.CreateNotificationRequest) (*domain.Notification, error)
	ListNotifications(userID, limit, offset int) ([]domain.Notification, error)
	CountUnread(userID int) (int, error)
	MarkRead(userID, notificationID int) error
	MarkAllRead(userID int) error
}
