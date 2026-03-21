package m65816

// Instruction defines a 65816 CPU instruction.
type Instruction struct {
	Name       string
	Unofficial bool

	// Addressing maps each supported addressing mode to its opcode and size info.
	Addressing map[AddressingMode]OpcodeInfo

	// Exactly one of these must be set.
	NoParamFunc func(c *CPU) error
	ParamFunc   func(c *CPU, params ...any) error
}

// OpcodeInfo contains the opcode byte and base instruction size.
type OpcodeInfo struct {
	Opcode   byte
	BaseSize byte // size when M=1 / X=1 (8-bit mode); some instructions grow by 1 in 16-bit mode
}

// HasAddressing returns true if the instruction supports any of the given modes.
func (ins *Instruction) HasAddressing(modes ...AddressingMode) bool {
	for _, m := range modes {
		if _, ok := ins.Addressing[m]; ok {
			return true
		}
	}
	return false
}

// Instruction name constants (sorted alphabetically).
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
	BraName = "bra"
	BrkName = "brk"
	BrlName = "brl"
	BvcName = "bvc"
	BvsName = "bvs"
	ClcName = "clc"
	CldName = "cld"
	CliName = "cli"
	ClvName = "clv"
	CmpName = "cmp"
	CopName = "cop"
	CpxName = "cpx"
	CpyName = "cpy"
	DecName = "dec"
	DexName = "dex"
	DeyName = "dey"
	EorName = "eor"
	IncName = "inc"
	InxName = "inx"
	InyName = "iny"
	JmlName = "jml"
	JmpName = "jmp"
	JslName = "jsl"
	JsrName = "jsr"
	LdaName = "lda"
	LdxName = "ldx"
	LdyName = "ldy"
	LsrName = "lsr"
	MvnName = "mvn"
	MvpName = "mvp"
	NopName = "nop"
	OraName = "ora"
	PeaName = "pea"
	PeiName = "pei"
	PerName = "per"
	PhaName = "pha"
	PhbName = "phb"
	PhdName = "phd"
	PhkName = "phk"
	PhpName = "php"
	PhxName = "phx"
	PhyName = "phy"
	PlaName = "pla"
	PlbName = "plb"
	PldName = "pld"
	PlpName = "plp"
	PlxName = "plx"
	PlyName = "ply"
	RepName = "rep"
	RolName = "rol"
	RorName = "ror"
	RtiName = "rti"
	RtlName = "rtl"
	RtsName = "rts"
	SbcName = "sbc"
	SecName = "sec"
	SedName = "sed"
	SeiName = "sei"
	SepName = "sep"
	StaName = "sta"
	StpName = "stp"
	StxName = "stx"
	StyName = "sty"
	StzName = "stz"
	TaxName = "tax"
	TayName = "tay"
	TcdName = "tcd"
	TcsName = "tcs"
	TdcName = "tdc"
	TrbName = "trb"
	TsbName = "tsb"
	TscName = "tsc"
	TsxName = "tsx"
	TxaName = "txa"
	TxsName = "txs"
	TxyName = "txy"
	TyaName = "tya"
	TyxName = "tyx"
	WaiName = "wai"
	WdmName = "wdm"
	XbaName = "xba"
	XceName = "xce"
)

// -- Instruction variable definitions --

// AdcInst - Add with Carry.
var AdcInst = &Instruction{
	Name: AdcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:                      {Opcode: 0x69, BaseSize: 2},
		DirectPageAddressing:                     {Opcode: 0x65, BaseSize: 2},
		DirectPageIndexedXAddressing:             {Opcode: 0x75, BaseSize: 2},
		DirectPageIndirectAddressing:             {Opcode: 0x72, BaseSize: 2},
		DirectPageIndexedXIndirectAddressing:     {Opcode: 0x61, BaseSize: 2},
		DirectPageIndirectIndexedYAddressing:     {Opcode: 0x71, BaseSize: 2},
		DirectPageIndirectLongAddressing:         {Opcode: 0x67, BaseSize: 2},
		DirectPageIndirectLongIndexedYAddressing: {Opcode: 0x77, BaseSize: 2},
		AbsoluteAddressing:                       {Opcode: 0x6D, BaseSize: 3},
		AbsoluteIndexedXAddressing:               {Opcode: 0x7D, BaseSize: 3},
		AbsoluteIndexedYAddressing:               {Opcode: 0x79, BaseSize: 3},
		AbsoluteLongAddressing:                   {Opcode: 0x6F, BaseSize: 4},
		AbsoluteLongIndexedXAddressing:           {Opcode: 0x7F, BaseSize: 4},
		StackRelativeAddressing:                  {Opcode: 0x63, BaseSize: 2},
		StackRelativeIndirectIndexedYAddressing:  {Opcode: 0x73, BaseSize: 2},
	},
	ParamFunc: adc,
}

