package m6502

import (
	. "github.com/retroenv/retrogolib/addressing"
	"github.com/retroenv/retrogolib/cpu"
)

// Adc - Add with Carry.
var Adc = &cpu.Instruction{
	Name: "adc",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0x69, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x65, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x75, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x6d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x7d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x79, Size: 3},
		IndirectXAddressing: {Opcode: 0x61, Size: 2},
		IndirectYAddressing: {Opcode: 0x71, Size: 2},
	},
}

// And - AND with accumulator.
var And = &cpu.Instruction{
	Name: "and",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0x29, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x25, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x35, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x2d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x3d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x39, Size: 3},
		IndirectXAddressing: {Opcode: 0x21, Size: 2},
		IndirectYAddressing: {Opcode: 0x31, Size: 2},
	},
}

// Asl - Arithmetic Shift Left.
var Asl = &cpu.Instruction{
	Name: "asl",
	Addressing: map[Mode]cpu.AddressingInfo{
		AccumulatorAddressing: {Opcode: 0x0a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0x06, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0x16, Size: 2},
		AbsoluteAddressing:    {Opcode: 0x0e, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0x1e, Size: 3},
	},
}

// Bcc - Branch if Carry Clear.
var Bcc = &cpu.Instruction{
	Name: "bcc",
	Addressing: map[Mode]cpu.AddressingInfo{
		RelativeAddressing: {Opcode: 0x90, Size: 2},
	},
}

// Bcs - Branch if Carry Set.
var Bcs = &cpu.Instruction{
	Name: "bcs",
	Addressing: map[Mode]cpu.AddressingInfo{
		RelativeAddressing: {Opcode: 0xb0, Size: 2},
	},
}

// Beq - Branch if Equal.
var Beq = &cpu.Instruction{
	Name: "beq",
	Addressing: map[Mode]cpu.AddressingInfo{
		RelativeAddressing: {Opcode: 0xf0, Size: 2},
	},
}

// Bit - Bit Test.
var Bit = &cpu.Instruction{
	Name: "bit",
	Addressing: map[Mode]cpu.AddressingInfo{
		ZeroPageAddressing: {Opcode: 0x24, Size: 2},
		AbsoluteAddressing: {Opcode: 0x2c, Size: 3},
	},
}

// Bmi - Branch if Minus.
var Bmi = &cpu.Instruction{
	Name: "bmi",
	Addressing: map[Mode]cpu.AddressingInfo{
		RelativeAddressing: {Opcode: 0x30, Size: 2},
	},
}

// Bne - Branch if Not Equal.
var Bne = &cpu.Instruction{
	Name: "bne",
	Addressing: map[Mode]cpu.AddressingInfo{
		RelativeAddressing: {Opcode: 0xd0, Size: 2},
	},
}

// Bpl - Branch if Positive.
var Bpl = &cpu.Instruction{
	Name: "bpl",
	Addressing: map[Mode]cpu.AddressingInfo{
		RelativeAddressing: {Opcode: 0x10, Size: 2},
	},
}

// Brk - Force Interrupt.
var Brk = &cpu.Instruction{
	Name: "brk",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x00, Size: 1},
	},
}

// Bvc - Branch if Overflow Clear.
var Bvc = &cpu.Instruction{
	Name: "bvc",
	Addressing: map[Mode]cpu.AddressingInfo{
		RelativeAddressing: {Opcode: 0x50, Size: 2},
	},
}

// Bvs - Branch if Overflow Set.
var Bvs = &cpu.Instruction{
	Name: "bvs",
	Addressing: map[Mode]cpu.AddressingInfo{
		RelativeAddressing: {Opcode: 0x70, Size: 2},
	},
}

// Clc - Clear Carry Flag.
var Clc = &cpu.Instruction{
	Name: "clc",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x18, Size: 1},
	},
}

// Cld - Clear Decimal Mode.
var Cld = &cpu.Instruction{
	Name: "cld",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xd8, Size: 1},
	},
}

// Cli - Clear Interrupt Disable.
var Cli = &cpu.Instruction{
	Name: "cli",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x58, Size: 1},
	},
}

// Clv - Clear Overflow Flag.
var Clv = &cpu.Instruction{
	Name: "clv",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xb8, Size: 1},
	},
}

// Cmp - Compare - compares the contents of A.
var Cmp = &cpu.Instruction{
	Name: "cmp",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0xc9, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xc5, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xd5, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xcd, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xdd, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xd9, Size: 3},
		IndirectXAddressing: {Opcode: 0xc1, Size: 2},
		IndirectYAddressing: {Opcode: 0xd1, Size: 2},
	},
}

