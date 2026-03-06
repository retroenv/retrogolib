package m65816

import "fmt"

// Step executes the next instruction and returns any error.
func (c *CPU) Step() error {
	if c.stopped {
		return nil // STP: halted until RESET
	}
	if c.waiting {
		return nil // WAI: waiting for interrupt
	}

	// Decode opcode at current PC
	opByte := c.memory.ReadByte(c.FullPC())
	op, ok := GetOpcodeInfo(opByte)
	if !ok {
		return fmt.Errorf("%w: 0x%02x at PB=%02X PC=%04X", ErrInvalidOpcode, opByte, c.PB, c.PC)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			PB:             c.PB,
			OpcodeOperands: []byte{opByte},
			Opcode:         op,
		}
	}

	c.cycles += uint64(op.Timing)

	ins := op.Instruction

	// -- No-param instructions --
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
			c.PC += uint16(c.instrSize(op))
		}
		return nil
	}

	// -- Param instructions --
	params, operands, pageCrossed, err := readOpParams(c, op.Addressing, op)
	if err != nil {
		return fmt.Errorf("reading params for %s: %w", ins.Name, err)
	}
	if c.opts.tracing {
		c.TraceStep.OpcodeOperands = append(c.TraceStep.OpcodeOperands, operands...)
		c.TraceStep.PageCrossed = pageCrossed
	}
	if pageCrossed && op.PageCrossCycle {
		// Branch (relative) page-crossing penalty only applies in emulation mode.
		// For all other addressing modes (abs,X etc.) the penalty applies always.
		if op.Addressing != RelativeAddressing || c.E {
			c.cycles++
		}
	}
	if c.opts.preExecutionHook != nil {
		c.opts.preExecutionHook(c, ins, params...)
	}

	instrLen := 1 + len(operands)

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