// AndInst - AND with Accumulator.
var AndInst = &Instruction{
	Name: AndName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:                      {Opcode: 0x29, BaseSize: 2},
		DirectPageAddressing:                     {Opcode: 0x25, BaseSize: 2},
		DirectPageIndexedXAddressing:             {Opcode: 0x35, BaseSize: 2},
		DirectPageIndirectAddressing:             {Opcode: 0x32, BaseSize: 2},
		DirectPageIndexedXIndirectAddressing:     {Opcode: 0x21, BaseSize: 2},
		DirectPageIndirectIndexedYAddressing:     {Opcode: 0x31, BaseSize: 2},
		DirectPageIndirectLongAddressing:         {Opcode: 0x27, BaseSize: 2},
		DirectPageIndirectLongIndexedYAddressing: {Opcode: 0x37, BaseSize: 2},
		AbsoluteAddressing:                       {Opcode: 0x2D, BaseSize: 3},
		AbsoluteIndexedXAddressing:               {Opcode: 0x3D, BaseSize: 3},
		AbsoluteIndexedYAddressing:               {Opcode: 0x39, BaseSize: 3},
		AbsoluteLongAddressing:                   {Opcode: 0x2F, BaseSize: 4},
		AbsoluteLongIndexedXAddressing:           {Opcode: 0x3F, BaseSize: 4},
		StackRelativeAddressing:                  {Opcode: 0x23, BaseSize: 2},
		StackRelativeIndirectIndexedYAddressing:  {Opcode: 0x33, BaseSize: 2},
	},
	ParamFunc: and,
}

// AslInst - Arithmetic Shift Left.
var AslInst = &Instruction{
	Name: AslName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AccumulatorAddressing:        {Opcode: 0x0A, BaseSize: 1},
		DirectPageAddressing:         {Opcode: 0x06, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0x16, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x0E, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0x1E, BaseSize: 3},
	},
	ParamFunc: asl,
}

// BccInst - Branch if Carry Clear.
var BccInst = &Instruction{
	Name:       BccName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x90, BaseSize: 2}},
	ParamFunc:  bcc,
}

// BcsInst - Branch if Carry Set.
var BcsInst = &Instruction{
	Name:       BcsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0xB0, BaseSize: 2}},
	ParamFunc:  bcs,
}

// BeqInst - Branch if Equal (Z=1).
var BeqInst = &Instruction{
	Name:       BeqName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0xF0, BaseSize: 2}},
	ParamFunc:  beq,
}

// BitInst - Bit Test.
var BitInst = &Instruction{
	Name: BitName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:          {Opcode: 0x89, BaseSize: 2},
		DirectPageAddressing:         {Opcode: 0x24, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0x34, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x2C, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0x3C, BaseSize: 3},
	},
	ParamFunc: bit,
}

// BmiInst - Branch if Minus (N=1).
var BmiInst = &Instruction{
	Name:       BmiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x30, BaseSize: 2}},
	ParamFunc:  bmi,
}

// BneInst - Branch if Not Equal (Z=0).
var BneInst = &Instruction{
	Name:       BneName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0xD0, BaseSize: 2}},
	ParamFunc:  bne,
}

// BplInst - Branch if Positive (N=0).
var BplInst = &Instruction{
	Name:       BplName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x10, BaseSize: 2}},
	ParamFunc:  bpl,
}

// BraInst - Branch Always.
var BraInst = &Instruction{
	Name:       BraName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x80, BaseSize: 2}},
	ParamFunc:  bra,
}

// BrkInst - Software Interrupt / Break.
var BrkInst = &Instruction{
	Name:        BrkName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x00, BaseSize: 2}},
	NoParamFunc: brk,
}

// BrlInst - Branch Long (16-bit offset).
var BrlInst = &Instruction{
	Name:       BrlName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Opcode: 0x82, BaseSize: 3}},
	ParamFunc:  brl,
}

// BvcInst - Branch if Overflow Clear.
var BvcInst = &Instruction{
	Name:       BvcName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x50, BaseSize: 2}},
	ParamFunc:  bvc,
}

// BvsInst - Branch if Overflow Set.
var BvsInst = &Instruction{
	Name:       BvsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x70, BaseSize: 2}},
	ParamFunc:  bvs,
}

// ClcInst - Clear Carry Flag.
var ClcInst = &Instruction{
	Name:        ClcName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x18, BaseSize: 1}},
	NoParamFunc: clc,
}

