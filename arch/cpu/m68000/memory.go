package m68000

// Memory defines the interface for 68000 memory access.
// The 68000 is big-endian: ReadWord reads [addr] as high byte, [addr+1] as low byte.
// Word and long accesses at odd addresses trigger an address error exception.
type Memory interface {
	// Read reads a byte from memory at the given address.
	Read(address uint32) uint8

	// ReadLong reads a 32-bit long word from memory at the given address (big-endian).
	ReadLong(address uint32) uint32

	// ReadWord reads a 16-bit word from memory at the given address (big-endian).
	ReadWord(address uint32) uint16

	// Write writes a byte to memory at the given address.
	Write(address uint32, value uint8)

	// WriteLong writes a 32-bit long word to memory at the given address (big-endian).
	WriteLong(address uint32, value uint32)

	// WriteWord writes a 16-bit word to memory at the given address (big-endian).
	WriteWord(address uint32, value uint16)
}

// Bus provides the full hardware interface for a 68000 system.
// It extends Memory with interrupt handling and reset notification.
type Bus interface {
	Memory

	// IRQAcknowledge is called when the CPU acknowledges an interrupt at the given level.
	// Returns the vector number for the interrupt.
	IRQAcknowledge(level uint8) uint32

	// IRQLevel returns the current interrupt priority level (0-7).
	IRQLevel() uint8

	// OnReset is called when the CPU executes the RESET instruction.
	OnReset()
}

// addressMask masks a 32-bit address to 24 bits for the 68000 address bus.
const addressMask = 0x00FFFFFF

// BasicMemory implements a simple 16MB flat memory space for the 68000.
type BasicMemory struct {
	data [0x1000000]uint8 // 16MB (24-bit address space)
}

// NewBasicMemory creates a new basic memory controller with flat 16MB address space.
func NewBasicMemory() *BasicMemory {
	return &BasicMemory{}
}

// Read reads a byte from memory at the given address.
func (mem *BasicMemory) Read(address uint32) uint8 {
	return mem.data[address&addressMask]
}

// ReadWord reads a 16-bit word from memory at the given address (big-endian).
func (mem *BasicMemory) ReadWord(address uint32) uint16 {
	addr := address & addressMask
	return uint16(mem.data[addr])<<8 | uint16(mem.data[addr+1])
}

// ReadLong reads a 32-bit long word from memory at the given address (big-endian).
func (mem *BasicMemory) ReadLong(address uint32) uint32 {
	addr := address & addressMask
	return uint32(mem.data[addr])<<24 |
		uint32(mem.data[addr+1])<<16 |
		uint32(mem.data[addr+2])<<8 |
		uint32(mem.data[addr+3])
}

// Write writes a byte to memory at the given address.
func (mem *BasicMemory) Write(address uint32, value uint8) {
	mem.data[address&addressMask] = value
}

// WriteWord writes a 16-bit word to memory at the given address (big-endian).
func (mem *BasicMemory) WriteWord(address uint32, value uint16) {
	addr := address & addressMask
	mem.data[addr] = uint8(value >> 8)
	mem.data[addr+1] = uint8(value)
}

// WriteLong writes a 32-bit long word to memory at the given address (big-endian).
func (mem *BasicMemory) WriteLong(address uint32, value uint32) {
	addr := address & addressMask
	mem.data[addr] = uint8(value >> 24)
	mem.data[addr+1] = uint8(value >> 16)
	mem.data[addr+2] = uint8(value >> 8)
	mem.data[addr+3] = uint8(value)
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

// Data returns a reference to the underlying memory array.
func (mem *BasicMemory) Data() *[0x1000000]uint8 {
	return &mem.data
}

// BasicBus wraps BasicMemory into a Bus implementation for simple use cases.
type BasicBus struct {
	Memory
	irqLevel uint8
}

// NewBasicBus creates a new basic bus wrapping the given memory.
func NewBasicBus(mem Memory) *BasicBus {
	return &BasicBus{Memory: mem}
}

// IRQAcknowledge acknowledges an interrupt and returns the autovector number.
func (b *BasicBus) IRQAcknowledge(level uint8) uint32 {
	return uint32(VectorAutoVector1) + uint32(level) - 1
}

// IRQLevel returns the current interrupt priority level.
func (b *BasicBus) IRQLevel() uint8 {
	return b.irqLevel
}

// OnReset handles the RESET instruction.
func (b *BasicBus) OnReset() {}
