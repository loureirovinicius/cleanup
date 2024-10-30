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
	logLevel = &slog.LevelVar{}
)

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}

func SetLevel(level string) {
	logLevel.Set(logLevels[level])
}

func Log(ctx context.Context, level string, msg string, args ...any) {
	level = strings.ToLower(level)

	logger.Log(ctx, logLevels[level], msg, args...)
}