// CldInst - Clear Decimal Flag.
var CldInst = &Instruction{
	Name:        CldName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xD8, BaseSize: 1}},
	NoParamFunc: cld,
}

// CliInst - Clear Interrupt Disable.
var CliInst = &Instruction{
	Name:        CliName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x58, BaseSize: 1}},
	NoParamFunc: cli,
}

// ClvInst - Clear Overflow Flag.
var ClvInst = &Instruction{
	Name:        ClvName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xB8, BaseSize: 1}},
	NoParamFunc: clv,
}

// CmpInst - Compare Accumulator.
var CmpInst = &Instruction{
	Name: CmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:                      {Opcode: 0xC9, BaseSize: 2},
		DirectPageAddressing:                     {Opcode: 0xC5, BaseSize: 2},
		DirectPageIndexedXAddressing:             {Opcode: 0xD5, BaseSize: 2},
		DirectPageIndirectAddressing:             {Opcode: 0xD2, BaseSize: 2},
		DirectPageIndexedXIndirectAddressing:     {Opcode: 0xC1, BaseSize: 2},
		DirectPageIndirectIndexedYAddressing:     {Opcode: 0xD1, BaseSize: 2},
		DirectPageIndirectLongAddressing:         {Opcode: 0xC7, BaseSize: 2},
		DirectPageIndirectLongIndexedYAddressing: {Opcode: 0xD7, BaseSize: 2},
		AbsoluteAddressing:                       {Opcode: 0xCD, BaseSize: 3},
		AbsoluteIndexedXAddressing:               {Opcode: 0xDD, BaseSize: 3},
		AbsoluteIndexedYAddressing:               {Opcode: 0xD9, BaseSize: 3},
		AbsoluteLongAddressing:                   {Opcode: 0xCF, BaseSize: 4},
		AbsoluteLongIndexedXAddressing:           {Opcode: 0xDF, BaseSize: 4},
		StackRelativeAddressing:                  {Opcode: 0xC3, BaseSize: 2},
		StackRelativeIndirectIndexedYAddressing:  {Opcode: 0xD3, BaseSize: 2},
	},
	ParamFunc: cmp,
}

// CopInst - Co-Processor Enable (software interrupt).
var CopInst = &Instruction{
	Name:        CopName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x02, BaseSize: 2}},
	NoParamFunc: cop,
}

// CpxInst - Compare X Register.
var CpxInst = &Instruction{
	Name: CpxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:  {Opcode: 0xE0, BaseSize: 2},
		DirectPageAddressing: {Opcode: 0xE4, BaseSize: 2},
		AbsoluteAddressing:   {Opcode: 0xEC, BaseSize: 3},
	},
	ParamFunc: cpx,
}

// CpyInst - Compare Y Register.
var CpyInst = &Instruction{
	Name: CpyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:  {Opcode: 0xC0, BaseSize: 2},
		DirectPageAddressing: {Opcode: 0xC4, BaseSize: 2},
		AbsoluteAddressing:   {Opcode: 0xCC, BaseSize: 3},
	},
	ParamFunc: cpy,
}

// DecInst - Decrement.
var DecInst = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AccumulatorAddressing:        {Opcode: 0x3A, BaseSize: 1},
		DirectPageAddressing:         {Opcode: 0xC6, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0xD6, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0xCE, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0xDE, BaseSize: 3},
	},
	ParamFunc: dec,
}

// DexInst - Decrement X.
var DexInst = &Instruction{
	Name:        DexName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xCA, BaseSize: 1}},
	NoParamFunc: dex,
}

// DeyInst - Decrement Y.
var DeyInst = &Instruction{
	Name:        DeyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x88, BaseSize: 1}},
	NoParamFunc: dey,
}

// EorInst - Exclusive OR with Accumulator.
var EorInst = &Instruction{
	Name: EorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:                      {Opcode: 0x49, BaseSize: 2},
		DirectPageAddressing:                     {Opcode: 0x45, BaseSize: 2},
		DirectPageIndexedXAddressing:             {Opcode: 0x55, BaseSize: 2},
		DirectPageIndirectAddressing:             {Opcode: 0x52, BaseSize: 2},
		DirectPageIndexedXIndirectAddressing:     {Opcode: 0x41, BaseSize: 2},
		DirectPageIndirectIndexedYAddressing:     {Opcode: 0x51, BaseSize: 2},
		DirectPageIndirectLongAddressing:         {Opcode: 0x47, BaseSize: 2},
		DirectPageIndirectLongIndexedYAddressing: {Opcode: 0x57, BaseSize: 2},
		AbsoluteAddressing:                       {Opcode: 0x4D, BaseSize: 3},
		AbsoluteIndexedXAddressing:               {Opcode: 0x5D, BaseSize: 3},
		AbsoluteIndexedYAddressing:               {Opcode: 0x59, BaseSize: 3},
		AbsoluteLongAddressing:                   {Opcode: 0x4F, BaseSize: 4},
		AbsoluteLongIndexedXAddressing:           {Opcode: 0x5F, BaseSize: 4},
		StackRelativeAddressing:                  {Opcode: 0x43, BaseSize: 2},
		StackRelativeIndirectIndexedYAddressing:  {Opcode: 0x53, BaseSize: 2},
	},
	ParamFunc: eor,
}

