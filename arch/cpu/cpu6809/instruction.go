package cpu6809

// Instruction defines a 6809 CPU instruction.
type Instruction struct {
	ID   OpcodeID
	Name string

	// Addressing maps each supported addressing mode to its opcode and size info.
	Addressing map[AddressingMode]OpcodeInfo

	// Exactly one of these must be set.
	noParamFunc func(c *CPU) error
	paramFunc   func(c *CPU, param any) error
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

// AbxInst - Add B to X (unsigned).
var AbxInst = &Instruction{
	ID:   Abx,
	Name: AbxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3A, Size: 1},
	},
	noParamFunc: abx,
}

// AdcaInst - Add with Carry to A.
var AdcaInst = &Instruction{
	ID:   Adca,
	Name: AdcaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x89, Size: 2},
		DirectAddressing:    {Opcode: 0x99, Size: 2},
		IndexedAddressing:   {Opcode: 0xA9, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB9, Size: 3},
	},
	paramFunc: adca,
}

// AdcbInst - Add with Carry to B.
var AdcbInst = &Instruction{
	ID:   Adcb,
	Name: AdcbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC9, Size: 2},
		DirectAddressing:    {Opcode: 0xD9, Size: 2},
		IndexedAddressing:   {Opcode: 0xE9, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF9, Size: 3},
	},
	paramFunc: adcb,
}

// AddaInst - Add to A.
var AddaInst = &Instruction{
	ID:   Adda,
	Name: AddaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x8B, Size: 2},
		DirectAddressing:    {Opcode: 0x9B, Size: 2},
		IndexedAddressing:   {Opcode: 0xAB, Size: 2},
		ExtendedAddressing:  {Opcode: 0xBB, Size: 3},
	},
	paramFunc: adda,
}

// AddbInst - Add to B.
var AddbInst = &Instruction{
	ID:   Addb,
	Name: AddbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xCB, Size: 2},
		DirectAddressing:    {Opcode: 0xDB, Size: 2},
		IndexedAddressing:   {Opcode: 0xEB, Size: 2},
		ExtendedAddressing:  {Opcode: 0xFB, Size: 3},
	},
	paramFunc: addb,
}

// AdddInst - Add to D (16-bit).
var AdddInst = &Instruction{
	ID:   Addd,
	Name: AdddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0xC3, Size: 3},
		DirectAddressing:      {Opcode: 0xD3, Size: 2},
		IndexedAddressing:     {Opcode: 0xE3, Size: 2},
		ExtendedAddressing:    {Opcode: 0xF3, Size: 3},
	},
	paramFunc: addd,
}

// AndaInst - AND with A.
var AndaInst = &Instruction{
	ID:   Anda,
	Name: AndaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x84, Size: 2},
		DirectAddressing:    {Opcode: 0x94, Size: 2},
		IndexedAddressing:   {Opcode: 0xA4, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB4, Size: 3},
	},
	paramFunc: anda,
}

// AndbInst - AND with B.
var AndbInst = &Instruction{
	ID:   Andb,
	Name: AndbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC4, Size: 2},
		DirectAddressing:    {Opcode: 0xD4, Size: 2},
		IndexedAddressing:   {Opcode: 0xE4, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF4, Size: 3},
	},
	paramFunc: andb,
}

// AndccInst - AND CC Register.
var AndccInst = &Instruction{
	ID:   Andcc,
	Name: AndccName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x1C, Size: 2},
	},
	paramFunc: andcc,
}

// AslInst - Arithmetic Shift Left (memory/inherent A/B).
var AslInst = &Instruction{
	ID:   Asl,
	Name: AslName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x08, Size: 2},
		IndexedAddressing:  {Opcode: 0x68, Size: 2},
		ExtendedAddressing: {Opcode: 0x78, Size: 3},
	},
	paramFunc: aslMem,
}

// AslaInst - Arithmetic Shift Left A (inherent).
var AslaInst = &Instruction{
	ID:   Asl,
	Name: AslName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x48, Size: 1},
	},
	noParamFunc: asla,
}

// AslbInst - Arithmetic Shift Left B (inherent).
var AslbInst = &Instruction{
	ID:   Asl,
	Name: AslName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x58, Size: 1},
	},
	noParamFunc: aslb,
}

// AsrInst - Arithmetic Shift Right (memory).
var AsrInst = &Instruction{
	ID:   Asr,
	Name: AsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x07, Size: 2},
		IndexedAddressing:  {Opcode: 0x67, Size: 2},
		ExtendedAddressing: {Opcode: 0x77, Size: 3},
	},
	paramFunc: asrMem,
}

// AsraInst - Arithmetic Shift Right A (inherent).
var AsraInst = &Instruction{
	ID:   Asr,
	Name: AsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x47, Size: 1},
	},
	noParamFunc: asra,
}

