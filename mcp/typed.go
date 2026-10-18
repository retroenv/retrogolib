package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// AddTool registers a typed tool. When the tool carries no input schema, one
// is generated from In. The output of h becomes the structured content, and
// the text content holds the JSON form of the output. Handler errors become
// isError results.
func AddTool[In, Out any](s *Server, t Tool, h func(ctx context.Context, req ToolRequest, in In) (Out, error)) error {
	if len(t.InputSchema) == 0 {
		schema, err := SchemaFor(reflect.TypeOf((*In)(nil)).Elem())
		if err != nil {
			return fmt.Errorf("generating input schema: %w", err)
		}
		t.InputSchema = schema
	}

	return s.AddTool(t, func(ctx context.Context, req ToolRequest) (*ToolResult, error) {
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
