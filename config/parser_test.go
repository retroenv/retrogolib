package config

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestParser_BasicParsing(t *testing.T) {
	data := `# Global comment
[section1]
key1 = "value1"
key2 = 42

[section2]
# Section comment
key3 = true`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Check sections exist
	assert.NotNil(t, config.sections["section1"])
	assert.NotNil(t, config.sections["section2"])

	// Check values
	section1 := config.sections["section1"]
	assert.Equal(t, "value1", section1["key1"].parsed)
	assert.Equal(t, 42, section1["key2"].parsed)

	section2 := config.sections["section2"]
	assert.True(t, section2["key3"].parsed.(bool))

	// Check comments
	assert.Len(t, config.comments, 2)
	assert.Equal(t, "Global comment", config.comments[0].Text)
	assert.Equal(t, "", config.comments[0].Section) // Global comment
	assert.Equal(t, "Section comment", config.comments[1].Text)
	assert.Equal(t, "section2", config.comments[1].Section)
}

func TestParser_EmptyLines(t *testing.T) {
	data := `
[section]

key = value

# Comment

`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Should handle empty lines gracefully
	assert.NotNil(t, config.sections["section"])
	assert.Equal(t, "value", config.sections["section"]["key"].parsed)

	// Check structure includes empty lines
	emptyLines := 0
	for _, element := range config.structure {
		if element.Type == emptyLineElement {
			emptyLines++
		}
	}
	assert.Greater(t, emptyLines, 0)
}

func TestParser_QuotedStrings(t *testing.T) {
	data := `[strings]
quoted = "hello world"
with_escapes = "line1\nline2\ttab"
with_quotes = "say \"hello\""`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	section := config.sections["strings"]
	assert.Equal(t, "hello world", section["quoted"].parsed)
	assert.Equal(t, "line1\nline2\ttab", section["with_escapes"].parsed)
	assert.Equal(t, "say \"hello\"", section["with_quotes"].parsed)
}

func TestParser_HexValues(t *testing.T) {
	data := `[hex]
lowercase = 0xff00
uppercase = 0xFF00
mixed = 0xAbCd`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	section := config.sections["hex"]
	assert.Equal(t, hexType, section["lowercase"].vtype)
	assert.Equal(t, 0xff00, section["lowercase"].parsed)
	assert.Equal(t, hexType, section["uppercase"].vtype)
	assert.Equal(t, 0xFF00, section["uppercase"].parsed)
	assert.Equal(t, hexType, section["mixed"].vtype)
	assert.Equal(t, 0xAbCd, section["mixed"].parsed)
}

func TestParser_FloatValues(t *testing.T) {
	data := `[floats]
simple = 3.14
scientific = 1.23e-4
zero = 0.0
negative = -2.5`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	section := config.sections["floats"]
	assert.Equal(t, floatType, section["simple"].vtype)
	assert.Equal(t, 3.14, section["simple"].parsed)
	assert.Equal(t, floatType, section["scientific"].vtype)
	assert.Equal(t, 1.23e-4, section["scientific"].parsed)
	assert.Equal(t, floatType, section["zero"].vtype)
	assert.Equal(t, 0.0, section["zero"].parsed)
	assert.Equal(t, floatType, section["negative"].vtype)
	assert.Equal(t, -2.5, section["negative"].parsed)
}

func TestParser_BooleanValues(t *testing.T) {
	data := `[booleans]
true_val = true
false_val = false`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	section := config.sections["booleans"]
	assert.Equal(t, boolType, section["true_val"].vtype)
	assert.True(t, section["true_val"].parsed.(bool))
	assert.Equal(t, boolType, section["false_val"].vtype)
	assert.False(t, section["false_val"].parsed.(bool))
}

func TestParser_InvalidSyntax(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"empty section name", "[]\nkey = value"},
		{"invalid key-value", "[section]\ninvalid line"},
		{"empty key", "[section]\n = value"},
		{"invalid quoted string", `[section]
key = "unterminated`},
		{"invalid hex", "[section]\nkey = 0xGGGG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadConfigBytes([]byte(tt.data))
			assert.Error(t, err)
		})
	}
}

func TestParser_StructureTracking(t *testing.T) {
	data := `# Comment 1
[section1]
# Comment 2
key1 = value1

# Comment 3
[section2]
key2 = value2`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Check structure elements are tracked in order
	assert.GreaterOrEqual(t, len(config.structure), 8)

	// Verify first few elements
	assert.Equal(t, commentElement, config.structure[0].Type)
	assert.Equal(t, "# Comment 1", config.structure[0].Content)

	assert.Equal(t, sectionElement, config.structure[1].Type)
	assert.Equal(t, "[section1]", config.structure[1].Content)

	assert.Equal(t, commentElement, config.structure[2].Type)
	assert.Equal(t, "# Comment 2", config.structure[2].Content)
	assert.Equal(t, "section1", config.structure[2].Section)

	assert.Equal(t, keyValueElement, config.structure[3].Type)
	assert.Equal(t, "key1", config.structure[3].Key)
	assert.Equal(t, "section1", config.structure[3].Section)
}

func TestParser_WhitespaceHandling(t *testing.T) {
	data := `   # Comment with leading spaces
  [  section  ]  
   key1   =   value1   
	key2	=	"value2"	`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Should handle whitespace correctly
	assert.NotNil(t, config.sections["section"])
	section := config.sections["section"]
	assert.Equal(t, "value1", section["key1"].parsed)
	assert.Equal(t, "value2", section["key2"].parsed)

	// Comment should be trimmed
	assert.Equal(t, "Comment with leading spaces", config.comments[0].Text)
}
