package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestUndocumentedNopPrefixes(t *testing.T) {
	tests := []struct {
		name         string
		program      []uint8
		expectedPC   uint16
		expectedSize uint8
		description  string
	}{
		{
			name:         "DD prefix with invalid opcode",
			program:      []uint8{0xDD, 0xFF}, // DD followed by invalid opcode
			expectedPC:   1,                   // Should advance by 1 (treat DD as NOP)
			expectedSize: 1,
			description:  "DD prefix followed by invalid opcode should act as 4-cycle NOP",
		},
		{
			name:         "FD prefix with invalid opcode",
			program:      []uint8{0xFD, 0xFF}, // FD followed by invalid opcode
			expectedPC:   1,                   // Should advance by 1 (treat FD as NOP)
			expectedSize: 1,
			description:  "FD prefix followed by invalid opcode should act as 4-cycle NOP",
		},
		{
			name:         "DD prefix at end of program",
			program:      []uint8{0xDD}, // DD at end
			expectedPC:   1,             // Should advance by 1
			expectedSize: 1,
			description:  "DD prefix alone should act as 4-cycle NOP",
		},
		{
			name:         "FD prefix at end of program",
			program:      []uint8{0xFD}, // FD at end
			expectedPC:   1,             // Should advance by 1
			expectedSize: 1,
			description:  "FD prefix alone should act as 4-cycle NOP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewBasicMemory()
			cpu, err := New(memory)
			assert.NoError(t, err)

			// Load program into memory
			for i, b := range tt.program {
				memory.Write(uint16(i), b)
			}

			// Reset CPU state
			cpu.PC = 0
			initialCycles := cpu.cycles

			// Execute one instruction
			err = cpu.Step()
			assert.NoError(t, err, tt.description)

			// Verify PC advanced correctly
			assert.Equal(t, tt.expectedPC, cpu.PC, "PC should advance by %d", tt.expectedSize)

			// Verify timing (4 cycles for undocumented NOP)
			expectedCycles := initialCycles + 4
			assert.Equal(t, expectedCycles, cpu.cycles, "Should consume 4 cycles")
		})
	}
}

func TestUndocumentedNopInstructionDefinitions(t *testing.T) {
	// Test NopUndoc1 (DD prefix alone)
	assert.Equal(t, Nop.Name, NopUndoc1.Name, "NopUndoc1 should have same name as regular NOP")
	assert.True(t, NopUndoc1.Unofficial, "NopUndoc1 should be marked as unofficial")

	opcodeInfo := NopUndoc1.Addressing[ImpliedAddressing]
	assert.Equal(t, PrefixDD, opcodeInfo.Opcode, "NopUndoc1 should use DD prefix")
	assert.Equal(t, byte(1), opcodeInfo.Size, "NopUndoc1 should be 1 byte")
	assert.Equal(t, byte(4), opcodeInfo.Cycles, "NopUndoc1 should take 4 cycles")

	// Test NopUndoc2 (FD prefix alone)
	assert.Equal(t, Nop.Name, NopUndoc2.Name, "NopUndoc2 should have same name as regular NOP")
	assert.True(t, NopUndoc2.Unofficial, "NopUndoc2 should be marked as unofficial")

	opcodeInfo2 := NopUndoc2.Addressing[ImpliedAddressing]
	assert.Equal(t, PrefixFD, opcodeInfo2.Opcode, "NopUndoc2 should use FD prefix")
	assert.Equal(t, byte(1), opcodeInfo2.Size, "NopUndoc2 should be 1 byte")
	assert.Equal(t, byte(4), opcodeInfo2.Cycles, "NopUndoc2 should take 4 cycles")
}

func TestUndocumentedInstructionMap(t *testing.T) {
	// Test that undocumented NOPs are in the map
	nopUndoc1, exists1 := UnofficialInstructions["nop_undoc_dd"]
	assert.True(t, exists1, "nop_undoc_dd should be in UnofficialInstructions map")
	assert.Equal(t, NopUndoc1, nopUndoc1, "Should return NopUndoc1 instance")

	nopUndoc2, exists2 := UnofficialInstructions["nop_undoc_fd"]
	assert.True(t, exists2, "nop_undoc_fd should be in UnofficialInstructions map")
	assert.Equal(t, NopUndoc2, nopUndoc2, "Should return NopUndoc2 instance")

	// Test IsUnofficialInstruction function
	assert.True(t, IsUnofficialInstruction("nop_undoc_dd"), "Should recognize nop_undoc_dd as unofficial")
	assert.True(t, IsUnofficialInstruction("nop_undoc_fd"), "Should recognize nop_undoc_fd as unofficial")
	assert.False(t, IsUnofficialInstruction("nop"), "Regular nop should not be considered unofficial")
}

func TestValidDDFDInstructionsStillWork(t *testing.T) {
	tests := []struct {
		name     string
		program  []uint8
		testFunc func(t *testing.T, cpu *CPU)
	}{
		{
			name:    "DD 21 (LD IX,nn) should work normally",
			program: []uint8{0xDD, 0x21, 0x34, 0x12}, // LD IX,$1234
			testFunc: func(t *testing.T, cpu *CPU) {
				t.Helper()
				assert.Equal(t, uint16(4), cpu.PC, "PC should advance by 4")
				assert.Equal(t, uint16(0x1234), cpu.IX, "IX should be loaded with $1234")
			},
		},
		{
			name:    "FD 21 (LD IY,nn) should work normally",
			program: []uint8{0xFD, 0x21, 0x56, 0x78}, // LD IY,$7856
			testFunc: func(t *testing.T, cpu *CPU) {
				t.Helper()
				assert.Equal(t, uint16(4), cpu.PC, "PC should advance by 4")
				assert.Equal(t, uint16(0x7856), cpu.IY, "IY should be loaded with $7856")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := NewBasicMemory()
			cpu, err := New(memory)
			assert.NoError(t, err)

			// Load program into memory
			for i, b := range tt.program {
				memory.Write(uint16(i), b)
			}

			// Reset CPU state
			cpu.PC = 0

			// Execute one instruction
			err = cpu.Step()
			assert.NoError(t, err)

			// Run custom test
			tt.testFunc(t, cpu)
		})
	}
}

func TestMultipleUndocumentedNops(t *testing.T) {
	memory := NewBasicMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test DD with invalid opcode - should act as 1-byte NOP
	memory.Write(0x0000, 0xDD) // DD prefix
	memory.Write(0x0001, 0xFF) // Invalid DD opcode

	cpu.PC = 0
	initialCycles := cpu.cycles

	// Execute DD undocumented NOP
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(1), cpu.PC, "DD+invalid should advance PC by 1")
	assert.Equal(t, initialCycles+4, cpu.cycles, "DD+invalid should consume 4 cycles")

	// Test FD with invalid opcode - should act as 1-byte NOP
	memory.Write(0x0001, 0xFD) // FD prefix at current PC
	memory.Write(0x0002, 0xFF) // Invalid FD opcode

	// Execute FD undocumented NOP
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(2), cpu.PC, "FD+invalid should advance PC by 1")
	assert.Equal(t, initialCycles+8, cpu.cycles, "FD+invalid should consume another 4 cycles")

	// Test regular NOP
	memory.Write(0x0002, 0x00) // Regular NOP at current PC

	// Execute regular NOP
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(3), cpu.PC, "Regular NOP should advance PC by 1")
	assert.Equal(t, initialCycles+12, cpu.cycles, "Regular NOP should consume another 4 cycles")
}
