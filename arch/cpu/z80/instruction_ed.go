package z80

// ED prefix instructions - Extended operations, block transfers, I/O

// EdNeg negates accumulator (NEG, ED prefix).
var EdNeg = &Instruction{
	Name: NegName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0x44, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB:  {Prefix: 0xED, Opcode: 0x4C, Size: 2, Cycles: 8}, // NEG (undocumented)
		RegD:  {Prefix: 0xED, Opcode: 0x54, Size: 2, Cycles: 8}, // NEG (undocumented)
		RegH:  {Prefix: 0xED, Opcode: 0x5C, Size: 2, Cycles: 8}, // NEG (undocumented)
		RegBC: {Prefix: 0xED, Opcode: 0x64, Size: 2, Cycles: 8}, // NEG (undocumented)
		RegHL: {Prefix: 0xED, Opcode: 0x6C, Size: 2, Cycles: 8}, // NEG (undocumented)
		RegIX: {Prefix: 0xED, Opcode: 0x74, Size: 2, Cycles: 8}, // NEG (undocumented)
		RegIY: {Prefix: 0xED, Opcode: 0x7C, Size: 2, Cycles: 8}, // NEG (undocumented)
	},
	NoParamFunc: edNeg,
}

// EdIm0 sets interrupt mode 0 (IM 0, ED prefix).
var EdIm0 = &Instruction{
	Name: ImName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Prefix: 0xED, Opcode: 0x46, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIM0: {Prefix: 0xED, Opcode: 0x46, Size: 2, Cycles: 8}, // IM 0
		RegI:   {Prefix: 0xED, Opcode: 0x66, Size: 2, Cycles: 8}, // IM 0 (undocumented)
	},
	ParamFunc: edIm0,
}

// EdIm1 sets interrupt mode 1 (IM 1, ED prefix).
var EdIm1 = &Instruction{
	Name: ImName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Prefix: 0xED, Opcode: 0x56, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIM1: {Prefix: 0xED, Opcode: 0x56, Size: 2, Cycles: 8}, // IM 1
		RegR:   {Prefix: 0xED, Opcode: 0x76, Size: 2, Cycles: 8}, // IM 1 (undocumented)
	},
	ParamFunc: edIm1,
}

// EdIm2 sets interrupt mode 2 (IM 2, ED prefix).
var EdIm2 = &Instruction{
	Name: ImName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Prefix: 0xED, Opcode: 0x5E, Size: 2, Cycles: 8},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIM2:  {Prefix: 0xED, Opcode: 0x5E, Size: 2, Cycles: 8}, // IM 2
		RegImm8: {Prefix: 0xED, Opcode: 0x7E, Size: 2, Cycles: 8}, // IM 2 (undocumented)
	},
	ParamFunc: edIm2,
}

// EdRetn returns from non-maskable interrupt (RETN, ED prefix).
var EdRetn = &Instruction{
	Name: RetnName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0x45, Size: 2, Cycles: 14},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegC: {Prefix: 0xED, Opcode: 0x55, Size: 2, Cycles: 14}, // RETN (undocumented)
		RegE: {Prefix: 0xED, Opcode: 0x65, Size: 2, Cycles: 14}, // RETN (undocumented)
		RegL: {Prefix: 0xED, Opcode: 0x75, Size: 2, Cycles: 14}, // RETN (undocumented)
	},
	NoParamFunc: edRetn,
}

// EdReti returns from maskable interrupt (RETI, ED prefix).
var EdReti = &Instruction{
	Name: RetiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0x4D, Size: 2, Cycles: 14},
	},
	NoParamFunc: edReti,
}

// EdRrd rotates right decimal (RRD, ED prefix).
var EdRrd = &Instruction{
	Name: RrdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0x67, Size: 2, Cycles: 18},
	},
	NoParamFunc: edRrd,
}

// EdRld rotates left decimal (RLD, ED prefix).
var EdRld = &Instruction{
	Name: RldName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0x6F, Size: 2, Cycles: 18},
	},
	NoParamFunc: edRld,
}

// EdAdcHlBc adds BC to HL with carry (ADC HL,BC, ED prefix).
var EdAdcHlBc = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x4A, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0xED, Opcode: 0x4A, Size: 2, Cycles: 15}, // ADC HL,BC
	},
	ParamFunc: edAdcHlBc,
}

// EdAdcHlDe adds DE to HL with carry (ADC HL,DE, ED prefix).
var EdAdcHlDe = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x5A, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegDE: {Prefix: 0xED, Opcode: 0x5A, Size: 2, Cycles: 15}, // ADC HL,DE
	},
	ParamFunc: edAdcHlDe,
}

