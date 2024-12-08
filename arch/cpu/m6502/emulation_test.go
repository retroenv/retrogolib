package m6502

import (
	"testing"

	. "github.com/retroenv/retrogolib/addressing"
	"github.com/retroenv/retrogolib/arch/nes"
	"github.com/retroenv/retrogolib/assert"
)

type cpuTest struct {
	Name  string
	Setup func(cpu *CPU)
	Check func(cpu *CPU)
}

const testIrqAddress = 0x9000

func cpuTestSetup() *CPU {
	memory := NewMemory(&testMemory{})
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
			cpu := cpuTestSetup()
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
				adc(cpu, 0xff)
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.A)
				assert.Equal(t, 1, cpu.Flags.C)
			},
		},
		{
			Name: "result 0x01",
			Setup: func(cpu *CPU) {
				adc(cpu, 1)
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
				adc(cpu, 0xff)
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
	cpu := cpuTestSetup()

	cpu.A = 0x12
	and(cpu, 2)

	assert.Equal(t, 2, cpu.A)
}

func TestAsl(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.A = 0b0000_0001
	asl(cpu)
	assert.Equal(t, 0b0000_0010, cpu.A)
	assert.Equal(t, 0, cpu.Flags.C)

	cpu.A = 0b1111_1110
	asl(cpu)
	assert.Equal(t, 0b1111_1100, cpu.A)
	assert.Equal(t, 1, cpu.Flags.C)

	cpu.memory.Write(1, 0b0000_0010)
	asl(cpu, Absolute(1))
	assert.Equal(t, 0b0000_0100, cpu.memory.Read(1))

	cpu.memory.Write(4, 0b0000_0010)
	cpu.X = 3
	asl(cpu, Absolute(1), cpu.X)
	assert.Equal(t, 0b0000_0100, cpu.memory.Read(4))
}

func TestBcc(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	bcc(cpu, Absolute(123))
	assert.Equal(t, 123, cpu.PC)

	cpu.PC = nes.CodeBaseAddress
	cpu.Flags.C = 1
	bcc(cpu, Absolute(123))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)
}

func TestBcs(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	bcs(cpu, Absolute(123))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)

	cpu.Flags.C = 1
	bcs(cpu, Absolute(123))
	assert.Equal(t, 123, cpu.PC)
}

func TestBeq(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	beq(cpu, Absolute(123))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)

	cpu.Flags.Z = 1
	beq(cpu, Absolute(123))
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
				bit(cpu, Absolute(0x100))
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
				bit(cpu, Absolute(0x100))
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
	cpu := cpuTestSetup()

	bmi(cpu, Absolute(123))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)

	cpu.Flags.N = 1
	bmi(cpu, Absolute(123))
	assert.Equal(t, 123, cpu.PC)
}

func TestBne(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	bne(cpu, Absolute(123))
	assert.Equal(t, 123, cpu.PC)

	cpu.PC = nes.CodeBaseAddress
	cpu.Flags.Z = 1
	bne(cpu, Absolute(123))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)
}

func TestBpl(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	bpl(cpu, Absolute(123))
	assert.Equal(t, 123, cpu.PC)

	cpu.PC = nes.CodeBaseAddress
	cpu.Flags.N = 1
	bpl(cpu, Absolute(123))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)
}

func TestBrk(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	brk(cpu)

	assert.Equal(t, testIrqAddress, cpu.PC)
}

func TestBvc(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	bvc(cpu, Absolute(123))
	assert.Equal(t, 123, cpu.PC)

	cpu.PC = nes.CodeBaseAddress
	cpu.Flags.V = 1
	bvc(cpu, Absolute(123))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)
}

func TestBvs(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	bvs(cpu, Absolute(123))
	assert.Equal(t, nes.CodeBaseAddress, cpu.PC)

	cpu.Flags.V = 1
	bvs(cpu, Absolute(123))
	assert.Equal(t, 123, cpu.PC)
}

func TestClc(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.Flags.C = 1
	clc(cpu)

	assert.Equal(t, 0, cpu.Flags.C)
}

func TestCld(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.Flags.D = 1
	cld(cpu)

	assert.Equal(t, 0, cpu.Flags.D)
}

func TestCli(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.Flags.I = 1
	cli(cpu)

	assert.Equal(t, 0, cpu.Flags.I)
}

func TestClv(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.Flags.V = 1
	clv(cpu)

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
				cmp(cpu, 1)
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
				cmp(cpu, Absolute(0x100))
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
				cpx(cpu, 1)
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
				cpx(cpu, Absolute(0x100))
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
				cpy(cpu, 1)
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
				cpy(cpu, Absolute(0x100))
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
				dec(cpu, ZeroPage(1))
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
				dec(cpu, ZeroPage(1), &cpu.X)
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 1, cpu.memory.Read(2))
			},
		},
		{
			Name: "absolute",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x101, 2)
				dec(cpu, Absolute(0x101))
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
				dec(cpu, Absolute(0x101), &cpu.X)
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
	cpu := cpuTestSetup()

	cpu.X = 2
	dex(cpu)

	assert.Equal(t, 1, cpu.X)
}

func TestDey(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.Y = 2
	dey(cpu)

	assert.Equal(t, 1, cpu.Y)
}

func TestEor(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	// TODO add test
	eor(cpu, 0)
}

