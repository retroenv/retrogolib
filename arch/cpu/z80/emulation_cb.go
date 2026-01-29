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
func cbRlc(c *CPU, _ ...any) error {
	c.applyCBOperation(c.rlc)
	return nil
}

// cbRrc implements CB 08-0F: RRC r.
func cbRrc(c *CPU, _ ...any) error {
	c.applyCBOperation(c.rrc)
	return nil
}

// cbRl implements CB 10-17: RL r.
func cbRl(c *CPU, _ ...any) error {
	c.applyCBOperation(c.rl)
	return nil
}

// cbRr implements CB 18-1F: RR r.
func cbRr(c *CPU, _ ...any) error {
	c.applyCBOperation(c.rr)
	return nil
}

// cbSla implements CB 20-27: SLA r.
func cbSla(c *CPU, _ ...any) error {
	c.applyCBOperation(c.sla)
	return nil
}

// cbSra implements CB 28-2F: SRA r.
func cbSra(c *CPU, _ ...any) error {
	c.applyCBOperation(c.sra)
	return nil
}

// cbSll implements CB 30-37: SLL r (undocumented shift left logical).
func cbSll(c *CPU, _ ...any) error {
	c.applyCBOperation(c.sll)
	return nil
}

// cbSrl implements CB 38-3F: SRL r.
func cbSrl(c *CPU, _ ...any) error {
	c.applyCBOperation(c.srl)
	return nil
}

// cbBit implements CB 40-7F: BIT n,r.
func cbBit(c *CPU, _ ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	bit := (opcodeByte >> 3) & 0x07
	reg := opcodeByte & 0x07

	var value uint8
	if reg == 6 { // BIT n,(HL)
		addr := c.hl()
		value = c.memory.Read(addr)
	} else {
		value = c.GetRegisterValue(reg)
	}

	c.bit(bit, value)
	return nil
}

// cbRes implements CB 80-BF: RES n,r.
func cbRes(c *CPU, _ ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	bit := (opcodeByte >> 3) & 0x07
	c.applyCBOperation(func(value uint8) uint8 {
		return c.res(bit, value)
	})
	return nil
}

// cbSet implements CB C0-FF: SET n,r.
func cbSet(c *CPU, _ ...any) error {
	opcodeByte := c.memory.Read(c.PC + 1)
	bit := (opcodeByte >> 3) & 0x07
	c.applyCBOperation(func(value uint8) uint8 {
		return c.setBit(bit, value)
	})
	return nil
}
