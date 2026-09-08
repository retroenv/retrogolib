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

// Step accepts a pending interrupt, performs a HALT idle cycle, or executes one instruction.
func (cpu *CPU) Step() error {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()

	if cpu.handleInterrupts() {
		return nil
	}

	if cpu.halted {
		cpu.incrementRefresh(false)
		cpu.cycles += 4
		return nil
	}

	cpu.eiPending = false
	cpu.lastWasLdAIR = false

	pcBeforeDecode := cpu.PC
	opcode, opcodeByte, err := cpu.decodeNextInstruction()
	if err != nil {
		return err
	}
	// Capture oldPC after decode: decode may advance PC past DD/FD prefix
	// when falling through to unprefixed instruction execution.
	oldPC := cpu.PC

	cpu.cycles += uint64(opcode.Timing)

	// Store current opcode for instruction functions to access
	cpu.currentOpcode = opcodeByte

	// Prefixed instructions (CB, ED, DD, FD) and DD/FD passthrough increment R twice.
	prefixed := opcodeByte == PrefixCB || opcodeByte == PrefixED ||
		opcodeByte == PrefixDD || opcodeByte == PrefixFD ||
		cpu.PC != pcBeforeDecode
	cpu.incrementRefresh(prefixed)

	if err := cpu.executeInstruction(opcode, opcodeByte, oldPC); err != nil {
		return err
	}
	cpu.q = cpu.GetFlags()
	return nil
}

// executeInstruction runs the decoded instruction and updates the program counter.
func (cpu *CPU) executeInstruction(opcode Opcode, opcodeByte byte, oldPC uint16) error {
	ins := opcode.Instruction
	if ins.NoParamFunc != nil {
		if cpu.opts.preExecutionHook != nil {
			cpu.opts.preExecutionHook(cpu, opcodeByte)
		}
		if err := ins.NoParamFunc(cpu); err != nil {
			return fmt.Errorf("executing no param instruction %s: %w", ins.Name, err)
		}
		cpu.updatePC(ins, oldPC, int(opcode.Size))
		return nil
	}

	params, operands, err := readOpParams(cpu, opcode.Addressing)
	if err != nil {
		return fmt.Errorf("reading opcode params: %w", err)
	}
	if cpu.opts.tracing {
		cpu.TraceStep.OpcodeOperands = append(cpu.TraceStep.OpcodeOperands, operands...)
	}
	if cpu.opts.preExecutionHook != nil {
		cpu.opts.preExecutionHook(cpu, opcodeByte, params...)
	}

	if err := ins.ParamFunc(cpu, params...); err != nil {
		return fmt.Errorf("executing param instruction %s: %w", ins.Name, err)
	}
	cpu.updatePC(ins, oldPC, int(opcode.Size))
	return nil
}

// incrementRefresh increments the memory refresh register R.
// Preserves bit 7 and increments the lower 7 bits. Prefixed instructions increment by 2.
func (cpu *CPU) incrementRefresh(prefixed bool) {
	inc := uint8(1)
	if prefixed {
		inc = 2
	}
	cpu.R = (cpu.R & 0x80) | ((cpu.R + inc) & 0x7F)
}

// decodeNextInstruction decodes the current instruction at the program counter.
func (cpu *CPU) decodeNextInstruction() (Opcode, uint8, error) {
	// Handle extended instruction prefixes first
	opcodeByte := cpu.bus.Read(cpu.PC)

	switch opcodeByte {
	case PrefixCB:
		// CB-prefixed instructions (bit operations)
		return cpu.decodeCBInstruction()

	case PrefixED:
		// ED-prefixed instructions (extended operations)
		return cpu.decodeEDInstruction()

	case PrefixDD:
		// DD-prefixed instructions (IX operations)
		return cpu.decodeDDInstruction()

	case PrefixFD:
		// FD-prefixed instructions (IY operations)
		return cpu.decodeFDInstruction()
	}

	// Single-byte instructions
	opcode := Opcodes[opcodeByte]
	if opcode.Instruction == nil {
		return Opcode{}, opcodeByte, fmt.Errorf("%w: opcode 0x%02x", ErrUnsupportedOpcode, opcodeByte)
	}

	if cpu.opts.tracing {
		cpu.TraceStep = TraceStep{
			PC:             cpu.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{opcodeByte},
		}
	}
	return opcode, opcodeByte, nil
}

// updatePC updates the program counter based on the instruction execution.
func (cpu *CPU) updatePC(ins *Instruction, oldPC uint16, amount int) {
	// Check if this is a jump instruction that always changes PC
	if ins != nil && isJumpInstruction(ins) {
		// Jump instructions handle PC themselves, don't modify it
		return
	}

	// Update PC only if the instruction execution did not change it
	if oldPC == cpu.PC {
		// PC unchanged, advance by instruction size
		cpu.PC += uint16(amount)
		return
	}

	// PC was changed by the instruction (e.g., conditional jump taken), don't modify it further
}

// decodeCBInstruction decodes CB-prefixed instructions (bit operations).
func (cpu *CPU) decodeCBInstruction() (Opcode, uint8, error) {
	opcodeByte := cpu.bus.Read(cpu.PC + 1) // Get the actual CB instruction

	opcode := CBOpcodes[opcodeByte]
	if opcode.Instruction == nil {
		return Opcode{}, PrefixCB, fmt.Errorf("%w: opcode CB %02X", ErrUnsupportedOpcode, opcodeByte)
	}

	if cpu.opts.tracing {
		cpu.TraceStep = TraceStep{
			PC:             cpu.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixCB, opcodeByte},
		}
	}

	return opcode, PrefixCB, nil
}

