package z80

// Instruction contains information about a Z80 CPU instruction.
type Instruction struct {
	Name       string // lowercased instruction name
	Unofficial bool   // unofficial instructions are not part of the original Z80 spec

	Addressing      map[AddressingMode]OpcodeInfo // addressing mode mapping to opcode info
	RegisterOpcodes map[RegisterParam]OpcodeInfo  // register-specific opcode mapping for disambiguating variants

	NoParamFunc func(c *CPU) error                // emulation function to execute when the instruction has no parameters
	ParamFunc   func(c *CPU, params ...any) error // emulation function to execute when the instruction has parameters
}

// HasAddressing returns whether the instruction has any of the passed addressing modes.
func (ins Instruction) HasAddressing(flags ...AddressingMode) bool {
	for _, flag := range flags {
		_, ok := ins.Addressing[flag]
		if ok {
			return ok
		}
	}
	return false
}

// GetOpcodeByRegister returns opcode info for a specific register parameter.
func (ins Instruction) GetOpcodeByRegister(register RegisterParam) (OpcodeInfo, bool) {
	if ins.RegisterOpcodes == nil {
		// Fall back to Addressing map if RegisterOpcodes is not defined
		for _, info := range ins.Addressing {
			return info, true
		}
		return OpcodeInfo{}, false
	}

	info, exists := ins.RegisterOpcodes[register]
	return info, exists
}

// GetAllRegisterVariants returns all register variants for this instruction.
func (ins Instruction) GetAllRegisterVariants() map[RegisterParam]OpcodeInfo {
	if ins.RegisterOpcodes == nil {
		return nil
	}

	variants := make(map[RegisterParam]OpcodeInfo)
	for reg, info := range ins.RegisterOpcodes {
		variants[reg] = info
	}
	return variants
}

// Nop - No Operation.
var Nop = &Instruction{
	Name: "nop",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x00, Size: 1, Cycles: 4},
	},
	NoParamFunc: nop,
}

// Halt - Halt execution.
var Halt = &Instruction{
	Name: "halt",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x76, Size: 1, Cycles: 4},
	},
	NoParamFunc: halt,
}

// LdImm8 - Load 8-bit immediate into register.
var LdImm8 = &Instruction{
	Name: "ld",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x3E, Size: 2, Cycles: 7}, // LD A,n (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Opcode: 0x06, Size: 2, Cycles: 7}, // LD B,n
		RegC: {Opcode: 0x0E, Size: 2, Cycles: 7}, // LD C,n
		RegD: {Opcode: 0x16, Size: 2, Cycles: 7}, // LD D,n
		RegE: {Opcode: 0x1E, Size: 2, Cycles: 7}, // LD E,n
		RegH: {Opcode: 0x26, Size: 2, Cycles: 7}, // LD H,n
		RegL: {Opcode: 0x2E, Size: 2, Cycles: 7}, // LD L,n
		RegA: {Opcode: 0x3E, Size: 2, Cycles: 7}, // LD A,n
	},
	ParamFunc: ldImm8,
}

// LdReg8 - Load between 8-bit registers.
var LdReg8 = &Instruction{
	Name: "ld",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x7F, Size: 1, Cycles: 4}, // LD A,A (base opcode, others calculated)
	},
	ParamFunc: ldReg8,
}

// IncReg8 - Increment 8-bit register.
var IncReg8 = &Instruction{
	Name: "inc",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x3C, Size: 1, Cycles: 4}, // INC A (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Opcode: 0x04, Size: 1, Cycles: 4}, // INC B
		RegC: {Opcode: 0x0C, Size: 1, Cycles: 4}, // INC C
		RegD: {Opcode: 0x14, Size: 1, Cycles: 4}, // INC D
		RegE: {Opcode: 0x1C, Size: 1, Cycles: 4}, // INC E
		RegH: {Opcode: 0x24, Size: 1, Cycles: 4}, // INC H
		RegL: {Opcode: 0x2C, Size: 1, Cycles: 4}, // INC L
		RegA: {Opcode: 0x3C, Size: 1, Cycles: 4}, // INC A
	},
	ParamFunc: incReg8,
}

// DecReg8 - Decrement 8-bit register.
var DecReg8 = &Instruction{
	Name: "dec",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x3D, Size: 1, Cycles: 4}, // DEC A (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Opcode: 0x05, Size: 1, Cycles: 4}, // DEC B
		RegC: {Opcode: 0x0D, Size: 1, Cycles: 4}, // DEC C
		RegD: {Opcode: 0x15, Size: 1, Cycles: 4}, // DEC D
		RegE: {Opcode: 0x1D, Size: 1, Cycles: 4}, // DEC E
		RegH: {Opcode: 0x25, Size: 1, Cycles: 4}, // DEC H
		RegL: {Opcode: 0x2D, Size: 1, Cycles: 4}, // DEC L
		RegA: {Opcode: 0x3D, Size: 1, Cycles: 4}, // DEC A
	},
	ParamFunc: decReg8,
}

// AddA - Add to accumulator.
var AddA = &Instruction{
	Name: "add",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0x87, Size: 1, Cycles: 4}, // ADD A,A (base opcode)
		ImmediateAddressing: {Opcode: 0xC6, Size: 2, Cycles: 7}, // ADD A,n
	},
	ParamFunc: addA,
}

// SubA - Subtract from accumulator.
var SubA = &Instruction{
	Name: "sub",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0x97, Size: 1, Cycles: 4}, // SUB A (base opcode)
		ImmediateAddressing: {Opcode: 0xD6, Size: 2, Cycles: 7}, // SUB n
	},
	ParamFunc: subA,
}

