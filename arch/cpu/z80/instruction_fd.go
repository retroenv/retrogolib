package z80

// FD prefix instructions - IY register operations

var FdIncIY = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xFD, Opcode: 0x23, Size: 2, Cycles: 10},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0x23, Size: 2, Cycles: 10}, // INC IY
	},
	NoParamFunc: fdIncIY,
}

// FdDecIY decrements IY register (DEC IY, FD prefix).
var FdDecIY = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xFD, Opcode: 0x2B, Size: 2, Cycles: 10},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0x2B, Size: 2, Cycles: 10}, // DEC IY
	},
	NoParamFunc: fdDecIY,
}

// FdLdIYnn loads immediate 16-bit value into IY (LD IY,nn, FD prefix).
var FdLdIYnn = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Prefix: 0xFD, Opcode: 0x21, Size: 4, Cycles: 14},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0x21, Size: 4, Cycles: 14}, // LD IY,nn
	},
	ParamFunc: fdLdIYnn,
}

// FdLdNnIY stores IY to memory address (LD (nn),IY, FD prefix).
var FdLdNnIY = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xFD, Opcode: 0x22, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0x22, Size: 4, Cycles: 20}, // LD (nn),IY
	},
	ParamFunc: fdLdNnIY,
}

// FdLdIYNn loads IY from memory address (LD IY,(nn), FD prefix).
var FdLdIYNn = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xFD, Opcode: 0x2A, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0x2A, Size: 4, Cycles: 20}, // LD IY,(nn)
	},
	ParamFunc: fdLdIYNn,
}

// FdAddIYBc adds BC to IY (ADD IY,BC, FD prefix).
var FdAddIYBc = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xFD, Opcode: 0x09, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0xFD, Opcode: 0x09, Size: 2, Cycles: 15}, // ADD IY,BC
	},
	ParamFunc: fdAddIYBc,
}

// FdAddIYDe adds DE to IY (ADD IY,DE, FD prefix).
var FdAddIYDe = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xFD, Opcode: 0x19, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegDE: {Prefix: 0xFD, Opcode: 0x19, Size: 2, Cycles: 15}, // ADD IY,DE
	},
	ParamFunc: fdAddIYDe,
}

// FdAddIYIY adds IY to IY (ADD IY,IY, FD prefix).
var FdAddIYIY = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xFD, Opcode: 0x29, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0x29, Size: 2, Cycles: 15}, // ADD IY,IY
	},
	ParamFunc: fdAddIYIY,
}

// FdAddIYSp adds SP to IY (ADD IY,SP, FD prefix).
var FdAddIYSp = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xFD, Opcode: 0x39, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegSP: {Prefix: 0xFD, Opcode: 0x39, Size: 2, Cycles: 15}, // ADD IY,SP
	},
	ParamFunc: fdAddIYSp,
}

// FdLdBIYd loads B from IY indexed memory (LD B,(IY+d), FD prefix).
var FdLdBIYd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x46, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0xFD, Opcode: 0x46, Size: 3, Cycles: 19}, // LD B,(IY+d)
	},
	ParamFunc: fdLdBIYd,
}

// FdLdCIYd loads C from IY indexed memory (LD C,(IY+d), FD prefix).
var FdLdCIYd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x4E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegC: {Prefix: 0xFD, Opcode: 0x4E, Size: 3, Cycles: 19}, // LD C,(IY+d)
	},
	ParamFunc: fdLdCIYd,
}

// FdLdDIYd loads D from IY indexed memory (LD D,(IY+d), FD prefix).
var FdLdDIYd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x56, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegD: {Prefix: 0xFD, Opcode: 0x56, Size: 3, Cycles: 19}, // LD D,(IY+d)
	},
	ParamFunc: fdLdDIYd,
}

// FdLdEIYd loads E from IY indexed memory (LD E,(IY+d), FD prefix).
var FdLdEIYd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x5E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegE: {Prefix: 0xFD, Opcode: 0x5E, Size: 3, Cycles: 19}, // LD E,(IY+d)
	},
	ParamFunc: fdLdEIYd,
}

// FdLdHIYd loads H from IY indexed memory (LD H,(IY+d), FD prefix).
var FdLdHIYd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x66, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegH: {Prefix: 0xFD, Opcode: 0x66, Size: 3, Cycles: 19}, // LD H,(IY+d)
	},
	ParamFunc: fdLdHIYd,
}

