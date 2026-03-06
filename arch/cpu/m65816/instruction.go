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

// Adc - Add with Carry.
var Adc = &Instruction{
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

// And - AND with Accumulator.
var And = &Instruction{
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

// Asl - Arithmetic Shift Left.
var Asl = &Instruction{
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

// Bcc - Branch if Carry Clear.
var Bcc = &Instruction{
	Name:       BccName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x90, BaseSize: 2}},
	ParamFunc:  bcc,
}

// Bcs - Branch if Carry Set.
var Bcs = &Instruction{
	Name:       BcsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0xB0, BaseSize: 2}},
	ParamFunc:  bcs,
}

// Beq - Branch if Equal (Z=1).
var Beq = &Instruction{
	Name:       BeqName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0xF0, BaseSize: 2}},
	ParamFunc:  beq,
}

// Bit - Bit Test.
var Bit = &Instruction{
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

// Bmi - Branch if Minus (N=1).
var Bmi = &Instruction{
	Name:       BmiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x30, BaseSize: 2}},
	ParamFunc:  bmi,
}

// Bne - Branch if Not Equal (Z=0).
var Bne = &Instruction{
	Name:       BneName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0xD0, BaseSize: 2}},
	ParamFunc:  bne,
}

// Bpl - Branch if Positive (N=0).
var Bpl = &Instruction{
	Name:       BplName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x10, BaseSize: 2}},
	ParamFunc:  bpl,
}

// Bra - Branch Always.
var Bra = &Instruction{
	Name:       BraName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x80, BaseSize: 2}},
	ParamFunc:  bra,
}

// Brk - Software Interrupt / Break.
var Brk = &Instruction{
	Name:        BrkName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x00, BaseSize: 2}},
	NoParamFunc: brk,
}

// Brl - Branch Long (16-bit offset).
var Brl = &Instruction{
	Name:       BrlName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Opcode: 0x82, BaseSize: 3}},
	ParamFunc:  brl,
}

// Bvc - Branch if Overflow Clear.
var Bvc = &Instruction{
	Name:       BvcName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x50, BaseSize: 2}},
	ParamFunc:  bvc,
}

// Bvs - Branch if Overflow Set.
var Bvs = &Instruction{
	Name:       BvsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x70, BaseSize: 2}},
	ParamFunc:  bvs,
}

// Clc - Clear Carry Flag.
var Clc = &Instruction{
	Name:        ClcName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x18, BaseSize: 1}},
	NoParamFunc: clc,
}

// Cld - Clear Decimal Flag.
var Cld = &Instruction{
	Name:        CldName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xD8, BaseSize: 1}},
	NoParamFunc: cld,
}

// Cli - Clear Interrupt Disable.
var Cli = &Instruction{
	Name:        CliName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x58, BaseSize: 1}},
	NoParamFunc: cli,
}

// Clv - Clear Overflow Flag.
var Clv = &Instruction{
	Name:        ClvName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xB8, BaseSize: 1}},
	NoParamFunc: clv,
}

// Cmp - Compare Accumulator.
var Cmp = &Instruction{
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

// Cop - Co-Processor Enable (software interrupt).
var Cop = &Instruction{
	Name:        CopName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x02, BaseSize: 2}},
	NoParamFunc: cop,
}

// Cpx - Compare X Register.
var Cpx = &Instruction{
	Name: CpxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:  {Opcode: 0xE0, BaseSize: 2},
		DirectPageAddressing: {Opcode: 0xE4, BaseSize: 2},
		AbsoluteAddressing:   {Opcode: 0xEC, BaseSize: 3},
	},
	ParamFunc: cpx,
}

// Cpy - Compare Y Register.
var Cpy = &Instruction{
	Name: CpyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing:  {Opcode: 0xC0, BaseSize: 2},
		DirectPageAddressing: {Opcode: 0xC4, BaseSize: 2},
		AbsoluteAddressing:   {Opcode: 0xCC, BaseSize: 3},
	},
	ParamFunc: cpy,
}

// Dec - Decrement.
var Dec = &Instruction{
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

// Dex - Decrement X.
var Dex = &Instruction{
	Name:        DexName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xCA, BaseSize: 1}},
	NoParamFunc: dex,
}

// Dey - Decrement Y.
var Dey = &Instruction{
	Name:        DeyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x88, BaseSize: 1}},
	NoParamFunc: dey,
}

// Eor - Exclusive OR with Accumulator.
var Eor = &Instruction{
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

// Inc - Increment.
var Inc = &Instruction{
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

// Inx - Increment X.
var Inx = &Instruction{
	Name:        InxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xE8, BaseSize: 1}},
	NoParamFunc: inx,
}

