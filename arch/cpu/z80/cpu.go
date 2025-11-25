package z80

import (
	"sync"

	"github.com/retroenv/retrogolib/arch"
)

// State represents complete CPU state for save/load and debugging.
type State struct {
	// Main 8-bit registers (can be paired as BC, DE, HL)
	A uint8 // Accumulator
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	// Shadow register set (accessed via EX AF,AF' and EXX)
	AltA uint8
	AltB uint8
	AltC uint8
	AltD uint8
	AltE uint8
	AltH uint8
	AltL uint8

	// Index registers for offset addressing
	IX uint16
	IY uint16

	// Program flow and memory
	SP uint16 // Stack pointer
	PC uint16 // Program counter
	I  uint8  // Interrupt vector base
	R  uint8  // Memory refresh counter

	Cycles     uint64
	Flags      Flags
	AltFlags   Flags
	Interrupts Interrupts

	Halted bool
}

// CPU represents a thread-safe Z80 microprocessor with full instruction set emulation.
type CPU struct {
	mu sync.RWMutex

	// Main 8-bit general purpose registers
	A uint8 // Accumulator (used in arithmetic/logic ops)
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	// Shadow register set
	AltA uint8
	AltB uint8
	AltC uint8
	AltD uint8
	AltE uint8
	AltH uint8
	AltL uint8

	// 16-bit index registers
	IX uint16
	IY uint16

	// Program control registers
	SP uint16 // Stack pointer
	PC uint16 // Program counter
	I  uint8  // Interrupt vector base register
	R  uint8  // Memory refresh register (auto-incremented)

	Flags    Flags // Main flag register
	AltFlags Flags // Shadow flag register

	cycles uint64
	halted bool

	// Interrupt control
	iff1 bool  // Interrupt enable flip-flop
	iff2 bool  // Backup of IFF1 for NMI handling
	im   uint8 // Interrupt mode: 0, 1, or 2

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

// inPortToRegister reads from port C to a register and sets flags.
func (c *CPU) inPortToRegister(regPtr *uint8) {
	value := c.readPort(c.C)
	*regPtr = value
	c.setSZP(value)
	c.setH(false)
	c.setN(false)
}

// applyCBOperation applies a CB prefix operation to a register or (HL).
// Used by rotate, shift, RES, and SET instructions.
func (c *CPU) applyCBOperation(operation func(uint8) uint8) {
	opcodeByte := c.memory.Read(c.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // Operation on (HL)
		addr := c.hl()
		value := c.memory.Read(addr)
		result := operation(value)
		c.memory.Write(addr, result)
	} else { // Operation on register
		value := c.GetRegisterValue(reg)
		result := operation(value)
		c.SetRegisterValue(reg, result)
	}
}

// calculateIndexedAddress extracts displacement from params and calculates indexed address.
// Used by DD (IX) and FD (IY) prefix instructions.
func (c *CPU) calculateIndexedAddress(indexReg uint16, params ...any) uint16 {
	displacement := int8(params[0].(uint8))
	return uint16(int32(indexReg) + int32(displacement))
}

// extractExtendedAddress extracts 16-bit address from instruction parameters (little-endian).
// Used by instructions that take a 16-bit address operand.
func extractExtendedAddress(params ...any) uint16 {
	return uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
}

// read16 reads a 16-bit value from memory at addr (little-endian).
func (c *CPU) read16(addr uint16) uint16 {
	low := c.memory.Read(addr)
	high := c.memory.Read(addr + 1)
	return uint16(high)<<8 | uint16(low)
}

// writeRegisterPair writes a register pair to memory at addr (little-endian).
func (c *CPU) writeRegisterPair(addr uint16, low, high uint8) {
	c.memory.Write(addr, low)
	c.memory.Write(addr+1, high)
}

// setLogicalFlags sets flags for logical operations (AND/OR/XOR).
// hFlag should be true for AND, false for OR/XOR.
func (c *CPU) setLogicalFlags(result uint8, hFlag bool) {
	c.setSZP(result)
	c.setH(hFlag)
	c.setN(false)
	c.setC(false)
}
