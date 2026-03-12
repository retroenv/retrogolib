package sm83

// checkCondition evaluates the condition encoded in bits 3-4 of the current opcode.
// Conditions: 0=NZ, 1=Z, 2=NC, 3=C.
func (c *CPU) checkCondition() bool {
	cond := (c.currentOpcode >> 3) & 0x03
	switch cond {
	case 0: // NZ
		return c.Flags.Z == 0
	case 1: // Z
		return c.Flags.Z == 1
	case 2: // NC
		return c.Flags.C == 0
	case 3: // C
		return c.Flags.C == 1
	}
	return false
}

// jpAbs performs an unconditional absolute jump (JP nn).
func jpAbs(c *CPU, params ...any) error {
	c.PC = uint16(params[0].(Extended))
	return nil
}

// jpCond performs a conditional absolute jump (JP cc,nn).
// If the condition is met, jumps and adds 1 extra cycle.
// If not taken, step.go advances PC by instruction size.
func jpCond(c *CPU, params ...any) error {
	addr := uint16(params[0].(Extended))
	if c.checkCondition() {
		c.PC = addr
		c.cycles++
	}
	return nil
}

// jpHL performs an unconditional jump to the address in HL (JP (HL)).
func jpHL(c *CPU, _ ...any) error {
	c.PC = c.hl()
	return nil
}

// jrRel performs an unconditional relative jump (JR e).
// Target = PC + 2 + signed offset.
func jrRel(c *CPU, params ...any) error {
	offset := int8(params[0].(Relative))
	c.PC = uint16(int32(c.PC) + 2 + int32(offset))
	return nil
}

// jrCond performs a conditional relative jump (JR cc,e).
// If the condition is met, jumps and adds 1 extra cycle.
// If not taken, step.go advances PC by instruction size.
func jrCond(c *CPU, params ...any) error {
	offset := int8(params[0].(Relative))
	if c.checkCondition() {
		c.PC = uint16(int32(c.PC) + 2 + int32(offset))
		c.cycles++
	}
	return nil
}

// call performs an unconditional call (CALL nn).
// Pushes the return address (PC+3) and jumps to the target.
func call(c *CPU, params ...any) error {
	addr := uint16(params[0].(Extended))
	c.push16(c.PC + 3)
	c.PC = addr
	return nil
}

// callCond performs a conditional call (CALL cc,nn).
// If the condition is met, pushes PC+3, jumps, and adds 3 extra cycles.
// If not taken, step.go advances PC by instruction size.
func callCond(c *CPU, params ...any) error {
	addr := uint16(params[0].(Extended))
	if c.checkCondition() {
		c.push16(c.PC + 3)
		c.PC = addr
		c.cycles += 3
	}
	return nil
}

// ret performs an unconditional return (RET).
func ret(c *CPU) error {
	c.PC = c.pop16()
	return nil
}

// retCond performs a conditional return (RET cc).
// If the condition is met, pops PC and adds 3 extra cycles.
// If not taken, step.go advances PC by instruction size.
func retCond(c *CPU) error {
	if c.checkCondition() {
		c.PC = c.pop16()
		c.cycles += 3
	}
	return nil
}

// reti performs a return and enables interrupts (RETI).
func reti(c *CPU) error {
	c.PC = c.pop16()
	c.ime = true
	return nil
}

// rst performs a restart (RST n).
// Pushes PC+1 and jumps to the vector encoded in bits 3-5 of the opcode.
func rst(c *CPU, _ ...any) error {
	c.push16(c.PC + 1)
	vector := uint16(c.currentOpcode & 0x38)
	c.PC = vector
	return nil
}
