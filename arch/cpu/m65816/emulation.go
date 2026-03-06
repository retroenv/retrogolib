package m65816

import (
	"fmt"
	"math"
)

// readOperand8 reads an 8-bit value from a param (immediate or memory).
func (c *CPU) readOperand8(param any) (uint8, error) {
	switch p := param.(type) {
	case Immediate8:
		return uint8(p), nil
	default:
		addr, err := c.resolveEA(param)
		if err != nil {
			return 0, err
		}
		return c.readMem8(addr), nil
	}
}

// readOperand16 reads a 16-bit value from a param (immediate or memory).
func (c *CPU) readOperand16(param any) (uint16, error) {
	switch p := param.(type) {
	case Immediate16:
		return uint16(p), nil
	default:
		addr, err := c.resolveEA(param)
		if err != nil {
			return 0, err
		}
		return c.readMem16(addr), nil
	}
}

// readOperandAcc reads a value using the current accumulator width (8 or 16 bit).
func (c *CPU) readOperandAcc(param any) (uint16, error) {
	if c.AccWidth() == 1 {
		v, err := c.readOperand8(param)
		return uint16(v), err
	}
	return c.readOperand16(param)
}

// readOperandIdx reads a value using the current index width (8 or 16 bit).
func (c *CPU) readOperandIdx(param any) (uint16, error) {
	if c.IdxWidth() == 1 {
		v, err := c.readOperand8(param)
		return uint16(v), err
	}
	return c.readOperand16(param)
}

// writeOperand8 writes an 8-bit value to a memory address resolved from param.
func (c *CPU) writeOperand8(param any, value uint8) error {
	addr, err := c.resolveEA(param)
	if err != nil {
		return err
	}
	c.writeMem8(addr, value)
	return nil
}

// writeOperand16 writes a 16-bit value to a memory address resolved from param.
func (c *CPU) writeOperand16(param any, value uint16) error {
	addr, err := c.resolveEA(param)
	if err != nil {
		return err
	}
	c.writeMem16(addr, value)
	return nil
}

// writeOperandAcc writes accumulator-width value to memory.
func (c *CPU) writeOperandAcc(param any, value uint16) error {
	if c.AccWidth() == 1 {
		return c.writeOperand8(param, uint8(value))
	}
	return c.writeOperand16(param, value)
}

// -- Core ALU instructions --

func adc(c *CPU, params ...any) error {
	val, err := c.readOperandAcc(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		a := uint8(c.C)
		sum := int(a) + int(uint8(val)) + int(c.Flags.C)
		result := uint8(sum)
		setFlag(&c.Flags.C, sum > math.MaxUint8)
		setFlag(&c.Flags.V, (a^uint8(val))&0x80 == 0 && (a^result)&0x80 != 0)
		c.C = uint16(c.B())<<8 | uint16(result)
		c.setZN8(result)
	} else {
		a := c.C
		sum := int32(a) + int32(val) + int32(c.Flags.C)
		result := uint16(sum)
		setFlag(&c.Flags.C, sum > math.MaxUint16)
		setFlag(&c.Flags.V, (a^val)&0x8000 == 0 && (a^result)&0x8000 != 0)
		c.C = result
		c.setZN16(result)
	}
	return nil
}

func and(c *CPU, params ...any) error {
	val, err := c.readOperandAcc(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		result := uint8(c.C) & uint8(val)
		c.C = uint16(c.B())<<8 | uint16(result)
		c.setZN8(result)
	} else {
		c.C &= val
		c.setZN16(c.C)
	}
	return nil
}

