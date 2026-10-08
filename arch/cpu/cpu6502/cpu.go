package cpu6502

import (
	"errors"
	"fmt"
	"sync"
)

const (
	initialCycles = 7
	initialFlags  = 0b0010_0100 // I and U flags are set after reset.

	// decimalCorrectionCycles is the extra cycle the 65C02 takes to correct a
	// decimal mode ADC or SBC result.
	decimalCorrectionCycles = 1

	// InitialStack is the stack pointer value after reset.
	InitialStack = 0xFD
)

// State represents complete 6502 CPU state for save/load and debugging.
type State struct {
	// Primary registers
	A  uint8  // Accumulator (arithmetic and logic operations)
	X  uint8  // X index register
	Y  uint8  // Y index register
	PC uint16 // Program counter captured when State was called
	SP uint8  // Stack pointer ($0100-$01FF)

	Cycles     uint64     // Total CPU cycles captured when State was called
	Flags      Flags      // Processor status flags
	Interrupts Interrupts // Pending and running interrupt state, not the source of a stack event
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
	PC uint16 // Live program counter; StackEvent.PC captures an execution boundary
	SP uint8  // Stack pointer ($0100-$01FF)

	Flags Flags // Processor status register

	// Cycle accounting
	cycles      uint64 // live total; executionCycle captures it before instruction timing is added
	stallCycles uint16 // DMA transfer stall cycles

	// Bus-cycle execution
	cycleActive bool // Bus-cycle execution is active inside an instruction or reset.
	jamCycle    byte // Initial bus sequence after KIL.
	jammed      bool // KIL stops instruction execution until reset.

	// Interrupt inputs
	irqLine    bool // IRQ input level
	triggerIrq bool // IRQ interrupt triggered
	triggerNmi bool // NMI interrupt triggered

	// Interrupt sampling
	irqPolled bool // IRQ selected by the instruction's last polling cycle.
	irqSample bool // IRQ input and mask at the end of the last bus cycle.
	nmiPolled bool // NMI selected by the instruction's last polling cycle.
	nmiSample bool // NMI request at the end of the last bus cycle.

	// Interrupt handler state
	irqRunning bool // IRQ/BRK handler active until RTI; not the source of the current stack operation
	nmiRunning bool // NMI handler active until RTI; not the source of the current stack operation

	opts      options
	TraceStep TraceStep // Trace step info (set if tracing enabled)

	branchTaken bool // Set by branch to distinguish a self-loop from a fallthrough.

	executionCycle uint64          // cycle count at the current execution boundary
	executionPC    uint16          // instruction address or interrupted program counter
	interrupt      InterruptSource // interrupt source for the current stack operation
	opcode         byte            // opcode for the current instruction, or zero for IRQ/NMI

	memory *Memory
}

// New creates a new CPU. It panics with ErrNilMemory when memory is nil.
func New(memory *Memory, options ...Option) *CPU {
	if memory == nil {
		panic(ErrNilMemory)
	}

	opts := newOptions(options...)
	c := &CPU{
		SP:     InitialStack,
		cycles: initialCycles,
		opts:   opts,
		memory: memory,
	}

	c.PC = memory.ReadWordBug(ResetAddress)

	c.setFlags(initialFlags)
	return c
}

// Cycles returns the live total CPU cycles executed since system start.
// StackEvent.Cycle instead captures this total at the start of the instruction
// or interrupt that produced the event.
func (c *CPU) Cycles() uint64 {
	return c.cycles
}

// StallCycles adds cycles during which Step does not execute an instruction.
// It is safe to call from a different goroutine than the one that calls Step.
func (c *CPU) StallCycles(cycles uint16) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stallCycles += cycles
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
			IrqTriggered: c.triggerIrq || c.irqLine,
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

	return nil
}

// Reset resets the CPU to its initial state while preserving memory.
func (c *CPU) Reset() {
	if c.opts.cycleHook != nil && c.opts.variant < Variant65C02 {
		c.resetCycles()
		return
	}
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
	c.irqLine = false
	c.irqRunning = false
	c.nmiRunning = false

	// Reset cycles
	c.cycles = initialCycles
	c.stallCycles = 0
	c.branchTaken = false
	c.TraceStep = TraceStep{}
	c.executionPC = 0
	c.executionCycle = 0
	c.opcode = 0
	c.interrupt = InterruptNone

	// Reload the reset vector.
	if c.memory != nil {
		c.PC = c.memory.ReadWordBug(ResetAddress)
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
	if !c.cycleActive {
		c.cycles++
	}
	return nil
}

// pop pops a byte from the stack and update the stack pointer.
func (c *CPU) pop() byte {
	// The 8-bit stack pointer wraps within page one, matching the hardware.
	before := c.SP
	c.SP++
	c.emitStackEvent(StackPull, before)
	if c.cycleActive {
		return c.readCycle(uint16(StackBase + int(c.SP)))
	}
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
	before := c.SP
	address := uint16(StackBase + int(c.SP))
	if c.cycleActive {
		c.writeCycle(address, value)
	} else {
		c.memory.Write(address, value)
	}
	c.SP--
	c.emitStackEvent(StackPush, before)
}

func (c *CPU) emitStackEvent(operation StackOperation, before byte) {
	if c.opts.stackEventHook == nil {
		return
	}

	c.opts.stackEventHook(StackEvent{
		After:     c.SP,
		Before:    before,
		Cycle:     c.executionCycle,
		Interrupt: c.interrupt,
		Opcode:    c.opcode,
		Operation: operation,
		PC:        c.executionPC,
	})
}

// push16 a word to the stack and update the stack pointer.
func (c *CPU) push16(value uint16) {
	high := byte(value >> 8)
	low := byte(value)
	c.push(high)
	c.push(low)
}

// consumeStallCycle takes one pending stall cycle. It reports false when no
// stall is pending. The caller accounts for the cycle.
func (c *CPU) consumeStallCycle() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stallCycles == 0 {
		return false
	}

	c.stallCycles--
	return true
}
