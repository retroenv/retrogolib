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
	AlrName = "alr" // Unofficial
	AncName = "anc" // Unofficial
	AndName = "and"
	AneName = "ane" // Unofficial - highly unstable (0x8B)
	ArrName = "arr" // Unofficial
	AslName = "asl"
	AxsName = "axs" // Unofficial
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
	KilName = "kil" // Unofficial - halts the CPU (0x02, 0x12, 0x22, 0x32, 0x42, 0x52, 0x62, 0x72, 0x92, 0xB2, 0xD2, 0xF2)
	LasName = "las" // Unofficial
	LaxName = "lax" // Unofficial
	LdaName = "lda"
	LdxName = "ldx"
	LdyName = "ldy"
	LsrName = "lsr"
	LxaName = "lxa" // Unofficial - highly unstable (0xAB)
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
	ShaName = "sha" // Unofficial - unstable (0x93, 0x9F)
	ShxName = "shx" // Unofficial - unstable (0x9E)
	ShyName = "shy" // Unofficial - unstable (0x9C)
	SloName = "slo" // Unofficial
	SreName = "sre" // Unofficial
	StaName = "sta"
	StxName = "stx"
	StyName = "sty"
	TasName = "tas" // Unofficial - unstable (0x9B)
	TaxName = "tax"
	TayName = "tay"
	TsxName = "tsx"
	TxaName = "txa"
	TxsName = "txs"
	TyaName = "tya"
)

// AdcInst - Add with Carry.
var AdcInst = &Instruction{
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

// AndInst - AND with accumulator.
var AndInst = &Instruction{
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

// AslInst - Arithmetic Shift Left.
var AslInst = &Instruction{
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

// BccInst - Branch if Carry Clear.
var BccInst = &Instruction{
	Name: BccName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x90, Size: 2},
	},
	ParamFunc: bcc,
}

// BcsInst - Branch if Carry Set.
var BcsInst = &Instruction{
	Name: BcsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0xb0, Size: 2},
	},
	ParamFunc: bcs,
}

// BeqInst - Branch if Equal.
var BeqInst = &Instruction{
	Name: BeqName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0xf0, Size: 2},
	},
	ParamFunc: beq,
}

// BitInst - BitInst Test.
var BitInst = &Instruction{
	Name: BitName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing: {Opcode: 0x24, Size: 2},
		AbsoluteAddressing: {Opcode: 0x2c, Size: 3},
	},
	ParamFunc: bit,
}

// BmiInst - Branch if Minus.
var BmiInst = &Instruction{
	Name: BmiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x30, Size: 2},
	},
	ParamFunc: bmi,
}

// BneInst - Branch if Not Equal.
var BneInst = &Instruction{
	Name: BneName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0xd0, Size: 2},
	},
	ParamFunc: bne,
}

// BplInst - Branch if Positive.
var BplInst = &Instruction{
	Name: BplName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x10, Size: 2},
	},
	ParamFunc: bpl,
}

// BrkInst - Force Interrupt.
var BrkInst = &Instruction{
	Name: BrkName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x00, Size: 2},
	},
	NoParamFunc: brk,
}

// BvcInst - Branch if Overflow Clear.
var BvcInst = &Instruction{
	Name: BvcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x50, Size: 2},
	},
	ParamFunc: bvc,
}

// BvsInst - Branch if Overflow Set.
var BvsInst = &Instruction{
	Name: BvsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RelativeAddressing: {Opcode: 0x70, Size: 2},
	},
	ParamFunc: bvs,
}

// ClcInst - Clear Carry Flag.
var ClcInst = &Instruction{
	Name: ClcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x18, Size: 1},
	},
	NoParamFunc: clc,
}

// CldInst - Clear Decimal Mode.
var CldInst = &Instruction{
	Name: CldName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xd8, Size: 1},
	},
	NoParamFunc: cld,
}

// CliInst - Clear Interrupt Disable.
var CliInst = &Instruction{
	Name: CliName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x58, Size: 1},
	},
	NoParamFunc: cli,
}

