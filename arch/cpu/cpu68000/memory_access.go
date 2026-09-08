package cpu68000

import "fmt"

// BusErrorHandler is an optional Bus capability for rejecting a memory transfer.
// OnBusAccess runs before the transfer, with a 24-bit address and a byte or word
// size. Long accesses are two word transfers, so a second-word fault preserves
// the completed first transfer. A non-nil error raises vector 2 in Step.
// Callbacks run under the CPU lock and must not call locking CPU methods.
type BusErrorHandler interface {
	OnBusAccess(address uint32, size OperandSize, write bool) error
}

func (cpu *CPU) checkAccess(address uint32, size OperandSize, write bool, space accessSpace) {
	if size != SizeByte && address&1 != 0 {
		cpu.raiseAccessFault(address, write, space, VectorAddressErr, ErrAddressError)
	}
	if handler, ok := cpu.bus.(BusErrorHandler); ok {
		if err := handler.OnBusAccess(address&addressMask, size, write); err != nil {
			cpu.raiseAccessFault(address, write, space, VectorBusError, fmt.Errorf("%w: %w", ErrBusError, err))
		}
	}
}

func (cpu *CPU) raiseAccessFault(address uint32, write bool, space accessSpace, vector int, cause error) {
	// Undefined upper SSW bits retain the instruction register on the 68000.
	status := cpu.instructionWord & 0xFFE0
	status |= 1
	if space == programSpace {
		status = status&^7 | 2
	}
	if cpu.IsSupervisor() {
		status |= 4
	}
	if !write {
		status |= 0x10
	}
	if cpu.exceptionAccess {
		status |= 8
	}
	pc := cpu.PC
	if space == dataSpace && !cpu.exceptionAccess {
		// Operand faults expose the instruction pipeline's current word address.
		pc = uint32(int32(pc) + cpu.operandPCOffset)
	}
	panic(&accessError{
		cause: cause, address: address, pc: pc, status: status,
		sr: cpu.GetSR(), word: cpu.instructionWord, vector: vector, cycles: cpu.accessCycles,
	})
}

func (cpu *CPU) readByte(address uint32) uint8 {
	cpu.checkAccess(address, SizeByte, false, dataSpace)
	cpu.accessCycles += 4
	return cpu.bus.Read(address & addressMask)
}

func (cpu *CPU) readBusWord(address uint32, space accessSpace) uint16 {
	cpu.checkAccess(address, SizeWord, false, space)
	cpu.accessCycles += 4
	return cpu.bus.ReadWord(address & addressMask)
}

func (cpu *CPU) readBusLong(address uint32) uint32 {
	high := uint32(cpu.readBusWord(address, dataSpace))
	return high<<16 | uint32(cpu.readBusWord(address+2, dataSpace))
}

func (cpu *CPU) writeByte(address uint32, value uint8) {
	cpu.checkAccess(address, SizeByte, true, dataSpace)
	cpu.accessCycles += 4
	cpu.bus.Write(address&addressMask, value)
}

func (cpu *CPU) writeBusWord(address uint32, value uint16) {
	cpu.checkAccess(address, SizeWord, true, dataSpace)
	cpu.accessCycles += 4
	cpu.bus.WriteWord(address&addressMask, value)
}

func (cpu *CPU) writeBusLong(address, value uint32) {
	cpu.writeBusWord(address, uint16(value>>16))
	cpu.writeBusWord(address+2, uint16(value))
}

func (cpu *CPU) readPredecrement(reg uint8, size OperandSize) (uint32, error) {
	if size != SizeLong {
		cpu.setRegA(reg, cpu.getRegA(reg)-incrementSize(reg, size))
		return cpu.readMemory(cpu.getRegA(reg), size)
	}
	// Multiprecision long operands decrement and transfer one word at a time.
	cpu.setRegA(reg, cpu.getRegA(reg)-2)
	low := uint32(cpu.readBusWord(cpu.getRegA(reg), dataSpace))
	cpu.setRegA(reg, cpu.getRegA(reg)-2)
	high := uint32(cpu.readBusWord(cpu.getRegA(reg), dataSpace))
	return high<<16 | low, nil
}

type accessSpace uint8

const (
	dataSpace accessSpace = iota
	programSpace
)

// accessError unwinds the current instruction before any later side effects.
// Only this private panic type is intercepted at CPU execution boundaries.
type accessError struct {
	cause   error
	address uint32
	pc      uint32
	status  uint16
	sr      uint16
	word    uint16
	vector  int
	cycles  uint64
}

func (flt *accessError) Error() string {
	return fmt.Sprintf("%v at address %08x", flt.cause, flt.address)
}

func (flt *accessError) Unwrap() error {
	return flt.cause
}

func catchAccessFault(action func() error) (err error) {
	defer func() {
		if value := recover(); value != nil {
			fault, ok := value.(*accessError)
			if !ok {
				panic(value)
			}
			err = fault
		}
	}()
	return action()
}
