package config

import (
	"time"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Session     SessionConfig
	RateLimit   RateLimitConfig
	Websocket   WebSocketConfig
	Frontend    FrontendConfig
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8000"),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./data/database/social.db"),
		},
		Session: SessionConfig{
			Duration: getEnvDuration("SESSION_DURATION", 24*time.Hour),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getIntEnv("RATE_LIMIT_RPM", 100),
			Enabled:           getBoolEnv("RATE_LIMIT_ENABLED", true),
		},
		Websocket: WebSocketConfig{
			ReadBufferSize:  getIntEnv("WS_READ_BUFFER_SIZE", 1024),
			WriteBufferSize: getIntEnv("WS_WRITE_BUFFER_SIZE", 1024),
			PingPeriod:      getEnvDuration("WS_PING_PERIOD", 54*time.Second),
			PongWait:        getEnvDuration("WS_PONG_WAIT", 60*time.Second),
			WriteWait:       getEnvDuration("WS_WRITE_WAIT", 10*time.Second),
		},
		Frontend: FrontendConfig{
			Path: getEnv("FRONTEND_PATH", "./frontend"),
		},
	}

	return cfg, nil
}
