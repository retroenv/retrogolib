package m6502

import (
	"testing"

	"github.com/retroenv/retrogolib/arch/nes"
	"github.com/retroenv/retrogolib/assert"
)

type cpuTest struct {
	Name  string
	Setup func(cpu *CPU)
	Check func(cpu *CPU)
}

const testIrqAddress = 0x9000

func cpuTestSetup(t *testing.T) *CPU {
	t.Helper()
	memory, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	memory.WriteWord(ResetAddress, nes.CodeBaseAddress)
	memory.WriteWord(IrqAddress, testIrqAddress)
	cpu := New(memory)
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

func TestAdc(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "result 0x00",
			Setup: func(cpu *CPU) {
				cpu.A = 2
				assert.NoError(t, adc(cpu, 0xff))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
			},
		},
		{
			Name: "result 0x01",
			Setup: func(cpu *CPU) {
				assert.NoError(t, adc(cpu, 1))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.A)
				assert.Equal(t, 0, cpu.Flags.C)
			},
		},
		{
			Name: "result 0x102",
			Setup: func(cpu *CPU) {
				cpu.A = 2
				cpu.Flags.C = 1
				assert.NoError(t, adc(cpu, 0xff))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 2, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestAnd(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.A = 0x12
	assert.NoError(t, and(cpu, 2))

	assert.Equal(t, 2, cpu.A)
}

func TestAsl(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.A = 0b0000_0001
	assert.NoError(t, asl(cpu))
	assert.Equal(t, 0b0000_0010, cpu.A)
	assert.Equal(t, 0, cpu.Flags.C)

	cpu.A = 0b1111_1110
	assert.NoError(t, asl(cpu))
	assert.Equal(t, 0b1111_1100, cpu.A)
	assert.Equal(t, 1, cpu.Flags.C)

	cpu.memory.Write(1, 0b0000_0010)
	assert.NoError(t, asl(cpu, Absolute(1)))
	assert.Equal(t, 0b0000_0100, cpu.memory.Read(1))

	cpu.memory.Write(4, 0b0000_0010)
	cpu.X = 3
	assert.NoError(t, asl(cpu, Absolute(1), cpu.X))
	assert.Equal(t, 0b0000_0100, cpu.memory.Read(4))
}

func TestBcc(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, bcc(cpu, Absolute(123)))
	assert.Equal(t, 123, cpu.PC)

	cpu.PC = nes.CodeBaseAddress
	cpu.Flags.C = 1
	assert.NoError(t, bcc(cpu, Absolute(123)))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)
}

func TestBcs(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, bcs(cpu, Absolute(123)))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)

	cpu.Flags.C = 1
	assert.NoError(t, bcs(cpu, Absolute(123)))
	assert.Equal(t, 123, cpu.PC)
}

func TestBeq(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, beq(cpu, Absolute(123)))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)

	cpu.Flags.Z = 1
	assert.NoError(t, beq(cpu, Absolute(123)))
	assert.Equal(t, 123, cpu.PC)
}

func TestBit(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "value 1",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x100, 1)
				cpu.A = 1
				assert.NoError(t, bit(cpu, Absolute(0x100)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.A)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.V)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "value 0xff",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x100, 0xff)
				cpu.A = 0xf0
				assert.NoError(t, bit(cpu, Absolute(0x100)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xf0, cpu.A)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 1, cpu.Flags.V)
				assert.Equal(t, 1, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestBmi(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, bmi(cpu, Absolute(123)))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)

	cpu.Flags.N = 1
	assert.NoError(t, bmi(cpu, Absolute(123)))
	assert.Equal(t, 123, cpu.PC)
}

func TestBne(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, bne(cpu, Absolute(123)))
	assert.Equal(t, 123, cpu.PC)

	cpu.PC = nes.CodeBaseAddress
	cpu.Flags.Z = 1
	assert.NoError(t, bne(cpu, Absolute(123)))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)
}

func TestBpl(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, bpl(cpu, Absolute(123)))
	assert.Equal(t, 123, cpu.PC)

	cpu.PC = nes.CodeBaseAddress
	cpu.Flags.N = 1
	assert.NoError(t, bpl(cpu, Absolute(123)))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)
}

