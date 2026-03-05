package z80

import "fmt"

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

	if c.halted {
		// CPU is halted, just advance cycles
		c.cycles += 4
		return nil
	}

	// Handle interrupts first
	c.handleInterrupts()

	pcBeforeDecode := c.PC
	opcode, opcodeByte, err := c.decodeNextInstruction()
	if err != nil {
		return err
	}
	// Capture oldPC after decode: decode may advance PC past DD/FD prefix
	// when falling through to unprefixed instruction execution.
	oldPC := c.PC

	c.cycles += uint64(opcode.Timing)

	// Store current opcode for instruction functions to access
	c.currentOpcode = opcodeByte

	// Increment refresh register.
	// Prefixed instructions (CB, ED, DD, FD) increment R twice:
	// once for the prefix byte fetch and once for the opcode byte fetch.
	// DD/FD passthrough (prefix consumed + unprefixed executed) also needs R+2.
	prefixed := opcodeByte == PrefixCB || opcodeByte == PrefixED ||
		opcodeByte == PrefixDD || opcodeByte == PrefixFD ||
		c.PC != pcBeforeDecode // DD/FD passthrough advanced PC
	if prefixed {
		c.R = (c.R & 0x80) | ((c.R + 2) & 0x7F)
	} else {
		c.R = (c.R & 0x80) | ((c.R + 1) & 0x7F)
	}

	ins := opcode.Instruction
	if ins.NoParamFunc != nil {
		if c.opts.preExecutionHook != nil {
			c.opts.preExecutionHook(c, opcodeByte)
		}

		if err := ins.NoParamFunc(c); err != nil {
			return fmt.Errorf("executing no param instruction %s: %w", ins.Name, err)
		}
		c.updatePC(ins, oldPC, int(opcode.Size))
		c.q = c.GetFlags()
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
	c.q = c.GetFlags()
	return nil
}

// decodeNextInstruction decodes the current instruction at the program counter.
func (c *CPU) decodeNextInstruction() (Opcode, uint8, error) {
	// Handle extended instruction prefixes first
	opcodeByte := c.memory.Read(c.PC)

	switch opcodeByte {
	case PrefixCB:
		// CB-prefixed instructions (bit operations)
		return c.decodeCBInstruction()

	case PrefixED:
		// ED-prefixed instructions (extended operations)
		return c.decodeEDInstruction()

	case PrefixDD:
		// DD-prefixed instructions (IX operations)
		return c.decodeDDInstruction()

	case PrefixFD:
		// FD-prefixed instructions (IY operations)
		return c.decodeFDInstruction()
	}

	// Single-byte instructions
	opcode := Opcodes[opcodeByte]
	if opcode.Instruction == nil {
		return Opcode{}, opcodeByte, fmt.Errorf("%w: opcode 0x%02x", ErrUnsupportedOpcode, opcodeByte)
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
	// Check if this is a jump instruction that always changes PC
	if ins != nil && isJumpInstruction(ins) {
		// Jump instructions handle PC themselves, don't modify it
		return
	}

	// Update PC only if the instruction execution did not change it
	if oldPC == c.PC {
		// PC unchanged, advance by instruction size
		c.PC += uint16(amount)
		return
	}

	// PC was changed by the instruction (e.g., conditional jump taken), don't modify it further
}

// isJumpInstruction checks if an instruction is an unconditional jump/branch instruction that always modifies PC.
// Conditional jumps (like DJNZ, conditional JR/JP) are not included since they may or may not change PC.
func isJumpInstruction(ins *Instruction) bool {
	if ins == nil {
		return false
	}
	// Check for specific unconditional jump instructions by comparing pointers
	// This is the most precise approach since conditional and unconditional variants have same names
	return ins == JpAbs || ins == JpCond || ins == JrRel || ins == JrCond ||
		ins == Call || ins == CallCond || ins == Ret || ins == RetCond ||
		ins == EdReti || ins == EdRetn || ins == edRetnAlias || ins == Rst || ins == JpIndirect ||
		ins == DdJpIX || ins == FdJpIY || ins == Djnz ||
		ins == EdLdir || ins == EdLddr || ins == EdCpir || ins == EdCpdr ||
		ins == EdInir || ins == EdIndr || ins == EdOtir || ins == EdOtdr
}

// decodeCBInstruction decodes CB-prefixed instructions (bit operations).
func (c *CPU) decodeCBInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual CB instruction

	opcode := CBOpcodes[opcodeByte]
	if opcode.Instruction == nil {
		return Opcode{}, PrefixCB, fmt.Errorf("%w: opcode CB %02X", ErrUnsupportedOpcode, opcodeByte)
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

// decodeEDInstruction decodes ED-prefixed instructions (extended operations).
func (c *CPU) decodeEDInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual ED instruction

	opcode := EDOpcodes[opcodeByte]
	if opcode.Instruction == nil {
		return Opcode{}, PrefixED, fmt.Errorf("%w: opcode ED %02X", ErrUnsupportedEDOpcode, opcodeByte)
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixED, opcodeByte},
		}
	}

	return opcode, PrefixED, nil
}

