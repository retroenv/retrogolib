package z80

// DD prefix instructions - IX register operations

var DdIncIX = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xDD, Opcode: 0x23, Size: 2, Cycles: 10},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0x23, Size: 2, Cycles: 10}, // INC IX
	},
	NoParamFunc: ddIncIX,
}

// DdDecIX decrements IX register (DEC IX, DD prefix).
var DdDecIX = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xDD, Opcode: 0x2B, Size: 2, Cycles: 10},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0x2B, Size: 2, Cycles: 10}, // DEC IX
	},
	NoParamFunc: ddDecIX,
}

// DdLdIXnn loads immediate 16-bit value into IX (LD IX,nn, DD prefix).
var DdLdIXnn = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Prefix: 0xDD, Opcode: 0x21, Size: 4, Cycles: 14},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0x21, Size: 4, Cycles: 14}, // LD IX,nn
	},
	ParamFunc: ddLdIXnn,
}

// DdLdNnIX stores IX to memory address (LD (nn),IX, DD prefix).
var DdLdNnIX = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xDD, Opcode: 0x22, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0x22, Size: 4, Cycles: 20}, // LD (nn),IX
	},
	ParamFunc: ddLdNnIX,
}

// DdLdIXNn loads IX from memory address (LD IX,(nn), DD prefix).
var DdLdIXNn = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xDD, Opcode: 0x2A, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0x2A, Size: 4, Cycles: 20}, // LD IX,(nn)
	},
	ParamFunc: ddLdIXNn,
}

// DdAddIXBc adds BC to IX (ADD IX,BC, DD prefix).
var DdAddIXBc = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xDD, Opcode: 0x09, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0xDD, Opcode: 0x09, Size: 2, Cycles: 15}, // ADD IX,BC
	},
	ParamFunc: ddAddIXBc,
}

// DdAddIXDe adds DE to IX (ADD IX,DE, DD prefix).
var DdAddIXDe = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xDD, Opcode: 0x19, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegDE: {Prefix: 0xDD, Opcode: 0x19, Size: 2, Cycles: 15}, // ADD IX,DE
	},
	ParamFunc: ddAddIXDe,
}

// DdAddIXIX adds IX to IX (ADD IX,IX, DD prefix).
var DdAddIXIX = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xDD, Opcode: 0x29, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0x29, Size: 2, Cycles: 15}, // ADD IX,IX
	},
	ParamFunc: ddAddIXIX,
}

// DdAddIXSp adds SP to IX (ADD IX,SP, DD prefix).
var DdAddIXSp = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xDD, Opcode: 0x39, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegSP: {Prefix: 0xDD, Opcode: 0x39, Size: 2, Cycles: 15}, // ADD IX,SP
	},
	ParamFunc: ddAddIXSp,
}

// DdLdBIXd loads B from IX indexed memory (LD B,(IX+d), DD prefix).
var DdLdBIXd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x46, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0xDD, Opcode: 0x46, Size: 3, Cycles: 19}, // LD B,(IX+d)
	},
	ParamFunc: ddLdBIXd,
}

// DdLdCIXd loads C from IX indexed memory (LD C,(IX+d), DD prefix).
var DdLdCIXd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x4E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegC: {Prefix: 0xDD, Opcode: 0x4E, Size: 3, Cycles: 19}, // LD C,(IX+d)
	},
	ParamFunc: ddLdCIXd,
}

// DdLdDIXd loads D from IX indexed memory (LD D,(IX+d), DD prefix).
var DdLdDIXd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x56, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegD: {Prefix: 0xDD, Opcode: 0x56, Size: 3, Cycles: 19}, // LD D,(IX+d)
	},
	ParamFunc: ddLdDIXd,
}

// DdLdEIXd loads E from IX indexed memory (LD E,(IX+d), DD prefix).
var DdLdEIXd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x5E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegE: {Prefix: 0xDD, Opcode: 0x5E, Size: 3, Cycles: 19}, // LD E,(IX+d)
	},
	ParamFunc: ddLdEIXd,
}

// DdLdHIXd loads H from IX indexed memory (LD H,(IX+d), DD prefix).
var DdLdHIXd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x66, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegH: {Prefix: 0xDD, Opcode: 0x66, Size: 3, Cycles: 19}, // LD H,(IX+d)
	},
	ParamFunc: ddLdHIXd,
}

// DdLdLIXd loads L from IX indexed memory (LD L,(IX+d), DD prefix).
var DdLdLIXd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x6E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegL: {Prefix: 0xDD, Opcode: 0x6E, Size: 3, Cycles: 19}, // LD L,(IX+d)
	},
	ParamFunc: ddLdLIXd,
}