// ClvInst - Clear Overflow Flag.
var ClvInst = &Instruction{
	Name: ClvName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xb8, Size: 1},
	},
	NoParamFunc: clv,
}

// CmpInst - Compare the contents of A.
var CmpInst = &Instruction{
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

// CpxInst - Compare the contents of X.
var CpxInst = &Instruction{
	Name: CpxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xe0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xe4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xec, Size: 3},
	},
	ParamFunc: cpx,
}

// CpyInst - Compare the contents of Y.
var CpyInst = &Instruction{
	Name: CpyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xc0, Size: 2},
		ZeroPageAddressing:  {Opcode: 0xc4, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xcc, Size: 3},
	},
	ParamFunc: cpy,
}

// DecInst - Decrement memory.
var DecInst = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xc6, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xd6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xce, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xde, Size: 3},
	},
	ParamFunc: dec,
}

// DexInst - Decrement X Register.
var DexInst = &Instruction{
	Name: DexName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xca, Size: 1},
	},
	NoParamFunc: dex,
}

// DeyInst - Decrement Y Register.
var DeyInst = &Instruction{
	Name: DeyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x88, Size: 1},
	},
	NoParamFunc: dey,
}

// EorInst - Exclusive OR - XOR.
var EorInst = &Instruction{
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

// IncInst - Increments memory.
var IncInst = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0xe6, Size: 2},
		ZeroPageXAddressing: {Opcode: 0xf6, Size: 2},
		AbsoluteAddressing:  {Opcode: 0xee, Size: 3},
		AbsoluteXAddressing: {Opcode: 0xfe, Size: 3},
	},
	ParamFunc: inc,
}

// InxInst - Increment X Register.
var InxInst = &Instruction{
	Name: InxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xe8, Size: 1},
	},
	NoParamFunc: inx,
}

// InyInst - Increment Y Register.
var InyInst = &Instruction{
	Name: InyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xc8, Size: 1},
	},
	NoParamFunc: iny,
}

// JmpInst - jump to address.
var JmpInst = &Instruction{
	Name: JmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing: {Opcode: 0x4c, Size: 3},
		IndirectAddressing: {Opcode: 0x6c},
	},
	ParamFunc: jmp,
}

// JsrInst - jump to subroutine.
var JsrInst = &Instruction{
	Name: JsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing: {Opcode: 0x20, Size: 3},
	},
	ParamFunc: jsr,
}

// KilInst - Kill/Jam: halts the CPU. Unofficial opcode that freezes the 6502.
// The test-visible effect is that PC advances by 1 (past the opcode byte).
var KilInst = &Instruction{
	Name: KilName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x02, Size: 1},
	},
	NoParamFunc: kil,
}

// LdaInst - Load Accumulator - load a byte into A.
var LdaInst = &Instruction{
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

// LdxInst - Load X Register - load a byte into X.
var LdxInst = &Instruction{
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

// LdyInst - Load Y Register - load a byte into Y.
var LdyInst = &Instruction{
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

// LsrInst - Logical Shift Right.
var LsrInst = &Instruction{
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

// NopInst - No Operation.
var NopInst = &Instruction{
	Name: NopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xea, Size: 1},
	},
	NoParamFunc: nop,
}

// OraInst - OR with Accumulator.
var OraInst = &Instruction{
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

// PhaInst - Push Accumulator - push A content to stack.
var PhaInst = &Instruction{
	Name: PhaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x48, Size: 1},
	},
	NoParamFunc: pha,
}

// PhpInst - Push Processor Status - push status flags to stack.
var PhpInst = &Instruction{
	Name: PhpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x08, Size: 1},
	},
	NoParamFunc: php,
}

// PlaInst - Pull Accumulator - pull A content from stack.
var PlaInst = &Instruction{
	Name: PlaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x68, Size: 1},
	},
	NoParamFunc: pla,
}

// PlpInst - Pull Processor Status - pull status flags from stack.
var PlpInst = &Instruction{
	Name: PlpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x28, Size: 1},
	},
	NoParamFunc: plp,
}