// Iny - Increment Y.
var Iny = &Instruction{
	Name:        InyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xC8, BaseSize: 1}},
	NoParamFunc: iny,
}

// Jml - Jump Long (sets PB).
var Jml = &Instruction{
	Name: JmlName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteLongAddressing:         {Opcode: 0x5C, BaseSize: 4},
		AbsoluteIndirectLongAddressing: {Opcode: 0xDC, BaseSize: 3},
	},
	ParamFunc: jml,
}

// Jmp - Jump.
var Jmp = &Instruction{
	Name: JmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing:                 {Opcode: 0x4C, BaseSize: 3},
		AbsoluteIndirectAddressing:         {Opcode: 0x6C, BaseSize: 3},
		AbsoluteIndexedXIndirectAddressing: {Opcode: 0x7C, BaseSize: 3},
	},
	ParamFunc: jmp,
}

// Jsl - Jump to Subroutine Long.
var Jsl = &Instruction{
	Name:       JslName,
	Addressing: map[AddressingMode]OpcodeInfo{AbsoluteLongAddressing: {Opcode: 0x22, BaseSize: 4}},
	ParamFunc:  jsl,
}

// Jsr - Jump to Subroutine.
var Jsr = &Instruction{
	Name: JsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		AbsoluteAddressing:                 {Opcode: 0x20, BaseSize: 3},
		AbsoluteIndexedXIndirectAddressing: {Opcode: 0xFC, BaseSize: 3},
	},
	ParamFunc: jsr,
}

// Lda - Load Accumulator.
var Lda = &Instruction{
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

// Ldx - Load X Register.
var Ldx = &Instruction{
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

// Ldy - Load Y Register.
var Ldy = &Instruction{
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

// Lsr - Logical Shift Right.
var Lsr = &Instruction{
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

// Mvn - Move Block Next (increment).
var Mvn = &Instruction{
	Name:       MvnName,
	Addressing: map[AddressingMode]OpcodeInfo{BlockMoveAddressing: {Opcode: 0x54, BaseSize: 3}},
	ParamFunc:  mvn,
}

// Mvp - Move Block Previous (decrement).
var Mvp = &Instruction{
	Name:       MvpName,
	Addressing: map[AddressingMode]OpcodeInfo{BlockMoveAddressing: {Opcode: 0x44, BaseSize: 3}},
	ParamFunc:  mvp,
}

// Nop - No Operation.
var Nop = &Instruction{
	Name:        NopName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xEA, BaseSize: 1}},
	NoParamFunc: nop,
}

// Ora - OR with Accumulator.
var Ora = &Instruction{
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

// Pea - Push Effective Absolute Address.
var Pea = &Instruction{
	Name:        PeaName,
	Addressing:  map[AddressingMode]OpcodeInfo{AbsoluteAddressing: {Opcode: 0xF4, BaseSize: 3}},
	NoParamFunc: pea,
}

// Pei - Push Effective Indirect Address.
var Pei = &Instruction{
	Name:        PeiName,
	Addressing:  map[AddressingMode]OpcodeInfo{DirectPageIndirectAddressing: {Opcode: 0xD4, BaseSize: 2}},
	NoParamFunc: pei,
}

// Per - Push Effective Relative Address.
var Per = &Instruction{
	Name:        PerName,
	Addressing:  map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Opcode: 0x62, BaseSize: 3}},
	NoParamFunc: per,
}

// Pha - Push Accumulator.
var Pha = &Instruction{
	Name:        PhaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x48, BaseSize: 1}},
	NoParamFunc: pha,
}

// Phb - Push Data Bank Register.
var Phb = &Instruction{
	Name:        PhbName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x8B, BaseSize: 1}},
	NoParamFunc: phb,
}

// Phd - Push Direct Page Register.
var Phd = &Instruction{
	Name:        PhdName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x0B, BaseSize: 1}},
	NoParamFunc: phd,
}

// Phk - Push Program Bank Register.
var Phk = &Instruction{
	Name:        PhkName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x4B, BaseSize: 1}},
	NoParamFunc: phk,
}

// Php - Push Processor Status.
var Php = &Instruction{
	Name:        PhpName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x08, BaseSize: 1}},
	NoParamFunc: php,
}

// Phx - Push X Register.
var Phx = &Instruction{
	Name:        PhxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xDA, BaseSize: 1}},
	NoParamFunc: phx,
}

// Phy - Push Y Register.
var Phy = &Instruction{
	Name:        PhyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x5A, BaseSize: 1}},
	NoParamFunc: phy,
}

// Pla - Pull Accumulator.
var Pla = &Instruction{
	Name:        PlaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x68, BaseSize: 1}},
	NoParamFunc: pla,
}

// Plb - Pull Data Bank Register.
var Plb = &Instruction{
	Name:        PlbName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xAB, BaseSize: 1}},
	NoParamFunc: plb,
}

