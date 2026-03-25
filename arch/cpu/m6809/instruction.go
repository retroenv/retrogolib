package m6809

// Instruction defines a 6809 CPU instruction.
type Instruction struct {
	Name       string
	Unofficial bool

	// Addressing maps each supported addressing mode to its opcode and size info.
	Addressing map[AddressingMode]OpcodeInfo

	// Exactly one of these must be set.
	NoParamFunc func(c *CPU) error
	ParamFunc   func(c *CPU, params ...any) error
}

// OpcodeInfo contains the opcode byte(s) and instruction size.
type OpcodeInfo struct {
	Prefix byte // Prefix byte (0x00 for base page, 0x10 for page 2, 0x11 for page 3)
	Opcode byte
	Size   byte // Total size in bytes including prefix
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
	AbxName   = "abx"
	AdcaName  = "adca"
	AdcbName  = "adcb"
	AddaName  = "adda"
	AddbName  = "addb"
	AdddName  = "addd"
	AndaName  = "anda"
	AndbName  = "andb"
	AndccName = "andcc"
	AslName   = "asl"
	AsrName   = "asr"
	BitaName  = "bita"
	BitbName  = "bitb"
	BccName   = "bcc"
	BcsName   = "bcs"
	BeqName   = "beq"
	BgeName   = "bge"
	BgtName   = "bgt"
	BhiName   = "bhi"
	BleName   = "ble"
	BlsName   = "bls"
	BltName   = "blt"
	BmiName   = "bmi"
	BneName   = "bne"
	BplName   = "bpl"
	BraName   = "bra"
	BrnName   = "brn"
	BsrName   = "bsr"
	BvcName   = "bvc"
	BvsName   = "bvs"
	ClrName   = "clr"
	CmpaName  = "cmpa"
	CmpbName  = "cmpb"
	CmpdName  = "cmpd"
	CmpsName  = "cmps"
	CmpuName  = "cmpu"
	CmpxName  = "cmpx"
	CmpyName  = "cmpy"
	ComName   = "com"
	CwaiName  = "cwai"
	DaaName   = "daa"
	DecName   = "dec"
	EoraName  = "eora"
	EorbName  = "eorb"
	ExgName   = "exg"
	IncName   = "inc"
	JmpName   = "jmp"
	JsrName   = "jsr"
	LbccName  = "lbcc"
	LbcsName  = "lbcs"
	LbeqName  = "lbeq"
	LbgeName  = "lbge"
	LbgtName  = "lbgt"
	LbhiName  = "lbhi"
	LbleName  = "lble"
	LblsName  = "lbls"
	LbltName  = "lblt"
	LbmiName  = "lbmi"
	LbneName  = "lbne"
	LbplName  = "lbpl"
	LbraName  = "lbra"
	LbrnName  = "lbrn"
	LbsrName  = "lbsr"
	LbvcName  = "lbvc"
	LbvsName  = "lbvs"
	LdaName   = "lda"
	LdbName   = "ldb"
	LddName   = "ldd"
	LdsName   = "lds"
	LduName   = "ldu"
	LdxName   = "ldx"
	LdyName   = "ldy"
	LeaxName  = "leax"
	LeayName  = "leay"
	LeasName  = "leas"
	LeauName  = "leau"
	LsrName   = "lsr"
	MulName   = "mul"
	NegName   = "neg"
	NopName   = "nop"
	OraName   = "ora"
	OrbName   = "orb"
	OrccName  = "orcc"
	PshsName  = "pshs"
	PshuName  = "pshu"
	PulsName  = "puls"
	PuluName  = "pulu"
	RolName   = "rol"
	RorName   = "ror"
	RtiName   = "rti"
	RtsName   = "rts"
	SbcaName  = "sbca"
	SbcbName  = "sbcb"
	SexName   = "sex"
	StaName   = "sta"
	StbName   = "stb"
	StdName   = "std"
	StsName   = "sts"
	StuName   = "stu"
	StxName   = "stx"
	StyName   = "sty"
	SubaName  = "suba"
	SubbName  = "subb"
	SubdName  = "subd"
	SwiName   = "swi"
	Swi2Name  = "swi2"
	Swi3Name  = "swi3"
	SyncName  = "sync"
	TfrName   = "tfr"
	TstName   = "tst"
)

// -- Instruction variable definitions (sorted alphabetically) --

// AbxInst - Add B to X (unsigned).
var AbxInst = &Instruction{
	Name: AbxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3A, Size: 1},
	},
	NoParamFunc: abx,
}

// AdcaInst - Add with Carry to A.
var AdcaInst = &Instruction{
	Name: AdcaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x89, Size: 2},
		DirectAddressing:    {Opcode: 0x99, Size: 2},
		IndexedAddressing:   {Opcode: 0xA9, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB9, Size: 3},
	},
	ParamFunc: adca,
}

// AdcbInst - Add with Carry to B.
var AdcbInst = &Instruction{
	Name: AdcbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC9, Size: 2},
		DirectAddressing:    {Opcode: 0xD9, Size: 2},
		IndexedAddressing:   {Opcode: 0xE9, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF9, Size: 3},
	},
	ParamFunc: adcb,
}

