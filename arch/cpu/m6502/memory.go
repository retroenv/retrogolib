package m6502

import (
	"errors"
	"fmt"
	"math"

	. "github.com/retroenv/retrogolib/addressing"
)

const (
	IrqAddress   = 0xFFFE
	NMIAddress   = 0xFFFA
	ResetAddress = 0xFFFC

	StackBase = 0x100
)

// BasicMemory represents a basic memory access interface.
type BasicMemory interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
}

// Memory represents an advanced memory access interface.
type Memory struct {
	BasicMemory
}

// NewMemory returns a new memory instance.
func NewMemory(mem BasicMemory) *Memory {
	return &Memory{BasicMemory: mem}
}

// ReadWord reads a word from a memory address.
func (m *Memory) ReadWord(address uint16) uint16 {
	low := uint16(m.Read(address))
	high := uint16(m.Read(address + 1))
	w := (high << 8) | low
	return w
}

// ReadWordBug reads a word from a memory address
// and emulates a 6502 bug that caused the low byte to wrap
// without incrementing the high byte.
func (m *Memory) ReadWordBug(address uint16) uint16 {
	low := uint16(m.Read(address))
	offset := (address & 0xFF00) | uint16(byte(address)+1)
	high := uint16(m.Read(offset))
	w := (high << 8) | low
	return w
}

// WriteWord writes a word to a memory address.
func (m *Memory) WriteWord(address, value uint16) {
	m.Write(address, byte(value))
	m.Write(address+1, byte(value>>8))
}

// WriteAddressModes writes to memory using different address modes:
// Absolute: the absolut memory address is used to write the value
// Absolute, X: the absolut memory address with offset from X is used
// Absolute, Y: the absolut memory address with offset from Y is used
// (Indirect, X): the absolut memory address to write the value to is read from (indirect address + X)
// (Indirect), Y: the pointer to the memory address is read from the indirect parameter and adjusted after
// reading it by adding Y. The value is written to this pointer.
func (m *Memory) WriteAddressModes(value byte, params ...any) error {
	param := params[0]
	var register any
	if len(params) > 1 {
		register = params[1]
	}

	switch address := param.(type) {
	case int:
		return m.writeMemoryAbsolute(address, value, register)
	case *uint8: // variable
		*address = value
	case Absolute, AbsoluteX, AbsoluteY:
		return m.writeMemoryAbsolute(address, value, register)
	case ZeroPage:
		return m.writeMemoryZeroPage(address, value, register)
	case Indirect, IndirectResolved:
		return m.writeMemoryIndirect(address, value, register)
	default:
		return fmt.Errorf("unsupported memory write addressing mode type %T", param)
	}
	return nil
}

func (m *Memory) writeMemoryIndirect(address any, value byte, register any) error {
	pointer, err := m.indirectMemoryPointer(address, register)
	if err != nil {
		return err
	}
	m.Write(pointer, value)
	return nil
}

func (m *Memory) writeMemoryAbsolute(address any, value byte, register any) error {
	if register == nil {
		return m.writeMemoryAbsoluteOffset(address, value, 0)
	}

	var offset uint16
	switch val := register.(type) {
	case *uint8: // X/Y register referenced in normal code
		offset = uint16(*val)
	case uint8: // X/Y register referenced in unit test as system.X
		offset = uint16(val)
	default:
		return fmt.Errorf("unsupported extra parameter type %T for absolute memory write", register)
	}

	return m.writeMemoryAbsoluteOffset(address, value, offset)
}

// Support 6502 bug, index will not leave zeropage when page boundary is crossed.
func (m *Memory) writeMemoryZeroPage(address ZeroPage, value byte, register any) error {
	if register == nil {
		return m.writeMemoryAbsoluteOffset(address, value, 0)
	}

	var offset byte
	switch val := register.(type) {
	case *uint8: // X/Y register referenced in normal code
		offset = *val
	case uint8: // X/Y register referenced in unit test as system.X
		offset = val
	default:
		return fmt.Errorf("unsupported extra parameter type %T for zero page memory write", register)
	}

	addr := uint16(byte(address) + offset)
	return m.writeMemoryAbsoluteOffset(addr, value, 0)
}

func (m *Memory) writeMemoryAbsoluteOffset(address any, value byte, offset uint16) error {
	switch addr := address.(type) {
	case int8:
		m.Write(uint16(addr)+offset, value)
	case uint8:
		m.Write(uint16(addr)+offset, value)
	case *uint8:
		*addr = value
	case uint16:
		m.Write(addr+offset, value)
	case *uint16:
		*addr = uint16(value)
	case int:
		m.Write(uint16(addr)+offset, value)
	case Absolute:
		m.Write(uint16(addr)+offset, value)
	case AbsoluteX:
		m.Write(uint16(addr)+offset, value)
	case AbsoluteY:
		m.Write(uint16(addr)+offset, value)
	case ZeroPage:
		m.Write(uint16(addr)+offset, value)
	default:
		return fmt.Errorf("unsupported address type %T for absolute memory write with register", address)
	}
	return nil
}

