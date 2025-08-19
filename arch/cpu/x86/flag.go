package x86

// Flags represents the x86 processor flags register.
// The flags register is a 16-bit register that contains status flags.
type Flags uint16

// x86 flag bit positions
const (
	FlagCarry      = 0  // CF - Carry flag
	FlagReserved1  = 1  // Reserved, always 1
	FlagParity     = 2  // PF - Parity flag
	FlagReserved3  = 3  // Reserved, always 0
	FlagAuxCarry   = 4  // AF - Auxiliary carry flag
	FlagReserved5  = 5  // Reserved, always 0
	FlagZero       = 6  // ZF - Zero flag
	FlagSign       = 7  // SF - Sign flag
	FlagTrap       = 8  // TF - Trap flag (single step)
	FlagInterrupt  = 9  // IF - Interrupt flag
	FlagDirection  = 10 // DF - Direction flag
	FlagOverflow   = 11 // OF - Overflow flag
	FlagIOPL0      = 12 // IOPL - I/O privilege level bit 0 (80286+)
	FlagIOPL1      = 13 // IOPL - I/O privilege level bit 1 (80286+)
	FlagNested     = 14 // NT - Nested task flag (80286+)
	FlagReserved15 = 15 // Reserved
)

// Flag masks for easy manipulation
const (
	MaskCarry     = 1 << FlagCarry
	MaskParity    = 1 << FlagParity
	MaskAuxCarry  = 1 << FlagAuxCarry
	MaskZero      = 1 << FlagZero
	MaskSign      = 1 << FlagSign
	MaskTrap      = 1 << FlagTrap
	MaskInterrupt = 1 << FlagInterrupt
	MaskDirection = 1 << FlagDirection
	MaskOverflow  = 1 << FlagOverflow
	MaskIOPL      = 3 << FlagIOPL0 // Both IOPL bits
	MaskNested    = 1 << FlagNested
	MaskReserved  = (1 << FlagReserved1) | (1 << FlagReserved15)
)

// Flag accessor methods

// GetCarry returns the carry flag (CF).
func (f Flags) GetCarry() bool {
	return (f & MaskCarry) != 0
}

// GetParity returns the parity flag (PF).
func (f Flags) GetParity() bool {
	return (f & MaskParity) != 0
}

// GetAuxCarry returns the auxiliary carry flag (AF).
func (f Flags) GetAuxCarry() bool {
	return (f & MaskAuxCarry) != 0
}

// GetZero returns the zero flag (ZF).
func (f Flags) GetZero() bool {
	return (f & MaskZero) != 0
}

// GetSign returns the sign flag (SF).
func (f Flags) GetSign() bool {
	return (f & MaskSign) != 0
}

// GetTrap returns the trap flag (TF).
func (f Flags) GetTrap() bool {
	return (f & MaskTrap) != 0
}

// GetInterrupt returns the interrupt flag (IF).
func (f Flags) GetInterrupt() bool {
	return (f & MaskInterrupt) != 0
}

// GetDirection returns the direction flag (DF).
func (f Flags) GetDirection() bool {
	return (f & MaskDirection) != 0
}

// GetOverflow returns the overflow flag (OF).
func (f Flags) GetOverflow() bool {
	return (f & MaskOverflow) != 0
}

// GetIOPL returns the I/O privilege level (IOPL).
func (f Flags) GetIOPL() uint8 {
	return uint8((f & MaskIOPL) >> FlagIOPL0)
}

// GetNested returns the nested task flag (NT).
func (f Flags) GetNested() bool {
	return (f & MaskNested) != 0
}

// Flag setter methods

// SetCarry sets or clears the carry flag (CF).
func (f *Flags) SetCarry(set bool) {
	if set {
		*f |= MaskCarry
	} else {
		*f &= ^Flags(MaskCarry)
	}
}

// SetParity sets or clears the parity flag (PF).
func (f *Flags) SetParity(set bool) {
	if set {
		*f |= MaskParity
	} else {
		*f &= ^Flags(MaskParity)
	}
}

// SetAuxCarry sets or clears the auxiliary carry flag (AF).
func (f *Flags) SetAuxCarry(set bool) {
	if set {
		*f |= MaskAuxCarry
	} else {
		*f &= ^Flags(MaskAuxCarry)
	}
}

// SetZero sets or clears the zero flag (ZF).
func (f *Flags) SetZero(set bool) {
	if set {
		*f |= MaskZero
	} else {
		*f &= ^Flags(MaskZero)
	}
}

// SetSign sets or clears the sign flag (SF).
func (f *Flags) SetSign(set bool) {
	if set {
		*f |= MaskSign
	} else {
		*f &= ^Flags(MaskSign)
	}
}

// SetTrap sets or clears the trap flag (TF).
func (f *Flags) SetTrap(set bool) {
	if set {
		*f |= MaskTrap
	} else {
		*f &= ^Flags(MaskTrap)
	}
}

// SetInterrupt sets or clears the interrupt flag (IF).
func (f *Flags) SetInterrupt(set bool) {
	if set {
		*f |= MaskInterrupt
	} else {
		*f &= ^Flags(MaskInterrupt)
	}
}

// SetDirection sets or clears the direction flag (DF).
func (f *Flags) SetDirection(set bool) {
	if set {
		*f |= MaskDirection
	} else {
		*f &= ^Flags(MaskDirection)
	}
}

// SetOverflow sets or clears the overflow flag (OF).
func (f *Flags) SetOverflow(set bool) {
	if set {
		*f |= MaskOverflow
	} else {
		*f &= ^Flags(MaskOverflow)
	}
}

// SetIOPL sets the I/O privilege level (IOPL).
func (f *Flags) SetIOPL(level uint8) {
	*f &= ^Flags(MaskIOPL)            // Clear current IOPL bits
	*f |= Flags(level&3) << FlagIOPL0 // Set new IOPL bits (mask to 2 bits)
}

// SetNested sets or clears the nested task flag (NT).
func (f *Flags) SetNested(set bool) {
	if set {
		*f |= MaskNested
	} else {
		*f &= ^Flags(MaskNested)
	}
}
