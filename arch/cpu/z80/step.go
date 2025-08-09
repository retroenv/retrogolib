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
func (cpu *CPU) Step() error {
	if cpu.halted {
		// CPU is halted, just advance cycles
		cpu.cycles += 4
		return nil
	}

	// Handle interrupts first
	if err := cpu.handleInterrupts(); err != nil {
		return err
	}

	oldPC := cpu.PC
	opcode, opcodeByte, err := cpu.decodeNextInstruction()
	if err != nil {
		return err
	}

	cpu.cycles += uint64(opcode.Timing)

	// Store current opcode for instruction functions to access
	cpu.currentOpcode = opcodeByte

	// Increment refresh register
	cpu.R = (cpu.R & 0x80) | ((cpu.R + 1) & 0x7F)

	ins := opcode.Instruction
	if ins.NoParamFunc != nil {
		if cpu.opts.preExecutionHook != nil {
			cpu.opts.preExecutionHook(c, opcodeByte)
		}

		if err := ins.NoParamFunc(c); err != nil {
			return fmt.Errorf("executing no param instruction %s: %w", ins.Name, err)
		}
		cpu.updatePC(ins, oldPC, int(opcode.Size))
		return nil
	}

	params, operands, err := readOpParams(c, opcode.Addressing)
	if err != nil {
		return fmt.Errorf("reading opcode params: %w", err)
	}
	if cpu.opts.tracing {
		cpu.TraceStep.OpcodeOperands = append(cpu.TraceStep.OpcodeOperands, operands...)
	}
	if cpu.opts.preExecutionHook != nil {
		cpu.opts.preExecutionHook(c, opcodeByte, params...)
	}

	opcodeLen := int(opcode.Size)

	if err := ins.ParamFunc(c, params...); err != nil {
		return fmt.Errorf("executing param instruction %s: %w", ins.Name, err)
	}
	cpu.updatePC(ins, oldPC, opcodeLen)
	return nil
}