func TestInc(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "zeropage",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(1, 1)
				inc(cpu, ZeroPage(1))
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
				inc(cpu, ZeroPage(1), &cpu.X)
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 2, cpu.memory.Read(2))
			},
		},
		{
			Name: "absolute",
			Setup: func(cpu *CPU) {
				cpu.memory.Write(0x101, 1)
				inc(cpu, Absolute(0x101))
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
				inc(cpu, Absolute(0x101), &cpu.X)
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
	cpu := cpuTestSetup()

	inx(cpu)

	assert.Equal(t, 1, cpu.X)
}

func TestIny(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	iny(cpu)

	assert.Equal(t, 1, cpu.Y)
}

func TestJmp(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "absolute",
			Setup: func(cpu *CPU) {
				jmp(cpu, Absolute(0x100))
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0x100, cpu.PC)
			},
		},
		{
			Name: "indirect",
			Setup: func(cpu *CPU) {
				cpu.memory.WriteWord(0x100, 0x200)
				jmp(cpu, Indirect(0x100))
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
	cpu := cpuTestSetup()
	jsr(cpu, Absolute(0x101))

	assert.Equal(t, InitialStack-2, cpu.SP)
	assert.Equal(t, 0x101, cpu.PC)
	w := cpu.pop16()
	assert.Equal(t, 0x8002, w)
}

func TestLda(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	lda(cpu, 1)

	assert.Equal(t, 1, cpu.A)
}

func TestLdx(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "immediate",
			Setup: func(cpu *CPU) {
				ldx(cpu, 1)
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
				ldx(cpu, ZeroPage(1), &cpu.Y)
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
				ldx(cpu, Absolute(0x101), &cpu.Y)
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
				ldy(cpu, 1)
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
				ldy(cpu, ZeroPage(1), &cpu.X)
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
				ldy(cpu, Absolute(0x101), &cpu.X)
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
				lsr(cpu)
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
				lsr(cpu)
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
				lsr(cpu, Absolute(0x101))
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
	cpu := cpuTestSetup()

	nop(cpu)
}

func TestOra(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	// TODO add test
	ora(cpu, 0)
}

func TestPha(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.A = 1
	pha(cpu)

	b := cpu.memory.Read(StackBase + InitialStack)
	assert.Equal(t, cpu.A, b)
	assert.Equal(t, StackBase+InitialStack-1, cpu.SP)
}

func TestPhp(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	php(cpu)

	b := cpu.memory.Read(StackBase + InitialStack)
	assert.Equal(t, 0b0011_0100, b)
}

func TestPla(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.SP = 1
	cpu.memory.Write(StackBase+2, 1)
	pla(cpu)

	assert.Equal(t, 1, cpu.A)
	assert.Equal(t, 2, cpu.SP)
}

func TestPlp(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.SP = 1
	cpu.memory.Write(StackBase+2, 1)
	plp(cpu)

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
				rol(cpu)
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
				rol(cpu)
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
				rol(cpu, Absolute(0x101))
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
				ror(cpu)
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
				ror(cpu)
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
				ror(cpu, Absolute(0x101))
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
	cpu := cpuTestSetup()

	rti(cpu)
}

func TestRts(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.push16(0x100)
	rts(cpu)
	assert.Equal(t, 0x101, cpu.PC)
}

func TestSbc(t *testing.T) {
	t.Parallel()
	tests := []cpuTest{
		{
			Name: "result 0xff C0",
			Setup: func(cpu *CPU) {
				cpu.A = 2
				sbc(cpu, 2)
			},
			Check: func(cpu *CPU) {
				assert.Equal(t, 0xff, cpu.A)
				assert.Equal(t, 0, cpu.Flags.C)
			},
		},
		{
			Name: "result 0xfe C0",
			Setup: func(cpu *CPU) {
				sbc(cpu, 1)
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
				sbc(cpu, 0)
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
	cpu := cpuTestSetup()

	sec(cpu)

	assert.Equal(t, 1, cpu.Flags.C)
}

func TestSed(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	sed(cpu)

	assert.Equal(t, 1, cpu.Flags.D)
}

func TestSei(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	sei(cpu)

	assert.Equal(t, 1, cpu.Flags.I)
}

func TestSta(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.A = 11
	sta(cpu, 0)

	b := cpu.memory.Read(0)
	assert.Equal(t, cpu.A, b)

	cpu.X = 0x22
	sta(cpu, Absolute(0), &cpu.X)

	b = cpu.memory.Read(0x22)
	assert.Equal(t, cpu.A, b)
}

func TestStx(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.X = 11
	stx(cpu, 0)

	b := cpu.memory.Read(0)
	assert.Equal(t, cpu.X, b)

	cpu.Y = 0x22
	stx(cpu, Absolute(0), &cpu.Y)

	b = cpu.memory.Read(0x22)
	assert.Equal(t, cpu.X, b)
}

func TestSty(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.Y = 11
	sty(cpu, 0)

	b := cpu.memory.Read(0)
	assert.Equal(t, cpu.Y, b)

	cpu.X = 0x22
	sty(cpu, Absolute(0), &cpu.X)

	b = cpu.memory.Read(0x22)
	assert.Equal(t, cpu.Y, b)
}

func TestTax(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.A = 2
	tax(cpu)

	assert.Equal(t, 2, cpu.X)
}

func TestTay(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.A = 2
	tay(cpu)

	assert.Equal(t, 2, cpu.Y)
}

func TestTsx(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	tsx(cpu)

	assert.Equal(t, InitialStack, cpu.SP)
	assert.Equal(t, InitialStack, cpu.X)
}

func TestTxa(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.X = 2
	txa(cpu)

	assert.Equal(t, 2, cpu.A)
}

func TestTxs(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.X = 2
	txs(cpu)

	assert.Equal(t, 2, cpu.SP)
}

func TestTya(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup()

	cpu.Y = 2
	tya(cpu)

	assert.Equal(t, 2, cpu.A)
}
