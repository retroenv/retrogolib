package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// AddTool registers a typed handler. If no input schema is supplied, it generates
// and compiles one from In, then checks arguments before calling h.
// The output becomes both structured content and JSON text. Handler errors
// become isError results.
func AddTool[In, Out any](s *Server, t Tool, h func(ctx context.Context, req ToolRequest, in In) (Out, error)) error {
	generatedSchema := len(t.InputSchema) == 0
	if generatedSchema {
		schema, err := SchemaFor(reflect.TypeOf((*In)(nil)).Elem())
		if err != nil {
			return fmt.Errorf("generating input schema: %w", err)
		}

		t.InputSchema = schema
	}

	var inputSchema *Schema
	if generatedSchema {
		var err error
		inputSchema, err = CompileSchema(t.InputSchema)
		if err != nil {
			return fmt.Errorf("compiling input schema: %w", err)
		}
	}

	return s.AddTool(t, func(ctx context.Context, req ToolRequest) (*ToolResult, error) {
		if generatedSchema {
			if err := validateToolArguments(inputSchema, req.Arguments); err != nil {
				return Errorf("%v", err), nil
			}
		}

		var in In
		if len(req.Arguments) > 0 {
			decoder := json.NewDecoder(bytes.NewReader(req.Arguments))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&in); err != nil {
				return Errorf("invalid arguments: %v", err), nil
			}
		}

		out, err := h(ctx, req, in)
		if err != nil {
			return Errorf("%v", err), nil
		}

		data, err := json.Marshal(out)
		if err != nil {
			return nil, fmt.Errorf("encoding tool result: %w", err)
		}

		return &ToolResult{
			Content:           []Content{Text(string(data))},
			StructuredContent: out,
		}, nil
	})
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
