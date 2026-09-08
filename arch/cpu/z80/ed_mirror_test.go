package z80

import (
	"fmt"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestEDNegMirrors(t *testing.T) {
	for _, opcode := range []byte{0x44, 0x4C, 0x54, 0x5C, 0x64, 0x6C, 0x74, 0x7C} {
		t.Run(fmt.Sprintf("%02x", opcode), func(t *testing.T) {
			cpu, err := New(NewBasicMemory())
			assert.NoError(t, err)
			cpu.A = 0x80
			cpu.bus.Write(0, PrefixED)
			cpu.bus.Write(1, opcode)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint8(0x80), cpu.A)
			assert.Equal(t, uint8(0x87), cpu.GetFlags())
			assert.Equal(t, uint16(2), cpu.PC)
			assert.Equal(t, uint8(2), cpu.R)
			assert.Equal(t, uint64(8), cpu.Cycles())
		})
	}
}

func TestEDInterruptModeMirrors(t *testing.T) {
	for _, tt := range []struct {
		opcode byte
		mode   InterruptMode
	}{
		{0x46, InterruptMode0}, {0x4E, InterruptMode0}, {0x66, InterruptMode0}, {0x6E, InterruptMode0},
		{0x56, InterruptMode1}, {0x76, InterruptMode1},
		{0x5E, InterruptMode2}, {0x7E, InterruptMode2},
	} {
		t.Run(fmt.Sprintf("%02x", tt.opcode), func(t *testing.T) {
			cpu, err := New(NewBasicMemory())
			assert.NoError(t, err)
			assert.NoError(t, cpu.SetInterruptMode((tt.mode+1)%3))
			cpu.setFlags(0xA5)
			cpu.bus.Write(0, PrefixED)
			cpu.bus.Write(1, tt.opcode)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, tt.mode, cpu.GetInterruptMode())
			assert.Equal(t, uint8(0xA5), cpu.GetFlags())
			assert.Equal(t, uint16(2), cpu.PC)
			assert.Equal(t, uint64(8), cpu.Cycles())
		})
	}
}

func TestEDReturnMirrorsAndNotification(t *testing.T) {
	for _, opcode := range []byte{0x45, 0x4D, 0x55, 0x5D, 0x65, 0x6D, 0x75, 0x7D} {
		t.Run(fmt.Sprintf("%02x", opcode), func(t *testing.T) {
			bus := &testBus{Memory: NewBasicMemory()}
			cpu, err := NewWithBus(bus, WithInitialSP(0x8000))
			assert.NoError(t, err)
			cpu.iff2 = true
			bus.Write(0, PrefixED)
			bus.Write(1, opcode)
			// A return target equal to the opcode PC must not auto-advance.
			bus.WriteWord(0x8000, 0)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint16(0), cpu.PC)
			assert.Equal(t, uint16(0), cpu.MEMPTR)
			assert.Equal(t, uint16(0x8002), cpu.SP)
			assert.True(t, cpu.iff1)
			assert.True(t, cpu.iff2)
			assert.Equal(t, uint64(14), cpu.Cycles())
			if opcode == 0x4D {
				assert.Equal(t, 1, bus.retiCalls)
			} else {
				assert.Equal(t, 0, bus.retiCalls)
			}
		})
	}
}

func TestFullBusPortAddresses(t *testing.T) {
	// Block output must expose B after decrement; inputs use its original value.
	for _, tt := range []struct {
		name  string
		code  []byte
		port  uint16
		value uint8
		read  bool
	}{
		{name: "IN immediate", code: []byte{0xDB, 0x34}, port: 0x5634, value: 0xA5, read: true},
		{name: "OUT immediate", code: []byte{0xD3, 0x34}, port: 0x5634, value: 0x56},
		{name: "IN B C", code: []byte{PrefixED, 0x40}, port: 0x1234, value: 0xA5, read: true},
		{name: "OUT C B", code: []byte{PrefixED, 0x41}, port: 0x1234, value: 0x12},
		{name: "IN flags", code: []byte{PrefixED, 0x70}, port: 0x1234, value: 0xA5, read: true},
		{name: "OUT zero", code: []byte{PrefixED, 0x71}, port: 0x1234},
		{name: "INI", code: []byte{PrefixED, 0xA2}, port: 0x1234, value: 0xA5, read: true},
		{name: "IND", code: []byte{PrefixED, 0xAA}, port: 0x1234, value: 0xA5, read: true},
		{name: "INIR", code: []byte{PrefixED, 0xB2}, port: 0x1234, value: 0xA5, read: true},
		{name: "INDR", code: []byte{PrefixED, 0xBA}, port: 0x1234, value: 0xA5, read: true},
		{name: "OUTI", code: []byte{PrefixED, 0xA3}, port: 0x1134, value: 0xA5},
		{name: "OUTD", code: []byte{PrefixED, 0xAB}, port: 0x1134, value: 0xA5},
		{name: "OTIR", code: []byte{PrefixED, 0xB3}, port: 0x1134, value: 0xA5},
		{name: "OTDR", code: []byte{PrefixED, 0xBB}, port: 0x1134, value: 0xA5},
	} {
		t.Run(tt.name, func(t *testing.T) {
			bus := &testBus{Memory: NewBasicMemory(), portValue: 0xA5}
			cpu, err := NewWithBus(bus)
			assert.NoError(t, err)
			cpu.A, cpu.B, cpu.C, cpu.H = 0x56, 0x12, 0x34, 0x80
			for i, value := range tt.code {
				bus.Write(uint16(i), value)
			}
			bus.Write(0x8000, 0xA5)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, []singleStepPort{{Address: tt.port, Value: tt.value, IsRead: tt.read}}, bus.ports)
		})
	}
}

func TestIgnoredPrefixResetsQ(t *testing.T) {
	// An ignored index prefix clears Q before SCF/CCF derive their X/Y flags.
	for _, prefix := range []byte{PrefixDD, PrefixFD} {
		for _, opcode := range []byte{0x37, 0x3F} {
			cpu, err := New(NewBasicMemory())
			assert.NoError(t, err)
			cpu.setFlags(0x28)
			cpu.q = 0x28
			cpu.bus.Write(0, prefix)
			cpu.bus.Write(1, opcode)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint8(0x29), cpu.GetFlags())
			assert.Equal(t, uint16(2), cpu.PC)
		}
	}
}