// IncInst - Increment.
var IncInst = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AccumulatorAddressing:        {Opcode: 0x1A, BaseSize: 1},
		DirectPageAddressing:         {Opcode: 0xE6, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0xF6, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0xEE, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0xFE, BaseSize: 3},
	},
	ParamFunc: inc,
}

// InxInst - Increment X.
var InxInst = &Instruction{
	Name:        InxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xE8, BaseSize: 1}},
	NoParamFunc: inx,
}

// InyInst - Increment Y.
var InyInst = &Instruction{
	Name:        InyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xC8, BaseSize: 1}},
	NoParamFunc: iny,
}

// JmlInst - Jump Long (sets PB).
var JmlInst = &Instruction{
	Name: JmlName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteLongAddressing:         {Opcode: 0x5C, BaseSize: 4},
		AbsoluteIndirectLongAddressing: {Opcode: 0xDC, BaseSize: 3},
	},
	ParamFunc: jml,
}

// JmpInst - Jump.
var JmpInst = &Instruction{
	Name: JmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing:                 {Opcode: 0x4C, BaseSize: 3},
		AbsoluteIndirectAddressing:         {Opcode: 0x6C, BaseSize: 3},
		AbsoluteIndexedXIndirectAddressing: {Opcode: 0x7C, BaseSize: 3},
	},
	ParamFunc: jmp,
}

// JslInst - Jump to Subroutine Long.
var JslInst = &Instruction{
	Name:       JslName,
	Addressing: map[AddressingMode]OpcodeInfo{AbsoluteLongAddressing: {Opcode: 0x22, BaseSize: 4}},
	ParamFunc:  jsl,
}

// JsrInst - Jump to Subroutine.
var JsrInst = &Instruction{
	Name: JsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing:                 {Opcode: 0x20, BaseSize: 3},
		AbsoluteIndexedXIndirectAddressing: {Opcode: 0xFC, BaseSize: 3},
	},
	ParamFunc: jsr,
}

// LdaInst - Load Accumulator.
var LdaInst = &Instruction{
	Name: LdaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:                      {Opcode: 0xA9, BaseSize: 2},
		DirectPageAddressing:                     {Opcode: 0xA5, BaseSize: 2},
		DirectPageIndexedXAddressing:             {Opcode: 0xB5, BaseSize: 2},
		DirectPageIndirectAddressing:             {Opcode: 0xB2, BaseSize: 2},
		DirectPageIndexedXIndirectAddressing:     {Opcode: 0xA1, BaseSize: 2},
		DirectPageIndirectIndexedYAddressing:     {Opcode: 0xB1, BaseSize: 2},
		DirectPageIndirectLongAddressing:         {Opcode: 0xA7, BaseSize: 2},
		DirectPageIndirectLongIndexedYAddressing: {Opcode: 0xB7, BaseSize: 2},
		AbsoluteAddressing:                       {Opcode: 0xAD, BaseSize: 3},
		AbsoluteIndexedXAddressing:               {Opcode: 0xBD, BaseSize: 3},
		AbsoluteIndexedYAddressing:               {Opcode: 0xB9, BaseSize: 3},
		AbsoluteLongAddressing:                   {Opcode: 0xAF, BaseSize: 4},
		AbsoluteLongIndexedXAddressing:           {Opcode: 0xBF, BaseSize: 4},
		StackRelativeAddressing:                  {Opcode: 0xA3, BaseSize: 2},
		StackRelativeIndirectIndexedYAddressing:  {Opcode: 0xB3, BaseSize: 2},
	},
	ParamFunc: lda,
}

// LdxInst - Load X Register.
var LdxInst = &Instruction{
	Name: LdxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:          {Opcode: 0xA2, BaseSize: 2},
		DirectPageAddressing:         {Opcode: 0xA6, BaseSize: 2},
		DirectPageIndexedYAddressing: {Opcode: 0xB6, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0xAE, BaseSize: 3},
		AbsoluteIndexedYAddressing:   {Opcode: 0xBE, BaseSize: 3},
	},
	ParamFunc: ldx,
}

