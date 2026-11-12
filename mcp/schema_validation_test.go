package mcp

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/mcp/jsonschema"
)

type unsupportedPattern string

type unsupportedPatternArgs struct {
	Name unsupportedPattern `json:"name"`
}

func (unsupportedPattern) JSONSchema() json.RawMessage {
	return json.RawMessage(`{"type":"string","pattern":".*"}`)
}

func echoUnsupportedPattern(_ context.Context, _ ToolRequest, in unsupportedPatternArgs) (string, error) {
	return string(in.Name), nil
}

func TestGeneratedToolRejectsUnsupportedKeywordAtRegistration(t *testing.T) {
	server := newTestServer(t)
	err := AddTool(server, Tool{Name: "unsupported_pattern"}, echoUnsupportedPattern)
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
}

func TestUncompiledSchemaDoesNotPanic(t *testing.T) {
	var absent *Schema
	assert.ErrorIs(t, absent.Validate(json.RawMessage(`null`)), jsonschema.ErrInvalidSchema)
	var empty Schema
	assert.ErrorIs(t, empty.Validate(json.RawMessage(`null`)), jsonschema.ErrInvalidSchema)
}

func TestGeneratedSchemaKeepsExactNumericTags(t *testing.T) {
	type limits struct {
		Value uint64 `json:"value" enum:"9007199254740993,18446744073709551615" minimum:"9007199254740993" maximum:"18446744073709551615"`
	}

	raw, err := SchemaFor(reflect.TypeFor[limits]())
	assert.NoError(t, err)
	schema, err := CompileSchema(raw)
	assert.NoError(t, err)
	for _, value := range []string{`{"value":9007199254740993}`, `{"value":18446744073709551615}`} {
		assert.NoError(t, schema.Validate(json.RawMessage(value)))
	}

	for _, value := range []string{`{"value":9007199254740992}`, `{"value":18446744073709551614}`} {
		assert.Error(t, schema.Validate(json.RawMessage(value)))
	}
}

func TestSchemaAnnotationsDoNotChangeValidation(t *testing.T) {
	raw := json.RawMessage(`{"$schema":"https://json-schema.org/draft/2020-12/schema","$comment":"test","title":"test","description":"test","default":null,"examples":[null],"readOnly":true,"writeOnly":false,"deprecated":false,"contentEncoding":"base64","contentMediaType":"image/png","x-mcp-header":"label"}`)
	schema, err := CompileSchema(raw)
	assert.NoError(t, err)
	assert.NoError(t, schema.Validate(json.RawMessage(`"not base64"`)))
}