// EdAdcHlHl adds HL to HL with carry (ADC HL,HL, ED prefix).
var EdAdcHlHl = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x6A, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegHL: {Prefix: 0xED, Opcode: 0x6A, Size: 2, Cycles: 15}, // ADC HL,HL
	},
	ParamFunc: edAdcHlHl,
}

// EdAdcHlSp adds SP to HL with carry (ADC HL,SP, ED prefix).
var EdAdcHlSp = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x7A, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegSP: {Prefix: 0xED, Opcode: 0x7A, Size: 2, Cycles: 15}, // ADC HL,SP
	},
	ParamFunc: edAdcHlSp,
}

// EdSbcHlBc subtracts BC from HL with carry (SBC HL,BC, ED prefix).
var EdSbcHlBc = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x42, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0xED, Opcode: 0x42, Size: 2, Cycles: 15}, // SBC HL,BC
	},
	ParamFunc: edSbcHlBc,
}

// EdSbcHlDe subtracts DE from HL with carry (SBC HL,DE, ED prefix).
var EdSbcHlDe = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x52, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegDE: {Prefix: 0xED, Opcode: 0x52, Size: 2, Cycles: 15}, // SBC HL,DE
	},
	ParamFunc: edSbcHlDe,
}

// EdSbcHlHl subtracts HL from HL with carry (SBC HL,HL, ED prefix).
var EdSbcHlHl = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x62, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegHL: {Prefix: 0xED, Opcode: 0x62, Size: 2, Cycles: 15}, // SBC HL,HL
	},
	ParamFunc: edSbcHlHl,
}

// EdSbcHlSp subtracts SP from HL with carry (SBC HL,SP, ED prefix).
var EdSbcHlSp = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x72, Size: 2, Cycles: 15},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegSP: {Prefix: 0xED, Opcode: 0x72, Size: 2, Cycles: 15}, // SBC HL,SP
	},
	ParamFunc: edSbcHlSp,
}

// EdLdIA loads accumulator into interrupt vector register (LD I,A, ED prefix).
var EdLdIA = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x47, Size: 2, Cycles: 9},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegI: {Prefix: 0xED, Opcode: 0x47, Size: 2, Cycles: 9}, // LD I,A
	},
	NoParamFunc: edLdIA,
}

// EdLdRA loads accumulator into memory refresh register (LD R,A, ED prefix).
var EdLdRA = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x4F, Size: 2, Cycles: 9},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegR: {Prefix: 0xED, Opcode: 0x4F, Size: 2, Cycles: 9}, // LD R,A
	},
	NoParamFunc: edLdRA,
}

// EdLdAI loads interrupt vector register into accumulator (LD A,I, ED prefix).
var EdLdAI = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x57, Size: 2, Cycles: 9},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xED, Opcode: 0x57, Size: 2, Cycles: 9}, // LD A,I
	},
	NoParamFunc: edLdAI,
}

// EdLdAR loads memory refresh register into accumulator (LD A,R, ED prefix).
var EdLdAR = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Prefix: 0xED, Opcode: 0x5F, Size: 2, Cycles: 9},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xED, Opcode: 0x5F, Size: 2, Cycles: 9}, // LD A,R
	},
	NoParamFunc: edLdAR,
}

// EdLdNnBc stores BC to memory address (LD (nn),BC, ED prefix).
var EdLdNnBc = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xED, Opcode: 0x43, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0xED, Opcode: 0x43, Size: 4, Cycles: 20}, // LD (nn),BC
	},
	ParamFunc: edLdNnBc,
}

// EdLdNnDe stores DE to memory address (LD (nn),DE, ED prefix).
var EdLdNnDe = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xED, Opcode: 0x53, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegDE: {Prefix: 0xED, Opcode: 0x53, Size: 4, Cycles: 20}, // LD (nn),DE
	},
	ParamFunc: edLdNnDe,
}

// EdLdNnHl stores HL to memory address (LD (nn),HL, ED prefix).
var EdLdNnHl = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xED, Opcode: 0x63, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegHL: {Prefix: 0xED, Opcode: 0x63, Size: 4, Cycles: 20}, // LD (nn),HL
	},
	ParamFunc: edLdNnHl,
}

// EdLdNnSp stores SP to memory address (LD (nn),SP, ED prefix).
var EdLdNnSp = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xED, Opcode: 0x73, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegSP: {Prefix: 0xED, Opcode: 0x73, Size: 4, Cycles: 20}, // LD (nn),SP
	},
	ParamFunc: edLdNnSp,
}

// EdLdBcNn loads BC from memory address (LD BC,(nn), ED prefix).
var EdLdBcNn = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xED, Opcode: 0x4B, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegBC: {Prefix: 0xED, Opcode: 0x4B, Size: 4, Cycles: 20}, // LD BC,(nn)
	},
	ParamFunc: edLdBcNn,
}

