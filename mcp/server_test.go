package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/log"
)

const (
	serverName    = "test-server"
	serverVersion = "1.2.3"
)

// testMeta is the common protocol metadata fragment of a request.
const testMeta = `"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{}`

func newTestServer(t *testing.T) *Server {
	t.Helper()
	return New(log.NewTestLogger(t), Options{
		Name:    serverName,
		Version: serverVersion,
	})
}

// runServe feeds the input lines to the server and returns the response lines.
func runServe(t *testing.T, s *Server, input string) []string {
	t.Helper()

	var out bytes.Buffer
	err := s.Serve(t.Context(), strings.NewReader(input), &out)
	assert.NoError(t, err)

	text := strings.TrimSpace(out.String())
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

// request builds one JSON-RPC request line with the required metadata.
func request(id, method, params string) string {
	return `{"jsonrpc":"2.0","id":` + id + `,"method":"` + method + `","params":{` + params + `}}`
}

// metaParams merges extra params with the required metadata.
func metaParams(extra string) string {
	if extra == "" {
		return `"_meta":{` + testMeta + `}`
	}
	return extra + `,"_meta":{` + testMeta + `}`
}

func TestDiscover(t *testing.T) {
	s := newTestServer(t)

	input := request(`"discover-1"`, methodDiscover, metaParams(""))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":"discover-1","result":{"resultType":"complete",` +
		`"supportedVersions":["2026-07-28"],"capabilities":{},"ttlMs":3600000,` +
		`"cacheScope":"public","_meta":{"io.modelcontextprotocol/serverInfo":` +
		`{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestDiscoverCapabilities(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{Name: "t"}, echoHandler))
	assert.NoError(t, s.AddResource(Resource{
		URI:  "mem://r",
		Name: "r",
	}, readHandler))
	assert.NoError(t, s.AddPrompt(Prompt{Name: "p"}, promptHandler))

	input := request(`1`, methodDiscover, metaParams(""))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"tools":{},"resources":{},"prompts":{}`)
}

