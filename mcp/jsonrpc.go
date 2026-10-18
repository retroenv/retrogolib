package mcp

import (
	"bytes"
	"encoding/json"
)

// message is the common wire envelope for JSON-RPC requests and notifications.
type message struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      json.RawMessage `json:"id"`
}

// resultEnvelope is a JSON-RPC result response.
type resultEnvelope struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result"`
}

// errorEnvelope is a JSON-RPC error response.
type errorEnvelope struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Error   *Error          `json:"error"`
}

// notificationEnvelope is a JSON-RPC notification sent by the server.
type notificationEnvelope struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

// hasID reports whether the message carries a request id.
func (m message) hasID() bool {
	return len(m.ID) > 0
}

// parseMessage decodes one JSON-RPC message and validates the wire shape.
// A nil error marks a valid request or notification.
func parseMessage(data []byte) (message, *Error) {
	var msg message
	if trimmed := bytes.TrimSpace(data); len(trimmed) > 0 && trimmed[0] == '[' {
		return msg, &Error{
			Code:    CodeInvalidRequest,
			Message: "Invalid Request: batch messages are not supported",
		}
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return msg, &Error{
			Code:    CodeParse,
			Message: "Parse error",
		}
	}
	if msg.JSONRPC != "" && msg.JSONRPC != "2.0" {
		return msg, &Error{
			Code:    CodeInvalidRequest,
			Message: "Invalid Request: unsupported JSON-RPC version",
		}
	}
	if msg.Method == "" {
		return msg, &Error{
			Code:    CodeInvalidRequest,
			Message: "Invalid Request: method missing",
		}
	}
	if msg.hasID() && !validID(msg.ID) {
		return msg, &Error{
			Code:    CodeInvalidRequest,
			Message: "Invalid Request: id must be a string or an integer",
		}
	}
	return msg, nil
}

// validID reports whether a request id is a JSON string or number, not null.
func validID(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return false
	}

	dec := json.NewDecoder(bytes.NewReader(trimmed))
	dec.UseNumber()

	var value any
	if err := dec.Decode(&value); err != nil {
		return false
	}

	switch value.(type) {
	case string, json.Number:
		return true

	default:
		return false
	}
}
