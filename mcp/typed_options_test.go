package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/mcp/jsonschema"
)

type typedContextKey struct{}

type typedResult struct {
	Sum int `json:"sum"`
}

type byteArgs struct {
	Value byte `json:"value"`
}

type brokenOutput struct{}

type typedSnapshotProbe struct {
	image    []byte
	resource *ResourceContents
	value    map[string]int
}

type typedExecutionProbe struct {
	active   bool
	cleanups int
	wraps    int

	calls         int
	contextInside bool
	handlerError  bool

	contentInside bool
	contents      int
}

func TestTypedToolChecksProvidedSchemaBeforeMiddleware(t *testing.T) {
	s := newTestServer(t)
	probe := &typedExecutionProbe{}
	assert.NoError(t, AddTool(s, Tool{
		Name: "checked",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"a":{"type":"integer","minimum":1},` +
			`"b":{"type":"integer"}},"required":["a","b"],"additionalProperties":false}`),
	}, probe.handle, TypedToolOptions[typedResult]{
		GenerateOutputSchema: true,

		Middleware: probe.wrap,
	}))
	for _, arguments := range []string{
		`{"a":0,"b":1}`, `{"a":1}`, `{"a":null,"b":1}`, `{"a":1,"b":1,"extra":1}`, `null`,
	} {
		result, err := s.tools[0].handler(t.Context(), ToolRequest{
			Arguments: json.RawMessage(arguments),
		})
		assert.NoError(t, err)
		assert.True(t, result.IsError)
	}
	assert.Equal(t, 0, probe.calls)
	assert.Equal(t, 0, probe.wraps)
}

func TestTypedToolMiddlewareIncludesImageContentAndJSON(t *testing.T) {
	s := newTestServer(t)
	probe := &typedExecutionProbe{}
	assert.NoError(t, AddTool(s, Tool{
		Name: "image",
	}, probe.handle, TypedToolOptions[typedResult]{
		Content:              probe.content,
		GenerateOutputSchema: true,

		Middleware: probe.wrap,
	}))
	result, err := s.tools[0].handler(t.Context(), ToolRequest{
		Arguments: json.RawMessage(`{"a":2,"b":3}`),
	})
	assert.NoError(t, err)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 2)
	assert.Equal(t, `{"sum":5}`, result.Content[0].Text)
	assert.Equal(t, "image/png", result.Content[1].MimeType)
	assert.Equal(t, []byte{1, 2}, result.Content[1].Data)
	assert.Equal(t, 1, probe.calls)
	assert.Equal(t, 1, probe.contents)
	assert.Equal(t, 1, probe.cleanups)
	assert.True(t, probe.contentInside)
	assert.True(t, probe.contextInside)
	data, err := json.Marshal(result.StructuredContent)
	assert.NoError(t, err)
	assert.Equal(t, result.Content[0].Text, string(data))
	validator, err := CompileSchema(s.tools[0].tool.OutputSchema)
	assert.NoError(t, err)
	assert.NoError(t, validator.Validate(data))
}

func TestTypedToolKeepsSuppliedOutputAndCleansUpOnFailure(t *testing.T) {
	for _, handlerError := range []bool{false, true} {
		s := newTestServer(t)
		probe := &typedExecutionProbe{
			handlerError: handlerError,
		}
		schema := json.RawMessage(`{"type":"object","properties":{"sum":{"type":"integer","enum":[5]}},"required":["sum"]}`)
		assert.NoError(t, AddTool(s, Tool{
			Name:         "failure",
			OutputSchema: schema,
		}, probe.handle, TypedToolOptions[typedResult]{
			Content:              probe.content,
			GenerateOutputSchema: true,

			Middleware: probe.wrap,
		}))
		assert.Equal(t, schema, s.tools[0].tool.OutputSchema)
		result := typedWireResult(t, s, "failure", `{"a":1,"b":2}`)
		assert.True(t, result.IsError)
		if handlerError {
			assert.Contains(t, result.Content[0].Text, "handler failed")
		} else {
			assert.Contains(t, result.Content[0].Text, "result does not match its schema")
		}

		assert.Equal(t, 0, probe.contents)
		assert.Equal(t, 1, probe.cleanups)
	}
}

func TestTypedToolOutputGenerationIsOptIn(t *testing.T) {
	s := newTestServer(t)
	assert.NoError(t, AddTool(s, Tool{
		Name: "scalar",
	}, addHandler))
	assert.Empty(t, s.tools[0].tool.OutputSchema)
	err := AddTool(s, Tool{
		Name: "scalar-output",
	}, addHandler, TypedToolOptions[int]{
		GenerateOutputSchema: true,
	})
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
}

func TestTypedToolDecodesBeforeMiddleware(t *testing.T) {
	s := newTestServer(t)
	probe := &typedExecutionProbe{}
	assert.NoError(t, AddTool(s, Tool{
		Name: "byte",
	}, probe.handleByte, TypedToolOptions[typedResult]{
		Middleware: probe.wrap,
	}))
	result, err := s.tools[0].handler(t.Context(), ToolRequest{
		Arguments: json.RawMessage(`{"value":256}`),
	})
	assert.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Equal(t, 0, probe.wraps)
	assert.Equal(t, 0, probe.calls)
}

func TestTypedToolCleansUpAfterEncodingError(t *testing.T) {
	s := newTestServer(t)
	probe := &typedExecutionProbe{}
	assert.NoError(t, AddTool(s, Tool{
		Name:         "encoding-error",
		OutputSchema: json.RawMessage(`{"type":"object"}`),
	}, brokenOutputHandler, TypedToolOptions[brokenOutput]{
		Middleware: probe.wrap,
	}))
	result := typedWireResult(t, s, "encoding-error", `{}`)
	assert.True(t, result.IsError)
	assert.Contains(t, result.Content[0].Text, "encoding tool result")
	assert.Equal(t, 1, probe.cleanups)
}

func TestTypedToolPreservesValidatedResultAfterCallbacks(t *testing.T) {
	// Callbacks can change the handler result after validation.
	// Both wire forms must keep the value that passed validation.
	for _, withContent := range []bool{false, true} {
		s := newTestServer(t)
		probe := &typedSnapshotProbe{
			image: []byte("image"),
			resource: &ResourceContents{
				URI:  "mem://snapshot",
				Blob: []byte("resource"),
			},
			value: map[string]int{"sum": 5},
		}
		opts := TypedToolOptions[map[string]int]{
			Middleware: probe.wrap,
		}
		if withContent {
			opts.Content = probe.content
		}

		schema := json.RawMessage(`{"type":"object","properties":{"sum":{"type":"integer","enum":[5]}},"required":["sum"]}`)
		assert.NoError(t, AddTool(s, Tool{
			Name:         "snapshot",
			OutputSchema: schema,
		}, probe.handle, opts))

		result := typedWireResult(t, s, "snapshot", `{}`)
		assert.False(t, result.IsError)
		assert.Empty(t, probe.value)
		data, err := json.Marshal(result.StructuredContent)
		assert.NoError(t, err)
		assert.Equal(t, `{"sum":5}`, string(data))
		assert.Equal(t, result.Content[0].Text, string(data))
		if withContent {
			assert.Len(t, result.Content, 3)
			assert.Equal(t, "image", string(result.Content[1].Data))
			assert.Equal(t, "mem://snapshot", result.Content[2].Resource.URI)
			assert.Equal(t, "resource", string(result.Content[2].Resource.Blob))
		}

		validator, err := CompileSchema(schema)
		assert.NoError(t, err)
		assert.NoError(t, validator.Validate(data))
	}
}

func TestTypedToolRejectsInvalidRegistrationOptions(t *testing.T) {
	s := newTestServer(t)
	err := AddTool(s, Tool{
		Name: "options",
	}, addHandler, TypedToolOptions[int]{}, TypedToolOptions[int]{})
	assert.Error(t, err)
	err = AddTool(s, Tool{
		Name:         "invalid-output",
		OutputSchema: json.RawMessage(`{"type":"object","unsupported":true}`),
	}, addHandler)
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
	assert.Empty(t, s.tools)
}

func (p *typedExecutionProbe) handle(ctx context.Context, _ ToolRequest, in addArgs) (typedResult, error) {
	p.calls++
	p.contextInside = ctx.Value(typedContextKey{}) == p
	if p.handlerError {
		return typedResult{}, errors.New("handler failed")
	}
	return typedResult{
		Sum: in.A + in.B,
	}, nil
}

func (p *typedExecutionProbe) handleByte(_ context.Context, _ ToolRequest, in byteArgs) (typedResult, error) {
	p.calls++
	return typedResult{
		Sum: int(in.Value),
	}, nil
}

func (p *typedExecutionProbe) content(_ typedResult) []Content {
	p.contents++
	p.contentInside = p.active
	return []Content{Image([]byte{1, 2}, "image/png")}
}

func (p *typedExecutionProbe) wrap(_ string, next ToolHandler) ToolHandler {
	p.wraps++
	return func(ctx context.Context, req ToolRequest) (*ToolResult, error) {
		p.active = true
		defer func() {
			p.active = false
			p.cleanups++
		}()
		return next(context.WithValue(ctx, typedContextKey{}, p), req)
	}
}

func (brokenOutput) MarshalJSON() ([]byte, error) {
	return nil, errors.New("cannot encode result")
}

func (p *typedSnapshotProbe) content(out map[string]int) []Content {
	out["sum"]++
	return []Content{
		Image(p.image, "image/png"),
		{
			Type:     contentTypeResource,
			Resource: p.resource,
		},
	}
}

func (p *typedSnapshotProbe) handle(context.Context, ToolRequest, struct{}) (map[string]int, error) {
	return p.value, nil
}

func (p *typedSnapshotProbe) wrap(_ string, next ToolHandler) ToolHandler {
	return func(ctx context.Context, req ToolRequest) (*ToolResult, error) {
		defer clear(p.value)
		defer clear(p.image)
		defer clear(p.resource.Blob)
		defer func() {
			p.resource.URI = ""
		}()
		return next(ctx, req)
	}
}

func brokenOutputHandler(_ context.Context, _ ToolRequest, _ struct{}) (brokenOutput, error) {
	return brokenOutput{}, nil
}

func typedWireResult(t *testing.T, server *Server, name, arguments string) *ToolResult {
	t.Helper()

	params := metaParams(`"name":"` + name + `","arguments":` + arguments)
	lines := runServe(t, server, request(`1`, methodToolsCall, params)+"\n")
	assert.Len(t, lines, 1)

	var response struct {
		Error  *Error      `json:"error"`
		Result *ToolResult `json:"result"`
	}
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &response))
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Result)
	return response.Result
}
