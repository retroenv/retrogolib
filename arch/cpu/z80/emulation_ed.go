package z80

// ED prefix instruction implementations - Extended instructions

// 16-bit ADC/SBC HL helper methods
func (cpu *CPU) adcHL(value uint16) {
	hl := uint16(cpu.H)<<8 | uint16(cpu.L)
	carry := uint32(cpu.Flags.C)
	result := uint32(hl) + uint32(value) + carry

	cpu.H = uint8(result >> 8)
	cpu.L = uint8(result)

	cpu.setSZ(uint8(result >> 8))
	cpu.setC(result > 0xFFFF)
	cpu.setH((hl&0x0FFF)+(value&0x0FFF)+uint16(carry) > 0x0FFF)
	cpu.setPOverflow(((hl ^ value ^ 0x8000) & (uint16(result) ^ hl) & 0x8000) != 0)
	cpu.setN(false)
}

func (cpu *CPU) sbcHL(value uint16) {
	hl := uint16(cpu.H)<<8 | uint16(cpu.L)
	carry := uint32(cpu.Flags.C)
	result := uint32(hl) - uint32(value) - carry

	cpu.H = uint8(result >> 8)
	cpu.L = uint8(result)

	cpu.setSZ(uint8(result >> 8))
	cpu.setC(result > 0xFFFF)
	cpu.setH((hl & 0x0FFF) < (value&0x0FFF)+uint16(carry))
	cpu.setPOverflow(((hl ^ value) & (hl ^ uint16(result)) & 0x8000) != 0)
	cpu.setN(true)
}

// edNeg implements ED 44: NEG.
func edNeg(c *CPU) error {
	cpu.A = cpu.neg(cpu.A)
	return nil
}

// Interrupt mode instructions
func edIm0(c *CPU, params ...any) error {
	cpu.im = 0
	return nil
}

func edIm1(c *CPU, params ...any) error {
	cpu.im = 1
	return nil
}

func edIm2(c *CPU, params ...any) error {
	cpu.im = 2
	return nil
}

// I/R register loads
func edLdIA(c *CPU) error {
	cpu.I = cpu.A
	return nil
}

func edLdRA(c *CPU) error {
	cpu.R = cpu.A
	return nil
}

func edLdAI(c *CPU) error {
	cpu.A = cpu.I
	cpu.setSZ(cpu.A)
	cpu.setH(false)
	cpu.setN(false)
	cpu.setPOverflow(cpu.iff2)
	return nil
}

func edLdAR(c *CPU) error {
	cpu.A = cpu.R
	cpu.setSZ(cpu.A)
	cpu.setH(false)
	cpu.setN(false)
	cpu.setPOverflow(cpu.iff2)
	return nil
}

// ED instruction implementations for 16-bit memory operations

// edLdNnBc implements ED 43: LD (nn),BC.
func edLdNnBc(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	cpu.memory.Write(addr, cpu.C)
	cpu.memory.Write(addr+1, cpu.B)
	return nil
}

// edLdNnDe implements ED 53: LD (nn),DE.
func edLdNnDe(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	cpu.memory.Write(addr, cpu.E)
	cpu.memory.Write(addr+1, cpu.D)
	return nil
}

// edLdNnHl implements ED 63: LD (nn),HL.
func edLdNnHl(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	cpu.memory.Write(addr, cpu.L)
	cpu.memory.Write(addr+1, cpu.H)
	return nil
}

// edLdNnSp implements ED 73: LD (nn),SP.
func edLdNnSp(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	cpu.memory.Write(addr, uint8(cpu.SP))
	cpu.memory.Write(addr+1, uint8(cpu.SP>>8))
	return nil
}

// edLdBcNn implements ED 4B: LD BC,(nn).
func edLdBcNn(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	cpu.C = cpu.memory.Read(addr)
	cpu.B = cpu.memory.Read(addr + 1)
	return nil
}

// edLdDeNn implements ED 5B: LD DE,(nn).
func edLdDeNn(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	cpu.E = cpu.memory.Read(addr)
	cpu.D = cpu.memory.Read(addr + 1)
	return nil
}

