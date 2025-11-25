package z80

// ED prefix instruction implementations - Extended instructions

// 16-bit ADC/SBC HL helper methods
func (c *CPU) adcHL(value uint16) {
	hl := uint16(c.H)<<8 | uint16(c.L)
	carry := uint32(c.Flags.C)
	result := uint32(hl) + uint32(value) + carry

	c.H = uint8(result >> 8)
	c.L = uint8(result)

	c.setSZ(uint8(result >> 8))
	c.setC(result > 0xFFFF)
	c.setH((hl&0x0FFF)+(value&0x0FFF)+uint16(carry) > 0x0FFF)
	c.setPOverflow(((hl ^ value ^ 0x8000) & (uint16(result) ^ hl) & 0x8000) != 0)
	c.setN(false)
}

func (c *CPU) sbcHL(value uint16) {
	hl := uint16(c.H)<<8 | uint16(c.L)
	carry := uint32(c.Flags.C)
	result := uint32(hl) - uint32(value) - carry

	c.H = uint8(result >> 8)
	c.L = uint8(result)

	c.setSZ(uint8(result >> 8))
	c.setC(result > 0xFFFF)
	c.setH((hl & 0x0FFF) < (value&0x0FFF)+uint16(carry))
	c.setPOverflow(((hl ^ value) & (hl ^ uint16(result)) & 0x8000) != 0)
	c.setN(true)
}

// edNeg implements ED 44: NEG.
func edNeg(c *CPU) error {
	c.A = c.neg(c.A)
	return nil
}

// Interrupt mode instructions
func edIm0(c *CPU, _ ...any) error {
	c.im = 0
	return nil
}

func edIm1(c *CPU, _ ...any) error {
	c.im = 1
	return nil
}

func edIm2(c *CPU, _ ...any) error {
	c.im = 2
	return nil
}

// I/R register loads
func edLdIA(c *CPU) error {
	c.I = c.A
	return nil
}

func edLdRA(c *CPU) error {
	c.R = c.A
	return nil
}

func edLdAI(c *CPU) error {
	c.A = c.I
	c.setSZ(c.A)
	c.setH(false)
	c.setN(false)
	c.setPOverflow(c.iff2)
	return nil
}

func edLdAR(c *CPU) error {
	c.A = c.R
	c.setSZ(c.A)
	c.setH(false)
	c.setN(false)
	c.setPOverflow(c.iff2)
	return nil
}

// ED instruction implementations for 16-bit memory operations

// edLdNnBc implements ED 43: LD (nn),BC.
func edLdNnBc(c *CPU, params ...any) error {
	addr := extractExtendedAddress(params...)
	c.writeRegisterPair(addr, c.C, c.B)
	return nil
}

// edLdNnDe implements ED 53: LD (nn),DE.
func edLdNnDe(c *CPU, params ...any) error {
	addr := extractExtendedAddress(params...)
	c.writeRegisterPair(addr, c.E, c.D)
	return nil
}

// edLdNnHl implements ED 63: LD (nn),HL.
func edLdNnHl(c *CPU, params ...any) error {
	addr := extractExtendedAddress(params...)
	c.writeRegisterPair(addr, c.L, c.H)
	return nil
}

// edLdNnSp implements ED 73: LD (nn),SP.
func edLdNnSp(c *CPU, params ...any) error {
	addr := extractExtendedAddress(params...)
	c.memory.Write(addr, uint8(c.SP))
	c.memory.Write(addr+1, uint8(c.SP>>8))
	return nil
}

// edLdBcNn implements ED 4B: LD BC,(nn).
func edLdBcNn(c *CPU, params ...any) error {
	addr := extractExtendedAddress(params...)
	value := c.read16(addr)
	c.setBC(value)
	return nil
}

// edLdDeNn implements ED 5B: LD DE,(nn).
func edLdDeNn(c *CPU, params ...any) error {
	addr := extractExtendedAddress(params...)
	value := c.read16(addr)
	c.setDE(value)
	return nil
}

// edLdHlNn implements ED 6B: LD HL,(nn).
func edLdHlNn(c *CPU, params ...any) error {
	addr := extractExtendedAddress(params...)
	value := c.read16(addr)
	c.setHL(value)
	return nil
}

// edLdSpNn implements ED 7B: LD SP,(nn).
func edLdSpNn(c *CPU, params ...any) error {
	addr := extractExtendedAddress(params...)
	c.SP = c.read16(addr)
	return nil
}

func edAdcHlBc(c *CPU, _ ...any) error { c.adcHL(uint16(c.B)<<8 | uint16(c.C)); return nil }
func edAdcHlDe(c *CPU, _ ...any) error { c.adcHL(uint16(c.D)<<8 | uint16(c.E)); return nil }
func edAdcHlHl(c *CPU, _ ...any) error { c.adcHL(uint16(c.H)<<8 | uint16(c.L)); return nil }
func edAdcHlSp(c *CPU, _ ...any) error { c.adcHL(c.SP); return nil }
func edSbcHlBc(c *CPU, _ ...any) error { c.sbcHL(uint16(c.B)<<8 | uint16(c.C)); return nil }
func edSbcHlDe(c *CPU, _ ...any) error { c.sbcHL(uint16(c.D)<<8 | uint16(c.E)); return nil }
func edSbcHlHl(c *CPU, _ ...any) error { c.sbcHL(uint16(c.H)<<8 | uint16(c.L)); return nil }
func edSbcHlSp(c *CPU, _ ...any) error { c.sbcHL(c.SP); return nil }