func TestDiscoverMissingMeta(t *testing.T) {
	s := newTestServer(t)

	input := `{"jsonrpc":"2.0","id":1,"method":"server/discover","params":{}}`
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"Invalid params: _meta missing"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestUnsupportedVersion(t *testing.T) {
	s := newTestServer(t)

	params := `"_meta":{"io.modelcontextprotocol/protocolVersion":"2025-11-25","io.modelcontextprotocol/clientCapabilities":{}}`
	input := request(`1`, methodDiscover, params)
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32022,"message":"Unsupported protocol version",` +
		`"data":{"supported":["2026-07-28"],"requested":"2025-11-25"}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestUnknownMethod(t *testing.T) {
	s := newTestServer(t)

	input := request(`1`, "unknown/method", metaParams(""))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"Method not found: unknown/method"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestLegacyInitialize(t *testing.T) {
	s := newTestServer(t)

	input := request(`0`, methodInitialize, `"protocolVersion":"2025-06-18","capabilities":{}`)
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":0,"error":{"code":-32022,"message":"Unsupported protocol version: ` +
		`this server supports only MCP 2026-07-28 and does not implement the initialize handshake",` +
		`"data":{"supported":["2026-07-28"],"requested":"2025-06-18"}}}`
	assert.Equal(t, []string{want}, lines)
}

func echoHandler(_ context.Context, req ToolRequest) (*ToolResult, error) {
	return &ToolResult{Content: []Content{Text(string(req.Arguments))}}, nil
}

func readHandler(_ context.Context, uri string) ([]ResourceContents, error) {
	return []ResourceContents{{URI: uri, MimeType: "text/plain", Text: "data"}}, nil
}

func promptHandler(_ context.Context, _ map[string]string) (*PromptResult, error) {
	return &PromptResult{Messages: []PromptMessage{{Role: "user", Content: Text("hi")}}}, nil
}

func TestToolCallSuccess(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{
		Name:        "echo",
		Description: "Echo arguments.",
	}, echoHandler))

	input := request(`1`, methodToolsCall, metaParams(`"name":"echo","arguments":{"x":1}`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"result":{"resultType":"complete","content":[{"type":"text",` +
		`"text":"{\"x\":1}"}],"isError":false,"_meta":{"io.modelcontextprotocol/serverInfo":` +
		`{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestToolIsError(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{Name: "fail"}, func(_ context.Context, _ ToolRequest) (*ToolResult, error) {
		return Errorf("something went wrong: %d", 42), nil
	}))

	input := request(`1`, methodToolsCall, metaParams(`"name":"fail"`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"result":{"resultType":"complete","content":[{"type":"text",` +
		`"text":"something went wrong: 42"}],"isError":true,"_meta":{"io.modelcontextprotocol/serverInfo":` +
		`{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestToolHandlerError(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{Name: "boom"}, func(_ context.Context, _ ToolRequest) (*ToolResult, error) {
		return nil, errors.New("boom")
	}))

	input := request(`1`, methodToolsCall, metaParams(`"name":"boom"`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"result":{"resultType":"complete","content":[{"type":"text",` +
		`"text":"boom"}],"isError":true,"_meta":{"io.modelcontextprotocol/serverInfo":` +
		`{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestToolHandlerPanic(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{Name: "panic"}, func(_ context.Context, _ ToolRequest) (*ToolResult, error) {
		panic("kaboom")
	}))

	input := request(`1`, methodToolsCall, metaParams(`"name":"panic"`))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"isError":true`)
	assert.Contains(t, lines[0], `"internal error: kaboom"`)
}

func TestToolUnknownTool(t *testing.T) {
	s := newTestServer(t)

	input := request(`1`, methodToolsCall, metaParams(`"name":"nope"`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"Unknown tool: nope"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestToolMissingName(t *testing.T) {
	s := newTestServer(t)

	input := request(`1`, methodToolsCall, metaParams(""))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"Invalid params: name missing"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestCancellation(t *testing.T) {
	s := newTestServer(t)

	cancelled := make(chan struct{})
	assert.NoError(t, s.AddTool(Tool{Name: "block"}, func(ctx context.Context, _ ToolRequest) (*ToolResult, error) {
		<-ctx.Done()
		close(cancelled)
		return nil, ctx.Err()
	}))

	input := request(`1`, methodToolsCall, metaParams(`"name":"block"`)) + "\n" +
		`{"jsonrpc":"2.0","method":"notifications/cancelled","params":{"requestId":1}}` + "\n"
	lines := runServe(t, s, input)

	assert.Empty(t, lines)
	select {
	case <-cancelled:
	default:
		t.Fatal("handler did not observe cancellation")
	}
}

func TestOversizedLine(t *testing.T) {
	s := New(log.NewTestLogger(t), Options{
		Name:           serverName,
		MaxMessageSize: 256,
	})

	oversized := request(`1`, methodDiscover, metaParams("")) + strings.Repeat(" ", 128<<10) + "\n"
	next := `{"jsonrpc":"2.0","id":2,"method":"unknown/method","params":{` + metaParams("") + `}}` + "\n"
	lines := runServe(t, s, oversized+next)

	assert.Equal(t, []string{
		`{"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"message too large"}}`,
		`{"jsonrpc":"2.0","id":2,"error":{"code":-32601,"message":"Method not found: unknown/method"}}`,
	}, lines)
}

func TestDuplicateInFlightID(t *testing.T) {
	s := newTestServer(t)

	release := make(chan struct{})
	assert.NoError(t, s.AddTool(Tool{Name: "block"}, func(_ context.Context, _ ToolRequest) (*ToolResult, error) {
		<-release
		return &ToolResult{Content: []Content{Text("done")}}, nil
	}))

	reader, writer := io.Pipe()
	out := &lineWriter{lines: make(chan string, 2)}
	done := make(chan error, 1)
	go func() {
		done <- s.Serve(t.Context(), reader, out)
	}()

	line := request(`1`, methodToolsCall, metaParams(`"name":"block"`)) + "\n"
	_, err := io.WriteString(writer, line+line)
	assert.NoError(t, err)

	want := `{"jsonrpc":"2.0","id":null,"error":{"code":-32600,` +
		`"message":"Invalid Request: id is in use by a request in progress"}}`
	assert.Equal(t, want, <-out.lines)

	close(release)
	assert.Contains(t, <-out.lines, `"id":1,"result"`)
	assert.NoError(t, writer.Close())
	assert.NoError(t, <-done)
}

// lineWriter sends each written message line to a channel.
type lineWriter struct {
	lines chan string
}

func (l *lineWriter) Write(p []byte) (int, error) {
	l.lines <- strings.TrimSpace(string(p))
	return len(p), nil
}

func TestResultEncodingFailure(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{Name: "nan"}, func(_ context.Context, _ ToolRequest) (*ToolResult, error) {
		return &ToolResult{StructuredContent: math.NaN()}, nil
	}))

	input := request(`1`, methodToolsCall, metaParams(`"name":"nan"`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"Internal error: the result cannot be encoded"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestToolDefaultInputSchema(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{Name: "plain"}, echoHandler))

	input := request(`1`, methodToolsList, metaParams(""))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `{"name":"plain","inputSchema":{"type":"object"}}`)
}

func TestAddToolInvalidSchema(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{name: "not an object", input: `[1]`},
		{name: "invalid json", input: `{`},
		{name: "root type not object", input: `{"type":"string"}`},
		{name: "output schema not an object", input: `{"type":"object"}`, output: `true`},
		{name: "header on array item", input: `{"type":"object","properties":{"list":{"type":"array",` +
			`"items":{"type":"object","properties":{"id":{"type":"string","x-mcp-header":"Id"}}}}}}`},
		{name: "header in composition", input: `{"type":"object","anyOf":[{"properties":` +
			`{"id":{"type":"string","x-mcp-header":"Id"}}}]}`},
		{name: "header on number", input: `{"type":"object","properties":{"ratio":{"type":"number",` +
			`"x-mcp-header":"Ratio"}}}`},
		{name: "header name not a token", input: `{"type":"object","properties":{"id":{"type":"string",` +
			`"x-mcp-header":"Bad Name"}}}`},
		{name: "header name duplicate", input: `{"type":"object","properties":{"a":{"type":"string",` +
			`"x-mcp-header":"Region"},"b":{"type":"string","x-mcp-header":"region"}}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestServer(t)
			err := s.AddTool(Tool{
				Name:         "bad",
				InputSchema:  []byte(tt.input),
				OutputSchema: []byte(tt.output),
			}, echoHandler)
			assert.ErrorIs(t, err, ErrInvalidSchema)
		})
	}
}

func TestContentEncoding(t *testing.T) {
	tests := []struct {
		name    string
		content Content
		want    string
	}{
		{name: "empty text", content: Text(""), want: `{"type":"text","text":""}`},
		{name: "image", content: Image([]byte{1, 2}, "image/png"), want: `{"type":"image","data":"AQI=","mimeType":"image/png"}`},
		{name: "image without data", content: Image(nil, "image/png"), want: `{"type":"image","data":"","mimeType":"image/png"}`},
		{name: "audio", content: Audio([]byte{1}, "audio/wav"), want: `{"type":"audio","data":"AQ==","mimeType":"audio/wav"}`},
		{name: "embedded empty text", content: EmbeddedText("mem://a", "", ""),
			want: `{"type":"resource","resource":{"uri":"mem://a","text":""}}`},
		{name: "embedded blob", content: Content{
			Type: "resource",
			Resource: &ResourceContents{
				URI:  "mem://b",
				Blob: []byte{},
			},
		}, want: `{"type":"resource","resource":{"uri":"mem://b","blob":""}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.content)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(data))
		})
	}
}

func TestInvalidJSON(t *testing.T) {
	s := newTestServer(t)

	lines := runServe(t, s, "not json\n")

	want := `{"jsonrpc":"2.0","id":null,"error":{"code":-32700,"message":"Parse error"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestBatchArray(t *testing.T) {
	s := newTestServer(t)

	input := `[{"jsonrpc":"2.0","id":1,"method":"server/discover"}]`
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"Invalid Request: batch messages are not supported"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestNullID(t *testing.T) {
	s := newTestServer(t)

	input := `{"jsonrpc":"2.0","id":null,"method":"server/discover","params":{` + metaParams("") + `}}`
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"Invalid Request: id must be a string or an integer"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestMissingID(t *testing.T) {
	s := newTestServer(t)

	input := `{"jsonrpc":"2.0","method":"tools/list","params":{` + metaParams("") + `}}`
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":null,"error":{"code":-32600,"message":"Invalid Request: id missing"}}`
	assert.Equal(t, []string{want}, lines)
}

func TestUnknownNotificationIgnored(t *testing.T) {
	s := newTestServer(t)

	input := `{"jsonrpc":"2.0","method":"notifications/unknown"}` + "\n"
	lines := runServe(t, s, input)

	assert.Empty(t, lines)
}

func TestIDEcho(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "string id", id: `"req-1"`},
		{name: "number id", id: `42`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestServer(t)

			input := request(tt.id, methodDiscover, metaParams(""))
			lines := runServe(t, s, input+"\n")

			assert.Len(t, lines, 1)
			assert.Contains(t, lines[0], `"id":`+tt.id)
		})
	}
}

func TestNoOutputAfterEOF(t *testing.T) {
	s := newTestServer(t)

	input := request(`1`, methodDiscover, metaParams("")) + "\n"
	lines := runServe(t, s, input)

	assert.Len(t, lines, 1)
}

func TestBlankLinesIgnored(t *testing.T) {
	s := newTestServer(t)

	input := "\n\n" + request(`1`, methodDiscover, metaParams("")) + "\n\n"
	lines := runServe(t, s, input)

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"resultType":"complete"`)
}
