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
func (c *CPU) SetCarry(value bool) {
	if value {
		c.Flags |= MaskCarry
	} else {
		c.Flags &= ^Flags(MaskCarry)
	}
}

// SetParity sets or clears the parity flag (PF).
func (c *CPU) SetParity(value bool) {
	if value {
		c.Flags |= MaskParity
	} else {
		c.Flags &= ^Flags(MaskParity)
	}
}

// SetAuxCarry sets or clears the auxiliary carry flag (AF).
func (c *CPU) SetAuxCarry(value bool) {
	if value {
		c.Flags |= MaskAuxCarry
	} else {
		c.Flags &= ^Flags(MaskAuxCarry)
	}
}

// SetZero sets or clears the zero flag (ZF).
func (c *CPU) SetZero(value bool) {
	if value {
		c.Flags |= MaskZero
	} else {
		c.Flags &= ^Flags(MaskZero)
	}
}

// SetSign sets or clears the sign flag (SF).
func (c *CPU) SetSign(value bool) {
	if value {
		c.Flags |= MaskSign
	} else {
		c.Flags &= ^Flags(MaskSign)
	}
}

// SetTrap sets or clears the trap flag (TF).
func (c *CPU) SetTrap(value bool) {
	if value {
		c.Flags |= MaskTrap
	} else {
		c.Flags &= ^Flags(MaskTrap)
	}
}

// SetInterrupt sets or clears the interrupt flag (IF).
func (c *CPU) SetInterrupt(value bool) {
	if value {
		c.Flags |= MaskInterrupt
	} else {
		c.Flags &= ^Flags(MaskInterrupt)
	}
}

// SetDirection sets or clears the direction flag (DF).
func (c *CPU) SetDirection(value bool) {
	if value {
		c.Flags |= MaskDirection
	} else {
		c.Flags &= ^Flags(MaskDirection)
	}
}

// SetOverflow sets or clears the overflow flag (OF).
func (c *CPU) SetOverflow(value bool) {
	if value {
		c.Flags |= MaskOverflow
	} else {
		c.Flags &= ^Flags(MaskOverflow)
	}
}

// SetIOPL sets the I/O privilege level (IOPL).
func (c *CPU) SetIOPL(value uint8) {
	c.Flags = (c.Flags &^ MaskIOPL) | (Flags(value&0x3) << FlagIOPL0)
}

// SetNested sets or clears the nested task flag (NT).
func (c *CPU) SetNested(value bool) {
	if value {
		c.Flags |= MaskNested
	} else {
		c.Flags &= ^Flags(MaskNested)
	}
}

// Utility methods for flag calculations

// SetSZP8 sets the Sign, Zero, and Parity flags based on the 8-bit result.
func (c *CPU) SetSZP8(result uint8) {
	c.SetSign((result & 0x80) != 0)
	c.SetZero(result == 0)
	c.SetParity(parity(result))
}

// SetSZP16 sets the Sign, Zero, and Parity flags based on a 16-bit result.
func (c *CPU) SetSZP16(result uint16) {
	c.SetSign((result & 0x8000) != 0)
	c.SetZero(result == 0)
	c.SetParity(parity(uint8(result))) // Parity only considers low byte
}

// parity calculates the parity of a byte (returns true for even parity).
func parity(value uint8) bool {
	// Count set bits
	count := 0
	for i := range 8 {
		if (value & (1 << i)) != 0 {
			count++
		}
	}
	return count%2 == 0
}

// GetFlags returns the flags register as a 16-bit value.
func (c *CPU) GetFlags() uint16 {
	return uint16(c.Flags)
}

// SetFlags sets the flags register from a 16-bit value.
func (c *CPU) SetFlags(value uint16) {
	// Preserve reserved bits
	c.Flags = Flags(value) | MaskReserved
}