// AddaInst - Add to A.
var AddaInst = &Instruction{
	Name: AddaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x8B, Size: 2},
		DirectAddressing:    {Opcode: 0x9B, Size: 2},
		IndexedAddressing:   {Opcode: 0xAB, Size: 2},
		ExtendedAddressing:  {Opcode: 0xBB, Size: 3},
	},
	ParamFunc: adda,
}

// AddbInst - Add to B.
var AddbInst = &Instruction{
	Name: AddbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xCB, Size: 2},
		DirectAddressing:    {Opcode: 0xDB, Size: 2},
		IndexedAddressing:   {Opcode: 0xEB, Size: 2},
		ExtendedAddressing:  {Opcode: 0xFB, Size: 3},
	},
	ParamFunc: addb,
}

// AdddInst - Add to D (16-bit).
var AdddInst = &Instruction{
	Name: AdddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0xC3, Size: 3},
		DirectAddressing:      {Opcode: 0xD3, Size: 2},
		IndexedAddressing:     {Opcode: 0xE3, Size: 2},
		ExtendedAddressing:    {Opcode: 0xF3, Size: 3},
	},
	ParamFunc: addd,
}

// AndaInst - AND with A.
var AndaInst = &Instruction{
	Name: AndaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x84, Size: 2},
		DirectAddressing:    {Opcode: 0x94, Size: 2},
		IndexedAddressing:   {Opcode: 0xA4, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB4, Size: 3},
	},
	ParamFunc: anda,
}

// AndbInst - AND with B.
var AndbInst = &Instruction{
	Name: AndbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC4, Size: 2},
		DirectAddressing:    {Opcode: 0xD4, Size: 2},
		IndexedAddressing:   {Opcode: 0xE4, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF4, Size: 3},
	},
	ParamFunc: andb,
}

// AndccInst - AND CC Register.
var AndccInst = &Instruction{
	Name: AndccName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x1C, Size: 2},
	},
	ParamFunc: andcc,
}

// AslInst - Arithmetic Shift Left (memory/inherent A/B).
var AslInst = &Instruction{
	Name: AslName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x08, Size: 2},
		IndexedAddressing:  {Opcode: 0x68, Size: 2},
		ExtendedAddressing: {Opcode: 0x78, Size: 3},
	},
	ParamFunc: aslMem,
}

// AslaInst - Arithmetic Shift Left A (inherent).
var AslaInst = &Instruction{
	Name: AslName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x48, Size: 1},
	},
	NoParamFunc: asla,
}

// AslbInst - Arithmetic Shift Left B (inherent).
var AslbInst = &Instruction{
	Name: AslName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x58, Size: 1},
	},
	NoParamFunc: aslb,
}

// AsrInst - Arithmetic Shift Right (memory).
var AsrInst = &Instruction{
	Name: AsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x07, Size: 2},
		IndexedAddressing:  {Opcode: 0x67, Size: 2},
		ExtendedAddressing: {Opcode: 0x77, Size: 3},
	},
	ParamFunc: asrMem,
}

// AsraInst - Arithmetic Shift Right A (inherent).
var AsraInst = &Instruction{
	Name: AsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x47, Size: 1},
	},
	NoParamFunc: asra,
}

// AsrbInst - Arithmetic Shift Right B (inherent).
var AsrbInst = &Instruction{
	Name: AsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x57, Size: 1},
	},
	NoParamFunc: asrb,
}

// BitaInst - Bit Test A.
var BitaInst = &Instruction{
	Name: BitaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x85, Size: 2},
		DirectAddressing:    {Opcode: 0x95, Size: 2},
		IndexedAddressing:   {Opcode: 0xA5, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB5, Size: 3},
	},
	ParamFunc: bita,
}

// BitbInst - Bit Test B.
var BitbInst = &Instruction{
	Name: BitbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC5, Size: 2},
		DirectAddressing:    {Opcode: 0xD5, Size: 2},
		IndexedAddressing:   {Opcode: 0xE5, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF5, Size: 3},
	},
	ParamFunc: bitb,
}

// BccInst - Branch if Carry Clear.
var BccInst = &Instruction{
	Name:       BccName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x24, Size: 2}},
	ParamFunc:  bccFn,
}

// BcsInst - Branch if Carry Set.
var BcsInst = &Instruction{
	Name:       BcsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x25, Size: 2}},
	ParamFunc:  bcsFn,
}

// BeqInst - Branch if Equal.
var BeqInst = &Instruction{
	Name:       BeqName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x27, Size: 2}},
	ParamFunc:  beqFn,
}

// BgeInst - Branch if Greater or Equal.
var BgeInst = &Instruction{
	Name:       BgeName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2C, Size: 2}},
	ParamFunc:  bgeFn,
}

// BgtInst - Branch if Greater Than.
var BgtInst = &Instruction{
	Name:       BgtName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2E, Size: 2}},
	ParamFunc:  bgtFn,
}

// BhiInst - Branch if Higher (unsigned).
var BhiInst = &Instruction{
	Name:       BhiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x22, Size: 2}},
	ParamFunc:  bhiFn,
}

// BleInst - Branch if Less or Equal.
var BleInst = &Instruction{
	Name:       BleName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2F, Size: 2}},
	ParamFunc:  bleFn,
}

