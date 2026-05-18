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

type UserHandler struct {
	authService   service.AuthServiceInterface
	postService   service.PostServiceInterface
	followService service.FollowServiceInterface
	logger        *logger.Logger
}

func NewUserHandler(authService service.AuthServiceInterface, postService service.PostServiceInterface, followService service.FollowServiceInterface, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		authService:   authService,
		postService:   postService,
		followService: followService,
		logger:        logger,
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	viewer, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.UserResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil || userID <= 0 {
		json.NewEncoder(w).Encode(domain.UserResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	profile, err := h.authService.GetUserProfile(viewer.ID, userID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.UserResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.UserResponse{
		Success: true,
		Message: "User retrieved successfully",
		User:    profile,
	})
}

func (h *UserHandler) ListUserPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	viewer, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostsResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["id"])
	if err != nil || targetID <= 0 {
		json.NewEncoder(w).Encode(domain.PostsResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	canView, err := h.authService.CanViewUser(viewer.ID, targetID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostsResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	if !canView {
		json.NewEncoder(w).Encode(domain.PostsResponse{
			Success: false,
			Message: "This profile is private",
		})
		return
	}

	limit, offset := parsePagination(r)
	posts, err := h.postService.GetPostsByUserID(targetID, viewer.ID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list user posts", "error", err, "targetID", targetID)
		json.NewEncoder(w).Encode(domain.PostsResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.PostsResponse{
		Success: true,
		Message: "Posts retrieved successfully",
		Posts:   posts,
	})
}
