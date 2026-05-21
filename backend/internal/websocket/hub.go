package websocket

import (
	"social-network/internal/repository"
	"social-network/packages/logger"
	"sync"
)

type Hub struct {
	clients    map[int]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	logger     *logger.Logger
	userRepo   repository.UserRepositoryInterface
}

func NewHub(logger *logger.Logger, userRepo repository.UserRepositoryInterface) *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mu:         sync.RWMutex{},
		logger:     logger,
		userRepo:   userRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClientCase(client)
		case client := <-h.unregister:
			h.unregisterClientCase(client)
		}
	}
}

func (h *Hub) registerClientCase(client *Client) {
	h.mu.Lock()
	h.clients[client.UserID] = client
	h.mu.Unlock()

	h.userRepo.UpdateLastSeen(client.UserID)
	if err := h.userRepo.SetOnline(client.UserID, true); err != nil {
		h.logger.Error("Failed to mark user online", "userID", client.UserID, "error", err)
	}

	h.broadcastUserStatus(client.UserID, true)

	h.sendOnlineUsers(client)

	h.logger.Info("Websocket client connected", "userID", client.UserID, "totalClients", len(h.clients))
}

func (h *Hub) unregisterClientCase(client *Client) {
	h.mu.Lock()
	removed := false
	if _, ok := h.clients[client.UserID]; ok {
		delete(h.clients, client.UserID)
		close(client.Send)
		removed = true
	}
	h.mu.Unlock()

	h.userRepo.UpdateLastSeen(client.UserID)
	if removed {
		if err := h.userRepo.SetOnline(client.UserID, false); err != nil {
			h.logger.Error("Failed to mark user offline", "userID", client.UserID, "error", err)
		}
	}

	h.broadcastUserStatus(client.UserID, false)

	h.logger.Info("Websocket client disconnected", "userID", client.UserID, "totalClients", len(h.clients))
}
