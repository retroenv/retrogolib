package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestStepConsumesStallCyclesBeforeExecuting(t *testing.T) {
	t.Parallel()

	cpu := newProgramCPU(t, VariantNMOS6502, 0x8000, 0xea)
	cpu.StallCycles(2)

	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x8000), cpu.PC)
	assert.Equal(t, initialCycles+1, cpu.cycles)

	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x8000), cpu.PC)
	assert.Equal(t, initialCycles+2, cpu.cycles)

	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x8001), cpu.PC)
	assert.Equal(t, initialCycles+4, cpu.cycles)
}

func TestPageCrossCycleDoesNotRequireTracing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts []Option
	}{
		{name: "without tracing"},
		{name: "with tracing", opts: []Option{WithTracing()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cpu := newProgramCPU(t, VariantNMOS6502, 0x8000, 0xbd, 0xff, 0x20)
			cpu.X = 1
			cpu.memory.Write(0x2100, 0x42)
			cpu.opts = newOptions(tt.opts...)

			assert.NoError(t, cpu.Step())

			assert.Equal(t, uint8(0x42), cpu.A)
			assert.Equal(t, initialCycles+5, cpu.cycles)
		})
	}
}

func TestAbsoluteXRMWTiming(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		variant    CPUVariant
		base       uint16
		wantCycles uint64
	}{
		{name: "NMOS same page", variant: VariantNMOS6502, base: 0x2000, wantCycles: 7},
		{name: "NMOS page cross", variant: VariantNMOS6502, base: 0x20ff, wantCycles: 7},
		{name: "65C02 same page", variant: Variant65C02, base: 0x2000, wantCycles: 6},
		{name: "65C02 page cross", variant: Variant65C02, base: 0x20ff, wantCycles: 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cpu := newProgramCPU(t, tt.variant, 0x8000, 0x5e, byte(tt.base), byte(tt.base>>8))
			cpu.X = 1
			cpu.memory.Write(tt.base+1, 0x04)

			assert.NoError(t, cpu.Step())

			assert.Equal(t, uint8(0x02), cpu.memory.Read(tt.base+1))
			assert.Equal(t, initialCycles+tt.wantCycles, cpu.cycles)
		})
	}
}

func Test65C02BranchTiming(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pc         uint16
		program    []byte
		wantPC     uint16
		wantCycles uint64
	}{
		{name: "BRA same page", pc: 0x8000, program: []byte{0x80, 0x01}, wantPC: 0x8003, wantCycles: 3},
		{name: "BRA page cross", pc: 0x80fd, program: []byte{0x80, 0x01}, wantPC: 0x8100, wantCycles: 4},
		{name: "BBR page cross", pc: 0x80fc, program: []byte{0x0f, 0x10, 0x01}, wantPC: 0x8100, wantCycles: 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cpu := newProgramCPU(t, Variant65C02, tt.pc, tt.program...)

			assert.NoError(t, cpu.Step())

			assert.Equal(t, tt.wantPC, cpu.PC)
			assert.Equal(t, initialCycles+tt.wantCycles, cpu.cycles)
		})
	}
}

func newProgramCPU(t *testing.T, variant CPUVariant, pc uint16, program ...byte) *CPU {
	t.Helper()

	memory, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	memory.WriteWord(ResetAddress, pc)
	for i, value := range program {
		memory.Write(pc+uint16(i), value)
	}
	return New(memory, WithVariant(variant))
}
