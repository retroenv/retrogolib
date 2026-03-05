package m68000

// AddressingMode specifies how an instruction accesses its operands.
// The 68000 uses 14 addressing modes encoded in 6-bit effective address fields.
type AddressingMode int

const (
	NoAddressing        AddressingMode = 0
	DataRegDirectMode   AddressingMode = 1 << iota // Dn
	AddrRegDirectMode                               // An
	AddrRegIndirectMode                             // (An)
	PostIncrementMode                               // (An)+
	PreDecrementMode                                // -(An)
	DisplacementMode                                // d16(An)
	IndexedMode                                     // d8(An,Xn)
	AbsShortMode                                    // (xxx).W
	AbsLongMode                                     // (xxx).L
	PCDisplacementMode                              // d16(PC)
	PCIndexedMode                                   // d8(PC,Xn)
	ImmediateMode                                   // #imm
	StatusRegMode                                   // SR/CCR (implicit)
	QuickImmediateMode                              // 3-bit or 8-bit in opcode
)

// OperandSize represents the size of an operand in bytes.
type OperandSize int

const (
	SizeByte OperandSize = 1
	SizeWord OperandSize = 2
	SizeLong OperandSize = 4
)

// sizeFromBits decodes the standard size field (bits 7-6) used by most instructions.
// 00=byte, 01=word, 10=long.
func sizeFromBits(bits uint16) OperandSize {
	switch bits {
	case 0:
		return SizeByte
	case 1:
		return SizeWord
	case 2:
		return SizeLong
	default:
		return 0
	}
}
