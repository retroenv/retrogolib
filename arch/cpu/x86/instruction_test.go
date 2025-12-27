package x86

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestInstruction_HasAddressing(t *testing.T) {
	tests := []struct {
		name     string
		inst     *Instruction
		modes    []AddressingMode
		expected bool
	}{
		{
			name:     "MovRegImm8 has immediate addressing",
			inst:     MovRegImm8,
			modes:    []AddressingMode{ImmediateAddressing},
			expected: true,
		},
		{
			name:     "MovRegImm8 does not have direct addressing",
			inst:     MovRegImm8,
			modes:    []AddressingMode{DirectAddressing},
			expected: false,
		},
		{
			name:     "IncReg16 has register addressing",
			inst:     IncReg16,
			modes:    []AddressingMode{RegisterAddressing},
			expected: true,
		},
		{
			name:     "Nop has implied addressing",
			inst:     Nop,
			modes:    []AddressingMode{ImpliedAddressing},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.inst.HasAddressing(tt.modes...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInstruction_GetOpcodeByRegister(t *testing.T) {
	tests := []struct {
		name     string
		inst     *Instruction
		register RegisterParam
		expected uint8
		hasInfo  bool
	}{
		{
			name:     "MovRegImm8 AL register",
			inst:     MovRegImm8,
			register: RegAL,
			expected: 0xB0,
			hasInfo:  true,
		},
		{
			name:     "MovRegImm8 BH register",
			inst:     MovRegImm8,
			register: RegBH,
			expected: 0xB7,
			hasInfo:  true,
		},
		{
			name:     "MovRegImm16 AX register",
			inst:     MovRegImm16,
			register: RegAX,
			expected: 0xB8,
			hasInfo:  true,
		},
		{
			name:     "IncReg16 CX register",
			inst:     IncReg16,
			register: RegCX,
			expected: 0x41,
			hasInfo:  true,
		},
		{
			name:     "IncReg16 invalid register",
			inst:     IncReg16,
			register: RegAL, // 8-bit register not valid for 16-bit inc
			expected: 0x00,
			hasInfo:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, hasInfo := tt.inst.GetOpcodeByRegister(tt.register)
			assert.Equal(t, tt.hasInfo, hasInfo)
			if hasInfo {
				assert.Equal(t, tt.expected, info.Opcode)
			}
		})
	}
}

func TestInstruction_GetOpcodeInfo(t *testing.T) {
	tests := []struct {
		name     string
		inst     *Instruction
		mode     AddressingMode
		expected uint8
		hasInfo  bool
	}{
		{
			name:     "MovRMReg8 ModR/M addressing",
			inst:     MovRMReg8,
			mode:     ModRMRegisterAddressing,
			expected: 0x88,
			hasInfo:  true,
		},
		{
			name:     "Jz relative addressing",
			inst:     Jz,
			mode:     RelativeAddressing,
			expected: 0x74,
			hasInfo:  true,
		},
		{
			name:     "Nop implied addressing",
			inst:     Nop,
			mode:     ImpliedAddressing,
			expected: 0x90,
			hasInfo:  true,
		},
		{
			name:     "Nop direct addressing (not supported)",
			inst:     Nop,
			mode:     DirectAddressing,
			expected: 0x00,
			hasInfo:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, hasInfo := tt.inst.GetOpcodeInfo(tt.mode)
			assert.Equal(t, tt.hasInfo, hasInfo)
			if hasInfo {
				assert.Equal(t, tt.expected, info.Opcode)
			}
		})
	}
}

func TestInstruction_SupportsRegister(t *testing.T) {
	tests := []struct {
		name     string
		inst     *Instruction
		register RegisterParam
		expected bool
	}{
		{
			name:     "PushReg16 supports AX",
			inst:     PushReg16,
			register: RegAX,
			expected: true,
		},
		{
			name:     "PushReg16 supports DI",
			inst:     PushReg16,
			register: RegDI,
			expected: true,
		},
		{
			name:     "PushReg16 does not support AL",
			inst:     PushReg16,
			register: RegAL,
			expected: false,
		},
		{
			name:     "MovRegImm8 supports BH",
			inst:     MovRegImm8,
			register: RegBH,
			expected: true,
		},
		{
			name:     "MovRegImm8 does not support AX",
			inst:     MovRegImm8,
			register: RegAX,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.inst.SupportsRegister(tt.register)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInstruction_GetSupportedAddressingModes(t *testing.T) {
	tests := []struct {
		name     string
		inst     *Instruction
		expected int // count of supported modes
	}{
		{
			name:     "MovRegImm8 has one addressing mode",
			inst:     MovRegImm8,
			expected: 1,
		},
		{
			name:     "MovMemImm8 has two addressing modes",
			inst:     MovMemImm8,
			expected: 2,
		},
		{
			name:     "Jmp has two addressing modes",
			inst:     Jmp,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modes := tt.inst.GetSupportedAddressingModes()
			assert.Equal(t, tt.expected, len(modes))
		})
	}
}

func TestInstruction_GetSupportedRegisters(t *testing.T) {
	tests := []struct {
		name     string
		inst     *Instruction
		expected int // count of supported registers
	}{
		{
			name:     "MovRegImm8 supports 8 registers",
			inst:     MovRegImm8,
			expected: 8,
		},
		{
			name:     "IncReg16 supports 8 registers",
			inst:     IncReg16,
			expected: 8,
		},
		{
			name:     "PushCS supports 1 register",
			inst:     PushCS,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registers := tt.inst.GetSupportedRegisters()
			assert.Equal(t, tt.expected, len(registers))
		})
	}
}

func TestInstruction_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		inst     *Instruction
		expected bool
	}{
		{
			name:     "MovRegImm8 is valid",
			inst:     MovRegImm8,
			expected: true,
		},
		{
			name:     "Nop is valid",
			inst:     Nop,
			expected: true,
		},
		{
			name:     "IncReg16 is valid",
			inst:     IncReg16,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.inst.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInstruction_GetAllRegisterVariants(t *testing.T) {
	t.Run("MovRegImm8 has all 8-bit register variants", func(t *testing.T) {
		variants := MovRegImm8.GetAllRegisterVariants()
		assert.Equal(t, 8, len(variants))

		// Check specific register mappings
		alInfo, exists := variants[RegAL]
		assert.True(t, exists)
		assert.Equal(t, uint16(0xB0), alInfo.Opcode)

		bhInfo, exists := variants[RegBH]
		assert.True(t, exists)
		assert.Equal(t, uint16(0xB7), bhInfo.Opcode)
	})

	t.Run("IncReg16 has all 16-bit register variants", func(t *testing.T) {
		variants := IncReg16.GetAllRegisterVariants()
		assert.Equal(t, 8, len(variants))

		// Check specific register mappings
		axInfo, exists := variants[RegAX]
		assert.True(t, exists)
		assert.Equal(t, uint16(0x40), axInfo.Opcode)

		diInfo, exists := variants[RegDI]
		assert.True(t, exists)
		assert.Equal(t, uint16(0x47), diInfo.Opcode)
	})
}

// Test specific instruction patterns
func TestArithmeticInstructions(t *testing.T) {
	t.Run("ADD instructions have correct opcodes", func(t *testing.T) {
		// Test ADD AL, imm8
		info, exists := AddALImm8.GetOpcodeInfo(ImmediateAddressing)
		assert.True(t, exists)
		assert.Equal(t, uint16(0x04), info.Opcode)
		assert.Equal(t, uint8(2), info.Size)
		assert.Equal(t, uint8(4), info.Cycles)
		assert.False(t, info.HasModRM)

		// Test ADD AX, imm16
		info, exists = AddAXImm16.GetOpcodeInfo(ImmediateAddressing)
		assert.True(t, exists)
		assert.Equal(t, uint16(0x05), info.Opcode)
		assert.Equal(t, uint8(3), info.Size)
		assert.Equal(t, uint8(4), info.Cycles)
		assert.False(t, info.HasModRM)

		// Test ADD r/m8, r8
		info, exists = AddRMReg8.GetOpcodeInfo(ModRMRegisterAddressing)
		assert.True(t, exists)
		assert.Equal(t, uint16(0x00), info.Opcode)
		assert.Equal(t, uint8(2), info.Size)
		assert.Equal(t, uint8(3), info.Cycles)
		assert.True(t, info.HasModRM)
	})
}

func TestJumpInstructions(t *testing.T) {
	t.Run("Conditional jump instructions have correct opcodes", func(t *testing.T) {
		tests := []struct {
			inst     *Instruction
			expected uint8
		}{
			{Jo, 0x70},  // JO
			{Jno, 0x71}, // JNO
			{Jz, 0x74},  // JZ/JE
			{Jnz, 0x75}, // JNZ/JNE
		}

		for _, tt := range tests {
			info, exists := tt.inst.GetOpcodeInfo(RelativeAddressing)
			assert.True(t, exists)
			assert.Equal(t, tt.expected, info.Opcode)
			assert.Equal(t, uint8(2), info.Size)
			assert.Equal(t, uint8(16), info.Cycles)
			assert.False(t, info.HasModRM)
		}
	})

	t.Run("Unconditional jump has correct opcodes", func(t *testing.T) {
		// Test JMP rel16
		info, exists := Jmp.GetOpcodeInfo(RelativeAddressing)
		assert.True(t, exists)
		assert.Equal(t, uint16(0xE9), info.Opcode)
		assert.Equal(t, uint8(3), info.Size)

		// Test JMP r/m16 (indirect)
		info, exists = Jmp.GetOpcodeInfo(ModRMRegisterAddressing)
		assert.True(t, exists)
		assert.Equal(t, uint16(0xFF), info.Opcode)
		assert.True(t, info.HasModRM)
	})
}

func TestStackInstructions(t *testing.T) {
	t.Run("PUSH register instructions have correct register mappings", func(t *testing.T) {
		expectedOpcodes := map[RegisterParam]uint8{
			RegAX: 0x50, RegCX: 0x51, RegDX: 0x52, RegBX: 0x53,
			RegSP: 0x54, RegBP: 0x55, RegSI: 0x56, RegDI: 0x57,
		}

		for reg, expectedOpcode := range expectedOpcodes {
			info, exists := PushReg16.GetOpcodeByRegister(reg)
			assert.True(t, exists)
			assert.Equal(t, expectedOpcode, info.Opcode)
			assert.Equal(t, uint8(1), info.Size)
			assert.Equal(t, uint8(11), info.Cycles)
			assert.False(t, info.HasModRM)
		}
	})

	t.Run("POP register instructions have correct register mappings", func(t *testing.T) {
		expectedOpcodes := map[RegisterParam]uint8{
			RegAX: 0x58, RegCX: 0x59, RegDX: 0x5A, RegBX: 0x5B,
			RegSP: 0x5C, RegBP: 0x5D, RegSI: 0x5E, RegDI: 0x5F,
		}

		for reg, expectedOpcode := range expectedOpcodes {
			info, exists := PopReg16.GetOpcodeByRegister(reg)
			assert.True(t, exists)
			assert.Equal(t, expectedOpcode, info.Opcode)
			assert.Equal(t, uint8(1), info.Size)
			assert.Equal(t, uint8(8), info.Cycles)
			assert.False(t, info.HasModRM)
		}
	})
}

func TestControlInstructions(t *testing.T) {
	t.Run("Control instructions have correct implied addressing", func(t *testing.T) {
		tests := []struct {
			inst     *Instruction
			expected uint8
			cycles   uint8
		}{
			{Nop, 0x90, 3},
			{Hlt, 0xF4, 2},
			{Clc, 0xF8, 2},
			{Stc, 0xF9, 2},
		}

		for _, tt := range tests {
			info, exists := tt.inst.GetOpcodeInfo(ImpliedAddressing)
			assert.True(t, exists)
			assert.Equal(t, tt.expected, info.Opcode)
			assert.Equal(t, uint8(1), info.Size)
			assert.Equal(t, tt.cycles, info.Cycles)
			assert.False(t, info.HasModRM)
		}
	})
}

func TestOpcodeInfo_ByteMethods(t *testing.T) {
	t.Run("Single-byte opcode (PUSHA)", func(t *testing.T) {
		info, exists := Pusha.GetOpcodeInfo(ImpliedAddressing)
		assert.True(t, exists)
		assert.Equal(t, uint16(0x60), info.Opcode)

		// Test helper methods
		assert.False(t, info.IsTwoByte())
		assert.Equal(t, uint8(0x60), info.PrimaryByte())
		assert.Equal(t, uint8(0), info.SecondaryByte())
	})

	t.Run("Two-byte opcode (BSF)", func(t *testing.T) {
		info, exists := Bsf.GetOpcodeInfo(ModRMRegisterAddressing)
		assert.True(t, exists)
		assert.Equal(t, uint16(0x0FBC), info.Opcode)

		// Test helper methods
		assert.True(t, info.IsTwoByte())
		assert.Equal(t, uint8(0x0F), info.PrimaryByte())   // Escape prefix
		assert.Equal(t, uint8(0xBC), info.SecondaryByte()) // Actual opcode
	})

	t.Run("Two-byte opcode (CMPXCHG)", func(t *testing.T) {
		info, exists := Cmpxchg.GetOpcodeInfo(ModRMRegisterAddressing)
		assert.True(t, exists)
		assert.Equal(t, uint16(0x0FB0), info.Opcode)

		// Test helper methods
		assert.True(t, info.IsTwoByte())
		assert.Equal(t, uint8(0x0F), info.PrimaryByte())
		assert.Equal(t, uint8(0xB0), info.SecondaryByte())
	})

	t.Run("Single-byte opcode boundary (0xFF)", func(t *testing.T) {
		// JMP r/m16 uses opcode 0xFF
		info, exists := Jmp.GetOpcodeInfo(ModRMRegisterAddressing)
		assert.True(t, exists)
		assert.Equal(t, uint16(0xFF), info.Opcode)

		// Should be single-byte
		assert.False(t, info.IsTwoByte())
		assert.Equal(t, uint8(0xFF), info.PrimaryByte())
		assert.Equal(t, uint8(0), info.SecondaryByte())
	})
}