func TestBrk(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, brk(cpu))

	assert.Equal(t, testIrqAddress, cpu.PC)
}

func TestBvc(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, bvc(cpu, Absolute(123)))
	assert.Equal(t, 123, cpu.PC)

	cpu.PC = nes.CodeBaseAddress
	cpu.Flags.V = 1
	assert.NoError(t, bvc(cpu, Absolute(123)))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)
}

func TestBvs(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, bvs(cpu, Absolute(123)))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)

	cpu.Flags.V = 1
	assert.NoError(t, bvs(cpu, Absolute(123)))
	assert.Equal(t, 123, cpu.PC)
}

func TestClc(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.Flags.C = 1
	assert.NoError(t, clc(cpu))

	assert.Equal(t, 0, cpu.Flags.C)
}

func TestCld(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.Flags.D = 1
	assert.NoError(t, cld(cpu))

	assert.Equal(t, 0, cpu.Flags.D)
}

func TestCli(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.Flags.I = 1
	assert.NoError(t, cli(cpu))

	assert.Equal(t, 0, cpu.Flags.I)
}

func TestClv(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.Flags.V = 1
	assert.NoError(t, clv(cpu))

	assert.Equal(t, 0, cpu.Flags.V)
}

func TestCmp(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "equal immediate",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x100, 1)
				cpu.A = 1
				assert.NoError(t, cmp(cpu, 1))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
				assert.Equal(t, 1, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "unequal absolute",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x100, 0xff)
				cpu.A = 1
				assert.NoError(t, cmp(cpu, Absolute(0x100)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.A)
				assert.Equal(t, 0, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestCpx(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "equal immediate",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x100, 1)
				cpu.X = 1
				assert.NoError(t, cpx(cpu, 1))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.X)
				assert.Equal(t, 1, cpu.Flags.C)
				assert.Equal(t, 1, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "unequal absolute",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x100, 0xff)
				cpu.X = 1
				assert.NoError(t, cpx(cpu, Absolute(0x100)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.X)
				assert.Equal(t, 0, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestCpy(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "equal immediate",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x100, 1)
				cpu.Y = 1
				assert.NoError(t, cpy(cpu, 1))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.Y)
				assert.Equal(t, 1, cpu.Flags.C)
				assert.Equal(t, 1, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "unequal absolute",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x100, 0xff)
				cpu.Y = 1
				assert.NoError(t, cpy(cpu, Absolute(0x100)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.Y)
				assert.Equal(t, 0, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestDec(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "zeropage",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(1, 2)
				assert.NoError(t, dec(cpu, ZeroPage(1)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.memory.Read(1))
			},
		},
		{
			Name: "zeropage x",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(2, 2)
				cpu.X = 1
				assert.NoError(t, dec(cpu, ZeroPage(1), &cpu.X))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.memory.Read(2))
			},
		},
		{
			Name: "absolute",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x101, 2)
				assert.NoError(t, dec(cpu, Absolute(0x101)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.memory.Read(0x101))
			},
		},
		{
			Name: "absolute x",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x102, 2)
				cpu.X = 1
				assert.NoError(t, dec(cpu, Absolute(0x101), &cpu.X))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.memory.Read(0x102))
			},
		},
	}
	runCPUTest(t, tests)
}

func TestDex(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.X = 2
	assert.NoError(t, dex(cpu))

	assert.Equal(t, 1, cpu.X)
}

func TestDey(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.Y = 2
	assert.NoError(t, dey(cpu))

	assert.Equal(t, 1, cpu.Y)
}

func TestEor(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// TODO add test
	assert.NoError(t, eor(cpu, 0))
}

func TestInc(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "zeropage",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(1, 1)
				assert.NoError(t, inc(cpu, ZeroPage(1)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 2, cpu.memory.Read(1))
			},
		},
		{
			Name: "zeropage x",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(2, 1)
				cpu.X = 1
				assert.NoError(t, inc(cpu, ZeroPage(1), &cpu.X))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 2, cpu.memory.Read(2))
			},
		},
		{
			Name: "absolute",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x101, 1)
				assert.NoError(t, inc(cpu, Absolute(0x101)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 2, cpu.memory.Read(0x101))
			},
		},
		{
			Name: "absolute x",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x102, 1)
				cpu.X = 1
				assert.NoError(t, inc(cpu, Absolute(0x101), &cpu.X))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 2, cpu.memory.Read(0x102))
			},
		},
	}
	runCPUTest(t, tests)
}

