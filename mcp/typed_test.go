package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type addArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

func addHandler(_ context.Context, _ ToolRequest, in addArgs) (int, error) {
	return in.A + in.B, nil
}

func TestAddToolTyped(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, AddTool(s, Tool{
		Name:        "add",
		Description: "Add two integers.",
	}, addHandler))

	input := request(`1`, methodToolsCall, metaParams(`"name":"add","arguments":{"a":2,"b":3}`))
	lines := runServe(t, s, input+"\n")

	want := `{"jsonrpc":"2.0","id":1,"result":{"resultType":"complete","content":[{"type":"text",` +
		`"text":"5"}],"structuredContent":5,"isError":false,"_meta":` +
		`{"io.modelcontextprotocol/serverInfo":{"name":"test-server","version":"1.2.3"}}}}`
	assert.Equal(t, []string{want}, lines)
}

func TestAddToolTypedGeneratedSchema(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, AddTool(s, Tool{
		Name:        "add",
		Description: "Add two integers.",
	}, addHandler))

	input := request(`1`, methodToolsList, metaParams(""))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)

	var response struct {
		Result struct {
			Tools []Tool `json:"tools"`
		} `json:"result"`
	}
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &response))

	assert.Len(t, response.Result.Tools, 1)
	assert.Equal(t, "add", response.Result.Tools[0].Name)

	var schema map[string]any
	assert.NoError(t, json.Unmarshal(response.Result.Tools[0].InputSchema, &schema))
	assert.Equal(t, []any{"a", "b"}, schema["required"])
}

func TestAddToolTypedProvidedSchema(t *testing.T) {
	s := newTestServer(t)
	provided := json.RawMessage(`{"type":"object"}`)
	assert.NoError(t, AddTool(s, Tool{
		Name:        "add",
		InputSchema: provided,
	}, addHandler))

	input := request(`1`, methodToolsList, metaParams(""))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)

	var response struct {
		Result struct {
			Tools []Tool `json:"tools"`
		} `json:"result"`
	}
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &response))

	assert.Len(t, response.Result.Tools, 1)
	assert.Equal(t, provided, response.Result.Tools[0].InputSchema)
}

func TestAddToolTypedInvalidArguments(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, AddTool(s, Tool{
		Name: "add",
	}, addHandler))

	input := request(`1`, methodToolsCall, metaParams(`"name":"add","arguments":{"a":1,"b":"text"}`))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"isError":true`)
	assert.Contains(t, lines[0], "invalid arguments")
}

func TestAddToolTypedUnknownField(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, AddTool(s, Tool{
		Name: "add",
	}, addHandler))

	input := request(`1`, methodToolsCall, metaParams(`"name":"add","arguments":{"a":1,"b":2,"c":3}`))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"isError":true`)
}

func TestAddToolTypedHandlerError(t *testing.T) {
	s := newTestServer(t)
	err := AddTool(s, Tool{Name: "fail"}, func(_ context.Context, _ ToolRequest, _ addArgs) (int, error) {
		return 0, errors.New("not available")
	})
	assert.NoError(t, err)

	input := request(`1`, methodToolsCall, metaParams(`"name":"fail"`))
	lines := runServe(t, s, input+"\n")

	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"isError":true`)
	assert.Contains(t, lines[0], "not available")
}

func TestAddToolTypedInvalidSchemaType(t *testing.T) {
	s := newTestServer(t)
	err := AddTool(s, Tool{
		Name: "bad",
	}, func(_ context.Context, _ ToolRequest, _ map[int]string) (string, error) {
		return "", nil
	})
	assert.ErrorIs(t, err, ErrInvalidSchema)
}

func TestAddToolDuplicate(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, s.AddTool(Tool{
		Name: "dup",
	}, echoHandler))
	assert.ErrorIs(t, s.AddTool(Tool{
		Name: "dup",
	}, echoHandler), ErrDuplicateTool)
}

func TestAddToolEmptyName(t *testing.T) {
	s := newTestServer(t)
	assert.ErrorIs(t, s.AddTool(Tool{}, echoHandler), ErrEmptyName)
}
