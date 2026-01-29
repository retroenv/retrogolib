// This file contains support for unofficial/undocumented Z80 CPU instructions.
// Reference: https://www.z80.info/z80undoc.htm

package z80

// Port I/O helper methods
func (c *CPU) readPort(port uint8) uint8 {
	if c.opts.ioHandler != nil {
		return c.opts.ioHandler.ReadPort(port)
	}
	return 0xFF
}

func (c *CPU) writePort(port uint8, value uint8) {
	if c.opts.ioHandler != nil {
		c.opts.ioHandler.WritePort(port, value)
	}
}

// SLL - Shift Left Logical (undocumented)
// Shifts left and sets bit 0 to 1 (unlike SLA which sets bit 0 to 0)
var SLL = &Instruction{
	Name:       SllName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x30, Size: 2, Cycles: 8}, // SLL r (CB 30+r)
	},
	ParamFunc: sll,
}

// INF/OUTF - Undocumented port instructions
// These behave like INI/IND/OUTI/OUTD but affect flags differently

// INF - Input and decrement with different flag behavior
var INF = &Instruction{
	Name:       InfName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xAA, Size: 2, Cycles: 16}, // ED AA
	},
	NoParamFunc: inf,
}

// OUTF - Output and decrement with different flag behavior
var OUTF = &Instruction{
	Name:       OutfName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xAB, Size: 2, Cycles: 16}, // ED AB
	},
	NoParamFunc: outf,
}

// Undocumented flag effects for various instructions
// Many Z80 instructions set the undocumented X and Y flags based on bits 3 and 5
// of the result or intermediate values

// The following instructions are documented but have undocumented flag effects:

// LDIX/LDIY - Undocumented effects on X and Y flags during LD (IX+d),r operations
// These are not separate instructions but rather undocumented behaviors

// BIT with IX/IY - When using BIT n,(IX+d) or BIT n,(IY+d),
// the X and Y flags are set from the high byte of IX/IY+d rather than
// from the tested memory location

// Some block instructions have undocumented behaviors:
// - CPI/CPD/CPIR/CPDR affect X and Y flags
// - LDI/LDD/LDIR/LDDR affect X and Y flags

// NopUndoc1 represents undocumented single-byte NOPs using DD prefix alone.
var NopUndoc1 = &Instruction{
	Name:       Nop.Name,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: PrefixDD, Size: 1, Cycles: 4}, // DD alone (partial IX prefix)
	},
	NoParamFunc: nop,
}

// NopUndoc2 represents undocumented single-byte NOPs using FD prefix alone.
var NopUndoc2 = &Instruction{
	Name:       Nop.Name,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: PrefixFD, Size: 1, Cycles: 4}, // FD alone (partial IY prefix)
	},
	NoParamFunc: nop,
}

// Undocumented register combinations
// Some combinations of prefix bytes and regular opcodes create undocumented behaviors

// Note: The Z80 has fewer truly undocumented instructions compared to the 6502
// Most "undocumented" behaviors are:
// 1. Partial prefix sequences that act as NOPs
// 2. Undocumented flag effects (X and Y flags)
// 3. Undocumented behaviors of some bit operations with index registers

// UnofficialInstructions maps undocumented instruction names to their definitions
var UnofficialInstructions = map[string]*Instruction{
	SllName:        SLL,
	InfName:        INF,
	OutfName:       OUTF,
	"nop_undoc_dd": NopUndoc1, // DD prefix alone
	"nop_undoc_fd": NopUndoc2, // FD prefix alone
}

// IsUnofficialInstruction returns true if the instruction name corresponds to
// an undocumented/unofficial Z80 instruction
func IsUnofficialInstruction(name string) bool {
	_, exists := UnofficialInstructions[name]
	return exists
}

// Emulation functions for undocumented instructions

// sll performs shift left logical (undocumented) - like SLA but sets bit 0 to 1
func sll(c *CPU, params ...any) error {
	if len(params) == 0 {
		return ErrMissingParameter
	}

	reg, ok := params[0].(Register)
	if !ok {
		return ErrInvalidParameterType
	}

	value := c.GetRegisterValue(uint8(reg))

	// Set carry flag from bit 7
	c.setC(value&0x80 != 0)

	// Shift left and set bit 0 to 1 (this is the undocumented behavior)
	result := (value << 1) | 0x01

	c.SetRegisterValue(uint8(reg), result)
	c.setSZP(result)
	c.setH(false)
	c.setN(false)

	return nil
}

// inf performs input and decrement (undocumented port instruction)
func inf(c *CPU) error {
	// Read from port C into memory location (HL)
	value := c.readPort(c.C)
	address := uint16(c.H)<<8 | uint16(c.L)
	c.memory.Write(address, value)

	// Decrement HL
	hl := address - 1
	c.H = uint8(hl >> 8)
	c.L = uint8(hl & 0xFF)

	// Decrement B
	c.B--

	// Set flags (undocumented behavior may differ from documented INI/IND)
	setFlag(&c.Flags.Z, c.B == 0)
	c.setN(true)

	return nil
}

// outf performs output and decrement (undocumented port instruction)
func outf(c *CPU) error {
	// Read from memory location (HL)
	address := uint16(c.H)<<8 | uint16(c.L)
	value := c.memory.Read(address)

	// Output to port C
	c.writePort(c.C, value)

	// Decrement HL
	hl := address - 1
	c.H = uint8(hl >> 8)
	c.L = uint8(hl & 0xFF)

	// Decrement B
	c.B--

	// Set flags (undocumented behavior may differ from documented OUTI/OUTD)
	setFlag(&c.Flags.Z, c.B == 0)
	c.setN(true)

	return nil
}
