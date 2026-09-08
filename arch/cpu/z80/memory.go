package z80

// Memory defines the interface for Z80 memory access.
// Different hardware implementations can provide their own memory controllers
// by implementing this interface (e.g., Game Boy MBC, MSX, Spectrum).
type Memory interface {
	// Read reads a byte from memory at the given address.
	Read(address uint16) uint8

	// ReadWord reads a 16-bit word from memory at the given address (little-endian).
	ReadWord(address uint16) uint16

	// Write writes a byte to memory at the given address.
	Write(address uint16, value uint8)

	// WriteWord writes a 16-bit word to memory at the given address (little-endian).
	WriteWord(address uint16, value uint16)
}

// Bus extends Memory with 16-bit port addresses and interrupt acknowledgment.
// NewWithBus uses this interface directly. New adapts Memory and WithIOHandler.
// Callbacks run under the CPU lock and must not call locking CPU methods.
type Bus interface {
	Memory

	// IRQData supplies an IM 0 opcode or the low byte of an IM 2 vector address.
	// It is sampled once per accepted IRQ in these modes; IM 1 ignores bus data.
	// IM 0 supports RST opcodes, with a RST 38h fallback for other values.
	IRQData() uint8

	// OnRETI notifies interrupt daisy-chain hardware after ED 4D executes.
	// RETN and its undocumented mirrors do not notify the host.
	OnRETI()

	// ReadPort receives A:n for immediate IN or B:C for register/block IN.
	// Block input uses B before its decrement.
	ReadPort(address uint16) uint8

	// WritePort receives A:n for immediate OUT or B:C for register/block OUT.
	// Block output uses B after its decrement.
	WritePort(address uint16, value uint8)
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

func (bus *legacyBusAdapter) ReadPort(address uint16) uint8 {
	if bus.ioHandler != nil {
		return bus.ioHandler.ReadPort(uint8(address))
	}
	return 0xFF
}

func (bus *legacyBusAdapter) WritePort(address uint16, value uint8) {
	if bus.ioHandler != nil {
		bus.ioHandler.WritePort(uint8(address), value)
	}
}

func (bus *legacyBusAdapter) IRQData() uint8 { return 0xFF }
func (bus *legacyBusAdapter) OnRETI()        {}