// Pld - Pull Direct Page Register.
var Pld = &Instruction{
	Name:        PldName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x2B, BaseSize: 1}},
	NoParamFunc: pld,
}

// Plp - Pull Processor Status.
var Plp = &Instruction{
	Name:        PlpName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x28, BaseSize: 1}},
	NoParamFunc: plp,
}

// Plx - Pull X Register.
var Plx = &Instruction{
	Name:        PlxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xFA, BaseSize: 1}},
	NoParamFunc: plx,
}

// Ply - Pull Y Register.
var Ply = &Instruction{
	Name:        PlyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x7A, BaseSize: 1}},
	NoParamFunc: ply,
}

// Rep - Reset Processor Status Bits.
var Rep = &Instruction{
	Name:       RepName,
	Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0xC2, BaseSize: 2}},
	ParamFunc:  rep,
}

// Rol - Rotate Left.
var Rol = &Instruction{
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

// Ror - Rotate Right.
var Ror = &Instruction{
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

// Rti - Return from Interrupt.
var Rti = &Instruction{
	Name:        RtiName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x40, BaseSize: 1}},
	NoParamFunc: rti,
}

// Rtl - Return from Subroutine Long.
var Rtl = &Instruction{
	Name:        RtlName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x6B, BaseSize: 1}},
	NoParamFunc: rtl,
}

// Rts - Return from Subroutine.
var Rts = &Instruction{
	Name:        RtsName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x60, BaseSize: 1}},
	NoParamFunc: rts,
}

// Sbc - Subtract with Carry.
var Sbc = &Instruction{
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

// Sec - Set Carry Flag.
var Sec = &Instruction{
	Name:        SecName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x38, BaseSize: 1}},
	NoParamFunc: sec,
}

// Sed - Set Decimal Flag.
var Sed = &Instruction{
	Name:        SedName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xF8, BaseSize: 1}},
	NoParamFunc: sed,
}

// Sei - Set Interrupt Disable.
var Sei = &Instruction{
	Name:        SeiName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x78, BaseSize: 1}},
	NoParamFunc: sei,
}

// Sep - Set Processor Status Bits.
var Sep = &Instruction{
	Name:       SepName,
	Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0xE2, BaseSize: 2}},
	ParamFunc:  sep,
}

// Sta - Store Accumulator.
var Sta = &Instruction{
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

// Stp - Stop the Processor.
var Stp = &Instruction{
	Name:        StpName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xDB, BaseSize: 1}},
	NoParamFunc: stp,
}

// Stx - Store X Register.
var Stx = &Instruction{
	Name: StxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing:         {Opcode: 0x86, BaseSize: 2},
		DirectPageIndexedYAddressing: {Opcode: 0x96, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x8E, BaseSize: 3},
	},
	ParamFunc: stx,
}

// Sty - Store Y Register.
var Sty = &Instruction{
	Name: StyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing:         {Opcode: 0x84, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0x94, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x8C, BaseSize: 3},
	},
	ParamFunc: sty,
}

// Stz - Store Zero.
var Stz = &Instruction{
	Name: StzName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing:         {Opcode: 0x64, BaseSize: 2},
		DirectPageIndexedXAddressing: {Opcode: 0x74, BaseSize: 2},
		AbsoluteAddressing:           {Opcode: 0x9C, BaseSize: 3},
		AbsoluteIndexedXAddressing:   {Opcode: 0x9E, BaseSize: 3},
	},
	ParamFunc: stz,
}

// Tax - Transfer A to X.
var Tax = &Instruction{
	Name:        TaxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xAA, BaseSize: 1}},
	NoParamFunc: tax,
}

// Tay - Transfer A to Y.
var Tay = &Instruction{
	Name:        TayName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xA8, BaseSize: 1}},
	NoParamFunc: tay,
}

// Tcd - Transfer C to Direct Page.
var Tcd = &Instruction{
	Name:        TcdName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x5B, BaseSize: 1}},
	NoParamFunc: tcd,
}

// Tcs - Transfer C to Stack Pointer.
var Tcs = &Instruction{
	Name:        TcsName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x1B, BaseSize: 1}},
	NoParamFunc: tcs,
}

// Tdc - Transfer Direct Page to C.
var Tdc = &Instruction{
	Name:        TdcName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x7B, BaseSize: 1}},
	NoParamFunc: tdc,
}

// Trb - Test and Reset Bits.
var Trb = &Instruction{
	Name: TrbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing: {Opcode: 0x14, BaseSize: 2},
		AbsoluteAddressing:   {Opcode: 0x1C, BaseSize: 3},
	},
	ParamFunc: trb,
}

