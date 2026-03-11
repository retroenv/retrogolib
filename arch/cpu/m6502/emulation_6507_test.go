package m6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func cpuTestSetup6507(t *testing.T) *CPU {
	t.Helper()
	memory, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	memory.WriteWord(ResetAddress, 0x1000) // Atari 2600 ROM starts at $1000
	cpu := New(memory, WithVariant(Variant6507))
	return cpu
}

func TestVariant6507LoadStore(t *testing.T) {
	t.Parallel()

	tests := []cpuTest{
		{
			Name: "LDA immediate",
			Setup: func(cpu *CPU) {
				assert.NoError(t, lda(cpu, 0x42))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x42, cpu.A)
			},
		},
		{
			Name: "LDX immediate",
			Setup: func(cpu *CPU) {
				assert.NoError(t, ldx(cpu, 0x10))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x10, cpu.X)
			},
		},
		{
			Name: "LDY immediate",
			Setup: func(cpu *CPU) {
				assert.NoError(t, ldy(cpu, 0x20))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x20, cpu.Y)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			cpu := cpuTestSetup6507(t)
			test.Setup(cpu)
			test.Check(cpu)
		})
	}
}

func TestVariant6507Arithmetic(t *testing.T) {
	t.Parallel()

	tests := []cpuTest{
		{
			Name: "ADC",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
				assert.NoError(t, adc(cpu, 0x20))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x30, cpu.A)
				assert.Equal(t, 0, cpu.Flags.C)
			},
		},
		{
			Name: "SBC",
			Setup: func(cpu *CPU) {
				cpu.A = 0x50
				cpu.Flags.C = 1
				assert.NoError(t, sbc(cpu, 0x10))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x40, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
			},
		},
		{
			Name: "INX/DEX",
			Setup: func(cpu *CPU) {
				cpu.X = 0x05
				assert.NoError(t, inx(cpu))
				assert.NoError(t, inx(cpu))
				assert.NoError(t, dex(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x06, cpu.X)
			},
		},
		{
			Name: "INY/DEY",
			Setup: func(cpu *CPU) {
				cpu.Y = 0x05
				assert.NoError(t, iny(cpu))
				assert.NoError(t, dey(cpu))
				assert.NoError(t, dey(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x04, cpu.Y)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			cpu := cpuTestSetup6507(t)
			test.Setup(cpu)
			test.Check(cpu)
		})
	}
}

func TestVariant6507BitwiseAndStack(t *testing.T) {
	t.Parallel()

	tests := []cpuTest{
		{
			Name: "TAX transfer",
			Setup: func(cpu *CPU) {
				cpu.A = 0xAB
				assert.NoError(t, tax(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xAB, cpu.X)
				assert.Equal(t, 0xAB, cpu.A)
			},
		},
		{
			Name: "PHA/PLA stack",
			Setup: func(cpu *CPU) {
				cpu.A = 0x42
				assert.NoError(t, pha(cpu))
				cpu.A = 0x00
				assert.NoError(t, pla(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x42, cpu.A)
			},
		},
		{
			Name: "AND/ORA/EOR bitwise",
			Setup: func(cpu *CPU) {
				cpu.A = 0xFF
				assert.NoError(t, and(cpu, 0x0F))
				assert.NoError(t, ora(cpu, 0xF0))
				assert.NoError(t, eor(cpu, 0x55))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xAA, cpu.A)
			},
		},
		{
			Name: "SEC flag control",
			Setup: func(cpu *CPU) {
				assert.NoError(t, sec(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.Flags.C)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			cpu := cpuTestSetup6507(t)
			test.Setup(cpu)
			test.Check(cpu)
		})
	}
}

func TestVariant6507IrqIsNoOp(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup6507(t)

	// Store initial state.
	initialPC := cpu.PC

	// TriggerIrq should be a no-op.
	cpu.TriggerIrq()
	assert.Equal(t, false, cpu.triggerIrq)

	// CheckInterrupts should not fire.
	fired := cpu.CheckInterrupts()
	assert.Equal(t, false, fired)
	assert.Equal(t, initialPC, cpu.PC)
}

func TestVariant6507NmiIsNoOp(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup6507(t)

	initialPC := cpu.PC

	// TriggerNMI should be a no-op.
	cpu.TriggerNMI()
	assert.Equal(t, false, cpu.triggerNmi)

	fired := cpu.CheckInterrupts()
	assert.Equal(t, false, fired)
	assert.Equal(t, initialPC, cpu.PC)
}

func TestVariant6507UsesNMOSOpcodeTable(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup6507(t)

	// The 6507 should use the NMOS 6502 opcode table, not the 65C02 table.
	// Verify by checking that the variant comparison works correctly.
	assert.True(t, cpu.opts.variant < Variant65C02,
		"6507 variant must be below 65C02 for opcode table selection")
	assert.Equal(t, Variant6507, cpu.opts.variant)
}

func TestVariant6507InterruptVsNMOS(t *testing.T) {
	t.Parallel()

	// NMOS 6502 should accept interrupts.
	nmosMem, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	nmosMem.WriteWord(ResetAddress, 0x8000)
	nmosMem.WriteWord(IrqAddress, testIrqAddress)
	nmos := New(nmosMem)

	nmos.TriggerIrq()
	assert.Equal(t, true, nmos.triggerIrq)

	nmos.TriggerNMI()
	assert.Equal(t, true, nmos.triggerNmi)

	// 6507 should reject interrupts.
	cpu6507 := cpuTestSetup6507(t)

	cpu6507.TriggerIrq()
	assert.Equal(t, false, cpu6507.triggerIrq)

	cpu6507.TriggerNMI()
	assert.Equal(t, false, cpu6507.triggerNmi)
}
