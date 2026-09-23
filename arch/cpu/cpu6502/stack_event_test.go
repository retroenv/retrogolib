package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type stackEventAccess struct {
	address uint16
	value   byte
	write   bool
}

type stackEventMemory struct {
	accesses []stackEventAccess
	testMemory
}

func (m *stackEventMemory) Read(address uint16) byte {
	value := m.b[address]
	m.accesses = append(m.accesses, stackEventAccess{
		address: address,
		value:   value,
	})
	if address == 0x8000 {
		m.b[address] = 0 // Model a read-to-clear mapped register.
	}
	return value
}

func (m *stackEventMemory) Write(address uint16, value byte) {
	m.accesses = append(m.accesses, stackEventAccess{
		address: address,
		value:   value,
		write:   true,
	})
	m.b[address] = value
}

func TestStackEventPreservesInterruptBusAccesses(t *testing.T) {
	t.Parallel()

	for _, source := range []InterruptSource{InterruptNMI, InterruptIRQ} {
		var buses [2]*stackEventMemory
		var states [2]State
		for i := range buses {
			bus := &stackEventMemory{}
			bus.b[ResetAddress+1] = 0x80
			bus.b[0x8000] = 0x48
			memory, err := NewMemory(bus)
			assert.NoError(t, err)
			var hook StackEventHook
			if i == 1 {
				hook = func(StackEvent) {}
			}
			cpu := New(memory, WithStackEventHook(hook))
			bus.accesses = nil
			if source == InterruptNMI {
				triggerNMI(cpu)
			} else {
				triggerIRQ(cpu)
			}
			assert.True(t, cpu.CheckInterrupts())
			buses[i] = bus
			states[i] = cpu.State()
		}
		assert.Equal(t, buses[0].accesses, buses[1].accesses)
		assert.Equal(t, byte(0x48), buses[1].b[0x8000])
		assert.Equal(t, buses[0].b, buses[1].b)
		assert.Equal(t, states[0], states[1])
	}
}

func TestStackEventInstructionCycles(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		count  int
		cycles uint64
		name   string
		opcode byte
	}{
		{name: "PHA", opcode: 0x48, count: 1, cycles: 3},
		{name: "PLA", opcode: 0x68, count: 1, cycles: 4},
		{name: "JSR", opcode: 0x20, count: 2, cycles: 6},
		{name: "RTS", opcode: 0x60, count: 2, cycles: 6},
		{name: "RTI", opcode: 0x40, count: 3, cycles: 6},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var events []StackEvent
			cpu := newEventTestCPU(t, []byte{0xea, test.opcode, 0x00, 0x90}, func(event StackEvent) {
				events = append(events, event)
			})
			assert.NoError(t, cpu.Step()) // Advance beyond the initial cycle count.
			start := cpu.Cycles()
			assert.NoError(t, cpu.Step())
			assert.Len(t, events, test.count)
			for _, event := range events {
				assert.Equal(t, start, event.Cycle)
				assert.Equal(t, uint16(0x8001), event.PC)
				assert.Equal(t, test.opcode, event.Opcode)
				assert.Equal(t, InterruptNone, event.Interrupt)
			}
			assert.Equal(t, start+test.cycles, cpu.Cycles())
		})
	}
}

func TestStackEventReportsPushWrap(t *testing.T) {
	t.Parallel()

	events := make([]StackEvent, 0, 1)
	cpu := newEventTestCPU(t, []byte{0x48}, func(event StackEvent) {
		events = append(events, event)
	})
	cpu.SP = 0x00

	assert.NoError(t, cpu.Step())
	assert.Len(t, events, 1)
	assert.Equal(t, StackPush, events[0].Operation)
	assert.Equal(t, byte(0x00), events[0].Before)
	assert.Equal(t, byte(0xff), events[0].After)
	assert.True(t, events[0].Wrapped())
	assert.Equal(t, uint16(0x8000), events[0].PC)
	assert.Equal(t, byte(0x48), events[0].Opcode)
}

func TestStackEventReportsPullWrap(t *testing.T) {
	t.Parallel()

	events := make([]StackEvent, 0, 1)
	cpu := newEventTestCPU(t, []byte{0x68}, func(event StackEvent) {
		events = append(events, event)
	})
	cpu.SP = 0xff

	assert.NoError(t, cpu.Step())
	assert.Len(t, events, 1)
	assert.Equal(t, StackPull, events[0].Operation)
	assert.Equal(t, byte(0xff), events[0].Before)
	assert.Equal(t, byte(0x00), events[0].After)
	assert.True(t, events[0].Wrapped())
}

func TestStackEventReportsInterruptSources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		interrupt bool
		name      string
		prepare   func(*CPU)
		source    InterruptSource
	}{
		{name: "BRK", source: InterruptBRK},
		{name: "NMI", prepare: triggerNMI, source: InterruptNMI, interrupt: true},
		{name: "IRQ", prepare: triggerIRQ, source: InterruptIRQ, interrupt: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			events := make([]StackEvent, 0, 3)
			cpu := newEventTestCPU(t, []byte{0xea, 0x00, 0x00}, func(event StackEvent) {
				events = append(events, event)
			})
			assert.NoError(t, cpu.Step())
			start := cpu.Cycles()
			if test.prepare != nil {
				test.prepare(cpu)
			}

			if test.interrupt {
				assert.True(t, cpu.CheckInterrupts())
			} else {
				assert.NoError(t, cpu.Step())
			}
			assert.Len(t, events, 3)
			for i, event := range events {
				assert.Equal(t, test.source, event.Interrupt)
				assert.Equal(t, StackPush, event.Operation)
				assert.Equal(t, byte(InitialStack-i), event.Before)
				assert.Equal(t, byte(InitialStack-i-1), event.After)
				assert.Equal(t, uint16(0x8001), event.PC)
				assert.Equal(t, byte(0), event.Opcode)
				assert.Equal(t, start, event.Cycle)
			}
			assert.Equal(t, start+7, cpu.Cycles())
			cpu.memory.Write(cpu.PC, 0x48) // Handler instruction must get fresh context.
			assert.NoError(t, cpu.Step())
			assert.Len(t, events, 4)
			assert.Equal(t, InterruptNone, events[3].Interrupt)
			assert.Equal(t, byte(0x48), events[3].Opcode)
			assert.Equal(t, start+7, events[3].Cycle)
		})
	}
}

func TestStackEventExcludesTXSAndReset(t *testing.T) {
	t.Parallel()

	events := make([]StackEvent, 0, 1)
	cpu := newEventTestCPU(t, []byte{0x9a}, func(event StackEvent) {
		events = append(events, event)
	})
	cpu.X = 0x40

	assert.NoError(t, cpu.Step())
	assert.Len(t, events, 0)
	assert.Equal(t, byte(0x40), cpu.SP)
	cpu.Reset()
	assert.Len(t, events, 0)
}

func triggerNMI(cpu *CPU) {
	cpu.TriggerNMI()
}

func triggerIRQ(cpu *CPU) {
	cpu.Flags.I = 0
	cpu.TriggerIrq()
}

func newEventTestCPU(t *testing.T, program []byte, hook StackEventHook) *CPU {
	t.Helper()

	bus := &testMemory{}
	copy(bus.b[0x8000:], program)
	bus.b[ResetAddress] = 0x00
	bus.b[ResetAddress+1] = 0x80
	memory, err := NewMemory(bus)
	assert.NoError(t, err)

	return New(memory, WithVariant(VariantNES6502), WithStackEventHook(hook))
}