// EdLdDeNn loads DE from memory address (LD DE,(nn), ED prefix).
var EdLdDeNn = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xED, Opcode: 0x5B, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegDE: {Prefix: 0xED, Opcode: 0x5B, Size: 4, Cycles: 20}, // LD DE,(nn)
	},
	ParamFunc: edLdDeNn,
}

// EdLdHlNn loads HL from memory address (LD HL,(nn), ED prefix).
var EdLdHlNn = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xED, Opcode: 0x6B, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegHL: {Prefix: 0xED, Opcode: 0x6B, Size: 4, Cycles: 20}, // LD HL,(nn)
	},
	ParamFunc: edLdHlNn,
}

// EdLdSpNn loads SP from memory address (LD SP,(nn), ED prefix).
var EdLdSpNn = &Instruction{
	Name: LdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ExtendedAddressing: {Prefix: 0xED, Opcode: 0x7B, Size: 4, Cycles: 20},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegSP: {Prefix: 0xED, Opcode: 0x7B, Size: 4, Cycles: 20}, // LD SP,(nn)
	},
	ParamFunc: edLdSpNn,
}

// EdLdi loads and increments (LDI, ED prefix).
var EdLdi = &Instruction{
	Name: LdiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xA0, Size: 2, Cycles: 16},
	},
	NoParamFunc: edLdi,
}

// EdLdd loads and decrements (LDD, ED prefix).
var EdLdd = &Instruction{
	Name: LddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xA8, Size: 2, Cycles: 16},
	},
	NoParamFunc: edLdd,
}

// EdLdir loads and increments with repeat (LDIR, ED prefix).
var EdLdir = &Instruction{
	Name: LdirName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xB0, Size: 2, Cycles: 16},
	},
	NoParamFunc: edLdir,
}

// EdLddr loads and decrements with repeat (LDDR, ED prefix).
var EdLddr = &Instruction{
	Name: LddrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xB8, Size: 2, Cycles: 16},
	},
	NoParamFunc: edLddr,
}

// EdCpi compares and increments (CPI, ED prefix).
var EdCpi = &Instruction{
	Name: CpiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xA1, Size: 2, Cycles: 16},
	},
	NoParamFunc: edCpi,
}

// EdCpd compares and decrements (CPD, ED prefix).
var EdCpd = &Instruction{
	Name: CpdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xA9, Size: 2, Cycles: 16},
	},
	NoParamFunc: edCpd,
}

// EdCpir compares and increments with repeat (CPIR, ED prefix).
var EdCpir = &Instruction{
	Name: CpirName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xB1, Size: 2, Cycles: 21},
	},
	NoParamFunc: edCpir,
}

// EdCpdr compares and decrements with repeat (CPDR, ED prefix).
var EdCpdr = &Instruction{
	Name: CpdrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xB9, Size: 2, Cycles: 21},
	},
	NoParamFunc: edCpdr,
}

// EdIni inputs and increments (INI, ED prefix).
var EdIni = &Instruction{
	Name: IniName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xA2, Size: 2, Cycles: 16},
	},
	NoParamFunc: edIni,
}

// EdInd inputs and decrements (IND, ED prefix).
var EdInd = &Instruction{
	Name: IndName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xAA, Size: 2, Cycles: 16},
	},
	NoParamFunc: edInd,
}

// EdInir inputs and increments with repeat (INIR, ED prefix).
var EdInir = &Instruction{
	Name: InirName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xB2, Size: 2, Cycles: 21},
	},
	NoParamFunc: edInir,
}

// EdIndr inputs and decrements with repeat (INDR, ED prefix).
var EdIndr = &Instruction{
	Name: IndrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xBA, Size: 2, Cycles: 21},
	},
	NoParamFunc: edIndr,
}

// EdOuti outputs and increments (OUTI, ED prefix).
var EdOuti = &Instruction{
	Name: OutiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xA3, Size: 2, Cycles: 16},
	},
	NoParamFunc: edOuti,
}

// EdOutd outputs and decrements (OUTD, ED prefix).
var EdOutd = &Instruction{
	Name: OutdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xAB, Size: 2, Cycles: 16},
	},
	NoParamFunc: edOutd,
}

// EdOtir outputs and increments with repeat (OTIR, ED prefix).
var EdOtir = &Instruction{
	Name: OtirName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xB3, Size: 2, Cycles: 21},
	},
	NoParamFunc: edOtir,
}

