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
	logger     *slog.Logger   // underlying slog logger instance
	handler    slog.Handler   // handler for processing log records
	callerInfo bool           // whether to include caller info in logs
	level      *slog.LevelVar // dynamic log level variable
}

// New returns a new Logger instance.
func New() *Logger {
	// Create default config with standard level
	cfg := Config{
		Level: DefaultLevel(),
	}

	return NewWithConfig(cfg)
}

// NewWithConfig creates a new logger for the given config.
// If no level is set in the config, it will use the default level of
// this package.
func NewWithConfig(cfg Config) *Logger {
	// Set up dynamic level variable
	level := &slog.LevelVar{}
	level.Set(cfg.Level)

	// Configure handler options
	opts := &slog.HandlerOptions{
		AddSource: cfg.CallerInfo,
		Level:     level,
	}

	// Default to stdout if no output specified
	var output io.Writer
	if cfg.Output == nil {
		output = os.Stdout
	} else {
		output = cfg.Output
	}

	// Use provided handler or create default console handler
	handler := cfg.Handler
	if handler == nil {
		// Set up console handler with custom formatting
		opts.ReplaceAttr = ReplaceLevelName
		consoleOpts := &ConsoleHandlerOptions{
			SlogOptions: opts,
			TimeFormat:  cfg.TimeFormat,
		}
		// Use default time format if none specified
		if cfg.TimeFormat == "" {
			consoleOpts.TimeFormat = DefaultTimeFormat
		}
		handler = NewConsoleHandler(output, consoleOpts)
	}

	// Create underlying slog logger
	l := slog.New(handler)
	// Wrap in our Logger struct
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
	// Create new logger with group namespace
	newLogger := l.logger.WithGroup(name)
	// Return wrapped logger with shared level
	return &Logger{
		logger: newLogger,
		level:  l.level,
	}
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields ...any) *Logger {
	// Create new logger with additional fields
	newLogger := l.logger.With(fields...)
	// Return wrapped logger with shared level
	return &Logger{
		logger: newLogger,
		level:  l.level,
	}
}

// Enabled reports whether l emits log records at the given context and level.
// nolint: contextcheck
func (l *Logger) Enabled(ctx context.Context, level Level) bool {
	// Use background context if none provided
	if ctx == nil {
		ctx = context.Background()
	}
	// Check if handler will process this level
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
func (l *Logger) Trace(msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= TraceLevel {
		l.Log(nil, TraceLevel, msg, args...)
	}
}

// TraceContext logs at TraceLevel with the given context.
func (l *Logger) TraceContext(ctx context.Context, msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= TraceLevel {
		l.Log(ctx, TraceLevel, msg, args...)
	}
}

// Debug logs at LevelDebug.
func (l *Logger) Debug(msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= DebugLevel {
		l.Log(nil, DebugLevel, msg, args...)
	}
}

// DebugContext logs at LevelDebug with the given context.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= DebugLevel {
		l.Log(ctx, DebugLevel, msg, args...)
	}
}

// Info logs at LevelInfo.
func (l *Logger) Info(msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= InfoLevel {
		l.Log(nil, InfoLevel, msg, args...)
	}
}

// InfoContext logs at LevelInfo with the given context.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= InfoLevel {
		l.Log(ctx, InfoLevel, msg, args...)
	}
}

// Warn logs at LevelWarn.
func (l *Logger) Warn(msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= WarnLevel {
		l.Log(nil, WarnLevel, msg, args...)
	}
}

// WarnContext logs at LevelWarn with the given context.
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= WarnLevel {
		l.Log(ctx, WarnLevel, msg, args...)
	}
}

// Error logs at LevelError.
func (l *Logger) Error(msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= ErrorLevel {
		l.Log(nil, ErrorLevel, msg, args...)
	}
}

// ErrorContext logs at LevelError with the given context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...Field) {
	// Early exit if logger nil or level too high
	if l != nil && l.level.Level() <= ErrorLevel {
		l.Log(ctx, ErrorLevel, msg, args...)
	}
}

// Fatal logs at FatalLevel.
func (l *Logger) Fatal(msg string, args ...Field) {
	// Always log fatal messages if logger exists
	if l != nil {
		l.Log(nil, FatalLevel, msg, args...)
	}
	// Exit the program after logging
	fatalExitFunc()
}

// FatalContext logs at FatalLevel with the given context.
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...Field) {
	// Always log fatal messages if logger exists
	if l != nil {
		l.Log(ctx, FatalLevel, msg, args...)
	}
	// Exit the program after logging
	fatalExitFunc()
}

// Log emits a log record with the current time and the given level and message.
// nolint: contextcheck
func (l *Logger) Log(ctx context.Context, level Level, msg string, args ...Field) {
	// Use background context if none provided
	if ctx == nil {
		ctx = context.Background()
	}

	// Skip if handler won't process this level
	if !l.handler.Enabled(ctx, level) {
		return
	}

	// Create log record with current time
	r := slog.Record{
		Time:    time.Now(),
		Message: msg,
		Level:   level,
	}

	// Add caller info if enabled
	if l.callerInfo {
		var pcs [1]uintptr
		// Skip 3 frames to get actual caller
		runtime.Callers(3, pcs[:])
		r.PC = pcs[0]
	}

	// Convert fields to interface slice
	fields := make([]any, len(args))
	for i, arg := range args {
		fields[i] = arg
	}
	// Add fields to record
	r.Add(fields...)
	// Send to handler (ignore error)
	_ = l.handler.Handle(ctx, r)
}

// fatalExitFunc defines the function to call when exiting due to a fatal log error.
// This is used in unit tests.
var fatalExitFunc = fatalExit

// fatalExit terminates the program with exit code 1.
func fatalExit() {
	os.Exit(1)
}