// FdLdLIYd loads L from IY indexed memory (LD L,(IY+d), FD prefix).
var FdLdLIYd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x6E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegL: {Prefix: 0xFD, Opcode: 0x6E, Size: 3, Cycles: 19}, // LD L,(IY+d)
	},
	ParamFunc: fdLdLIYd,
}

// FdLdAIYd loads A from IY indexed memory (LD A,(IY+d), FD prefix).
var FdLdAIYd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x7E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0x7E, Size: 3, Cycles: 19}, // LD A,(IY+d)
	},
	ParamFunc: fdLdAIYd,
}

// FdLdIYdB stores B to IY indexed memory (LD (IY+d),B, FD prefix).
var FdLdIYdB = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x70, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0xFD, Opcode: 0x70, Size: 3, Cycles: 19}, // LD (IY+d),B
	},
	ParamFunc: fdLdIYdB,
}

// FdLdIYdC stores C to IY indexed memory (LD (IY+d),C, FD prefix).
var FdLdIYdC = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x71, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegC: {Prefix: 0xFD, Opcode: 0x71, Size: 3, Cycles: 19}, // LD (IY+d),C
	},
	ParamFunc: fdLdIYdC,
}

// FdLdIYdD stores D to IY indexed memory (LD (IY+d),D, FD prefix).
var FdLdIYdD = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x72, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegD: {Prefix: 0xFD, Opcode: 0x72, Size: 3, Cycles: 19}, // LD (IY+d),D
	},
	ParamFunc: fdLdIYdD,
}

// FdLdIYdE stores E to IY indexed memory (LD (IY+d),E, FD prefix).
var FdLdIYdE = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x73, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegE: {Prefix: 0xFD, Opcode: 0x73, Size: 3, Cycles: 19}, // LD (IY+d),E
	},
	ParamFunc: fdLdIYdE,
}

// FdLdIYdH stores H to IY indexed memory (LD (IY+d),H, FD prefix).
var FdLdIYdH = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x74, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegH: {Prefix: 0xFD, Opcode: 0x74, Size: 3, Cycles: 19}, // LD (IY+d),H
	},
	ParamFunc: fdLdIYdH,
}

// FdLdIYdL stores L to IY indexed memory (LD (IY+d),L, FD prefix).
var FdLdIYdL = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x75, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegL: {Prefix: 0xFD, Opcode: 0x75, Size: 3, Cycles: 19}, // LD (IY+d),L
	},
	ParamFunc: fdLdIYdL,
}

// FdLdIYdA stores A to IY indexed memory (LD (IY+d),A, FD prefix).
var FdLdIYdA = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x77, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0x77, Size: 3, Cycles: 19}, // LD (IY+d),A
	},
	ParamFunc: fdLdIYdA,
}

// FdLdIYdN stores immediate to IY indexed memory (LD (IY+d),n, FD prefix).
var FdLdIYdN = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Prefix: 0xFD, Opcode: 0x36, Size: 4, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYIndirect: {Prefix: 0xFD, Opcode: 0x36, Size: 4, Cycles: 19}, // LD (IY+d),n
	},
	ParamFunc: fdLdIYdN,
}

// FdIncIYd increments IY indexed memory (INC (IY+d), FD prefix).
var FdIncIYd = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x34, Size: 3, Cycles: 23},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYIndirect: {Prefix: 0xFD, Opcode: 0x34, Size: 3, Cycles: 23}, // INC (IY+d)
	},
	ParamFunc: fdIncIYd,
}

// FdDecIYd decrements IY indexed memory (DEC (IY+d), FD prefix).
var FdDecIYd = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x35, Size: 3, Cycles: 23},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYIndirect: {Prefix: 0xFD, Opcode: 0x35, Size: 3, Cycles: 23}, // DEC (IY+d)
	},
	ParamFunc: fdDecIYd,
}

// FdAddAIYd adds IY indexed memory to A (ADD A,(IY+d), FD prefix).
var FdAddAIYd = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x86, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0x86, Size: 3, Cycles: 19}, // ADD A,(IY+d)
	},
	ParamFunc: fdAddAIYd,
}

// FdAdcAIYd adds IY indexed memory to A with carry (ADC A,(IY+d), FD prefix).
var FdAdcAIYd = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x8E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0x8E, Size: 3, Cycles: 19}, // ADC A,(IY+d)
	},
	ParamFunc: fdAdcAIYd,
}

// FdSubAIYd subtracts IY indexed memory from A (SUB (IY+d), FD prefix).
var FdSubAIYd = &Instruction{
	Name: SubName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x96, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0x96, Size: 3, Cycles: 19}, // SUB (IY+d)
	},
	ParamFunc: fdSubAIYd,
}