// Cpx - Compare X Register - compares the contents of X.
var Cpx = &cpu.Instruction{
	Name: "cpx",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0xe0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xe4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xec, Size: 3},
	},
}

// Cpy - Compare Y Register - compares the contents of Y.
var Cpy = &cpu.Instruction{
	Name: "cpy",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0xc0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xc4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xcc, Size: 3},
	},
}

// Dec - Decrement memory.
var Dec = &cpu.Instruction{
	Name: "dec",
	Addressing: map[Mode]cpu.AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0xc6, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xd6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xce, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xde, Size: 3},
	},
}

// Dex - Decrement X Register.
var Dex = &cpu.Instruction{
	Name: "dex",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xca, Size: 1},
	},
}

// Dey - Decrement Y Register.
var Dey = &cpu.Instruction{
	Name: "dey",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x88, Size: 1},
	},
}

// Eor - Exclusive OR - XOR.
var Eor = &cpu.Instruction{
	Name: "eor",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0x49, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x45, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x55, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x4d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x5d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x59, Size: 3},
		IndirectXAddressing: {Opcode: 0x41, Size: 2},
		IndirectYAddressing: {Opcode: 0x51, Size: 2},
	},
}

// Inc - Increments memory.
var Inc = &cpu.Instruction{
	Name: "inc",
	Addressing: map[Mode]cpu.AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0xe6, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xf6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xee, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xfe, Size: 3},
	},
}

// Inx - Increment X Register.
var Inx = &cpu.Instruction{
	Name: "inx",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xe8, Size: 1},
	},
}

// Iny - Increment Y Register.
var Iny = &cpu.Instruction{
	Name: "iny",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xc8, Size: 1},
	},
}

// Jmp - jump to address.
var Jmp = &cpu.Instruction{
	Name: "jmp",
	Addressing: map[Mode]cpu.AddressingInfo{
		AbsoluteAddressing: {Opcode: 0x4c, Size: 3},
		IndirectAddressing: {Opcode: 0x6c},
	},
}

// Jsr - jump to subroutine.
var Jsr = &cpu.Instruction{
	Name: "jsr",
	Addressing: map[Mode]cpu.AddressingInfo{
		AbsoluteAddressing: {Opcode: 0x20, Size: 3},
	},
}

// Lda - Load Accumulator - load a byte into A.
var Lda = &cpu.Instruction{
	Name: "lda",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0xa9, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xa5, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xb5, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xad, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xbd, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xb9, Size: 3},
		IndirectXAddressing: {Opcode: 0xa1, Size: 2},
		IndirectYAddressing: {Opcode: 0xb1, Size: 2},
	},
}

// Ldx - Load X Register - load a byte into X.
var Ldx = &cpu.Instruction{
	Name: "ldx",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0xa2, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xa6, Size: 2},
		ZeroPageYAddressing: {Opcode: 0xb6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xae, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xbe, Size: 3},
	},
}

// Ldy - Load Y Register - load a byte into Y.
var Ldy = &cpu.Instruction{
	Name: "ldy",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0xa0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xa4, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xb4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xac, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xbc, Size: 3},
	},
}

// Lsr - Logical Shift Right.
var Lsr = &cpu.Instruction{
	Name: "lsr",
	Addressing: map[Mode]cpu.AddressingInfo{
		AccumulatorAddressing: {Opcode: 0x4a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0x46, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0x56, Size: 2},
		AbsoluteAddressing:    {Opcode: 0x4e, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0x5e, Size: 3},
	},
}

// Nop - No Operation.
var Nop = &cpu.Instruction{
	Name: "nop",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xea, Size: 1},
	},
}

// Ora - OR with Accumulator.
var Ora = &cpu.Instruction{
	Name: "ora",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0x09, Size: 2},
		ZeroPageAddressing:  {Opcode: 0x05, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x15, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x0d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x1d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x19, Size: 3},
		IndirectXAddressing: {Opcode: 0x01, Size: 2},
		IndirectYAddressing: {Opcode: 0x11, Size: 2},
	},
}

// Pha - Push Accumulator - push A content to stack.
var Pha = &cpu.Instruction{
	Name: "pha",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x48, Size: 1},
	},
}

// Php - Push Processor Status - push status flags to stack.
var Php = &cpu.Instruction{
	Name: "php",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x08, Size: 1},
	},
}

