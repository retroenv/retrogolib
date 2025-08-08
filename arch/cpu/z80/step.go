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
		return c.decodeFDInstruction()
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
			OpcodeOperands: []byte{0xCB, opcodeByte},
		}
	}

	return opcode, 0xCB, nil
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
	timing := byte(8)

	switch {
	case opcodeByte <= 0x07: // RLC r
		instruction = &Instruction{Name: "rlc", ParamFunc: cbRlc}
	case opcodeByte <= 0x0F: // RRC r
		instruction = &Instruction{Name: "rrc", ParamFunc: cbRrc}
	case opcodeByte <= 0x17: // RL r
		instruction = &Instruction{Name: "rl", ParamFunc: cbRl}
	case opcodeByte <= 0x1F: // RR r
		instruction = &Instruction{Name: "rr", ParamFunc: cbRr}
	case opcodeByte <= 0x27: // SLA r
		instruction = &Instruction{Name: "sla", ParamFunc: cbSla}
	case opcodeByte <= 0x2F: // SRA r
		instruction = &Instruction{Name: "sra", ParamFunc: cbSra}
	case opcodeByte <= 0x37: // SLL r (undocumented)
		instruction = &Instruction{Name: "sll", ParamFunc: cbSll}
	default: // SRL r (0x38-0x3F)
		instruction = &Instruction{Name: "srl", ParamFunc: cbSrl}
	}

	// Special timing for (HL) operations
	if reg == 6 {
		timing = 15
	}

	return instruction, timing
}

// decodeCBBit handles CB BIT instructions (0x40-0x7F).
func (c *CPU) decodeCBBit(reg uint8) (*Instruction, byte) {
	instruction := &Instruction{Name: "bit", ParamFunc: cbBit}
	timing := byte(8)

	if reg == 6 { // BIT n,(HL)
		timing = 12
	}

	return instruction, timing
}

// decodeCBRes handles CB RES instructions (0x80-0xBF).
func (c *CPU) decodeCBRes(reg uint8) (*Instruction, byte) {
	instruction := &Instruction{Name: "res", ParamFunc: cbRes}
	timing := byte(8)

	if reg == 6 { // RES n,(HL)
		timing = 15
	}

	return instruction, timing
}

// decodeCBSet handles CB SET instructions (0xC0-0xFF).
func (c *CPU) decodeCBSet(reg uint8) (*Instruction, byte) {
	instruction := &Instruction{Name: "set", ParamFunc: cbSet}
	timing := byte(8)

	if reg == 6 { // SET n,(HL)
		timing = 15
	}

	return instruction, timing
}

// decodeEDInstruction decodes ED-prefixed instructions (extended operations).
func (c *CPU) decodeEDInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual ED instruction

	instruction, timing, size, err := c.decodeEDInstructionType(opcodeByte)
	if err != nil {
		return Opcode{}, 0xED, err
	}

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        size,
		Timing:      timing,
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{0xED, opcodeByte},
		}
	}

	return opcode, 0xED, nil
}

// decodeEDInstructionType determines the instruction, timing, and size for ED-prefixed opcodes.
func (c *CPU) decodeEDInstructionType(opcodeByte uint8) (*Instruction, byte, byte, error) {
	// Group instructions by functionality to reduce complexity
	if instruction, timing, size := c.decodeEDBasicInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := c.decodeEDArithmeticInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := c.decodeEDLoadInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := c.decodeEDBlockInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := c.decodeEDIOInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	return nil, 0, 0, fmt.Errorf("unimplemented ED instruction: ED %02X", opcodeByte)
}

