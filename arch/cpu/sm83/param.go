package sm83

import "fmt"

// readOpParams reads the instruction operand bytes based on addressing mode.
// Returns the typed parameter and raw operand bytes.
func readOpParams(c *CPU, addressing AddressingMode) ([]any, []byte, error) {
	switch addressing {
	case ImpliedAddressing, RegisterAddressing:
		return nil, nil, nil

	case ImmediateAddressing:
		return readImmediateParam(c)

	case ExtendedAddressing:
		low := c.memory.Read(c.PC + 1)
		high := c.memory.Read(c.PC + 2)
		addr := uint16(high)<<8 | uint16(low)
		return []any{Extended(addr)}, []byte{low, high}, nil

	case RegisterIndirectAddressing:
		return nil, nil, nil

	case RelativeAddressing:
		offset := c.memory.Read(c.PC + 1)
		return []any{Relative(int8(offset))}, []byte{offset}, nil

	case BitAddressing:
		return nil, nil, nil

	default:
		return nil, nil, fmt.Errorf("%w: %d", ErrUnsupportedAddressingMode, addressing)
	}
}

// readImmediateParam reads an 8-bit or 16-bit immediate value based on instruction size.
func readImmediateParam(c *CPU) ([]any, []byte, error) {
	opcode := c.memory.Read(c.PC)
	opcodeInfo := Opcodes[opcode]

	if opcodeInfo.Size == 3 {
		// 16-bit immediate (3-byte instruction: opcode + low byte + high byte)
		low := c.memory.Read(c.PC + 1)
		high := c.memory.Read(c.PC + 2)
		value := uint16(high)<<8 | uint16(low)
		return []any{Immediate16(value)}, []byte{low, high}, nil
	}

	// 8-bit immediate (2-byte instruction: opcode + value)
	value := c.memory.Read(c.PC + 1)
	return []any{Immediate8(value)}, []byte{value}, nil
}
