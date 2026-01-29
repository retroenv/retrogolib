package z80

// AddressingMode specifies how an instruction accesses its operands.
// Multiple modes can be combined using bitwise OR for instructions that support variants.
//
// Available modes:
//   - ImpliedAddressing: No operands (NOP, HALT)
//   - RegisterAddressing: Direct register access (LD A,B)
//   - ImmediateAddressing: Constant value (LD A,n)
//   - ExtendedAddressing: Absolute memory address (LD A,(nn))
//   - RegisterIndirectAddressing: Memory via register pointer (LD A,(HL))
//   - RelativeAddressing: PC-relative offset (JR e)
//   - BitAddressing: Bit manipulation (BIT n,r)
//   - PortAddressing: I/O port access (IN A,(n))
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
	PortAddressing
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

// RegisterIndirect represents register indirect addressing.
type RegisterIndirect uint16

// Relative represents relative jump addresses.
type Relative int8

// Bit represents bit positions (0-7).
type Bit uint8

// Port represents I/O port addresses.
type Port uint8

// RegisterParam constants are used as keys in opcode lookup maps to identify
// specific instruction variants. These enable bidirectional mapping between
// opcodes and instructions for disassembly and code generation.
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

	// 16-bit register pairs and special registers
	RegBC
	RegDE
	RegHL
	RegSP // Stack pointer
	RegAF // Accumulator and flags
	RegIX // Index register X
	RegIY // Index register Y

	// Indirect addressing through register pairs
	RegHLIndirect
	RegBCIndirect
	RegDEIndirect
	RegSPIndirect
	RegIXIndirect
	RegIYIndirect

	// Immediate values and addressing
	RegImm8
	RegImm16 // 16-bit immediate value
	RegAddr  // Absolute memory address
	RegRel   // PC-relative offset for branches

	// RST instruction restart vectors
	RegRst00 // Call address 0x00
	RegRst08 // Call address 0x08
	RegRst10 // Call address 0x10
	RegRst18 // Call address 0x18
	RegRst20 // Call address 0x20
	RegRst28 // Call address 0x28
	RegRst30 // Call address 0x30
	RegRst38 // Call address 0x38

	// Conditional jump/call flags
	RegCondNZ // Non-zero (Z flag clear)
	RegCondZ  // Zero (Z flag set)
	RegCondNC // No carry (C flag clear)
	RegCondC  // Carry (C flag set)
	RegCondPO // Parity odd (P/V flag clear)
	RegCondPE // Parity even (P/V flag set)
	RegCondP  // Positive (S flag clear)
	RegCondM  // Negative (S flag set)

	// Memory load operation variants
	RegLoadBC    // Load A from (BC)
	RegLoadDE    // Load A from (DE)
	RegLoadHLB   // Load B from (HL)
	RegLoadHLC   // Load C from (HL)
	RegLoadHLD   // Load D from (HL)
	RegLoadHLE   // Load E from (HL)
	RegLoadHLH   // Load H from (HL)
	RegLoadHLL   // Load L from (HL)
	RegLoadHLA   // Load A from (HL)
	RegLoadExtHL // Load HL from (nn)
	RegLoadExtA  // Load A from (nn)
	RegStoreExtA // Store A to (nn)

	// Z80-specific special registers
	RegI // Interrupt vector base
	RegR // Memory refresh counter

	// Interrupt modes
	RegIM0 // Mode 0: 8080-compatible
	RegIM1 // Mode 1: RST 38H
	RegIM2 // Mode 2: Vectored interrupts
)

// RegisterParam identifies a specific register or addressing mode variant for opcode mapping.
type RegisterParam uint8

// registerNames provides register parameter string representations.
var registerNames = [...]string{
	RegNone:       "",
	RegB:          "b",
	RegC:          "c",
	RegD:          "d",
	RegE:          "e",
	RegH:          "h",
	RegL:          "l",
	RegA:          "a",
	RegBC:         "bc",
	RegDE:         "de",
	RegHL:         "hl",
	RegSP:         "sp",
	RegAF:         "af",
	RegIX:         "ix",
	RegIY:         "iy",
	RegHLIndirect: "(hl)",
	RegBCIndirect: "(bc)",
	RegDEIndirect: "(de)",
	RegSPIndirect: "(sp)",
	RegIXIndirect: "(ix)",
	RegIYIndirect: "(iy)",
	RegImm8:       "n",
	RegImm16:      "nn",
	RegAddr:       "(nn)",
	RegRel:        "e",
	RegRst00:      "00h",
	RegRst08:      "08h",
	RegRst10:      "10h",
	RegRst18:      "18h",
	RegRst20:      "20h",
	RegRst28:      "28h",
	RegRst30:      "30h",
	RegRst38:      "38h",
	RegCondNZ:     "nz",
	RegCondZ:      "z",
	RegCondNC:     "nc",
	RegCondC:      "c",
	RegCondPO:     "po",
	RegCondPE:     "pe",
	RegCondP:      "p",
	RegCondM:      "m",
	RegLoadBC:     "a,(bc)",
	RegLoadDE:     "a,(de)",
	RegLoadHLB:    "b,(hl)",
	RegLoadHLC:    "c,(hl)",
	RegLoadHLD:    "d,(hl)",
	RegLoadHLE:    "e,(hl)",
	RegLoadHLH:    "h,(hl)",
	RegLoadHLL:    "l,(hl)",
	RegLoadHLA:    "a,(hl)",
	RegLoadExtHL:  "hl,(nn)",
	RegLoadExtA:   "a,(nn)",
	RegStoreExtA:  "(nn),a",
	RegI:          "i",
	RegR:          "r",
	RegIM0:        "0",
	RegIM1:        "1",
	RegIM2:        "2",
}

// String returns the register parameter name.
func (r RegisterParam) String() string {
	if int(r) < len(registerNames) {
		return registerNames[r]
	}
	return "unknown"
}
