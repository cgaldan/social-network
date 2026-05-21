package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/internal/domain"
	"social-network/internal/service"
	"social-network/packages/logger"
	"strconv"

	"github.com/gorilla/mux"
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

func (h *ConversationHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ConversationsResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	limit, offset := parsePagination(r)

	conversations, err := h.convService.ListConversations(user.ID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list conversations", "error", err, "userID", user.ID)
		json.NewEncoder(w).Encode(domain.ConversationsResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.ConversationsResponse{
		Success:       true,
		Message:       "Conversations retrieved successfully",
		Conversations: conversations,
	})
}

func (h *ConversationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.ConversationResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	convID, err := strconv.Atoi(vars["id"])
	if err != nil || convID <= 0 {
		json.NewEncoder(w).Encode(domain.ConversationResponse{
			Success: false,
			Message: "Invalid conversation ID",
		})
		return
	}

	if err := h.convService.MarkRead(user.ID, convID); err != nil {
		h.logger.Error("Failed to mark conversation read", "error", err, "userID", user.ID, "convID", convID)
		json.NewEncoder(w).Encode(domain.ConversationResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.ConversationResponse{
		Success: true,
		Message: "Conversation marked as read",
	})
}
