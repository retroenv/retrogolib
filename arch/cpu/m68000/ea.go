package m68000

import "fmt"

// EffectiveAddress represents a decoded effective address with its resolved location.
type EffectiveAddress struct {
	Mode    uint8       // EA mode (0-7)
	Reg     uint8       // EA register (0-7)
	Size    OperandSize // Operand size
	Address uint32      // Resolved memory address (for memory modes)
	Value   uint32      // Immediate or register value (for non-memory modes)
}

// decodeEA decodes an effective address from the mode and register fields.
// It reads extension words from the instruction stream as needed.
func (c *CPU) decodeEA(mode, reg uint8, size OperandSize) (EffectiveAddress, error) {
	ea := EffectiveAddress{
		Mode: mode,
		Reg:  reg,
		Size: size,
	}

	switch mode {
	case 0: // Data register direct: Dn
		ea.Value = c.getRegD(reg, size)
		return ea, nil

	case 1: // Address register direct: An
		ea.Value = c.getRegA(reg)
		return ea, nil

	case 2: // Address register indirect: (An)
		ea.Address = c.getRegA(reg)
		return ea, nil

	case 3: // Postincrement: (An)+
		ea.Address = c.getRegA(reg)
		c.setRegA(reg, c.getRegA(reg)+incrementSize(reg, size))
		return ea, nil

	case 4: // Predecrement: -(An)
		c.setRegA(reg, c.getRegA(reg)-incrementSize(reg, size))
		ea.Address = c.getRegA(reg)
		return ea, nil

	case 5: // Displacement: d16(An)
		disp := int16(c.readWord())
		ea.Address = uint32(int32(c.getRegA(reg)) + int32(disp))
		return ea, nil

	case 6: // Indexed: d8(An,Xn)
		ea.Address = c.decodeIndexed(c.getRegA(reg))
		return ea, nil

	case 7: // Extended modes based on register field
		return c.decodeEAMode7(ea)

	default:
		return ea, fmt.Errorf("%w: mode %d reg %d", ErrInvalidAddressMode, mode, reg)
	}
}

// decodeEAMode7 handles the extended addressing modes (mode 7, reg 0-4).
func (c *CPU) decodeEAMode7(ea EffectiveAddress) (EffectiveAddress, error) {
	switch ea.Reg {
	case 0: // Absolute short: (xxx).W
		addr := int16(c.readWord())
		ea.Address = uint32(int32(addr)) & addressMask
		return ea, nil

	case 1: // Absolute long: (xxx).L
		ea.Address = c.readLong() & addressMask
		return ea, nil

	case 2: // PC displacement: d16(PC)
		pcBefore := c.PC
		disp := int16(c.readWord())
		ea.Address = uint32(int32(pcBefore) + int32(disp))
		return ea, nil

	case 3: // PC indexed: d8(PC,Xn)
		pcBefore := c.PC
		ea.Address = c.decodeIndexed(pcBefore)
		return ea, nil

	case 4: // Immediate: #imm
		ea.Value = c.readImmediate(ea.Size)
		return ea, nil

	default:
		return ea, fmt.Errorf("%w: mode 7 reg %d", ErrInvalidAddressMode, ea.Reg)
	}
}

// decodeIndexed decodes an indexed extension word and returns the computed address.
// Extension word format: D/A | Reg | W/L | 0 | 0 | 0 | displacement(8 bits).
func (c *CPU) decodeIndexed(baseAddr uint32) uint32 {
	ext := c.readWord()

	disp := int8(ext & 0xFF)
	indexReg := (ext >> 12) & 7
	isAddrReg := ext&0x8000 != 0
	isLong := ext&0x0800 != 0

	var indexValue int32

	if isAddrReg {
		indexValue = int32(c.getRegA(uint8(indexReg)))
	} else {
		indexValue = int32(c.D[indexReg])
	}

	if !isLong {
		indexValue = int32(int16(indexValue))
	}

	return uint32(int32(baseAddr) + indexValue + int32(disp))
}

// readEA reads the value at an effective address.
func (c *CPU) readEA(ea EffectiveAddress) (uint32, error) {
	switch ea.Mode {
	case 0: // Data register direct
		return c.getRegD(ea.Reg, ea.Size), nil

	case 1: // Address register direct
		return c.getRegA(ea.Reg), nil

	case 2, 3, 4, 5, 6: // Memory modes
		return c.readMemory(ea.Address, ea.Size)

	case 7:
		if ea.Reg == 4 { // Immediate
			return ea.Value, nil
		}
		return c.readMemory(ea.Address, ea.Size)

	default:
		return 0, fmt.Errorf("%w: read mode %d", ErrInvalidAddressMode, ea.Mode)
	}
}

// writeEA writes a value to an effective address.
func (c *CPU) writeEA(ea EffectiveAddress, value uint32) error {
	switch ea.Mode {
	case 0: // Data register direct
		c.setRegD(ea.Reg, value, ea.Size)
		return nil

	case 1: // Address register direct
		c.setRegA(ea.Reg, value)
		return nil

	case 2, 3, 4, 5, 6: // Memory modes
		return c.writeMemory(ea.Address, value, ea.Size)

	case 7:
		if ea.Reg <= 1 { // Absolute short/long
			return c.writeMemory(ea.Address, value, ea.Size)
		}
		return fmt.Errorf("%w: write mode 7 reg %d", ErrInvalidAddressMode, ea.Reg)

	default:
		return fmt.Errorf("%w: write mode %d", ErrInvalidAddressMode, ea.Mode)
	}
}

// readMemory reads a value from memory at the given address with the given size.
func (c *CPU) readMemory(addr uint32, size OperandSize) (uint32, error) {
	addr &= addressMask
	switch size {
	case SizeByte:
		return uint32(c.bus.Read(addr)), nil
	case SizeWord:
		return uint32(c.bus.ReadWord(addr)), nil
	case SizeLong:
		return c.bus.ReadLong(addr), nil
	default:
		return 0, ErrInvalidOperandSize
	}
}

// writeMemory writes a value to memory at the given address with the given size.
func (c *CPU) writeMemory(addr, value uint32, size OperandSize) error {
	addr &= addressMask
	switch size {
	case SizeByte:
		c.bus.Write(addr, uint8(value))
	case SizeWord:
		c.bus.WriteWord(addr, uint16(value))
	case SizeLong:
		c.bus.WriteLong(addr, value)
	default:
		return ErrInvalidOperandSize
	}
	return nil
}
