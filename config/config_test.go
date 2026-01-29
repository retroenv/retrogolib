package config

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type TestConfig struct {
	CPU     string  `config:"emulation.cpu"`
	Speed   int     `config:"emulation.speed"`
	Debug   bool    `config:"emulation.debug"`
	Timeout float64 `config:"network.timeout"`
	Address int     `config:"nes.code_base_address"`
}

type NestedConfig struct {
	Emulation EmulationConfig `config:"emulation"`
	Network   NetworkConfig   `config:"network"`
}

type EmulationConfig struct {
	CPU   string `config:"cpu"`
	Speed int    `config:"speed"`
	Debug bool   `config:"debug"`
}

type NetworkConfig struct {
	Timeout float64 `config:"timeout"`
	Port    int     `config:"port"`
}

func TestLoad_Success(t *testing.T) {
	data := `[emulation]
cpu = "6502"
speed = 1789773
debug = true

[network]
timeout = 5.5

[nes]
code_base_address = 0x8000`

	var cfg TestConfig
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "6502", cfg.CPU)
	assert.Equal(t, 1789773, cfg.Speed)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 5.5, cfg.Timeout)
	assert.Equal(t, 0x8000, cfg.Address)
}

func TestLoad_NestedStruct(t *testing.T) {
	data := `[emulation]
cpu = "6502"
speed = 1789773
debug = false

[network]
timeout = 2.5
port = 8080`

	var cfg NestedConfig
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)
	assert.Equal(t, "6502", cfg.Emulation.CPU)
	assert.Equal(t, 1789773, cfg.Emulation.Speed)
	assert.False(t, cfg.Emulation.Debug)
	assert.Equal(t, 2.5, cfg.Network.Timeout)
	assert.Equal(t, 8080, cfg.Network.Port)
}

func TestLoadConfig_CommentPreservation(t *testing.T) {
	data := `# RetroGoLib Configuration
# Main emulation settings

[emulation]
# CPU type - currently supports 6502
cpu = "6502"
# Clock speed in Hz
speed = 1789773
debug = false

# Network settings
[network]
timeout = 5.0`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Check comments are preserved (should be 5: global comments + section comments)
	assert.Len(t, config.comments, 5)
	assert.Equal(t, "RetroGoLib Configuration", config.comments[0].Text)
	assert.Equal(t, "Main emulation settings", config.comments[1].Text)
	assert.Equal(t, "CPU type - currently supports 6502", config.comments[2].Text)
	assert.Equal(t, "emulation", config.comments[2].Section)

	// Check structure is preserved
	assert.Greater(t, len(config.structure), 10)
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	original := TestConfig{
		CPU:     "6502",
		Speed:   2000000,
		Debug:   true,
		Timeout: 10.5,
		Address: 0x8000,
	}

	// Create config and marshal
	config, err := LoadConfigBytes([]byte(`[emulation]
cpu = "old"
speed = 1000
debug = false

[network]
timeout = 1.0

[nes]
code_base_address = 0x6000`))
	assert.NoError(t, err)

	err = config.Marshal(&original)
	assert.NoError(t, err)

	// Unmarshal back
	var loaded TestConfig
	err = config.Unmarshal(&loaded)
	assert.NoError(t, err)

	// Verify round-trip preservation
	assert.Equal(t, original, loaded)
}

func TestSaveBytes_CommentPreservation(t *testing.T) {
	data := `# RetroGoLib Configuration
# Main emulation settings

[emulation]
# CPU type
cpu = "6502"
speed = 1789773
debug = false`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Modify a value
	cfg := TestConfig{
		CPU:   "6502",
		Speed: 3579546, // Double the speed
		Debug: true,    // Toggle debug
	}
	err = config.Marshal(&cfg)
	assert.NoError(t, err)

	// Save to bytes
	result, err := config.SaveBytes()
	assert.NoError(t, err)

	resultStr := string(result)

	// Verify comments are preserved
	assert.Contains(t, resultStr, "# RetroGoLib Configuration")
	assert.Contains(t, resultStr, "# Main emulation settings")
	assert.Contains(t, resultStr, "# CPU type")

	// Verify values are updated
	assert.Contains(t, resultStr, "speed = 3579546")
	assert.Contains(t, resultStr, "debug = true")

	// Verify CPU unchanged
	assert.Contains(t, resultStr, `cpu = "6502"`)
}

func TestUnmarshalError_InvalidType(t *testing.T) {
	type InvalidConfig struct {
		Name string `config:"test.name"`
	}

	data := `[test]
name = 42` // Number instead of string

	var cfg InvalidConfig
	err := LoadBytes([]byte(data), &cfg)

	var unmarshalErr *UnmarshalError
	assert.ErrorAs(t, err, &unmarshalErr)
	assert.Equal(t, "Name", unmarshalErr.Field)
	assert.Equal(t, "test", unmarshalErr.Section)
	assert.Equal(t, "name", unmarshalErr.Key)
	assert.ErrorIs(t, unmarshalErr.Err, ErrTypeMismatch)
}

func TestParseError_InvalidFormat(t *testing.T) {
	data := `[emulation]
invalid line without equals`

	_, err := LoadConfigBytes([]byte(data))

	var parseErr *ParseError
	assert.ErrorAs(t, err, &parseErr)
	assert.Equal(t, 2, parseErr.Line)
}

func TestParser_ValueTypes(t *testing.T) {
	data := `[types]
str1 = "quoted string"
str2 = unquoted_string
int_val = 42
hex_val = 0xFF00
bool_true = true
bool_false = false
float_val = 3.14159`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	section := config.sections["types"]
	assert.NotNil(t, section)

	// Test string types
	assert.Equal(t, stringType, section["str1"].vtype)
	assert.Equal(t, "quoted string", section["str1"].parsed)
	assert.Equal(t, stringType, section["str2"].vtype)
	assert.Equal(t, "unquoted_string", section["str2"].parsed)

	// Test numeric types
	assert.Equal(t, intType, section["int_val"].vtype)
	assert.Equal(t, 42, section["int_val"].parsed)
	assert.Equal(t, hexType, section["hex_val"].vtype)
	assert.Equal(t, 0xFF00, section["hex_val"].parsed)

	// Test boolean types
	assert.Equal(t, boolType, section["bool_true"].vtype)
	assert.True(t, section["bool_true"].parsed.(bool))
	assert.Equal(t, boolType, section["bool_false"].vtype)
	assert.False(t, section["bool_false"].parsed.(bool))

	// Test float type
	assert.Equal(t, floatType, section["float_val"].vtype)
	assert.Equal(t, 3.14159, section["float_val"].parsed)
}

