package x86

import (
	"sync"

	"github.com/retroenv/retrogolib/arch"
)

// State contains the current state of the CPU.
// Used for save/load functionality and debugging.
type State struct {
	// General purpose registers (16-bit)
	AX uint16 // accumulator
	BX uint16 // base
	CX uint16 // count
	DX uint16 // data

	// Index and pointer registers
	SI uint16 // source index
	DI uint16 // destination index
	BP uint16 // base pointer
	SP uint16 // stack pointer

	// Segment registers
	CS uint16 // code segment
	DS uint16 // data segment
	ES uint16 // extra segment
	SS uint16 // stack segment

	// Instruction pointer and flags
	IP    uint16 // instruction pointer
	Flags Flags  // processor flags

	Cycles     uint64
	Interrupts Interrupts

	Halted bool
}

// CPU represents an Intel x86 (8086/8088) microprocessor with full instruction set emulation.
// Thread-safe through mutex locks for concurrent access.
type CPU struct {
	mu sync.RWMutex

	// General purpose registers (16-bit)
	AX uint16 // accumulator (AH:AL)
	BX uint16 // base register (BH:BL)
	CX uint16 // count register (CH:CL)
	DX uint16 // data register (DH:DL)

	// Index and pointer registers
	SI uint16 // source index
	DI uint16 // destination index
	BP uint16 // base pointer
	SP uint16 // stack pointer

	// Segment registers
	CS uint16 // code segment
	DS uint16 // data segment
	ES uint16 // extra segment
	SS uint16 // stack segment

	// Instruction pointer and flags
	IP    uint16 // instruction pointer
	Flags Flags  // processor flags

	cycles uint64
	halted bool

	// Interrupt handling
	interruptsEnabled bool
	triggerInt        bool
	intVector         uint8

	opts      Options
	TraceStep TraceStep // trace step info, set if tracing is enabled

	memory *Memory
}

// Interrupts holds the current interrupt state.
type Interrupts struct {
	Enabled      bool
	IntTriggered bool
	Vector       uint8
}

// CPU initialization constants
const (
	initialCycles = 0
	initialFlags  = 0x0002 // Reserved bit 1 is always set
)

// New creates a new x86 CPU.
func New(memory *Memory, options ...Option) (*CPU, error) {
	if memory == nil {
		return nil, ErrNilMemory
	}

	opts := NewOptions(options...)

	// Set default values for generic system if no system type specified
	if opts.initialIP == 0 && opts.initialSP == 0 && opts.systemType == "" {
		opts.systemType = string(arch.Generic)
		opts.initialIP = 0x0000
		opts.initialSP = 0xFFFE
		opts.initialCS = 0xF000 // Traditional BIOS start segment
		opts.initialSS = 0x0000
	}

	c := &CPU{
		IP:                opts.initialIP,
		SP:                opts.initialSP,
		CS:                opts.initialCS,
		SS:                opts.initialSS,
		DS:                opts.initialDS,
		ES:                opts.initialES,
		cycles:            initialCycles,
		Flags:             Flags(initialFlags),
		opts:              opts,
		memory:            memory,
		interruptsEnabled: false, // interrupts disabled by default
	}

	return c, nil
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

	return State{
		AX:     c.AX,
		BX:     c.BX,
		CX:     c.CX,
		DX:     c.DX,
		SI:     c.SI,
		DI:     c.DI,
		BP:     c.BP,
		SP:     c.SP,
		CS:     c.CS,
		DS:     c.DS,
		ES:     c.ES,
		SS:     c.SS,
		IP:     c.IP,
		Flags:  c.Flags,
		Cycles: c.cycles,
		Interrupts: Interrupts{
			Enabled:      c.interruptsEnabled,
			IntTriggered: c.triggerInt,
			Vector:       c.intVector,
		},
		Halted: c.halted,
	}
}

// Memory returns the CPU memory.
func (c *CPU) Memory() *Memory {
	return c.memory
}

// 8-bit register accessors for high/low bytes

// AL returns the low byte of AX.
func (c *CPU) AL() uint8 {
	return uint8(c.AX)
}

// AH returns the high byte of AX.
func (c *CPU) AH() uint8 {
	return uint8(c.AX >> 8)
}

// BL returns the low byte of BX.
func (c *CPU) BL() uint8 {
	return uint8(c.BX)
}

// BH returns the high byte of BX.
func (c *CPU) BH() uint8 {
	return uint8(c.BX >> 8)
}

