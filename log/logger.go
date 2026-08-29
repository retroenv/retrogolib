package log

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Logger provides fast, leveled, structured logging. All methods are safe
// for concurrent use.
type Logger struct {
	logger     *slog.Logger
	callerInfo bool
	level      *slog.LevelVar
}

// New returns a Logger using the current default level and console output.
func New() *Logger {
	return NewWithConfig(Config{
		Level: DefaultLevel(),
	})
}

// NewWithConfig creates a logger with the given configuration.
func NewWithConfig(cfg Config) *Logger {
	level := &slog.LevelVar{}
	level.Set(cfg.Level)

	opts := &slog.HandlerOptions{
		AddSource: cfg.CallerInfo,
		Level:     level,
	}

	output := cfg.Output
	if cfg.Output == nil {
		output = os.Stdout
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

	return &Logger{
		logger:     slog.New(handler),
		level:      level,
		callerInfo: cfg.CallerInfo,
	}
}

// Named scopes subsequent fields under name. Repeated calls create nested
// groups. The child shares the parent's dynamic level.
func (l *Logger) Named(name string) *Logger {
	if !l.valid() {
		return l
	}

	return &Logger{
		logger:     l.logger.WithGroup(name),
		callerInfo: l.callerInfo,
		level:      l.level,
	}
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa. The child shares the
// parent's dynamic level.
func (l *Logger) With(fields ...any) *Logger {
	if !l.valid() {
		return l
	}

	return &Logger{
		logger:     l.logger.With(fields...),
		callerInfo: l.callerInfo,
		level:      l.level,
	}
}

// Enabled reports whether l emits log records at the given context and level.
//
//nolint:contextcheck // A nil context is intentionally normalized for convenience.
func (l *Logger) Enabled(ctx context.Context, level Level) bool {
	if !l.valid() || l.level.Level() > level {
		return false
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return l.logger.Enabled(ctx, level)
}

// Level returns the minimum enabled log level.
func (l *Logger) Level() Level {
	if !l.valid() {
		return InfoLevel
	}

	return l.level.Level()
}

// SetLevel alters the logging level.
func (l *Logger) SetLevel(level Level) {
	if !l.valid() {
		return
	}

	l.level.Set(level)
}

// Trace logs at TraceLevel.
func (l *Logger) Trace(msg string, args ...Field) {
	l.log(nil, TraceLevel, msg, args...)
}

// TraceContext logs at TraceLevel with the given context.
func (l *Logger) TraceContext(ctx context.Context, msg string, args ...Field) {
	l.log(ctx, TraceLevel, msg, args...)
}

// Debug logs at DebugLevel.
func (l *Logger) Debug(msg string, args ...Field) {
	l.log(nil, DebugLevel, msg, args...)
}

// DebugContext logs at DebugLevel with the given context.
func (l *Logger) DebugContext(ctx context.Context, msg string, args ...Field) {
	l.log(ctx, DebugLevel, msg, args...)
}

// Info logs at InfoLevel.
func (l *Logger) Info(msg string, args ...Field) {
	l.log(nil, InfoLevel, msg, args...)
}

// InfoContext logs at InfoLevel with the given context.
func (l *Logger) InfoContext(ctx context.Context, msg string, args ...Field) {
	l.log(ctx, InfoLevel, msg, args...)
}

// Warn logs at WarnLevel.
func (l *Logger) Warn(msg string, args ...Field) {
	l.log(nil, WarnLevel, msg, args...)
}

// WarnContext logs at WarnLevel with the given context.
func (l *Logger) WarnContext(ctx context.Context, msg string, args ...Field) {
	l.log(ctx, WarnLevel, msg, args...)
}

// Error logs at ErrorLevel.
func (l *Logger) Error(msg string, args ...Field) {
	l.log(nil, ErrorLevel, msg, args...)
}

// ErrorContext logs at ErrorLevel with the given context.
func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...Field) {
	l.log(ctx, ErrorLevel, msg, args...)
}

// Fatal logs at FatalLevel, then terminates the process with exit code 1.
func (l *Logger) Fatal(msg string, args ...Field) {
	l.log(nil, FatalLevel, msg, args...)
	fatalExitFunc()
}

// FatalContext logs at FatalLevel with the given context, then terminates the
// process with exit code 1.
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...Field) {
	l.log(ctx, FatalLevel, msg, args...)
	fatalExitFunc()
}

// Log emits a log record with the current time and the given level and message.
//
//nolint:contextcheck // A nil context is intentionally normalized for convenience.
func (l *Logger) Log(ctx context.Context, level Level, msg string, args ...Field) {
	l.log(ctx, level, msg, args...)
}

//nolint:contextcheck // A nil context is intentionally normalized for convenience.
func (l *Logger) log(ctx context.Context, level Level, msg string, args ...Field) {
	if !l.Enabled(ctx, level) {
		return
	}

	if ctx == nil {
		ctx = context.Background()
	}

	r := slog.Record{
		Time:    time.Now(),
		Message: msg,
		Level:   level,
	}

	if l.callerInfo {
		var pcs [1]uintptr
		// Skip runtime.Callers, log, and the public logging method.
		if runtime.Callers(3, pcs[:]) == 1 {
			r.PC = pcs[0]
		}
	}

	r.AddAttrs(args...)
	// Logging methods cannot return handler errors without breaking their API.
	_ = l.logger.Handler().Handle(ctx, r)
}

func (l *Logger) valid() bool {
	return l != nil && l.logger != nil && l.level != nil
}

// fatalExitFunc defines the function to call when exiting due to a fatal log error.
// This is used in unit tests.
var fatalExitFunc = fatalExit

// fatalExit terminates the program with exit code 1.
func fatalExit() {
	os.Exit(1)
}
