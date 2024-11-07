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
	logLevel  = &slog.LevelVar{}
	logFormat = map[string]slog.Handler{
		"text": slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
		"json": slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}),
	}
)

// Initialize logger
func InitializeLogger(level string, format string) {
	// Set log level
	lvl, ok := logLevels[level]
	if !ok {
		lvl = logLevels["info"] // default/fallback level
	}
	logLevel.Set(lvl)

	// Set log format
	fmt, ok := logFormat[format]
	if !ok {
		fmt = logFormat["text"] // default/fallback format
	}
	logger = slog.New(fmt)
}

// Log whatever is being passed in parameters
func Log(ctx context.Context, level string, msg string, args ...any) {
	level = strings.ToLower(level)

	lvl, ok := logLevels[level]
	if !ok {
		lvl = slog.LevelInfo
	}
	logger.Log(ctx, lvl, msg, args...)
}
