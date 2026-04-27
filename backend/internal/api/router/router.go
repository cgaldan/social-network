package router

import (
	"net/http"
	"social-network/internal/api/handlers"
	"social-network/internal/api/middleware"
	"social-network/internal/config"
	"social-network/internal/service"
	"social-network/internal/websocket"
	"social-network/packages/logger"

	"github.com/gorilla/mux"
)

func NewRouter(services *service.Services, config *config.Config, hub *websocket.Hub, logger *logger.Logger) *mux.Router {
	r := mux.NewRouter()

	authHandler := handlers.NewAuthHandler(services.Auth, logger)
	postHandler := handlers.NewPostHandler(services.Post, services.Auth, services.Content, logger)
	commentHandler := handlers.NewCommentHandler(services.Comment, services.Auth, logger)
	websocketHandler := handlers.NewWebSocketHandler(hub, services.Auth, logger)
	followHandler := handlers.NewFollowHandler(services.Follow, services.Auth, logger)
	conversationHandler := handlers.NewConversationHandler(services.Conversation, services.Auth, logger)
	messageHandler := handlers.NewMessageHandler(services.Message, services.Auth, logger)
	groupHandler := handlers.NewGroupHandler(services.Group, services.Auth, logger)
	notificationHandler := handlers.NewNotificationHandler(services.Notification, services.Auth, logger)
	healthHandler := handlers.NewHealthHandler("1.0.0")

	r.HandleFunc("/health", healthHandler.Health).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	api.HandleFunc("/auth/me", authHandler.GetCurrentUser).Methods("GET")

	// Post routes
	api.HandleFunc("/posts", postHandler.GetPosts).Methods("GET")
	api.HandleFunc("/posts", postHandler.CreatePost).Methods("POST")
	api.HandleFunc("/posts/{id}", postHandler.GetPostByID).Methods("GET")
	api.HandleFunc("/posts/{id}/comments", commentHandler.CreateComment).Methods("POST")

	// Follow routes
	api.HandleFunc("/follow/{id}", followHandler.FollowUser).Methods("POST")

	// Chat routes
	api.HandleFunc("/conversations/direct", conversationHandler.CreateDirectConversation).Methods("POST")
	api.HandleFunc("/messages", messageHandler.SendMessage).Methods("POST")

	// Group routes
	api.HandleFunc("/groups", groupHandler.ListGroups).Methods("GET")
	api.HandleFunc("/groups", groupHandler.CreateGroup).Methods("POST")
	api.HandleFunc("/groups/join", groupHandler.JoinGroup).Methods("POST")
	api.HandleFunc("/groups/join/{id}/accept", groupHandler.AcceptGroupJoinRequest).Methods("POST")
	api.HandleFunc("/groups/join/{id}/decline", groupHandler.DeclineGroupJoinRequest).Methods("POST")
	api.HandleFunc("/groups/invitations", groupHandler.InviteToGroup).Methods("POST")
	api.HandleFunc("/groups/invitations/{id}/accept", groupHandler.AcceptGroupInvitation).Methods("POST")
	api.HandleFunc("/groups/invitations/{id}/decline", groupHandler.DeclineGroupInvitation).Methods("POST")
	api.HandleFunc("/groups/{id}/posts", postHandler.GetGroupPosts).Methods("GET")
	api.HandleFunc("/groups/{id}/posts", postHandler.CreateGroupPost).Methods("POST")
	api.HandleFunc("/groups/{id}/events", groupHandler.ListGroupEvents).Methods("GET")
	api.HandleFunc("/groups/{id}/events", groupHandler.CreateGroupEvent).Methods("POST")
	api.HandleFunc("/groups/{id}/events/{eventId}/rsvp", groupHandler.SetGroupEventRSVP).Methods("POST")

	// Notification routes
	api.HandleFunc("/notifications", notificationHandler.ListNotifications).Methods("GET")
	api.HandleFunc("/notifications/unread-count", notificationHandler.GetUnreadCount).Methods("GET")
	api.HandleFunc("/notifications/{id}/read", notificationHandler.MarkRead).Methods("POST")
	api.HandleFunc("/notifications/read-all", notificationHandler.MarkAllRead).Methods("POST")

	// Websocket routes
	r.HandleFunc("/ws", websocketHandler.HandleWebSocket)

	frontendPath := "../frontend"
	if config.Environment == "production" {
		frontendPath = config.Frontend.Path
	}
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(frontendPath))))

	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware(config))
	r.Use(middleware.RateLimiterMiddleware(config))

	return r
}
