package logging

import (
	"log/slog"
	"os"
	"strings"
)

// Config del logger
func Setup(levelName string) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parseLevel(levelName)})
	slog.SetDefault(slog.New(handler))
}

func parseLevel(levelName string) slog.Level {
	switch strings.ToLower(levelName) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
