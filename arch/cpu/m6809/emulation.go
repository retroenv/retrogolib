package m6809

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
		return c.memory.Read(addr), nil
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
		return c.memory.ReadWord(addr), nil
	}
}

// -- 8-bit ALU operations --

// add8 performs an 8-bit addition and sets H, N, Z, V, C flags.
func (c *CPU) add8(a, b uint8, carry uint8) uint8 {
	sum := uint16(a) + uint16(b) + uint16(carry)
	result := uint8(sum)

	setFlag(&c.Flags.H, (a^b^result)&0x10 != 0)
	setFlag(&c.Flags.C, sum > 0xFF)
	setFlag(&c.Flags.V, (a^b)&0x80 == 0 && (a^result)&0x80 != 0)
	c.setZN8(result)

	return result
}

// sub8 performs an 8-bit subtraction and sets N, Z, V, C flags.
func (c *CPU) sub8(a, b uint8, carry uint8) uint8 {
	diff := int16(a) - int16(b) - int16(carry)
	result := uint8(diff)

	setFlag(&c.Flags.C, diff < 0)
	setFlag(&c.Flags.V, (a^b)&0x80 != 0 && (a^result)&0x80 != 0)
	c.setZN8(result)

	return result
}

// -- Instruction handlers --

func adca(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.A = c.add8(c.A, val, c.Flags.C)
	return nil
}

func adcb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.B = c.add8(c.B, val, c.Flags.C)
	return nil
}

func adda(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.A = c.add8(c.A, val, 0)
	return nil
}

func addb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.B = c.add8(c.B, val, 0)
	return nil
}

func addd(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	d := c.D()
	sum := uint32(d) + uint32(val)
	result := uint16(sum)

	setFlag(&c.Flags.C, sum > 0xFFFF)
	setFlag(&c.Flags.V, (d^val)&0x8000 == 0 && (d^result)&0x8000 != 0)
	c.setZN16(result)
	c.SetD(result)
	return nil
}

func anda(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.A &= val
	c.setZN8(c.A)
	c.Flags.V = 0
	return nil
}

func andb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.B &= val
	c.setZN8(c.B)
	c.Flags.V = 0
	return nil
}

func andcc(c *CPU, params ...any) error {
	mask := uint8(params[0].(Immediate8))
	c.SetCC(c.GetCC() & mask)
	return nil
}

func orcc(c *CPU, params ...any) error {
	mask := uint8(params[0].(Immediate8))
	c.SetCC(c.GetCC() | mask)
	return nil
}

func bita(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	result := c.A & val
	c.setZN8(result)
	c.Flags.V = 0
	return nil
}

func bitb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	result := c.B & val
	c.setZN8(result)
	c.Flags.V = 0
	return nil
}

func cmpa(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.compare8(c.A, val)
	return nil
}

func cmpb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.compare8(c.B, val)
	return nil
}

func cmpd(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.compare16(c.D(), val)
	return nil
}

func cmps(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.compare16(c.S, val)
	return nil
}

func cmpu(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.compare16(c.U, val)
	return nil
}

func cmpx(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.compare16(c.X, val)
	return nil
}

func cmpy(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.compare16(c.Y, val)
	return nil
}

func eora(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.A ^= val
	c.setZN8(c.A)
	c.Flags.V = 0
	return nil
}

func eorb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.B ^= val
	c.setZN8(c.B)
	c.Flags.V = 0
	return nil
}

func ora(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.A |= val
	c.setZN8(c.A)
	c.Flags.V = 0
	return nil
}

func orb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.B |= val
	c.setZN8(c.B)
	c.Flags.V = 0
	return nil
}

func suba(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.A = c.sub8(c.A, val, 0)
	return nil
}

func subb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.B = c.sub8(c.B, val, 0)
	return nil
}

func subd(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	d := c.D()
	diff := int32(d) - int32(val)
	result := uint16(diff)

	setFlag(&c.Flags.C, diff < 0)
	setFlag(&c.Flags.V, (d^val)&0x8000 != 0 && (d^result)&0x8000 != 0)
	c.setZN16(result)
	c.SetD(result)
	return nil
}

func sbca(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.A = c.sub8(c.A, val, c.Flags.C)
	return nil
}

func sbcb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.B = c.sub8(c.B, val, c.Flags.C)
	return nil
}

// -- Load and Store --

func lda(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.A = val
	c.setZN8(c.A)
	c.Flags.V = 0
	return nil
}

