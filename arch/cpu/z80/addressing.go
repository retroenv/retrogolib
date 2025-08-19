package z80

// AddressingMode defines Z80 instruction address modes.
type AddressingMode int

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

type (
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

// Z80 register constants for instruction encoding.
const (
	RegNone RegisterParam = iota

	RegB
	RegC
	RegD
	RegE
	RegH
	RegL
	RegA

	RegBC
	RegDE
	RegHL
	RegSP
	RegAF
	RegIX
	RegIY

	RegHLIndirect
	RegBCIndirect
	RegDEIndirect
	RegSPIndirect
	RegIXIndirect
	RegIYIndirect

	RegImm8
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
}

// String returns the register parameter name.
func (r RegisterParam) String() string {
	if int(r) < len(registerNames) {
		return registerNames[r]
	}
	return "unknown"
}
