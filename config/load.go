package config

import (
	"fmt"
	"io"
	"os"
	"slices"

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

// Open parses a configuration file using options and remembers its path for Save.
// The file is closed before Open returns.
func Open(filename string, options Options) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	defer func() { _ = file.Close() }()

	config, err := Parse(file, options)
	if err != nil {
		return nil, err
	}

	config.filename = filename
	return config, nil
}

// Parse reads configuration using options, limited to 10 MiB and 100,000 lines.
// It does not close reader or associate a filename with the result.
// Call Unmarshal on the result to populate a struct, or Entries to read dynamic keys.
func Parse(reader io.Reader, options Options) (*Config, error) {
	// Read one extra byte to distinguish oversized input from input exactly at the limit.
	data, err := io.ReadAll(io.LimitReader(reader, maxConfigSize+1))
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	return parseConfig(data, options)
}

func parseConfig(data []byte, options Options) (*Config, error) {
	if len(data) > maxConfigSize {
		return nil, fmt.Errorf("%w: %d bytes exceeds limit of %d bytes", ErrConfigTooLarge, len(data), maxConfigSize)
	}

	// Keep the parser's options independent of later changes to the caller's slice.
	options.LiteralSections = slices.Clone(options.LiteralSections)
	if options.CommentPrefixes == "" {
		options.CommentPrefixes = "#"
	}
	config := &Config{
		options:  options,
		sections: make(map[string]Section),
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
