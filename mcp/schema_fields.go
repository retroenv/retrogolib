package mcp

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/retroenv/retrogolib/mcp/jsonschema"
	"github.com/retroenv/retrogolib/set"
)

type schemaField struct {
	field     reflect.StructField
	name      string
	depth     int
	tagged    bool
	omitEmpty bool
}

type schemaFieldSelection struct {
	index     int
	depth     int
	tagged    bool
	ambiguous bool
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
			fields = append(fields, children...)

			continue
		}
		tagged := name != ""
		if name == "" {
			name = field.Name
		}

		fields = append(fields, schemaField{
			field:     field,
			name:      name,
			depth:     depth,
			tagged:    tagged,
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
				index:  index,
				depth:  field.depth,
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
