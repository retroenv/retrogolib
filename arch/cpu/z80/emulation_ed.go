package z80

import "math/bits"

// ED prefix instruction implementations - Extended instructions

// 16-bit ADC/SBC HL helper methods
func (c *CPU) adcHL(value uint16) {
	hl := uint16(c.H)<<8 | uint16(c.L)
	c.MEMPTR = hl + 1
	carry := uint32(c.Flags.C)
	result := uint32(hl) + uint32(value) + carry

	c.H = uint8(result >> 8)
	c.L = uint8(result)

	r16 := uint16(result)
	c.setS(uint8(r16 >> 8))
	setFlag(&c.Flags.Z, r16 == 0)
	c.setXY(uint8(r16 >> 8))
	c.setC(result > 0xFFFF)
	c.setH((hl&0x0FFF)+(value&0x0FFF)+uint16(carry) > 0x0FFF)
	c.setPOverflow(((hl ^ value ^ 0x8000) & (r16 ^ hl) & 0x8000) != 0)
	c.setN(false)
}

func (c *CPU) sbcHL(value uint16) {
	hl := uint16(c.H)<<8 | uint16(c.L)
	c.MEMPTR = hl + 1
	carry := uint32(c.Flags.C)
	result := uint32(hl) - uint32(value) - carry

	c.H = uint8(result >> 8)
	c.L = uint8(result)

	r16 := uint16(result)
	c.setS(uint8(r16 >> 8))
	setFlag(&c.Flags.Z, r16 == 0)
	c.setXY(uint8(r16 >> 8))
	c.setC(result > 0xFFFF)
	c.setH((hl & 0x0FFF) < (value&0x0FFF)+uint16(carry))
	c.setPOverflow(((hl ^ value) & (hl ^ r16) & 0x8000) != 0)
	c.setN(true)
}

// setIOBlockFlags sets the full undocumented flag behavior for INI/IND/OUTI/OUTD.
// value is the byte transferred, k is value + secondary (where secondary depends on instruction).
func (c *CPU) setIOBlockFlags(value uint8, k uint16) {
	c.setS(c.B)
	c.setZ(c.B)
	c.setXY(c.B)
	setFlag(&c.Flags.N, value&0x80 != 0) // N = bit 7 of transferred value
	carry := k > 255
	c.setH(carry)
	c.setC(carry)
	// P/V = parity of ((k & 7) ^ B)
	c.setP(uint8(k&7) ^ c.B)
}

