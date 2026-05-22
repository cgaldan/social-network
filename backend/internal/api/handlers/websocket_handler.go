package handlers

import (
	"net/http"
	"social-network/internal/api/middleware"
	"social-network/internal/config"
	"social-network/internal/service"
	"social-network/internal/websocket"
	"social-network/packages/logger"

	gorillaws "github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub         *websocket.Hub
	authService service.AuthServiceInterface
	logger      *logger.Logger
	upgrader    gorillaws.Upgrader
	wsConfig    config.WebSocketConfig
}

func NewWebSocketHandler(hub *websocket.Hub, authService service.AuthServiceInterface, logger *logger.Logger, wsConfig config.WebSocketConfig, allowedOrigins []string) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
		logger:      logger,
		upgrader: gorillaws.Upgrader{
			ReadBufferSize:  wsConfig.ReadBufferSize,
			WriteBufferSize: wsConfig.WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return middleware.IsAllowedOrigin(r.Header.Get("Origin"), allowedOrigins)
			},
		},
		wsConfig: wsConfig,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing authorization token", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.ValidateSession(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := websocket.NewClient(h.hub, conn, user.ID, h.logger, h.wsConfig)

	h.hub.RegisterClientToHub(client)

	go client.ReadPump()
	go client.WritePump()
}
