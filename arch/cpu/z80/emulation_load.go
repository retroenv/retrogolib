package z80

import "fmt"

// ldImm8 loads an 8-bit immediate value into a register.
func ldImm8(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	imm, ok := params[0].(Immediate8)
	if !ok {
		return ErrInvalidParameterType
	}

	// Use the stored current opcode
	opcode := c.currentOpcode
	switch opcode {
	case 0x06: // LD B,n
		c.B = uint8(imm)
	case 0x0E: // LD C,n
		c.C = uint8(imm)
	case 0x16: // LD D,n
		c.D = uint8(imm)
	case 0x1E: // LD E,n
		c.E = uint8(imm)
	case 0x26: // LD H,n
		c.H = uint8(imm)
	case 0x2E: // LD L,n
		c.L = uint8(imm)
	case 0x3E: // LD A,n
		c.A = uint8(imm)
	default:
		// Default to A register for unknown opcodes
		c.A = uint8(imm)
	}
	return nil
}

// ldReg8 loads between 8-bit registers.
// params[0] = Register(bits 0-2) = source, params[1] = Register(bits 3-5) = destination.
func ldReg8(c *CPU, params ...any) error {
	if len(params) < 2 {
		return ErrMissingParameter
	}

	src, ok1 := params[0].(Register)
	dst, ok2 := params[1].(Register)

	if !ok1 || !ok2 {
		return ErrInvalidParameterType
	}

	value := c.GetRegisterValue(uint8(src))
	c.SetRegisterValue(uint8(dst), value)
	return nil
}

// ldReg16 loads a 16-bit immediate value into a register pair.
func ldReg16(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	imm, ok := params[0].(Immediate16)
	if !ok {
		return ErrInvalidParameterType
	}

	value := uint16(imm)
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

// ldIndirect loads between register pairs and memory locations.
func ldIndirect(c *CPU, _ ...any) error {
	// Use the stored current opcode
	opcode := c.currentOpcode

	switch opcode {
	case 0x02: // LD (BC),A - store A at (BC)
		addr := c.bc()
		c.memory.Write(addr, c.A)
		c.MEMPTR = (addr+1)&0xFF | uint16(c.A)<<8
	case 0x0A: // LD A,(BC) - load A from (BC)
		addr := c.bc()
		c.A = c.memory.Read(addr)
		c.MEMPTR = addr + 1
	case 0x12: // LD (DE),A - store A at (DE)
		addr := c.de()
		c.memory.Write(addr, c.A)
		c.MEMPTR = (addr+1)&0xFF | uint16(c.A)<<8
	case 0x1A: // LD A,(DE) - load A from (DE)
		addr := c.de()
		c.A = c.memory.Read(addr)
		c.MEMPTR = addr + 1
	default:
		return fmt.Errorf("unsupported indirect load opcode: 0x%02X", opcode)
	}
	return nil
}

// ldExtended loads using extended addressing.
func ldExtended(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}
	address, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	addr := uint16(address)
	opcode := c.currentOpcode
	switch opcode {
	case 0x22: // LD (nn),HL - store HL to memory address nn
		c.memory.WriteWord(addr, c.hl())
		c.MEMPTR = addr + 1
	case 0x2A: // LD HL,(nn) - load HL from memory address nn
		value := c.memory.ReadWord(addr)
		c.setHL(value)
		c.MEMPTR = addr + 1
	case 0x32: // LD (nn),A - store A to memory address nn
		c.memory.Write(addr, c.A)
		c.MEMPTR = (addr+1)&0xFF | uint16(c.A)<<8
	case 0x3A: // LD A,(nn) - load A from memory address nn
		c.A = c.memory.Read(addr)
		c.MEMPTR = addr + 1
	default:
		return fmt.Errorf("unsupported ldExtended opcode: 0x%02X", opcode)
	}
	return nil
}

// ldIndirectImm loads immediate value to indirect memory location.
func ldIndirectImm(c *CPU, _ ...any) error {
	// LD (HL),n - load immediate byte to memory at HL
	immediate := c.memory.Read(c.PC + 1)
	address := c.hl()
	c.memory.Write(address, immediate)
	return nil
}

// ldSp loads SP from HL.
func ldSp(c *CPU, _ ...any) error {
	c.SP = uint16(c.H)<<8 | uint16(c.L)
	return nil
}

// pushReg16 pushes 16-bit register to stack.
func pushReg16(c *CPU, _ ...any) error {
	opcode := c.currentOpcode
	var value uint16
	switch opcode {
	case 0xC5: // PUSH BC
		value = c.bc()
	case 0xD5: // PUSH DE
		value = c.de()
	case 0xE5: // PUSH HL
		value = c.hl()
	case 0xF5: // PUSH AF
		value = c.af()
	default:
		return fmt.Errorf("unsupported pushReg16 opcode: 0x%02X", opcode)
	}
	c.push16(value)
	return nil
}

// popReg16 pops 16-bit register from stack.
func popReg16(c *CPU, _ ...any) error {
	value := c.pop16()
	opcode := c.currentOpcode
	switch opcode {
	case 0xC1: // POP BC
		c.setBC(value)
	case 0xD1: // POP DE
		c.setDE(value)
	case 0xE1: // POP HL
		c.setHL(value)
	case 0xF1: // POP AF
		c.setAF(value)
	default:
		return fmt.Errorf("unsupported popReg16 opcode: 0x%02X", opcode)
	}
	return nil
}

// exAf exchanges AF with AF'.
func exAf(c *CPU) error {
	c.A, c.AltA = c.AltA, c.A
	c.Flags, c.AltFlags = c.AltFlags, c.Flags
	return nil
}

// exx exchanges register pairs.
func exx(c *CPU) error {
	c.B, c.AltB = c.AltB, c.B
	c.C, c.AltC = c.AltC, c.C
	c.D, c.AltD = c.AltD, c.D
	c.E, c.AltE = c.AltE, c.E
	c.H, c.AltH = c.AltH, c.H
	c.L, c.AltL = c.AltL, c.L
	return nil
}

// exSp exchanges top of stack with register pair.
func exSp(c *CPU, _ ...any) error {
	// EX (SP),HL - Exchange HL with word at top of stack
	low := c.memory.Read(c.SP)
	high := c.memory.Read(c.SP + 1)
	stackValue := uint16(high)<<8 | uint16(low)

	hlValue := c.hl()

	c.memory.Write(c.SP, uint8(hlValue))
	c.memory.Write(c.SP+1, uint8(hlValue>>8))

	c.setHL(stackValue)
	c.MEMPTR = stackValue

	return nil
}

// exDeHl exchanges DE with HL.
func exDeHl(c *CPU) error {
	c.D, c.H = c.H, c.D
	c.E, c.L = c.L, c.E
	return nil
}