// decodeNextInstruction decodes the current instruction at the program counter.
func (cpu *CPU) decodeNextInstruction() (Opcode, uint8, error) {
	// Handle extended instruction prefixes first
	opcodeByte := cpu.memory.Read(cpu.PC)

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
func (cpu *CPU) decodeCBInstruction() (Opcode, uint8, error) {
	opcodeByte := cpu.memory.Read(cpu.PC + 1) // Get the actual CB instruction

	instruction, timing := cpu.decodeCBInstructionType(opcodeByte)

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        2,
		Timing:      timing,
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

// decodeCBInstructionType determines the instruction and timing for CB-prefixed opcodes.
func (cpu *CPU) decodeCBInstructionType(opcodeByte uint8) (*Instruction, byte) {
	// CB instructions follow a pattern:
	// 00-07: RLC r    08-0F: RRC r    10-17: RL r     18-1F: RR r
	// 20-27: SLA r    28-2F: SRA r    30-37: SLL r    38-3F: SRL r
	// 40-7F: BIT n,r  80-BF: RES n,r  C0-FF: SET n,r

	reg := opcodeByte & 0x07 // Lower 3 bits determine register

	switch {
	case opcodeByte <= 0x3F:
		return cpu.decodeCBRotateShift(opcodeByte, reg)
	case opcodeByte <= 0x7F:
		return cpu.decodeCBBit(reg)
	case opcodeByte <= 0xBF:
		return cpu.decodeCBRes(reg)
	default:
		return cpu.decodeCBSet(reg)
	}
}

// decodeCBRotateShift handles CB rotate/shift instructions (0x00-0x3F).
func (cpu *CPU) decodeCBRotateShift(opcodeByte, reg uint8) (*Instruction, byte) {
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
func (cpu *CPU) decodeCBBit(reg uint8) (*Instruction, byte) {
	timing := byte(8)
	if reg == 6 { // BIT n,(HL)
		timing = 12
	}
	return CBBit, timing
}

// decodeCBRes handles CB RES instructions (0x80-0xBF).
func (cpu *CPU) decodeCBRes(reg uint8) (*Instruction, byte) {
	timing := byte(8)
	if reg == 6 { // RES n,(HL)
		timing = 15
	}
	return CBRes, timing
}

// decodeCBSet handles CB SET instructions (0xC0-0xFF).
func (cpu *CPU) decodeCBSet(reg uint8) (*Instruction, byte) {
	timing := byte(8)
	if reg == 6 { // SET n,(HL)
		timing = 15
	}
	return CBSet, timing
}

// decodeEDInstruction decodes ED-prefixed instructions (extended operations).
func (cpu *CPU) decodeEDInstruction() (Opcode, uint8, error) {
	opcodeByte := cpu.memory.Read(cpu.PC + 1) // Get the actual ED instruction

	instruction, timing, size, err := cpu.decodeEDInstructionType(opcodeByte)
	if err != nil {
		return Opcode{}, PrefixED, err
	}

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        size,
		Timing:      timing,
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

// decodeEDInstructionType determines the instruction, timing, and size for ED-prefixed opcodes.
func (cpu *CPU) decodeEDInstructionType(opcodeByte uint8) (*Instruction, byte, byte, error) {
	// Group instructions by functionality to reduce complexity
	if instruction, timing, size := cpu.decodeEDBasicInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeEDArithmeticInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeEDLoadInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeEDBlockInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeEDIOInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	return nil, 0, 0, fmt.Errorf("%w: opcode ED %02X", ErrUnsupportedEDOpcode, opcodeByte)
}

// decodeEDBasicInstructions handles basic ED instructions (NEG, IM, RETN, RETI, RRD, RLD).
func (cpu *CPU) decodeEDBasicInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// NEG - Negate Accumulator
	case 0x44, 0x4C, 0x54, 0x5C, 0x64, 0x6C, 0x74, 0x7C:
		return EdNeg, 8, 2

	// IM - Set Interrupt Mode
	case 0x46, 0x66: // IM 0
		return EdIm0, 8, 2
	case 0x56, 0x76: // IM 1
		return EdIm1, 8, 2
	case 0x5E, 0x7E: // IM 2
		return EdIm2, 8, 2

	// RETN/RETI
	case 0x45, 0x55, 0x65, 0x75: // RETN (multiple opcodes)
		return EdRetn, 14, 2
	case 0x4D: // RETI
		return EdReti, 14, 2

	// RRD/RLD
	case 0x67: // RRD
		return EdRrd, 18, 2
	case 0x6F: // RLD
		return EdRld, 18, 2
	}

	return nil, 0, 0
}

// decodeEDArithmeticInstructions handles ED arithmetic instructions (ADC HL,rr / SBC HL,rr).
func (cpu *CPU) decodeEDArithmeticInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// ADC HL,rr
	case 0x4A: // ADC HL,BC
		return EdAdcHlBc, 15, 2
	case 0x5A: // ADC HL,DE
		return EdAdcHlDe, 15, 2
	case 0x6A: // ADC HL,HL
		return EdAdcHlHl, 15, 2
	case 0x7A: // ADC HL,SP
		return EdAdcHlSp, 15, 2

	// SBC HL,rr
	case 0x42: // SBC HL,BC
		return EdSbcHlBc, 15, 2
	case 0x52: // SBC HL,DE
		return EdSbcHlDe, 15, 2
	case 0x62: // SBC HL,HL
		return EdSbcHlHl, 15, 2
	case 0x72: // SBC HL,SP
		return EdSbcHlSp, 15, 2
	}

	return nil, 0, 0
}

// decodeEDLoadInstructions handles ED load instructions.
func (cpu *CPU) decodeEDLoadInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// LD I,A / LD R,A
	case 0x47: // LD I,A
		return EdLdIA, 9, 2
	case 0x4F: // LD R,A
		return EdLdRA, 9, 2

	// LD A,I / LD A,R
	case 0x57: // LD A,I
		return EdLdAI, 9, 2
	case 0x5F: // LD A,R
		return EdLdAR, 9, 2

	// LD (nn),rr
	case 0x43: // LD (nn),BC
		return EdLdNnBc, 20, 4
	case 0x53: // LD (nn),DE
		return EdLdNnDe, 20, 4
	case 0x63: // LD (nn),HL
		return EdLdNnHl, 20, 4
	case 0x73: // LD (nn),SP
		return EdLdNnSp, 20, 4

	// LD rr,(nn)
	case 0x4B: // LD BC,(nn)
		return EdLdBcNn, 20, 4
	case 0x5B: // LD DE,(nn)
		return EdLdDeNn, 20, 4
	case 0x6B: // LD HL,(nn)
		return EdLdHlNn, 20, 4
	case 0x7B: // LD SP,(nn)
		return EdLdSpNn, 20, 4
	}

	return nil, 0, 0
}

