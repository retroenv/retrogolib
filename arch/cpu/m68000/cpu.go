package m68000

import (
	"sync"
)

// State represents complete CPU state for save/load and debugging.
type State struct {
	D      [8]uint32 // Data registers D0-D7
	A      [7]uint32 // Address registers A0-A6
	USP    uint32    // User stack pointer
	SSP    uint32    // Supervisor stack pointer
	SP     uint32    // Active stack pointer
	PC     uint32    // Program counter
	SR     uint16    // Status register
	Flags  Flags     // CCR flags
	Cycles uint64
	Halted bool
}

// CPU represents a thread-safe Motorola 68000 microprocessor with full instruction set emulation.
type CPU struct {
	mu sync.RWMutex

	D [8]uint32 // Data registers D0-D7
	A [7]uint32 // Address registers A0-A6

	USP uint32 // User stack pointer
	SSP uint32 // Supervisor stack pointer
	sp  uint32 // Active stack pointer (USP or SSP based on mode)
	PC  uint32 // Program counter

	Flags Flags // CCR flags

	cycles  uint64
	halted  bool
	stopped bool // STOP instruction state

	sr  uint16 // Status register system byte (high byte)
	bus Bus

	opts      Options
	TraceStep TraceStep // trace step info, set if tracing is enabled
}

// New creates a new 68000 CPU with a bus interface.
// On reset, the 68000 loads the initial SSP from vector 0 and initial PC from vector 1.
func New(bus Bus, options ...Option) (*CPU, error) {
	if bus == nil {
		return nil, ErrNilBus
	}

	opts := NewOptions(options...)

	c := &CPU{
		bus:  bus,
		opts: opts,
	}

	// Start in supervisor mode with interrupts masked.
	c.sr = MaskSupervisor | MaskIPM

	if opts.initialPC != 0 || opts.initialSP != 0 {
		c.PC = opts.initialPC
		c.sp = opts.initialSP
		c.SSP = opts.initialSP
	} else {
		// Standard 68000 reset: load SSP from vector 0, PC from vector 1.
		c.SSP = bus.ReadLong(0x000000)
		c.sp = c.SSP
		c.PC = bus.ReadLong(0x000004)
	}

	return c, nil
}

// A7 returns the active stack pointer based on the current privilege mode.
func (c *CPU) A7() uint32 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sp
}

// Bus returns the attached bus interface.
//
//nolint:ireturn
func (c *CPU) Bus() Bus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bus
}

// Cycles returns total CPU cycles executed.
func (c *CPU) Cycles() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cycles
}

// Halt stops CPU execution.
func (c *CPU) Halt() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.halted = true
}

// Halted returns CPU halt state.
func (c *CPU) Halted() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.halted
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
		D:      c.D,
		A:      c.A,
		USP:    c.USP,
		SSP:    c.SSP,
		SP:     c.sp,
		PC:     c.PC,
		SR:     c.GetSR(),
		Flags:  c.Flags,
		Cycles: c.cycles,
		Halted: c.halted,
	}
}

// push16 pushes a 16-bit word onto the stack (big-endian, predecrement).
func (c *CPU) push16(value uint16) {
	c.sp -= 2
	c.bus.WriteWord(c.sp, value)
}

// push32 pushes a 32-bit long word onto the stack (big-endian, predecrement).
func (c *CPU) push32(value uint32) {
	c.sp -= 4
	c.bus.WriteLong(c.sp, value)
}

// pop16 pops a 16-bit word from the stack (postincrement).
func (c *CPU) pop16() uint16 {
	value := c.bus.ReadWord(c.sp)
	c.sp += 2
	return value
}

// pop32 pops a 32-bit long word from the stack (postincrement).
func (c *CPU) pop32() uint32 {
	value := c.bus.ReadLong(c.sp)
	c.sp += 4
	return value
}

// readWord reads a word from the instruction stream and advances PC.
func (c *CPU) readWord() uint16 {
	value := c.bus.ReadWord(c.PC)
	c.PC += 2
	return value
}

// readLong reads a long from the instruction stream and advances PC.
func (c *CPU) readLong() uint32 {
	value := c.bus.ReadLong(c.PC)
	c.PC += 4
	return value
}

// readImmediate reads an immediate value from the instruction stream.
// Byte-sized immediates occupy the low byte of a word.
func (c *CPU) readImmediate(size OperandSize) uint32 {
	switch size {
	case SizeByte:
		w := c.readWord()
		return uint32(w & 0xFF)
	case SizeWord:
		return uint32(c.readWord())
	case SizeLong:
		return c.readLong()
	default:
		return 0
	}
}

// getRegD returns the value of data register Dn masked to the given size.
func (c *CPU) getRegD(reg uint8, size OperandSize) uint32 {
	return maskValue(c.D[reg], size)
}

// setRegD sets the data register Dn, preserving upper bits for byte/word operations.
func (c *CPU) setRegD(reg uint8, value uint32, size OperandSize) {
	switch size {
	case SizeByte:
		c.D[reg] = (c.D[reg] & 0xFFFFFF00) | (value & 0xFF)
	case SizeWord:
		c.D[reg] = (c.D[reg] & 0xFFFF0000) | (value & 0xFFFF)
	case SizeLong:
		c.D[reg] = value
	}
}

// getRegA returns the value of address register An (A0-A6 or A7/SP).
func (c *CPU) getRegA(reg uint8) uint32 {
	if reg == 7 {
		return c.sp
	}
	return c.A[reg]
}

// setRegA sets the value of address register An (A0-A6 or A7/SP).
func (c *CPU) setRegA(reg uint8, value uint32) {
	if reg == 7 {
		c.sp = value
	} else {
		c.A[reg] = value
	}
}

// incrementSize returns the increment amount for the given size.
// For A7 with byte size, returns 2 to maintain word alignment.
func incrementSize(reg uint8, size OperandSize) uint32 {
	if size == SizeByte && reg == 7 {
		return 2
	}
	return uint32(size)
}
