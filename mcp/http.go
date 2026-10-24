package mcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/retroenv/retrogolib/log"
)

// Streamable HTTP header names.
const (
	headerMethod          = "Mcp-Method"
	headerName            = "Mcp-Name"
	headerProtocolVersion = "MCP-Protocol-Version"
)

// Base64 sentinel markers for encoded header values.
const (
	base64HeaderPrefix = "=?base64?"
	base64HeaderSuffix = "?="
)

// Handler returns an http.Handler that serves the streamable HTTP transport.
// The handler accepts POST only, validates the required headers against the
// body, and answers with a single JSON object. It accepts every notification
// with status 202 and ignores it. The client cancels a request when it closes
// the connection. Progress notifications are not sent over HTTP, so
// ToolRequest.Progress does nothing on this transport.
//
//nolint:ireturn // Required by net/http.
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(s.serveHTTP)
}

// serveHTTP processes one streamable HTTP request.
func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.originAllowed(r.Header.Get("Origin")) {
		http.Error(w, "forbidden origin", http.StatusForbidden)
		return
	}

	ctx := r.Context()
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, int64(s.opts.MaxMessageSize)))
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			s.writeHTTPError(ctx, w, http.StatusRequestEntityTooLarge, nil, &Error{
				Code:    CodeInvalidRequest,
				Message: "message too large",
			})
			return
		}

		s.writeHTTPError(ctx, w, http.StatusBadRequest, nil, &Error{
			Code:    CodeInvalidRequest,
			Message: "Invalid Request: reading the body failed",
		})

		return
	}

	msg, rpcErr := parseMessage(body)
	if rpcErr != nil {
		s.writeHTTPError(ctx, w, http.StatusBadRequest, nil, rpcErr)
		return
	}

	if !msg.hasID() {
		// The protocol defines no client notifications on this transport.
		// The client closes the connection to cancel a request.
		w.WriteHeader(http.StatusAccepted)
		return
	}

	s.serveHTTPRequest(ctx, w, r.Header, msg)
}

