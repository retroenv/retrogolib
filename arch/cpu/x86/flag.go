package x86

// x86 flag bit positions
const (
	FlagCarry     = 0  // CF - Carry flag
	FlagParity    = 2  // PF - Parity flag
	FlagAuxCarry  = 4  // AF - Auxiliary carry flag
	FlagZero      = 6  // ZF - Zero flag
	FlagSign      = 7  // SF - Sign flag
	FlagTrap      = 8  // TF - Trap flag (single step)
	FlagInterrupt = 9  // IF - Interrupt flag
	FlagDirection = 10 // DF - Direction flag
	FlagOverflow  = 11 // OF - Overflow flag
	FlagIOPL0     = 12 // IOPL - I/O privilege level bit 0 (80286+)
	FlagIOPL1     = 13 // IOPL - I/O privilege level bit 1 (80286+)
	FlagNested    = 14 // NT - Nested task flag (80286+)
)
