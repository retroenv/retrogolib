package sm83

// nop does nothing.
func nop(_ *CPU) error {
	return nil
}

// halt halts the CPU.
func halt(c *CPU) error {
	c.halted = true
	return nil
}

// stop enters low-power standby mode.
func stop(c *CPU) error {
	c.halted = true
	return nil
}

// getALUOperand returns the operand for ALU operations.
// If an Immediate8 parameter is provided, it returns that value.
// Otherwise, it reads the register encoded in the lower 3 bits of the current opcode.
func (c *CPU) getALUOperand(params []any) uint8 {
	if len(params) > 0 {
		if v, ok := params[0].(Immediate8); ok {
			return uint8(v)
		}
	}
	return c.GetRegisterValue(c.currentOpcode & 0x07)
}

// addA performs ADD A,r or ADD A,n.
func addA(c *CPU, params ...any) error {
	operand := c.getALUOperand(params)
	result16 := uint16(c.A) + uint16(operand)
	result := uint8(result16)

	c.setZ(result)
	c.setN(false)
	c.setH((c.A&0x0F)+(operand&0x0F) > 0x0F)
	c.setC(result16 > 0xFF)
	c.A = result
	return nil
}

// adcA performs ADC A,r or ADC A,n (add with carry).
func adcA(c *CPU, params ...any) error {
	operand := c.getALUOperand(params)
	carry := c.Flags.C
	result16 := uint16(c.A) + uint16(operand) + uint16(carry)
	result := uint8(result16)

	c.setZ(result)
	c.setN(false)
	c.setH((c.A&0x0F)+(operand&0x0F)+carry > 0x0F)
	c.setC(result16 > 0xFF)
	c.A = result
	return nil
}

// subA performs SUB r or SUB n.
func subA(c *CPU, params ...any) error {
	operand := c.getALUOperand(params)
	result := c.A - operand

	c.setZ(result)
	c.setN(true)
	c.setH((c.A & 0x0F) < (operand & 0x0F))
	c.setC(c.A < operand)
	c.A = result
	return nil
}

// sbcA performs SBC A,r or SBC A,n (subtract with carry).
func sbcA(c *CPU, params ...any) error {
	operand := c.getALUOperand(params)
	carry := c.Flags.C
	result16 := uint16(c.A) - uint16(operand) - uint16(carry)
	result := uint8(result16)

	c.setZ(result)
	c.setN(true)
	c.setH((c.A & 0x0F) < (operand&0x0F)+carry)
	c.setC(result16 > 0xFF)
	c.A = result
	return nil
}

// andA performs AND r or AND n.
func andA(c *CPU, params ...any) error {
	operand := c.getALUOperand(params)
	c.A &= operand

	c.setZ(c.A)
	c.setN(false)
	c.setH(true)
	c.setC(false)
	return nil
}

// orA performs OR r or OR n.
func orA(c *CPU, params ...any) error {
	operand := c.getALUOperand(params)
	c.A |= operand

	c.setZ(c.A)
	c.setN(false)
	c.setH(false)
	c.setC(false)
	return nil
}

// xorA performs XOR r or XOR n.
func xorA(c *CPU, params ...any) error {
	operand := c.getALUOperand(params)
	c.A ^= operand

	c.setZ(c.A)
	c.setN(false)
	c.setH(false)
	c.setC(false)
	return nil
}

// cpA performs CP r or CP n (compare, like SUB but doesn't store result).
func cpA(c *CPU, params ...any) error {
	operand := c.getALUOperand(params)
	result := c.A - operand

	c.setZ(result)
	c.setN(true)
	c.setH((c.A & 0x0F) < (operand & 0x0F))
	c.setC(c.A < operand)
	return nil
}

// incReg8 increments an 8-bit register (INC r).
// Register encoded in bits 3-5 of the opcode.
func incReg8(c *CPU, _ ...any) error {
	reg := (c.currentOpcode >> 3) & 0x07
	value := c.GetRegisterValue(reg)
	result := value + 1

	c.setZ(result)
	c.setN(false)
	c.setH((value & 0x0F) == 0x0F)
	c.SetRegisterValue(reg, result)
	return nil
}

// decReg8 decrements an 8-bit register (DEC r).
// Register encoded in bits 3-5 of the opcode.
func decReg8(c *CPU, _ ...any) error {
	reg := (c.currentOpcode >> 3) & 0x07
	value := c.GetRegisterValue(reg)
	result := value - 1

	c.setZ(result)
	c.setN(true)
	c.setH((value & 0x0F) == 0x00)
	c.SetRegisterValue(reg, result)
	return nil
}

// incReg16 increments a 16-bit register pair (INC rr).
// Register pair encoded in bits 4-5 of the opcode. No flags affected.
func incReg16(c *CPU, _ ...any) error {
	switch (c.currentOpcode >> 4) & 0x03 {
	case 0:
		c.setBC(c.bc() + 1)
	case 1:
		c.setDE(c.de() + 1)
	case 2:
		c.setHL(c.hl() + 1)
	case 3:
		c.SP++
	}
	return nil
}

// decReg16 decrements a 16-bit register pair (DEC rr).
// Register pair encoded in bits 4-5 of the opcode. No flags affected.
func decReg16(c *CPU, _ ...any) error {
	switch (c.currentOpcode >> 4) & 0x03 {
	case 0:
		c.setBC(c.bc() - 1)
	case 1:
		c.setDE(c.de() - 1)
	case 2:
		c.setHL(c.hl() - 1)
	case 3:
		c.SP--
	}
	return nil
}

