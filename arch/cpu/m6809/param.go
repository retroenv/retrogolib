package m6809

import "fmt"

// readOpParams reads the instruction operand bytes for the given addressing mode.
// Returns the decoded params, raw operand bytes, and error.
func readOpParams(c *CPU, mode AddressingMode, baseOffset uint16) ([]any, []byte, error) {
	switch mode {
	case ImpliedAddressing:
		return nil, nil, nil

	case ImmediateAddressing:
		b := c.fetchByte(baseOffset)
		return []any{Immediate8(b)}, []byte{b}, nil

	case Immediate16Addressing:
		hi := c.fetchByte(baseOffset)
		lo := c.fetchByte(baseOffset + 1)
		val := uint16(hi)<<8 | uint16(lo)
		return []any{Immediate16(val)}, []byte{hi, lo}, nil

	case DirectAddressing:
		b := c.fetchByte(baseOffset)
		return []any{DirectPage(b)}, []byte{b}, nil

	case ExtendedAddressing:
		hi := c.fetchByte(baseOffset)
		lo := c.fetchByte(baseOffset + 1)
		addr := uint16(hi)<<8 | uint16(lo)
		return []any{Extended16(addr)}, []byte{hi, lo}, nil

	case IndexedAddressing:
		return readIndexedParam(c, baseOffset)

	case RelativeAddressing:
		offset := int8(c.fetchByte(baseOffset))
		// Branch target: PC + instruction size (2) + signed offset
		target := uint16(int32(c.PC) + 2 + int32(offset))
		return []any{target}, []byte{uint8(offset)}, nil

	case RelativeLongAddressing:
		return readRelativeLongParam(c, baseOffset)

	case RegisterAddressing:
		b := c.fetchByte(baseOffset)
		return []any{RegisterPair(b)}, []byte{b}, nil

	case StackAddressing:
		b := c.fetchByte(baseOffset)
		return []any{StackMask(b)}, []byte{b}, nil

	default:
		return nil, nil, fmt.Errorf("%w: mode 0x%x", ErrUnsupportedAddressingMode, mode)
	}
}

// readRelativeLongParam reads a 16-bit relative offset.
// The base offset varies: base page instructions (LBRA $16, LBSR $17) use offset 1;
// page 2 prefixed instructions use offset 2 (prefix + opcode already consumed).
func readRelativeLongParam(c *CPU, baseOffset uint16) ([]any, []byte, error) {
	hi := c.fetchByte(baseOffset)
	lo := c.fetchByte(baseOffset + 1)
	offset := int16(uint16(hi)<<8 | uint16(lo))

	// instrSize = baseOffset + 2 (for the offset bytes themselves)
	instrSize := int32(baseOffset) + 2
	target := uint16(int32(c.PC) + instrSize + int32(offset))

	return []any{target}, []byte{hi, lo}, nil
}

// readIndexedParam decodes the 6809 indexed addressing postbyte.
// Returns the decoded param, raw bytes, and error.
func readIndexedParam(c *CPU, baseOffset uint16) ([]any, []byte, error) {
	postbyte := c.fetchByte(baseOffset)
	operands := make([]byte, 1, 4)
	operands[0] = postbyte

	// Get the register from bits 5-6
	reg := c.indexedRegister(postbyte)

	// Bit 7 = 0: 5-bit constant offset
	if postbyte&0x80 == 0 {
		offset := int8(postbyte & 0x1F)
		// Sign-extend the 5-bit value
		if offset&0x10 != 0 {
			offset |= -0x20 // sign extend from 5 bits
		}
		addr := uint16(int32(reg) + int32(offset))
		return []any{IndexedAddr(addr)}, operands, nil
	}

	// Mode 0x0F with indirect is "extended indirect [n]" which is self-contained.
	// All other indirect modes need the indirection applied after address calculation.
	mode := postbyte & 0x0F
	indirect := postbyte&0x10 != 0

	if mode == 0x0F && indirect {
		return readExtendedIndirect(c, baseOffset, operands)
	}

	addr, extra, err := c.decodeIndexedMode(postbyte, baseOffset, reg)
	if err != nil {
		return nil, nil, err
	}
	operands = append(operands, extra...)

	// Handle indirection (bit 4 set)
	if indirect {
		addr = c.memory.ReadWord(addr)
	}

	return []any{IndexedAddr(addr)}, operands, nil
}

