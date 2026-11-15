package mcp

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/mcp/jsonschema"
)

type nullableOutputValue string

func TestOutputSchemaMatchesNilAndNestedJSONValues(t *testing.T) {
	type child struct {
		Values []int `json:"values"`
	}
	type result struct {
		Array   [2]int            `json:"array"`
		Bytes   []byte            `json:"bytes"`
		Child   *child            `json:"child"`
		Entries map[string]*child `json:"entries"`
		Matrix  [][]*int          `json:"matrix"`
		Quoted  *int              `json:"quoted,string" enum:"42"`
		Values  []int             `json:"values"`
	}
	value := 42
	for _, out := range []result{
		{},
		{
			Bytes: []byte{1, 2},
			Child: &child{},
			Entries: map[string]*child{
				"nil": nil,
				"set": {},
			},
			Matrix: [][]*int{nil, {nil, &value}},
			Quoted: &value,
			Values: []int{1, 2},
		},
	} {
		schema, err := OutputSchemaFor(reflect.TypeFor[result]())
		assert.NoError(t, err)
		validator, err := CompileSchema(schema)
		assert.NoError(t, err)
		data, err := json.Marshal(out)
		assert.NoError(t, err)
		assert.NoError(t, validator.Validate(data))
	}
	input, err := SchemaFor(reflect.TypeFor[result]())
	assert.NoError(t, err)
	validator, err := CompileSchema(input)
	assert.NoError(t, err)
	data, err := json.Marshal(result{})
	assert.NoError(t, err)
	assert.Error(t, validator.Validate(data))
}

func TestOutputSchemaKeepsArraysAndScalarConstraints(t *testing.T) {
	type result struct {
		Array  [1]int `json:"array"`
		Choice *int   `json:"choice" enum:"42"`
	}
	schema, err := OutputSchemaFor(reflect.TypeFor[result]())
	assert.NoError(t, err)
	validator, err := CompileSchema(schema)
	assert.NoError(t, err)
	assert.NoError(t, validator.Validate(json.RawMessage(`{"array":[0],"choice":null}`)))
	assert.Error(t, validator.Validate(json.RawMessage(`{"array":null,"choice":null}`)))
	assert.Error(t, validator.Validate(json.RawMessage(`{"array":[0],"choice":43}`)))
}

func TestOutputSchemaPreservesExactCustomNumbers(t *testing.T) {
	type result struct {
		Value *exactSchemaValue `json:"value"`
	}
	schema, err := OutputSchemaFor(reflect.TypeFor[result]())
	assert.NoError(t, err)
	validator, err := CompileSchema(schema)
	assert.NoError(t, err)
	assert.NoError(t, validator.Validate(json.RawMessage(`{"value":null}`)))
	assert.NoError(t, validator.Validate(json.RawMessage(`{"value":9007199254740993}`)))
	assert.Error(t, validator.Validate(json.RawMessage(`{"value":9007199254740992}`)))
}

func TestOutputSchemaObjectRootAndExistingNullBranch(t *testing.T) {
	schema, err := OutputSchemaFor(reflect.TypeFor[map[string][]int]())
	assert.NoError(t, err)
	validator, err := CompileSchema(schema)
	assert.NoError(t, err)
	assert.NoError(t, validator.Validate(json.RawMessage(`{"values":null}`)))
	assert.Error(t, validator.Validate(json.RawMessage(`null`)))
	_, err = OutputSchemaFor(reflect.TypeFor[int]())
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
	_, err = OutputSchemaFor(nil)
	assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)

	type result struct {
		Value *nullableOutputValue `json:"value"`
	}
	schema, err = OutputSchemaFor(reflect.TypeFor[result]())
	assert.NoError(t, err)
	validator, err = CompileSchema(schema)
	assert.NoError(t, err)
	assert.NoError(t, validator.Validate(json.RawMessage(`{"value":null}`)))
}

func TestOutputSchemaAllowsOmittedEmbeddedPointerFields(t *testing.T) {
	type embedded struct {
		Count int `json:"count"`
	}
	type result struct {
		*embedded
	}
	schema, err := OutputSchemaFor(reflect.TypeFor[result]())
	assert.NoError(t, err)
	validator, err := CompileSchema(schema)
	assert.NoError(t, err)
	for _, value := range []result{
		{},
		{
			embedded: &embedded{
				Count: 42,
			},
		},
	} {
		data, err := json.Marshal(value)
		assert.NoError(t, err)
		assert.NoError(t, validator.Validate(data))
	}
}

func TestSchemaRejectsRecursiveContainers(t *testing.T) {
	// Container cycles must return an error before the stack is exhausted.
	type recursiveMap map[string]recursiveMap
	type recursiveSlice []recursiveSlice
	type recursiveArray [1]*recursiveArray
	type recursivePointer *recursivePointer
	type recursiveMapSlice map[string][]recursiveMapSlice

	for _, typ := range []reflect.Type{
		reflect.TypeFor[recursiveMap](),
		reflect.TypeFor[recursiveSlice](),
		reflect.TypeFor[recursiveArray](),
		reflect.TypeFor[recursivePointer](),
		reflect.TypeFor[recursiveMapSlice](),
	} {
		t.Run(typ.String(), func(t *testing.T) {
			_, err := OutputSchemaFor(typ)
			assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)

			wrapped := reflect.StructOf([]reflect.StructField{
				{
					Name: "Value",
					Type: typ,
				},
			})
			_, err = SchemaFor(wrapped)
			assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
			_, err = OutputSchemaFor(wrapped)
			assert.ErrorIs(t, err, jsonschema.ErrInvalidSchema)
		})
	}
}

func (nullableOutputValue) JSONSchema() json.RawMessage {
	return json.RawMessage(`{"oneOf":[{"type":"string"},{"type":"null"}]}`)
}
