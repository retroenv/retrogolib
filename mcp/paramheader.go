package mcp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/retroenv/retrogolib/set"
)

// Schema keys used by tool schema validation.
const (
	schemaKeyMCPHeader  = "x-mcp-header"
	schemaKeyProperties = "properties"
	schemaKeyType       = "type"
)

// headerParamPrefix is the prefix of the header names that mirror tool
// arguments on the streamable HTTP transport.
const headerParamPrefix = "Mcp-Param-"

// errNotObject reports a schema that is not a JSON object.
var errNotObject = errors.New("schema is not a JSON object")

// defaultInputSchema is the input schema of a tool that does not supply one.
var defaultInputSchema = json.RawMessage(`{"type":"object"}`)

// paramHeader is one tool argument that the client mirrors into an
// Mcp-Param-{name} HTTP header.
type paramHeader struct {
	name string
	path []string
}

// newParamHeader validates one x-mcp-header annotation.
func newParamHeader(schema map[string]any, value any, path []string, static bool) (paramHeader, error) {
	name, ok := value.(string)
	if !ok || !validToken(name) {
		return paramHeader{}, fmt.Errorf("%w: %s value %v is not a valid header name", ErrInvalidSchema,
			schemaKeyMCPHeader, value)
	}
	if !static || len(path) == 0 {
		return paramHeader{}, fmt.Errorf("%w: %s %q is not on a property that is statically reachable from the root",
			ErrInvalidSchema, schemaKeyMCPHeader, name)
	}

	switch schema[schemaKeyType] {
	case "string", "integer", "boolean":
	default:
		return paramHeader{}, fmt.Errorf("%w: %s %q must be on a string, integer, or boolean property",
			ErrInvalidSchema, schemaKeyMCPHeader, name)
	}

	return paramHeader{
		name: name,
		path: path,
	}, nil
}

// prepareToolSchemas sets the default input schema, validates the tool
// schemas, and returns the x-mcp-header annotations of the input schema.
func prepareToolSchemas(t *Tool) ([]paramHeader, error) {
	if len(t.InputSchema) == 0 {
		t.InputSchema = defaultInputSchema
	}

	root, err := decodeSchemaObject(t.InputSchema)
	if err != nil {
		return nil, fmt.Errorf("%w: input schema: %w", ErrInvalidSchema, err)
	}
	if root[schemaKeyType] != "object" {
		return nil, fmt.Errorf(`%w: input schema root must have "type": "object"`, ErrInvalidSchema)
	}

	if len(t.OutputSchema) > 0 {
		if _, err := decodeSchemaObject(t.OutputSchema); err != nil {
			return nil, fmt.Errorf("%w: output schema: %w", ErrInvalidSchema, err)
		}
	}

	var headers []paramHeader
	if err := collectParamHeaders(root, nil, true, &headers); err != nil {
		return nil, err
	}

	seen := set.New[string]()
	for _, header := range headers {
		key := strings.ToLower(header.name)
		if seen.Contains(key) {
			return nil, fmt.Errorf("%w: duplicate %s value %q", ErrInvalidSchema, schemaKeyMCPHeader, header.name)
		}
		seen.Add(key)
	}
	return headers, nil
}

// decodeSchemaObject decodes a JSON Schema that must be a JSON object.
func decodeSchemaObject(raw json.RawMessage) (map[string]any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()

	var schema map[string]any
	if err := dec.Decode(&schema); err != nil {
		return nil, fmt.Errorf("decoding schema: %w", err)
	}
	if schema == nil {
		return nil, errNotObject
	}
	return schema, nil
}

// collectParamHeaders finds the x-mcp-header annotations in a schema node. An
// annotation is valid only on a property that the root reaches through a
// chain of "properties" keys. The static flag records this condition.
func collectParamHeaders(node any, path []string, static bool, headers *[]paramHeader) error {
	switch n := node.(type) {
	case []any:
		for _, child := range n {
			if err := collectParamHeaders(child, nil, false, headers); err != nil {
				return err
			}
		}

	case map[string]any:
		return collectObjectHeaders(n, path, static, headers)
	}
	return nil
}

// collectObjectHeaders finds the x-mcp-header annotations in a schema object
// and in its children.
func collectObjectHeaders(node map[string]any, path []string, static bool, headers *[]paramHeader) error {
	if value, ok := node[schemaKeyMCPHeader]; ok {
		header, err := newParamHeader(node, value, path, static)
		if err != nil {
			return err
		}
		*headers = append(*headers, header)
	}

	for _, key := range slices.Sorted(maps.Keys(node)) {
		properties, ok := node[key].(map[string]any)
		if key != schemaKeyProperties || !static || !ok {
			if err := collectParamHeaders(node[key], nil, false, headers); err != nil {
				return err
			}
			continue
		}

		for _, name := range slices.Sorted(maps.Keys(properties)) {
			childPath := append(slices.Clone(path), name)
			if err := collectParamHeaders(properties[name], childPath, true, headers); err != nil {
				return err
			}
		}
	}
	return nil
}

// validToken reports whether s matches the HTTP field-name token syntax of
// RFC 9110.
func validToken(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case strings.ContainsRune("!#$%&'*+-.^_`|~", r):
		default:
			return false
		}
	}
	return true
}

// validateParamHeaders checks the Mcp-Param-{name} headers of a tools/call
// request against the argument values in the body. A header is required only
// when the arguments hold a value that is not null at the annotated path.
func validateParamHeaders(h http.Header, headers []paramHeader, arguments json.RawMessage) *Error {
	if len(headers) == 0 {
		return nil
	}

	var args map[string]any
	dec := json.NewDecoder(bytes.NewReader(arguments))
	dec.UseNumber()
	if err := dec.Decode(&args); err != nil {
		args = nil // the tool handler reports invalid arguments
	}

	for _, header := range headers {
		headerName := headerParamPrefix + header.name
		raw := h.Get(headerName)
		if raw != "" && !validHeaderValue(raw) {
			return headerMismatchError("malformed " + headerName + " header")
		}

		value, ok := lookupArgument(args, header.path)
		if !ok || value == nil {
			continue
		}
		if raw == "" {
			return headerMismatchError("missing " + headerName + " header")
		}
		if !headerValueMatches(decodeHeaderValue(raw), value) {
			return headerMismatchError(headerName + " header does not match the request body")
		}
	}
	return nil
}

// lookupArgument returns the argument value at a property path.
func lookupArgument(args map[string]any, path []string) (any, bool) {
	var value any = args
	for _, key := range path {
		object, ok := value.(map[string]any)
		if !ok {
			return nil, false
		}
		if value, ok = object[key]; !ok {
			return nil, false
		}
	}
	return value, true
}

// headerValueMatches compares a decoded header value with a body value.
// Numbers compare by value, so that "42.0" matches 42.
func headerValueMatches(header string, value any) bool {
	switch v := value.(type) {
	case string:
		return header == v

	case bool:
		return header == strconv.FormatBool(v)

	case json.Number:
		bodyNumber, err := v.Float64()
		if err != nil {
			return false
		}
		headerNumber, err := strconv.ParseFloat(header, 64)
		return err == nil && headerNumber == bodyNumber

	default:
		return false
	}
}

// headerMismatchError builds a HeaderMismatch protocol error.
func headerMismatchError(detail string) *Error {
	return &Error{
		Code:    CodeHeaderMismatch,
		Message: "Header mismatch: " + detail,
	}
}
