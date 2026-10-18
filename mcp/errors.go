package mcp

import "errors"

// Sentinel errors returned by registration and transport functions.
var (
	ErrDuplicateTool      = errors.New("duplicate tool")
	ErrDuplicateResource  = errors.New("duplicate resource")
	ErrDuplicatePrompt    = errors.New("duplicate prompt")
	ErrEmptyName          = errors.New("empty name")
	ErrInvalidSchema      = errors.New("invalid schema")
	ErrMessageTooLarge    = errors.New("message too large")
	ErrUnsupportedVersion = errors.New("unsupported protocol version")
)

// ErrResourceNotFound is the error that a ResourceHandler returns when no
// resource exists for the URI. The server then sends the protocol error for a
// missing resource instead of an internal error.
var ErrResourceNotFound = errors.New("resource not found")

// JSON-RPC 2.0 and MCP error codes.
const (
	CodeParse                      = -32700
	CodeInvalidRequest             = -32600
	CodeMethodNotFound             = -32601
	CodeInvalidParams              = -32602
	CodeInternal                   = -32603
	CodeHeaderMismatch             = -32020
	CodeMissingClientCapability    = -32021
	CodeUnsupportedProtocolVersion = -32022
)

// Error is a JSON-RPC error response.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Message
}

// Unwrap maps protocol errors to their sentinels for errors.Is checks.
func (e *Error) Unwrap() error {
	if e.Code == CodeUnsupportedProtocolVersion {
		return ErrUnsupportedVersion
	}
	return nil
}
