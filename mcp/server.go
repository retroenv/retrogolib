package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/retroenv/retrogolib/log"
)

// Method names implemented by the server.
const (
	methodCancelled              = "notifications/cancelled"
	methodDiscover               = "server/discover"
	methodInitialize             = "initialize"
	methodProgress               = "notifications/progress"
	methodPromptsGet             = "prompts/get"
	methodPromptsList            = "prompts/list"
	methodResourcesList          = "resources/list"
	methodResourcesRead          = "resources/read"
	methodResourcesTemplatesList = "resources/templates/list"
	methodToolsCall              = "tools/call"
	methodToolsList              = "tools/list"
)

// Wire and default values.
const (
	cacheScopePublic      = "public"
	defaultMaxMessageSize = 16 << 20 // 16 MiB
	resultTypeComplete    = "complete"
)

// defaultListTTL is the ttlMs value for list and read results when the
// Options value is zero.
const defaultListTTL = time.Hour

// Options configures the server.
type Options struct {
	// Name is the server name reported in serverInfo and discover results.
	Name string

	// Version is the server version reported in serverInfo and discover results.
	Version string

	// Instructions is optional natural-language guidance for agents.
	Instructions string

	// MaxMessageSize limits one JSON-RPC message in bytes. Zero means 16 MiB.
	MaxMessageSize int

	// ListTTL is the ttlMs value attached to list and read results. Zero means one hour.
	ListTTL time.Duration

	// AllowedOrigins lists additional allowed Origin values for the HTTP transport.
	// localhost and 127.0.0.1 on any port are always allowed.
	AllowedOrigins []string
}

// New creates a server with the given options. The logger is nil-safe.
func New(logger *log.Logger, opts Options) *Server {
	if opts.MaxMessageSize <= 0 {
		opts.MaxMessageSize = defaultMaxMessageSize
	}
	if opts.ListTTL <= 0 {
		opts.ListTTL = defaultListTTL
	}

	var serverInfo json.RawMessage
	if opts.Name != "" {
		serverInfo, _ = json.Marshal(map[string]any{
			metaKeyServerInfo: Implementation{
				Name:    opts.Name,
				Version: opts.Version,
			},
		})
	}

	return &Server{
		logger:     logger,
		opts:       opts,
		serverInfo: serverInfo,
	}
}

type toolEntry struct {
	tool    Tool
	handler ToolHandler
	headers []paramHeader
}

type resourceEntry struct {
	resource Resource
	handler  ResourceHandler
}

type templateEntry struct {
	template ResourceTemplate
	handler  ResourceHandler
}

type promptEntry struct {
	prompt  Prompt
	handler PromptHandler
}

// Server holds the registered features and answers MCP requests.
type Server struct {
	logger     *log.Logger
	opts       Options
	serverInfo json.RawMessage

	mu        sync.RWMutex
	tools     []toolEntry
	resources []resourceEntry
	templates []templateEntry
	prompts   []promptEntry
}

// AddTool registers a tool with its call handler. When the tool has no input
// schema, AddTool sets a schema that accepts any object. It returns
// ErrInvalidSchema when a schema is not a JSON object, when the input schema
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
		tool:    t,
		handler: h,
		headers: headers,
	})
	return nil
}

