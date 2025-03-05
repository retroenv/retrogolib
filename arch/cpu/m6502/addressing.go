package m6502

// AddressingMode defines an address mode.
type AddressingMode int

// addressing modes.
const (
	NoAddressing      AddressingMode = 0
	ImpliedAddressing AddressingMode = 1 << iota
	AccumulatorAddressing
	ImmediateAddressing
	AbsoluteAddressing
	ZeroPageAddressing
	AbsoluteXAddressing
	ZeroPageXAddressing
	AbsoluteYAddressing
	ZeroPageYAddressing
	IndirectAddressing
	IndirectXAddressing
	IndirectYAddressing
	RelativeAddressing
)

// AccessMode defines an address access mode.
type AccessMode int

// address accessing modes.
const (
	NoAccess        AccessMode = 0
	ReadAccess      AccessMode = 1
	WriteAccess     AccessMode = 2
	ReadWriteAccess AccessMode = 3
)

// AccessModeConstant is used to specify for every memory address what access mode applies to it.
// A memory address like 0x4017 has a different meaning depending on the type of access.
type AccessModeConstant struct {
	Constant string
	Mode     AccessMode
}

// internal types
type (
	// AbsoluteX defines absolute addressing using the X register
	AbsoluteX uint16
	// AbsoluteY defines absolute addressing using the Y register
	AbsoluteY uint16
	// IndirectX defines indirect addressing using the X register
	IndirectX uint16
	// IndirectY defines indirect addressing using the Y register
	IndirectY uint16
	// ZeroPageX defines zeropage addressing using the X register
	ZeroPageX uint8
	// ZeroPageY defines zeropage addressing using the Y register
	ZeroPageY uint8
)

// ZeroPage indicates that the parameter for the instruction is addressing
// the zero page.
type ZeroPage uint8

// Absolute indicates that the parameter for the instruction is an
// absolute address.
type Absolute uint16

// Indirect indicates that the parameter for the instruction is using
// indirect addressing using an address and an optional X or Y register.
// For usage with a register, the indirect address is a byte and refers
// to the zero page.
type Indirect uint16

// IndirectResolved indicates that the parameter for the instruction is using
// indirect addressing using an address and an optional X or Y register.
// The final address including the memory read bug has been resolved.
type IndirectResolved uint16

// Accumulator indicates that the parameter for the instruction is the
// accumulator.
type Accumulator int
