package z80

// AddressingMode defines an address mode for Z80 instructions.
type AddressingMode int

// Z80 addressing modes.
const (
	NoAddressing               AddressingMode = 0
	ImpliedAddressing          AddressingMode = 1 << iota
	RegisterAddressing                        // r, rr
	ImmediateAddressing                       // n, nn
	ExtendedAddressing                        // (nn)
	RegisterIndirectAddressing                // (rr), (r)
	RelativeAddressing                        // e
	BitAddressing                             // b,r / b,(HL) / b,(IX+d) / b,(IY+d)
	PortAddressing                            // (n), (C)
)

// Register types for Z80 addressing
type (
	// Register represents an 8-bit register (A, B, C, D, E, H, L)
	Register uint8
)

// Immediate8 indicates that the parameter is an 8-bit immediate value.
type Immediate8 uint8

// Immediate16 indicates that the parameter is a 16-bit immediate value.
type Immediate16 uint16

// Extended indicates that the parameter is a 16-bit absolute address.
type Extended uint16

// RegisterIndirect indicates that the parameter uses register indirect addressing.
type RegisterIndirect uint16

// Relative indicates that the parameter is a relative address for jumps.
type Relative int8

// Bit indicates that the parameter specifies a bit number (0-7).
type Bit uint8

// Port indicates that the parameter is a port address for I/O operations.
type Port uint8

// Z80 Register constants for opcode mapping
// These are used as keys in the RegisterOpcodes map to differentiate
// between instructions that target different registers.
const (
	// 8-bit registers
	RegB RegisterParam = "b"
	RegC RegisterParam = "c"
	RegD RegisterParam = "d"
	RegE RegisterParam = "e"
	RegH RegisterParam = "h"
	RegL RegisterParam = "l"
	RegA RegisterParam = "a"
	
	// 16-bit register pairs
	RegBC RegisterParam = "bc"
	RegDE RegisterParam = "de"
	RegHL RegisterParam = "hl"
	RegSP RegisterParam = "sp"
	RegAF RegisterParam = "af"
	RegIX RegisterParam = "ix"
	RegIY RegisterParam = "iy"
	
	// Special register references
	RegHLIndirect RegisterParam = "(hl)"
	RegBCIndirect RegisterParam = "(bc)"
	RegDEIndirect RegisterParam = "(de)"
	RegSPIndirect RegisterParam = "(sp)"
	RegIXIndirect RegisterParam = "(ix)"
	RegIYIndirect RegisterParam = "(iy)"
	
	// Immediate value placeholders
	RegImm8  RegisterParam = "n"   // 8-bit immediate
	RegImm16 RegisterParam = "nn"  // 16-bit immediate
	RegAddr  RegisterParam = "(nn)" // 16-bit address
	RegRel   RegisterParam = "e"    // relative address
	
	// Special values for RST instruction
	RegRst00 RegisterParam = "00h"
	RegRst08 RegisterParam = "08h"
	RegRst10 RegisterParam = "10h"
	RegRst18 RegisterParam = "18h"
	RegRst20 RegisterParam = "20h"
	RegRst28 RegisterParam = "28h"
	RegRst30 RegisterParam = "30h"
	RegRst38 RegisterParam = "38h"
)

// RegisterParam represents a register parameter for opcode mapping.
type RegisterParam string
