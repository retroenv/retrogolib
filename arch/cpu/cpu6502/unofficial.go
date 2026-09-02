// This file contains support for unofficial CPU instructions.
// Reference https://www.nesdev.org/wiki/Programming_with_unofficial_opcodes

package cpu6502

// DcpInst ...
var DcpInst = &Instruction{
	Name:       DcpName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xc7, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xd7, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xcf, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xdf, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xdb, Size: 3},
		IndirectXAddressing: {Opcode: 0xc3, Size: 2},
		IndirectYAddressing: {Opcode: 0xd3, Size: 2},
	},
	ParamFunc: dcp,
}

// IscInst ...
var IscInst = &Instruction{
	Name:       IscName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xe7, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xf7, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xef, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xff, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xfb, Size: 3},
		IndirectXAddressing: {Opcode: 0xe3, Size: 2},
		IndirectYAddressing: {Opcode: 0xf3, Size: 2},
	},
	ParamFunc: isc,
}

// LasInst - AND memory with SP, store result in A, X, and SP.
// Also known as LAR or LAE.
var LasInst = &Instruction{
	Name:       LasName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteYAddressing: {Opcode: 0xbb, Size: 3},
	},
	ParamFunc: las,
}

// LaxInst ...
var LaxInst = &Instruction{
	Name:       LaxName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xa7, Size: 2},
		ZeroPageYAddressing: {Opcode: 0xb7, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xaf, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xbf, Size: 3},
		IndirectXAddressing: {Opcode: 0xa3, Size: 2},
		IndirectYAddressing: {Opcode: 0xb3, Size: 2},
	},
	ParamFunc: lax,
}

// NopUnofficialInst ...
var NopUnofficialInst = &Instruction{
	Name:       NopName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing:   {Opcode: 0x1a, Size: 1},
		ImmediateAddressing: {Opcode: 0x80, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x04, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x14, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x0c, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x1c, Size: 3},
	},
	ParamFunc: nopUnofficial,
}

// RlaInst ...
var RlaInst = &Instruction{
	Name:       RlaName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x27, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x37, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x2f, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x3f, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x3b, Size: 3},
		IndirectXAddressing: {Opcode: 0x23, Size: 2},
		IndirectYAddressing: {Opcode: 0x33, Size: 2},
	},
	ParamFunc: rla,
}

// RraInst ...
var RraInst = &Instruction{
	Name:       RraName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x67, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x77, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x6f, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x7f, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x7b, Size: 3},
		IndirectXAddressing: {Opcode: 0x63, Size: 2},
		IndirectYAddressing: {Opcode: 0x73, Size: 2},
	},
	ParamFunc: rra,
}

// SaxInst ...
var SaxInst = &Instruction{
	Name:       SaxName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x87, Size: 2},
		ZeroPageYAddressing: {Opcode: 0x97, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8f, Size: 3},
		IndirectXAddressing: {Opcode: 0x83, Size: 2},
	},
	ParamFunc: sax,
}

// SbcUnofficialInst ...
var SbcUnofficialInst = &Instruction{
	Name:       SbcName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xeb, Size: 2},
	},
	ParamFunc: sbc,
}

// ShaInst - Store A AND X AND (addr_hi + 1).
// Also known as AHX or AXA. Unstable: address corruption occurs on page cross.
var ShaInst = &Instruction{
	Name:       ShaName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		IndirectYAddressing: {Opcode: 0x93, Size: 2},
		AbsoluteYAddressing: {Opcode: 0x9f, Size: 3},
	},
	ParamFunc: sha,
}

// ShxInst - Store X AND (addr_hi + 1).
// Also known as SXA or XAS. Unstable: address corruption occurs on page cross.
var ShxInst = &Instruction{
	Name:       ShxName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteYAddressing: {Opcode: 0x9e, Size: 3},
	},
	ParamFunc: shx,
}

// ShyInst - Store Y AND (addr_hi + 1).
// Also known as SYA or SAY. Unstable: address corruption occurs on page cross.
var ShyInst = &Instruction{
	Name:       ShyName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteXAddressing: {Opcode: 0x9c, Size: 3},
	},
	ParamFunc: shy,
}

// SloInst ...
var SloInst = &Instruction{
	Name:       SloName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x07, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x17, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x0f, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x1f, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x1b, Size: 3},
		IndirectXAddressing: {Opcode: 0x03, Size: 2},
		IndirectYAddressing: {Opcode: 0x13, Size: 2},
	},
	ParamFunc: slo,
}

// SreInst ...
var SreInst = &Instruction{
	Name:       SreName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x47, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x57, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x4f, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x5f, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x5b, Size: 3},
		IndirectXAddressing: {Opcode: 0x43, Size: 2},
		IndirectYAddressing: {Opcode: 0x53, Size: 2},
	},
	ParamFunc: sre,
}

// TasInst - Transfer A AND X to SP, then store SP AND (addr_hi + 1).
// Also known as XAS or SHS. Unstable: corrupts SP; address corruption on page cross.
var TasInst = &Instruction{
	Name:       TasName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteYAddressing: {Opcode: 0x9b, Size: 3},
	},
	ParamFunc: tas,
}

// AlrInst - AND with accumulator, then LSR.
var AlrInst = &Instruction{
	Name:       AlrName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x4b, Size: 2},
	},
	ParamFunc: alr,
}

// AncInst - AND with accumulator, copy N flag to C flag.
var AncInst = &Instruction{
	Name:       AncName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x0b, Size: 2},
	},
	ParamFunc: anc,
}

// AncUnofficialInst - Alternate opcode for ANC (same behavior as AncInst).
var AncUnofficialInst = &Instruction{
	Name:       AncName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x2b, Size: 2},
	},
	ParamFunc: anc,
}

// AneInst - OR accumulator with magic constant 0xFF, AND with X and immediate, store in A.
// Also known as XAA. Highly unstable: the magic constant varies by chip and environment.
// Reference: https://www.nesdev.org/wiki/Visual6502wiki/6502_Opcode_8B_(XAA,_ANE)
var AneInst = &Instruction{
	Name:       AneName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x8b, Size: 2},
	},
	ParamFunc: ane,
}

// ArrInst - AND with accumulator, then ROR.
var ArrInst = &Instruction{
	Name:       ArrName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x6b, Size: 2},
	},
	ParamFunc: arr,
}

// AxsInst - (A AND X) minus immediate, store in X.
var AxsInst = &Instruction{
	Name:       AxsName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xcb, Size: 2},
	},
	ParamFunc: axs,
}

// LxaInst - OR accumulator with magic constant 0xFF, AND with immediate, store in A and X.
// Also known as ATX or OAL. Highly unstable: the magic constant varies by chip and environment.
// Reference: https://www.nesdev.org/wiki/CPU_unofficial_opcodes
var LxaInst = &Instruction{
	Name:       LxaName,
	Unofficial: true,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xab, Size: 2},
	},
	ParamFunc: lxa,
}