func TestInx(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, inx(cpu))

	assert.Equal(t, 1, cpu.X)
}

func TestIny(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, iny(cpu))

	assert.Equal(t, 1, cpu.Y)
}

func TestJmp(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "absolute",
			Setup: func(cpu *CPU) {
				assert.NoError(t, jmp(cpu, Absolute(0x100)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x100, cpu.PC)
			},
		},
		{
			Name: "indirect",
			Setup: func(cpu *CPU) {
				cpu.memory.WriteWord(0x100, 0x200)
				assert.NoError(t, jmp(cpu, Indirect(0x100)))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x200, cpu.PC)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestJsr(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)
	assert.NoError(t, jsr(cpu, Absolute(0x101)))

	assert.Equal(t, InitialStack-2, cpu.SP)
	assert.Equal(t, 0x101, cpu.PC)
	w := cpu.pop16()
	assert.Equal(t, 0x8002, w)
}

func TestLda(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, lda(cpu, 1))

	assert.Equal(t, 1, cpu.A)
}

func TestLdx(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "immediate",
			Setup: func(cpu *CPU) {
				assert.NoError(t, ldx(cpu, 1))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.X)
			},
		},
		{
			Name: "zeropage y",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(2, 8)
				cpu.Y = 1
				assert.NoError(t, ldx(cpu, ZeroPage(1), &cpu.Y))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 8, cpu.X)
			},
		},
		{
			Name: "absolute y",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x102, 8)
				cpu.Y = 1
				assert.NoError(t, ldx(cpu, Absolute(0x101), &cpu.Y))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 8, cpu.X)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestLdy(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "immediate",
			Setup: func(cpu *CPU) {
				assert.NoError(t, ldy(cpu, 1))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.Y)
			},
		},
		{
			Name: "zeropage x",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(2, 8)
				cpu.X = 1
				assert.NoError(t, ldy(cpu, ZeroPage(1), &cpu.X))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 8, cpu.Y)
			},
		},
		{
			Name: "absolute x",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x102, 8)
				cpu.X = 1
				assert.NoError(t, ldy(cpu, Absolute(0x101), &cpu.X))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 8, cpu.Y)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestLsr(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "value 0b0000_0010 accumulator",
			Setup: func(cpu *CPU) {
				cpu.A = 0b0000_0010
				assert.NoError(t, lsr(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0b0000_0001, cpu.A)
				assert.Equal(t, 0, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "value 0b0111_1111 accumulator",
			Setup: func(cpu *CPU) {
				cpu.A = 0b0111_1111
				assert.NoError(t, lsr(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0b0011_1111, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "value 0b0111_1111 absolute",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x101, 0b0111_1111)
				assert.NoError(t, lsr(cpu, Absolute(0x101)))
			},
			Check: func(cpu *CPU) {
				b := cpu.memory.Read(0x101)
				assert.Equal(t, 0b0011_1111, b)
				assert.Equal(t, 0, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestNop(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, nop(cpu))
}

func TestOra(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// TODO add test
	assert.NoError(t, ora(cpu, 0))
}

func TestPha(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.A = 1
	assert.NoError(t, pha(cpu))

	b := cpu.memory.Read(StackBase + InitialStack)
	assert.Equal(t, cpu.A, b)
	assert.Equal(t, StackBase+InitialStack-1, cpu.SP)
}

func TestPhp(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, php(cpu))

	b := cpu.memory.Read(StackBase + InitialStack)
	assert.Equal(t, 0b0011_0100, b)
}

func TestPla(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.SP = 1
	cpu.memory.Write(StackBase+2, 1)
	assert.NoError(t, pla(cpu))

	assert.Equal(t, 1, cpu.A)
	assert.Equal(t, 2, cpu.SP)
}

func TestPlp(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.SP = 1
	cpu.memory.Write(StackBase+2, 1)
	assert.NoError(t, plp(cpu))

	assert.Equal(t, 0b0010_0001, cpu.GetFlags())
	assert.Equal(t, 2, cpu.SP)
}

func TestRol(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "value 0b0000_0010 accumulator",
			Setup: func(cpu *CPU) {
				cpu.A = 0b0000_0010
				assert.NoError(t, rol(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0b0000_0100, cpu.A)
				assert.Equal(t, 0, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "value 0b1111_1110 accumulator C0",
			Setup: func(cpu *CPU) {
				cpu.A = 0b1111_1110
				assert.NoError(t, rol(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0b1111_1100, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 1, cpu.Flags.N)
			},
		},
		{
			Name: "value 0b1111_1110 absolute C1",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x101, 0b1111_1110)
				cpu.Flags.C = 1
				assert.NoError(t, rol(cpu, Absolute(0x101)))
			},
			Check: func(cpu *CPU) {
				b := cpu.memory.Read(0x101)
				assert.Equal(t, 0b1111_1101, b)
				assert.Equal(t, 0, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 1, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestRor(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "value 0b0000_0010 accumulator",
			Setup: func(cpu *CPU) {
				cpu.A = 0b0000_0010
				assert.NoError(t, ror(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0b0000_0001, cpu.A)
				assert.Equal(t, 0, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "value 0b0111_1111 accumulator C0",
			Setup: func(cpu *CPU) {
				cpu.A = 0b0111_1111
				assert.NoError(t, ror(cpu))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0b0011_1111, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "value 0b0111_1111 absolute C1",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x101, 0b0111_1111)
				cpu.Flags.C = 1
				assert.NoError(t, ror(cpu, Absolute(0x101)))
			},
			Check: func(cpu *CPU) {
				b := cpu.memory.Read(0x101)
				assert.Equal(t, 0b1011_1111, b)
				assert.Equal(t, 0, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 1, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestRti(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, rti(cpu))
}

func TestRts(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.push16(0x100)
	assert.NoError(t, rts(cpu))
	assert.Equal(t, 0x101, cpu.PC)
}

func TestSbc(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "result 0xff C0",
			Setup: func(cpu *CPU) {
				cpu.A = 2
				assert.NoError(t, sbc(cpu, 2))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xff, cpu.A)
				assert.Equal(t, 0, cpu.Flags.C)
			},
		},
		{
			Name: "result 0xfe C0",
			Setup: func(cpu *CPU) {
				assert.NoError(t, sbc(cpu, 1))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xfe, cpu.A)
				assert.Equal(t, 0, cpu.Flags.C)
			},
		},
		{
			Name: "result 0x00 C1",
			Setup: func(cpu *CPU) {
				cpu.Flags.C = 1
				assert.NoError(t, sbc(cpu, 0))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x00, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestSec(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, sec(cpu))

	assert.Equal(t, 1, cpu.Flags.C)
}

func TestSed(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, sed(cpu))

	assert.Equal(t, 1, cpu.Flags.D)
}

func TestSei(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, sei(cpu))

	assert.Equal(t, 1, cpu.Flags.I)
}

func TestSta(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.A = 11
	assert.NoError(t, sta(cpu, 0))

	b := cpu.memory.Read(0)
	assert.Equal(t, cpu.A, b)

	cpu.X = 0x22
	assert.NoError(t, sta(cpu, Absolute(0), &cpu.X))

	b = cpu.memory.Read(0x22)
	assert.Equal(t, cpu.A, b)
}

func TestStx(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.X = 11
	assert.NoError(t, stx(cpu, 0))

	b := cpu.memory.Read(0)
	assert.Equal(t, cpu.X, b)

	cpu.Y = 0x22
	assert.NoError(t, stx(cpu, Absolute(0), &cpu.Y))

	b = cpu.memory.Read(0x22)
	assert.Equal(t, cpu.X, b)
}

func TestSty(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.Y = 11
	assert.NoError(t, sty(cpu, 0))

	b := cpu.memory.Read(0)
	assert.Equal(t, cpu.Y, b)

	cpu.X = 0x22
	assert.NoError(t, sty(cpu, Absolute(0), &cpu.X))

	b = cpu.memory.Read(0x22)
	assert.Equal(t, cpu.Y, b)
}

func TestTax(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.A = 2
	assert.NoError(t, tax(cpu))

	assert.Equal(t, 2, cpu.X)
}

func TestTay(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.A = 2
	assert.NoError(t, tay(cpu))

	assert.Equal(t, 2, cpu.Y)
}

func TestTsx(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	assert.NoError(t, tsx(cpu))

	assert.Equal(t, InitialStack, cpu.SP)
	assert.Equal(t, InitialStack, cpu.X)
}

func TestTxa(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.X = 2
	assert.NoError(t, txa(cpu))

	assert.Equal(t, 2, cpu.A)
}

func TestTxs(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.X = 2
	assert.NoError(t, txs(cpu))

	assert.Equal(t, 2, cpu.SP)
}

func TestTya(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	cpu.Y = 2
	assert.NoError(t, tya(cpu))

	assert.Equal(t, 2, cpu.A)
}

// Tests for validation functions

func TestGetFlags(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Set all flags to 1
	cpu.Flags.C = 1
	cpu.Flags.Z = 1
	cpu.Flags.I = 1
	cpu.Flags.D = 1
	cpu.Flags.B = 1
	cpu.Flags.U = 1
	cpu.Flags.V = 1
	cpu.Flags.N = 1

	result := cpu.GetFlags()
	assert.Equal(t, 0xFF, result)

	// Test with mixed flags
	cpu.Flags.C = 1
	cpu.Flags.Z = 0
	cpu.Flags.I = 1
	cpu.Flags.D = 0
	cpu.Flags.B = 1
	cpu.Flags.U = 0
	cpu.Flags.V = 1
	cpu.Flags.N = 0

	result = cpu.GetFlags()
	expected := uint8(0b01010101) // C=1, Z=0, I=1, D=0, B=1, U=0, V=1, N=0
	assert.Equal(t, expected, result)
}

func TestSetZN(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "zero value",
			Setup: func(cpu *CPU) {
				cpu.setZN(0x00)
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
		{
			Name: "negative value",
			Setup: func(cpu *CPU) {
				cpu.setZN(0x80)
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 1, cpu.Flags.N)
			},
		},
		{
			Name: "positive value",
			Setup: func(cpu *CPU) {
				cpu.setZN(0x7F)
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0, cpu.Flags.Z)
				assert.Equal(t, 0, cpu.Flags.N)
			},
		},
	}
	runCPUTest(t, tests)
}

func TestValidateState(t *testing.T) {
	t.Parallel()

	// Test valid state
	cpu := cpuTestSetup(t)
	assert.NoError(t, cpu.ValidateState())

	// Test invalid flag values
	cpu = cpuTestSetup(t)
	cpu.Flags.C = 2
	assert.ErrorContains(t, cpu.ValidateState(), "invalid flag values")

	// Test nil memory
	cpu = cpuTestSetup(t)
	cpu.memory = nil
	assert.ErrorContains(t, cpu.ValidateState(), "CPU memory is nil")
}

func TestCPUReset(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Modify CPU state
	cpu.A = 0x42
	cpu.X = 0x24
	cpu.Y = 0x84
	cpu.Flags.C = 1
	cpu.Flags.Z = 1
	cpu.cycles = 1000

	// Reset and verify state
	cpu.Reset()
	assert.Equal(t, 0, cpu.A)
	assert.Equal(t, 0, cpu.X)
	assert.Equal(t, 0, cpu.Y)
	assert.Equal(t, InitialStack, cpu.SP)
	assert.Equal(t, initialCycles, cpu.cycles)
	assert.Equal(t, initialFlags, cpu.GetFlags())
}

func TestGetInstructionCount(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)

	// Test initial instruction count
	count := cpu.GetInstructionCount()
	assert.Equal(t, initialCycles/4, count)

	// Add cycles and test again
	cpu.cycles = 100
	count = cpu.GetInstructionCount()
	assert.Equal(t, uint64(25), count)
}