// decodeDDInstruction decodes DD-prefixed instructions (IX operations).
func (c *CPU) decodeDDInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual DD instruction

	// Handle DD CB prefix first
	if opcodeByte == PrefixCB {
		return c.decodeDDCBInstruction()
	}

	opcode := DDOpcodes[opcodeByte]
	if opcode.Instruction == nil {
		// Undocumented behavior: DD prefix with no IX-specific instruction
		// executes the unprefixed instruction with 4 extra T-states.
		// Advance PC past the DD prefix so param readers see the correct offsets.
		c.PC++
		unprefixed := Opcodes[opcodeByte]
		if unprefixed.Instruction == nil {
			return Opcode{}, PrefixDD, fmt.Errorf("%w: opcode DD %02X", ErrUnsupportedOpcode, opcodeByte)
		}
		unprefixed.Timing += 4
		return unprefixed, opcodeByte, nil
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixDD, opcodeByte},
		}
	}

	return opcode, PrefixDD, nil
}

// decodeDDCBInstruction decodes DD CB prefixed instructions (IX bit operations).
func (c *CPU) decodeDDCBInstruction() (Opcode, uint8, error) {
	displacement := int8(c.memory.Read(c.PC + 2)) // Get displacement
	opcodeByte := c.memory.Read(c.PC + 3)         // Get bit operation

	var instruction *Instruction
	var timing byte = 23 // All DDCB operations take 23 T-states

	switch {
	case opcodeByte <= 0x3F: // Rotate/shift operations
		instruction = DdcbShift
	case opcodeByte <= 0x7F: // BIT operations
		instruction = DdcbBit
	case opcodeByte <= 0xBF: // RES operations
		instruction = DdcbRes
	default: // SET operations (0xC0-0xFF)
		instruction = DdcbSet
	}

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        4,
		Timing:      timing,
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixDD, PrefixCB, uint8(displacement), opcodeByte},
		}
	}

	return opcode, PrefixDD, nil
}

// decodeFDInstruction decodes FD-prefixed instructions (IY operations).
func (c *CPU) decodeFDInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual FD instruction

	// Handle FD CB prefix first
	if opcodeByte == PrefixCB {
		return c.decodeFDCBInstruction()
	}

	opcode := FDOpcodes[opcodeByte]
	if opcode.Instruction == nil {
		// Undocumented behavior: FD prefix with no IY-specific instruction
		// executes the unprefixed instruction with 4 extra T-states.
		// Advance PC past the FD prefix so param readers see the correct offsets.
		c.PC++
		unprefixed := Opcodes[opcodeByte]
		if unprefixed.Instruction == nil {
			return Opcode{}, PrefixFD, fmt.Errorf("%w: opcode FD %02X", ErrUnsupportedOpcode, opcodeByte)
		}
		unprefixed.Timing += 4
		return unprefixed, opcodeByte, nil
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixFD, opcodeByte},
		}
	}

	return opcode, PrefixFD, nil
}

// decodeFDCBInstruction decodes FD CB prefixed instructions (IY bit operations).
func (c *CPU) decodeFDCBInstruction() (Opcode, uint8, error) {
	displacement := int8(c.memory.Read(c.PC + 2)) // Get displacement
	opcodeByte := c.memory.Read(c.PC + 3)         // Get bit operation

	var instruction *Instruction
	var timing byte = 23 // All FDCB operations take 23 T-states

	switch {
	case opcodeByte <= 0x3F: // Rotate/shift operations
		instruction = FdcbShift
	case opcodeByte <= 0x7F: // BIT operations
		instruction = FdcbBit
	case opcodeByte <= 0xBF: // RES operations
		instruction = FdcbRes
	default: // SET operations (0xC0-0xFF)
		instruction = FdcbSet
	}

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        4,
		Timing:      timing,
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixFD, PrefixCB, uint8(displacement), opcodeByte},
		}
	}

	return opcode, PrefixFD, nil
}

// handleInterrupts processes pending interrupts.
func (c *CPU) handleInterrupts() {
	// Non-maskable interrupt has the highest priority
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
		return
	}

	// Maskable interrupt
	if c.triggerIrq && c.iff1 {
		c.triggerIrq = false
		c.halted = false
		c.iff1 = false
		c.iff2 = false

		// Save current PC
		c.push16(c.PC)

		// NOTE: Interrupt handling is simplified. In real hardware:
		// - IM 0: Device places instruction on data bus, CPU executes it
		// - IM 2: Device provides low byte of vector address on data bus
		// This implementation assumes RST 38H (0xFF) for IM 0 and reads
		// the vector low byte from 0xFFFF for IM 2.
		switch c.im {
		case 0:
			// Simplified: assumes RST 38H instruction on data bus
			c.PC = 0x0038
			c.cycles += 13
		case 1:
			c.PC = 0x0038
			c.cycles += 13
		case 2:
			// Simplified: reads vector low byte from 0xFFFF instead of data bus
			vector := uint16(c.I)<<8 | uint16(c.memory.Read(0xFFFF))
			c.PC = c.memory.ReadWord(vector)
			c.cycles += 19
		}
	}
}