// BlsInst - Branch if Lower or Same (unsigned).
var BlsInst = &Instruction{
	Name:       BlsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x23, Size: 2}},
	ParamFunc:  blsFn,
}

// BltInst - Branch if Less Than.
var BltInst = &Instruction{
	Name:       BltName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2D, Size: 2}},
	ParamFunc:  bltFn,
}

// BmiInst - Branch if Minus.
var BmiInst = &Instruction{
	Name:       BmiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2B, Size: 2}},
	ParamFunc:  bmiFn,
}

// BneInst - Branch if Not Equal.
var BneInst = &Instruction{
	Name:       BneName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x26, Size: 2}},
	ParamFunc:  bneFn,
}

// BplInst - Branch if Plus.
var BplInst = &Instruction{
	Name:       BplName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2A, Size: 2}},
	ParamFunc:  bplFn,
}

// BraInst - Branch Always.
var BraInst = &Instruction{
	Name:       BraName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x20, Size: 2}},
	ParamFunc:  braFn,
}

// BrnInst - Branch Never.
var BrnInst = &Instruction{
	Name:       BrnName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x21, Size: 2}},
	ParamFunc:  brnFn,
}

// BsrInst - Branch to Subroutine.
var BsrInst = &Instruction{
	Name:       BsrName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x8D, Size: 2}},
	ParamFunc:  bsrFn,
}

// BvcInst - Branch if Overflow Clear.
var BvcInst = &Instruction{
	Name:       BvcName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x28, Size: 2}},
	ParamFunc:  bvcFn,
}

// BvsInst - Branch if Overflow Set.
var BvsInst = &Instruction{
	Name:       BvsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x29, Size: 2}},
	ParamFunc:  bvsFn,
}

// ClrInst - Clear (memory).
var ClrInst = &Instruction{
	Name: ClrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0F, Size: 2},
		IndexedAddressing:  {Opcode: 0x6F, Size: 2},
		ExtendedAddressing: {Opcode: 0x7F, Size: 3},
	},
	ParamFunc: clrMem,
}

// ClraInst - Clear A (inherent).
var ClraInst = &Instruction{
	Name: ClrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x4F, Size: 1},
	},
	NoParamFunc: clra,
}

// ClrbInst - Clear B (inherent).
var ClrbInst = &Instruction{
	Name: ClrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5F, Size: 1},
	},
	NoParamFunc: clrb,
}

// CmpaInst - Compare A.
var CmpaInst = &Instruction{
	Name: CmpaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x81, Size: 2},
		DirectAddressing:    {Opcode: 0x91, Size: 2},
		IndexedAddressing:   {Opcode: 0xA1, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB1, Size: 3},
	},
	ParamFunc: cmpa,
}

// CmpbInst - Compare B.
var CmpbInst = &Instruction{
	Name: CmpbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC1, Size: 2},
		DirectAddressing:    {Opcode: 0xD1, Size: 2},
		IndexedAddressing:   {Opcode: 0xE1, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF1, Size: 3},
	},
	ParamFunc: cmpb,
}

// CmpdInst - Compare D (16-bit, page 2).
var CmpdInst = &Instruction{
	Name: CmpdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x10, Opcode: 0x83, Size: 4},
		DirectAddressing:      {Prefix: 0x10, Opcode: 0x93, Size: 3},
		IndexedAddressing:     {Prefix: 0x10, Opcode: 0xA3, Size: 3},
		ExtendedAddressing:    {Prefix: 0x10, Opcode: 0xB3, Size: 4},
	},
	ParamFunc: cmpd,
}

// CmpsInst - Compare S (16-bit, page 3).
var CmpsInst = &Instruction{
	Name: CmpsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x11, Opcode: 0x8C, Size: 4},
		DirectAddressing:      {Prefix: 0x11, Opcode: 0x9C, Size: 3},
		IndexedAddressing:     {Prefix: 0x11, Opcode: 0xAC, Size: 3},
		ExtendedAddressing:    {Prefix: 0x11, Opcode: 0xBC, Size: 4},
	},
	ParamFunc: cmps,
}

// CmpuInst - Compare U (16-bit, page 3).
var CmpuInst = &Instruction{
	Name: CmpuName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x11, Opcode: 0x83, Size: 4},
		DirectAddressing:      {Prefix: 0x11, Opcode: 0x93, Size: 3},
		IndexedAddressing:     {Prefix: 0x11, Opcode: 0xA3, Size: 3},
		ExtendedAddressing:    {Prefix: 0x11, Opcode: 0xB3, Size: 4},
	},
	ParamFunc: cmpu,
}

// CmpxInst - Compare X (16-bit).
var CmpxInst = &Instruction{
	Name: CmpxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0x8C, Size: 3},
		DirectAddressing:      {Opcode: 0x9C, Size: 2},
		IndexedAddressing:     {Opcode: 0xAC, Size: 2},
		ExtendedAddressing:    {Opcode: 0xBC, Size: 3},
	},
	ParamFunc: cmpx,
}

// CmpyInst - Compare Y (16-bit, page 2).
var CmpyInst = &Instruction{
	Name: CmpyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x10, Opcode: 0x8C, Size: 4},
		DirectAddressing:      {Prefix: 0x10, Opcode: 0x9C, Size: 3},
		IndexedAddressing:     {Prefix: 0x10, Opcode: 0xAC, Size: 3},
		ExtendedAddressing:    {Prefix: 0x10, Opcode: 0xBC, Size: 4},
	},
	ParamFunc: cmpy,
}