// adjustIORepeatFlags applies additional flag modifications for repeat I/O instructions
// (INIR/INDR/OTIR/OTDR) when the instruction will repeat (B != 0).
// The flags PF and HF undergo additional transformations based on carry and data bit 7.
func (c *CPU) adjustIORepeatFlags(value uint8, k uint16) {
	carry := k > 255
	dataBit7 := value&0x80 != 0

	switch {
	case carry && dataBit7:
		c.Flags.P ^= parityByte((c.B - 1) & 0x07)
		c.Flags.P ^= 1
		setFlag(&c.Flags.H, (c.B&0x0F) == 0x00)
	case carry:
		c.Flags.P ^= parityByte((c.B + 1) & 0x07)
		c.Flags.P ^= 1
		setFlag(&c.Flags.H, (c.B&0x0F) == 0x0F)
	default:
		c.Flags.P ^= parityByte(c.B & 0x07)
		c.Flags.P ^= 1
	}
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
func edLdNnBc(c *CPU, _ ...any) error {
	addr := c.read16(c.PC + 2)
	c.writeRegisterPair(addr, c.C, c.B)
	c.MEMPTR = addr + 1
	return nil
}

// edLdNnDe implements ED 53: LD (nn),DE.
func edLdNnDe(c *CPU, _ ...any) error {
	addr := c.read16(c.PC + 2)
	c.writeRegisterPair(addr, c.E, c.D)
	c.MEMPTR = addr + 1
	return nil
}

// edLdNnHl implements ED 63: LD (nn),HL.
func edLdNnHl(c *CPU, _ ...any) error {
	addr := c.read16(c.PC + 2)
	c.writeRegisterPair(addr, c.L, c.H)
	c.MEMPTR = addr + 1
	return nil
}

// edLdNnSp implements ED 73: LD (nn),SP.
func edLdNnSp(c *CPU, _ ...any) error {
	addr := c.read16(c.PC + 2)
	c.memory.Write(addr, uint8(c.SP))
	c.memory.Write(addr+1, uint8(c.SP>>8))
	c.MEMPTR = addr + 1
	return nil
}

// edLdBcNn implements ED 4B: LD BC,(nn).
func edLdBcNn(c *CPU, _ ...any) error {
	addr := c.read16(c.PC + 2)
	value := c.read16(addr)
	c.setBC(value)
	c.MEMPTR = addr + 1
	return nil
}

// edLdDeNn implements ED 5B: LD DE,(nn).
func edLdDeNn(c *CPU, _ ...any) error {
	addr := c.read16(c.PC + 2)
	value := c.read16(addr)
	c.setDE(value)
	c.MEMPTR = addr + 1
	return nil
}

// edLdHlNn implements ED 6B: LD HL,(nn).
func edLdHlNn(c *CPU, _ ...any) error {
	addr := c.read16(c.PC + 2)
	value := c.read16(addr)
	c.setHL(value)
	c.MEMPTR = addr + 1
	return nil
}

// edLdSpNn implements ED 7B: LD SP,(nn).
func edLdSpNn(c *CPU, _ ...any) error {
	addr := c.read16(c.PC + 2)
	c.SP = c.read16(addr)
	c.MEMPTR = addr + 1
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

	transferred := c.memory.Read(hl)
	c.memory.Write(de, transferred)
	c.setHL(hl + 1)
	c.setDE(de + 1)
	c.setBC(bc - 1)

	c.setH(false)
	c.setN(false)
	c.setPOverflow(bc != 1)

	// Undocumented: X/Y from (transferred + A)
	n := transferred + c.A
	c.Flags.X = (n >> 3) & 1 // bit 3
	c.Flags.Y = (n >> 1) & 1 // bit 1 → Y (bit 5 position)
	return nil
}

// edLdd implements ED A8: LDD (HL),(DE), DEC HL, DEC DE, DEC BC.
func edLdd(c *CPU) error {
	hl := c.hl()
	de := c.de()
	bc := c.bc()

	transferred := c.memory.Read(hl)
	c.memory.Write(de, transferred)
	c.setHL(hl - 1)
	c.setDE(de - 1)
	c.setBC(bc - 1)

	c.setH(false)
	c.setN(false)
	c.setPOverflow(bc != 1)

	n := transferred + c.A
	c.Flags.X = (n >> 3) & 1
	c.Flags.Y = (n >> 1) & 1
	return nil
}

// edLdir implements ED B0: LDIR - Execute one LDI iteration.
// If BC != 0 after, PC stays at LDIR for the next Step to repeat.
func edLdir(c *CPU) error {
	if err := edLdi(c); err != nil {
		return err
	}
	if c.bc() != 0 {
		c.cycles += 5
		c.MEMPTR = c.PC + 1
		c.setXY(uint8(c.PC >> 8))
	} else {
		c.PC += 2
	}
	return nil
}

// edLddr implements ED B8: LDDR - Execute one LDD iteration.
// If BC != 0 after, PC stays at LDDR for the next Step to repeat.
func edLddr(c *CPU) error {
	if err := edLdd(c); err != nil {
		return err
	}
	if c.bc() != 0 {
		c.cycles += 5
		c.MEMPTR = c.PC + 1
		c.setXY(uint8(c.PC >> 8))
	} else {
		c.PC += 2
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
	c.MEMPTR++

	hf := (c.A & 0x0F) < (memValue & 0x0F)
	c.setS(result)
	c.setZ(result)
	c.setH(hf)
	c.setPOverflow(bc != 1)
	c.setN(true)

	// Undocumented: X/Y from (A - operand - HF)
	n := result
	if hf {
		n--
	}
	c.Flags.X = (n >> 3) & 1
	c.Flags.Y = (n >> 1) & 1
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
	c.MEMPTR--

	hf := (c.A & 0x0F) < (memValue & 0x0F)
	c.setS(result)
	c.setZ(result)
	c.setH(hf)
	c.setPOverflow(bc != 1)
	c.setN(true)

	n := result
	if hf {
		n--
	}
	c.Flags.X = (n >> 3) & 1
	c.Flags.Y = (n >> 1) & 1
	return nil
}

// edCpir implements ED B1: CPIR - Execute one CPI iteration.
// If BC != 0 and no match, PC stays at CPIR for the next Step to repeat.
func edCpir(c *CPU) error {
	if err := edCpi(c); err != nil {
		return err
	}
	if c.bc() != 0 && c.Flags.Z == 0 {
		c.cycles += 5
		c.MEMPTR = c.PC + 1
		c.setXY(uint8(c.PC >> 8))
	} else {
		c.PC += 2
	}
	return nil
}

// edCpdr implements ED B9: CPDR - Execute one CPD iteration.
// If BC != 0 and no match, PC stays at CPDR for the next Step to repeat.
func edCpdr(c *CPU) error {
	if err := edCpd(c); err != nil {
		return err
	}
	if c.bc() != 0 && c.Flags.Z == 0 {
		c.cycles += 5
		c.MEMPTR = c.PC + 1
		c.setXY(uint8(c.PC >> 8))
	} else {
		c.PC += 2
	}
	return nil
}

// ED I/O block operations

// parityByte returns 1 for even parity, 0 for odd parity.
func parityByte(v uint8) uint8 {
	if bits.OnesCount8(v)%2 == 0 {
		return 1
	}
	return 0
}

// edIni implements ED A2: INI - IN (HL),(C), INC HL, DEC B.
func edIni(c *CPU) error {
	c.MEMPTR = c.bc() + 1
	hl := c.hl()
	value := c.readPort(c.C)

	c.memory.Write(hl, value)
	c.setHL(hl + 1)
	c.B--

	k := uint16(value) + uint16((c.C+1)&0xFF)
	c.setIOBlockFlags(value, k)
	return nil
}

// edInd implements ED AA: IND - IN (HL),(C), DEC HL, DEC B.
func edInd(c *CPU) error {
	c.MEMPTR = c.bc() - 1
	hl := c.hl()
	value := c.readPort(c.C)

	c.memory.Write(hl, value)
	c.setHL(hl - 1)
	c.B--

	k := uint16(value) + uint16((c.C-1)&0xFF)
	c.setIOBlockFlags(value, k)
	return nil
}

// edInir implements ED B2: INIR - Execute one INI iteration.
// If B != 0 after, PC stays at INIR for the next Step to repeat.
func edInir(c *CPU) error {
	if err := edIni(c); err != nil {
		return err
	}
	if c.B != 0 {
		c.cycles += 5
		c.MEMPTR = c.PC + 1
		c.setXY(uint8(c.PC >> 8))
		// Reconstruct port value from memory (INI wrote it to HL-1 since HL was incremented).
		portVal := c.memory.Read(c.hl() - 1)
		k := uint16(portVal) + uint16((c.C+1)&0xFF)
		c.adjustIORepeatFlags(portVal, k)
	} else {
		c.PC += 2
	}
	return nil
}

// edIndr implements ED BA: INDR - Execute one IND iteration.
// If B != 0 after, PC stays at INDR for the next Step to repeat.
func edIndr(c *CPU) error {
	if err := edInd(c); err != nil {
		return err
	}
	if c.B != 0 {
		c.cycles += 5
		c.MEMPTR = c.PC + 1
		c.setXY(uint8(c.PC >> 8))
		// Reconstruct port value from memory (IND wrote it to HL+1 since HL was decremented).
		portVal := c.memory.Read(c.hl() + 1)
		k := uint16(portVal) + uint16((c.C-1)&0xFF)
		c.adjustIORepeatFlags(portVal, k)
	} else {
		c.PC += 2
	}
	return nil
}

// edOuti implements ED A3: OUTI - OUT (C),(HL), INC HL, DEC B.
func edOuti(c *CPU) error {
	hl := c.hl()
	value := c.memory.Read(hl)

	c.writePort(c.C, value)
	c.setHL(hl + 1)
	c.B--

	c.MEMPTR = c.bc() + 1
	lAfter := uint8(hl + 1) // L after increment
	k := uint16(value) + uint16(lAfter)
	c.setIOBlockFlags(value, k)
	return nil
}

// edOutd implements ED AB: OUTD - OUT (C),(HL), DEC HL, DEC B.
func edOutd(c *CPU) error {
	hl := c.hl()
	value := c.memory.Read(hl)

	c.writePort(c.C, value)
	c.setHL(hl - 1)
	c.B--

	c.MEMPTR = c.bc() - 1
	lAfter := uint8(hl - 1) // L after decrement
	k := uint16(value) + uint16(lAfter)
	c.setIOBlockFlags(value, k)
	return nil
}

// edOtir implements ED B3: OTIR - Execute one OUTI iteration.
// If B != 0 after, PC stays at OTIR for the next Step to repeat.
func edOtir(c *CPU) error {
	if err := edOuti(c); err != nil {
		return err
	}
	if c.B != 0 {
		c.cycles += 5
		c.MEMPTR = c.PC + 1
		c.setXY(uint8(c.PC >> 8))
		// Reconstruct: OUTI read from (HL-1) since HL was incremented.
		value := c.memory.Read(c.hl() - 1)
		lAfter := uint8(c.hl()) // L after OUTI incremented HL
		k := uint16(value) + uint16(lAfter)
		c.adjustIORepeatFlags(value, k)
	} else {
		c.PC += 2
	}
	return nil
}

// edOtdr implements ED BB: OTDR - Execute one OUTD iteration.
// If B != 0 after, PC stays at OTDR for the next Step to repeat.
func edOtdr(c *CPU) error {
	if err := edOutd(c); err != nil {
		return err
	}
	if c.B != 0 {
		c.cycles += 5
		c.MEMPTR = c.PC + 1
		c.setXY(uint8(c.PC >> 8))
		// Reconstruct: OUTD read from (HL+1) since HL was decremented.
		value := c.memory.Read(c.hl() + 1)
		lAfter := uint8(c.hl()) // L after OUTD decremented HL
		k := uint16(value) + uint16(lAfter)
		c.adjustIORepeatFlags(value, k)
	} else {
		c.PC += 2
	}
	return nil
}

// ED I/O operations with C register

// edInFC implements ED 70: IN F,(C) - undocumented.
// Reads port C, sets flags from result, discards the value.
func edInFC(c *CPU, _ ...any) error {
	c.MEMPTR = c.bc() + 1
	value := c.readPort(c.C)
	c.setSZP(value)
	c.setH(false)
	c.setN(false)
	return nil
}

// edOut0C implements ED 71: OUT (C),0 - undocumented.
// Outputs 0 to port C.
func edOut0C(c *CPU, _ ...any) error {
	c.MEMPTR = c.bc() + 1
	c.writePort(c.C, 0)
	return nil
}

// edNop implements undocumented ED NOP instructions (ED 77, ED 7F).
func edNop(_ *CPU, _ ...any) error {
	return nil
}

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
	c.MEMPTR = c.bc() + 1
	c.writePort(c.C, c.B)
	return nil
}

// edOutCC implements ED 49: OUT (C),C.
func edOutCC(c *CPU, _ ...any) error {
	c.MEMPTR = c.bc() + 1
	c.writePort(c.C, c.C)
	return nil
}

// edOutCD implements ED 51: OUT (C),D.
func edOutCD(c *CPU, _ ...any) error {
	c.MEMPTR = c.bc() + 1
	c.writePort(c.C, c.D)
	return nil
}

// edOutCE implements ED 59: OUT (C),E.
func edOutCE(c *CPU, _ ...any) error {
	c.MEMPTR = c.bc() + 1
	c.writePort(c.C, c.E)
	return nil
}

// edOutCH implements ED 61: OUT (C),H.
func edOutCH(c *CPU, _ ...any) error {
	c.MEMPTR = c.bc() + 1
	c.writePort(c.C, c.H)
	return nil
}

// edOutCL implements ED 69: OUT (C),L.
func edOutCL(c *CPU, _ ...any) error {
	c.MEMPTR = c.bc() + 1
	c.writePort(c.C, c.L)
	return nil
}

// edOutCA implements ED 79: OUT (C),A.
func edOutCA(c *CPU, _ ...any) error {
	c.MEMPTR = c.bc() + 1
	c.writePort(c.C, c.A)
	return nil
}

// Return and rotate operations
func edRetn(c *CPU) error {
	c.iff1 = c.iff2
	c.PC = c.pop16()
	c.MEMPTR = c.PC
	return nil
}

func edReti(c *CPU) error {
	c.iff1 = c.iff2
	c.PC = c.pop16()
	c.MEMPTR = c.PC
	return nil
}

// edRrd implements ED 67: RRD - Rotate Right Decimal.
func edRrd(c *CPU) error {
	hl := c.hl()
	c.MEMPTR = hl + 1
	memValue := c.memory.Read(hl)

	lowNibbleA := c.A & 0x0F
	c.A = (c.A & 0xF0) | (memValue & 0x0F)
	c.memory.Write(hl, (lowNibbleA<<4)|(memValue>>4))

	c.setSZP(c.A)
	c.setH(false)
	c.setN(false)
	return nil
}

// edRld implements ED 6F: RLD - Rotate Left Decimal.
func edRld(c *CPU) error {
	hl := c.hl()
	c.MEMPTR = hl + 1
	memValue := c.memory.Read(hl)

	lowNibbleA := c.A & 0x0F
	highNibbleMem := memValue >> 4
	c.A = (c.A & 0xF0) | highNibbleMem
	c.memory.Write(hl, ((memValue&0x0F)<<4)|lowNibbleA)

	c.setSZP(c.A)
	c.setH(false)
	c.setN(false)
	return nil
}
