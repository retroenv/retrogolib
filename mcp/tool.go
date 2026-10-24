package mcp

import (
	"context"
	"encoding/json"
	"fmt"
)

// Tool describes a tool the server exposes to agents. InputSchema must be a
// JSON Schema object with "type": "object" at the root. When InputSchema is
// empty, AddTool sets a schema that accepts any object.
type Tool struct {
	Name         string           `json:"name"`
	Title        string           `json:"title,omitempty"`
	Description  string           `json:"description,omitempty"`
	InputSchema  json.RawMessage  `json:"inputSchema"`
	OutputSchema json.RawMessage  `json:"outputSchema,omitempty"`
	Annotations  *ToolAnnotations `json:"annotations,omitempty"`
}

// ToolAnnotations carries optional hints about tool behavior.
type ToolAnnotations struct {
	Title           string `json:"title,omitempty"`
	ReadOnlyHint    bool   `json:"readOnlyHint,omitempty"`
	IdempotentHint  bool   `json:"idempotentHint,omitempty"`
	DestructiveHint *bool  `json:"destructiveHint,omitempty"`
	OpenWorldHint   *bool  `json:"openWorldHint,omitempty"`
}

// ToolHandler executes one tool call.
type ToolHandler func(ctx context.Context, req ToolRequest) (*ToolResult, error)

// ToolRequest carries the parsed arguments and metadata of a tools/call request.
type ToolRequest struct {
	Arguments json.RawMessage
	Meta      Meta

	// Progress sends a notifications/progress message for the request. It is
	// never nil. It does nothing when the request has no progress token, on
	// the HTTP transport, and after the request ends. A negative total means
	// that the total is not known.
	Progress func(current, total float64, message string)
}

// ToolResult is the result of a tool execution.
type ToolResult struct {
	Content           []Content `json:"content"`
	StructuredContent any       `json:"structuredContent,omitempty"`
	IsError           bool      `json:"isError"`
}

// Content types.
const (
	contentTypeAudio    = "audio"
	contentTypeImage    = "image"
	contentTypeResource = "resource"
	contentTypeText     = "text"
)

// Content is one content item in a tool result or prompt message.
type Content struct {
	Type     string            `json:"type"`
	Text     string            `json:"text,omitempty"`
	Data     []byte            `json:"data,omitempty"`
	MimeType string            `json:"mimeType,omitempty"`
	Resource *ResourceContents `json:"resource,omitempty"`
}

// MarshalJSON encodes the content item with the fields that its type
// requires. An empty text or empty data stays in the output, because the
// protocol schema requires these fields.
func (c Content) MarshalJSON() ([]byte, error) {
	var v any

	switch c.Type {
	case contentTypeText:
		v = struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			Type: c.Type,
			Text: c.Text,
		}

	case contentTypeImage, contentTypeAudio:
		data := c.Data
		if data == nil {
			data = []byte{}
		}

		v = struct {
			Type     string `json:"type"`
			Data     []byte `json:"data"`
			MimeType string `json:"mimeType"`
		}{
			Type:     c.Type,
			Data:     data,
			MimeType: c.MimeType,
		}

	case contentTypeResource:
		v = struct {
			Type     string            `json:"type"`
			Resource *ResourceContents `json:"resource"`
		}{
			Type:     c.Type,
			Resource: c.Resource,
		}

	default:
		type plain Content
		v = plain(c)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encoding content: %w", err)
	}
	return data, nil
}

// Text builds a text content item.
func Text(s string) Content {
	return Content{
		Type: contentTypeText,
		Text: s,
	}
}

// Image builds an image content item from raw image bytes.
func Image(data []byte, mimeType string) Content {
	return Content{
		Type:     contentTypeImage,
		Data:     data,
		MimeType: mimeType,
	}
}

// Audio builds an audio content item from raw audio bytes.
func Audio(data []byte, mimeType string) Content {
	return Content{
		Type:     contentTypeAudio,
		Data:     data,
		MimeType: mimeType,
	}
}

// EmbeddedText builds an embedded text resource content item.
func EmbeddedText(uri, mimeType, text string) Content {
	return Content{
		Type: contentTypeResource,
		Resource: &ResourceContents{
			URI:      uri,
			MimeType: mimeType,
			Text:     text,
		},
	}
}

// Errorf builds a tool result that reports a tool execution error.
func Errorf(format string, args ...any) *ToolResult {
	return &ToolResult{
		Content: []Content{Text(fmt.Sprintf(format, args...))},
		IsError: true,
	}
}

// toolCallResult is the wire form of a tools/call result.
type toolCallResult struct {
	ResultType        string          `json:"resultType"`
	Content           []Content       `json:"content"`
	StructuredContent any             `json:"structuredContent,omitempty"`
	IsError           bool            `json:"isError"`
	Meta              json.RawMessage `json:"_meta,omitempty"`
}