// edLdHlNn implements ED 6B: LD HL,(nn).
func edLdHlNn(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	cpu.L = cpu.memory.Read(addr)
	cpu.H = cpu.memory.Read(addr + 1)
	return nil
}

// edLdSpNn implements ED 7B: LD SP,(nn).
func edLdSpNn(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	low := cpu.memory.Read(addr)
	high := cpu.memory.Read(addr + 1)
	cpu.SP = uint16(high)<<8 | uint16(low)
	return nil
}

func edAdcHlBc(c *CPU, params ...any) error { cpu.adcHL(uint16(cpu.B)<<8 | uint16(cpu.C)); return nil }
func edAdcHlDe(c *CPU, params ...any) error { cpu.adcHL(uint16(cpu.D)<<8 | uint16(cpu.E)); return nil }
func edAdcHlHl(c *CPU, params ...any) error { cpu.adcHL(uint16(cpu.H)<<8 | uint16(cpu.L)); return nil }
func edAdcHlSp(c *CPU, params ...any) error { cpu.adcHL(cpu.SP); return nil }
func edSbcHlBc(c *CPU, params ...any) error { cpu.sbcHL(uint16(cpu.B)<<8 | uint16(cpu.C)); return nil }
func edSbcHlDe(c *CPU, params ...any) error { cpu.sbcHL(uint16(cpu.D)<<8 | uint16(cpu.E)); return nil }
func edSbcHlHl(c *CPU, params ...any) error { cpu.sbcHL(uint16(cpu.H)<<8 | uint16(cpu.L)); return nil }
func edSbcHlSp(c *CPU, params ...any) error { cpu.sbcHL(cpu.SP); return nil }

// ED block transfer and search operations

// edLdi implements ED A0: LDI (HL),(DE), INC HL, INC DE, DEC BC.
func edLdi(c *CPU) error {
	hl := cpu.HL()
	de := cpu.DE()
	bc := cpu.BC()

	cpu.memory.Write(de, cpu.memory.Read(hl))
	cpu.setHL(hl + 1)
	cpu.setDE(de + 1)
	cpu.setBC(bc - 1)

	cpu.setH(false)
	cpu.setN(false)
	cpu.setPOverflow(bc != 1) // P/V flag set if BC != 0 after decrement
	return nil
}

// edLdd implements ED A8: LDD (HL),(DE), DEC HL, DEC DE, DEC BC.
func edLdd(c *CPU) error {
	hl := cpu.HL()
	de := cpu.DE()
	bc := cpu.BC()

	cpu.memory.Write(de, cpu.memory.Read(hl))
	cpu.setHL(hl - 1)
	cpu.setDE(de - 1)
	cpu.setBC(bc - 1)

	cpu.setH(false)
	cpu.setN(false)
	cpu.setPOverflow(bc != 1) // P/V flag set if BC != 0 after decrement
	return nil
}

// edLdir implements ED B0: LDIR - Repeat LDI until BC=0.
func edLdir(c *CPU) error {
	for cpu.BC() != 0 {
		if err := edLdi(c); err != nil {
			return err
		}
	}
	return nil
}

// edLddr implements ED B8: LDDR - Repeat LDD until BC=0.
func edLddr(c *CPU) error {
	for cpu.BC() != 0 {
		if err := edLdd(c); err != nil {
			return err
		}
	}
	return nil
}

// edCpi implements ED A1: CPI - Compare A with (HL), INC HL, DEC BC.
func edCpi(c *CPU) error {
	hl := cpu.HL()
	bc := cpu.BC()
	memValue := cpu.memory.Read(hl)

	result := cpu.A - memValue
	cpu.setHL(hl + 1)
	cpu.setBC(bc - 1)

	cpu.setSZ(result)
	cpu.setH((cpu.A & 0x0F) < (memValue & 0x0F))
	cpu.setPOverflow(bc != 1) // P/V flag set if BC != 0 after decrement
	cpu.setN(true)
	return nil
}

