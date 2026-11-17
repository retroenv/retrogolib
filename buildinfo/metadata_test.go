package buildinfo_test

import (
	"runtime"
	"runtime/debug"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/buildinfo"
)

// Test markers are separate from the implementation constants.
const (
	testDevelopmentVersion       = "dev"
	testDevelopmentModuleVersion = "(devel)"
	testModuleVersion            = "v1.2.3"
	testOverrideVersion          = "v2.0.0"
	testRevision                 = "abcdef0123456789"
	testOverrideCommit           = "abc"
	testBuildDate                = "release-date"
	testVCSTime                  = "2026-01-02T03:04:05Z"
	testVCSModifiedSetting       = "vcs.modified"
)

func TestResolveVersions(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		date    string
		info    *debug.BuildInfo
		want    buildinfo.Info
	}{
		{
			name: "no metadata",
			want: buildinfo.Info{Version: testDevelopmentVersion},
		},
		{
			name:    "explicit values without build information",
			version: testOverrideVersion,
			commit:  testOverrideCommit,
			date:    testBuildDate,
			want: buildinfo.Info{
				Version:   testOverrideVersion,
				Commit:    testOverrideCommit,
				BuildDate: testBuildDate,
			},
		},
		{
			name: "empty build information",
			info: &debug.BuildInfo{},
			want: buildinfo.Info{Version: testDevelopmentVersion},
		},
		{
			name: "development module",
			info: &debug.BuildInfo{Main: debug.Module{Version: testDevelopmentModuleVersion}},
			want: buildinfo.Info{Version: testDevelopmentVersion},
		},
		{
			name: "installed module",
			info: &debug.BuildInfo{Main: debug.Module{Version: testModuleVersion}},
			want: buildinfo.Info{Version: testModuleVersion},
		},
		{
			name:    "dev uses module version",
			version: testDevelopmentVersion,
			info:    &debug.BuildInfo{Main: debug.Module{Version: testModuleVersion}},
			want:    buildinfo.Info{Version: testModuleVersion},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, buildinfo.Resolve(test.version, test.commit, test.date, test.info))
		})
	}
}

func TestResolveOverrides(t *testing.T) {
	info := testBuildInfo()
	tests := []struct {
		name    string
		version string
		commit  string
		date    string
		info    *debug.BuildInfo
		want    buildinfo.Info
	}{
		{
			name:    "explicit version preserves other defaults",
			version: testOverrideVersion,
			info:    info,
			want: buildinfo.Info{
				Version:  testOverrideVersion,
				Commit:   testRevision,
				VCSTime:  testVCSTime,
				Modified: true,
			},
		},
		{
			name:   "explicit commit preserves other defaults",
			commit: testOverrideCommit,
			info:   info,
			want: buildinfo.Info{
				Version:  testModuleVersion,
				Commit:   testOverrideCommit,
				VCSTime:  testVCSTime,
				Modified: true,
			},
		},
		{
			name: "explicit date preserves VCS time",
			date: testBuildDate,
			info: info,
			want: buildinfo.Info{
				Version:   testModuleVersion,
				Commit:    testRevision,
				BuildDate: testBuildDate,
				VCSTime:   testVCSTime,
				Modified:  true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, buildinfo.Resolve(test.version, test.commit, test.date, test.info))
		})
	}
	assert.Equal(t, testModuleVersion, info.Main.Version)
	assert.Equal(t, testRevision, info.Settings[0].Value)
}

func TestResolveExplicitValues(t *testing.T) {
	info := testBuildInfo()
	want := buildinfo.Info{
		Version:   testOverrideVersion,
		Commit:    testOverrideCommit,
		BuildDate: testBuildDate,
		VCSTime:   testVCSTime,
		Modified:  true,
	}
	assert.Equal(t, want, buildinfo.Resolve(testOverrideVersion, testOverrideCommit, testBuildDate, info))
}

func TestResolveVCS(t *testing.T) {
	info := testBuildInfo()
	tests := []struct {
		name    string
		version string
		commit  string
		date    string
		info    *debug.BuildInfo
		want    buildinfo.Info
	}{
		{
			name: "local modified build",
			info: &debug.BuildInfo{
				Main:     debug.Module{Version: testDevelopmentModuleVersion},
				Settings: info.Settings,
			},
			want: buildinfo.Info{
				Version:  testDevelopmentVersion,
				Commit:   testRevision,
				VCSTime:  testVCSTime,
				Modified: true,
			},
		},
		{
			name: "clean checkout",
			info: &debug.BuildInfo{Settings: []debug.BuildSetting{
				{
					Key:   testVCSModifiedSetting,
					Value: "false",
				},
			}},
			want: buildinfo.Info{Version: testDevelopmentVersion},
		},
		{
			name: "invalid modified setting",
			info: &debug.BuildInfo{Settings: []debug.BuildSetting{
				{
					Key:   testVCSModifiedSetting,
					Value: "invalid",
				},
			}},
			want: buildinfo.Info{Version: testDevelopmentVersion},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, buildinfo.Resolve(test.version, test.commit, test.date, test.info))
		})
	}
	assert.Equal(t, testModuleVersion, info.Main.Version)
	assert.Equal(t, testRevision, info.Settings[0].Value)
}

func TestRead(t *testing.T) {
	info, _ := debug.ReadBuildInfo()
	want := buildinfo.Resolve(testDevelopmentVersion, "", "", info)
	assert.Equal(t, want, buildinfo.Read(testDevelopmentVersion, "", ""))
	metadata := buildinfo.Read(testOverrideVersion, testOverrideCommit, testBuildDate)
	assert.Equal(t, testOverrideVersion, metadata.Version)
	assert.Equal(t, testOverrideCommit, metadata.Commit)
	assert.Equal(t, testBuildDate, metadata.BuildDate)
}

func TestVersion_FormatCompatibility(t *testing.T) {
	want := "v1.0.0 commit: abc built at: release-date built with: " + runtime.Version()
	assert.Equal(t, want, buildinfo.Version(testVersion, testOverrideCommit, testBuildDate))
}

func testBuildInfo() *debug.BuildInfo {
	return &debug.BuildInfo{
		Main: debug.Module{Version: testModuleVersion},
		Settings: []debug.BuildSetting{
			{
				Key:   "vcs.revision",
				Value: testRevision,
			},
			{
				Key:   "vcs.time",
				Value: testVCSTime,
			},
			{
				Key:   testVCSModifiedSetting,
				Value: "true",
			},
			{
				Key:   "unrelated",
				Value: "ignored",
			},
		},
	}
}
