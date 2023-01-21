package cpu

import . "github.com/retroenv/retrogolib/nes/addressing"

// AddressingInfo contains the opcode and timing info for an instruction addressing mode.
type AddressingInfo struct {
	Opcode byte
}

// Instruction contains information about a NES CPU instruction.
type Instruction struct {
	Name       string
	Unofficial bool

	// instruction has no parameters
	NoParamFunc func()
	// instruction has parameters
	ParamFunc func(params ...any)

	// maps addressing mode to cpu cycles
	Addressing map[Mode]AddressingInfo
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
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0x69},
		ZeroPageAddressing:  {Opcode: 0x65},
		ZeroPageXAddressing: {Opcode: 0x75},
		AbsoluteAddressing:  {Opcode: 0x6d},
		AbsoluteXAddressing: {Opcode: 0x7d},
		AbsoluteYAddressing: {Opcode: 0x79},
		IndirectXAddressing: {Opcode: 0x61},
		IndirectYAddressing: {Opcode: 0x71},
	},
}

// And - AND with accumulator.
var And = &Instruction{
	Name: "and",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0x29},
		ZeroPageAddressing:  {Opcode: 0x25},
		ZeroPageXAddressing: {Opcode: 0x35},
		AbsoluteAddressing:  {Opcode: 0x2d},
		AbsoluteXAddressing: {Opcode: 0x3d},
		AbsoluteYAddressing: {Opcode: 0x39},
		IndirectXAddressing: {Opcode: 0x21},
		IndirectYAddressing: {Opcode: 0x31},
	},
}

// Asl - Arithmetic Shift Left.
var Asl = &Instruction{
	Name: "asl",
	Addressing: map[Mode]AddressingInfo{
		AccumulatorAddressing: {Opcode: 0x0a},
		ZeroPageAddressing:    {Opcode: 0x06},
		ZeroPageXAddressing:   {Opcode: 0x16},
		AbsoluteAddressing:    {Opcode: 0x0e},
		AbsoluteXAddressing:   {Opcode: 0x1e},
	},
}

// Bcc - Branch if Carry Clear.
var Bcc = &Instruction{
	Name: "bcc",
	Addressing: map[Mode]AddressingInfo{
		RelativeAddressing: {Opcode: 0x90},
	},
}

// Bcs - Branch if Carry Set.
var Bcs = &Instruction{
	Name: "bcs",
	Addressing: map[Mode]AddressingInfo{
		RelativeAddressing: {Opcode: 0xb0},
	},
}

// Beq - Branch if Equal.
var Beq = &Instruction{
	Name: "beq",
	Addressing: map[Mode]AddressingInfo{
		RelativeAddressing: {Opcode: 0xf0},
	},
}

// Bit - Bit Test.
var Bit = &Instruction{
	Name: "bit",
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing: {Opcode: 0x24},
		AbsoluteAddressing: {Opcode: 0x2c},
	},
}

// Bmi - Branch if Minus.
var Bmi = &Instruction{
	Name: "bmi",
	Addressing: map[Mode]AddressingInfo{
		RelativeAddressing: {Opcode: 0x30},
	},
}

// Bne - Branch if Not Equal.
var Bne = &Instruction{
	Name: "bne",
	Addressing: map[Mode]AddressingInfo{
		RelativeAddressing: {Opcode: 0xd0},
	},
}

// Bpl - Branch if Positive.
var Bpl = &Instruction{
	Name: "bpl",
	Addressing: map[Mode]AddressingInfo{
		RelativeAddressing: {Opcode: 0x10},
	},
}

// Brk - Force Interrupt.
var Brk = &Instruction{
	Name: "brk",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x00},
	},
}

// Bvc - Branch if Overflow Clear.
var Bvc = &Instruction{
	Name: "bvc",
	Addressing: map[Mode]AddressingInfo{
		RelativeAddressing: {Opcode: 0x50},
	},
}

// Bvs - Branch if Overflow Set.
var Bvs = &Instruction{
	Name: "bvs",
	Addressing: map[Mode]AddressingInfo{
		RelativeAddressing: {Opcode: 0x70},
	},
}

// Clc - Clear Carry Flag.
var Clc = &Instruction{
	Name: "clc",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x18},
	},
}

// Cld - Clear Decimal Mode.
var Cld = &Instruction{
	Name: "cld",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xd8},
	},
}

// Cli - Clear Interrupt Disable.
var Cli = &Instruction{
	Name: "cli",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x58},
	},
}

// Clv - Clear Overflow Flag.
var Clv = &Instruction{
	Name: "clv",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xb8},
	},
}

// Cmp - Compare - compares the contents of A.
var Cmp = &Instruction{
	Name: "cmp",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0xc9},
		ZeroPageAddressing:  {Opcode: 0xc5},
		ZeroPageXAddressing: {Opcode: 0xd5},
		AbsoluteAddressing:  {Opcode: 0xcd},
		AbsoluteXAddressing: {Opcode: 0xdd},
		AbsoluteYAddressing: {Opcode: 0xd9},
		IndirectXAddressing: {Opcode: 0xc1},
		IndirectYAddressing: {Opcode: 0xd1},
	},
}

