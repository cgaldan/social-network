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

type NotificationHandler struct {
	notificationService service.NotificationServiceInterface
	authService         service.AuthServiceInterface
	logger              *logger.Logger
}

func NewNotificationHandler(notificationService service.NotificationServiceInterface, authService service.AuthServiceInterface, logger *logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		authService:         authService,
		logger:              logger,
	}
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.NotificationsResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	limit, offset := parsePagination(r)

	notifications, err := h.notificationService.ListNotifications(user.ID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list notifications", "error", err, "userID", user.ID)
		json.NewEncoder(w).Encode(domain.NotificationsResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.NotificationsResponse{
		Success:       true,
		Message:       "Notifications retrieved successfully",
		Notifications: notifications,
	})
}

func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.NotificationUnreadCountResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	count, err := h.notificationService.CountUnread(user.ID)
	if err != nil {
		h.logger.Error("Failed to count unread notifications", "error", err, "userID", user.ID)
		json.NewEncoder(w).Encode(domain.NotificationUnreadCountResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.NotificationUnreadCountResponse{
		Success:     true,
		Message:     "Unread notification count retrieved successfully",
		UnreadCount: count,
	})
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.NotificationResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	notificationID, err := strconv.Atoi(vars["id"])
	if err != nil || notificationID <= 0 {
		json.NewEncoder(w).Encode(domain.NotificationResponse{
			Success: false,
			Message: "Invalid notification ID",
		})
		return
	}

	if err := h.notificationService.MarkRead(user.ID, notificationID); err != nil {
		h.logger.Error("Failed to mark notification read", "error", err, "userID", user.ID, "notificationID", notificationID)
		json.NewEncoder(w).Encode(domain.NotificationResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.NotificationResponse{
		Success: true,
		Message: "Notification marked read successfully",
	})
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.NotificationResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	if err := h.notificationService.MarkAllRead(user.ID); err != nil {
		h.logger.Error("Failed to mark all notifications read", "error", err, "userID", user.ID)
		json.NewEncoder(w).Encode(domain.NotificationResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.NotificationResponse{
		Success: true,
		Message: "Notifications marked read successfully",
	})
}

func parsePagination(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	if limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}

	offset := 0
	if offsetStr != "" {
		if offsetNum, err := strconv.Atoi(offsetStr); err == nil && offsetNum >= 0 {
			offset = offsetNum
		}
	}

	return limit, offset
}
