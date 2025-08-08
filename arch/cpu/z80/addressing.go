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
	IndexedAddressing                         // (IX+d), (IY+d)
	RelativeAddressing                        // e
	BitAddressing                             // b,r / b,(HL) / b,(IX+d) / b,(IY+d)
	PortAddressing                            // (n), (C)
)

// Register types for Z80 addressing
type (
	// Register8 represents an 8-bit register (A, B, C, D, E, H, L)
	Register8 uint8
	// Register16 represents a 16-bit register (BC, DE, HL, SP, IX, IY)
	Register16 uint16
	// IndexRegister represents IX or IY with displacement
	IndexRegister struct {
		Register uint16 // IX or IY value
		Offset   int8   // displacement (-128 to +127)
	}
)

// Immediate8 indicates that the parameter is an 8-bit immediate value.
type Immediate8 uint8

// Immediate16 indicates that the parameter is a 16-bit immediate value.
type Immediate16 uint16

// Extended indicates that the parameter is a 16-bit absolute address.
type Extended uint16

// RegisterIndirect indicates that the parameter uses register indirect addressing.
type RegisterIndirect uint16

// Indexed indicates that the parameter uses indexed addressing (IX+d) or (IY+d).
type Indexed struct {
	Base   uint16 // IX or IY register value
	Offset int8   // displacement
}

// Relative indicates that the parameter is a relative address for jumps.
type Relative int8

// Bit indicates that the parameter specifies a bit number (0-7).
type Bit uint8

// Port indicates that the parameter is a port address for I/O operations.
type Port uint8