func asl(c *CPU, params ...any) error {
	_, isAcc := params[0].(Accumulator)
	if isAcc {
		if c.AccWidth() == 1 {
			a := uint8(c.C)
			setFlag(&c.Flags.C, a&0x80 != 0)
			a <<= 1
			c.C = uint16(c.B())<<8 | uint16(a)
			c.setZN8(a)
		} else {
			setFlag(&c.Flags.C, c.C&0x8000 != 0)
			c.C <<= 1
			c.setZN16(c.C)
		}
		return nil
	}
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		v := c.readMem8(addr)
		setFlag(&c.Flags.C, v&0x80 != 0)
		v <<= 1
		c.writeMem8(addr, v)
		c.setZN8(v)
	} else {
		v := c.readMem16(addr)
		setFlag(&c.Flags.C, v&0x8000 != 0)
		v <<= 1
		c.writeMem16(addr, v)
		c.setZN16(v)
	}
	return nil
}

func bit(c *CPU, params ...any) error {
	val, err := c.readOperandAcc(params[0])
	if err != nil {
		return err
	}
	// BIT immediate only sets Z; non-immediate also sets N and V from memory
	_, isImm8 := params[0].(Immediate8)
	_, isImm16 := params[0].(Immediate16)
	isImm := isImm8 || isImm16

	if c.AccWidth() == 1 {
		b := uint8(val)
		setFlag(&c.Flags.Z, uint8(c.C)&b == 0)
		if !isImm {
			setFlag(&c.Flags.N, b&0x80 != 0)
			setFlag(&c.Flags.V, b&0x40 != 0)
		}
	} else {
		setFlag(&c.Flags.Z, c.C&val == 0)
		if !isImm {
			setFlag(&c.Flags.N, val&0x8000 != 0)
			setFlag(&c.Flags.V, val&0x4000 != 0)
		}
	}
	return nil
}

func cmp(c *CPU, params ...any) error {
	val, err := c.readOperandAcc(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		c.compare8(uint8(c.C), uint8(val))
	} else {
		c.compare16(c.C, val)
	}
	return nil
}

func cpx(c *CPU, params ...any) error {
	val, err := c.readOperandIdx(params[0])
	if err != nil {
		return err
	}
	if c.IdxWidth() == 1 {
		c.compare8(uint8(c.X), uint8(val))
	} else {
		c.compare16(c.X, val)
	}
	return nil
}

func cpy(c *CPU, params ...any) error {
	val, err := c.readOperandIdx(params[0])
	if err != nil {
		return err
	}
	if c.IdxWidth() == 1 {
		c.compare8(uint8(c.Y), uint8(val))
	} else {
		c.compare16(c.Y, val)
	}
	return nil
}

func dec(c *CPU, params ...any) error {
	_, isAcc := params[0].(Accumulator)
	if isAcc {
		if c.AccWidth() == 1 {
			a := uint8(c.C) - 1
			c.C = uint16(c.B())<<8 | uint16(a)
			c.setZN8(a)
		} else {
			c.C--
			c.setZN16(c.C)
		}
		return nil
	}
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		v := c.readMem8(addr) - 1
		c.writeMem8(addr, v)
		c.setZN8(v)
	} else {
		v := c.readMem16(addr) - 1
		c.writeMem16(addr, v)
		c.setZN16(v)
	}
	return nil
}

func dex(c *CPU) error {
	if c.IdxWidth() == 1 {
		v := uint8(c.X) - 1
		c.X = uint16(c.X&0xFF00) | uint16(v)
		c.setZN8(v)
	} else {
		c.X--
		c.setZN16(c.X)
	}
	return nil
}

func dey(c *CPU) error {
	if c.IdxWidth() == 1 {
		v := uint8(c.Y) - 1
		c.Y = uint16(c.Y&0xFF00) | uint16(v)
		c.setZN8(v)
	} else {
		c.Y--
		c.setZN16(c.Y)
	}
	return nil
}

func eor(c *CPU, params ...any) error {
	val, err := c.readOperandAcc(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		result := uint8(c.C) ^ uint8(val)
		c.C = uint16(c.B())<<8 | uint16(result)
		c.setZN8(result)
	} else {
		c.C ^= val
		c.setZN16(c.C)
	}
	return nil
}

