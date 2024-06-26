package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

var (
	logLevels = map[string]slog.Level{
		"info":  slog.LevelInfo,
		"error": slog.LevelError,
		"debug": slog.LevelDebug,
	}
)

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
}

func Log(ctx context.Context, level string, msg string, args ...any) {
	level = strings.ToLower(level)

	logger.Log(ctx, logLevels[level], msg, args...)
}