// decodeEDBasicInstructions handles basic ED instructions (NEG, IM, RETN, RETI, RRD, RLD).
func (c *CPU) decodeEDBasicInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// NEG - Negate Accumulator
	case 0x44, 0x4C, 0x54, 0x5C, 0x64, 0x6C, 0x74, 0x7C:
		return &Instruction{Name: "neg", NoParamFunc: edNeg}, 8, 2

	// IM - Set Interrupt Mode
	case 0x46, 0x66: // IM 0
		return &Instruction{Name: "im", ParamFunc: edIm0}, 8, 2
	case 0x56, 0x76: // IM 1
		return &Instruction{Name: "im", ParamFunc: edIm1}, 8, 2
	case 0x5E, 0x7E: // IM 2
		return &Instruction{Name: "im", ParamFunc: edIm2}, 8, 2

	// RETN/RETI
	case 0x45, 0x55, 0x65, 0x75: // RETN (multiple opcodes)
		return &Instruction{Name: "retn", NoParamFunc: edRetn}, 14, 2
	case 0x4D: // RETI
		return &Instruction{Name: "reti", NoParamFunc: edReti}, 14, 2

	// RRD/RLD
	case 0x67: // RRD
		return &Instruction{Name: "rrd", NoParamFunc: edRrd}, 18, 2
	case 0x6F: // RLD
		return &Instruction{Name: "rld", NoParamFunc: edRld}, 18, 2
	}

	return nil, 0, 0
}

// decodeEDArithmeticInstructions handles ED arithmetic instructions (ADC HL,rr / SBC HL,rr).
func (c *CPU) decodeEDArithmeticInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// ADC HL,rr
	case 0x4A: // ADC HL,BC
		return &Instruction{Name: "adc", ParamFunc: edAdcHlBc}, 15, 2
	case 0x5A: // ADC HL,DE
		return &Instruction{Name: "adc", ParamFunc: edAdcHlDe}, 15, 2
	case 0x6A: // ADC HL,HL
		return &Instruction{Name: "adc", ParamFunc: edAdcHlHl}, 15, 2
	case 0x7A: // ADC HL,SP
		return &Instruction{Name: "adc", ParamFunc: edAdcHlSp}, 15, 2

	// SBC HL,rr
	case 0x42: // SBC HL,BC
		return &Instruction{Name: "sbc", ParamFunc: edSbcHlBc}, 15, 2
	case 0x52: // SBC HL,DE
		return &Instruction{Name: "sbc", ParamFunc: edSbcHlDe}, 15, 2
	case 0x62: // SBC HL,HL
		return &Instruction{Name: "sbc", ParamFunc: edSbcHlHl}, 15, 2
	case 0x72: // SBC HL,SP
		return &Instruction{Name: "sbc", ParamFunc: edSbcHlSp}, 15, 2
	}

	return nil, 0, 0
}

// decodeEDLoadInstructions handles ED load instructions.
func (c *CPU) decodeEDLoadInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// LD I,A / LD R,A
	case 0x47: // LD I,A
		return &Instruction{Name: "ld", NoParamFunc: edLdIA}, 9, 2
	case 0x4F: // LD R,A
		return &Instruction{Name: "ld", NoParamFunc: edLdRA}, 9, 2

	// LD A,I / LD A,R
	case 0x57: // LD A,I
		return &Instruction{Name: "ld", NoParamFunc: edLdAI}, 9, 2
	case 0x5F: // LD A,R
		return &Instruction{Name: "ld", NoParamFunc: edLdAR}, 9, 2

	// LD (nn),rr
	case 0x43: // LD (nn),BC
		return &Instruction{Name: "ld", ParamFunc: edLdNnBc}, 20, 4
	case 0x53: // LD (nn),DE
		return &Instruction{Name: "ld", ParamFunc: edLdNnDe}, 20, 4
	case 0x63: // LD (nn),HL
		return &Instruction{Name: "ld", ParamFunc: edLdNnHl}, 20, 4
	case 0x73: // LD (nn),SP
		return &Instruction{Name: "ld", ParamFunc: edLdNnSp}, 20, 4

	// LD rr,(nn)
	case 0x4B: // LD BC,(nn)
		return &Instruction{Name: "ld", ParamFunc: edLdBcNn}, 20, 4
	case 0x5B: // LD DE,(nn)
		return &Instruction{Name: "ld", ParamFunc: edLdDeNn}, 20, 4
	case 0x6B: // LD HL,(nn)
		return &Instruction{Name: "ld", ParamFunc: edLdHlNn}, 20, 4
	case 0x7B: // LD SP,(nn)
		return &Instruction{Name: "ld", ParamFunc: edLdSpNn}, 20, 4
	}

	return nil, 0, 0
}

