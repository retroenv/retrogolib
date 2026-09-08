package z80

// Undocumented FD prefix instructions - IYH/IYL half-register operations

// IYH returns the high byte of IY.
func (cpu *CPU) IYH() uint8 { return uint8(cpu.IY >> 8) }

// IYL returns the low byte of IY.
func (cpu *CPU) IYL() uint8 { return uint8(cpu.IY) }

// SetIYH sets the high byte of IY.
func (cpu *CPU) SetIYH(v uint8) { cpu.IY = uint16(v)<<8 | uint16(cpu.IYL()) }

// SetIYL sets the low byte of IY.
func (cpu *CPU) SetIYL(v uint8) { cpu.IY = uint16(cpu.IYH())<<8 | uint16(v) }

// INC/DEC IYH/IYL

func fdIncIYH(cpu *CPU) error { cpu.SetIYH(cpu.inc8(cpu.IYH())); return nil }
func fdDecIYH(cpu *CPU) error { cpu.SetIYH(cpu.dec8(cpu.IYH())); return nil }
func fdIncIYL(cpu *CPU) error { cpu.SetIYL(cpu.inc8(cpu.IYL())); return nil }
func fdDecIYL(cpu *CPU) error { cpu.SetIYL(cpu.dec8(cpu.IYL())); return nil }

// LD IYH/IYL,n

func fdLdIYHn(cpu *CPU, _ ...any) error {
	cpu.SetIYH(cpu.bus.Read(cpu.PC + 2))
	return nil
}

func fdLdIYLn(c *CPU, _ ...any) error {
	c.SetIYL(c.bus.Read(c.PC + 2))
	return nil
}

// LD r,IYH / LD r,IYL

func fdLdBIYH(c *CPU) error { c.B = c.IYH(); return nil }
func fdLdBIYL(c *CPU) error { c.B = c.IYL(); return nil }
func fdLdCIYH(c *CPU) error { c.C = c.IYH(); return nil }
func fdLdCIYL(c *CPU) error { c.C = c.IYL(); return nil }
func fdLdDIYH(c *CPU) error { c.D = c.IYH(); return nil }
func fdLdDIYL(c *CPU) error { c.D = c.IYL(); return nil }
func fdLdEIYH(c *CPU) error { c.E = c.IYH(); return nil }
func fdLdEIYL(c *CPU) error { c.E = c.IYL(); return nil }
func fdLdAIYH(c *CPU) error { c.A = c.IYH(); return nil }
func fdLdAIYL(c *CPU) error { c.A = c.IYL(); return nil }

// LD IYH,r / LD IYL,r

func fdLdIYHB(c *CPU) error   { c.SetIYH(c.B); return nil }
func fdLdIYHC(c *CPU) error   { c.SetIYH(c.C); return nil }
func fdLdIYHD(c *CPU) error   { c.SetIYH(c.D); return nil }
func fdLdIYHE(c *CPU) error   { c.SetIYH(c.E); return nil }
func fdLdIYHIYH(_ *CPU) error { return nil } // LD IYH,IYH = NOP
func fdLdIYHIYL(c *CPU) error { c.SetIYH(c.IYL()); return nil }
func fdLdIYHA(c *CPU) error   { c.SetIYH(c.A); return nil }

func fdLdIYLB(c *CPU) error   { c.SetIYL(c.B); return nil }
func fdLdIYLC(c *CPU) error   { c.SetIYL(c.C); return nil }
func fdLdIYLD(c *CPU) error   { c.SetIYL(c.D); return nil }
func fdLdIYLE(c *CPU) error   { c.SetIYL(c.E); return nil }
func fdLdIYLIYH(c *CPU) error { c.SetIYL(c.IYH()); return nil }
func fdLdIYLIYL(_ *CPU) error { return nil } // LD IYL,IYL = NOP
func fdLdIYLA(c *CPU) error   { c.SetIYL(c.A); return nil }

// Arithmetic operations with IYH/IYL

func fdAddAIYH(c *CPU) error { c.A = c.add8(c.A, c.IYH()); return nil }
func fdAddAIYL(c *CPU) error { c.A = c.add8(c.A, c.IYL()); return nil }
func fdAdcAIYH(c *CPU) error { c.A = c.adc(c.A, c.IYH()); return nil }
func fdAdcAIYL(c *CPU) error { c.A = c.adc(c.A, c.IYL()); return nil }
func fdSubIYH(c *CPU) error  { c.A = c.sub8(c.A, c.IYH()); return nil }
func fdSubIYL(c *CPU) error  { c.A = c.sub8(c.A, c.IYL()); return nil }
func fdSbcAIYH(c *CPU) error { c.A = c.sbc(c.A, c.IYH()); return nil }
func fdSbcAIYL(c *CPU) error { c.A = c.sbc(c.A, c.IYL()); return nil }
func fdAndIYH(c *CPU) error  { c.A = c.and8(c.A, c.IYH()); return nil }
func fdAndIYL(c *CPU) error  { c.A = c.and8(c.A, c.IYL()); return nil }
func fdXorIYH(c *CPU) error  { c.A = c.xor8(c.A, c.IYH()); return nil }
func fdXorIYL(c *CPU) error  { c.A = c.xor8(c.A, c.IYL()); return nil }
func fdOrIYH(c *CPU) error   { c.A = c.or8(c.A, c.IYH()); return nil }
func fdOrIYL(c *CPU) error   { c.A = c.or8(c.A, c.IYL()); return nil }
func fdCpIYH(c *CPU) error   { c.cp(c.A, c.IYH()); return nil }
func fdCpIYL(c *CPU) error   { c.cp(c.A, c.IYL()); return nil }
