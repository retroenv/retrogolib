package m6809

// AddressingMode specifies how a 6809 instruction accesses its operands.
type AddressingMode int

const (
	NoAddressing AddressingMode = 0

	ImpliedAddressing      AddressingMode = 1 << iota // no operand (NOP, ABX)
	ImmediateAddressing                               // #imm8 (LDA #$42)
	Immediate16Addressing                             // #imm16 (LDD #$1234)
	DirectAddressing                                  // dp (LDA <$10)
	ExtendedAddressing                                // abs16 (LDA $1234)
	IndexedAddressing                                 // indexed (LDA ,X)
	RelativeAddressing                                // 8-bit signed offset (BRA)
	RelativeLongAddressing                            // 16-bit signed offset (LBRA)
	RegisterAddressing                                // register pair (TFR/EXG)
	StackAddressing                                   // register bitmask (PSH/PUL)
)

// Typed operand types for safe dispatch in param readers and handlers.
type (
	Immediate8  uint8  // 8-bit immediate value
	Immediate16 uint16 // 16-bit immediate value

	DirectPage uint8  // 8-bit direct page offset
	Extended16 uint16 // 16-bit absolute address

	IndexedAddr uint16 // resolved indexed effective address

	RegisterPair uint8 // TFR/EXG register pair postbyte
	StackMask    uint8 // PSH/PUL register bitmask postbyte
)