// EdOtdr outputs and decrements with repeat (OTDR, ED prefix).
var EdOtdr = &Instruction{
	Name: OtdrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0xED, Opcode: 0xBB, Size: 2, Cycles: 21},
	},
	NoParamFunc: edOtdr,
}

// EdInBC inputs to B from port C (IN B,(C), ED prefix).
var EdInBC = &Instruction{
	Name: InName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x40, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0xED, Opcode: 0x40, Size: 2, Cycles: 12}, // IN B,(C)
	},
	ParamFunc: edInBC,
}

// EdInCC inputs to C from port C (IN C,(C), ED prefix).
var EdInCC = &Instruction{
	Name: InName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x48, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegC: {Prefix: 0xED, Opcode: 0x48, Size: 2, Cycles: 12}, // IN C,(C)
	},
	ParamFunc: edInCC,
}

// EdInDC inputs to D from port C (IN D,(C), ED prefix).
var EdInDC = &Instruction{
	Name: InName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x50, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegD: {Prefix: 0xED, Opcode: 0x50, Size: 2, Cycles: 12}, // IN D,(C)
	},
	ParamFunc: edInDC,
}

// EdInEC inputs to E from port C (IN E,(C), ED prefix).
var EdInEC = &Instruction{
	Name: InName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x58, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegE: {Prefix: 0xED, Opcode: 0x58, Size: 2, Cycles: 12}, // IN E,(C)
	},
	ParamFunc: edInEC,
}

// EdInHC inputs to H from port C (IN H,(C), ED prefix).
var EdInHC = &Instruction{
	Name: InName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x60, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegH: {Prefix: 0xED, Opcode: 0x60, Size: 2, Cycles: 12}, // IN H,(C)
	},
	ParamFunc: edInHC,
}

// EdInLC inputs to L from port C (IN L,(C), ED prefix).
var EdInLC = &Instruction{
	Name: InName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x68, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegL: {Prefix: 0xED, Opcode: 0x68, Size: 2, Cycles: 12}, // IN L,(C)
	},
	ParamFunc: edInLC,
}

// EdInAC inputs to A from port C (IN A,(C), ED prefix).
var EdInAC = &Instruction{
	Name: InName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x78, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xED, Opcode: 0x78, Size: 2, Cycles: 12}, // IN A,(C)
	},
	ParamFunc: edInAC,
}

// EdOutCB outputs B to port C (OUT (C),B, ED prefix).
var EdOutCB = &Instruction{
	Name: OutName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x41, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegB: {Prefix: 0xED, Opcode: 0x41, Size: 2, Cycles: 12}, // OUT (C),B
	},
	ParamFunc: edOutCB,
}

// EdOutCC outputs C to port C (OUT (C),C, ED prefix).
var EdOutCC = &Instruction{
	Name: OutName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x49, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegC: {Prefix: 0xED, Opcode: 0x49, Size: 2, Cycles: 12}, // OUT (C),C
	},
	ParamFunc: edOutCC,
}

// EdOutCD outputs D to port C (OUT (C),D, ED prefix).
var EdOutCD = &Instruction{
	Name: OutName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x51, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegD: {Prefix: 0xED, Opcode: 0x51, Size: 2, Cycles: 12}, // OUT (C),D
	},
	ParamFunc: edOutCD,
}

// EdOutCE outputs E to port C (OUT (C),E, ED prefix).
var EdOutCE = &Instruction{
	Name: OutName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x59, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegE: {Prefix: 0xED, Opcode: 0x59, Size: 2, Cycles: 12}, // OUT (C),E
	},
	ParamFunc: edOutCE,
}

// EdOutCH outputs H to port C (OUT (C),H, ED prefix).
var EdOutCH = &Instruction{
	Name: OutName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x61, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegH: {Prefix: 0xED, Opcode: 0x61, Size: 2, Cycles: 12}, // OUT (C),H
	},
	ParamFunc: edOutCH,
}

// EdOutCL outputs L to port C (OUT (C),L, ED prefix).
var EdOutCL = &Instruction{
	Name: OutName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x69, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegL: {Prefix: 0xED, Opcode: 0x69, Size: 2, Cycles: 12}, // OUT (C),L
	},
	ParamFunc: edOutCL,
}

// EdOutCA outputs A to port C (OUT (C),A, ED prefix).
var EdOutCA = &Instruction{
	Name: OutName,
	Addressing: map[AddressingMode]OpcodeInfo{
		PortAddressing: {Prefix: 0xED, Opcode: 0x79, Size: 2, Cycles: 12},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegA: {Prefix: 0xED, Opcode: 0x79, Size: 2, Cycles: 12}, // OUT (C),A
	},
	ParamFunc: edOutCA,
}

// DdIncIX increments IX register (INC IX, DD prefix).