// serveHTTPRequest processes a request POST and writes its response.
func (s *Server) serveHTTPRequest(ctx context.Context, w http.ResponseWriter, header http.Header, msg message) {
	meta, rpcErr := parseRequestMeta(msg)
	if rpcErr != nil {
		s.writeHTTPError(ctx, w, http.StatusBadRequest, msg.ID, rpcErr)
		return
	}

	if rpcErr := validateHTTPHeaders(header, msg, meta); rpcErr != nil {
		s.writeHTTPError(ctx, w, http.StatusBadRequest, msg.ID, rpcErr)
		return
	}
	if rpcErr := s.validateToolParamHeaders(header, msg); rpcErr != nil {
		s.writeHTTPError(ctx, w, http.StatusBadRequest, msg.ID, rpcErr)
		return
	}

	result, rpcErr := s.call(ctx, msg.Method, msg.Params, meta, nil)
	if rpcErr != nil {
		s.writeHTTPError(ctx, w, httpStatusForError(rpcErr), msg.ID, rpcErr)
		return
	}

	data, rpcErr := s.encodeResult(ctx, msg.ID, result)
	if rpcErr != nil {
		s.writeHTTPError(ctx, w, httpStatusForError(rpcErr), msg.ID, rpcErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(append(data, '\n')); err != nil {
		s.logger.WarnContext(ctx, "Writing response", log.Err(err))
	}
}

// validateToolParamHeaders checks the Mcp-Param-{name} headers of a
// tools/call request against the x-mcp-header annotations of the tool.
func (s *Server) validateToolParamHeaders(h http.Header, msg message) *Error {
	if msg.Method != methodToolsCall {
		return nil
	}

	var body struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	_ = json.Unmarshal(msg.Params, &body) // the call reports invalid params

	entry, ok := s.findTool(body.Name)
	if !ok {
		return nil
	}
	return validateParamHeaders(h, entry.headers, body.Arguments)
}

// writeHTTPError writes a JSON-RPC error with the matching HTTP status.
func (s *Server) writeHTTPError(ctx context.Context, w http.ResponseWriter, status int, id json.RawMessage,
	rpcErr *Error) {

	data, err := json.Marshal(errorEnvelope{
		JSONRPC: "2.0",
		ID:      id,
		Error:   rpcErr,
	})
	if err != nil {
		s.logger.WarnContext(ctx, "Encoding error response", log.Err(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(append(data, '\n')); err != nil {
		s.logger.WarnContext(ctx, "Writing error response", log.Err(err))
	}
}

// originAllowed checks the Origin header against the configured allowlist.
// Requests without an Origin header are allowed. localhost and 127.0.0.1 on
// any port are always allowed.
func (s *Server) originAllowed(origin string) bool {
	if origin == "" {
		return true
	}
	for _, allowed := range s.opts.AllowedOrigins {
		if allowed == origin {
			return true
		}
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	switch u.Scheme {
	case "http", "https":
	default:
		return false
	}

	switch u.Hostname() {
	case "localhost", "127.0.0.1", "::1":
		return true

	default:
		return false
	}
}

// httpStatusForError maps a JSON-RPC error code to an HTTP status.
func httpStatusForError(rpcErr *Error) int {
	switch rpcErr.Code {
	case CodeMethodNotFound:
		return http.StatusNotFound

	case CodeInternal:
		return http.StatusInternalServerError

	default:
		return http.StatusBadRequest
	}
}

// validateHTTPHeaders checks the required streamable HTTP headers against the
// request body.
func validateHTTPHeaders(h http.Header, msg message, meta Meta) *Error {
	version := h.Get(headerProtocolVersion)
	if version == "" || !validHeaderValue(version) {
		return headerMismatchError("missing or malformed " + headerProtocolVersion + " header")
	}
	if version != meta.ProtocolVersion {
		return headerMismatchError(headerProtocolVersion + " header does not match the request body")
	}

	method := h.Get(headerMethod)
	if method == "" || !validHeaderValue(method) {
		return headerMismatchError("missing or malformed " + headerMethod + " header")
	}
	if method != msg.Method {
		return headerMismatchError(headerMethod + " header does not match the request body")
	}

	nameValue, required, rpcErr := headerNameValue(msg)
	if rpcErr != nil {
		return rpcErr
	}
	if !required {
		return nil
	}

	header := h.Get(headerName)
	if header == "" || !validHeaderValue(header) {
		return headerMismatchError("missing or malformed " + headerName + " header")
	}
	if decodeHeaderValue(header) != nameValue {
		return headerMismatchError(headerName + " header does not match the request body")
	}
	return nil
}

// headerNameValue returns the body value the Mcp-Name header must carry for
// the request method.
func headerNameValue(msg message) (string, bool, *Error) {
	switch msg.Method {
	case methodToolsCall, methodPromptsGet:
		var body struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(msg.Params, &body); err != nil || body.Name == "" {
			return "", false, &Error{
				Code:    CodeInvalidParams,
				Message: "Invalid params: name missing",
			}
		}
		return body.Name, true, nil

	case methodResourcesRead:
		var body struct {
			URI string `json:"uri"`
		}
		if err := json.Unmarshal(msg.Params, &body); err != nil || body.URI == "" {
			return "", false, &Error{
				Code:    CodeInvalidParams,
				Message: "Invalid params: uri missing",
			}
		}
		return body.URI, true, nil

	default:
		return "", false, nil
	}
}

// validHeaderValue checks that a header value uses allowed ASCII characters.
func validHeaderValue(v string) bool {
	for _, r := range v {
		if r == '\t' || r == ' ' || (r >= 0x21 && r <= 0x7E) {
			continue
		}
		return false
	}

	return true
}

// decodeHeaderValue decodes a base64 sentinel header value.
func decodeHeaderValue(v string) string {
	if !strings.HasPrefix(v, base64HeaderPrefix) || !strings.HasSuffix(v, base64HeaderSuffix) {
		return v
	}

	encoded := strings.TrimSuffix(strings.TrimPrefix(v, base64HeaderPrefix), base64HeaderSuffix)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return v
	}
	return string(decoded)
}
