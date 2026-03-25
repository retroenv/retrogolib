package m6809

import "github.com/retroenv/retrogolib/set"

// MaxOpcodeSize is the maximum instruction size in bytes (prefix + opcode + postbyte + 16-bit operand).
const MaxOpcodeSize = 5

// Prefix bytes for extended opcode pages.
const (
	Prefix10 = 0x10 // Page 2 prefix
	Prefix11 = 0x11 // Page 3 prefix
)

// Opcode contains decoded instruction information for a single opcode byte.
type Opcode struct {
	Instruction *Instruction
	Addressing  AddressingMode
	Timing      byte // Base cycle count
	Size        byte // Base size in bytes (not including indexed postbyte extra)
}

// ReadsMemory returns true if this opcode reads from memory.
func (op Opcode) ReadsMemory(memReadInstructions set.Set[string]) bool {
	switch op.Addressing {
	case ImmediateAddressing, Immediate16Addressing, ImpliedAddressing, RelativeAddressing, RelativeLongAddressing:
		return false
	}
	return memReadInstructions.Contains(op.Instruction.Name)
}

// WritesMemory returns true if this opcode writes to memory.
func (op Opcode) WritesMemory(memWriteInstructions set.Set[string]) bool {
	switch op.Addressing {
	case ImmediateAddressing, Immediate16Addressing, ImpliedAddressing, RelativeAddressing, RelativeLongAddressing:
		return false
	}
	return memWriteInstructions.Contains(op.Instruction.Name)
}

// IsBranching returns true if this opcode is a branching instruction.
func (op Opcode) IsBranching(branchInstructions set.Set[string]) bool {
	return branchInstructions.Contains(op.Instruction.Name)
}

// GetOpcodeInfo returns opcode information for the given byte.
// Returns the Opcode and true if defined, or zero Opcode and false if undefined.
func GetOpcodeInfo(b uint8) (Opcode, bool) {
	op := Opcodes[b]
	if op.Instruction == nil {
		return Opcode{}, false
	}
	return op, true
}

// GetPage2OpcodeInfo returns opcode information for a page 2 ($10 prefix) opcode.
func GetPage2OpcodeInfo(b uint8) (Opcode, bool) {
	op := OpcodesPage2[b]
	if op.Instruction == nil {
		return Opcode{}, false
	}
	return op, true
}

// GetPage3OpcodeInfo returns opcode information for a page 3 ($11 prefix) opcode.
func GetPage3OpcodeInfo(b uint8) (Opcode, bool) {
	op := OpcodesPage3[b]
	if op.Instruction == nil {
		return Opcode{}, false
	}
	return op, true
}