func ldb(c *CPU, params ...any) error {
	val, err := c.readOperand8(params[0])
	if err != nil {
		return err
	}
	c.B = val
	c.setZN8(c.B)
	c.Flags.V = 0
	return nil
}

func ldd(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.SetD(val)
	c.setZN16(val)
	c.Flags.V = 0
	return nil
}

func lds(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.S = val
	c.setZN16(c.S)
	c.Flags.V = 0
	return nil
}

func ldu(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.U = val
	c.setZN16(c.U)
	c.Flags.V = 0
	return nil
}

func ldx(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.X = val
	c.setZN16(c.X)
	c.Flags.V = 0
	return nil
}

func ldy(c *CPU, params ...any) error {
	val, err := c.readOperand16(params[0])
	if err != nil {
		return err
	}
	c.Y = val
	c.setZN16(c.Y)
	c.Flags.V = 0
	return nil
}

func sta(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	c.memory.Write(addr, c.A)
	c.setZN8(c.A)
	c.Flags.V = 0
	return nil
}

func stb(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	c.memory.Write(addr, c.B)
	c.setZN8(c.B)
	c.Flags.V = 0
	return nil
}

func std(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	d := c.D()
	c.memory.WriteWord(addr, d)
	c.setZN16(d)
	c.Flags.V = 0
	return nil
}

func sts(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	c.memory.WriteWord(addr, c.S)
	c.setZN16(c.S)
	c.Flags.V = 0
	return nil
}

func stu(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	c.memory.WriteWord(addr, c.U)
	c.setZN16(c.U)
	c.Flags.V = 0
	return nil
}

func stx(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	c.memory.WriteWord(addr, c.X)
	c.setZN16(c.X)
	c.Flags.V = 0
	return nil
}

func sty(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	c.memory.WriteWord(addr, c.Y)
	c.setZN16(c.Y)
	c.Flags.V = 0
	return nil
}

// -- LEA instructions --

func leax(c *CPU, params ...any) error {
	c.X = uint16(params[0].(IndexedAddr))
	setFlag(&c.Flags.Z, c.X == 0)
	return nil
}

func leay(c *CPU, params ...any) error {
	c.Y = uint16(params[0].(IndexedAddr))
	setFlag(&c.Flags.Z, c.Y == 0)
	return nil
}

func leas(c *CPU, params ...any) error {
	c.S = uint16(params[0].(IndexedAddr))
	return nil
}

func leau(c *CPU, params ...any) error {
	c.U = uint16(params[0].(IndexedAddr))
	return nil
}

// -- Shift and Rotate (memory) --

func neg8(c *CPU, val uint8) uint8 {
	result := uint8(-int8(val))
	setFlag(&c.Flags.C, val != 0)
	setFlag(&c.Flags.V, val == 0x80)
	c.setZN8(result)
	return result
}

func negMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	c.memory.Write(addr, neg8(c, val))
	return nil
}

func nega(c *CPU) error {
	c.A = neg8(c, c.A)
	return nil
}

func negb(c *CPU) error {
	c.B = neg8(c, c.B)
	return nil
}

func com8(c *CPU, val uint8) uint8 {
	result := ^val
	c.setZN8(result)
	c.Flags.V = 0
	c.Flags.C = 1
	return result
}

func comMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	c.memory.Write(addr, com8(c, val))
	return nil
}

func coma(c *CPU) error {
	c.A = com8(c, c.A)
	return nil
}

func comb(c *CPU) error {
	c.B = com8(c, c.B)
	return nil
}

func lsr8(c *CPU, val uint8) uint8 {
	setFlag(&c.Flags.C, val&0x01 != 0)
	result := val >> 1
	c.setZN8(result)
	return result
}

func lsrMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	c.memory.Write(addr, lsr8(c, val))
	return nil
}

func lsra(c *CPU) error {
	c.A = lsr8(c, c.A)
	return nil
}

func lsrb(c *CPU) error {
	c.B = lsr8(c, c.B)
	return nil
}

func ror8(c *CPU, val uint8) uint8 {
	oldC := c.Flags.C
	setFlag(&c.Flags.C, val&0x01 != 0)
	result := val >> 1
	if oldC != 0 {
		result |= 0x80
	}
	c.setZN8(result)
	return result
}

func rorMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	c.memory.Write(addr, ror8(c, val))
	return nil
}

func rora(c *CPU) error {
	c.A = ror8(c, c.A)
	return nil
}

func rorb(c *CPU) error {
	c.B = ror8(c, c.B)
	return nil
}