// decodeEDBlockInstructions handles ED block transfer and search instructions.
func (cpu *CPU) decodeEDBlockInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// Block transfer instructions
	case 0xA0: // LDI
		return EdLdi, 16, 2
	case 0xA8: // LDD
		return EdLdd, 16, 2
	case 0xB0: // LDIR
		return EdLdir, 21, 2
	case 0xB8: // LDDR
		return EdLddr, 21, 2

	// Block search instructions
	case 0xA1: // CPI
		return EdCpi, 16, 2
	case 0xA9: // CPD
		return EdCpd, 16, 2
	case 0xB1: // CPIR
		return EdCpir, 21, 2
	case 0xB9: // CPDR
		return EdCpdr, 21, 2
	}

	return nil, 0, 0
}

// decodeEDIOInstructions handles ED I/O instructions.
func (cpu *CPU) decodeEDIOInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	// Block I/O instructions
	if instruction, timing, size := cpu.decodeEDBlockIO(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	// I/O with register - IN r,(C)
	if instruction, timing, size := cpu.decodeEDInputInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	// I/O with register - OUT (C),r
	if instruction, timing, size := cpu.decodeEDOutputInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	return nil, 0, 0
}

// decodeEDBlockIO handles ED block I/O instructions.
func (cpu *CPU) decodeEDBlockIO(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0xA2: // INI
		return EdIni, 16, 2
	case 0xAA: // IND
		return EdInd, 16, 2
	case 0xB2: // INIR
		return EdInir, 21, 2
	case 0xBA: // INDR
		return EdIndr, 21, 2
	case 0xA3: // OUTI
		return EdOuti, 16, 2
	case 0xAB: // OUTD
		return EdOutd, 16, 2
	case 0xB3: // OTIR
		return EdOtir, 21, 2
	case 0xBB: // OTDR
		return EdOtdr, 21, 2
	}

	return nil, 0, 0
}

// decodeEDInputInstructions handles ED IN r,(C) instructions.
func (cpu *CPU) decodeEDInputInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x40: // IN B,(C)
		return EdInBC, 12, 2
	case 0x48: // IN C,(C)
		return EdInCC, 12, 2
	case 0x50: // IN D,(C)
		return EdInDC, 12, 2
	case 0x58: // IN E,(C)
		return EdInEC, 12, 2
	case 0x60: // IN H,(C)
		return EdInHC, 12, 2
	case 0x68: // IN L,(C)
		return EdInLC, 12, 2
	case 0x78: // IN A,(C)
		return EdInAC, 12, 2
	}

	return nil, 0, 0
}

// decodeEDOutputInstructions handles ED OUT (C),r instructions.
func (cpu *CPU) decodeEDOutputInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x41: // OUT (C),B
		return EdOutCB, 12, 2
	case 0x49: // OUT (C),C
		return EdOutCC, 12, 2
	case 0x51: // OUT (C),D
		return EdOutCD, 12, 2
	case 0x59: // OUT (C),E
		return EdOutCE, 12, 2
	case 0x61: // OUT (C),H
		return EdOutCH, 12, 2
	case 0x69: // OUT (C),L
		return EdOutCL, 12, 2
	case 0x79: // OUT (C),A
		return EdOutCA, 12, 2
	}

	return nil, 0, 0
}

