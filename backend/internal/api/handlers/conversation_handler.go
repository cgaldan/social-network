package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/internal/domain"
	"social-network/internal/service"
	"social-network/packages/logger"
)

type ConversationHandler struct {
	convService service.ConversationServiceInterface
	authService service.AuthServiceInterface
	logger      *logger.Logger
}

func NewConversationHandler(convService service.ConversationServiceInterface, authService service.AuthServiceInterface, logger *logger.Logger) *ConversationHandler {
	return &ConversationHandler{
		convService: convService,
		authService: authService,
		logger:      logger,
	}
}

func (h *ConversationHandler) CreateDirectConversation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	sender, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ConversationResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req domain.DirectConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.ConversationResponse{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	req.SenderID = sender.ID

	if req.ReceiverID == 0 {
		json.NewEncoder(w).Encode(domain.ConversationResponse{
			Success: false,
			Message: "Receiver ID is required",
		})
		return
	}

	conv, err := h.convService.CreateDirectConversation(req)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ConversationResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.ConversationResponse{
		Success:      true,
		Message:      "Conversation created successfully",
		Conversation: conv,
	})
}
