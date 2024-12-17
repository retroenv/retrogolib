package chip8

import "fmt"

// cls clears the display.
func cls(c *CPU, _ uint16) error {
	c.Display = [displayWidth * displayHeight]bool{}
	c.RedrawScreen = true
	c.PC += 2
	return nil
}

// ret returns from a subroutine.
func ret(c *CPU, _ uint16) error {
	c.PC = c.Stack[c.SP]
	c.SP--
	c.PC += 2
	return nil
}

// jp jumps to an address and optionally adds V0 to the address.
func jp(c *CPU, param uint16) error {
	mode := (param & 0xF000) >> 12
	addr := param & 0x0FFF

	switch mode {
	case 0x1: // JP addr
		c.PC = addr
	case 0xb: // JP V0, addr
		c.PC = addr + uint16(c.V[0])
	default:
		return fmt.Errorf("invalid mode for jp: %04X", mode)
	}

	return nil
}

// call calls a subroutine.
func call(c *CPU, param uint16) error {
	c.SP++
	c.Stack[c.SP] = c.PC
	c.PC = param & 0x0FFF
	return nil
}

// se skips the next instruction if the register equals a value/register.
func se(c *CPU, param uint16) error {
	mode := (param & 0xF000) >> 12
	reg := (param & 0x0F00) >> 8

	switch mode {
	case 0x3: // SE Vx, byte
		value := byte(param & 0x00FF)
		c.updatePC(c.V[reg] == value)

	case 0x5: // SE Vx, Vy
		reg2 := (param & 0x00F0) >> 4
		c.updatePC(c.V[reg] == c.V[reg2])

	default:
		return fmt.Errorf("invalid mode for se: %04X", mode)
	}
	return nil
}

// sne skips the next instruction if the register does not equal a value/register.
func sne(c *CPU, param uint16) error {
	mode := (param & 0xF000) >> 12
	reg := (param & 0x0F00) >> 8

	switch mode {
	case 0x4: // SNE Vx, byte
		value := byte(param & 0x00FF)
		c.updatePC(c.V[reg] != value)

	case 0x9: // SNE Vx, Vy
		reg2 := (param & 0x00F0) >> 4
		c.updatePC(c.V[reg] != c.V[reg2])

	default:
		return fmt.Errorf("invalid mode for sne: %04X", mode)
	}
	return nil
}

// or performs a bitwise OR operation on two registers.
func or(c *CPU, param uint16) error {
	reg1 := (param & 0x0F00) >> 8
	reg2 := (param & 0x00F0) >> 4
	c.V[reg1] |= c.V[reg2]
	c.PC += 2
	return nil
}

// xor performs a bitwise XOR operation on two registers.
func xor(c *CPU, param uint16) error {
	reg1 := (param & 0x0F00) >> 8
	reg2 := (param & 0x00F0) >> 4
	c.V[reg1] ^= c.V[reg2]
	c.PC += 2
	return nil
}

// add adds a value/register to a register.
func add(c *CPU, param uint16) error {
	mode := (param & 0xF000) >> 12
	reg := (param & 0x0F00) >> 8
	value := byte(param & 0x00FF)

	switch {
	case mode == 0x7: // ADD Vx, byte
		c.V[reg] += value

	case mode == 0x8: // ADD Vx, Vy
		reg2 := (param & 0x00F0) >> 4
		c.V[reg] += c.V[reg2]

	case mode == 0xf && value == 0x1e: // ADD I, Vx
		c.I += uint16(c.V[reg])

	default:
		return fmt.Errorf("invalid mode for add: %04X", mode)
	}

	c.PC += 2
	return nil
}

// sub subtracts a value/register from a register.
func sub(c *CPU, param uint16) error {
	reg1 := (param & 0x0F00) >> 8
	reg2 := (param & 0x00F0) >> 4

	if c.V[reg1] > c.V[reg2] {
		c.V[0xf] = 1
	} else {
		c.V[0xf] = 0
	}

	c.V[reg1] -= c.V[reg2]

	c.PC += 2
	return nil
}

// ld loads a value/register into a register.
func ld(c *CPU, param uint16) error {
	mode := (param & 0xF000) >> 12
	reg := (param & 0x0F00) >> 8
	value := byte(param & 0x00FF)

	switch mode {
	case 0x6: // LD Vx, byte
		c.V[reg] = value

	case 0x8: // LD Vx, Vy
		reg2 := (param & 0x00F0) >> 4
		c.V[reg] = c.V[reg2]

	case 0xa: // LD I, addr
		addr := param & 0x0FFF
		c.I = addr

	case 0xf:
		return ldF(c, param)

	default:
		return fmt.Errorf("invalid mode for ld: %04X", mode)
	}

	c.PC += 2
	return nil
}