// AsrbInst - Arithmetic Shift Right B (inherent).
var AsrbInst = &Instruction{
	ID:   Asr,
	Name: AsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x57, Size: 1},
	},
	noParamFunc: asrb,
}

// BitaInst - Bit Test A.
var BitaInst = &Instruction{
	ID:   Bita,
	Name: BitaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x85, Size: 2},
		DirectAddressing:    {Opcode: 0x95, Size: 2},
		IndexedAddressing:   {Opcode: 0xA5, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB5, Size: 3},
	},
	paramFunc: bita,
}

// BitbInst - Bit Test B.
var BitbInst = &Instruction{
	ID:   Bitb,
	Name: BitbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC5, Size: 2},
		DirectAddressing:    {Opcode: 0xD5, Size: 2},
		IndexedAddressing:   {Opcode: 0xE5, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF5, Size: 3},
	},
	paramFunc: bitb,
}

// BccInst - Branch if Carry Clear.
var BccInst = &Instruction{
	ID:         Bcc,
	Name:       BccName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x24, Size: 2}},
	paramFunc:  bccFn,
}

// BcsInst - Branch if Carry Set.
var BcsInst = &Instruction{
	ID:         Bcs,
	Name:       BcsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x25, Size: 2}},
	paramFunc:  bcsFn,
}

// BeqInst - Branch if Equal.
var BeqInst = &Instruction{
	ID:         Beq,
	Name:       BeqName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x27, Size: 2}},
	paramFunc:  beqFn,
}

// BgeInst - Branch if Greater or Equal.
var BgeInst = &Instruction{
	ID:         Bge,
	Name:       BgeName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2C, Size: 2}},
	paramFunc:  bgeFn,
}

// BgtInst - Branch if Greater Than.
var BgtInst = &Instruction{
	ID:         Bgt,
	Name:       BgtName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2E, Size: 2}},
	paramFunc:  bgtFn,
}

// BhiInst - Branch if Higher (unsigned).
var BhiInst = &Instruction{
	ID:         Bhi,
	Name:       BhiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x22, Size: 2}},
	paramFunc:  bhiFn,
}

// BleInst - Branch if Less or Equal.
var BleInst = &Instruction{
	ID:         Ble,
	Name:       BleName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2F, Size: 2}},
	paramFunc:  bleFn,
}

// BlsInst - Branch if Lower or Same (unsigned).
var BlsInst = &Instruction{
	ID:         Bls,
	Name:       BlsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x23, Size: 2}},
	paramFunc:  blsFn,
}

// BltInst - Branch if Less Than.
var BltInst = &Instruction{
	ID:         Blt,
	Name:       BltName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2D, Size: 2}},
	paramFunc:  bltFn,
}

// BmiInst - Branch if Minus.
var BmiInst = &Instruction{
	ID:         Bmi,
	Name:       BmiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2B, Size: 2}},
	paramFunc:  bmiFn,
}

// BneInst - Branch if Not Equal.
var BneInst = &Instruction{
	ID:         Bne,
	Name:       BneName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x26, Size: 2}},
	paramFunc:  bneFn,
}

// BplInst - Branch if Plus.
var BplInst = &Instruction{
	ID:         Bpl,
	Name:       BplName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x2A, Size: 2}},
	paramFunc:  bplFn,
}

// BraInst - Branch Always.
var BraInst = &Instruction{
	ID:         Bra,
	Name:       BraName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x20, Size: 2}},
	paramFunc:  braFn,
}

// BrnInst - Branch Never.
var BrnInst = &Instruction{
	ID:         Brn,
	Name:       BrnName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x21, Size: 2}},
	paramFunc:  brnFn,
}

// BsrInst - Branch to Subroutine.
var BsrInst = &Instruction{
	ID:         Bsr,
	Name:       BsrName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x8D, Size: 2}},
	paramFunc:  bsrFn,
}

// BvcInst - Branch if Overflow Clear.
var BvcInst = &Instruction{
	ID:         Bvc,
	Name:       BvcName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x28, Size: 2}},
	paramFunc:  bvcFn,
}

// BvsInst - Branch if Overflow Set.
var BvsInst = &Instruction{
	ID:         Bvs,
	Name:       BvsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeAddressing: {Opcode: 0x29, Size: 2}},
	paramFunc:  bvsFn,
}

// ClrInst - Clear (memory).
var ClrInst = &Instruction{
	ID:   Clr,
	Name: ClrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0F, Size: 2},
		IndexedAddressing:  {Opcode: 0x6F, Size: 2},
		ExtendedAddressing: {Opcode: 0x7F, Size: 3},
	},
	paramFunc: clrMem,
}

// ClraInst - Clear A (inherent).
var ClraInst = &Instruction{
	ID:   Clr,
	Name: ClrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x4F, Size: 1},
	},
	noParamFunc: clra,
}

