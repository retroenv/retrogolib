package sm83

// ldImm8 loads an 8-bit immediate value into a register (LD r,n).
// Register encoded in bits 3-5 of the opcode.
func ldImm8(c *CPU, params ...any) error {
	value := uint8(params[0].(Immediate8))
	reg := (c.currentOpcode >> 3) & 0x07
	c.SetRegisterValue(reg, value)
	return nil
}

// ldReg8 loads between 8-bit registers (LD r,r').
// Destination from bits 3-5, source from bits 0-2 of the opcode.
func ldReg8(c *CPU, _ ...any) error {
	dst := (c.currentOpcode >> 3) & 0x07
	src := c.currentOpcode & 0x07
	c.SetRegisterValue(dst, c.GetRegisterValue(src))
	return nil
}

// ldReg16 loads a 16-bit immediate value into a register pair (LD rr,nn).
// Register pair encoded in bits 4-5 of the opcode.
func ldReg16(c *CPU, params ...any) error {
	value := uint16(params[0].(Immediate16))
	switch (c.currentOpcode >> 4) & 0x03 {
	case 0:
		c.setBC(value)
	case 1:
		c.setDE(value)
	case 2:
		c.setHL(value)
	case 3:
		c.SP = value
	}
	return nil
}

// ldIndirect performs indirect loads through BC or DE register pairs.
// 0x02: LD (BC),A  0x0A: LD A,(BC)  0x12: LD (DE),A  0x1A: LD A,(DE)
func ldIndirect(c *CPU, _ ...any) error {
	switch c.currentOpcode {
	case 0x02: // LD (BC),A
		c.memory.Write(c.bc(), c.A)
	case 0x0A: // LD A,(BC)
		c.A = c.memory.Read(c.bc())
	case 0x12: // LD (DE),A
		c.memory.Write(c.de(), c.A)
	case 0x1A: // LD A,(DE)
		c.A = c.memory.Read(c.de())
	}
	return nil
}

// ldHLPlusA stores A at (HL), then increments HL. LD (HL+),A.
func ldHLPlusA(c *CPU) error {
	addr := c.hl()
	c.memory.Write(addr, c.A)
	c.setHL(addr + 1)
	return nil
}

// ldAHLPlus loads A from (HL), then increments HL. LD A,(HL+).
func ldAHLPlus(c *CPU) error {
	addr := c.hl()
	c.A = c.memory.Read(addr)
	c.setHL(addr + 1)
	return nil
}

// ldHLMinusA stores A at (HL), then decrements HL. LD (HL-),A.
func ldHLMinusA(c *CPU) error {
	addr := c.hl()
	c.memory.Write(addr, c.A)
	c.setHL(addr - 1)
	return nil
}

// ldAHLMinus loads A from (HL), then decrements HL. LD A,(HL-).
func ldAHLMinus(c *CPU) error {
	addr := c.hl()
	c.A = c.memory.Read(addr)
	c.setHL(addr - 1)
	return nil
}

// ldAddrSP stores SP at address nn (little-endian). LD (nn),SP.
func ldAddrSP(c *CPU, params ...any) error {
	addr := uint16(params[0].(Extended))
	c.memory.Write(addr, uint8(c.SP))
	c.memory.Write(addr+1, uint8(c.SP>>8))
	return nil
}

// ldIndirectImm stores an immediate byte at (HL). LD (HL),n.
func ldIndirectImm(c *CPU, params ...any) error {
	value := uint8(params[0].(Immediate8))
	c.memory.Write(c.hl(), value)
	return nil
}

// ldSPHL copies HL to SP. LD SP,HL.
func ldSPHL(c *CPU) error {
	c.SP = c.hl()
	return nil
}

// ldAddrA stores A at address nn. LD (nn),A.
func ldAddrA(c *CPU, params ...any) error {
	addr := uint16(params[0].(Extended))
	c.memory.Write(addr, c.A)
	return nil
}

// ldAAddr loads A from address nn. LD A,(nn).
func ldAAddr(c *CPU, params ...any) error {
	addr := uint16(params[0].(Extended))
	c.A = c.memory.Read(addr)
	return nil
}

// ldHLSPOffset loads HL with SP + signed 8-bit offset. LD HL,SP+e.
// Z=0, N=0, H and C computed on the low byte of SP + unsigned offset.
func ldHLSPOffset(c *CPU, params ...any) error {
	offset := int8(params[0].(Immediate8))

	c.setH((c.SP&0x0F)+(uint16(uint8(offset))&0x0F) > 0x0F)
	c.setC((c.SP&0xFF)+uint16(uint8(offset)) > 0xFF)

	c.setHL(uint16(int32(c.SP) + int32(offset)))

	c.Flags.Z = 0
	c.setN(false)
	return nil
}

// ldhNA stores A at $FF00+n. LDH (n),A.
func ldhNA(c *CPU, params ...any) error {
	addr := 0xFF00 + uint16(params[0].(Immediate8))
	c.memory.Write(addr, c.A)
	return nil
}

// ldhAN loads A from $FF00+n. LDH A,(n).
func ldhAN(c *CPU, params ...any) error {
	addr := 0xFF00 + uint16(params[0].(Immediate8))
	c.A = c.memory.Read(addr)
	return nil
}

// ldCA stores A at $FF00+C. LD (C),A.
func ldCA(c *CPU) error {
	c.memory.Write(0xFF00+uint16(c.C), c.A)
	return nil
}

// ldAC loads A from $FF00+C. LD A,(C).
func ldAC(c *CPU) error {
	c.A = c.memory.Read(0xFF00 + uint16(c.C))
	return nil
}

// pushReg16 pushes a 16-bit register pair onto the stack.
// Register pair encoded in bits 4-5: 0=BC, 1=DE, 2=HL, 3=AF.
func pushReg16(c *CPU, _ ...any) error {
	switch (c.currentOpcode >> 4) & 0x03 {
	case 0:
		c.push16(c.bc())
	case 1:
		c.push16(c.de())
	case 2:
		c.push16(c.hl())
	case 3:
		c.push16(c.af())
	}
	return nil
}

// popReg16 pops a 16-bit register pair from the stack.
// Register pair encoded in bits 4-5: 0=BC, 1=DE, 2=HL, 3=AF.
// POP AF masks the lower 4 bits of F (always 0 on SM83).
func popReg16(c *CPU, _ ...any) error {
	value := c.pop16()
	switch (c.currentOpcode >> 4) & 0x03 {
	case 0:
		c.setBC(value)
	case 1:
		c.setDE(value)
	case 2:
		c.setHL(value)
	case 3:
		c.setAF(value & 0xFFF0) // Lower 4 bits of F always 0
	}
	return nil
}
