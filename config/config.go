package config

import (
	"iter"
	"strings"
)

// Options controls parsing. The zero value preserves case-insensitive names,
// typed values, hash comments, and rejection of repeated sections.
type Options struct {
	CaseSensitive         bool
	RawValues             bool     // Preserve values as strings, including quotes and escape sequences.
	CommentPrefixes       string   // Individual comment characters; empty defaults to "#".
	InlineComments        bool     // Recognize comment characters outside quoted values and after section headers.
	LiteralSections       []string // Preserve comment characters in these sections' values; follows CaseSensitive.
	AllowRepeatedSections bool     // Reopening a section still rejects duplicate keys.
}

// Entry is a parsed key/value pair with its source location.
type Entry struct {
	Section string
	Key     string
	Value   Value
	Line    int
}

// SectionInfo identifies a section header in source order.
type SectionInfo struct {
	Name string
	Line int
}

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
	Text    string // Comment text without its prefix character
	Section string // Section this comment belongs to (empty for global)
}

// StructureElement represents an element in the original file structure.
type StructureElement struct {
	InlineComment string      // Trailing comment retained when inline comments are enabled.
	Type          ElementType // Comment, Section, KeyValue, EmptyLine
	Line          int         // Original line number
	Content       string      // Original content
	Section       string      // Current section context
	Key           string      // Key name (for KeyValue elements)
}

// Config represents a loaded configuration with sections and values.
type Config struct {
	options   Options
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

// Entries iterates loaded entries in source order, using current values.
// Values added by Marshal without a source location are not included.
func (c *Config) Entries() iter.Seq[Entry] {
	return func(yield func(Entry) bool) {
		for _, element := range c.structure {
			if element.Type != keyValueElement {
				continue
			}
			value, ok := c.sections[element.Section][element.Key]
			if ok && !yield(Entry{Section: element.Section, Key: element.Key, Value: value, Line: element.Line}) {
				return
			}
		}
	}
}

// Sections iterates loaded section headers in source order, including reopened sections.
func (c *Config) Sections() iter.Seq[SectionInfo] {
	return func(yield func(SectionInfo) bool) {
		for _, element := range c.structure {
			if element.Type == sectionElement && !yield(SectionInfo{Name: element.Section, Line: element.Line}) {
				return
			}
		}
	}
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

func (c *Config) normalizeName(name string) string {
	if c.options.CaseSensitive {
		return name
	}
	return strings.ToLower(name)
}