// readExtendedIndirect handles the extended indirect [n] addressing mode ($9F postbyte).
func readExtendedIndirect(c *CPU, baseOffset uint16, operands []byte) ([]any, []byte, error) {
	hi := c.fetchByte(baseOffset + 1)
	lo := c.fetchByte(baseOffset + 2)
	operands = append(operands, hi, lo)
	addr := uint16(hi)<<8 | uint16(lo)
	addr = c.memory.ReadWord(addr)
	return []any{IndexedAddr(addr)}, operands, nil
}

// decodeIndexedMode decodes the complex indexed addressing modes (bit 7 = 1).
// Returns the effective address, any extra operand bytes, and error.
// The extended indirect mode (0x0F) is handled separately by readExtendedIndirect.
func (c *CPU) decodeIndexedMode(postbyte uint8, baseOffset uint16, reg uint16) (uint16, []byte, error) { //nolint:cyclop,funlen
	mode := postbyte & 0x0F

	switch mode {
	case 0x00: // ,R+
		c.setIndexedRegister(postbyte, reg+1)
		return reg, nil, nil

	case 0x01: // ,R++
		c.setIndexedRegister(postbyte, reg+2)
		return reg, nil, nil

	case 0x02: // ,-R
		c.setIndexedRegister(postbyte, reg-1)
		return c.indexedRegister(postbyte), nil, nil

	case 0x03: // ,--R
		c.setIndexedRegister(postbyte, reg-2)
		return c.indexedRegister(postbyte), nil, nil

	case 0x04: // ,R (no offset)
		return reg, nil, nil

	case 0x05: // B,R
		return uint16(int32(reg) + int32(int8(c.B))), nil, nil

	case 0x06: // A,R
		return uint16(int32(reg) + int32(int8(c.A))), nil, nil

	case 0x08: // n,R (8-bit offset)
		b := c.fetchByte(baseOffset + 1)
		return uint16(int32(reg) + int32(int8(b))), []byte{b}, nil

	case 0x09: // n,R (16-bit offset)
		hi := c.fetchByte(baseOffset + 1)
		lo := c.fetchByte(baseOffset + 2)
		offset := int16(uint16(hi)<<8 | uint16(lo))
		return uint16(int32(reg) + int32(offset)), []byte{hi, lo}, nil

	case 0x0B: // D,R
		return uint16(int32(reg) + int32(int16(c.D()))), nil, nil

	case 0x0C: // n,PCR (8-bit offset)
		b := c.fetchByte(baseOffset + 1)
		pc := c.PC + baseOffset + 2
		return uint16(int32(pc) + int32(int8(b))), []byte{b}, nil

	case 0x0D: // n,PCR (16-bit offset)
		hi := c.fetchByte(baseOffset + 1)
		lo := c.fetchByte(baseOffset + 2)
		offset := int16(uint16(hi)<<8 | uint16(lo))
		pc := c.PC + baseOffset + 3
		return uint16(int32(pc) + int32(offset)), []byte{hi, lo}, nil

	default:
		return 0, nil, fmt.Errorf("%w: postbyte 0x%02X", ErrInvalidIndexPostbyte, postbyte)
	}
}

// indexedRegister returns the value of the register specified in the postbyte (bits 5-6).
func (c *CPU) indexedRegister(postbyte uint8) uint16 {
	switch (postbyte >> 5) & 0x03 {
	case 0x00:
		return c.X
	case 0x01:
		return c.Y
	case 0x02:
		return c.U
	case 0x03:
		return c.S
	}
	return 0
}

// setIndexedRegister sets the register specified in the postbyte (bits 5-6).
func (c *CPU) setIndexedRegister(postbyte uint8, value uint16) {
	switch (postbyte >> 5) & 0x03 {
	case 0x00:
		c.X = value
	case 0x01:
		c.Y = value
	case 0x02:
		c.U = value
	case 0x03:
		c.S = value
	}
}

// resolveEA resolves an effective address parameter to a readable 16-bit address.
func (c *CPU) resolveEA(param any) (uint16, error) {
	switch p := param.(type) {
	case DirectPage:
		return c.dpAddr(uint8(p)), nil
	case Extended16:
		return uint16(p), nil
	case IndexedAddr:
		return uint16(p), nil
	default:
		return 0, fmt.Errorf("%w: type %T", ErrUnsupportedAddressingMode, param)
	}
}
