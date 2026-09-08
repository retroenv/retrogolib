package cpu65816

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestStepServicesInterrupt(t *testing.T) {
	// Step used to ignore pending interrupts, and emulation-mode entry kept PB.
	tests := []struct {
		name       string
		native     bool
		irq        bool
		vector     uint32
		cycles     uint64
		stackBytes uint16
	}{
		{name: "emulation NMI", vector: VectorEmuNMI, cycles: 7, stackBytes: 3},
		{name: "emulation IRQ", irq: true, vector: VectorEmuIRQ, cycles: 7, stackBytes: 3},
		{name: "native NMI", native: true, vector: VectorNativeNMI, cycles: 8, stackBytes: 4},
		{name: "native IRQ", native: true, irq: true, vector: VectorNativeIRQ, cycles: 8, stackBytes: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, mem := newTestCPU(t)
			cpu.E = !tt.native
			cpu.PB = 2
			cpu.Flags.I = 0
			cpu.Flags.D = 1
			mem.WriteWord(tt.vector, 0x9000)
			mem.Write(0x9000, 0x40) // RTI must execute on the next Step.
			cycles, sp := cpu.Cycles(), cpu.SP

			if tt.irq {
				cpu.TriggerIRQ()
			} else {
				cpu.TriggerNMI()
			}
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint16(0x9000), cpu.PC)
			assert.Equal(t, uint8(0), cpu.PB)
			assert.Equal(t, cycles+tt.cycles, cpu.Cycles())
			assert.Equal(t, sp-tt.stackBytes, cpu.SP)
			assert.Equal(t, uint8(1), cpu.Flags.I)
			assert.Equal(t, uint8(0), cpu.Flags.D)
			assert.Equal(t, uint16(0x8000), mem.ReadWord(uint32(cpu.SP)+2))
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint16(0x8000), cpu.PC)
			assert.Equal(t, sp, cpu.SP)
			assert.Equal(t, uint8(1), cpu.Flags.D)
			if tt.native {
				assert.Equal(t, uint8(2), cpu.PB)
			}
		})
	}
}

func TestWAIWakesForMaskedIRQ(t *testing.T) {
	// A request preceding WAI must also release the wait on the next Step.
	for _, pendingBefore := range []bool{false, true} {
		cpu, mem := newTestCPU(t)
		mem.Write(0x8000, 0xCB) // WAI.
		mem.Write(0x8001, 0x58) // CLI.
		mem.WriteWord(VectorEmuIRQ, 0x9000)
		if pendingBefore {
			cpu.TriggerIRQ()
		}
		assert.NoError(t, cpu.Step())
		assert.True(t, cpu.waiting)
		if !pendingBefore {
			cpu.TriggerIRQ()
		}
		assert.NoError(t, cpu.Step())
		assert.False(t, cpu.waiting)
		assert.Equal(t, uint16(0x8002), cpu.PC)
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint16(0x9000), cpu.PC)
	}
}

func TestSTPIgnoresInterruptsUntilReset(t *testing.T) {
	// Adding interrupt dispatch to Step must preserve STP's reset-only exit.
	cpu, mem := newTestCPU(t)
	mem.Write(0x8000, 0xDB) // STP.
	assert.NoError(t, cpu.Step())
	cycles := cpu.Cycles()
	cpu.TriggerNMI()
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x8001), cpu.PC)
	assert.Equal(t, cycles, cpu.Cycles())
	cpu.Reset()
	assert.False(t, cpu.stopped)
	assert.False(t, cpu.CheckInterrupts())
}
