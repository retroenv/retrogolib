package cpu65816

import (
	"errors"
	"sync"
)

// State represents a complete snapshot of the 65816 CPU state.
type State struct {
	C  uint16 // Accumulator (full 16-bit C register)
	X  uint16 // Index X
	Y  uint16 // Index Y
	SP uint16 // Stack pointer
	DP uint16 // Direct Page register
	DB uint8  // Data Bank register
	PB uint8  // Program Bank register
	PC uint16 // Program counter (within bank)
	P  uint8  // Processor status (from Flags.Get())
	E  bool   // Emulation flag

	Cycles uint64
}

// CPU represents a thread-safe WDC 65C816 microprocessor.
type CPU struct {
	mu sync.RWMutex

	// Registers
	C  uint16 // Accumulator (full 16-bit; A = low byte, B = high byte)
	X  uint16 // Index register X
	Y  uint16 // Index register Y
	SP uint16 // Stack pointer
	DP uint16 // Direct Page register
	DB uint8  // Data Bank register
	PB uint8  // Program Bank register
	PC uint16 // Program counter (within current program bank)

	Flags Flags // Processor status register (P)
	E     bool  // Emulation flag (toggled by XCE)

	cycles    uint64
	stopped   bool // STP instruction state
	waiting   bool // WAI instruction state
	pcChanged bool // set by instructions that explicitly set PC (branches, jumps)

	// Interrupt control
	triggerNMI bool
	triggerIRQ bool
	nmiRunning bool
	irqRunning bool

	memory *Memory
	opts   options

	TraceStep TraceStep // set when tracing is enabled
}

// TraceStep holds information for instruction tracing.
type TraceStep struct {
	PC             uint16
	PB             uint8
	OpcodeOperands []byte
	Opcode         Opcode
	PageCrossed    bool
}

const (
	initialCycles = 7
	// Initial SP is $01FF in emulation mode (stack page 1).
	initialSP = 0x01FF
)

// New creates a new 65816 CPU, reads the reset vector, and initializes registers.
func New(memory *Memory, opts ...Option) (*CPU, error) {
	if memory == nil {
		return nil, errors.New("memory cannot be nil")
	}

	c := &CPU{
		SP:     initialSP,
		cycles: initialCycles,
		memory: memory,
		opts:   newOptions(opts...),
		E:      true, // Start in emulation mode
	}

	// Force M=1, X=1 in emulation mode
	c.Flags.M = 1
	c.Flags.X = 1
	c.Flags.I = 1

	// Read reset vector (emulation mode vector at $FFFC)
	resetVec := memory.ReadVector(VectorEmuRESET)
	c.PC = resetVec

	return c, nil
}

// A returns the low byte of the accumulator (8-bit mode or low half of 16-bit).
func (c *CPU) A() uint8 { return uint8(c.C) }

// B returns the high byte of the accumulator.
func (c *CPU) B() uint8 { return uint8(c.C >> 8) }

// FullPC returns the 24-bit effective program address (PB:PC).
func (c *CPU) FullPC() uint32 {
	return bank24(c.PB, c.PC)
}

// AccWidth returns the current accumulator width in bytes (1 or 2).
func (c *CPU) AccWidth() int {
	if c.E || c.Flags.M != 0 {
		return 1
	}
	return 2
}

// IdxWidth returns the current index register width in bytes (1 or 2).
func (c *CPU) IdxWidth() int {
	if c.E || c.Flags.X != 0 {
		return 1
	}
	return 2
}

// Cycles returns the total number of cycles executed.
func (c *CPU) Cycles() uint64 { return c.cycles }

// State returns a snapshot of the current CPU state.
func (c *CPU) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return State{
		C:      c.C,
		X:      c.X,
		Y:      c.Y,
		SP:     c.SP,
		DP:     c.DP,
		DB:     c.DB,
		PB:     c.PB,
		PC:     c.PC,
		P:      c.Flags.Get(),
		E:      c.E,
		Cycles: c.cycles,
	}
}

// ValidateState checks that CPU state is consistent.
func (c *CPU) ValidateState() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.memory == nil {
		return errors.New("CPU memory is nil")
	}
	return nil
}

// Reset resets the CPU to its initial post-reset state.
func (c *CPU) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.C = 0
	c.X = 0
	c.Y = 0
	c.SP = initialSP
	c.DP = 0
	c.DB = 0
	c.PB = 0
	c.E = true
	c.Flags = Flags{M: 1, X: 1, I: 1}
	c.cycles = initialCycles
	c.stopped = false
	c.waiting = false
	c.triggerNMI = false
	c.triggerIRQ = false
	c.nmiRunning = false
	c.irqRunning = false

	if c.memory != nil {
		c.PC = c.memory.ReadVector(VectorEmuRESET)
	}
}

// Memory returns the CPU's memory.
func (c *CPU) Memory() *Memory { return c.memory }

// GetP returns the current processor status byte.
func (c *CPU) GetP() uint8 {
	return c.Flags.Get()
}

// SetP sets the processor status from a byte, handling E-mode constraints.
func (c *CPU) SetP(p uint8) {
	c.Flags.Set(p)
	if c.E {
		// Emulation mode forces M=1, X=1
		c.Flags.M = 1
		c.Flags.X = 1
	} else if c.Flags.X != 0 {
		// X flag transition to 8-bit: zero high bytes of X and Y
		c.X &= 0x00FF
		c.Y &= 0x00FF
	}
}