// CL returns the low byte of CX.
func (c *CPU) CL() uint8 {
	return uint8(c.CX)
}

// CH returns the high byte of CX.
func (c *CPU) CH() uint8 {
	return uint8(c.CX >> 8)
}

// DL returns the low byte of DX.
func (c *CPU) DL() uint8 {
	return uint8(c.DX)
}

// DH returns the high byte of DX.
func (c *CPU) DH() uint8 {
	return uint8(c.DX >> 8)
}

// SetAL sets the low byte of AX.
func (c *CPU) SetAL(value uint8) {
	c.AX = (c.AX & 0xFF00) | uint16(value)
}

// SetAH sets the high byte of AX.
func (c *CPU) SetAH(value uint8) {
	c.AX = (c.AX & 0x00FF) | (uint16(value) << 8)
}

// SetBL sets the low byte of BX.
func (c *CPU) SetBL(value uint8) {
	c.BX = (c.BX & 0xFF00) | uint16(value)
}

// SetBH sets the high byte of BX.
func (c *CPU) SetBH(value uint8) {
	c.BX = (c.BX & 0x00FF) | (uint16(value) << 8)
}

// SetCL sets the low byte of CX.
func (c *CPU) SetCL(value uint8) {
	c.CX = (c.CX & 0xFF00) | uint16(value)
}

// SetCH sets the high byte of CX.
func (c *CPU) SetCH(value uint8) {
	c.CX = (c.CX & 0x00FF) | (uint16(value) << 8)
}

// SetDL sets the low byte of DX.
func (c *CPU) SetDL(value uint8) {
	c.DX = (c.DX & 0xFF00) | uint16(value)
}

// SetDH sets the high byte of DX.
func (c *CPU) SetDH(value uint8) {
	c.DX = (c.DX & 0x00FF) | (uint16(value) << 8)
}

// Segment register accessors

// SetCS sets the code segment register.
func (c *CPU) SetCS(value uint16) {
	c.CS = value
}

// SetDS sets the data segment register.
func (c *CPU) SetDS(value uint16) {
	c.DS = value
}

// SetES sets the extra segment register.
func (c *CPU) SetES(value uint16) {
	c.ES = value
}

// SetSS sets the stack segment register.
func (c *CPU) SetSS(value uint16) {
	c.SS = value
}

// SetIP sets the instruction pointer.
func (c *CPU) SetIP(value uint16) {
	c.IP = value
}

// TriggerInterrupt triggers a software interrupt.
func (c *CPU) TriggerInterrupt(vector uint8) {
	c.triggerInt = true
	c.intVector = vector
}

// EnableInterrupts enables maskable interrupts.
func (c *CPU) EnableInterrupts() {
	c.interruptsEnabled = true
}

// DisableInterrupts disables maskable interrupts.
func (c *CPU) DisableInterrupts() {
	c.interruptsEnabled = false
}

// CalculateAddress calculates the linear address from segment:offset.
func (c *CPU) CalculateAddress(segment, offset uint16) uint32 {
	return uint32(segment)<<4 + uint32(offset)
}

// GetCSIP returns the current code segment:instruction pointer address.
func (c *CPU) GetCSIP() uint32 {
	return c.CalculateAddress(c.CS, c.IP)
}

// GetSSBP returns the current stack segment:base pointer address.
func (c *CPU) GetSSBP() uint32 {
	return c.CalculateAddress(c.SS, c.BP)
}

// GetSSSP returns the current stack segment:stack pointer address.
func (c *CPU) GetSSSP() uint32 {
	return c.CalculateAddress(c.SS, c.SP)
}

// push8 pushes a byte to the stack.
func (c *CPU) push8(value uint8) {
	c.SP--
	addr := c.CalculateAddress(c.SS, c.SP)
	c.memory.Write8(addr, value)
}

// push16 pushes a word to the stack.
func (c *CPU) push16(value uint16) {
	c.SP -= 2
	addr := c.CalculateAddress(c.SS, c.SP)
	c.memory.Write16(addr, value)
}

// pop8 pops a byte from the stack.
func (c *CPU) pop8() uint8 {
	addr := c.CalculateAddress(c.SS, c.SP)
	value := c.memory.Read8(addr)
	c.SP++
	return value
}

// pop16 pops a word from the stack.
func (c *CPU) pop16() uint16 {
	addr := c.CalculateAddress(c.SS, c.SP)
	value := c.memory.Read16(addr)
	c.SP += 2
	return value
}
