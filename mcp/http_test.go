package mcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func newHTTPServer(t *testing.T) *httptest.Server {
	t.Helper()

	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{Name: "echo"}, echoHandler))

	ts := httptest.NewServer(s.Handler())
	t.Cleanup(ts.Close)
	return ts
}

// httpHeaders builds the standard headers for a request.
func httpHeaders(method, name string) map[string]string {
	headers := map[string]string{
		headerProtocolVersion: ProtocolVersion,
		headerMethod:          method,
	}
	if name != "" {
		headers[headerName] = name
	}
	return headers
}

func post(t *testing.T, ts *httptest.Server, body string, headers map[string]string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, ts.URL, strings.NewReader(body))
	assert.NoError(t, err)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, resp.Body.Close())
	})
	return resp
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	data, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	return strings.TrimSpace(string(data))
}

func TestHTTPDiscover(t *testing.T) {
	ts := newHTTPServer(t)

	body := request(`1`, methodDiscover, metaParams(""))
	resp := post(t, ts, body, httpHeaders(methodDiscover, ""))

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Contains(t, readBody(t, resp), `"resultType":"complete"`)
}

func TestHTTPMissingVersionHeader(t *testing.T) {
	ts := newHTTPServer(t)

	headers := httpHeaders(methodDiscover, "")
	delete(headers, headerProtocolVersion)

	body := request(`1`, methodDiscover, metaParams(""))
	resp := post(t, ts, body, headers)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"code":-32020`)
}

func TestHTTPVersionHeaderMismatch(t *testing.T) {
	ts := newHTTPServer(t)

	headers := httpHeaders(methodDiscover, "")
	headers[headerProtocolVersion] = "2025-11-25"

	body := request(`1`, methodDiscover, metaParams(""))
	resp := post(t, ts, body, headers)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"code":-32020`)
}

func TestHTTPMissingMethodHeader(t *testing.T) {
	ts := newHTTPServer(t)

	headers := httpHeaders(methodDiscover, "")
	delete(headers, headerMethod)

	body := request(`1`, methodDiscover, metaParams(""))
	resp := post(t, ts, body, headers)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"code":-32020`)
}

func TestHTTPMethodHeaderMismatch(t *testing.T) {
	ts := newHTTPServer(t)

	headers := httpHeaders(methodToolsList, "")

	body := request(`1`, methodDiscover, metaParams(""))
	resp := post(t, ts, body, headers)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"code":-32020`)
}

func TestHTTPMissingNameHeader(t *testing.T) {
	ts := newHTTPServer(t)

	headers := httpHeaders(methodToolsCall, "")

	body := request(`1`, methodToolsCall, metaParams(`"name":"echo"`))
	resp := post(t, ts, body, headers)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"code":-32020`)
}

func TestHTTPNameHeaderMismatch(t *testing.T) {
	ts := newHTTPServer(t)

	body := request(`1`, methodToolsCall, metaParams(`"name":"echo"`))
	resp := post(t, ts, body, httpHeaders(methodToolsCall, "other"))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"code":-32020`)
}

func TestHTTPBase64NameHeader(t *testing.T) {
	ts := newHTTPServer(t)

	headers := httpHeaders(methodToolsCall, "")
	headers[headerName] = base64HeaderPrefix + base64.StdEncoding.EncodeToString([]byte("echo")) + base64HeaderSuffix

	body := request(`1`, methodToolsCall, metaParams(`"name":"echo","arguments":{"x":1}`))
	resp := post(t, ts, body, headers)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"isError":false`)
}

func TestHTTPWrongOrigin(t *testing.T) {
	ts := newHTTPServer(t)

	headers := httpHeaders(methodDiscover, "")
	headers["Origin"] = "https://evil.example"

	body := request(`1`, methodDiscover, metaParams(""))
	resp := post(t, ts, body, headers)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestHTTPAllowedOrigin(t *testing.T) {
	s := newTestServer(t)
	s.opts.AllowedOrigins = []string{"https://good.example"}

	ts := httptest.NewServer(s.Handler())
	t.Cleanup(ts.Close)

	headers := httpHeaders(methodDiscover, "")
	headers["Origin"] = "https://good.example"

	body := request(`1`, methodDiscover, metaParams(""))
	resp := post(t, ts, body, headers)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTPGet(t *testing.T) {
	ts := newHTTPServer(t)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, ts.URL, nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, resp.Body.Close())
	}()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.Equal(t, http.MethodPost, resp.Header.Get("Allow"))
}

func TestHTTPNotification(t *testing.T) {
	ts := newHTTPServer(t)

	body := `{"jsonrpc":"2.0","method":"notifications/cancelled","params":{"requestId":1}}`
	resp := post(t, ts, body, httpHeaders(methodCancelled, ""))

	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	assert.Empty(t, readBody(t, resp))
}

func TestHTTPUnknownMethod(t *testing.T) {
	ts := newHTTPServer(t)

	body := request(`1`, "unknown/method", metaParams(""))
	resp := post(t, ts, body, httpHeaders("unknown/method", ""))

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"code":-32601`)
}

func TestHTTPLegacyInitialize(t *testing.T) {
	ts := newHTTPServer(t)

	body := request(`0`, methodInitialize, `"protocolVersion":"2025-11-25","capabilities":{}`)
	resp := post(t, ts, body, nil)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	text := readBody(t, resp)
	assert.Contains(t, text, `"code":-32022`)
	assert.Contains(t, text, `"supported":["2026-07-28"],"requested":"2025-11-25"`)
}