// ComInst - Complement (memory).
var ComInst = &Instruction{
	Name: ComName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x03, Size: 2},
		IndexedAddressing:  {Opcode: 0x63, Size: 2},
		ExtendedAddressing: {Opcode: 0x73, Size: 3},
	},
	ParamFunc: comMem,
}

// ComaInst - Complement A (inherent).
var ComaInst = &Instruction{
	Name: ComName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x43, Size: 1},
	},
	NoParamFunc: coma,
}

// CombInst - Complement B (inherent).
var CombInst = &Instruction{
	Name: ComName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x53, Size: 1},
	},
	NoParamFunc: comb,
}

// CwaiInst - AND CC then Wait for Interrupt.
var CwaiInst = &Instruction{
	Name: CwaiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x3C, Size: 2},
	},
	ParamFunc: cwaiFn,
}

// DaaInst - Decimal Adjust A.
var DaaInst = &Instruction{
	Name: DaaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x19, Size: 1},
	},
	NoParamFunc: daaFn,
}

// DecInst - Decrement (memory).
var DecInst = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0A, Size: 2},
		IndexedAddressing:  {Opcode: 0x6A, Size: 2},
		ExtendedAddressing: {Opcode: 0x7A, Size: 3},
	},
	ParamFunc: decMem,
}

// DecaInst - Decrement A (inherent).
var DecaInst = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x4A, Size: 1},
	},
	NoParamFunc: deca,
}

// DecbInst - Decrement B (inherent).
var DecbInst = &Instruction{
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5A, Size: 1},
	},
	NoParamFunc: decb,
}

// EoraInst - Exclusive OR with A.
var EoraInst = &Instruction{
	Name: EoraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x88, Size: 2},
		DirectAddressing:    {Opcode: 0x98, Size: 2},
		IndexedAddressing:   {Opcode: 0xA8, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB8, Size: 3},
	},
	ParamFunc: eora,
}

// EorbInst - Exclusive OR with B.
var EorbInst = &Instruction{
	Name: EorbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC8, Size: 2},
		DirectAddressing:    {Opcode: 0xD8, Size: 2},
		IndexedAddressing:   {Opcode: 0xE8, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF8, Size: 3},
	},
	ParamFunc: eorb,
}

// ExgInst - Exchange Registers.
var ExgInst = &Instruction{
	Name: ExgName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x1E, Size: 2},
	},
	ParamFunc: exgFn,
}

// IncInst - Increment (memory).
var IncInst = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0C, Size: 2},
		IndexedAddressing:  {Opcode: 0x6C, Size: 2},
		ExtendedAddressing: {Opcode: 0x7C, Size: 3},
	},
	ParamFunc: incMem,
}

// IncaInst - Increment A (inherent).
var IncaInst = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x4C, Size: 1},
	},
	NoParamFunc: inca,
}

// IncbInst - Increment B (inherent).
var IncbInst = &Instruction{
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5C, Size: 1},
	},
	NoParamFunc: incb,
}

// JmpInst - Jump.
var JmpInst = &Instruction{
	Name: JmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0E, Size: 2},
		IndexedAddressing:  {Opcode: 0x6E, Size: 2},
		ExtendedAddressing: {Opcode: 0x7E, Size: 3},
	},
	ParamFunc: jmpFn,
}

// JsrInst - Jump to Subroutine.
var JsrInst = &Instruction{
	Name: JsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x9D, Size: 2},
		IndexedAddressing:  {Opcode: 0xAD, Size: 2},
		ExtendedAddressing: {Opcode: 0xBD, Size: 3},
	},
	ParamFunc: jsrFn,
}

// LbccInst - Long Branch if Carry Clear (page 2).
var LbccInst = &Instruction{
	Name:       LbccName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x24, Size: 4}},
	ParamFunc:  lbccFn,
}

// LbcsInst - Long Branch if Carry Set (page 2).
var LbcsInst = &Instruction{
	Name:       LbcsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x25, Size: 4}},
	ParamFunc:  lbcsFn,
}

// LbeqInst - Long Branch if Equal (page 2).
var LbeqInst = &Instruction{
	Name:       LbeqName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x27, Size: 4}},
	ParamFunc:  lbeqFn,
}

// LbgeInst - Long Branch if Greater or Equal (page 2).
var LbgeInst = &Instruction{
	Name:       LbgeName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2C, Size: 4}},
	ParamFunc:  lbgeFn,
}

// LbgtInst - Long Branch if Greater Than (page 2).
var LbgtInst = &Instruction{
	Name:       LbgtName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2E, Size: 4}},
	ParamFunc:  lbgtFn,
}

// LbhiInst - Long Branch if Higher (page 2).
var LbhiInst = &Instruction{
	Name:       LbhiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x22, Size: 4}},
	ParamFunc:  lbhiFn,
}

// LbleInst - Long Branch if Less or Equal (page 2).
var LbleInst = &Instruction{
	Name:       LbleName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2F, Size: 4}},
	ParamFunc:  lbleFn,
}

