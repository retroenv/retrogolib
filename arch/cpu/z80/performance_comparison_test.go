package z80

import (
	"testing"
	"unsafe"
)

// Performance comparison between string and int-based RegisterParam

type StringParam string
type IntParam uint8

const (
	StrRegB  StringParam = "b"
	StrRegC  StringParam = "c"
	StrRegHL StringParam = "hl"
	StrImm8  StringParam = "n"
)

const (
	IntRegB IntParam = iota + 1
	IntRegC
	IntRegHL
	IntImm8
)

func TestPerformanceComparison(t *testing.T) {
	t.Run("Memory Usage", func(t *testing.T) {
		t.Logf("String param size: %d bytes", unsafe.Sizeof(StringParam("")))
		t.Logf("Int param size: %d bytes", unsafe.Sizeof(IntParam(0)))

		// String includes 16 bytes (pointer + length on 64-bit)
		// Int is just 1 byte
		savings := unsafe.Sizeof(StringParam("")) - unsafe.Sizeof(IntParam(0))
		t.Logf("Memory savings per register field: %d bytes", savings)

		// For Opcode struct with 3 register fields * 256 opcodes
		totalSavings := savings * 3 * 256
		t.Logf("Total memory savings for all opcodes: %d bytes", totalSavings)
	})

	t.Run("Comparison Speed", func(t *testing.T) {
		// String comparison requires content comparison
		strA, strB := StrRegB, StrRegC
		equal1 := strA == strB

		// Int comparison is single instruction
		intA, intB := IntRegB, IntRegC
		equal2 := intA == intB

		t.Logf("String comparison result: %t", equal1)
		t.Logf("Int comparison result: %t", equal2)

		// Int comparison is much faster in CPU cycles
	})
}

// Demonstrate CPU emulation efficiency gains
func TestCPUEmulationEfficiency(t *testing.T) {
	t.Run("Register Array Access", func(t *testing.T) {
		// With int-based params, can use direct array indexing
		registers := [8]byte{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80}

		// Fast array access using int as index
		reg := IntRegC
		if reg > 0 && int(reg) < len(registers) {
			value := registers[reg]
			t.Logf("Register %d value: 0x%02X (direct array access)", reg, value)
		}

		// String-based would require map lookup (much slower)
		registerMap := map[StringParam]byte{
			StrRegB: 0x20,
			StrRegC: 0x30,
		}
		if value, exists := registerMap[StrRegC]; exists {
			t.Logf("String-based lookup value: 0x%02X (map lookup required)", value)
		}
	})

	t.Run("Switch Statement Optimization", func(t *testing.T) {
		// Int-based switch can be optimized to jump table by compiler
		reg := IntRegB
		var cycles int
		switch reg {
		case IntRegB:
			cycles = 4
		case IntRegC:
			cycles = 4
		case IntRegHL:
			cycles = 6
		default:
			cycles = 0
		}
		t.Logf("Int-based switch: %d cycles (jump table)", cycles)

		// String-based switch requires string comparisons
		strReg := StrRegB
		switch strReg {
		case StrRegB:
			cycles = 4
		case StrRegC:
			cycles = 4
		case StrRegHL:
			cycles = 6
		default:
			cycles = 0
		}
		t.Logf("String-based switch: %d cycles (string comparisons)", cycles)
	})
}

// Benchmark comparison
func BenchmarkStringComparison(b *testing.B) {
	reg1, reg2 := StrRegB, StrRegC
	for range b.N {
		_ = reg1 == reg2
	}
}

func BenchmarkIntComparison(b *testing.B) {
	reg1, reg2 := IntRegB, IntRegC
	for range b.N {
		_ = reg1 == reg2
	}
}

func BenchmarkStringSwitch(b *testing.B) {
	reg := StrRegB
	var result int
	for range b.N {
		switch reg {
		case StrRegB:
			result = 1
		case StrRegC:
			result = 2
		case StrRegHL:
			result = 3
		default:
			result = 0
		}
	}
	_ = result
}

func BenchmarkIntSwitch(b *testing.B) {
	reg := IntRegB
	var result int
	for range b.N {
		switch reg {
		case IntRegB:
			result = 1
		case IntRegC:
			result = 2
		case IntRegHL:
			result = 3
		default:
			result = 0
		}
	}
	_ = result
}
