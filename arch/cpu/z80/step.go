package z80

import (
	"errors"
	"fmt"
)

// TraceStep contains information about the current step for tracing.
type TraceStep struct {
	PC     uint16
	Opcode uint8

	// Register values before execution
	A, B, C, D, E, H, L uint8
	SP, IX, IY          uint16
	Flags               uint8

	// Instruction info
	InstructionName string
	CyclesTaken     uint8
}

// Step executes a single instruction and returns any error.
func (c *CPU) Step() error {
	if c.halted {
		// CPU is halted, just advance cycles
		c.cycles += 4
		return nil
	}

	// Validate memory is available
	if c.memory == nil {
		return errors.New("z80: memory is nil")
	}

	// Handle interrupts first
	if err := c.handleInterrupts(); err != nil {
		return err
	}

	// Save trace information if tracing is enabled
	if c.opts.TraceExecution {
		c.TraceStep = TraceStep{
			PC:    c.PC,
			A:     c.A,
			B:     c.B,
			C:     c.C,
			D:     c.D,
			E:     c.E,
			H:     c.H,
			L:     c.L,
			SP:    c.SP,
			IX:    c.IX,
			IY:    c.IY,
			Flags: c.GetFlags(),
		}
	}

	// Fetch instruction with bounds checking
	// PC should be within valid memory range
	opcode := c.memory.Read(c.PC)
	if c.opts.TraceExecution {
		c.TraceStep.Opcode = opcode
	}

	// Increment refresh register
	c.R = (c.R & 0x80) | ((c.R + 1) & 0x7F)

	// Execute instruction
	cycles, err := c.executeInstruction(opcode)
	if err != nil {
		return fmt.Errorf("error executing instruction at PC=0x%04X, opcode=0x%02X: %w", c.PC, opcode, err)
	}

	c.cycles += uint64(cycles)
	if c.opts.TraceExecution {
		c.TraceStep.CyclesTaken = cycles
	}

	return nil
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
			// For Game Boy, this is typically RST 40h
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

// executeInstruction executes a single instruction and returns the number of cycles taken.
func (c *CPU) executeInstruction(opcode uint8) (uint8, error) {
	// Handle extended instruction prefixes first
	switch opcode {
	case 0xCB:
		return c.executeCBInstruction()
	case 0xED:
		return c.executeEDInstruction()
	case 0xDD:
		return c.executeDDInstruction()
	case 0xFD:
		return c.executeFDInstruction()
	}

	// Handle single-byte instructions
	return c.executeSingleByteInstruction(opcode)
}

// executeSingleByteInstruction handles non-prefixed Z80 instructions.
func (c *CPU) executeSingleByteInstruction(opcode uint8) (uint8, error) {
	switch {
	case opcode <= 0x0F:
		return c.executeBasicInstructions(opcode)
	case opcode == 0x10:
		return c.executeDJNZ()
	case opcode == 0x76:
		return c.executeHALT()
	case opcode >= 0xC0:
		return c.executeControlInstructions(opcode)
	default:
		return 4, fmt.Errorf("unimplemented opcode: 0x%02X", opcode)
	}
}

// executeBasicInstructions handles opcodes 0x00-0x0F.
func (c *CPU) executeBasicInstructions(opcode uint8) (uint8, error) {
	switch {
	case opcode <= 0x07:
		return c.executeBasicInstructions0x00To0x07(opcode)
	case opcode <= 0x0F:
		return c.executeBasicInstructions0x08To0x0F(opcode)
	}
	return 4, fmt.Errorf("unimplemented basic opcode: 0x%02X", opcode)
}

// executeBasicInstructions0x00To0x07 handles opcodes 0x00-0x07.
func (c *CPU) executeBasicInstructions0x00To0x07(opcode uint8) (uint8, error) {
	switch opcode {
	case 0x00: // NOP
		c.PC++
		return 4, nil
	case 0x01: // LD BC,nn
		nn := c.memory.ReadWord(c.PC + 1)
		c.setBC(nn)
		c.PC += 3
		return 10, nil
	case 0x02: // LD (BC),A
		c.memory.Write(c.BC(), c.A)
		c.PC++
		return 7, nil
	case 0x03: // INC BC
		c.setBC(c.BC() + 1)
		c.PC++
		return 6, nil
	case 0x04: // INC B
		c.B = c.inc8(c.B)
		c.PC++
		return 4, nil
	case 0x05: // DEC B
		c.B = c.dec8(c.B)
		c.PC++
		return 4, nil
	case 0x06: // LD B,n
		c.B = c.memory.Read(c.PC + 1)
		c.PC += 2
		return 7, nil
	case 0x07: // RLCA
		c.A = c.rlca(c.A)
		c.PC++
		return 4, nil
	}
	return 4, fmt.Errorf("unimplemented basic opcode: 0x%02X", opcode)
}

// executeBasicInstructions0x08To0x0F handles opcodes 0x08-0x0F.
func (c *CPU) executeBasicInstructions0x08To0x0F(opcode uint8) (uint8, error) {
	switch opcode {
	case 0x08: // EX AF,AF'
		c.exchangeAF()
		c.PC++
		return 4, nil
	case 0x09: // ADD HL,BC
		c.setHL(c.add16(c.HL(), c.BC()))
		c.PC++
		return 11, nil
	case 0x0A: // LD A,(BC)
		c.A = c.memory.Read(c.BC())
		c.PC++
		return 7, nil
	case 0x0B: // DEC BC
		c.setBC(c.BC() - 1)
		c.PC++
		return 6, nil
	case 0x0C: // INC C
		c.C = c.inc8(c.C)
		c.PC++
		return 4, nil
	case 0x0D: // DEC C
		c.C = c.dec8(c.C)
		c.PC++
		return 4, nil
	case 0x0E: // LD C,n
		c.C = c.memory.Read(c.PC + 1)
		c.PC += 2
		return 7, nil
	case 0x0F: // RRCA
		c.A = c.rrca(c.A)
		c.PC++
		return 4, nil
	}
	return 4, fmt.Errorf("unimplemented basic opcode: 0x%02X", opcode)
}

// executeDJNZ handles the DJNZ instruction.
func (c *CPU) executeDJNZ() (uint8, error) {
	c.B--
	if c.B != 0 {
		offset := int8(c.memory.Read(c.PC + 1))
		c.PC = uint16(int32(c.PC) + int32(offset) + 2)
		return 13, nil
	}
	c.PC += 2
	return 8, nil
}

// executeHALT handles the HALT instruction.
func (c *CPU) executeHALT() (uint8, error) {
	c.halted = true
	c.PC++
	return 4, nil
}

// executeControlInstructions handles control flow instructions (0xC0-0xFF range).
func (c *CPU) executeControlInstructions(opcode uint8) (uint8, error) {
	switch opcode {
	case 0xC3: // JP nn
		c.PC = c.memory.ReadWord(c.PC + 1)
		return 10, nil
	case 0xCD: // CALL nn
		c.push16(c.PC + 3)
		c.PC = c.memory.ReadWord(c.PC + 1)
		return 17, nil
	case 0xC9: // RET
		c.PC = c.pop16()
		return 10, nil
	case 0xF3: // DI
		c.iff1 = false
		c.iff2 = false
		c.PC++
		return 4, nil
	case 0xFB: // EI
		c.iff1 = true
		c.iff2 = true
		c.PC++
		return 4, nil
	case 0xFF: // RST 38h
		c.push16(c.PC + 1)
		c.PC = 0x0038
		return 11, nil
	}
	return 4, fmt.Errorf("unimplemented control opcode: 0x%02X", opcode)
}

// executeCBInstruction executes CB-prefixed instructions.
func (c *CPU) executeCBInstruction() (uint8, error) {
	opcode := c.memory.Read(c.PC + 1)
	c.PC += 2

	// CB instructions are mostly bit operations
	switch opcode {
	case 0x00: // RLC B
		c.B = c.rlc(c.B)
		return 8, nil
	case 0x01: // RLC C
		c.C = c.rlc(c.C)
		return 8, nil
	// Add more CB instructions as needed
	default:
		return 8, fmt.Errorf("unimplemented CB instruction: 0x%02X", opcode)
	}
}

// executeEDInstruction executes ED-prefixed instructions.
func (c *CPU) executeEDInstruction() (uint8, error) {
	opcode := c.memory.Read(c.PC + 1)
	c.PC += 2

	switch opcode {
	case 0x44: // NEG
		c.A = c.neg(c.A)
		return 8, nil
	case 0x46: // IM 0
		c.im = 0
		return 8, nil
	case 0x56: // IM 1
		c.im = 1
		return 8, nil
	case 0x5E: // IM 2
		c.im = 2
		return 8, nil
	// Add more ED instructions as needed
	default:
		return 8, fmt.Errorf("unimplemented ED instruction: 0x%02X", opcode)
	}
}

// executeDDInstruction executes DD-prefixed instructions (IX).
func (c *CPU) executeDDInstruction() (uint8, error) {
	opcode := c.memory.Read(c.PC + 1)
	c.PC += 2

	// DD instructions work with IX register
	switch opcode {
	case 0x21: // LD IX,nn
		c.IX = c.memory.ReadWord(c.PC)
		c.PC += 2
		return 14, nil
	// Add more DD instructions as needed
	default:
		return 8, fmt.Errorf("unimplemented DD instruction: 0x%02X", opcode)
	}
}

// executeFDInstruction executes FD-prefixed instructions (IY).
func (c *CPU) executeFDInstruction() (uint8, error) {
	opcode := c.memory.Read(c.PC + 1)
	c.PC += 2

	// FD instructions work with IY register
	switch opcode {
	case 0x21: // LD IY,nn
		c.IY = c.memory.ReadWord(c.PC)
		c.PC += 2
		return 14, nil
	// Add more FD instructions as needed
	default:
		return 8, fmt.Errorf("unimplemented FD instruction: 0x%02X", opcode)
	}
}
