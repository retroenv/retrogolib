package x86

// CPU represents x86 instruction processing capabilities.
type CPU struct {
	CS uint16 // code segment
	DS uint16 // data segment
	ES uint16 // extra segment
	SS uint16 // stack segment

	AX uint16 // accumulator (AH:AL)
	BX uint16 // base register (BH:BL)
	CX uint16 // count register (CH:CL)
	DX uint16 // data register (DH:DL)

	SI uint16 // source index
	DI uint16 // destination index
	BP uint16 // base pointer
	SP uint16 // stack pointer

	memory *Memory
}

// New creates a new x86 instruction processor.
func New(memory *Memory) (*CPU, error) {
	if memory == nil {
		return nil, ErrNilMemory
	}

	c := &CPU{
		CS: 0x1000,
		DS: 0x1000,
		ES: 0x1000,
		SS: 0x2000,
		SP: 0x1000,

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
