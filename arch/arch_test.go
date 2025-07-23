package arch

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestArchitecture_String(t *testing.T) {
	tests := []struct {
		name string
		arch Architecture
		want string
	}{
		{
			name: "M6502",
			arch: M6502,
			want: "6502",
		},
		{
			name: "Z80",
			arch: Z80,
			want: "z80",
		},
		{
			name: "CHIP8",
			arch: CHIP8,
			want: "chip8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arch.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArchitecture_IsValid(t *testing.T) {
	tests := []struct {
		name string
		arch Architecture
		want bool
	}{
		{
			name: "M6502 is valid",
			arch: M6502,
			want: true,
		},
		{
			name: "Z80 is valid",
			arch: Z80,
			want: true,
		},
		{
			name: "CHIP8 is valid",
			arch: CHIP8,
			want: true,
		},
		{
			name: "empty string is invalid",
			arch: Architecture(""),
			want: false,
		},
		{
			name: "random string is invalid",
			arch: Architecture("invalid"),
			want: false,
		},
		{
			name: "case sensitive - uppercase Z80 is invalid",
			arch: Architecture("Z80"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arch.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   Architecture
		wantOk bool
	}{
		{
			name:   "valid 6502",
			input:  "6502",
			want:   M6502,
			wantOk: true,
		},
		{
			name:   "valid z80",
			input:  "z80",
			want:   Z80,
			wantOk: true,
		},
		{
			name:   "valid chip8",
			input:  "chip8",
			want:   CHIP8,
			wantOk: true,
		},
		{
			name:   "invalid architecture",
			input:  "invalid",
			want:   "",
			wantOk: false,
		},
		{
			name:   "empty string",
			input:  "",
			want:   "",
			wantOk: false,
		},
		{
			name:   "case sensitive - uppercase",
			input:  "Z80",
			want:   "",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := FromString(tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestSupportedArchitectures(t *testing.T) {
	got := SupportedArchitectures()
	expected := []Architecture{M6502, Z80, CHIP8}

	assert.Equal(t, len(expected), len(got))

	// Check that all expected architectures are present
	for _, expectedArch := range expected {
		found := false
		for _, gotArch := range got {
			if gotArch == expectedArch {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected architecture %s not found in supported architectures", expectedArch)
	}
}

func TestConstants(t *testing.T) {
	// Verify the constant values are as expected
	assert.Equal(t, "6502", string(M6502))
	assert.Equal(t, "z80", string(Z80))
	assert.Equal(t, "chip8", string(CHIP8))
}

// Integration test to ensure all supported architectures are valid
func TestAllSupportedArchitecturesAreValid(t *testing.T) {
	supported := SupportedArchitectures()
	for _, arch := range supported {
		assert.True(t, arch.IsValid(), "Supported architecture %s should be valid", arch)
	}
}

// Integration test to ensure FromString works for all supported architectures
func TestFromStringWorksForAllSupported(t *testing.T) {
	supported := SupportedArchitectures()
	for _, arch := range supported {
		got, ok := FromString(arch.String())
		assert.True(t, ok, "FromString should work for supported architecture %s", arch)
		assert.Equal(t, arch, got)
	}
}
