package cpu6502

import (
	"errors"
	"fmt"
	"sync"
)

const (
	initialCycles = 7
	initialFlags  = 0b0010_0100 // I and U flags are set after reset.

	// InitialStack is the stack pointer value after reset.
	InitialStack = 0xFD
)

// State represents complete 6502 CPU state for save/load and debugging.
type State struct {
	// Primary registers
	A  uint8  // Accumulator (arithmetic and logic operations)
	X  uint8  // X index register
	Y  uint8  // Y index register
	PC uint16 // Program counter
	SP uint8  // Stack pointer ($0100-$01FF)

	Cycles     uint64     // Total CPU cycles executed
	Flags      Flags      // Processor status flags
	Interrupts Interrupts // Interrupt state
}

// CPU represents a 6502 microprocessor with full instruction set emulation.
// Instruction execution must be driven by one goroutine; interrupt requests may
// be triggered concurrently.
type CPU struct {
	mu sync.RWMutex

	// Primary registers
	A  uint8  // Accumulator (arithmetic and logic operations)
	X  uint8  // X index register
	Y  uint8  // Y index register
	PC uint16 // Program counter
	SP uint8  // Stack pointer ($0100-$01FF)

	Flags Flags // Processor status register

	cycles      uint64
	stallCycles uint16 // DMA transfer stall cycles

	// Interrupt control
	triggerIrq bool // IRQ interrupt triggered
	triggerNmi bool // NMI interrupt triggered
	irqRunning bool // IRQ handler executing
	nmiRunning bool // NMI handler executing

	// Interrupt vectors
	irqAddress uint16 // IRQ/BRK handler address (from $FFFE-$FFFF)
	nmiAddress uint16 // NMI handler address (from $FFFA-$FFFB)

	opts      Options
	TraceStep TraceStep // Trace step info (set if tracing enabled)

	branchTaken bool // Set by branch to distinguish a self-loop from a fallthrough.

	memory *Memory
}

// New creates a new CPU. It panics with ErrNilMemory when memory is nil.
func New(memory *Memory, options ...Option) *CPU {
	if memory == nil {
		panic(ErrNilMemory)
	}

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
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stallCycles = cycles
}

// State returns the current state of the CPU.
func (c *CPU) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return State{
		A:      c.A,
		X:      c.X,
		Y:      c.Y,
		PC:     c.PC,
		SP:     c.SP,
		Cycles: c.cycles,
		Flags:  c.Flags,
		Interrupts: Interrupts{
			NMITriggered: c.triggerNmi,
			NMIRunning:   c.nmiRunning,
			IrqTriggered: c.triggerIrq,
			IrqRunning:   c.irqRunning,
		},
	}
}

// Memory returns the CPU memory.
func (c *CPU) Memory() *Memory {
	return c.memory
}

// ValidateState performs comprehensive validation of CPU state.
// Returns an error if the CPU state is invalid or corrupted.
func (c *CPU) ValidateState() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Validate flags are 0 or 1
	if c.Flags.C > 1 || c.Flags.Z > 1 || c.Flags.I > 1 || c.Flags.D > 1 ||
		c.Flags.B > 1 || c.Flags.U > 1 || c.Flags.V > 1 || c.Flags.N > 1 {

		return errors.New("invalid flag values: flags must be 0 or 1")
	}

	// Validate memory is not nil
	if c.memory == nil {
		return ErrNilMemory
	}

	// Validate interrupt addresses are reasonable
	if c.nmiAddress == 0 && c.irqAddress == 0 {
		return errors.New("both interrupt vectors are zero")
	}

	return nil
}

// Reset resets the CPU to its initial state while preserving memory.
func (c *CPU) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Reset registers
	c.A = 0
	c.X = 0
	c.Y = 0
	c.SP = InitialStack

	// Reset flags to initial state
	c.setFlags(initialFlags)

	// Reset interrupt state
	c.triggerIrq = false
	c.triggerNmi = false
	c.irqRunning = false
	c.nmiRunning = false

	// Reset cycles
	c.cycles = initialCycles
	c.stallCycles = 0
	c.branchTaken = false
	c.TraceStep = TraceStep{}

	// Reload interrupt vectors
	if c.memory != nil {
		c.nmiAddress = c.memory.ReadWordBug(NMIAddress)
		c.PC = c.memory.ReadWordBug(ResetAddress)
		c.irqAddress = c.memory.ReadWordBug(IrqAddress)
	}
}

// GetInstructionCount returns the approximate number of instructions executed
// based on cycle count and average cycles per instruction.
func (c *CPU) GetInstructionCount() uint64 {
	const averageCyclesPerInstruction = 4
	return c.cycles / averageCyclesPerInstruction
}

// execute branch jump if the branching op result is true.
func (c *CPU) branch(branchTo bool, params ...any) error {
	if len(params) == 0 {
		return ErrMissingParameter
	}

	addr, ok := params[0].(Absolute)
	if !ok {
		return fmt.Errorf("%w: branch target type %T", ErrInvalidParameterType, params[0])
	}
	if !branchTo {
		return nil
	}

	c.PC = uint16(addr)
	c.branchTaken = true
	c.cycles++
	return nil
}

// pop pops a byte from the stack and update the stack pointer.
func (c *CPU) pop() byte {
	// The 8-bit stack pointer wraps within page one, matching the hardware.
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

func (c *CPU) consumeStallCycle() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stallCycles == 0 {
		return false
	}

	c.stallCycles--
	c.cycles++
	return true
}