// LblsInst - Long Branch if Lower or Same (page 2).
var LblsInst = &Instruction{
	Name:       LblsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x23, Size: 4}},
	ParamFunc:  lblsFn,
}

// LbltInst - Long Branch if Less Than (page 2).
var LbltInst = &Instruction{
	Name:       LbltName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2D, Size: 4}},
	ParamFunc:  lbltFn,
}

// LbmiInst - Long Branch if Minus (page 2).
var LbmiInst = &Instruction{
	Name:       LbmiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2B, Size: 4}},
	ParamFunc:  lbmiFn,
}

// LbneInst - Long Branch if Not Equal (page 2).
var LbneInst = &Instruction{
	Name:       LbneName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x26, Size: 4}},
	ParamFunc:  lbneFn,
}

// LbplInst - Long Branch if Plus (page 2).
var LbplInst = &Instruction{
	Name:       LbplName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2A, Size: 4}},
	ParamFunc:  lbplFn,
}

// LbraInst - Long Branch Always.
var LbraInst = &Instruction{
	Name:       LbraName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Opcode: 0x16, Size: 3}},
	ParamFunc:  lbraFn,
}

// LbrnInst - Long Branch Never (page 2).
var LbrnInst = &Instruction{
	Name:       LbrnName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x21, Size: 4}},
	ParamFunc:  lbrnFn,
}

// LbsrInst - Long Branch to Subroutine.
var LbsrInst = &Instruction{
	Name:       LbsrName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Opcode: 0x17, Size: 3}},
	ParamFunc:  lbsrFn,
}

// LbvcInst - Long Branch if Overflow Clear (page 2).
var LbvcInst = &Instruction{
	Name:       LbvcName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x28, Size: 4}},
	ParamFunc:  lbvcFn,
}

// LbvsInst - Long Branch if Overflow Set (page 2).
var LbvsInst = &Instruction{
	Name:       LbvsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x29, Size: 4}},
	ParamFunc:  lbvsFn,
}

// LdaInst - Load A.
var LdaInst = &Instruction{
	Name: LdaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x86, Size: 2},
		DirectAddressing:    {Opcode: 0x96, Size: 2},
		IndexedAddressing:   {Opcode: 0xA6, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB6, Size: 3},
	},
	ParamFunc: lda,
}

// LdbInst - Load B.
var LdbInst = &Instruction{
	Name: LdbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC6, Size: 2},
		DirectAddressing:    {Opcode: 0xD6, Size: 2},
		IndexedAddressing:   {Opcode: 0xE6, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF6, Size: 3},
	},
	ParamFunc: ldb,
}

// LddInst - Load D (16-bit).
var LddInst = &Instruction{
	Name: LddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0xCC, Size: 3},
		DirectAddressing:      {Opcode: 0xDC, Size: 2},
		IndexedAddressing:     {Opcode: 0xEC, Size: 2},
		ExtendedAddressing:    {Opcode: 0xFC, Size: 3},
	},
	ParamFunc: ldd,
}

// LdsInst - Load S (16-bit, page 2).
var LdsInst = &Instruction{
	Name: LdsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x10, Opcode: 0xCE, Size: 4},
		DirectAddressing:      {Prefix: 0x10, Opcode: 0xDE, Size: 3},
		IndexedAddressing:     {Prefix: 0x10, Opcode: 0xEE, Size: 3},
		ExtendedAddressing:    {Prefix: 0x10, Opcode: 0xFE, Size: 4},
	},
	ParamFunc: lds,
}

// LduInst - Load U (16-bit).
var LduInst = &Instruction{
	Name: LduName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0xCE, Size: 3},
		DirectAddressing:      {Opcode: 0xDE, Size: 2},
		IndexedAddressing:     {Opcode: 0xEE, Size: 2},
		ExtendedAddressing:    {Opcode: 0xFE, Size: 3},
	},
	ParamFunc: ldu,
}

// LdxInst - Load X (16-bit).
var LdxInst = &Instruction{
	Name: LdxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0x8E, Size: 3},
		DirectAddressing:      {Opcode: 0x9E, Size: 2},
		IndexedAddressing:     {Opcode: 0xAE, Size: 2},
		ExtendedAddressing:    {Opcode: 0xBE, Size: 3},
	},
	ParamFunc: ldx,
}

// LdyInst - Load Y (16-bit, page 2).
var LdyInst = &Instruction{
	Name: LdyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x10, Opcode: 0x8E, Size: 4},
		DirectAddressing:      {Prefix: 0x10, Opcode: 0x9E, Size: 3},
		IndexedAddressing:     {Prefix: 0x10, Opcode: 0xAE, Size: 3},
		ExtendedAddressing:    {Prefix: 0x10, Opcode: 0xBE, Size: 4},
	},
	ParamFunc: ldy,
}

// LeaxInst - Load Effective Address into X.
var LeaxInst = &Instruction{
	Name: LeaxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		IndexedAddressing: {Opcode: 0x30, Size: 2},
	},
	ParamFunc: leax,
}

// LeayInst - Load Effective Address into Y.
var LeayInst = &Instruction{
	Name: LeayName,
	Addressing: map[AddressingMode]OpcodeInfo{
		IndexedAddressing: {Opcode: 0x31, Size: 2},
	},
	ParamFunc: leay,
}

