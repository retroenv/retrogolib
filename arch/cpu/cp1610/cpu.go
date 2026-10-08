package cp1610

import "sync"

// Register numbers used for program flow.
const (
	StackPointer   = 6
	ProgramCounter = 7
)

// State is a snapshot of the CPU registers and status.
type State struct {
	R                 [8]uint16
	Flags             Flags
	Cycles            uint64
	Halted            bool
	IRQPending        bool
	LastInterruptible bool
}

// CPU executes CP1610 instructions against a caller-supplied memory bus.
// TriggerIRQ can run concurrently with Step. Callers must serialize other access.
type CPU struct {
	irqMu sync.Mutex

	R     [8]uint16
	Flags Flags

	memory            *Memory
	cycles            uint64
	halted            bool
	irqPending        bool
	lastInterruptible bool
}

// New creates a CPU in its Intellivision reset state.
func New(memory *Memory) (*CPU, error) {
	if memory == nil || memory.BasicMemory == nil {
		return nil, ErrNilMemory
	}
	c := &CPU{memory: memory}
	c.Reset()
	return c, nil
}

// Cycles returns the number of executed CP1610 clock cycles.
func (c *CPU) Cycles() uint64 { return c.cycles }

// Halted reports whether HLT stopped the CPU.
func (c *CPU) Halted() bool { return c.halted }

// Memory returns the CPU memory bus.
func (c *CPU) Memory() *Memory { return c.memory }

// State returns the current CPU state.
func (c *CPU) State() State {
	c.irqMu.Lock()
	defer c.irqMu.Unlock()
	return State{
		R:                 c.R,
		Flags:             c.Flags,
		Cycles:            c.cycles,
		Halted:            c.halted,
		IRQPending:        c.irqPending,
		LastInterruptible: c.lastInterruptible,
	}
}

// ValidateState checks whether the CPU has a memory bus.
func (c *CPU) ValidateState() error {
	if c.memory == nil || c.memory.BasicMemory == nil {
		return ErrNilMemory
	}
	return nil
}

// Reset starts execution at the Intellivision Executive ROM entry address.
func (c *CPU) Reset() {
	c.R = [8]uint16{}
	c.R[ProgramCounter] = ResetAddress
	c.Flags = Flags{}
	c.cycles = 0
	c.halted = false
	c.lastInterruptible = false
	c.irqMu.Lock()
	c.irqPending = false
	c.irqMu.Unlock()
}
