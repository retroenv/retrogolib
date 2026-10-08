package cpu65816

import "fmt"

// Branch and jump instructions.

// branch performs a relative branch to the precomputed target when taken.
// A taken branch adds one cycle. In emulation mode it adds one more cycle
// when the target is in a different page than the next instruction.
func (c *CPU) branch(taken bool, params []any) error {
	target, err := operand[BranchTarget](params)
	if err != nil {
		return err
	}
	if !taken {
		return nil
	}

	nextPC := c.PC + 2
	c.cycles++
	if c.E && uint16(target)&0xFF00 != nextPC&0xFF00 {
		c.cycles++
	}
	c.PC = uint16(target)
	c.pcChanged = true
	return nil
}

func bcc(c *CPU, params ...any) error {
	return c.branch(c.Flags.C == 0, params)
}

func bcs(c *CPU, params ...any) error {
	return c.branch(c.Flags.C != 0, params)
}

func beq(c *CPU, params ...any) error {
	return c.branch(c.Flags.Z != 0, params)
}

func bmi(c *CPU, params ...any) error {
	return c.branch(c.Flags.N != 0, params)
}

func bne(c *CPU, params ...any) error {
	return c.branch(c.Flags.Z == 0, params)
}

func bpl(c *CPU, params ...any) error {
	return c.branch(c.Flags.N == 0, params)
}

func bra(c *CPU, params ...any) error {
	return c.branch(true, params)
}

// brl - Branch Long: always taken, 16-bit offset already resolved to the target.
func brl(c *CPU, params ...any) error {
	target, err := operand[BranchTarget](params)
	if err != nil {
		return err
	}
	c.PC = uint16(target)
	c.pcChanged = true
	return nil
}

func bvc(c *CPU, params ...any) error {
	return c.branch(c.Flags.V == 0, params)
}

func bvs(c *CPU, params ...any) error {
	return c.branch(c.Flags.V != 0, params)
}

// jmp - Jump (same bank).
func jmp(c *CPU, params ...any) error {
	if len(params) == 0 {
		return ErrMissingParameter
	}
	switch p := params[0].(type) {
	case Absolute16:
		c.PC = uint16(p)
	case AbsIndirect:
		c.PC = uint16(p)
	case AbsIndirectX:
		c.PC = uint16(p)
	default:
		return fmt.Errorf("%w: jump target type %T", ErrInvalidParameterType, params[0])
	}
	c.pcChanged = true
	return nil
}

// jml - Jump Long (sets PB).
func jml(c *CPU, params ...any) error {
	target, err := operand[AbsLong](params)
	if err != nil {
		return err
	}
	c.PB = uint8(uint32(target) >> 16)
	c.PC = uint16(target)
	c.pcChanged = true
	return nil
}

// jsr - Jump to Subroutine.
// The return address is the last byte of the 3-byte instruction (PC+2).
func jsr(c *CPU, params ...any) error {
	if len(params) == 0 {
		return ErrMissingParameter
	}
	var target uint16
	switch p := params[0].(type) {
	case Absolute16:
		target = uint16(p)
	case AbsIndirectX:
		target = uint16(p)
	default:
		return fmt.Errorf("%w: subroutine target type %T", ErrInvalidParameterType, params[0])
	}

	c.push16(c.PC + 2)
	c.PC = target
	c.pcChanged = true
	return nil
}

// jsl - Jump to Subroutine Long.
// The return address is the last byte of the 4-byte instruction (PC+3).
// 65816-native: uses full 16-bit SP (no page-1 wrap between bytes).
func jsl(c *CPU, params ...any) error {
	target, err := operand[AbsLong](params)
	if err != nil {
		return err
	}

	c.push8raw(c.PB)
	c.push16raw(c.PC + 3)
	c.fixEmuSP()
	c.PB = uint8(uint32(target) >> 16)
	c.PC = uint16(target)
	c.pcChanged = true
	return nil
}
