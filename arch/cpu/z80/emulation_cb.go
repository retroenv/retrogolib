package z80

// CB prefix instruction implementations - Bit operations

// getCBTiming calculates timing for CB-prefixed instructions based on operation and register.
func getCBTiming(opcodeByte, reg uint8) byte {
	switch {
	case opcodeByte <= 0x3F && reg == 6: // Rotate/shift (HL)
		return 15
	case opcodeByte <= 0x7F && reg == 6: // BIT n,(HL)
		return 12
	case opcodeByte >= 0x80 && reg == 6: // RES/SET n,(HL)
		return 15
	default:
		return 8
	}
}

// cbRlc implements CB 00-07: RLC r.
func cbRlc(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // RLC (HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.rlc(value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.rlc(value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}

// cbRrc implements CB 08-0F: RRC r.
func cbRrc(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // RRC (HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.rrc(value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.rrc(value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}

// cbRl implements CB 10-17: RL r.
func cbRl(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // RL (HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.rl(value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.rl(value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}

// cbRr implements CB 18-1F: RR r.
func cbRr(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // RR (HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.rr(value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.rr(value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}

// cbSla implements CB 20-27: SLA r.
func cbSla(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // SLA (HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.sla(value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.sla(value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}

// cbSra implements CB 28-2F: SRA r.
func cbSra(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // SRA (HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.sra(value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.sra(value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}

// cbSll implements CB 30-37: SLL r (undocumented shift left logical).
func cbSll(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // SLL (HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.sll(value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.sll(value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}

// cbSrl implements CB 38-3F: SRL r.
func cbSrl(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	reg := opcodeByte & 0x07

	if reg == 6 { // SRL (HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.srl(value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.srl(value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}

// cbBit implements CB 40-7F: BIT n,r.
func cbBit(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	bit := (opcodeByte >> 3) & 0x07
	reg := opcodeByte & 0x07

	var value uint8
	if reg == 6 { // BIT n,(HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value = c.memory.Read(addr)
	} else {
		value = c.GetRegisterValue(reg)
	}

	c.bit(bit, value)
	return nil
}

// cbRes implements CB 80-BF: RES n,r.
func cbRes(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	bit := (opcodeByte >> 3) & 0x07
	reg := opcodeByte & 0x07

	if reg == 6 { // RES n,(HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.res(bit, value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.res(bit, value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}

// cbSet implements CB C0-FF: SET n,r.
func cbSet(c *CPU, params ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	bit := (opcodeByte >> 3) & 0x07
	reg := opcodeByte & 0x07

	if reg == 6 { // SET n,(HL)
		addr := uint16(c.H)<<8 | uint16(c.L)
		value := c.memory.Read(addr)
		result := c.setBit(bit, value)
		c.memory.Write(addr, result)
	} else {
		value := c.GetRegisterValue(reg)
		result := c.setBit(bit, value)
		c.SetRegisterValue(reg, result)
	}
	return nil
}