// FdSbcAIYd subtracts IY indexed memory from A with carry (SBC A,(IY+d), FD prefix).
var FdSbcAIYd = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0x9E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0x9E, Size: 3, Cycles: 19}, // SBC A,(IY+d)
	},
	ParamFunc: fdSbcAIYd,
}

// FdAndAIYd performs logical AND with IY indexed memory (AND (IY+d), FD prefix).
var FdAndAIYd = &Instruction{
	Name: AndName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0xA6, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0xA6, Size: 3, Cycles: 19}, // AND (IY+d)
	},
	ParamFunc: fdAndAIYd,
}

// FdXorAIYd performs logical XOR with IY indexed memory (XOR (IY+d), FD prefix).
var FdXorAIYd = &Instruction{
	Name: XorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0xAE, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0xAE, Size: 3, Cycles: 19}, // XOR (IY+d)
	},
	ParamFunc: fdXorAIYd,
}

// FdOrAIYd performs logical OR with IY indexed memory (OR (IY+d), FD prefix).
var FdOrAIYd = &Instruction{
	Name: OrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0xB6, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0xB6, Size: 3, Cycles: 19}, // OR (IY+d)
	},
	ParamFunc: fdOrAIYd,
}

// FdCpAIYd compares A with IY indexed memory (CP (IY+d), FD prefix).
var FdCpAIYd = &Instruction{
	Name: CpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0xBE, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xFD, Opcode: 0xBE, Size: 3, Cycles: 19}, // CP (IY+d)
	},
	ParamFunc: fdCpAIYd,
}

// FdJpIY jumps to address in IY (JP (IY), FD prefix).
var FdJpIY = &Instruction{
	Name: JpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0xE9, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0xE9, Size: 2, Cycles: 8}, // JP (IY)
	},
	NoParamFunc: fdJpIY,
}

// FdExSpIY exchanges IY with top of stack (EX (SP),IY, FD prefix).
var FdExSpIY = &Instruction{
	Name: ExName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xFD, Opcode: 0xE3, Size: 2, Cycles: 23},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0xE3, Size: 2, Cycles: 23}, // EX (SP),IY
	},
	NoParamFunc: fdExSpIY,
}

// FdPushIY pushes IY onto stack (PUSH IY, FD prefix).
var FdPushIY = &Instruction{
	Name: PushName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xFD, Opcode: 0xE5, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0xE5, Size: 2, Cycles: 15}, // PUSH IY
	},
	NoParamFunc: fdPushIY,
}

// FdPopIY pops IY from stack (POP IY, FD prefix).
var FdPopIY = &Instruction{
	Name: PopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xFD, Opcode: 0xE1, Size: 2, Cycles: 14},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIY: {Prefix: 0xFD, Opcode: 0xE1, Size: 2, Cycles: 14}, // POP IY
	},
	NoParamFunc: fdPopIY,
}

// FdcbShift performs shift/rotate operations on IY indexed memory (FDCB prefix).
var FdcbShift = &Instruction{
	Name: FdcbShiftName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xFD, Opcode: 0x00, Size: 4, Cycles: 23},
	},
	ParamFunc: fdcbShift,
}

// FdcbBit tests bit in IY indexed memory (BIT b,(IY+d), FDCB prefix).
var FdcbBit = &Instruction{
	Name: BitName,
	Addressing: map[AddressingMode]OpcodeInfo{
		BitAddressing: {Prefix: 0xFD, Opcode: 0x40, Size: 4, Cycles: 23},
	},
	ParamFunc: fdcbBit,
}

// FdcbRes resets bit in IY indexed memory (RES b,(IY+d), FDCB prefix).
var FdcbRes = &Instruction{
	Name: ResName,
	Addressing: map[AddressingMode]OpcodeInfo{
		BitAddressing: {Prefix: 0xFD, Opcode: 0x80, Size: 4, Cycles: 23},
	},
	ParamFunc: fdcbRes,
}

// FdcbSet sets bit in IY indexed memory (SET b,(IY+d), FDCB prefix).
var FdcbSet = &Instruction{
	Name: SetName,
	Addressing: map[AddressingMode]OpcodeInfo{
		BitAddressing: {Prefix: 0xFD, Opcode: 0xC0, Size: 4, Cycles: 23},
	},
	ParamFunc: fdcbSet,
}
