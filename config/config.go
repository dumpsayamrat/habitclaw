package config

import "os"

type Config struct {
	Port     string
	DBPath   string
	Password string
	UserID   string
	LogLevel string
}

func Load() Config {
	return Config{
		Port:     getEnv("HABITCLAW_PORT", "3000"),
		DBPath:   getEnv("HABITCLAW_DB_PATH", "./habitclaw.db"),
		Password: getEnv("HABITCLAW_PASSWORD", ""),
		UserID:   getEnv("HABITCLAW_USER_ID", "local"),
		LogLevel: getEnv("HABITCLAW_LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