// edCpd implements ED A9: CPD - Compare A with (HL), DEC HL, DEC BC.
func edCpd(c *CPU) error {
	hl := cpu.HL()
	bc := cpu.BC()
	memValue := cpu.memory.Read(hl)

	result := cpu.A - memValue
	cpu.setHL(hl - 1)
	cpu.setBC(bc - 1)

	cpu.setSZ(result)
	cpu.setH((cpu.A & 0x0F) < (memValue & 0x0F))
	cpu.setPOverflow(bc != 1) // P/V flag set if BC != 0 after decrement
	cpu.setN(true)
	return nil
}

// edCpir implements ED B1: CPIR - Repeat CPI until BC=0 or match found.
func edCpir(c *CPU) error {
	for cpu.BC() != 0 {
		if err := edCpi(c); err != nil {
			return err
		}
		if cpu.Flags.Z != 0 {
			break // Match found
		}
	}
	return nil
}

// edCpdr implements ED B9: CPDR - Repeat CPD until BC=0 or match found.
func edCpdr(c *CPU) error {
	for cpu.BC() != 0 {
		if err := edCpd(c); err != nil {
			return err
		}
		if cpu.Flags.Z != 0 {
			break // Match found
		}
	}
	return nil
}

// ED I/O block operations

// edIni implements ED A2: INI - IN (HL),(C), INC HL, DEC B.
func edIni(c *CPU) error {
	hl := cpu.HL()
	port := cpu.C
	value := cpu.readPort(port)

	cpu.memory.Write(hl, value)
	cpu.setHL(hl + 1)
	cpu.B--

	cpu.setZ(cpu.B)
	cpu.setN(true)
	return nil
}

// edInd implements ED AA: IND - IN (HL),(C), DEC HL, DEC B.
func edInd(c *CPU) error {
	hl := cpu.HL()
	port := cpu.C
	value := cpu.readPort(port)

	cpu.memory.Write(hl, value)
	cpu.setHL(hl - 1)
	cpu.B--

	cpu.setZ(cpu.B)
	cpu.setN(true)
	return nil
}

// edInir implements ED B2: INIR - Repeat INI until B=0.
func edInir(c *CPU) error {
	for cpu.B != 0 {
		if err := edIni(c); err != nil {
			return err
		}
	}
	return nil
}

// edIndr implements ED BA: INDR - Repeat IND until B=0.
func edIndr(c *CPU) error {
	for cpu.B != 0 {
		if err := edInd(c); err != nil {
			return err
		}
	}
	return nil
}

// edOuti implements ED A3: OUTI - OUT (C),(HL), INC HL, DEC B.
func edOuti(c *CPU) error {
	hl := cpu.HL()
	port := cpu.C
	value := cpu.memory.Read(hl)

	cpu.writePort(port, value)
	cpu.setHL(hl + 1)
	cpu.B--

	cpu.setZ(cpu.B)
	cpu.setN(true)
	return nil
}

// edOutd implements ED AB: OUTD - OUT (C),(HL), DEC HL, DEC B.
func edOutd(c *CPU) error {
	hl := cpu.HL()
	port := cpu.C
	value := cpu.memory.Read(hl)

	cpu.writePort(port, value)
	cpu.setHL(hl - 1)
	cpu.B--

	cpu.setZ(cpu.B)
	cpu.setN(true)
	return nil
}

// edOtir implements ED B3: OTIR - Repeat OUTI until B=0.
func edOtir(c *CPU) error {
	for cpu.B != 0 {
		if err := edOuti(c); err != nil {
			return err
		}
	}
	return nil
}

// edOtdr implements ED BB: OTDR - Repeat OUTD until B=0.
func edOtdr(c *CPU) error {
	for cpu.B != 0 {
		if err := edOutd(c); err != nil {
			return err
		}
	}
	return nil
}

// ED I/O operations with C register

// edInBC implements ED 40: IN B,(C).
func edInBC(c *CPU, params ...any) error {
	value := cpu.readPort(cpu.C)
	cpu.B = value
	cpu.setSZP(value)
	cpu.setH(false)
	cpu.setN(false)
	return nil
}

// edInCC implements ED 48: IN C,(C).
func edInCC(c *CPU, params ...any) error {
	value := cpu.readPort(cpu.C)
	cpu.C = value
	cpu.setSZP(value)
	cpu.setH(false)
	cpu.setN(false)
	return nil
}

