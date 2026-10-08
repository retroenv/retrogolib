package sm83

// AddressingMode specifies how an instruction accesses its operands.
type AddressingMode int

const (
	NoAddressing      AddressingMode = 0
	ImpliedAddressing AddressingMode = 1 << iota
	RegisterAddressing
	ImmediateAddressing
	ExtendedAddressing
	RegisterIndirectAddressing
	RelativeAddressing
	BitAddressing
)

type (
	// Register represents an 8-bit CPU register.
	Register uint8
)

// Immediate8 represents 8-bit immediate values.
type Immediate8 uint8

// Immediate16 represents 16-bit immediate values.
type Immediate16 uint16

// Extended represents 16-bit absolute addresses.
type Extended uint16

// Relative represents relative jump addresses.
type Relative int8

// Bit represents bit positions (0-7).
type Bit uint8

// RegisterParam identifies a specific register or addressing mode variant for opcode mapping.
type RegisterParam uint8

// RegisterParam constants.
const (
	RegNone RegisterParam = iota

	// 8-bit general purpose registers
	RegB
	RegC
	RegD
	RegE
	RegH
	RegL
	RegA

	// 16-bit register pairs and special registers.
	RegBC
	RegDE
	RegHL
	RegSP // Stack pointer
	RegAF // Accumulator and flags

	// Indirect addressing through register pairs.
	RegHLIndirect
	RegBCIndirect
	RegDEIndirect

	// RST instruction restart vectors.
	RegRst00 // Call address 0x00
	RegRst08 // Call address 0x08
	RegRst10 // Call address 0x10
	RegRst18 // Call address 0x18
	RegRst20 // Call address 0x20
	RegRst28 // Call address 0x28
	RegRst30 // Call address 0x30
	RegRst38 // Call address 0x38

	// Conditional jump/call flags (SM83 has only 4 conditions).
	RegCondNZ // Non-zero (Z flag clear)
	RegCondZ  // Zero (Z flag set)
	RegCondNC // No carry (C flag clear)
	RegCondC  // Carry (C flag set)

	// Memory load operation variants.
	RegLoadBC  // Load A from (BC)
	RegLoadDE  // Load A from (DE)
	RegLoadHLB // Load B from (HL)
	RegLoadHLC // Load C from (HL)
	RegLoadHLD // Load D from (HL)
	RegLoadHLE // Load E from (HL)
	RegLoadHLH // Load H from (HL)
	RegLoadHLL // Load L from (HL)
	RegLoadHLA // Load A from (HL)

	// LDH high memory variants ($FF00 page).
	RegHighMem       // LDH (n),A — store A to $FF00+n
	RegLoadHighMem   // LDH A,(n) — load A from $FF00+n
	RegCIndirect     // LDH (C),A — store A to $FF00+C
	RegLoadCIndirect // LDH A,(C) — load A from $FF00+C
)

// registerNames provides register parameter string representations.
var registerNames = [...]string{
	RegNone:          "",
	RegB:             "b",
	RegC:             "c",
	RegD:             "d",
	RegE:             "e",
	RegH:             "h",
	RegL:             "l",
	RegA:             "a",
	RegBC:            "bc",
	RegDE:            "de",
	RegHL:            "hl",
	RegSP:            "sp",
	RegAF:            "af",
	RegHLIndirect:    "(hl)",
	RegBCIndirect:    "(bc)",
	RegDEIndirect:    "(de)",
	RegRst00:         "00h",
	RegRst08:         "08h",
	RegRst10:         "10h",
	RegRst18:         "18h",
	RegRst20:         "20h",
	RegRst28:         "28h",
	RegRst30:         "30h",
	RegRst38:         "38h",
	RegCondNZ:        "nz",
	RegCondZ:         "z",
	RegCondNC:        "nc",
	RegCondC:         "c",
	RegLoadBC:        "a,(bc)",
	RegLoadDE:        "a,(de)",
	RegLoadHLB:       "b,(hl)",
	RegLoadHLC:       "c,(hl)",
	RegLoadHLD:       "d,(hl)",
	RegLoadHLE:       "e,(hl)",
	RegLoadHLH:       "h,(hl)",
	RegLoadHLL:       "l,(hl)",
	RegLoadHLA:       "a,(hl)",
	RegHighMem:       "($ff00+n),a",
	RegLoadHighMem:   "a,($ff00+n)",
	RegCIndirect:     "($ff00+c),a",
	RegLoadCIndirect: "a,($ff00+c)",
}

// String returns the register parameter name.
func (r RegisterParam) String() string {
	if int(r) < len(registerNames) {
		return registerNames[r]
	}
	return "unknown"
}
