package websocket

import (
	"encoding/json"
	"social-network/internal/config"
	"social-network/packages/logger"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub        *Hub
	Connection *websocket.Conn
	Send       chan []byte
	UserID     int
	Logger     *logger.Logger
	WsConfig   config.WebSocketConfig
}

type WsMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func NewClient(hub *Hub, connection *websocket.Conn, userID int, logger *logger.Logger, wsConfig config.WebSocketConfig) *Client {
	return &Client{
		Hub:        hub,
		Connection: connection,
		Send:       make(chan []byte, 256),
		UserID:     userID,
		Logger:     logger,
		WsConfig:   wsConfig,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Connection.Close()
	}()
	c.Connection.SetReadDeadline(time.Now().Add(c.WsConfig.PongWait))
	c.Connection.SetPongHandler(func(string) error {
		c.Connection.SetReadDeadline(time.Now().Add(c.WsConfig.PongWait))
		return nil
	})

	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Logger.Error("WebSocket read error", "error", err, "userID", c.UserID)
			}
			break
		}

		var wsMessage WsMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			c.Logger.Warn("Failed to unmarshal message", "error", err)
			continue
		}

		switch wsMessage.Type {
		case "ping":
			pong := WsMessage{
				Type: "pong",
			}
			data, _ := json.Marshal(pong)
			c.Send <- data
		default:
			c.Logger.Debug("Unknown message type", "type", wsMessage.Type)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(c.WsConfig.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Connection.SetWriteDeadline(time.Now().Add(c.WsConfig.WriteWait))
			if !ok {
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(message)

			if err := writer.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(c.WsConfig.WriteWait))
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
