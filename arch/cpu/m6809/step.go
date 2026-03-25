package m6809

import "fmt"

// Step executes the next instruction and returns any error.
func (c *CPU) Step() error {
	if c.waiting {
		return nil // SYNC: waiting for interrupt
	}

	// Check for prefix bytes
	opByte := c.memory.Read(c.PC)

	switch opByte {
	case Prefix10:
		return c.stepPage2()
	case Prefix11:
		return c.stepPage3()
	default:
		return c.stepBase(opByte)
	}
}

// stepBase executes a base page instruction.
func (c *CPU) stepBase(opByte uint8) error {
	op, ok := GetOpcodeInfo(opByte)
	if !ok {
		return fmt.Errorf("%w: 0x%02x at PC=%04X", ErrInvalidOpcode, opByte, c.PC)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			OpcodeOperands: []byte{opByte},
			Opcode:         op,
		}
	}

	c.cycles += uint64(op.Timing)
	return c.executeInstruction(op, 1)
}

// stepPage2 executes a page 2 ($10 prefix) instruction.
func (c *CPU) stepPage2() error {
	opByte := c.memory.Read(c.PC + 1)
	op, ok := GetPage2OpcodeInfo(opByte)
	if !ok {
		return fmt.Errorf("%w: 0x10 0x%02x at PC=%04X", ErrInvalidOpcode, opByte, c.PC)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			OpcodeOperands: []byte{Prefix10, opByte},
			Opcode:         op,
		}
	}

	c.cycles += uint64(op.Timing)
	return c.executeInstruction(op, 2)
}

// stepPage3 executes a page 3 ($11 prefix) instruction.
func (c *CPU) stepPage3() error {
	opByte := c.memory.Read(c.PC + 1)
	op, ok := GetPage3OpcodeInfo(opByte)
	if !ok {
		return fmt.Errorf("%w: 0x11 0x%02x at PC=%04X", ErrInvalidOpcode, opByte, c.PC)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			OpcodeOperands: []byte{Prefix11, opByte},
			Opcode:         op,
		}
	}

	c.cycles += uint64(op.Timing)
	return c.executeInstruction(op, 2)
}

// executeInstruction dispatches execution to the instruction handler.
// baseOffset is the offset from PC to the first operand byte (1 for base page, 2 for prefixed).
func (c *CPU) executeInstruction(op Opcode, baseOffset uint16) error {
	ins := op.Instruction

	// No-param instructions
	if ins.NoParamFunc != nil {
		if c.opts.preExecutionHook != nil {
			c.opts.preExecutionHook(c, ins)
		}
		c.pcChanged = false
		if err := ins.NoParamFunc(c); err != nil {
			return fmt.Errorf("executing %s: %w", ins.Name, err)
		}
		// Advance PC unless the instruction explicitly changed it
		if !c.pcChanged {
			c.PC += uint16(op.Size)
		}
		return nil
	}

	// Param instructions
	params, operands, err := readOpParams(c, op.Addressing, baseOffset)
	if err != nil {
		return fmt.Errorf("reading params for %s: %w", ins.Name, err)
	}
	if c.opts.tracing {
		c.TraceStep.OpcodeOperands = append(c.TraceStep.OpcodeOperands, operands...)
	}
	if c.opts.preExecutionHook != nil {
		c.opts.preExecutionHook(c, ins, params...)
	}

	instrLen := int(baseOffset) + len(operands)

	// Set nextPC so JSR/BSR handlers can push the correct return address.
	c.nextPC = c.PC + uint16(instrLen)
	c.pcChanged = false
	if err := ins.ParamFunc(c, params...); err != nil {
		return fmt.Errorf("executing %s: %w", ins.Name, err)
	}

	// Advance PC unless the instruction explicitly changed it (branch/jump)
	if !c.pcChanged {
		c.PC += uint16(instrLen)
	}
	return nil
}
