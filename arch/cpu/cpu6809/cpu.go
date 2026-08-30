package cpu6809

import "sync"

// State represents a complete snapshot of the 6809 CPU state.
type State struct {
	A  uint8  // Accumulator A
	B  uint8  // Accumulator B
	X  uint16 // Index register X
	Y  uint16 // Index register Y
	S  uint16 // System stack pointer
	U  uint16 // User stack pointer
	DP uint8  // Direct page register
	PC uint16 // Program counter
	CC uint8  // Condition codes (from Flags.Get())

	Cycles uint64
}

// CPU represents a Motorola 6809 microprocessor.
// Interrupt triggers may run concurrently with Step; callers must serialize all other access.
type CPU struct {
	interruptMu sync.Mutex

	// Registers
	A  uint8  // Accumulator A (high byte of D)
	B  uint8  // Accumulator B (low byte of D)
	X  uint16 // Index register X
	Y  uint16 // Index register Y
	S  uint16 // System stack pointer
	U  uint16 // User stack pointer
	DP uint8  // Direct page register
	PC uint16 // Program counter

	Flags Flags // Condition code register (CC)

	cycles    uint64
	waitMode  waitMode
	pcChanged bool   // set by instructions that explicitly set PC (branches, jumps)
	nextPC    uint16 // address following the current instruction

	// Interrupt control
	triggerNMI  bool
	triggerIRQ  bool
	triggerFIRQ bool

	memory *Memory
	opts   options

	TraceStep TraceStep // set when tracing is enabled
}

// TraceStep holds information for instruction tracing.
type TraceStep struct {
	PC             uint16
	OpcodeOperands []byte
	Opcode         Opcode
}

// New creates a new 6809 CPU, reads the reset vector, and initializes registers.
func New(memory *Memory, opts ...Option) (*CPU, error) {
	if memory == nil || memory.BasicMemory == nil {
		return nil, ErrNilMemory
	}

	c := &CPU{
		memory: memory,
		opts:   newOptions(opts...),
	}

	// Reset starts with both maskable interrupt classes disabled.
	c.Flags.I = 1
	c.Flags.F = 1

	resetVec := memory.ReadVector(VectorRESET)
	c.PC = resetVec

	return c, nil
}

// D returns the 16-bit D register (A:B, A=high byte, B=low byte).
func (c *CPU) D() uint16 {
	return uint16(c.A)<<8 | uint16(c.B)
}

// SetD sets the 16-bit D register (A:B).
func (c *CPU) SetD(d uint16) {
	c.A = uint8(d >> 8)
	c.B = uint8(d)
}

// Cycles returns the total number of cycles executed.
func (c *CPU) Cycles() uint64 { return c.cycles }

// State returns a snapshot of the current CPU state.
func (c *CPU) State() State {
	return State{
		A:      c.A,
		B:      c.B,
		X:      c.X,
		Y:      c.Y,
		S:      c.S,
		U:      c.U,
		DP:     c.DP,
		PC:     c.PC,
		CC:     c.Flags.Get(),
		Cycles: c.cycles,
	}
}

// ValidateState checks that CPU state is consistent.
func (c *CPU) ValidateState() error {
	if c.memory == nil || c.memory.BasicMemory == nil {
		return ErrNilMemory
	}
	return nil
}

// Reset resets the CPU to its initial post-reset state.
func (c *CPU) Reset() {
	c.A = 0
	c.B = 0
	c.X = 0
	c.Y = 0
	c.S = 0
	c.U = 0
	c.DP = 0
	c.Flags = Flags{I: 1, F: 1}
	c.cycles = 0

	c.interruptMu.Lock()
	c.waitMode = waitNone
	c.triggerNMI = false
	c.triggerIRQ = false
	c.triggerFIRQ = false
	c.interruptMu.Unlock()

	if c.memory != nil && c.memory.BasicMemory != nil {
		c.PC = c.memory.ReadVector(VectorRESET)
	}
}

// Memory returns the CPU's memory.
func (c *CPU) Memory() *Memory { return c.memory }

// GetCC returns the current condition code register byte.
func (c *CPU) GetCC() uint8 {
	return c.Flags.Get()
}

// SetCC sets the condition code register from a byte.
func (c *CPU) SetCC(cc uint8) {
	c.Flags.Set(cc)
}

func (c *CPU) setWaitMode(mode waitMode) {
	c.interruptMu.Lock()
	c.waitMode = mode
	c.interruptMu.Unlock()
}

func (c *CPU) isWaiting() bool {
	c.interruptMu.Lock()
	defer c.interruptMu.Unlock()

	return c.waitMode != waitNone
}

// pushS8 pushes a byte onto the system stack (S) and decrements S.
func (c *CPU) pushS8(value uint8) {
	c.S--
	c.memory.Write(c.S, value)
}

// pushS16 pushes the low byte first so the decremented stack points at the high byte.
func (c *CPU) pushS16(value uint16) {
	c.pushS8(uint8(value))
	c.pushS8(uint8(value >> 8))
}

// popS8 pops a byte from the system stack and increments S.
func (c *CPU) popS8() uint8 {
	value := c.memory.Read(c.S)
	c.S++
	return value
}

// popS16 pops a 16-bit word from the system stack (high byte first).
func (c *CPU) popS16() uint16 {
	hi := uint16(c.popS8())
	lo := uint16(c.popS8())
	return hi<<8 | lo
}

// pushU8 pushes a byte onto the user stack (U) and decrements U.
func (c *CPU) pushU8(value uint8) {
	c.U--
	c.memory.Write(c.U, value)
}

// pushU16 pushes the low byte first so the decremented stack points at the high byte.
func (c *CPU) pushU16(value uint16) {
	c.pushU8(uint8(value))
	c.pushU8(uint8(value >> 8))
}

// popU8 pops a byte from the user stack and increments U.
func (c *CPU) popU8() uint8 {
	value := c.memory.Read(c.U)
	c.U++
	return value
}

// popU16 pops a 16-bit word from the user stack.
func (c *CPU) popU16() uint16 {
	hi := uint16(c.popU8())
	lo := uint16(c.popU8())
	return hi<<8 | lo
}

// dpAddr forms a 16-bit address using the direct page register.
func (c *CPU) dpAddr(offset uint8) uint16 {
	return uint16(c.DP)<<8 | uint16(offset)
}

// fetchByte reads the next byte from PC without advancing PC.
func (c *CPU) fetchByte(offset uint16) uint8 {
	return c.memory.Read(c.PC + offset)
}

// branch performs a relative branch if the condition is true.
// addr is the pre-computed absolute branch target address.
func (c *CPU) branch(taken bool, addr uint16) {
	if !taken {
		return
	}
	c.PC = addr
	c.pcChanged = true
}
