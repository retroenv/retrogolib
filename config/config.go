package config

// ValueType represents the type of configuration value.
type ValueType int

const (
	stringType ValueType = iota
	intType
	boolType
	floatType
	hexType
)

// ElementType represents the type of structural element.
type ElementType int

const (
	commentElement ElementType = iota
	sectionElement
	keyValueElement
	emptyLineElement
)

// Value represents a configuration value with type information.
type Value struct {
	Raw    string
	parsed any
	vtype  ValueType
}

// Section represents a configuration section with key-value pairs.
type Section map[string]Value

// Comment represents a comment in the configuration file.
type Comment struct {
	Line    int    // Line number where comment appears
	Text    string // Comment text without # prefix
	Section string // Section this comment belongs to (empty for global)
}

// StructureElement represents an element in the original file structure.
type StructureElement struct {
	Type    ElementType // Comment, Section, KeyValue, EmptyLine
	Line    int         // Original line number
	Content string      // Original content
	Section string      // Current section context
	Key     string      // Key name (for KeyValue elements)
}

// Config represents a loaded configuration with sections and values.
type Config struct {
	sections  map[string]Section
	filename  string
	comments  []Comment          // Preserved comments from original file
	structure []StructureElement // Original file structure for write operations
}

// TagInfo contains parsed tag information including default values and required flag.
type TagInfo struct {
	Section      string
	Key          string
	DefaultValue string
	HasDefault   bool
	Required     bool
}

// String returns the string representation of ValueType.
func (vt ValueType) String() string {
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
