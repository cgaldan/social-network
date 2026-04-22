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
