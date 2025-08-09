package z80

import (
	"testing"
	"unsafe"
)

// Performance benchmarks for critical Z80 CPU operations

// CPU Core Operation Benchmarks

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

// Memory Operation Benchmarks

func BenchmarkMemory_Read(b *testing.B) {
	memory := NewMemory()
	memory.Write(0x1000, 0x42)

	b.ResetTimer()
	for range b.N {
		memory.Read(0x1000)
	}
}

func BenchmarkMemory_Write(b *testing.B) {
	memory := NewMemory()

	b.ResetTimer()
	for range b.N {
		memory.Write(0x1000, 0x42)
	}
}

func BenchmarkMemory_ReadWord(b *testing.B) {
	memory := NewMemory()
	memory.WriteWord(0x1000, 0x1234)

	b.ResetTimer()
	for range b.N {
		memory.ReadWord(0x1000)
	}
}

func BenchmarkMemory_WriteWord(b *testing.B) {
	memory := NewMemory()

	b.ResetTimer()
	for range b.N {
		memory.WriteWord(0x1000, 0x1234)
	}
}

func BenchmarkMemory_LoadROM(b *testing.B) {
	rom := make([]byte, 0x8000) // 32KB ROM
	for i := range rom {
		rom[i] = uint8(i & 0xFF)
	}

	b.ResetTimer()
	for range b.N {
		memory := NewMemory()
		memory.LoadROM(rom)
	}
}

// Performance Comparison Benchmarks

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

func BenchmarkStringParamComparison(b *testing.B) {
	strA, strB := StrRegB, StrRegC

	b.ResetTimer()
	for range b.N {
		_ = strA == strB
	}
}

func BenchmarkIntParamComparison(b *testing.B) {
	intA, intB := IntRegB, IntRegC

	b.ResetTimer()
	for range b.N {
		_ = intA == intB
	}
}

func BenchmarkRegisterArrayAccess(b *testing.B) {
	// Simulate register array access pattern used in CPU emulation
	registers := [8]uint8{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80}
	regIndex := 2 // RegC equivalent

	b.ResetTimer()
	for range b.N {
		_ = registers[regIndex]
	}
}

func BenchmarkStringMapLookup(b *testing.B) {
	// Simulate string-based register lookup
	registerMap := map[string]uint8{
		"b": 0x10,
		"c": 0x20,
		"d": 0x30,
		"e": 0x40,
		"h": 0x50,
		"l": 0x60,
		"a": 0x70,
		"f": 0x80,
	}

	b.ResetTimer()
	for range b.N {
		_ = registerMap["c"]
	}
}

func BenchmarkIntSwitchStatement(b *testing.B) {
	regValue := IntRegC
	var result uint8

	b.ResetTimer()
	for range b.N {
		switch regValue {
		case IntRegB:
			result = 0x10
		case IntRegC:
			result = 0x20
		case IntRegHL:
			result = 0x30
		case IntImm8:
			result = 0x40
		}
	}
	_ = result
}

func BenchmarkStringSwitchStatement(b *testing.B) {
	regValue := StrRegC
	var result uint8

	b.ResetTimer()
	for range b.N {
		switch regValue {
		case StrRegB:
			result = 0x10
		case StrRegC:
			result = 0x20
		case StrRegHL:
			result = 0x30
		case StrImm8:
			result = 0x40
		}
	}
	_ = result
}

// Memory usage comparison
func TestMemoryUsageComparison(t *testing.T) {
	t.Run("Memory Usage Analysis", func(t *testing.T) {
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
}
