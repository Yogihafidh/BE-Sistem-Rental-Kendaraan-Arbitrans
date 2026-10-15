package logger

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger initializes and returns an idiomatic Go standard library structured logger (slog).
// In production, it emits JSON logs with Info level. In development, it emits readable text logs with Debug level.
func NewLogger(env string) *slog.Logger {
	var handler slog.Handler

	isProd := strings.EqualFold(env, "production")

	opts := &slog.HandlerOptions{
		Level: func() slog.Level {
			if isProd {
				return slog.LevelInfo
			}
			return slog.LevelDebug
		}(),
		AddSource: !isProd,
	}

	if isProd {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}
