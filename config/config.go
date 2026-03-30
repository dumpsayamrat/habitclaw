package config

import "os"

type Config struct {
	Port     string
	DBPath   string
	DBType   string
	DBDSN    string
	Password string
	UserID   string
	LogLevel string
}

func Load() Config {
	return Config{
		Port:     getEnv("HABITCLAW_PORT", "3000"),
		DBPath:   getEnv("HABITCLAW_DB_PATH", "./habitclaw.db"),
		DBType:   getEnv("HABITCLAW_DB_TYPE", "sqlite"),
		DBDSN:    getEnv("HABITCLAW_DB_DSN", ""),
		Password: getEnv("HABITCLAW_PASSWORD", ""),
		UserID:   getEnv("HABITCLAW_USER_ID", "local"),
		LogLevel: getEnv("HABITCLAW_LOG_LEVEL", "info"),
	}
}

// DSN returns the appropriate connection string based on database type.
func (c Config) DSN() string {
	if c.DBType == "sqlite" {
		return c.DBPath
	}
	return c.DBDSN
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
