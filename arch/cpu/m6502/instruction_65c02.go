// This file contains 65C02 specific instruction definitions.

package m6502

// 65C02 instruction name constants (sorted alphabetically).
const (
	BbrName = "bbr"
	BbsName = "bbs"
	BraName = "bra"
	PhxName = "phx"
	PhyName = "phy"
	PlxName = "plx"
	PlyName = "ply"
	RmbName = "rmb"
	SmbName = "smb"
	StzName = "stz"
	TrbName = "trb"
	TsbName = "tsb"
)

// BraInst - Branch Always.
var BraInst = &Instruction{
	Name: BraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x80, Size: 2},
	},
	ParamFunc: bra,
}

// Bit65C02Inst extends the BIT instruction with additional addressing modes for the 65C02.
// The base BIT instruction (zp, abs) is defined in instruction.go.
var Bit65C02Inst = &Instruction{
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

// Dec65C02Inst extends the DEC instruction with accumulator addressing for the 65C02.
var Dec65C02Inst = &Instruction{
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

// Inc65C02Inst extends the INC instruction with accumulator addressing for the 65C02.
var Inc65C02Inst = &Instruction{
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

// Jmp65C02Inst extends the JMP instruction with absolute indexed indirect for the 65C02.
var Jmp65C02Inst = &Instruction{
	Name: JmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing:          {Opcode: 0x4c, Size: 3},
		IndirectAddressing:          {Opcode: 0x6c},
		AbsoluteXIndirectAddressing: {Opcode: 0x7c, Size: 3},
	},
	ParamFunc: jmp65c02,
}

// PhxInst - Push X Register.
var PhxInst = &Instruction{
	Name: PhxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xda, Size: 1},
	},
	NoParamFunc: phx,
}

// PhyInst - Push Y Register.
var PhyInst = &Instruction{
	Name: PhyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5a, Size: 1},
	},
	NoParamFunc: phy,
}

// PlxInst - Pull X Register.
var PlxInst = &Instruction{
	Name: PlxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xfa, Size: 1},
	},
	NoParamFunc: plx,
}

// PlyInst - Pull Y Register.
var PlyInst = &Instruction{
	Name: PlyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x7a, Size: 1},
	},
	NoParamFunc: ply,
}

// StzInst - Store Zero.
var StzInst = &Instruction{
	Name: StzName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x64, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x74, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x9c, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x9e, Size: 3},
	},
	ParamFunc: stz,
}

// TrbInst - Test and Reset Bits.
var TrbInst = &Instruction{
	Name: TrbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing: {Opcode: 0x14, Size: 2},
		AbsoluteAddressing: {Opcode: 0x1c, Size: 3},
	},
	ParamFunc: trb,
}

// TsbInst - Test and Set Bits.
var TsbInst = &Instruction{
	Name: TsbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing: {Opcode: 0x04, Size: 2},
		AbsoluteAddressing: {Opcode: 0x0c, Size: 3},
	},
	ParamFunc: tsb,
}

// 65C02 instruction variants for existing instructions with new zero page indirect mode.
// These extend ORA, AND, EOR, ADC, STA, LDA, CMP, SBC with (zp) addressing.

// Ora65C02Inst extends ORA with zero page indirect addressing.
var Ora65C02Inst = &Instruction{
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

// And65C02Inst extends AND with zero page indirect addressing.
var And65C02Inst = &Instruction{
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

// Eor65C02Inst extends EOR with zero page indirect addressing.
var Eor65C02Inst = &Instruction{
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

// Adc65C02Inst extends ADC with zero page indirect addressing.
var Adc65C02Inst = &Instruction{
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

// Sta65C02Inst extends STA with zero page indirect addressing.
var Sta65C02Inst = &Instruction{
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

// Lda65C02Inst extends LDA with zero page indirect addressing.
var Lda65C02Inst = &Instruction{
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

// Cmp65C02Inst extends CMP with zero page indirect addressing.
var Cmp65C02Inst = &Instruction{
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

// Sbc65C02Inst extends SBC with zero page indirect addressing.
var Sbc65C02Inst = &Instruction{
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

// Rockwell 65C02 extensions: RMB, SMB, BBR, BBS.
// These are present in Rockwell and WDC 65C02 variants but not the original GTE/Synertek 65C02.

// Rmb0-Rmb7 reset a specific bit in a zero-page memory location.
var (
	Rmb0 = &Instruction{Name: RmbName + "0", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x07, Size: 2}}, ParamFunc: rmbFunc(0)}
	Rmb1 = &Instruction{Name: RmbName + "1", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x17, Size: 2}}, ParamFunc: rmbFunc(1)}
	Rmb2 = &Instruction{Name: RmbName + "2", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x27, Size: 2}}, ParamFunc: rmbFunc(2)}
	Rmb3 = &Instruction{Name: RmbName + "3", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x37, Size: 2}}, ParamFunc: rmbFunc(3)}
	Rmb4 = &Instruction{Name: RmbName + "4", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x47, Size: 2}}, ParamFunc: rmbFunc(4)}
	Rmb5 = &Instruction{Name: RmbName + "5", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x57, Size: 2}}, ParamFunc: rmbFunc(5)}
	Rmb6 = &Instruction{Name: RmbName + "6", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x67, Size: 2}}, ParamFunc: rmbFunc(6)}
	Rmb7 = &Instruction{Name: RmbName + "7", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x77, Size: 2}}, ParamFunc: rmbFunc(7)}
)