// DdLdAIXd loads A from IX indexed memory (LD A,(IX+d), DD prefix).
var DdLdAIXd = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x7E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0x7E, Size: 3, Cycles: 19}, // LD A,(IX+d)
	},
	ParamFunc: ddLdAIXd,
}

// DdLdIXdB stores B to IX indexed memory (LD (IX+d),B, DD prefix).
var DdLdIXdB = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x70, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0xDD, Opcode: 0x70, Size: 3, Cycles: 19}, // LD (IX+d),B
	},
	ParamFunc: ddLdIXdB,
}

// DdLdIXdC stores C to IX indexed memory (LD (IX+d),C, DD prefix).
var DdLdIXdC = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x71, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegC: {Prefix: 0xDD, Opcode: 0x71, Size: 3, Cycles: 19}, // LD (IX+d),C
	},
	ParamFunc: ddLdIXdC,
}

// DdLdIXdD stores D to IX indexed memory (LD (IX+d),D, DD prefix).
var DdLdIXdD = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x72, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegD: {Prefix: 0xDD, Opcode: 0x72, Size: 3, Cycles: 19}, // LD (IX+d),D
	},
	ParamFunc: ddLdIXdD,
}

// DdLdIXdE stores E to IX indexed memory (LD (IX+d),E, DD prefix).
var DdLdIXdE = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x73, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegE: {Prefix: 0xDD, Opcode: 0x73, Size: 3, Cycles: 19}, // LD (IX+d),E
	},
	ParamFunc: ddLdIXdE,
}

// DdLdIXdH stores H to IX indexed memory (LD (IX+d),H, DD prefix).
var DdLdIXdH = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x74, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegH: {Prefix: 0xDD, Opcode: 0x74, Size: 3, Cycles: 19}, // LD (IX+d),H
	},
	ParamFunc: ddLdIXdH,
}

// DdLdIXdL stores L to IX indexed memory (LD (IX+d),L, DD prefix).
var DdLdIXdL = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x75, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegL: {Prefix: 0xDD, Opcode: 0x75, Size: 3, Cycles: 19}, // LD (IX+d),L
	},
	ParamFunc: ddLdIXdL,
}

// DdLdIXdA stores A to IX indexed memory (LD (IX+d),A, DD prefix).
var DdLdIXdA = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x77, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0x77, Size: 3, Cycles: 19}, // LD (IX+d),A
	},
	ParamFunc: ddLdIXdA,
}

// DdLdIXdN stores immediate to IX indexed memory (LD (IX+d),n, DD prefix).
var DdLdIXdN = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Prefix: 0xDD, Opcode: 0x36, Size: 4, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegImm8: {Prefix: 0xDD, Opcode: 0x36, Size: 4, Cycles: 19}, // LD (IX+d),n
	},
	ParamFunc: ddLdIXdN,
}

// DdIncIXd increments IX indexed memory (INC (IX+d), DD prefix).
var DdIncIXd = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x34, Size: 3, Cycles: 23},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXIndirect: {Prefix: 0xDD, Opcode: 0x34, Size: 3, Cycles: 23}, // INC (IX+d)
	},
	ParamFunc: ddIncIXd,
}

// DdDecIXd decrements IX indexed memory (DEC (IX+d), DD prefix).
var DdDecIXd = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x35, Size: 3, Cycles: 23},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXIndirect: {Prefix: 0xDD, Opcode: 0x35, Size: 3, Cycles: 23}, // DEC (IX+d)
	},
	ParamFunc: ddDecIXd,
}

// DdAddAIXd adds IX indexed memory to A (ADD A,(IX+d), DD prefix).
var DdAddAIXd = &Instruction{
	Name: AddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x86, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0x86, Size: 3, Cycles: 19}, // ADD A,(IX+d)
	},
	ParamFunc: ddAddAIXd,
}

// DdAdcAIXd adds IX indexed memory to A with carry (ADC A,(IX+d), DD prefix).
var DdAdcAIXd = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x8E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0x8E, Size: 3, Cycles: 19}, // ADC A,(IX+d)
	},
	ParamFunc: ddAdcAIXd,
}

// DdSubAIXd subtracts IX indexed memory from A (SUB (IX+d), DD prefix).
var DdSubAIXd = &Instruction{
	Name: SubName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x96, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0x96, Size: 3, Cycles: 19}, // SUB (IX+d)
	},
	ParamFunc: ddSubAIXd,
}

