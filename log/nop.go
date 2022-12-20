package log

import (
	"io"
)

// NewNop creates a no-op logger which never writes logs to the output.
// Useful for tests.
func NewNop() *Logger {
	cfg := Config{
		Output: io.Discard,
		Level:  Level(100),
	}

	return NewWithConfig(cfg)
}