// Opcodes is the base page opcode table for the 6809.
// Reference: Motorola MC6809 Programming Manual.
var Opcodes = [256]Opcode{
	// $00 - $0F: Direct page operations
	{Instruction: NegInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $00 NEG direct
	{}, // $01 (illegal)
	{}, // $02 (illegal)
	{Instruction: ComInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $03 COM direct
	{Instruction: LsrInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $04 LSR direct
	{}, // $05 (illegal)
	{Instruction: RorInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $06 ROR direct
	{Instruction: AsrInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $07 ASR direct
	{Instruction: AslInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $08 ASL/LSL direct
	{Instruction: RolInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $09 ROL direct
	{Instruction: DecInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $0A DEC direct
	{}, // $0B (illegal)
	{Instruction: IncInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $0C INC direct
	{Instruction: TstInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $0D TST direct
	{Instruction: JmpInst, Addressing: DirectAddressing, Timing: 3, Size: 2}, // $0E JMP direct
	{Instruction: ClrInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $0F CLR direct

	// $10 - $1F: Miscellaneous / prefix
	{}, // $10 prefix page 2 (handled in Step)
	{}, // $11 prefix page 3 (handled in Step)
	{Instruction: NopInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1},  // $12 NOP
	{Instruction: SyncInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1}, // $13 SYNC
	{}, // $14 (illegal)
	{}, // $15 (illegal)
	{Instruction: LbraInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 3}, // $16 LBRA
	{Instruction: LbsrInst, Addressing: RelativeLongAddressing, Timing: 9, Size: 3}, // $17 LBSR
	{}, // $18 (illegal)
	{Instruction: DaaInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1},    // $19 DAA
	{Instruction: OrccInst, Addressing: ImmediateAddressing, Timing: 3, Size: 2}, // $1A ORCC
	{}, // $1B (illegal)
	{Instruction: AndccInst, Addressing: ImmediateAddressing, Timing: 3, Size: 2}, // $1C ANDCC
	{Instruction: SexInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1},     // $1D SEX
	{Instruction: ExgInst, Addressing: RegisterAddressing, Timing: 8, Size: 2},    // $1E EXG
	{Instruction: TfrInst, Addressing: RegisterAddressing, Timing: 6, Size: 2},    // $1F TFR

	// $20 - $2F: Short branches
	{Instruction: BraInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $20 BRA
	{Instruction: BrnInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $21 BRN
	{Instruction: BhiInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $22 BHI
	{Instruction: BlsInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $23 BLS
	{Instruction: BccInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $24 BCC/BHS
	{Instruction: BcsInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $25 BCS/BLO
	{Instruction: BneInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $26 BNE
	{Instruction: BeqInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $27 BEQ
	{Instruction: BvcInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $28 BVC
	{Instruction: BvsInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $29 BVS
	{Instruction: BplInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $2A BPL
	{Instruction: BmiInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $2B BMI
	{Instruction: BgeInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $2C BGE
	{Instruction: BltInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $2D BLT
	{Instruction: BgtInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $2E BGT
	{Instruction: BleInst, Addressing: RelativeAddressing, Timing: 3, Size: 2}, // $2F BLE

	// $30 - $3F: Indexed / stack / system
	{Instruction: LeaxInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $30 LEAX
	{Instruction: LeayInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $31 LEAY
	{Instruction: LeasInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $32 LEAS
	{Instruction: LeauInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $33 LEAU
	{Instruction: PshsInst, Addressing: StackAddressing, Timing: 5, Size: 2},   // $34 PSHS
	{Instruction: PulsInst, Addressing: StackAddressing, Timing: 5, Size: 2},   // $35 PULS
	{Instruction: PshuInst, Addressing: StackAddressing, Timing: 5, Size: 2},   // $36 PSHU
	{Instruction: PuluInst, Addressing: StackAddressing, Timing: 5, Size: 2},   // $37 PULU
	{}, // $38 (illegal)
	{Instruction: RtsInst, Addressing: ImpliedAddressing, Timing: 5, Size: 1},     // $39 RTS
	{Instruction: AbxInst, Addressing: ImpliedAddressing, Timing: 3, Size: 1},     // $3A ABX
	{Instruction: RtiInst, Addressing: ImpliedAddressing, Timing: 6, Size: 1},     // $3B RTI
	{Instruction: CwaiInst, Addressing: ImmediateAddressing, Timing: 20, Size: 2}, // $3C CWAI
	{Instruction: MulInst, Addressing: ImpliedAddressing, Timing: 11, Size: 1},    // $3D MUL
	{}, // $3E (illegal)
	{Instruction: SwiInst, Addressing: ImpliedAddressing, Timing: 19, Size: 1}, // $3F SWI

	// $40 - $4F: Inherent A operations
	{Instruction: NegaInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $40 NEGA
	{}, // $41 (illegal)
	{}, // $42 (illegal)
	{Instruction: ComaInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $43 COMA
	{Instruction: LsraInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $44 LSRA
	{}, // $45 (illegal)
	{Instruction: RoraInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $46 RORA
	{Instruction: AsraInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $47 ASRA
	{Instruction: AslaInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $48 ASLA/LSLA
	{Instruction: RolaInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $49 ROLA
	{Instruction: DecaInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $4A DECA
	{}, // $4B (illegal)
	{Instruction: IncaInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $4C INCA
	{Instruction: TstaInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $4D TSTA
	{}, // $4E (illegal)
	{Instruction: ClraInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $4F CLRA

	// $50 - $5F: Inherent B operations
	{Instruction: NegbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $50 NEGB
	{}, // $51 (illegal)
	{}, // $52 (illegal)
	{Instruction: CombInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $53 COMB
	{Instruction: LsrbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $54 LSRB
	{}, // $55 (illegal)
	{Instruction: RorbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $56 RORB
	{Instruction: AsrbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $57 ASRB
	{Instruction: AslbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $58 ASLB/LSLB
	{Instruction: RolbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $59 ROLB
	{Instruction: DecbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $5A DECB
	{}, // $5B (illegal)
	{Instruction: IncbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $5C INCB
	{Instruction: TstbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $5D TSTB
	{}, // $5E (illegal)
	{Instruction: ClrbInst, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // $5F CLRB

	// $60 - $6F: Indexed memory operations
	{Instruction: NegInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $60 NEG indexed
	{}, // $61 (illegal)
	{}, // $62 (illegal)
	{Instruction: ComInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $63 COM indexed
	{Instruction: LsrInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $64 LSR indexed
	{}, // $65 (illegal)
	{Instruction: RorInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $66 ROR indexed
	{Instruction: AsrInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $67 ASR indexed
	{Instruction: AslInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $68 ASL indexed
	{Instruction: RolInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $69 ROL indexed
	{Instruction: DecInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $6A DEC indexed
	{}, // $6B (illegal)
	{Instruction: IncInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $6C INC indexed
	{Instruction: TstInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $6D TST indexed
	{Instruction: JmpInst, Addressing: IndexedAddressing, Timing: 3, Size: 2}, // $6E JMP indexed
	{Instruction: ClrInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $6F CLR indexed

	// $70 - $7F: Extended memory operations
	{Instruction: NegInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $70 NEG extended
	{}, // $71 (illegal)
	{}, // $72 (illegal)
	{Instruction: ComInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $73 COM extended
	{Instruction: LsrInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $74 LSR extended
	{}, // $75 (illegal)
	{Instruction: RorInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $76 ROR extended
	{Instruction: AsrInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $77 ASR extended
	{Instruction: AslInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $78 ASL extended
	{Instruction: RolInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $79 ROL extended
	{Instruction: DecInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $7A DEC extended
	{}, // $7B (illegal)
	{Instruction: IncInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $7C INC extended
	{Instruction: TstInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $7D TST extended
	{Instruction: JmpInst, Addressing: ExtendedAddressing, Timing: 4, Size: 3}, // $7E JMP extended
	{Instruction: ClrInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $7F CLR extended

	// $80 - $8F: A register immediate/relative
	{Instruction: SubaInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $80 SUBA imm
	{Instruction: CmpaInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $81 CMPA imm
	{Instruction: SbcaInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $82 SBCA imm
	{Instruction: SubdInst, Addressing: Immediate16Addressing, Timing: 4, Size: 3}, // $83 SUBD imm16
	{Instruction: AndaInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $84 ANDA imm
	{Instruction: BitaInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $85 BITA imm
	{Instruction: LdaInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},    // $86 LDA imm
	{}, // $87 (illegal)
	{Instruction: EoraInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $88 EORA imm
	{Instruction: AdcaInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $89 ADCA imm
	{Instruction: OraInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},    // $8A ORA imm
	{Instruction: AddaInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $8B ADDA imm
	{Instruction: CmpxInst, Addressing: Immediate16Addressing, Timing: 4, Size: 3}, // $8C CMPX imm16
	{Instruction: BsrInst, Addressing: RelativeAddressing, Timing: 7, Size: 2},     // $8D BSR rel
	{Instruction: LdxInst, Addressing: Immediate16Addressing, Timing: 3, Size: 3},  // $8E LDX imm16
	{}, // $8F (illegal)

	// $90 - $9F: A register direct page
	{Instruction: SubaInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $90 SUBA direct
	{Instruction: CmpaInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $91 CMPA direct
	{Instruction: SbcaInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $92 SBCA direct
	{Instruction: SubdInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $93 SUBD direct
	{Instruction: AndaInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $94 ANDA direct
	{Instruction: BitaInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $95 BITA direct
	{Instruction: LdaInst, Addressing: DirectAddressing, Timing: 4, Size: 2},  // $96 LDA direct
	{Instruction: StaInst, Addressing: DirectAddressing, Timing: 4, Size: 2},  // $97 STA direct
	{Instruction: EoraInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $98 EORA direct
	{Instruction: AdcaInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $99 ADCA direct
	{Instruction: OraInst, Addressing: DirectAddressing, Timing: 4, Size: 2},  // $9A ORA direct
	{Instruction: AddaInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $9B ADDA direct
	{Instruction: CmpxInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $9C CMPX direct
	{Instruction: JsrInst, Addressing: DirectAddressing, Timing: 7, Size: 2},  // $9D JSR direct
	{Instruction: LdxInst, Addressing: DirectAddressing, Timing: 5, Size: 2},  // $9E LDX direct
	{Instruction: StxInst, Addressing: DirectAddressing, Timing: 5, Size: 2},  // $9F STX direct

	// $A0 - $AF: A register indexed
	{Instruction: SubaInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $A0 SUBA indexed
	{Instruction: CmpaInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $A1 CMPA indexed
	{Instruction: SbcaInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $A2 SBCA indexed
	{Instruction: SubdInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $A3 SUBD indexed
	{Instruction: AndaInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $A4 ANDA indexed
	{Instruction: BitaInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $A5 BITA indexed
	{Instruction: LdaInst, Addressing: IndexedAddressing, Timing: 4, Size: 2},  // $A6 LDA indexed
	{Instruction: StaInst, Addressing: IndexedAddressing, Timing: 4, Size: 2},  // $A7 STA indexed
	{Instruction: EoraInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $A8 EORA indexed
	{Instruction: AdcaInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $A9 ADCA indexed
	{Instruction: OraInst, Addressing: IndexedAddressing, Timing: 4, Size: 2},  // $AA ORA indexed
	{Instruction: AddaInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $AB ADDA indexed
	{Instruction: CmpxInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $AC CMPX indexed
	{Instruction: JsrInst, Addressing: IndexedAddressing, Timing: 7, Size: 2},  // $AD JSR indexed
	{Instruction: LdxInst, Addressing: IndexedAddressing, Timing: 5, Size: 2},  // $AE LDX indexed
	{Instruction: StxInst, Addressing: IndexedAddressing, Timing: 5, Size: 2},  // $AF STX indexed

	// $B0 - $BF: A register extended
	{Instruction: SubaInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $B0 SUBA extended
	{Instruction: CmpaInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $B1 CMPA extended
	{Instruction: SbcaInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $B2 SBCA extended
	{Instruction: SubdInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $B3 SUBD extended
	{Instruction: AndaInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $B4 ANDA extended
	{Instruction: BitaInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $B5 BITA extended
	{Instruction: LdaInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3},  // $B6 LDA extended
	{Instruction: StaInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3},  // $B7 STA extended
	{Instruction: EoraInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $B8 EORA extended
	{Instruction: AdcaInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $B9 ADCA extended
	{Instruction: OraInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3},  // $BA ORA extended
	{Instruction: AddaInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $BB ADDA extended
	{Instruction: CmpxInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $BC CMPX extended
	{Instruction: JsrInst, Addressing: ExtendedAddressing, Timing: 8, Size: 3},  // $BD JSR extended
	{Instruction: LdxInst, Addressing: ExtendedAddressing, Timing: 6, Size: 3},  // $BE LDX extended
	{Instruction: StxInst, Addressing: ExtendedAddressing, Timing: 6, Size: 3},  // $BF STX extended

	// $C0 - $CF: B register immediate
	{Instruction: SubbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $C0 SUBB imm
	{Instruction: CmpbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $C1 CMPB imm
	{Instruction: SbcbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $C2 SBCB imm
	{Instruction: AdddInst, Addressing: Immediate16Addressing, Timing: 4, Size: 3}, // $C3 ADDD imm16
	{Instruction: AndbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $C4 ANDB imm
	{Instruction: BitbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $C5 BITB imm
	{Instruction: LdbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},    // $C6 LDB imm
	{}, // $C7 (illegal)
	{Instruction: EorbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},  // $C8 EORB imm
	{Instruction: AdcbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},  // $C9 ADCB imm
	{Instruction: OrbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // $CA ORB imm
	{Instruction: AddbInst, Addressing: ImmediateAddressing, Timing: 2, Size: 2},  // $CB ADDB imm
	{Instruction: LddInst, Addressing: Immediate16Addressing, Timing: 3, Size: 3}, // $CC LDD imm16
	{}, // $CD (illegal)
	{Instruction: LduInst, Addressing: Immediate16Addressing, Timing: 3, Size: 3}, // $CE LDU imm16
	{}, // $CF (illegal)

	// $D0 - $DF: B register direct page
	{Instruction: SubbInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $D0 SUBB direct
	{Instruction: CmpbInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $D1 CMPB direct
	{Instruction: SbcbInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $D2 SBCB direct
	{Instruction: AdddInst, Addressing: DirectAddressing, Timing: 6, Size: 2}, // $D3 ADDD direct
	{Instruction: AndbInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $D4 ANDB direct
	{Instruction: BitbInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $D5 BITB direct
	{Instruction: LdbInst, Addressing: DirectAddressing, Timing: 4, Size: 2},  // $D6 LDB direct
	{Instruction: StbInst, Addressing: DirectAddressing, Timing: 4, Size: 2},  // $D7 STB direct
	{Instruction: EorbInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $D8 EORB direct
	{Instruction: AdcbInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $D9 ADCB direct
	{Instruction: OrbInst, Addressing: DirectAddressing, Timing: 4, Size: 2},  // $DA ORB direct
	{Instruction: AddbInst, Addressing: DirectAddressing, Timing: 4, Size: 2}, // $DB ADDB direct
	{Instruction: LddInst, Addressing: DirectAddressing, Timing: 5, Size: 2},  // $DC LDD direct
	{Instruction: StdInst, Addressing: DirectAddressing, Timing: 5, Size: 2},  // $DD STD direct
	{Instruction: LduInst, Addressing: DirectAddressing, Timing: 5, Size: 2},  // $DE LDU direct
	{Instruction: StuInst, Addressing: DirectAddressing, Timing: 5, Size: 2},  // $DF STU direct

	// $E0 - $EF: B register indexed
	{Instruction: SubbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $E0 SUBB indexed
	{Instruction: CmpbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $E1 CMPB indexed
	{Instruction: SbcbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $E2 SBCB indexed
	{Instruction: AdddInst, Addressing: IndexedAddressing, Timing: 6, Size: 2}, // $E3 ADDD indexed
	{Instruction: AndbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $E4 ANDB indexed
	{Instruction: BitbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $E5 BITB indexed
	{Instruction: LdbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2},  // $E6 LDB indexed
	{Instruction: StbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2},  // $E7 STB indexed
	{Instruction: EorbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $E8 EORB indexed
	{Instruction: AdcbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $E9 ADCB indexed
	{Instruction: OrbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2},  // $EA ORB indexed
	{Instruction: AddbInst, Addressing: IndexedAddressing, Timing: 4, Size: 2}, // $EB ADDB indexed
	{Instruction: LddInst, Addressing: IndexedAddressing, Timing: 5, Size: 2},  // $EC LDD indexed
	{Instruction: StdInst, Addressing: IndexedAddressing, Timing: 5, Size: 2},  // $ED STD indexed
	{Instruction: LduInst, Addressing: IndexedAddressing, Timing: 5, Size: 2},  // $EE LDU indexed
	{Instruction: StuInst, Addressing: IndexedAddressing, Timing: 5, Size: 2},  // $EF STU indexed

	// $F0 - $FF: B register extended
	{Instruction: SubbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $F0 SUBB extended
	{Instruction: CmpbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $F1 CMPB extended
	{Instruction: SbcbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $F2 SBCB extended
	{Instruction: AdddInst, Addressing: ExtendedAddressing, Timing: 7, Size: 3}, // $F3 ADDD extended
	{Instruction: AndbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $F4 ANDB extended
	{Instruction: BitbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $F5 BITB extended
	{Instruction: LdbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3},  // $F6 LDB extended
	{Instruction: StbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3},  // $F7 STB extended
	{Instruction: EorbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $F8 EORB extended
	{Instruction: AdcbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $F9 ADCB extended
	{Instruction: OrbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3},  // $FA ORB extended
	{Instruction: AddbInst, Addressing: ExtendedAddressing, Timing: 5, Size: 3}, // $FB ADDB extended
	{Instruction: LddInst, Addressing: ExtendedAddressing, Timing: 6, Size: 3},  // $FC LDD extended
	{Instruction: StdInst, Addressing: ExtendedAddressing, Timing: 6, Size: 3},  // $FD STD extended
	{Instruction: LduInst, Addressing: ExtendedAddressing, Timing: 6, Size: 3},  // $FE LDU extended
	{Instruction: StuInst, Addressing: ExtendedAddressing, Timing: 6, Size: 3},  // $FF STU extended
}
