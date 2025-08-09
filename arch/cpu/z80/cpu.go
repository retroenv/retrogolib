package z80

import (
	"sync"

	"github.com/retroenv/retrogolib/arch"
)

// State contains the current state of the CPU.
// Used for save/load functionality and debugging.
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
	AltA uint8
	AltB uint8
	AltC uint8
	AltD uint8
	AltE uint8
	AltH uint8
	AltL uint8

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
	AltFlags   Flags // alternate flags
	Interrupts Interrupts

	Halted bool
}

// CPU represents a Z80 microprocessor with full instruction set emulation.
// Thread-safe through mutex locks for concurrent access.
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
	AltA uint8
	AltB uint8
	AltC uint8
	AltD uint8
	AltE uint8
	AltH uint8
	AltL uint8

	// Index registers
	IX uint16
	IY uint16

	// Special registers
	SP uint16 // stack pointer
	PC uint16 // program counter
	I  uint8  // interrupt vector
	R  uint8  // refresh register

	Flags    Flags
	AltFlags Flags // alternate flags

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

	memory *Memory
}

// Interrupts holds the current interrupt state.
// Used for interrupt management and state serialization.
type Interrupts struct {
	IFF1         bool
	IFF2         bool
	IM           uint8
	NMITriggered bool
	IrqTriggered bool
}

// CPU initialization constants
const (
	initialCycles = 0
)

// New creates a new Z80 CPU.
func New(memory *Memory, options ...Option) (*CPU, error) {
	if memory == nil {
		return nil, ErrNilMemory
	}

	opts := NewOptions(options...)

	// Set default values for generic system if no system type specified
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
		iff1:   false, // interrupts disabled by default
		iff2:   false,
		im:     0, // interrupt mode 0 by default
	}

	return c, nil
}

// Cycles returns the amount of CPU cycles executed since system start.
func (cpu *CPU) Cycles() uint64 {
	return cpu.cycles
}

// Halted returns whether the CPU is in halted state.
func (cpu *CPU) Halted() bool {
	return cpu.halted
}

// Halt puts the CPU into halted state.
func (cpu *CPU) Halt() {
	cpu.halted = true
}

// Resume resumes the CPU from halted state.
func (cpu *CPU) Resume() {
	cpu.halted = false
}

// State returns the current state of the CPU.
func (cpu *CPU) State() State {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()

	return State{
		A:        cpu.A,
		B:        cpu.B,
		C:        cpu.C,
		D:        cpu.D,
		E:        cpu.E,
		H:        cpu.H,
		L:        cpu.L,
		AltA:     cpu.AltA,
		AltB:     cpu.AltB,
		AltC:     cpu.AltC,
		AltD:     cpu.AltD,
		AltE:     cpu.AltE,
		AltH:     cpu.AltH,
		AltL:     cpu.AltL,
		IX:       cpu.IX,
		IY:       cpu.IY,
		SP:       cpu.SP,
		PC:       cpu.PC,
		I:        cpu.I,
		R:        cpu.R,
		Cycles:   cpu.cycles,
		Flags:    cpu.Flags,
		AltFlags: cpu.AltFlags,
		Interrupts: Interrupts{
			IFF1:         cpu.iff1,
			IFF2:         cpu.iff2,
			IM:           cpu.im,
			NMITriggered: cpu.triggerNmi,
			IrqTriggered: cpu.triggerIrq,
		},
		Halted: cpu.halted,
	}
}

// Memory returns the CPU memory.
func (cpu *CPU) Memory() *Memory {
	return cpu.memory
}

// BC returns the BC register pair as a 16-bit value.
func (cpu *CPU) BC() uint16 {
	return uint16(cpu.B)<<8 | uint16(cpu.C)
}

// DE returns the DE register pair as a 16-bit value.
func (cpu *CPU) DE() uint16 {
	return uint16(cpu.D)<<8 | uint16(cpu.E)
}

// HL returns the HL register pair as a 16-bit value.
func (cpu *CPU) HL() uint16 {
	return uint16(cpu.H)<<8 | uint16(cpu.L)
}

// AF returns the AF register pair as a 16-bit value.
func (cpu *CPU) AF() uint16 {
	return uint16(cpu.A)<<8 | uint16(cpu.GetFlags())
}

// TriggerNMI triggers a non-maskable interrupt.
func (cpu *CPU) TriggerNMI() {
	cpu.triggerNmi = true
}

// TriggerIRQ triggers a maskable interrupt.
func (cpu *CPU) TriggerIRQ() {
	cpu.triggerIrq = true
}

// setBC sets the BC register pair from a 16-bit value.
func (cpu *CPU) setBC(value uint16) {
	cpu.B = uint8(value >> 8)
	cpu.C = uint8(value)
}

// setDE sets the DE register pair from a 16-bit value.
func (cpu *CPU) setDE(value uint16) {
	cpu.D = uint8(value >> 8)
	cpu.E = uint8(value)
}

// setHL sets the HL register pair from a 16-bit value.
func (cpu *CPU) setHL(value uint16) {
	cpu.H = uint8(value >> 8)
	cpu.L = uint8(value)
}

// setAF sets the AF register pair from a 16-bit value.
func (cpu *CPU) setAF(value uint16) {
	cpu.A = uint8(value >> 8)
	cpu.setFlags(uint8(value))
}

// exchange exchanges the main and alternate register sets.
func (cpu *CPU) exchange() {
	cpu.A, cpu.AltA = cpu.AltA, cpu.A
	cpu.B, cpu.AltB = cpu.AltB, cpu.B
	cpu.C, cpu.AltC = cpu.AltC, cpu.C
	cpu.D, cpu.AltD = cpu.AltD, cpu.D
	cpu.E, cpu.AltE = cpu.AltE, cpu.E
	cpu.H, cpu.AltH = cpu.AltH, cpu.H
	cpu.L, cpu.AltL = cpu.AltL, cpu.L
	cpu.Flags, cpu.AltFlags = cpu.AltFlags, cpu.Flags
}

// exchangeAF exchanges only the AF and AF' registers.
func (cpu *CPU) exchangeAF() {
	cpu.A, cpu.AltA = cpu.AltA, cpu.A
	cpu.Flags, cpu.AltFlags = cpu.AltFlags, cpu.Flags
}

// pop pops a byte from the stack and updates the stack pointer.
func (cpu *CPU) pop() uint8 {
	value := cpu.memory.Read(cpu.SP)
	cpu.SP++
	return value
}

// pop16 pops a word from the stack and updates the stack pointer.
func (cpu *CPU) pop16() uint16 {
	low := uint16(cpu.pop())
	high := uint16(cpu.pop())
	return high<<8 | low
}

// push pushes a byte to the stack and updates the stack pointer.
func (cpu *CPU) push(value uint8) {
	cpu.SP--
	cpu.memory.Write(cpu.SP, value)
}

// push16 pushes a word to the stack and updates the stack pointer.
func (cpu *CPU) push16(value uint16) {
	high := uint8(value >> 8)
	low := uint8(value)
	cpu.push(high)
	cpu.push(low)
}
