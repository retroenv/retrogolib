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
	cacheScopePublic             = "public"
	defaultMaxMessageSize        = 16 << 20 // 16 MiB
	defaultMaxConcurrentRequests = 64
	resultTypeComplete           = "complete"
)

// defaultListTTL is the ttlMs value for list and read results when the
// Options value is zero.
const defaultListTTL = time.Hour

// Options configures the server.
type Options struct {
	// CancelOnEOF cancels in-flight requests when the stdio input closes.
	// The default waits for those requests to finish.
	CancelOnEOF bool

	// Name is the server name reported in serverInfo and discover results.
	Name string

	// Version is the server version reported in serverInfo and discover results.
	Version string

	// Instructions is optional natural-language guidance for agents.
	Instructions string

	// MaxMessageSize limits one JSON-RPC message in bytes. Zero means 16 MiB.
	MaxMessageSize int

	// MaxConcurrentRequests limits active requests across both transports.
	// Values less than 1 use 64.
	MaxConcurrentRequests int

	// ListTTL is the ttlMs value attached to list and read results. Zero means one hour.
	ListTTL time.Duration

	// AllowedOrigins lists additional allowed Origin values for the HTTP transport.
	// localhost and 127.0.0.1 on any port are always allowed.
	AllowedOrigins []string
}

// Server holds the registered features and answers MCP requests.
type Server struct {
	logger     *log.Logger
	opts       Options
	serverInfo json.RawMessage

	requests chan struct{}

	mu        sync.RWMutex
	prompts   []promptEntry
	resources []resourceEntry
	templates []templateEntry
	tools     []toolEntry
}

// New creates a server with the given options. The logger is nil-safe.
func New(logger *log.Logger, opts Options) *Server {
	if opts.MaxMessageSize <= 0 {
		opts.MaxMessageSize = defaultMaxMessageSize
	}

	if opts.MaxConcurrentRequests <= 0 {
		opts.MaxConcurrentRequests = defaultMaxConcurrentRequests
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

		requests: make(chan struct{}, opts.MaxConcurrentRequests),
	}
}

func (s *Server) tryAcquireRequest() bool {
	select {
	case s.requests <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Server) releaseRequest() {
	<-s.requests
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

// present returns an empty object when the feature group is registered.
func present(registered bool) *struct{} {
	if !registered {
		return nil
	}
	return &struct{}{}
}