// LeasInst - Load Effective Address into S.
var LeasInst = &Instruction{
	Name: LeasName,
	Addressing: map[AddressingMode]OpcodeInfo{
		IndexedAddressing: {Opcode: 0x32, Size: 2},
	},
	ParamFunc: leas,
}

// LeauInst - Load Effective Address into U.
var LeauInst = &Instruction{
	Name: LeauName,
	Addressing: map[AddressingMode]OpcodeInfo{
		IndexedAddressing: {Opcode: 0x33, Size: 2},
	},
	ParamFunc: leau,
}

// LsrInst - Logical Shift Right (memory).
var LsrInst = &Instruction{
	Name: LsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x04, Size: 2},
		IndexedAddressing:  {Opcode: 0x64, Size: 2},
		ExtendedAddressing: {Opcode: 0x74, Size: 3},
	},
	ParamFunc: lsrMem,
}

// LsraInst - Logical Shift Right A (inherent).
var LsraInst = &Instruction{
	Name: LsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x44, Size: 1},
	},
	NoParamFunc: lsra,
}

// LsrbInst - Logical Shift Right B (inherent).
var LsrbInst = &Instruction{
	Name: LsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x54, Size: 1},
	},
	NoParamFunc: lsrb,
}

// MulInst - Multiply (unsigned A*B -> D).
var MulInst = &Instruction{
	Name: MulName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3D, Size: 1},
	},
	NoParamFunc: mulFn,
}

// NegInst - Negate (memory).
var NegInst = &Instruction{
	Name: NegName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x00, Size: 2},
		IndexedAddressing:  {Opcode: 0x60, Size: 2},
		ExtendedAddressing: {Opcode: 0x70, Size: 3},
	},
	ParamFunc: negMem,
}

// NegaInst - Negate A (inherent).
var NegaInst = &Instruction{
	Name: NegName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x40, Size: 1},
	},
	NoParamFunc: nega,
}

// NegbInst - Negate B (inherent).
var NegbInst = &Instruction{
	Name: NegName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x50, Size: 1},
	},
	NoParamFunc: negb,
}

// NopInst - No Operation.
var NopInst = &Instruction{
	Name: NopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x12, Size: 1},
	},
	NoParamFunc: nop,
}

// OraInst - OR with A.
var OraInst = &Instruction{
	Name: OraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x8A, Size: 2},
		DirectAddressing:    {Opcode: 0x9A, Size: 2},
		IndexedAddressing:   {Opcode: 0xAA, Size: 2},
		ExtendedAddressing:  {Opcode: 0xBA, Size: 3},
	},
	ParamFunc: ora,
}

// OrbInst - OR with B.
var OrbInst = &Instruction{
	Name: OrbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xCA, Size: 2},
		DirectAddressing:    {Opcode: 0xDA, Size: 2},
		IndexedAddressing:   {Opcode: 0xEA, Size: 2},
		ExtendedAddressing:  {Opcode: 0xFA, Size: 3},
	},
	ParamFunc: orb,
}

// OrccInst - OR CC Register.
var OrccInst = &Instruction{
	Name: OrccName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x1A, Size: 2},
	},
	ParamFunc: orcc,
}

// PshsInst - Push registers onto S stack.
var PshsInst = &Instruction{
	Name: PshsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		StackAddressing: {Opcode: 0x34, Size: 2},
	},
	ParamFunc: pshsFn,
}

// PshuInst - Push registers onto U stack.
var PshuInst = &Instruction{
	Name: PshuName,
	Addressing: map[AddressingMode]OpcodeInfo{
		StackAddressing: {Opcode: 0x36, Size: 2},
	},
	ParamFunc: pshuFn,
}

// PulsInst - Pull registers from S stack.
var PulsInst = &Instruction{
	Name: PulsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		StackAddressing: {Opcode: 0x35, Size: 2},
	},
	ParamFunc: pulsFn,
}

// PuluInst - Pull registers from U stack.
var PuluInst = &Instruction{
	Name: PuluName,
	Addressing: map[AddressingMode]OpcodeInfo{
		StackAddressing: {Opcode: 0x37, Size: 2},
	},
	ParamFunc: puluFn,
}

// RolInst - Rotate Left (memory).
var RolInst = &Instruction{
	Name: RolName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x09, Size: 2},
		IndexedAddressing:  {Opcode: 0x69, Size: 2},
		ExtendedAddressing: {Opcode: 0x79, Size: 3},
	},
	ParamFunc: rolMem,
}

// RolaInst - Rotate Left A (inherent).
var RolaInst = &Instruction{
	Name: RolName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x49, Size: 1},
	},
	NoParamFunc: rola,
}

// RolbInst - Rotate Left B (inherent).
var RolbInst = &Instruction{
	Name: RolName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x59, Size: 1},
	},
	NoParamFunc: rolb,
}

// RorInst - Rotate Right (memory).
var RorInst = &Instruction{
	Name: RorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x06, Size: 2},
		IndexedAddressing:  {Opcode: 0x66, Size: 2},
		ExtendedAddressing: {Opcode: 0x76, Size: 3},
	},
	ParamFunc: rorMem,
}

