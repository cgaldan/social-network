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

type FollowHandler struct {
	followService service.FollowServiceInterface
	authService   service.AuthServiceInterface
	logger        *logger.Logger
}

func NewFollowHandler(followService service.FollowServiceInterface, authService service.AuthServiceInterface, logger *logger.Logger) *FollowHandler {
	return &FollowHandler{
		followService: followService,
		authService:   authService,
		logger:        logger,
	}
}

func (h *FollowHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		json.NewEncoder(w).Encode(domain.FollowResponse{
			Success: false,
			Message: "Missing authorization token",
		})
		return
	}

	user, err := h.authService.ValidateSession(tokenString)
	if err != nil {
		json.NewEncoder(w).Encode(domain.FollowResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	followeeID, err := strconv.Atoi(vars["id"])
	if err != nil || followeeID <= 0 {
		json.NewEncoder(w).Encode(domain.FollowResponse{
			Success: false,
			Message: "Invalid followee ID",
		})
		return
	}

	var followRequest domain.FollowRequest
	err = json.NewDecoder(r.Body).Decode(&followRequest)
	if err != nil {
		json.NewEncoder(w).Encode(domain.FollowResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if followRequest.FollowerID != user.ID && followRequest.FollowerID != 0 {
		json.NewEncoder(w).Encode(domain.FollowResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	followRequest.FollowerID = user.ID
	followRequest.FolloweeID = followeeID

	status, err := h.followService.FollowUser(followRequest)
	if err != nil {
		json.NewEncoder(w).Encode(domain.FollowResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.FollowResponse{
		Success: true,
		Message: "Follow request processed",
		Status:  status,
	})
}
