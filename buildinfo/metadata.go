package buildinfo

import (
	"runtime/debug"
	"strings"
)

const (
	developmentVersion       = "dev"
	developmentModuleVersion = "(devel)"
	vcsRevisionSetting       = "vcs.revision"
	vcsTimeSetting           = "vcs.time"
	vcsModifiedSetting       = "vcs.modified"
	vcsModifiedTrue          = "true"
)

// Options supplies metadata that takes precedence over Go build information.
type Options struct {
	Version    string
	Commit     string
	CommitTime string
	BuildDate  string
}

// Info contains release metadata and Go build information.
type Info struct {
	Version string
	Commit  string
	// BuildDate is supplied by the caller. VCSTime is the VCS commit time.
	BuildDate string
	VCSTime   string
	Modified  bool
}

// Read reads metadata from the running binary. Explicit values take precedence.
// An empty version or "dev" uses the module version when it is available.
func Read(options Options) Info {
	info, _ := debug.ReadBuildInfo()
	return Resolve(options, info)
}

// Resolve combines explicit values with Go build information.
// An empty version or "dev" uses the module version, except for "(devel)".
// If no release version is available, Version is "dev".
// BuildDate and VCSTime stay separate because a commit time is not a build date.
// An explicit commit time takes precedence over the embedded VCS time.
// An overridden commit uses the embedded time only if it matches the revision.
// Modified describes the source checkout, including when explicit values are set.
// A nil info supplies no defaults. The function does not change info.
func Resolve(options Options, info *debug.BuildInfo) Info {
	metadata := Info{
		Version:   options.Version,
		Commit:    options.Commit,
		BuildDate: options.BuildDate,
		VCSTime:   options.CommitTime,
	}
	if metadata.Version == "" {
		metadata.Version = developmentVersion
	}

	if info == nil {
		return metadata
	}
	if metadata.Version == developmentVersion && info.Main.Version != "" && info.Main.Version != developmentModuleVersion {
		metadata.Version = info.Main.Version
	}

	var revision, commitTime string
	for _, setting := range info.Settings {
		switch setting.Key {
		case vcsRevisionSetting:
			revision = setting.Value
			if metadata.Commit == "" {
				metadata.Commit = setting.Value
			}

		case vcsTimeSetting:
			commitTime = setting.Value
		case vcsModifiedSetting:
			metadata.Modified = setting.Value == vcsModifiedTrue
		}
	}

	if options.CommitTime == "" && (options.Commit == "" || strings.HasPrefix(revision, options.Commit)) {
		metadata.VCSTime = commitTime
	}

	return metadata
}
