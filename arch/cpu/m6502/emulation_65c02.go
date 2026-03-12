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
// For indexed modes (zp,X and abs,X), the index register is applied.
func bit65c02(c *CPU, params ...any) error {
	param := params[0]
	// Check if immediate mode (int type from paramReaderImmediate)
	if val, ok := param.(int); ok {
		value := uint8(val)
		c.setZ(value & c.A)
		return nil
	}
	value, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.setV((value>>6)&1 == 1)
	c.setZ(value & c.A)
	c.setN(value)
	return nil
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

// rmbFunc returns a Reset Memory Bit handler for the given bit number (0-7).
// Reads a zero-page byte, clears the specified bit, and writes it back.
// Flags are not affected.
func rmbFunc(bit uint8) func(c *CPU, params ...any) error {
	mask := ^(uint8(1) << bit)
	return func(c *CPU, params ...any) error {
		addr := uint16(params[0].(Absolute))
		v := c.memory.Read(addr)
		c.memory.Write(addr, v&mask)
		return nil
	}
}

// smbFunc returns a Set Memory Bit handler for the given bit number (0-7).
// Reads a zero-page byte, sets the specified bit, and writes it back.
// Flags are not affected.
func smbFunc(bit uint8) func(c *CPU, params ...any) error {
	mask := uint8(1) << bit
	return func(c *CPU, params ...any) error {
		addr := uint16(params[0].(Absolute))
		v := c.memory.Read(addr)
		c.memory.Write(addr, v|mask)
		return nil
	}
}

// bbrFunc returns a Branch on Bit Reset handler for the given bit number (0-7).
// Branches to the target address if the specified bit of the zero-page byte is 0.
func bbrFunc(bit uint8) func(c *CPU, params ...any) error {
	return func(c *CPU, params ...any) error {
		zpAddr := uint16(params[0].(ZeroPage))
		v := c.memory.Read(zpAddr)
		if v&(1<<bit) == 0 {
			c.branch(true, params[1])
		}
		return nil
	}
}

// bbsFunc returns a Branch on Bit Set handler for the given bit number (0-7).
// Branches to the target address if the specified bit of the zero-page byte is 1.
func bbsFunc(bit uint8) func(c *CPU, params ...any) error {
	return func(c *CPU, params ...any) error {
		zpAddr := uint16(params[0].(ZeroPage))
		v := c.memory.Read(zpAddr)
		if v&(1<<bit) != 0 {
			c.branch(true, params[1])
		}
		return nil
	}
}
