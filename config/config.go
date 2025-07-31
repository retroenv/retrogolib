// Package config provides configuration file parsing and struct marshaling capabilities.
//
// The config package supports INI-style configuration files with sections, key-value pairs,
// and comments. It provides automatic marshaling and unmarshalling between configuration
// files and Go structs using struct tags.
//
// Basic usage:
//
//	type Config struct {
//		Name string `config:"general.name"`
//		Port int    `config:"server.port,default=8080"`
//	}
//
//	var cfg Config
//	err := config.Load("config.ini", &cfg)
package config

// valueType represents the type of configuration value.
type valueType int

const (
	stringType valueType = iota
	intType
	boolType
	floatType
	hexType
)

// elementType represents the type of structural element.
type elementType int

const (
	commentElement elementType = iota
	sectionElement
	keyValueElement
	emptyLineElement
)

// Config represents a loaded configuration with sections and values.
type Config struct {
	sections  map[string]section
	filename  string
	comments  []comment          // Preserved comments from original file
	structure []structureElement // Original file structure for write operations
}

// section represents a configuration section with key-value pairs.
type section map[string]value

// value represents a configuration value with type information.
type value struct {
	Raw    string
	parsed any
	vtype  valueType
}

// comment represents a comment in the configuration file.
type comment struct {
	Line    int    // Line number where comment appears
	Text    string // Comment text without # prefix
	Section string // Section this comment belongs to (empty for global)
}

// structureElement represents elements in the original file structure.
type structureElement struct {
	Type    elementType // Comment, Section, KeyValue, EmptyLine
	Line    int         // Original line number
	Content string      // Original content
	Section string      // Current section context
	Key     string      // Key name (for KeyValue elements)
}

// tagInfo contains parsed tag information including default values and required flag.
type tagInfo struct {
	Section      string
	Key          string
	DefaultValue string
	HasDefault   bool
	Required     bool
}

// String returns the string representation of valueType.
func (vt valueType) String() string {
	switch vt {
	case stringType:
		return "string"
	case intType:
		return "int"
	case boolType:
		return "bool"
	case floatType:
		return "float"
	case hexType:
		return "hex"
	default:
		return "unknown"
	}
}
