package z80

import (
	"sync"

	"github.com/retroenv/retrogolib/arch"
)

// State represents CPU state for serialization and debugging.
type State struct {
	A uint8
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	AltA uint8
	AltB uint8
	AltC uint8
	AltD uint8
	AltE uint8
	AltH uint8
	AltL uint8

	IX uint16
	IY uint16

	SP uint16
	PC uint16
	I  uint8
	R  uint8

	Cycles     uint64
	Flags      Flags
	AltFlags   Flags
	Interrupts Interrupts

	Halted bool
}

// CPU represents a thread-safe Z80 microprocessor.
type CPU struct {
	mu sync.RWMutex

	A uint8
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	AltA uint8
	AltB uint8
	AltC uint8
	AltD uint8
	AltE uint8
	AltH uint8
	AltL uint8

	IX uint16
	IY uint16

	SP uint16
	PC uint16
	I  uint8
	R  uint8

	Flags    Flags
	AltFlags Flags

	cycles uint64
	halted bool

	// Interrupt handling
	iff1 bool  // interrupt flip-flop 1
	iff2 bool  // interrupt flip-flop 2
	im   uint8 // interrupt mode (0, 1, or 2)

	triggerIrq bool
	triggerNmi bool

	opts      Options
	TraceStep TraceStep // trace step info, set if tracing is enabled

	currentOpcode uint8 // opcode being executed (for instruction functions to access)

	memory Memory
}

// Interrupts holds the current interrupt state.
type Interrupts struct {
	IFF1         bool
	IFF2         bool
	IM           uint8
	NMITriggered bool
	IrqTriggered bool
}

const (
	initialCycles = 0
)

// New creates a new Z80 CPU with a memory controller.
// This allows different hardware implementations (Game Boy, MSX, ZX Spectrum, etc.)
// to provide their own memory mapping logic.
func New(memory Memory, options ...Option) (*CPU, error) {
	if memory == nil {
		return nil, ErrNilMemory
	}

	opts := NewOptions(options...)

	// Default to generic system
	if opts.initialPC == 0 && opts.initialSP == 0 && opts.systemType == "" {
		opts.systemType = arch.Generic
		opts.initialPC = 0x0000
		opts.initialSP = 0xFFFF
	}

	c := &CPU{
		PC:     opts.initialPC,
		SP:     opts.initialSP,
		cycles: initialCycles,
		opts:   opts,
		memory: memory,
		iff1:   false,
		iff2:   false,
		im:     0,
	}

	return c, nil
}

// Cycles returns total CPU cycles executed.
func (c *CPU) Cycles() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cycles
}

// Halted returns CPU halt state.
func (c *CPU) Halted() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.halted
}

// Halt stops CPU execution.
func (c *CPU) Halt() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.halted = true
}

// Resume continues CPU execution.
func (c *CPU) Resume() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.halted = false
}

// State returns complete CPU state.
func (c *CPU) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return State{
		A:        c.A,
		B:        c.B,
		C:        c.C,
		D:        c.D,
		E:        c.E,
		H:        c.H,
		L:        c.L,
		AltA:     c.AltA,
		AltB:     c.AltB,
		AltC:     c.AltC,
		AltD:     c.AltD,
		AltE:     c.AltE,
		AltH:     c.AltH,
		AltL:     c.AltL,
		IX:       c.IX,
		IY:       c.IY,
		SP:       c.SP,
		PC:       c.PC,
		I:        c.I,
		R:        c.R,
		Cycles:   c.cycles,
		Flags:    c.Flags,
		AltFlags: c.AltFlags,
		Interrupts: Interrupts{
			IFF1:         c.iff1,
			IFF2:         c.iff2,
			IM:           c.im,
			NMITriggered: c.triggerNmi,
			IrqTriggered: c.triggerIrq,
		},
		Halted: c.halted,
	}
}

// Memory returns the attached memory controller.
//
//nolint:ireturn // Returning interface is intentional for flexibility
func (c *CPU) Memory() Memory {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.memory
}

// BC returns the BC register pair as a 16-bit value.
func (c *CPU) BC() uint16 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bc()
}

// DE returns the DE register pair as a 16-bit value.
func (c *CPU) DE() uint16 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.de()
}

// HL returns the HL register pair as a 16-bit value.
func (c *CPU) HL() uint16 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hl()
}

// AF returns the AF register pair as a 16-bit value.
func (c *CPU) AF() uint16 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.af()
}

// TriggerNMI triggers a non-maskable interrupt.
func (c *CPU) TriggerNMI() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.triggerNmi = true
}

// TriggerIRQ triggers a maskable interrupt.
func (c *CPU) TriggerIRQ() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.triggerIrq = true
}

// bc returns the BC register pair as a 16-bit value (internal, no lock).
func (c *CPU) bc() uint16 {
	return uint16(c.B)<<8 | uint16(c.C)
}

// de returns the DE register pair as a 16-bit value (internal, no lock).
func (c *CPU) de() uint16 {
	return uint16(c.D)<<8 | uint16(c.E)
}

// hl returns the HL register pair as a 16-bit value (internal, no lock).
func (c *CPU) hl() uint16 {
	return uint16(c.H)<<8 | uint16(c.L)
}

// af returns the AF register pair as a 16-bit value (internal, no lock).
func (c *CPU) af() uint16 {
	return uint16(c.A)<<8 | uint16(c.GetFlags())
}

// setBC sets the BC register pair from a 16-bit value.
func (c *CPU) setBC(value uint16) {
	c.B = uint8(value >> 8)
	c.C = uint8(value)
}

// setDE sets the DE register pair from a 16-bit value.
func (c *CPU) setDE(value uint16) {
	c.D = uint8(value >> 8)
	c.E = uint8(value)
}

// setHL sets the HL register pair from a 16-bit value.
func (c *CPU) setHL(value uint16) {
	c.H = uint8(value >> 8)
	c.L = uint8(value)
}

// setAF sets the AF register pair from a 16-bit value.
func (c *CPU) setAF(value uint16) {
	c.A = uint8(value >> 8)
	c.setFlags(uint8(value))
}

// pop pops a byte from the stack and updates the stack pointer.
func (c *CPU) pop() uint8 {
	value := c.memory.Read(c.SP)
	c.SP++
	return value
}

// pop16 pops a word from the stack and updates the stack pointer.
func (c *CPU) pop16() uint16 {
	low := uint16(c.pop())
	high := uint16(c.pop())
	return high<<8 | low
}

// push pushes a byte to the stack and updates the stack pointer.
func (c *CPU) push(value uint8) {
	c.SP--
	c.memory.Write(c.SP, value)
}

// push16 pushes a word to the stack and updates the stack pointer.
func (c *CPU) push16(value uint16) {
	high := uint8(value >> 8)
	low := uint8(value)
	c.push(high)
	c.push(low)
}