// LdyInst - Load Y Register.
var LdyInst = &Instruction{
	Name: LdyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:          {Opcode: 0xA0, BaseSize: 2},
		DirectPageAddressing:         {Opcode: 0xA4, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0xB4, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0xAC, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0xBC, BaseSize: 3},
	},
	ParamFunc: ldy,
}

// LsrInst - Logical Shift Right.
var LsrInst = &Instruction{
	Name: LsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AccumulatorAddressing:        {Opcode: 0x4A, BaseSize: 1},
		DirectPageAddressing:         {Opcode: 0x46, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0x56, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x4E, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0x5E, BaseSize: 3},
	},
	ParamFunc: lsr,
}

// MvnInst - Move Block Next (increment).
var MvnInst = &Instruction{
	Name:       MvnName,
	Addressing: map[AddressingMode]OpcodeInfo{BlockMoveAddressing: {Opcode: 0x54, BaseSize: 3}},
	ParamFunc:  mvn,
}

// MvpInst - Move Block Previous (decrement).
var MvpInst = &Instruction{
	Name:       MvpName,
	Addressing: map[AddressingMode]OpcodeInfo{BlockMoveAddressing: {Opcode: 0x44, BaseSize: 3}},
	ParamFunc:  mvp,
}

// NopInst - No Operation.
var NopInst = &Instruction{
	Name:        NopName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xEA, BaseSize: 1}},
	NoParamFunc: nop,
}

// OraInst - OR with Accumulator.
var OraInst = &Instruction{
	Name: OraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:                      {Opcode: 0x09, BaseSize: 2},
		DirectPageAddressing:                     {Opcode: 0x05, BaseSize: 2},
		DirectPageIndexedXAddressing:             {Opcode: 0x15, BaseSize: 2},
		DirectPageIndirectAddressing:             {Opcode: 0x12, BaseSize: 2},
		DirectPageIndexedXIndirectAddressing:     {Opcode: 0x01, BaseSize: 2},
		DirectPageIndirectIndexedYAddressing:     {Opcode: 0x11, BaseSize: 2},
		DirectPageIndirectLongAddressing:         {Opcode: 0x07, BaseSize: 2},
		DirectPageIndirectLongIndexedYAddressing: {Opcode: 0x17, BaseSize: 2},
		AbsoluteAddressing:                       {Opcode: 0x0D, BaseSize: 3},
		AbsoluteIndexedXAddressing:               {Opcode: 0x1D, BaseSize: 3},
		AbsoluteIndexedYAddressing:               {Opcode: 0x19, BaseSize: 3},
		AbsoluteLongAddressing:                   {Opcode: 0x0F, BaseSize: 4},
		AbsoluteLongIndexedXAddressing:           {Opcode: 0x1F, BaseSize: 4},
		StackRelativeAddressing:                  {Opcode: 0x03, BaseSize: 2},
		StackRelativeIndirectIndexedYAddressing:  {Opcode: 0x13, BaseSize: 2},
	},
	ParamFunc: ora,
}

// PeaInst - Push Effective Absolute Address.
var PeaInst = &Instruction{
	Name:        PeaName,
	Addressing:  map[AddressingMode]OpcodeInfo{AbsoluteAddressing: {Opcode: 0xF4, BaseSize: 3}},
	NoParamFunc: pea,
}

// PeiInst - Push Effective Indirect Address.
var PeiInst = &Instruction{
	Name:        PeiName,
	Addressing:  map[AddressingMode]OpcodeInfo{DirectPageIndirectAddressing: {Opcode: 0xD4, BaseSize: 2}},
	NoParamFunc: pei,
}

// PerInst - Push Effective Relative Address.
var PerInst = &Instruction{
	Name:        PerName,
	Addressing:  map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Opcode: 0x62, BaseSize: 3}},
	NoParamFunc: per,
}

// PhaInst - Push Accumulator.
var PhaInst = &Instruction{
	Name:        PhaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x48, BaseSize: 1}},
	NoParamFunc: pha,
}

// PhbInst - Push Data Bank Register.
var PhbInst = &Instruction{
	Name:        PhbName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x8B, BaseSize: 1}},
	NoParamFunc: phb,
}

// PhdInst - Push Direct Page Register.
var PhdInst = &Instruction{
	Name:        PhdName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x0B, BaseSize: 1}},
	NoParamFunc: phd,
}

// PhkInst - Push Program Bank Register.
var PhkInst = &Instruction{
	Name:        PhkName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x4B, BaseSize: 1}},
	NoParamFunc: phk,
}

