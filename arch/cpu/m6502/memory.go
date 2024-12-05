package m6502

const (
	StackBase = 0x100
)

// BasicMemory represents a basic memory access interface.
type BasicMemory interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
}

// Memory represents an advanced memory access interface.
type Memory interface {
	BasicMemory

	ReadAbsolute(address any, register any) byte
	ReadAddressModes(immediate bool, params ...any) byte
	ReadWord(address uint16) uint16
	WriteAddressModes(value byte, params ...any)
}

// pop pops a byte from the stack and update the stack pointer.
func (c *CPU) pop() byte {
	c.SP++
	return c.memory.Read(uint16(StackBase + int(c.SP)))
}

// pop16 pops a word from the stack and updates the stack pointer.
func (c *CPU) pop16() uint16 {
	low := uint16(c.pop())
	high := uint16(c.pop())
	return high<<8 | low
}

// push a value to the stack and update the stack pointer.
func (c *CPU) push(value byte) {
	c.memory.Write(uint16(StackBase+int(c.SP)), value)
	c.SP--
}

// push16 a word to the stack and update the stack pointer.
func (c *CPU) push16(value uint16) {
	high := byte(value >> 8)
	low := byte(value)
	c.push(high)
	c.push(low)
}
