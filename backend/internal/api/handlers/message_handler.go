package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/internal/domain"
	"social-network/internal/service"
	"social-network/packages/logger"
	"strconv"
	// WITH GORILLA PKG IMPLEMENTATION
	// "github.com/gorilla/mux"
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

func (h *MessageHandler) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// WITHOUT GORILLA PKG IMPLEMENTATION
	messageIDStr := r.PathValue("id")
	messageID, err := strconv.Atoi(messageIDStr)

	// WITH GORILLA PKG IMPLEMENTATION
	// vars := mux.Vars(r)
	// messageID, err := strconv.Atoi(vars["id"])

	if err != nil || messageID <= 0 {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Invalid message ID",
		})
		return
	}

	var req domain.UpdateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Invalid JSON",
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

	updated, err := h.messageService.UpdateMessage(messageID, user.ID, req.Content)
	if err != nil {
		h.logger.Error("Failed to update message", "error", err, "messageID", messageID, "userID", user.ID)
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.MessageResponse{
		Success: true,
		Message: "Message updated",
		Msg:     updated,
	})
}

func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// WITHOUT GORILLA PKG IMPLEMENTATION
	messageIDStr := r.PathValue("id")
	messageID, err := strconv.Atoi(messageIDStr)

	// WITH GORILLA PKG IMPLEMENTATION
	// vars := mux.Vars(r)
	// messageID, err := strconv.Atoi(vars["id"])

	if err != nil || messageID <= 0 {
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: "Invalid message ID",
		})
		return
	}

	if err := h.messageService.DeleteMessage(messageID, user.ID); err != nil {
		h.logger.Error("Failed to delete message", "error", err, "messageID", messageID, "userID", user.ID)
		json.NewEncoder(w).Encode(domain.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.MessageResponse{
		Success: true,
		Message: "Message deleted",
	})
}

func (h *MessageHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.MessagesResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// WITHOUT GORILLA PKG IMPLEMENTATION
	convIDStr := r.PathValue("id")
	convID, err := strconv.Atoi(convIDStr)

	// WITH GORILLA PKG IMPLEMENTATION
	// vars := mux.Vars(r)
	// convID, err := strconv.Atoi(vars["id"])

	if err != nil || convID <= 0 {
		json.NewEncoder(w).Encode(domain.MessagesResponse{
			Success: false,
			Message: "Invalid conversation ID",
		})
		return
	}

	limit, _ := parsePagination(r)
	beforeID := 0
	if v := r.URL.Query().Get("before_id"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			beforeID = n
		}
	}

	messages, err := h.messageService.ListMessages(convID, user.ID, limit+1, beforeID)
	if err != nil {
		h.logger.Error("Failed to list messages", "error", err, "convID", convID, "userID", user.ID)
		json.NewEncoder(w).Encode(domain.MessagesResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	items, hasMore := trimPage(messages, limit)
	json.NewEncoder(w).Encode(domain.MessagesResponse{
		Success:  true,
		Message:  "Messages retrieved successfully",
		Messages: items,
		HasMore:  hasMore,
	})
}
