package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/retroenv/retrogolib/set"
)

const maxSchemaDepth = 64

const (
	keywordType                 = "type"
	keywordProperties           = "properties"
	keywordItems                = "items"
	keywordAdditionalProperties = "additionalProperties"
	keywordOneOf                = "oneOf"
	keywordRequired             = "required"
	keywordEnum                 = "enum"
	keywordMinimum              = "minimum"
	keywordMaximum              = "maximum"
	keywordMinItems             = "minItems"
	keywordMaxItems             = "maxItems"
	keywordUniqueItems          = "uniqueItems"

	typeObject  = "object"
	typeArray   = "array"
	typeString  = "string"
	typeBoolean = "boolean"
	typeInteger = "integer"
	typeNumber  = "number"
	typeNull    = "null"
)

// ErrInvalidSchema identifies a schema that the validator cannot compile.
var ErrInvalidSchema = errors.New("invalid schema")

// Schema is an immutable validator. It is safe for concurrent use.
type Schema struct {
	root any
}

// Validate checks one complete JSON value against the schema.
func (schema *Schema) Validate(raw json.RawMessage) error {
	if schema == nil || schema.root == nil {
		return fmt.Errorf("%w: schema has not been compiled", ErrInvalidSchema)
	}

	value, err := DecodeValue(raw)
	if err != nil {
		return fmt.Errorf("decoding JSON value: %w", err)
	}

	return validateSchemaValue(schema.root, value, "value")
}

// Compile checks and copies a schema for validation.
// Unsupported keywords and schemas with more than 64 levels return ErrInvalidSchema.
func Compile(raw json.RawMessage) (*Schema, error) {
	value, err := DecodeValue(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidSchema, err)
	}

	if err := checkSchemaNode(value, "schema", 0); err != nil {
		return nil, err
	}

	return &Schema{root: value}, nil
}

func checkSchemaNode(value any, path string, depth int) error {
	if depth >= maxSchemaDepth {
		return fmt.Errorf("%w: %s exceeds the schema depth limit", ErrInvalidSchema, path)
	}

	if _, ok := value.(bool); ok {
		return nil
	}

	node, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("%w: %s must be an object or boolean", ErrInvalidSchema, path)
	}

	for keyword, item := range node {
		if err := checkSchemaKeyword(keyword, item, path+"."+keyword, depth); err != nil {
			return err
		}
	}

	return nil
}

func checkSchemaKeyword(keyword string, value any, path string, depth int) error {
	switch keyword {
	case keywordType:
		return checkSchemaTypes(value, path)

	case keywordProperties:
		return checkSchemaProperties(value, path, depth)

	case keywordItems, keywordAdditionalProperties:
		return checkSchemaNode(value, path, depth+1)

	case keywordOneOf:
		return checkSchemaAlternatives(value, path, depth)

	case keywordRequired:
		return checkSchemaRequired(value, path)

	case keywordEnum:
		return checkSchemaEnum(value, path)

	case keywordMinimum, keywordMaximum:
		if _, ok := value.(json.Number); ok {
			return nil
		}

	case keywordMinItems, keywordMaxItems:
		return checkSchemaItemBound(value, path)

	case keywordUniqueItems:
		if _, ok := value.(bool); ok {
			return nil
		}

	default:
		return checkSchemaAnnotation(keyword, value, path)
	}

	return fmt.Errorf("%w: invalid value for %s", ErrInvalidSchema, path)
}

func checkSchemaItemBound(value any, path string) error {
	if raw, ok := value.(json.Number); ok {
		number := parseSchemaNumber(raw)
		if !number.negative && number.exponent.Sign() >= 0 {
			return nil
		}
	}

	return fmt.Errorf("%w: %s must be a nonnegative integer", ErrInvalidSchema, path)
}

func checkSchemaTypes(value any, path string) error {
	values := []any{value}
	if list, ok := value.([]any); ok {
		values = list
	}

	if len(values) == 0 {
		return fmt.Errorf("%w: %s needs at least one type", ErrInvalidSchema, path)
	}

	seen := set.New[string]()
	for _, item := range values {
		kind, ok := item.(string)
		if !ok || !knownSchemaType(kind) || seen.Contains(kind) {
			return fmt.Errorf("%w: invalid or duplicate type in %s", ErrInvalidSchema, path)
		}

		seen.Add(kind)
	}

	return nil
}

func checkSchemaProperties(value any, path string, depth int) error {
	properties, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("%w: %s must be an object", ErrInvalidSchema, path)
	}

	for name, child := range properties {
		if err := checkSchemaNode(child, path+"."+name, depth+1); err != nil {
			return err
		}
	}

	return nil
}

func checkSchemaAlternatives(value any, path string, depth int) error {
	children, ok := value.([]any)
	if !ok || len(children) == 0 {
		return fmt.Errorf("%w: %s must be a nonempty schema array", ErrInvalidSchema, path)
	}

	for index, child := range children {
		if err := checkSchemaNode(child, fmt.Sprintf("%s[%d]", path, index), depth+1); err != nil {
			return err
		}
	}

	return nil
}

func checkSchemaRequired(value any, path string) error {
	items, ok := value.([]any)
	if !ok {
		return fmt.Errorf("%w: %s must be a string array", ErrInvalidSchema, path)
	}

	seen := set.New[string]()
	for _, item := range items {
		name, ok := item.(string)
		if !ok || seen.Contains(name) {
			return fmt.Errorf("%w: invalid or duplicate property in %s", ErrInvalidSchema, path)
		}

		seen.Add(name)
	}

	return nil
}

func checkSchemaEnum(value any, path string) error {
	items, ok := value.([]any)
	if !ok || len(items) == 0 {
		return fmt.Errorf("%w: %s must be a nonempty array", ErrInvalidSchema, path)
	}

	if !uniqueJSONValues(items) {
		return fmt.Errorf("%w: %s contains duplicate values", ErrInvalidSchema, path)
	}

	return nil
}

func checkSchemaAnnotation(keyword string, value any, path string) error {
	switch keyword {
	case "$schema":
		if value == "https://json-schema.org/draft/2020-12/schema" ||
			value == "https://json-schema.org/draft/2020-12/schema#" {

			return nil
		}

	case "$comment", "title", "description", "contentEncoding", "contentMediaType", "x-mcp-header":
		if _, ok := value.(string); ok {
			return nil
		}

	case "readOnly", "writeOnly", "deprecated":
		if _, ok := value.(bool); ok {
			return nil
		}

	case "examples":
		if _, ok := value.([]any); ok {
			return nil
		}

	case "default":
		return nil

	default:
		return fmt.Errorf("%w: unsupported keyword %s", ErrInvalidSchema, path)
	}

	return fmt.Errorf("%w: invalid annotation %s", ErrInvalidSchema, path)
}
