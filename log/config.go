package log

import (
	"io"
	"log/slog"
)

// DefaultTimeFormat is a slimmer default time format used if no other time format is specified.
const DefaultTimeFormat = "2006-01-02 15:04:05"

// Config represents configuration for a logger.
type Config struct {
	// CallerInfo adds a ("source", "file:line") attribute to the output
	// indicating the source code position of the log statement.
	CallerInfo bool

	Level Level

	Output io.Writer

	// Handler handles log records produced by a Logger..
	Handler slog.Handler

	// TimeFormat defines the time format to use, defaults to "2006-01-02 15:04:05"
	// Outputting of time can be disabled with - for the console handler.
	TimeFormat string
}

// DefaultConfig returns the default config. The returned config can be adjusted
// and used to create a logger with custom config using the NewWithConfig() function.
func DefaultConfig() Config {
	cfg := Config{
		Level:      DefaultLevel(),
		TimeFormat: DefaultTimeFormat,
	}
	return cfg
}
