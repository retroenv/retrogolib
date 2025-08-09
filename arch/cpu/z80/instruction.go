package z80

// Instruction contains information about a Z80 CPU instruction.
type Instruction struct {
	Name       string // lowercased instruction name
	Unofficial bool   // unofficial instructions are not part of the original Z80 spec

	Addressing      map[AddressingMode]OpcodeInfo       // addressing mode mapping to opcode info
	RegisterOpcodes map[RegisterParam]OpcodeInfo        // register-specific opcode mapping for disambiguating variants

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
// This method provides the disambiguated opcode information that was previously 
// handled by the separate OpcodeMap.
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
// This replaces the functionality from OpcodeMap.GetInstructionVariants.
func (ins Instruction) GetAllRegisterVariants() map[RegisterParam]OpcodeInfo {
	if ins.RegisterOpcodes == nil {
		return nil
	}
	
	// Return a copy to prevent external modification
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
	CBRlc = &Instruction{Name: "rlc", ParamFunc: cbRlc}
	CBRrc = &Instruction{Name: "rrc", ParamFunc: cbRrc}
	CBRl  = &Instruction{Name: "rl", ParamFunc: cbRl}
	CBRr  = &Instruction{Name: "rr", ParamFunc: cbRr}
	CBSla = &Instruction{Name: "sla", ParamFunc: cbSla}
	CBSra = &Instruction{Name: "sra", ParamFunc: cbSra}
	CBSll = &Instruction{Name: SLL.Name, ParamFunc: cbSll} // undocumented
	CBSrl = &Instruction{Name: "srl", ParamFunc: cbSrl}
	CBBit = &Instruction{Name: "bit", ParamFunc: cbBit}
	CBRes = &Instruction{Name: "res", ParamFunc: cbRes}
	CBSet = &Instruction{Name: "set", ParamFunc: cbSet}
)

// ED-prefixed instructions (extended operations)
var (
	EdNeg  = &Instruction{Name: "neg", NoParamFunc: edNeg}
	EdIm0  = &Instruction{Name: "im", ParamFunc: edIm0}
	EdIm1  = &Instruction{Name: "im", ParamFunc: edIm1}
	EdIm2  = &Instruction{Name: "im", ParamFunc: edIm2}
	EdRetn = &Instruction{Name: "retn", NoParamFunc: edRetn}
	EdReti = &Instruction{Name: "reti", NoParamFunc: edReti}
	EdRrd  = &Instruction{Name: "rrd", NoParamFunc: edRrd}
	EdRld  = &Instruction{Name: "rld", NoParamFunc: edRld}

	// ED arithmetic instructions
	EdAdcHlBc = &Instruction{Name: "adc", ParamFunc: edAdcHlBc}
	EdAdcHlDe = &Instruction{Name: "adc", ParamFunc: edAdcHlDe}
	EdAdcHlHl = &Instruction{Name: "adc", ParamFunc: edAdcHlHl}
	EdAdcHlSp = &Instruction{Name: "adc", ParamFunc: edAdcHlSp}
	EdSbcHlBc = &Instruction{Name: "sbc", ParamFunc: edSbcHlBc}
	EdSbcHlDe = &Instruction{Name: "sbc", ParamFunc: edSbcHlDe}
	EdSbcHlHl = &Instruction{Name: "sbc", ParamFunc: edSbcHlHl}
	EdSbcHlSp = &Instruction{Name: "sbc", ParamFunc: edSbcHlSp}

	// ED load instructions
	EdLdIA = &Instruction{Name: "ld", NoParamFunc: edLdIA}
	EdLdRA = &Instruction{Name: "ld", NoParamFunc: edLdRA}
	EdLdAI = &Instruction{Name: "ld", NoParamFunc: edLdAI}
	EdLdAR = &Instruction{Name: "ld", NoParamFunc: edLdAR}

	EdLdNnBc = &Instruction{Name: "ld", ParamFunc: edLdNnBc}
	EdLdNnDe = &Instruction{Name: "ld", ParamFunc: edLdNnDe}
	EdLdNnHl = &Instruction{Name: "ld", ParamFunc: edLdNnHl}
	EdLdNnSp = &Instruction{Name: "ld", ParamFunc: edLdNnSp}
	EdLdBcNn = &Instruction{Name: "ld", ParamFunc: edLdBcNn}
	EdLdDeNn = &Instruction{Name: "ld", ParamFunc: edLdDeNn}
	EdLdHlNn = &Instruction{Name: "ld", ParamFunc: edLdHlNn}
	EdLdSpNn = &Instruction{Name: "ld", ParamFunc: edLdSpNn}

	// ED block instructions
	EdLdi  = &Instruction{Name: "ldi", NoParamFunc: edLdi}
	EdLdd  = &Instruction{Name: "ldd", NoParamFunc: edLdd}
	EdLdir = &Instruction{Name: "ldir", NoParamFunc: edLdir}
	EdLddr = &Instruction{Name: "lddr", NoParamFunc: edLddr}
	EdCpi  = &Instruction{Name: "cpi", NoParamFunc: edCpi}
	EdCpd  = &Instruction{Name: "cpd", NoParamFunc: edCpd}
	EdCpir = &Instruction{Name: "cpir", NoParamFunc: edCpir}
	EdCpdr = &Instruction{Name: "cpdr", NoParamFunc: edCpdr}

	// ED I/O instructions
	EdIni  = &Instruction{Name: "ini", NoParamFunc: edIni}
	EdInd  = &Instruction{Name: "ind", NoParamFunc: edInd}
	EdInir = &Instruction{Name: "inir", NoParamFunc: edInir}
	EdIndr = &Instruction{Name: "indr", NoParamFunc: edIndr}
	EdOuti = &Instruction{Name: "outi", NoParamFunc: edOuti}
	EdOutd = &Instruction{Name: "outd", NoParamFunc: edOutd}
	EdOtir = &Instruction{Name: "otir", NoParamFunc: edOtir}
	EdOtdr = &Instruction{Name: "otdr", NoParamFunc: edOtdr}

	EdInBC = &Instruction{Name: "in", ParamFunc: edInBC}
	EdInCC = &Instruction{Name: "in", ParamFunc: edInCC}
	EdInDC = &Instruction{Name: "in", ParamFunc: edInDC}
	EdInEC = &Instruction{Name: "in", ParamFunc: edInEC}
	EdInHC = &Instruction{Name: "in", ParamFunc: edInHC}
	EdInLC = &Instruction{Name: "in", ParamFunc: edInLC}
	EdInAC = &Instruction{Name: "in", ParamFunc: edInAC}

	EdOutCB = &Instruction{Name: "out", ParamFunc: edOutCB}
	EdOutCC = &Instruction{Name: "out", ParamFunc: edOutCC}
	EdOutCD = &Instruction{Name: "out", ParamFunc: edOutCD}
	EdOutCE = &Instruction{Name: "out", ParamFunc: edOutCE}
	EdOutCH = &Instruction{Name: "out", ParamFunc: edOutCH}
	EdOutCL = &Instruction{Name: "out", ParamFunc: edOutCL}
	EdOutCA = &Instruction{Name: "out", ParamFunc: edOutCA}
)

// DD-prefixed instructions (IX operations)
var (
	DdIncIX  = &Instruction{Name: "inc", NoParamFunc: ddIncIX}
	DdDecIX  = &Instruction{Name: "dec", NoParamFunc: ddDecIX}
	DdLdIXnn = &Instruction{Name: "ld", ParamFunc: ddLdIXnn}
	DdLdNnIX = &Instruction{Name: "ld", ParamFunc: ddLdNnIX}
	DdLdIXNn = &Instruction{Name: "ld", ParamFunc: ddLdIXNn}

	DdAddIXBc = &Instruction{Name: "add", ParamFunc: ddAddIXBc}
	DdAddIXDe = &Instruction{Name: "add", ParamFunc: ddAddIXDe}
	DdAddIXIX = &Instruction{Name: "add", ParamFunc: ddAddIXIX}
	DdAddIXSp = &Instruction{Name: "add", ParamFunc: ddAddIXSp}

	DdLdBIXd = &Instruction{Name: "ld", ParamFunc: ddLdBIXd}
	DdLdCIXd = &Instruction{Name: "ld", ParamFunc: ddLdCIXd}
	DdLdDIXd = &Instruction{Name: "ld", ParamFunc: ddLdDIXd}
	DdLdEIXd = &Instruction{Name: "ld", ParamFunc: ddLdEIXd}
	DdLdHIXd = &Instruction{Name: "ld", ParamFunc: ddLdHIXd}
	DdLdLIXd = &Instruction{Name: "ld", ParamFunc: ddLdLIXd}
	DdLdAIXd = &Instruction{Name: "ld", ParamFunc: ddLdAIXd}

	DdLdIXdB = &Instruction{Name: "ld", ParamFunc: ddLdIXdB}
	DdLdIXdC = &Instruction{Name: "ld", ParamFunc: ddLdIXdC}
	DdLdIXdD = &Instruction{Name: "ld", ParamFunc: ddLdIXdD}
	DdLdIXdE = &Instruction{Name: "ld", ParamFunc: ddLdIXdE}
	DdLdIXdH = &Instruction{Name: "ld", ParamFunc: ddLdIXdH}
	DdLdIXdL = &Instruction{Name: "ld", ParamFunc: ddLdIXdL}
	DdLdIXdA = &Instruction{Name: "ld", ParamFunc: ddLdIXdA}
	DdLdIXdN = &Instruction{Name: "ld", ParamFunc: ddLdIXdN}

	DdIncIXd = &Instruction{Name: "inc", ParamFunc: ddIncIXd}
	DdDecIXd = &Instruction{Name: "dec", ParamFunc: ddDecIXd}

	DdAddAIXd = &Instruction{Name: "add", ParamFunc: ddAddAIXd}
	DdAdcAIXd = &Instruction{Name: "adc", ParamFunc: ddAdcAIXd}
	DdSubAIXd = &Instruction{Name: "sub", ParamFunc: ddSubAIXd}
	DdSbcAIXd = &Instruction{Name: "sbc", ParamFunc: ddSbcAIXd}
	DdAndAIXd = &Instruction{Name: "and", ParamFunc: ddAndAIXd}
	DdXorAIXd = &Instruction{Name: "xor", ParamFunc: ddXorAIXd}
	DdOrAIXd  = &Instruction{Name: "or", ParamFunc: ddOrAIXd}
	DdCpAIXd  = &Instruction{Name: "cp", ParamFunc: ddCpAIXd}

	DdJpIX    = &Instruction{Name: "jp", NoParamFunc: ddJpIX}
	DdExSpIX  = &Instruction{Name: "ex", NoParamFunc: ddExSpIX}
	DdPushIX  = &Instruction{Name: "push", NoParamFunc: ddPushIX}
	DdPopIX   = &Instruction{Name: "pop", NoParamFunc: ddPopIX}
	DdcbShift = &Instruction{Name: "ddcb-shift", ParamFunc: ddcbShift}
	DdcbBit   = &Instruction{Name: "bit", ParamFunc: ddcbBit}
	DdcbRes   = &Instruction{Name: "res", ParamFunc: ddcbRes}
	DdcbSet   = &Instruction{Name: "set", ParamFunc: ddcbSet}
)

// FD-prefixed instructions (IY operations)
var (
	FdIncIY  = &Instruction{Name: "inc", NoParamFunc: fdIncIY}
	FdDecIY  = &Instruction{Name: "dec", NoParamFunc: fdDecIY}
	FdLdIYnn = &Instruction{Name: "ld", ParamFunc: fdLdIYnn}
	FdLdNnIY = &Instruction{Name: "ld", ParamFunc: fdLdNnIY}
	FdLdIYNn = &Instruction{Name: "ld", ParamFunc: fdLdIYNn}

	FdAddIYBc = &Instruction{Name: "add", ParamFunc: fdAddIYBc}
	FdAddIYDe = &Instruction{Name: "add", ParamFunc: fdAddIYDe}
	FdAddIYIY = &Instruction{Name: "add", ParamFunc: fdAddIYIY}
	FdAddIYSp = &Instruction{Name: "add", ParamFunc: fdAddIYSp}

	FdLdBIYd = &Instruction{Name: "ld", ParamFunc: fdLdBIYd}
	FdLdCIYd = &Instruction{Name: "ld", ParamFunc: fdLdCIYd}
	FdLdDIYd = &Instruction{Name: "ld", ParamFunc: fdLdDIYd}
	FdLdEIYd = &Instruction{Name: "ld", ParamFunc: fdLdEIYd}
	FdLdHIYd = &Instruction{Name: "ld", ParamFunc: fdLdHIYd}
	FdLdLIYd = &Instruction{Name: "ld", ParamFunc: fdLdLIYd}
	FdLdAIYd = &Instruction{Name: "ld", ParamFunc: fdLdAIYd}

	FdLdIYdB = &Instruction{Name: "ld", ParamFunc: fdLdIYdB}
	FdLdIYdC = &Instruction{Name: "ld", ParamFunc: fdLdIYdC}
	FdLdIYdD = &Instruction{Name: "ld", ParamFunc: fdLdIYdD}
	FdLdIYdE = &Instruction{Name: "ld", ParamFunc: fdLdIYdE}
	FdLdIYdH = &Instruction{Name: "ld", ParamFunc: fdLdIYdH}
	FdLdIYdL = &Instruction{Name: "ld", ParamFunc: fdLdIYdL}
	FdLdIYdA = &Instruction{Name: "ld", ParamFunc: fdLdIYdA}
	FdLdIYdN = &Instruction{Name: "ld", ParamFunc: fdLdIYdN}

	FdIncIYd = &Instruction{Name: "inc", ParamFunc: fdIncIYd}
	FdDecIYd = &Instruction{Name: "dec", ParamFunc: fdDecIYd}

	// FD arithmetic instructions
	FdAddAIYd = &Instruction{Name: "add", ParamFunc: fdAddAIYd}
	FdAdcAIYd = &Instruction{Name: "adc", ParamFunc: fdAdcAIYd}
	FdSubAIYd = &Instruction{Name: "sub", ParamFunc: fdSubAIYd}
	FdSbcAIYd = &Instruction{Name: "sbc", ParamFunc: fdSbcAIYd}
	FdAndAIYd = &Instruction{Name: "and", ParamFunc: fdAndAIYd}
	FdXorAIYd = &Instruction{Name: "xor", ParamFunc: fdXorAIYd}
	FdOrAIYd  = &Instruction{Name: "or", ParamFunc: fdOrAIYd}
	FdCpAIYd  = &Instruction{Name: "cp", ParamFunc: fdCpAIYd}

	FdJpIY    = &Instruction{Name: "jp", NoParamFunc: fdJpIY}
	FdExSpIY  = &Instruction{Name: "ex", NoParamFunc: fdExSpIY}
	FdPushIY  = &Instruction{Name: "push", NoParamFunc: fdPushIY}
	FdPopIY   = &Instruction{Name: "pop", NoParamFunc: fdPopIY}
	FdcbShift = &Instruction{Name: "fdcb-shift", ParamFunc: fdcbShift}
	FdcbBit   = &Instruction{Name: "bit", ParamFunc: fdcbBit}
	FdcbRes   = &Instruction{Name: "res", ParamFunc: fdcbRes}
	FdcbSet   = &Instruction{Name: "set", ParamFunc: fdcbSet}
)
