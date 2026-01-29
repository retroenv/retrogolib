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
	return ins == JpAbs || ins == JrRel || ins == Call || ins == Ret || ins == EdReti || ins == EdRetn ||
		ins == Rst || ins == JpIndirect ||
		ins == DdJpIX || ins == FdJpIY
}

// decodeCBInstruction decodes CB-prefixed instructions (bit operations).
func (c *CPU) decodeCBInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual CB instruction

	instruction, timing := c.decodeCBInstructionType(opcodeByte)

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        2,
		Timing:      timing,
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

// decodeCBInstructionType determines the instruction and timing for CB-prefixed opcodes.
func (c *CPU) decodeCBInstructionType(opcodeByte uint8) (*Instruction, byte) {
	// CB instructions follow a pattern:
	// 00-07: RLC r    08-0F: RRC r    10-17: RL r     18-1F: RR r
	// 20-27: SLA r    28-2F: SRA r    30-37: SLL r    38-3F: SRL r
	// 40-7F: BIT n,r  80-BF: RES n,r  C0-FF: SET n,r

	reg := opcodeByte & 0x07 // Lower 3 bits determine register

	switch {
	case opcodeByte <= 0x3F:
		return c.decodeCBRotateShift(opcodeByte, reg)
	case opcodeByte <= 0x7F:
		return c.decodeCBBit(reg)
	case opcodeByte <= 0xBF:
		return c.decodeCBRes(reg)
	default:
		return c.decodeCBSet(reg)
	}
}

// decodeCBRotateShift handles CB rotate/shift instructions (0x00-0x3F).
func (c *CPU) decodeCBRotateShift(opcodeByte, reg uint8) (*Instruction, byte) {
	var instruction *Instruction

	switch {
	case opcodeByte <= 0x07: // RLC r
		instruction = CBRlc
	case opcodeByte <= 0x0F: // RRC r
		instruction = CBRrc
	case opcodeByte <= 0x17: // RL r
		instruction = CBRl
	case opcodeByte <= 0x1F: // RR r
		instruction = CBRr
	case opcodeByte <= 0x27: // SLA r
		instruction = CBSla
	case opcodeByte <= 0x2F: // SRA r
		instruction = CBSra
	case opcodeByte <= 0x37: // SLL r (undocumented)
		instruction = CBSll
	default: // SRL r (0x38-0x3F)
		instruction = CBSrl
	}

	// Use helper function for timing calculation
	timing := getCBTiming(opcodeByte, reg)
	return instruction, timing
}

// decodeCBBit handles CB BIT instructions (0x40-0x7F).
func (c *CPU) decodeCBBit(reg uint8) (*Instruction, byte) {
	timing := byte(8)
	if reg == 6 { // BIT n,(HL)
		timing = 12
	}
	return CBBit, timing
}

// decodeCBRes handles CB RES instructions (0x80-0xBF).
func (c *CPU) decodeCBRes(reg uint8) (*Instruction, byte) {
	timing := byte(8)
	if reg == 6 { // RES n,(HL)
		timing = 15
	}
	return CBRes, timing
}

// decodeCBSet handles CB SET instructions (0xC0-0xFF).
func (c *CPU) decodeCBSet(reg uint8) (*Instruction, byte) {
	timing := byte(8)
	if reg == 6 { // SET n,(HL)
		timing = 15
	}
	return CBSet, timing
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
		// Handle undocumented behavior: DD prefix alone acts as 4-cycle NOP
		// This occurs when DD is followed by an invalid IX instruction
		opcode = Opcode{
			Instruction: Nop,
			Addressing:  ImpliedAddressing,
			Timing:      4,
			Size:        1,
		}
		// Return nil error since we're handling it as undocumented NOP
		return opcode, PrefixDD, nil
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
		// Handle undocumented behavior: FD prefix alone acts as 4-cycle NOP
		// This occurs when FD is followed by an invalid IY instruction
		opcode = Opcode{
			Instruction: Nop,
			Addressing:  ImpliedAddressing,
			Timing:      4,
			Size:        1,
		}
		// Return nil error since we're handling it as undocumented NOP
		return opcode, PrefixFD, nil
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