// ClrbInst - Clear B (inherent).
var ClrbInst = &Instruction{
	ID:   Clr,
	Name: ClrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5F, Size: 1},
	},
	noParamFunc: clrb,
}

// CmpaInst - Compare A.
var CmpaInst = &Instruction{
	ID:   Cmpa,
	Name: CmpaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x81, Size: 2},
		DirectAddressing:    {Opcode: 0x91, Size: 2},
		IndexedAddressing:   {Opcode: 0xA1, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB1, Size: 3},
	},
	paramFunc: cmpa,
}

// CmpbInst - Compare B.
var CmpbInst = &Instruction{
	ID:   Cmpb,
	Name: CmpbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC1, Size: 2},
		DirectAddressing:    {Opcode: 0xD1, Size: 2},
		IndexedAddressing:   {Opcode: 0xE1, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF1, Size: 3},
	},
	paramFunc: cmpb,
}

// CmpdInst - Compare D (16-bit, page 2).
var CmpdInst = &Instruction{
	ID:   Cmpd,
	Name: CmpdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x10, Opcode: 0x83, Size: 4},
		DirectAddressing:      {Prefix: 0x10, Opcode: 0x93, Size: 3},
		IndexedAddressing:     {Prefix: 0x10, Opcode: 0xA3, Size: 3},
		ExtendedAddressing:    {Prefix: 0x10, Opcode: 0xB3, Size: 4},
	},
	paramFunc: cmpd,
}

// CmpsInst - Compare S (16-bit, page 3).
var CmpsInst = &Instruction{
	ID:   Cmps,
	Name: CmpsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x11, Opcode: 0x8C, Size: 4},
		DirectAddressing:      {Prefix: 0x11, Opcode: 0x9C, Size: 3},
		IndexedAddressing:     {Prefix: 0x11, Opcode: 0xAC, Size: 3},
		ExtendedAddressing:    {Prefix: 0x11, Opcode: 0xBC, Size: 4},
	},
	paramFunc: cmps,
}

// CmpuInst - Compare U (16-bit, page 3).
var CmpuInst = &Instruction{
	ID:   Cmpu,
	Name: CmpuName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x11, Opcode: 0x83, Size: 4},
		DirectAddressing:      {Prefix: 0x11, Opcode: 0x93, Size: 3},
		IndexedAddressing:     {Prefix: 0x11, Opcode: 0xA3, Size: 3},
		ExtendedAddressing:    {Prefix: 0x11, Opcode: 0xB3, Size: 4},
	},
	paramFunc: cmpu,
}

// CmpxInst - Compare X (16-bit).
var CmpxInst = &Instruction{
	ID:   Cmpx,
	Name: CmpxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0x8C, Size: 3},
		DirectAddressing:      {Opcode: 0x9C, Size: 2},
		IndexedAddressing:     {Opcode: 0xAC, Size: 2},
		ExtendedAddressing:    {Opcode: 0xBC, Size: 3},
	},
	paramFunc: cmpx,
}

// CmpyInst - Compare Y (16-bit, page 2).
var CmpyInst = &Instruction{
	ID:   Cmpy,
	Name: CmpyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x10, Opcode: 0x8C, Size: 4},
		DirectAddressing:      {Prefix: 0x10, Opcode: 0x9C, Size: 3},
		IndexedAddressing:     {Prefix: 0x10, Opcode: 0xAC, Size: 3},
		ExtendedAddressing:    {Prefix: 0x10, Opcode: 0xBC, Size: 4},
	},
	paramFunc: cmpy,
}

// ComInst - Complement (memory).
var ComInst = &Instruction{
	ID:   Com,
	Name: ComName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x03, Size: 2},
		IndexedAddressing:  {Opcode: 0x63, Size: 2},
		ExtendedAddressing: {Opcode: 0x73, Size: 3},
	},
	paramFunc: comMem,
}

// ComaInst - Complement A (inherent).
var ComaInst = &Instruction{
	ID:   Com,
	Name: ComName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x43, Size: 1},
	},
	noParamFunc: coma,
}

// CombInst - Complement B (inherent).
var CombInst = &Instruction{
	ID:   Com,
	Name: ComName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x53, Size: 1},
	},
	noParamFunc: comb,
}

// CwaiInst - AND CC then Wait for Interrupt.
var CwaiInst = &Instruction{
	ID:   Cwai,
	Name: CwaiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x3C, Size: 2},
	},
	paramFunc: cwaiFn,
}

// DaaInst - Decimal Adjust A.
var DaaInst = &Instruction{
	ID:   Daa,
	Name: DaaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x19, Size: 1},
	},
	noParamFunc: daaFn,
}