func inc(c *CPU, params ...any) error {
	_, isAcc := params[0].(Accumulator)
	if isAcc {
		if c.AccWidth() == 1 {
			a := uint8(c.C) + 1
			c.C = uint16(c.B())<<8 | uint16(a)
			c.setZN8(a)
		} else {
			c.C++
			c.setZN16(c.C)
		}
		return nil
	}
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		v := c.readMem8(addr) + 1
		c.writeMem8(addr, v)
		c.setZN8(v)
	} else {
		v := c.readMem16(addr) + 1
		c.writeMem16(addr, v)
		c.setZN16(v)
	}
	return nil
}

func inx(c *CPU) error {
	if c.IdxWidth() == 1 {
		v := uint8(c.X) + 1
		c.X = uint16(c.X&0xFF00) | uint16(v)
		c.setZN8(v)
	} else {
		c.X++
		c.setZN16(c.X)
	}
	return nil
}

func iny(c *CPU) error {
	if c.IdxWidth() == 1 {
		v := uint8(c.Y) + 1
		c.Y = uint16(c.Y&0xFF00) | uint16(v)
		c.setZN8(v)
	} else {
		c.Y++
		c.setZN16(c.Y)
	}
	return nil
}

func lsr(c *CPU, params ...any) error {
	_, isAcc := params[0].(Accumulator)
	if isAcc {
		if c.AccWidth() == 1 {
			a := uint8(c.C)
			setFlag(&c.Flags.C, a&0x01 != 0)
			a >>= 1
			c.C = uint16(c.B())<<8 | uint16(a)
			c.setZN8(a)
		} else {
			setFlag(&c.Flags.C, c.C&0x0001 != 0)
			c.C >>= 1
			c.setZN16(c.C)
		}
		return nil
	}
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		v := c.readMem8(addr)
		setFlag(&c.Flags.C, v&0x01 != 0)
		v >>= 1
		c.writeMem8(addr, v)
		c.setZN8(v)
	} else {
		v := c.readMem16(addr)
		setFlag(&c.Flags.C, v&0x0001 != 0)
		v >>= 1
		c.writeMem16(addr, v)
		c.setZN16(v)
	}
	return nil
}

func nop(_ *CPU) error { return nil }

func ora(c *CPU, params ...any) error {
	val, err := c.readOperandAcc(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		result := uint8(c.C) | uint8(val)
		c.C = uint16(c.B())<<8 | uint16(result)
		c.setZN8(result)
	} else {
		c.C |= val
		c.setZN16(c.C)
	}
	return nil
}

func rol(c *CPU, params ...any) error {
	carry := c.Flags.C
	_, isAcc := params[0].(Accumulator)
	if isAcc {
		if c.AccWidth() == 1 {
			a := uint8(c.C)
			setFlag(&c.Flags.C, a&0x80 != 0)
			a = (a << 1) | carry
			c.C = uint16(c.B())<<8 | uint16(a)
			c.setZN8(a)
		} else {
			setFlag(&c.Flags.C, c.C&0x8000 != 0)
			c.C = (c.C << 1) | uint16(carry)
			c.setZN16(c.C)
		}
		return nil
	}
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		v := c.readMem8(addr)
		setFlag(&c.Flags.C, v&0x80 != 0)
		v = (v << 1) | carry
		c.writeMem8(addr, v)
		c.setZN8(v)
	} else {
		v := c.readMem16(addr)
		setFlag(&c.Flags.C, v&0x8000 != 0)
		v = (v << 1) | uint16(carry)
		c.writeMem16(addr, v)
		c.setZN16(v)
	}
	return nil
}

