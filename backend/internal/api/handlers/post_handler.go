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

type PostHandler struct {
	postService    service.PostServiceInterface
	authService    service.AuthServiceInterface
	contentService service.ContentServiceInterface
	logger         *logger.Logger
}

func NewPostHandler(postService service.PostServiceInterface, authService service.AuthServiceInterface, contentService service.ContentServiceInterface, logger *logger.Logger) *PostHandler {
	return &PostHandler{
		postService:    postService,
		authService:    authService,
		contentService: contentService,
		logger:         logger,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var postData domain.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	post, err := h.contentService.CreatePost(user.ID, postData)
	if err != nil {
		h.logger.Error("Failed to create post", "error", err, "userID", user.ID)
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.PostDetailResponse{
		Success: true,
		Message: "Post created successfully",
		Post:    post,
	})
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	category := r.URL.Query().Get("category")
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

	posts, err := h.postService.ListPosts(category, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list posts", "error", err)
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

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// WITH GORILLA PKG IMPLEMENTATION
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])

	if err != nil || postID <= 0 {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	post, err := h.postService.GetPostByID(user.ID, postID)
	if err != nil {
		h.logger.Error("Failed to get post by ID", "error", err, "postID", postID)
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.PostDetailResponse{
		Success: true,
		Message: "Post retrieved successfully",
		Post:    post,
	})
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil || postID <= 0 {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	var body domain.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	post, err := h.contentService.UpdatePost(user.ID, postID, body)
	if err != nil {
		h.logger.Error("Failed to update post", "error", err, "userID", user.ID, "postID", postID)
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.PostDetailResponse{
		Success: true,
		Message: "Post updated successfully",
		Post:    post,
	})
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil || postID <= 0 {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	if err := h.contentService.DeletePost(user.ID, postID); err != nil {
		h.logger.Error("Failed to delete post", "error", err, "userID", user.ID, "postID", postID)
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.PostDetailResponse{
		Success: true,
		Message: "Post deleted successfully",
	})
}

func (h *PostHandler) CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["id"])

	if err != nil || groupID <= 0 {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Invalid group ID",
		})
		return
	}

	var postData domain.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	post, err := h.contentService.CreateGroupPost(user.ID, groupID, postData)
	if err != nil {
		h.logger.Error("Failed to create group post", "error", err, "userID", user.ID, "groupID", groupID)
		json.NewEncoder(w).Encode(domain.PostDetailResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.PostDetailResponse{
		Success: true,
		Message: "Post created successfully",
		Post:    post,
	})
}

func (h *PostHandler) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.PostsResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	groupID, err := strconv.Atoi(vars["id"])
	if err != nil || groupID <= 0 {
		json.NewEncoder(w).Encode(domain.PostsResponse{
			Success: false,
			Message: "Invalid group ID",
		})
		return
	}

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

	posts, err := h.postService.ListPostsByGroupID(user.ID, groupID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list group posts", "error", err, "userID", user.ID, "groupID", groupID)
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
