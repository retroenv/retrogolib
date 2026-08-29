package log

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
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
	originalExitFunc := fatalExitFunc
	defer func() { fatalExitFunc = originalExitFunc }()
	fatalExitFunc = func() {
		exited = true
	}

	logger.Fatal("Something bad happened", Err(errors.New("network error")))

	assert.True(t, exited)
	output := buf.String()
	assert.Equal(t, "FATAL   Something bad happened {\"error\":\"network error\"}\n", output)
}

func TestLoggerTrace(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer

	cfg.Level = TraceLevel
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	logger.Trace("Something happened")

	output := buf.String()
	assert.Equal(t, "TRACE   Something happened\n", output)
}

func TestLoggerChildrenPreserveConfiguration(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.CallerInfo = true
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	child := logger.With(String("component", "cpu")).Named("step")
	child.Info("Executed", Int("cycles", 2))

	output := buf.String()
	assert.Contains(t, output, "logger_test.go")
	assert.Contains(t, output, "Executed")
	assert.Contains(t, output, `"component":"cpu"`)
	assert.Contains(t, output, `"step":{"cycles":2}`)

	buf.Reset()
	child.SetLevel(ErrorLevel)
	logger.Info("Filtered")
	logger.Error("Visible")
	assert.NotContains(t, buf.String(), "Filtered")
	assert.Contains(t, buf.String(), "Visible")
}

func TestLoggerLogUsesHandlerAndLevel(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: DebugLevel})
	logger := NewWithConfig(Config{Level: InfoLevel, Handler: handler})

	logger.Log(nil, DebugLevel, "Filtered")
	logger.Log(nil, InfoLevel, "Visible")

	assert.NotContains(t, buf.String(), "Filtered")
	assert.Contains(t, buf.String(), "Visible")
}

func TestLoggerNilSafety(t *testing.T) {
	var logger *Logger

	assert.False(t, logger.Enabled(nil, InfoLevel))
	assert.Equal(t, InfoLevel, logger.Level())
	logger.SetLevel(DebugLevel)
	logger.Trace("Trace")
	logger.Debug("Debug")
	logger.Info("Info")
	logger.Warn("Warn")
	logger.Error("Error")
	logger.Log(nil, InfoLevel, "Log")
	assert.Nil(t, logger.Named("child"))
	assert.Nil(t, logger.With(String("key", "value")))

	zero := &Logger{}
	assert.False(t, zero.Enabled(nil, InfoLevel))
	zero.Info("Info")
}

func TestLoggerDoesNotResolveDisabledFields(t *testing.T) {
	logger := NewNop()
	called := false
	logger.Debug("Disabled", StringFunc("value", func() string {
		called = true
		return "computed"
	}))

	assert.False(t, called)
}

func TestLoggerCaller(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer

	cfg.CallerInfo = true
	cfg.Level = TraceLevel
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)

	logger.Trace("Something happened")

	output := buf.String()
	assert.True(t, strings.Contains(output, "TRACE"))
	assert.True(t, strings.Contains(output, "logger_test.go"))
	assert.True(t, strings.Contains(output, "Something happened\n"))
}
