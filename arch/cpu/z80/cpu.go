package z80

import (
	"sync"

	"github.com/retroenv/retrogolib/arch"
)

// State is a register and execution-status snapshot for debugging.
// It does not include every internal latch needed for save/restore.
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
	SP     uint16 // Stack pointer
	PC     uint16 // Program counter
	I      uint8  // Interrupt vector base
	R      uint8  // Memory refresh counter
	MEMPTR uint16 // Internal WZ register, its high byte provides the X/Y flag bits for BIT n,(HL) and indexed BIT instructions

	Cycles     uint64
	Flags      Flags
	AltFlags   Flags
	Interrupts Interrupts

	Halted bool
}

// CPU emulates a Z80 microprocessor.
// Step, CheckInterrupts, and interrupt triggers are synchronized. Callers must
// serialize direct register access and interrupt configuration with execution.
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
	SP     uint16 // Stack pointer
	PC     uint16 // Program counter
	I      uint8  // Interrupt vector base register
	R      uint8  // Memory refresh register (auto-incremented)
	MEMPTR uint16 // Internal WZ register, its high byte provides the X/Y flag bits for BIT n,(HL) and indexed BIT instructions

	Flags    Flags // Main flag register
	AltFlags Flags // Shadow flag register

	cycles uint64
	halted bool

	// Interrupt control
	iff1 bool          // Interrupt enable flip-flop
	iff2 bool          // Backup of IFF1 for NMI handling
	im   InterruptMode // Interrupt mode: 0, 1, or 2

	triggerIrq bool
	triggerNmi bool

	opts      options
	TraceStep TraceStep // trace step info, set if tracing is enabled

	currentOpcode uint8 // opcode being executed (for instruction functions to access)

	// Q register: tracks previous flag state so SCF/CCF can set the
	// undocumented X/Y flag bits (bits 3 and 5) correctly.
	// After each instruction, Q captures the flags byte. SCF/CCF then
	// compute X/Y as: (A | (F & ~Q)) & 0x28.
	q uint8

	// lastWasLdAIR tracks if the previous instruction was LD A,I or LD A,R.
	// Used to emulate the Zilog NMOS bug where P/V is reset if an IRQ fires
	// during these instructions.
	lastWasLdAIR bool
	eiPending    bool // IRQ acceptance waits until the instruction after EI completes.

	bus Bus
}

// Interrupts holds the current interrupt state.
type Interrupts struct {
	IFF1         bool
	IFF2         bool
	IM           uint8
	NMITriggered bool
	IrqTriggered bool
}

// New creates a new Z80 CPU with a memory controller.
// For simple use cases where full bus emulation is not needed.
// I/O ports can optionally be handled via WithIOHandler.
func New(memory Memory, options ...Option) (*CPU, error) {
	if memory == nil {
		return nil, ErrNilMemory
	}

	opts := newOptions(options...)
	bus := &legacyBusAdapter{Memory: memory, ioHandler: opts.ioHandler}
	return newCPU(bus, opts)
}

// NewWithBus creates a new Z80 CPU with a full bus interface.
// The Bus interface provides memory, I/O ports, and interrupt acknowledgment,
// enabling accurate emulation of systems with interrupt daisy chains and
// full 16-bit I/O port addressing. WithIOHandler applies only to New.
// Bus callbacks execute under the CPU lock and must not call locking CPU methods.
func NewWithBus(bus Bus, options ...Option) (*CPU, error) {
	if bus == nil {
		return nil, ErrNilMemory
	}

	opts := newOptions(options...)
	return newCPU(bus, opts)
}

// Cycles returns total CPU cycles executed.
func (cpu *CPU) Cycles() uint64 {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.cycles
}

// Halted returns CPU halt state.
func (cpu *CPU) Halted() bool {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.halted
}

// Halt stops CPU execution.
func (cpu *CPU) Halt() {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	cpu.halted = true
}

// Resume continues CPU execution.
func (cpu *CPU) Resume() {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	cpu.halted = false
}

// State returns a register and execution-status snapshot.
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
		MEMPTR:   cpu.MEMPTR,
		Cycles:   cpu.cycles,
		Flags:    cpu.Flags,
		AltFlags: cpu.AltFlags,
		Interrupts: Interrupts{
			IFF1:         cpu.iff1,
			IFF2:         cpu.iff2,
			IM:           uint8(cpu.im),
			NMITriggered: cpu.triggerNmi,
			IrqTriggered: cpu.triggerIrq,
		},
		Halted: cpu.halted,
	}
}

// Memory returns the attached memory controller.
//
//nolint:ireturn // intentional: Memory is the public API interface
func (cpu *CPU) Memory() Memory {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.bus
}

// Bus returns the attached bus interface.
//
//nolint:ireturn // intentional: Bus is the public API interface
func (cpu *CPU) Bus() Bus {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.bus
}

