package cpu6809

import "fmt"

var indexedModeCycleCounts = [16]uint64{
	0x00: 2,
	0x01: 3,
	0x02: 2,
	0x03: 3,
	0x05: 1,
	0x06: 1,
	0x08: 1,
	0x09: 4,
	0x0B: 4,
	0x0C: 1,
	0x0D: 5,
}

// decodeIndexedMode decodes indexed modes whose postbyte has bit seven set.
func (c *CPU) decodeIndexedMode(postbyte uint8, baseOffset, reg uint16) (uint16, []byte, error) {
	mode := postbyte & 0x0F

	switch mode {
	case 0x00, 0x01, 0x02, 0x03, 0x04:
		return c.decodeIndexedRegisterMode(postbyte, reg, mode)
	case 0x05, 0x06, 0x0B:
		return c.decodeIndexedAccumulatorMode(reg, mode)
	case 0x08, 0x09, 0x0C, 0x0D:
		return c.decodeIndexedOffsetMode(baseOffset, reg, mode)
	default:
		return 0, nil, fmt.Errorf("%w: postbyte 0x%02X", ErrInvalidIndexPostbyte, postbyte)
	}
}

func (c *CPU) decodeIndexedRegisterMode(postbyte uint8, reg uint16, mode uint8) (uint16, []byte, error) {
	switch mode {
	case 0x00: // ,R+
		c.setIndexedRegister(postbyte, reg+1)
		return reg, nil, nil
	case 0x01: // ,R++
		c.setIndexedRegister(postbyte, reg+2)
		return reg, nil, nil
	case 0x02: // ,-R
		reg--
		c.setIndexedRegister(postbyte, reg)
		return reg, nil, nil
	case 0x03: // ,--R
		reg -= 2
		c.setIndexedRegister(postbyte, reg)
		return reg, nil, nil
	default: // ,R
		return reg, nil, nil
	}
}

func (c *CPU) decodeIndexedAccumulatorMode(reg uint16, mode uint8) (uint16, []byte, error) {
	var offset int32
	switch mode {
	case 0x05: // B,R
		offset = int32(int8(c.B))
	case 0x06: // A,R
		offset = int32(int8(c.A))
	default: // D,R
		offset = int32(int16(c.D()))
	}

	return uint16(int32(reg) + offset), nil, nil
}

func (c *CPU) decodeIndexedOffsetMode(baseOffset, reg uint16, mode uint8) (uint16, []byte, error) {
	switch mode {
	case 0x08: // n,R (8-bit offset)
		value := c.fetchByte(baseOffset + 1)
		return uint16(int32(reg) + int32(int8(value))), []byte{value}, nil
	case 0x09: // n,R (16-bit offset)
		value, operands := c.fetchOperandWord(baseOffset + 1)
		return uint16(int32(reg) + int32(int16(value))), operands, nil
	case 0x0C: // n,PCR (8-bit offset)
		value := c.fetchByte(baseOffset + 1)
		pc := c.PC + baseOffset + 2
		return uint16(int32(pc) + int32(int8(value))), []byte{value}, nil
	default: // n,PCR (16-bit offset)
		value, operands := c.fetchOperandWord(baseOffset + 1)
		pc := c.PC + baseOffset + 3
		return uint16(int32(pc) + int32(int16(value))), operands, nil
	}
}

// indexedRegister returns the register selected by postbyte bits five and six.
func (c *CPU) indexedRegister(postbyte uint8) uint16 {
	switch (postbyte >> 5) & 0x03 {
	case 0x00:
		return c.X
	case 0x01:
		return c.Y
	case 0x02:
		return c.U
	default:
		return c.S
	}
}

// setIndexedRegister updates the register selected by postbyte bits five and six.
func (c *CPU) setIndexedRegister(postbyte uint8, value uint16) {
	switch (postbyte >> 5) & 0x03 {
	case 0x00:
		c.X = value
	case 0x01:
		c.Y = value
	case 0x02:
		c.U = value
	default:
		c.S = value
	}
}

// resolveEA resolves an effective-address operand.
func (c *CPU) resolveEA(param any) (uint16, error) {
	switch value := param.(type) {
	case DirectPage:
		return c.dpAddr(uint8(value)), nil
	case Extended16:
		return uint16(value), nil
	case IndexedAddr:
		return uint16(value), nil
	default:
		return 0, fmt.Errorf("%w: type %T", ErrUnsupportedAddressingMode, param)
	}
}

func (c *CPU) fetchOperandWord(offset uint16) (uint16, []byte) {
	hi := c.fetchByte(offset)
	lo := c.fetchByte(offset + 1)

	return uint16(hi)<<8 | uint16(lo), []byte{hi, lo}
}

