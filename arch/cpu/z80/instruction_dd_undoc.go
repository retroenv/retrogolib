package z80

// Undocumented DD prefix instructions - IXH/IXL operations

var DdIncIXH = &Instruction{
	Name:       IncName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0x24,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddIncIXH,
}
var DdDecIXH = &Instruction{
	Name:       DecName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0x25,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddDecIXH,
}
var DdIncIXL = &Instruction{
	Name:       IncName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0x2C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddIncIXL,
}
var DdDecIXL = &Instruction{
	Name:       DecName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0x2D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddDecIXL,
}

var DdLdIXHn = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {
			Prefix: PrefixDD,
			Opcode: 0x26,
			Size:   3,
			Cycles: 11,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0x26,
			Size:   3,
			Cycles: 11,
		},
	},
	ParamFunc: ddLdIXHn,
}
var DdLdIXLn = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {
			Prefix: PrefixDD,
			Opcode: 0x2E,
			Size:   3,
			Cycles: 11,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0x2E,
			Size:   3,
			Cycles: 11,
		},
	},
	ParamFunc: ddLdIXLn,
}

var DdLdBIXH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegB, RegIXH}: {
			Prefix: PrefixDD,
			Opcode: 0x44,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdBIXH,
}
var DdLdBIXL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegB, RegIXL}: {
			Prefix: PrefixDD,
			Opcode: 0x45,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdBIXL,
}
var DdLdCIXH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegC, RegIXH}: {
			Prefix: PrefixDD,
			Opcode: 0x4C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdCIXH,
}
var DdLdCIXL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegC, RegIXL}: {
			Prefix: PrefixDD,
			Opcode: 0x4D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdCIXL,
}
var DdLdDIXH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegD, RegIXH}: {
			Prefix: PrefixDD,
			Opcode: 0x54,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdDIXH,
}
var DdLdDIXL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegD, RegIXL}: {
			Prefix: PrefixDD,
			Opcode: 0x55,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdDIXL,
}
var DdLdEIXH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegE, RegIXH}: {
			Prefix: PrefixDD,
			Opcode: 0x5C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdEIXH,
}
var DdLdEIXL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegE, RegIXL}: {
			Prefix: PrefixDD,
			Opcode: 0x5D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdEIXL,
}
var DdLdAIXH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegA, RegIXH}: {
			Prefix: PrefixDD,
			Opcode: 0x7C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdAIXH,
}
var DdLdAIXL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegA, RegIXL}: {
			Prefix: PrefixDD,
			Opcode: 0x7D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdAIXL,
}

var DdLdIXHB = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXH, RegB}: {
			Prefix: PrefixDD,
			Opcode: 0x60,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXHB,
}
var DdLdIXHC = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXH, RegC}: {
			Prefix: PrefixDD,
			Opcode: 0x61,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXHC,
}
var DdLdIXHD = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXH, RegD}: {
			Prefix: PrefixDD,
			Opcode: 0x62,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXHD,
}
var DdLdIXHE = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXH, RegE}: {
			Prefix: PrefixDD,
			Opcode: 0x63,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXHE,
}
var DdLdIXHIXH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXH, RegIXH}: {
			Prefix: PrefixDD,
			Opcode: 0x64,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXHIXH,
}
var DdLdIXHIXL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXH, RegIXL}: {
			Prefix: PrefixDD,
			Opcode: 0x65,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXHIXL,
}
var DdLdIXHA = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXH, RegA}: {
			Prefix: PrefixDD,
			Opcode: 0x67,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXHA,
}

var DdLdIXLB = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXL, RegB}: {
			Prefix: PrefixDD,
			Opcode: 0x68,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXLB,
}
var DdLdIXLC = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXL, RegC}: {
			Prefix: PrefixDD,
			Opcode: 0x69,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXLC,
}
var DdLdIXLD = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXL, RegD}: {
			Prefix: PrefixDD,
			Opcode: 0x6A,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXLD,
}
var DdLdIXLE = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXL, RegE}: {
			Prefix: PrefixDD,
			Opcode: 0x6B,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXLE,
}
var DdLdIXLIXH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXL, RegIXH}: {
			Prefix: PrefixDD,
			Opcode: 0x6C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXLIXH,
}
var DdLdIXLIXL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXL, RegIXL}: {
			Prefix: PrefixDD,
			Opcode: 0x6D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXLIXL,
}
var DdLdIXLA = &Instruction{
	Name:       LdName,
	Unofficial: true,
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIXL, RegA}: {
			Prefix: PrefixDD,
			Opcode: 0x6F,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddLdIXLA,
}

var DdAddAIXH = &Instruction{
	Name:       AddName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0x84,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddAddAIXH,
}
var DdAddAIXL = &Instruction{
	Name:       AddName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0x85,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddAddAIXL,
}
var DdAdcAIXH = &Instruction{
	Name:       AdcName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0x8C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddAdcAIXH,
}
var DdAdcAIXL = &Instruction{
	Name:       AdcName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0x8D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddAdcAIXL,
}
var DdSubIXH = &Instruction{
	Name:       SubName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0x94,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddSubIXH,
}
var DdSubIXL = &Instruction{
	Name:       SubName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0x95,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddSubIXL,
}
var DdSbcAIXH = &Instruction{
	Name:       SbcName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0x9C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddSbcAIXH,
}
var DdSbcAIXL = &Instruction{
	Name:       SbcName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0x9D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddSbcAIXL,
}
var DdAndIXH = &Instruction{
	Name:       AndName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0xA4,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddAndIXH,
}
var DdAndIXL = &Instruction{
	Name:       AndName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0xA5,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddAndIXL,
}
var DdXorIXH = &Instruction{
	Name:       XorName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0xAC,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddXorIXH,
}
var DdXorIXL = &Instruction{
	Name:       XorName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0xAD,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddXorIXL,
}
var DdOrIXH = &Instruction{
	Name:       OrName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0xB4,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddOrIXH,
}
var DdOrIXL = &Instruction{
	Name:       OrName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0xB5,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddOrIXL,
}
var DdCpIXH = &Instruction{
	Name:       CpName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXH: {
			Prefix: PrefixDD,
			Opcode: 0xBC,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddCpIXH,
}
var DdCpIXL = &Instruction{
	Name:       CpName,
	Unofficial: true,
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIXL: {
			Prefix: PrefixDD,
			Opcode: 0xBD,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: ddCpIXL,
}
