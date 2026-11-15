package mcp

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/retroenv/retrogolib/mcp/jsonschema"
	"github.com/retroenv/retrogolib/set"
)

// Interface types that change the JSON encoding of a type.
var (
	jsonMarshalerType      = reflect.TypeFor[json.Marshaler]()
	jsonSchemaProviderType = reflect.TypeFor[jsonSchemaProvider]()
	textMarshalerType      = reflect.TypeFor[encoding.TextMarshaler]()
)

// SchemaFor generates a JSON Schema 2020-12 object for a struct type.
// Field names come from the json tag and then from the Go name. The
// description, enum, minimum, and maximum tags add matching schema keywords,
// and the header tag adds an x-mcp-header annotation. The JSON string option
// produces a string schema for scalar fields. Numeric bound tags are not
// supported with this option. Fields are required unless they are pointers,
// use omitempty or omitzero, or set required:"false". The tag
// required:"true" makes any field required. A type can supply its own schema
// with a JSONSchema() json.RawMessage method. Byte slices are base64 strings,
// types with a MarshalText method are strings, and types with a MarshalJSON
// method accept any value. Recursive types and tag values that are not valid
// return jsonschema.ErrInvalidSchema.
func SchemaFor(t reflect.Type) (json.RawMessage, error) {
	return generateSchema(t, false)
}

// OutputSchemaFor generates a schema for a struct or string-keyed map result.
// Nil pointers, slices, maps, and byte slices can encode as null. Fixed arrays
// stay arrays. Input schemas from SchemaFor keep their current null rules.
// A result must have an object schema at the root.
func OutputSchemaFor(t reflect.Type) (json.RawMessage, error) {
	return generateSchema(t, true)
}

// schemaBuilder applies the shared reflection rules to inputs or outputs.
type schemaBuilder struct {
	output bool

	seen set.Set[reflect.Type]
}

// jsonSchemaProvider lets a type supply its own JSON Schema.
type jsonSchemaProvider interface {
	JSONSchema() json.RawMessage
}

// schemaForType builds the schema value for one type.
func (b schemaBuilder) schemaForType(t reflect.Type) (map[string]any, error) {
	schema, err := b.nonNullSchema(t)
	if err != nil {
		return nil, err
	}
	return b.allowNull(t, schema)
}

func (b schemaBuilder) nonNullSchema(t reflect.Type) (map[string]any, error) {
	if b.seen.Contains(t) {
		return nil, fmt.Errorf("%w: recursive type %s", jsonschema.ErrInvalidSchema, t)
	}
	b.seen.Add(t)
	defer b.seen.Remove(t)

	if t.Kind() == reflect.Pointer {
		return b.nonNullSchema(t.Elem())
	}

	if schema, ok, err := encoderSchema(t); ok {
		return schema, err
	}

	switch t.Kind() {
	case reflect.String:
		return map[string]any{"type": "string"}, nil

	case reflect.Bool:
		return map[string]any{"type": "boolean"}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]any{"type": "integer"}, nil

	case reflect.Float32, reflect.Float64:
		return map[string]any{"type": "number"}, nil

	case reflect.Slice, reflect.Array:
		return b.listSchema(t)

	case reflect.Map:
		return b.mapSchema(t)

	case reflect.Interface:
		return map[string]any{}, nil

	case reflect.Struct:
		return b.structSchema(t)

	default:
		return nil, fmt.Errorf("%w: unsupported type %s", jsonschema.ErrInvalidSchema, t)
	}
}

// listSchema builds the schema for a slice or array type. A byte slice is a
// base64 string, as encoding/json encodes it.
func (b schemaBuilder) listSchema(t reflect.Type) (map[string]any, error) {
	if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return map[string]any{"type": "string", "contentEncoding": "base64"}, nil
	}

	items, err := b.schemaForType(t.Elem())
	if err != nil {
		return nil, err
	}
	return map[string]any{"type": "array", "items": items}, nil
}

// mapSchema builds the schema for a map type with string keys.
func (b schemaBuilder) mapSchema(t reflect.Type) (map[string]any, error) {
	if t.Key().Kind() != reflect.String {
		return nil, fmt.Errorf("%w: map key must be a string", jsonschema.ErrInvalidSchema)
	}

	values, err := b.schemaForType(t.Elem())
	if err != nil {
		return nil, err
	}
	return map[string]any{"type": "object", "additionalProperties": values}, nil
}

