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
	RegB // B register - general purpose 8-bit
	RegC // C register - general purpose 8-bit
	RegD // D register - general purpose 8-bit
	RegE // E register - general purpose 8-bit
	RegH // H register - high byte of HL register pair
	RegL // L register - low byte of HL register pair
	RegA // A register - accumulator, primary register for arithmetic operations

	// 16-bit register pairs
	RegBC // BC register pair - combination of B and C registers
	RegDE // DE register pair - combination of D and E registers
	RegHL // HL register pair - combination of H and L registers, commonly used as memory pointer
	RegSP // SP register - stack pointer, points to top of stack
	RegAF // AF register pair - combination of A (accumulator) and F (flags) registers
	RegIX // IX register - 16-bit index register for indexed addressing
	RegIY // IY register - 16-bit index register for indexed addressing

	// Special register references (indirect addressing)
	RegHLIndirect // (HL) - memory location pointed to by HL register pair
	RegBCIndirect // (BC) - memory location pointed to by BC register pair
	RegDEIndirect // (DE) - memory location pointed to by DE register pair
	RegSPIndirect // (SP) - memory location pointed to by SP register
	RegIXIndirect // (IX) - memory location pointed to by IX register
	RegIYIndirect // (IY) - memory location pointed to by IY register

	// Immediate value placeholders
	RegImm8  // n - 8-bit immediate value
	RegImm16 // nn - 16-bit immediate value
	RegAddr  // (nn) - 16-bit absolute address
	RegRel   // e - relative address for branch instructions

	// Special values for RST instruction (restart vectors)
	RegRst00 // RST 00H - restart at address 0x00
	RegRst08 // RST 08H - restart at address 0x08
	RegRst10 // RST 10H - restart at address 0x10
	RegRst18 // RST 18H - restart at address 0x18
	RegRst20 // RST 20H - restart at address 0x20
	RegRst28 // RST 28H - restart at address 0x28
	RegRst30 // RST 30H - restart at address 0x30
	RegRst38 // RST 38H - restart at address 0x38
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