// decodeEDBlockInstructions handles ED block transfer and search instructions.
func (c *CPU) decodeEDBlockInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// Block transfer instructions
	case 0xA0: // LDI
		return &Instruction{Name: "ldi", NoParamFunc: edLdi}, 16, 2
	case 0xA8: // LDD
		return &Instruction{Name: "ldd", NoParamFunc: edLdd}, 16, 2
	case 0xB0: // LDIR
		return &Instruction{Name: "ldir", NoParamFunc: edLdir}, 21, 2
	case 0xB8: // LDDR
		return &Instruction{Name: "lddr", NoParamFunc: edLddr}, 21, 2

	// Block search instructions
	case 0xA1: // CPI
		return &Instruction{Name: "cpi", NoParamFunc: edCpi}, 16, 2
	case 0xA9: // CPD
		return &Instruction{Name: "cpd", NoParamFunc: edCpd}, 16, 2
	case 0xB1: // CPIR
		return &Instruction{Name: "cpir", NoParamFunc: edCpir}, 21, 2
	case 0xB9: // CPDR
		return &Instruction{Name: "cpdr", NoParamFunc: edCpdr}, 21, 2
	}

	return nil, 0, 0
}

// decodeEDIOInstructions handles ED I/O instructions.
func (c *CPU) decodeEDIOInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	// Block I/O instructions
	if instruction, timing, size := c.decodeEDBlockIO(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	// I/O with register - IN r,(C)
	if instruction, timing, size := c.decodeEDInputInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	// I/O with register - OUT (C),r
	if instruction, timing, size := c.decodeEDOutputInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	return nil, 0, 0
}

// decodeEDBlockIO handles ED block I/O instructions.
func (c *CPU) decodeEDBlockIO(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0xA2: // INI
		return &Instruction{Name: "ini", NoParamFunc: edIni}, 16, 2
	case 0xAA: // IND
		return &Instruction{Name: "ind", NoParamFunc: edInd}, 16, 2
	case 0xB2: // INIR
		return &Instruction{Name: "inir", NoParamFunc: edInir}, 21, 2
	case 0xBA: // INDR
		return &Instruction{Name: "indr", NoParamFunc: edIndr}, 21, 2
	case 0xA3: // OUTI
		return &Instruction{Name: "outi", NoParamFunc: edOuti}, 16, 2
	case 0xAB: // OUTD
		return &Instruction{Name: "outd", NoParamFunc: edOutd}, 16, 2
	case 0xB3: // OTIR
		return &Instruction{Name: "otir", NoParamFunc: edOtir}, 21, 2
	case 0xBB: // OTDR
		return &Instruction{Name: "otdr", NoParamFunc: edOtdr}, 21, 2
	}

	return nil, 0, 0
}

// decodeEDInputInstructions handles ED IN r,(C) instructions.
func (c *CPU) decodeEDInputInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x40: // IN B,(C)
		return &Instruction{Name: "in", ParamFunc: edInBC}, 12, 2
	case 0x48: // IN C,(C)
		return &Instruction{Name: "in", ParamFunc: edInCC}, 12, 2
	case 0x50: // IN D,(C)
		return &Instruction{Name: "in", ParamFunc: edInDC}, 12, 2
	case 0x58: // IN E,(C)
		return &Instruction{Name: "in", ParamFunc: edInEC}, 12, 2
	case 0x60: // IN H,(C)
		return &Instruction{Name: "in", ParamFunc: edInHC}, 12, 2
	case 0x68: // IN L,(C)
		return &Instruction{Name: "in", ParamFunc: edInLC}, 12, 2
	case 0x78: // IN A,(C)
		return &Instruction{Name: "in", ParamFunc: edInAC}, 12, 2
	}

	return nil, 0, 0
}

// decodeEDOutputInstructions handles ED OUT (C),r instructions.
func (c *CPU) decodeEDOutputInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x41: // OUT (C),B
		return &Instruction{Name: "out", ParamFunc: edOutCB}, 12, 2
	case 0x49: // OUT (C),C
		return &Instruction{Name: "out", ParamFunc: edOutCC}, 12, 2
	case 0x51: // OUT (C),D
		return &Instruction{Name: "out", ParamFunc: edOutCD}, 12, 2
	case 0x59: // OUT (C),E
		return &Instruction{Name: "out", ParamFunc: edOutCE}, 12, 2
	case 0x61: // OUT (C),H
		return &Instruction{Name: "out", ParamFunc: edOutCH}, 12, 2
	case 0x69: // OUT (C),L
		return &Instruction{Name: "out", ParamFunc: edOutCL}, 12, 2
	case 0x79: // OUT (C),A
		return &Instruction{Name: "out", ParamFunc: edOutCA}, 12, 2
	}

	return nil, 0, 0
}

