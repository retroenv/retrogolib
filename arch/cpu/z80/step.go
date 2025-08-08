package z80

import (
	"fmt"
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
	if c.halted {
		// CPU is halted, just advance cycles
		c.cycles += 4
		return nil
	}

	// Handle interrupts first
	if err := c.handleInterrupts(); err != nil {
		return err
	}

	oldPC := c.PC
	opcode, opcodeByte, err := c.decodeNextInstruction()
	if err != nil {
		return err
	}

	c.cycles += uint64(opcode.Timing)

	// Store current opcode for instruction functions to access
	c.currentOpcode = opcodeByte

	// Increment refresh register
	c.R = (c.R & 0x80) | ((c.R + 1) & 0x7F)

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

	opcodeLen := int(opcode.Size)

	if err := ins.ParamFunc(c, params...); err != nil {
		return fmt.Errorf("executing param instruction %s: %w", ins.Name, err)
	}
	c.updatePC(ins, oldPC, opcodeLen)
	return nil
}

// decodeNextInstruction decodes the current instruction at the program counter.
func (c *CPU) decodeNextInstruction() (Opcode, uint8, error) {
	// Handle extended instruction prefixes first
	opcodeByte := c.memory.Read(c.PC)

	switch opcodeByte {
	case 0xCB:
		// CB-prefixed instructions (bit operations)
		return c.decodeCBInstruction()

	case 0xED:
		// ED-prefixed instructions (extended operations)
		return c.decodeEDInstruction()

	case 0xDD:
		// DD-prefixed instructions (IX operations)
		return c.decodeDDInstruction()

	case 0xFD:
		// FD-prefixed instructions (IY operations)
		return Opcode{}, opcodeByte, fmt.Errorf("unimplemented FD-prefixed instruction: 0x%02X", opcodeByte)
	}

	// Single-byte instructions
	opcode := Opcodes[opcodeByte]
	if opcode.Instruction == nil {
		return Opcode{}, opcodeByte, fmt.Errorf("%w: opcode 0x%02x", ErrUnsupportedAddressingMode, opcodeByte)
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

// updatePC updates the program counter based on the instruction execution.
func (c *CPU) updatePC(ins *Instruction, oldPC uint16, amount int) {
	// Update PC only if the instruction execution did not change it
	if oldPC == c.PC {
		if ins.Name == JpAbs.Name || ins.Name == JpCond.Name {
			return // endless loop detected
		}

		c.PC += uint16(amount)
		return
	}

	// For relative branches, we might need to account for branch timing
}

// decodeCBInstruction decodes CB-prefixed instructions (bit operations).
func (c *CPU) decodeCBInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual CB instruction

	// For CB 00 - RLC B
	if opcodeByte == 0x00 {
		opcode := Opcode{
			Instruction: &Instruction{
				Name:        "rlc",
				NoParamFunc: rlcB,
			},
			Size:   2,
			Timing: 8,
		}
		return opcode, 0xCB, nil
	}

	return Opcode{}, 0xCB, fmt.Errorf("unimplemented CB instruction: CB %02X", opcodeByte)
}

// decodeEDInstruction decodes ED-prefixed instructions (extended operations).
func (c *CPU) decodeEDInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual ED instruction

	// For ED 44 - NEG
	if opcodeByte == 0x44 {
		opcode := Opcode{
			Instruction: &Instruction{
				Name:        "neg",
				NoParamFunc: negA,
			},
			Size:   2,
			Timing: 8,
		}
		return opcode, 0xED, nil
	}

	return Opcode{}, 0xED, fmt.Errorf("unimplemented ED instruction: ED %02X", opcodeByte)
}

// decodeDDInstruction decodes DD-prefixed instructions (IX operations).
func (c *CPU) decodeDDInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual DD instruction

	// For DD 21 - LD IX,nn
	if opcodeByte == 0x21 {
		opcode := Opcode{
			Instruction: &Instruction{
				Name:      "ld",
				ParamFunc: ldIXnn,
			},
			Addressing: ImmediateAddressing,
			Size:       4,
			Timing:     14,
		}
		return opcode, 0xDD, nil
	}

	return Opcode{}, 0xDD, fmt.Errorf("unimplemented DD instruction: DD %02X", opcodeByte)
}

// handleInterrupts processes pending interrupts.
func (c *CPU) handleInterrupts() error {
	// Non-maskable interrupt has highest priority
	if c.triggerNmi {
		c.triggerNmi = false
		c.halted = false

		// Save current PC
		c.push16(c.PC)

		// Jump to NMI vector
		c.PC = 0x0066
		c.iff2 = c.iff1
		c.iff1 = false

		c.cycles += 11
		return nil
	}

	// Maskable interrupt
	if c.triggerIrq && c.iff1 {
		c.triggerIrq = false
		c.halted = false
		c.iff1 = false
		c.iff2 = false

		// Save current PC
		c.push16(c.PC)

		switch c.im {
		case 0:
			// Interrupt mode 0: Execute instruction on data bus (usually RST)
			c.PC = 0x0040
			c.cycles += 13
		case 1:
			// Interrupt mode 1: Jump to 0x0038
			c.PC = 0x0038
			c.cycles += 13
		case 2:
			// Interrupt mode 2: Vector table lookup
			vector := uint16(c.I)<<8 | uint16(c.memory.Read(0xFFFF))
			c.PC = c.memory.ReadWord(vector)
			c.cycles += 19
		}

		return nil
	}

	return nil
}
