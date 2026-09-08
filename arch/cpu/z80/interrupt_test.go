package z80

import (
	"fmt"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type stackMappedIRQBus struct {
	*testBus
	address uint16
}

func (bus *stackMappedIRQBus) IRQData() uint8 {
	return bus.Read(bus.address)
}

func TestInterruptEntryPoints(t *testing.T) {
	// CheckInterrupts used a different vector source, while Step also ran the ISR.
	for _, direct := range []bool{false, true} {
		for _, tt := range []struct {
			name         string
			mode         InterruptMode
			data         uint8
			vector       uint16
			cycles       uint64
			acknowledges int
		}{
			{name: "IM0", mode: InterruptMode0, data: 0xCF, vector: 8, cycles: 13, acknowledges: 1},
			{name: "IM0 fallback", mode: InterruptMode0, data: 0, vector: 0x38, cycles: 13, acknowledges: 1},
			{name: "IM1", mode: InterruptMode1, vector: 0x38, cycles: 13},
			{name: "IM2", mode: InterruptMode2, data: 0xA5, vector: 0x4567, cycles: 19, acknowledges: 1},
		} {
			t.Run(fmt.Sprintf("%s/direct=%t", tt.name, direct), func(t *testing.T) {
				bus := &testBus{Memory: NewBasicMemory(), irqData: tt.data}
				cpu, err := NewWithBus(bus, WithInitialPC(0x1234), WithInitialSP(0x8000))
				assert.NoError(t, err)
				assert.NoError(t, cpu.SetInterruptMode(tt.mode))
				cpu.I, cpu.R, cpu.A = 0xAB, 0xFF, 0x42
				bus.WriteWord(0xABA5, tt.vector)
				bus.Write(0xFFFF, 0x10)    // Must not supply IM2's vector byte.
				bus.Write(tt.vector, 0x3C) // INC A must wait for the next Step.
				cpu.EnableInterrupts()
				cpu.TriggerIRQ()
				if direct {
					assert.True(t, cpu.CheckInterrupts())
				} else {
					assert.NoError(t, cpu.Step())
				}
				assert.Equal(t, tt.vector, cpu.PC)
				assert.Equal(t, tt.vector, cpu.MEMPTR)
				assert.Equal(t, uint16(0x7FFE), cpu.SP)
				assert.Equal(t, uint16(0x1234), bus.ReadWord(cpu.SP))
				assert.Equal(t, tt.cycles, cpu.Cycles())
				assert.Equal(t, uint8(0x80), cpu.R)
				assert.Equal(t, uint8(0x42), cpu.A)
				assert.Equal(t, tt.acknowledges, bus.irqCalls)
				assert.False(t, cpu.iff1)
				assert.False(t, cpu.iff2)
				assert.False(t, cpu.CheckInterrupts())
			})
		}
	}
}

func TestIRQDataSampledBeforeStackWrites(t *testing.T) {
	// Stack writes previously could change a mapped device before IRQData was sampled.
	for _, mode := range []InterruptMode{InterruptMode0, InterruptMode2} {
		bus := &stackMappedIRQBus{testBus: &testBus{Memory: NewBasicMemory()}, address: 0x7FFF}
		cpu, err := NewWithBus(bus, WithInitialPC(0x1234), WithInitialSP(0x8000))
		assert.NoError(t, err)
		assert.NoError(t, cpu.SetInterruptMode(mode))
		cpu.I = 0xAB
		bus.Write(bus.address, 0xCF)
		bus.WriteWord(0xABCF, 8)
		cpu.EnableInterrupts()
		cpu.TriggerIRQ()
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint16(8), cpu.PC)
		assert.Equal(t, uint8(0x12), bus.Read(bus.address))
	}
}

func TestIM0RestartVectors(t *testing.T) {
	for vector := uint8(0); vector < 0x40; vector += 8 {
		bus := &testBus{Memory: NewBasicMemory(), irqData: 0xC7 | vector}
		cpu, err := NewWithBus(bus, WithInitialPC(0x1234))
		assert.NoError(t, err)
		cpu.EnableInterrupts()
		cpu.TriggerIRQ()
		assert.True(t, cpu.CheckInterrupts())
		assert.Equal(t, uint16(vector), cpu.PC)
		assert.Equal(t, uint64(13), cpu.Cycles())
	}
}

func TestIM2VectorWrapsMemory(t *testing.T) {
	bus := &testBus{Memory: NewBasicMemory(), irqData: 0xFF}
	cpu, err := NewWithBus(bus, WithInitialSP(0x8000))
	assert.NoError(t, err)
	cpu.I = 0xFF
	bus.WriteWord(0xFFFF, 0x3456)
	assert.NoError(t, cpu.SetInterruptMode(InterruptMode2))
	cpu.EnableInterrupts()
	cpu.TriggerIRQ()
	assert.True(t, cpu.CheckInterrupts())
	assert.Equal(t, uint16(0x3456), cpu.PC)
}

