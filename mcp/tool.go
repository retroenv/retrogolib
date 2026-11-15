package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/retroenv/retrogolib/log"
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

// AddTool registers a tool with its call handler. When the tool has no input
// schema, AddTool sets a schema that accepts any object. It returns
// jsonschema.ErrInvalidSchema when a schema is not a JSON object, when the input schema
// root does not have "type": "object", or when an x-mcp-header annotation is
// not valid.
func (s *Server) AddTool(t Tool, h ToolHandler) error {
	if t.Name == "" {
		return ErrEmptyName
	}

	headers, err := prepareToolSchemas(&t)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range s.tools {
		if entry.tool.Name == t.Name {
			return ErrDuplicateTool
		}
	}

	s.tools = append(s.tools, toolEntry{
		handler: h,
		headers: headers,
		tool:    t,
	})

	return nil
}

// toolListResult builds the tools/list result in registration order.
func (s *Server) toolListResult() *toolListResult {
	s.mu.RLock()
	tools := make([]Tool, 0, len(s.tools))
	for _, entry := range s.tools {
		tools = append(tools, entry.tool)
	}

	s.mu.RUnlock()

	return &toolListResult{
		ResultType: resultTypeComplete,
		Tools:      tools,
		TTLMs:      s.opts.ListTTL.Milliseconds(),
		CacheScope: cacheScopePublic,
		Meta:       s.serverInfo,
	}
}

// callToolResult executes the named tool and builds its wire result.
func (s *Server) callToolResult(ctx context.Context, params json.RawMessage, meta Meta,
	conn *connection) (any, *Error) {

	var body struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &body); err != nil {
		return nil, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params",
		}
	}
	if body.Name == "" {
		return nil, &Error{
			Code:    CodeInvalidParams,
			Message: "Invalid params: name missing",
		}
	}

	entry, ok := s.findTool(body.Name)
	if !ok {
		return nil, &Error{
			Code:    CodeInvalidParams,
			Message: "Unknown tool: " + body.Name,
		}
	}
	req := ToolRequest{
		Arguments: body.Arguments,
		Meta:      meta,
		Progress:  conn.progressFn(ctx, meta.ProgressToken),
	}
	result := s.executeTool(ctx, entry, req)

	content := result.Content
	if content == nil {
		content = []Content{}
	}

	return &toolCallResult{
		ResultType:        resultTypeComplete,
		Content:           content,
		StructuredContent: result.StructuredContent,
		IsError:           result.IsError,
		Meta:              s.serverInfo,
	}, nil
}

// executeTool runs a tool handler and converts errors and panics into
// isError tool results.
func (s *Server) executeTool(ctx context.Context, entry toolEntry, req ToolRequest) (result *ToolResult) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.WarnContext(ctx, "Tool handler panicked",
				log.String("tool", entry.tool.Name),
				log.Object("panic", r),
				log.String("stack", string(debug.Stack())),
			)
			result = Errorf("internal error: %v", r)
		}
	}()

	res, err := entry.handler(ctx, req)
	if err != nil {
		return Errorf("%v", err)
	}
	if res == nil {
		return &ToolResult{}
	}
	return res
}

// findTool locates a registered tool.
func (s *Server) findTool(name string) (toolEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, entry := range s.tools {
		if entry.tool.Name == name {
			return entry, true
		}
	}

	return toolEntry{}, false
}

// Errorf builds a tool result that reports a tool execution error.
func Errorf(format string, args ...any) *ToolResult {
	return &ToolResult{
		Content: []Content{Text(fmt.Sprintf(format, args...))},
		IsError: true,
	}
}

type toolEntry struct {
	handler ToolHandler
	headers []paramHeader
	tool    Tool
}

// toolListResult is the wire form of a tools/list result.
type toolListResult struct {
	ResultType string          `json:"resultType"`
	Tools      []Tool          `json:"tools"`
	TTLMs      int64           `json:"ttlMs"`
	CacheScope string          `json:"cacheScope"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// toolCallResult is the wire form of a tools/call result.
type toolCallResult struct {
	ResultType        string          `json:"resultType"`
	Content           []Content       `json:"content"`
	StructuredContent any             `json:"structuredContent,omitempty"`
	IsError           bool            `json:"isError"`
	Meta              json.RawMessage `json:"_meta,omitempty"`
}
