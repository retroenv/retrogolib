package m6502

import (
	"sync"

	. "github.com/retroenv/retrogolib/addressing"
)

// State contains the current state of the CPU.
type State struct {
	A          uint8
	X          uint8
	Y          uint8
	PC         uint16
	SP         uint8
	Cycles     uint64
	Flags      Flags
	Interrupts Interrupts
}

type CPU struct {
	mu sync.RWMutex

	A     uint8  // accumulator
	X     uint8  // x register
	Y     uint8  // y register
	PC    uint16 // program counter
	SP    uint8  // stack pointer
	Flags Flags

	cycles      uint64
	stallCycles uint16 // TODO stall cycles, use a Step() function

	triggerIrq bool
	triggerNmi bool

	irqRunning bool
	nmiRunning bool

	irqAddress uint16
	nmiAddress uint16

	opts      Options
	TraceStep TraceStep // trace step info, set if tracing is enabled

	memory *Memory
}

const (
	initialCycles = 7
	initialFlags  = 0b0010_0100 // I and U flags are 1, the rest 0
	InitialStack  = 0xFD
)

// New creates a new CPU.
func New(memory *Memory, options ...Option) *CPU {
	opts := NewOptions(options...)
	c := &CPU{
		SP:     InitialStack,
		cycles: initialCycles,
		opts:   opts,
		memory: memory,
	}

	// read interrupt handler addresses
	c.nmiAddress = memory.ReadWordBug(NMIAddress)
	c.PC = memory.ReadWordBug(ResetAddress)
	c.irqAddress = memory.ReadWordBug(IrqAddress)

	c.setFlags(initialFlags)
	return c
}

// Cycles returns the amount of CPU cycles executed since system start.
func (c *CPU) Cycles() uint64 {
	return c.cycles
}

// StallCycles stalls the CPU for the given amount of cycles. This is used for DMA transfer in the PPU.
func (c *CPU) StallCycles(cycles uint16) {
	c.stallCycles = cycles
}

// State returns the current state of the CPU.
func (c *CPU) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()

	state := State{
		A:      c.A,
		X:      c.X,
		Y:      c.Y,
		PC:     c.PC,
		SP:     c.SP,
		Cycles: c.cycles,
		Flags: Flags{
			C: c.Flags.C,
			Z: c.Flags.Z,
			I: c.Flags.I,
			D: c.Flags.D,
			B: c.Flags.B,
			V: c.Flags.V,
			N: c.Flags.N,
		},
		Interrupts: Interrupts{
			NMITriggered: c.triggerNmi,
			NMIRunning:   c.nmiRunning,
			IrqTriggered: c.triggerIrq,
			IrqRunning:   c.irqRunning,
		},
	}
	return state
}

// Memory returns the CPU memory.
func (c *CPU) Memory() *Memory {
	return c.memory
}

// execute branch jump if the branching op result is true.
func (c *CPU) branch(branchTo bool, param any) {
	if !branchTo {
		return
	}

	addr := param.(Absolute)

	c.PC = uint16(addr)
	c.cycles++
}

// pop pops a byte from the stack and update the stack pointer.
func (c *CPU) pop() byte {
	c.SP++
	return c.memory.Read(uint16(StackBase + int(c.SP)))
}

// pop16 pops a word from the stack and updates the stack pointer.
func (c *CPU) pop16() uint16 {
	low := uint16(c.pop())
	high := uint16(c.pop())
	return high<<8 | low
}

// push a value to the stack and update the stack pointer.
func (c *CPU) push(value byte) {
	c.memory.Write(uint16(StackBase+int(c.SP)), value)
	c.SP--
}

// push16 a word to the stack and update the stack pointer.
func (c *CPU) push16(value uint16) {
	high := byte(value >> 8)
	low := byte(value)
	c.push(high)
	c.push(low)
}