// AndA - AND with accumulator.
var AndA = &Instruction{
	Name: "and",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0xA7, Size: 1, Cycles: 4}, // AND A (base opcode)
		ImmediateAddressing: {Opcode: 0xE6, Size: 2, Cycles: 7}, // AND n
	},
	ParamFunc: andA,
}

// OrA - OR with accumulator.
var OrA = &Instruction{
	Name: "or",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0xB7, Size: 1, Cycles: 4}, // OR A (base opcode)
		ImmediateAddressing: {Opcode: 0xF6, Size: 2, Cycles: 7}, // OR n
	},
	ParamFunc: orA,
}

// XorA - XOR with accumulator.
var XorA = &Instruction{
	Name: "xor",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0xAF, Size: 1, Cycles: 4}, // XOR A (base opcode)
		ImmediateAddressing: {Opcode: 0xEE, Size: 2, Cycles: 7}, // XOR n
	},
	ParamFunc: xorA,
}

// CpA - Compare with accumulator.
var CpA = &Instruction{
	Name: "cp",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0xBF, Size: 1, Cycles: 4}, // CP A (base opcode)
		ImmediateAddressing: {Opcode: 0xFE, Size: 2, Cycles: 7}, // CP n
	},
	ParamFunc: cpA,
}

// JpAbs - Jump absolute.
var JpAbs = &Instruction{
	Name: "jp",
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xC3, Size: 3, Cycles: 10}, // JP nn
	},
	ParamFunc: jpAbs,
}

// JrRel - Jump relative.
var JrRel = &Instruction{
	Name: "jr",
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x18, Size: 2, Cycles: 12}, // JR e
	},
	ParamFunc: jrRel,
}

// Additional Z80 instructions

// LdReg16 - Load 16-bit register with immediate value.
var LdReg16 = &Instruction{
	Name: "ld",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x01, Size: 3, Cycles: 10}, // LD BC,nn (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0x01, Size: 3, Cycles: 10}, // LD BC,nn
		RegDE: {Opcode: 0x11, Size: 3, Cycles: 10}, // LD DE,nn
		RegHL: {Opcode: 0x21, Size: 3, Cycles: 10}, // LD HL,nn
		RegSP: {Opcode: 0x31, Size: 3, Cycles: 10}, // LD SP,nn
	},
	ParamFunc: ldReg16,
}

// LdIndirect - Load indirect (register pair to memory or memory to register).
var LdIndirect = &Instruction{
	Name: "ld",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0x02, Size: 1, Cycles: 7}, // LD (BC),A
	},
	ParamFunc: ldIndirect,
}

// IncReg16 - Increment 16-bit register.
var IncReg16 = &Instruction{
	Name: "inc",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x03, Size: 1, Cycles: 6}, // INC BC (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0x03, Size: 1, Cycles: 6}, // INC BC
		RegDE: {Opcode: 0x13, Size: 1, Cycles: 6}, // INC DE
		RegHL: {Opcode: 0x23, Size: 1, Cycles: 6}, // INC HL
		RegSP: {Opcode: 0x33, Size: 1, Cycles: 6}, // INC SP
	},
	ParamFunc: incReg16,
}

// DecReg16 - Decrement 16-bit register.
var DecReg16 = &Instruction{
	Name: "dec",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x0B, Size: 1, Cycles: 6}, // DEC BC (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0x0B, Size: 1, Cycles: 6}, // DEC BC
		RegDE: {Opcode: 0x1B, Size: 1, Cycles: 6}, // DEC DE
		RegHL: {Opcode: 0x2B, Size: 1, Cycles: 6}, // DEC HL
		RegSP: {Opcode: 0x3B, Size: 1, Cycles: 6}, // DEC SP
	},
	ParamFunc: decReg16,
}

// Rlca - Rotate Left Circular Accumulator.
var Rlca = &Instruction{
	Name: "rlca",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x07, Size: 1, Cycles: 4},
	},
	NoParamFunc: rlca,
}

// Rrca - Rotate Right Circular Accumulator.
var Rrca = &Instruction{
	Name: "rrca",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x0F, Size: 1, Cycles: 4},
	},
	NoParamFunc: rrca,
}

// Rla - Rotate Left Accumulator through carry.
var Rla = &Instruction{
	Name: "rla",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x17, Size: 1, Cycles: 4},
	},
	NoParamFunc: rla,
}

// Rra - Rotate Right Accumulator through carry.
var Rra = &Instruction{
	Name: "rra",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x1F, Size: 1, Cycles: 4},
	},
	NoParamFunc: rra,
}

// ExAf - Exchange AF with AF'.
var ExAf = &Instruction{
	Name: "ex",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x08, Size: 1, Cycles: 4}, // EX AF,AF'
	},
	NoParamFunc: exAf,
}

// AddHl - Add register pair to HL.
var AddHl = &Instruction{
	Name: "add",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x09, Size: 1, Cycles: 11}, // ADD HL,BC
	},
	ParamFunc: addHl,
}

// Djnz - Decrement B and Jump if Not Zero.
var Djnz = &Instruction{
	Name: "djnz",
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x10, Size: 2, Cycles: 8},
	},
	ParamFunc: djnz,
}

// JrCond - Conditional Jump Relative.
var JrCond = &Instruction{
	Name: "jr",
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x20, Size: 2, Cycles: 7}, // JR NZ,e
	},
	ParamFunc: jrCond,
}

// LdExtended - Load using extended addressing (nn).
var LdExtended = &Instruction{
	Name: "ld",
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0x22, Size: 3, Cycles: 16}, // LD (nn),HL
	},
	ParamFunc: ldExtended,
}

// Daa - Decimal Adjust Accumulator.
var Daa = &Instruction{
	Name: "daa",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x27, Size: 1, Cycles: 4},
	},
	NoParamFunc: daa,
}

