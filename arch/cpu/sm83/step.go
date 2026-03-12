package sm83

import (
	"fmt"

	"github.com/retroenv/retrogolib/set"
)

// TraceStep contains all info needed to print a trace step.
type TraceStep struct {
	PC             uint16 // program counter
	OpcodeOperands []byte // instruction opcode and operand bytes
	Opcode         Opcode

	CustomData string // custom data field that can be used in the pre execution hook
}

// Step executes the next instruction in the CPU.
func (c *CPU) Step() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Handle interrupts first
	c.HandleInterrupts()

	if c.halted {
		// CPU is halted, just advance cycles
		c.cycles++
		return nil
	}

	// Capture pending IME enable from previous EI instruction.
	// EI enables interrupts after the NEXT instruction executes,
	// so the instruction immediately after EI runs with interrupts still disabled.
	pendingIME := c.imeDelay
	if pendingIME {
		c.imeDelay = false
	}

	opcode, opcodeByte, err := c.decodeNextInstruction()
	if err != nil {
		return err
	}
	oldPC := c.PC

	c.cycles += uint64(opcode.Timing)
	c.currentOpcode = opcodeByte

	if err := c.executeInstruction(opcode, opcodeByte, oldPC); err != nil {
		return err
	}

	// Enable IME after instruction if EI was the previous instruction
	if pendingIME {
		c.ime = true
	}

	// Handle HALT bug: if HALT bug is active, PC doesn't increment
	if c.haltBug {
		c.haltBug = false
		c.PC = oldPC
	}

	return nil
}

// executeInstruction runs the decoded instruction and updates the program counter.
func (c *CPU) executeInstruction(opcode Opcode, opcodeByte byte, oldPC uint16) error {
	ins := opcode.Instruction
	if ins.NoParamFunc != nil {
		if c.opts.preExecutionHook != nil {
			c.opts.preExecutionHook(c, opcodeByte)
		}
		if err := ins.NoParamFunc(c); err != nil {
			return fmt.Errorf("executing no param instruction %s: %w", ins.Name, err)
		}
		c.updatePC(ins, oldPC, int(opcode.Size))
		return nil
	}

	params, operands, err := readOpParams(c, opcode.Addressing)
	if err != nil {
		return fmt.Errorf("reading opcode params: %w", err)
	}
	if c.opts.tracing {
		c.TraceStep.OpcodeOperands = append(c.TraceStep.OpcodeOperands, operands...)
	}
	if c.opts.preExecutionHook != nil {
		c.opts.preExecutionHook(c, opcodeByte, params...)
	}

	if err := ins.ParamFunc(c, params...); err != nil {
		return fmt.Errorf("executing param instruction %s: %w", ins.Name, err)
	}
	c.updatePC(ins, oldPC, int(opcode.Size))
	return nil
}

// decodeNextInstruction decodes the current instruction at the program counter.
func (c *CPU) decodeNextInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC)

	if opcodeByte == PrefixCB {
		return c.decodeCBInstruction()
	}

	opcode := Opcodes[opcodeByte]
	if opcode.Instruction == nil {
		return Opcode{}, opcodeByte, fmt.Errorf("%w: 0x%02X", ErrIllegalOpcode, opcodeByte)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{opcodeByte},
		}
	}
	return opcode, opcodeByte, nil
}

// decodeCBInstruction decodes CB-prefixed instructions (bit operations).
func (c *CPU) decodeCBInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1)

	opcode := CBOpcodes[opcodeByte]
	if opcode.Instruction == nil {
		return Opcode{}, PrefixCB, fmt.Errorf("%w: CB %02X", ErrUnsupportedOpcode, opcodeByte)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixCB, opcodeByte},
		}
	}

	return opcode, PrefixCB, nil
}

// updatePC updates the program counter based on the instruction execution.
func (c *CPU) updatePC(ins *Instruction, oldPC uint16, amount int) {
	// Check if this is a jump instruction that always changes PC
	if ins != nil && isJumpInstruction(ins) {
		// Jump instructions handle PC themselves, don't modify it
		return
	}

	// Update PC only if the instruction execution did not change it
	if oldPC == c.PC {
		// PC unchanged, advance by instruction size
		c.PC += uint16(amount)
	}

	// PC was changed by the instruction (e.g., conditional jump taken), don't modify it further
}

// jumpInstructions is a lookup set of instructions that always modify PC.
var jumpInstructions = set.Set[*Instruction]{
	Call:     {},
	CallCond: {},
	JpAbs:    {},
	JpCond:   {},
	JpHL:     {},
	JrRel:    {},
	JrCond:   {},
	Ret:      {},
	RetCond:  {},
	Reti:     {},
	Rst:      {},
}

// isJumpInstruction checks if an instruction is a jump/branch instruction that always modifies PC.
func isJumpInstruction(ins *Instruction) bool {
	return ins != nil && jumpInstructions.Contains(ins)
}
