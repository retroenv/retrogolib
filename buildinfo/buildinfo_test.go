package buildinfo_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/buildinfo"
)

const testVersion = "v1.0.0"

func TestVersion_Empty(t *testing.T) {
	result := buildinfo.Version("", "", "")

	// Should always include Go version
	goVersion := runtime.Version()
	assert.True(t, strings.Contains(result, goVersion), "Should contain Go version")
	assert.True(t, strings.Contains(result, "built with:"), "Should contain 'built with:' prefix")

	// Should not contain commit or date info
	assert.False(t, strings.Contains(result, "commit:"), "Should not contain commit info")
	assert.False(t, strings.Contains(result, "built at:"), "Should not contain build date info")
}

func TestVersion_OnlyVersion(t *testing.T) {
	version := testVersion
	result := buildinfo.Version(version, "", "")

	assert.True(t, strings.HasPrefix(result, version), "Should start with version")
	assert.True(t, strings.Contains(result, "built with:"), "Should contain Go version")
	assert.False(t, strings.Contains(result, "commit:"), "Should not contain commit info")
	assert.False(t, strings.Contains(result, "built at:"), "Should not contain build date info")
}

func TestVersion_VersionAndCommit(t *testing.T) {
	version := testVersion
	commit := "abc123def"
	result := buildinfo.Version(version, commit, "")

	assert.True(t, strings.HasPrefix(result, version), "Should start with version")
	assert.True(t, strings.Contains(result, "commit: "+commit), "Should contain commit hash")
	assert.True(t, strings.Contains(result, "built with:"), "Should contain Go version")
	assert.False(t, strings.Contains(result, "built at:"), "Should not contain build date info")
}

func TestVersion_VersionAndDate(t *testing.T) {
	version := testVersion
	date := "2024-01-15T10:30:00Z"
	result := buildinfo.Version(version, "", date)

	assert.True(t, strings.HasPrefix(result, version), "Should start with version")
	assert.True(t, strings.Contains(result, "built at: "+date), "Should contain build date")
	assert.True(t, strings.Contains(result, "built with:"), "Should contain Go version")
	assert.False(t, strings.Contains(result, "commit:"), "Should not contain commit info")
}

func TestVersion_AllFields(t *testing.T) {
	version := "v1.2.3"
	commit := "deadbeef"
	date := "2024-01-15T10:30:00Z"
	result := buildinfo.Version(version, commit, date)

	assert.True(t, strings.HasPrefix(result, version), "Should start with version")
	assert.True(t, strings.Contains(result, "commit: "+commit), "Should contain commit hash")
	assert.True(t, strings.Contains(result, "built at: "+date), "Should contain build date")
	assert.True(t, strings.Contains(result, "built with:"), "Should contain Go version")

	// Verify order: version, commit, date, go version
	versionIdx := strings.Index(result, version)
	commitIdx := strings.Index(result, "commit:")
	dateIdx := strings.Index(result, "built at:")
	goIdx := strings.Index(result, "built with:")

	assert.True(t, versionIdx < commitIdx, "Version should come before commit")
	assert.True(t, commitIdx < dateIdx, "Commit should come before date")
	assert.True(t, dateIdx < goIdx, "Date should come before Go version")
}

func TestVersion_EmptyStrings(t *testing.T) {
	// Test with explicit empty strings
	result := buildinfo.Version("", "", "")

	goVersion := runtime.Version()
	expected := " built with: " + goVersion
	assert.Equal(t, expected, result, "Should only contain Go version when all fields empty")
}

func TestVersion_WhitespaceCommit(t *testing.T) {
	version := testVersion
	commit := "   "
	result := buildinfo.Version(version, commit, "")

	// Whitespace-only commit should be treated as non-empty
	assert.True(t, strings.Contains(result, "commit: "+commit), "Should include whitespace commit")
}

func TestVersion_SpecialCharacters(t *testing.T) {
	version := "v1.0.0-beta+special"
	commit := "abc123-def456"
	date := "2024-01-15T10:30:00+02:00"
	result := buildinfo.Version(version, commit, date)

	assert.True(t, strings.Contains(result, version), "Should handle special version characters")
	assert.True(t, strings.Contains(result, commit), "Should handle hyphenated commit")
	assert.True(t, strings.Contains(result, date), "Should handle timezone in date")
}

func TestVersion_LongStrings(t *testing.T) {
	version := strings.Repeat("v", 100)
	commit := strings.Repeat("a", 200)
	date := strings.Repeat("2024", 50)
	result := buildinfo.Version(version, commit, date)

	assert.True(t, strings.Contains(result, version), "Should handle long version string")
	assert.True(t, strings.Contains(result, commit), "Should handle long commit string")
	assert.True(t, strings.Contains(result, date), "Should handle long date string")
	assert.True(t, len(result) > 300, "Result should be appropriately long")
}
