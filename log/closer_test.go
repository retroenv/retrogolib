package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/assert"
)

const (
	closingFailedMsg  = "closing failed"
	closingTimeoutMsg = "closing with timeout"
	benchmarkCloseMsg = "benchmark close"
	benchmarkMultiMsg = "benchmark multi close"
	benchmarkErrorMsg = "benchmark close with error"
)

type testCloser struct {
	err error
}

func (t testCloser) Close() error {
	return t.err
}

type testCloserCtx struct {
	err       error
	sleepTime time.Duration
}

func (t testCloserCtx) Close(ctx context.Context) error {
	if t.sleepTime > 0 {
		select {
		case <-time.After(t.sleepTime):
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		}
	}
	return t.err
}

func TestLoggerCloser(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	closer := testCloser{}
	msg := closingFailedMsg

	// Test successful close (no error)
	logger.Closer(closer, msg)
	output := buf.String()
	assert.NotContains(t, output, "ERROR")
	assert.NotContains(t, output, msg)

	// Test error close
	errMsg := "failure"
	closer.err = errors.New(errMsg)
	logger.Closer(closer, msg)
	output = buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, msg)
	assert.Contains(t, output, errMsg)
}

func TestLoggerCloserIgnoresExpectedErrors(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	msg := closingFailedMsg

	expectedErrors := []error{
		os.ErrClosed,
		net.ErrClosed,
		io.EOF,
		syscall.EBADF,
		syscall.EINVAL,
		&net.OpError{Err: errors.New("use of closed network connection")},
		&net.OpError{Err: errors.New("broken pipe")},
		&net.OpError{Err: errors.New("connection reset by peer")},
	}

	for _, expectedErr := range expectedErrors {
		buf.Reset()
		closer := testCloser{err: expectedErr}
		logger.Closer(closer, msg)
		output := buf.String()
		assert.NotContains(t, output, "ERROR", "Expected error %v should be ignored", expectedErr)
	}
}

func TestLoggerCloserCtx(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	ctx := context.Background()
	closer := testCloserCtx{}
	msg := "closing failed"

	// Test successful close (no error)
	logger.CloserCtx(ctx, closer, msg)
	output := buf.String()
	assert.NotContains(t, output, "ERROR")
	assert.NotContains(t, output, msg)

	// Test error close
	errMsg := "failure"
	closer.err = errors.New(errMsg)
	logger.CloserCtx(ctx, closer, msg)
	output = buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, msg)
	assert.Contains(t, output, errMsg)
}

func TestLoggerCloserCtxTimeout(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Closer that takes longer than timeout
	closer := testCloserCtx{sleepTime: 50 * time.Millisecond}
	msg := closingTimeoutMsg

	logger.CloserCtx(ctx, closer, msg)
	output := buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, msg)
	assert.Contains(t, output, "context deadline exceeded")
	assert.Contains(t, output, "reason")
}

func TestLoggerCloserCtxCanceled(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	// Closer that respects cancellation
	closer := testCloserCtx{sleepTime: 50 * time.Millisecond}
	msg := "closing with cancellation"

	// Cancel immediately to trigger cancellation
	cancel()

	logger.CloserCtx(ctx, closer, msg)
	output := buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, msg)
	assert.Contains(t, output, "context canceled")
	assert.Contains(t, output, "reason")
}

func TestLoggerMultiCloser(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	msg := "closing multiple resources"

	// Test with all successful closers
	closers := []io.Closer{
		testCloser{},
		testCloser{},
		testCloser{},
	}
	logger.MultiCloser(msg, closers...)
	output := buf.String()
	assert.NotContains(t, output, "ERROR")

	// Test with some failing closers
	buf.Reset()
	closers = []io.Closer{
		testCloser{},
		testCloser{err: errors.New("first failure")},
		testCloser{},
		testCloser{err: errors.New("second failure")},
	}
	logger.MultiCloser(msg, closers...)
	output = buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, "first failure")
	assert.Contains(t, output, "second failure")
	assert.Contains(t, output, "closer_index")
	assert.Contains(t, output, "1") // First failure at index 1
	assert.Contains(t, output, "3") // Second failure at index 3
}

func TestLoggerMultiCloserWithNilClosers(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	msg := "closing with nil closers"

	// Test with nil closers mixed in
	closers := []io.Closer{
		testCloser{},
		nil,
		testCloser{err: errors.New("failure")},
		nil,
	}
	logger.MultiCloser(msg, closers...)
	output := buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, "failure")
	assert.Contains(t, output, "closer_index")
	assert.Contains(t, output, "2") // Failure at index 2
}