// decodeEDInstruction decodes ED-prefixed instructions (extended operations).
func (cpu *CPU) decodeEDInstruction() (Opcode, uint8, error) {
	opcodeByte := cpu.bus.Read(cpu.PC + 1) // Get the actual ED instruction

	opcode := EDOpcodes[opcodeByte]
	if opcode.Instruction == nil {
		return Opcode{}, PrefixED, fmt.Errorf("%w: opcode ED %02X", ErrUnsupportedEDOpcode, opcodeByte)
	}

	if cpu.opts.tracing {
		cpu.TraceStep = TraceStep{
			PC:             cpu.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixED, opcodeByte},
		}
	}

	return opcode, PrefixED, nil
}

// decodeDDInstruction decodes DD-prefixed instructions (IX operations).
func (cpu *CPU) decodeDDInstruction() (Opcode, uint8, error) {
	opcodeByte := cpu.bus.Read(cpu.PC + 1) // Get the actual DD instruction

	// Handle DD CB prefix first
	if opcodeByte == PrefixCB {
		return cpu.decodeDDCBInstruction()
	}

	opcode := DDOpcodes[opcodeByte]
	if opcode.Instruction == nil {
		// Undocumented behavior: DD prefix with no IX-specific instruction
		// executes the unprefixed instruction with 4 extra T-states.
		// Advance PC past the DD prefix so param readers see the correct offsets.
		cpu.PC++
		unprefixed := Opcodes[opcodeByte]
		if unprefixed.Instruction == nil {
			return Opcode{}, PrefixDD, fmt.Errorf("%w: opcode DD %02X", ErrUnsupportedOpcode, opcodeByte)
		}
		cpu.q = 0 // An ignored prefix does not modify flags.
		unprefixed.Timing += 4
		return unprefixed, opcodeByte, nil
	}

	if cpu.opts.tracing {
		cpu.TraceStep = TraceStep{
			PC:             cpu.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixDD, opcodeByte},
		}
	}

	return opcode, PrefixDD, nil
}

// decodeDDCBInstruction decodes DD CB prefixed instructions (IX bit operations).
func (cpu *CPU) decodeDDCBInstruction() (Opcode, uint8, error) {
	displacement := int8(cpu.bus.Read(cpu.PC + 2)) // Get displacement
	opcodeByte := cpu.bus.Read(cpu.PC + 3)         // Get bit operation

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

	if cpu.opts.tracing {
		cpu.TraceStep = TraceStep{
			PC:             cpu.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixDD, PrefixCB, uint8(displacement), opcodeByte},
		}
	}

	return opcode, PrefixDD, nil
}

// decodeFDInstruction decodes FD-prefixed instructions (IY operations).
func (cpu *CPU) decodeFDInstruction() (Opcode, uint8, error) {
	opcodeByte := cpu.bus.Read(cpu.PC + 1) // Get the actual FD instruction

	// Handle FD CB prefix first
	if opcodeByte == PrefixCB {
		return cpu.decodeFDCBInstruction()
	}

	opcode := FDOpcodes[opcodeByte]
	if opcode.Instruction == nil {
		// Undocumented behavior: FD prefix with no IY-specific instruction
		// executes the unprefixed instruction with 4 extra T-states.
		// Advance PC past the FD prefix so param readers see the correct offsets.
		cpu.PC++
		unprefixed := Opcodes[opcodeByte]
		if unprefixed.Instruction == nil {
			return Opcode{}, PrefixFD, fmt.Errorf("%w: opcode FD %02X", ErrUnsupportedOpcode, opcodeByte)
		}
		cpu.q = 0 // An ignored prefix does not modify flags.
		unprefixed.Timing += 4
		return unprefixed, opcodeByte, nil
	}

	if cpu.opts.tracing {
		cpu.TraceStep = TraceStep{
			PC:             cpu.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixFD, opcodeByte},
		}
	}

	return opcode, PrefixFD, nil
}

// decodeFDCBInstruction decodes FD CB prefixed instructions (IY bit operations).
func (cpu *CPU) decodeFDCBInstruction() (Opcode, uint8, error) {
	displacement := int8(cpu.bus.Read(cpu.PC + 2)) // Get displacement
	opcodeByte := cpu.bus.Read(cpu.PC + 3)         // Get bit operation

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

	if cpu.opts.tracing {
		cpu.TraceStep = TraceStep{
			PC:             cpu.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{PrefixFD, PrefixCB, uint8(displacement), opcodeByte},
		}
	}

	return opcode, PrefixFD, nil
}

// jumpInstructions is a lookup set of instructions that always modify PC.
// These include all jump, call, return, and repeat block instructions.
var jumpInstructions = set.Set[*Instruction]{
	CallInst:    {},
	CallCond:    {},
	DdJpIX:      {},
	DjnzInst:    {},
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
	RetInst:     {},
	RetCond:     {},
	RstInst:     {},
	edRetnAlias: {},
}

// isJumpInstruction checks if an instruction is a jump/branch instruction that always modifies PC.
func isJumpInstruction(ins *Instruction) bool {
	return ins != nil && jumpInstructions.Contains(ins)
}
