package sm83

import (
	"sync"
)

// State represents complete CPU state for save/load and debugging.
type State struct {
	// 8-bit registers
	A uint8 // Accumulator
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	// 16-bit registers
	SP uint16 // Stack pointer
	PC uint16 // Program counter

	Cycles uint64
	Flags  Flags

	IME      bool // Interrupt Master Enable
	IMEDelay bool // IME becomes set after the next instruction (EI executed)
	HaltBug  bool // The next opcode fetch does not increment PC
	Halted   bool
	Stopped  bool
}

// CPU represents an SM83 microprocessor.
// Step and the accessor methods lock the CPU. The exported register fields and
// the register access methods do not lock, so callers must serialize them with Step.
type CPU struct {
	mu sync.RWMutex

	// 8-bit general purpose registers
	A uint8 // Accumulator
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	// 16-bit registers
	SP uint16 // Stack pointer
	PC uint16 // Program counter

	Flags Flags // Flag register

	cycles  uint64
	halted  bool
	stopped bool // STOP executed, only a pending joypad interrupt releases it

	// Interrupt control.
	ime      bool // Interrupt Master Enable
	imeDelay bool // IME is enabled after the instruction following EI
	haltBug  bool // HALT bug: PC fails to increment after HALT with IME=0 and pending interrupt

	opts options

	currentOpcode uint8 // opcode being executed

	memory Memory

	TraceStep TraceStep // trace step info populated when tracing is enabled
}

// New creates a new SM83 CPU with a memory controller.
func New(memory Memory, options ...Option) (*CPU, error) {
	if memory == nil {
		return nil, ErrNilMemory
	}

	opts := newOptions(options...)

	c := &CPU{
		PC:     opts.initialPC,
		SP:     opts.initialSP,
		opts:   opts,
		memory: memory,
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

// Stopped returns whether a STOP instruction waits for a joypad interrupt.
func (c *CPU) Stopped() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stopped
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
		SP:       c.SP,
		PC:       c.PC,
		Cycles:   c.cycles,
		Flags:    c.Flags,
		IME:      c.ime,
		IMEDelay: c.imeDelay,
		HaltBug:  c.haltBug,
		Halted:   c.halted,
		Stopped:  c.stopped,
	}
}

// SetState restores a complete CPU state from a snapshot.
func (c *CPU) SetState(s State) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.A, c.B, c.C, c.D, c.E, c.H, c.L = s.A, s.B, s.C, s.D, s.E, s.H, s.L
	c.SP = s.SP
	c.PC = s.PC
	c.cycles = s.Cycles
	c.Flags = s.Flags
	c.ime = s.IME
	c.imeDelay = s.IMEDelay
	c.haltBug = s.HaltBug
	c.halted = s.Halted
	c.stopped = s.Stopped
}

// Memory returns the attached memory controller.
//
//nolint:ireturn // intentional: Memory is the public API interface
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

// GetRegisterValue returns the value of a register by its 3-bit encoding.
// Encoding: B=0, C=1, D=2, E=3, H=4, L=5, (HL)=6, A=7.
func (c *CPU) GetRegisterValue(reg uint8) uint8 {
	switch reg {
	case 0:
		return c.B
	case 1:
		return c.C
	case 2:
		return c.D
	case 3:
		return c.E
	case 4:
		return c.H
	case 5:
		return c.L
	case 6:
		return c.memory.Read(c.hl())
	case 7:
		return c.A
	}
	return 0
}

// SetRegisterValue sets the value of a register by its 3-bit encoding.
func (c *CPU) SetRegisterValue(reg, value uint8) {
	switch reg {
	case 0:
		c.B = value
	case 1:
		c.C = value
	case 2:
		c.D = value
	case 3:
		c.E = value
	case 4:
		c.H = value
	case 5:
		c.L = value
	case 6:
		c.memory.Write(c.hl(), value)
	case 7:
		c.A = value
	}
}

// bc returns the BC register pair (internal, no lock).
func (c *CPU) bc() uint16 {
	return uint16(c.B)<<8 | uint16(c.C)
}

// de returns the DE register pair (internal, no lock).
func (c *CPU) de() uint16 {
	return uint16(c.D)<<8 | uint16(c.E)
}

// hl returns the HL register pair (internal, no lock).
func (c *CPU) hl() uint16 {
	return uint16(c.H)<<8 | uint16(c.L)
}

// af returns the AF register pair (internal, no lock).
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

// pop pops a byte from the stack.
func (c *CPU) pop() uint8 {
	value := c.memory.Read(c.SP)
	c.SP++
	return value
}

// pop16 pops a word from the stack.
func (c *CPU) pop16() uint16 {
	low := uint16(c.pop())
	high := uint16(c.pop())
	return high<<8 | low
}

// push pushes a byte to the stack.
func (c *CPU) push(value uint8) {
	c.SP--
	c.memory.Write(c.SP, value)
}

// push16 pushes a word to the stack.
func (c *CPU) push16(value uint16) {
	high := uint8(value >> 8)
	low := uint8(value)
	c.push(high)
	c.push(low)
}