// decodeDDInstruction decodes DD-prefixed instructions (IX operations).
func (c *CPU) decodeDDInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual DD instruction

	// Handle DD CB prefix first
	if opcodeByte == 0xCB {
		return c.decodeDDCBInstruction()
	}

	instruction, timing, size, err := c.decodeDDInstructionType(opcodeByte)
	if err != nil {
		return Opcode{}, 0xDD, err
	}

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        size,
		Timing:      timing,
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{0xDD, opcodeByte},
		}
	}

	return opcode, 0xDD, nil
}

// decodeDDInstructionType determines the instruction, timing, and size for DD-prefixed opcodes.
func (c *CPU) decodeDDInstructionType(opcodeByte uint8) (*Instruction, byte, byte, error) {
	// Group instructions by functionality to reduce complexity
	if instruction, timing, size := c.decodeDDBasicInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := c.decodeDDLoadInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := c.decodeDDArithmeticInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := c.decodeDDStackInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	return nil, 0, 0, fmt.Errorf("unimplemented DD instruction: DD %02X", opcodeByte)
}

// decodeDDBasicInstructions handles basic DD instructions (INC/DEC IX, ADD IX,rr).
func (c *CPU) decodeDDBasicInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// INC IX / DEC IX
	case 0x23:
		return &Instruction{Name: "inc", NoParamFunc: ddIncIX}, 10, 2
	case 0x2B:
		return &Instruction{Name: "dec", NoParamFunc: ddDecIX}, 10, 2

	// ADD IX,rr
	case 0x09: // ADD IX,BC
		return &Instruction{Name: "add", ParamFunc: ddAddIXBc}, 15, 2
	case 0x19: // ADD IX,DE
		return &Instruction{Name: "add", ParamFunc: ddAddIXDe}, 15, 2
	case 0x29: // ADD IX,IX
		return &Instruction{Name: "add", ParamFunc: ddAddIXIX}, 15, 2
	case 0x39: // ADD IX,SP
		return &Instruction{Name: "add", ParamFunc: ddAddIXSp}, 15, 2

	// INC/DEC (IX+d)
	case 0x34: // INC (IX+d)
		return &Instruction{Name: "inc", ParamFunc: ddIncIXd}, 23, 3
	case 0x35: // DEC (IX+d)
		return &Instruction{Name: "dec", ParamFunc: ddDecIXd}, 23, 3
	}

	return nil, 0, 0
}

// decodeDDLoadInstructions handles DD load instructions.
func (c *CPU) decodeDDLoadInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// LD IX,nn
	case 0x21:
		return &Instruction{Name: "ld", ParamFunc: ddLdIXnn}, 14, 4

	// LD (nn),IX / LD IX,(nn)
	case 0x22:
		return &Instruction{Name: "ld", ParamFunc: ddLdNnIX}, 20, 4
	case 0x2A:
		return &Instruction{Name: "ld", ParamFunc: ddLdIXNn}, 20, 4

	// LD (IX+d),n
	case 0x36:
		return &Instruction{Name: "ld", ParamFunc: ddLdIXdN}, 19, 4
	}

	// LD r,(IX+d) - Load register from (IX+d)
	if instruction, timing, size := c.decodeDDLoadFromIX(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	// LD (IX+d),r - Load (IX+d) from register
	if instruction, timing, size := c.decodeDDLoadToIX(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	return nil, 0, 0
}

// decodeDDLoadFromIX handles LD r,(IX+d) instructions.
func (c *CPU) decodeDDLoadFromIX(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x46: // LD B,(IX+d)
		return &Instruction{Name: "ld", ParamFunc: ddLdBIXd}, 19, 3
	case 0x4E: // LD C,(IX+d)
		return &Instruction{Name: "ld", ParamFunc: ddLdCIXd}, 19, 3
	case 0x56: // LD D,(IX+d)
		return &Instruction{Name: "ld", ParamFunc: ddLdDIXd}, 19, 3
	case 0x5E: // LD E,(IX+d)
		return &Instruction{Name: "ld", ParamFunc: ddLdEIXd}, 19, 3
	case 0x66: // LD H,(IX+d)
		return &Instruction{Name: "ld", ParamFunc: ddLdHIXd}, 19, 3
	case 0x6E: // LD L,(IX+d)
		return &Instruction{Name: "ld", ParamFunc: ddLdLIXd}, 19, 3
	case 0x7E: // LD A,(IX+d)
		return &Instruction{Name: "ld", ParamFunc: ddLdAIXd}, 19, 3
	}

	return nil, 0, 0
}

