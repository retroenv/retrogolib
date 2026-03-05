package m6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func cpuTestSetup65C02(t *testing.T) *CPU {
	t.Helper()
	memory, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	memory.WriteWord(ResetAddress, 0x8000)
	memory.WriteWord(IrqAddress, testIrqAddress)
	cpu := New(memory, WithVariant(Variant65C02))
	return cpu
}

// TestBra tests the Branch Always instruction.
func TestBra(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup65C02(t)

	assert.NoError(t, bra(cpu, Absolute(0x1234)))
	assert.Equal(t, 0x1234, cpu.PC)
}

// TestPhx tests the Push X Register instruction.
func TestPhx(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup65C02(t)

	cpu.X = 0x42
	assert.NoError(t, phx(cpu))

	b := cpu.memory.Read(StackBase + InitialStack)
	assert.Equal(t, 0x42, b)
	assert.Equal(t, InitialStack-1, cpu.SP)
}

// TestPhy tests the Push Y Register instruction.
func TestPhy(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup65C02(t)

	cpu.Y = 0x42
	assert.NoError(t, phy(cpu))

	b := cpu.memory.Read(StackBase + InitialStack)
	assert.Equal(t, 0x42, b)
	assert.Equal(t, InitialStack-1, cpu.SP)
}

// TestPlx tests the Pull X Register instruction.
func TestPlx(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "basic pull",
			Setup: func(cpu *CPU) {
				cpu.SP = 1
				cpu.memory.Write(StackBase+2, 0x42)
				assert.NoError(t, plx(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x42, cpu.X)
				assert.Equal(t, 2, cpu.SP)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "pull zero sets Z flag",
			Setup: func(cpu *CPU) {
				cpu.SP = 1
				cpu.memory.Write(StackBase+2, 0x00)
				assert.NoError(t, plx(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0, cpu.X)
				assert.Equal(t, 1, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "pull negative sets N flag",
			Setup: func(cpu *CPU) {
				cpu.SP = 1
				cpu.memory.Write(StackBase+2, 0x80)
				assert.NoError(t, plx(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x80, cpu.X)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 1, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

// TestPly tests the Pull Y Register instruction.
func TestPly(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "basic pull",
			Setup: func(cpu *CPU) {
				cpu.SP = 1
				cpu.memory.Write(StackBase+2, 0x42)
				assert.NoError(t, ply(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x42, cpu.Y)
				assert.Equal(t, 2, cpu.SP)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "pull zero sets Z flag",
			Setup: func(cpu *CPU) {
				cpu.SP = 1
				cpu.memory.Write(StackBase+2, 0x00)
				assert.NoError(t, ply(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0, cpu.Y)
				assert.Equal(t, 1, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

// TestStz tests the Store Zero instruction.
func TestStz(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "zero page",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x10, 0xFF)
				assert.NoError(t, stz(cpu, Absolute(0x10)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0, cpu.memory.Read(0x10))
			},
		},
		{
			Name: "absolute",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x1234, 0xFF)
				assert.NoError(t, stz(cpu, Absolute(0x1234)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0, cpu.memory.Read(0x1234))
			},
		},
		{
			Name: "zero page x",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x15, 0xFF)
				cpu.X = 5
				assert.NoError(t, stz(cpu, ZeroPage(0x10), &cpu.X))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0, cpu.memory.Read(0x15))
			},
		},
		{
			Name: "absolute x",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x1239, 0xFF)
				cpu.X = 5
				assert.NoError(t, stz(cpu, Absolute(0x1234), &cpu.X))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0, cpu.memory.Read(0x1239))
			},
		},
	}
	runCPUTest(t, tests)
}

// TestTrb tests the Test and Reset Bits instruction.
func TestTrb(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "basic reset bits",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x10, 0xFF)
				cpu.A = 0x0F
				assert.NoError(t, trb(cpu, Absolute(0x10)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xF0, cpu.memory.Read(0x10))
				assert.Equal(t, 0, cpu.Flags.Z) // A & original != 0
			},
		},
		{
			Name: "Z flag set when AND is zero",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x10, 0xF0)
				cpu.A = 0x0F
				assert.NoError(t, trb(cpu, Absolute(0x10)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xF0, cpu.memory.Read(0x10))
				assert.Equal(t, 1, cpu.Flags.Z) // A & original == 0
			},
		},
		{
			Name: "all bits reset",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x10, 0xFF)
				cpu.A = 0xFF
				assert.NoError(t, trb(cpu, Absolute(0x10)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x00, cpu.memory.Read(0x10))
				assert.Equal(t, 0, cpu.Flags.Z) // A & original != 0
			},
		},
	}
	runCPUTest(t, tests)
}

