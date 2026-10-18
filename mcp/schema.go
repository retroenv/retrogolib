package mcp

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/retroenv/retrogolib/set"
)

// SchemaFor generates a JSON Schema 2020-12 object for a struct type.
// Field names come from the json tag and then from the Go name. The
// description, enum, minimum, and maximum tags add matching schema keywords,
// and the header tag adds an x-mcp-header annotation. Fields are required
// unless they are pointers, carry omitempty, or set required:"false". The tag
// required:"true" makes any field required. A type can supply its own schema
// with a JSONSchema() json.RawMessage method. Byte slices are base64 strings,
// types with a MarshalText method are strings, and types with a MarshalJSON
// method accept any value. Recursive types and tag values that are not valid
// return ErrInvalidSchema.
func SchemaFor(t reflect.Type) (json.RawMessage, error) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: input must be a struct type", ErrInvalidSchema)
	}

	schema, err := schemaForType(t, set.New[reflect.Type]())
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("encoding schema: %w", err)
	}
	return data, nil
}

// jsonSchemaProvider lets a type supply its own JSON Schema.
type jsonSchemaProvider interface {
	JSONSchema() json.RawMessage
}

// Interface types that change the JSON encoding of a type.
var (
	jsonMarshalerType      = reflect.TypeFor[json.Marshaler]()
	jsonSchemaProviderType = reflect.TypeFor[jsonSchemaProvider]()
	textMarshalerType      = reflect.TypeFor[encoding.TextMarshaler]()
)

// schemaForType builds the schema value for one type.
func schemaForType(t reflect.Type, seen set.Set[reflect.Type]) (map[string]any, error) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
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
		return listSchema(t, seen)

	case reflect.Map:
		return mapSchema(t, seen)

	case reflect.Interface:
		return map[string]any{}, nil

	case reflect.Struct:
		return structSchema(t, seen)

	default:
		return nil, fmt.Errorf("%w: unsupported type %s", ErrInvalidSchema, t)
	}
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

// listSchema builds the schema for a slice or array type. A byte slice is a
// base64 string, as encoding/json encodes it.
func listSchema(t reflect.Type, seen set.Set[reflect.Type]) (map[string]any, error) {
	if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return map[string]any{"type": "string", "contentEncoding": "base64"}, nil
	}

	items, err := schemaForType(t.Elem(), seen)
	if err != nil {
		return nil, err
	}
	return map[string]any{"type": "array", "items": items}, nil
}

// mapSchema builds the schema for a map type with string keys.
func mapSchema(t reflect.Type, seen set.Set[reflect.Type]) (map[string]any, error) {
	if t.Key().Kind() != reflect.String {
		return nil, fmt.Errorf("%w: map key must be a string", ErrInvalidSchema)
	}

	values, err := schemaForType(t.Elem(), seen)
	if err != nil {
		return nil, err
	}
	return map[string]any{"type": "object", "additionalProperties": values}, nil
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
	var schema map[string]any
	if err := json.Unmarshal(raw, &schema); err != nil || schema == nil {
		return nil, fmt.Errorf("%w: JSONSchema of %s does not return a JSON object", ErrInvalidSchema, t)
	}
	return schema, nil
}

// structSchema builds an inline object schema for a struct type.
func structSchema(t reflect.Type, seen set.Set[reflect.Type]) (map[string]any, error) {
	properties := map[string]any{}
	required := []string{}
	if err := appendStructFields(t, seen, properties, &required); err != nil {
		return nil, err
	}

	result := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		result["required"] = required
	}
	return result, nil
}

// appendStructFields adds the schemas of the struct fields to properties. The
// fields of embedded structs without a json name are promoted, as
// encoding/json does.
func appendStructFields(t reflect.Type, seen set.Set[reflect.Type], properties map[string]any,
	required *[]string) error {

	if seen.Contains(t) {
		return fmt.Errorf("%w: recursive type %s", ErrInvalidSchema, t)
	}
	seen.Add(t)
	defer seen.Remove(t)

	for i := range t.NumField() {
		field := t.Field(i)
		name, omitEmpty, skip := parseJSONTag(field.Tag.Get("json"))
		if skip {
			continue
		}

		if field.Anonymous && name == "" {
			embedded := field.Type
			for embedded.Kind() == reflect.Pointer {
				embedded = embedded.Elem()
			}
			if embedded.Kind() == reflect.Struct {
				if err := appendStructFields(embedded, seen, properties, required); err != nil {
					return err
				}
				continue
			}
		}
		if !field.IsExported() {
			continue
		}
		if name == "" {
			name = field.Name
		}

		schema, err := fieldSchema(field, seen)
		if err != nil {
			return err
		}

		properties[name] = schema
		if fieldRequired(field, omitEmpty) {
			*required = append(*required, name)
		}
	}
	return nil
}

// fieldSchema builds the schema for one struct field and adds the keyword
// tags.
func fieldSchema(field reflect.StructField, seen set.Set[reflect.Type]) (map[string]any, error) {
	schema, err := schemaForType(field.Type, seen)
	if err != nil {
		return nil, err
	}
	if err := applyKeywords(schema, field); err != nil {
		return nil, fmt.Errorf("%w: field %s: %w", ErrInvalidSchema, field.Name, err)
	}
	return schema, nil
}

// applyKeywords adds the keyword tags to a field schema.
func applyKeywords(schema map[string]any, field reflect.StructField) error {
	if desc := field.Tag.Get("description"); desc != "" {
		schema["description"] = desc
	}
	if header := field.Tag.Get("header"); header != "" {
		schema[schemaKeyMCPHeader] = header
	}
	if enum := field.Tag.Get("enum"); enum != "" {
		values, err := parseEnum(field.Type, enum)
		if err != nil {
			return err
		}
		schema["enum"] = values
	}

	for _, keyword := range [...]string{"minimum", "maximum"} {
		tag := field.Tag.Get(keyword)
		if tag == "" {
			continue
		}
		value, err := strconv.ParseFloat(tag, 64)
		if err != nil {
			return fmt.Errorf("%s tag %q is not a number", keyword, tag)
		}
		schema[keyword] = value
	}
	return nil
}

// parseEnum splits an enum tag and converts the values to numbers for
// numeric field types.
func parseEnum(t reflect.Type, tag string) ([]any, error) {
	values := strings.Split(tag, ",")
	result := make([]any, 0, len(values))

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	numeric := t.Kind() >= reflect.Int && t.Kind() <= reflect.Float64

	for _, value := range values {
		if !numeric {
			result = append(result, value)
			continue
		}

		number, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("enum value %q is not a number", value)
		}
		result = append(result, number)
	}
	return result, nil
}

// fieldRequired reports whether a struct field belongs in the required list.
func fieldRequired(field reflect.StructField, omitEmpty bool) bool {
	switch field.Tag.Get("required") {
	case "true":
		return true

	case "false":
		return false

	default:
		return !omitEmpty && field.Type.Kind() != reflect.Pointer
	}
}

// parseJSONTag splits a json tag into its name and options.
func parseJSONTag(tag string) (name string, omitEmpty, skip bool) {
	if tag == "-" {
		return "", false, true
	}

	parts := strings.Split(tag, ",")
	for _, option := range parts[1:] {
		if option == "omitempty" || option == "omitzero" {
			omitEmpty = true
		}
	}
	return parts[0], omitEmpty, false
}