// incIndirect increments the byte at (HL).
func incIndirect(c *CPU, _ ...any) error {
	addr := c.hl()
	value := c.memory.Read(addr)
	result := value + 1

	c.setZ(result)
	c.setN(false)
	c.setH((value & 0x0F) == 0x0F)
	c.memory.Write(addr, result)
	return nil
}

// decIndirect decrements the byte at (HL).
func decIndirect(c *CPU, _ ...any) error {
	addr := c.hl()
	value := c.memory.Read(addr)
	result := value - 1

	c.setZ(result)
	c.setN(true)
	c.setH((value & 0x0F) == 0x00)
	c.memory.Write(addr, result)
	return nil
}

// addHL performs ADD HL,rr.
// Register pair encoded in bits 4-5 of the opcode.
// Z flag unchanged, N=0, H from bit 11, C from bit 15.
func addHL(c *CPU, _ ...any) error {
	hl := c.hl()

	var value uint16
	switch (c.currentOpcode >> 4) & 0x03 {
	case 0:
		value = c.bc()
	case 1:
		value = c.de()
	case 2:
		value = hl
	case 3:
		value = c.SP
	}

	result32 := uint32(hl) + uint32(value)
	result := uint16(result32)

	c.setN(false)
	c.setH((hl&0x0FFF)+(value&0x0FFF) > 0x0FFF)
	c.setC(result32 > 0xFFFF)
	c.setHL(result)
	return nil
}

// rlcaFunc performs RLCA (rotate A left circular).
// Z=0, N=0, H=0, C=old bit 7.
func rlcaFunc(c *CPU) error {
	carry := (c.A & 0x80) >> 7
	c.A = (c.A << 1) | carry

	c.Flags.Z = 0
	c.setN(false)
	c.setH(false)
	c.setC(carry != 0)
	return nil
}

// rrcaFunc performs RRCA (rotate A right circular).
// Z=0, N=0, H=0, C=old bit 0.
func rrcaFunc(c *CPU) error {
	carry := c.A & 0x01
	c.A = (c.A >> 1) | (carry << 7)

	c.Flags.Z = 0
	c.setN(false)
	c.setH(false)
	c.setC(carry != 0)
	return nil
}

// rlaFunc performs RLA (rotate A left through carry).
// Z=0, N=0, H=0, C=old bit 7.
func rlaFunc(c *CPU) error {
	newCarry := c.A >> 7
	c.A = (c.A << 1) | c.Flags.C

	c.Flags.Z = 0
	c.setN(false)
	c.setH(false)
	c.setC(newCarry != 0)
	return nil
}

// rraFunc performs RRA (rotate A right through carry).
// Z=0, N=0, H=0, C=old bit 0.
func rraFunc(c *CPU) error {
	newCarry := c.A & 0x01
	c.A = (c.A >> 1) | (c.Flags.C << 7)

	c.Flags.Z = 0
	c.setN(false)
	c.setH(false)
	c.setC(newCarry != 0)
	return nil
}

// daa performs Decimal Adjust Accumulator.
// SM83 DAA adjusts A based on the previous arithmetic operation (N flag).
func daa(c *CPU) error {
	correction := uint8(0)
	carry := c.Flags.C != 0

	if c.Flags.N == 0 {
		carry = daaAfterAdd(c, &correction, carry)
	} else {
		daaAfterSub(c, &correction, carry)
	}

	c.setZ(c.A)
	c.setH(false)
	c.setC(carry)
	// N flag unchanged
	return nil
}

// daaAfterAdd applies DAA correction after ADD/ADC operations.
func daaAfterAdd(c *CPU, correction *uint8, carry bool) bool {
	if carry || c.A > 0x99 {
		*correction |= 0x60
		carry = true
	}
	if c.Flags.H != 0 || (c.A&0x0F) > 9 {
		*correction |= 0x06
	}
	c.A += *correction
	return carry
}

// daaAfterSub applies DAA correction after SUB/SBC operations.
func daaAfterSub(c *CPU, correction *uint8, carry bool) {
	if carry {
		*correction |= 0x60
	}
	if c.Flags.H != 0 {
		*correction |= 0x06
	}
	c.A -= *correction
}

// cpl complements the accumulator (A = ~A).
func cpl(c *CPU) error {
	c.A = ^c.A
	c.setN(true)
	c.setH(true)
	return nil
}

// scf sets the carry flag.
func scf(c *CPU) error {
	c.setN(false)
	c.setH(false)
	c.setC(true)
	return nil
}

// ccf complements the carry flag.
func ccf(c *CPU) error {
	c.setN(false)
	c.setH(false)
	c.setC(c.Flags.C == 0)
	return nil
}

// di disables interrupts.
func di(c *CPU) error {
	c.ime = false
	c.imeDelay = false
	return nil
}

// ei enables interrupts (delayed by one instruction).
func ei(c *CPU) error {
	c.imeDelay = true
	return nil
}

// addSPE performs ADD SP,e (SP += signed 8-bit offset).
// Z=0, N=0, H and C computed on the low byte of SP + unsigned offset.
func addSPE(c *CPU, params ...any) error {
	offset := int8(params[0].(Immediate8))

	c.setH((c.SP&0x0F)+(uint16(uint8(offset))&0x0F) > 0x0F)
	c.setC((c.SP&0xFF)+uint16(uint8(offset)) > 0xFF)

	c.SP = uint16(int32(c.SP) + int32(offset))

	c.Flags.Z = 0
	c.setN(false)
	return nil
}
