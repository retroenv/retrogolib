package m6502

import (
	"fmt"

	. "github.com/retroenv/retrogolib/addressing"
)

// TraceStep contains all info needed to print a trace step.
type TraceStep struct {
	PC             uint16
	Opcode         []byte
	Addressing     Mode
	Timing         byte
	PageCrossCycle bool
	PageCrossed    bool
	Unofficial     bool
	Instruction    string
}

// Step executes the next instruction in the CPU.
func (c *CPU) Step() (*Instruction, []any, error) {
	oldPC := c.PC
	opcode, err := c.decodeNextInstruction()
	if err != nil {
		return nil, nil, err
	}

	ins := opcode.Instruction
	if ins.NoParamFunc != nil {
		ins.NoParamFunc(c)
		c.updatePC(ins, oldPC, 1)
		return ins, nil, nil
	}

	params, opcodes, pageCrossed := readOpParams(c, opcode.Addressing, true)

	c.cycles += uint64(opcode.Timing)
	if pageCrossed {
		c.cycles++
	}

	opcodeLen := len(opcodes) + 1

	if c.opts.tracing {
		c.TraceStep.Opcode = append(c.TraceStep.Opcode, opcodes...)
		c.TraceStep.PageCrossed = pageCrossed
		c.TraceStep.Unofficial = ins.Unofficial
	}

	ins.ParamFunc(c, params...)
	c.updatePC(ins, oldPC, opcodeLen)
	return ins, params, nil
}

// decodeNextInstruction decodes the current instruction at the program counter.
func (c *CPU) decodeNextInstruction() (Opcode, error) {
	b := c.memory.Read(c.PC)
	opcode := Opcodes[b]
	if opcode.Instruction == nil {
		return Opcode{}, fmt.Errorf("unsupported opcode %00x", b)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         []byte{b},
			Addressing:     opcode.Addressing,
			Timing:         opcode.Timing,
			PageCrossCycle: opcode.PageCrossCycle,
			PageCrossed:    false,
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
	} else {
		// page crossing is measured based on the start of the instruction that follows the
		// current instruction
		nextAddress := oldPC + uint16(amount-1)
		pageCrossed := c.PC&0xff00 != nextAddress&0xff00
		if pageCrossed {
			c.accountBranchingPageCrossCycle(ins)
		}
	}
}

// accountBranchingPageCrossCycle accounts for a branch page crossing extra CPU cycle.
func (c *CPU) accountBranchingPageCrossCycle(ins *Instruction) {
	if _, ok := BranchingInstructions[ins.Name]; !ok {
		return
	}
	if ins.Name != Jmp.Name && ins.Name != Jsr.Name {
		c.cycles++
	}
}
