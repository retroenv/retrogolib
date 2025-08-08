package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestFlags_GetFlags(t *testing.T) {
	tests := []struct {
		name  string
		flags Flags
		want  uint8
	}{
		{
			name:  "all flags clear",
			flags: Flags{},
			want:  0x00,
		},
		{
			name:  "carry flag set",
			flags: Flags{C: 1},
			want:  0x01,
		},
		{
			name:  "zero flag set",
			flags: Flags{Z: 1},
			want:  0x40,
		},
		{
			name:  "sign flag set",
			flags: Flags{S: 1},
			want:  0x80,
		},
		{
			name:  "all flags set",
			flags: Flags{C: 1, N: 1, P: 1, X: 1, H: 1, Y: 1, Z: 1, S: 1},
			want:  0xFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := &CPU{Flags: tt.flags}
			got := cpu.GetFlags()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCPU_SetFlags(t *testing.T) {
	tests := []struct {
		name  string
		input uint8
		want  Flags
	}{
		{
			name:  "all clear",
			input: 0x00,
			want:  Flags{},
		},
		{
			name:  "carry set",
			input: 0x01,
			want:  Flags{C: 1},
		},
		{
			name:  "zero set",
			input: 0x40,
			want:  Flags{Z: 1},
		},
		{
			name:  "sign set",
			input: 0x80,
			want:  Flags{S: 1},
		},
		{
			name:  "all set",
			input: 0xFF,
			want:  Flags{C: 1, N: 1, P: 1, X: 1, H: 1, Y: 1, Z: 1, S: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := &CPU{}
			cpu.setFlags(tt.input)
			assert.Equal(t, tt.want, cpu.Flags)
		})
	}
}

func TestCalculateParity(t *testing.T) {
	tests := []struct {
		name  string
		value uint8
		want  bool
	}{
		{"zero has even parity", 0x00, true},
		{"0x01 has odd parity", 0x01, false},
		{"0x03 has even parity", 0x03, true},
		{"0x07 has odd parity", 0x07, false},
		{"0xFF has even parity", 0xFF, true},
		{"0x80 has odd parity", 0x80, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateParity(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Benchmark critical flag operations for performance validation
func BenchmarkCPU_GetFlags(b *testing.B) {
	cpu := &CPU{
		Flags: Flags{C: 1, N: 1, P: 1, X: 1, H: 1, Y: 1, Z: 1, S: 1},
	}

	b.ResetTimer()
	for range b.N {
		cpu.GetFlags()
	}
}

func BenchmarkCPU_SetFlags(b *testing.B) {
	cpu := &CPU{}

	b.ResetTimer()
	for range b.N {
		cpu.setFlags(0xFF)
	}
}

func BenchmarkCalculateParity(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		calculateParity(0x55) // Alternating bits for realistic calculation
	}
}

func BenchmarkSetFlag(b *testing.B) {
	var flag uint8

	b.ResetTimer()
	for i := range b.N {
		setFlag(&flag, i%2 == 0)
	}
}
