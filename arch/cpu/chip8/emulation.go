package chip8

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// cls clears the display.
func cls(c *CPU, _ uint16) error {
	for i := range c.Display {
		c.Display[i] = 0
	}
	c.RedrawScreen = true
	c.PC += 2
	return nil
}

// ret returns from a subroutine.
func ret(c *CPU, _ uint16) error {
	if c.SP == 0 {
		return fmt.Errorf("%w: cannot return from subroutine", ErrStackUnderflow)
	}
	c.SP--
	c.PC = c.Stack[c.SP]
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
	if c.SP >= uint8(len(c.Stack)) {
		return fmt.Errorf("%w: cannot call subroutine", ErrStackOverflow)
	}
	c.Stack[c.SP] = c.PC
	c.SP++
	c.PC = param & 0x0FFF
	return nil
}

// se skips the next instruction if the register equals a value/register.
func se(c *CPU, param uint16) error {
	mode := (param & 0xF000) >> 12
	reg := (param & 0x0F00) >> 8
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}

	switch mode {
	case 0x3: // SE Vx, byte
		value := byte(param & 0x00FF)
		c.updatePC(c.V[reg] == value)

	case 0x5: // SE Vx, Vy
		reg2 := (param & 0x00F0) >> 4
		if reg2 > 15 {
			return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg2)
		}
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
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}

	switch mode {
	case 0x4: // SNE Vx, byte
		value := byte(param & 0x00FF)
		c.updatePC(c.V[reg] != value)

	case 0x9: // SNE Vx, Vy
		reg2 := (param & 0x00F0) >> 4
		if reg2 > 15 {
			return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg2)
		}
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
	if reg1 > 15 || reg2 > 15 {
		return fmt.Errorf("%w: 0x%X, 0x%X", ErrRegisterOutOfBounds, reg1, reg2)
	}
	c.V[reg1] |= c.V[reg2]
	c.PC += 2
	return nil
}

// xor performs a bitwise XOR operation on two registers.
func xor(c *CPU, param uint16) error {
	reg1 := (param & 0x0F00) >> 8
	reg2 := (param & 0x00F0) >> 4
	if reg1 > 15 || reg2 > 15 {
		return fmt.Errorf("%w: 0x%X, 0x%X", ErrRegisterOutOfBounds, reg1, reg2)
	}
	c.V[reg1] ^= c.V[reg2]
	c.PC += 2
	return nil
}

// add adds a value/register to a register.
func add(c *CPU, param uint16) error {
	mode := (param & 0xF000) >> 12
	reg := (param & 0x0F00) >> 8
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}
	value := byte(param & 0x00FF)

	switch {
	case mode == 0x7: // ADD Vx, byte
		c.V[reg] += value

	case mode == 0x8: // ADD Vx, Vy
		reg2 := (param & 0x00F0) >> 4
		if reg2 > 15 {
			return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg2)
		}

		if uint16(c.V[reg])+uint16(c.V[reg2]) > math.MaxUint8 {
			c.V[0xf] = 1
		} else {
			c.V[0xf] = 0
		}

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
	if reg1 > 15 || reg2 > 15 {
		return fmt.Errorf("%w: 0x%X, 0x%X", ErrRegisterOutOfBounds, reg1, reg2)
	}

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
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}
	value := byte(param & 0x00FF)

	switch mode {
	case 0x6: // LD Vx, byte
		c.V[reg] = value

	case 0x8: // LD Vx, Vy
		reg2 := (param & 0x00F0) >> 4
		if reg2 > 15 {
			return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg2)
		}
		c.V[reg] = c.V[reg2]

	case 0xa: // LD I, addr
		c.I = param & 0x0FFF

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
	reg := (param & 0x0F00) >> 8
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}

	switch value {
	case 0x07: // LD Vx, DT
		c.V[reg] = c.DelayTimer
	case 0x0a: // LD Vx, K
		return c.ldVxK(reg)
	case 0x15: // LD DT, Vx
		c.DelayTimer = c.V[reg]
	case 0x18: // LD ST, Vx
		c.SoundTimer = c.V[reg]
	case 0x29: // LD F, Vx
		return c.ldFVx(reg)
	case 0x33: // LD B, Vx
		return c.ldBVx(reg)
	case 0x55: // LD [I], Vx
		return c.ldIVx(reg)
	case 0x65: // LD Vx, [I]
		return c.ldVxI(reg)
	default:
		return fmt.Errorf("invalid value for ldF: %04X", value)
	}

	c.PC += 2
	return nil
}

// ldVxK implements LD Vx, K instruction (wait for key press)
func (c *CPU) ldVxK(reg uint16) error {
	keyPressed := -1
	for i, isKeyPressed := range c.Key {
		if isKeyPressed {
			keyPressed = i
			break
		}
	}
	if keyPressed == -1 {
		return nil // do not update program counter and wait for a key press
	}
	c.V[reg] = byte(keyPressed)
	c.PC += 2
	return nil
}

// ldFVx implements LD F, Vx instruction (set I to font location)
func (c *CPU) ldFVx(reg uint16) error {
	fontIndex := c.V[reg]
	if fontIndex > 15 {
		return fmt.Errorf("%w: 0x%X", ErrFontIndexOutOfBounds, fontIndex)
	}
	c.I = uint16(fontIndex) * 0x5
	c.PC += 2
	return nil
}