func ror(c *CPU, params ...any) error {
	carry := c.Flags.C
	_, isAcc := params[0].(Accumulator)
	if isAcc {
		if c.AccWidth() == 1 {
			a := uint8(c.C)
			setFlag(&c.Flags.C, a&0x01 != 0)
			a = (a >> 1) | (carry << 7)
			c.C = uint16(c.B())<<8 | uint16(a)
			c.setZN8(a)
		} else {
			setFlag(&c.Flags.C, c.C&0x0001 != 0)
			c.C = (c.C >> 1) | (uint16(carry) << 15)
			c.setZN16(c.C)
		}
		return nil
	}
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		v := c.readMem8(addr)
		setFlag(&c.Flags.C, v&0x01 != 0)
		v = (v >> 1) | (carry << 7)
		c.writeMem8(addr, v)
		c.setZN8(v)
	} else {
		v := c.readMem16(addr)
		setFlag(&c.Flags.C, v&0x0001 != 0)
		v = (v >> 1) | (uint16(carry) << 15)
		c.writeMem16(addr, v)
		c.setZN16(v)
	}
	return nil
}

func sbc(c *CPU, params ...any) error {
	val, err := c.readOperandAcc(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		a := uint8(c.C)
		v := uint8(val)
		diff := int(a) - int(v) - (1 - int(c.Flags.C))
		result := uint8(diff)
		setFlag(&c.Flags.C, diff >= 0)
		setFlag(&c.Flags.V, (a^v)&0x80 != 0 && (a^result)&0x80 != 0)
		c.C = uint16(c.B())<<8 | uint16(result)
		c.setZN8(result)
	} else {
		a := c.C
		diff := int32(a) - int32(val) - int32(1-c.Flags.C)
		result := uint16(diff)
		setFlag(&c.Flags.C, diff >= 0)
		setFlag(&c.Flags.V, (a^val)&0x8000 != 0 && (a^result)&0x8000 != 0)
		c.C = result
		c.setZN16(result)
	}
	return nil
}

// TSB/TRB
func tsb(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		mem := c.readMem8(addr)
		setFlag(&c.Flags.Z, uint8(c.C)&mem == 0)
		c.writeMem8(addr, mem|uint8(c.C))
	} else {
		mem := c.readMem16(addr)
		setFlag(&c.Flags.Z, c.C&mem == 0)
		c.writeMem16(addr, mem|c.C)
	}
	return nil
}

func trb(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		mem := c.readMem8(addr)
		setFlag(&c.Flags.Z, uint8(c.C)&mem == 0)
		c.writeMem8(addr, mem&^uint8(c.C))
	} else {
		mem := c.readMem16(addr)
		setFlag(&c.Flags.Z, c.C&mem == 0)
		c.writeMem16(addr, mem&^c.C)
	}
	return nil
}

// mvn - Move Block Next (increment addresses).
func mvn(c *CPU, params ...any) error {
	bm := params[0].(BlockMove)
	c.DB = bm.Dst
	for c.C != 0xFFFF {
		src := bank24(bm.Src, c.X)
		dst := bank24(bm.Dst, c.Y)
		c.writeMem8(dst, c.readMem8(src))
		c.X++
		c.Y++
		c.C--
		c.cycles += 7
	}
	// Final byte
	src := bank24(bm.Src, c.X)
	dst := bank24(bm.Dst, c.Y)
	c.writeMem8(dst, c.readMem8(src))
	c.X++
	c.Y++
	c.C--
	return nil
}

// mvp - Move Block Previous (decrement addresses).
func mvp(c *CPU, params ...any) error {
	bm := params[0].(BlockMove)
	c.DB = bm.Dst
	for c.C != 0xFFFF {
		src := bank24(bm.Src, c.X)
		dst := bank24(bm.Dst, c.Y)
		c.writeMem8(dst, c.readMem8(src))
		c.X--
		c.Y--
		c.C--
		c.cycles += 7
	}
	src := bank24(bm.Src, c.X)
	dst := bank24(bm.Dst, c.Y)
	c.writeMem8(dst, c.readMem8(src))
	c.X--
	c.Y--
	c.C--
	return nil
}

// wdm - WDM reserved (2-byte NOP, ignore operand).
func wdm(_ *CPU, _ ...any) error { return nil }

// resolveEA helper for accumulator (should not be called).
func (c *CPU) errImmResolve(param any) error {
	return fmt.Errorf("%w: cannot resolve immediate as address (type %T)", ErrUnsupportedAddressingMode, param)
}
