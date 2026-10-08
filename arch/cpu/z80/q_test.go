package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestQResetsAfterInstructionWithoutFlagWrite(t *testing.T) {
	// Q latched F after every instruction, so SCF after NOP lost the X/Y bits.
	cpu, err := New(NewBasicMemory())
	assert.NoError(t, err)
	cpu.setFlags(0x28)
	cpu.q = 0x28
	cpu.bus.Write(0, 0x00) // NOP
	cpu.bus.Write(1, 0x37) // SCF
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint8(0), cpu.q)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint8(0x29), cpu.GetFlags())
	assert.Equal(t, uint8(0x29), cpu.q)
}

func TestQFollowsFlagWritesWithoutValueChange(t *testing.T) {
	// ALU flag writes set Q even when F keeps its value. Loading F directly with
	// POP AF or EX AF,AF' is not a flag write and resets Q, as the corpus records.
	for _, tt := range []struct {
		name string
		code []byte
		q    uint8
	}{
		{name: "AND A", code: []byte{0xA7}, q: 0x54},
		{name: "EX AF,AF'", code: []byte{0x08}},
		{name: "POP AF", code: []byte{0xF1}},
		{name: "LD A,B", code: []byte{0x78}},
		{name: "EXX", code: []byte{0xD9}},
		{name: "DD NOP", code: []byte{PrefixDD, 0x00}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cpu, err := New(NewBasicMemory(), WithInitialSP(0x8000))
			assert.NoError(t, err)
			cpu.setFlags(0x54)
			cpu.AltFlags = cpu.Flags
			cpu.bus.WriteWord(0x8000, 0x0054)
			cpu.q = 0xFF
			for i, value := range tt.code {
				cpu.bus.Write(uint16(i), value)
			}
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint8(0x54), cpu.GetFlags())
			assert.Equal(t, tt.q, cpu.q)
		})
	}
}

func TestQResetsDuringHaltIdle(t *testing.T) {
	cpu, err := New(NewBasicMemory())
	assert.NoError(t, err)
	cpu.bus.Write(0, 0xA7) // AND A
	cpu.bus.Write(1, 0x76) // HALT
	assert.NoError(t, cpu.Step())
	assert.NotEqual(t, uint8(0), cpu.q)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint8(0), cpu.q)
	cpu.q = 0x44
	assert.NoError(t, cpu.Step()) // Idle cycle.
	assert.Equal(t, uint8(0), cpu.q)
}
