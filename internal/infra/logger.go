package infra

import (
	"log/slog"
	"os"
	"strings"
)

type LogConfig struct {
	Level  slog.Level
	Format string
}

func (r *Resources) initLogger() {
	handler := getHandler(r.Env.LogFormat, &slog.HandlerOptions{
		Level: parseLevel(r.Env.LogLevel),
	})

	r.Logger = slog.New(handler)
	slog.SetDefault(r.Logger)
}

func parseLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}
	return slog.LevelInfo
}

func getHandler(f string, opts *slog.HandlerOptions) slog.Handler {
	switch f {
	case "json":
		return slog.NewJSONHandler(os.Stdout, opts)
	default:
		return slog.NewTextHandler(os.Stdout, opts)
	}
}
