package z80

// Undocumented DD prefix instructions - IXH/IXL half-register operations

// IXH returns the high byte of IX.
func (c *CPU) IXH() uint8 { return uint8(c.IX >> 8) }

// IXL returns the low byte of IX.
func (c *CPU) IXL() uint8 { return uint8(c.IX) }

// SetIXH sets the high byte of IX.
func (c *CPU) SetIXH(v uint8) { c.IX = uint16(v)<<8 | uint16(c.IXL()) }

// SetIXL sets the low byte of IX.
func (c *CPU) SetIXL(v uint8) { c.IX = uint16(c.IXH())<<8 | uint16(v) }

// INC/DEC IXH/IXL

func ddIncIXH(c *CPU) error { c.SetIXH(c.inc8(c.IXH())); return nil }
func ddDecIXH(c *CPU) error { c.SetIXH(c.dec8(c.IXH())); return nil }
func ddIncIXL(c *CPU) error { c.SetIXL(c.inc8(c.IXL())); return nil }
func ddDecIXL(c *CPU) error { c.SetIXL(c.dec8(c.IXL())); return nil }

// LD IXH/IXL,n

func ddLdIXHn(c *CPU, _ ...any) error {
	c.SetIXH(c.memory.Read(c.PC + 2))
	return nil
}

func ddLdIXLn(c *CPU, _ ...any) error {
	c.SetIXL(c.memory.Read(c.PC + 2))
	return nil
}

// LD r,IXH / LD r,IXL

func ddLdBIXH(c *CPU) error { c.B = c.IXH(); return nil }
func ddLdBIXL(c *CPU) error { c.B = c.IXL(); return nil }
func ddLdCIXH(c *CPU) error { c.C = c.IXH(); return nil }
func ddLdCIXL(c *CPU) error { c.C = c.IXL(); return nil }
func ddLdDIXH(c *CPU) error { c.D = c.IXH(); return nil }
func ddLdDIXL(c *CPU) error { c.D = c.IXL(); return nil }
func ddLdEIXH(c *CPU) error { c.E = c.IXH(); return nil }
func ddLdEIXL(c *CPU) error { c.E = c.IXL(); return nil }
func ddLdAIXH(c *CPU) error { c.A = c.IXH(); return nil }
func ddLdAIXL(c *CPU) error { c.A = c.IXL(); return nil }

// LD IXH,r / LD IXL,r

func ddLdIXHB(c *CPU) error   { c.SetIXH(c.B); return nil }
func ddLdIXHC(c *CPU) error   { c.SetIXH(c.C); return nil }
func ddLdIXHD(c *CPU) error   { c.SetIXH(c.D); return nil }
func ddLdIXHE(c *CPU) error   { c.SetIXH(c.E); return nil }
func ddLdIXHIXH(c *CPU) error { return nil } // LD IXH,IXH = NOP
func ddLdIXHIXL(c *CPU) error { c.SetIXH(c.IXL()); return nil }
func ddLdIXHA(c *CPU) error   { c.SetIXH(c.A); return nil }

func ddLdIXLB(c *CPU) error   { c.SetIXL(c.B); return nil }
func ddLdIXLC(c *CPU) error   { c.SetIXL(c.C); return nil }
func ddLdIXLD(c *CPU) error   { c.SetIXL(c.D); return nil }
func ddLdIXLE(c *CPU) error   { c.SetIXL(c.E); return nil }
func ddLdIXLIXH(c *CPU) error { c.SetIXL(c.IXH()); return nil }
func ddLdIXLIXL(c *CPU) error { return nil } // LD IXL,IXL = NOP
func ddLdIXLA(c *CPU) error   { c.SetIXL(c.A); return nil }

// Arithmetic operations with IXH/IXL

func ddAddAIXH(c *CPU) error { c.A = c.add8(c.A, c.IXH()); return nil }
func ddAddAIXL(c *CPU) error { c.A = c.add8(c.A, c.IXL()); return nil }
func ddAdcAIXH(c *CPU) error { c.A = c.adc(c.A, c.IXH()); return nil }
func ddAdcAIXL(c *CPU) error { c.A = c.adc(c.A, c.IXL()); return nil }
func ddSubIXH(c *CPU) error  { c.A = c.sub8(c.A, c.IXH()); return nil }
func ddSubIXL(c *CPU) error  { c.A = c.sub8(c.A, c.IXL()); return nil }
func ddSbcAIXH(c *CPU) error { c.A = c.sbc(c.A, c.IXH()); return nil }
func ddSbcAIXL(c *CPU) error { c.A = c.sbc(c.A, c.IXL()); return nil }
func ddAndIXH(c *CPU) error  { c.A = c.and8(c.A, c.IXH()); return nil }
func ddAndIXL(c *CPU) error  { c.A = c.and8(c.A, c.IXL()); return nil }
func ddXorIXH(c *CPU) error  { c.A = c.xor8(c.A, c.IXH()); return nil }
func ddXorIXL(c *CPU) error  { c.A = c.xor8(c.A, c.IXL()); return nil }
func ddOrIXH(c *CPU) error   { c.A = c.or8(c.A, c.IXH()); return nil }
func ddOrIXL(c *CPU) error   { c.A = c.or8(c.A, c.IXL()); return nil }
func ddCpIXH(c *CPU) error   { c.cp(c.A, c.IXH()); return nil }
func ddCpIXL(c *CPU) error   { c.cp(c.A, c.IXL()); return nil }
