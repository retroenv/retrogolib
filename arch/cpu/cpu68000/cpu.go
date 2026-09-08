package cpu68000

import (
	"sync"
)

// State is a register and execution-status snapshot for debugging.
// It is not a complete save/restore format.
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

// CPU emulates a Motorola 68000 microprocessor.
// Step, Reset, and interrupt triggers are synchronized. Callers must serialize
// direct register access and status-register updates with execution.
type CPU struct {
	mu sync.RWMutex

	D [8]uint32 // Data registers D0-D7
	A [7]uint32 // Address registers A0-A6

	USP uint32 // User stack pointer
	SSP uint32 // Supervisor stack pointer
	sp  uint32 // Active stack pointer (USP or SSP based on mode)
	PC  uint32 // Program counter

	Flags Flags // CCR flags

	cycles      uint64
	halted      bool
	stopped     bool  // STOP instruction state
	pendingIRQ  uint8 // Highest queued interrupt request, retained while masked.
	faultHalted bool  // Double fault; only external reset can resume execution.

	instructionPC   uint32
	instructionWord uint16
	stepCycles      uint64
	exceptionAccess bool
	exceptionRaised bool
	operandPCOffset int32
	accessCycles    uint64

	sr  uint16 // Status register system byte (high byte)
	bus Bus

	opts      options
	TraceStep TraceStep // trace step info, set if tracing is enabled
}

// New creates a new 68000 CPU with a bus interface.
// On reset, the 68000 loads the initial SSP from vector 0 and initial PC from vector 1.
func New(bus Bus, options ...Option) (*CPU, error) {
	if bus == nil {
		return nil, ErrNilBus
	}

	opts := newOptions(options...)

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
		if err := c.Reset(); err != nil {
			return nil, err
		}
		c.cycles = 0
	}

	return c, nil
}

// A7 returns the active stack pointer based on the current privilege mode.
func (cpu *CPU) A7() uint32 {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.sp
}

// Bus returns the attached bus interface.
//
//nolint:ireturn // intentional: Bus is the public API interface
func (cpu *CPU) Bus() Bus {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.bus
}

// Cycles returns total CPU cycles executed.
func (cpu *CPU) Cycles() uint64 {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.cycles
}

// Halt stops CPU execution.
func (cpu *CPU) Halt() {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	cpu.halted = true
}

// Halted returns CPU halt state.
func (cpu *CPU) Halted() bool {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.halted
}

// Resume continues CPU execution.
func (cpu *CPU) Resume() {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	if !cpu.faultHalted {
		cpu.halted = false
	}
}

// Reset performs an external reset, loading SSP and PC from the reset vectors.
// A failed vector read leaves the CPU halted and returns the access error.
func (cpu *CPU) Reset() error {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	cpu.sr = MaskSupervisor | MaskIPM
	cpu.Flags = Flags{}
	cpu.pendingIRQ = 0
	cpu.halted, cpu.stopped, cpu.faultHalted = false, false, false
	cpu.exceptionAccess = true
	cpu.exceptionRaised = false
	cpu.instructionWord = 0
	cpu.accessCycles = 0
	cpu.cycles = 40
	err := catchAccessFault(func() error {
		cpu.sp = cpu.readBusLong(VectorResetSSP * 4)
		cpu.SSP = cpu.sp
		cpu.PC = cpu.readBusLong(VectorResetPC * 4)
		return nil
	})
	cpu.exceptionAccess = false
	if err != nil {
		cpu.halted, cpu.faultHalted = true, true
	}
	return err
}

// State returns a register and execution-status snapshot.
func (cpu *CPU) State() State {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()

	return State{
		D:      cpu.D,
		A:      cpu.A,
		USP:    cpu.USP,
		SSP:    cpu.SSP,
		SP:     cpu.sp,
		PC:     cpu.PC,
		SR:     cpu.GetSR(),
		Flags:  cpu.Flags,
		Cycles: cpu.cycles,
		Halted: cpu.halted,
	}
}

// push16 pushes a 16-bit word onto the stack (big-endian, predecrement).
func (cpu *CPU) push16(value uint16) {
	cpu.sp -= 2
	cpu.writeBusWord(cpu.sp, value)
}

// push32 pushes a 32-bit long word onto the stack (big-endian, predecrement).
func (cpu *CPU) push32(value uint32) {
	cpu.sp -= 4
	cpu.writeBusLong(cpu.sp, value)
}

// pop16 pops a 16-bit word from the stack (postincrement).
func (cpu *CPU) pop16() uint16 {
	value := cpu.readBusWord(cpu.sp, dataSpace)
	cpu.sp += 2
	return value
}

// pop32 pops a 32-bit long word from the stack (postincrement).
func (cpu *CPU) pop32() uint32 {
	value := cpu.readBusLong(cpu.sp)
	cpu.sp += 4
	return value
}

// readWord reads a word from the instruction stream and advances PC.
func (cpu *CPU) readWord() uint16 {
	value := cpu.readBusWord(cpu.PC, programSpace)
	cpu.PC += 2
	return value
}

// readLong reads a long from the instruction stream and advances PC.
func (cpu *CPU) readLong() uint32 {
	high := uint32(cpu.readWord())
	return high<<16 | uint32(cpu.readWord())
}

// readImmediate reads an immediate value from the instruction stream.
// Byte-sized immediates occupy the low byte of a word.
func (cpu *CPU) readImmediate(size OperandSize) uint32 {
	switch size {
	case SizeByte:
		w := cpu.readWord()
		return uint32(w & 0xFF)
	case SizeWord:
		return uint32(cpu.readWord())
	case SizeLong:
		return cpu.readLong()
	default:
		return 0
	}
}

// getRegD returns the value of data register Dn masked to the given size.
func (cpu *CPU) getRegD(reg uint8, size OperandSize) uint32 {
	return maskValue(cpu.D[reg], size)
}

// setRegD sets the data register Dn, preserving upper bits for byte/word operations.
func (cpu *CPU) setRegD(reg uint8, value uint32, size OperandSize) {
	switch size {
	case SizeByte:
		cpu.D[reg] = (cpu.D[reg] & 0xFFFFFF00) | (value & 0xFF)
	case SizeWord:
		cpu.D[reg] = (cpu.D[reg] & 0xFFFF0000) | (value & 0xFFFF)
	case SizeLong:
		cpu.D[reg] = value
	}
}

// getRegA returns the value of address register An (A0-A6 or A7/SP).
func (cpu *CPU) getRegA(reg uint8) uint32 {
	if reg == 7 {
		return cpu.sp
	}
	return cpu.A[reg]
}

// setRegA sets the value of address register An (A0-A6 or A7/SP).
func (cpu *CPU) setRegA(reg uint8, value uint32) {
	if reg == 7 {
		cpu.sp = value
	} else {
		cpu.A[reg] = value
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
