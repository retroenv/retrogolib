package log

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNew(t *testing.T) {
	prev := DefaultLevel()
	SetDefaultLevel(DebugLevel)
	defer SetDefaultLevel(prev)

	logger := New()

	assert.True(t, logger.Enabled(context.TODO(), DebugLevel))
}

func TestLoggerFatal(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer

	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	exited := false
	fatalExitFunc = func() {
		exited = true
	}

	logger.Fatal("something bad happened", Err(errors.New("network error")))

	assert.True(t, exited)
	output := buf.String()
	assert.Equal(t, "FATAL   something bad happened {\"error\":\"network error\"}\n", output)
}

func TestLoggerTrace(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer

	cfg.Level = TraceLevel
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	exited := false
	fatalExitFunc = func() {
		exited = true
	}

	logger.Trace("something happened")

	assert.False(t, exited)
	output := buf.String()
	assert.Equal(t, "TRACE   something happened\n", output)
}