// ED block transfer and search operations

// edLdi implements ED A0: LDI (HL),(DE), INC HL, INC DE, DEC BC.
func edLdi(c *CPU) error {
	hl := c.hl()
	de := c.de()
	bc := c.bc()

	c.memory.Write(de, c.memory.Read(hl))
	c.setHL(hl + 1)
	c.setDE(de + 1)
	c.setBC(bc - 1)

	c.setH(false)
	c.setN(false)
	c.setPOverflow(bc != 1) // P/V flag set if BC != 0 after decrement
	return nil
}

// edLdd implements ED A8: LDD (HL),(DE), DEC HL, DEC DE, DEC BC.
func edLdd(c *CPU) error {
	hl := c.hl()
	de := c.de()
	bc := c.bc()

	c.memory.Write(de, c.memory.Read(hl))
	c.setHL(hl - 1)
	c.setDE(de - 1)
	c.setBC(bc - 1)

	c.setH(false)
	c.setN(false)
	c.setPOverflow(bc != 1) // P/V flag set if BC != 0 after decrement
	return nil
}

// edLdir implements ED B0: LDIR - Repeat LDI until BC=0.
func edLdir(c *CPU) error {
	for c.bc() != 0 {
		if err := edLdi(c); err != nil {
			return err
		}
	}
	return nil
}

// edLddr implements ED B8: LDDR - Repeat LDD until BC=0.
func edLddr(c *CPU) error {
	for c.bc() != 0 {
		if err := edLdd(c); err != nil {
			return err
		}
	}
	return nil
}

// edCpi implements ED A1: CPI - Compare A with (HL), INC HL, DEC BC.
func edCpi(c *CPU) error {
	hl := c.hl()
	bc := c.bc()
	memValue := c.memory.Read(hl)

	result := c.A - memValue
	c.setHL(hl + 1)
	c.setBC(bc - 1)

	c.setSZ(result)
	c.setH((c.A & 0x0F) < (memValue & 0x0F))
	c.setPOverflow(bc != 1) // P/V flag set if BC != 0 after decrement
	c.setN(true)
	return nil
}

// edCpd implements ED A9: CPD - Compare A with (HL), DEC HL, DEC BC.
func edCpd(c *CPU) error {
	hl := c.hl()
	bc := c.bc()
	memValue := c.memory.Read(hl)

	result := c.A - memValue
	c.setHL(hl - 1)
	c.setBC(bc - 1)

	c.setSZ(result)
	c.setH((c.A & 0x0F) < (memValue & 0x0F))
	c.setPOverflow(bc != 1) // P/V flag set if BC != 0 after decrement
	c.setN(true)
	return nil
}

// edCpir implements ED B1: CPIR - Repeat CPI until BC=0 or match found.
func edCpir(c *CPU) error {
	for c.bc() != 0 {
		if err := edCpi(c); err != nil {
			return err
		}
		if c.Flags.Z != 0 {
			break // Match found
		}
	}
	return nil
}

// edCpdr implements ED B9: CPDR - Repeat CPD until BC=0 or match found.
func edCpdr(c *CPU) error {
	for c.bc() != 0 {
		if err := edCpd(c); err != nil {
			return err
		}
		if c.Flags.Z != 0 {
			break // Match found
		}
	}
	return nil
}

// ED I/O block operations

// edIni implements ED A2: INI - IN (HL),(C), INC HL, DEC B.
func edIni(c *CPU) error {
	hl := c.hl()
	port := c.C
	value := c.readPort(port)

	c.memory.Write(hl, value)
	c.setHL(hl + 1)
	c.B--

	c.setZ(c.B)
	c.setN(true)
	return nil
}

// edInd implements ED AA: IND - IN (HL),(C), DEC HL, DEC B.
func edInd(c *CPU) error {
	hl := c.hl()
	port := c.C
	value := c.readPort(port)

	c.memory.Write(hl, value)
	c.setHL(hl - 1)
	c.B--

	c.setZ(c.B)
	c.setN(true)
	return nil
}

// edInir implements ED B2: INIR - Repeat INI until B=0.
func edInir(c *CPU) error {
	for c.B != 0 {
		if err := edIni(c); err != nil {
			return err
		}
	}
	return nil
}

// edIndr implements ED BA: INDR - Repeat IND until B=0.
func edIndr(c *CPU) error {
	for c.B != 0 {
		if err := edInd(c); err != nil {
			return err
		}
	}
	return nil
}

