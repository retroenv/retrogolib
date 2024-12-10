package m6502

import (
	. "github.com/retroenv/retrogolib/addressing"
)

// Instruction contains information about a CPU instruction.
type Instruction struct {
	Name       string // lowercased instruction name
	Unofficial bool   // unofficial instructions are not part of the original 6502 spec

	Addressing map[Mode]OpcodeInfo // addressing mode mapping to opcode info

	NoParamFunc func(c *CPU) error                // emulation function to execute when the instruction has no parameters
	ParamFunc   func(c *CPU, params ...any) error // emulation function to execute when the instruction has parameters
}

// HasAddressing returns whether the instruction has any of the passed addressing modes.
func (ins Instruction) HasAddressing(flags ...Mode) bool {
	for _, flag := range flags {
		_, ok := ins.Addressing[flag]
		if ok {
			return ok
		}
	}
	return false
}

// Adc - Add with Carry.
var Adc = &Instruction{
	Name: "adc",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x69, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x65, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x75, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x6d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x7d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x79, Size: 3},
		IndirectXAddressing: {Opcode: 0x61, Size: 2},
		IndirectYAddressing: {Opcode: 0x71, Size: 2},
	},
	ParamFunc: adc,
}

// And - AND with accumulator.
var And = &Instruction{
	Name: "and",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x29, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x25, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x35, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x2d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x3d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x39, Size: 3},
		IndirectXAddressing: {Opcode: 0x21, Size: 2},
		IndirectYAddressing: {Opcode: 0x31, Size: 2},
	},
	ParamFunc: and,
}

// Asl - Arithmetic Shift Left.
var Asl = &Instruction{
	Name: "asl",
	Addressing: map[Mode]OpcodeInfo{
		AccumulatorAddressing: {Opcode: 0x0a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0x06, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0x16, Size: 2},
		AbsoluteAddressing:    {Opcode: 0x0e, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0x1e, Size: 3},
	},
	ParamFunc: asl,
}

// Bcc - Branch if Carry Clear.
var Bcc = &Instruction{
	Name: "bcc",
	Addressing: map[Mode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x90, Size: 2},
	},
	ParamFunc: bcc,
}

// Bcs - Branch if Carry Set.
var Bcs = &Instruction{
	Name: "bcs",
	Addressing: map[Mode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0xb0, Size: 2},
	},
	ParamFunc: bcs,
}

// Beq - Branch if Equal.
var Beq = &Instruction{
	Name: "beq",
	Addressing: map[Mode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0xf0, Size: 2},
	},
	ParamFunc: beq,
}

// Bit - Bit Test.
var Bit = &Instruction{
	Name: "bit",
	Addressing: map[Mode]OpcodeInfo{
		ZeroPageAddressing: {Opcode: 0x24, Size: 2},
		AbsoluteAddressing: {Opcode: 0x2c, Size: 3},
	},
	ParamFunc: bit,
}

// Bmi - Branch if Minus.
var Bmi = &Instruction{
	Name: "bmi",
	Addressing: map[Mode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x30, Size: 2},
	},
	ParamFunc: bmi,
}

// Bne - Branch if Not Equal.
var Bne = &Instruction{
	Name: "bne",
	Addressing: map[Mode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0xd0, Size: 2},
	},
	ParamFunc: bne,
}

// Bpl - Branch if Positive.
var Bpl = &Instruction{
	Name: "bpl",
	Addressing: map[Mode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x10, Size: 2},
	},
	ParamFunc: bpl,
}

// Brk - Force Interrupt.
var Brk = &Instruction{
	Name: "brk",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x00, Size: 1},
	},
	NoParamFunc: brk,
}

// Bvc - Branch if Overflow Clear.
var Bvc = &Instruction{
	Name: "bvc",
	Addressing: map[Mode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x50, Size: 2},
	},
	ParamFunc: bvc,
}

