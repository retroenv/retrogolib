package mcp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/retroenv/retrogolib/mcp/jsonschema"
	"github.com/retroenv/retrogolib/set"
)

type schemaField struct {
	depth  int
	field  reflect.StructField
	name   string
	tagged bool

	indirect  bool
	omitEmpty bool
}

type schemaFieldSelection struct {
	ambiguous bool
	depth     int
	index     int
	tagged    bool
}

func (b schemaBuilder) appendStructFields(t reflect.Type, properties map[string]any, required *[]string) error {
	fields, err := schemaFields(t, set.New[reflect.Type](), 0)
	if err != nil {
		return err
	}
	selected := selectSchemaFields(fields)
	for index, field := range fields {
		if selected[field.name] != index {
			continue
		}
		schema, err := b.fieldSchema(field.field)
		if err != nil {
			return err
		}
		properties[field.name] = schema
		if b.required(field) {
			*required = append(*required, field.name)
		}
	}

	return nil
}

func (b schemaBuilder) required(field schemaField) bool {
	if b.output && field.indirect && field.field.Tag.Get("required") != "true" {
		return false
	}
	return fieldRequired(field.field, field.omitEmpty)
}

// fieldSchema builds the schema for one struct field and adds the keyword
// tags.
func (b schemaBuilder) fieldSchema(field reflect.StructField) (map[string]any, error) {
	schema, err := b.nonNullSchema(field.Type)
	if err != nil {
		return nil, err
	}
	if err := applyKeywords(schema, field); err != nil {
		return nil, fmt.Errorf("%w: field %s: %w", jsonschema.ErrInvalidSchema, field.Name, err)
	}
	if slices.Contains(strings.Split(field.Tag.Get("json"), ",")[1:], "string") {
		quoted, err := quotedFieldSchema(schema, field.Type)
		if err != nil {
			return nil, fmt.Errorf("%w: field %s: %w", jsonschema.ErrInvalidSchema, field.Name, err)
		}
		return b.allowNull(field.Type, quoted)
	}
	return b.allowNull(field.Type, schema)
}

func schemaFields(t reflect.Type, seen set.Set[reflect.Type], depth int) ([]schemaField, error) {
	if seen.Contains(t) {
		return nil, fmt.Errorf("%w: recursive type %s", jsonschema.ErrInvalidSchema, t)
	}
	seen.Add(t)
	defer seen.Remove(t)

	var fields []schemaField
	for index := range t.NumField() {
		field := t.Field(index)
		embedded := field.Type
		if embedded.Kind() == reflect.Pointer {
			embedded = embedded.Elem()
		}

		if !field.IsExported() && (!field.Anonymous || embedded.Kind() != reflect.Struct) {
			continue
		}

		name, omitEmpty, skip := parseJSONTag(field.Tag.Get("json"))
		if skip {
			continue
		}
		if field.Anonymous && name == "" && embedded.Kind() == reflect.Struct {
			children, err := schemaFields(embedded, seen, depth+1)
			if err != nil {
				return nil, err
			}
			if field.Type.Kind() == reflect.Pointer {
				for index := range children {
					children[index].indirect = true
				}
			}

			fields = append(fields, children...)

			continue
		}
		tagged := name != ""
		if name == "" {
			name = field.Name
		}

		fields = append(fields, schemaField{
			depth:  depth,
			field:  field,
			name:   name,
			tagged: tagged,

			omitEmpty: omitEmpty,
		})
	}

	return fields, nil
}

// selectSchemaFields selects one field per JSON name. Shallower fields take
// precedence, then fields with a JSON name tag. Fields that still tie are excluded.
func selectSchemaFields(fields []schemaField) map[string]int {
	selections := make(map[string]schemaFieldSelection)
	for index, field := range fields {
		best, exists := selections[field.name]
		switch {
		case !exists || field.depth < best.depth || (field.depth == best.depth && field.tagged && !best.tagged):
			selections[field.name] = schemaFieldSelection{
				depth:  field.depth,
				index:  index,
				tagged: field.tagged,
			}

		case field.depth == best.depth && field.tagged == best.tagged:
			best.ambiguous = true
			selections[field.name] = best
		}
	}

	selected := make(map[string]int, len(selections))
	for name, best := range selections {
		index := best.index
		if best.ambiguous {
			index = -1
		}

		selected[name] = index
	}

	return selected
}

func quotedFieldSchema(schema map[string]any, t reflect.Type) (map[string]any, error) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.String:
	default:
		return schema, nil
	}

	if _, ok := schema["minimum"]; ok {
		return nil, fmt.Errorf("%w: minimum is not supported with the JSON string option", jsonschema.ErrInvalidSchema)
	}
	if _, ok := schema["maximum"]; ok {
		return nil, fmt.Errorf("%w: maximum is not supported with the JSON string option", jsonschema.ErrInvalidSchema)
	}
	schema["type"] = "string"
	if values, ok := schema["enum"].([]any); ok {
		for index, value := range values {
			if t.Kind() == reflect.Bool {
				values[index] = fmt.Sprint(value)
				continue
			}
			encoded, err := json.Marshal(value)
			if err != nil {
				return nil, fmt.Errorf("%w: encoding quoted enum: %w", jsonschema.ErrInvalidSchema, err)
			}
			values[index] = string(encoded)
		}
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
		value, err := schemaTagNumber(tag)
		if err != nil {
			return fmt.Errorf("%s tag %q is not a number: %w", keyword, tag, err)
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

		number, err := schemaTagNumber(value)
		if err != nil {
			return nil, fmt.Errorf("enum value %q is not a number: %w", value, err)
		}
		result = append(result, number)
	}

	return result, nil
}

// schemaTagNumber keeps the exact decimal text from a numeric field tag.
func schemaTagNumber(text string) (json.Number, error) {
	value, err := jsonschema.DecodeValue(json.RawMessage(text))
	if err != nil {
		return "", fmt.Errorf("reading numeric tag: %w", err)
	}

	number, ok := value.(json.Number)
	if !ok {
		return "", fmt.Errorf("%w: expected a JSON number", jsonschema.ErrInvalidSchema)
	}

	return number, nil
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
