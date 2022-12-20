package log

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNewTestLogger(t *testing.T) {
	logger := NewTestLogger(t)
	assert.Equal(t, DebugLevel, logger.level.Level())
}
