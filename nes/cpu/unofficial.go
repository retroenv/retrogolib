// This file contains support for unofficial CPU instructions.
// https://www.nesdev.org/wiki/Programming_with_unofficial_opcodes

package cpu

import . "github.com/retroenv/retrogolib/nes/addressing"

// Dcp ...
var Dcp = &Instruction{
	Name:       "dcp",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0xc7},
		ZeroPageXAddressing: {Opcode: 0xd7},
		AbsoluteAddressing:  {Opcode: 0xcf},
		AbsoluteXAddressing: {Opcode: 0xdf},
		AbsoluteYAddressing: {Opcode: 0xdb},
		IndirectXAddressing: {Opcode: 0xc3},
		IndirectYAddressing: {Opcode: 0xd3},
	},
}

// Isc ...
var Isc = &Instruction{
	Name:       "isc",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0xe7},
		ZeroPageXAddressing: {Opcode: 0xf7},
		AbsoluteAddressing:  {Opcode: 0xef},
		AbsoluteXAddressing: {Opcode: 0xff},
		AbsoluteYAddressing: {Opcode: 0xfb},
		IndirectXAddressing: {Opcode: 0xe3},
		IndirectYAddressing: {Opcode: 0xf3},
	},
}

// Lax ...
var Lax = &Instruction{
	Name:       "lax",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0xa7},
		ZeroPageYAddressing: {Opcode: 0xb7},
		AbsoluteAddressing:  {Opcode: 0xaf},
		AbsoluteYAddressing: {Opcode: 0xbf},
		IndirectXAddressing: {Opcode: 0xa3},
		IndirectYAddressing: {Opcode: 0xb3},
	},
}

// NopUnofficial ...
var NopUnofficial = &Instruction{
	Name:       "nop",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing:   {Opcode: 0x1a},
		ImmediateAddressing: {Opcode: 0x80},
		ZeroPageAddressing:  {Opcode: 0x04},
		ZeroPageXAddressing: {Opcode: 0x14},
		AbsoluteAddressing:  {Opcode: 0x0c},
		AbsoluteXAddressing: {Opcode: 0x1c},
	},
}

// Rla ...
var Rla = &Instruction{
	Name:       "rla",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x27},
		ZeroPageXAddressing: {Opcode: 0x37},
		AbsoluteAddressing:  {Opcode: 0x2f},
		AbsoluteXAddressing: {Opcode: 0x3f},
		AbsoluteYAddressing: {Opcode: 0x3b},
		IndirectXAddressing: {Opcode: 0x23},
		IndirectYAddressing: {Opcode: 0x33},
	},
}

// Rra ...
var Rra = &Instruction{
	Name:       "rra",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x67},
		ZeroPageXAddressing: {Opcode: 0x77},
		AbsoluteAddressing:  {Opcode: 0x6f},
		AbsoluteXAddressing: {Opcode: 0x7f},
		AbsoluteYAddressing: {Opcode: 0x7b},
		IndirectXAddressing: {Opcode: 0x63},
		IndirectYAddressing: {Opcode: 0x73},
	},
}

// Sax ...
var Sax = &Instruction{
	Name:       "sax",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x87},
		ZeroPageYAddressing: {Opcode: 0x97},
		AbsoluteAddressing:  {Opcode: 0x8f},
		IndirectXAddressing: {Opcode: 0x83},
	},
}

// SbcUnofficial ...
var SbcUnofficial = &Instruction{
	Name:       "sbc",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0xeb},
	},
}

// Slo ...
var Slo = &Instruction{
	Name:       "slo",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x07},
		ZeroPageXAddressing: {Opcode: 0x17},
		AbsoluteAddressing:  {Opcode: 0x0f},
		AbsoluteXAddressing: {Opcode: 0x1f},
		AbsoluteYAddressing: {Opcode: 0x1b},
		IndirectXAddressing: {Opcode: 0x03},
		IndirectYAddressing: {Opcode: 0x13},
	},
}

// Sre ...
var Sre = &Instruction{
	Name:       "sre",
	Unofficial: true,
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x47},
		ZeroPageXAddressing: {Opcode: 0x57},
		AbsoluteAddressing:  {Opcode: 0x4f},
		AbsoluteXAddressing: {Opcode: 0x5f},
		AbsoluteYAddressing: {Opcode: 0x5b},
		IndirectXAddressing: {Opcode: 0x43},
		IndirectYAddressing: {Opcode: 0x53},
	},
}