// RolInst - Rotate Left.
var RolInst = &Instruction{
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

// RorInst - Rotate Right.
var RorInst = &Instruction{
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

// RtiInst - Return from Interrupt.
var RtiInst = &Instruction{
	Name: RtiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x40, Size: 1},
	},
	NoParamFunc: rti,
}

// RtsInst - return from subroutine.
var RtsInst = &Instruction{
	Name: RtsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x60, Size: 1},
	},
	NoParamFunc: rts,
}

// SbcInst - subtract with Carry.
var SbcInst = &Instruction{
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

// SecInst - Set Carry Flag.
var SecInst = &Instruction{
	Name: SecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x38, Size: 1},
	},
	NoParamFunc: sec,
}

// SedInst - Set Decimal Flag.
var SedInst = &Instruction{
	Name: SedName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xf8, Size: 1},
	},
	NoParamFunc: sed,
}

// SeiInst - Set Interrupt Disable.
var SeiInst = &Instruction{
	Name: SeiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x78, Size: 1},
	},
	NoParamFunc: sei,
}

// StaInst - Store Accumulator.
var StaInst = &Instruction{
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

// StxInst - Store X Register.
var StxInst = &Instruction{
	Name: StxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x86, Size: 2},
		ZeroPageYAddressing: {Opcode: 0x96, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8e, Size: 3},
	},
	ParamFunc: stx,
}

// StyInst - Store Y Register.
var StyInst = &Instruction{
	Name: StyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ZeroPageAddressing:  {Opcode: 0x84, Size: 2},
		ZeroPageXAddressing: {Opcode: 0x94, Size: 2},
		AbsoluteAddressing:  {Opcode: 0x8c, Size: 3},
	},
	ParamFunc: sty,
}

// TaxInst - Transfer Accumulator to X.
var TaxInst = &Instruction{
	Name: TaxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xaa, Size: 1},
	},
	NoParamFunc: tax,
}

// TayInst - Transfer Accumulator to Y.
var TayInst = &Instruction{
	Name: TayName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xa8, Size: 1},
	},
	NoParamFunc: tay,
}

// TsxInst - Transfer Stack Pointer to X.
var TsxInst = &Instruction{
	Name: TsxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xba, Size: 1},
	},
	NoParamFunc: tsx,
}

// TxaInst - Transfer X to Accumulator.
var TxaInst = &Instruction{
	Name: TxaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x8a, Size: 1},
	},
	NoParamFunc: txa,
}

// TxsInst - Transfer X to Stack Pointer.
var TxsInst = &Instruction{
	Name: TxsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x9a, Size: 1},
	},
	NoParamFunc: txs,
}

// TyaInst - Transfer Y to Accumulator.
var TyaInst = &Instruction{
	Name: TyaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x98, Size: 1},
	},
	NoParamFunc: tya,
}

