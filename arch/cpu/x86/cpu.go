package x86

// CPU represents x86 instruction processing capabilities for assembler/disassembler use.
// Maintains minimal register state for address calculations.
type CPU struct {
	// Segment registers for address calculations
	CS uint16 // code segment
	DS uint16 // data segment
	ES uint16 // extra segment
	SS uint16 // stack segment

	// General purpose registers for addressing modes
	AX uint16 // accumulator (AH:AL)
	BX uint16 // base register (BH:BL)
	CX uint16 // count register (CH:CL)
	DX uint16 // data register (DH:DL)

	// Index and pointer registers for addressing modes
	SI uint16 // source index
	DI uint16 // destination index
	BP uint16 // base pointer
	SP uint16 // stack pointer

	opts   Options
	memory *Memory
}

// New creates a new x86 instruction processor.
func New(memory *Memory, options ...Option) (*CPU, error) {
	if memory == nil {
		return nil, ErrNilMemory
	}

	opts := NewOptions(options...)

	c := &CPU{
		// Initialize with standard segment values
		CS: 0x1000, // Typical code segment
		DS: 0x1000, // Data segment
		ES: 0x1000, // Extra segment
		SS: 0x2000, // Stack segment
		SP: 0x1000, // Stack pointer

		opts:   opts,
		memory: memory,
	}

	return c, nil
}

// Memory returns the CPU memory.
func (c *CPU) Memory() *Memory {
	return c.memory
}

// CalculateAddress calculates the linear address from segment:offset.
func (c *CPU) CalculateAddress(segment, offset uint16) uint32 {
	return uint32(segment)<<4 + uint32(offset)
}