// PhpInst - Push Processor Status.
var PhpInst = &Instruction{
	Name:        PhpName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x08, BaseSize: 1}},
	NoParamFunc: php,
}

// PhxInst - Push X Register.
var PhxInst = &Instruction{
	Name:        PhxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xDA, BaseSize: 1}},
	NoParamFunc: phx,
}

// PhyInst - Push Y Register.
var PhyInst = &Instruction{
	Name:        PhyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x5A, BaseSize: 1}},
	NoParamFunc: phy,
}

// PlaInst - Pull Accumulator.
var PlaInst = &Instruction{
	Name:        PlaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x68, BaseSize: 1}},
	NoParamFunc: pla,
}

// PlbInst - Pull Data Bank Register.
var PlbInst = &Instruction{
	Name:        PlbName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xAB, BaseSize: 1}},
	NoParamFunc: plb,
}

// PldInst - Pull Direct Page Register.
var PldInst = &Instruction{
	Name:        PldName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x2B, BaseSize: 1}},
	NoParamFunc: pld,
}

// PlpInst - Pull Processor Status.
var PlpInst = &Instruction{
	Name:        PlpName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x28, BaseSize: 1}},
	NoParamFunc: plp,
}

// PlxInst - Pull X Register.
var PlxInst = &Instruction{
	Name:        PlxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xFA, BaseSize: 1}},
	NoParamFunc: plx,
}

// PlyInst - Pull Y Register.
var PlyInst = &Instruction{
	Name:        PlyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x7A, BaseSize: 1}},
	NoParamFunc: ply,
}

// RepInst - Reset Processor Status Bits.
var RepInst = &Instruction{
	Name:       RepName,
	Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0xC2, BaseSize: 2}},
	ParamFunc:  rep,
}

// RolInst - Rotate Left.
var RolInst = &Instruction{
	Name: RolName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AccumulatorAddressing:        {Opcode: 0x2A, BaseSize: 1},
		DirectPageAddressing:         {Opcode: 0x26, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0x36, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x2E, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0x3E, BaseSize: 3},
	},
	ParamFunc: rol,
}

// RorInst - Rotate Right.
var RorInst = &Instruction{
	Name: RorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AccumulatorAddressing:        {Opcode: 0x6A, BaseSize: 1},
		DirectPageAddressing:         {Opcode: 0x66, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0x76, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x6E, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0x7E, BaseSize: 3},
	},
	ParamFunc: ror,
}

// RtiInst - Return from Interrupt.
var RtiInst = &Instruction{
	Name:        RtiName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x40, BaseSize: 1}},
	NoParamFunc: rti,
}

// RtlInst - Return from Subroutine Long.
var RtlInst = &Instruction{
	Name:        RtlName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x6B, BaseSize: 1}},
	NoParamFunc: rtl,
}

// RtsInst - Return from Subroutine.
var RtsInst = &Instruction{
	Name:        RtsName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x60, BaseSize: 1}},
	NoParamFunc: rts,
}

// SbcInst - Subtract with Carry.
var SbcInst = &Instruction{
	Name: SbcName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:                      {Opcode: 0xE9, BaseSize: 2},
		DirectPageAddressing:                     {Opcode: 0xE5, BaseSize: 2},
		DirectPageIndexedXAddressing:             {Opcode: 0xF5, BaseSize: 2},
		DirectPageIndirectAddressing:             {Opcode: 0xF2, BaseSize: 2},
		DirectPageIndexedXIndirectAddressing:     {Opcode: 0xE1, BaseSize: 2},
		DirectPageIndirectIndexedYAddressing:     {Opcode: 0xF1, BaseSize: 2},
		DirectPageIndirectLongAddressing:         {Opcode: 0xE7, BaseSize: 2},
		DirectPageIndirectLongIndexedYAddressing: {Opcode: 0xF7, BaseSize: 2},
		AbsoluteAddressing:                       {Opcode: 0xED, BaseSize: 3},
		AbsoluteIndexedXAddressing:               {Opcode: 0xFD, BaseSize: 3},
		AbsoluteIndexedYAddressing:               {Opcode: 0xF9, BaseSize: 3},
		AbsoluteLongAddressing:                   {Opcode: 0xEF, BaseSize: 4},
		AbsoluteLongIndexedXAddressing:           {Opcode: 0xFF, BaseSize: 4},
		StackRelativeAddressing:                  {Opcode: 0xE3, BaseSize: 2},
		StackRelativeIndirectIndexedYAddressing:  {Opcode: 0xF3, BaseSize: 2},
	},
	ParamFunc: sbc,
}

