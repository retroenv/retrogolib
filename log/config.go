package log

import (
	"io"

	"golang.org/x/exp/slog"
)

// Config represents configuration for a logger.
type Config struct {
	// CallerInfo adds a ("source", "file:line") attribute to the output
	// indicating the source code position of the log statement.
	CallerInfo bool

	Level Level

	Output io.Writer

	// Handler handles log records produced by a Logger..
	Handler slog.Handler
}

// DefaultConfig returns the default config. The returned config can be adjusted
// and used to create a logger with custom config using the NewWithConfig() function.
func DefaultConfig() Config {
	cfg := Config{
		Level: DefaultLevel(),
	}
	return cfg
}
