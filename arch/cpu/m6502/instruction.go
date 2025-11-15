package m6502

// Instruction defines a 6502 CPU instruction with its opcodes and execution logic.
// Instructions support multiple addressing modes through the Addressing map that
// enables opcode lookup for disassembly and code generation.
type Instruction struct {
	Name       string // Instruction mnemonic (lowercase)
	Unofficial bool   // True for undocumented opcodes not in original 6502 spec

	// Opcode lookup map for addressing mode to opcode mapping
	Addressing map[AddressingMode]OpcodeInfo // Maps addressing mode to opcode info

	// Execution handlers - exactly one must be set
	NoParamFunc func(c *CPU) error                // Handler for implied addressing (NOP, DEX)
	ParamFunc   func(c *CPU, params ...any) error // Handler for parameterized instructions (LDA, JMP)
}

// HasAddressing returns whether the instruction has any of the passed addressing modes.
func (ins Instruction) HasAddressing(flags ...AddressingMode) bool {
	for _, flag := range flags {
		_, ok := ins.Addressing[flag]
		if ok {
			return ok
		}
	}
	return false
}

// Instruction name constants for easy access by external packages.
const (
	AdcName = "adc"
	AndName = "and"
	AslName = "asl"
	BccName = "bcc"
	BcsName = "bcs"
	BeqName = "beq"
	BitName = "bit"
	BmiName = "bmi"
	BneName = "bne"
	BplName = "bpl"
	BrkName = "brk"
	BvcName = "bvc"
	BvsName = "bvs"
	ClcName = "clc"
	CldName = "cld"
	CliName = "cli"
	ClvName = "clv"
	CmpName = "cmp"
	CpxName = "cpx"
	CpyName = "cpy"
	DcpName = "dcp" // Unofficial
	DecName = "dec"
	DexName = "dex"
	DeyName = "dey"
	EorName = "eor"
	IncName = "inc"
	InxName = "inx"
	InyName = "iny"
	IscName = "isc" // Unofficial
	JmpName = "jmp"
	JsrName = "jsr"
	LaxName = "lax" // Unofficial
	LdaName = "lda"
	LdxName = "ldx"
	LdyName = "ldy"
	LsrName = "lsr"
	NopName = "nop"
	OraName = "ora"
	PhaName = "pha"
	PhpName = "php"
	PlaName = "pla"
	PlpName = "plp"
	RlaName = "rla" // Unofficial
	RolName = "rol"
	RorName = "ror"
	RraName = "rra" // Unofficial
	RtiName = "rti"
	RtsName = "rts"
	SaxName = "sax" // Unofficial
	SbcName = "sbc"
	SecName = "sec"
	SedName = "sed"
	SeiName = "sei"
	SloName = "slo" // Unofficial
	SreName = "sre" // Unofficial
	StaName = "sta"
	StxName = "stx"
	StyName = "sty"
	TaxName = "tax"
	TayName = "tay"
	TsxName = "tsx"
	TxaName = "txa"
	TxsName = "txs"
	TyaName = "tya"
)

// Adc - Add with Carry.
var Adc = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: AndName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: AslName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: BccName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x90, Size: 2},
	},
	ParamFunc: bcc,
}

// Bcs - Branch if Carry Set.
var Bcs = &Instruction{
	Name: BcsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0xb0, Size: 2},
	},
	ParamFunc: bcs,
}

// Beq - Branch if Equal.
var Beq = &Instruction{
	Name: BeqName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0xf0, Size: 2},
	},
	ParamFunc: beq,
}

// Bit - Bit Test.
var Bit = &Instruction{
	Name: BitName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing: {Opcode: 0x24, Size: 2},
		AbsoluteAddressing: {Opcode: 0x2c, Size: 3},
	},
	ParamFunc: bit,
}

// Bmi - Branch if Minus.
var Bmi = &Instruction{
	Name: BmiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x30, Size: 2},
	},
	ParamFunc: bmi,
}

// Bne - Branch if Not Equal.
var Bne = &Instruction{
	Name: BneName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0xd0, Size: 2},
	},
	ParamFunc: bne,
}