// AddResource registers a resource with its read handler.
func (s *Server) AddResource(r Resource, h ResourceHandler) error {
	if r.URI == "" {
		return ErrEmptyName
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range s.resources {
		if entry.resource.URI == r.URI {
			return ErrDuplicateResource
		}
	}
	s.resources = append(s.resources, resourceEntry{
		resource: r,
		handler:  h,
	})
	return nil
}

// AddResourceTemplate registers a URI template with its read handler.
// Templates support {name} segments only (RFC 6570 level 1).
func (s *Server) AddResourceTemplate(t ResourceTemplate, h ResourceHandler) error {
	if t.URITemplate == "" {
		return ErrEmptyName
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range s.templates {
		if entry.template.URITemplate == t.URITemplate {
			return ErrDuplicateResource
		}
	}
	s.templates = append(s.templates, templateEntry{
		template: t,
		handler:  h,
	})
	return nil
}

// AddPrompt registers a prompt with its render handler.
func (s *Server) AddPrompt(p Prompt, h PromptHandler) error {
	if p.Name == "" {
		return ErrEmptyName
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range s.prompts {
		if entry.prompt.Name == p.Name {
			return ErrDuplicatePrompt
		}
	}
	s.prompts = append(s.prompts, promptEntry{
		prompt:  p,
		handler: h,
	})
	return nil
}

// dispatch validates _meta and routes one request to its handler.
func (s *Server) dispatch(ctx context.Context, conn *connection, msg message) {
	meta, rpcErr := parseRequestMeta(msg)
	if rpcErr != nil {
		conn.writeError(ctx, msg.ID, rpcErr)
		return
	}

	result, rpcErr := s.call(ctx, msg.Method, msg.Params, meta, conn)
	if ctx.Err() != nil {
		return // the request was cancelled: send no further messages
	}
	if rpcErr != nil {
		conn.writeError(ctx, msg.ID, rpcErr)
		return
	}
	conn.writeResult(ctx, msg.ID, result)
}

// encodeResult encodes a result response. When the result cannot be encoded,
// it records a warning and returns an internal error.
func (s *Server) encodeResult(ctx context.Context, id json.RawMessage, result any) ([]byte, *Error) {
	data, err := json.Marshal(resultEnvelope{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
	if err != nil {
		s.logger.WarnContext(ctx, "Encoding result", log.Err(err))
		return nil, &Error{
			Code:    CodeInternal,
			Message: "Internal error: the result cannot be encoded",
		}
	}
	return data, nil
}

// call routes a request to the matching feature handler. A panic in a
// resource or prompt handler becomes an internal error.
func (s *Server) call(ctx context.Context, method string, params json.RawMessage, meta Meta,
	conn *connection) (result any, rpcErr *Error) {

	defer func() {
		if r := recover(); r != nil {
			s.logger.WarnContext(ctx, "Request handler panicked",
				log.String("method", method),
				log.Object("panic", r),
				log.String("stack", string(debug.Stack())),
			)
			result = nil
			rpcErr = &Error{
				Code:    CodeInternal,
				Message: "Internal error",
				Data:    map[string]any{"error": fmt.Sprint(r)},
			}
		}
	}()

	switch method {
	case methodDiscover:
		return s.discoverResult(), nil

	case methodToolsList:
		return s.toolListResult(), nil

	case methodToolsCall:
		return s.callToolResult(ctx, params, meta, conn)

	case methodPromptsList:
		return s.promptListResult(), nil

	case methodPromptsGet:
		return s.getPromptResult(ctx, params)

	case methodResourcesList:
		return s.resourceListResult(), nil

	case methodResourcesTemplatesList:
		return s.templateListResult(), nil

	case methodResourcesRead:
		return s.readResourceResult(ctx, params)

	default:
		return nil, &Error{
			Code:    CodeMethodNotFound,
			Message: "Method not found: " + method,
		}
	}
}

// discoverResult builds the server/discover result from the registered features.
func (s *Server) discoverResult() *discoverResult {
	s.mu.RLock()
	caps := capabilities{
		Tools:     present(len(s.tools) > 0),
		Resources: present(len(s.resources) > 0 || len(s.templates) > 0),
		Prompts:   present(len(s.prompts) > 0),
	}
	s.mu.RUnlock()

	return &discoverResult{
		ResultType:        resultTypeComplete,
		SupportedVersions: []string{ProtocolVersion},
		Capabilities:      caps,
		Instructions:      s.opts.Instructions,
		TTLMs:             s.opts.ListTTL.Milliseconds(),
		CacheScope:        cacheScopePublic,
		Meta:              s.serverInfo,
	}
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

// findResource locates the handler for a resource URI, checking static
// resources first and then templates in registration order.
func (s *Server) findResource(uri string) (ResourceHandler, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, entry := range s.resources {
		if entry.resource.URI == uri {
			return entry.handler, true
		}
	}
	for _, entry := range s.templates {
		if matchURITemplate(entry.template.URITemplate, uri) {
			return entry.handler, true
		}
	}
	return nil, false
}

// findPrompt locates a registered prompt.
func (s *Server) findPrompt(name string) (promptEntry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, entry := range s.prompts {
		if entry.prompt.Name == name {
			return entry, true
		}
	}
	return promptEntry{}, false
}

// capabilities lists the registered feature groups.
type capabilities struct {
	Tools     *struct{} `json:"tools,omitempty"`
	Resources *struct{} `json:"resources,omitempty"`
	Prompts   *struct{} `json:"prompts,omitempty"`
}

// discoverResult is the wire form of a server/discover result.
type discoverResult struct {
	ResultType        string          `json:"resultType"`
	SupportedVersions []string        `json:"supportedVersions"`
	Capabilities      capabilities    `json:"capabilities"`
	Instructions      string          `json:"instructions,omitempty"`
	TTLMs             int64           `json:"ttlMs"`
	CacheScope        string          `json:"cacheScope"`
	Meta              json.RawMessage `json:"_meta,omitempty"`
}

// toolListResult is the wire form of a tools/list result.
type toolListResult struct {
	ResultType string          `json:"resultType"`
	Tools      []Tool          `json:"tools"`
	TTLMs      int64           `json:"ttlMs"`
	CacheScope string          `json:"cacheScope"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// present returns an empty object when the feature group is registered.
func present(registered bool) *struct{} {
	if !registered {
		return nil
	}
	return &struct{}{}
}