// SecInst - Set Carry Flag.
var SecInst = &Instruction{
	Name:        SecName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x38, BaseSize: 1}},
	NoParamFunc: sec,
}

// SedInst - Set Decimal Flag.
var SedInst = &Instruction{
	Name:        SedName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xF8, BaseSize: 1}},
	NoParamFunc: sed,
}

// SeiInst - Set Interrupt Disable.
var SeiInst = &Instruction{
	Name:        SeiName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x78, BaseSize: 1}},
	NoParamFunc: sei,
}

// SepInst - Set Processor Status Bits.
var SepInst = &Instruction{
	Name:       SepName,
	Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0xE2, BaseSize: 2}},
	ParamFunc:  sep,
}

// StaInst - Store Accumulator.
var StaInst = &Instruction{
	Name: StaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing:                     {Opcode: 0x85, BaseSize: 2},
		DirectPageIndexedXAddressing:             {Opcode: 0x95, BaseSize: 2},
		DirectPageIndirectAddressing:             {Opcode: 0x92, BaseSize: 2},
		DirectPageIndexedXIndirectAddressing:     {Opcode: 0x81, BaseSize: 2},
		DirectPageIndirectIndexedYAddressing:     {Opcode: 0x91, BaseSize: 2},
		DirectPageIndirectLongAddressing:         {Opcode: 0x87, BaseSize: 2},
		DirectPageIndirectLongIndexedYAddressing: {Opcode: 0x97, BaseSize: 2},
		AbsoluteAddressing:                       {Opcode: 0x8D, BaseSize: 3},
		AbsoluteIndexedXAddressing:               {Opcode: 0x9D, BaseSize: 3},
		AbsoluteIndexedYAddressing:               {Opcode: 0x99, BaseSize: 3},
		AbsoluteLongAddressing:                   {Opcode: 0x8F, BaseSize: 4},
		AbsoluteLongIndexedXAddressing:           {Opcode: 0x9F, BaseSize: 4},
		StackRelativeAddressing:                  {Opcode: 0x83, BaseSize: 2},
		StackRelativeIndirectIndexedYAddressing:  {Opcode: 0x93, BaseSize: 2},
	},
	ParamFunc: sta,
}

// StpInst - Stop the Processor.
var StpInst = &Instruction{
	Name:        StpName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xDB, BaseSize: 1}},
	NoParamFunc: stp,
}

// StxInst - Store X Register.
var StxInst = &Instruction{
	Name: StxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing:         {Opcode: 0x86, BaseSize: 2},
		DirectPageIndexedYAddressing: {Opcode: 0x96, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x8E, BaseSize: 3},
	},
	ParamFunc: stx,
}

// StyInst - Store Y Register.
var StyInst = &Instruction{
	Name: StyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing:         {Opcode: 0x84, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0x94, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x8C, BaseSize: 3},
	},
	ParamFunc: sty,
}

// StzInst - Store Zero.
var StzInst = &Instruction{
	Name: StzName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing:         {Opcode: 0x64, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0x74, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x9C, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0x9E, BaseSize: 3},
	},
	ParamFunc: stz,
}

// TaxInst - Transfer A to X.
var TaxInst = &Instruction{
	Name:        TaxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xAA, BaseSize: 1}},
	NoParamFunc: tax,
}

// TayInst - Transfer A to Y.
var TayInst = &Instruction{
	Name:        TayName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xA8, BaseSize: 1}},
	NoParamFunc: tay,
}

// TcdInst - Transfer C to Direct Page.
var TcdInst = &Instruction{
	Name:        TcdName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x5B, BaseSize: 1}},
	NoParamFunc: tcd,
}

// TcsInst - Transfer C to Stack Pointer.
var TcsInst = &Instruction{
	Name:        TcsName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x1B, BaseSize: 1}},
	NoParamFunc: tcs,
}

// TdcInst - Transfer Direct Page to C.
var TdcInst = &Instruction{
	Name:        TdcName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x7B, BaseSize: 1}},
	NoParamFunc: tdc,
}

// TrbInst - Test and Reset Bits.
var TrbInst = &Instruction{
	Name: TrbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing: {Opcode: 0x14, BaseSize: 2},
		AbsoluteAddressing:   {Opcode: 0x1C, BaseSize: 3},
	},
	ParamFunc: trb,
}

// TsbInst - Test and Set Bits.
var TsbInst = &Instruction{
	Name: TsbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing: {Opcode: 0x04, BaseSize: 2},
		AbsoluteAddressing:   {Opcode: 0x0C, BaseSize: 3},
	},
	ParamFunc: tsb,
}