// RoraInst - Rotate Right A (inherent).
var RoraInst = &Instruction{
	Name: RorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x46, Size: 1},
	},
	NoParamFunc: rora,
}

// RorbInst - Rotate Right B (inherent).
var RorbInst = &Instruction{
	Name: RorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x56, Size: 1},
	},
	NoParamFunc: rorb,
}

// RtiInst - Return from Interrupt.
var RtiInst = &Instruction{
	Name: RtiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3B, Size: 1},
	},
	NoParamFunc: rtiFn,
}

// RtsInst - Return from Subroutine.
var RtsInst = &Instruction{
	Name: RtsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x39, Size: 1},
	},
	NoParamFunc: rtsFn,
}

// SbcaInst - Subtract with Carry from A.
var SbcaInst = &Instruction{
	Name: SbcaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x82, Size: 2},
		DirectAddressing:    {Opcode: 0x92, Size: 2},
		IndexedAddressing:   {Opcode: 0xA2, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB2, Size: 3},
	},
	ParamFunc: sbca,
}

// SbcbInst - Subtract with Carry from B.
var SbcbInst = &Instruction{
	Name: SbcbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC2, Size: 2},
		DirectAddressing:    {Opcode: 0xD2, Size: 2},
		IndexedAddressing:   {Opcode: 0xE2, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF2, Size: 3},
	},
	ParamFunc: sbcb,
}

// SexInst - Sign Extend B into A.
var SexInst = &Instruction{
	Name: SexName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x1D, Size: 1},
	},
	NoParamFunc: sexFn,
}

// StaInst - Store A.
var StaInst = &Instruction{
	Name: StaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x97, Size: 2},
		IndexedAddressing:  {Opcode: 0xA7, Size: 2},
		ExtendedAddressing: {Opcode: 0xB7, Size: 3},
	},
	ParamFunc: sta,
}

// StbInst - Store B.
var StbInst = &Instruction{
	Name: StbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0xD7, Size: 2},
		IndexedAddressing:  {Opcode: 0xE7, Size: 2},
		ExtendedAddressing: {Opcode: 0xF7, Size: 3},
	},
	ParamFunc: stb,
}

// StdInst - Store D (16-bit).
var StdInst = &Instruction{
	Name: StdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0xDD, Size: 2},
		IndexedAddressing:  {Opcode: 0xED, Size: 2},
		ExtendedAddressing: {Opcode: 0xFD, Size: 3},
	},
	ParamFunc: std,
}

// StsInst - Store S (16-bit, page 2).
var StsInst = &Instruction{
	Name: StsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Prefix: 0x10, Opcode: 0xDF, Size: 3},
		IndexedAddressing:  {Prefix: 0x10, Opcode: 0xEF, Size: 3},
		ExtendedAddressing: {Prefix: 0x10, Opcode: 0xFF, Size: 4},
	},
	ParamFunc: sts,
}

// StuInst - Store U (16-bit).
var StuInst = &Instruction{
	Name: StuName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0xDF, Size: 2},
		IndexedAddressing:  {Opcode: 0xEF, Size: 2},
		ExtendedAddressing: {Opcode: 0xFF, Size: 3},
	},
	ParamFunc: stu,
}

// StxInst - Store X (16-bit).
var StxInst = &Instruction{
	Name: StxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x9F, Size: 2},
		IndexedAddressing:  {Opcode: 0xAF, Size: 2},
		ExtendedAddressing: {Opcode: 0xBF, Size: 3},
	},
	ParamFunc: stx,
}

// StyInst - Store Y (16-bit, page 2).
var StyInst = &Instruction{
	Name: StyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Prefix: 0x10, Opcode: 0x9F, Size: 3},
		IndexedAddressing:  {Prefix: 0x10, Opcode: 0xAF, Size: 3},
		ExtendedAddressing: {Prefix: 0x10, Opcode: 0xBF, Size: 4},
	},
	ParamFunc: sty,
}

// SubaInst - Subtract from A.
var SubaInst = &Instruction{
	Name: SubaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x80, Size: 2},
		DirectAddressing:    {Opcode: 0x90, Size: 2},
		IndexedAddressing:   {Opcode: 0xA0, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB0, Size: 3},
	},
	ParamFunc: suba,
}

// SubbInst - Subtract from B.
var SubbInst = &Instruction{
	Name: SubbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC0, Size: 2},
		DirectAddressing:    {Opcode: 0xD0, Size: 2},
		IndexedAddressing:   {Opcode: 0xE0, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF0, Size: 3},
	},
	ParamFunc: subb,
}

// SubdInst - Subtract from D (16-bit).
var SubdInst = &Instruction{
	Name: SubdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0x83, Size: 3},
		DirectAddressing:      {Opcode: 0x93, Size: 2},
		IndexedAddressing:     {Opcode: 0xA3, Size: 2},
		ExtendedAddressing:    {Opcode: 0xB3, Size: 3},
	},
	ParamFunc: subd,
}

// SwiInst - Software Interrupt 1.
var SwiInst = &Instruction{
	Name: SwiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3F, Size: 1},
	},
	NoParamFunc: swiFn,
}

// Swi2Inst - Software Interrupt 2 (page 2).
var Swi2Inst = &Instruction{
	Name: Swi2Name,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x10, Opcode: 0x3F, Size: 2},
	},
	NoParamFunc: swi2Fn,
}