func assignParam[T any](param any, target *T) error {
	value, ok := param.(T)
	if ok {
		*target = value
		return nil
	}

	return fmt.Errorf("%w: got %T, want %T", ErrInvalidParameterType, param, *target)
}

func immediate8Param(param any) (Immediate8, error) {
	var value Immediate8
	err := assignParam(param, &value)

	return value, err
}

func indexedAddrParam(param any) (IndexedAddr, error) {
	var value IndexedAddr
	err := assignParam(param, &value)

	return value, err
}

func branchTargetParam(param any) (uint16, error) {
	var value uint16
	err := assignParam(param, &value)

	return value, err
}

func registerPairParam(param any) (RegisterPair, error) {
	var value RegisterPair
	err := assignParam(param, &value)

	return value, err
}

func stackMaskParam(param any) (StackMask, error) {
	var value StackMask
	err := assignParam(param, &value)

	return value, err
}

func indexedModeCycles(postbyte uint8) uint64 {
	if postbyte&0x80 == 0 {
		return 1
	}
	if postbyte == 0x9F {
		return 5
	}

	cycles := indexedModeCycleCounts[postbyte&0x0F]
	if postbyte&0x10 != 0 {
		cycles += 3
	}

	return cycles
}

// readOpParam decodes the instruction's single operand and returns its raw bytes.
func readOpParam(c *CPU, mode AddressingMode, baseOffset uint16) (any, []byte, error) {
	switch mode {
	case ImpliedAddressing:
		return nil, nil, nil
	case ImmediateAddressing:
		value := c.fetchByte(baseOffset)
		return Immediate8(value), []byte{value}, nil
	case Immediate16Addressing:
		value, operands := c.fetchOperandWord(baseOffset)
		return Immediate16(value), operands, nil
	case DirectAddressing:
		value := c.fetchByte(baseOffset)
		return DirectPage(value), []byte{value}, nil
	case ExtendedAddressing:
		value, operands := c.fetchOperandWord(baseOffset)
		return Extended16(value), operands, nil
	case IndexedAddressing:
		return readIndexedParam(c, baseOffset)
	case RelativeAddressing:
		offset := int8(c.fetchByte(baseOffset))
		target := uint16(int32(c.PC) + int32(baseOffset) + 1 + int32(offset))
		return target, []byte{uint8(offset)}, nil
	case RelativeLongAddressing:
		return readRelativeLongParam(c, baseOffset)
	case RegisterAddressing:
		value := c.fetchByte(baseOffset)
		return RegisterPair(value), []byte{value}, nil
	case StackAddressing:
		value := c.fetchByte(baseOffset)
		return StackMask(value), []byte{value}, nil
	default:
		return nil, nil, fmt.Errorf("%w: mode 0x%x", ErrUnsupportedAddressingMode, mode)
	}
}

func readRelativeLongParam(c *CPU, baseOffset uint16) (any, []byte, error) {
	value, operands := c.fetchOperandWord(baseOffset)
	target := uint16(int32(c.PC) + int32(baseOffset) + 2 + int32(int16(value)))

	return target, operands, nil
}

// readIndexedParam decodes the 6809 indexed-addressing postbyte.
func readIndexedParam(c *CPU, baseOffset uint16) (any, []byte, error) {
	postbyte := c.fetchByte(baseOffset)
	operands := make([]byte, 1, 4)
	operands[0] = postbyte
	reg := c.indexedRegister(postbyte)

	if postbyte&0x80 == 0 {
		offset := int8(postbyte & 0x1F)
		if offset&0x10 != 0 {
			offset |= -0x20
		}
		addr := uint16(int32(reg) + int32(offset))

		return IndexedAddr(addr), operands, nil
	}

	mode := postbyte & 0x0F
	indirect := postbyte&0x10 != 0
	if mode == 0x0F && indirect {
		return readExtendedIndirect(c, baseOffset, operands)
	}
	// The 6809 does not encode indirect one-byte auto-increment/decrement.
	if indirect && (mode == 0x00 || mode == 0x02) {
		return nil, nil, fmt.Errorf("%w: postbyte 0x%02X", ErrInvalidIndexPostbyte, postbyte)
	}

	addr, extra, err := c.decodeIndexedMode(postbyte, baseOffset, reg)
	if err != nil {
		return nil, nil, err
	}
	operands = append(operands, extra...)
	if indirect {
		addr = c.memory.ReadWord(addr)
	}

	return IndexedAddr(addr), operands, nil
}

func readExtendedIndirect(c *CPU, baseOffset uint16, operands []byte) (any, []byte, error) {
	addr, extra := c.fetchOperandWord(baseOffset + 1)
	operands = append(operands, extra...)

	return IndexedAddr(c.memory.ReadWord(addr)), operands, nil
}
