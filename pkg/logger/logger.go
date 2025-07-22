package logger

import (
	"fmt"
	"io"
	"log/slog"

	"os"
	"strings"
)

const (
	production  = "prod"
	development = "dev"
	mock        = "mock"
)

type Logger interface {
	Debug(msg string)
	DebugF(format string, args ...any)
	DebugW(msg string, args map[string]any)

	Info(msg string)
	InfoF(format string, args ...any)
	InfoW(msg string, args map[string]any)

	Warn(msg string)
	WarnF(format string, args ...any)
	WarnW(msg string, args map[string]any)

	Error(msg string)
	ErrorF(format string, args ...any)
	ErrorW(msg string, args map[string]any)

	WithAttrs(attrs map[string]any) Logger
}

type logger struct {
	log *slog.Logger
}

func New(mode string) (Logger, error) {
	switch strings.ToLower(mode) {
	case production:
		return &logger{
			log: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		}, nil
	case development:
		return &logger{
			log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		}, nil
	case mock:
		return &logger{
			log: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug})),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported logger mode: %s", mode)
	}
}

func (l *logger) Debug(msg string) {
	l.log.Debug(msg)
}

func (l *logger) DebugF(format string, args ...any) {
	l.log.Debug(fmt.Sprintf(format, args...))
}

func (l *logger) DebugW(msg string, args map[string]any) {
	l.log.Debug(msg, convertMapToSlogAttrs(args)...)
}

func (l *logger) Info(msg string) {
	l.log.Info(msg)
}

func (l *logger) InfoF(format string, args ...any) {
	l.log.Info(fmt.Sprintf(format, args...))
}

func (l *logger) InfoW(msg string, args map[string]any) {
	l.log.Info(msg, convertMapToSlogAttrs(args)...)
}

func (l *logger) Warn(msg string) {
	l.log.Warn(msg)
}

func (l *logger) WarnF(format string, args ...any) {
	l.log.Warn(fmt.Sprintf(format, args...))
}

func (l *logger) WarnW(msg string, args map[string]any) {
	l.log.Warn(msg, convertMapToSlogAttrs(args)...)
}

func (l *logger) Error(msg string) {
	l.log.Error(msg)
}

func (l *logger) ErrorF(format string, args ...any) {
	l.log.Error(fmt.Sprintf(format, args...))
}

func (l *logger) ErrorW(msg string, args map[string]any) {
	l.log.Error(msg, convertMapToSlogAttrs(args)...)
}

func (l *logger) WithAttrs(attrs map[string]any) Logger {
	return &logger{
		log: l.log.With(convertMapToSlogAttrs(attrs)...),
	}
}

func convertMapToSlogAttrs(args map[string]any) []any {
	attrs := make([]any, 0, len(args)*2)

	for k, v := range args {
		attrs = append(attrs, slog.Any(k, v))
	}

	return attrs
}