// Cpx - Compare X Register - compares the contents of X.
var Cpx = &Instruction{
	Name: "cpx",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0xe0},
		ZeroPageAddressing:  {Opcode: 0xe4},
		AbsoluteAddressing:  {Opcode: 0xec},
	},
}

// Cpy - Compare Y Register - compares the contents of Y.
var Cpy = &Instruction{
	Name: "cpy",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0xc0},
		ZeroPageAddressing:  {Opcode: 0xc4},
		AbsoluteAddressing:  {Opcode: 0xcc},
	},
}

// Dec - Decrement memory.
var Dec = &Instruction{
	Name: "Dec",
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0xc6},
		ZeroPageXAddressing: {Opcode: 0xd6},
		AbsoluteAddressing:  {Opcode: 0xce},
		AbsoluteXAddressing: {Opcode: 0xde},
	},
}

// Dex - Decrement X Register.
var Dex = &Instruction{
	Name: "dex",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xca},
	},
}

// Dey - Decrement Y Register.
var Dey = &Instruction{
	Name: "dey",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x88},
	},
}

// Eor - Exclusive OR - XOR.
var Eor = &Instruction{
	Name: "eor",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0x49},
		ZeroPageAddressing:  {Opcode: 0x45},
		ZeroPageXAddressing: {Opcode: 0x55},
		AbsoluteAddressing:  {Opcode: 0x4d},
		AbsoluteXAddressing: {Opcode: 0x5d},
		AbsoluteYAddressing: {Opcode: 0x59},
		IndirectXAddressing: {Opcode: 0x41},
		IndirectYAddressing: {Opcode: 0x51},
	},
}

// Inc - Increments memory.
var Inc = &Instruction{
	Name: "inc",
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0xe6},
		ZeroPageXAddressing: {Opcode: 0xf6},
		AbsoluteAddressing:  {Opcode: 0xee},
		AbsoluteXAddressing: {Opcode: 0xfe},
	},
}

// Inx - Increment X Register.
var Inx = &Instruction{
	Name: "inx",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xe8},
	},
}

// Iny - Increment Y Register.
var Iny = &Instruction{
	Name: "iny",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xc8},
	},
}

// Jmp - jump to address.
var Jmp = &Instruction{
	Name: "jmp",
	Addressing: map[Mode]AddressingInfo{
		AbsoluteAddressing: {Opcode: 0x4c},
		IndirectAddressing: {Opcode: 0x6c},
	},
}

// Jsr - jump to subroutine.
var Jsr = &Instruction{
	Name: "jsr",
	Addressing: map[Mode]AddressingInfo{
		AbsoluteAddressing: {Opcode: 0x20},
	},
}

// Lda - Load Accumulator - load a byte into A.
var Lda = &Instruction{
	Name: "lda",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0xa9},
		ZeroPageAddressing:  {Opcode: 0xa5},
		ZeroPageXAddressing: {Opcode: 0xb5},
		AbsoluteAddressing:  {Opcode: 0xad},
		AbsoluteXAddressing: {Opcode: 0xbd},
		AbsoluteYAddressing: {Opcode: 0xb9},
		IndirectXAddressing: {Opcode: 0xa1},
		IndirectYAddressing: {Opcode: 0xb1},
	},
}

// Ldx - Load X Register - load a byte into X.
var Ldx = &Instruction{
	Name: "ldx",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0xa2},
		ZeroPageAddressing:  {Opcode: 0xa6},
		ZeroPageYAddressing: {Opcode: 0xb6},
		AbsoluteAddressing:  {Opcode: 0xae},
		AbsoluteYAddressing: {Opcode: 0xbe},
	},
}

// Ldy - Load Y Register - load a byte into Y.
var Ldy = &Instruction{
	Name: "ldy",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0xa0},
		ZeroPageAddressing:  {Opcode: 0xa4},
		ZeroPageXAddressing: {Opcode: 0xb4},
		AbsoluteAddressing:  {Opcode: 0xac},
		AbsoluteXAddressing: {Opcode: 0xbc},
	},
}

// Lsr - Logical Shift Right.
var Lsr = &Instruction{
	Name: "lsr",
	Addressing: map[Mode]AddressingInfo{
		AccumulatorAddressing: {Opcode: 0x4a},
		ZeroPageAddressing:    {Opcode: 0x46},
		ZeroPageXAddressing:   {Opcode: 0x56},
		AbsoluteAddressing:    {Opcode: 0x4e},
		AbsoluteXAddressing:   {Opcode: 0x5e},
	},
}

// Nop - No Operation.
var Nop = &Instruction{
	Name: "nop",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xea},
	},
}

