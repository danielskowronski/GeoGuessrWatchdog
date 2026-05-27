package config

import (
	"log/slog"
	"strings"
)

type DatabaseConfig struct {
	URL string `koanf:"url"`
}

type LoggingConfig struct {
	Level string `koanf:"level"`
}

func (c LoggingConfig) GetLevel() slog.Level {
	switch strings.ToLower(c.Level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
