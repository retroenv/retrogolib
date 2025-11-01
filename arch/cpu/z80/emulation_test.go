package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type cpuTest struct {
	Name  string
	Setup func(cpu *CPU)
	Check func(cpu *CPU)
}

func cpuTestSetup(t *testing.T) *CPU {
	t.Helper()
	memory := NewBasicMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)
	return cpu
}

func runCPUTest(t *testing.T, tests []cpuTest) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			cpu := cpuTestSetup(t)
			test.Setup(cpu)
			test.Check(cpu)
		})
	}
}

func TestInc8(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "normal increment",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
			},
			Check: func(cpu *CPU) {
				result := cpu.inc8(cpu.A)
				assert.Equal(t, uint8(0x11), result)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
				assert.Equal(t, uint8(0), cpu.Flags.S)
				assert.Equal(t, uint8(0), cpu.Flags.H)
				assert.Equal(t, uint8(0), cpu.Flags.N)
			},
		},
		{
			Name: "zero result",
			Setup: func(cpu *CPU) {
				cpu.A = 0xFF
			},
			Check: func(cpu *CPU) {
				result := cpu.inc8(cpu.A)
				assert.Equal(t, uint8(0x00), result)
				assert.Equal(t, uint8(1), cpu.Flags.Z)
				assert.Equal(t, uint8(1), cpu.Flags.H)
			},
		},
		{
			Name: "overflow",
			Setup: func(cpu *CPU) {
				cpu.A = 0x7F
			},
			Check: func(cpu *CPU) {
				result := cpu.inc8(cpu.A)
				assert.Equal(t, uint8(0x80), result)
				assert.Equal(t, uint8(1), cpu.Flags.S)
				assert.Equal(t, uint8(1), cpu.Flags.P)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestDec8(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "normal decrement",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
			},
			Check: func(cpu *CPU) {
				result := cpu.dec8(cpu.A)
				assert.Equal(t, uint8(0x0F), result)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
				assert.Equal(t, uint8(0), cpu.Flags.S)
				assert.Equal(t, uint8(1), cpu.Flags.N)
			},
		},
		{
			Name: "zero result",
			Setup: func(cpu *CPU) {
				cpu.A = 0x01
			},
			Check: func(cpu *CPU) {
				result := cpu.dec8(cpu.A)
				assert.Equal(t, uint8(0x00), result)
				assert.Equal(t, uint8(1), cpu.Flags.Z)
			},
		},
		{
			Name: "underflow",
			Setup: func(cpu *CPU) {
				cpu.A = 0x80
			},
			Check: func(cpu *CPU) {
				result := cpu.dec8(cpu.A)
				assert.Equal(t, uint8(0x7F), result)
				assert.Equal(t, uint8(0), cpu.Flags.S)
				assert.Equal(t, uint8(1), cpu.Flags.P)
			},
		},
		{
			Name: "half carry",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
			},
			Check: func(cpu *CPU) {
				cpu.dec8(cpu.A)
				assert.Equal(t, uint8(1), cpu.Flags.H)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestAdd8(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "normal addition",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
			},
			Check: func(cpu *CPU) {
				result := cpu.add8(cpu.A, 0x20)
				assert.Equal(t, uint8(0x30), result)
				assert.Equal(t, uint8(0), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
				assert.Equal(t, uint8(0), cpu.Flags.N)
			},
		},
		{
			Name: "carry addition",
			Setup: func(cpu *CPU) {
				cpu.A = 0xFF
			},
			Check: func(cpu *CPU) {
				result := cpu.add8(cpu.A, 0x02)
				assert.Equal(t, uint8(0x01), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
			},
		},
		{
			Name: "zero result",
			Setup: func(cpu *CPU) {
				cpu.A = 0x80
			},
			Check: func(cpu *CPU) {
				result := cpu.add8(cpu.A, 0x80)
				assert.Equal(t, uint8(0x00), result)
				assert.Equal(t, uint8(1), cpu.Flags.Z)
				assert.Equal(t, uint8(1), cpu.Flags.C)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestSub8(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "normal subtraction",
			Setup: func(cpu *CPU) {
				cpu.A = 0x30
			},
			Check: func(cpu *CPU) {
				result := cpu.sub8(cpu.A, 0x10)
				assert.Equal(t, uint8(0x20), result)
				assert.Equal(t, uint8(0), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
				assert.Equal(t, uint8(1), cpu.Flags.N)
			},
		},
		{
			Name: "zero result",
			Setup: func(cpu *CPU) {
				cpu.A = 0x20
			},
			Check: func(cpu *CPU) {
				result := cpu.sub8(cpu.A, 0x20)
				assert.Equal(t, uint8(0x00), result)
				assert.Equal(t, uint8(1), cpu.Flags.Z)
				assert.Equal(t, uint8(1), cpu.Flags.N)
			},
		},
		{
			Name: "borrow subtraction",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
			},
			Check: func(cpu *CPU) {
				result := cpu.sub8(cpu.A, 0x20)
				assert.Equal(t, uint8(0xF0), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
				assert.Equal(t, uint8(1), cpu.Flags.S)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestArithmeticWithCarry(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "adc with carry clear",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
				cpu.Flags.C = 0
			},
			Check: func(cpu *CPU) {
				result := cpu.adc(cpu.A, 0x20)
				assert.Equal(t, uint8(0x30), result)
				assert.Equal(t, uint8(0), cpu.Flags.C)
			},
		},
		{
			Name: "adc with carry set",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
				cpu.Flags.C = 1
			},
			Check: func(cpu *CPU) {
				result := cpu.adc(cpu.A, 0x20)
				assert.Equal(t, uint8(0x31), result)
				assert.Equal(t, uint8(0), cpu.Flags.C)
			},
		},
		{
			Name: "sbc with carry clear",
			Setup: func(cpu *CPU) {
				cpu.A = 0x30
				cpu.Flags.C = 0
			},
			Check: func(cpu *CPU) {
				result := cpu.sbc(cpu.A, 0x10)
				assert.Equal(t, uint8(0x20), result)
				assert.Equal(t, uint8(0), cpu.Flags.C)
			},
		},
		{
			Name: "sbc with carry set",
			Setup: func(cpu *CPU) {
				cpu.A = 0x30
				cpu.Flags.C = 1
			},
			Check: func(cpu *CPU) {
				result := cpu.sbc(cpu.A, 0x10)
				assert.Equal(t, uint8(0x1F), result)
				assert.Equal(t, uint8(0), cpu.Flags.C)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestNegation(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "negate positive",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
			},
			Check: func(cpu *CPU) {
				result := cpu.neg(cpu.A)
				assert.Equal(t, uint8(0xF0), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
				assert.Equal(t, uint8(1), cpu.Flags.S)
				assert.Equal(t, uint8(1), cpu.Flags.N)
			},
		},
		{
			Name: "negate zero",
			Setup: func(cpu *CPU) {
				cpu.A = 0x00
			},
			Check: func(cpu *CPU) {
				result := cpu.neg(cpu.A)
				assert.Equal(t, uint8(0x00), result)
				assert.Equal(t, uint8(1), cpu.Flags.Z)
				assert.Equal(t, uint8(0), cpu.Flags.C)
			},
		},
		{
			Name: "negate 0x80",
			Setup: func(cpu *CPU) {
				cpu.A = 0x80
			},
			Check: func(cpu *CPU) {
				result := cpu.neg(cpu.A)
				assert.Equal(t, uint8(0x80), result)
				assert.Equal(t, uint8(1), cpu.Flags.S)
				assert.Equal(t, uint8(1), cpu.Flags.P)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestCompare(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "equal values",
			Setup: func(cpu *CPU) {
				cpu.A = 0x20
			},
			Check: func(cpu *CPU) {
				cpu.cp(cpu.A, 0x20)
				assert.Equal(t, uint8(1), cpu.Flags.Z)
				assert.Equal(t, uint8(0), cpu.Flags.C)
				assert.Equal(t, uint8(1), cpu.Flags.N)
			},
		},
		{
			Name: "A greater than operand",
			Setup: func(cpu *CPU) {
				cpu.A = 0x30
			},
			Check: func(cpu *CPU) {
				cpu.cp(cpu.A, 0x20)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
				assert.Equal(t, uint8(0), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.S)
			},
		},
		{
			Name: "A less than operand",
			Setup: func(cpu *CPU) {
				cpu.A = 0x10
			},
			Check: func(cpu *CPU) {
				cpu.cp(cpu.A, 0x20)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
				assert.Equal(t, uint8(1), cpu.Flags.C)
				assert.Equal(t, uint8(1), cpu.Flags.S)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestLogicalOperations(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "AND operation",
			Setup: func(cpu *CPU) {
				cpu.A = 0xFF
			},
			Check: func(cpu *CPU) {
				result := cpu.and8(cpu.A, 0x0F)
				assert.Equal(t, uint8(0x0F), result)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
				assert.Equal(t, uint8(0), cpu.Flags.C)
				assert.Equal(t, uint8(1), cpu.Flags.H)
				assert.Equal(t, uint8(0), cpu.Flags.N)
			},
		},
		{
			Name: "OR operation",
			Setup: func(cpu *CPU) {
				cpu.A = 0x0F
			},
			Check: func(cpu *CPU) {
				result := cpu.or8(cpu.A, 0xF0)
				assert.Equal(t, uint8(0xFF), result)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
				assert.Equal(t, uint8(1), cpu.Flags.S)
				assert.Equal(t, uint8(0), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.H)
				assert.Equal(t, uint8(0), cpu.Flags.N)
			},
		},
		{
			Name: "XOR operation",
			Setup: func(cpu *CPU) {
				cpu.A = 0xFF
			},
			Check: func(cpu *CPU) {
				result := cpu.xor8(cpu.A, 0xFF)
				assert.Equal(t, uint8(0x00), result)
				assert.Equal(t, uint8(1), cpu.Flags.Z)
				assert.Equal(t, uint8(0), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.H)
				assert.Equal(t, uint8(0), cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestBitOperations(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Test BIT instruction - test bit 7 of value 0x80
	cpu.bit(7, 0x80)
	assert.Equal(t, uint8(0), cpu.Flags.Z, "Bit 7 of 0x80 should be set, Z flag should be clear")
	assert.Equal(t, uint8(1), cpu.Flags.H, "Half carry should be set for BIT instruction")
	assert.Equal(t, uint8(0), cpu.Flags.N, "N flag should be clear for BIT instruction")

	// Test BIT instruction - test bit 0 of value 0x80
	cpu.bit(0, 0x80)
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Bit 0 of 0x80 should be clear, Z flag should be set")

	// Test SET instruction - set bit 3 of value 0x00
	result := cpu.setBit(3, 0x00)
	assert.Equal(t, uint8(0x08), result, "Setting bit 3 of 0x00 should give 0x08")

	// Test RES instruction - reset bit 3 of value 0xFF
	result = cpu.res(3, 0xFF)
	assert.Equal(t, uint8(0xF7), result, "Resetting bit 3 of 0xFF should give 0xF7")
}

func TestRotateOperations(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "RLC operation",
			Setup: func(cpu *CPU) {
				cpu.A = 0x85
			},
			Check: func(cpu *CPU) {
				result := cpu.rlc(cpu.A)
				assert.Equal(t, uint8(0x0B), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
				assert.Equal(t, uint8(0), cpu.Flags.H)
				assert.Equal(t, uint8(0), cpu.Flags.N)
			},
		},
		{
			Name: "RRC operation",
			Setup: func(cpu *CPU) {
				cpu.A = 0x85
			},
			Check: func(cpu *CPU) {
				result := cpu.rrc(cpu.A)
				assert.Equal(t, uint8(0xC2), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
			},
		},
		{
			Name: "RL operation with carry set",
			Setup: func(cpu *CPU) {
				cpu.A = 0x85
				cpu.Flags.C = 1
			},
			Check: func(cpu *CPU) {
				result := cpu.rl(cpu.A)
				assert.Equal(t, uint8(0x0B), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
			},
		},
		{
			Name: "RR operation with carry set",
			Setup: func(cpu *CPU) {
				cpu.A = 0x85
				cpu.Flags.C = 1
			},
			Check: func(cpu *CPU) {
				result := cpu.rr(cpu.A)
				assert.Equal(t, uint8(0xC2), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestShiftOperations(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "SLA operation",
			Setup: func(cpu *CPU) {
				cpu.A = 0x85
			},
			Check: func(cpu *CPU) {
				result := cpu.sla(cpu.A)
				assert.Equal(t, uint8(0x0A), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
			},
		},
		{
			Name: "SRA operation",
			Setup: func(cpu *CPU) {
				cpu.A = 0x85
			},
			Check: func(cpu *CPU) {
				result := cpu.sra(cpu.A)
				assert.Equal(t, uint8(0xC2), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
			},
		},
		{
			Name: "SRL operation",
			Setup: func(cpu *CPU) {
				cpu.A = 0x85
			},
			Check: func(cpu *CPU) {
				result := cpu.srl(cpu.A)
				assert.Equal(t, uint8(0x42), result)
				assert.Equal(t, uint8(1), cpu.Flags.C)
				assert.Equal(t, uint8(0), cpu.Flags.Z)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestRegisterPairs(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Test 16-bit increment
	cpu.B, cpu.C = 0xFF, 0xFF
	bc := cpu.bc()
	result := cpu.inc16(bc)
	assert.Equal(t, uint16(0x0000), result, "Increment of 0xFFFF should wrap to 0x0000")

	// Test 16-bit decrement
	cpu.H, cpu.L = 0x00, 0x00
	hl := cpu.hl()
	result = cpu.dec16(hl)
	assert.Equal(t, uint16(0xFFFF), result, "Decrement of 0x0000 should wrap to 0xFFFF")

	// Test 16-bit addition
	cpu.H, cpu.L = 0x12, 0x34
	hl = cpu.hl()
	result = cpu.addHL(hl, 0x5678)
	assert.Equal(t, uint16(0x68AC), result, "0x1234 + 0x5678 should equal 0x68AC")
	assert.Equal(t, uint8(0), cpu.Flags.C, "Should not set carry for this addition")
	assert.Equal(t, uint8(0), cpu.Flags.N, "N flag should be clear for ADD HL")

	// Test 16-bit addition with carry
	cpu.H, cpu.L = 0xFF, 0xFF
	hl = cpu.hl()
	result = cpu.addHL(hl, 0x0001)
	assert.Equal(t, uint16(0x0000), result, "0xFFFF + 0x0001 should wrap to 0x0000")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Should set carry for this addition")
}

func TestStackOperations(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Test PUSH operation
	cpu.SP = 0x1000
	cpu.B, cpu.C = 0x12, 0x34
	bc := cpu.bc()
	cpu.push16(bc)
	assert.Equal(t, uint16(0x0FFE), cpu.SP, "SP should decrement by 2 after PUSH")
	assert.Equal(t, uint8(0x34), cpu.Memory().Read(0x0FFE), "Low byte should be at SP")
	assert.Equal(t, uint8(0x12), cpu.Memory().Read(0x0FFF), "High byte should be at SP+1")

	// Test POP operation
	cpu.SP = 0x0FFE
	cpu.Memory().Write(0x0FFE, 0x78) // Low byte at SP
	cpu.Memory().Write(0x0FFF, 0x56) // High byte at SP+1
	result := cpu.pop16()
	assert.Equal(t, uint16(0x5678), result, "POP should return correct 16-bit value")
	assert.Equal(t, uint16(0x1000), cpu.SP, "SP should increment by 2 after POP")
}

func TestExchange(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Test EX DE,HL
	cpu.D, cpu.E = 0x12, 0x34
	cpu.H, cpu.L = 0x56, 0x78
	cpu.exDEHL()
	assert.Equal(t, uint8(0x56), cpu.D, "D should have H's value")
	assert.Equal(t, uint8(0x78), cpu.E, "E should have L's value")
	assert.Equal(t, uint8(0x12), cpu.H, "H should have D's value")
	assert.Equal(t, uint8(0x34), cpu.L, "L should have E's value")

	// Test EXX (exchange BC, DE, HL with shadow registers)
	cpu.B, cpu.C = 0x11, 0x22
	cpu.D, cpu.E = 0x33, 0x44
	cpu.H, cpu.L = 0x55, 0x66
	cpu.AltB, cpu.AltC = 0xAA, 0xBB
	cpu.AltD, cpu.AltE = 0xCC, 0xDD
	cpu.AltH, cpu.AltL = 0xEE, 0xFF

	cpu.exx()

	assert.Equal(t, uint8(0xAA), cpu.B, "B should have alt B's value")
	assert.Equal(t, uint8(0xBB), cpu.C, "C should have alt C's value")
	assert.Equal(t, uint8(0xCC), cpu.D, "D should have alt D's value")
	assert.Equal(t, uint8(0xDD), cpu.E, "E should have alt E's value")
	assert.Equal(t, uint8(0xEE), cpu.H, "H should have alt H's value")
	assert.Equal(t, uint8(0xFF), cpu.L, "L should have alt L's value")
	assert.Equal(t, uint8(0x11), cpu.AltB, "Alt B should have B's value")
	assert.Equal(t, uint8(0x22), cpu.AltC, "Alt C should have C's value")
}

func TestExchangeAF(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Test EX AF,AF'
	cpu.A = 0x12
	cpu.setFlags(0x34)
	cpu.AltA = 0x56
	cpu.setFlagsFromUint8(&cpu.AltFlags, 0x78)

	cpu.exAF()

	assert.Equal(t, uint8(0x56), cpu.A, "A should have alt A's value")
	assert.Equal(t, uint8(0x78), cpu.getFlags(), "F should have alt F's value")
	assert.Equal(t, uint8(0x12), cpu.AltA, "Alt A should have A's value")
	assert.Equal(t, uint8(0x34), cpu.getFlagsAsUint8(cpu.AltFlags), "Alt F should have F's value")
}

func TestLoadInstructionExecution(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Test LD r,r' (register to register)
	cpu.B = 0x42
	cpu.ldReg8(&cpu.A, cpu.B)
	assert.Equal(t, uint8(0x42), cpu.A, "A should have B's value after LD A,B")

	// Test LD r,n (immediate to register)
	cpu.ldImm8(&cpu.C, 0x84)
	assert.Equal(t, uint8(0x84), cpu.C, "C should have immediate value after LD C,n")

	// Test LD r,(HL) (memory to register)
	cpu.H, cpu.L = 0x20, 0x00
	cpu.Memory().Write(0x2000, 0x99)
	hl := cpu.hl()
	cpu.ldMemToReg8(&cpu.D, hl)
	assert.Equal(t, uint8(0x99), cpu.D, "D should have memory value after LD D,(HL)")

	// Test LD (HL),r (register to memory)
	cpu.E = 0x77
	hl = cpu.hl()
	cpu.ldRegToMem8(hl, cpu.E)
	assert.Equal(t, uint8(0x77), cpu.Memory().Read(0x2000), "Memory should have E's value after LD (HL),E")
}

func TestMemoryOperations(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Test LDI (Load and increment)
	cpu.H, cpu.L = 0x11, 0x11
	cpu.D, cpu.E = 0x22, 0x22
	cpu.B, cpu.C = 0x00, 0x03
	cpu.Memory().Write(0x1111, 0x88)

	cpu.ldi()

	assert.Equal(t, uint8(0x88), cpu.Memory().Read(0x2222), "LDI should copy byte from (HL) to (DE)")
	assert.Equal(t, uint16(0x1112), cpu.hl(), "HL should increment after LDI")
	assert.Equal(t, uint16(0x2223), cpu.de(), "DE should increment after LDI")
	assert.Equal(t, uint16(0x0002), cpu.bc(), "BC should decrement after LDI")
	assert.Equal(t, uint8(1), cpu.Flags.P, "P/V flag should indicate BC != 0")

	// Test LDD (Load and decrement)
	cpu.H, cpu.L = 0x11, 0x11
	cpu.D, cpu.E = 0x22, 0x22
	cpu.B, cpu.C = 0x00, 0x01
	cpu.Memory().Write(0x1111, 0x66)

	cpu.ldd()

	assert.Equal(t, uint8(0x66), cpu.Memory().Read(0x2222), "LDD should copy byte from (HL) to (DE)")
	assert.Equal(t, uint16(0x1110), cpu.hl(), "HL should decrement after LDD")
	assert.Equal(t, uint16(0x2221), cpu.de(), "DE should decrement after LDD")
	assert.Equal(t, uint16(0x0000), cpu.bc(), "BC should decrement after LDD")
	assert.Equal(t, uint8(0), cpu.Flags.P, "P/V flag should indicate BC == 0")
}

func TestJumpInstructions(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "JP unconditional",
			Setup: func(cpu *CPU) {
				cpu.PC = 0x1000
			},
			Check: func(cpu *CPU) {
				cpu.jp(0x2000)
				assert.Equal(t, uint16(0x2000), cpu.PC)
			},
		},
		{
			Name: "JR positive offset",
			Setup: func(cpu *CPU) {
				cpu.PC = 0x1000
			},
			Check: func(cpu *CPU) {
				cpu.jr(10)
				assert.Equal(t, uint16(0x100A), cpu.PC)
			},
		},
		{
			Name: "JR negative offset",
			Setup: func(cpu *CPU) {
				cpu.PC = 0x1000
			},
			Check: func(cpu *CPU) {
				cpu.jr(-10) // 0xF6 in two's complement
				assert.Equal(t, uint16(0x0FF6), cpu.PC)
			},
		},
		{
			Name: "JP conditional taken",
			Setup: func(cpu *CPU) {
				cpu.PC = 0x1000
				cpu.Flags.Z = 1
			},
			Check: func(cpu *CPU) {
				cpu.jpZ(0x2000)
				assert.Equal(t, uint16(0x2000), cpu.PC)
			},
		},
		{
			Name: "JP conditional not taken",
			Setup: func(cpu *CPU) {
				cpu.PC = 0x1000
				cpu.Flags.Z = 0
			},
			Check: func(cpu *CPU) {
				cpu.jpZ(0x2000)
				assert.Equal(t, uint16(0x1000), cpu.PC)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestDJNZ(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "DJNZ jump taken",
			Setup: func(cpu *CPU) {
				cpu.PC = 0x1000
				cpu.B = 0x02
			},
			Check: func(cpu *CPU) {
				cpu.djnz(5)
				assert.Equal(t, uint8(0x01), cpu.B)
				assert.Equal(t, uint16(0x1005), cpu.PC)
			},
		},
		{
			Name: "DJNZ jump not taken",
			Setup: func(cpu *CPU) {
				cpu.PC = 0x1000
				cpu.B = 0x01
			},
			Check: func(cpu *CPU) {
				cpu.djnz(5)
				assert.Equal(t, uint8(0x00), cpu.B)
				assert.Equal(t, uint16(0x1000), cpu.PC)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestEndlessLoopDetection(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.PC = 0x0100

	// Create an endless loop: JR -2 (0x18 0xFE)
	cpu.Memory().Write(0x0100, 0x18) // JR relative
	cpu.Memory().Write(0x0101, 0xFE) // -2 (jumps to 0x0100 + 2 + (-2) = 0x0100)

	initialCycles := cpu.cycles
	initialPC := cpu.PC

	// Execute the jump instruction once
	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error for JR")

	// Verify that PC jumped back to the same address (endless loop)
	assert.Equal(t, initialPC, cpu.PC, "PC should jump back to itself (endless loop)")
	assert.Greater(t, cpu.cycles, initialCycles, "Cycles should have advanced")

	// Execute the same instruction again to confirm it's truly an endless loop
	secondCycles := cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error for second iteration")
	assert.Equal(t, initialPC, cpu.PC, "PC should still be at the same address")
	assert.Greater(t, cpu.cycles, secondCycles, "Cycles should have advanced again")

	// Test another endless loop pattern: JP 0x0200 to itself
	cpu.PC = 0x0200
	cpu.Memory().Write(0x0200, 0xC3) // JP absolute
	cpu.Memory().Write(0x0201, 0x00) // Low byte of address (0x0200)
	cpu.Memory().Write(0x0202, 0x02) // High byte of address

	thirdCycles := cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error for JP")
	assert.Equal(t, uint16(0x0200), cpu.PC, "PC should jump to same address (endless loop)")
	assert.Greater(t, cpu.cycles, thirdCycles, "Cycles should have advanced for JP")
}