// TestTsb tests the Test and Set Bits instruction.
func TestTsb(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "basic set bits",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x10, 0xF0)
				cpu.A = 0x0F
				assert.NoError(t, tsb(cpu, Absolute(0x10)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xFF, cpu.memory.Read(0x10))
				assert.Equal(t, 1, cpu.Flags.Z) // A & original == 0
			},
		},
		{
			Name: "Z flag clear when AND is non-zero",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x10, 0xFF)
				cpu.A = 0x0F
				assert.NoError(t, tsb(cpu, Absolute(0x10)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xFF, cpu.memory.Read(0x10))
				assert.Equal(t, 0, cpu.Flags.Z) // A & original != 0
			},
		},
	}
	runCPUTest(t, tests)
}

// TestIncAccumulator tests the INC A instruction (65C02).
func TestIncAccumulator(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "basic increment",
			Setup: func(cpu *CPU) {
				cpu.A = 0x41
				assert.NoError(t, inc65c02(cpu, Accumulator(0)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x42, cpu.A)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "increment wraps to zero",
			Setup: func(cpu *CPU) {
				cpu.A = 0xFF
				assert.NoError(t, inc65c02(cpu, Accumulator(0)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x00, cpu.A)
				assert.Equal(t, 1, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "increment sets negative",
			Setup: func(cpu *CPU) {
				cpu.A = 0x7F
				assert.NoError(t, inc65c02(cpu, Accumulator(0)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x80, cpu.A)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 1, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

// TestDecAccumulator tests the DEC A instruction (65C02).
func TestDecAccumulator(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "basic decrement",
			Setup: func(cpu *CPU) {
				cpu.A = 0x42
				assert.NoError(t, dec65c02(cpu, Accumulator(0)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x41, cpu.A)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "decrement to zero",
			Setup: func(cpu *CPU) {
				cpu.A = 0x01
				assert.NoError(t, dec65c02(cpu, Accumulator(0)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x00, cpu.A)
				assert.Equal(t, 1, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "decrement wraps to 0xFF",
			Setup: func(cpu *CPU) {
				cpu.A = 0x00
				assert.NoError(t, dec65c02(cpu, Accumulator(0)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xFF, cpu.A)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 1, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

// TestBit65C02Immediate tests BIT with immediate addressing on 65C02.
// In immediate mode, only Z flag is affected (V and N are not modified).
func TestBit65C02Immediate(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "immediate zero result",
			Setup: func(cpu *CPU) {
				cpu.A = 0x0F
				cpu.Flags.V = 1 // Should remain unchanged
				cpu.Flags.N = 1 // Should remain unchanged
				assert.NoError(t, bit65c02(cpu, int(0xF0)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.Flags.Z) // A & imm == 0
				assert.Equal(t, 1, cpu.Flags.V) // Unchanged
				assert.Equal(t, 1, cpu.Flags.N) // Unchanged
			},
		},
		{
			Name: "immediate non-zero result",
			Setup: func(cpu *CPU) {
				cpu.A = 0xFF
				cpu.Flags.V = 0 // Should remain unchanged
				cpu.Flags.N = 0 // Should remain unchanged
				assert.NoError(t, bit65c02(cpu, int(0x0F)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0, cpu.Flags.Z) // A & imm != 0
				assert.Equal(t, 0, cpu.Flags.V) // Unchanged
				assert.Equal(t, 0, cpu.Flags.N) // Unchanged
			},
		},
	}
	runCPUTest(t, tests)
}

// TestJmp65C02IndirectPageBug tests that JMP (abs) works correctly on 65C02
// (no page boundary bug like NMOS).
func TestJmp65C02IndirectPageBug(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup65C02(t)

	// On NMOS, JMP ($02FF) would read low byte from $02FF and high byte from $0200.
	// On 65C02, it correctly reads from $02FF and $0300.
	cpu.memory.Write(0x02FF, 0x34)
	cpu.memory.Write(0x0300, 0x12) // 65C02 reads this correctly
	cpu.memory.Write(0x0200, 0x56) // NMOS would read this instead

	assert.NoError(t, jmp65c02(cpu, Indirect(0x02FF)))
	assert.Equal(t, 0x1234, cpu.PC) // Should be 0x1234, not 0x5634
}

// TestJmp65C02AbsoluteXIndirect tests JMP (abs,X) on 65C02.
func TestJmp65C02AbsoluteXIndirect(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup65C02(t)

	// Set up a jump table: abs = 0x1000, X = 4, so reads from 0x1004
	cpu.memory.Write(0x1004, 0x00)
	cpu.memory.Write(0x1005, 0x20) // Target = 0x2000

	assert.NoError(t, jmp65c02(cpu, AbsoluteXIndirect(0x1000), Absolute(0x2000)))
	assert.Equal(t, 0x2000, cpu.PC)
}

// TestLdaZeroPageIndirect tests LDA with zero page indirect addressing (65C02).
func TestLdaZeroPageIndirect(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup65C02(t)

	// Set up zero page pointer at $10 pointing to $1234
	cpu.memory.Write(0x10, 0x34) // Low byte
	cpu.memory.Write(0x11, 0x12) // High byte
	// Store value at $1234
	cpu.memory.Write(0x1234, 0x42)

	assert.NoError(t, lda(cpu, ZeroPageIndirect(0x10), IndirectResolved(0x1234)))
	assert.Equal(t, 0x42, cpu.A)
}

// TestStaZeroPageIndirect tests STA with zero page indirect addressing (65C02).
func TestStaZeroPageIndirect(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup65C02(t)

	// Set up zero page pointer at $10 pointing to $1234
	cpu.memory.Write(0x10, 0x34)
	cpu.memory.Write(0x11, 0x12)

	cpu.A = 0x42
	assert.NoError(t, sta(cpu, ZeroPageIndirect(0x10), IndirectResolved(0x1234)))
	assert.Equal(t, 0x42, cpu.memory.Read(0x1234))
}

// TestBrkClearsDecimalFlag65C02 tests that BRK clears the D flag on 65C02.
func TestBrkClearsDecimalFlag65C02(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup65C02(t)

	cpu.PC = 0x1000
	cpu.SP = 0xFD
	cpu.Flags.D = 1

	assert.NoError(t, brk(cpu))

	assert.Equal(t, 0, cpu.Flags.D) // D flag cleared on 65C02
	assert.Equal(t, 1, cpu.Flags.I) // I flag set
}

// TestInterruptClearsDecimalFlag65C02 tests that interrupts clear D flag on 65C02.
func TestInterruptClearsDecimalFlag65C02(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup65C02(t)

	cpu.PC = 0x1000
	cpu.SP = 0xFD
	cpu.Flags.D = 1

	cpu.executeInterrupt(testIrqAddress)

	assert.Equal(t, 0, cpu.Flags.D) // D flag cleared on 65C02
	assert.Equal(t, testIrqAddress, cpu.PC)
}

// TestBrkPreservesDecimalFlagNMOS tests that BRK does NOT clear D flag on NMOS.
func TestBrkPreservesDecimalFlagNMOS(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t) // NMOS variant

	cpu.PC = 0x1000
	cpu.SP = 0xFD
	cpu.Flags.D = 1

	assert.NoError(t, brk(cpu))

	assert.Equal(t, 1, cpu.Flags.D) // D flag preserved on NMOS
}

// Test65C02OpcodeTableCompleteness ensures all 256 entries have instructions.
func Test65C02OpcodeTableCompleteness(t *testing.T) {
	t.Parallel()

	for i, opcode := range Opcodes65C02 {
		assert.NotNil(t, opcode.Instruction,
			"Opcode 0x%02X has nil instruction in 65C02 table", i)
		assert.True(t, opcode.Timing > 0,
			"Opcode 0x%02X has zero timing in 65C02 table", i)
	}
}

// Test65C02OpcodeTableNoUnofficialInstructions ensures the 65C02 table
// has no unofficial NMOS instructions (they should all be NOPs).
func Test65C02OpcodeTableNoUnofficialInstructions(t *testing.T) {
	t.Parallel()

	// These are NMOS unofficial instructions that should NOT appear in 65C02 table
	unofficialNames := map[string]bool{
		AlrName: true,
		AncName: true,
		AneName: true,
		ArrName: true,
		AxsName: true,
		DcpName: true,
		IscName: true,
		LasName: true,
		LaxName: true,
		LxaName: true,
		RlaName: true,
		RraName: true,
		SaxName: true,
		ShaName: true,
		ShxName: true,
		ShyName: true,
		SloName: true,
		SreName: true,
		TasName: true,
	}

	for i, opcode := range Opcodes65C02 {
		if opcode.Instruction == nil {
			continue
		}
		_, isUnofficial := unofficialNames[opcode.Instruction.Name]
		assert.False(t, isUnofficial,
			"Opcode 0x%02X uses unofficial instruction %s in 65C02 table", i, opcode.Instruction.Name)
	}
}

// Test65C02VariantSelection ensures the correct opcode table is selected.
func Test65C02VariantSelection(t *testing.T) {
	t.Parallel()

	// NMOS CPU
	nmosMemory, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	nmosMemory.WriteWord(ResetAddress, 0x8000)
	nmosMemory.WriteWord(IrqAddress, testIrqAddress)
	nmosCPU := New(nmosMemory)

	// 65C02 CPU
	cmosMemory, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	cmosMemory.WriteWord(ResetAddress, 0x8000)
	cmosMemory.WriteWord(IrqAddress, testIrqAddress)
	cmosCPU := New(cmosMemory, WithVariant(Variant65C02))

	// Write opcode 0x1A (NOP on NMOS, INC A on 65C02) at PC
	nmosMemory.Write(0x8000, 0x1A)
	cmosMemory.Write(0x8000, 0x1A)

	// NMOS: 0x1A is NOP (unofficial)
	nmosCPU.A = 0x42
	assert.NoError(t, nmosCPU.Step())
	assert.Equal(t, 0x42, nmosCPU.A) // A unchanged (NOP)

	// 65C02: 0x1A is INC A
	cmosCPU.A = 0x42
	assert.NoError(t, cmosCPU.Step())
	assert.Equal(t, 0x43, cmosCPU.A) // A incremented
}

// Test65C02NewInstructionOpcodes verifies key 65C02 opcodes are at correct positions.
func Test65C02NewInstructionOpcodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		opcode byte
		name   string
	}{
		{0x04, TsbName},
		{0x0c, TsbName},
		{0x12, OraName},
		{0x14, TrbName},
		{0x1a, IncName},
		{0x1c, TrbName},
		{0x32, AndName},
		{0x34, BitName},
		{0x3a, DecName},
		{0x3c, BitName},
		{0x52, EorName},
		{0x5a, PhyName},
		{0x64, StzName},
		{0x72, AdcName},
		{0x74, StzName},
		{0x7a, PlyName},
		{0x7c, JmpName},
		{0x80, BraName},
		{0x89, BitName},
		{0x92, StaName},
		{0x9c, StzName},
		{0x9e, StzName},
		{0xb2, LdaName},
		{0xd2, CmpName},
		{0xda, PhxName},
		{0xf2, SbcName},
		{0xfa, PlxName},
	}

	for _, tt := range tests {
		op := Opcodes65C02[tt.opcode]
		assert.NotNil(t, op.Instruction,
			"Opcode 0x%02X should have instruction", tt.opcode)
		assert.Equal(t, tt.name, op.Instruction.Name,
			"Opcode 0x%02X should be %s, got %s", tt.opcode, tt.name, op.Instruction.Name)
	}
}

// Test65C02VerifyOpcodes ensures bidirectional opcode mapping for 65C02 table.
func Test65C02VerifyOpcodes(t *testing.T) {
	t.Parallel()

	for b, op := range Opcodes65C02 {
		ins := op.Instruction
		if ins == nil {
			continue
		}
		// Skip NOP variants used for replaced unofficial opcodes
		if ins == Nop65C02 {
			continue
		}
		// Skip the official NOP at 0xEA
		if ins == Nop {
			continue
		}

		info, ok := ins.Addressing[op.Addressing]
		if !ok {
			continue // Some instructions share handlers but have different addressing maps
		}
		assert.Equal(t, b, info.Opcode,
			"Opcode mismatch for instruction %s with addressing %d at position 0x%02X", ins.Name, op.Addressing, b)
	}
}