// BC returns the BC register pair as a 16-bit value.
func (cpu *CPU) BC() uint16 {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.bc()
}

// DE returns the DE register pair as a 16-bit value.
func (cpu *CPU) DE() uint16 {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.de()
}

// HL returns the HL register pair as a 16-bit value.
func (cpu *CPU) HL() uint16 {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.hl()
}

// AF returns the AF register pair as a 16-bit value.
func (cpu *CPU) AF() uint16 {
	cpu.mu.RLock()
	defer cpu.mu.RUnlock()
	return cpu.af()
}

// TriggerNMI triggers a non-maskable interrupt.
func (cpu *CPU) TriggerNMI() {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	cpu.triggerNmi = true
}

// TriggerIRQ triggers a maskable interrupt.
func (cpu *CPU) TriggerIRQ() {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	cpu.triggerIrq = true
}

// bc returns the BC register pair as a 16-bit value (internal, no lock).
func (cpu *CPU) bc() uint16 {
	return uint16(cpu.B)<<8 | uint16(cpu.C)
}

// de returns the DE register pair as a 16-bit value (internal, no lock).
func (cpu *CPU) de() uint16 {
	return uint16(cpu.D)<<8 | uint16(cpu.E)
}

// hl returns the HL register pair as a 16-bit value (internal, no lock).
func (cpu *CPU) hl() uint16 {
	return uint16(cpu.H)<<8 | uint16(cpu.L)
}

// af returns the AF register pair as a 16-bit value (internal, no lock).
func (cpu *CPU) af() uint16 {
	return uint16(cpu.A)<<8 | uint16(cpu.GetFlags())
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

// pop pops a byte from the stack and updates the stack pointer.
func (cpu *CPU) pop() uint8 {
	value := cpu.bus.Read(cpu.SP)
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
	cpu.bus.Write(cpu.SP, value)
}

// push16 pushes a word to the stack and updates the stack pointer.
func (cpu *CPU) push16(value uint16) {
	high := uint8(value >> 8)
	low := uint8(value)
	cpu.push(high)
	cpu.push(low)
}

// inPortToRegister reads from port C to a register and sets flags.
func (cpu *CPU) inPortToRegister(regPtr *uint8) {
	cpu.MEMPTR = cpu.bc() + 1
	value := cpu.bus.ReadPort(cpu.bc())
	*regPtr = value
	cpu.setSZP(value)
	cpu.setH(false)
	cpu.setN(false)
}

// applyCBOperation applies a CB prefix operation to a register or (HL).
// Used by rotate, shift, RES, and SET instructions.
func (cpu *CPU) applyCBOperation(operation func(uint8) uint8) {
	opcodeByte := cpu.bus.Read(cpu.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // Operation on (HL)
		addr := cpu.hl()
		value := cpu.bus.Read(addr)
		result := operation(value)
		cpu.bus.Write(addr, result)
	} else { // Operation on register
		value := cpu.GetRegisterValue(reg)
		result := operation(value)
		cpu.SetRegisterValue(reg, result)
	}
}

// calculateIndexedAddress reads displacement byte from memory at PC+2 and calculates indexed address.
// Used by DD (IX) and FD (IY) prefix instructions where PC points to the prefix byte.
// Also sets MEMPTR to the calculated address.
func (cpu *CPU) calculateIndexedAddress(indexReg uint16, _ ...any) uint16 {
	displacement := int8(cpu.bus.Read(cpu.PC + 2))
	addr := uint16(int32(indexReg) + int32(displacement))
	cpu.MEMPTR = addr
	return addr
}

// read16 reads a 16-bit value from memory at addr (little-endian).
func (cpu *CPU) read16(addr uint16) uint16 {
	low := cpu.bus.Read(addr)
	high := cpu.bus.Read(addr + 1)
	return uint16(high)<<8 | uint16(low)
}

// writeRegisterPair writes a register pair to memory at addr (little-endian).
func (cpu *CPU) writeRegisterPair(addr uint16, low, high uint8) {
	cpu.bus.Write(addr, low)
	cpu.bus.Write(addr+1, high)
}

// setLogicalFlags sets flags for logical operations (AND/OR/XOR).
// hFlag should be true for AND, false for OR/XOR.
func (cpu *CPU) setLogicalFlags(result uint8, hFlag bool) {
	cpu.setSZP(result)
	cpu.setH(hFlag)
	cpu.setN(false)
	cpu.setC(false)
}

func newCPU(bus Bus, opts options) (*CPU, error) {
	// Default to generic system
	if opts.initialPC == 0 && opts.initialSP == 0 && opts.systemType == "" {
		opts.systemType = arch.Generic
		opts.initialPC = 0x0000
		opts.initialSP = 0xFFFF
	}

	c := &CPU{
		PC:   opts.initialPC,
		SP:   opts.initialSP,
		opts: opts,
		bus:  bus,
	}

	return c, nil
}