// Smb0-Smb7 set a specific bit in a zero-page memory location.
var (
	Smb0 = &Instruction{Name: SmbName + "0", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x87, Size: 2}}, ParamFunc: smbFunc(0)}
	Smb1 = &Instruction{Name: SmbName + "1", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0x97, Size: 2}}, ParamFunc: smbFunc(1)}
	Smb2 = &Instruction{Name: SmbName + "2", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0xa7, Size: 2}}, ParamFunc: smbFunc(2)}
	Smb3 = &Instruction{Name: SmbName + "3", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0xb7, Size: 2}}, ParamFunc: smbFunc(3)}
	Smb4 = &Instruction{Name: SmbName + "4", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0xc7, Size: 2}}, ParamFunc: smbFunc(4)}
	Smb5 = &Instruction{Name: SmbName + "5", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0xd7, Size: 2}}, ParamFunc: smbFunc(5)}
	Smb6 = &Instruction{Name: SmbName + "6", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0xe7, Size: 2}}, ParamFunc: smbFunc(6)}
	Smb7 = &Instruction{Name: SmbName + "7", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageAddressing: {Opcode: 0xf7, Size: 2}}, ParamFunc: smbFunc(7)}
)

// Bbr0-Bbr7 branch if the specified bit of a zero-page byte is reset (0).
var (
	Bbr0 = &Instruction{Name: BbrName + "0", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x0f, Size: 3}}, ParamFunc: bbrFunc(0)}
	Bbr1 = &Instruction{Name: BbrName + "1", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x1f, Size: 3}}, ParamFunc: bbrFunc(1)}
	Bbr2 = &Instruction{Name: BbrName + "2", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x2f, Size: 3}}, ParamFunc: bbrFunc(2)}
	Bbr3 = &Instruction{Name: BbrName + "3", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x3f, Size: 3}}, ParamFunc: bbrFunc(3)}
	Bbr4 = &Instruction{Name: BbrName + "4", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x4f, Size: 3}}, ParamFunc: bbrFunc(4)}
	Bbr5 = &Instruction{Name: BbrName + "5", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x5f, Size: 3}}, ParamFunc: bbrFunc(5)}
	Bbr6 = &Instruction{Name: BbrName + "6", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x6f, Size: 3}}, ParamFunc: bbrFunc(6)}
	Bbr7 = &Instruction{Name: BbrName + "7", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x7f, Size: 3}}, ParamFunc: bbrFunc(7)}
)

// Bbs0-Bbs7 branch if the specified bit of a zero-page byte is set (1).
var (
	Bbs0 = &Instruction{Name: BbsName + "0", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x8f, Size: 3}}, ParamFunc: bbsFunc(0)}
	Bbs1 = &Instruction{Name: BbsName + "1", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0x9f, Size: 3}}, ParamFunc: bbsFunc(1)}
	Bbs2 = &Instruction{Name: BbsName + "2", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0xaf, Size: 3}}, ParamFunc: bbsFunc(2)}
	Bbs3 = &Instruction{Name: BbsName + "3", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0xbf, Size: 3}}, ParamFunc: bbsFunc(3)}
	Bbs4 = &Instruction{Name: BbsName + "4", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0xcf, Size: 3}}, ParamFunc: bbsFunc(4)}
	Bbs5 = &Instruction{Name: BbsName + "5", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0xdf, Size: 3}}, ParamFunc: bbsFunc(5)}
	Bbs6 = &Instruction{Name: BbsName + "6", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0xef, Size: 3}}, ParamFunc: bbsFunc(6)}
	Bbs7 = &Instruction{Name: BbsName + "7", Addressing: map[AddressingMode]OpcodeInfo{ZeroPageRelativeAddressing: {Opcode: 0xff, Size: 3}}, ParamFunc: bbsFunc(7)}
)