// Bvs - Branch if Overflow Set.
var Bvs = &Instruction{
	Name: "bvs",
	Addressing: map[Mode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x70, Size: 2},
	},
	ParamFunc: bvs,
}

// Clc - Clear Carry Flag.
var Clc = &Instruction{
	Name: "clc",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x18, Size: 1},
	},
	NoParamFunc: clc,
}

// Cld - Clear Decimal Mode.
var Cld = &Instruction{
	Name: "cld",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xd8, Size: 1},
	},
	NoParamFunc: cld,
}

// Cli - Clear Interrupt Disable.
var Cli = &Instruction{
	Name: "cli",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x58, Size: 1},
	},
	NoParamFunc: cli,
}

// Clv - Clear Overflow Flag.
var Clv = &Instruction{
	Name: "clv",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xb8, Size: 1},
	},
	NoParamFunc: clv,
}

// Cmp - Compare the contents of A.
var Cmp = &Instruction{
	Name: "cmp",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xc9, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xc5, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xd5, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xcd, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xdd, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xd9, Size: 3},
		IndirectXAddressing: {Opcode: 0xc1, Size: 2},
		IndirectYAddressing: {Opcode: 0xd1, Size: 2},
	},
	ParamFunc: cmp,
}

// Cpx - Compare the contents of X.
var Cpx = &Instruction{
	Name: "cpx",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xe0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xe4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xec, Size: 3},
	},
	ParamFunc: cpx,
}

// Cpy - Compare the contents of Y.
var Cpy = &Instruction{
	Name: "cpy",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xc0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xc4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xcc, Size: 3},
	},
	ParamFunc: cpy,
}

// Dec - Decrement memory.
var Dec = &Instruction{
	Name: "dec",
	Addressing: map[Mode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xc6, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xd6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xce, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xde, Size: 3},
	},
	ParamFunc: dec,
}

// Dex - Decrement X Register.
var Dex = &Instruction{
	Name: "dex",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xca, Size: 1},
	},
	NoParamFunc: dex,
}

// Dey - Decrement Y Register.
var Dey = &Instruction{
	Name: "dey",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x88, Size: 1},
	},
	NoParamFunc: dey,
}

// Eor - Exclusive OR - XOR.
var Eor = &Instruction{
	Name: "eor",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x49, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x45, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x55, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x4d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x5d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x59, Size: 3},
		IndirectXAddressing: {Opcode: 0x41, Size: 2},
		IndirectYAddressing: {Opcode: 0x51, Size: 2},
	},
	ParamFunc: eor,
}

// Inc - Increments memory.
var Inc = &Instruction{
	Name: "inc",
	Addressing: map[Mode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xe6, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xf6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xee, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xfe, Size: 3},
	},
	ParamFunc: inc,
}

// Inx - Increment X Register.
var Inx = &Instruction{
	Name: "inx",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xe8, Size: 1},
	},
	NoParamFunc: inx,
}

// Iny - Increment Y Register.
var Iny = &Instruction{
	Name: "iny",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xc8, Size: 1},
	},
	NoParamFunc: iny,
}

// Jmp - jump to address.
var Jmp = &Instruction{
	Name: "jmp",
	Addressing: map[Mode]OpcodeInfo{
		AbsoluteAddressing: {Opcode: 0x4c, Size: 3},
		IndirectAddressing: {Opcode: 0x6c},
	},
	ParamFunc: jmp,
}

// Jsr - jump to subroutine.
var Jsr = &Instruction{
	Name: "jsr",
	Addressing: map[Mode]OpcodeInfo{
		AbsoluteAddressing: {Opcode: 0x20, Size: 3},
	},
	ParamFunc: jsr,
}

// Lda - Load Accumulator - load a byte into A.
var Lda = &Instruction{
	Name: "lda",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xa9, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xa5, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xb5, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xad, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xbd, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xb9, Size: 3},
		IndirectXAddressing: {Opcode: 0xa1, Size: 2},
		IndirectYAddressing: {Opcode: 0xb1, Size: 2},
	},
	ParamFunc: lda,
}

