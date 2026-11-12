package mcp

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/mcp/jsonschema"
)

// schemaFor generates the schema for a type and decodes it into a map.
func schemaFor(t *testing.T, typ any) map[string]any {
	t.Helper()

	raw, err := SchemaFor(reflect.TypeOf(typ))
	assert.NoError(t, err)

	var schema map[string]any
	assert.NoError(t, json.Unmarshal(raw, &schema))
	return schema
}

func TestSchemaKinds(t *testing.T) {
	type allKinds struct {
		Text    string          `json:"text"`
		Enabled bool            `json:"enabled"`
		Count   int             `json:"count"`
		Byte    uint8           `json:"byte"`
		Ratio   float64         `json:"ratio"`
		List    []string        `json:"list"`
		Pairs   map[string]int  `json:"pairs"`
		Nested  struct{ X int } `json:"nested"`
		Ptr     *int            `json:"ptr"`
		Raw     json.RawMessage `json:"raw"`
		Any     any             `json:"any"`
	}

	schema := schemaFor(t, allKinds{})
	properties := schema["properties"].(map[string]any)

	assert.Equal(t, map[string]any{"type": "string"}, properties["text"])
	assert.Equal(t, map[string]any{"type": "boolean"}, properties["enabled"])
	assert.Equal(t, map[string]any{"type": "integer"}, properties["count"])
	assert.Equal(t, map[string]any{"type": "integer"}, properties["byte"])
	assert.Equal(t, map[string]any{"type": "number"}, properties["ratio"])
	assert.Equal(t, map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, properties["list"])
	assert.Equal(t, map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "integer"}}, properties["pairs"])
	assert.Equal(t, map[string]any{
		"additionalProperties": false,
		"properties":           map[string]any{"X": map[string]any{"type": "integer"}},
		"required":             []any{"X"},
		"type":                 "object",
	}, properties["nested"])
	assert.Equal(t, map[string]any{"type": "integer"}, properties["ptr"])
	assert.Equal(t, map[string]any{}, properties["raw"])
	assert.Equal(t, map[string]any{}, properties["any"])
}

func TestSchemaRequiredRules(t *testing.T) {
	type rules struct {
		Plain    int    `json:"plain"`
		Optional *int   `json:"optional"`
		Omit     string `json:"omit,omitempty"`
		Relaxed  string `json:"relaxed" required:"false"`
		Renamed  string `json:"renamed"`
		Skipped  string `json:"-"`
	}

	schema := schemaFor(t, rules{})
	assert.Equal(t, []any{"plain", "renamed"}, schema["required"])

	properties := schema["properties"].(map[string]any)
	assert.Len(t, properties, 5)
	assert.NotNil(t, properties["omit"])
	assert.NotNil(t, properties["relaxed"])
	assert.Nil(t, properties["Skipped"])
}

func TestSchemaKeywords(t *testing.T) {
	type keywords struct {
		Level int    `json:"level" description:"severity level" enum:"0,1,2" minimum:"0" maximum:"2"`
		Name  string `json:"name" description:"display name" enum:"small,large"`
	}

	schema := schemaFor(t, keywords{})
	properties := schema["properties"].(map[string]any)

	assert.Equal(t, map[string]any{
		"type":        "integer",
		"description": "severity level",
		"enum":        []any{float64(0), float64(1), float64(2)},
		"minimum":     float64(0),
		"maximum":     float64(2),
	}, properties["level"])

	assert.Equal(t, map[string]any{
		"type":        "string",
		"description": "display name",
		"enum":        []any{"small", "large"},
	}, properties["name"])
}

type customSchemaType uint16

// JSONSchema supplies the schema for the type.
func (c customSchemaType) JSONSchema() json.RawMessage {
	return json.RawMessage(`{"type":["integer","string"],"pattern":"^\\$[0-9A-F]{1,4}$"}`)
}

func TestSchemaCustomHook(t *testing.T) {
	type withCustom struct {
		Addr customSchemaType `json:"addr"`
	}

	schema := schemaFor(t, withCustom{})
	properties := schema["properties"].(map[string]any)

	assert.Equal(t, map[string]any{
		"type":    []any{"integer", "string"},
		"pattern": `^\$[0-9A-F]{1,4}$`,
	}, properties["addr"])
}

func TestSchemaRecursiveType(t *testing.T) {
	type node struct {
		Next *node `json:"next"`
	}

	_, err := SchemaFor(reflect.TypeOf(node{}))
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
}

func TestSchemaNonStructInput(t *testing.T) {
	_, err := SchemaFor(reflect.TypeOf(42))
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
}

func TestSchemaMapWithNonStringKey(t *testing.T) {
	type withMap struct {
		Lookup map[int]string `json:"lookup"`
	}

	_, err := SchemaFor(reflect.TypeOf(withMap{}))
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
}

func TestSchemaEmbeddedField(t *testing.T) {
	type inner struct {
		Value int `json:"value"`
	}
	type outer struct {
		inner
		Name string `json:"name"`
	}

	schema := schemaFor(t, outer{})
	properties := schema["properties"].(map[string]any)

	assert.Equal(t, map[string]any{"type": "integer"}, properties["value"])
	assert.Equal(t, map[string]any{"type": "string"}, properties["name"])
	assert.Equal(t, []any{"value", "name"}, schema["required"])
}

func TestSchemaEmbeddedVariants(t *testing.T) {
	type optional struct {
		Note string `json:"note,omitempty"`
	}
	type Base struct {
		ID int `json:"id"`
	}
	type Label string
	type outer struct {
		optional
		*Base
		Label
		Name string `json:"name"`
	}

	schema := schemaFor(t, outer{})
	properties := schema["properties"].(map[string]any)

	assert.Equal(t, map[string]any{"type": "string"}, properties["note"])
	assert.Equal(t, map[string]any{"type": "integer"}, properties["id"])
	assert.Equal(t, map[string]any{"type": "string"}, properties["Label"])
	assert.Equal(t, []any{"id", "Label", "name"}, schema["required"])
}

