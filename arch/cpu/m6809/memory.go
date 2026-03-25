package m6809

import "errors"

// Interrupt vector addresses (all big-endian, 16-bit address space).
const (
	VectorReserved = uint16(0xFFF0)
	VectorSWI3     = uint16(0xFFF2)
	VectorSWI2     = uint16(0xFFF4)
	VectorFIRQ     = uint16(0xFFF6)
	VectorIRQ      = uint16(0xFFF8)
	VectorSWI      = uint16(0xFFFA)
	VectorNMI      = uint16(0xFFFC)
	VectorRESET    = uint16(0xFFFE)
)

// BasicMemory defines the interface required for 6809 memory access.
// All addresses are 16-bit. Byte order is big-endian (6809 native byte order).
type BasicMemory interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
	// ReadWord reads two bytes in big-endian order.
	ReadWord(address uint16) uint16
	// WriteWord writes two bytes in big-endian order.
	WriteWord(address uint16, value uint16)
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

// ReadVector reads a 16-bit interrupt vector from the given address.
func (m *Memory) ReadVector(address uint16) uint16 {
	return m.ReadWord(address)
}