// DecInst - Decrement (memory).
var DecInst = &Instruction{
	ID:   Dec,
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0A, Size: 2},
		IndexedAddressing:  {Opcode: 0x6A, Size: 2},
		ExtendedAddressing: {Opcode: 0x7A, Size: 3},
	},
	paramFunc: decMem,
}

// DecaInst - Decrement A (inherent).
var DecaInst = &Instruction{
	ID:   Dec,
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x4A, Size: 1},
	},
	noParamFunc: deca,
}

// DecbInst - Decrement B (inherent).
var DecbInst = &Instruction{
	ID:   Dec,
	Name: DecName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5A, Size: 1},
	},
	noParamFunc: decb,
}

// EoraInst - Exclusive OR with A.
var EoraInst = &Instruction{
	ID:   Eora,
	Name: EoraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x88, Size: 2},
		DirectAddressing:    {Opcode: 0x98, Size: 2},
		IndexedAddressing:   {Opcode: 0xA8, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB8, Size: 3},
	},
	paramFunc: eora,
}

// EorbInst - Exclusive OR with B.
var EorbInst = &Instruction{
	ID:   Eorb,
	Name: EorbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC8, Size: 2},
		DirectAddressing:    {Opcode: 0xD8, Size: 2},
		IndexedAddressing:   {Opcode: 0xE8, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF8, Size: 3},
	},
	paramFunc: eorb,
}

// ExgInst - Exchange Registers.
var ExgInst = &Instruction{
	ID:   Exg,
	Name: ExgName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x1E, Size: 2},
	},
	paramFunc: exgFn,
}

// IncInst - Increment (memory).
var IncInst = &Instruction{
	ID:   Inc,
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0C, Size: 2},
		IndexedAddressing:  {Opcode: 0x6C, Size: 2},
		ExtendedAddressing: {Opcode: 0x7C, Size: 3},
	},
	paramFunc: incMem,
}

// IncaInst - Increment A (inherent).
var IncaInst = &Instruction{
	ID:   Inc,
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x4C, Size: 1},
	},
	noParamFunc: inca,
}

// IncbInst - Increment B (inherent).
var IncbInst = &Instruction{
	ID:   Inc,
	Name: IncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5C, Size: 1},
	},
	noParamFunc: incb,
}

// JmpInst - Jump.
var JmpInst = &Instruction{
	ID:   Jmp,
	Name: JmpName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0E, Size: 2},
		IndexedAddressing:  {Opcode: 0x6E, Size: 2},
		ExtendedAddressing: {Opcode: 0x7E, Size: 3},
	},
	paramFunc: jmpFn,
}

// JsrInst - Jump to Subroutine.
var JsrInst = &Instruction{
	ID:   Jsr,
	Name: JsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x9D, Size: 2},
		IndexedAddressing:  {Opcode: 0xAD, Size: 2},
		ExtendedAddressing: {Opcode: 0xBD, Size: 3},
	},
	paramFunc: jsrFn,
}

// LbccInst - Long Branch if Carry Clear (page 2).
var LbccInst = &Instruction{
	ID:         Lbcc,
	Name:       LbccName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x24, Size: 4}},
	paramFunc:  bccFn,
}

// LbcsInst - Long Branch if Carry Set (page 2).
var LbcsInst = &Instruction{
	ID:         Lbcs,
	Name:       LbcsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x25, Size: 4}},
	paramFunc:  bcsFn,
}

// LbeqInst - Long Branch if Equal (page 2).
var LbeqInst = &Instruction{
	ID:         Lbeq,
	Name:       LbeqName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x27, Size: 4}},
	paramFunc:  beqFn,
}

// LbgeInst - Long Branch if Greater or Equal (page 2).
var LbgeInst = &Instruction{
	ID:         Lbge,
	Name:       LbgeName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2C, Size: 4}},
	paramFunc:  bgeFn,
}

// LbgtInst - Long Branch if Greater Than (page 2).
var LbgtInst = &Instruction{
	ID:         Lbgt,
	Name:       LbgtName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2E, Size: 4}},
	paramFunc:  bgtFn,
}

// LbhiInst - Long Branch if Higher (page 2).
var LbhiInst = &Instruction{
	ID:         Lbhi,
	Name:       LbhiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x22, Size: 4}},
	paramFunc:  bhiFn,
}

// LbleInst - Long Branch if Less or Equal (page 2).
var LbleInst = &Instruction{
	ID:         Lble,
	Name:       LbleName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2F, Size: 4}},
	paramFunc:  bleFn,
}

// LblsInst - Long Branch if Lower or Same (page 2).
var LblsInst = &Instruction{
	ID:         Lbls,
	Name:       LblsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x23, Size: 4}},
	paramFunc:  blsFn,
}