func asr8(c *CPU, val uint8) uint8 {
	setFlag(&c.Flags.C, val&0x01 != 0)
	result := (val >> 1) | (val & 0x80) // preserve sign bit
	c.setZN8(result)
	return result
}

func asrMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	c.memory.Write(addr, asr8(c, val))
	return nil
}

func asra(c *CPU) error {
	c.A = asr8(c, c.A)
	return nil
}

func asrb(c *CPU) error {
	c.B = asr8(c, c.B)
	return nil
}

func asl8(c *CPU, val uint8) uint8 {
	setFlag(&c.Flags.C, val&0x80 != 0)
	result := val << 1
	setFlag(&c.Flags.V, (val^result)&0x80 != 0)
	c.setZN8(result)
	return result
}

func aslMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	c.memory.Write(addr, asl8(c, val))
	return nil
}

func asla(c *CPU) error {
	c.A = asl8(c, c.A)
	return nil
}

func aslb(c *CPU) error {
	c.B = asl8(c, c.B)
	return nil
}

func rol8(c *CPU, val uint8) uint8 {
	oldC := c.Flags.C
	setFlag(&c.Flags.C, val&0x80 != 0)
	result := val << 1
	if oldC != 0 {
		result |= 0x01
	}
	setFlag(&c.Flags.V, (val^result)&0x80 != 0)
	c.setZN8(result)
	return result
}

func rolMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	c.memory.Write(addr, rol8(c, val))
	return nil
}

func rola(c *CPU) error {
	c.A = rol8(c, c.A)
	return nil
}

func rolb(c *CPU) error {
	c.B = rol8(c, c.B)
	return nil
}

func dec8(c *CPU, val uint8) uint8 {
	result := val - 1
	setFlag(&c.Flags.V, val == 0x80)
	c.setZN8(result)
	return result
}

func decMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	c.memory.Write(addr, dec8(c, val))
	return nil
}

func deca(c *CPU) error {
	c.A = dec8(c, c.A)
	return nil
}

func decb(c *CPU) error {
	c.B = dec8(c, c.B)
	return nil
}

func inc8(c *CPU, val uint8) uint8 {
	result := val + 1
	setFlag(&c.Flags.V, val == 0x7F)
	c.setZN8(result)
	return result
}

func incMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	c.memory.Write(addr, inc8(c, val))
	return nil
}

func inca(c *CPU) error {
	c.A = inc8(c, c.A)
	return nil
}

func incb(c *CPU) error {
	c.B = inc8(c, c.B)
	return nil
}

func tst8(c *CPU, val uint8) {
	c.setZN8(val)
	c.Flags.V = 0
}

func tstMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	val := c.memory.Read(addr)
	tst8(c, val)
	return nil
}

func tsta(c *CPU) error {
	tst8(c, c.A)
	return nil
}

func tstb(c *CPU) error {
	tst8(c, c.B)
	return nil
}

func clr8(c *CPU) {
	c.Flags.N = 0
	c.Flags.Z = 1
	c.Flags.V = 0
	c.Flags.C = 0
}

func clrMem(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	c.memory.Write(addr, 0)
	clr8(c)
	return nil
}

func clra(c *CPU) error {
	c.A = 0
	clr8(c)
	return nil
}

func clrb(c *CPU) error {
	c.B = 0
	clr8(c)
	return nil
}

// -- Miscellaneous --

// nop - No Operation.
func nop(_ *CPU) error { return nil }

// abx - Add B to X (unsigned).
func abx(c *CPU) error {
	c.X += uint16(c.B)
	return nil
}

// sexFn - Sign Extend B into A.
func sexFn(c *CPU) error {
	if c.B&0x80 != 0 {
		c.A = 0xFF
	} else {
		c.A = 0x00
	}
	c.setZN16(c.D())
	return nil
}

// mulFn - Unsigned multiply A * B -> D.
func mulFn(c *CPU) error {
	result := uint16(c.A) * uint16(c.B)
	c.SetD(result)
	setFlag(&c.Flags.Z, result == 0)
	setFlag(&c.Flags.C, c.B&0x80 != 0) // bit 7 of result low byte
	return nil
}

// daaFn - Decimal Adjust A.
func daaFn(c *CPU) error {
	a := uint16(c.A)
	cf := uint16(0)

	if c.Flags.H != 0 || (a&0x0F) > 9 {
		cf |= 0x06
	}
	if c.Flags.C != 0 || a > 0x99 {
		cf |= 0x60
	}

	a += cf
	c.A = uint8(a)
	c.setZN8(c.A)
	if a > 0xFF {
		c.Flags.C = 1
	}
	return nil
}