func TestLoggerMultiCloserCtx(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	ctx := context.Background()
	msg := "closing multiple context resources"

	// Test with all successful closers
	closers := []closerCtx{
		testCloserCtx{},
		testCloserCtx{},
		testCloserCtx{},
	}
	logger.MultiCloserCtx(ctx, msg, closers...)
	output := buf.String()
	assert.NotContains(t, output, "ERROR")

	// Test with some failing closers
	buf.Reset()
	closers = []closerCtx{
		testCloserCtx{},
		testCloserCtx{err: errors.New("first failure")},
		testCloserCtx{},
		testCloserCtx{err: errors.New("second failure")},
	}
	logger.MultiCloserCtx(ctx, msg, closers...)
	output = buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, "first failure")
	assert.Contains(t, output, "second failure")
	assert.Contains(t, output, "closer_index")
}

func TestLoggerMultiCloserCtxTimeout(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	msg := closingTimeoutMsg

	// Test with some closers that timeout
	closers := []closerCtx{
		testCloserCtx{},
		testCloserCtx{sleepTime: 50 * time.Millisecond}, // This will timeout
		testCloserCtx{},
	}
	logger.MultiCloserCtx(ctx, msg, closers...)
	output := buf.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, "context deadline exceeded")
	assert.Contains(t, output, "closer_index")
	assert.Contains(t, output, "1") // Timeout at index 1
}

func TestShouldIgnoreCloseError(t *testing.T) {
	logger := New()

	// Test nil error
	assert.True(t, logger.shouldIgnoreCloseError(nil))

	// Test expected errors
	expectedErrors := []error{
		os.ErrClosed,
		net.ErrClosed,
		io.EOF,
		syscall.EBADF,
		syscall.EINVAL,
	}

	for _, err := range expectedErrors {
		assert.True(t, logger.shouldIgnoreCloseError(err), "Should ignore %v", err)
	}

	// Test network operation errors
	networkErrors := []*net.OpError{
		{Err: errors.New("use of closed network connection")},
		{Err: errors.New("broken pipe")},
		{Err: errors.New("connection reset by peer")},
	}

	for _, err := range networkErrors {
		assert.True(t, logger.shouldIgnoreCloseError(err), "Should ignore %v", err)
	}

	// Test errors that should NOT be ignored
	unexpectedErrors := []error{
		errors.New("unexpected error"),
		&net.OpError{Err: errors.New("some other network error")},
	}

	for _, err := range unexpectedErrors {
		assert.False(t, logger.shouldIgnoreCloseError(err), "Should NOT ignore %v", err)
	}
}

func TestLoggerCloserIntegrationWithRealTypes(t *testing.T) {
	cfg := DefaultConfig()
	var buf bytes.Buffer
	cfg.Output = &buf
	cfg.TimeFormat = "-"

	logger := NewWithConfig(cfg)

	// Create a simple closer that behaves like a real resource

	type realCloser struct {
		closed bool
	}

	rc := &realCloser{}
	closer := &struct {
		*realCloser
		io.Closer
	}{
		realCloser: rc,
		Closer: closerFunc(func() error {
			if rc.closed {
				return os.ErrClosed
			}
			rc.closed = true
			return nil
		}),
	}

	// First close should succeed without logging
	logger.Closer(closer, "closing real resource")
	output := buf.String()
	assert.NotContains(t, output, "ERROR")

	// Second close should be ignored (os.ErrClosed)
	buf.Reset()
	logger.Closer(closer, "closing already closed resource")
	output = buf.String()
	assert.NotContains(t, output, "ERROR")
}

// closerFunc is a function type that implements io.Closer
type closerFunc func() error

func (f closerFunc) Close() error {
	return f()
}

func BenchmarkLoggerCloser(b *testing.B) {
	logger := New()
	closer := testCloser{}

	b.ResetTimer()
	for range b.N {
		logger.Closer(closer, benchmarkCloseMsg)
	}
}

func BenchmarkLoggerCloserWithError(b *testing.B) {
	cfg := DefaultConfig()
	cfg.Output = io.Discard // Discard output for pure benchmark
	logger := NewWithConfig(cfg)
	closer := testCloser{err: errors.New("benchmark error")}

	b.ResetTimer()
	for range b.N {
		logger.Closer(closer, benchmarkErrorMsg)
	}
}

func BenchmarkLoggerMultiCloser(b *testing.B) {
	logger := New()
	closers := []io.Closer{
		testCloser{},
		testCloser{},
		testCloser{},
	}

	b.ResetTimer()
	for range b.N {
		logger.MultiCloser(benchmarkMultiMsg, closers...)
	}
}
