package arch

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestSystem_String(t *testing.T) {
	tests := []struct {
		name   string
		system System
		want   string
	}{
		{
			name:   "CHIP8System",
			system: CHIP8System,
			want:   "chip8",
		},
		{
			name:   "DOS",
			system: DOS,
			want:   "dos",
		},
		{
			name:   "GameBoy",
			system: GameBoy,
			want:   "gameboy",
		},
		{
			name:   "NES",
			system: NES,
			want:   "nes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.system.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSystem_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		system System
		want   bool
	}{
		{
			name:   "CHIP8System is valid",
			system: CHIP8System,
			want:   true,
		},
		{
			name:   "DOS is valid",
			system: DOS,
			want:   true,
		},
		{
			name:   "GameBoy is valid",
			system: GameBoy,
			want:   true,
		},
		{
			name:   "NES is valid",
			system: NES,
			want:   true,
		},
		{
			name:   "empty string is invalid",
			system: System(""),
			want:   false,
		},
		{
			name:   "random string is invalid",
			system: System("invalid"),
			want:   false,
		},
		{
			name:   "uppercase DOS is invalid (IsValid is case-sensitive)",
			system: System("DOS"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.system.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSystemFromString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   System
		wantOk bool
	}{
		{"valid chip8", "chip8", CHIP8System, true},
		{"valid dos", "dos", DOS, true},
		{"valid gameboy", "gameboy", GameBoy, true},
		{"valid nes", "nes", NES, true},
		{"invalid system", "invalid", "", false},
		{"empty string", "", "", false},
		{"uppercase DOS now valid (case-insensitive)", "DOS", DOS, true},
		{"mixed case GameBoy now valid (case-insensitive)", "GAMEBOY", GameBoy, true},
		{"mixed case NES now valid (case-insensitive)", "NES", NES, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := SystemFromString(tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestSupportedSystems(t *testing.T) {
	got := SupportedSystems()
	expected := []System{CHIP8System, DOS, GameBoy, NES}

	assert.Equal(t, len(expected), len(got))

	// Check that all expected systems are present
	for _, expectedSys := range expected {
		found := false
		for _, gotSys := range got {
			if gotSys == expectedSys {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected system %s not found in supported systems", expectedSys)
	}

	// Verify no unexpected systems are present
	for _, gotSys := range got {
		found := false
		for _, expectedSys := range expected {
			if gotSys == expectedSys {
				found = true
				break
			}
		}
		assert.True(t, found, "Unexpected system %s found in supported systems", gotSys)
	}
}

func TestSystemConstants(t *testing.T) {
	// Verify the constant values are as expected
	assert.Equal(t, "chip8", string(CHIP8System))
	assert.Equal(t, "dos", string(DOS))
	assert.Equal(t, "gameboy", string(GameBoy))
	assert.Equal(t, "nes", string(NES))
}

// Integration test to ensure all supported systems are valid
func TestAllSupportedSystemsAreValid(t *testing.T) {
	supported := SupportedSystems()
	for _, sys := range supported {
		assert.True(t, sys.IsValid(), "Supported system %s should be valid", sys)
	}
}

// Integration test to ensure SystemFromString works for all supported systems
func TestSystemFromStringWorksForAllSupported(t *testing.T) {
	supported := SupportedSystems()
	for _, sys := range supported {
		got, ok := SystemFromString(sys.String())
		assert.True(t, ok, "SystemFromString should work for supported system %s", sys)
		assert.Equal(t, sys, got)
	}
}
