package x86

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestOpcodes(t *testing.T) {
	assert.Equal(t, 256, len(Opcodes))

	// ADD r/m8, r8
	assert.Equal(t, "add", Opcodes[0x00].Instruction.Name)
	assert.Equal(t, ModRMRegisterAddressing, Opcodes[0x00].Addressing)

	// NOP
	assert.Equal(t, "nop", Opcodes[0x90].Instruction.Name)

	// MOV r/m8, r8
	assert.Equal(t, "mov", Opcodes[0x88].Instruction.Name)
}

func TestInstructions(t *testing.T) {
	assert.NotNil(t, Instructions["mov"])
	assert.NotNil(t, Instructions["add"])
	assert.NotNil(t, Instructions["nop"])
	assert.Nil(t, Instructions["invalid"])
}

func TestInstruction_GetOpcodeInfo(t *testing.T) {
	// Immediate addressing
	info, ok := AddALImm8.GetOpcodeInfo(ImmediateAddressing)
	assert.True(t, ok)
	assert.Equal(t, uint16(0x04), info.Opcode)

	// Implied addressing
	info, ok = Nop.GetOpcodeInfo(ImpliedAddressing)
	assert.True(t, ok)
	assert.Equal(t, uint16(0x90), info.Opcode)

	// Not supported
	_, ok = Nop.GetOpcodeInfo(DirectAddressing)
	assert.False(t, ok)
}

func TestInstruction_GetOpcodeByRegister(t *testing.T) {
	// 8-bit register
	info, ok := MovRegImm8.GetOpcodeByRegister(RegAL)
	assert.True(t, ok)
	assert.Equal(t, uint8(0xB0), info.Opcode)

	// 16-bit register
	info, ok = PushReg16.GetOpcodeByRegister(RegAX)
	assert.True(t, ok)
	assert.Equal(t, uint8(0x50), info.Opcode)

	// Invalid register
	_, ok = PushReg16.GetOpcodeByRegister(RegAL)
	assert.False(t, ok)
}

func TestInstruction_HasAddressing(t *testing.T) {
	assert.True(t, Nop.HasAddressing(ImpliedAddressing))
	assert.False(t, Nop.HasAddressing(DirectAddressing))
	assert.True(t, Jmp.HasAddressing(RelativeAddressing, ModRMRegisterAddressing))
}

func TestOpcodeInfo_TwoByte(t *testing.T) {
	// Single-byte
	info, _ := Nop.GetOpcodeInfo(ImpliedAddressing)
	assert.False(t, info.IsTwoByte())
	assert.Equal(t, uint8(0x90), info.PrimaryByte())

	// Two-byte (0x0F prefix)
	info, _ = Bsf.GetOpcodeInfo(ModRMRegisterAddressing)
	assert.True(t, info.IsTwoByte())
	assert.Equal(t, uint8(0x0F), info.PrimaryByte())
	assert.Equal(t, uint8(0xBC), info.SecondaryByte())
}

func TestModRM(t *testing.T) {
	// FromByte
	var m ModRM
	m.FromByte(0xC0) // mod=3, reg=0, rm=0
	assert.Equal(t, uint8(3), m.Mod)
	assert.Equal(t, uint8(0), m.Reg)
	assert.Equal(t, uint8(0), m.RM)

	// ToByte
	m2 := NewModRM(2, 5, 3)
	assert.Equal(t, uint8(0xAB), m2.ToByte()) // 10_101_011
}