// Cpl - Complement Accumulator.
var Cpl = &Instruction{
	Name: "cpl",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x2F, Size: 1, Cycles: 4},
	},
	NoParamFunc: cpl,
}

// IncIndirect - Increment indirect memory location.
var IncIndirect = &Instruction{
	Name: "inc",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0x34, Size: 1, Cycles: 11}, // INC (HL)
	},
	ParamFunc: incIndirect,
}

// DecIndirect - Decrement indirect memory location.
var DecIndirect = &Instruction{
	Name: "dec",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0x35, Size: 1, Cycles: 11}, // DEC (HL)
	},
	ParamFunc: decIndirect,
}

// LdIndirectImm - Load immediate to indirect memory location.
var LdIndirectImm = &Instruction{
	Name: "ld",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0x36, Size: 2, Cycles: 10}, // LD (HL),n
	},
	ParamFunc: ldIndirectImm,
}

// Scf - Set Carry Flag.
var Scf = &Instruction{
	Name: "scf",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x37, Size: 1, Cycles: 4},
	},
	NoParamFunc: scf,
}

// Ccf - Complement Carry Flag.
var Ccf = &Instruction{
	Name: "ccf",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3F, Size: 1, Cycles: 4},
	},
	NoParamFunc: ccf,
}

// AdcA - Add with Carry to Accumulator.
var AdcA = &Instruction{
	Name: "adc",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0x8F, Size: 1, Cycles: 4}, // ADC A,A
		ImmediateAddressing: {Opcode: 0xCE, Size: 2, Cycles: 7}, // ADC A,n
	},
	ParamFunc: adcA,
}

// SbcA - Subtract with Carry from Accumulator.
var SbcA = &Instruction{
	Name: "sbc",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing:  {Opcode: 0x9F, Size: 1, Cycles: 4}, // SBC A,A
		ImmediateAddressing: {Opcode: 0xDE, Size: 2, Cycles: 7}, // SBC A,n
	},
	ParamFunc: sbcA,
}

// RetCond - Conditional Return.
var RetCond = &Instruction{
	Name: "ret",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xC0, Size: 1, Cycles: 5}, // RET NZ
	},
	NoParamFunc: retCond,
}

// PopReg16 - Pop 16-bit register from stack.
var PopReg16 = &Instruction{
	Name: "pop",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0xC1, Size: 1, Cycles: 10}, // POP BC (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0xC1, Size: 1, Cycles: 10}, // POP BC
		RegDE: {Opcode: 0xD1, Size: 1, Cycles: 10}, // POP DE
		RegHL: {Opcode: 0xE1, Size: 1, Cycles: 10}, // POP HL
		RegAF: {Opcode: 0xF1, Size: 1, Cycles: 10}, // POP AF
	},
	ParamFunc: popReg16,
}

// JpCond - Conditional Jump.
var JpCond = &Instruction{
	Name: "jp",
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xC2, Size: 3, Cycles: 10}, // JP NZ,nn
	},
	ParamFunc: jpCond,
}

// CallCond - Conditional Call.
var CallCond = &Instruction{
	Name: "call",
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xC4, Size: 3, Cycles: 10}, // CALL NZ,nn
	},
	ParamFunc: callCond,
}

// PushReg16 - Push 16-bit register to stack.
var PushReg16 = &Instruction{
	Name: "push",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0xC5, Size: 1, Cycles: 11}, // PUSH BC (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Opcode: 0xC5, Size: 1, Cycles: 11}, // PUSH BC
		RegDE: {Opcode: 0xD5, Size: 1, Cycles: 11}, // PUSH DE
		RegHL: {Opcode: 0xE5, Size: 1, Cycles: 11}, // PUSH HL
		RegAF: {Opcode: 0xF5, Size: 1, Cycles: 11}, // PUSH AF
	},
	ParamFunc: pushReg16,
}

// Rst - Restart (call to fixed address).
var Rst = &Instruction{
	Name: "rst",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xC7, Size: 1, Cycles: 11}, // RST 00H (base opcode)
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegRst00: {Opcode: 0xC7, Size: 1, Cycles: 11}, // RST 00H
		RegRst08: {Opcode: 0xCF, Size: 1, Cycles: 11}, // RST 08H
		RegRst10: {Opcode: 0xD7, Size: 1, Cycles: 11}, // RST 10H
		RegRst18: {Opcode: 0xDF, Size: 1, Cycles: 11}, // RST 18H
		RegRst20: {Opcode: 0xE7, Size: 1, Cycles: 11}, // RST 20H
		RegRst28: {Opcode: 0xEF, Size: 1, Cycles: 11}, // RST 28H
		RegRst30: {Opcode: 0xF7, Size: 1, Cycles: 11}, // RST 30H
		RegRst38: {Opcode: 0xFF, Size: 1, Cycles: 11}, // RST 38H
	},
	ParamFunc: rst,
}

// Ret - Return from subroutine.
var Ret = &Instruction{
	Name: "ret",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xC9, Size: 1, Cycles: 10},
	},
	NoParamFunc: ret,
}

// Call - Call subroutine.
var Call = &Instruction{
	Name: "call",
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Opcode: 0xCD, Size: 3, Cycles: 17},
	},
	ParamFunc: call,
}

// OutPort - Output to port.
var OutPort = &Instruction{
	Name: "out",
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Opcode: 0xD3, Size: 2, Cycles: 11}, // OUT (n),A
	},
	ParamFunc: outPort,
}

// InPort - Input from port.
var InPort = &Instruction{
	Name: "in",
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Opcode: 0xDB, Size: 2, Cycles: 11}, // IN A,(n)
	},
	ParamFunc: inPort,
}