// decodeDDLoadToIX handles LD (IX+d),r instructions.
func (c *CPU) decodeDDLoadToIX(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x70: // LD (IX+d),B
		return &Instruction{Name: "ld", ParamFunc: ddLdIXdB}, 19, 3
	case 0x71: // LD (IX+d),C
		return &Instruction{Name: "ld", ParamFunc: ddLdIXdC}, 19, 3
	case 0x72: // LD (IX+d),D
		return &Instruction{Name: "ld", ParamFunc: ddLdIXdD}, 19, 3
	case 0x73: // LD (IX+d),E
		return &Instruction{Name: "ld", ParamFunc: ddLdIXdE}, 19, 3
	case 0x74: // LD (IX+d),H
		return &Instruction{Name: "ld", ParamFunc: ddLdIXdH}, 19, 3
	case 0x75: // LD (IX+d),L
		return &Instruction{Name: "ld", ParamFunc: ddLdIXdL}, 19, 3
	case 0x77: // LD (IX+d),A
		return &Instruction{Name: "ld", ParamFunc: ddLdIXdA}, 19, 3
	}

	return nil, 0, 0
}

// decodeDDArithmeticInstructions handles DD arithmetic instructions with (IX+d).
func (c *CPU) decodeDDArithmeticInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x86: // ADD A,(IX+d)
		return &Instruction{Name: "add", ParamFunc: ddAddAIXd}, 19, 3
	case 0x8E: // ADC A,(IX+d)
		return &Instruction{Name: "adc", ParamFunc: ddAdcAIXd}, 19, 3
	case 0x96: // SUB (IX+d)
		return &Instruction{Name: "sub", ParamFunc: ddSubAIXd}, 19, 3
	case 0x9E: // SBC A,(IX+d)
		return &Instruction{Name: "sbc", ParamFunc: ddSbcAIXd}, 19, 3
	case 0xA6: // AND (IX+d)
		return &Instruction{Name: "and", ParamFunc: ddAndAIXd}, 19, 3
	case 0xAE: // XOR (IX+d)
		return &Instruction{Name: "xor", ParamFunc: ddXorAIXd}, 19, 3
	case 0xB6: // OR (IX+d)
		return &Instruction{Name: "or", ParamFunc: ddOrAIXd}, 19, 3
	case 0xBE: // CP (IX+d)
		return &Instruction{Name: "cp", ParamFunc: ddCpAIXd}, 19, 3
	}

	return nil, 0, 0
}

// decodeDDStackInstructions handles DD stack and jump instructions.
func (c *CPU) decodeDDStackInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// JP (IX)
	case 0xE9:
		return &Instruction{Name: "jp", NoParamFunc: ddJpIX}, 8, 2

	// EX (SP),IX
	case 0xE3:
		return &Instruction{Name: "ex", NoParamFunc: ddExSpIX}, 23, 2

	// PUSH IX / POP IX
	case 0xE5:
		return &Instruction{Name: "push", NoParamFunc: ddPushIX}, 15, 2
	case 0xE1:
		return &Instruction{Name: "pop", NoParamFunc: ddPopIX}, 14, 2
	}

	return nil, 0, 0
}

// decodeDDCBInstruction decodes DD CB prefixed instructions (IX bit operations).
func (c *CPU) decodeDDCBInstruction() (Opcode, uint8, error) {
	displacement := int8(c.memory.Read(c.PC + 2)) // Get displacement
	opcodeByte := c.memory.Read(c.PC + 3)         // Get bit operation

	var instruction *Instruction
	var timing byte = 23 // All DDCB operations take 23 T-states

	switch {
	case opcodeByte <= 0x3F: // Rotate/shift operations
		instruction = &Instruction{Name: "ddcb-shift", ParamFunc: ddcbShift}
	case opcodeByte <= 0x7F: // BIT operations
		instruction = &Instruction{Name: "bit", ParamFunc: ddcbBit}
	case opcodeByte <= 0xBF: // RES operations
		instruction = &Instruction{Name: "res", ParamFunc: ddcbRes}
	default: // SET operations (0xC0-0xFF)
		instruction = &Instruction{Name: "set", ParamFunc: ddcbSet}
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
			OpcodeOperands: []byte{0xDD, 0xCB, uint8(displacement), opcodeByte},
		}
	}

	return opcode, 0xDD, nil
}

