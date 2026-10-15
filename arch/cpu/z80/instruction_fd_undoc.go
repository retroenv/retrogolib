package z80

// Undocumented FD prefix instructions - IYH/IYL operations

var FdIncIYH = &Instruction{
	Name:       IncName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x24,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0x24,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdIncIYH,
}
var FdDecIYH = &Instruction{
	Name:       DecName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x25,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0x25,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdDecIYH,
}
var FdIncIYL = &Instruction{
	Name:       IncName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x2C,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0x2C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdIncIYL,
}
var FdDecIYL = &Instruction{
	Name:       DecName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x2D,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0x2D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdDecIYL,
}

var FdLdIYHn = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x26,
			Size:   3,
			Cycles: 11,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0x26,
			Size:   3,
			Cycles: 11,
		},
	},
	ParamFunc: fdLdIYHn,
}
var FdLdIYLn = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x2E,
			Size:   3,
			Cycles: 11,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0x2E,
			Size:   3,
			Cycles: 11,
		},
	},
	ParamFunc: fdLdIYLn,
}

var FdLdBIYH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x44,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegB, RegIYH}: {
			Prefix: PrefixFD,
			Opcode: 0x44,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdBIYH,
}
var FdLdBIYL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x45,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegB, RegIYL}: {
			Prefix: PrefixFD,
			Opcode: 0x45,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdBIYL,
}
var FdLdCIYH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x4C,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegC, RegIYH}: {
			Prefix: PrefixFD,
			Opcode: 0x4C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdCIYH,
}
var FdLdCIYL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x4D,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegC, RegIYL}: {
			Prefix: PrefixFD,
			Opcode: 0x4D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdCIYL,
}
var FdLdDIYH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x54,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegD, RegIYH}: {
			Prefix: PrefixFD,
			Opcode: 0x54,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdDIYH,
}
var FdLdDIYL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x55,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegD, RegIYL}: {
			Prefix: PrefixFD,
			Opcode: 0x55,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdDIYL,
}
var FdLdEIYH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x5C,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegE, RegIYH}: {
			Prefix: PrefixFD,
			Opcode: 0x5C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdEIYH,
}
var FdLdEIYL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x5D,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegE, RegIYL}: {
			Prefix: PrefixFD,
			Opcode: 0x5D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdEIYL,
}
var FdLdAIYH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x7C,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegA, RegIYH}: {
			Prefix: PrefixFD,
			Opcode: 0x7C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdAIYH,
}
var FdLdAIYL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x7D,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegA, RegIYL}: {
			Prefix: PrefixFD,
			Opcode: 0x7D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdAIYL,
}

var FdLdIYHB = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x60,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYH, RegB}: {
			Prefix: PrefixFD,
			Opcode: 0x60,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYHB,
}
var FdLdIYHC = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x61,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYH, RegC}: {
			Prefix: PrefixFD,
			Opcode: 0x61,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYHC,
}
var FdLdIYHD = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x62,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYH, RegD}: {
			Prefix: PrefixFD,
			Opcode: 0x62,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYHD,
}
var FdLdIYHE = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x63,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYH, RegE}: {
			Prefix: PrefixFD,
			Opcode: 0x63,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYHE,
}
var FdLdIYHIYH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x64,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYH, RegIYH}: {
			Prefix: PrefixFD,
			Opcode: 0x64,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYHIYH,
}
var FdLdIYHIYL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x65,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYH, RegIYL}: {
			Prefix: PrefixFD,
			Opcode: 0x65,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYHIYL,
}
var FdLdIYHA = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x67,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYH, RegA}: {
			Prefix: PrefixFD,
			Opcode: 0x67,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYHA,
}

var FdLdIYLB = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x68,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYL, RegB}: {
			Prefix: PrefixFD,
			Opcode: 0x68,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYLB,
}
var FdLdIYLC = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x69,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYL, RegC}: {
			Prefix: PrefixFD,
			Opcode: 0x69,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYLC,
}
var FdLdIYLD = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x6A,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYL, RegD}: {
			Prefix: PrefixFD,
			Opcode: 0x6A,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYLD,
}
var FdLdIYLE = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x6B,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYL, RegE}: {
			Prefix: PrefixFD,
			Opcode: 0x6B,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYLE,
}
var FdLdIYLIYH = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x6C,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYL, RegIYH}: {
			Prefix: PrefixFD,
			Opcode: 0x6C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYLIYH,
}
var FdLdIYLIYL = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x6D,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYL, RegIYL}: {
			Prefix: PrefixFD,
			Opcode: 0x6D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYLIYL,
}
var FdLdIYLA = &Instruction{
	Name:       LdName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x6F,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterPairOpcodes: map[[2]RegisterParam]OpcodeInfo{
		{RegIYL, RegA}: {
			Prefix: PrefixFD,
			Opcode: 0x6F,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdLdIYLA,
}

var FdAddAIYH = &Instruction{
	Name:       AddName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x84,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0x84,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdAddAIYH,
}
var FdAddAIYL = &Instruction{
	Name:       AddName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x85,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0x85,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdAddAIYL,
}
var FdAdcAIYH = &Instruction{
	Name:       AdcName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x8C,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0x8C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdAdcAIYH,
}
var FdAdcAIYL = &Instruction{
	Name:       AdcName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x8D,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0x8D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdAdcAIYL,
}
var FdSubIYH = &Instruction{
	Name:       SubName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x94,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0x94,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdSubIYH,
}
var FdSubIYL = &Instruction{
	Name:       SubName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x95,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0x95,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdSubIYL,
}
var FdSbcAIYH = &Instruction{
	Name:       SbcName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x9C,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0x9C,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdSbcAIYH,
}
var FdSbcAIYL = &Instruction{
	Name:       SbcName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0x9D,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0x9D,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdSbcAIYL,
}
var FdAndIYH = &Instruction{
	Name:       AndName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0xA4,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0xA4,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdAndIYH,
}
var FdAndIYL = &Instruction{
	Name:       AndName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0xA5,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0xA5,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdAndIYL,
}
var FdXorIYH = &Instruction{
	Name:       XorName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0xAC,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0xAC,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdXorIYH,
}
var FdXorIYL = &Instruction{
	Name:       XorName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0xAD,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0xAD,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdXorIYL,
}
var FdOrIYH = &Instruction{
	Name:       OrName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0xB4,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0xB4,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdOrIYH,
}
var FdOrIYL = &Instruction{
	Name:       OrName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0xB5,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0xB5,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdOrIYL,
}
var FdCpIYH = &Instruction{
	Name:       CpName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0xBC,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYH: {
			Prefix: PrefixFD,
			Opcode: 0xBC,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdCpIYH,
}
var FdCpIYL = &Instruction{
	Name:       CpName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {
			Prefix: PrefixFD,
			Opcode: 0xBD,
			Size:   2,
			Cycles: 8,
		},
	},
	RegisterOpcodes: map[RegisterParam]OpcodeInfo{
		RegIYL: {
			Prefix: PrefixFD,
			Opcode: 0xBD,
			Size:   2,
			Cycles: 8,
		},
	},
	NoParamFunc: fdCpIYL,
}
