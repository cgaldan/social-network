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

type CommentHandler struct {
	commentService service.CommentServiceInterface
	authService    service.AuthServiceInterface
	logger         *logger.Logger
}

func NewCommentHandler(commentService service.CommentServiceInterface, authService service.AuthServiceInterface, logger *logger.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		authService:    authService,
		logger:         logger,
	}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(tokenString)
	if err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// WITH GORILLA PKG IMPLEMANTATION
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])

	if err != nil || postID <= 0 {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var commentRequest domain.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&commentRequest); err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	comment, err := h.commentService.CreateComment(user.ID, postID, commentRequest)
	if err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.CommentResponse{
		Success: true,
		Message: "Comment created successfully",
		Comment: comment,
	})
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(tokenString)
	if err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil || postID <= 0 {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	commentID, err := strconv.Atoi(vars["commentId"])
	if err != nil || commentID <= 0 {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Invalid comment ID",
		})
		return
	}

	var body domain.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	comment, err := h.commentService.UpdateComment(user.ID, postID, commentID, body)
	if err != nil {
		h.logger.Error("Failed to update comment", "error", err, "userID", user.ID, "postID", postID, "commentID", commentID)
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.CommentResponse{
		Success: true,
		Message: "Comment updated successfully",
		Comment: comment,
	})
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(tokenString)
	if err != nil {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil || postID <= 0 {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	commentID, err := strconv.Atoi(vars["commentId"])
	if err != nil || commentID <= 0 {
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: "Invalid comment ID",
		})
		return
	}

	if err := h.commentService.DeleteComment(user.ID, postID, commentID); err != nil {
		h.logger.Error("Failed to delete comment", "error", err, "userID", user.ID, "postID", postID, "commentID", commentID)
		json.NewEncoder(w).Encode(domain.CommentResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.CommentResponse{
		Success: true,
		Message: "Comment deleted successfully",
	})
}