// Ldx - Load X Register - load a byte into X.
var Ldx = &Instruction{
	Name: "ldx",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xa2, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xa6, Size: 2},
		ZeroPageYAddressing: {Opcode: 0xb6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xae, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xbe, Size: 3},
	},
	ParamFunc: ldx,
}

// Ldy - Load Y Register - load a byte into Y.
var Ldy = &Instruction{
	Name: "ldy",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xa0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xa4, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xb4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xac, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xbc, Size: 3},
	},
	ParamFunc: ldy,
}

// Lsr - Logical Shift Right.
var Lsr = &Instruction{
	Name: "lsr",
	Addressing: map[Mode]OpcodeInfo{
		AccumulatorAddressing: {Opcode: 0x4a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0x46, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0x56, Size: 2},
		AbsoluteAddressing:    {Opcode: 0x4e, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0x5e, Size: 3},
	},
	ParamFunc: lsr,
}

// Nop - No Operation.
var Nop = &Instruction{
	Name: "nop",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xea, Size: 1},
	},
	NoParamFunc: nop,
}

// Ora - OR with Accumulator.
var Ora = &Instruction{
	Name: "ora",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x09, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x05, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x15, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x0d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x1d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x19, Size: 3},
		IndirectXAddressing: {Opcode: 0x01, Size: 2},
		IndirectYAddressing: {Opcode: 0x11, Size: 2},
	},
	ParamFunc: ora,
}

// Pha - Push Accumulator - push A content to stack.
var Pha = &Instruction{
	Name: "pha",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x48, Size: 1},
	},
	NoParamFunc: pha,
}

// Php - Push Processor Status - push status flags to stack.
var Php = &Instruction{
	Name: "php",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x08, Size: 1},
	},
	NoParamFunc: php,
}

// Pla - Pull Accumulator - pull A content from stack.
var Pla = &Instruction{
	Name: "pla",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x68, Size: 1},
	},
	NoParamFunc: pla,
}

// Plp - Pull Processor Status - pull status flags from stack.
var Plp = &Instruction{
	Name: "plp",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x28, Size: 1},
	},
	NoParamFunc: plp,
}

// Rol - Rotate Left.
var Rol = &Instruction{
	Name: "rol",
	Addressing: map[Mode]OpcodeInfo{
		AccumulatorAddressing: {Opcode: 0x2a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0x26, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0x36, Size: 2},
		AbsoluteAddressing:    {Opcode: 0x2e, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0x3e, Size: 3},
	},
	ParamFunc: rol,
}

// Ror - Rotate Right.
var Ror = &Instruction{
	Name: "ror",
	Addressing: map[Mode]OpcodeInfo{
		AccumulatorAddressing: {Opcode: 0x6a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0x66, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0x76, Size: 2},
		AbsoluteAddressing:    {Opcode: 0x6e, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0x7e, Size: 3},
	},
	ParamFunc: ror,
}

// Rti - Return from Interrupt.
var Rti = &Instruction{
	Name: "rti",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x40, Size: 1},
	},
	NoParamFunc: rti,
}

// Rts - return from subroutine.
var Rts = &Instruction{
	Name: "rts",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x60, Size: 1},
	},
	NoParamFunc: rts,
}

// Sbc - subtract with Carry.
var Sbc = &Instruction{
	Name: "sbc",
	Addressing: map[Mode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xe9, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xe5, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xf5, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xed, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xfd, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xf9, Size: 3},
		IndirectXAddressing: {Opcode: 0xe1, Size: 2},
		IndirectYAddressing: {Opcode: 0xf1, Size: 2},
	},
	ParamFunc: sbc,
}

// Sec - Set Carry Flag.
var Sec = &Instruction{
	Name: "sec",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x38, Size: 1},
	},
	NoParamFunc: sec,
}

// Sed - Set Decimal Flag.
var Sed = &Instruction{
	Name: "sed",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xf8, Size: 1},
	},
	NoParamFunc: sed,
}

