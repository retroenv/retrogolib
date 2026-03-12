package sm83

// Memory defines the interface for SM83 memory access.
// The SM83 has a 16-bit address bus (64KB addressable space).
type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, value uint8)
	ReadWord(addr uint16) uint16
	WriteWord(addr uint16, value uint16)
}

// BasicMemory provides a simple flat 64KB memory implementation.
type BasicMemory struct {
	data [65536]byte
}

// NewBasicMemory creates a new BasicMemory instance.
func NewBasicMemory() *BasicMemory {
	return &BasicMemory{}
}

// Read reads a byte from the specified address.
func (m *BasicMemory) Read(addr uint16) uint8 {
	return m.data[addr]
}

// Write writes a byte to the specified address.
func (m *BasicMemory) Write(addr uint16, value uint8) {
	m.data[addr] = value
}

// ReadWord reads a 16-bit word from the specified address (little-endian).
func (m *BasicMemory) ReadWord(addr uint16) uint16 {
	low := uint16(m.data[addr])
	high := uint16(m.data[addr+1])
	return high<<8 | low
}

// WriteWord writes a 16-bit word to the specified address (little-endian).
func (m *BasicMemory) WriteWord(addr uint16, value uint16) {
	m.data[addr] = uint8(value)
	m.data[addr+1] = uint8(value >> 8)
}

// LoadProgram loads program data into memory starting at the specified address.
func (m *BasicMemory) LoadProgram(addr uint16, data []byte) {
	copy(m.data[addr:], data)
}