// decodeFDInstruction decodes FD-prefixed instructions (IY operations).
func (c *CPU) decodeFDInstruction() (Opcode, uint8, error) {
	opcodeByte := c.memory.Read(c.PC + 1) // Get the actual FD instruction

	// Handle FD CB prefix first
	if opcodeByte == 0xCB {
		return c.decodeFDCBInstruction()
	}

	instruction, timing, size, err := c.decodeFDInstructionType(opcodeByte)
	if err != nil {
		return Opcode{}, 0xFD, err
	}

	opcode := Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Size:        size,
		Timing:      timing,
	}

	if c.opts.tracing {
		c.TraceStep = TraceStep{
			PC:             c.PC,
			Opcode:         opcode,
			OpcodeOperands: []byte{0xFD, opcodeByte},
		}
	}

	return opcode, 0xFD, nil
}

// decodeFDInstructionType determines the instruction, timing, and size for FD-prefixed opcodes.
func (c *CPU) decodeFDInstructionType(opcodeByte uint8) (*Instruction, byte, byte, error) {
	// Group instructions by functionality to reduce complexity
	if instruction, timing, size := c.decodeFDBasicInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := c.decodeFDLoadInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	if instruction, timing, size := c.decodeFDStackInstructions(opcodeByte); instruction != nil {
		return instruction, timing, size, nil
	}

	return nil, 0, 0, fmt.Errorf("unimplemented FD instruction: FD %02X", opcodeByte)
}

// decodeFDBasicInstructions handles basic FD instructions (INC/DEC IY, ADD IY,rr).
func (c *CPU) decodeFDBasicInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// INC IY / DEC IY
	case 0x23:
		return &Instruction{Name: "inc", NoParamFunc: fdIncIY}, 10, 2
	case 0x2B:
		return &Instruction{Name: "dec", NoParamFunc: fdDecIY}, 10, 2

	// ADD IY,rr
	case 0x09: // ADD IY,BC
		return &Instruction{Name: "add", ParamFunc: fdAddIYBc}, 15, 2
	case 0x19: // ADD IY,DE
		return &Instruction{Name: "add", ParamFunc: fdAddIYDe}, 15, 2
	case 0x29: // ADD IY,IY
		return &Instruction{Name: "add", ParamFunc: fdAddIYIY}, 15, 2
	case 0x39: // ADD IY,SP
		return &Instruction{Name: "add", ParamFunc: fdAddIYSp}, 15, 2

	// INC/DEC (IY+d)
	case 0x34: // INC (IY+d)
		return &Instruction{Name: "inc", ParamFunc: fdIncIYd}, 23, 3
	case 0x35: // DEC (IY+d)
		return &Instruction{Name: "dec", ParamFunc: fdDecIYd}, 23, 3
	}

	return nil, 0, 0
}

// decodeFDLoadInstructions handles FD load instructions.
func (c *CPU) decodeFDLoadInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// LD IY,nn
	case 0x21:
		return &Instruction{Name: "ld", ParamFunc: fdLdIYnn}, 14, 4

	// LD (nn),IY / LD IY,(nn)
	case 0x22:
		return &Instruction{Name: "ld", ParamFunc: fdLdNnIY}, 20, 4
	case 0x2A:
		return &Instruction{Name: "ld", ParamFunc: fdLdIYNn}, 20, 4

	// LD (IY+d),n
	case 0x36:
		return &Instruction{Name: "ld", ParamFunc: fdLdIYdN}, 19, 4
	}

	// LD r,(IY+d) - Load register from (IY+d)
	if instruction, timing, size := c.decodeFDLoadFromIY(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	// LD (IY+d),r - Load (IY+d) from register
	if instruction, timing, size := c.decodeFDLoadToIY(opcodeByte); instruction != nil {
		return instruction, timing, size
	}

	return nil, 0, 0
}

