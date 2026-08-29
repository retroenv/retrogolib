package arch

import (
	"slices"
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
			name: "CPU6502",
			arch: CPU6502,
			want: "6502",
		},
		{
			name: "CPU68000",
			arch: CPU68000,
			want: "68000",
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
		{
			name: "X86",
			arch: X86,
			want: "x86",
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
			name: "CPU6502 is valid",
			arch: CPU6502,
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
			name: "X86 is valid",
			arch: X86,
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
			name: "uppercase Z80 is invalid (IsValid is case-sensitive)",
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
		{"valid 6502", "6502", CPU6502, true},
		{"valid 68000", "68000", CPU68000, true},
		{"valid z80", "z80", Z80, true},
		{"valid chip8", "chip8", CHIP8, true},
		{"valid x86", "x86", X86, true},
		{"invalid architecture", "invalid", "", false},
		{"legacy m68000 identifier is invalid", "m68000", "", false},
		{"empty string", "", "", false},
		{"uppercase Z80 now valid (case-insensitive)", "Z80", Z80, true},
		{"mixed case CHIP8 now valid (case-insensitive)", "CHIP8", CHIP8, true},
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
	expected := []Architecture{CHIP8, CPU6502, CPU65C02, CPU65816, CPU6809, CPU68000, SM83, X86, Z80}

	assert.Len(t, expected, len(got))

	// Check that all expected architectures are present
	for _, expectedArch := range expected {
		found := slices.Contains(got, expectedArch)
		assert.True(t, found, "Expected architecture %s not found in supported architectures", expectedArch)
	}

	// Verify no unexpected architectures are present
	for _, gotArch := range got {
		found := slices.Contains(expected, gotArch)
		assert.True(t, found, "Unexpected architecture %s found in supported architectures", gotArch)
	}
}

func TestConstants(t *testing.T) {
	// Verify the constant values are as expected
	assert.Equal(t, "6502", string(CPU6502))
	assert.Equal(t, "65c02", string(CPU65C02))
	assert.Equal(t, "65816", string(CPU65816))
	assert.Equal(t, "6809", string(CPU6809))
	assert.Equal(t, "68000", string(CPU68000))
	assert.Equal(t, "chip8", string(CHIP8))
	assert.Equal(t, "x86", string(X86))
	assert.Equal(t, "z80", string(Z80))
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
