package sm83

import "fmt"

// readOpParams reads the instruction operand bytes based on addressing mode.
// Returns the typed parameter and raw operand bytes.
func readOpParams(c *CPU, addressing AddressingMode) ([]any, []byte, error) {
	switch addressing {
	case ImpliedAddressing, RegisterAddressing:
		return nil, nil, nil

	case ImmediateAddressing:
		value := c.memory.Read(c.PC + 1)
		return []any{Immediate8(value)}, []byte{value}, nil

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
