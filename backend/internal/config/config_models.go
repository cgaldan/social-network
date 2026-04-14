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
