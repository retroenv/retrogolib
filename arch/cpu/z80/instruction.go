package z80

// Instruction contains information about a Z80 CPU instruction.
type Instruction struct {
	Name       string // lowercased instruction name
	Unofficial bool   // unofficial instructions are not part of the original Z80 spec

	Addressing map[AddressingMode]OpcodeInfo // addressing mode mapping to opcode info

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
		ImmediateAddressing: {Opcode: 0x3E, Size: 2, Cycles: 7}, // LD A,n
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
	ParamFunc: incReg8,
}

// DecReg8 - Decrement 8-bit register.
var DecReg8 = &Instruction{
	Name: "dec",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x3D, Size: 1, Cycles: 4}, // DEC A (base opcode)
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
		ImmediateAddressing: {Opcode: 0x01, Size: 3, Cycles: 10}, // LD BC,nn
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
		RegisterAddressing: {Opcode: 0x03, Size: 1, Cycles: 6}, // INC BC
	},
	ParamFunc: incReg16,
}

// DecReg16 - Decrement 16-bit register.
var DecReg16 = &Instruction{
	Name: "dec",
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x0B, Size: 1, Cycles: 6}, // DEC BC
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
		RegisterAddressing: {Opcode: 0xC1, Size: 1, Cycles: 10}, // POP BC
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
		RegisterAddressing: {Opcode: 0xC5, Size: 1, Cycles: 11}, // PUSH BC
	},
	ParamFunc: pushReg16,
}

// Rst - Restart (call to fixed address).
var Rst = &Instruction{
	Name: "rst",
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xC7, Size: 1, Cycles: 11}, // RST 00H
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

// Instructions maps instruction names to their information struct.
var Instructions = map[string]*Instruction{
	"nop":  Nop,
	"halt": Halt,
	"ld":   LdImm8, // Primary LD instruction (others can be added later)
	"inc":  IncReg8,
	"dec":  DecReg8,
	"add":  AddA,
	"sub":  SubA,
	"and":  AndA,
	"or":   OrA,
	"xor":  XorA,
	"cp":   CpA,
	"jp":   JpAbs,
	"jr":   JrRel,
}
