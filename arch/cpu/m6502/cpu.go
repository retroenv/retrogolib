package m6502

import (
	"sync"

	. "github.com/retroenv/retrogolib/addressing"
)

// Flags contains the status flags of the CPU.
// Bit No.   7   6   5   4   3   2   1   0
// Flag      S   V       B   D   I   Z   C
type Flags struct {
	C uint8 // carry flag
	Z uint8 // zero flag
	I uint8 // interrupt disable flag
	D uint8 // decimal mode flag
	B uint8 // break command flag
	U uint8 // unused flag
	V uint8 // overflow flag
	N uint8 // negative flag
}

// Interrupts contains the CPU interrupt info.
type Interrupts struct {
	NMITriggered bool
	NMIRunning   bool
	IrqTriggered bool
	IrqRunning   bool
}

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

	memory Memory
}

const (
	initialCycles = 7
	initialFlags  = 0b0010_0100 // I and U flags are 1, the rest 0
	InitialStack  = 0xFD
)

// New creates a new CPU.
func New(memory Memory) *CPU {
	c := &CPU{
		SP:     InitialStack,
		cycles: initialCycles,
		memory: memory,
	}

	// read interrupt handler addresses
	c.nmiAddress = memory.ReadWordBug(0xFFFA)
	c.PC = memory.ReadWordBug(0xFFFC)
	c.irqAddress = memory.ReadWordBug(0xFFFE)

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

// TriggerIrq causes a interrupt request to occur on the next cycle.
func (c *CPU) TriggerIrq() {
	c.triggerIrq = true
}

// TriggerNMI causes a non-maskable interrupt to occur on the next cycle.
func (c *CPU) TriggerNMI() {
	c.triggerNmi = true
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

// execute branch jump if the branching op result is true.
func (c *CPU) branch(branchTo bool, param any) {
	if !branchTo {
		return
	}

	addr := param.(Absolute)

	c.PC = uint16(addr)
	c.cycles++
}

func (c *CPU) irq() {
	c.mu.Lock()
	c.triggerIrq = false
	c.irqRunning = true
	c.mu.Unlock()

	c.executeInterrupt(c.irqAddress)
}

func (c *CPU) executeInterrupt(funAddress uint16) {
	c.push16(c.PC)
	php(c)

	if funAddress != 0 {
		c.Flags.I = 1
		c.cycles += 7
		c.PC = funAddress
	}
}

func (c *CPU) setFlags(flags uint8) {
	c.Flags.C = (flags >> 0) & 1
	c.Flags.Z = (flags >> 1) & 1
	c.Flags.I = (flags >> 2) & 1
	c.Flags.D = (flags >> 3) & 1
	c.Flags.B = (flags >> 4) & 1
	c.Flags.U = (flags >> 5) & 1
	c.Flags.V = (flags >> 6) & 1
	c.Flags.N = (flags >> 7) & 1
}

// GetFlags returns the current state of flags as byte.
func (c *CPU) GetFlags() uint8 {
	var f byte
	f |= c.Flags.C << 0
	f |= c.Flags.Z << 1
	f |= c.Flags.I << 2
	f |= c.Flags.D << 3
	f |= c.Flags.B << 4
	f |= c.Flags.U << 5
	f |= c.Flags.V << 6
	f |= c.Flags.N << 7
	return f
}

// setZ - set the zero flag if the argument is zero.
func (c *CPU) setZ(value uint8) {
	if value == 0 {
		c.Flags.Z = 1
	} else {
		c.Flags.Z = 0
	}
}

// setN - set the negative flag if the argument is negative (high bit is set).
func (c *CPU) setN(value uint8) {
	if value&0x80 != 0 {
		c.Flags.N = 1
	} else {
		c.Flags.N = 0
	}
}

// setV - set the overflow flag.
func (c *CPU) setV(set bool) {
	if set {
		c.Flags.V = 1
	} else {
		c.Flags.V = 0
	}
}

func (c *CPU) setZN(value uint8) {
	c.setZ(value)
	c.setN(value)
}

func (c *CPU) compare(a, b byte) {
	c.setZN(a - b)
	if a >= b {
		c.Flags.C = 1
	} else {
		c.Flags.C = 0
	}
}