// Exx - Exchange register pairs.
var Exx = &Instruction{
	Name: "exx",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xD9, Size: 1, Cycles: 4},
	},
	NoParamFunc: exx,
}

// ExSp - Exchange top of stack with register pair.
var ExSp = &Instruction{
	Name: "ex",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0xE3, Size: 1, Cycles: 19}, // EX (SP),HL
	},
	ParamFunc: exSp,
}

// JpIndirect - Jump indirect.
var JpIndirect = &Instruction{
	Name: "jp",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Opcode: 0xE9, Size: 1, Cycles: 4}, // JP (HL)
	},
	ParamFunc: jpIndirect,
}

// ExDeHl - Exchange DE with HL.
var ExDeHl = &Instruction{
	Name: "ex",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xEB, Size: 1, Cycles: 4}, // EX DE,HL
	},
	NoParamFunc: exDeHl,
}

// Di - Disable Interrupts.
var Di = &Instruction{
	Name: "di",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xF3, Size: 1, Cycles: 4},
	},
	NoParamFunc: di,
}

// Ei - Enable Interrupts.
var Ei = &Instruction{
	Name: "ei",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xFB, Size: 1, Cycles: 4},
	},
	NoParamFunc: ei,
}

// LdSp - Load SP from HL.
var LdSp = &Instruction{
	Name: "ld",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0xF9, Size: 1, Cycles: 6}, // LD SP,HL
	},
	ParamFunc: ldSp,
}

// CB-prefixed instructions (bit operations)
var (
	CBRlc = &Instruction{Name: "rlc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x00, Size: 2, Cycles: 8}}, RegisterOpcodes: map[RegisterParam]OpcodeInfo{RegB: {Opcode: 0x00, Size: 2, Cycles: 8}, RegC: {Opcode: 0x01, Size: 2, Cycles: 8}, RegD: {Opcode: 0x02, Size: 2, Cycles: 8}, RegE: {Opcode: 0x03, Size: 2, Cycles: 8}, RegH: {Opcode: 0x04, Size: 2, Cycles: 8}, RegL: {Opcode: 0x05, Size: 2, Cycles: 8}, RegHLIndirect: {Opcode: 0x06, Size: 2, Cycles: 15}, RegA: {Opcode: 0x07, Size: 2, Cycles: 8}}, ParamFunc: cbRlc}
	CBRrc = &Instruction{Name: "rrc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x08, Size: 2, Cycles: 8}}, RegisterOpcodes: map[RegisterParam]OpcodeInfo{RegB: {Opcode: 0x08, Size: 2, Cycles: 8}, RegC: {Opcode: 0x09, Size: 2, Cycles: 8}, RegD: {Opcode: 0x0A, Size: 2, Cycles: 8}, RegE: {Opcode: 0x0B, Size: 2, Cycles: 8}, RegH: {Opcode: 0x0C, Size: 2, Cycles: 8}, RegL: {Opcode: 0x0D, Size: 2, Cycles: 8}, RegHLIndirect: {Opcode: 0x0E, Size: 2, Cycles: 15}, RegA: {Opcode: 0x0F, Size: 2, Cycles: 8}}, ParamFunc: cbRrc}
	CBRl  = &Instruction{Name: "rl", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x10, Size: 2, Cycles: 8}}, RegisterOpcodes: map[RegisterParam]OpcodeInfo{RegB: {Opcode: 0x10, Size: 2, Cycles: 8}, RegC: {Opcode: 0x11, Size: 2, Cycles: 8}, RegD: {Opcode: 0x12, Size: 2, Cycles: 8}, RegE: {Opcode: 0x13, Size: 2, Cycles: 8}, RegH: {Opcode: 0x14, Size: 2, Cycles: 8}, RegL: {Opcode: 0x15, Size: 2, Cycles: 8}, RegHLIndirect: {Opcode: 0x16, Size: 2, Cycles: 15}, RegA: {Opcode: 0x17, Size: 2, Cycles: 8}}, ParamFunc: cbRl}
	CBRr  = &Instruction{Name: "rr", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x18, Size: 2, Cycles: 8}}, RegisterOpcodes: map[RegisterParam]OpcodeInfo{RegB: {Opcode: 0x18, Size: 2, Cycles: 8}, RegC: {Opcode: 0x19, Size: 2, Cycles: 8}, RegD: {Opcode: 0x1A, Size: 2, Cycles: 8}, RegE: {Opcode: 0x1B, Size: 2, Cycles: 8}, RegH: {Opcode: 0x1C, Size: 2, Cycles: 8}, RegL: {Opcode: 0x1D, Size: 2, Cycles: 8}, RegHLIndirect: {Opcode: 0x1E, Size: 2, Cycles: 15}, RegA: {Opcode: 0x1F, Size: 2, Cycles: 8}}, ParamFunc: cbRr}
	CBSla = &Instruction{Name: "sla", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x20, Size: 2, Cycles: 8}}, RegisterOpcodes: map[RegisterParam]OpcodeInfo{RegB: {Opcode: 0x20, Size: 2, Cycles: 8}, RegC: {Opcode: 0x21, Size: 2, Cycles: 8}, RegD: {Opcode: 0x22, Size: 2, Cycles: 8}, RegE: {Opcode: 0x23, Size: 2, Cycles: 8}, RegH: {Opcode: 0x24, Size: 2, Cycles: 8}, RegL: {Opcode: 0x25, Size: 2, Cycles: 8}, RegHLIndirect: {Opcode: 0x26, Size: 2, Cycles: 15}, RegA: {Opcode: 0x27, Size: 2, Cycles: 8}}, ParamFunc: cbSla}
	CBSra = &Instruction{Name: "sra", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x28, Size: 2, Cycles: 8}}, RegisterOpcodes: map[RegisterParam]OpcodeInfo{RegB: {Opcode: 0x28, Size: 2, Cycles: 8}, RegC: {Opcode: 0x29, Size: 2, Cycles: 8}, RegD: {Opcode: 0x2A, Size: 2, Cycles: 8}, RegE: {Opcode: 0x2B, Size: 2, Cycles: 8}, RegH: {Opcode: 0x2C, Size: 2, Cycles: 8}, RegL: {Opcode: 0x2D, Size: 2, Cycles: 8}, RegHLIndirect: {Opcode: 0x2E, Size: 2, Cycles: 15}, RegA: {Opcode: 0x2F, Size: 2, Cycles: 8}}, ParamFunc: cbSra}
	CBSll = &Instruction{Name: SLL.Name, Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x30, Size: 2, Cycles: 8}}, RegisterOpcodes: map[RegisterParam]OpcodeInfo{RegB: {Opcode: 0x30, Size: 2, Cycles: 8}, RegC: {Opcode: 0x31, Size: 2, Cycles: 8}, RegD: {Opcode: 0x32, Size: 2, Cycles: 8}, RegE: {Opcode: 0x33, Size: 2, Cycles: 8}, RegH: {Opcode: 0x34, Size: 2, Cycles: 8}, RegL: {Opcode: 0x35, Size: 2, Cycles: 8}, RegHLIndirect: {Opcode: 0x36, Size: 2, Cycles: 15}, RegA: {Opcode: 0x37, Size: 2, Cycles: 8}}, ParamFunc: cbSll} // undocumented
	CBSrl = &Instruction{Name: "srl", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x38, Size: 2, Cycles: 8}}, RegisterOpcodes: map[RegisterParam]OpcodeInfo{RegB: {Opcode: 0x38, Size: 2, Cycles: 8}, RegC: {Opcode: 0x39, Size: 2, Cycles: 8}, RegD: {Opcode: 0x3A, Size: 2, Cycles: 8}, RegE: {Opcode: 0x3B, Size: 2, Cycles: 8}, RegH: {Opcode: 0x3C, Size: 2, Cycles: 8}, RegL: {Opcode: 0x3D, Size: 2, Cycles: 8}, RegHLIndirect: {Opcode: 0x3E, Size: 2, Cycles: 15}, RegA: {Opcode: 0x3F, Size: 2, Cycles: 8}}, ParamFunc: cbSrl}
	CBBit = &Instruction{Name: "bit", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x40, Size: 2, Cycles: 8}}, ParamFunc: cbBit}
	CBRes = &Instruction{Name: "res", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x80, Size: 2, Cycles: 8}}, ParamFunc: cbRes}
	CBSet = &Instruction{Name: "set", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0xC0, Size: 2, Cycles: 8}}, ParamFunc: cbSet}
)

