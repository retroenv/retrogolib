package cpu68000

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestInstructionTiming(t *testing.T) {
	// Timing formerly used fixed decoder values regardless of operands or outcome.
	tests := []struct {
		name   string
		words  []uint16
		cycles uint64
	}{
		{name: "MOVE word indirect", words: []uint16{0x3010}, cycles: 8},
		{name: "MOVE long indirect", words: []uint16{0x2010}, cycles: 12},
		{name: "MOVE long displacement", words: []uint16{0x2028, 0}, cycles: 16},
		{name: "MOVE long indexed", words: []uint16{0x2030, 0}, cycles: 18},
		{name: "MOVE long absolute", words: []uint16{0x2039, 0, 0x3000}, cycles: 20},
		{name: "MOVE long immediate", words: []uint16{0x203C, 0, 1}, cycles: 12},
		{name: "MOVE long destination predecrement", words: []uint16{0x2100}, cycles: 12},
		{name: "ADD long register", words: []uint16{0xD081}, cycles: 8},
		{name: "ADD long memory source", words: []uint16{0xD090}, cycles: 14},
		{name: "ADD long memory destination", words: []uint16{0xD190}, cycles: 20},
		{name: "ADDI word memory", words: []uint16{0x0650, 1}, cycles: 16},
		{name: "ANDI long register", words: []uint16{0x0280, 0, 1}, cycles: 14},
		{name: "BNE taken", words: []uint16{0x6602}, cycles: 10},
		{name: "BEQ short not taken", words: []uint16{0x6702}, cycles: 8},
		{name: "BEQ word not taken", words: []uint16{0x6700, 2}, cycles: 12},
		{name: "DBF expired", words: []uint16{0x51C8, 0xFFFE}, cycles: 14},
		{name: "DBT condition true", words: []uint16{0x50C8, 0xFFFE}, cycles: 12},
		{name: "ST register", words: []uint16{0x50C0}, cycles: 6},
		{name: "SF register", words: []uint16{0x51C0}, cycles: 4},
		{name: "LSL word immediate eight", words: []uint16{0xE148}, cycles: 22},
		{name: "LSL long register zero", words: []uint16{0xE3A8}, cycles: 8},
		{name: "MOVEM long two registers", words: []uint16{0x48D0, 3}, cycles: 24},
		{name: "MOVEM word two registers read", words: []uint16{0x4C90, 3}, cycles: 20},
		{name: "LEA indirect", words: []uint16{0x43D0}, cycles: 4},
		{name: "LEA indexed", words: []uint16{0x43F0, 0}, cycles: 12},
		{name: "JMP displacement", words: []uint16{0x4EE8, 0}, cycles: 10},
		{name: "JSR absolute long", words: []uint16{0x4EB9, 0, 0x3000}, cycles: 20},
		{name: "MULU zero", words: []uint16{0xC0FC, 0}, cycles: 42},
		{name: "MULU all ones", words: []uint16{0xC0FC, 0xFFFF}, cycles: 74},
		{name: "MULS all ones", words: []uint16{0xC1FC, 0xFFFF}, cycles: 44},
		{name: "MULS alternating", words: []uint16{0xC1FC, 0x5555}, cycles: 74},
		{name: "divide by zero", words: []uint16{0x80D0}, cycles: 42},
		{name: "TRAP", words: []uint16{0x4E40}, cycles: 34},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t)
			cpu.A[0] = 0x3000
			for i, word := range tt.words {
				cpu.bus.WriteWord(cpu.PC+uint32(i)*2, word)
			}
			assert.NoError(t, cpu.Step())
			assert.Equal(t, tt.cycles, cpu.Cycles())
		})
	}
}

func TestDivisionTiming(t *testing.T) {
	tests := []struct {
		name     string
		opcode   uint16
		dividend uint32
		divisor  uint32
		cycles   uint64
	}{
		{name: "unsigned zero quotient", opcode: 0x80C1, divisor: 1, cycles: 136},
		{name: "unsigned overflow", opcode: 0x80C1, dividend: 0x10000, divisor: 1, cycles: 10},
		{name: "signed zero quotient", opcode: 0x81C1, divisor: 1, cycles: 150},
		{name: "signed negative dividend", opcode: 0x81C1, dividend: 0xFFFFFFFF, divisor: 1, cycles: 156},
		{name: "signed absolute overflow", opcode: 0x81C1, dividend: 0x80000000, divisor: 0xFFFF, cycles: 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t)
			cpu.D[0], cpu.D[1] = tt.dividend, tt.divisor
			cpu.bus.WriteWord(cpu.PC, tt.opcode)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, tt.cycles, cpu.Cycles())
		})
	}
}
