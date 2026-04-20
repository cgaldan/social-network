package config

import "time"

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Path string
}

type SessionConfig struct {
	Duration time.Duration
}
type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	PingPeriod      time.Duration
	PongWait        time.Duration
	WriteWait       time.Duration
}

type FrontendConfig struct {
	Path string
}
