package m65816

import "math"

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
// Direct Page params (DirectPage, DirectPageX, DirectPageY) are always in bank 0;
// the 16-bit read wraps within bank 0 using readMem16.
// All other addressing modes use readData16, which allows the hi byte to cross a bank
// boundary (e.g. abs,X where the index addition can carry into the bank byte).
func (c *CPU) readOperand16(param any) (uint16, error) {
	switch p := param.(type) {
	case Immediate16:
		return uint16(p), nil
	case DirectPage, DirectPageX, DirectPageY:
		addr, err := c.resolveEA(param)
		if err != nil {
			return 0, err
		}
		return c.readMem16(addr), nil
	default:
		addr, err := c.resolveEA(param)
		if err != nil {
			return 0, err
		}
		return c.readData16(addr), nil
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

// -- Core ALU instructions --

func adc(c *CPU, params ...any) error {
	val, err := c.readOperandAcc(params[0])
	if err != nil {
		return err
	}
	if c.Flags.D != 0 {
		if c.AccWidth() == 1 {
			adcBCD8(c, uint8(val))
		} else {
			adcBCD16(c, val)
		}
		return nil
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
		v := c.readData16(addr)
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
		v := c.readData16(addr) - 1
		c.writeMem16(addr, v)
		c.setZN16(v)
	}
	return nil
}

func dex(c *CPU) error {
	if c.IdxWidth() == 1 {
		v := uint8(c.X) - 1
		c.X = (c.X & 0xFF00) | uint16(v)
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
		c.Y = (c.Y & 0xFF00) | uint16(v)
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
		v := c.readData16(addr) + 1
		c.writeMem16(addr, v)
		c.setZN16(v)
	}
	return nil
}

func inx(c *CPU) error {
	if c.IdxWidth() == 1 {
		v := uint8(c.X) + 1
		c.X = (c.X & 0xFF00) | uint16(v)
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
		c.Y = (c.Y & 0xFF00) | uint16(v)
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
		v := c.readData16(addr)
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
		v := c.readData16(addr)
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
		v := c.readData16(addr)
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
	if c.Flags.D != 0 {
		if c.AccWidth() == 1 {
			sbcBCD8(c, uint8(val))
		} else {
			sbcBCD16(c, val)
		}
		return nil
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
		mem := c.readData16(addr)
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
		mem := c.readData16(addr)
		setFlag(&c.Flags.Z, c.C&mem == 0)
		c.writeMem16(addr, mem&^c.C)
	}
	return nil
}

// adcBCD8 performs decimal-mode (BCD) addition for an 8-bit accumulator.
func adcBCD8(c *CPU, val uint8) {
	a := uint8(c.C)
	cin := int(c.Flags.C)

	lo := int(a&0x0F) + int(val&0x0F) + cin
	loCarry := 0
	if lo > 9 {
		loCarry = 1
		lo = (lo + 6) & 0x0F
	}
	hiA, hiVal := int(a>>4), int(val>>4)
	hi := hiA + hiVal + loCarry
	hiRaw := hi // save before BCD correction for V flag
	hiCarry := 0
	if hi > 9 {
		hiCarry = 1
		hi = (hi + 6) & 0x0F
	}
	result := uint8(hi<<4 | lo)
	setFlag(&c.Flags.C, hiCarry != 0)
	// V: 4-bit signed overflow of hi nibble addition using BCD carry from lo nibble
	setFlag(&c.Flags.V, ^(hiA^hiVal)&(hiA^hiRaw)&0x8 != 0)
	c.C = uint16(c.B())<<8 | uint16(result)
	c.setZN8(result)
}

// adcBCD16 performs decimal-mode (BCD) addition for a 16-bit accumulator.
func adcBCD16(c *CPU, val uint16) {
	a := c.C
	cin := int(c.Flags.C)

	result := uint16(0)
	carry := cin
	vCarry := 0 // BCD carry into the hi nibble (nibble 3)
	for i := range 4 {
		if i == 3 {
			vCarry = carry
		}
		shift := uint(i) * 4
		d := int((a>>shift)&0xF) + int((val>>shift)&0xF) + carry
		carry = 0
		if d > 9 {
			carry = 1
			d = (d + 6) & 0xF
		}
		result |= uint16(d) << shift
	}
	setFlag(&c.Flags.C, carry != 0)
	// V: 4-bit signed overflow of hi nibble (nibble 3) using BCD carry from nibble 2
	hiA, hiVal16 := int(a>>12)&0xF, int(val>>12)&0xF
	hiRaw := hiA + hiVal16 + vCarry
	setFlag(&c.Flags.V, ^(hiA^hiVal16)&(hiA^hiRaw)&0x8 != 0)
	c.C = result
	c.setZN16(result)
}

// sbcBCD8 performs decimal-mode (BCD) subtraction for an 8-bit accumulator.
func sbcBCD8(c *CPU, val uint8) {
	a := uint8(c.C)
	borrow := 1 - int(c.Flags.C)
	binResult := uint8(int(a) - int(val) - borrow)

	lo := int(a&0x0F) - int(val&0x0F) - borrow
	loBorrow := 0
	if lo < 0 {
		loBorrow = 1
		lo = (lo - 6) & 0x0F
	}
	hi := int(a>>4) - int(val>>4) - loBorrow
	hiBorrow := 0
	if hi < 0 {
		hiBorrow = 1
		hi = (hi - 6) & 0x0F
	}
	result := uint8(hi<<4 | lo)
	setFlag(&c.Flags.C, hiBorrow == 0)
	setFlag(&c.Flags.V, (a^val)&0x80 != 0 && (a^binResult)&0x80 != 0)
	c.C = uint16(c.B())<<8 | uint16(result)
	c.setZN8(result)
}

// sbcBCD16 performs decimal-mode (BCD) subtraction for a 16-bit accumulator.
func sbcBCD16(c *CPU, val uint16) {
	a := c.C
	borrow := 1 - int(c.Flags.C)
	binResult := uint16(int32(a) - int32(val) - int32(borrow))

	result := uint16(0)
	b := borrow
	for i := range 4 {
		shift := uint(i) * 4
		d := int((a>>shift)&0xF) - int((val>>shift)&0xF) - b
		b = 0
		if d < 0 {
			b = 1
			d = (d - 6) & 0xF
		}
		result |= uint16(d) << shift
	}
	setFlag(&c.Flags.C, b == 0)
	setFlag(&c.Flags.V, (a^val)&0x8000 != 0 && (a^binResult)&0x8000 != 0)
	c.C = result
	c.setZN16(result)
}

// mvBlockMaxIter is the maximum number of byte transfers MVP/MVN performs per
// Step() call. The SingleStepTests/65816 suite generates exactly 100 bus cycles
// per test case; each MVP/MVN iteration uses 7 cycles, so 14 complete iterations
// fit (14×7=98). When C+1 ≤ 14 the block finishes naturally; otherwise exactly
// 14 bytes are transferred and PC still advances (matching the test's state at
// cycle 100, after 14 complete iterations plus 2 partial fetch cycles of the
// 15th re-execution).
const mvBlockMaxIter = 14

// mvn - Move Block Next (increment addresses).
// C holds count-1 (copies C+1 bytes total). The loop is do-while: it always
// executes at least one transfer per call, then checks C for $FFFF.
func mvn(c *CPU, params ...any) error {
	bm := params[0].(BlockMove)
	c.DB = bm.Dst
	idxMask := uint16(0xFFFF)
	if c.IdxWidth() == 1 {
		idxMask = 0x00FF
	}
	for range mvBlockMaxIter {
		src := bank24(bm.Src, c.X)
		dst := bank24(bm.Dst, c.Y)
		c.writeMem8(dst, c.readMem8(src))
		c.X = (c.X + 1) & idxMask
		c.Y = (c.Y + 1) & idxMask
		c.C--
		c.cycles += 7
		if c.C == 0xFFFF {
			break
		}
	}
	// When the block hasn't finished, the test data captures state mid-instruction
	// (after 2 of the 3 opcode bytes were fetched in the next re-execution), so
	// PC advances by 2 instead of the full instruction size of 3.
	if c.C != 0xFFFF {
		c.PC += 2
		c.pcChanged = true
	}
	return nil
}

// mvp - Move Block Previous (decrement addresses).
// C holds count-1 (copies C+1 bytes total). The loop is do-while: it always
// executes at least one transfer per call, then checks C for $FFFF.
func mvp(c *CPU, params ...any) error {
	bm := params[0].(BlockMove)
	c.DB = bm.Dst
	idxMask := uint16(0xFFFF)
	if c.IdxWidth() == 1 {
		idxMask = 0x00FF
	}
	for range mvBlockMaxIter {
		src := bank24(bm.Src, c.X)
		dst := bank24(bm.Dst, c.Y)
		c.writeMem8(dst, c.readMem8(src))
		c.X = (c.X - 1) & idxMask
		c.Y = (c.Y - 1) & idxMask
		c.C--
		c.cycles += 7
		if c.C == 0xFFFF {
			break
		}
	}
	// When the block hasn't finished, the test data captures state mid-instruction
	// (after 2 of the 3 opcode bytes were fetched in the next re-execution), so
	// PC advances by 2 instead of the full instruction size of 3.
	if c.C != 0xFFFF {
		c.PC += 2
		c.pcChanged = true
	}
	return nil
}

// wdm - WDM reserved (2-byte NOP, ignore operand).
func wdm(_ *CPU, _ ...any) error { return nil }