// decodeDDInstruction decodes DD-prefixed instructions (IX operations).
func (cpu *CPU) decodeDDInstruction() (Opcode, uint8, error) {
	opcodeByte := cpu.memory.Read(cpu.PC + 1) // Get the actual DD instruction

	// Handle DD CB prefix first
	if opcodeByte == PrefixCB {
		return cpu.decodeDDCBInstruction()
	}

	instruction, timing, size, err := cpu.decodeDDInstructionType(opcodeByte)
	if err != nil {
		// Handle undocumented behavior: DD prefix alone acts as 4-cycle NOP
		// This occurs when DD is followed by an invalid IX instruction
		opcode := Opcode{
			Instruction: Nop,
			Addressing:  ImpliedAddressing,
			Timing:      4,
			Size:        1,
		}
		// Return nil error since we're handling it as undocumented NOP
		return opcode, PrefixDD, nil //nolint:nilerr
	}

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        size,
		Timing:      timing,
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

// decodeDDInstructionType determines the instruction, timing, and size for DD-prefixed opcodes.
func (cpu *CPU) decodeDDInstructionType(opcodeByte uint8) (*Instruction, byte, byte, error) {
	// Group instructions by functionality to reduce complexity
	if instruction, timing, size := cpu.decodeDDBasicInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeDDLoadInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeDDArithmeticInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeDDStackInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	return nil, 0, 0, fmt.Errorf("%w: opcode DD %02X", ErrUnsupportedDDOpcode, opcodeByte)
}

// decodeDDBasicInstructions handles basic DD instructions (INC/DEC IX, ADD IX,rr).
func (cpu *CPU) decodeDDBasicInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// INC IX / DEC IX
	case 0x23:
		return DdIncIX, 10, 2
	case 0x2B:
		return DdDecIX, 10, 2

	// ADD IX,rr
	case 0x09: // ADD IX,BC
		return DdAddIXBc, 15, 2
	case 0x19: // ADD IX,DE
		return DdAddIXDe, 15, 2
	case 0x29: // ADD IX,IX
		return DdAddIXIX, 15, 2
	case 0x39: // ADD IX,SP
		return DdAddIXSp, 15, 2

	// INC/DEC (IX+d)
	case 0x34: // INC (IX+d)
		return DdIncIXd, 23, 3
	case 0x35: // DEC (IX+d)
		return DdDecIXd, 23, 3
	}

	return nil, 0, 0
}

// decodeDDLoadInstructions handles DD load instructions.
func (cpu *CPU) decodeDDLoadInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// LD IX,nn
	case 0x21:
		return DdLdIXnn, 14, 4

	// LD (nn),IX / LD IX,(nn)
	case 0x22:
		return DdLdNnIX, 20, 4
	case 0x2A:
		return DdLdIXNn, 20, 4

	// LD (IX+d),n
	case 0x36:
		return DdLdIXdN, 19, 4
	}

	// LD r,(IX+d) - Load register from (IX+d)
	if instruction, timing, size := cpu.decodeDDLoadFromIX(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	// LD (IX+d),r - Load (IX+d) from register
	if instruction, timing, size := cpu.decodeDDLoadToIX(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	return nil, 0, 0
}

// decodeDDLoadFromIX handles LD r,(IX+d) instructions.
func (cpu *CPU) decodeDDLoadFromIX(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x46: // LD B,(IX+d)
		return DdLdBIXd, 19, 3
	case 0x4E: // LD C,(IX+d)
		return DdLdCIXd, 19, 3
	case 0x56: // LD D,(IX+d)
		return DdLdDIXd, 19, 3
	case 0x5E: // LD E,(IX+d)
		return DdLdEIXd, 19, 3
	case 0x66: // LD H,(IX+d)
		return DdLdHIXd, 19, 3
	case 0x6E: // LD L,(IX+d)
		return DdLdLIXd, 19, 3
	case 0x7E: // LD A,(IX+d)
		return DdLdAIXd, 19, 3
	}

	return nil, 0, 0
}

// decodeDDLoadToIX handles LD (IX+d),r instructions.
func (cpu *CPU) decodeDDLoadToIX(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x70: // LD (IX+d),B
		return DdLdIXdB, 19, 3
	case 0x71: // LD (IX+d),C
		return DdLdIXdC, 19, 3
	case 0x72: // LD (IX+d),D
		return DdLdIXdD, 19, 3
	case 0x73: // LD (IX+d),E
		return DdLdIXdE, 19, 3
	case 0x74: // LD (IX+d),H
		return DdLdIXdH, 19, 3
	case 0x75: // LD (IX+d),L
		return DdLdIXdL, 19, 3
	case 0x77: // LD (IX+d),A
		return DdLdIXdA, 19, 3
	}

	return nil, 0, 0
}

// decodeDDArithmeticInstructions handles DD arithmetic instructions with (IX+d).
func (cpu *CPU) decodeDDArithmeticInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x86: // ADD A,(IX+d)
		return DdAddAIXd, 19, 3
	case 0x8E: // ADC A,(IX+d)
		return DdAdcAIXd, 19, 3
	case 0x96: // SUB (IX+d)
		return DdSubAIXd, 19, 3
	case 0x9E: // SBC A,(IX+d)
		return DdSbcAIXd, 19, 3
	case 0xA6: // AND (IX+d)
		return DdAndAIXd, 19, 3
	case 0xAE: // XOR (IX+d)
		return DdXorAIXd, 19, 3
	case 0xB6: // OR (IX+d)
		return DdOrAIXd, 19, 3
	case 0xBE: // CP (IX+d)
		return DdCpAIXd, 19, 3
	}

	return nil, 0, 0
}