func TestHaltWakesForInterrupt(t *testing.T) {
	// HALT previously returned before checking interrupts and never resumed.
	cpu, err := New(NewBasicMemory(), WithInitialPC(0x1000))
	assert.NoError(t, err)
	cpu.bus.Write(cpu.PC, 0x76)
	assert.NoError(t, cpu.Step())
	cpu.TriggerIRQ()
	assert.NoError(t, cpu.Step())
	assert.True(t, cpu.Halted())
	assert.Equal(t, uint8(2), cpu.R)
	cpu.EnableInterrupts()
	assert.NoError(t, cpu.Step())
	assert.False(t, cpu.Halted())
	assert.Equal(t, uint16(0x38), cpu.PC)
	assert.Equal(t, uint16(0x1001), cpu.bus.ReadWord(cpu.SP))
	assert.Equal(t, uint64(21), cpu.Cycles())
}

func TestNMIKeepsIFF2AndPendingIRQ(t *testing.T) {
	// Nested NMIs must not overwrite the original interrupt-enable state.
	cpu, err := New(NewBasicMemory(), WithInitialPC(0x1000))
	assert.NoError(t, err)
	cpu.EnableInterrupts()
	cpu.Halt()
	cpu.TriggerIRQ()
	cpu.TriggerNMI()
	assert.NoError(t, cpu.Step())
	assert.False(t, cpu.Halted())
	assert.Equal(t, uint16(0x66), cpu.PC)
	assert.Equal(t, uint16(0x66), cpu.MEMPTR)
	assert.True(t, cpu.iff2)
	assert.False(t, cpu.iff1)
	assert.True(t, cpu.triggerIrq)
	assert.Equal(t, uint64(11), cpu.Cycles())
	cpu.TriggerNMI()
	assert.True(t, cpu.CheckInterrupts())
	assert.True(t, cpu.iff2)
	assert.Equal(t, uint64(22), cpu.Cycles())
}

func TestEIDelaysIRQByOneInstruction(t *testing.T) {
	// A pending IRQ must allow the instruction after EI to finish first.
	for _, next := range []byte{0, 0xF3, 0x76, 0xFB} { // NOP, DI, HALT, EI.
		t.Run(fmt.Sprintf("%02x", next), func(t *testing.T) {
			cpu, err := New(NewBasicMemory(), WithInitialPC(0x1000))
			assert.NoError(t, err)
			cpu.bus.Write(0x1000, 0xFB)
			cpu.bus.Write(0x1001, next)
			cpu.TriggerIRQ()
			assert.NoError(t, cpu.Step())
			assert.False(t, cpu.CheckInterrupts())
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint16(0x1002), cpu.PC)
			assert.NoError(t, cpu.Step())
			if next == 0xF3 || next == 0xFB {
				assert.Equal(t, uint16(0x1003), cpu.PC)
			} else {
				assert.Equal(t, uint16(0x38), cpu.PC)
			}
		})
	}
}

func TestLdAIRInterruptParity(t *testing.T) {
	for _, opcode := range []byte{0x57, 0x5F} {
		for _, direct := range []bool{false, true} {
			cpu, err := New(NewBasicMemory(), WithInitialPC(0x1000))
			assert.NoError(t, err)
			cpu.bus.Write(0x1000, PrefixED)
			cpu.bus.Write(0x1001, opcode)
			cpu.EnableInterrupts()
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint8(1), cpu.Flags.P)
			cpu.TriggerIRQ()
			if direct {
				assert.True(t, cpu.CheckInterrupts())
			} else {
				assert.NoError(t, cpu.Step())
			}
			assert.Equal(t, uint8(0), cpu.Flags.P)
		}
	}
}

func TestLdAIRParitySurvivesLaterInterrupt(t *testing.T) {
	cpu, err := New(NewBasicMemory(), WithInitialPC(0x1000))
	assert.NoError(t, err)
	cpu.bus.Write(0x1000, PrefixED)
	cpu.bus.Write(0x1001, 0x57)
	cpu.EnableInterrupts()
	assert.NoError(t, cpu.Step())
	assert.NoError(t, cpu.Step()) // Intervening NOP ends the quirk's window.
	cpu.TriggerIRQ()
	assert.True(t, cpu.CheckInterrupts())
	assert.Equal(t, uint8(1), cpu.Flags.P)
}