// structSchema builds an inline object schema for a struct type.
func (b schemaBuilder) structSchema(t reflect.Type) (map[string]any, error) {
	properties := map[string]any{}
	required := []string{}
	if err := b.appendStructFields(t, properties, &required); err != nil {
		return nil, err
	}

	result := map[string]any{
		"additionalProperties": false,
		"properties":           properties,
		"type":                 "object",
	}
	if len(required) > 0 {
		result["required"] = required
	}

	return result, nil
}

// allowNull keeps field constraints on the non-null branch of an output.
func (b schemaBuilder) allowNull(t reflect.Type, schema map[string]any) (map[string]any, error) {
	if !b.output || !nilableJSONType(t) {
		return schema, nil
	}
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("encoding nullable output schema: %w", err)
	}
	validator, err := CompileSchema(data)
	if err != nil {
		return nil, err
	}
	if validator.Validate(json.RawMessage(`null`)) == nil {
		return schema, nil
	}
	return map[string]any{
		"oneOf": []any{schema, map[string]any{"type": "null"}},
	}, nil
}

func generateSchema(t reflect.Type, output bool) (json.RawMessage, error) {
	root, err := schemaRootType(t, output)
	if err != nil {
		return nil, err
	}
	builder := schemaBuilder{
		output: output,

		seen: set.New[reflect.Type](),
	}
	// The root is an object even when its Go type is a map or pointer.
	schema, err := builder.nonNullSchema(root)
	if err != nil {
		return nil, err
	}
	if output && schema["type"] != "object" {
		return nil, fmt.Errorf("%w: output schema root must be an object", jsonschema.ErrInvalidSchema)
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("encoding schema: %w", err)
	}
	return data, nil
}

func schemaRootType(t reflect.Type, output bool) (reflect.Type, error) {
	if t == nil {
		return nil, fmt.Errorf("%w: type must not be nil", jsonschema.ErrInvalidSchema)
	}
	seen := set.New[reflect.Type]()

	for t.Kind() == reflect.Pointer {
		if seen.Contains(t) {
			return nil, fmt.Errorf("%w: recursive type %s", jsonschema.ErrInvalidSchema, t)
		}
		seen.Add(t)
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct && (!output || t.Kind() != reflect.Map) {
		return nil, fmt.Errorf("%w: root must be a struct, or a map for output", jsonschema.ErrInvalidSchema)
	}
	return t, nil
}

// encoderSchema returns the schema of a type that controls its own JSON
// encoding. The ok result is false for other types.
func encoderSchema(t reflect.Type) (schema map[string]any, ok bool, err error) {
	switch {
	case implements(t, jsonSchemaProviderType):
		schema, err = customSchema(t)
		return schema, true, err

	case implements(t, jsonMarshalerType):
		return map[string]any{}, true, nil // any JSON value

	case implements(t, textMarshalerType):
		return map[string]any{"type": "string"}, true, nil

	default:
		return nil, false, nil
	}
}

// implements reports whether the type or its pointer implements iface.
func implements(t, iface reflect.Type) bool {
	return t.Implements(iface) || reflect.PointerTo(t).Implements(iface)
}

// customSchema decodes the schema that a type supplies through the
// JSONSchema hook. The hook receives a zero value.
func customSchema(t reflect.Type) (map[string]any, error) {
	value := reflect.New(t)
	if !t.Implements(jsonSchemaProviderType) {
		return decodeCustomSchema(t, value.Interface().(jsonSchemaProvider).JSONSchema())
	}
	return decodeCustomSchema(t, value.Elem().Interface().(jsonSchemaProvider).JSONSchema())
}

// decodeCustomSchema decodes a type-supplied schema into a map, so that the
// field tags can add keywords to it.
func decodeCustomSchema(t reflect.Type, raw json.RawMessage) (map[string]any, error) {
	value, err := jsonschema.DecodeValue(raw)
	schema, ok := value.(map[string]any)
	if err != nil || !ok {
		return nil, fmt.Errorf("%w: JSONSchema of %s does not return a JSON object", jsonschema.ErrInvalidSchema, t)
	}
	return schema, nil
}

func nilableJSONType(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice || t.Kind() == reflect.Map
}
