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

// Bus provides the full hardware interface for a Z80 system.
// It extends Memory with I/O port access and interrupt acknowledgment.
// For simple use cases (tests, basic emulation), use Memory + WithIOHandler instead.
type Bus interface {
	Memory

	// ReadPort reads from an I/O port. The full 16-bit address is provided
	// because the Z80 places register data on the upper address lines:
	//   - IN A,(n):    address = A<<8 | n
	//   - IN r,(C):    address = B<<8 | C
	//   - INI/IND/etc: address = B<<8 | C
	ReadPort(address uint16) uint8

	// WritePort writes to an I/O port with full 16-bit address.
	WritePort(address uint16, value uint8)

	// IRQData returns the byte placed on the data bus during interrupt acknowledge.
	// For IM 0, this should be an instruction opcode (typically RST n, e.g. 0xFF for RST 38h).
	// For IM 2, this is the low byte of the interrupt vector table address.
	IRQData() uint8

	// OnRETI is called when a RETI instruction executes.
	// Hardware (e.g., Z80 PIO/CTC daisy chain) monitors the bus for RETI
	// to manage interrupt priority.
	OnRETI()
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

// legacyBusAdapter wraps Memory + IOHandler into a Bus for backward compatibility.
type legacyBusAdapter struct {
	Memory
	ioHandler IOHandler
}

func (a *legacyBusAdapter) ReadPort(address uint16) uint8 {
	if a.ioHandler != nil {
		return a.ioHandler.ReadPort(uint8(address))
	}
	return 0xFF
}

func (a *legacyBusAdapter) WritePort(address uint16, value uint8) {
	if a.ioHandler != nil {
		a.ioHandler.WritePort(uint8(address), value)
	}
}

func (a *legacyBusAdapter) IRQData() uint8 { return 0xFF }
func (a *legacyBusAdapter) OnRETI()        {}