// Ora - OR with Accumulator.
var Ora = &Instruction{
	Name: "ora",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0x09},
		ZeroPageAddressing:  {Opcode: 0x05},
		ZeroPageXAddressing: {Opcode: 0x15},
		AbsoluteAddressing:  {Opcode: 0x0d},
		AbsoluteXAddressing: {Opcode: 0x1d},
		AbsoluteYAddressing: {Opcode: 0x19},
		IndirectXAddressing: {Opcode: 0x01},
		IndirectYAddressing: {Opcode: 0x11},
	},
}

// Pha - Push Accumulator - push A content to stack.
var Pha = &Instruction{
	Name: "pha",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x48},
	},
}

// Php - Push Processor Status - push status flags to stack.
var Php = &Instruction{
	Name: "php",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x08},
	},
}

// Pla - Pull Accumulator - pull A content from stack.
var Pla = &Instruction{
	Name: "pla",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x68},
	},
}

// Plp - Pull Processor Status - pull status flags from stack.
var Plp = &Instruction{
	Name: "plp",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x28},
	},
}

// Rol - Rotate Left.
var Rol = &Instruction{
	Name: "rol",
	Addressing: map[Mode]AddressingInfo{
		AccumulatorAddressing: {Opcode: 0x2a},
		ZeroPageAddressing:    {Opcode: 0x26},
		ZeroPageXAddressing:   {Opcode: 0x36},
		AbsoluteAddressing:    {Opcode: 0x2e},
		AbsoluteXAddressing:   {Opcode: 0x3e},
	},
}

// Ror - Rotate Right.
var Ror = &Instruction{
	Name: "ror",
	Addressing: map[Mode]AddressingInfo{
		AccumulatorAddressing: {Opcode: 0x6a},
		ZeroPageAddressing:    {Opcode: 0x66},
		ZeroPageXAddressing:   {Opcode: 0x76},
		AbsoluteAddressing:    {Opcode: 0x6e},
		AbsoluteXAddressing:   {Opcode: 0x7e},
	},
}

// Rti - Return from Interrupt.
var Rti = &Instruction{
	Name: "rti",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x40},
	},
}

// Rts - return from subroutine.
var Rts = &Instruction{
	Name: "rts",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x60},
	},
}

// Sbc - subtract with Carry.
var Sbc = &Instruction{
	Name: "sbc",
	Addressing: map[Mode]AddressingInfo{
		ImmediateAddressing: {Opcode: 0xe9},
		ZeroPageAddressing:  {Opcode: 0xe5},
		ZeroPageXAddressing: {Opcode: 0xf5},
		AbsoluteAddressing:  {Opcode: 0xed},
		AbsoluteXAddressing: {Opcode: 0xfd},
		AbsoluteYAddressing: {Opcode: 0xf9},
		IndirectXAddressing: {Opcode: 0xe1},
		IndirectYAddressing: {Opcode: 0xf1},
	},
}

// Sec - Set Carry Flag.
var Sec = &Instruction{
	Name: "sec",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x38},
	},
}

// Sed - Set Decimal Flag.
var Sed = &Instruction{
	Name: "sed",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xf8},
	},
}

// Sei - Set Interrupt Disable.
var Sei = &Instruction{
	Name: "sei",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x78},
	},
}

// Sta - Store Accumulator.
var Sta = &Instruction{
	Name: "sta",
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x85},
		ZeroPageXAddressing: {Opcode: 0x95},
		AbsoluteAddressing:  {Opcode: 0x8d},
		AbsoluteXAddressing: {Opcode: 0x9d},
		AbsoluteYAddressing: {Opcode: 0x99},
		IndirectXAddressing: {Opcode: 0x81},
		IndirectYAddressing: {Opcode: 0x91},
	},
}

// Stx - Store X Register.
var Stx = &Instruction{
	Name: "stx",
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x86},
		ZeroPageYAddressing: {Opcode: 0x96},
		AbsoluteAddressing:  {Opcode: 0x8e},
	},
}

// Sty - Store Y Register.
var Sty = &Instruction{
	Name: "sty",
	Addressing: map[Mode]AddressingInfo{
		ZeroPageAddressing:  {Opcode: 0x84},
		ZeroPageXAddressing: {Opcode: 0x94},
		AbsoluteAddressing:  {Opcode: 0x8c},
	},
}

// Tax - Transfer Accumulator to X.
var Tax = &Instruction{
	Name: "tax",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xaa},
	},
}

// Tay - Transfer Accumulator to Y.
var Tay = &Instruction{
	Name: "tay",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xa8},
	},
}

// Tsx - Transfer Stack Pointer to X.
var Tsx = &Instruction{
	Name: "tsx",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0xba},
	},
}

// Txa - Transfer X to Accumulator.
var Txa = &Instruction{
	Name: "txa",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x8a},
	},
}

// Txs - Transfer X to Stack Pointer.
var Txs = &Instruction{
	Name: "txs",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x9a},
	},
}

// Tya - Transfer Y to Accumulator.
var Tya = &Instruction{
	Name: "tya",
	Addressing: map[Mode]AddressingInfo{
		ImpliedAddressing: {Opcode: 0x98},
	},
}

// Instructions maps instruction names to NES CPU instruction information.
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
	"Dec": Dec,
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