// ED-prefixed instructions (extended operations)
var (
	EdNeg  = &Instruction{Name: "neg", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x44, Size: 2, Cycles: 8}}, NoParamFunc: edNeg}
	EdIm0  = &Instruction{Name: "im", Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x46, Size: 2, Cycles: 8}}, ParamFunc: edIm0}
	EdIm1  = &Instruction{Name: "im", Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x56, Size: 2, Cycles: 8}}, ParamFunc: edIm1}
	EdIm2  = &Instruction{Name: "im", Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x5E, Size: 2, Cycles: 8}}, ParamFunc: edIm2}
	EdRetn = &Instruction{Name: "retn", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x45, Size: 2, Cycles: 14}}, NoParamFunc: edRetn}
	EdReti = &Instruction{Name: "reti", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x4D, Size: 2, Cycles: 14}}, NoParamFunc: edReti}
	EdRrd  = &Instruction{Name: "rrd", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x67, Size: 2, Cycles: 18}}, NoParamFunc: edRrd}
	EdRld  = &Instruction{Name: "rld", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x6F, Size: 2, Cycles: 18}}, NoParamFunc: edRld}

	// ED arithmetic instructions
	EdAdcHlBc = &Instruction{Name: "adc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x4A, Size: 2, Cycles: 15}}, ParamFunc: edAdcHlBc}
	EdAdcHlDe = &Instruction{Name: "adc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x5A, Size: 2, Cycles: 15}}, ParamFunc: edAdcHlDe}
	EdAdcHlHl = &Instruction{Name: "adc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x6A, Size: 2, Cycles: 15}}, ParamFunc: edAdcHlHl}
	EdAdcHlSp = &Instruction{Name: "adc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x7A, Size: 2, Cycles: 15}}, ParamFunc: edAdcHlSp}
	EdSbcHlBc = &Instruction{Name: "sbc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x42, Size: 2, Cycles: 15}}, ParamFunc: edSbcHlBc}
	EdSbcHlDe = &Instruction{Name: "sbc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x52, Size: 2, Cycles: 15}}, ParamFunc: edSbcHlDe}
	EdSbcHlHl = &Instruction{Name: "sbc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x62, Size: 2, Cycles: 15}}, ParamFunc: edSbcHlHl}
	EdSbcHlSp = &Instruction{Name: "sbc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x72, Size: 2, Cycles: 15}}, ParamFunc: edSbcHlSp}

	// ED load instructions
	EdLdIA = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x47, Size: 2, Cycles: 9}}, NoParamFunc: edLdIA}
	EdLdRA = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x4F, Size: 2, Cycles: 9}}, NoParamFunc: edLdRA}
	EdLdAI = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x57, Size: 2, Cycles: 9}}, NoParamFunc: edLdAI}
	EdLdAR = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x5F, Size: 2, Cycles: 9}}, NoParamFunc: edLdAR}

	EdLdNnBc = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x43, Size: 4, Cycles: 20}}, ParamFunc: edLdNnBc}
	EdLdNnDe = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x53, Size: 4, Cycles: 20}}, ParamFunc: edLdNnDe}
	EdLdNnHl = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x63, Size: 4, Cycles: 20}}, ParamFunc: edLdNnHl}
	EdLdNnSp = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x73, Size: 4, Cycles: 20}}, ParamFunc: edLdNnSp}
	EdLdBcNn = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x4B, Size: 4, Cycles: 20}}, ParamFunc: edLdBcNn}
	EdLdDeNn = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x5B, Size: 4, Cycles: 20}}, ParamFunc: edLdDeNn}
	EdLdHlNn = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x6B, Size: 4, Cycles: 20}}, ParamFunc: edLdHlNn}
	EdLdSpNn = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x7B, Size: 4, Cycles: 20}}, ParamFunc: edLdSpNn}

	// ED block instructions
	EdLdi  = &Instruction{Name: "ldi", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xA0, Size: 2, Cycles: 16}}, NoParamFunc: edLdi}
	EdLdd  = &Instruction{Name: "ldd", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xA8, Size: 2, Cycles: 16}}, NoParamFunc: edLdd}
	EdLdir = &Instruction{Name: "ldir", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xB0, Size: 2, Cycles: 16}}, NoParamFunc: edLdir}
	EdLddr = &Instruction{Name: "lddr", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xB8, Size: 2, Cycles: 16}}, NoParamFunc: edLddr}
	EdCpi  = &Instruction{Name: "cpi", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xA1, Size: 2, Cycles: 16}}, NoParamFunc: edCpi}
	EdCpd  = &Instruction{Name: "cpd", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xA9, Size: 2, Cycles: 16}}, NoParamFunc: edCpd}
	EdCpir = &Instruction{Name: "cpir", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xB1, Size: 2, Cycles: 21}}, NoParamFunc: edCpir}
	EdCpdr = &Instruction{Name: "cpdr", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xB9, Size: 2, Cycles: 21}}, NoParamFunc: edCpdr}

	// ED I/O instructions
	EdIni  = &Instruction{Name: "ini", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xA2, Size: 2, Cycles: 16}}, NoParamFunc: edIni}
	EdInd  = &Instruction{Name: "ind", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xAA, Size: 2, Cycles: 16}}, NoParamFunc: edInd}
	EdInir = &Instruction{Name: "inir", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xB2, Size: 2, Cycles: 21}}, NoParamFunc: edInir}
	EdIndr = &Instruction{Name: "indr", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xBA, Size: 2, Cycles: 21}}, NoParamFunc: edIndr}
	EdOuti = &Instruction{Name: "outi", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xA3, Size: 2, Cycles: 16}}, NoParamFunc: edOuti}
	EdOutd = &Instruction{Name: "outd", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xAB, Size: 2, Cycles: 16}}, NoParamFunc: edOutd}
	EdOtir = &Instruction{Name: "otir", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xB3, Size: 2, Cycles: 21}}, NoParamFunc: edOtir}
	EdOtdr = &Instruction{Name: "otdr", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xBB, Size: 2, Cycles: 21}}, NoParamFunc: edOtdr}

	EdInBC = &Instruction{Name: "in", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x40, Size: 2, Cycles: 12}}, ParamFunc: edInBC}
	EdInCC = &Instruction{Name: "in", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x48, Size: 2, Cycles: 12}}, ParamFunc: edInCC}
	EdInDC = &Instruction{Name: "in", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x50, Size: 2, Cycles: 12}}, ParamFunc: edInDC}
	EdInEC = &Instruction{Name: "in", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x58, Size: 2, Cycles: 12}}, ParamFunc: edInEC}
	EdInHC = &Instruction{Name: "in", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x60, Size: 2, Cycles: 12}}, ParamFunc: edInHC}
	EdInLC = &Instruction{Name: "in", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x68, Size: 2, Cycles: 12}}, ParamFunc: edInLC}
	EdInAC = &Instruction{Name: "in", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x78, Size: 2, Cycles: 12}}, ParamFunc: edInAC}

	EdOutCB = &Instruction{Name: "out", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x41, Size: 2, Cycles: 12}}, ParamFunc: edOutCB}
	EdOutCC = &Instruction{Name: "out", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x49, Size: 2, Cycles: 12}}, ParamFunc: edOutCC}
	EdOutCD = &Instruction{Name: "out", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x51, Size: 2, Cycles: 12}}, ParamFunc: edOutCD}
	EdOutCE = &Instruction{Name: "out", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x59, Size: 2, Cycles: 12}}, ParamFunc: edOutCE}
	EdOutCH = &Instruction{Name: "out", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x61, Size: 2, Cycles: 12}}, ParamFunc: edOutCH}
	EdOutCL = &Instruction{Name: "out", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x69, Size: 2, Cycles: 12}}, ParamFunc: edOutCL}
	EdOutCA = &Instruction{Name: "out", Addressing: map[AddressingMode]OpcodeInfo{PortAddressing: {Opcode: 0x79, Size: 2, Cycles: 12}}, ParamFunc: edOutCA}
)

