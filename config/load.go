package config

import (
	"fmt"
	"os"

	"github.com/retroenv/retrogolib/set"
)

const (
	// maxConfigSize limits configuration file size to prevent memory exhaustion (10MB).
	maxConfigSize = 10 * 1024 * 1024
	// maxLines limits the number of lines to prevent memory exhaustion.
	maxLines = 100000
	// maxNameLength limits the maximum length for section and key names.
	maxNameLength = 256
	// avgElementSize is the estimated average characters per structure element for buffer sizing.
	avgElementSize = 40
	// configFilePermissions defines the file permissions for saved configuration files.
	configFilePermissions = 0644
)

// Load loads configuration from a file and unmarshalls it into the provided struct.
func Load(filename string, v any) error {
	config, err := LoadConfig(filename)
	if err != nil {
		return err
	}
	return config.Unmarshal(v)
}

// LoadBytes loads configuration from byte slice and unmarshalls it into the provided struct.
func LoadBytes(data []byte, v any) error {
	config, err := LoadConfigBytes(data)
	if err != nil {
		return err
	}
	return config.Unmarshal(v)
}

// LoadConfig loads configuration from a file, returning a Config for advanced operations.
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	config, err := LoadConfigBytes(data)
	if err != nil {
		return nil, err
	}

	config.filename = filename
	return config, nil
}

// LoadConfigBytes loads configuration from byte slice, returning a Config for advanced operations.
func LoadConfigBytes(data []byte) (*Config, error) {
	if len(data) > maxConfigSize {
		return nil, fmt.Errorf("%w: %d bytes exceeds limit of %d bytes", ErrConfigTooLarge, len(data), maxConfigSize)
	}

	config := &Config{
		sections: make(map[string]section),
	}

	parser := &parser{
		data:           data,
		config:         config,
		currentSection: "",
		seenItems:      set.New[string](),
		itemLines:      make(map[string]int),
	}

	if err := parser.parse(); err != nil {
		return nil, err
	}

	return config, nil
}
