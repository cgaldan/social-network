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

func (h *GroupHandler) InviteToGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	inviter, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req domain.InviteToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.GroupResponse{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	req.InviterID = inviter.ID

	err = h.groupService.CreateGroupInvitation(req.GroupID, req.InviterID, req.InviteeID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.GroupResponse{
		Success: true,
		Message: "Invitation created successfully",
	})
}

func (h *GroupHandler) AcceptGroupInvitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	invitationID, err := strconv.Atoi(vars["id"])

	if err != nil || invitationID <= 0 {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: "Invalid invitation ID",
		})
		return
	}

	invitation, err := h.groupService.GetGroupInvitationByID(invitationID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: "Invitation not found",
		})
		return
	}

	err = h.groupService.AcceptGroupInvitation(user.ID, invitation)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	invitation, err = h.groupService.GetGroupInvitationByID(invitationID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: "Invitation not found",
		})
		return
	}

	json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
		Success:    true,
		Message:    "Invitation accepted successfully",
		Invitation: invitation,
	})
}

func (h *GroupHandler) DeclineGroupInvitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	invitationID, err := strconv.Atoi(vars["id"])

	if err != nil || invitationID <= 0 {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: "Invalid invitation ID",
		})
		return
	}

	invitation, err := h.groupService.GetGroupInvitationByID(invitationID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: "Invitation not found",
		})
		return
	}

	err = h.groupService.DeclineGroupInvitation(user.ID, invitation)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	invitation, err = h.groupService.GetGroupInvitationByID(invitationID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
			Success: false,
			Message: "Invitation not found",
		})
		return
	}

	json.NewEncoder(w).Encode(domain.GroupInvitationResponse{
		Success:    true,
		Message:    "Invitation declined successfully",
		Invitation: invitation,
	})
}

func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var req domain.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(domain.GroupResponse{
			Success: false,
			Message: "Invalid JSON",
		})
		return
	}

	err = h.groupService.CreateGroupJoinRequest(req.GroupID, user.ID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
		Success: true,
		Message: "Joine group request created successfully",
	})
}

func (h *GroupHandler) AcceptGroupJoinRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	answerer, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	requestID, err := strconv.Atoi(vars["id"])
	if err != nil || requestID <= 0 {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: "Invalid request ID",
		})
		return
	}

	request, err := h.groupService.GetGroupJoinRequestByID(requestID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: "Request not found",
		})
		return
	}

	err = h.groupService.AcceptGroupJoinRequest(answerer.ID, request)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	request, err = h.groupService.GetGroupJoinRequestByID(requestID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: "Request not found",
		})
		return
	}

	json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
		Success: true,
		Message: "Request accepted successfully",
		Request: request,
	})
}

func (h *GroupHandler) DeclineGroupJoinRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	answerer, err := h.authService.ValidateSession(token)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	vars := mux.Vars(r)
	requestID, err := strconv.Atoi(vars["id"])
	if err != nil || requestID <= 0 {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: "Invalid request ID",
		})
		return
	}

	request, err := h.groupService.GetGroupJoinRequestByID(requestID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: "Request not found",
		})
		return
	}

	err = h.groupService.DeclineGroupJoinRequest(answerer.ID, request)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	request, err = h.groupService.GetGroupJoinRequestByID(requestID)
	if err != nil {
		json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
			Success: false,
			Message: "Request not found",
		})
		return
	}

	json.NewEncoder(w).Encode(domain.GroupJoinRequestResponse{
		Success: true,
		Message: "Request declined successfully",
		Request: request,
	})
}