// LbltInst - Long Branch if Less Than (page 2).
var LbltInst = &Instruction{
	ID:         Lblt,
	Name:       LbltName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2D, Size: 4}},
	paramFunc:  bltFn,
}

// LbmiInst - Long Branch if Minus (page 2).
var LbmiInst = &Instruction{
	ID:         Lbmi,
	Name:       LbmiName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2B, Size: 4}},
	paramFunc:  bmiFn,
}

// LbneInst - Long Branch if Not Equal (page 2).
var LbneInst = &Instruction{
	ID:         Lbne,
	Name:       LbneName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x26, Size: 4}},
	paramFunc:  bneFn,
}

// LbplInst - Long Branch if Plus (page 2).
var LbplInst = &Instruction{
	ID:         Lbpl,
	Name:       LbplName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x2A, Size: 4}},
	paramFunc:  bplFn,
}

// LbraInst - Long Branch Always.
var LbraInst = &Instruction{
	ID:         Lbra,
	Name:       LbraName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Opcode: 0x16, Size: 3}},
	paramFunc:  braFn,
}

// LbrnInst - Long Branch Never (page 2).
var LbrnInst = &Instruction{
	ID:         Lbrn,
	Name:       LbrnName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x21, Size: 4}},
	paramFunc:  brnFn,
}

// LbsrInst - Long Branch to Subroutine.
var LbsrInst = &Instruction{
	ID:         Lbsr,
	Name:       LbsrName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Opcode: 0x17, Size: 3}},
	paramFunc:  bsrFn,
}

// LbvcInst - Long Branch if Overflow Clear (page 2).
var LbvcInst = &Instruction{
	ID:         Lbvc,
	Name:       LbvcName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x28, Size: 4}},
	paramFunc:  bvcFn,
}

// LbvsInst - Long Branch if Overflow Set (page 2).
var LbvsInst = &Instruction{
	ID:         Lbvs,
	Name:       LbvsName,
	Addressing: map[AddressingMode]OpcodeInfo{RelativeLongAddressing: {Prefix: 0x10, Opcode: 0x29, Size: 4}},
	paramFunc:  bvsFn,
}

// LdaInst - Load A.
var LdaInst = &Instruction{
	ID:   Lda,
	Name: LdaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x86, Size: 2},
		DirectAddressing:    {Opcode: 0x96, Size: 2},
		IndexedAddressing:   {Opcode: 0xA6, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB6, Size: 3},
	},
	paramFunc: lda,
}

// LdbInst - Load B.
var LdbInst = &Instruction{
	ID:   Ldb,
	Name: LdbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC6, Size: 2},
		DirectAddressing:    {Opcode: 0xD6, Size: 2},
		IndexedAddressing:   {Opcode: 0xE6, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF6, Size: 3},
	},
	paramFunc: ldb,
}

// LddInst - Load D (16-bit).
var LddInst = &Instruction{
	ID:   Ldd,
	Name: LddName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0xCC, Size: 3},
		DirectAddressing:      {Opcode: 0xDC, Size: 2},
		IndexedAddressing:     {Opcode: 0xEC, Size: 2},
		ExtendedAddressing:    {Opcode: 0xFC, Size: 3},
	},
	paramFunc: ldd,
}

// LdsInst - Load S (16-bit, page 2).
var LdsInst = &Instruction{
	ID:   Lds,
	Name: LdsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x10, Opcode: 0xCE, Size: 4},
		DirectAddressing:      {Prefix: 0x10, Opcode: 0xDE, Size: 3},
		IndexedAddressing:     {Prefix: 0x10, Opcode: 0xEE, Size: 3},
		ExtendedAddressing:    {Prefix: 0x10, Opcode: 0xFE, Size: 4},
	},
	paramFunc: lds,
}

// LduInst - Load U (16-bit).
var LduInst = &Instruction{
	ID:   Ldu,
	Name: LduName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0xCE, Size: 3},
		DirectAddressing:      {Opcode: 0xDE, Size: 2},
		IndexedAddressing:     {Opcode: 0xEE, Size: 2},
		ExtendedAddressing:    {Opcode: 0xFE, Size: 3},
	},
	paramFunc: ldu,
}

// LdxInst - Load X (16-bit).
var LdxInst = &Instruction{
	ID:   Ldx,
	Name: LdxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0x8E, Size: 3},
		DirectAddressing:      {Opcode: 0x9E, Size: 2},
		IndexedAddressing:     {Opcode: 0xAE, Size: 2},
		ExtendedAddressing:    {Opcode: 0xBE, Size: 3},
	},
	paramFunc: ldx,
}