// Instructions maps instruction names to their information struct.
var Instructions = map[string]*Instruction{
	AdcName: AdcInst,
	AndName: AndInst,
	AneName: AneInst,
	AslName: AslInst,
	BccName: BccInst,
	BcsName: BcsInst,
	BeqName: BeqInst,
	BitName: BitInst,
	BmiName: BmiInst,
	BneName: BneInst,
	BplName: BplInst,
	BraName: BraInst,
	BrkName: BrkInst,
	BvcName: BvcInst,
	BvsName: BvsInst,
	ClcName: ClcInst,
	CldName: CldInst,
	CliName: CliInst,
	ClvName: ClvInst,
	CmpName: CmpInst,
	CpxName: CpxInst,
	CpyName: CpyInst,
	DcpName: DcpInst,
	DecName: DecInst,
	DexName: DexInst,
	DeyName: DeyInst,
	EorName: EorInst,
	IncName: IncInst,
	InxName: InxInst,
	InyName: InyInst,
	IscName: IscInst,
	JmpName: JmpInst,
	JsrName: JsrInst,
	KilName: KilInst,
	LasName: LasInst,
	LaxName: LaxInst,
	LdaName: LdaInst,
	LdxName: LdxInst,
	LdyName: LdyInst,
	LsrName: LsrInst,
	LxaName: LxaInst,
	NopName: NopInst,
	OraName: OraInst,
	PhaName: PhaInst,
	PhpName: PhpInst,
	PhxName: PhxInst,
	PhyName: PhyInst,
	PlaName: PlaInst,
	PlpName: PlpInst,
	PlxName: PlxInst,
	PlyName: PlyInst,
	RlaName: RlaInst,
	RolName: RolInst,
	RorName: RorInst,
	RraName: RraInst,
	RtiName: RtiInst,
	RtsName: RtsInst,
	SaxName: SaxInst,
	SbcName: SbcInst,
	SecName: SecInst,
	SedName: SedInst,
	SeiName: SeiInst,
	ShaName: ShaInst,
	ShxName: ShxInst,
	ShyName: ShyInst,
	SloName: SloInst,
	SreName: SreInst,
	StaName: StaInst,
	StxName: StxInst,
	StyName: StyInst,
	StzName: StzInst,
	TasName: TasInst,
	TaxName: TaxInst,
	TayName: TayInst,
	TsxName: TsxInst,
	TxaName: TxaInst,
	TrbName: TrbInst,
	TsbName: TsbInst,
	TxsName: TxsInst,
	TyaName: TyaInst,
}

// InstructionsByID maps OpcodeID to *Instruction for O(1) lookup by numeric ID.
// Index 0 (InvalidOpcodeID) is nil. Alr/Anc/Arr/Axs have OpcodeIDs but no Instructions entry.
var InstructionsByID = [OpcodeIDMax + 1]*Instruction{
	Adc: AdcInst,
	And: AndInst,
	Ane: AneInst,
	Asl: AslInst,
	Bcc: BccInst,
	Bcs: BcsInst,
	Beq: BeqInst,
	Bit: BitInst,
	Bmi: BmiInst,
	Bne: BneInst,
	Bpl: BplInst,
	Bra: BraInst,
	Brk: BrkInst,
	Bvc: BvcInst,
	Bvs: BvsInst,
	Clc: ClcInst,
	Cld: CldInst,
	Cli: CliInst,
	Clv: ClvInst,
	Cmp: CmpInst,
	Cpx: CpxInst,
	Cpy: CpyInst,
	Dcp: DcpInst,
	Dec: DecInst,
	Dex: DexInst,
	Dey: DeyInst,
	Eor: EorInst,
	Inc: IncInst,
	Inx: InxInst,
	Iny: InyInst,
	Isc: IscInst,
	Jmp: JmpInst,
	Jsr: JsrInst,
	Kil: KilInst,
	Las: LasInst,
	Lax: LaxInst,
	Lda: LdaInst,
	Ldx: LdxInst,
	Ldy: LdyInst,
	Lsr: LsrInst,
	Lxa: LxaInst,
	Nop: NopInst,
	Ora: OraInst,
	Pha: PhaInst,
	Php: PhpInst,
	Phx: PhxInst,
	Phy: PhyInst,
	Pla: PlaInst,
	Plp: PlpInst,
	Plx: PlxInst,
	Ply: PlyInst,
	Rla: RlaInst,
	Rol: RolInst,
	Ror: RorInst,
	Rra: RraInst,
	Rti: RtiInst,
	Rts: RtsInst,
	Sax: SaxInst,
	Sbc: SbcInst,
	Sec: SecInst,
	Sed: SedInst,
	Sei: SeiInst,
	Sha: ShaInst,
	Shx: ShxInst,
	Shy: ShyInst,
	Slo: SloInst,
	Sre: SreInst,
	Sta: StaInst,
	Stx: StxInst,
	Sty: StyInst,
	Stz: StzInst,
	Tas: TasInst,
	Tax: TaxInst,
	Tay: TayInst,
	Trb: TrbInst,
	Tsb: TsbInst,
	Tsx: TsxInst,
	Txa: TxaInst,
	Txs: TxsInst,
	Tya: TyaInst,
}
