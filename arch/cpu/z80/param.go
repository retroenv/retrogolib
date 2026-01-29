package z80

import (
	"fmt"
)

type paramReaderFunc func(c *CPU) ([]any, []byte)

// paramReader maps addressing modes to their parameter reading functions.
var paramReader = map[AddressingMode]paramReaderFunc{
	ImpliedAddressing:          paramReaderImplied,
	RegisterAddressing:         paramReaderRegister,
	ImmediateAddressing:        paramReaderImmediate,
	ExtendedAddressing:         paramReaderExtended,
	RegisterIndirectAddressing: paramReaderRegisterIndirect,
	RelativeAddressing:         paramReaderRelative,
	BitAddressing:              paramReaderBit,
	PortAddressing:             paramReaderPort,
}

// GetRegisterValue returns the value of a register by its encoding number (0-7).
// Register 6 accesses memory at (HL) instead of a direct register.
// Returns 0 for invalid register numbers.
func (c *CPU) GetRegisterValue(reg uint8) uint8 {
	switch reg {
	case 0:
		return c.B
	case 1:
		return c.C
	case 2:
		return c.D
	case 3:
		return c.E
	case 4:
		return c.H
	case 5:
		return c.L
	case 6:
		return c.memory.Read(c.hl())
	case 7:
		return c.A
	default:
		return 0
	}
}

// SetRegisterValue sets the value of a register by its encoding number (0-7).
// Register 6 writes to memory at (HL) instead of a direct register.
// Invalid register numbers are silently ignored.
func (c *CPU) SetRegisterValue(reg uint8, value uint8) {
	switch reg {
	case 0:
		c.B = value
	case 1:
		c.C = value
	case 2:
		c.D = value
	case 3:
		c.E = value
	case 4:
		c.H = value
	case 5:
		c.L = value
	case 6:
		c.memory.Write(c.hl(), value)
	case 7:
		c.A = value
	}
}

// readOpParams reads the opcode parameters after the first opcode byte
// and translates it into emulator specific types.
func readOpParams(c *CPU, addressing AddressingMode) ([]any, []byte, error) {
	fun, ok := paramReader[addressing]
	if !ok {
		return nil, nil, fmt.Errorf("%w: mode %02x", ErrUnsupportedAddressingMode, addressing)
	}

	params, opcodes := fun(c)
	return params, opcodes, nil
}

func paramReaderImplied(_ *CPU) ([]any, []byte) {
	return nil, nil
}

func paramReaderRegister(c *CPU) ([]any, []byte) {
	// Register addressing is encoded in the opcode itself
	// The specific register is determined by the opcode bits
	opcode := c.memory.Read(c.PC)

	// Extract register from opcode (varies by instruction)
	// For most instructions, bits 0-2 specify source, bits 3-5 specify destination
	srcReg := opcode & 0x07
	dstReg := (opcode >> 3) & 0x07

	params := []any{Register(srcReg), Register(dstReg)}
	return params, nil
}

func paramReaderImmediate(c *CPU) ([]any, []byte) {
	// Check if this is a DD-prefixed instruction
	prefixByte := c.memory.Read(c.PC)
	if prefixByte == PrefixDD {
		// For DD-prefixed instructions, read the 16-bit immediate after the DD XX prefix
		b1 := c.memory.Read(c.PC + 2) // Low byte (after DD 21)
		b2 := c.memory.Read(c.PC + 3) // High byte
		value := uint16(b2)<<8 | uint16(b1)
		params := []any{Immediate16(value)}
		opcodes := [2]uint8{b1, b2}
		return params, opcodes[:]
	}

	opcode := c.memory.Read(c.PC)
	opcodeInfo := Opcodes[opcode]

	if opcodeInfo.Size == 3 {
		// 16-bit immediate (3-byte instruction: opcode + low byte + high byte)
		b1 := c.memory.Read(c.PC + 1) // Low byte
		b2 := c.memory.Read(c.PC + 2) // High byte
		value := uint16(b2)<<8 | uint16(b1)
		params := []any{Immediate16(value)}
		opcodes := [2]uint8{b1, b2}
		return params, opcodes[:]
	} else {
		// 8-bit immediate (2-byte instruction: opcode + immediate)
		b := c.memory.Read(c.PC + 1)
		params := []any{Immediate8(b)}
		opcodes := [1]uint8{b}
		return params, opcodes[:]
	}
}

func paramReaderExtended(c *CPU) ([]any, []byte) {
	b1 := c.memory.Read(c.PC + 1) // Low byte
	b2 := c.memory.Read(c.PC + 2) // High byte

	address := uint16(b2)<<8 | uint16(b1)
	params := []any{Extended(address)}
	opcodes := [2]uint8{b1, b2}
	return params, opcodes[:]
}

func paramReaderRegisterIndirect(c *CPU) ([]any, []byte) {
	// Register indirect addressing - the register pair is encoded in the opcode
	opcode := c.memory.Read(c.PC)

	// Determine register pair from opcode
	var regPair uint16
	switch (opcode >> 4) & 0x03 {
	case 0: // BC
		regPair = uint16(c.B)<<8 | uint16(c.C)
	case 1: // DE
		regPair = uint16(c.D)<<8 | uint16(c.E)
	case 2: // HL
		regPair = uint16(c.H)<<8 | uint16(c.L)
	case 3: // SP
		regPair = c.SP
	}

	params := []any{RegisterIndirect(regPair)}
	return params, nil
}

func paramReaderRelative(c *CPU) ([]any, []byte) {
	offset := int8(c.memory.Read(c.PC + 1))

	params := []any{Relative(offset)}
	opcodes := [1]uint8{uint8(offset)}
	return params, opcodes[:]
}

func paramReaderBit(c *CPU) ([]any, []byte) {
	// Bit addressing - bit number is encoded in the opcode
	opcode := c.memory.Read(c.PC)

	bitNum := (opcode >> 3) & 0x07 // Bits 3-5 contain bit number
	targetReg := opcode & 0x07     // Bits 0-2 contain target register

	params := []any{Bit(bitNum), Register(targetReg)}
	return params, nil
}

func paramReaderPort(c *CPU) ([]any, []byte) {
	// Port addressing can be immediate (n) or register indirect (C)
	opcode := c.memory.Read(c.PC)

	if opcode == 0xDB || opcode == 0xD3 { // IN A,(n) or OUT (n),A
		portAddr := c.memory.Read(c.PC + 1)
		params := []any{Port(portAddr)}
		opcodes := [1]uint8{portAddr}
		return params, opcodes[:]
	}

	// Port (C) - use C register as port address
	params := []any{Port(c.C)}
	return params, nil
}
