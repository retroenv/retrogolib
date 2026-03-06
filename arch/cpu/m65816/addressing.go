// Package m65816 addressing mode definitions.
package m65816

// AddressingMode specifies how a 65816 instruction accesses its operands.
type AddressingMode int

const (
	NoAddressing AddressingMode = 0

	// Inherited from 65C02
	ImpliedAddressing     AddressingMode = 1 << iota // no operand (NOP, TAX)
	AccumulatorAddressing                            // A register (ASL A)
	ImmediateAddressing                              // #imm (LDA #$42) - size varies with M/X

	// Direct Page (replaces zero page; base address from DP register)
	DirectPageAddressing                 // dp (LDA dp)
	DirectPageIndexedXAddressing         // dp,X
	DirectPageIndexedYAddressing         // dp,Y
	DirectPageIndirectAddressing         // (dp)
	DirectPageIndexedXIndirectAddressing // (dp,X)
	DirectPageIndirectIndexedYAddressing // (dp),Y

	// Direct Page Long (24-bit pointer in direct page)
	DirectPageIndirectLongAddressing         // [dp]
	DirectPageIndirectLongIndexedYAddressing // [dp],Y

	// Absolute (16-bit address in current data bank)
	AbsoluteAddressing                 // abs
	AbsoluteIndexedXAddressing         // abs,X
	AbsoluteIndexedYAddressing         // abs,Y
	AbsoluteIndirectAddressing         // (abs) -- JMP only
	AbsoluteIndexedXIndirectAddressing // (abs,X) -- JMP/JSR only

	// Absolute Long (24-bit address)
	AbsoluteLongAddressing         // al
	AbsoluteLongIndexedXAddressing // al,X

	// Absolute Long Indirect (24-bit pointer at absolute address)
	AbsoluteIndirectLongAddressing // [abs] -- JML only

	// Stack Relative
	StackRelativeAddressing                 // sr,S
	StackRelativeIndirectIndexedYAddressing // (sr,S),Y

	// Relative
	RelativeAddressing     // 8-bit signed offset (branches)
	RelativeLongAddressing // 16-bit signed offset (BRL, PER)

	// Block Move
	BlockMoveAddressing // srcBank,dstBank (MVN, MVP)
)

// Typed operand types for safe dispatch in param readers and handlers.
type (
	Immediate8  uint8  // 8-bit immediate value
	Immediate16 uint16 // 16-bit immediate value

	DirectPage     uint8  // 8-bit direct page offset
	DirectPageX    uint8  // direct page + X
	DirectPageY    uint8  // direct page + Y
	DPIndirect     uint32 // resolved address from (dp) indirect
	DPIndirectX    uint32 // resolved address from (dp,X) indirect
	DPIndirectY    uint32 // resolved address from (dp),Y
	DPIndirectLong uint32 // resolved address from [dp]
	DPIndLongY     uint32 // resolved address from [dp],Y

	Absolute16  uint16 // 16-bit absolute (DB:addr)
	AbsoluteX16 uint32 // resolved abs,X
	AbsoluteY16 uint32 // resolved abs,Y
	AbsLong     uint32 // 24-bit absolute long
	AbsLongX    uint32 // 24-bit absolute long + X

	StackRel uint8  // stack-relative offset
	SRIndY   uint32 // resolved (sr,S),Y address

	RelOffset  int8  // 8-bit branch offset (already resolved to absolute)
	LongOffset int16 // 16-bit branch offset (already resolved to absolute)

	BlockMove struct{ Src, Dst uint8 } // source and destination banks for MVN/MVP
)

// Accumulator is a sentinel type for accumulator-mode instructions.
type Accumulator struct{}