// Sei - Set Interrupt Disable.
var Sei = &Instruction{
	Name: "sei",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x78, Size: 1},
	},
	NoParamFunc: sei,
}

// Sta - Store Accumulator.
var Sta = &Instruction{
	Name: "sta",
	Addressing: map[Mode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x85, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x95, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x9d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x99, Size: 3},
		IndirectXAddressing: {Opcode: 0x81, Size: 2},
		IndirectYAddressing: {Opcode: 0x91, Size: 2},
	},
	ParamFunc: sta,
}

// Stx - Store X Register.
var Stx = &Instruction{
	Name: "stx",
	Addressing: map[Mode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x86, Size: 2},
		ZeroPageYAddressing: {Opcode: 0x96, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8e, Size: 3},
	},
	ParamFunc: stx,
}

// Sty - Store Y Register.
var Sty = &Instruction{
	Name: "sty",
	Addressing: map[Mode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x84, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x94, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8c, Size: 3},
	},
	ParamFunc: sty,
}

// Tax - Transfer Accumulator to X.
var Tax = &Instruction{
	Name: "tax",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xaa, Size: 1},
	},
	NoParamFunc: tax,
}

// Tay - Transfer Accumulator to Y.
var Tay = &Instruction{
	Name: "tay",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xa8, Size: 1},
	},
	NoParamFunc: tay,
}

// Tsx - Transfer Stack Pointer to X.
var Tsx = &Instruction{
	Name: "tsx",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xba, Size: 1},
	},
	NoParamFunc: tsx,
}

// Txa - Transfer X to Accumulator.
var Txa = &Instruction{
	Name: "txa",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x8a, Size: 1},
	},
	NoParamFunc: txa,
}

// Txs - Transfer X to Stack Pointer.
var Txs = &Instruction{
	Name: "txs",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x9a, Size: 1},
	},
	NoParamFunc: txs,
}

// Tya - Transfer Y to Accumulator.
var Tya = &Instruction{
	Name: "tya",
	Addressing: map[Mode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x98, Size: 1},
	},
	NoParamFunc: tya,
}

// Instructions maps instruction names to their information struct.
var Instructions = map[string]*Instruction{
	"adc": Adc,
	"and": And,
	"asl": Asl,
	"bcc": Bcc,
	"bcs": Bcs,
	"beq": Beq,
	"bit": Bit,
	"bmi": Bmi,
	"bne": Bne,
	"bpl": Bpl,
	"brk": Brk,
	"bvc": Bvc,
	"bvs": Bvs,
	"clc": Clc,
	"cld": Cld,
	"cli": Cli,
	"clv": Clv,
	"cmp": Cmp,
	"cpx": Cpx,
	"cpy": Cpy,
	"dcp": Dcp,
	"dec": Dec,
	"dex": Dex,
	"dey": Dey,
	"eor": Eor,
	"inc": Inc,
	"inx": Inx,
	"iny": Iny,
	"isc": Isc,
	"jmp": Jmp,
	"jsr": Jsr,
	"lax": Lax,
	"lda": Lda,
	"ldx": Ldx,
	"ldy": Ldy,
	"lsr": Lsr,
	"nop": Nop,
	"ora": Ora,
	"pha": Pha,
	"php": Php,
	"pla": Pla,
	"plp": Plp,
	"rla": Rla,
	"rol": Rol,
	"ror": Ror,
	"rra": Rra,
	"rti": Rti,
	"rts": Rts,
	"sax": Sax,
	"sbc": Sbc,
	"sec": Sec,
	"sed": Sed,
	"sei": Sei,
	"slo": Slo,
	"sre": Sre,
	"sta": Sta,
	"stx": Stx,
	"sty": Sty,
	"tax": Tax,
	"tay": Tay,
	"tsx": Tsx,
	"txa": Txa,
	"txs": Txs,
	"tya": Tya,
}
