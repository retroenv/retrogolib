package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestUndocumentedDDFDPassthrough(t *testing.T) {
	t.Parallel()

	// DD/FD prefix with non-IX/IY opcode falls through to unprefixed instruction
	// with 4 extra T-states. DD 00 = NOP (4+4=8 cycles, PC advances by 2).
	tests := []struct {
		name           string
		program        []uint8
		expectedPC     uint16
		expectedCycles uint64
		description    string
	}{
		{
			name:           "DD prefix with NOP",
			program:        []uint8{0xDD, 0x00},
			expectedPC:     2,
			expectedCycles: 8,
			description:    "DD prefix + NOP should execute NOP with 4 extra cycles",
		},
		{
			name:           "FD prefix with NOP",
			program:        []uint8{0xFD, 0x00},
			expectedPC:     2,
			expectedCycles: 8,
			description:    "FD prefix + NOP should execute NOP with 4 extra cycles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			memory := NewBasicMemory()
			cpu, err := New(memory)
			assert.NoError(t, err)

			for i, b := range tt.program {
				memory.Write(uint16(i), b)
			}

			cpu.PC = 0
			initialCycles := cpu.cycles

			err = cpu.Step()
			assert.NoError(t, err, tt.description)

			assert.Equal(t, tt.expectedPC, cpu.PC, "PC should advance correctly")
			assert.Equal(t, initialCycles+tt.expectedCycles, cpu.cycles, "cycles should match")
		})
	}
}

func TestUndocumentedNopInstructionDefinitions(t *testing.T) {
	// Test NopUndoc1 (DD prefix alone)
	assert.Equal(t, NopInst.Name, NopUndoc1.Name, "NopUndoc1 should have same name as regular NOP")
	assert.True(t, NopUndoc1.Unofficial, "NopUndoc1 should be marked as unofficial")

	opcodeInfo := NopUndoc1.Addressing[ImpliedAddressing]
	assert.Equal(t, PrefixDD, opcodeInfo.Opcode, "NopUndoc1 should use DD prefix")
	assert.Equal(t, byte(1), opcodeInfo.Size, "NopUndoc1 should be 1 byte")
	assert.Equal(t, byte(4), opcodeInfo.Cycles, "NopUndoc1 should take 4 cycles")

	// Test NopUndoc2 (FD prefix alone)
	assert.Equal(t, NopInst.Name, NopUndoc2.Name, "NopUndoc2 should have same name as regular NOP")
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

func TestDDFDPassthroughSequence(t *testing.T) {
	t.Parallel()

	memory := NewBasicMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// DD 00 = NOP with DD prefix (8 cycles, advances 2)
	memory.Write(0x0000, 0xDD)
	memory.Write(0x0001, 0x00)
	// FD 00 = NOP with FD prefix (8 cycles, advances 2)
	memory.Write(0x0002, 0xFD)
	memory.Write(0x0003, 0x00)
	// Regular NOP (4 cycles, advances 1)
	memory.Write(0x0004, 0x00)

	cpu.PC = 0
	initialCycles := cpu.cycles

	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(2), cpu.PC, "DD+NOP should advance PC by 2")
	assert.Equal(t, initialCycles+8, cpu.cycles, "DD+NOP should consume 8 cycles")

	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(4), cpu.PC, "FD+NOP should advance PC by 2")
	assert.Equal(t, initialCycles+16, cpu.cycles, "FD+NOP should consume another 8 cycles")

	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(5), cpu.PC, "Regular NOP should advance PC by 1")
	assert.Equal(t, initialCycles+20, cpu.cycles, "Regular NOP should consume another 4 cycles")
}