// Swi3Inst - Software Interrupt 3 (page 3).
var Swi3Inst = &Instruction{
	Name: Swi3Name,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x11, Opcode: 0x3F, Size: 2},
	},
	NoParamFunc: swi3Fn,
}

// SyncInst - Synchronize with Interrupt.
var SyncInst = &Instruction{
	Name: SyncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x13, Size: 1},
	},
	NoParamFunc: syncFn,
}

// TfrInst - Transfer Register to Register.
var TfrInst = &Instruction{
	Name: TfrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x1F, Size: 2},
	},
	ParamFunc: tfrFn,
}

// TstInst - Test (memory).
var TstInst = &Instruction{
	Name: TstName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0D, Size: 2},
		IndexedAddressing:  {Opcode: 0x6D, Size: 2},
		ExtendedAddressing: {Opcode: 0x7D, Size: 3},
	},
	ParamFunc: tstMem,
}

// TstaInst - Test A (inherent).
var TstaInst = &Instruction{
	Name: TstName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x4D, Size: 1},
	},
	NoParamFunc: tsta,
}

// TstbInst - Test B (inherent).
var TstbInst = &Instruction{
	Name: TstName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5D, Size: 1},
	},
	NoParamFunc: tstb,
}

// Instructions maps instruction names to their definitions.
var Instructions = map[string]*Instruction{
	AbxName:   AbxInst,
	AdcaName:  AdcaInst,
	AdcbName:  AdcbInst,
	AddaName:  AddaInst,
	AddbName:  AddbInst,
	AdddName:  AdddInst,
	AndaName:  AndaInst,
	AndbName:  AndbInst,
	AndccName: AndccInst,
	AslName:   AslInst,
	AsrName:   AsrInst,
	BitaName:  BitaInst,
	BitbName:  BitbInst,
	BccName:   BccInst,
	BcsName:   BcsInst,
	BeqName:   BeqInst,
	BgeName:   BgeInst,
	BgtName:   BgtInst,
	BhiName:   BhiInst,
	BleName:   BleInst,
	BlsName:   BlsInst,
	BltName:   BltInst,
	BmiName:   BmiInst,
	BneName:   BneInst,
	BplName:   BplInst,
	BraName:   BraInst,
	BrnName:   BrnInst,
	BsrName:   BsrInst,
	BvcName:   BvcInst,
	BvsName:   BvsInst,
	ClrName:   ClrInst,
	CmpaName:  CmpaInst,
	CmpbName:  CmpbInst,
	CmpdName:  CmpdInst,
	CmpsName:  CmpsInst,
	CmpuName:  CmpuInst,
	CmpxName:  CmpxInst,
	CmpyName:  CmpyInst,
	ComName:   ComInst,
	CwaiName:  CwaiInst,
	DaaName:   DaaInst,
	DecName:   DecInst,
	EoraName:  EoraInst,
	EorbName:  EorbInst,
	ExgName:   ExgInst,
	IncName:   IncInst,
	JmpName:   JmpInst,
	JsrName:   JsrInst,
	LbccName:  LbccInst,
	LbcsName:  LbcsInst,
	LbeqName:  LbeqInst,
	LbgeName:  LbgeInst,
	LbgtName:  LbgtInst,
	LbhiName:  LbhiInst,
	LbleName:  LbleInst,
	LblsName:  LblsInst,
	LbltName:  LbltInst,
	LbmiName:  LbmiInst,
	LbneName:  LbneInst,
	LbplName:  LbplInst,
	LbraName:  LbraInst,
	LbrnName:  LbrnInst,
	LbsrName:  LbsrInst,
	LbvcName:  LbvcInst,
	LbvsName:  LbvsInst,
	LdaName:   LdaInst,
	LdbName:   LdbInst,
	LddName:   LddInst,
	LdsName:   LdsInst,
	LduName:   LduInst,
	LdxName:   LdxInst,
	LdyName:   LdyInst,
	LeaxName:  LeaxInst,
	LeayName:  LeayInst,
	LeasName:  LeasInst,
	LeauName:  LeauInst,
	LsrName:   LsrInst,
	MulName:   MulInst,
	NegName:   NegInst,
	NopName:   NopInst,
	OraName:   OraInst,
	OrbName:   OrbInst,
	OrccName:  OrccInst,
	PshsName:  PshsInst,
	PshuName:  PshuInst,
	PulsName:  PulsInst,
	PuluName:  PuluInst,
	RolName:   RolInst,
	RorName:   RorInst,
	RtiName:   RtiInst,
	RtsName:   RtsInst,
	SbcaName:  SbcaInst,
	SbcbName:  SbcbInst,
	SexName:   SexInst,
	StaName:   StaInst,
	StbName:   StbInst,
	StdName:   StdInst,
	StsName:   StsInst,
	StuName:   StuInst,
	StxName:   StxInst,
	StyName:   StyInst,
	SubaName:  SubaInst,
	SubbName:  SubbInst,
	SubdName:  SubdInst,
	SwiName:   SwiInst,
	Swi2Name:  Swi2Inst,
	Swi3Name:  Swi3Inst,
	SyncName:  SyncInst,
	TfrName:   TfrInst,
	TstName:   TstInst,
}
