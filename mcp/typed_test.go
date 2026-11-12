package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/mcp/jsonschema"
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

func TestAddToolTypedMissingRequiredArgument(t *testing.T) {
	s := newTestServer(t)
	called := false
	assert.NoError(t, AddTool(s, Tool{Name: "add"}, func(_ context.Context, _ ToolRequest, in addArgs) (int, error) {
		called = true
		return in.A + in.B, nil
	}))

	input := request(`1`, methodToolsCall, metaParams(`"name":"add","arguments":{"a":1}`))
	lines := runServe(t, s, input+"\n")
	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"isError":true`)
	assert.False(t, called)
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

	input := request(`1`, methodToolsCall, metaParams(`"name":"fail","arguments":{"a":0,"b":0}`))
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
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
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

func TestAddToolTypedQuotedFields(t *testing.T) {
	type args struct {
		Count    int     `json:"count,string" enum:"42,43"`
		Enabled  bool    `json:"enabled,string"`
		Ratio    float64 `json:"ratio,string"`
		Text     string  `json:"text,string" enum:"hello"`
		Optional *int    `json:"optional,string"`
	}
	s := newTestServer(t)
	calls := 0
	assert.NoError(t, AddTool(s, Tool{Name: "quoted"}, func(_ context.Context, _ ToolRequest, in args) (int, error) {
		calls++
		assert.Equal(t, 42, in.Count)
		assert.True(t, in.Enabled)
		assert.Equal(t, 1.5, in.Ratio)
		assert.Equal(t, "hello", in.Text)
		assert.NotNil(t, in.Optional)
		assert.Equal(t, 7, *in.Optional)
		return in.Count, nil
	}))
	for _, test := range []struct {
		arguments string
		isError   bool
	}{
		{`{"count":"42","enabled":"true","ratio":"1.5","text":"\"hello\"","optional":"7"}`, false},
		{`{"count":42,"enabled":"true","ratio":"1.5","text":"\"hello\""}`, true},
		{`{"count":"44","enabled":"true","ratio":"1.5","text":"\"hello\""}`, true},
		{`{"count":"42","enabled":"invalid","ratio":"1.5","text":"\"hello\""}`, true},
	} {
		lines := runServe(t, s, request(`1`, methodToolsCall, metaParams(`"name":"quoted","arguments":`+test.arguments))+"\n")
		assert.Len(t, lines, 1)
		var response struct {
			Result ToolResult `json:"result"`
		}
		assert.NoError(t, json.Unmarshal([]byte(lines[0]), &response))
		assert.Equal(t, test.isError, response.Result.IsError)
	}
	assert.Equal(t, 1, calls)
}

type exactSchemaValue uint64

func (exactSchemaValue) JSONSchema() json.RawMessage {
	return json.RawMessage(`{"type":"integer","enum":[9007199254740993]}`)
}

func TestAddToolTypedCustomSchemaPrecision(t *testing.T) {
	type args struct {
		Value exactSchemaValue `json:"value"`
	}
	s := newTestServer(t)
	calls := 0
	assert.NoError(t, AddTool(s, Tool{Name: "exact"}, func(_ context.Context, _ ToolRequest, in args) (uint64, error) {
		calls++
		assert.Equal(t, uint64(9007199254740993), uint64(in.Value))
		return uint64(in.Value), nil
	}))
	for _, test := range []struct {
		arguments string
		isError   bool
	}{
		{`{"value":9007199254740993}`, false},
		{`{"value":9007199254740992}`, true},
	} {
		lines := runServe(t, s, request(`1`, methodToolsCall, metaParams(`"name":"exact","arguments":`+test.arguments))+"\n")
		assert.Len(t, lines, 1)
		var response struct {
			Result ToolResult `json:"result"`
		}
		assert.NoError(t, json.Unmarshal([]byte(lines[0]), &response))
		assert.Equal(t, test.isError, response.Result.IsError)
	}
	assert.Equal(t, 1, calls)
}

func TestAddToolTypedShadowedField(t *testing.T) {
	type Base struct {
		Count string `json:"count"`
	}
	type args struct {
		Base
		Count int `json:"count"`
	}
	s := newTestServer(t)
	assert.NoError(t, AddTool(s, Tool{Name: "shadow"}, func(_ context.Context, _ ToolRequest, in args) (int, error) {
		assert.Equal(t, "", in.Base.Count)
		return in.Count, nil
	}))
	lines := runServe(t, s, request(`1`, methodToolsCall, metaParams(`"name":"shadow","arguments":{"count":42}`))+"\n")
	assert.Len(t, lines, 1)
	assert.Contains(t, lines[0], `"structuredContent":42`)
	assert.Contains(t, lines[0], `"isError":false`)
}
