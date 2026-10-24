package mcp

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// TestPinnedProtocolSchema verifies that the pinned reference copy of the
// official protocol schema is valid JSON. The copy serves as a reference for
// wire shapes; the package does not validate against it at runtime.
func TestPinnedProtocolSchema(t *testing.T) {
	data, err := os.ReadFile("testdata/schema/2026-07-28/schema.json")
	assert.NoError(t, err)

	var schema map[string]any
	assert.NoError(t, json.Unmarshal(data, &schema))

	assert.Equal(t, "https://json-schema.org/draft/2020-12/schema", schema["$schema"])
	assert.NotNil(t, schema["$defs"])
}
