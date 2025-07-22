package log

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"syscall"
)

// Closer calls the closer function and if an error gets returned it logs an error.
// This function is useful when using patterns like defer resp.Body.Close() which now become:
// defer logger.Closer(resp.Body, "closing body").
// It filters out common expected errors like os.ErrClosed and network connection errors.
func (l *Logger) Closer(closer io.Closer, msg string) {
	err := closer.Close()
	if l.shouldIgnoreCloseError(err) {
		return
	}

	l.Error(msg, Err(err))
}

// closerCtx is the interface that wraps the extended Close method.
type closerCtx interface {
	Close(ctx context.Context) error
}

// CloserCtx calls the closer function and if an error gets returned it logs an error.
// It respects context deadlines and cancellation, logging timeout errors appropriately.
func (l *Logger) CloserCtx(ctx context.Context, closer closerCtx, msg string) {
	err := closer.Close(ctx)
	if l.shouldIgnoreCloseError(err) {
		return
	}

	// Add context information for timeout/cancellation errors
	if errors.Is(err, context.DeadlineExceeded) {
		l.ErrorContext(ctx, msg, Err(err), String("reason", "context deadline exceeded"))
		return
	}
	if errors.Is(err, context.Canceled) {
		l.ErrorContext(ctx, msg, Err(err), String("reason", "context canceled"))
		return
	}

	l.ErrorContext(ctx, msg, Err(err))
}

// MultiCloser calls multiple closer functions and logs any errors.
// It continues closing all resources even if some fail, logging each error separately.
func (l *Logger) MultiCloser(msg string, closers ...io.Closer) {
	for i, closer := range closers {
		if closer == nil {
			continue
		}
		err := closer.Close()
		if l.shouldIgnoreCloseError(err) {
			continue
		}

		l.Error(msg, Err(err), Int("closer_index", i))
	}
}

// MultiCloserCtx calls multiple context-aware closer functions and logs any errors.
// It continues closing all resources even if some fail, logging each error separately.
func (l *Logger) MultiCloserCtx(ctx context.Context, msg string, closers ...closerCtx) {
	for i, closer := range closers {
		if closer == nil {
			continue
		}

		err := closer.Close(ctx)
		if l.shouldIgnoreCloseError(err) {
			continue
		}

		// Add context information for timeout/cancellation errors
		if errors.Is(err, context.DeadlineExceeded) {
			l.ErrorContext(ctx, msg, Err(err), Int("closer_index", i), String("reason", "context deadline exceeded"))
			continue
		}
		if errors.Is(err, context.Canceled) {
			l.ErrorContext(ctx, msg, Err(err), Int("closer_index", i), String("reason", "context canceled"))
			continue
		}

		l.ErrorContext(ctx, msg, Err(err), Int("closer_index", i))
	}
}

// shouldIgnoreCloseError returns true for errors that are expected and should not be logged.
func (l *Logger) shouldIgnoreCloseError(err error) bool {
	if err == nil {
		return true
	}

	// Common expected errors that should not be logged
	if errors.Is(err, os.ErrClosed) ||
		errors.Is(err, net.ErrClosed) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, syscall.EBADF) ||
		errors.Is(err, syscall.EINVAL) {

		return true
	}

	// Check for "use of closed network connection" error string
	// This is necessary because Go's net package doesn't always wrap these properly
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if opErr.Err != nil {
			errStr := opErr.Err.Error()
			if errStr == "use of closed network connection" ||
				errStr == "broken pipe" ||
				errStr == "connection reset by peer" {

				return true
			}
		}
	}

	return false
}
