package logger

import (
	"context"
	"io"
	"log/slog"
	"strings"
)

var (
	logger    *slog.Logger
	logLevels = map[string]slog.Level{
		"info":  slog.LevelInfo,
		"error": slog.LevelError,
		"debug": slog.LevelDebug,
	}
	logLevel = &slog.LevelVar{}
)

// Initialize logger
func InitializeLogger(level string, format string, dst io.Writer) {

	logFormat := map[string]slog.Handler{
		"text": slog.NewTextHandler(dst, &slog.HandlerOptions{Level: logLevel}),
		"json": slog.NewJSONHandler(dst, &slog.HandlerOptions{Level: logLevel}),
	}

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