// Bpl - Branch if Positive.
var Bpl = &Instruction{
	Name: BplName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x10, Size: 2},
	},
	ParamFunc: bpl,
}

// Brk - Force Interrupt.
var Brk = &Instruction{
	Name: BrkName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x00, Size: 2},
	},
	NoParamFunc: brk,
}

// Bvc - Branch if Overflow Clear.
var Bvc = &Instruction{
	Name: BvcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x50, Size: 2},
	},
	ParamFunc: bvc,
}

// Bvs - Branch if Overflow Set.
var Bvs = &Instruction{
	Name: BvsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x70, Size: 2},
	},
	ParamFunc: bvs,
}

// Clc - Clear Carry Flag.
var Clc = &Instruction{
	Name: ClcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x18, Size: 1},
	},
	NoParamFunc: clc,
}

// Cld - Clear Decimal Mode.
var Cld = &Instruction{
	Name: CldName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xd8, Size: 1},
	},
	NoParamFunc: cld,
}

// Cli - Clear Interrupt Disable.
var Cli = &Instruction{
	Name: CliName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x58, Size: 1},
	},
	NoParamFunc: cli,
}

// Clv - Clear Overflow Flag.
var Clv = &Instruction{
	Name: ClvName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xb8, Size: 1},
	},
	NoParamFunc: clv,
}

// Cmp - Compare the contents of A.
var Cmp = &Instruction{
	Name: CmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: CpxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xe0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xe4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xec, Size: 3},
	},
	ParamFunc: cpx,
}

// Cpy - Compare the contents of Y.
var Cpy = &Instruction{
	Name: CpyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xc0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xc4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xcc, Size: 3},
	},
	ParamFunc: cpy,
}

// Dec - Decrement memory.
var Dec = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xc6, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xd6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xce, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xde, Size: 3},
	},
	ParamFunc: dec,
}

// Dex - Decrement X Register.
var Dex = &Instruction{
	Name: DexName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xca, Size: 1},
	},
	NoParamFunc: dex,
}

// Dey - Decrement Y Register.
var Dey = &Instruction{
	Name: DeyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x88, Size: 1},
	},
	NoParamFunc: dey,
}

// Eor - Exclusive OR - XOR.
var Eor = &Instruction{
	Name: EorName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xe6, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xf6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xee, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xfe, Size: 3},
	},
	ParamFunc: inc,
}

// Inx - Increment X Register.
var Inx = &Instruction{
	Name: InxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xe8, Size: 1},
	},
	NoParamFunc: inx,
}

// Iny - Increment Y Register.
var Iny = &Instruction{
	Name: InyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xc8, Size: 1},
	},
	NoParamFunc: iny,
}

// Jmp - jump to address.
var Jmp = &Instruction{
	Name: JmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing: {Opcode: 0x4c, Size: 3},
		IndirectAddressing: {Opcode: 0x6c},
	},
	ParamFunc: jmp,
}

// Jsr - jump to subroutine.
var Jsr = &Instruction{
	Name: JsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing: {Opcode: 0x20, Size: 3},
	},
	ParamFunc: jsr,
}

// Lda - Load Accumulator - load a byte into A.
var Lda = &Instruction{
	Name: LdaName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: LdxName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: LdyName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: LsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: NopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xea, Size: 1},
	},
	NoParamFunc: nop,
}

// Ora - OR with Accumulator.
var Ora = &Instruction{
	Name: OraName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: PhaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x48, Size: 1},
	},
	NoParamFunc: pha,
}

// Php - Push Processor Status - push status flags to stack.
var Php = &Instruction{
	Name: PhpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x08, Size: 1},
	},
	NoParamFunc: php,
}

// Pla - Pull Accumulator - pull A content from stack.
var Pla = &Instruction{
	Name: PlaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x68, Size: 1},
	},
	NoParamFunc: pla,
}

// Plp - Pull Processor Status - pull status flags from stack.
var Plp = &Instruction{
	Name: PlpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x28, Size: 1},
	},
	NoParamFunc: plp,
}

