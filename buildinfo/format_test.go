package buildinfo_test

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/buildinfo"
)

func TestFormatMultiline(t *testing.T) {
	tests := []struct {
		name string
		info buildinfo.Info
		want string
	}{
		{
			name: "release metadata",
			info: buildinfo.Info{
				Version: "1.0.14",
				Commit:  testRevision,
				VCSTime: testVCSTime,
			},
			want: "tool version: v1.0.14\ncommit hash: abcdef0\ncommit time: 2026-01-02T03:04:05Z\n",
		},
		{
			name: "modified local build",
			info: buildinfo.Info{
				Version:  testDevelopmentVersion,
				Commit:   testOverrideCommit,
				Modified: true,
			},
			want: "tool version: dev\ncommit hash: abc\nsource: modified\n",
		},
		{
			name: "module and build date",
			info: buildinfo.Info{
				Version:   testModuleVersion,
				BuildDate: testBuildDate,
			},
			want: "tool version: v1.2.3\nbuild date: release-date\n",
		},
		{
			name: "numeric hash version",
			info: buildinfo.Info{Version: "123abcd"},
			want: "tool version: 123abcd\n",
		},
		{
			name: "empty values",
			want: "tool version: \n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, buildinfo.FormatMultiline("tool", test.info))
		})
	}
}