// decodeDDStackInstructions handles DD stack and jump instructions.
func (cpu *CPU) decodeDDStackInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// JP (IX)
	case 0xE9:
		return DdJpIX, 8, 2

	// EX (SP),IX
	case 0xE3:
		return DdExSpIX, 23, 2

	// PUSH IX / POP IX
	case 0xE5:
		return DdPushIX, 15, 2
	case 0xE1:
		return DdPopIX, 14, 2
	}

	return nil, 0, 0
}

// decodeDDCBInstruction decodes DD CB prefixed instructions (IX bit operations).
func (cpu *CPU) decodeDDCBInstruction() (Opcode, uint8, error) {
	displacement := int8(cpu.memory.Read(cpu.PC + 2)) // Get displacement
	opcodeByte := cpu.memory.Read(cpu.PC + 3)         // Get bit operation

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
	opcodeByte := cpu.memory.Read(cpu.PC + 1) // Get the actual FD instruction

	// Handle FD CB prefix first
	if opcodeByte == PrefixCB {
		return cpu.decodeFDCBInstruction()
	}

	instruction, timing, size, err := cpu.decodeFDInstructionType(opcodeByte)
	if err != nil {
		// Handle undocumented behavior: FD prefix alone acts as 4-cycle NOP
		// This occurs when FD is followed by an invalid IY instruction
		opcode := Opcode{
			Instruction: Nop,
			Addressing:  ImpliedAddressing,
			Timing:      4,
			Size:        1,
		}
		// Return nil error since we're handling it as undocumented NOP
		return opcode, PrefixFD, nil //nolint:nilerr
	}

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        size,
		Timing:      timing,
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

// decodeFDInstructionType determines the instruction, timing, and size for FD-prefixed opcodes.
func (cpu *CPU) decodeFDInstructionType(opcodeByte uint8) (*Instruction, byte, byte, error) {
	// Group instructions by functionality to reduce complexity
	if instruction, timing, size := cpu.decodeFDBasicInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeFDLoadInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeFDArithmeticInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := cpu.decodeFDStackInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	return nil, 0, 0, fmt.Errorf("%w: opcode FD %02X", ErrUnsupportedFDOpcode, opcodeByte)
}

// decodeFDBasicInstructions handles basic FD instructions (INC/DEC IY, ADD IY,rr).
func (cpu *CPU) decodeFDBasicInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// INC IY / DEC IY
	case 0x23:
		return FdIncIY, 10, 2
	case 0x2B:
		return FdDecIY, 10, 2

	// ADD IY,rr
	case 0x09: // ADD IY,BC
		return FdAddIYBc, 15, 2
	case 0x19: // ADD IY,DE
		return FdAddIYDe, 15, 2
	case 0x29: // ADD IY,IY
		return FdAddIYIY, 15, 2
	case 0x39: // ADD IY,SP
		return FdAddIYSp, 15, 2

	// INC/DEC (IY+d)
	case 0x34: // INC (IY+d)
		return FdIncIYd, 23, 3
	case 0x35: // DEC (IY+d)
		return FdDecIYd, 23, 3
	}

	return nil, 0, 0
}

// decodeFDLoadInstructions handles FD load instructions.
func (cpu *CPU) decodeFDLoadInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// LD IY,nn
	case 0x21:
		return FdLdIYnn, 14, 4

	// LD (nn),IY / LD IY,(nn)
	case 0x22:
		return FdLdNnIY, 20, 4
	case 0x2A:
		return FdLdIYNn, 20, 4

	// LD (IY+d),n
	case 0x36:
		return FdLdIYdN, 19, 4
	}

	// LD r,(IY+d) - Load register from (IY+d)
	if instruction, timing, size := cpu.decodeFDLoadFromIY(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	// LD (IY+d),r - Load (IY+d) from register
	if instruction, timing, size := cpu.decodeFDLoadToIY(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	return nil, 0, 0
}

// decodeFDLoadFromIY handles LD r,(IY+d) instructions.
func (cpu *CPU) decodeFDLoadFromIY(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x46: // LD B,(IY+d)
		return FdLdBIYd, 19, 3
	case 0x4E: // LD C,(IY+d)
		return FdLdCIYd, 19, 3
	case 0x56: // LD D,(IY+d)
		return FdLdDIYd, 19, 3
	case 0x5E: // LD E,(IY+d)
		return FdLdEIYd, 19, 3
	case 0x66: // LD H,(IY+d)
		return FdLdHIYd, 19, 3
	case 0x6E: // LD L,(IY+d)
		return FdLdLIYd, 19, 3
	case 0x7E: // LD A,(IY+d)
		return FdLdAIYd, 19, 3
	}

	return nil, 0, 0
}

