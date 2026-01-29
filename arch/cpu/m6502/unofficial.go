// This file contains support for unofficial CPU instructions.
// Reference https://www.nesdev.org/wiki/Programming_with_unofficial_opcodes

package m6502

// Dcp ...
var Dcp = &Instruction{
	Name:       DcpName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xc7},
		ZeroPageXAddressing: {Opcode: 0xd7},
		AbsoluteAddressing:  {Opcode: 0xcf},
		AbsoluteXAddressing: {Opcode: 0xdf},
		AbsoluteYAddressing: {Opcode: 0xdb},
		IndirectXAddressing: {Opcode: 0xc3},
		IndirectYAddressing: {Opcode: 0xd3},
	},
	ParamFunc: dcp,
}

// Isc ...
var Isc = &Instruction{
	Name:       IscName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xe7},
		ZeroPageXAddressing: {Opcode: 0xf7},
		AbsoluteAddressing:  {Opcode: 0xef},
		AbsoluteXAddressing: {Opcode: 0xff},
		AbsoluteYAddressing: {Opcode: 0xfb},
		IndirectXAddressing: {Opcode: 0xe3},
		IndirectYAddressing: {Opcode: 0xf3},
	},
	ParamFunc: isc,
}

// Lax ...
var Lax = &Instruction{
	Name:       LaxName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xa7},
		ZeroPageYAddressing: {Opcode: 0xb7},
		AbsoluteAddressing:  {Opcode: 0xaf},
		AbsoluteYAddressing: {Opcode: 0xbf},
		IndirectXAddressing: {Opcode: 0xa3},
		IndirectYAddressing: {Opcode: 0xb3},
	},
	ParamFunc: lax,
}

// NopUnofficial ...
var NopUnofficial = &Instruction{
	Name:       NopName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing:   {Opcode: 0x1a},
		ImmediateAddressing: {Opcode: 0x80},
		ZeroPageAddressing:  {Opcode: 0x04},
		ZeroPageXAddressing: {Opcode: 0x14},
		AbsoluteAddressing:  {Opcode: 0x0c},
		AbsoluteXAddressing: {Opcode: 0x1c},
	},
	ParamFunc: nopUnofficial,
}

// Rla ...
var Rla = &Instruction{
	Name:       RlaName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x27},
		ZeroPageXAddressing: {Opcode: 0x37},
		AbsoluteAddressing:  {Opcode: 0x2f},
		AbsoluteXAddressing: {Opcode: 0x3f},
		AbsoluteYAddressing: {Opcode: 0x3b},
		IndirectXAddressing: {Opcode: 0x23},
		IndirectYAddressing: {Opcode: 0x33},
	},
	ParamFunc: rla,
}

// Rra ...
var Rra = &Instruction{
	Name:       RraName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x67},
		ZeroPageXAddressing: {Opcode: 0x77},
		AbsoluteAddressing:  {Opcode: 0x6f},
		AbsoluteXAddressing: {Opcode: 0x7f},
		AbsoluteYAddressing: {Opcode: 0x7b},
		IndirectXAddressing: {Opcode: 0x63},
		IndirectYAddressing: {Opcode: 0x73},
	},
	ParamFunc: rra,
}

// Sax ...
var Sax = &Instruction{
	Name:       SaxName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x87},
		ZeroPageYAddressing: {Opcode: 0x97},
		AbsoluteAddressing:  {Opcode: 0x8f},
		IndirectXAddressing: {Opcode: 0x83},
	},
	ParamFunc: sax,
}

// SbcUnofficial ...
var SbcUnofficial = &Instruction{
	Name:       SbcName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xeb},
	},
	ParamFunc: sbc,
}

// Slo ...
var Slo = &Instruction{
	Name:       SloName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x07},
		ZeroPageXAddressing: {Opcode: 0x17},
		AbsoluteAddressing:  {Opcode: 0x0f},
		AbsoluteXAddressing: {Opcode: 0x1f},
		AbsoluteYAddressing: {Opcode: 0x1b},
		IndirectXAddressing: {Opcode: 0x03},
		IndirectYAddressing: {Opcode: 0x13},
	},
	ParamFunc: slo,
}

// Sre ...
var Sre = &Instruction{
	Name:       SreName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x47},
		ZeroPageXAddressing: {Opcode: 0x57},
		AbsoluteAddressing:  {Opcode: 0x4f},
		AbsoluteXAddressing: {Opcode: 0x5f},
		AbsoluteYAddressing: {Opcode: 0x5b},
		IndirectXAddressing: {Opcode: 0x43},
		IndirectYAddressing: {Opcode: 0x53},
	},
	ParamFunc: sre,
}

// Alr - AND with accumulator, then LSR.
var Alr = &Instruction{
	Name:       AlrName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x4b},
	},
	ParamFunc: alr,
}

// Anc - AND with accumulator, copy N flag to C flag.
var Anc = &Instruction{
	Name:       AncName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x0b},
	},
	ParamFunc: anc,
}

// AncUnofficial - Alternate opcode for ANC (same behavior as Anc).
var AncUnofficial = &Instruction{
	Name:       AncName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x2b},
	},
	ParamFunc: anc,
}

// Arr - AND with accumulator, then ROR.
var Arr = &Instruction{
	Name:       ArrName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x6b},
	},
	ParamFunc: arr,
}

// Axs - (A AND X) minus immediate, store in X.
var Axs = &Instruction{
	Name:       AxsName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xcb},
	},
	ParamFunc: axs,
}