// ldBVx implements LD B, Vx instruction (store BCD representation)
func (c *CPU) ldBVx(reg uint16) error {
	if c.I+2 >= uint16(len(c.Memory)) {
		return fmt.Errorf("%w: I=0x%03X", ErrMemoryOutOfBounds, c.I)
	}
	bcd := c.V[reg]
	for i := 2; i >= 0; i-- {
		c.Memory[c.I+uint16(i)] = bcd % 10
		bcd /= 10
	}
	c.PC += 2
	return nil
}

// ldIVx implements LD [I], Vx instruction (store registers V0 through Vx in memory)
func (c *CPU) ldIVx(reg uint16) error {
	if c.I+reg >= uint16(len(c.Memory)) {
		return fmt.Errorf("%w: I=0x%03X, reg=0x%X", ErrMemoryOutOfBounds, c.I, reg)
	}
	for i := uint16(0); i <= reg; i++ {
		c.Memory[c.I+i] = c.V[i]
	}
	c.PC += 2
	return nil
}

// ldVxI implements LD Vx, [I] instruction (read registers V0 through Vx from memory)
func (c *CPU) ldVxI(reg uint16) error {
	if c.I+reg >= uint16(len(c.Memory)) {
		return fmt.Errorf("%w: I=0x%03X, reg=0x%X", ErrMemoryOutOfBounds, c.I, reg)
	}
	for i := uint16(0); i <= reg; i++ {
		c.V[i] = c.Memory[c.I+i]
	}
	c.PC += 2
	return nil
}

// and performs a bitwise AND operation on two registers.
func and(c *CPU, param uint16) error {
	reg1 := (param & 0x0F00) >> 8
	reg2 := (param & 0x00F0) >> 4
	if reg1 > 15 || reg2 > 15 {
		return fmt.Errorf("%w: 0x%X, 0x%X", ErrRegisterOutOfBounds, reg1, reg2)
	}
	c.V[reg1] &= c.V[reg2]
	c.PC += 2
	return nil
}

// drw displays n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
func drw(c *CPU, param uint16) error {
	reg1 := (param & 0x0F00) >> 8
	reg2 := (param & 0x00F0) >> 4
	if reg1 > 15 || reg2 > 15 {
		return fmt.Errorf("%w: 0x%X, 0x%X", ErrRegisterOutOfBounds, reg1, reg2)
	}

	x := uint16(c.V[reg1]) % displayWidth
	y := uint16(c.V[reg2]) % displayHeight
	height := param & 0x000F

	if c.I+height-1 >= uint16(len(c.Memory)) {
		return fmt.Errorf("%w: sprite I=0x%03X, height=%d", ErrMemoryOutOfBounds, c.I, height)
	}

	c.V[0xf] = 0

	for yLine := range height {
		if y+yLine >= displayHeight {
			break // Stop drawing if we go past screen boundary
		}
		sprite := c.Memory[c.I+yLine]

		for xLine := range uint16(8) {
			if x+xLine >= displayWidth {
				break // Stop drawing if we go past screen boundary
			}
			if (sprite & (0x80 >> xLine)) != 0 {
				index := (x + xLine) + (y+yLine)*displayWidth
				if index >= uint16(len(c.Display)) {
					return fmt.Errorf("%w: %d", ErrDisplayOutOfBounds, index)
				}
				if c.Display[index] == 1 {
					c.V[0xf] = 1
				}
				c.Display[index] ^= 1
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
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}
	value := byte(param & 0x00FF)
	c.V[reg] = byte(rand.IntN(256)) & value
	c.PC += 2
	return nil
}

// shl shifts a register left by one.
func shl(c *CPU, param uint16) error {
	reg := (param & 0x0F00) >> 8
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}
	c.V[0xf] = c.V[reg] >> 7
	c.V[reg] <<= 1
	c.PC += 2
	return nil
}

// shr shifts a register right by one.
func shr(c *CPU, param uint16) error {
	reg := (param & 0x0F00) >> 8
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}
	c.V[0xf] = c.V[reg] & 0x1
	c.V[reg] >>= 1
	c.PC += 2
	return nil
}

// skp skips the next instruction if the key with the value of Vx is pressed.
func skp(c *CPU, param uint16) error {
	reg := (param & 0x0F00) >> 8
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}
	keyIndex := c.V[reg]
	if keyIndex >= 16 {
		return fmt.Errorf("%w: 0x%X", ErrKeyIndexOutOfBounds, keyIndex)
	}
	if c.Key[keyIndex] {
		c.PC += 4
	} else {
		c.PC += 2
	}
	return nil
}

// sknp skips the next instruction if the key with the value of Vx is not pressed.
func sknp(c *CPU, param uint16) error {
	reg := (param & 0x0F00) >> 8
	if reg > 15 {
		return fmt.Errorf("%w: 0x%X", ErrRegisterOutOfBounds, reg)
	}
	keyIndex := c.V[reg]
	if keyIndex >= 16 {
		return fmt.Errorf("%w: 0x%X", ErrKeyIndexOutOfBounds, keyIndex)
	}
	if !c.Key[keyIndex] {
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
	if reg1 > 15 || reg2 > 15 {
		return fmt.Errorf("%w: 0x%X, 0x%X", ErrRegisterOutOfBounds, reg1, reg2)
	}

	if c.V[reg2] > c.V[reg1] {
		c.V[0xf] = 1
	} else {
		c.V[0xf] = 0
	}

	c.V[reg1] = c.V[reg2] - c.V[reg1]

	c.PC += 2
	return nil
}
