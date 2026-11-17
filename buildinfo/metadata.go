package buildinfo

import "runtime/debug"

const (
	developmentVersion       = "dev"
	developmentModuleVersion = "(devel)"
	vcsRevisionSetting       = "vcs.revision"
	vcsTimeSetting           = "vcs.time"
	vcsModifiedSetting       = "vcs.modified"
	vcsModifiedTrue          = "true"
)

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
func Read(version, commit, date string) Info {
	info, _ := debug.ReadBuildInfo()
	return Resolve(version, commit, date, info)
}

// Resolve combines explicit values with Go build information.
// An empty version or "dev" uses the module version, except for "(devel)".
// If no release version is available, Version is "dev".
// BuildDate and VCSTime stay separate because a commit time is not a build date.
// Modified describes the source checkout, including when explicit values are set.
// A nil info supplies no defaults. The function does not change info.
func Resolve(version, commit, date string, info *debug.BuildInfo) Info {
	metadata := Info{
		Version:   version,
		Commit:    commit,
		BuildDate: date,
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

	for _, setting := range info.Settings {
		switch setting.Key {
		case vcsRevisionSetting:
			if metadata.Commit == "" {
				metadata.Commit = setting.Value
			}

		case vcsTimeSetting:
			metadata.VCSTime = setting.Value
		case vcsModifiedSetting:
			metadata.Modified = setting.Value == vcsModifiedTrue
		}
	}

	return metadata
}