func TestHTTPNotificationWithoutHeaders(t *testing.T) {
	ts := newHTTPServer(t)

	body := `{"jsonrpc":"2.0","method":"notifications/initialized"}`
	resp := post(t, ts, body, nil)

	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	assert.Empty(t, readBody(t, resp))
}

func TestHTTPBodyTooLarge(t *testing.T) {
	s := New(nil, Options{
		Name:           serverName,
		MaxMessageSize: 32,
	})
	ts := httptest.NewServer(s.Handler())
	t.Cleanup(ts.Close)

	body := request(`1`, methodDiscover, metaParams(""))
	resp := post(t, ts, body, httpHeaders(methodDiscover, ""))

	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode)
}

func TestHTTPParamHeaders(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{
		Name: "query",
		InputSchema: json.RawMessage(`{"type":"object","properties":{` +
			`"region":{"type":"string","x-mcp-header":"Region"},` +
			`"limit":{"type":"integer","x-mcp-header":"Limit"},` +
			`"options":{"type":"object","properties":{"dry":{"type":"boolean","x-mcp-header":"Dry"}}}}}`),
	}, echoHandler))

	ts := httptest.NewServer(s.Handler())
	t.Cleanup(ts.Close)

	tests := []struct {
		name      string
		arguments string
		headers   map[string]string
		status    int
	}{
		{name: "all headers match", arguments: `{"region":"us-west1","limit":42,"options":{"dry":true}}`,
			headers: map[string]string{"Mcp-Param-Region": "us-west1", "Mcp-Param-Limit": "42.0", "Mcp-Param-Dry": "true"},
			status:  http.StatusOK},
		{name: "base64 value", arguments: `{"region":"Grüße"}`,
			headers: map[string]string{"Mcp-Param-Region": "=?base64?R3LDvMOfZQ==?="},
			status:  http.StatusOK},
		{name: "absent and null values need no header", arguments: `{"region":null}`,
			status: http.StatusOK},
		{name: "missing header", arguments: `{"region":"us-west1"}`,
			status: http.StatusBadRequest},
		{name: "value mismatch", arguments: `{"limit":42}`,
			headers: map[string]string{"Mcp-Param-Limit": "43"},
			status:  http.StatusBadRequest},
		{name: "nested value mismatch", arguments: `{"options":{"dry":false}}`,
			headers: map[string]string{"Mcp-Param-Dry": "true"},
			status:  http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := httpHeaders(methodToolsCall, "query")
			for key, value := range tt.headers {
				headers[key] = value
			}

			body := request(`1`, methodToolsCall, metaParams(`"name":"query","arguments":`+tt.arguments))
			resp := post(t, ts, body, headers)

			assert.Equal(t, tt.status, resp.StatusCode)
			if tt.status == http.StatusBadRequest {
				assert.Contains(t, readBody(t, resp), `"code":-32020`)
			}
		})
	}
}

func TestHTTPUnsupportedVersion(t *testing.T) {
	ts := newHTTPServer(t)

	headers := httpHeaders(methodDiscover, "")
	headers[headerProtocolVersion] = "2025-11-25"

	params := `"_meta":{"io.modelcontextprotocol/protocolVersion":"2025-11-25","io.modelcontextprotocol/clientCapabilities":{}}`
	body := request(`1`, methodDiscover, params)
	resp := post(t, ts, body, headers)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"code":-32022`)
}

func TestHTTPInvalidJSON(t *testing.T) {
	ts := newHTTPServer(t)

	resp := post(t, ts, "not json", httpHeaders(methodDiscover, ""))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, readBody(t, resp), `"code":-32700`)
}

func TestToolProgress(t *testing.T) {
	s := newTestServer(t)

	assert.NoError(t, s.AddTool(Tool{Name: "work"}, func(_ context.Context, req ToolRequest) (*ToolResult, error) {
		if req.Progress != nil {
			req.Progress(0.5, 1, "half")
		}
		return &ToolResult{Content: []Content{Text("done")}}, nil
	}))

	meta := `"progressToken":"tok-1","io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{}`
	input := request(`1`, methodToolsCall, `"name":"work","_meta":{`+meta+`}`)
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 2)
	assert.Equal(t, `{"jsonrpc":"2.0","method":"notifications/progress","params":{`+
		`"message":"half","progress":0.5,"progressToken":"tok-1","total":1}}`, lines[0])
	assert.Contains(t, lines[1], `"text":"done"`)
}

func TestToolNoProgressToken(t *testing.T) {
	s := newTestServer(t)

	assert.NoError(t, s.AddTool(Tool{Name: "work"}, func(_ context.Context, req ToolRequest) (*ToolResult, error) {
		req.Progress(0.5, 1, "half")
		return &ToolResult{Content: []Content{Text("done")}}, nil
	}))

	input := request(`1`, methodToolsCall, metaParams(`"name":"work"`))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"text":"done"`)
}

func TestHTTPToolProgressIgnored(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{Name: "work"}, func(_ context.Context, req ToolRequest) (*ToolResult, error) {
		req.Progress(0.5, 1, "half")
		return &ToolResult{Content: []Content{Text("done")}}, nil
	}))
	ts := httptest.NewServer(s.Handler())
	t.Cleanup(ts.Close)

	meta := `"progressToken":1,` + testMeta
	body := request(`1`, methodToolsCall, `"name":"work","_meta":{`+meta+`}`)
	resp := post(t, ts, body, httpHeaders(methodToolsCall, "work"))

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	text := readBody(t, resp)
	assert.Contains(t, text, `"text":"done"`)
	assert.Contains(t, text, `"isError":false`)
}