// decodeFDLoadToIY handles LD (IY+d),r instructions.
func (cpu *CPU) decodeFDLoadToIY(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x70: // LD (IY+d),B
		return FdLdIYdB, 19, 3
	case 0x71: // LD (IY+d),C
		return FdLdIYdC, 19, 3
	case 0x72: // LD (IY+d),D
		return FdLdIYdD, 19, 3
	case 0x73: // LD (IY+d),E
		return FdLdIYdE, 19, 3
	case 0x74: // LD (IY+d),H
		return FdLdIYdH, 19, 3
	case 0x75: // LD (IY+d),L
		return FdLdIYdL, 19, 3
	case 0x77: // LD (IY+d),A
		return FdLdIYdA, 19, 3
	}

	return nil, 0, 0
}

// decodeFDArithmeticInstructions handles FD arithmetic instructions with (IY+d).
func (cpu *CPU) decodeFDArithmeticInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x86: // ADD A,(IY+d)
		return FdAddAIYd, 19, 3
	case 0x8E: // ADC A,(IY+d)
		return FdAdcAIYd, 19, 3
	case 0x96: // SUB (IY+d)
		return FdSubAIYd, 19, 3
	case 0x9E: // SBC A,(IY+d)
		return FdSbcAIYd, 19, 3
	case 0xA6: // AND (IY+d)
		return FdAndAIYd, 19, 3
	case 0xAE: // XOR (IY+d)
		return FdXorAIYd, 19, 3
	case 0xB6: // OR (IY+d)
		return FdOrAIYd, 19, 3
	case 0xBE: // CP (IY+d)
		return FdCpAIYd, 19, 3
	}

	return nil, 0, 0
}

// decodeFDStackInstructions handles FD stack and jump instructions.
func (cpu *CPU) decodeFDStackInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// JP (IY)
	case 0xE9:
		return FdJpIY, 8, 2

	// EX (SP),IY
	case 0xE3:
		return FdExSpIY, 23, 2

	// PUSH IY / POP IY
	case 0xE5:
		return FdPushIY, 15, 2
	case 0xE1:
		return FdPopIY, 14, 2
	}

	return nil, 0, 0
}

// decodeFDCBInstruction decodes FD CB prefixed instructions (IY bit operations).
func (cpu *CPU) decodeFDCBInstruction() (Opcode, uint8, error) {
	displacement := int8(cpu.memory.Read(cpu.PC + 2)) // Get displacement
	opcodeByte := cpu.memory.Read(cpu.PC + 3)         // Get bit operation

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

// handleInterrupts processes pending interrupts.
func (cpu *CPU) handleInterrupts() error {
	// Non-maskable interrupt has highest priority
	if cpu.triggerNmi {
		cpu.triggerNmi = false
		cpu.halted = false

		// Save current PC
		cpu.push16(cpu.PC)

		// Jump to NMI vector
		cpu.PC = 0x0066
		cpu.iff2 = cpu.iff1
		cpu.iff1 = false

		cpu.cycles += 11
		return nil
	}

	// Maskable interrupt
	if cpu.triggerIrq && cpu.iff1 {
		cpu.triggerIrq = false
		cpu.halted = false
		cpu.iff1 = false
		cpu.iff2 = false

		// Save current PC
		cpu.push16(cpu.PC)

		switch cpu.im {
		case 0:
			// Interrupt mode 0: Execute instruction on data bus (usually RST)
			cpu.PC = 0x0040
			cpu.cycles += 13
		case 1:
			// Interrupt mode 1: Jump to 0x0038
			cpu.PC = 0x0038
			cpu.cycles += 13
		case 2:
			// Interrupt mode 2: Vector table lookup
			vector := uint16(cpu.I)<<8 | uint16(cpu.memory.Read(0xFFFF))
			cpu.PC = cpu.memory.ReadWord(vector)
			cpu.cycles += 19
		}

		return nil
	}

	return nil
}
