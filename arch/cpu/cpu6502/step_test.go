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

	// Only the 65C02 shift and rotate instructions save a cycle without a
	// page crossing. INC and DEC keep the fixed NMOS timing.
	tests := []struct {
		name       string
		variant    CPUVariant
		opcode     byte
		base       uint16
		wantValue  uint8
		wantCycles uint64
	}{
		{name: "NMOS LSR same page", variant: VariantNMOS6502, opcode: 0x5e, base: 0x2000, wantValue: 0x02, wantCycles: 7},
		{name: "NMOS LSR page cross", variant: VariantNMOS6502, opcode: 0x5e, base: 0x20ff, wantValue: 0x02, wantCycles: 7},
		{name: "NMOS INC same page", variant: VariantNMOS6502, opcode: 0xfe, base: 0x2000, wantValue: 0x05, wantCycles: 7},
		{name: "NMOS DEC page cross", variant: VariantNMOS6502, opcode: 0xde, base: 0x20ff, wantValue: 0x03, wantCycles: 7},
		{name: "65C02 LSR same page", variant: Variant65C02, opcode: 0x5e, base: 0x2000, wantValue: 0x02, wantCycles: 6},
		{name: "65C02 LSR page cross", variant: Variant65C02, opcode: 0x5e, base: 0x20ff, wantValue: 0x02, wantCycles: 7},
		{name: "65C02 INC same page", variant: Variant65C02, opcode: 0xfe, base: 0x2000, wantValue: 0x05, wantCycles: 7},
		{name: "65C02 INC page cross", variant: Variant65C02, opcode: 0xfe, base: 0x20ff, wantValue: 0x05, wantCycles: 7},
		{name: "65C02 DEC same page", variant: Variant65C02, opcode: 0xde, base: 0x2000, wantValue: 0x03, wantCycles: 7},
		{name: "65C02 DEC page cross", variant: Variant65C02, opcode: 0xde, base: 0x20ff, wantValue: 0x03, wantCycles: 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cpu := newProgramCPU(t, tt.variant, 0x8000, tt.opcode, byte(tt.base), byte(tt.base>>8))
			cpu.X = 1
			cpu.memory.Write(tt.base+1, 0x04)

			assert.NoError(t, cpu.Step())

			assert.Equal(t, tt.wantValue, cpu.memory.Read(tt.base+1))
			assert.Equal(t, initialCycles+tt.wantCycles, cpu.cycles)
		})
	}
}

func Test65C02DecimalArithmeticTiming(t *testing.T) {
	t.Parallel()

	// The 65C02 takes one more cycle to correct a decimal ADC or SBC result.
	tests := []struct {
		name       string
		variant    CPUVariant
		opcode     byte
		decimal    uint8
		wantCycles uint64
	}{
		{name: "NMOS ADC binary", variant: VariantNMOS6502, opcode: 0x69, wantCycles: 2},
		{name: "NMOS ADC decimal", variant: VariantNMOS6502, opcode: 0x69, decimal: 1, wantCycles: 2},
		{name: "NMOS SBC decimal", variant: VariantNMOS6502, opcode: 0xe9, decimal: 1, wantCycles: 2},
		{name: "65C02 ADC binary", variant: Variant65C02, opcode: 0x69, wantCycles: 2},
		{name: "65C02 ADC decimal", variant: Variant65C02, opcode: 0x69, decimal: 1, wantCycles: 3},
		{name: "65C02 SBC binary", variant: Variant65C02, opcode: 0xe9, wantCycles: 2},
		{name: "65C02 SBC decimal", variant: Variant65C02, opcode: 0xe9, decimal: 1, wantCycles: 3},
		{name: "65C02 ADC abs decimal", variant: Variant65C02, opcode: 0x6d, decimal: 1, wantCycles: 5},
		{name: "65C02 SBC abs decimal", variant: Variant65C02, opcode: 0xed, decimal: 1, wantCycles: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cpu := newProgramCPU(t, tt.variant, 0x8000, tt.opcode, 0x05, 0x00)
			cpu.A = 0x01
			cpu.Flags.D = tt.decimal

			assert.NoError(t, cpu.Step())

			assert.Equal(t, initialCycles+tt.wantCycles, cpu.cycles)
		})
	}
}

func Test65C02NopAbsoluteTiming(t *testing.T) {
	t.Parallel()

	// The W65C02S datasheet lists $5C as a 3-byte, 8-cycle NOP. The emulator-generated
	// wdc65c02 single-step corpus records 4 cycles; the datasheet value is used here.
	cpu := newProgramCPU(t, Variant65C02, 0x8000, 0x5c, 0x34, 0x12)

	assert.NoError(t, cpu.Step())

	assert.Equal(t, uint16(0x8003), cpu.PC)
	assert.Equal(t, initialCycles+8, cpu.cycles)
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

func TestSelfBranchPageCrossTiming(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		variant    CPUVariant
		pc         uint16
		program    []byte
		wantCycles uint64
	}{
		{name: "NMOS same page", pc: 0x80fd, program: []byte{0xd0, 0xfe}, wantCycles: 3},
		{name: "NMOS last two bytes", pc: 0x80fe, program: []byte{0xd0, 0xfe}, wantCycles: 4},
		{name: "NMOS operand on next page", pc: 0x80ff, program: []byte{0xd0, 0xfe}, wantCycles: 4},
		{name: "NMOS address wrap", pc: 0xfffe, program: []byte{0xd0, 0xfe}, wantCycles: 4},
		{name: "65C02 BRA", variant: Variant65C02, pc: 0x80fe, program: []byte{0x80, 0xfe}, wantCycles: 4},
		{name: "65C02 BBR", variant: Variant65C02, pc: 0x80fd, program: []byte{0x0f, 0x10, 0xfd}, wantCycles: 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cpu := newProgramCPU(t, tt.variant, tt.pc, tt.program...)
			assert.NoError(t, cpu.Step())

			// The page comparison uses the following instruction's address,
			// even when the taken branch returns to its own opcode.
			assert.Equal(t, tt.pc, cpu.PC)
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
