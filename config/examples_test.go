package config_test

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/retroenv/retrogolib/config"
)

// Example demonstrates basic configuration loading.
func ExampleLoad() {
	configData := `[emulation]
cpu = "6502"
speed = 1789773
debug = true

[nes]
code_base_address = 0x8000`

	type AppConfig struct {
		CPU          string `config:"emulation.cpu"`
		Speed        int    `config:"emulation.speed"`
		Debug        bool   `config:"emulation.debug"`
		CodeBaseAddr int    `config:"nes.code_base_address"`
	}

	var cfg AppConfig
	if err := config.LoadBytes([]byte(configData), &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("CPU: %s, Speed: %d, Debug: %t, Address: 0x%X\n",
		cfg.CPU, cfg.Speed, cfg.Debug, cfg.CodeBaseAddr)

	// Output: CPU: 6502, Speed: 1789773, Debug: true, Address: 0x8000
}

// Example demonstrates nested struct configuration.
func ExampleLoad_nestedStruct() {
	configData := `[emulation]
cpu = "6502"
speed = 1789773

[logging]
level = "info"
output = "console"`

	type EmulationConfig struct {
		CPU   string `config:"cpu"`
		Speed int    `config:"speed"`
	}

	type LoggingConfig struct {
		Level  string `config:"level"`
		Output string `config:"output"`
	}

	type AppConfig struct {
		Emulation EmulationConfig `config:"emulation"`
		Logging   LoggingConfig   `config:"logging"`
	}

	var cfg AppConfig
	if err := config.LoadBytes([]byte(configData), &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("CPU: %s, Log Level: %s\n", cfg.Emulation.CPU, cfg.Logging.Level)

	// Output: CPU: 6502, Log Level: info
}

// Example demonstrates comment preservation during write operations.
func ExampleConfig_Marshal() {
	configData := `# RetroGoLib Configuration
# Main emulation settings

[emulation]
# CPU type - currently supports 6502
cpu = "6502"
# Clock speed in Hz
speed = 1789773
debug = false`

	// Load config preserving comments
	configObj, err := config.LoadConfigBytes([]byte(configData))
	if err != nil {
		log.Fatal(err)
	}

	type AppConfig struct {
		CPU   string `config:"emulation.cpu"`
		Speed int    `config:"emulation.speed"`
		Debug bool   `config:"emulation.debug"`
	}

	// Load current values
	var cfg AppConfig
	if err := configObj.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	// Modify values
	cfg.Speed *= 2   // Double the speed
	cfg.Debug = true // Enable debug

	// Marshal back preserving comments
	if err := configObj.Marshal(&cfg); err != nil {
		log.Fatal(err)
	}

	// Save to bytes
	result, err := configObj.SaveBytes()
	if err != nil {
		log.Fatal(err)
	}

	output := string(result)

	// Comments are preserved
	fmt.Println(strings.Contains(output, "# RetroGoLib Configuration"))
	fmt.Println(strings.Contains(output, "# CPU type - currently supports 6502"))

	// Values are updated
	fmt.Println(strings.Contains(output, "speed = 3579546"))
	fmt.Println(strings.Contains(output, "debug = true"))

	// Output: true
	// true
	// true
	// true
}

// Example demonstrates adding new configuration sections and keys.
func ExampleConfig_Marshal_addNew() {
	configData := `[existing]
old_key = "old_value"`

	configObj, err := config.LoadConfigBytes([]byte(configData))
	if err != nil {
		log.Fatal(err)
	}

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

	if err := configObj.Marshal(&extended); err != nil {
		log.Fatal(err)
	}

	result, err := configObj.SaveBytes()
	if err != nil {
		log.Fatal(err)
	}

	output := string(result)

	// Old content preserved
	fmt.Println(strings.Contains(output, `old_key = "old_value"`))

	// New key added to existing section
	fmt.Println(strings.Contains(output, "new_key = \"added_key\""))

	// New section added
	fmt.Println(strings.Contains(output, "[new_section]"))
	fmt.Println(strings.Contains(output, "value = 100"))

	// Output: true
	// true
	// true
	// true
}

// Example demonstrates error handling for type mismatches.
func ExampleUnmarshalError() {
	configData := `[test]
name = 42` // Number instead of string

	type Config struct {
		Name string `config:"test.name"`
	}

	var cfg Config
	err := config.LoadBytes([]byte(configData), &cfg)
	if err != nil {
		var unmarshalErr *config.UnmarshalError
		if errors.As(err, &unmarshalErr) {
			fmt.Printf("Field: %s, Section: %s, Key: %s\n",
				unmarshalErr.Field, unmarshalErr.Section, unmarshalErr.Key)
		}
	}

	// Output: Field: Name, Section: test, Key: name
}

// Example demonstrates configuration value types.
func Example_valueTypes() {
	configData := `[types]
str_quoted = "hello world"
str_unquoted = simple_value
integer = 42
hex_value = 0xFF00
boolean = true
float_val = 3.14159`

	configObj, err := config.LoadConfigBytes([]byte(configData))
	if err != nil {
		log.Fatal(err)
	}

	type TypesConfig struct {
		StrQuoted   string  `config:"types.str_quoted"`
		StrUnquoted string  `config:"types.str_unquoted"`
		Integer     int     `config:"types.integer"`
		HexValue    int     `config:"types.hex_value"`
		Boolean     bool    `config:"types.boolean"`
		FloatVal    float64 `config:"types.float_val"`
	}

	var cfg TypesConfig
	if err := configObj.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("String: %s\n", cfg.StrQuoted)
	fmt.Printf("Unquoted: %s\n", cfg.StrUnquoted)
	fmt.Printf("Integer: %d\n", cfg.Integer)
	fmt.Printf("Hex: 0x%X\n", cfg.HexValue)
	fmt.Printf("Boolean: %t\n", cfg.Boolean)
	fmt.Printf("Float: %g\n", cfg.FloatVal)

	// Output: String: hello world
	// Unquoted: simple_value
	// Integer: 42
	// Hex: 0xFF00
	// Boolean: true
	// Float: 3.14159
}

func ExampleConfig_defaultValues() {
	// Configuration struct with default values
	type GameConfig struct {
		CPU      string  `config:"emulation.cpu,default=6502"`
		Speed    int     `config:"emulation.speed,default=1789773"`
		Debug    bool    `config:"emulation.debug,default=false"`
		Volume   float64 `config:"audio.volume,default=0.8"`
		BaseAddr int     `config:"memory.base_address,default=0x8000"`
		Optional string  `config:"optional.setting"` // No default
	}

	// Example 1: Empty config file - all defaults applied
	emptyConfig := ``

	var cfg1 GameConfig
	if err := config.LoadBytes([]byte(emptyConfig), &cfg1); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Empty config - defaults applied:\n")
	fmt.Printf("CPU: %s, Speed: %d, Debug: %t\n", cfg1.CPU, cfg1.Speed, cfg1.Debug)
	fmt.Printf("Volume: %.1f, BaseAddr: 0x%X\n", cfg1.Volume, cfg1.BaseAddr)
	fmt.Printf("Optional (no default): '%s'\n", cfg1.Optional)

	// Example 2: Partial config - mix of provided values and defaults
	partialConfig := `[emulation]
cpu = "z80"
debug = true

[audio]
volume = 0.5`

	var cfg2 GameConfig
	if err := config.LoadBytes([]byte(partialConfig), &cfg2); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nPartial config - mix of values and defaults:\n")
	fmt.Printf("CPU: %s (provided), Speed: %d (default), Debug: %t (provided)\n",
		cfg2.CPU, cfg2.Speed, cfg2.Debug)
	fmt.Printf("Volume: %.1f (provided), BaseAddr: 0x%X (default)\n",
		cfg2.Volume, cfg2.BaseAddr)

	// Output:
	// Empty config - defaults applied:
	// CPU: 6502, Speed: 1789773, Debug: false
	// Volume: 0.8, BaseAddr: 0x8000
	// Optional (no default): ''
	//
	// Partial config - mix of values and defaults:
	// CPU: z80 (provided), Speed: 1789773 (default), Debug: true (provided)
	// Volume: 0.5 (provided), BaseAddr: 0x8000 (default)
}

func ExampleConfig_requiredFields() {
	// Configuration struct with required fields
	type ServerConfig struct {
		DatabaseURL string `config:"db.url,required"`
		APIKey      string `config:"api.key,required"`
		Port        int    `config:"server.port,required,default=8080"`
		Debug       bool   `config:"debug"`                  // Optional
		LogLevel    string `config:"log.level,default=info"` // Optional with default
	}

	// Example 1: Valid configuration with all required fields
	validConfig := `[db]
url = "postgres://localhost/myapp"

[api]
key = "secret-api-key-123"

[server]
port = 3000`

	var cfg1 ServerConfig
	if err := config.LoadBytes([]byte(validConfig), &cfg1); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Valid config loaded successfully:\n")
	fmt.Printf("Database: %s\n", cfg1.DatabaseURL)
	fmt.Printf("API Key: %s\n", cfg1.APIKey)
	fmt.Printf("Port: %d\n", cfg1.Port)
	fmt.Printf("Debug: %t (optional, zero-value)\n", cfg1.Debug)
	fmt.Printf("Log Level: %s (optional, default applied)\n", cfg1.LogLevel)

	// Example 2: Invalid configuration - missing required field
	invalidConfig := `[db]
url = "postgres://localhost/myapp"

# Missing required api.key
[server]
port = 3000`

	var cfg2 ServerConfig
	if err := config.LoadBytes([]byte(invalidConfig), &cfg2); err != nil {
		var unmarshalErr *config.UnmarshalError
		if errors.As(err, &unmarshalErr) && errors.Is(unmarshalErr.Err, config.ErrRequiredField) {
			fmt.Printf("\nRequired field validation failed:\n")
			fmt.Printf("Missing field: %s in section: %s\n", unmarshalErr.Key, unmarshalErr.Section)
		}
	}

	// Output:
	// Valid config loaded successfully:
	// Database: postgres://localhost/myapp
	// API Key: secret-api-key-123
	// Port: 3000
	// Debug: false (optional, zero-value)
	// Log Level: info (optional, default applied)
	//
	// Required field validation failed:
	// Missing field: key in section: api
}

// ExampleLoad_automaticFieldMapping demonstrates automatic field mapping for untagged fields.
func ExampleLoad_automaticFieldMapping() {
	configData := `name = "retro-emulator"
port = 8080
debug = true

[database]
host = "localhost"
timeout = 30

[cache]
host = "redis-server"
port = 6379`

	// Struct with mixed tagged and untagged fields
	type DatabaseConfig struct {
		Host    string // Automatically maps to database.host
		Timeout int    // Automatically maps to database.timeout
	}

	type CacheConfig struct {
		Host string // Automatically maps to cache.host
		Port int    // Automatically maps to cache.port
	}

	type AppConfig struct {
		Name     string         // Automatically maps to root-level name
		Port     int            // Automatically maps to root-level port
		Debug    bool           // Automatically maps to root-level debug
		Database DatabaseConfig // Uses "database" as section name
		Cache    CacheConfig    // Uses "cache" as section name
	}

	var cfg AppConfig
	if err := config.LoadBytes([]byte(configData), &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("App: %s, Port: %d, Debug: %t\n", cfg.Name, cfg.Port, cfg.Debug)
	fmt.Printf("Database: %s (timeout: %ds)\n", cfg.Database.Host, cfg.Database.Timeout)
	fmt.Printf("Cache: %s:%d\n", cfg.Cache.Host, cfg.Cache.Port)

	// Output: App: retro-emulator, Port: 8080, Debug: true
	// Database: localhost (timeout: 30s)
	// Cache: redis-server:6379
}
