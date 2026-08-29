package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const consoleLevelWidth = 8

var _ slog.Handler = &ConsoleHandler{}

// ConsoleHandlerOptions configures a ConsoleHandler.
// Nil SlogOptions uses the slog defaults.
type ConsoleHandlerOptions struct {
	// SlogOptions controls level filtering, source reporting, and attribute replacement.
	SlogOptions *slog.HandlerOptions

	// TimeFormat controls timestamps. Its zero value uses time.RFC3339; "-" disables timestamps.
	TimeFormat string
}

// consoleOutput keeps the text prefix and JSON attributes atomic across handler clones.
type consoleOutput struct {
	mu sync.Mutex
	w  io.Writer
}

// ConsoleHandler formats records for human-readable console output.
type ConsoleHandler struct {
	opts            ConsoleHandlerOptions
	internalHandler slog.Handler
	// Inherited attributes require a JSON suffix even when a record has none.
	hasAttrs bool
	output   *consoleOutput
}

// NewConsoleHandler returns a new console handler.
func NewConsoleHandler(w io.Writer, opts *ConsoleHandlerOptions) *ConsoleHandler {
	var slogOpts slog.HandlerOptions
	timeFormat := time.RFC3339
	if opts != nil {
		if opts.SlogOptions != nil {
			slogOpts = *opts.SlogOptions
		}
		if opts.TimeFormat != "" {
			timeFormat = opts.TimeFormat
		}
	}

	internalOpts := slog.HandlerOptions{
		AddSource: slogOpts.AddSource,
		Level:     slogOpts.Level,
	}

	internalOpts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if slogOpts.AddSource && a.Key == slog.SourceKey {
			return slog.Attr{}
		}

		if slogOpts.ReplaceAttr != nil {
			return slogOpts.ReplaceAttr(groups, a)
		}
		return a
	}

	stableOpts := ConsoleHandlerOptions{
		SlogOptions: &slogOpts,
		TimeFormat:  timeFormat,
	}

	return &ConsoleHandler{
		opts:            stableOpts,
		internalHandler: slog.NewJSONHandler(w, &internalOpts),
		output:          &consoleOutput{w: w},
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
	// Most console records fit without growing the buffer again.
	buf.Grow(256)

	if h.opts.TimeFormat != "-" {
		buf.WriteString(r.Time.Format(h.opts.TimeFormat))
		buf.WriteString("  ")
	}

	writeConsoleLevel(&buf, r.Level)

	if h.opts.SlogOptions.AddSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		frame, _ := fs.Next()
		if frame.File != "" {
			buf.WriteString(frame.File)
			buf.WriteRune(':')
			buf.WriteString(strconv.Itoa(frame.Line))
			buf.WriteRune(' ')
		}
	}

	buf.WriteString(r.Message)

	hasEntries := h.hasAttrs
	r.Attrs(func(a slog.Attr) bool {
		if !a.Equal(slog.Attr{}) {
			hasEntries = true
			return false
		}
		return true
	})
	if hasEntries {
		buf.WriteRune(' ')
	} else {
		buf.WriteRune('\n')
	}

	h.output.mu.Lock()
	defer h.output.mu.Unlock()

	_, err := h.output.w.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("writing log prefix: %w", err)
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
//
//nolint:ireturn // Required by slog.Handler.
func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ConsoleHandler{
		opts:            h.opts,
		internalHandler: h.internalHandler.WithAttrs(attrs),
		hasAttrs:        h.hasAttrs || hasAttrs(attrs),
		output:          h.output,
	}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
//
//nolint:ireturn // Required by slog.Handler.
func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return &ConsoleHandler{
		opts:            h.opts,
		internalHandler: h.internalHandler.WithGroup(name),
		hasAttrs:        h.hasAttrs,
		output:          h.output,
	}
}

func writeConsoleLevel(buf *bytes.Buffer, level slog.Level) {
	text := level.String()
	switch level {
	case TraceLevel:
		text = "TRACE"
	case FatalLevel:
		text = "FATAL"
	}

	buf.WriteString(text)
	for range max(consoleLevelWidth-len(text), 1) {
		buf.WriteByte(' ')
	}
}

func hasAttrs(attrs []slog.Attr) bool {
	for _, attr := range attrs {
		if !attr.Equal(slog.Attr{}) {
			return true
		}
	}

	return false
}
