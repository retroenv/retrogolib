package jsonschema

import (
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

var compiledSchemaCases = []struct {
	name    string
	schema  string
	valid   []string
	invalid []string
}{
	{
		name:    "closed object and required fields",
		schema:  `{"type":"object","properties":{"name":{"type":"string"}},"required":["name"],"additionalProperties":false}`,
		valid:   []string{`{"name":"test"}`},
		invalid: []string{`{}`, `{"name":null}`, `{"name":1}`, `{"name":"test","extra":true}`, `null`, `[]`},
	},
	{
		name:    "schema for additional properties",
		schema:  `{"type":"object","additionalProperties":{"type":"integer"}}`,
		valid:   []string{`{}`, `{"a":1,"b":2}`},
		invalid: []string{`{"a":"one"}`, `{"a":null}`},
	},
	{
		name:    "type union",
		schema:  `{"type":["integer","string"]}`,
		valid:   []string{`1`, `1.0`, `"0x100"`},
		invalid: []string{`null`, `true`, `1.5`, `[]`, `{}`},
	},
	{
		name:    "nullable boolean",
		schema:  `{"type":["boolean","null"]}`,
		valid:   []string{`true`, `false`, `null`},
		invalid: []string{`0`, `"true"`},
	},
	{
		name:    "integer precision above float64 range",
		schema:  `{"type":"integer","minimum":9007199254740993,"maximum":9007199254740993}`,
		valid:   []string{`9007199254740993`, `9007199254740993.0`},
		invalid: []string{`9007199254740992`, `9007199254740994`, `9007199254740993.00000000001`},
	},
	{
		name:    "uint64 bounds",
		schema:  `{"type":"integer","minimum":0,"maximum":18446744073709551615}`,
		valid:   []string{`0`, `-0`, `18446744073709551615`},
		invalid: []string{`-1`, `18446744073709551616`, `1.0000000000000000000001`},
	},
	{
		name:    "decimal bounds",
		schema:  `{"type":"number","minimum":-1.20,"maximum":0.15}`,
		valid:   []string{`-1.2`, `-1.1999`, `0`, `0.1`, `0.15`, `15e-2`, `1e-1000000000`},
		invalid: []string{`-1.20001`, `0.150000000000000000001`, `1e1000000000`},
	},
	{
		name:    "integer exponents do not expand powers",
		schema:  `{"type":"integer"}`,
		valid:   []string{`1e1000000000`, `0e-1000000000`, `100e-2`, `1.00e+000`},
		invalid: []string{`1e-1000000000`, `10e-2`},
	},
	{
		name:    "enum uses JSON numeric equality",
		schema:  `{"enum":[1,{"a":[2]}]}`,
		valid:   []string{`1.0`, `1e0`, `{"a":[2.00]}`},
		invalid: []string{`"1"`, `1.01`, `{"a":[3]}`},
	},
	{
		name:    "array bounds and item enum",
		schema:  `{"type":"array","minItems":1,"maxItems":3,"uniqueItems":true,"items":{"enum":["mapper","system","menu"]}}`,
		valid:   []string{`["mapper"]`, `["mapper","system","menu"]`},
		invalid: []string{`[]`, `["cpu"]`, `["mapper","mapper"]`, `[null]`, `["mapper","system","menu","mapper"]`},
	},
	{
		name:    "array constraints without items",
		schema:  `{"type":"array","minItems":1,"maxItems":2,"uniqueItems":true}`,
		valid:   []string{`[1,"1"]`, `[{"a":1},{"a":2}]`},
		invalid: []string{`[]`, `[1,2,3]`, `[1,1.0]`, `[{"a":1,"b":[2]},{"b":[2.0],"a":1e0}]`},
	},
	{
		name:    "oneOf counts zero and multiple matches",
		schema:  `{"type":"object","oneOf":[{"required":["ok"]},{"required":["error"]}]}`,
		valid:   []string{`{"ok":true}`, `{"error":"test"}`},
		invalid: []string{`{}`, `{"ok":true,"error":"test"}`, `null`},
	},
	{
		name:    "oneOf overlapping numeric types",
		schema:  `{"oneOf":[{"type":"integer"},{"type":"number"}]}`,
		valid:   []string{`1.5`},
		invalid: []string{`1`, `"one"`},
	},
	{
		name:    "boolean property schema",
		schema:  `{"type":"object","properties":{"blocked":false,"free":true}}`,
		valid:   []string{`{}`, `{"free":null}`},
		invalid: []string{`{"blocked":null}`},
	},
}

func TestCompiledSchemaValues(t *testing.T) {
	for _, test := range compiledSchemaCases {
		t.Run(test.name, func(t *testing.T) {
			schema, err := Compile(json.RawMessage(test.schema))
			assert.NoError(t, err)
			for _, raw := range test.valid {
				t.Run("valid "+raw, func(t *testing.T) { assert.NoError(t, schema.Validate(json.RawMessage(raw))) })
			}
			for _, raw := range test.invalid {
				t.Run("invalid "+raw, func(t *testing.T) { assert.Error(t, schema.Validate(json.RawMessage(raw))) })
			}
		})
	}
}

func TestCompileRejectsUnsupportedAndMalformedSchemas(t *testing.T) {
	for _, raw := range []string{
		`null`, `[]`, `1`, `{} {}`, `{"type":"object"} garbage`,
		`{"pattern":".*"}`, `{"$ref":"https://example.invalid/schema"}`,
		`{"properties":{"name":{"type":"string","pattern":".*"}}}`,
		`{"$schema":"http://json-schema.org/draft-07/schema#"}`,
		`{"type":"unknown"}`, `{"type":[]}`, `{"type":["string","string"]}`, `{"type":null}`,
		`{"properties":[]}`, `{"required":null}`, `{"required":[1]}`, `{"required":["a","a"]}`,
		`{"minimum":"0"}`, `{"maximum":null}`, `{"minItems":-1}`, `{"maxItems":1.5}`,
		`{"uniqueItems":1}`, `{"items":[]}`, `{"additionalProperties":null}`,
		`{"oneOf":[]}`, `{"oneOf":[null]}`, `{"enum":[]}`, `{"enum":[1,1.0]}`,
		`{"description":null}`, `{"examples":false}`, `{"readOnly":"true"}`,
	} {
		t.Run(raw, func(t *testing.T) {
			_, err := Compile(json.RawMessage(raw))
			assert.ErrorIs(t, err, ErrInvalidSchema)
		})
	}
	deep := strings.Repeat(`{"items":`, maxSchemaDepth) + `true` + strings.Repeat(`}`, maxSchemaDepth)
	_, err := Compile(json.RawMessage(deep))
	assert.ErrorIs(t, err, ErrInvalidSchema)
}

func TestSchemaValidatesOneCompleteJSONValue(t *testing.T) {
	schema, err := Compile(json.RawMessage(`true`))
	assert.NoError(t, err)
	for _, raw := range []string{`null`, `true`, `0`, `"text"`, `[]`, `{}`} {
		assert.NoError(t, schema.Validate(json.RawMessage(raw)))
	}
	for _, raw := range []string{``, `null null`, `{} garbage`, `NaN`, `1e`, `{"a":}`} {
		assert.Error(t, schema.Validate(json.RawMessage(raw)))
	}
	denied, err := Compile(json.RawMessage(`false`))
	assert.NoError(t, err)
	assert.Error(t, denied.Validate(json.RawMessage(`null`)))
}

func TestCompiledSchemaOwnsItsDataAndSupportsConcurrentValidation(t *testing.T) {
	raw := json.RawMessage(`{"type":"integer","minimum":0,"maximum":100}`)
	schema, err := Compile(raw)
	assert.NoError(t, err)
	for index := range raw {
		raw[index] = ' '
	}
	var work sync.WaitGroup
	for range 8 {
		work.Go(func() {
			assert.NoError(t, schema.Validate(json.RawMessage(`50`)))
			assert.Error(t, schema.Validate(json.RawMessage(`101`)))
		})
	}
	work.Wait()
}
