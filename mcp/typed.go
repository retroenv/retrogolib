package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// ToolMiddleware wraps execution after argument validation and decoding.
// It is applied once per valid request. Its handler includes result encoding,
// output validation, and content creation. It must be safe for concurrent use.
type ToolMiddleware func(name string, next ToolHandler) ToolHandler

// TypedToolOptions configures typed tool results and execution.
// AddTool accepts at most one options value.
type TypedToolOptions[Out any] struct {
	// Content adds items after the JSON text item. It runs only after a
	// successful handler result passes its output schema. The result copies
	// its content buffers before middleware cleanup.
	Content func(Out) []Content
	// GenerateOutputSchema derives a schema from Out when OutputSchema is
	// absent. Leave it false for existing tools with scalar results.
	GenerateOutputSchema bool

	Middleware ToolMiddleware
}

// AddTool registers a typed handler and checks every call against its input
// schema, whether supplied or generated. A supplied output schema is checked
// on each success. Options can generate it and add content or middleware.
// The output becomes structured content and JSON text. Handler errors become
// isError results. The server also returns encoding and output validation
// errors as isError results.
func AddTool[In, Out any](s *Server, t Tool, h func(context.Context, ToolRequest, In) (Out, error),
	options ...TypedToolOptions[Out]) error {

	if len(options) > 1 {
		return errors.New("expected at most one typed tool options value")
	}
	var opts TypedToolOptions[Out]
	if len(options) == 1 {
		opts = options[0]
	}

	if len(t.InputSchema) == 0 {
		schema, err := SchemaFor(reflect.TypeFor[In]())
		if err != nil {
			return fmt.Errorf("generating input schema: %w", err)
		}
		t.InputSchema = schema
	}

	input, err := CompileSchema(t.InputSchema)
	if err != nil {
		return fmt.Errorf("compiling input schema: %w", err)
	}
	if len(t.OutputSchema) == 0 && opts.GenerateOutputSchema {
		t.OutputSchema, err = OutputSchemaFor(reflect.TypeFor[Out]())
		if err != nil {
			return fmt.Errorf("generating output schema: %w", err)
		}
	}

	var output *Schema
	if len(t.OutputSchema) > 0 {
		output, err = CompileSchema(t.OutputSchema)
		if err != nil {
			return fmt.Errorf("compiling output schema: %w", err)
		}
	}

	typed := &typedTool[In, Out]{
		handler:    h,
		middleware: opts.Middleware,
		name:       t.Name,

		input:  input,
		output: output,

		content: opts.Content,
	}

	return s.AddTool(t, typed.call)
}

type typedTool[In, Out any] struct {
	handler    func(context.Context, ToolRequest, In) (Out, error)
	middleware ToolMiddleware
	name       string

	input  *Schema
	output *Schema

	content func(Out) []Content
}

type typedCall[In, Out any] struct {
	input In
	tool  *typedTool[In, Out]
}

func (t *typedTool[In, Out]) call(ctx context.Context, req ToolRequest) (*ToolResult, error) {
	if err := validateToolArguments(t.input, req.Arguments); err != nil {
		return Errorf("%v", err), nil
	}
	var in In
	if len(req.Arguments) > 0 {
		decoder := json.NewDecoder(bytes.NewReader(req.Arguments))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&in); err != nil {
			return Errorf("invalid arguments: %v", err), nil
		}
	}

	call := typedCall[In, Out]{
		input: in,
		tool:  t,
	}
	handler := call.execute
	if t.middleware != nil {
		handler = t.middleware(t.name, handler)
	}

	return handler(ctx, req)
}

func (t *typedTool[In, Out]) result(out Out) (*ToolResult, error) {
	data, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("encoding tool result: %w", err)
	}
	if t.output != nil {
		if err := t.output.Validate(data); err != nil {
			return nil, fmt.Errorf("result does not match its schema: %w", err)
		}
	}

	content := []Content{Text(string(data))}
	if t.content != nil {
		for _, item := range t.content(out) {
			content = append(content, item.clone())
		}
	}

	return &ToolResult{
		Content:           content,
		StructuredContent: json.RawMessage(data),
	}, nil
}

func (c typedCall[In, Out]) execute(ctx context.Context, req ToolRequest) (*ToolResult, error) {
	out, err := c.tool.handler(ctx, req, c.input)
	if err != nil {
		return Errorf("%v", err), nil
	}
	return c.tool.result(out)
}

func validateToolArguments(schema *Schema, argumentsRaw json.RawMessage) error {
	if len(argumentsRaw) == 0 {
		argumentsRaw = json.RawMessage(`{}`)
	}

	if !bytes.HasPrefix(bytes.TrimSpace(argumentsRaw), []byte{'{'}) {
		return errors.New("invalid arguments: expected an object")
	}
	if err := schema.Validate(argumentsRaw); err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}
	return nil
}
