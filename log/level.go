package log

import (
	"log/slog"
	"sync/atomic"
)

// Log levels.
const (
	// TraceLevel logs are typically voluminous, and are usually disabled in
	// production.
	TraceLevel = slog.LevelDebug << 1

	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = slog.LevelDebug

	// InfoLevel is the default logging priority.
	InfoLevel = slog.LevelInfo

	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = slog.LevelWarn

	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = slog.LevelError

	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = slog.LevelError << 1
)

// Level is a logging priority. Higher levels are more important.
type Level = slog.Level

var (
	defaultLevel   = uintptr(InfoLevel)
	fatalLevelText = slog.StringValue("FATAL")
	traceLevelText = slog.StringValue("TRACE")
)

// DefaultLevel returns the current default level for all loggers
// newly created with New().
func DefaultLevel() Level {
	return Level(atomic.LoadUintptr(&defaultLevel))
}

// SetDefaultLevel sets the default level for all newly created loggers.
func SetDefaultLevel(level Level) {
	atomic.StoreUintptr(&defaultLevel, uintptr(level))
}

// ReplaceLevelName sets custom defined level names for outputting.
func ReplaceLevelName(_ []string, a slog.Attr) slog.Attr {
	if a.Key != slog.LevelKey {
		return a
	}

	level, ok := a.Value.Any().(slog.Level)
	if !ok {
		return a
	}

	switch level {
	case TraceLevel:
		a.Value = traceLevelText
	case FatalLevel:
		a.Value = fatalLevelText
	}

	return a
}