// Rol - Rotate Left.
var Rol = &Instruction{
	Name: RolName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: RorName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: RtiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x40, Size: 1},
	},
	NoParamFunc: rti,
}

// Rts - return from subroutine.
var Rts = &Instruction{
	Name: RtsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x60, Size: 1},
	},
	NoParamFunc: rts,
}

// Sbc - subtract with Carry.
var Sbc = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: SecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x38, Size: 1},
	},
	NoParamFunc: sec,
}

// Sed - Set Decimal Flag.
var Sed = &Instruction{
	Name: SedName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xf8, Size: 1},
	},
	NoParamFunc: sed,
}

// Sei - Set Interrupt Disable.
var Sei = &Instruction{
	Name: SeiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x78, Size: 1},
	},
	NoParamFunc: sei,
}

// Sta - Store Accumulator.
var Sta = &Instruction{
	Name: StaName,
	Addressing: map[AddressingMode]OpcodeInfo{
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
	Name: StxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x86, Size: 2},
		ZeroPageYAddressing: {Opcode: 0x96, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8e, Size: 3},
	},
	ParamFunc: stx,
}

// Sty - Store Y Register.
var Sty = &Instruction{
	Name: StyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x84, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x94, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8c, Size: 3},
	},
	ParamFunc: sty,
}

// Tax - Transfer Accumulator to X.
var Tax = &Instruction{
	Name: TaxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xaa, Size: 1},
	},
	NoParamFunc: tax,
}

// Tay - Transfer Accumulator to Y.
var Tay = &Instruction{
	Name: TayName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xa8, Size: 1},
	},
	NoParamFunc: tay,
}

// Tsx - Transfer Stack Pointer to X.
var Tsx = &Instruction{
	Name: TsxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xba, Size: 1},
	},
	NoParamFunc: tsx,
}

// Txa - Transfer X to Accumulator.
var Txa = &Instruction{
	Name: TxaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x8a, Size: 1},
	},
	NoParamFunc: txa,
}

// Txs - Transfer X to Stack Pointer.
var Txs = &Instruction{
	Name: TxsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x9a, Size: 1},
	},
	NoParamFunc: txs,
}

// Tya - Transfer Y to Accumulator.
var Tya = &Instruction{
	Name: TyaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x98, Size: 1},
	},
	NoParamFunc: tya,
}

// Instructions maps instruction names to their information struct.
var Instructions = map[string]*Instruction{
	AdcName: Adc,
	AndName: And,
	AslName: Asl,
	BccName: Bcc,
	BcsName: Bcs,
	BeqName: Beq,
	BitName: Bit,
	BmiName: Bmi,
	BneName: Bne,
	BplName: Bpl,
	BrkName: Brk,
	BvcName: Bvc,
	BvsName: Bvs,
	ClcName: Clc,
	CldName: Cld,
	CliName: Cli,
	ClvName: Clv,
	CmpName: Cmp,
	CpxName: Cpx,
	CpyName: Cpy,
	DcpName: Dcp,
	DecName: Dec,
	DexName: Dex,
	DeyName: Dey,
	EorName: Eor,
	IncName: Inc,
	InxName: Inx,
	InyName: Iny,
	IscName: Isc,
	JmpName: Jmp,
	JsrName: Jsr,
	LaxName: Lax,
	LdaName: Lda,
	LdxName: Ldx,
	LdyName: Ldy,
	LsrName: Lsr,
	NopName: Nop,
	OraName: Ora,
	PhaName: Pha,
	PhpName: Php,
	PlaName: Pla,
	PlpName: Plp,
	RlaName: Rla,
	RolName: Rol,
	RorName: Ror,
	RraName: Rra,
	RtiName: Rti,
	RtsName: Rts,
	SaxName: Sax,
	SbcName: Sbc,
	SecName: Sec,
	SedName: Sed,
	SeiName: Sei,
	SloName: Slo,
	SreName: Sre,
	StaName: Sta,
	StxName: Stx,
	StyName: Sty,
	TaxName: Tax,
	TayName: Tay,
	TsxName: Tsx,
	TxaName: Txa,
	TxsName: Txs,
	TyaName: Tya,
}