func ldF(c *CPU, param uint16) error {
	value := byte(param & 0x00FF)

	switch value {
	case 0x07: // LD Vx, DT
		reg := (param & 0x0F00) >> 8
		c.V[reg] = c.DelayTimer

	case 0x0a: // LD Vx, K
		// TODO
		panic("not implemented")

	case 0x15: // LD DT, Vx
		reg := (param & 0x0F00) >> 8
		c.DelayTimer = c.V[reg]

	case 0x18: // LD ST, Vx
		reg := (param & 0x0F00) >> 8
		c.SoundTimer = c.V[reg]

	case 0x29: // LD F, Vx
		// TODO
		panic("not implemented")

	case 0x33: // LD B, Vx
		// TODO
		panic("not implemented")

	case 0x55: // LD [I], Vx
		reg := (param & 0x0F00) >> 8
		for i := uint16(0); i <= reg; i++ {
			c.Memory[c.I+i] = c.V[i]
		}

	case 0x65: // LD Vx, [I]
		reg := (param & 0x0F00) >> 8
		for i := uint16(0); i <= reg; i++ {
			c.V[i] = c.Memory[c.I+i]
		}

	default:
		return fmt.Errorf("invalid value for ldF: %04X", value)
	}

	c.PC += 2
	return nil
}

// and performs a bitwise AND operation on two registers.
func and(c *CPU, param uint16) error {
	reg1 := (param & 0x0F00) >> 8
	reg2 := (param & 0x00F0) >> 4
	c.V[reg1] &= c.V[reg2]
	c.PC += 2
	return nil
}

// drw displays n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
func drw(c *CPU, param uint16) error {
	x := uint16(c.V[(param&0x0F00)>>8])
	y := uint16(c.V[(param&0x00F0)>>4])
	n := param & 0x000F

	c.V[0xf] = 0
	for yline := range n {
		pixel := c.Memory[c.I+yline]
		for xline := range uint16(8) {
			if (pixel & (0x80 >> xline)) != 0 {
				if c.Display[(x+xline)+(y+yline)*displayWidth] {
					c.V[0xf] = 1
				}
				c.Display[(x+xline)+(y+yline)*displayWidth] = !c.Display[(x+xline)+(y+yline)*displayWidth]
			}
		}
	}

	c.RedrawScreen = true
	c.PC += 2
	return nil
}

// rnd generates a random number and performs a bitwise AND operation on it.
func rnd(c *CPU, param uint16) error {
	reg := (param & 0x0F00) >> 8
	value := byte(param & 0x00FF)
	c.V[reg] = byte(c.rnd.Int63()) & value
	c.PC += 2
	return nil
}

// shl shifts a register left by one.
func shl(c *CPU, param uint16) error {
	reg := (param & 0x0F00) >> 8
	c.V[0xf] = c.V[reg] >> 7
	c.V[reg] <<= 1
	c.PC += 2
	return nil
}

// shr shifts a register right by one.
func shr(c *CPU, param uint16) error {
	reg := (param & 0x0F00) >> 8
	c.V[0xf] = c.V[reg] & 0x1
	c.V[reg] >>= 1
	c.PC += 2
	return nil
}

// skp skips the next instruction if the key with the value of Vx is pressed.
func skp(c *CPU, param uint16) error {
	reg := (param & 0x0F00) >> 8
	if c.Key[c.V[reg]] {
		c.PC += 4
	} else {
		c.PC += 2
	}
	return nil
}

// sknp skips the next instruction if the key with the value of Vx is not pressed.
func sknp(c *CPU, param uint16) error {
	reg := (param & 0x0F00) >> 8
	if !c.Key[c.V[reg]] {
		c.PC += 4
	} else {
		c.PC += 2
	}
	return nil
}

// subn subtracts a register from another register
func subn(c *CPU, param uint16) error {
	reg1 := (param & 0x0F00) >> 8
	reg2 := (param & 0x00F0) >> 4

	if c.V[reg2] > c.V[reg1] {
		c.V[0xf] = 1
	} else {
		c.V[0xf] = 0
	}

	c.V[reg1] = c.V[reg2] - c.V[reg1]

	c.PC += 2
	return nil
}