// decodeFDLoadFromIY handles LD r,(IY+d) instructions.
func (c *CPU) decodeFDLoadFromIY(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x46: // LD B,(IY+d)
		return &Instruction{Name: "ld", ParamFunc: fdLdBIYd}, 19, 3
	case 0x4E: // LD C,(IY+d)
		return &Instruction{Name: "ld", ParamFunc: fdLdCIYd}, 19, 3
	case 0x56: // LD D,(IY+d)
		return &Instruction{Name: "ld", ParamFunc: fdLdDIYd}, 19, 3
	case 0x5E: // LD E,(IY+d)
		return &Instruction{Name: "ld", ParamFunc: fdLdEIYd}, 19, 3
	case 0x66: // LD H,(IY+d)
		return &Instruction{Name: "ld", ParamFunc: fdLdHIYd}, 19, 3
	case 0x6E: // LD L,(IY+d)
		return &Instruction{Name: "ld", ParamFunc: fdLdLIYd}, 19, 3
	case 0x7E: // LD A,(IY+d)
		return &Instruction{Name: "ld", ParamFunc: fdLdAIYd}, 19, 3
	}

	return nil, 0, 0
}

// decodeFDLoadToIY handles LD (IY+d),r instructions.
func (c *CPU) decodeFDLoadToIY(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	case 0x70: // LD (IY+d),B
		return &Instruction{Name: "ld", ParamFunc: fdLdIYdB}, 19, 3
	case 0x71: // LD (IY+d),C
		return &Instruction{Name: "ld", ParamFunc: fdLdIYdC}, 19, 3
	case 0x72: // LD (IY+d),D
		return &Instruction{Name: "ld", ParamFunc: fdLdIYdD}, 19, 3
	case 0x73: // LD (IY+d),E
		return &Instruction{Name: "ld", ParamFunc: fdLdIYdE}, 19, 3
	case 0x74: // LD (IY+d),H
		return &Instruction{Name: "ld", ParamFunc: fdLdIYdH}, 19, 3
	case 0x75: // LD (IY+d),L
		return &Instruction{Name: "ld", ParamFunc: fdLdIYdL}, 19, 3
	case 0x77: // LD (IY+d),A
		return &Instruction{Name: "ld", ParamFunc: fdLdIYdA}, 19, 3
	}

	return nil, 0, 0
}

// decodeFDStackInstructions handles FD stack and jump instructions.
func (c *CPU) decodeFDStackInstructions(opcodeByte uint8) (*Instruction, byte, byte) {
	switch opcodeByte {
	// JP (IY)
	case 0xE9:
		return &Instruction{Name: "jp", NoParamFunc: fdJpIY}, 8, 2

	// EX (SP),IY
	case 0xE3:
		return &Instruction{Name: "ex", NoParamFunc: fdExSpIY}, 23, 2

	// PUSH IY / POP IY
	case 0xE5:
		return &Instruction{Name: "push", NoParamFunc: fdPushIY}, 15, 2
	case 0xE1:
		return &Instruction{Name: "pop", NoParamFunc: fdPopIY}, 14, 2
	}

	return nil, 0, 0
}

// decodeFDCBInstruction decodes FD CB prefixed instructions (IY bit operations).
func (c *CPU) decodeFDCBInstruction() (Opcode, uint8, error) {
	displacement := int8(c.memory.Read(c.PC + 2)) // Get displacement
	opcodeByte := c.memory.Read(c.PC + 3)         // Get bit operation

	var instruction *Instruction
	var timing byte = 23 // All FDCB operations take 23 T-states

	switch {
	case opcodeByte <= 0x3F: // Rotate/shift operations
		instruction = &Instruction{Name: "fdcb-shift", ParamFunc: fdcbShift}
	case opcodeByte <= 0x7F: // BIT operations
		instruction = &Instruction{Name: "bit", ParamFunc: fdcbBit}
	case opcodeByte <= 0xBF: // RES operations
		instruction = &Instruction{Name: "res", ParamFunc: fdcbRes}
	default: // SET operations (0xC0-0xFF)
		instruction = &Instruction{Name: "set", ParamFunc: fdcbSet}
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
			OpcodeOperands: []byte{0xFD, 0xCB, uint8(displacement), opcodeByte},
		}
	}

	return opcode, 0xFD, nil
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
