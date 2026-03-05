// This file contains 65C02 specific instruction definitions.

package m6502

// 65C02 instruction name constants (sorted alphabetically).
const (
	BraName = "bra"
	PhxName = "phx"
	PhyName = "phy"
	PlxName = "plx"
	PlyName = "ply"
	StzName = "stz"
	TrbName = "trb"
	TsbName = "tsb"
)

// Bra - Branch Always.
var Bra = &Instruction{
	Name: BraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x80, Size: 2},
	},
	ParamFunc: bra,
}

// Bit65C02 extends the BIT instruction with additional addressing modes for the 65C02.
// The base BIT instruction (zp, abs) is defined in instruction.go.
var Bit65C02 = &Instruction{
	Name: BitName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x89, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x24, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x34, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x2c, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x3c, Size: 3},
	},
	ParamFunc: bit65c02,
}

// Dec65C02 extends the DEC instruction with accumulator addressing for the 65C02.
var Dec65C02 = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AccumulatorAddressing: {Opcode: 0x3a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0xc6, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0xd6, Size: 2},
		AbsoluteAddressing:    {Opcode: 0xce, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0xde, Size: 3},
	},
	ParamFunc: dec65c02,
}

// Inc65C02 extends the INC instruction with accumulator addressing for the 65C02.
var Inc65C02 = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AccumulatorAddressing: {Opcode: 0x1a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0xe6, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0xf6, Size: 2},
		AbsoluteAddressing:    {Opcode: 0xee, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0xfe, Size: 3},
	},
	ParamFunc: inc65c02,
}

// Jmp65C02 extends the JMP instruction with absolute indexed indirect for the 65C02.
var Jmp65C02 = &Instruction{
	Name: JmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing:          {Opcode: 0x4c, Size: 3},
		IndirectAddressing:          {Opcode: 0x6c},
		AbsoluteXIndirectAddressing: {Opcode: 0x7c, Size: 3},
	},
	ParamFunc: jmp65c02,
}

// Phx - Push X Register.
var Phx = &Instruction{
	Name: PhxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xda, Size: 1},
	},
	NoParamFunc: phx,
}

// Phy - Push Y Register.
var Phy = &Instruction{
	Name: PhyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5a, Size: 1},
	},
	NoParamFunc: phy,
}

// Plx - Pull X Register.
var Plx = &Instruction{
	Name: PlxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xfa, Size: 1},
	},
	NoParamFunc: plx,
}

// Ply - Pull Y Register.
var Ply = &Instruction{
	Name: PlyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x7a, Size: 1},
	},
	NoParamFunc: ply,
}

// Stz - Store Zero.
var Stz = &Instruction{
	Name: StzName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x64, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x74, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x9c, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x9e, Size: 3},
	},
	ParamFunc: stz,
}

// Trb - Test and Reset Bits.
var Trb = &Instruction{
	Name: TrbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing: {Opcode: 0x14, Size: 2},
		AbsoluteAddressing: {Opcode: 0x1c, Size: 3},
	},
	ParamFunc: trb,
}

// Tsb - Test and Set Bits.
var Tsb = &Instruction{
	Name: TsbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing: {Opcode: 0x04, Size: 2},
		AbsoluteAddressing: {Opcode: 0x0c, Size: 3},
	},
	ParamFunc: tsb,
}

// 65C02 instruction variants for existing instructions with new zero page indirect mode.
// These extend ORA, AND, EOR, ADC, STA, LDA, CMP, SBC with (zp) addressing.

// Ora65C02 extends ORA with zero page indirect addressing.
var Ora65C02 = &Instruction{
	Name: OraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:        {Opcode: 0x09, Size: 2},
		ZeroPageAddressing:         {Opcode: 0x05, Size: 2},
		ZeroPageXAddressing:        {Opcode: 0x15, Size: 2},
		AbsoluteAddressing:         {Opcode: 0x0d, Size: 3},
		AbsoluteXAddressing:        {Opcode: 0x1d, Size: 3},
		AbsoluteYAddressing:        {Opcode: 0x19, Size: 3},
		IndirectXAddressing:        {Opcode: 0x01, Size: 2},
		IndirectYAddressing:        {Opcode: 0x11, Size: 2},
		ZeroPageIndirectAddressing: {Opcode: 0x12, Size: 2},
	},
	ParamFunc: ora,
}

// And65C02 extends AND with zero page indirect addressing.
var And65C02 = &Instruction{
	Name: AndName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:        {Opcode: 0x29, Size: 2},
		ZeroPageAddressing:         {Opcode: 0x25, Size: 2},
		ZeroPageXAddressing:        {Opcode: 0x35, Size: 2},
		AbsoluteAddressing:         {Opcode: 0x2d, Size: 3},
		AbsoluteXAddressing:        {Opcode: 0x3d, Size: 3},
		AbsoluteYAddressing:        {Opcode: 0x39, Size: 3},
		IndirectXAddressing:        {Opcode: 0x21, Size: 2},
		IndirectYAddressing:        {Opcode: 0x31, Size: 2},
		ZeroPageIndirectAddressing: {Opcode: 0x32, Size: 2},
	},
	ParamFunc: and,
}

