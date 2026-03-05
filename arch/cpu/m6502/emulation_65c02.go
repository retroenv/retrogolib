// This file contains 65C02 specific instruction handlers.

package m6502

import "fmt"

// bra - Branch Always.
func bra(c *CPU, params ...any) error {
	c.branch(true, params[0])
	return nil
}

// bit65c02 - Bit Test (65C02 extended).
// For immediate mode, only Z flag is affected (V and N are not modified).
func bit65c02(c *CPU, params ...any) error {
	param := params[0]
	// Check if immediate mode (int type from paramReaderImmediate)
	if val, ok := param.(int); ok {
		value := uint8(val)
		c.setZ(value & c.A)
		return nil
	}
	// For non-immediate modes, behave like standard BIT
	return bit(c, params...)
}

// dec65c02 - Decrement (65C02 with accumulator mode).
func dec65c02(c *CPU, params ...any) error {
	if hasAccumulatorParam(params...) {
		c.A--
		c.setZN(c.A)
		return nil
	}
	return dec(c, params...)
}

// inc65c02 - Increment (65C02 with accumulator mode).
func inc65c02(c *CPU, params ...any) error {
	if hasAccumulatorParam(params...) {
		c.A++
		c.setZN(c.A)
		return nil
	}
	return inc(c, params...)
}

// jmp65c02 - Jump (65C02 with absolute indexed indirect).
func jmp65c02(c *CPU, params ...any) error {
	param := params[0]
	switch address := param.(type) {
	case Absolute:
		c.PC = uint16(address)
	case Indirect:
		// 65C02 fixes the page boundary bug
		c.PC = c.memory.ReadWord(uint16(address))
	case AbsoluteXIndirect:
		// The resolved address is in the second parameter
		addr, ok := params[1].(Absolute)
		if !ok {
			return fmt.Errorf("%w: jmp (abs,X) resolved type %T", ErrInvalidParameterType, params[1])
		}
		c.PC = uint16(addr)
	default:
		return fmt.Errorf("%w: jmp mode type %T", ErrUnsupportedAddressingMode, param)
	}
	return nil
}

// phx - Push X Register.
func phx(c *CPU) error {
	c.push(c.X)
	return nil
}

// phy - Push Y Register.
func phy(c *CPU) error {
	c.push(c.Y)
	return nil
}

// plx - Pull X Register.
func plx(c *CPU) error {
	c.X = c.pop()
	c.setZN(c.X)
	return nil
}

// ply - Pull Y Register.
func ply(c *CPU) error {
	c.Y = c.pop()
	c.setZN(c.Y)
	return nil
}

// stz - Store Zero to memory.
func stz(c *CPU, params ...any) error {
	return c.memory.WriteAddressModes(0, params...)
}

// trb - Test and Reset Bits.
// AND A with memory to set Z flag, then clear bits in memory where A has 1s.
func trb(c *CPU, params ...any) error {
	value, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.setZ(value & c.A)
	result := value & ^c.A
	return c.memory.WriteAddressModes(result, params...)
}

// tsb - Test and Set Bits.
// AND A with memory to set Z flag, then set bits in memory where A has 1s.
func tsb(c *CPU, params ...any) error {
	value, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.setZ(value & c.A)
	result := value | c.A
	return c.memory.WriteAddressModes(result, params...)
}
