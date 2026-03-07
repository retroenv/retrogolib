package m65816

import "errors"

// Interrupt vector addresses (all in bank $00).
const (
	// Native mode vectors
	VectorNativeCOP   = uint32(0xFFE4)
	VectorNativeBRK   = uint32(0xFFE6)
	VectorNativeABORT = uint32(0xFFE8)
	VectorNativeNMI   = uint32(0xFFEA)
	VectorNativeIRQ   = uint32(0xFFEE)

	// Emulation mode vectors
	VectorEmuCOP   = uint32(0xFFF4)
	VectorEmuABORT = uint32(0xFFF8)
	VectorEmuNMI   = uint32(0xFFFA)
	VectorEmuRESET = uint32(0xFFFC)
	VectorEmuIRQ   = uint32(0xFFFE) // also BRK in emulation mode
)

// BasicMemory defines the interface required for 65816 memory access.
// All addresses are 24-bit (uint32 masked to low 24 bits).
// Byte order is little-endian (65816 native byte order).
type BasicMemory interface {
	Read(address uint32) uint8
	Write(address uint32, value uint8)
	// ReadWord reads two bytes in little-endian order.
	ReadWord(address uint32) uint16
	// WriteWord writes two bytes in little-endian order.
	WriteWord(address uint32, value uint16)
}

// Memory wraps BasicMemory with higher-level helper methods.
type Memory struct {
	BasicMemory
}

// NewMemory creates a new Memory wrapper.
func NewMemory(mem BasicMemory) (*Memory, error) {
	if mem == nil {
		return nil, errors.New("BasicMemory cannot be nil")
	}
	return &Memory{BasicMemory: mem}, nil
}

// ReadLong reads a 24-bit (3-byte) value in little-endian order.
func (m *Memory) ReadLong(address uint32) uint32 {
	lo := uint32(m.Read(address))
	mid := uint32(m.Read(address + 1))
	hi := uint32(m.Read(address + 2))
	return hi<<16 | mid<<8 | lo
}

// WriteLong writes a 24-bit (3-byte) value in little-endian order.
func (m *Memory) WriteLong(address uint32, value uint32) {
	m.Write(address, uint8(value))
	m.Write(address+1, uint8(value>>8))
	m.Write(address+2, uint8(value>>16))
}

// ReadVector reads a 16-bit interrupt vector from the given address.
func (m *Memory) ReadVector(address uint32) uint16 {
	return m.ReadWord(address)
}

// bank24 forms a 24-bit address from an 8-bit bank and a 16-bit offset.
func bank24(bank uint8, offset uint16) uint32 {
	return uint32(bank)<<16 | uint32(offset)
}
