package cpu6809

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestStepServicesPendingNMI(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.S = 0x0200
	mem.WriteWord(VectorNMI, 0x9000)
	mem.data[0x9000] = 0x3B // RTI

	cpu.TriggerNMI()
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x9000), cpu.PC)
	assert.Equal(t, uint16(0x01F4), cpu.S)

	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8000), cpu.PC)
	assert.Equal(t, uint16(0x0200), cpu.S)
}

func TestMaskedIRQRemainsPending(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.S = 0x0200
	mem.WriteWord(VectorIRQ, 0x9000)
	mem.data[0x8000] = 0x12 // NOP

	cpu.TriggerIRQ()
	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8001), cpu.PC)

	cpu.Flags.I = 0
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x9000), cpu.PC)
}

func TestSYNCWakesForMaskedInterrupt(t *testing.T) {
	tests := []struct {
		name          string
		triggerBefore bool
	}{
		{name: "interrupt already pending", triggerBefore: true},
		{name: "interrupt arrives while waiting"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, mem := newTestCPU(t)
			mem.data[0x8000] = 0x13 // SYNC
			mem.data[0x8001] = 0x12 // NOP

			if tt.triggerBefore {
				cpu.TriggerIRQ()
			}
			err := cpu.Step()
			assert.NoError(t, err)
			assert.Equal(t, uint16(0x8001), cpu.PC)

			if !tt.triggerBefore {
				cpu.TriggerIRQ()
			}
			err = cpu.Step()
			assert.NoError(t, err)
			assert.Equal(t, uint16(0x8002), cpu.PC)
		})
	}
}

func TestCWAIDoesNotStackStateTwice(t *testing.T) {
	cpu, mem := newTestCPU(t)
	cpu.S = 0x0200
	mem.WriteWord(VectorIRQ, 0x9000)
	mem.data[0x8000] = 0x3C // CWAI
	mem.data[0x8001] = 0xEF // clear I, retain F
	mem.data[0x9000] = 0x3B // RTI

	err := cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8002), cpu.PC)
	assert.Equal(t, uint16(0x01F4), cpu.S)

	cpu.TriggerIRQ()
	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x9000), cpu.PC)
	assert.Equal(t, uint16(0x01F4), cpu.S)

	err = cpu.Step()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x8002), cpu.PC)
	assert.Equal(t, uint16(0x0200), cpu.S)
}

func TestInterruptEntryCycles(t *testing.T) {
	// Hardware interrupt entry used to consume zero cycles, freezing host timing.
	tests := []struct {
		name       string
		trigger    func(*CPU)
		vector     uint16
		cycles     uint64
		stackBytes uint16
	}{
		{name: "NMI", trigger: (*CPU).TriggerNMI, vector: VectorNMI, cycles: 19, stackBytes: 12},
		{name: "IRQ", trigger: (*CPU).TriggerIRQ, vector: VectorIRQ, cycles: 19, stackBytes: 12},
		{name: "FIRQ", trigger: (*CPU).TriggerFIRQ, vector: VectorFIRQ, cycles: 10, stackBytes: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, mem := newTestCPU(t)
			cpu.S = 0x200
			cpu.Flags.I, cpu.Flags.F = 0, 0
			mem.WriteWord(tt.vector, 0x9000)
			tt.trigger(cpu)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint16(0x9000), cpu.PC)
			assert.Equal(t, tt.cycles, cpu.Cycles())
			assert.Equal(t, uint16(0x200)-tt.stackBytes, cpu.S)
		})
	}
}
