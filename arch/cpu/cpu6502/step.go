package cpu6502

import (
	"fmt"
)

// TraceStep contains all info needed to print a trace step.
// It is populated only for decoded instructions when tracing is enabled. It does
// not describe hardware interrupt entry; StackEvent carries its own context so
// stack hooks work independently of tracing.
type TraceStep struct {
	// PC is the instruction address captured during decoding.
	PC uint16
	// OpcodeOperands contains the raw opcode byte followed by its operand bytes.
	OpcodeOperands []byte
	// Opcode contains the decoded instruction and addressing metadata.
	Opcode Opcode

	// CustomData contains caller-defined trace data set by the pre-execution hook.
	CustomData string
	// PageCrossed reports whether operand address resolution crossed a page.
	PageCrossed bool
}

// Step consumes one stall cycle or executes the next instruction. It does not
// service interrupts; call CheckInterrupts at an instruction boundary.
func (c *CPU) Step() error {
	if c.opts.cycleHook != nil {
		return c.stepCycles()
	}
	return c.stepInstruction()
}

func (c *CPU) stepInstruction() error {
	if c.consumeStallCycle() {
		c.cycles++
		return nil
	}

	c.branchTaken = false
	oldPC := c.PC
	opcode, err := c.decodeNextInstruction()
	if err != nil {
		return err
	}

	c.cycles += uint64(opcode.Timing)

	ins := opcode.Instruction
	if ins == BrkInst {
		c.interrupt = InterruptBRK
	}
	if ins.NoParamFunc != nil {
		if c.opts.tracing {
			c.TraceStep.PageCrossed = false
		}
		if c.opts.preExecutionHook != nil {
			c.opts.preExecutionHook(c, ins)
		}

		if err := ins.NoParamFunc(c); err != nil {
			return fmt.Errorf("executing no param instruction %s: %w", ins.Name, err)
		}

		// Determine instruction size from the opcode table's addressing mode,
		// not the instruction's own addressing map (which may differ for shared NOPs).
		size := addressingModeSize(opcode.Addressing)
		c.updatePC(ins, oldPC, size)
		return nil
	}

	params, operands, pageCrossed, err := readOpParams(c, opcode.Addressing)
	if err != nil {
		return fmt.Errorf("reading opcode params: %w", err)
	}
	if c.opts.tracing {
		c.TraceStep.OpcodeOperands = append(c.TraceStep.OpcodeOperands, operands...)
		c.TraceStep.PageCrossed = pageCrossed
	}
	if c.opts.preExecutionHook != nil {
		c.opts.preExecutionHook(c, ins, params...)
	}

	if pageCrossed && opcode.PageCrossCycle {
		c.cycles++
	}

	opcodeLen := len(operands) + 1

	if err := ins.ParamFunc(c, params...); err != nil {
		return fmt.Errorf("executing param instruction %s: %w", ins.Name, err)
	}
	c.updatePC(ins, oldPC, opcodeLen)
	return nil
}

// decodeNextInstruction decodes the current instruction at the program counter.
func (c *CPU) decodeNextInstruction() (Opcode, error) {
	b := c.memory.Read(c.PC)
	c.executionPC = c.PC
	c.executionCycle = c.cycles
	c.opcode = b
	c.interrupt = InterruptNone

	var opcode Opcode
	switch {
	case c.opts.variant == VariantSynertek65C02:
		opcode = OpcodesSynertek65C02[b]
	case c.opts.variant >= Variant65C02:
		opcode = Opcodes65C02[b]
	default:
		opcode = Opcodes[b]
	}
	if opcode.Instruction == nil {
		return Opcode{}, fmt.Errorf("%w: 0x%02x at PC=0x%04x", ErrUnknownOpcode, b, c.PC)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{b},
		}
	}
	return opcode, nil
}

// updatePC updates the program counter based on the instruction execution.
func (c *CPU) updatePC(ins *Instruction, oldPC uint16, amount int) {
	if c.branchTaken {
		// Taken branches compare the target with the following instruction,
		// including self-loops and address-space wraparound.
		nextAddress := oldPC + uint16(amount)
		if c.PC&0xff00 != nextAddress&0xff00 {
			c.cycles++
		}
		return
	}

	// update PC only if the instruction execution did not change it
	if oldPC == c.PC {
		// If the instruction explicitly sets PC (JMP, JSR, RTI, RTS, BRK) but landed on the
		// same address, don't advance further (self-referencing jump or interrupt to same addr).
		if ins.Name == JmpInst.Name || ins.Name == JsrInst.Name {
			return
		}
		if _, ok := NotExecutingFollowingOpcodeInstructions[ins.Name]; ok {
			return
		}
		c.PC += uint16(amount)
	}
}
