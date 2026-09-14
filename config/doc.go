// Package config provides configuration loading and saving with comment preservation.
//
// The config package implements a simple, zero-dependency configuration format
// similar to TOML/INI that supports sections, comments, and basic types while
// preserving all original formatting and comments during write operations.
//
// # Basic Usage
//
// Open a configuration file and unmarshal it into a struct:
//
//	type AppConfig struct {
//	    CPU   string `config:"emulation.cpu"`
//	    Speed int    `config:"emulation.speed"`
//	    Debug bool   `config:"emulation.debug"`
//	}
//
//	var cfg AppConfig
//	document, err := config.Open("app.conf", config.Options{})
//	if err != nil {
//	    return err
//	}
//	if err := document.Unmarshal(&cfg); err != nil {
//	    return err
//	}
//
// # Parser Options and Dynamic Entries
//
// Open accepts a filename; Parse accepts an io.Reader, including strings.NewReader
// or bytes.NewReader for in-memory input. Both accept Options, whose zero value
// retains case-insensitive defaults. Options{CaseSensitive: true} preserves and compares
// section names, keys, struct tags, and automatic field names exactly.
//
// Entries and Sections iterate loaded source entries and headers in file order,
// with line numbers. They do not expose mutable internal maps or include new
// fields added by Marshal without a source location.
//
// For INI-style symbol files, RawValues retains quotes and escape sequences;
// CommentPrefixes: ";#" enables either comment character; InlineComments allows
// trailing comments outside double quotes. LiteralSections exempts selected
// sections' values from inline comments, using the configured case policy.
// AllowRepeatedSections permits reopening a section but still rejects duplicate
// keys. The zero Options value keeps the original parsing behavior. Parse
// accepts a leading UTF-8 BOM and enforces the existing 10 MiB/100,000-line limits.
//
// # Comment Preservation
//
// For write operations that preserve comments and formatting:
//
//	configObj, err := config.Open("app.conf", config.Options{})
//	if err != nil {
//	    return err
//	}
//
//	// Modify configuration
//	var cfg AppConfig
//	configObj.Unmarshal(&cfg)
//	cfg.Speed = 2000000
//	configObj.Marshal(&cfg)
//
//	// Save preserving all original comments
//	configObj.Save()
//
// # Configuration Format
//
// The configuration format supports sections, key-value pairs, and comments:
//
//	# RetroGoLib Configuration
//	# Comments start with #
//
//	[emulation]
//	cpu = "6502"
//	speed = 1789773
//	debug = true
//
//	[nes]
//	code_base_address = 0x8000
//	ram_end_address = 0x0FFF
//
// # Supported Types
//
// - String: quoted or unquoted values
// - Integer: decimal numbers
// - Boolean: true/false
// - Float: decimal numbers with dot
// - Hex: numbers with 0x prefix
//
// # Struct Tags
//
// Use config struct tags to map fields to configuration keys:
//
//	type Config struct {
//	    Name    string `config:"section.key"`     // Maps to [section] key = value
//	    Timeout int    `config:"timeout"`         // Maps to current section
//	    Section struct {
//	        Value string `config:"value"`         // Maps to nested section
//	    } `config:"section"`
//	}
//
// # Automatic Field Mapping
//
// Fields without config tags are automatically mapped using lowercase field names:
//
//	type AppConfig struct {
//	    Name     string          // Automatically maps to global.name
//	    Port     int             // Automatically maps to global.port
//	    Database DatabaseConfig  // Uses "database" as section name
//	}
//
//	type DatabaseConfig struct {
//	    Host    string // Automatically maps to database.host
//	    Timeout int    // Automatically maps to database.timeout
//	}
//
// Automatic mapping rules:
// - Simple fields map to "global.fieldname" (all lowercase)
// - Nested struct fields use the field name as the section name
// - Fields with `config:"-"` are ignored
// - Explicit struct tags override automatic mapping
//
// # Default Values
//
// Specify default values in struct tags for fields that may be missing from configuration:
//
//	type AppConfig struct {
//	    CPU        string  `config:"emulation.cpu,default=6502"`
//	    Speed      int     `config:"emulation.speed,default=1789773"`
//	    Debug      bool    `config:"emulation.debug,default=false"`
//	    Volume     float64 `config:"audio.volume,default=0.8"`
//	    BaseAddr   int     `config:"nes.base_address,default=0x8000"`
//	    Required   string  `config:"required.field"`                    // No default, must be present
//	}
//
// Default values are applied when:
// - The configuration section doesn't exist
// - The configuration key is missing from its section
// - An empty configuration file is loaded
//
// Supported default value formats:
// - String: any text value
// - Integer: decimal numbers (e.g., "123")
// - Boolean: "true" or "false"
// - Float: decimal numbers with dot (e.g., "3.14")
// - Hex: hexadecimal with 0x prefix (e.g., "0xFF")
//
// # Required Fields
//
// Mark fields as required to enforce their presence during configuration loading:
//
//	type AppConfig struct {
//	    DatabaseURL string  `config:"db.url,required"`
//	    APIKey      string  `config:"api.key,required"`
//	    Port        int     `config:"server.port,required,default=8080"`
//	    Optional    string  `config:"optional.setting"`                    // Not required
//	}
//
// Required field validation:
// - Returns UnmarshalError with ErrRequiredField for missing required fields
// - Works with default values (field is required but uses default if missing)
// - Validates during Config.Unmarshal()
// - Provides clear error messages indicating which field and section are missing
//
// Example error handling:
//
//	var cfg AppConfig
//	if err := document.Unmarshal(&cfg); err != nil {
//	    var unmarshalErr *config.UnmarshalError
//	    if errors.As(err, &unmarshalErr) && errors.Is(unmarshalErr.Err, config.ErrRequiredField) {
//	        log.Printf("Required field missing: %s in section %s", unmarshalErr.Key, unmarshalErr.Section)
//	    }
//	}
package config
