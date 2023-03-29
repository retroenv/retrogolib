package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

// consoleLevelString translates a level to a padded string ready for printing on the console.
var consoleLevelString = map[Level]string{
	TraceLevel: "TRACE   ",
	DebugLevel: "DEBUG   ",
	InfoLevel:  "INFO    ",
	WarnLevel:  "WARN    ",
	ErrorLevel: "ERROR   ",
	FatalLevel: "FATAL   ",
}

// ConsoleHandler formats the logger output in a better human-readable way.
type ConsoleHandler struct {
	opts            ConsoleHandlerOptions
	internalHandler slog.Handler

	mu sync.Mutex
	w  io.Writer
}

// ConsoleHandlerOptions are options for a ConsoleHandler.
// A zero HandlerOptions consists entirely of default values.
type ConsoleHandlerOptions struct {
	SlogOptions slog.HandlerOptions

	TimeFormat string
}

// NewConsoleHandler returns a new console handler with default options.
func NewConsoleHandler(w io.Writer) *ConsoleHandler {
	opts := &ConsoleHandlerOptions{}
	return opts.NewConsoleHandler(w)
}

// NewConsoleHandler returns a new console handler based on the set options.
func (opts *ConsoleHandlerOptions) NewConsoleHandler(w io.Writer) *ConsoleHandler {
	internalOpts := opts.SlogOptions
	timeFormat := opts.TimeFormat
	if timeFormat == "" {
		opts.TimeFormat = time.RFC3339
	}

	internalOpts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if opts.SlogOptions.AddSource && a.Key == slog.SourceKey {
			return slog.Attr{}
		}

		rep := opts.SlogOptions.ReplaceAttr
		if rep != nil {
			return rep(groups, a)
		}
		return a
	}

	return &ConsoleHandler{
		opts:            *opts,
		w:               w,
		internalHandler: internalOpts.NewJSONHandler(w),
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *ConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.internalHandler.Enabled(ctx, level)
}

// Handle handles the Record.
func (h *ConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	var buf bytes.Buffer

	if h.opts.TimeFormat != "-" {
		buf.WriteString(r.Time.Format(h.opts.TimeFormat))
		buf.WriteString("  ")
	}

	buf.WriteString(consoleLevelString[r.Level])

	if h.opts.SlogOptions.AddSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		frame, _ := fs.Next()
		if frame.File != "" {
			buf.WriteString(frame.File)
			buf.WriteRune(':')
			buf.Write([]byte(strconv.Itoa(frame.Line)))
			buf.WriteRune(' ')
		}
	}

	buf.WriteString(r.Message)

	hasEntries := false
	r.Attrs(func(a slog.Attr) {
		if a.Key != "" {
			hasEntries = true
		}
	})
	if hasEntries {
		buf.WriteRune(' ')
	} else {
		buf.WriteRune('\n')
	}

	h.mu.Lock()
	_, err := h.w.Write(buf.Bytes())
	h.mu.Unlock()

	if err != nil {
		return fmt.Errorf("writing to buffer: %w", err)
	}

	if hasEntries {
		if err := h.internalHandler.Handle(ctx, r); err != nil {
			return fmt.Errorf("handling record: %w", err)
		}
	}
	return nil
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// nolint: ireturn
func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ConsoleHandler{
		opts:            h.opts,
		internalHandler: h.internalHandler.WithAttrs(attrs),
		w:               h.w,
	}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// nolint: ireturn
func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return &ConsoleHandler{
		opts:            h.opts,
		internalHandler: h.internalHandler.WithGroup(name),
		w:               h.w,
	}
}
