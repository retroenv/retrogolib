// Package log provides logging functionality.
package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Logger provides fast, leveled, structured logging. All methods are safe
// for concurrent use.
type Logger struct {
	logger     *slog.Logger
	handler    slog.Handler
	callerInfo bool
	level      *slog.LevelVar
}

// New returns a new Logger instance.
func New() *Logger {
	cfg := Config{
		Level: DefaultLevel(),
	}

	return NewWithConfig(cfg)
}

// NewWithConfig creates a new logger for the given config.
// If no level is set in the config, it will use the default level of
// this package.
func NewWithConfig(cfg Config) *Logger {
	level := &slog.LevelVar{}
	level.Set(cfg.Level)

	opts := &slog.HandlerOptions{
		AddSource: cfg.CallerInfo,
		Level:     level,
	}

	var output io.Writer
	if cfg.Output == nil {
		output = os.Stdout
	} else {
		output = cfg.Output
	}

	handler := cfg.Handler
	if handler == nil {
		opts.ReplaceAttr = ReplaceLevelName
		consoleOpts := &ConsoleHandlerOptions{
			SlogOptions: opts,
			TimeFormat:  cfg.TimeFormat,
		}
		if cfg.TimeFormat == "" {
			consoleOpts.TimeFormat = DefaultTimeFormat
		}
		handler = NewConsoleHandler(output, consoleOpts)
	}

	l := slog.New(handler)
	logger := &Logger{
		logger:     l,
		handler:    handler,
		level:      level,
		callerInfo: cfg.CallerInfo,
	}
	return logger
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (l *Logger) Named(name string) *Logger {
	newLogger := l.logger.WithGroup(name)
	return &Logger{
		logger: newLogger,
		level:  l.level,
	}
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields ...any) *Logger {
	newLogger := l.logger.With(fields...)
	return &Logger{
		logger: newLogger,
		level:  l.level,
	}
}

// Enabled reports whether l emits log records at the given context and level.
// nolint: contextcheck
func (l *Logger) Enabled(ctx context.Context, level Level) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	return l.handler.Enabled(ctx, level)
}

// Level returns the minimum enabled log level.
func (l *Logger) Level() Level {
	return l.level.Level()
}

// SetLevel alters the logging level.
func (l *Logger) SetLevel(level Level) {
	l.level.Set(level)
}

// Trace logs at TraceLevel.
func (l *Logger) Trace(msg string, args ...any) {
	l.Log(nil, TraceLevel, msg, args...)
}

// TraceContext logs at TraceLevel with the given context.
func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, TraceLevel, msg, args...)
}

// Debug logs at LevelDebug.
func (l *Logger) Debug(msg string, args ...any) {
	l.Log(nil, DebugLevel, msg, args...)
}

// DebugContext logs at LevelDebug with the given context.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, DebugLevel, msg, args...)
}

// Info logs at LevelInfo.
func (l *Logger) Info(msg string, args ...any) {
	l.Log(nil, InfoLevel, msg, args...)
}

// InfoContext logs at LevelInfo with the given context.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, InfoLevel, msg, args...)
}

// Warn logs at LevelWarn.
func (l *Logger) Warn(msg string, args ...any) {
	l.Log(nil, WarnLevel, msg, args...)
}

// WarnContext logs at LevelWarn with the given context.
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, WarnLevel, msg, args...)
}

// Error logs at LevelError.
func (l *Logger) Error(msg string, args ...any) {
	l.Log(nil, ErrorLevel, msg, args...)
}

// ErrorContext logs at LevelError with the given context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, ErrorLevel, msg, args...)
}

// Fatal logs at FatalLevel.
func (l *Logger) Fatal(msg string, args ...any) {
	l.Log(nil, FatalLevel, msg, args...)
	fatalExitFunc()
}

// FatalContext logs at FatalLevel with the given context.
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, FatalLevel, msg, args...)
	fatalExitFunc()
}

// Log emits a log record with the current time and the given level and message.
// nolint: contextcheck
func (l *Logger) Log(ctx context.Context, level Level, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}

	if !l.handler.Enabled(ctx, level) {
		return
	}

	r := slog.Record{
		Time:    time.Now(),
		Message: msg,
		Level:   level,
	}

	if l.callerInfo {
		var pcs [1]uintptr
		runtime.Callers(3, pcs[:])
		r.PC = pcs[0]
	}

	r.Add(args...)
	_ = l.handler.Handle(ctx, r)
}

// fatalExitFunc defines the function to call when exiting due to a fatal log error.
// This is used in unit tests.
var fatalExitFunc = fatalExit

func fatalExit() {
	os.Exit(1)
}
