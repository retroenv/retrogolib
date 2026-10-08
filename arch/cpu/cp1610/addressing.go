package cp1610

// AddressingMode identifies how an instruction gets its operand.
type AddressingMode uint8

const (
	ImpliedAddressing AddressingMode = iota
	RegisterAddressing
	DirectAddressing
	IndirectAddressing
	ImmediateAddressing
	RelativeAddressing
	SpecialAddressing
)

// OperandKey identifies the register and mode fields of an opcode.
// Source and Destination are register numbers. Variant is a shift count flag,
// branch direction flag, or alternate form for an instruction.
type OperandKey struct {
	Addressing  AddressingMode
	Source      uint8
	Destination uint8
	Variant     uint8
}