// DD-prefixed instructions (IX operations)
var (
	DdIncIX  = &Instruction{Name: "inc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x23, Size: 2, Cycles: 10}}, NoParamFunc: ddIncIX}
	DdDecIX  = &Instruction{Name: "dec", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x2B, Size: 2, Cycles: 10}}, NoParamFunc: ddDecIX}
	DdLdIXnn = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x21, Size: 4, Cycles: 14}}, ParamFunc: ddLdIXnn}
	DdLdNnIX = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x22, Size: 4, Cycles: 20}}, ParamFunc: ddLdNnIX}
	DdLdIXNn = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x2A, Size: 4, Cycles: 20}}, ParamFunc: ddLdIXNn}

	DdAddIXBc = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x09, Size: 2, Cycles: 15}}, ParamFunc: ddAddIXBc}
	DdAddIXDe = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x19, Size: 2, Cycles: 15}}, ParamFunc: ddAddIXDe}
	DdAddIXIX = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x29, Size: 2, Cycles: 15}}, ParamFunc: ddAddIXIX}
	DdAddIXSp = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x39, Size: 2, Cycles: 15}}, ParamFunc: ddAddIXSp}

	DdLdBIXd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x46, Size: 3, Cycles: 19}}, ParamFunc: ddLdBIXd}
	DdLdCIXd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x4E, Size: 3, Cycles: 19}}, ParamFunc: ddLdCIXd}
	DdLdDIXd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x56, Size: 3, Cycles: 19}}, ParamFunc: ddLdDIXd}
	DdLdEIXd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x5E, Size: 3, Cycles: 19}}, ParamFunc: ddLdEIXd}
	DdLdHIXd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x66, Size: 3, Cycles: 19}}, ParamFunc: ddLdHIXd}
	DdLdLIXd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x6E, Size: 3, Cycles: 19}}, ParamFunc: ddLdLIXd}
	DdLdAIXd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x7E, Size: 3, Cycles: 19}}, ParamFunc: ddLdAIXd}

	DdLdIXdB = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x70, Size: 3, Cycles: 19}}, ParamFunc: ddLdIXdB}
	DdLdIXdC = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x71, Size: 3, Cycles: 19}}, ParamFunc: ddLdIXdC}
	DdLdIXdD = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x72, Size: 3, Cycles: 19}}, ParamFunc: ddLdIXdD}
	DdLdIXdE = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x73, Size: 3, Cycles: 19}}, ParamFunc: ddLdIXdE}
	DdLdIXdH = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x74, Size: 3, Cycles: 19}}, ParamFunc: ddLdIXdH}
	DdLdIXdL = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x75, Size: 3, Cycles: 19}}, ParamFunc: ddLdIXdL}
	DdLdIXdA = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x77, Size: 3, Cycles: 19}}, ParamFunc: ddLdIXdA}
	DdLdIXdN = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x36, Size: 4, Cycles: 19}}, ParamFunc: ddLdIXdN}

	DdIncIXd = &Instruction{Name: "inc", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x34, Size: 3, Cycles: 23}}, ParamFunc: ddIncIXd}
	DdDecIXd = &Instruction{Name: "dec", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x35, Size: 3, Cycles: 23}}, ParamFunc: ddDecIXd}

	DdAddAIXd = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x86, Size: 3, Cycles: 19}}, ParamFunc: ddAddAIXd}
	DdAdcAIXd = &Instruction{Name: "adc", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x8E, Size: 3, Cycles: 19}}, ParamFunc: ddAdcAIXd}
	DdSubAIXd = &Instruction{Name: "sub", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x96, Size: 3, Cycles: 19}}, ParamFunc: ddSubAIXd}
	DdSbcAIXd = &Instruction{Name: "sbc", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x9E, Size: 3, Cycles: 19}}, ParamFunc: ddSbcAIXd}
	DdAndAIXd = &Instruction{Name: "and", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xA6, Size: 3, Cycles: 19}}, ParamFunc: ddAndAIXd}
	DdXorAIXd = &Instruction{Name: "xor", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xAE, Size: 3, Cycles: 19}}, ParamFunc: ddXorAIXd}
	DdOrAIXd  = &Instruction{Name: "or", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xB6, Size: 3, Cycles: 19}}, ParamFunc: ddOrAIXd}
	DdCpAIXd  = &Instruction{Name: "cp", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xBE, Size: 3, Cycles: 19}}, ParamFunc: ddCpAIXd}

	DdJpIX    = &Instruction{Name: "jp", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xE9, Size: 2, Cycles: 8}}, NoParamFunc: ddJpIX}
	DdExSpIX  = &Instruction{Name: "ex", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xE3, Size: 2, Cycles: 23}}, NoParamFunc: ddExSpIX}
	DdPushIX  = &Instruction{Name: "push", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0xE5, Size: 2, Cycles: 15}}, NoParamFunc: ddPushIX}
	DdPopIX   = &Instruction{Name: "pop", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0xE1, Size: 2, Cycles: 14}}, NoParamFunc: ddPopIX}
	DdcbShift = &Instruction{Name: "ddcb-shift", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x00, Size: 4, Cycles: 23}}, ParamFunc: ddcbShift}
	DdcbBit   = &Instruction{Name: "bit", Addressing: map[AddressingMode]OpcodeInfo{BitAddressing: {Opcode: 0x40, Size: 4, Cycles: 23}}, ParamFunc: ddcbBit}
	DdcbRes   = &Instruction{Name: "res", Addressing: map[AddressingMode]OpcodeInfo{BitAddressing: {Opcode: 0x80, Size: 4, Cycles: 23}}, ParamFunc: ddcbRes}
	DdcbSet   = &Instruction{Name: "set", Addressing: map[AddressingMode]OpcodeInfo{BitAddressing: {Opcode: 0xC0, Size: 4, Cycles: 23}}, ParamFunc: ddcbSet}
)