// ReadAddressModes reads memory using different address modes:
// Immediate: if immediate is true and the passed first param fits into a byte, it's immediate value is returned
// Absolute: the absolut memory address is used to read the value
// Absolute, X: the absolut memory address with offset from X is used
// Absolute, Y: the absolut memory address with offset from Y is used
// (Indirect, X): the absolut memory address to write the value to is read from (indirect address + X)
// (Indirect), Y: the pointer to the memory address is read from the indirect parameter and adjusted after
// reading it by adding Y. The value is read from this pointer.
func (m *Memory) ReadAddressModes(immediate bool, params ...any) (byte, error) {
	param := params[0]
	var register any
	if len(params) > 1 {
		register = params[1]
	}

	switch address := param.(type) {
	case int:
		if immediate && register == nil && address <= math.MaxUint8 {
			return uint8(address), nil // immediate, not an address
		}
		return m.ReadAbsolute(address, register)
	case uint8:
		return address, nil // immediate, not an address
	case *uint8: // variable
		return *address, nil
	case Absolute, AbsoluteX, AbsoluteY:
		return m.ReadAbsolute(address, register)
	case ZeroPage:
		return m.ReadMemoryZeroPage(address, register)
	case Indirect, IndirectResolved:
		return m.readMemoryIndirect(address, register)
	default:
		return 0, fmt.Errorf("unsupported memory read addressing mode type %T", param)
	}
}

// ReadAbsolute reads a byte from an address using absolute addressing.
func (m *Memory) ReadAbsolute(address any, register any) (byte, error) {
	if register == nil {
		return m.readAbsoluteOffset(address, 0)
	}

	var offset uint16
	switch val := register.(type) {
	case *uint8: // X/Y register referenced in normal code
		offset = uint16(*val)
	case uint8: // X/Y register referenced in unit test as system.X
		offset = uint16(val)
	default:
		return 0, fmt.Errorf("unsupported extra parameter type %T for absolute memory read", register)
	}
	return m.readAbsoluteOffset(address, offset)
}

// ReadMemoryZeroPage reads a byte from an address in zeropage using absolute addressing.
// Support 6502 bug, index will not leave zeropage when page boundary is crossed.
func (m *Memory) ReadMemoryZeroPage(address ZeroPage, register any) (byte, error) {
	if register == nil {
		return m.readAbsoluteOffset(address, 0)
	}

	var offset byte
	switch val := register.(type) {
	case *uint8: // X/Y register referenced in normal code
		offset = *val
	case uint8: // X/Y register referenced in unit test as system.X
		offset = val
	default:
		return 0, fmt.Errorf("unsupported extra parameter type %T for zero page memory read", register)
	}
	addr := uint16(byte(address) + offset)
	return m.readAbsoluteOffset(addr, 0)
}

func (m *Memory) readAbsoluteOffset(address any, offset uint16) (byte, error) {
	switch addr := address.(type) {
	case *uint8:
		if offset != 0 {
			return 0, errors.New("memory pointer read with offset is not supported")
		}
		return *addr, nil
	case uint16:
		return m.Read(addr + offset), nil
	case int:
		return m.Read(uint16(addr) + offset), nil
	case Absolute:
		return m.Read(uint16(addr) + offset), nil
	case AbsoluteX:
		val := m.Read(uint16(addr))
		val += byte(offset)
		return val, nil
	case AbsoluteY:
		val := m.Read(uint16(addr))
		val += byte(offset)
		return val, nil
	case ZeroPage:
		return m.Read(uint16(addr) + offset), nil
	default:
		return 0, fmt.Errorf("unsupported address type %T for absolute memory write", address)
	}
}

func (m *Memory) readMemoryIndirect(address any, register any) (byte, error) {
	pointer, err := m.indirectMemoryPointer(address, register)
	if err != nil {
		return 0, err
	}
	return m.Read(pointer), nil
}

func (m *Memory) indirectMemoryPointer(addressParam any, register any) (uint16, error) {
	if register == nil {
		return 0, errors.New("register parameter missing for indirect memory addressing")
	}

	_, ok := register.(*uint8)
	if !ok {
		return 0, fmt.Errorf("unsupported extra parameter type %T for indirect memory addressing", register)
	}

	address := uint16(addressParam.(IndirectResolved))
	return address, nil
}
