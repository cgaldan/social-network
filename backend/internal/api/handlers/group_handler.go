package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/internal/domain"
	"social-network/internal/service"
	"social-network/packages/logger"
)

type GroupHandler struct {
	groupService service.GroupServiceInterface
	authService  service.AuthServiceInterface
	logger       *logger.Logger
}

func NewGroupHandler(groupService service.GroupServiceInterface, authService service.AuthServiceInterface, logger *logger.Logger) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
		authService:  authService,
		logger:       logger,
	}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	creator, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req domain.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.GroupResponse{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	req.CreatorID = creator.ID

	group, err := h.groupService.CreateGroup(&domain.Group{
		CreatorID:   req.CreatorID,
		Title:       req.Title,
		Description: req.Description,
	})

	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.GroupResponse{
		Success: true,
		Message: "Group created successfully",
		Group:   group,
	})
}
