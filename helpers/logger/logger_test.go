package logger

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

// Helper function to capture stdout
func captureOutput(f func(dst io.Writer)) string {
	var buf bytes.Buffer
	f(&buf)
	return buf.String()
}

func TestInitializeLogger(t *testing.T) {
	tests := []struct {
		level     string
		format    string
		expectOut string
	}{
		{"info", "text", "INFO"},
		{"debug", "json", "DEBUG"},
		{"error", "text", "ERROR"},
		{"invalid", "invalid", "INFO"}, // fallback to defaults
	}

	for _, input := range tests {
		t.Run(input.level+"_"+input.format, func(t *testing.T) {
			output := captureOutput(func(dst io.Writer) {
				InitializeLogger(input.level, input.format, dst)
				Log(context.Background(), input.level, "test message")
			})

			if !strings.Contains(output, input.expectOut) {
				t.Errorf("expected log level %s, got %s", input.expectOut, output)
			}
		})
	}
}

func TestLog(t *testing.T) {

	tests := []struct {
		level     string
		message   string
		args      []any
		expectOut string
	}{
		{"info", "info message", nil, "INFO"},
		{"debug", "debug message", nil, ""}, // won't log at info level
		{"error", "error message", nil, "ERROR"},
		{"warn", "warn message", nil, "INFO"}, // logs as info (fallback)
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			output := captureOutput(func(dst io.Writer) {
				InitializeLogger("info", "text", dst)
				Log(context.Background(), tt.level, tt.message, tt.args...)
			})

			if tt.expectOut != "" && !strings.Contains(output, tt.expectOut) {
				t.Errorf("expected output to contain %s, got %s", tt.expectOut, output)
			} else if tt.expectOut == "" && output != "" {
				t.Errorf("expected no output, got %s", output)
			}
		})
	}
}