// DdSbcAIXd subtracts IX indexed memory from A with carry (SBC A,(IX+d), DD prefix).
var DdSbcAIXd = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0x9E, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0x9E, Size: 3, Cycles: 19}, // SBC A,(IX+d)
	},
	ParamFunc: ddSbcAIXd,
}

// DdAndAIXd performs logical AND with IX indexed memory (AND (IX+d), DD prefix).
var DdAndAIXd = &Instruction{
	Name: AndName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0xA6, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0xA6, Size: 3, Cycles: 19}, // AND (IX+d)
	},
	ParamFunc: ddAndAIXd,
}

// DdXorAIXd performs logical XOR with IX indexed memory (XOR (IX+d), DD prefix).
var DdXorAIXd = &Instruction{
	Name: XorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0xAE, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0xAE, Size: 3, Cycles: 19}, // XOR (IX+d)
	},
	ParamFunc: ddXorAIXd,
}

// DdOrAIXd performs logical OR with IX indexed memory (OR (IX+d), DD prefix).
var DdOrAIXd = &Instruction{
	Name: OrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0xB6, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0xB6, Size: 3, Cycles: 19}, // OR (IX+d)
	},
	ParamFunc: ddOrAIXd,
}

// DdCpAIXd compares A with IX indexed memory (CP (IX+d), DD prefix).
var DdCpAIXd = &Instruction{
	Name: CpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0xBE, Size: 3, Cycles: 19},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xDD, Opcode: 0xBE, Size: 3, Cycles: 19}, // CP (IX+d)
	},
	ParamFunc: ddCpAIXd,
}

// DdJpIX jumps to address in IX (JP (IX), DD prefix).
var DdJpIX = &Instruction{
	Name: JpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0xE9, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0xE9, Size: 2, Cycles: 8}, // JP (IX)
	},
	NoParamFunc: ddJpIX,
}

// DdExSpIX exchanges IX with top of stack (EX (SP),IX, DD prefix).
var DdExSpIX = &Instruction{
	Name: ExName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterIndirectAddressing: {Prefix: 0xDD, Opcode: 0xE3, Size: 2, Cycles: 23},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0xE3, Size: 2, Cycles: 23}, // EX (SP),IX
	},
	NoParamFunc: ddExSpIX,
}

// DdPushIX pushes IX onto stack (PUSH IX, DD prefix).
var DdPushIX = &Instruction{
	Name: PushName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xDD, Opcode: 0xE5, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0xE5, Size: 2, Cycles: 15}, // PUSH IX
	},
	NoParamFunc: ddPushIX,
}

// DdPopIX pops IX from stack (POP IX, DD prefix).
var DdPopIX = &Instruction{
	Name: PopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xDD, Opcode: 0xE1, Size: 2, Cycles: 14},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIX: {Prefix: 0xDD, Opcode: 0xE1, Size: 2, Cycles: 14}, // POP IX
	},
	NoParamFunc: ddPopIX,
}

// DdcbShift performs shift/rotate operations on IX indexed memory (DDCB prefix).
var DdcbShift = &Instruction{
	Name: DdcbShiftName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xDD, Opcode: 0x00, Size: 4, Cycles: 23},
	},
	ParamFunc: ddcbShift,
}

// DdcbBit tests bit in IX indexed memory (BIT b,(IX+d), DDCB prefix).
var DdcbBit = &Instruction{
	Name: BitName,
	Addressing: map[AddressingMode]OpcodeInfo{
		BitAddressing: {Prefix: 0xDD, Opcode: 0x40, Size: 4, Cycles: 23},
	},
	ParamFunc: ddcbBit,
}

// DdcbRes resets bit in IX indexed memory (RES b,(IX+d), DDCB prefix).
var DdcbRes = &Instruction{
	Name: ResName,
	Addressing: map[AddressingMode]OpcodeInfo{
		BitAddressing: {Prefix: 0xDD, Opcode: 0x80, Size: 4, Cycles: 23},
	},
	ParamFunc: ddcbRes,
}

// DdcbSet sets bit in IX indexed memory (SET b,(IX+d), DDCB prefix).
var DdcbSet = &Instruction{
	Name: SetName,
	Addressing: map[AddressingMode]OpcodeInfo{
		BitAddressing: {Prefix: 0xDD, Opcode: 0xC0, Size: 4, Cycles: 23},
	},
	ParamFunc: ddcbSet,
}

// FdIncIY increments IY register (INC IY, FD prefix).
