package m6502

import (
	"fmt"
)

// TraceStep contains all info needed to print a trace step.
type TraceStep struct {
	PC             uint16 // program counter
	OpcodeOperands []byte // instruction opcode and operand bytes
	Opcode         Opcode

	CustomData  string // custom data field that can be used in the pre execution hook
	PageCrossed bool
}

// Step executes the next instruction in the CPU.
func (c *CPU) Step() error {
	oldPC := c.PC
	opcode, err := c.decodeNextInstruction()
	if err != nil {
		return err
	}

	c.cycles += uint64(opcode.Timing)

	ins := opcode.Instruction
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

		// Get the correct instruction size from the opcode info
		size := 1
		for _, info := range ins.Addressing {
			size = int(info.Size)
			break
		}
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

	if pageCrossed && c.TraceStep.Opcode.PageCrossCycle {
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
	opcode := Opcodes[b]
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
	// update PC only if the instruction execution did not change it
	if oldPC == c.PC {
		if ins.Name == Jmp.Name {
			return // endless loop detected
		}

		c.PC += uint16(amount)
		return
	}

	// page crossing is measured based on the start of the instruction that follows the
	// current instruction
	nextAddress := oldPC + uint16(amount)
	pageCrossed := c.PC&0xff00 != nextAddress&0xff00
	if !pageCrossed {
		return
	}
	if _, ok := BranchingInstructions[ins.Name]; !ok {
		return
	}

	// account for a branch page crossing extra CPU cycle.
	if ins.Name != Jmp.Name && ins.Name != Jsr.Name {
		c.cycles++
	}
}
