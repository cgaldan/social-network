package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/internal/domain"
	"social-network/internal/service"
	"social-network/packages/logger"
)

type MessageHandler struct {
	messageService service.MessageServiceInterface
	authService    service.AuthServiceInterface
	logger         *logger.Logger
}

func NewMessageHandler(messageService service.MessageServiceInterface, authService service.AuthServiceInterface, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		authService:    authService,
		logger:         logger,
	}
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	sender, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req domain.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	if req.ConversationID == 0 {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Conversation ID is required",
		})
		return
	}

	if req.Content == "" {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Content is required",
		})
		return
	}

	message, err := h.messageService.SendMessage(req.ConversationID, sender.ID, req.Content)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.MessageResponse{
		Success: true,
		Message: "Message sent",
		Msg:     message,
	})
}