// Tsb - Test and Set Bits.
var Tsb = &Instruction{
	Name: TsbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectPageAddressing: {Opcode: 0x04, BaseSize: 2},
		AbsoluteAddressing:   {Opcode: 0x0C, BaseSize: 3},
	},
	ParamFunc: tsb,
}

// Tsc - Transfer Stack Pointer to C.
var Tsc = &Instruction{
	Name:        TscName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x3B, BaseSize: 1}},
	NoParamFunc: tsc,
}

// Tsx - Transfer SP to X.
var Tsx = &Instruction{
	Name:        TsxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xBA, BaseSize: 1}},
	NoParamFunc: tsx,
}

// Txa - Transfer X to A.
var Txa = &Instruction{
	Name:        TxaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x8A, BaseSize: 1}},
	NoParamFunc: txa,
}

// Txs - Transfer X to SP.
var Txs = &Instruction{
	Name:        TxsName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x9A, BaseSize: 1}},
	NoParamFunc: txs,
}

// Txy - Transfer X to Y.
var Txy = &Instruction{
	Name:        TxyName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x9B, BaseSize: 1}},
	NoParamFunc: txy,
}

// Tya - Transfer Y to A.
var Tya = &Instruction{
	Name:        TyaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0x98, BaseSize: 1}},
	NoParamFunc: tya,
}

// Tyx - Transfer Y to X.
var Tyx = &Instruction{
	Name:        TyxName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xBB, BaseSize: 1}},
	NoParamFunc: tyx,
}

// Wai - Wait for Interrupt.
var Wai = &Instruction{
	Name:        WaiName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xCB, BaseSize: 1}},
	NoParamFunc: wai,
}

// Wdm - Reserved/WDM (2-byte NOP).
var Wdm = &Instruction{
	Name:       WdmName,
	Addressing: map[AddressingMode]OpcodeInfo{ImmediateAddressing: {Opcode: 0x42, BaseSize: 2}},
	ParamFunc:  wdm,
}

// Xba - Exchange B and A (swap accumulator bytes).
var Xba = &Instruction{
	Name:        XbaName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xEB, BaseSize: 1}},
	NoParamFunc: xba,
}

// Xce - Exchange Carry and Emulation flags.
var Xce = &Instruction{
	Name:        XceName,
	Addressing:  map[AddressingMode]OpcodeInfo{ImpliedAddressing: {Opcode: 0xFB, BaseSize: 1}},
	NoParamFunc: xce,
}

// Instructions maps instruction names to their definitions.
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
	BraName: Bra,
	BrkName: Brk,
	BrlName: Brl,
	BvcName: Bvc,
	BvsName: Bvs,
	ClcName: Clc,
	CldName: Cld,
	CliName: Cli,
	ClvName: Clv,
	CmpName: Cmp,
	CopName: Cop,
	CpxName: Cpx,
	CpyName: Cpy,
	DecName: Dec,
	DexName: Dex,
	DeyName: Dey,
	EorName: Eor,
	IncName: Inc,
	InxName: Inx,
	InyName: Iny,
	JmlName: Jml,
	JmpName: Jmp,
	JslName: Jsl,
	JsrName: Jsr,
	LdaName: Lda,
	LdxName: Ldx,
	LdyName: Ldy,
	LsrName: Lsr,
	MvnName: Mvn,
	MvpName: Mvp,
	NopName: Nop,
	OraName: Ora,
	PeaName: Pea,
	PeiName: Pei,
	PerName: Per,
	PhaName: Pha,
	PhbName: Phb,
	PhdName: Phd,
	PhkName: Phk,
	PhpName: Php,
	PhxName: Phx,
	PhyName: Phy,
	PlaName: Pla,
	PlbName: Plb,
	PldName: Pld,
	PlpName: Plp,
	PlxName: Plx,
	PlyName: Ply,
	RepName: Rep,
	RolName: Rol,
	RorName: Ror,
	RtiName: Rti,
	RtlName: Rtl,
	RtsName: Rts,
	SbcName: Sbc,
	SecName: Sec,
	SedName: Sed,
	SeiName: Sei,
	SepName: Sep,
	StaName: Sta,
	StpName: Stp,
	StxName: Stx,
	StyName: Sty,
	StzName: Stz,
	TaxName: Tax,
	TayName: Tay,
	TcdName: Tcd,
	TcsName: Tcs,
	TdcName: Tdc,
	TrbName: Trb,
	TsbName: Tsb,
	TscName: Tsc,
	TsxName: Tsx,
	TxaName: Txa,
	TxsName: Txs,
	TxyName: Txy,
	TyaName: Tya,
	TyxName: Tyx,
	WaiName: Wai,
	WdmName: Wdm,
	XbaName: Xba,
	XceName: Xce,
}
