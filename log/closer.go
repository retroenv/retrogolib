package log

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"syscall"
)

// Closer closes a resource and logs unexpected errors. It is useful for
// deferred cleanup:
//
//	defer logger.Closer(resp.Body, "Closing response body")
//
// Common close errors such as os.ErrClosed and network disconnects are ignored.
func (l *Logger) Closer(closer io.Closer, msg string) {
	if closer == nil {
		return
	}

	err := closer.Close()
	if shouldIgnoreCloseError(err) {
		return
	}

	l.Error(msg, Err(err))
}

// CloserCtx closes a context-aware resource and logs unexpected errors.
// Deadline and cancellation errors include a "reason" field in the log record.
func (l *Logger) CloserCtx(ctx context.Context, closer closerCtx, msg string) {
	if closer == nil {
		return
	}

	err := closer.Close(ctx)
	if shouldIgnoreCloseError(err) {
		return
	}

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

// MultiCloser attempts every resource even after a close fails, logging each
// unexpected error with its position in closers.
func (l *Logger) MultiCloser(msg string, closers ...io.Closer) {
	for i, closer := range closers {
		if closer == nil {
			continue
		}
		err := closer.Close()
		if shouldIgnoreCloseError(err) {
			continue
		}

		l.Error(msg, Err(err), Int("closer_index", i))
	}
}

// MultiCloserCtx attempts every context-aware resource even after a close
// fails, logging each unexpected error with its position in closers.
func (l *Logger) MultiCloserCtx(ctx context.Context, msg string, closers ...closerCtx) {
	for i, closer := range closers {
		if closer == nil {
			continue
		}

		err := closer.Close(ctx)
		if shouldIgnoreCloseError(err) {
			continue
		}

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

type closerCtx interface {
	Close(ctx context.Context) error
}

func shouldIgnoreCloseError(err error) bool {
	if err == nil {
		return true
	}

	if errors.Is(err, os.ErrClosed) || errors.Is(err, net.ErrClosed) {
		return true
	}
	if errors.Is(err, io.EOF) || errors.Is(err, syscall.EBADF) || errors.Is(err, syscall.EINVAL) {
		return true
	}

	// Some net.OpError values expose only platform error text instead of a sentinel.
	var opErr *net.OpError
	if errors.As(err, &opErr) && opErr.Err != nil {
		switch opErr.Err.Error() {
		case "use of closed network connection", "broken pipe", "connection reset by peer":
			return true
		}
	}

	return false
}
