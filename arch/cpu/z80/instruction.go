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
