// Package buildinfo formats build information that is embedded into the binaries.
package buildinfo

import (
	"runtime"
	"strings"
)

// Version builds a version string based on binary release information.
func Version(version, commit, date string) string {
	buf := strings.Builder{}
	buf.WriteString(version)

	if commit != "" {
		buf.WriteString(" commit: " + commit)
	}
	if date != "" {
		buf.WriteString(" built at: " + date)
	}
	goVersion := runtime.Version()
	buf.WriteString(" built with: " + goVersion)
	return buf.String()
}
