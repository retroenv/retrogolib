package mcp

import (
	"encoding/json"

	"github.com/retroenv/retrogolib/mcp/jsonschema"
)

// Schema validates JSON values against a compiled schema.
// It is safe for concurrent use.
type Schema struct {
	validator *jsonschema.Schema
}

// Validate checks one complete JSON value without rounding numbers.
func (schema *Schema) Validate(raw json.RawMessage) error {
	var validator *jsonschema.Schema
	if schema != nil {
		validator = schema.validator
	}

	return validator.Validate(raw) //nolint:wrapcheck // Preserve the public validation error contract.
}

// CompileSchema compiles a schema with [jsonschema.Compile].
// See package jsonschema for supported keywords and limits.
func CompileSchema(raw json.RawMessage) (*Schema, error) {
	validator, err := jsonschema.Compile(raw)
	if err != nil {
		return nil, err //nolint:wrapcheck // Preserve the public compilation error contract.
	}
	return &Schema{validator: validator}, nil
}
