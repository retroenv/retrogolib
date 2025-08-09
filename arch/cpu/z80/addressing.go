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
// Using uint8 values for memory efficiency and performance.
const (
	RegNone RegisterParam = iota // No register / empty

	// 8-bit registers
	RegB
	RegC
	RegD
	RegE
	RegH
	RegL
	RegA

	// 16-bit register pairs
	RegBC
	RegDE
	RegHL
	RegSP
	RegAF
	RegIX
	RegIY

	// Special register references (indirect addressing)
	RegHLIndirect
	RegBCIndirect
	RegDEIndirect
	RegSPIndirect
	RegIXIndirect
	RegIYIndirect

	// Immediate value placeholders
	RegImm8  // 8-bit immediate
	RegImm16 // 16-bit immediate
	RegAddr  // 16-bit address
	RegRel   // relative address

	// Special values for RST instruction
	RegRst00
	RegRst08
	RegRst10
	RegRst18
	RegRst20
	RegRst28
	RegRst30
	RegRst38
)

// RegisterParam represents a register parameter for opcode mapping.
// Using uint8 for memory efficiency and performance.
type RegisterParam uint8

// registerNames provides a lookup table for register parameter string representations.
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
}

// String returns the human-readable representation of the register parameter.
func (r RegisterParam) String() string {
	if int(r) < len(registerNames) {
		return registerNames[r]
	}
	return "unknown"
}
