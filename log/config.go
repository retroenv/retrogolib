package log

import (
	"io"
	"log/slog"
)

// DefaultTimeFormat is the compact timestamp format used by the default console handler.
const DefaultTimeFormat = "2006-01-02 15:04:05"

// Config represents configuration for a logger.
type Config struct {
	// CallerInfo adds a ("source", "file:line") attribute to the output
	// indicating the source code position of the log statement.
	CallerInfo bool

	// Level is the minimum enabled level. Its zero value is InfoLevel.
	Level Level

	// Output receives records from the default console handler. Nil uses os.Stdout.
	Output io.Writer

	// Handler overrides the default console handler when non-nil.
	Handler slog.Handler

	// TimeFormat controls timestamps in the default console handler. It defaults
	// to DefaultTimeFormat; "-" disables timestamps.
	TimeFormat string
}

// DefaultConfig returns an adjustable Config that can be passed to NewWithConfig.
func DefaultConfig() Config {
	return Config{
		Level:      DefaultLevel(),
		TimeFormat: DefaultTimeFormat,
	}
}
