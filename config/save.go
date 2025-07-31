package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/retroenv/retrogolib/set"
)

// Save saves the configuration to the original file, preserving comments.
func (c *Config) Save() error {
	if c.filename == "" {
		return ErrWriteOnly
	}
	return c.SaveAs(c.filename)
}

// SaveAs saves the configuration to a new file, preserving comments.
func (c *Config) SaveAs(filename string) error {
	data, err := c.SaveBytes()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, data, configFilePermissions); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}
	return nil
}

// SaveBytes returns the configuration as bytes with preserved comments.
func (c *Config) SaveBytes() ([]byte, error) {
	// Estimate buffer size based on structure elements
	estimatedSize := len(c.structure) * avgElementSize
	var buf strings.Builder
	buf.Grow(estimatedSize)

	for _, element := range c.structure {
		switch element.Type {
		case commentElement, emptyLineElement:
			// Preserve comments and empty lines as-is
			buf.WriteString(element.Content)
			buf.WriteByte('\n')

		case sectionElement:
			// Preserve section headers
			buf.WriteString(element.Content)
			buf.WriteByte('\n')

		case keyValueElement:
			// Update value while preserving formatting
			if section, exists := c.sections[element.Section]; exists {
				if value, exists := section[element.Key]; exists {
					// Reconstruct line with updated value
					buf.WriteString(c.formatKeyValue(element.Key, value))
					buf.WriteByte('\n')
				} else {
					// Key was removed, skip this line
					continue
				}
			}
		}
	}

	// Add any new sections/keys that weren't in original file
	c.appendNewContent(&buf)

	return []byte(buf.String()), nil
}

// formatKeyValue formats a key-value pair matching original style.
func (c *Config) formatKeyValue(key string, value value) string {
	switch value.vtype {
	case stringType:
		// Always quote strings to maintain consistency and handle spaces
		return fmt.Sprintf("%s = %q", key, value.Raw)
	case hexType:
		if val, ok := value.parsed.(int); ok {
			return fmt.Sprintf("%s = 0x%X", key, val)
		}
	case boolType:
		return fmt.Sprintf("%s = %t", key, value.parsed.(bool))
	case intType:
		return fmt.Sprintf("%s = %d", key, value.parsed.(int))
	case floatType:
		return fmt.Sprintf("%s = %g", key, value.parsed.(float64))
	}
	return fmt.Sprintf("%s = %s", key, value.Raw)
}

// appendNewContent adds new sections/keys that weren't in the original file.
func (c *Config) appendNewContent(buf *strings.Builder) {
	existingSections := set.New[string]()
	existingKeys := make(map[string]set.Set[string], len(c.sections))

	// Track what already exists in the structure
	for _, element := range c.structure {
		switch element.Type {
		case sectionElement:
			existingSections.Add(element.Section)
			if existingKeys[element.Section] == nil {
				existingKeys[element.Section] = set.New[string]()
			}
		case keyValueElement:
			if existingKeys[element.Section] == nil {
				existingKeys[element.Section] = set.New[string]()
			}
			existingKeys[element.Section].Add(element.Key)
		}
	}

	// Add new sections and keys
	for sectionName, section := range c.sections {
		if !existingSections.Contains(sectionName) {
			// New section
			_, _ = fmt.Fprintf(buf, "\n[%s]\n", sectionName)
			for key, value := range section {
				buf.WriteString(c.formatKeyValue(key, value))
				buf.WriteByte('\n')
			}
		} else {
			// Existing section, check for new keys
			hasNewKeys := false
			for key := range section {
				if existingKeys[sectionName] == nil || !existingKeys[sectionName].Contains(key) {
					if !hasNewKeys {
						// Add spacing before new keys in existing section
						buf.WriteByte('\n')
						hasNewKeys = true
					}
					buf.WriteString(c.formatKeyValue(key, section[key]))
					buf.WriteByte('\n')
				}
			}
		}
	}
}
