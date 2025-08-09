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