// edInDC implements ED 50: IN D,(C).
func edInDC(c *CPU, params ...any) error {
	value := cpu.readPort(cpu.C)
	cpu.D = value
	cpu.setSZP(value)
	cpu.setH(false)
	cpu.setN(false)
	return nil
}

// edInEC implements ED 58: IN E,(C).
func edInEC(c *CPU, params ...any) error {
	value := cpu.readPort(cpu.C)
	cpu.E = value
	cpu.setSZP(value)
	cpu.setH(false)
	cpu.setN(false)
	return nil
}

// edInHC implements ED 60: IN H,(C).
func edInHC(c *CPU, params ...any) error {
	value := cpu.readPort(cpu.C)
	cpu.H = value
	cpu.setSZP(value)
	cpu.setH(false)
	cpu.setN(false)
	return nil
}

// edInLC implements ED 68: IN L,(C).
func edInLC(c *CPU, params ...any) error {
	value := cpu.readPort(cpu.C)
	cpu.L = value
	cpu.setSZP(value)
	cpu.setH(false)
	cpu.setN(false)
	return nil
}

// edInAC implements ED 78: IN A,(C).
func edInAC(c *CPU, params ...any) error {
	value := cpu.readPort(cpu.C)
	cpu.A = value
	cpu.setSZP(value)
	cpu.setH(false)
	cpu.setN(false)
	return nil
}

// edOutCB implements ED 41: OUT (C),B.
func edOutCB(c *CPU, params ...any) error {
	cpu.writePort(cpu.C, cpu.B)
	return nil
}

// edOutCC implements ED 49: OUT (C),C.
func edOutCC(c *CPU, params ...any) error {
	cpu.writePort(cpu.C, cpu.C)
	return nil
}

// edOutCD implements ED 51: OUT (C),D.
func edOutCD(c *CPU, params ...any) error {
	cpu.writePort(cpu.C, cpu.D)
	return nil
}

// edOutCE implements ED 59: OUT (C),E.
func edOutCE(c *CPU, params ...any) error {
	cpu.writePort(cpu.C, cpu.E)
	return nil
}

// edOutCH implements ED 61: OUT (C),H.
func edOutCH(c *CPU, params ...any) error {
	cpu.writePort(cpu.C, cpu.H)
	return nil
}

// edOutCL implements ED 69: OUT (C),L.
func edOutCL(c *CPU, params ...any) error {
	cpu.writePort(cpu.C, cpu.L)
	return nil
}

// edOutCA implements ED 79: OUT (C),A.
func edOutCA(c *CPU, params ...any) error {
	cpu.writePort(cpu.C, cpu.A)
	return nil
}

// Return and rotate operations
func edRetn(c *CPU) error {
	cpu.iff1 = cpu.iff2
	cpu.PC = cpu.pop16()
	return nil
}

func edReti(c *CPU) error {
	cpu.iff1 = cpu.iff2
	cpu.PC = cpu.pop16()
	return nil
}

// edRrd implements ED 67: RRD - Rotate Right Decimal.
// The contents of A and (HL) are rotated right 4 bits.
func edRrd(c *CPU) error {
	hl := cpu.HL()
	memValue := cpu.memory.Read(hl)

	// Rotate A and (HL) right 4 bits
	lowNibbleA := cpu.A & 0x0F

	cpu.A = (cpu.A & 0xF0) | (memValue & 0x0F)
	cpu.memory.Write(hl, (lowNibbleA<<4)|(memValue>>4))

	cpu.setSZP(cpu.A)
	cpu.setH(false)
	cpu.setN(false)
	return nil
}

// edRld implements ED 6F: RLD - Rotate Left Decimal.
// The contents of A and (HL) are rotated left 4 bits.
func edRld(c *CPU) error {
	hl := cpu.HL()
	memValue := cpu.memory.Read(hl)

	// Rotate A and (HL) left 4 bits
	lowNibbleA := cpu.A & 0x0F
	highNibbleMem := memValue >> 4

	cpu.A = (cpu.A & 0xF0) | highNibbleMem
	cpu.memory.Write(hl, ((memValue&0x0F)<<4)|lowNibbleA)

	cpu.setSZP(cpu.A)
	cpu.setH(false)
	cpu.setN(false)
	return nil
}
