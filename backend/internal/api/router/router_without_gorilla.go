package router

import (
	"net/http"
	"social-network/internal/api/handlers"
	"social-network/internal/api/middleware"
	"social-network/internal/config"
	"social-network/internal/service"
	"social-network/internal/websocket"
	"social-network/packages/logger"
)

func NewRouterNoGorilla(services *service.Services, config *config.Config, hub *websocket.Hub, logger *logger.Logger) http.Handler {
	mux := http.NewServeMux()

	authHandler := handlers.NewAuthHandler(services.Auth, logger)
	postHandler := handlers.NewPostHandler(services.Post, services.Auth, services.Content, logger)
	commentHandler := handlers.NewCommentHandler(services.Comment, services.Auth, logger)
	websocketHandler := handlers.NewWebSocketHandler(hub, services.Auth, logger, config.Websocket, config.CORS.AllowedOrigins)
	followHandler := handlers.NewFollowHandler(services.Follow, services.Auth, logger)
	conversationHandler := handlers.NewConversationHandler(services.Conversation, services.Auth, logger)
	messageHandler := handlers.NewMessageHandler(services.Message, services.Auth, logger)
	groupHandler := handlers.NewGroupHandler(services.Group, services.Auth, logger)
	notificationHandler := handlers.NewNotificationHandler(services.Notification, services.Auth, logger)
	userHandler := handlers.NewUserHandler(services.Auth, services.Post, services.Follow, logger)
	uploadHandler := handlers.NewUploadHandler(services.Auth, logger, config.Upload)
	healthHandler := handlers.NewHealthHandler("1.0.0")

	// Health
	mux.HandleFunc("GET /health", healthHandler.Health)

	// Auth routes
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
	mux.HandleFunc("GET /api/auth/me", authHandler.GetCurrentUser)
	mux.HandleFunc("PUT /api/auth/me", authHandler.UpdateCurrentUser)
	mux.HandleFunc("DELETE /api/auth/me", authHandler.DeleteCurrentUser)

	// Post routes
	mux.HandleFunc("GET /api/posts", postHandler.GetPosts)
	mux.HandleFunc("POST /api/posts", postHandler.CreatePost)
	mux.HandleFunc("GET /api/posts/{id}", postHandler.GetPostByID)
	mux.HandleFunc("PUT /api/posts/{id}", postHandler.UpdatePost)
	mux.HandleFunc("DELETE /api/posts/{id}", postHandler.DeletePost)
	mux.HandleFunc("GET /api/posts/{id}/comments", commentHandler.ListComments)
	mux.HandleFunc("POST /api/posts/{id}/comments", commentHandler.CreateComment)
	mux.HandleFunc("PUT /api/posts/{id}/comments/{commentId}", commentHandler.UpdateComment)
	mux.HandleFunc("DELETE /api/posts/{id}/comments/{commentId}", commentHandler.DeleteComment)

	// Follow routes
	mux.HandleFunc("GET /api/follow/requests", followHandler.ListFollowRequests)
	mux.HandleFunc("POST /api/follow/{id}", followHandler.FollowUser)
	mux.HandleFunc("POST /api/follow/{id}/accept", followHandler.AcceptFollowRequest)
	mux.HandleFunc("POST /api/follow/{id}/decline", followHandler.DeclineFollowRequest)
	mux.HandleFunc("POST /api/follow/{id}/unfollow", followHandler.UnfollowUser)
	mux.HandleFunc("POST /api/follow/{id}/remove", followHandler.RemoveFollower)

	// User routes
	mux.HandleFunc("GET /api/users", userHandler.ListUsers)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetUser)
	mux.HandleFunc("GET /api/users/{id}/posts", userHandler.ListUserPosts)
	mux.HandleFunc("GET /api/users/{id}/followers", userHandler.ListUserFollowers)
	mux.HandleFunc("GET /api/users/{id}/following", userHandler.ListUserFollowing)

	// Chat routes
	mux.HandleFunc("GET /api/conversations", conversationHandler.ListConversations)
	mux.HandleFunc("POST /api/conversations/direct", conversationHandler.CreateDirectConversation)
	mux.HandleFunc("GET /api/conversations/{id}/messages", messageHandler.ListMessages)
	mux.HandleFunc("POST /api/conversations/{id}/read", conversationHandler.MarkRead)
	mux.HandleFunc("POST /api/messages", messageHandler.SendMessage)
	mux.HandleFunc("PUT /api/messages/{id}", messageHandler.UpdateMessage)
	mux.HandleFunc("DELETE /api/messages/{id}", messageHandler.DeleteMessage)

	// Group routes
	mux.HandleFunc("GET /api/groups", groupHandler.ListGroups)
	mux.HandleFunc("POST /api/groups", groupHandler.CreateGroup)
	mux.HandleFunc("POST /api/groups/join", groupHandler.JoinGroup)
	mux.HandleFunc("POST /api/groups/join/{id}/accept", groupHandler.AcceptGroupJoinRequest)
	mux.HandleFunc("POST /api/groups/join/{id}/decline", groupHandler.DeclineGroupJoinRequest)
	mux.HandleFunc("POST /api/groups/invitations", groupHandler.InviteToGroup)
	mux.HandleFunc("POST /api/groups/invitations/{id}/accept", groupHandler.AcceptGroupInvitation)
	mux.HandleFunc("POST /api/groups/invitations/{id}/decline", groupHandler.DeclineGroupInvitation)
	mux.HandleFunc("GET /api/groups/{id}", groupHandler.GetGroup)
	mux.HandleFunc("PUT /api/groups/{id}", groupHandler.UpdateGroup)
	mux.HandleFunc("DELETE /api/groups/{id}", groupHandler.DeleteGroup)
	mux.HandleFunc("POST /api/groups/{id}/leave", groupHandler.LeaveGroup)
	mux.HandleFunc("PUT /api/groups/{id}/avatar", groupHandler.UpdateGroupAvatar)
	mux.HandleFunc("GET /api/groups/{id}/posts", postHandler.GetGroupPosts)
	mux.HandleFunc("POST /api/groups/{id}/posts", postHandler.CreateGroupPost)
	mux.HandleFunc("GET /api/groups/{id}/events", groupHandler.ListGroupEvents)
	mux.HandleFunc("POST /api/groups/{id}/events", groupHandler.CreateGroupEvent)
	mux.HandleFunc("PUT /api/groups/{id}/events/{eventId}", groupHandler.UpdateGroupEvent)
	mux.HandleFunc("DELETE /api/groups/{id}/events/{eventId}", groupHandler.DeleteGroupEvent)
	mux.HandleFunc("POST /api/groups/{id}/events/{eventId}/rsvp", groupHandler.SetGroupEventRSVP)

	// Upload routes
	mux.HandleFunc("POST /api/uploads", uploadHandler.UploadImage)

	// Notification routes
	mux.HandleFunc("GET /api/notifications", notificationHandler.ListNotifications)
	mux.HandleFunc("GET /api/notifications/unread-count", notificationHandler.GetUnreadCount)
	mux.HandleFunc("POST /api/notifications/{id}/read", notificationHandler.MarkRead)
	mux.HandleFunc("POST /api/notifications/read-all", notificationHandler.MarkAllRead)

	// Websocket route
	mux.HandleFunc("/ws", websocketHandler.HandleWebSocket)

	// Static uploads (served from disk)
	mux.Handle(handlers.UploadURLPrefix, http.StripPrefix(handlers.UploadURLPrefix, http.FileServer(http.Dir(config.Upload.UploadPath))))

	// Frontend
	mux.Handle("/", http.FileServer(http.Dir(config.Frontend.Path)))

	// Middleware chain (outermost first).
	var handler http.Handler = mux
	handler = middleware.RateLimiterMiddleware(config)(handler)
	handler = middleware.CORSMiddleware(config)(handler)
	handler = middleware.SecurityHeadersMiddleware()(handler)
	handler = middleware.LoggingMiddleware(logger)(handler)
	handler = middleware.RecoveryMiddleware(logger)(handler)

	return handler
}
