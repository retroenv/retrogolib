package x86

// AddressingMode represents different x86 addressing modes.
type AddressingMode uint8

// x86 addressing modes for DOS-era instructions.
const (
	// Basic addressing modes
	ImpliedAddressing          AddressingMode = iota // No operands (NOP, CLC)
	RegisterAddressing                               // Register operand (INC AX)
	ImmediateAddressing                              // Immediate operand (MOV AX, 1234h)
	DirectAddressing                                 // Direct memory (MOV AX, [1234h])
	RegisterIndirectAddressing                       // Register indirect ([BX])

	// Indexed addressing modes
	IndexedAddressing      // Base + index ([BX+SI])
	BasedIndexedAddressing // Base + index + displacement ([BX+SI+08h])
	RelativeAddressing     // PC-relative (JMP rel8, JMP rel16)

	// Special addressing modes
	SegmentOffsetAddressing // Segment:offset (JMP 1234h:5678h)
	PortAddressing          // I/O port (IN AL, 21h)
	StringAddressing        // String operations (MOVSB)
	StackAddressing         // Stack operations (PUSH, POP)

	// ModR/M based addressing modes
	ModRMRegisterAddressing  // ModR/M register-to-register
	ModRMMemoryAddressing    // ModR/M memory addressing
	ModRMImmediateAddressing // ModR/M with immediate operand
)

// AddressingModeNames provides string representations of addressing modes.
var AddressingModeNames = map[AddressingMode]string{
	ImpliedAddressing:          "implied",
	RegisterAddressing:         "register",
	ImmediateAddressing:        "immediate",
	DirectAddressing:           "direct",
	RegisterIndirectAddressing: "register_indirect",
	IndexedAddressing:          "indexed",
	BasedIndexedAddressing:     "based_indexed",
	RelativeAddressing:         "relative",
	SegmentOffsetAddressing:    "segment_offset",
	PortAddressing:             "port",
	StringAddressing:           "string",
	StackAddressing:            "stack",
	ModRMRegisterAddressing:    "modrm_register",
	ModRMMemoryAddressing:      "modrm_memory",
	ModRMImmediateAddressing:   "modrm_immediate",
}

// String returns the string representation of the addressing mode.
func (am AddressingMode) String() string {
	if name, exists := AddressingModeNames[am]; exists {
		return name
	}
	return "unknown"
}

// RegisterParam represents different register parameters used in instructions.
type RegisterParam uint8

// x86 register parameters for instruction operands.
const (
	// 8-bit registers
	RegAL RegisterParam = iota
	RegCL
	RegDL
	RegBL
	RegAH
	RegCH
	RegDH
	RegBH

	// 16-bit general purpose registers
	RegAX
	RegCX
	RegDX
	RegBX
	RegSP
	RegBP
	RegSI
	RegDI

	// Segment registers
	RegES
	RegCS
	RegSS
	RegDS

	// Memory references
	RegBXSIRef // [BX+SI]
	RegBXDIRef // [BX+DI]
	RegBPSIRef // [BP+SI]
	RegBPDIRef // [BP+DI]
	RegSIRef   // [SI]
	RegDIRef   // [DI]
	RegBPRef   // [BP]
	RegBXRef   // [BX]

	// Immediate values
	RegImm8    // 8-bit immediate
	RegImm16   // 16-bit immediate
	RegRel8    // 8-bit relative
	RegRel16   // 16-bit relative
	RegPtr1616 // 16:16 far pointer
	RegMem     // Memory operand
	RegPort    // I/O port

	// Special registers/flags
	RegFlags // FLAGS register
	RegIP    // Instruction Pointer
)

// RegisterParamNames provides string representations of register parameters.
var RegisterParamNames = map[RegisterParam]string{
	// 8-bit registers
	RegAL: "al", RegCL: "cl", RegDL: "dl", RegBL: "bl",
	RegAH: "ah", RegCH: "ch", RegDH: "dh", RegBH: "bh",

	// 16-bit registers
	RegAX: "ax", RegCX: "cx", RegDX: "dx", RegBX: "bx",
	RegSP: "sp", RegBP: "bp", RegSI: "si", RegDI: "di",

	// Segment registers
	RegES: "es", RegCS: "cs", RegSS: "ss", RegDS: "ds",

	// Memory references
	RegBXSIRef: "[bx+si]", RegBXDIRef: "[bx+di]",
	RegBPSIRef: "[bp+si]", RegBPDIRef: "[bp+di]",
	RegSIRef: "[si]", RegDIRef: "[di]",
	RegBPRef: "[bp]", RegBXRef: "[bx]",

	// Immediate and special
	RegImm8: "imm8", RegImm16: "imm16",
	RegRel8: "rel8", RegRel16: "rel16",
	RegPtr1616: "ptr16:16", RegMem: "mem",
	RegPort: "port", RegFlags: "flags", RegIP: "ip",
}

// String returns the string representation of the register parameter.
func (rp RegisterParam) String() string {
	if name, exists := RegisterParamNames[rp]; exists {
		return name
	}
	return "unknown"
}

// Is8Bit returns true if the register parameter represents an 8-bit register.
func (rp RegisterParam) Is8Bit() bool {
	return rp >= RegAL && rp <= RegBH
}

// Is16Bit returns true if the register parameter represents a 16-bit register.
func (rp RegisterParam) Is16Bit() bool {
	return rp >= RegAX && rp <= RegDI
}