// LdyInst - Load Y (16-bit, page 2).
var LdyInst = &Instruction{
	ID:   Ldy,
	Name: LdyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Prefix: 0x10, Opcode: 0x8E, Size: 4},
		DirectAddressing:      {Prefix: 0x10, Opcode: 0x9E, Size: 3},
		IndexedAddressing:     {Prefix: 0x10, Opcode: 0xAE, Size: 3},
		ExtendedAddressing:    {Prefix: 0x10, Opcode: 0xBE, Size: 4},
	},
	paramFunc: ldy,
}

// LeaxInst - Load Effective Address into X.
var LeaxInst = &Instruction{
	ID:   Leax,
	Name: LeaxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		IndexedAddressing: {Opcode: 0x30, Size: 2},
	},
	paramFunc: leax,
}

// LeayInst - Load Effective Address into Y.
var LeayInst = &Instruction{
	ID:   Leay,
	Name: LeayName,
	Addressing: map[AddressingMode]OpcodeInfo{
		IndexedAddressing: {Opcode: 0x31, Size: 2},
	},
	paramFunc: leay,
}

// LeasInst - Load Effective Address into S.
var LeasInst = &Instruction{
	ID:   Leas,
	Name: LeasName,
	Addressing: map[AddressingMode]OpcodeInfo{
		IndexedAddressing: {Opcode: 0x32, Size: 2},
	},
	paramFunc: leas,
}

// LeauInst - Load Effective Address into U.
var LeauInst = &Instruction{
	ID:   Leau,
	Name: LeauName,
	Addressing: map[AddressingMode]OpcodeInfo{
		IndexedAddressing: {Opcode: 0x33, Size: 2},
	},
	paramFunc: leau,
}

// LsrInst - Logical Shift Right (memory).
var LsrInst = &Instruction{
	ID:   Lsr,
	Name: LsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x04, Size: 2},
		IndexedAddressing:  {Opcode: 0x64, Size: 2},
		ExtendedAddressing: {Opcode: 0x74, Size: 3},
	},
	paramFunc: lsrMem,
}

// LsraInst - Logical Shift Right A (inherent).
var LsraInst = &Instruction{
	ID:   Lsr,
	Name: LsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x44, Size: 1},
	},
	noParamFunc: lsra,
}

// LsrbInst - Logical Shift Right B (inherent).
var LsrbInst = &Instruction{
	ID:   Lsr,
	Name: LsrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x54, Size: 1},
	},
	noParamFunc: lsrb,
}

// MulInst - Multiply (unsigned A*B -> D).
var MulInst = &Instruction{
	ID:   Mul,
	Name: MulName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3D, Size: 1},
	},
	noParamFunc: mulFn,
}

// NegInst - Negate (memory).
var NegInst = &Instruction{
	ID:   Neg,
	Name: NegName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x00, Size: 2},
		IndexedAddressing:  {Opcode: 0x60, Size: 2},
		ExtendedAddressing: {Opcode: 0x70, Size: 3},
	},
	paramFunc: negMem,
}

// NegaInst - Negate A (inherent).
var NegaInst = &Instruction{
	ID:   Neg,
	Name: NegName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x40, Size: 1},
	},
	noParamFunc: nega,
}

// NegbInst - Negate B (inherent).
var NegbInst = &Instruction{
	ID:   Neg,
	Name: NegName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x50, Size: 1},
	},
	noParamFunc: negb,
}

// NopInst - No Operation.
var NopInst = &Instruction{
	ID:   Nop,
	Name: NopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x12, Size: 1},
	},
	noParamFunc: nop,
}

// OraInst - OR with A.
var OraInst = &Instruction{
	ID:   Ora,
	Name: OraName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x8A, Size: 2},
		DirectAddressing:    {Opcode: 0x9A, Size: 2},
		IndexedAddressing:   {Opcode: 0xAA, Size: 2},
		ExtendedAddressing:  {Opcode: 0xBA, Size: 3},
	},
	paramFunc: ora,
}

// OrbInst - OR with B.
var OrbInst = &Instruction{
	ID:   Orb,
	Name: OrbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xCA, Size: 2},
		DirectAddressing:    {Opcode: 0xDA, Size: 2},
		IndexedAddressing:   {Opcode: 0xEA, Size: 2},
		ExtendedAddressing:  {Opcode: 0xFA, Size: 3},
	},
	paramFunc: orb,
}

// OrccInst - OR CC Register.
var OrccInst = &Instruction{
	ID:   Orcc,
	Name: OrccName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x1A, Size: 2},
	},
	paramFunc: orcc,
}

// PshsInst - Push registers onto S stack.
var PshsInst = &Instruction{
	ID:   Pshs,
	Name: PshsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		StackAddressing: {Opcode: 0x34, Size: 2},
	},
	paramFunc: pshsFn,
}

// PshuInst - Push registers onto U stack.
var PshuInst = &Instruction{
	ID:   Pshu,
	Name: PshuName,
	Addressing: map[AddressingMode]OpcodeInfo{
		StackAddressing: {Opcode: 0x36, Size: 2},
	},
	paramFunc: pshuFn,
}

