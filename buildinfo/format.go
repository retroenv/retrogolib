package buildinfo

import (
	"fmt"
	"strings"
)

const shortCommitLength = 7

// FormatMultiline formats metadata with one value per line.
// It adds a v prefix to numeric versions and shows up to seven commit characters.
// Empty optional fields are omitted. The returned string ends with a newline.
func FormatMultiline(name string, info Info) string {
	var output strings.Builder
	fmt.Fprintf(&output, "%s version: %s\n", name, displayVersion(info.Version))
	if info.Commit != "" {
		fmt.Fprintf(&output, "commit hash: %s\n", info.Commit[:min(len(info.Commit), shortCommitLength)])
	}

	if info.VCSTime != "" {
		fmt.Fprintf(&output, "commit time: %s\n", info.VCSTime)
	}

	if info.BuildDate != "" {
		fmt.Fprintf(&output, "build date: %s\n", info.BuildDate)
	}

	if info.Modified {
		output.WriteString("source: modified\n")
	}

	return output.String()
}

func displayVersion(version string) string {
	major, remainder, found := strings.Cut(version, ".")
	if found && major != "" && strings.Trim(major, "0123456789") == "" &&
		len(remainder) > 0 && remainder[0] >= '0' && remainder[0] <= '9' {

		return "v" + version
	}

	return version
}