// IsSegment returns true if the register parameter represents a segment register.
func (rp RegisterParam) IsSegment() bool {
	return rp >= RegES && rp <= RegDS
}

// IsMemoryRef returns true if the register parameter represents a memory reference.
func (rp RegisterParam) IsMemoryRef() bool {
	return rp >= RegBXSIRef && rp <= RegBXRef
}

// IsImmediate returns true if the register parameter represents an immediate value.
func (rp RegisterParam) IsImmediate() bool {
	return rp == RegImm8 || rp == RegImm16 || rp == RegRel8 || rp == RegRel16 || rp == RegPtr1616
}

// GetRegisterSize returns the size in bytes of the register (1 for 8-bit, 2 for 16-bit).
func (rp RegisterParam) GetRegisterSize() int {
	if rp.Is8Bit() {
		return 1
	}
	if rp.Is16Bit() || rp.IsSegment() {
		return 2
	}
	return 0 // Unknown or special register
}

// ModRM represents the ModR/M byte used in x86 instruction encoding.
type ModRM struct {
	Mod uint8 // Mode field (bits 7-6)
	Reg uint8 // Register field (bits 5-3)
	RM  uint8 // R/M field (bits 2-0)
}

// NewModRM creates a ModR/M byte from its components.
func NewModRM(mod, reg, rm uint8) ModRM {
	return ModRM{
		Mod: mod & 0x03,
		Reg: reg & 0x07,
		RM:  rm & 0x07,
	}
}

// FromByte creates a ModR/M from a raw byte value.
func (m *ModRM) FromByte(value uint8) {
	m.Mod = (value >> 6) & 0x03
	m.Reg = (value >> 3) & 0x07
	m.RM = value & 0x07
}

// ToByte converts the ModR/M to a raw byte value.
func (m ModRM) ToByte() uint8 {
	return (m.Mod << 6) | (m.Reg << 3) | m.RM
}

// GetEffectiveAddress calculates the effective address based on ModR/M and SIB.
// This is a simplified version for basic addressing modes.
func (c *CPU) GetEffectiveAddress(modrm ModRM, displacement int16, segment uint16) uint32 {
	var offset uint16

	switch modrm.Mod {
	case 0: // No displacement (except special cases)
		switch modrm.RM {
		case 0: // [BX + SI]
			offset = c.BX + c.SI
		case 1: // [BX + DI]
			offset = c.BX + c.DI
		case 2: // [BP + SI] (SS segment default)
			offset = c.BP + c.SI
			if segment == 0 {
				segment = c.SS
			}
		case 3: // [BP + DI] (SS segment default)
			offset = c.BP + c.DI
			if segment == 0 {
				segment = c.SS
			}
		case 4: // [SI]
			offset = c.SI
		case 5: // [DI]
			offset = c.DI
		case 6: // Direct address (16-bit displacement)
			offset = uint16(displacement)
		case 7: // [BX]
			offset = c.BX
		}
	case 1: // 8-bit displacement
		switch modrm.RM {
		case 0: // [BX + SI + disp8]
			offset = c.BX + c.SI + uint16(int8(displacement))
		case 1: // [BX + DI + disp8]
			offset = c.BX + c.DI + uint16(int8(displacement))
		case 2: // [BP + SI + disp8] (SS segment default)
			offset = c.BP + c.SI + uint16(int8(displacement))
			if segment == 0 {
				segment = c.SS
			}
		case 3: // [BP + DI + disp8] (SS segment default)
			offset = c.BP + c.DI + uint16(int8(displacement))
			if segment == 0 {
				segment = c.SS
			}
		case 4: // [SI + disp8]
			offset = c.SI + uint16(int8(displacement))
		case 5: // [DI + disp8]
			offset = c.DI + uint16(int8(displacement))
		case 6: // [BP + disp8] (SS segment default)
			offset = c.BP + uint16(int8(displacement))
			if segment == 0 {
				segment = c.SS
			}
		case 7: // [BX + disp8]
			offset = c.BX + uint16(int8(displacement))
		}
	case 2: // 16-bit displacement
		switch modrm.RM {
		case 0: // [BX + SI + disp16]
			offset = c.BX + c.SI + uint16(displacement)
		case 1: // [BX + DI + disp16]
			offset = c.BX + c.DI + uint16(displacement)
		case 2: // [BP + SI + disp16] (SS segment default)
			offset = c.BP + c.SI + uint16(displacement)
			if segment == 0 {
				segment = c.SS
			}
		case 3: // [BP + DI + disp16] (SS segment default)
			offset = c.BP + c.DI + uint16(displacement)
			if segment == 0 {
				segment = c.SS
			}
		case 4: // [SI + disp16]
			offset = c.SI + uint16(displacement)
		case 5: // [DI + disp16]
			offset = c.DI + uint16(displacement)
		case 6: // [BP + disp16] (SS segment default)
			offset = c.BP + uint16(displacement)
			if segment == 0 {
				segment = c.SS
			}
		case 7: // [BX + disp16]
			offset = c.BX + uint16(displacement)
		}
	case 3: // Register addressing (not memory)
		return 0 // This should not be called for register addressing
	}

	// Use DS as default segment if no specific segment was provided
	if segment == 0 {
		segment = c.DS
	}

	return c.CalculateAddress(segment, offset)
}