// TscInst - Transfer Stack Pointer to C.
var TscInst = &Instruction{
	Name:        TscName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x3B, BaseSize: 1}},
	NoParamFunc: tsc,
}

// TsxInst - Transfer SP to X.
var TsxInst = &Instruction{
	Name:        TsxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xBA, BaseSize: 1}},
	NoParamFunc: tsx,
}

// TxaInst - Transfer X to A.
var TxaInst = &Instruction{
	Name:        TxaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x8A, BaseSize: 1}},
	NoParamFunc: txa,
}

// TxsInst - Transfer X to SP.
var TxsInst = &Instruction{
	Name:        TxsName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x9A, BaseSize: 1}},
	NoParamFunc: txs,
}

// TxyInst - Transfer X to Y.
var TxyInst = &Instruction{
	Name:        TxyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x9B, BaseSize: 1}},
	NoParamFunc: txy,
}

// TyaInst - Transfer Y to A.
var TyaInst = &Instruction{
	Name:        TyaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x98, BaseSize: 1}},
	NoParamFunc: tya,
}

// TyxInst - Transfer Y to X.
var TyxInst = &Instruction{
	Name:        TyxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xBB, BaseSize: 1}},
	NoParamFunc: tyx,
}

// WaiInst - Wait for Interrupt.
var WaiInst = &Instruction{
	Name:        WaiName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xCB, BaseSize: 1}},
	NoParamFunc: wai,
}

// WdmInst - Reserved/WDM (2-byte NOP).
var WdmInst = &Instruction{
	Name:       WdmName,
	Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x42, BaseSize: 2}},
	ParamFunc:  wdm,
}

// XbaInst - Exchange B and A (swap accumulator bytes).
var XbaInst = &Instruction{
	Name:        XbaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xEB, BaseSize: 1}},
	NoParamFunc: xba,
}

// XceInst - Exchange Carry and Emulation flags.
var XceInst = &Instruction{
	Name:        XceName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xFB, BaseSize: 1}},
	NoParamFunc: xce,
}

// Instructions maps instruction names to their definitions.
var Instructions = map[string]*Instruction{
	AdcName: AdcInst,
	AndName: AndInst,
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
	BrlName: BrlInst,
	BvcName: BvcInst,
	BvsName: BvsInst,
	ClcName: ClcInst,
	CldName: CldInst,
	CliName: CliInst,
	ClvName: ClvInst,
	CmpName: CmpInst,
	CopName: CopInst,
	CpxName: CpxInst,
	CpyName: CpyInst,
	DecName: DecInst,
	DexName: DexInst,
	DeyName: DeyInst,
	EorName: EorInst,
	IncName: IncInst,
	InxName: InxInst,
	InyName: InyInst,
	JmlName: JmlInst,
	JmpName: JmpInst,
	JslName: JslInst,
	JsrName: JsrInst,
	LdaName: LdaInst,
	LdxName: LdxInst,
	LdyName: LdyInst,
	LsrName: LsrInst,
	MvnName: MvnInst,
	MvpName: MvpInst,
	NopName: NopInst,
	OraName: OraInst,
	PeaName: PeaInst,
	PeiName: PeiInst,
	PerName: PerInst,
	PhaName: PhaInst,
	PhbName: PhbInst,
	PhdName: PhdInst,
	PhkName: PhkInst,
	PhpName: PhpInst,
	PhxName: PhxInst,
	PhyName: PhyInst,
	PlaName: PlaInst,
	PlbName: PlbInst,
	PldName: PldInst,
	PlpName: PlpInst,
	PlxName: PlxInst,
	PlyName: PlyInst,
	RepName: RepInst,
	RolName: RolInst,
	RorName: RorInst,
	RtiName: RtiInst,
	RtlName: RtlInst,
	RtsName: RtsInst,
	SbcName: SbcInst,
	SecName: SecInst,
	SedName: SedInst,
	SeiName: SeiInst,
	SepName: SepInst,
	StaName: StaInst,
	StpName: StpInst,
	StxName: StxInst,
	StyName: StyInst,
	StzName: StzInst,
	TaxName: TaxInst,
	TayName: TayInst,
	TcdName: TcdInst,
	TcsName: TcsInst,
	TdcName: TdcInst,
	TrbName: TrbInst,
	TsbName: TsbInst,
	TscName: TscInst,
	TsxName: TsxInst,
	TxaName: TxaInst,
	TxsName: TxsInst,
	TxyName: TxyInst,
	TyaName: TyaInst,
	TyxName: TyxInst,
	WaiName: WaiInst,
	WdmName: WdmInst,
	XbaName: XbaInst,
	XceName: XceInst,
}