// PulsInst - Pull registers from S stack.
var PulsInst = &Instruction{
	ID:   Puls,
	Name: PulsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		StackAddressing: {Opcode: 0x35, Size: 2},
	},
	paramFunc: pulsFn,
}

// PuluInst - Pull registers from U stack.
var PuluInst = &Instruction{
	ID:   Pulu,
	Name: PuluName,
	Addressing: map[AddressingMode]OpcodeInfo{
		StackAddressing: {Opcode: 0x37, Size: 2},
	},
	paramFunc: puluFn,
}

// RolInst - Rotate Left (memory).
var RolInst = &Instruction{
	ID:   Rol,
	Name: RolName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x09, Size: 2},
		IndexedAddressing:  {Opcode: 0x69, Size: 2},
		ExtendedAddressing: {Opcode: 0x79, Size: 3},
	},
	paramFunc: rolMem,
}

// RolaInst - Rotate Left A (inherent).
var RolaInst = &Instruction{
	ID:   Rol,
	Name: RolName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x49, Size: 1},
	},
	noParamFunc: rola,
}

// RolbInst - Rotate Left B (inherent).
var RolbInst = &Instruction{
	ID:   Rol,
	Name: RolName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x59, Size: 1},
	},
	noParamFunc: rolb,
}

// RorInst - Rotate Right (memory).
var RorInst = &Instruction{
	ID:   Ror,
	Name: RorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x06, Size: 2},
		IndexedAddressing:  {Opcode: 0x66, Size: 2},
		ExtendedAddressing: {Opcode: 0x76, Size: 3},
	},
	paramFunc: rorMem,
}

// RoraInst - Rotate Right A (inherent).
var RoraInst = &Instruction{
	ID:   Ror,
	Name: RorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x46, Size: 1},
	},
	noParamFunc: rora,
}

// RorbInst - Rotate Right B (inherent).
var RorbInst = &Instruction{
	ID:   Ror,
	Name: RorName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x56, Size: 1},
	},
	noParamFunc: rorb,
}

// RtiInst - Return from Interrupt.
var RtiInst = &Instruction{
	ID:   Rti,
	Name: RtiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3B, Size: 1},
	},
	noParamFunc: rtiFn,
}

// RtsInst - Return from Subroutine.
var RtsInst = &Instruction{
	ID:   Rts,
	Name: RtsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x39, Size: 1},
	},
	noParamFunc: rtsFn,
}

// SbcaInst - Subtract with Carry from A.
var SbcaInst = &Instruction{
	ID:   Sbca,
	Name: SbcaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x82, Size: 2},
		DirectAddressing:    {Opcode: 0x92, Size: 2},
		IndexedAddressing:   {Opcode: 0xA2, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB2, Size: 3},
	},
	paramFunc: sbca,
}

// SbcbInst - Subtract with Carry from B.
var SbcbInst = &Instruction{
	ID:   Sbcb,
	Name: SbcbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC2, Size: 2},
		DirectAddressing:    {Opcode: 0xD2, Size: 2},
		IndexedAddressing:   {Opcode: 0xE2, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF2, Size: 3},
	},
	paramFunc: sbcb,
}

// SexInst - Sign Extend B into A.
var SexInst = &Instruction{
	ID:   Sex,
	Name: SexName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x1D, Size: 1},
	},
	noParamFunc: sexFn,
}

// StaInst - Store A.
var StaInst = &Instruction{
	ID:   Sta,
	Name: StaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x97, Size: 2},
		IndexedAddressing:  {Opcode: 0xA7, Size: 2},
		ExtendedAddressing: {Opcode: 0xB7, Size: 3},
	},
	paramFunc: sta,
}

// StbInst - Store B.
var StbInst = &Instruction{
	ID:   Stb,
	Name: StbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0xD7, Size: 2},
		IndexedAddressing:  {Opcode: 0xE7, Size: 2},
		ExtendedAddressing: {Opcode: 0xF7, Size: 3},
	},
	paramFunc: stb,
}

// StdInst - Store D (16-bit).
var StdInst = &Instruction{
	ID:   Std,
	Name: StdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0xDD, Size: 2},
		IndexedAddressing:  {Opcode: 0xED, Size: 2},
		ExtendedAddressing: {Opcode: 0xFD, Size: 3},
	},
	paramFunc: std,
}

// StsInst - Store S (16-bit, page 2).
var StsInst = &Instruction{
	ID:   Sts,
	Name: StsName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Prefix: 0x10, Opcode: 0xDF, Size: 3},
		IndexedAddressing:  {Prefix: 0x10, Opcode: 0xEF, Size: 3},
		ExtendedAddressing: {Prefix: 0x10, Opcode: 0xFF, Size: 4},
	},
	paramFunc: sts,
}

