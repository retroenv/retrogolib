package cpu68000

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
func (cpu *CPU) decodeEA(mode, reg uint8, size OperandSize) (EffectiveAddress, error) {
	ea := EffectiveAddress{
		Mode: mode,
		Reg:  reg,
		Size: size,
	}

	switch mode {
	case 0: // Data register direct: Dn
		ea.Value = cpu.getRegD(reg, size)
		return ea, nil

	case 1: // Address register direct: An
		ea.Value = cpu.getRegA(reg)
		return ea, nil

	case 2: // Address register indirect: (An)
		ea.Address = cpu.getRegA(reg)
		return ea, nil

	case 3: // Postincrement: (An)+
		ea.Address = cpu.getRegA(reg)
		cpu.setRegA(reg, cpu.getRegA(reg)+incrementSize(reg, size))
		return ea, nil

	case 4: // Predecrement: -(An)
		cpu.accessCycles += 2
		cpu.setRegA(reg, cpu.getRegA(reg)-incrementSize(reg, size))
		ea.Address = cpu.getRegA(reg)
		return ea, nil

	case 5: // Displacement: d16(An)
		disp := int16(cpu.readWord())
		ea.Address = uint32(int32(cpu.getRegA(reg)) + int32(disp))
		return ea, nil

	case 6: // Indexed: d8(An,Xn)
		ea.Address = cpu.decodeIndexed(cpu.getRegA(reg))
		return ea, nil

	case 7: // Extended modes based on register field
		return cpu.decodeEAMode7(ea)

	default:
		return ea, fmt.Errorf("%w: mode %d reg %d", ErrInvalidAddressMode, mode, reg)
	}
}

// decodeEAMode7 handles the extended addressing modes (mode 7, reg 0-4).
func (cpu *CPU) decodeEAMode7(ea EffectiveAddress) (EffectiveAddress, error) {
	switch ea.Reg {
	case 0: // Absolute short: (xxx).W
		addr := int16(cpu.readWord())
		ea.Address = uint32(int32(addr))
		return ea, nil

	case 1: // Absolute long: (xxx).L
		ea.Address = cpu.readLong()
		return ea, nil

	case 2: // PC displacement: d16(PC)
		pcBefore := cpu.PC
		disp := int16(cpu.readWord())
		ea.Address = uint32(int32(pcBefore) + int32(disp))
		return ea, nil

	case 3: // PC indexed: d8(PC,Xn)
		pcBefore := cpu.PC
		ea.Address = cpu.decodeIndexed(pcBefore)
		return ea, nil

	case 4: // Immediate: #imm
		ea.Value = cpu.readImmediate(ea.Size)
		return ea, nil

	default:
		return ea, fmt.Errorf("%w: mode 7 reg %d", ErrInvalidAddressMode, ea.Reg)
	}
}

// decodeIndexed decodes an indexed extension word and returns the computed address.
// Extension word format: D/A | Reg | W/L | 0 | 0 | 0 | displacement(8 bits).
func (cpu *CPU) decodeIndexed(baseAddr uint32) uint32 {
	ext := cpu.readWord()
	cpu.accessCycles += 2

	disp := int8(ext & 0xFF)
	indexReg := (ext >> 12) & 7
	isAddrReg := ext&0x8000 != 0
	isLong := ext&0x0800 != 0

	var indexValue int32

	if isAddrReg {
		indexValue = int32(cpu.getRegA(uint8(indexReg)))
	} else {
		indexValue = int32(cpu.D[indexReg])
	}

	if !isLong {
		indexValue = int32(int16(indexValue))
	}

	return uint32(int32(baseAddr) + indexValue + int32(disp))
}

// readEA reads the value at an effective address.
func (cpu *CPU) readEA(ea EffectiveAddress) (uint32, error) {
	switch ea.Mode {
	case 0: // Data register direct
		return cpu.getRegD(ea.Reg, ea.Size), nil

	case 1: // Address register direct
		return cpu.getRegA(ea.Reg), nil

	case 2, 3, 4, 5, 6: // Memory modes
		return cpu.readMemory(ea.Address, ea.Size)

	case 7:
		if ea.Reg == 4 { // Immediate
			return ea.Value, nil
		}
		return cpu.readMemory(ea.Address, ea.Size)

	default:
		return 0, fmt.Errorf("%w: read mode %d", ErrInvalidAddressMode, ea.Mode)
	}
}

// writeEA writes a value to an effective address.
func (cpu *CPU) writeEA(ea EffectiveAddress, value uint32) error {
	switch ea.Mode {
	case 0: // Data register direct
		cpu.setRegD(ea.Reg, value, ea.Size)
		return nil

	case 1: // Address register direct
		cpu.setRegA(ea.Reg, value)
		return nil

	case 2, 3, 4, 5, 6: // Memory modes
		return cpu.writeMemory(ea.Address, value, ea.Size)

	case 7:
		if ea.Reg <= 1 { // Absolute short/long
			return cpu.writeMemory(ea.Address, value, ea.Size)
		}
		return fmt.Errorf("%w: write mode 7 reg %d", ErrInvalidAddressMode, ea.Reg)

	default:
		return fmt.Errorf("%w: write mode %d", ErrInvalidAddressMode, ea.Mode)
	}
}

// readMemory reads a value from memory at the given address with the given size.
func (cpu *CPU) readMemory(addr uint32, size OperandSize) (uint32, error) {
	switch size {
	case SizeByte:
		return uint32(cpu.readByte(addr)), nil
	case SizeWord:
		return uint32(cpu.readBusWord(addr, dataSpace)), nil
	case SizeLong:
		return cpu.readBusLong(addr), nil
	default:
		return 0, ErrInvalidOperandSize
	}
}

// writeMemory writes a value to memory at the given address with the given size.
func (cpu *CPU) writeMemory(addr, value uint32, size OperandSize) error {
	switch size {
	case SizeByte:
		cpu.writeByte(addr, uint8(value))
	case SizeWord:
		cpu.writeBusWord(addr, uint16(value))
	case SizeLong:
		cpu.writeBusLong(addr, value)
	default:
		return ErrInvalidOperandSize
	}
	return nil
}