type textValue struct{}

// MarshalText encodes the value as text.
func (textValue) MarshalText() ([]byte, error) {
	return []byte("text"), nil
}

func TestSchemaSpecialTypes(t *testing.T) {
	type special struct {
		Bytes  []byte                       `json:"bytes"`
		Text   textValue                    `json:"text"`
		Addrs  []customSchemaType           `json:"addrs"`
		Lookup map[string]*customSchemaType `json:"lookup"`
		Addr   customSchemaType             `json:"addr" description:"CPU address" header:"Addr"`
	}

	schema := schemaFor(t, special{})
	properties := schema["properties"].(map[string]any)
	addrSchema := map[string]any{
		"type":    []any{"integer", "string"},
		"pattern": `^\$[0-9A-F]{1,4}$`,
	}

	assert.Equal(t, map[string]any{"type": "string", "contentEncoding": "base64"}, properties["bytes"])
	assert.Equal(t, map[string]any{"type": "string"}, properties["text"])
	assert.Equal(t, map[string]any{"type": "array", "items": addrSchema}, properties["addrs"])
	assert.Equal(t, map[string]any{"type": "object", "additionalProperties": addrSchema}, properties["lookup"])
	assert.Equal(t, map[string]any{
		"type":         []any{"integer", "string"},
		"pattern":      `^\$[0-9A-F]{1,4}$`,
		"description":  "CPU address",
		"x-mcp-header": "Addr",
	}, properties["addr"])
}

func TestSchemaRequiredTrue(t *testing.T) {
	type rules struct {
		Pointer *int   `json:"pointer" required:"true"`
		Omit    string `json:"omit,omitempty" required:"true"`
	}

	schema := schemaFor(t, rules{})
	assert.Equal(t, []any{"pointer", "omit"}, schema["required"])
}

func TestSchemaInvalidTags(t *testing.T) {
	type badEnum struct {
		Level int `json:"level" enum:"0,high"`
	}
	type badMinimum struct {
		Level int `json:"level" minimum:"low"`
	}

	for _, typ := range []reflect.Type{reflect.TypeFor[badEnum](), reflect.TypeFor[badMinimum]()} {
		_, err := SchemaFor(typ)
		assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
	}
}

func TestSchemaRecursiveSlice(t *testing.T) {
	type node struct {
		Children []node `json:"children"`
	}

	_, err := SchemaFor(reflect.TypeFor[node]())
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
}

func TestSchemaEmbeddedFieldSelection(t *testing.T) {
	type Base struct {
		Value string `json:"value"`
	}
	type Shadow struct {
		Base
		Value int `json:"value"`
	}
	type ReverseShadow struct {
		Value int `json:"value"`
		Base
	}
	type Other struct {
		Value int `json:"value"`
	}
	type Untagged struct{ Value int }
	type AnotherUntagged struct{ Value int }
	type Conflict struct {
		Untagged
		AnotherUntagged
	}
	type Tagged struct {
		Different string `json:"Value"`
	}
	type TaggedWins struct {
		Untagged
		Tagged
	}
	type Deep struct{ Other }
	type ShallowWins struct {
		Deep
		Base
	}
	for _, test := range []struct {
		name       string
		value      any
		properties map[string]any
		required   []any
	}{
		{"outer field", Shadow{}, map[string]any{"value": map[string]any{"type": "integer"}}, []any{"value"}},
		{"outer field before embedding", ReverseShadow{}, map[string]any{"value": map[string]any{"type": "integer"}}, []any{"value"}},
		{"ambiguous fields", Conflict{}, map[string]any{}, nil},
		{"tag wins", TaggedWins{}, map[string]any{"Value": map[string]any{"type": "string"}}, []any{"Value"}},
		{"shallow wins", ShallowWins{}, map[string]any{"value": map[string]any{"type": "string"}}, []any{"value"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			schema := schemaFor(t, test.value)
			assert.Equal(t, test.properties, schema["properties"])
			if test.required == nil {
				assert.Nil(t, schema["required"])
			} else {
				assert.Equal(t, test.required, schema["required"])
			}
			raw, err := json.Marshal(schema)
			assert.NoError(t, err)
			_, err = CompileSchema(raw)
			assert.NoError(t, err)
		})
	}
}

func TestSchemaQuotedFieldTypes(t *testing.T) {
	type fields struct {
		Count   int    `json:"count,string" enum:"42,43"`
		Enabled bool   `json:"enabled,string" enum:"true,false"`
		Text    string `json:"text,string" enum:"hello"`
	}
	properties := schemaFor(t, fields{})["properties"].(map[string]any)
	assert.Equal(t, map[string]any{"enum": []any{"42", "43"}, "type": "string"}, properties["count"])
	assert.Equal(t, map[string]any{"enum": []any{"true", "false"}, "type": "string"}, properties["enabled"])
	assert.Equal(t, map[string]any{"enum": []any{`"hello"`}, "type": "string"}, properties["text"])
}

func TestSchemaRejectsQuotedNumericBounds(t *testing.T) {
	type minimum struct {
		Count int `json:"count,string" minimum:"1"`
	}
	type maximum struct {
		Count int `json:"count,string" maximum:"1"`
	}
	for _, typ := range []reflect.Type{reflect.TypeFor[minimum](), reflect.TypeFor[maximum]()} {
		_, err := SchemaFor(typ)
		assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
	}
}