// edOuti implements ED A3: OUTI - OUT (C),(HL), INC HL, DEC B.
func edOuti(c *CPU) error {
	hl := c.hl()
	port := c.C
	value := c.memory.Read(hl)

	c.writePort(port, value)
	c.setHL(hl + 1)
	c.B--

	c.setZ(c.B)
	c.setN(true)
	return nil
}

// edOutd implements ED AB: OUTD - OUT (C),(HL), DEC HL, DEC B.
func edOutd(c *CPU) error {
	hl := c.hl()
	port := c.C
	value := c.memory.Read(hl)

	c.writePort(port, value)
	c.setHL(hl - 1)
	c.B--

	c.setZ(c.B)
	c.setN(true)
	return nil
}

// edOtir implements ED B3: OTIR - Repeat OUTI until B=0.
func edOtir(c *CPU) error {
	for c.B != 0 {
		if err := edOuti(c); err != nil {
			return err
		}
	}
	return nil
}

// edOtdr implements ED BB: OTDR - Repeat OUTD until B=0.
func edOtdr(c *CPU) error {
	for c.B != 0 {
		if err := edOutd(c); err != nil {
			return err
		}
	}
	return nil
}

// ED I/O operations with C register

// edInBC implements ED 40: IN B,(C).
func edInBC(c *CPU, _ ...any) error {
	c.inPortToRegister(&c.B)
	return nil
}

// edInCC implements ED 48: IN C,(C).
func edInCC(c *CPU, _ ...any) error {
	c.inPortToRegister(&c.C)
	return nil
}

// edInDC implements ED 50: IN D,(C).
func edInDC(c *CPU, _ ...any) error {
	c.inPortToRegister(&c.D)
	return nil
}

// edInEC implements ED 58: IN E,(C).
func edInEC(c *CPU, _ ...any) error {
	c.inPortToRegister(&c.E)
	return nil
}

// edInHC implements ED 60: IN H,(C).
func edInHC(c *CPU, _ ...any) error {
	c.inPortToRegister(&c.H)
	return nil
}

// edInLC implements ED 68: IN L,(C).
func edInLC(c *CPU, _ ...any) error {
	c.inPortToRegister(&c.L)
	return nil
}

// edInAC implements ED 78: IN A,(C).
func edInAC(c *CPU, _ ...any) error {
	c.inPortToRegister(&c.A)
	return nil
}

// edOutCB implements ED 41: OUT (C),B.
func edOutCB(c *CPU, _ ...any) error {
	c.writePort(c.C, c.B)
	return nil
}

// edOutCC implements ED 49: OUT (C),C.
func edOutCC(c *CPU, _ ...any) error {
	c.writePort(c.C, c.C)
	return nil
}

// edOutCD implements ED 51: OUT (C),D.
func edOutCD(c *CPU, _ ...any) error {
	c.writePort(c.C, c.D)
	return nil
}

// edOutCE implements ED 59: OUT (C),E.
func edOutCE(c *CPU, _ ...any) error {
	c.writePort(c.C, c.E)
	return nil
}

// edOutCH implements ED 61: OUT (C),H.
func edOutCH(c *CPU, _ ...any) error {
	c.writePort(c.C, c.H)
	return nil
}

// edOutCL implements ED 69: OUT (C),L.
func edOutCL(c *CPU, _ ...any) error {
	c.writePort(c.C, c.L)
	return nil
}

// edOutCA implements ED 79: OUT (C),A.
func edOutCA(c *CPU, _ ...any) error {
	c.writePort(c.C, c.A)
	return nil
}

// Return and rotate operations
func edRetn(c *CPU) error {
	c.iff1 = c.iff2
	c.PC = c.pop16()
	return nil
}

func edReti(c *CPU) error {
	c.iff1 = c.iff2
	c.PC = c.pop16()
	return nil
}

// edRrd implements ED 67: RRD - Rotate Right Decimal.
// The contents of A and (HL) are rotated right 4 bits.
func edRrd(c *CPU) error {
	hl := c.hl()
	memValue := c.memory.Read(hl)

	// Rotate A and (HL) right 4 bits
	lowNibbleA := c.A & 0x0F

	c.A = (c.A & 0xF0) | (memValue & 0x0F)
	c.memory.Write(hl, (lowNibbleA<<4)|(memValue>>4))

	c.setSZP(c.A)
	c.setH(false)
	c.setN(false)
	return nil
}

// edRld implements ED 6F: RLD - Rotate Left Decimal.
// The contents of A and (HL) are rotated left 4 bits.
func edRld(c *CPU) error {
	hl := c.hl()
	memValue := c.memory.Read(hl)

	// Rotate A and (HL) left 4 bits
	lowNibbleA := c.A & 0x0F
	highNibbleMem := memValue >> 4

	c.A = (c.A & 0xF0) | highNibbleMem
	c.memory.Write(hl, ((memValue&0x0F)<<4)|lowNibbleA)

	c.setSZP(c.A)
	c.setH(false)
	c.setN(false)
	return nil
}
