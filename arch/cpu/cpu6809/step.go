package cpu6809

import "fmt"

// Step services a pending interrupt or executes the next instruction.
func (c *CPU) Step() error {
	if c.CheckInterrupts() {
		return nil
	}
	if c.isWaiting() {
		return nil
	}

	decoded, err := c.decodeOpcode()
	if err != nil {
		return err
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			OpcodeOperands: append([]byte(nil), decoded.opcodeBytes[:decoded.baseOffset]...),
			Opcode:         decoded.opcode,
		}
	}

	c.cycles += uint64(decoded.opcode.Timing)

	return c.executeInstruction(decoded.opcode, decoded.baseOffset)
}

func (c *CPU) decodeOpcode() (decodedOpcode, error) {
	first := c.memory.Read(c.PC)
	decoded := decodedOpcode{opcodeBytes: [2]byte{first}, baseOffset: 1}

	switch first {
	case Prefix10:
		return c.decodePrefixedOpcode(decoded, GetPage2OpcodeInfo)
	case Prefix11:
		return c.decodePrefixedOpcode(decoded, GetPage3OpcodeInfo)
	default:
		op, ok := GetOpcodeInfo(first)
		if !ok {
			return decodedOpcode{}, fmt.Errorf("%w: 0x%02x at PC=%04X", ErrInvalidOpcode, first, c.PC)
		}
		decoded.opcode = op

		return decoded, nil
	}
}

func (c *CPU) decodePrefixedOpcode(
	decoded decodedOpcode,
	lookup func(uint8) (Opcode, bool),
) (decodedOpcode, error) {

	second := c.memory.Read(c.PC + 1)
	op, ok := lookup(second)
	if !ok {
		return decodedOpcode{}, fmt.Errorf(
			"%w: 0x%02x 0x%02x at PC=%04X",
			ErrInvalidOpcode,
			decoded.opcodeBytes[0],
			second,
			c.PC,
		)
	}

	decoded.opcode = op
	decoded.opcodeBytes[1] = second
	decoded.baseOffset = 2

	return decoded, nil
}

// executeInstruction dispatches execution to the instruction handler.
// baseOffset is the byte offset from PC to the operand.
func (c *CPU) executeInstruction(op Opcode, baseOffset uint16) error {
	ins := op.Instruction

	if ins.noParamFunc != nil {
		if c.opts.preExecutionHook != nil {
			c.opts.preExecutionHook(c, ins)
		}

		c.nextPC = c.PC + uint16(op.Size)
		c.pcChanged = false
		if err := ins.noParamFunc(c); err != nil {
			return fmt.Errorf("executing %s: %w", ins.Name, err)
		}
		if !c.pcChanged {
			c.PC = c.nextPC
		}

		return nil
	}

	param, operands, err := readOpParam(c, op.Addressing, baseOffset)
	if err != nil {
		return fmt.Errorf("reading params for %s: %w", ins.Name, err)
	}
	if c.opts.tracing {
		c.TraceStep.OpcodeOperands = append(c.TraceStep.OpcodeOperands, operands...)
	}
	if op.Addressing == IndexedAddressing {
		c.cycles += indexedModeCycles(operands[0])
	}
	if c.opts.preExecutionHook != nil {
		c.opts.preExecutionHook(c, ins, param)
	}

	c.nextPC = c.PC + baseOffset + uint16(len(operands))
	c.pcChanged = false
	if err := ins.paramFunc(c, param); err != nil {
		return fmt.Errorf("executing %s: %w", ins.Name, err)
	}
	if op.BranchTakenCycle && c.pcChanged {
		c.cycles++
	}
	if !c.pcChanged {
		c.PC = c.nextPC
	}

	return nil
}