func TestNewContent_Addition(t *testing.T) {
	data := `[existing]
old_key = "old_value"`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Add new section and key
	type ExtendedConfig struct {
		OldKey string `config:"existing.old_key"`
		NewKey string `config:"existing.new_key"`
		NewVal int    `config:"new_section.value"`
	}

	extended := ExtendedConfig{
		OldKey: "old_value",
		NewKey: "added_key",
		NewVal: 100,
	}

	err = config.Marshal(&extended)
	assert.NoError(t, err)

	result, err := config.SaveBytes()
	assert.NoError(t, err)

	resultStr := string(result)

	// Verify old content preserved
	assert.Contains(t, resultStr, "old_key = \"old_value\"")

	// Verify new key added to existing section
	assert.Contains(t, resultStr, "new_key = \"added_key\"")

	// Verify new section added
	assert.Contains(t, resultStr, "[new_section]")
	assert.Contains(t, resultStr, "value = 100")
}

func TestDefaultValues(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected DefaultTestConfig
	}{
		{
			name: "all defaults applied",
			data: `# empty config file`,
			expected: DefaultTestConfig{
				StringField:  "default_string",
				IntField:     42,
				BoolField:    true,
				FloatField:   3.14,
				HexField:     255,
				NoDefaultStr: "",
				NoDefaultInt: 0,
			},
		},
		{
			name: "partial config with defaults",
			data: `[test]
string_field = "custom_value"
int_field = 100`,
			expected: DefaultTestConfig{
				StringField:  "custom_value",
				IntField:     100,
				BoolField:    true, // default
				FloatField:   3.14, // default
				HexField:     255,  // default
				NoDefaultStr: "",   // no default
				NoDefaultInt: 0,    // no default
			},
		},
		{
			name: "all fields provided",
			data: `[test]
string_field = "custom_string"
int_field = 200
bool_field = false
float_field = 2.71
hex_field = 0x100
no_default_str = "provided"
no_default_int = 999`,
			expected: DefaultTestConfig{
				StringField:  "custom_string",
				IntField:     200,
				BoolField:    false,
				FloatField:   2.71,
				HexField:     256,
				NoDefaultStr: "provided",
				NoDefaultInt: 999,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg DefaultTestConfig
			err := LoadBytes([]byte(tt.data), &cfg)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestDefaultValueTypes(t *testing.T) {
	tests := []struct {
		name        string
		fieldTag    string
		expectedErr bool
	}{
		{"valid string default", "test.field,default=hello", false},
		{"valid int default", "test.field,default=123", false},
		{"valid bool default true", "test.field,default=true", false},
		{"valid bool default false", "test.field,default=false", false},
		{"valid float default", "test.field,default=1.23", false},
		{"valid hex default", "test.field,default=0xFF", false},
		{"invalid int default", "test.field,default=notanumber", true},
		{"invalid bool default", "test.field,default=maybe", true},
		{"invalid float default", "test.field,default=notafloat", true},
		{"invalid hex default", "test.field,default=0xGG", true},
	}

	config := &Config{sections: make(map[string]section)}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagInfo := config.parseTag(tt.fieldTag, "")

			if !tagInfo.HasDefault {
				t.Fatal("Expected tag to have default value")
			}

			var fieldType reflect.Type
			switch {
			case strings.Contains(tt.name, "string"):
				fieldType = reflect.TypeOf("")
			case strings.Contains(tt.name, "int") || strings.Contains(tt.name, "hex"):
				fieldType = reflect.TypeOf(0)
			case strings.Contains(tt.name, "bool"):
				fieldType = reflect.TypeOf(true)
			case strings.Contains(tt.name, "float"):
				fieldType = reflect.TypeOf(0.0)
			}

			_, err := config.parseDefaultValue(tagInfo.DefaultValue, fieldType)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTagInfoParsing_Basic(t *testing.T) {
	config := &Config{}

	tests := []struct {
		name     string
		tag      string
		parent   string
		expected tagInfo
	}{
		{
			name:   "simple key with default",
			tag:    "key,default=value",
			parent: "section",
			expected: tagInfo{
				Section:      "section",
				Key:          "key",
				DefaultValue: "value",
				HasDefault:   true,
			},
		},
		{
			name:   "section.key with default",
			tag:    "test.field,default=123",
			parent: "",
			expected: tagInfo{
				Section:      "test",
				Key:          "field",
				DefaultValue: "123",
				HasDefault:   true,
			},
		},
		{
			name:   "no default value",
			tag:    "test.field",
			parent: "",
			expected: tagInfo{
				Section:      "test",
				Key:          "field",
				DefaultValue: "",
				HasDefault:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.parseTag(tt.tag, tt.parent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTagInfoParsing_Advanced(t *testing.T) {
	config := &Config{}

	tests := []struct {
		name     string
		tag      string
		parent   string
		expected tagInfo
	}{
		{
			name:   "root level with default",
			tag:    "global_key,default=global_value",
			parent: "",
			expected: tagInfo{
				Section:      "", // Root level uses empty section
				Key:          "global_key",
				DefaultValue: "global_value",
				HasDefault:   true,
			},
		},
		{
			name:   "whitespace handling",
			tag:    " section.key , default=spaced value ",
			parent: "",
			expected: tagInfo{
				Section:      "section",
				Key:          "key",
				DefaultValue: "spaced value ",
				HasDefault:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.parseTag(tt.tag, tt.parent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTagInfoParsing_EdgeCases(t *testing.T) {
	config := &Config{}

	edgeTests := []struct {
		name     string
		tag      string
		parent   string
		expected tagInfo
	}{
		{
			name:   "empty default value",
			tag:    "key,default=",
			parent: "section",
			expected: tagInfo{
				Section:      "section",
				Key:          "key",
				DefaultValue: "",
				HasDefault:   true,
			},
		},
	}

	for _, tt := range edgeTests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.parseTag(tt.tag, tt.parent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

type DefaultTestConfig struct {
	StringField  string  `config:"test.string_field,default=default_string"`
	IntField     int     `config:"test.int_field,default=42"`
	BoolField    bool    `config:"test.bool_field,default=true"`
	FloatField   float64 `config:"test.float_field,default=3.14"`
	HexField     int     `config:"test.hex_field,default=0xFF"`
	NoDefaultStr string  `config:"test.no_default_str"`
	NoDefaultInt int     `config:"test.no_default_int"`
}

func TestRequiredFields_Success(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "all required fields provided",
			data: `[app]
database_url = "postgres://localhost/test"
api_key = "secret123"
port = 8080`,
		},
		{
			name: "missing optional field is ok",
			data: `[app]
database_url = "postgres://localhost/test"
api_key = "secret123"`,
		},
		{
			name: "required with default - missing field uses default",
			data: `[app]
database_url = "postgres://localhost/test"
api_key = "secret123"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg RequiredTestConfig
			err := LoadBytes([]byte(tt.data), &cfg)
			assert.NoError(t, err)
			// Verify required fields are populated
			assert.NotEmpty(t, cfg.DatabaseURL)
			assert.NotEmpty(t, cfg.APIKey)
		})
	}
}

func TestRequiredFields_Errors(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		errorField string
	}{
		{
			name: "missing required database_url",
			data: `[app]
api_key = "secret123"
port = 8080`,
			errorField: "DatabaseURL",
		},
		{
			name: "missing required api_key",
			data: `[app]
database_url = "postgres://localhost/test"
port = 8080`,
			errorField: "APIKey",
		},
		{
			name:       "empty config with required fields",
			data:       `# empty config`,
			errorField: "DatabaseURL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg RequiredTestConfig
			err := LoadBytes([]byte(tt.data), &cfg)
			assert.Error(t, err)
			var unmarshalErr *UnmarshalError
			assert.ErrorAs(t, err, &unmarshalErr)
			assert.Equal(t, tt.errorField, unmarshalErr.Field)
			assert.ErrorIs(t, unmarshalErr.Err, ErrRequiredField)
		})
	}
}

func TestRequiredWithDefault(t *testing.T) {
	// Test that required + default works correctly
	data := `[app]
database_url = "postgres://localhost/test"
api_key = "secret123"`

	var cfg RequiredTestConfig
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	// Required fields should be present
	assert.Equal(t, "postgres://localhost/test", cfg.DatabaseURL)
	assert.Equal(t, "secret123", cfg.APIKey)
	// Required with default should use default
	assert.Equal(t, 8080, cfg.Port)
	// Optional should be zero value
	assert.Equal(t, "", cfg.OptionalSetting)
}

func TestTagParsingRequired(t *testing.T) {
	config := &Config{}

	tests := []struct {
		name     string
		tag      string
		expected tagInfo
	}{
		{
			name: "required flag only",
			tag:  "db.url,required",
			expected: tagInfo{
				Section:      "db",
				Key:          "url",
				DefaultValue: "",
				HasDefault:   false,
				Required:     true,
			},
		},
		{
			name: "required with default",
			tag:  "server.port,required,default=8080",
			expected: tagInfo{
				Section:      "server",
				Key:          "port",
				DefaultValue: "8080",
				HasDefault:   true,
				Required:     true,
			},
		},
		{
			name: "default with required (different order)",
			tag:  "api.key,default=dev-key,required",
			expected: tagInfo{
				Section:      "api",
				Key:          "key",
				DefaultValue: "dev-key",
				HasDefault:   true,
				Required:     true,
			},
		},
		{
			name: "no required flag",
			tag:  "optional.setting,default=value",
			expected: tagInfo{
				Section:      "optional",
				Key:          "setting",
				DefaultValue: "value",
				HasDefault:   true,
				Required:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.parseTag(tt.tag, "")
			assert.Equal(t, tt.expected, result)
		})
	}
}

type RequiredTestConfig struct {
	DatabaseURL     string `config:"app.database_url,required"`
	APIKey          string `config:"app.api_key,required"`
	Port            int    `config:"app.port,required,default=8080"`
	OptionalSetting string `config:"app.optional_setting"`
}

func TestDuplicateSections_ErrorsBehavior(t *testing.T) {
	data := `[database]
host = "localhost"
port = 5432

[server]
port = 8080
debug = true

[database]
user = "admin"
password = "secret"
port = 3306`

	_, err := LoadConfigBytes([]byte(data))
	assert.ErrorIs(t, err, ErrDuplicateSection)
	assert.ErrorContains(t, err, "section 'database' first defined at line 1")
}

// TestDuplicateSections_ErrorOnDuplicate tests that duplicate sections cause errors.
func TestDuplicateSections_ErrorOnDuplicate(t *testing.T) {
	data := `[config]
key1 = "value1"

[config]
key2 = "value2"`

	_, err := LoadConfigBytes([]byte(data))
	assert.ErrorIs(t, err, ErrDuplicateSection)
	assert.ErrorContains(t, err, "section 'config' first defined at line 1")
}

// TestDuplicateKeys_ErrorOnDuplicate tests that duplicate keys within a section cause errors.
func TestDuplicateKeys_ErrorOnDuplicate(t *testing.T) {
	data := `[settings]
timeout = 30
retries = 3
timeout = 60`

	_, err := LoadConfigBytes([]byte(data))
	assert.ErrorIs(t, err, ErrDuplicateKey)
	assert.ErrorContains(t, err, "key 'timeout' in section 'settings' first defined at line 2")
}

// TestValidConfig_NoDuplicates tests that valid configs without duplicates still work.
func TestValidConfig_NoDuplicates(t *testing.T) {
	data := `[app]
name = "myapp"
version = "1.0"

[database]
host = "localhost"
port = 5432`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Verify app section
	appSection := config.sections["app"]
	assert.Len(t, appSection, 2)
	assert.Equal(t, "myapp", appSection["name"].parsed)
	assert.Equal(t, "1.0", appSection["version"].parsed)

	// Verify database section
	dbSection := config.sections["database"]
	assert.Len(t, dbSection, 2)
	assert.Equal(t, "localhost", dbSection["host"].parsed)
	assert.Equal(t, 5432, dbSection["port"].parsed)

	// Test save/reload round-trip
	savedData, err := config.SaveBytes()
	assert.NoError(t, err)

	reloadedConfig, err := LoadConfigBytes(savedData)
	assert.NoError(t, err)
	assert.Equal(t, len(config.sections), len(reloadedConfig.sections))
}

// TestDuplicateSection_Variations tests different patterns of duplicate sections.
func TestDuplicateSection_Variations(t *testing.T) {
	tests := []struct {
		name          string
		data          string
		expectedLine  int
		sectionName   string
		duplicateLine int
	}{
		{
			name: "duplicate_immediately_after",
			data: `[section1]
key1 = "value1"
[section1]
key2 = "value2"`,
			expectedLine:  1,
			sectionName:   "section1",
			duplicateLine: 3,
		},
		{
			name: "duplicate_with_gap",
			data: `[config]
setting = "value"

# Comment between
[other]
key = "value"

[config]
duplicate = "error"`,
			expectedLine:  1,
			sectionName:   "config",
			duplicateLine: 8,
		},
		{
			name: "multiple_sections_same_duplicate",
			data: `[app]
name = "test"
[database]
host = "localhost"
[app]
version = "1.0"`,
			expectedLine:  1,
			sectionName:   "app",
			duplicateLine: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadConfigBytes([]byte(tt.data))
			assert.ErrorIs(t, err, ErrDuplicateSection)
			assert.ErrorContains(t, err, tt.sectionName)
			assert.ErrorContains(t, err, "first defined at line")
		})
	}
}

// TestDuplicateKey_Variations tests different patterns of duplicate keys.
func TestDuplicateKey_Variations(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		expectedKey string
		section     string
	}{
		{
			name: "duplicate_key_immediate",
			data: `[section]
key = "value1"
key = "value2"`,
			expectedKey: "key",
			section:     "section",
		},
		{
			name: "duplicate_key_with_gap",
			data: `[config]
setting1 = "value1"
setting2 = "value2"
# Comment
setting1 = "duplicate"`,
			expectedKey: "setting1",
			section:     "config",
		},
		{
			name: "exact_duplicate_key",
			data: `[test]
key = "value1"
key = "value2"`,
			expectedKey: "key",
			section:     "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadConfigBytes([]byte(tt.data))
			assert.ErrorIs(t, err, ErrDuplicateKey)
			assert.ErrorContains(t, err, tt.expectedKey)
			assert.ErrorContains(t, err, tt.section)
			assert.ErrorContains(t, err, "first defined at line")
		})
	}
}

// TestCaseInsensitiveKeys ensures keys with different cases are treated as duplicates.
func TestCaseInsensitiveKeys(t *testing.T) {
	data := `[test]
Key = "value1"
KEY = "value2"`

	_, err := LoadConfigBytes([]byte(data))
	assert.ErrorIs(t, err, ErrDuplicateKey)
	assert.ErrorContains(t, err, "key 'KEY' in section 'test' first defined at line")
}

// TestCaseInsensitiveAccess ensures keys can be accessed case-insensitively.
func TestCaseInsensitiveAccess(t *testing.T) {
	data := `[Database]
Host = "localhost"
PORT = 5432

[Logging]
level = "info"`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Sections are normalized to lowercase
	databaseSection := config.sections["database"]
	assert.NotNil(t, databaseSection)
	assert.Equal(t, "localhost", databaseSection["host"].parsed)
	assert.Equal(t, 5432, databaseSection["port"].parsed)

	loggingSection := config.sections["logging"]
	assert.NotNil(t, loggingSection)
	assert.Equal(t, "info", loggingSection["level"].parsed)
}

// TestNoDuplicates_DifferentSections ensures keys can be duplicated across different sections.
func TestNoDuplicates_DifferentSections(t *testing.T) {
	data := `[database]
host = "localhost"
port = 5432

[redis]
host = "redis-server"
port = 6379

[app]
port = 8080`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Verify same key names in different sections work fine
	assert.Equal(t, "localhost", config.sections["database"]["host"].parsed)
	assert.Equal(t, "redis-server", config.sections["redis"]["host"].parsed)

	assert.Equal(t, 5432, config.sections["database"]["port"].parsed)
	assert.Equal(t, 6379, config.sections["redis"]["port"].parsed)
	assert.Equal(t, 8080, config.sections["app"]["port"].parsed)
}

// TestComplexValidConfig ensures complex but valid configs still work.
func TestComplexValidConfig(t *testing.T) {
	data := `# Global config
[database]
host = "localhost"
port = 5432
user = "admin"
password = "secret"

# Server configuration
[server]
host = "0.0.0.0"
port = 8080
debug = true
max_connections = 100

# Logging configuration  
[logging]
level = "info"
file = "/var/log/app.log"
rotate = true

# Cache configuration
[cache]
enabled = true
ttl = 3600
max_size = 1000`

	config, err := LoadConfigBytes([]byte(data))
	assert.NoError(t, err)

	// Verify all sections loaded correctly
	assert.Len(t, config.sections, 4)

	// Spot check some values
	assert.Equal(t, "localhost", config.sections["database"]["host"].parsed)
	assert.Equal(t, 8080, config.sections["server"]["port"].parsed)
	assert.Equal(t, "info", config.sections["logging"]["level"].parsed)
	assert.True(t, config.sections["cache"]["enabled"].parsed.(bool))
}

// TestEdgeCases_EmptyButValid tests edge cases that should be valid.
func TestEdgeCases_EmptyButValid(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "single_section_single_key",
			data: `[test]
key = "value"`,
		},
		{
			name: "empty_sections_different_names",
			data: `[section1]

[section2]

[section3]`,
		},
		{
			name: "sections_with_only_comments",
			data: `[app]
# This section has only comments

[database]  
# Another comment-only section`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadConfigBytes([]byte(tt.data))
			assert.NoError(t, err)
		})
	}
}

// TestParseError_Integration tests that duplicate errors are properly wrapped in ParseError.
func TestParseError_DuplicateIntegration(t *testing.T) {
	data := `[valid]
key = "value"

[valid]
duplicate = "error"`

	_, err := LoadConfigBytes([]byte(data))
	assert.Error(t, err)

	var parseErr *ParseError
	assert.ErrorAs(t, err, &parseErr)
	assert.Equal(t, 4, parseErr.Line) // Line where duplicate section occurs
	assert.ErrorIs(t, parseErr.Err, ErrDuplicateSection)
}

// TestSecurityValidation tests security-related input validation using constants.
func TestSecurityValidation(t *testing.T) {
	t.Run("maxConfigSize", func(t *testing.T) {
		// Create data larger than maxConfigSize
		largeData := make([]byte, maxConfigSize+1)
		for i := range largeData {
			largeData[i] = 'a'
		}

		_, err := LoadConfigBytes(largeData)
		assert.ErrorIs(t, err, ErrConfigTooLarge)
	})

	t.Run("maxLines", func(t *testing.T) {
		// Create config with too many lines
		var buf strings.Builder
		buf.WriteString("[section]\n")
		for i := range maxLines + 100 {
			buf.WriteString(fmt.Sprintf("key%d = value%d\n", i, i))
		}

		_, err := LoadConfigBytes([]byte(buf.String()))
		assert.ErrorIs(t, err, ErrTooManyLines)
	})

	t.Run("LongSectionName", func(t *testing.T) {
		longName := strings.Repeat("a", maxNameLength+1)
		data := fmt.Sprintf("[%s]\nkey = value", longName)

		_, err := LoadConfigBytes([]byte(data))
		assert.ErrorIs(t, err, ErrSectionNameTooLong)
	})

	t.Run("LongKeyName", func(t *testing.T) {
		longKey := strings.Repeat("a", maxNameLength+1)
		data := fmt.Sprintf("[section]\n%s = value", longKey)

		_, err := LoadConfigBytes([]byte(data))
		assert.ErrorIs(t, err, ErrKeyNameTooLong)
	})
}

// TestConstants verifies the values of configuration constants.
func TestConstants(t *testing.T) {
	assert.Equal(t, 10*1024*1024, maxConfigSize)
	assert.Equal(t, 100000, maxLines)
	assert.Equal(t, 256, maxNameLength)
	assert.Equal(t, 40, avgElementSize)
	assert.Equal(t, 0644, configFilePermissions)
}

// TestAutomaticFieldMapping tests that fields without config tags are automatically mapped
func TestAutomaticFieldMapping_SimpleFields(t *testing.T) {
	data := `name = "test"
port = 8080
enabled = true
timeout = 30

[database]
host = "localhost"`

	type Config struct {
		Name    string // Should map to root-level name
		Port    int    // Should map to root-level port
		Enabled bool   // Should map to root-level enabled
		Timeout int    // Should map to root-level timeout
		Host    string `config:"database.host"` // Explicit tag for comparison
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "test", cfg.Name)
	assert.Equal(t, 8080, cfg.Port)
	assert.True(t, cfg.Enabled)
	assert.Equal(t, 30, cfg.Timeout)
	assert.Equal(t, "localhost", cfg.Host)
}

// TestAutomaticFieldMapping_NestedStructs tests automatic mapping with nested structs
func TestAutomaticFieldMapping_NestedStructs(t *testing.T) {
	data := `[database]
host = "localhost"
port = 5432
timeout = 30

[cache]
host = "redis-server"
port = 6379
ttl = 3600`

	type DatabaseConfig struct {
		Host    string // Should map to database.host
		Port    int    // Should map to database.port
		Timeout int    // Should map to database.timeout
	}

	type CacheConfig struct {
		Host string // Should map to cache.host
		Port int    // Should map to cache.port
		TTL  int    // Should map to cache.ttl (lowercase)
	}

	type Config struct {
		Database DatabaseConfig // Should use "database" as section
		Cache    CacheConfig    // Should use "cache" as section
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, 30, cfg.Database.Timeout)

	assert.Equal(t, "redis-server", cfg.Cache.Host)
	assert.Equal(t, 6379, cfg.Cache.Port)
	assert.Equal(t, 3600, cfg.Cache.TTL)
}

// TestAutomaticFieldMapping_MixedTagsAndAuto tests mixing explicit tags with automatic mapping
func TestAutomaticFieldMapping_MixedTagsAndAuto(t *testing.T) {
	data := `[server]
name = "web-server"
port = 8080
debug = true

[database]
connection_string = "postgres://localhost/db"`

	type Config struct {
		ServerName string `config:"server.name"` // Explicit tag
		Port       int    // Should map to global.port (not found, so zero value)
		Debug      bool   `config:"server.debug"` // Explicit tag
		ConnStr    string // Should map to global.connstr (not found, zero value)
		ServerPort int    `config:"server.port"`                // Explicit tag
		DBConn     string `config:"database.connection_string"` // Explicit tag
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "web-server", cfg.ServerName)
	assert.Equal(t, 0, cfg.Port) // Not found in global section
	assert.True(t, cfg.Debug)
	assert.Equal(t, "", cfg.ConnStr) // Not found in global section
	assert.Equal(t, 8080, cfg.ServerPort)
	assert.Equal(t, "postgres://localhost/db", cfg.DBConn)
}

// TestAutomaticFieldMapping_IgnoreTag tests that fields with "-" tag are ignored
func TestAutomaticFieldMapping_IgnoreTag(t *testing.T) {
	data := `name = "test"
secret = "should-not-load"`

	type Config struct {
		Name   string // Should map to root-level name
		Secret string `config:"-"` // Should be ignored
		Port   int    // Should map to root-level port (not found, zero value)
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "test", cfg.Name)
	assert.Equal(t, "", cfg.Secret) // Should remain zero value
	assert.Equal(t, 0, cfg.Port)    // Not found, zero value
}

// TestAutomaticFieldMapping_CaseInsensitive tests that field names are converted to lowercase
func TestAutomaticFieldMapping_CaseInsensitive(t *testing.T) {
	data := `myfield = "test"
anotherfield = 42

[section]
mixedcase = true`

	type Config struct {
		MyField      string // Should map to root-level myfield
		AnotherField int    // Should map to root-level anotherfield
		MixedCase    bool   `config:"section.mixedcase"` // Explicit for comparison
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "test", cfg.MyField)
	assert.Equal(t, 42, cfg.AnotherField)
	assert.True(t, cfg.MixedCase)
}

// TestAutomaticFieldMapping_Marshal tests that automatic mapping works for marshalling too
func TestAutomaticFieldMapping_Marshal(t *testing.T) {
	type TestConfig struct {
		Name    string
		Port    int
		Enabled bool
	}

	cfg := TestConfig{
		Name:    "test-server",
		Port:    9090,
		Enabled: true,
	}

	config := &Config{sections: make(map[string]section)}
	err := config.Marshal(&cfg)
	assert.NoError(t, err)

	// Check that values were marshalled to the correct sections/keys (root level = empty section)
	assert.NotNil(t, config.sections[""])
	assert.Equal(t, "test-server", config.sections[""]["name"].parsed)
	assert.Equal(t, 9090, config.sections[""]["port"].parsed)
	assert.True(t, config.sections[""]["enabled"].parsed.(bool))
}

// TestGenerateFieldTag tests the generateFieldTag helper function
func TestGenerateFieldTag(t *testing.T) {
	config := &Config{}

	// Test simple field with no parent section (root level)
	tag := config.generateFieldTag("Name", "", false)
	assert.Equal(t, "name", tag)

	// Test simple field with parent section
	tag = config.generateFieldTag("Port", "server", false)
	assert.Equal(t, "server.port", tag)

	// Test struct field (should use field name as section)
	tag = config.generateFieldTag("Database", "", true)
	assert.Equal(t, "database", tag)

	// Test struct field with parent (parent should be ignored for structs)
	tag = config.generateFieldTag("Cache", "app", true)
	assert.Equal(t, "cache", tag)

	// Test field name case conversion
	tag = config.generateFieldTag("MyLongFieldName", "section", false)
	assert.Equal(t, "section.mylongfieldname", tag)
}

// ============================================================================
// Deep Nesting Tests (3+ Levels Deep)
// ============================================================================

// TestDeepNesting_ThreeLevels tests 3-level nested struct configuration
func TestDeepNesting_ThreeLevels(t *testing.T) {
	data := `[system]
cpu_type = "6502"
memory_size = 64

[system.cpu]
frequency = 1789773
debug_mode = true

[system.cpu.cache]
enabled = true
size = 8192
policy = "write-through"`

	// Define 3-level nested structs
	type CacheConfig struct {
		Enabled bool   `config:"enabled"`
		Size    int    `config:"size"`
		Policy  string `config:"policy"`
	}

	type CPUConfig struct {
		Frequency int         `config:"frequency"`
		DebugMode bool        `config:"debug_mode"`
		Cache     CacheConfig `config:"cache"`
	}

	type SystemConfig struct {
		CPUType    string    `config:"cpu_type"`
		MemorySize int       `config:"memory_size"`
		CPU        CPUConfig `config:"cpu"`
	}

	type Config struct {
		System SystemConfig `config:"system"`
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	// Verify 3-level deep access
	assert.Equal(t, "6502", cfg.System.CPUType)
	assert.Equal(t, 64, cfg.System.MemorySize)
	assert.Equal(t, 1789773, cfg.System.CPU.Frequency)
	assert.True(t, cfg.System.CPU.DebugMode)
	assert.True(t, cfg.System.CPU.Cache.Enabled)
	assert.Equal(t, 8192, cfg.System.CPU.Cache.Size)
	assert.Equal(t, "write-through", cfg.System.CPU.Cache.Policy)
}

// TestDeepNesting_FourLevels tests 4-level nested struct configuration
func TestDeepNesting_FourLevels(t *testing.T) {
	data := getFourLevelTestData()
	cfg := getFourLevelConfig()

	err := LoadBytes([]byte(data), cfg)
	assert.NoError(t, err)

	verifyFourLevelResults(t, cfg)
}

func getFourLevelTestData() string {
	return `[emulator]
name = "RetroGo"
version = "1.0"

[emulator.console]
type = "nes"
region = "ntsc"

[emulator.console.cartridge]
mapper = 0
prg_banks = 2
chr_banks = 1

[emulator.console.cartridge.header]
valid = true
trainer = false
battery = true
fourscreen = false`
}

func getFourLevelConfig() *fourLevelConfig {
	return &fourLevelConfig{}
}

func verifyFourLevelResults(t *testing.T, cfg *fourLevelConfig) {
	t.Helper()
	// Verify 4-level deep access
	assert.Equal(t, "RetroGo", cfg.Emulator.Name)
	assert.Equal(t, "1.0", cfg.Emulator.Version)
	assert.Equal(t, "nes", cfg.Emulator.Console.Type)
	assert.Equal(t, "ntsc", cfg.Emulator.Console.Region)
	assert.Equal(t, 0, cfg.Emulator.Console.Cartridge.Mapper)
	assert.Equal(t, 2, cfg.Emulator.Console.Cartridge.PRGBanks)
	assert.Equal(t, 1, cfg.Emulator.Console.Cartridge.CHRBanks)
	assert.True(t, cfg.Emulator.Console.Cartridge.Header.Valid)
	assert.False(t, cfg.Emulator.Console.Cartridge.Header.Trainer)
	assert.True(t, cfg.Emulator.Console.Cartridge.Header.Battery)
	assert.False(t, cfg.Emulator.Console.Cartridge.Header.Fourscreen)
}

// Define 4-level nested structs
type fourLevelHeaderConfig struct {
	Valid      bool `config:"valid"`
	Trainer    bool `config:"trainer"`
	Battery    bool `config:"battery"`
	Fourscreen bool `config:"fourscreen"`
}

type fourLevelCartridgeConfig struct {
	Mapper   int                   `config:"mapper"`
	PRGBanks int                   `config:"prg_banks"`
	CHRBanks int                   `config:"chr_banks"`
	Header   fourLevelHeaderConfig `config:"header"`
}

type fourLevelConsoleConfig struct {
	Type      string                   `config:"type"`
	Region    string                   `config:"region"`
	Cartridge fourLevelCartridgeConfig `config:"cartridge"`
}

type fourLevelEmulatorConfig struct {
	Name    string                 `config:"name"`
	Version string                 `config:"version"`
	Console fourLevelConsoleConfig `config:"console"`
}

type fourLevelConfig struct {
	Emulator fourLevelEmulatorConfig `config:"emulator"`
}

// TestDeepNesting_AutomaticMapping tests automatic field mapping in deep nesting
func TestDeepNesting_AutomaticMapping(t *testing.T) {
	data := `[system]
name = "retro-system"
enabled = true

[cpu]
type = "6502"
speed = 2000000

[cpu.registers]
accumulator = 0
xindex = 0
yindex = 0

[cpu.registers.flags]
carry = false
zero = true
interrupt = false
decimal = false`

	// Define structs with automatic field mapping (no explicit config tags)
	type FlagsConfig struct {
		Carry     bool // Automatically maps to cpu.registers.flags.carry
		Zero      bool // Automatically maps to cpu.registers.flags.zero
		Interrupt bool // Automatically maps to cpu.registers.flags.interrupt
		Decimal   bool // Automatically maps to cpu.registers.flags.decimal
	}

	type RegistersConfig struct {
		Accumulator int         // Automatically maps to cpu.registers.accumulator
		XIndex      int         // Automatically maps to cpu.registers.x_index
		YIndex      int         // Automatically maps to cpu.registers.y_index
		Flags       FlagsConfig // Uses "flags" as subsection
	}

	type CPUConfig struct {
		Type      string          // Automatically maps to cpu.type
		Speed     int             // Automatically maps to cpu.speed
		Registers RegistersConfig // Uses "registers" as subsection
	}

	type SystemConfig struct {
		Name    string // Automatically maps to system.name
		Enabled bool   // Automatically maps to system.enabled
	}

	type Config struct {
		System SystemConfig // Uses "system" as section
		CPU    CPUConfig    // Uses "cpu" as section
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	// Verify automatic mapping works at all levels
	assert.Equal(t, "retro-system", cfg.System.Name)
	assert.True(t, cfg.System.Enabled)
	assert.Equal(t, "6502", cfg.CPU.Type)
	assert.Equal(t, 2000000, cfg.CPU.Speed)
	assert.Equal(t, 0, cfg.CPU.Registers.Accumulator)
	assert.Equal(t, 0, cfg.CPU.Registers.XIndex)
	assert.Equal(t, 0, cfg.CPU.Registers.YIndex)
	assert.False(t, cfg.CPU.Registers.Flags.Carry)
	assert.True(t, cfg.CPU.Registers.Flags.Zero)
	assert.False(t, cfg.CPU.Registers.Flags.Interrupt)
	assert.False(t, cfg.CPU.Registers.Flags.Decimal)
}

// TestDeepNesting_MixedTagsAndAuto tests mixing explicit tags with automatic mapping in deep structures
func TestDeepNesting_MixedTagsAndAuto(t *testing.T) {
	data := getMixedTagsTestData()
	cfg := getMixedTagsConfig()

	err := LoadBytes([]byte(data), cfg)
	assert.NoError(t, err)

	verifyMixedTagsResults(t, cfg)
}

// OpenGLConfig represents OpenGL configuration with mixed tag types
type OpenGLConfig struct {
	Version     string `config:"version"`
	CoreProfile bool   // Automatic: graphics.opengl.coreprofile
}

// GraphicsConfig represents graphics configuration with nested OpenGL
type GraphicsConfig struct {
	Width  int          `config:"width"`
	Height int          `config:"height"`
	VSync  bool         // Automatic: graphics.vsync
	OpenGL OpenGLConfig `config:"opengl"`
}

// ChannelsConfig represents audio channels configuration
type ChannelsConfig struct {
	MasterVolume float64 `config:"mastervolume"`
	SFXVolume    float64 // Automatic: audio.channels.sfxvolume
}

// AudioConfig represents audio configuration with nested channels
type AudioConfig struct {
	SampleRate int            `config:"sample_rate"`
	BufferSize int            // Automatic: audio.buffersize
	Channels   ChannelsConfig // Automatic: uses "channels" subsection
}

// AppConfig represents application configuration
type AppConfig struct {
	Name  string // Automatic: app.name
	Debug bool   `config:"app.debug"` // Explicit tag
}

// mixedTagsConfig represents configuration with mixed tagging approaches
type mixedTagsConfig struct {
	App      AppConfig      // Automatic: uses "app" section
	Graphics GraphicsConfig `config:"graphics"`
	Audio    AudioConfig    // Automatic: uses "audio" section
}

func getMixedTagsTestData() string {
	return `[app]
name = "RetroGoLib"
debug = true

[graphics]
width = 256
height = 240
vsync = true

[graphics.opengl]
version = "3.3"
coreprofile = true

[audio]
sample_rate = 44100
buffersize = 1024

[audio.channels]
mastervolume = 0.8
sfxvolume = 0.6`
}

func getMixedTagsConfig() *mixedTagsConfig {
	return &mixedTagsConfig{}
}

func verifyMixedTagsResults(t *testing.T, cfg *mixedTagsConfig) {
	t.Helper()
	// Verify mixed mapping works correctly
	assert.Equal(t, "RetroGoLib", cfg.App.Name)
	assert.True(t, cfg.App.Debug)
	assert.Equal(t, 256, cfg.Graphics.Width)
	assert.Equal(t, 240, cfg.Graphics.Height)
	assert.True(t, cfg.Graphics.VSync)
	assert.Equal(t, "3.3", cfg.Graphics.OpenGL.Version)
	assert.True(t, cfg.Graphics.OpenGL.CoreProfile)
	assert.Equal(t, 44100, cfg.Audio.SampleRate)
	assert.Equal(t, 1024, cfg.Audio.BufferSize)
	assert.Equal(t, 0.8, cfg.Audio.Channels.MasterVolume)
	assert.Equal(t, 0.6, cfg.Audio.Channels.SFXVolume)
}

// TestDeepNesting_WithDefaults tests deep nesting with default values
func TestDeepNesting_WithDefaults(t *testing.T) {
	// Minimal config - most values will use defaults
	data := `[compiler]
target = "6502"

[compiler.optimization]
level = 2`

	type DebugConfig struct {
		Enabled    bool   `config:"enabled,default=false"`
		OutputFile string `config:"output_file,default=debug.log"`
		Verbose    bool   `config:"verbose,default=true"`
	}

	type OptimizationConfig struct {
		Level       int         `config:"level,default=1"`
		Inline      bool        `config:"inline,default=true"`
		DeadCode    bool        `config:"dead_code,default=true"`
		DebugConfig DebugConfig `config:"debug"`
	}

	type CompilerConfig struct {
		Target       string             `config:"target,default=generic"`
		OutputDir    string             `config:"output_dir,default=./build"`
		Optimization OptimizationConfig `config:"optimization"`
	}

	type Config struct {
		Compiler CompilerConfig `config:"compiler"`
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	// Values from config
	assert.Equal(t, "6502", cfg.Compiler.Target)
	assert.Equal(t, 2, cfg.Compiler.Optimization.Level)

	// Default values applied
	assert.Equal(t, "./build", cfg.Compiler.OutputDir)
	assert.True(t, cfg.Compiler.Optimization.Inline)
	assert.True(t, cfg.Compiler.Optimization.DeadCode)
	assert.False(t, cfg.Compiler.Optimization.DebugConfig.Enabled)
	assert.Equal(t, "debug.log", cfg.Compiler.Optimization.DebugConfig.OutputFile)
	assert.True(t, cfg.Compiler.Optimization.DebugConfig.Verbose)
}

// TestDeepNesting_EmptyStructs tests deep nesting with empty intermediate structs
func TestDeepNesting_EmptyStructs(t *testing.T) {
	data := `[platform]
name = "NES"

[platform.memory.regions.prg]
start = 0x8000
size = 0x8000

[platform.memory.regions.chr]
start = 0x0000
size = 0x2000`

	type PRGConfig struct {
		Start int `config:"start"`
		Size  int `config:"size"`
	}

	type CHRConfig struct {
		Start int `config:"start"`
		Size  int `config:"size"`
	}

	type RegionsConfig struct {
		PRG PRGConfig `config:"prg"`
		CHR CHRConfig `config:"chr"`
	}

	type MemoryConfig struct {
		Regions RegionsConfig `config:"regions"`
	}

	type PlatformConfig struct {
		Name   string       `config:"name"`
		Memory MemoryConfig `config:"memory"`
	}

	type Config struct {
		Platform PlatformConfig `config:"platform"`
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "NES", cfg.Platform.Name)
	assert.Equal(t, 0x8000, cfg.Platform.Memory.Regions.PRG.Start)
	assert.Equal(t, 0x8000, cfg.Platform.Memory.Regions.PRG.Size)
	assert.Equal(t, 0x0000, cfg.Platform.Memory.Regions.CHR.Start)
	assert.Equal(t, 0x2000, cfg.Platform.Memory.Regions.CHR.Size)
}

// TestDeepNesting_Marshal tests marshaling deep nested structs back to config
func TestDeepNesting_Marshal(t *testing.T) {
	type CacheConfig struct {
		Enabled bool   `config:"enabled"`
		Size    int    `config:"size"`
		Policy  string `config:"policy"`
	}

	type CPUConfig struct {
		Frequency int         `config:"frequency"`
		Cache     CacheConfig `config:"cache"`
	}

	type SystemConfig struct {
		Name string    `config:"name"`
		CPU  CPUConfig `config:"cpu"`
	}

	type Config struct {
		System SystemConfig `config:"system"`
	}

	// Create nested config
	cfg := Config{
		System: SystemConfig{
			Name: "RetroSystem",
			CPU: CPUConfig{
				Frequency: 1789773,
				Cache: CacheConfig{
					Enabled: true,
					Size:    8192,
					Policy:  "write-back",
				},
			},
		},
	}

	// Create empty config and marshal struct into it
	configObj, err := LoadConfigBytes([]byte("# empty config"))
	assert.NoError(t, err)

	err = configObj.Marshal(&cfg)
	assert.NoError(t, err)

	// Test round-trip: unmarshal back and verify
	var reloaded Config
	err = configObj.Unmarshal(&reloaded)
	assert.NoError(t, err)
	assert.Equal(t, cfg, reloaded)

	// Save and verify config can be written back
	result, err := configObj.SaveBytes()
	assert.NoError(t, err)

	// Verify the saved config contains the expected sections and values
	configStr := string(result)
	assert.Contains(t, configStr, "[system]")
	assert.Contains(t, configStr, "[system.cpu]")
	assert.Contains(t, configStr, "[system.cpu.cache]")
	assert.Contains(t, configStr, `name = "RetroSystem"`)
	assert.Contains(t, configStr, "frequency = 1789773")
	assert.Contains(t, configStr, "enabled = true")
	assert.Contains(t, configStr, "size = 8192")
	assert.Contains(t, configStr, `policy = "write-back"`)
}

// TestRootLevelKeys tests support for root-level keys without requiring a global section
func TestRootLevelKeys(t *testing.T) {
	data := `# Root level configuration
name = "RetroGoLib"
version = "1.0.0"
debug = true
port = 8080
timeout = 30.5

[database]
host = "localhost"
port = 5432

[logging]
level = "info"`

	type Config struct {
		Name    string  `config:"name"`
		Version string  `config:"version"`
		Debug   bool    `config:"debug"`
		Port    int     `config:"port"`
		Timeout float64 `config:"timeout"`

		Database struct {
			Host string `config:"host"`
			Port int    `config:"port"`
		} `config:"database"`

		Logging struct {
			Level string `config:"level"`
		} `config:"logging"`
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	// Verify root-level keys
	assert.Equal(t, "RetroGoLib", cfg.Name)
	assert.Equal(t, "1.0.0", cfg.Version)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, 30.5, cfg.Timeout)

	// Verify section-based keys still work
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "info", cfg.Logging.Level)
}

// TestRootLevelKeysWithAutoMapping tests root-level keys with automatic field mapping
func TestRootLevelKeysWithAutoMapping(t *testing.T) {
	data := `title = "RetroGo Emulator"
enabled = true
maxfps = 60

[graphics]
width = 256
height = 240`

	type Config struct {
		Title   string // Automatic mapping to root-level "title"
		Enabled bool   // Automatic mapping to root-level "enabled"
		MaxFPS  int    // Automatic mapping to root-level "maxfps" (lowercase conversion)

		Graphics struct {
			Width  int // Automatic mapping to "graphics.width"
			Height int // Automatic mapping to "graphics.height"
		} // Automatic mapping to "graphics" section
	}

	var cfg Config
	err := LoadBytes([]byte(data), &cfg)
	assert.NoError(t, err)

	// Verify automatic mapping for root-level keys
	assert.Equal(t, "RetroGo Emulator", cfg.Title)
	assert.True(t, cfg.Enabled)
	assert.Equal(t, 60, cfg.MaxFPS)

	// Verify section-based automatic mapping still works
	assert.Equal(t, 256, cfg.Graphics.Width)
	assert.Equal(t, 240, cfg.Graphics.Height)
}

// TestRootLevelKeysMarshal tests marshaling structs with root-level keys
func TestRootLevelKeysMarshal(t *testing.T) {
	type Config struct {
		AppName string `config:"app_name"`
		Version string `config:"version"`
		Debug   bool   `config:"debug"`

		Server struct {
			Host string `config:"host"`
			Port int    `config:"port"`
		} `config:"server"`
	}

	cfg := Config{
		AppName: "TestApp",
		Version: "2.0.0",
		Debug:   false,
		Server: struct {
			Host string `config:"host"`
			Port int    `config:"port"`
		}{
			Host: "0.0.0.0",
			Port: 9000,
		},
	}

	// Create empty config and marshal struct into it
	configObj, err := LoadConfigBytes([]byte("# empty config"))
	assert.NoError(t, err)

	err = configObj.Marshal(&cfg)
	assert.NoError(t, err)

	// Verify root-level keys are written correctly
	result, err := configObj.SaveBytes()
	assert.NoError(t, err)

	configStr := string(result)
	assert.Contains(t, configStr, `app_name = "TestApp"`)
	assert.Contains(t, configStr, `version = "2.0.0"`)
	assert.Contains(t, configStr, "debug = false")
	assert.Contains(t, configStr, "[server]")
	assert.Contains(t, configStr, `host = "0.0.0.0"`)
	assert.Contains(t, configStr, "port = 9000")

	// Test round-trip
	var reloaded Config
	err = configObj.Unmarshal(&reloaded)
	assert.NoError(t, err)
	assert.Equal(t, cfg, reloaded)
}

// TestRootLevelKeysIntegration tests comprehensive root-level key functionality with edge cases
func TestRootLevelKeysIntegration(t *testing.T) {
	data := getIntegrationTestData()
	cfg := getIntegrationConfig()

	err := LoadBytes([]byte(data), cfg)
	assert.NoError(t, err)

	verifyIntegrationResults(t, cfg)
	testIntegrationMarshal(t, cfg)
}

func getIntegrationTestData() string {
	return `# Configuration with root-level keys and sections
app_name = "RetroGoLib"
version = "1.0.0"
debug = true
max_connections = 100
timeout = 30.5

[database]
host = "localhost"
port = 5432
ssl = true

[cache]
enabled = true
ttl = 3600`
}

type integrationDatabaseConfig struct {
	Host string `config:"host"`
	Port int    `config:"port"`
	SSL  bool   `config:"ssl"`
}

type integrationCacheConfig struct {
	Enabled bool `config:"enabled"`
	TTL     int  `config:"ttl"`
}

type integrationConfig struct {
	AppName        string                    `config:"app_name"`
	Version        string                    `config:"version"`
	Debug          bool                      `config:"debug"`
	MaxConnections int                       `config:"max_connections"`
	Timeout        float64                   `config:"timeout"`
	Database       integrationDatabaseConfig `config:"database"`
	Cache          integrationCacheConfig    `config:"cache"`
}

func getIntegrationConfig() *integrationConfig {
	return &integrationConfig{}
}

func verifyIntegrationResults(t *testing.T, cfg *integrationConfig) {
	t.Helper()
	// Verify root-level keys work correctly
	assert.Equal(t, "RetroGoLib", cfg.AppName)
	assert.Equal(t, "1.0.0", cfg.Version)
	assert.True(t, cfg.Debug)
	assert.Equal(t, 100, cfg.MaxConnections)
	assert.Equal(t, 30.5, cfg.Timeout)

	// Verify section keys still work
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.True(t, cfg.Database.SSL)
	assert.True(t, cfg.Cache.Enabled)
	assert.Equal(t, 3600, cfg.Cache.TTL)
}

func testIntegrationMarshal(t *testing.T, cfg *integrationConfig) {
	t.Helper()
	// Test marshaling back
	configObj, err := LoadConfigBytes([]byte("# empty"))
	assert.NoError(t, err)

	err = configObj.Marshal(cfg)
	assert.NoError(t, err)

	result, err := configObj.SaveBytes()
	assert.NoError(t, err)

	resultStr := string(result)
	// Root-level keys should appear at the top
	assert.Contains(t, resultStr, `app_name = "RetroGoLib"`)
	assert.Contains(t, resultStr, `version = "1.0.0"`)
	assert.Contains(t, resultStr, "debug = true")
	assert.Contains(t, resultStr, "max_connections = 100")
	assert.Contains(t, resultStr, "timeout = 30.5")

	// Section keys should be in their sections
	assert.Contains(t, resultStr, "[database]")
	assert.Contains(t, resultStr, "[cache]")
	assert.Contains(t, resultStr, `host = "localhost"`)
	assert.Contains(t, resultStr, "enabled = true")
}
