package z80

// Memory defines the interface for Z80 memory access.
// Different hardware implementations can provide their own memory controllers
// by implementing this interface (e.g., Game Boy MBC, MSX, Spectrum).
type Memory interface {
	// Read reads a byte from memory at the given address.
	Read(address uint16) uint8

	// Write writes a byte to memory at the given address.
	Write(address uint16, value uint8)

	// ReadWord reads a 16-bit word from memory at the given address (little-endian).
	ReadWord(address uint16) uint16

	// WriteWord writes a 16-bit word to memory at the given address (little-endian).
	WriteWord(address uint16, value uint16)
}

// BasicMemory implements a simple 64KB flat memory space with no banking.
// This is suitable for basic Z80 systems without memory mappers.
type BasicMemory struct {
	data [0x10000]uint8
}

// NewBasicMemory creates a new basic memory controller with flat 64KB address space.
func NewBasicMemory() *BasicMemory {
	return &BasicMemory{}
}

// Read reads a byte from memory at the given address.
func (mem *BasicMemory) Read(address uint16) uint8 {
	return mem.data[address]
}

// Write writes a byte to memory at the given address.
func (mem *BasicMemory) Write(address uint16, value uint8) {
	mem.data[address] = value
}

// ReadWord reads a 16-bit word from memory at the given address (little-endian).
func (mem *BasicMemory) ReadWord(address uint16) uint16 {
	low := uint16(mem.Read(address))
	high := uint16(mem.Read(address + 1))
	return high<<8 | low
}

// WriteWord writes a 16-bit word to memory at the given address (little-endian).
func (mem *BasicMemory) WriteWord(address uint16, value uint16) {
	mem.Write(address, uint8(value))
	mem.Write(address+1, uint8(value>>8))
}

// LoadROM loads ROM data into memory starting at address 0.
func (mem *BasicMemory) LoadROM(data []byte) {
	if data == nil {
		return
	}

	n := min(len(data), len(mem.data))
	if n > 0 {
		copy(mem.data[:n], data[:n])
	}
}

// LoadProgram loads program data into memory starting at address 0.
// This is an alias for LoadROM for backward compatibility.
func (mem *BasicMemory) LoadProgram(data []byte) {
	mem.LoadROM(data)
}

// Data returns a reference to the underlying memory array.
// This is useful for direct memory access in testing or debugging.
func (mem *BasicMemory) Data() *[0x10000]uint8 {
	return &mem.data
}
