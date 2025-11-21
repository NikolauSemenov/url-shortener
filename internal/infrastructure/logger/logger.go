package logger

import (
	"log/slog"
	"os"
	"url-shortener/internal/ports"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Logger struct {
	l *slog.Logger
}

func NewLogger(env string) *Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return &Logger{l: log}
}

func (s *Logger) Info(msg string, args ...any) {
	s.l.Info(msg, args...)
}

func (s *Logger) Error(msg string, args ...any) {
	s.l.Error(msg, args...)
}

func (s *Logger) With(args ...any) ports.Logger {
	return &Logger{
		l: s.l.With(args...),
	}
}
