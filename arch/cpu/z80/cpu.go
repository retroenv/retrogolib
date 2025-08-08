package z80

import (
	"sync"
)

// State contains the current state of the CPU.
type State struct {
	// Main registers
	A uint8
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	// Alternate registers
	A_ uint8
	B_ uint8
	C_ uint8
	D_ uint8
	E_ uint8
	H_ uint8
	L_ uint8

	// Index registers
	IX uint16
	IY uint16

	// Special registers
	SP uint16 // stack pointer
	PC uint16 // program counter
	I  uint8  // interrupt vector
	R  uint8  // refresh register

	Cycles     uint64
	Flags      Flags
	Flags_     Flags // alternate flags
	Interrupts Interrupts

	Halted bool
}

type CPU struct {
	mu sync.RWMutex

	// Main registers
	A uint8 // accumulator
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	// Alternate registers (shadow registers)
	A_ uint8
	B_ uint8
	C_ uint8
	D_ uint8
	E_ uint8
	H_ uint8
	L_ uint8

	// Index registers
	IX uint16
	IY uint16

	// Special registers
	SP uint16 // stack pointer
	PC uint16 // program counter
	I  uint8  // interrupt vector
	R  uint8  // refresh register

	Flags  Flags
	Flags_ Flags // alternate flags

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

	memory *Memory
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
	InitialStack  = 0xFFFE
)

// New creates a new Z80 CPU.
func New(memory *Memory, options ...Option) *CPU {
	opts := NewOptions(options...)
	c := &CPU{
		SP:     InitialStack,
		cycles: initialCycles,
		opts:   opts,
		memory: memory,
		iff1:   false, // interrupts disabled by default
		iff2:   false,
		im:     0, // interrupt mode 0 by default
	}

	// Initialize PC to reset vector (Game Boy starts at 0x0100)
	c.PC = 0x0100

	return c
}

// Cycles returns the amount of CPU cycles executed since system start.
func (c *CPU) Cycles() uint64 {
	return c.cycles
}

// Halted returns whether the CPU is in halted state.
func (c *CPU) Halted() bool {
	return c.halted
}

// Halt puts the CPU into halted state.
func (c *CPU) Halt() {
	c.halted = true
}

// Resume resumes the CPU from halted state.
func (c *CPU) Resume() {
	c.halted = false
}

// State returns the current state of the CPU.
func (c *CPU) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := State{
		A:      c.A,
		B:      c.B,
		C:      c.C,
		D:      c.D,
		E:      c.E,
		H:      c.H,
		L:      c.L,
		A_:     c.A_,
		B_:     c.B_,
		C_:     c.C_,
		D_:     c.D_,
		E_:     c.E_,
		H_:     c.H_,
		L_:     c.L_,
		IX:     c.IX,
		IY:     c.IY,
		SP:     c.SP,
		PC:     c.PC,
		I:      c.I,
		R:      c.R,
		Cycles: c.cycles,
		Flags: Flags{
			C: c.Flags.C,
			N: c.Flags.N,
			P: c.Flags.P,
			X: c.Flags.X,
			H: c.Flags.H,
			Y: c.Flags.Y,
			Z: c.Flags.Z,
			S: c.Flags.S,
		},
		Flags_: Flags{
			C: c.Flags_.C,
			N: c.Flags_.N,
			P: c.Flags_.P,
			X: c.Flags_.X,
			H: c.Flags_.H,
			Y: c.Flags_.Y,
			Z: c.Flags_.Z,
			S: c.Flags_.S,
		},
		Interrupts: Interrupts{
			IFF1:         c.iff1,
			IFF2:         c.iff2,
			IM:           c.im,
			NMITriggered: c.triggerNmi,
			IrqTriggered: c.triggerIrq,
		},
		Halted: c.halted,
	}
	return state
}

// Memory returns the CPU memory.
func (c *CPU) Memory() *Memory {
	return c.memory
}

// BC returns the BC register pair as a 16-bit value.
func (c *CPU) BC() uint16 {
	return uint16(c.B)<<8 | uint16(c.C)
}

// DE returns the DE register pair as a 16-bit value.
func (c *CPU) DE() uint16 {
	return uint16(c.D)<<8 | uint16(c.E)
}

// HL returns the HL register pair as a 16-bit value.
func (c *CPU) HL() uint16 {
	return uint16(c.H)<<8 | uint16(c.L)
}

// AF returns the AF register pair as a 16-bit value.
func (c *CPU) AF() uint16 {
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

// exchange exchanges the main and alternate register sets.
func (c *CPU) exchange() {
	c.A, c.A_ = c.A_, c.A
	c.B, c.B_ = c.B_, c.B
	c.C, c.C_ = c.C_, c.C
	c.D, c.D_ = c.D_, c.D
	c.E, c.E_ = c.E_, c.E
	c.H, c.H_ = c.H_, c.H
	c.L, c.L_ = c.L_, c.L
	c.Flags, c.Flags_ = c.Flags_, c.Flags
}

// exchangeAF exchanges only the AF and AF' registers.
func (c *CPU) exchangeAF() {
	c.A, c.A_ = c.A_, c.A
	c.Flags, c.Flags_ = c.Flags_, c.Flags
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

// TriggerNMI triggers a non-maskable interrupt.
func (c *CPU) TriggerNMI() {
	c.triggerNmi = true
}

// TriggerIRQ triggers a maskable interrupt.
func (c *CPU) TriggerIRQ() {
	c.triggerIrq = true
}
