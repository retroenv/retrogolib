package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/retroenv/retrogolib/log"
)

// Serve reads one JSON-RPC message per line from r and writes the responses
// to w. stdout carries MCP messages only, so callers must log to stderr.
// It returns when r reaches EOF, after the in-flight handlers have stopped.
// A message that exceeds MaxMessageSize gets an error response, and Serve
// then continues with the next line. The cancellation of ctx stops the
// in-flight handlers, but it does not stop a read that blocks on r.
func (s *Server) Serve(ctx context.Context, r io.Reader, w io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn := &connection{
		server:   s,
		w:        w,
		inflight: map[string]context.CancelFunc{},
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	reader := bufio.NewReaderSize(r, 64<<10)
	for {
		line, err := readLine(reader, s.opts.MaxMessageSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(line) > 0 {
					s.handleMessage(ctx, conn, &wg, line)
				}

				return nil
			}
			if errors.Is(err, ErrMessageTooLarge) {
				conn.writeError(ctx, nil, &Error{
					Code:    CodeInvalidRequest,
					Message: "message too large",
				})
				continue
			}
			return fmt.Errorf("reading input: %w", err)
		}

		s.handleMessage(ctx, conn, &wg, line)
	}
}

// handleMessage processes one JSON-RPC message line.
func (s *Server) handleMessage(ctx context.Context, conn *connection, wg *sync.WaitGroup, line []byte) {
	if len(bytes.TrimSpace(line)) == 0 {
		return
	}

	msg, rpcErr := parseMessage(line)
	if rpcErr != nil {
		conn.writeError(ctx, nil, rpcErr)
		return
	}

	if !msg.hasID() {
		switch {
		case msg.Method == methodCancelled:
			conn.cancelRequest(ctx, msg.Params)

		case strings.HasPrefix(msg.Method, "notifications/"):
			// Unknown notifications get no response.

		default:
			conn.writeError(ctx, nil, &Error{
				Code:    CodeInvalidRequest,
				Message: "Invalid Request: id missing",
			})
		}

		return
	}

	reqCtx, cancelReq := context.WithCancel(ctx)
	if !conn.track(string(msg.ID), cancelReq) {
		cancelReq()
		conn.writeError(ctx, nil, &Error{
			Code:    CodeInvalidRequest,
			Message: "Invalid Request: id is in use by a request in progress",
		})
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancelReq()
		defer conn.untrack(string(msg.ID))
		s.dispatch(reqCtx, conn, msg)
	}()
}

// connection carries the per-connection state of the stdio transport.
type connection struct {
	server *Server
	w      io.Writer

	writeMu  sync.Mutex
	cancelMu sync.Mutex
	inflight map[string]context.CancelFunc
}

// track registers the cancellation function of an in-flight request. It
// returns false when a request with the same id is in progress.
func (c *connection) track(id string, cancel context.CancelFunc) bool {
	c.cancelMu.Lock()
	defer c.cancelMu.Unlock()

	if _, ok := c.inflight[id]; ok {
		return false
	}
	c.inflight[id] = cancel

	return true
}

// untrack removes the cancellation entry of a finished request.
func (c *connection) untrack(id string) {
	c.cancelMu.Lock()
	delete(c.inflight, id)
	c.cancelMu.Unlock()
}

// cancelRequest cancels the in-flight request referenced by a
// notifications/cancelled message.
func (c *connection) cancelRequest(ctx context.Context, params json.RawMessage) {
	var body struct {
		RequestID json.RawMessage `json:"requestId"`
	}
	if err := json.Unmarshal(params, &body); err != nil {
		c.server.logger.DebugContext(ctx, "Cancelling request: invalid params", log.Err(err))
		return
	}

	c.cancelMu.Lock()
	cancel, ok := c.inflight[string(body.RequestID)]
	c.cancelMu.Unlock()

	if ok {
		cancel()
	}
}

// progressFn builds the progress callback for a request. The callback does
// nothing without a connection or a progress token, and after the request
// context ends.
func (c *connection) progressFn(ctx context.Context, token json.RawMessage) func(float64, float64, string) {
	if c == nil || len(token) == 0 {
		return func(float64, float64, string) {}
	}

	return func(current, total float64, message string) {
		if ctx.Err() != nil {
			return
		}
		c.writeProgress(ctx, token, current, total, message)
	}
}

// writeProgress sends a notifications/progress message for the given token.
func (c *connection) writeProgress(ctx context.Context, token json.RawMessage, current, total float64, message string) {
	params := map[string]any{
		"progressToken": token,
		"progress":      current,
	}
	if total >= 0 {
		params["total"] = total
	}

	if message != "" {
		params["message"] = message
	}

	c.writeJSON(ctx, notificationEnvelope{
		JSONRPC: "2.0",
		Method:  methodProgress,
		Params:  params,
	})
}

// writeResult serializes and writes one result response. When the result
// cannot be encoded, it writes an internal error, so that the client does not
// wait for a response that never comes.
func (c *connection) writeResult(ctx context.Context, id json.RawMessage, result any) {
	data, rpcErr := c.server.encodeResult(ctx, id, result)
	if rpcErr != nil {
		c.writeError(ctx, id, rpcErr)
		return
	}
	c.writeLine(ctx, data)
}

// writeError serializes and writes one error response.
func (c *connection) writeError(ctx context.Context, id json.RawMessage, rpcErr *Error) {
	c.writeJSON(ctx, errorEnvelope{
		JSONRPC: "2.0",
		ID:      id,
		Error:   rpcErr,
	})
}

// writeJSON serializes one message and writes it as a single line.
func (c *connection) writeJSON(ctx context.Context, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		c.server.logger.WarnContext(ctx, "Encoding message", log.Err(err))
		return
	}
	c.writeLine(ctx, data)
}

// writeLine writes one encoded message and its line terminator.
func (c *connection) writeLine(ctx context.Context, data []byte) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if _, err := c.w.Write(append(data, '\n')); err != nil {
		c.server.logger.WarnContext(ctx, "Writing response", log.Err(err))
	}
}

// readLine reads one line from the reader and enforces the message size
// limit. It returns the line without the trailing newline. On EOF with a
// partial line, it returns the partial line together with io.EOF.
func readLine(reader *bufio.Reader, maxSize int) ([]byte, error) {
	var line []byte
	for {
		part, err := reader.ReadSlice('\n')
		if err != nil && !errors.Is(err, bufio.ErrBufferFull) && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("reading message: %w", err)
		}

		length := len(line) + len(part)
		if len(part) > 0 && part[len(part)-1] == '\n' {
			length--
		}

		if length > maxSize {
			if errors.Is(err, bufio.ErrBufferFull) {
				discardLine(reader)
			}

			return nil, ErrMessageTooLarge
		}

		line = append(line, part...)
		switch {
		case err == nil:
			return bytes.TrimSuffix(line, []byte{'\n'}), nil

		case errors.Is(err, io.EOF):
			return line, io.EOF

		default:
			// bufio.ErrBufferFull: read the next slice of the same line.
		}
	}
}

// discardLine skips the rest of the current line. A read error stays for the
// next read to report.
func discardLine(reader *bufio.Reader) {
	for {
		_, err := reader.ReadSlice('\n')
		if !errors.Is(err, bufio.ErrBufferFull) {
			return
		}
	}
}