// Pla - Pull Accumulator - pull A content from stack.
var Pla = &cpu.Instruction{
	Name: "pla",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x68, Size: 1},
	},
}

// Plp - Pull Processor Status - pull status flags from stack.
var Plp = &cpu.Instruction{
	Name: "plp",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x28, Size: 1},
	},
}

// Rol - Rotate Left.
var Rol = &cpu.Instruction{
	Name: "rol",
	Addressing: map[Mode]cpu.AddressingInfo{
		AccumulatorAddressing: {Opcode: 0x2a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0x26, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0x36, Size: 2},
		AbsoluteAddressing:    {Opcode: 0x2e, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0x3e, Size: 3},
	},
}

// Ror - Rotate Right.
var Ror = &cpu.Instruction{
	Name: "ror",
	Addressing: map[Mode]cpu.AddressingInfo{
		AccumulatorAddressing: {Opcode: 0x6a, Size: 1},
		ZeroPageAddressing:    {Opcode: 0x66, Size: 2},
		ZeroPageXAddressing:   {Opcode: 0x76, Size: 2},
		AbsoluteAddressing:    {Opcode: 0x6e, Size: 3},
		AbsoluteXAddressing:   {Opcode: 0x7e, Size: 3},
	},
}

// Rti - Return from Interrupt.
var Rti = &cpu.Instruction{
	Name: "rti",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x40, Size: 1},
	},
}

// Rts - return from subroutine.
var Rts = &cpu.Instruction{
	Name: "rts",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x60, Size: 1},
	},
}

// Sbc - subtract with Carry.
var Sbc = &cpu.Instruction{
	Name: "sbc",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImmediateAddressing: {Opcode: 0xe9, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xe5, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xf5, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xed, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xfd, Size: 3},
		AbsoluteYAddressing: {Opcode: 0xf9, Size: 3},
		IndirectXAddressing: {Opcode: 0xe1, Size: 2},
		IndirectYAddressing: {Opcode: 0xf1, Size: 2},
	},
}

// Sec - Set Carry Flag.
var Sec = &cpu.Instruction{
	Name: "sec",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x38, Size: 1},
	},
}

// Sed - Set Decimal Flag.
var Sed = &cpu.Instruction{
	Name: "sed",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xf8, Size: 1},
	},
}

// Sei - Set Interrupt Disable.
var Sei = &cpu.Instruction{
	Name: "sei",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x78, Size: 1},
	},
}

// Sta - Store Accumulator.
var Sta = &cpu.Instruction{
	Name: "sta",
	Addressing: map[Mode]cpu.AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x85, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x95, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8d, Size: 3},
		AbsoluteXAddressing: {Opcode: 0x9d, Size: 3},
		AbsoluteYAddressing: {Opcode: 0x99, Size: 3},
		IndirectXAddressing: {Opcode: 0x81, Size: 2},
		IndirectYAddressing: {Opcode: 0x91, Size: 2},
	},
}

// Stx - Store X Register.
var Stx = &cpu.Instruction{
	Name: "stx",
	Addressing: map[Mode]cpu.AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x86, Size: 2},
		ZeroPageYAddressing: {Opcode: 0x96, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8e, Size: 3},
	},
}

// Sty - Store Y Register.
var Sty = &cpu.Instruction{
	Name: "sty",
	Addressing: map[Mode]cpu.AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x84, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x94, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8c, Size: 3},
	},
}

// Tax - Transfer Accumulator to X.
var Tax = &cpu.Instruction{
	Name: "tax",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xaa, Size: 1},
	},
}

// Tay - Transfer Accumulator to Y.
var Tay = &cpu.Instruction{
	Name: "tay",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xa8, Size: 1},
	},
}

// Tsx - Transfer Stack Pointer to X.
var Tsx = &cpu.Instruction{
	Name: "tsx",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0xba, Size: 1},
	},
}

// Txa - Transfer X to Accumulator.
var Txa = &cpu.Instruction{
	Name: "txa",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x8a, Size: 1},
	},
}

// Txs - Transfer X to Stack Pointer.
var Txs = &cpu.Instruction{
	Name: "txs",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x9a, Size: 1},
	},
}

// Tya - Transfer Y to Accumulator.
var Tya = &cpu.Instruction{
	Name: "tya",
	Addressing: map[Mode]cpu.AddressingInfo{
		ImpliedAddressing: {Opcode: 0x98, Size: 1},
	},
}

// Instructions maps instruction names to NES CPU instruction information.
var Instructions = map[string]*cpu.Instruction{
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