// FD-prefixed instructions (IY operations)
var (
	FdIncIY  = &Instruction{Name: "inc", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x23, Size: 2, Cycles: 10}}, NoParamFunc: fdIncIY}
	FdDecIY  = &Instruction{Name: "dec", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x2B, Size: 2, Cycles: 10}}, NoParamFunc: fdDecIY}
	FdLdIYnn = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x21, Size: 4, Cycles: 14}}, ParamFunc: fdLdIYnn}
	FdLdNnIY = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x22, Size: 4, Cycles: 20}}, ParamFunc: fdLdNnIY}
	FdLdIYNn = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ExtendedAddressing: {Opcode: 0x2A, Size: 4, Cycles: 20}}, ParamFunc: fdLdIYNn}

	FdAddIYBc = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x09, Size: 2, Cycles: 15}}, ParamFunc: fdAddIYBc}
	FdAddIYDe = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x19, Size: 2, Cycles: 15}}, ParamFunc: fdAddIYDe}
	FdAddIYIY = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x29, Size: 2, Cycles: 15}}, ParamFunc: fdAddIYIY}
	FdAddIYSp = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0x39, Size: 2, Cycles: 15}}, ParamFunc: fdAddIYSp}

	FdLdBIYd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x46, Size: 3, Cycles: 19}}, ParamFunc: fdLdBIYd}
	FdLdCIYd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x4E, Size: 3, Cycles: 19}}, ParamFunc: fdLdCIYd}
	FdLdDIYd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x56, Size: 3, Cycles: 19}}, ParamFunc: fdLdDIYd}
	FdLdEIYd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x5E, Size: 3, Cycles: 19}}, ParamFunc: fdLdEIYd}
	FdLdHIYd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x66, Size: 3, Cycles: 19}}, ParamFunc: fdLdHIYd}
	FdLdLIYd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x6E, Size: 3, Cycles: 19}}, ParamFunc: fdLdLIYd}
	FdLdAIYd = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x7E, Size: 3, Cycles: 19}}, ParamFunc: fdLdAIYd}

	FdLdIYdB = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x70, Size: 3, Cycles: 19}}, ParamFunc: fdLdIYdB}
	FdLdIYdC = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x71, Size: 3, Cycles: 19}}, ParamFunc: fdLdIYdC}
	FdLdIYdD = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x72, Size: 3, Cycles: 19}}, ParamFunc: fdLdIYdD}
	FdLdIYdE = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x73, Size: 3, Cycles: 19}}, ParamFunc: fdLdIYdE}
	FdLdIYdH = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x74, Size: 3, Cycles: 19}}, ParamFunc: fdLdIYdH}
	FdLdIYdL = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x75, Size: 3, Cycles: 19}}, ParamFunc: fdLdIYdL}
	FdLdIYdA = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x77, Size: 3, Cycles: 19}}, ParamFunc: fdLdIYdA}
	FdLdIYdN = &Instruction{Name: "ld", Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x36, Size: 4, Cycles: 19}}, ParamFunc: fdLdIYdN}

	FdIncIYd = &Instruction{Name: "inc", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x34, Size: 3, Cycles: 23}}, ParamFunc: fdIncIYd}
	FdDecIYd = &Instruction{Name: "dec", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x35, Size: 3, Cycles: 23}}, ParamFunc: fdDecIYd}

	// FD arithmetic instructions
	FdAddAIYd = &Instruction{Name: "add", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x86, Size: 3, Cycles: 19}}, ParamFunc: fdAddAIYd}
	FdAdcAIYd = &Instruction{Name: "adc", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x8E, Size: 3, Cycles: 19}}, ParamFunc: fdAdcAIYd}
	FdSubAIYd = &Instruction{Name: "sub", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x96, Size: 3, Cycles: 19}}, ParamFunc: fdSubAIYd}
	FdSbcAIYd = &Instruction{Name: "sbc", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0x9E, Size: 3, Cycles: 19}}, ParamFunc: fdSbcAIYd}
	FdAndAIYd = &Instruction{Name: "and", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xA6, Size: 3, Cycles: 19}}, ParamFunc: fdAndAIYd}
	FdXorAIYd = &Instruction{Name: "xor", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xAE, Size: 3, Cycles: 19}}, ParamFunc: fdXorAIYd}
	FdOrAIYd  = &Instruction{Name: "or", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xB6, Size: 3, Cycles: 19}}, ParamFunc: fdOrAIYd}
	FdCpAIYd  = &Instruction{Name: "cp", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xBE, Size: 3, Cycles: 19}}, ParamFunc: fdCpAIYd}

	FdJpIY    = &Instruction{Name: "jp", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xE9, Size: 2, Cycles: 8}}, NoParamFunc: fdJpIY}
	FdExSpIY  = &Instruction{Name: "ex", Addressing: map[AddressingMode]OpcodeInfo{RegisterIndirectAddressing: {Opcode: 0xE3, Size: 2, Cycles: 23}}, NoParamFunc: fdExSpIY}
	FdPushIY  = &Instruction{Name: "push", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0xE5, Size: 2, Cycles: 15}}, NoParamFunc: fdPushIY}
	FdPopIY   = &Instruction{Name: "pop", Addressing: map[AddressingMode]OpcodeInfo{RegisterAddressing: {Opcode: 0xE1, Size: 2, Cycles: 14}}, NoParamFunc: fdPopIY}
	FdcbShift = &Instruction{Name: "fdcb-shift", Addressing: map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x00, Size: 4, Cycles: 23}}, ParamFunc: fdcbShift}
	FdcbBit   = &Instruction{Name: "bit", Addressing: map[AddressingMode]OpcodeInfo{BitAddressing: {Opcode: 0x40, Size: 4, Cycles: 23}}, ParamFunc: fdcbBit}
	FdcbRes   = &Instruction{Name: "res", Addressing: map[AddressingMode]OpcodeInfo{BitAddressing: {Opcode: 0x80, Size: 4, Cycles: 23}}, ParamFunc: fdcbRes}
	FdcbSet   = &Instruction{Name: "set", Addressing: map[AddressingMode]OpcodeInfo{BitAddressing: {Opcode: 0xC0, Size: 4, Cycles: 23}}, ParamFunc: fdcbSet}
)
