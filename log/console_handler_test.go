package log

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/assert"
)

func TestConsoleHandlerWithAttrsAndGroup(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(&buf, &ConsoleHandlerOptions{TimeFormat: "-"})
	logger := slog.New(handler).
		With("component", "cpu").
		WithGroup("step").
		With("cycles", 2)

	logger.Info("Executed")

	assert.Equal(t, "INFO    Executed {\"component\":\"cpu\",\"step\":{\"cycles\":2}}\n", buf.String())
}

func TestConsoleHandlerAllowsNilSlogOptions(t *testing.T) {
	var buf bytes.Buffer
	opts := &ConsoleHandlerOptions{TimeFormat: "-"}
	handler := NewConsoleHandler(&buf, opts)

	slog.New(handler).Info("Ready")

	assert.Equal(t, "INFO    Ready\n", buf.String())
	assert.Equal(t, "-", opts.TimeFormat)
	assert.Nil(t, opts.SlogOptions)
}

func TestConsoleHandlerFormatsCustomLevel(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(&buf, &ConsoleHandlerOptions{TimeFormat: "-"})
	record := slog.NewRecord(time.Time{}, slog.Level(2), "custom", 0)

	err := handler.Handle(context.Background(), record)

	assert.NoError(t, err)
	assert.Equal(t, "INFO+2  custom\n", buf.String())
}

func TestConsoleHandlerSerializesChildWrites(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(&buf, &ConsoleHandlerOptions{TimeFormat: "-"})
	logger := slog.New(handler)

	const records = 100
	var wg sync.WaitGroup
	for i := range records {
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.With("worker", i).Info("Running", "record", i)
		}()
	}
	wg.Wait()

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, records)
	for _, line := range lines {
		assert.True(t, strings.HasPrefix(line, "INFO    Running {"), "malformed line: "+line)
		assert.True(t, strings.HasSuffix(line, "}"), "malformed line: "+line)
	}
}