// Eor65C02 extends EOR with zero page indirect addressing.
var Eor65C02 = &Instruction{
	Name: EorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:        {Opcode: 0x49, Size: 2},
		ZeroPageAddressing:         {Opcode: 0x45, Size: 2},
		ZeroPageXAddressing:        {Opcode: 0x55, Size: 2},
		AbsoluteAddressing:         {Opcode: 0x4d, Size: 3},
		AbsoluteXAddressing:        {Opcode: 0x5d, Size: 3},
		AbsoluteYAddressing:        {Opcode: 0x59, Size: 3},
		IndirectXAddressing:        {Opcode: 0x41, Size: 2},
		IndirectYAddressing:        {Opcode: 0x51, Size: 2},
		ZeroPageIndirectAddressing: {Opcode: 0x52, Size: 2},
	},
	ParamFunc: eor,
}

// Adc65C02 extends ADC with zero page indirect addressing.
var Adc65C02 = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:        {Opcode: 0x69, Size: 2},
		ZeroPageAddressing:         {Opcode: 0x65, Size: 2},
		ZeroPageXAddressing:        {Opcode: 0x75, Size: 2},
		AbsoluteAddressing:         {Opcode: 0x6d, Size: 3},
		AbsoluteXAddressing:        {Opcode: 0x7d, Size: 3},
		AbsoluteYAddressing:        {Opcode: 0x79, Size: 3},
		IndirectXAddressing:        {Opcode: 0x61, Size: 2},
		IndirectYAddressing:        {Opcode: 0x71, Size: 2},
		ZeroPageIndirectAddressing: {Opcode: 0x72, Size: 2},
	},
	ParamFunc: adc,
}

// Sta65C02 extends STA with zero page indirect addressing.
var Sta65C02 = &Instruction{
	Name: StaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:         {Opcode: 0x85, Size: 2},
		ZeroPageXAddressing:        {Opcode: 0x95, Size: 2},
		AbsoluteAddressing:         {Opcode: 0x8d, Size: 3},
		AbsoluteXAddressing:        {Opcode: 0x9d, Size: 3},
		AbsoluteYAddressing:        {Opcode: 0x99, Size: 3},
		IndirectXAddressing:        {Opcode: 0x81, Size: 2},
		IndirectYAddressing:        {Opcode: 0x91, Size: 2},
		ZeroPageIndirectAddressing: {Opcode: 0x92, Size: 2},
	},
	ParamFunc: sta,
}

// Lda65C02 extends LDA with zero page indirect addressing.
var Lda65C02 = &Instruction{
	Name: LdaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:        {Opcode: 0xa9, Size: 2},
		ZeroPageAddressing:         {Opcode: 0xa5, Size: 2},
		ZeroPageXAddressing:        {Opcode: 0xb5, Size: 2},
		AbsoluteAddressing:         {Opcode: 0xad, Size: 3},
		AbsoluteXAddressing:        {Opcode: 0xbd, Size: 3},
		AbsoluteYAddressing:        {Opcode: 0xb9, Size: 3},
		IndirectXAddressing:        {Opcode: 0xa1, Size: 2},
		IndirectYAddressing:        {Opcode: 0xb1, Size: 2},
		ZeroPageIndirectAddressing: {Opcode: 0xb2, Size: 2},
	},
	ParamFunc: lda,
}

// Cmp65C02 extends CMP with zero page indirect addressing.
var Cmp65C02 = &Instruction{
	Name: CmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:        {Opcode: 0xc9, Size: 2},
		ZeroPageAddressing:         {Opcode: 0xc5, Size: 2},
		ZeroPageXAddressing:        {Opcode: 0xd5, Size: 2},
		AbsoluteAddressing:         {Opcode: 0xcd, Size: 3},
		AbsoluteXAddressing:        {Opcode: 0xdd, Size: 3},
		AbsoluteYAddressing:        {Opcode: 0xd9, Size: 3},
		IndirectXAddressing:        {Opcode: 0xc1, Size: 2},
		IndirectYAddressing:        {Opcode: 0xd1, Size: 2},
		ZeroPageIndirectAddressing: {Opcode: 0xd2, Size: 2},
	},
	ParamFunc: cmp,
}

// Sbc65C02 extends SBC with zero page indirect addressing.
var Sbc65C02 = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:        {Opcode: 0xe9, Size: 2},
		ZeroPageAddressing:         {Opcode: 0xe5, Size: 2},
		ZeroPageXAddressing:        {Opcode: 0xf5, Size: 2},
		AbsoluteAddressing:         {Opcode: 0xed, Size: 3},
		AbsoluteXAddressing:        {Opcode: 0xfd, Size: 3},
		AbsoluteYAddressing:        {Opcode: 0xf9, Size: 3},
		IndirectXAddressing:        {Opcode: 0xe1, Size: 2},
		IndirectYAddressing:        {Opcode: 0xf1, Size: 2},
		ZeroPageIndirectAddressing: {Opcode: 0xf2, Size: 2},
	},
	ParamFunc: sbc,
}
