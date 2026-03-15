package z80

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

	if c.halted {
		// CPU is halted, just advance cycles
		c.cycles += 4
		return nil
	}

	// Handle interrupts first
	c.handleInterrupts()

	c.lastWasLdAIR = false

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

	// Prefixed instructions (CB, ED, DD, FD) and DD/FD passthrough increment R twice.
	prefixed := opcodeByte == PrefixCB || opcodeByte == PrefixED ||
		opcodeByte == PrefixDD || opcodeByte == PrefixFD ||
		c.PC != pcBeforeDecode
	c.incrementRefresh(prefixed)

	if err := c.executeInstruction(opcode, opcodeByte, oldPC); err != nil {
		return err
	}
	c.q = c.GetFlags()
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

// incrementRefresh increments the memory refresh register R.
// Preserves bit 7 and increments the lower 7 bits. Prefixed instructions increment by 2.
func (c *CPU) incrementRefresh(prefixed bool) {
	inc := uint8(1)
	if prefixed {
		inc = 2
	}
	c.R = (c.R & 0x80) | ((c.R + inc) & 0x7F)
}

// decodeNextInstruction decodes the current instruction at the program counter.
func (c *CPU) decodeNextInstruction() (Opcode, uint8, error) {
	// Handle extended instruction prefixes first
	opcodeByte := c.bus.Read(c.PC)

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

// decodeCBInstruction decodes CB-prefixed instructions (bit operations).
func (c *CPU) decodeCBInstruction() (Opcode, uint8, error) {
	opcodeByte := c.bus.Read(c.PC + 1) // Get the actual CB instruction

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
	opcodeByte := c.bus.Read(c.PC + 1) // Get the actual ED instruction

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
	opcodeByte := c.bus.Read(c.PC + 1) // Get the actual DD instruction

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
	displacement := int8(c.bus.Read(c.PC + 2)) // Get displacement
	opcodeByte := c.bus.Read(c.PC + 3)         // Get bit operation

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
	opcodeByte := c.bus.Read(c.PC + 1) // Get the actual FD instruction

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
	displacement := int8(c.bus.Read(c.PC + 2)) // Get displacement
	opcodeByte := c.bus.Read(c.PC + 3)         // Get bit operation

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

		// Zilog NMOS bug: LD A,I and LD A,R set P/V from IFF2, but if an
		// interrupt is accepted immediately after, P/V is reset to 0.
		if c.lastWasLdAIR {
			c.Flags.P = 0
		}

		c.iff1 = false
		c.iff2 = false

		// Save current PC
		c.push16(c.PC)

		switch c.im {
		case 0:
			// IM 0: read instruction opcode from data bus via Bus.IRQData().
			// In practice this is almost always a RST instruction.
			dataBusValue := c.bus.IRQData()
			if dataBusValue&0xC7 == 0xC7 {
				// RST instruction: extract vector from bits 3-5
				vector := uint16(dataBusValue & 0x38)
				c.PC = vector
				c.MEMPTR = vector
			} else {
				// Fallback for non-RST instructions on the bus.
				// Full arbitrary instruction execution is not yet supported.
				c.PC = 0x0038
				c.MEMPTR = 0x0038
			}
			c.cycles += 13
		case 1:
			c.PC = 0x0038
			c.MEMPTR = 0x0038
			c.cycles += 13
		case 2:
			// IM 2: read vector low byte from data bus, combine with I register.
			vectorLow := c.bus.IRQData()
			vectorAddr := uint16(c.I)<<8 | uint16(vectorLow)
			c.PC = c.bus.ReadWord(vectorAddr)
			c.MEMPTR = c.PC
			c.cycles += 19
		}
	}
}

// jumpInstructions is a lookup set of instructions that always modify PC.
// These include all jump, call, return, and repeat block instructions.
var jumpInstructions = set.Set[*Instruction]{
	Call:        {},
	CallCond:    {},
	DdJpIX:      {},
	Djnz:        {},
	EdCpdr:      {},
	EdCpir:      {},
	EdIndr:      {},
	EdInir:      {},
	EdLddr:      {},
	EdLdir:      {},
	EdOtdr:      {},
	EdOtir:      {},
	EdReti:      {},
	EdRetn:      {},
	FdJpIY:      {},
	JpAbs:       {},
	JpCond:      {},
	JpIndirect:  {},
	JrCond:      {},
	JrRel:       {},
	Ret:         {},
	RetCond:     {},
	Rst:         {},
	edRetnAlias: {},
}

// isJumpInstruction checks if an instruction is a jump/branch instruction that always modifies PC.
func isJumpInstruction(ins *Instruction) bool {
	return ins != nil && jumpInstructions.Contains(ins)
}
