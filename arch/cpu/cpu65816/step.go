package cpu65816

import "fmt"

// Step services a pending interrupt or executes the next instruction.
// While the CPU is stopped by STP or waits in WAI, each Step consumes one
// idle cycle so that cycle-driven callers keep making progress.
func (c *CPU) Step() error {
	if c.stopped {
		c.cycles++
		return nil
	}
	if c.CheckInterrupts() {
		return nil
	}
	if c.waiting {
		c.cycles++
		return nil
	}

	opByte := c.memory.Read(c.FullPC())
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

	ins := op.Instruction
	if ins.NoParamFunc != nil {
		return c.executeNoParamInstruction(op)
	}
	return c.executeParamInstruction(op)
}

// executeNoParamInstruction handles instructions that decode their own operand bytes.
func (c *CPU) executeNoParamInstruction(op Opcode) error {
	ins := op.Instruction
	c.cycles += uint64(op.Timing) + c.extraCycles(op, false)

	if c.opts.preExecutionHook != nil {
		c.opts.preExecutionHook(c, ins)
	}

	c.pcChanged = false
	if err := ins.NoParamFunc(c); err != nil {
		return fmt.Errorf("executing %s: %w", ins.Name, err)
	}
	if !c.pcChanged {
		c.PC += uint16(c.instrSize(op))
	}
	return nil
}

// executeParamInstruction handles instructions that require decoded address/value parameters.
func (c *CPU) executeParamInstruction(op Opcode) error {
	ins := op.Instruction

	params, operands, pageCrossed, err := readOpParams(c, op.Addressing, op)
	if err != nil {
		return fmt.Errorf("reading params for %s: %w", ins.Name, err)
	}
	if c.opts.tracing {
		c.TraceStep.OpcodeOperands = append(c.TraceStep.OpcodeOperands, operands...)
		c.TraceStep.PageCrossed = pageCrossed
	}
	c.cycles += uint64(op.Timing) + c.extraCycles(op, pageCrossed)

	if c.opts.preExecutionHook != nil {
		c.opts.preExecutionHook(c, ins, params...)
	}

	instrLen := 1 + len(operands)

	c.pcChanged = false
	if err := ins.ParamFunc(c, params...); err != nil {
		return fmt.Errorf("executing %s: %w", ins.Name, err)
	}
	if !c.pcChanged {
		c.PC += uint16(instrLen)
	}
	return nil
}

// extraCycles applies the cycle rules of an opcode to the current CPU mode.
// It must run before the instruction executes because the instruction can
// change the flags that the rules depend on.
func (c *CPU) extraCycles(op Opcode, pageCrossed bool) uint64 {
	var extra uint64
	rules := op.Cycles

	if rules&CycleM != 0 && c.AccWidth() == 2 {
		extra++
	}
	if rules&CycleRMW != 0 && c.AccWidth() == 2 {
		extra += 2
	}
	if rules&CycleX != 0 && c.IdxWidth() == 2 {
		extra++
	}
	if rules&CycleDL != 0 && c.DP&0xFF != 0 {
		extra++
	}
	if rules&CycleIndex != 0 && (pageCrossed || c.IdxWidth() == 2) {
		extra++
	}
	if rules&CycleNative != 0 && !c.E {
		extra++
	}
	return extra
}

// instrSize accounts for immediate operands whose width follows M or X.
func (c *CPU) instrSize(op Opcode) int {
	size := int(op.Instruction.Addressing[op.Addressing].BaseSize)
	switch op.WidthFlag {
	case WidthM:
		if c.AccWidth() == 2 {
			size++
		}

	case WidthX:
		if c.IdxWidth() == 2 {
			size++
		}
	}
	return size
}