// StuInst - Store U (16-bit).
var StuInst = &Instruction{
	ID:   Stu,
	Name: StuName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0xDF, Size: 2},
		IndexedAddressing:  {Opcode: 0xEF, Size: 2},
		ExtendedAddressing: {Opcode: 0xFF, Size: 3},
	},
	paramFunc: stu,
}

// StxInst - Store X (16-bit).
var StxInst = &Instruction{
	ID:   Stx,
	Name: StxName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x9F, Size: 2},
		IndexedAddressing:  {Opcode: 0xAF, Size: 2},
		ExtendedAddressing: {Opcode: 0xBF, Size: 3},
	},
	paramFunc: stx,
}

// StyInst - Store Y (16-bit, page 2).
var StyInst = &Instruction{
	ID:   Sty,
	Name: StyName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Prefix: 0x10, Opcode: 0x9F, Size: 3},
		IndexedAddressing:  {Prefix: 0x10, Opcode: 0xAF, Size: 3},
		ExtendedAddressing: {Prefix: 0x10, Opcode: 0xBF, Size: 4},
	},
	paramFunc: sty,
}

// SubaInst - Subtract from A.
var SubaInst = &Instruction{
	ID:   Suba,
	Name: SubaName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0x80, Size: 2},
		DirectAddressing:    {Opcode: 0x90, Size: 2},
		IndexedAddressing:   {Opcode: 0xA0, Size: 2},
		ExtendedAddressing:  {Opcode: 0xB0, Size: 3},
	},
	paramFunc: suba,
}

// SubbInst - Subtract from B.
var SubbInst = &Instruction{
	ID:   Subb,
	Name: SubbName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImmediateAddressing: {Opcode: 0xC0, Size: 2},
		DirectAddressing:    {Opcode: 0xD0, Size: 2},
		IndexedAddressing:   {Opcode: 0xE0, Size: 2},
		ExtendedAddressing:  {Opcode: 0xF0, Size: 3},
	},
	paramFunc: subb,
}

// SubdInst - Subtract from D (16-bit).
var SubdInst = &Instruction{
	ID:   Subd,
	Name: SubdName,
	Addressing: map[AddressingMode]OpcodeInfo{
		Immediate16Addressing: {Opcode: 0x83, Size: 3},
		DirectAddressing:      {Opcode: 0x93, Size: 2},
		IndexedAddressing:     {Opcode: 0xA3, Size: 2},
		ExtendedAddressing:    {Opcode: 0xB3, Size: 3},
	},
	paramFunc: subd,
}

// SwiInst - Software Interrupt 1.
var SwiInst = &Instruction{
	ID:   Swi,
	Name: SwiName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x3F, Size: 1},
	},
	noParamFunc: swiFn,
}

// Swi2Inst - Software Interrupt 2 (page 2).
var Swi2Inst = &Instruction{
	ID:   Swi2,
	Name: Swi2Name,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x10, Opcode: 0x3F, Size: 2},
	},
	noParamFunc: swi2Fn,
}

// Swi3Inst - Software Interrupt 3 (page 3).
var Swi3Inst = &Instruction{
	ID:   Swi3,
	Name: Swi3Name,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Prefix: 0x11, Opcode: 0x3F, Size: 2},
	},
	noParamFunc: swi3Fn,
}

// SyncInst - Synchronize with Interrupt.
var SyncInst = &Instruction{
	ID:   Sync,
	Name: SyncName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x13, Size: 1},
	},
	noParamFunc: syncFn,
}

// TfrInst - Transfer Register to Register.
var TfrInst = &Instruction{
	ID:   Tfr,
	Name: TfrName,
	Addressing: map[AddressingMode]OpcodeInfo{
		RegisterAddressing: {Opcode: 0x1F, Size: 2},
	},
	paramFunc: tfrFn,
}

// TstInst - Test (memory).
var TstInst = &Instruction{
	ID:   Tst,
	Name: TstName,
	Addressing: map[AddressingMode]OpcodeInfo{
		DirectAddressing:   {Opcode: 0x0D, Size: 2},
		IndexedAddressing:  {Opcode: 0x6D, Size: 2},
		ExtendedAddressing: {Opcode: 0x7D, Size: 3},
	},
	paramFunc: tstMem,
}

// TstaInst - Test A (inherent).
var TstaInst = &Instruction{
	ID:   Tst,
	Name: TstName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x4D, Size: 1},
	},
	noParamFunc: tsta,
}

// TstbInst - Test B (inherent).
var TstbInst = &Instruction{
	ID:   Tst,
	Name: TstName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0x5D, Size: 1},
	},
	noParamFunc: tstb,
}
