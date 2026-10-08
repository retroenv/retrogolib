package cpu65816

import "github.com/retroenv/retrogolib/set"

// MaxOpcodeSize is the maximum instruction size in bytes (abs long + 16-bit immediate).
const MaxOpcodeSize = 4

// WidthFlag indicates which processor flag controls an instruction's operand width.
type WidthFlag byte

const (
	WidthNone WidthFlag = iota // Fixed size regardless of flags
	WidthM                     // Size increases by 1 when M=0 (16-bit accumulator)
	WidthX                     // Size increases by 1 when X=0 (16-bit index registers)
)

// CycleRule is a set of conditions that add cycles to the base timing of an opcode.
// The rules follow the W65C816S datasheet notes for the instruction table.
type CycleRule uint8

const (
	CycleM      CycleRule = 1 << iota // Add 1 cycle when M=0 (16-bit accumulator access).
	CycleRMW                          // Add 2 cycles when M=0 (16-bit read-modify-write).
	CycleX                            // Add 1 cycle when X=0 (16-bit index register access).
	CycleDL                           // Add 1 cycle when the low byte of DP is not zero.
	CycleIndex                        // Add 1 cycle when X=0 or the index crosses a page boundary.
	CycleNative                       // Add 1 cycle in native mode (E=0).
)

// Opcode contains decoded instruction information for a single opcode byte.
type Opcode struct {
	Instruction *Instruction
	Addressing  AddressingMode
	Timing      byte      // Base cycle count with 8-bit M and X, DL=0, and no page crossing
	Cycles      CycleRule // Conditions that add cycles to Timing
	WidthFlag   WidthFlag // Which flag affects operand width
}

// ReadsMemory returns true if this opcode reads from memory.
func (op Opcode) ReadsMemory(memReadInstructions set.Set[string]) bool {
	switch op.Addressing {
	case ImmediateAddressing, ImpliedAddressing, AccumulatorAddressing, RelativeAddressing, RelativeLongAddressing:
		return false
	}
	return memReadInstructions.Contains(op.Instruction.Name)
}

// WritesMemory returns true if this opcode writes to memory.
func (op Opcode) WritesMemory(memWriteInstructions set.Set[string]) bool {
	switch op.Addressing {
	case ImmediateAddressing, ImpliedAddressing, AccumulatorAddressing, RelativeAddressing, RelativeLongAddressing:
		return false
	}
	return memWriteInstructions.Contains(op.Instruction.Name)
}

// ReadWritesMemory returns true if this opcode reads and writes the same memory operand.
func (op Opcode) ReadWritesMemory(memReadWriteInstructions set.Set[string]) bool {
	switch op.Addressing {
	case ImmediateAddressing, ImpliedAddressing, AccumulatorAddressing, RelativeAddressing, RelativeLongAddressing:
		return false
	}
	return memReadWriteInstructions.Contains(op.Instruction.Name)
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

// Opcodes is the complete 65816 opcode table.
// All 256 entries are defined (the 65816 has no unused opcodes).
// Timing values are base cycles with 8-bit M and X flags; Cycles lists the
// conditions that add cycles. Branch instructions add their taken and
// page-cross cycles in the branch handler.
// Reference: WDC W65C816S Datasheet, Programming the 65816 (Eyes & Lichty).
var Opcodes = [256]Opcode{
	// $00 - $0F
	{Instruction: BrkInst, Addressing: ImmediateAddressing, Timing: 7, Cycles: CycleNative},                       // $00 BRK
	{Instruction: OraInst, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6, Cycles: CycleM | CycleDL}, // $01 ORA (dp,X)
	{Instruction: CopInst, Addressing: ImmediateAddressing, Timing: 7, Cycles: CycleNative},                       // $02 COP
	{Instruction: OraInst, Addressing: StackRelativeAddressing, Timing: 4, Cycles: CycleM},                        // $03 ORA sr,S
	{Instruction: TsbInst, Addressing: DirectPageAddressing, Timing: 5, Cycles: CycleRMW | CycleDL},               // $04 TSB dp
	{Instruction: OraInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $05 ORA dp
	{Instruction: AslInst, Addressing: DirectPageAddressing, Timing: 5, Cycles: CycleRMW | CycleDL},               // $06 ASL dp
	{Instruction: OraInst, Addressing: DirectPageIndirectLongAddressing, Timing: 6, Cycles: CycleM | CycleDL},     // $07 ORA [dp]
	{Instruction: PhpInst, Addressing: ImpliedAddressing, Timing: 3},                                              // $08 PHP
	{Instruction: OraInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleM, WidthFlag: WidthM},         // $09 ORA #imm
	{Instruction: AslInst, Addressing: AccumulatorAddressing, Timing: 2},                                          // $0A ASL A
	{Instruction: PhdInst, Addressing: ImpliedAddressing, Timing: 4},                                              // $0B PHD
	{Instruction: TsbInst, Addressing: AbsoluteAddressing, Timing: 6, Cycles: CycleRMW},                           // $0C TSB abs
	{Instruction: OraInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                             // $0D ORA abs
	{Instruction: AslInst, Addressing: AbsoluteAddressing, Timing: 6, Cycles: CycleRMW},                           // $0E ASL abs
	{Instruction: OraInst, Addressing: AbsoluteLongAddressing, Timing: 5, Cycles: CycleM},                         // $0F ORA al

	// $10 - $1F
	{Instruction: BplInst, Addressing: RelativeAddressing, Timing: 2},                                                          // $10 BPL rel
	{Instruction: OraInst, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, Cycles: CycleM | CycleDL | CycleIndex}, // $11 ORA (dp),Y
	{Instruction: OraInst, Addressing: DirectPageIndirectAddressing, Timing: 5, Cycles: CycleM | CycleDL},                      // $12 ORA (dp)
	{Instruction: OraInst, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7, Cycles: CycleM},                     // $13 ORA (sr,S),Y
	{Instruction: TrbInst, Addressing: DirectPageAddressing, Timing: 5, Cycles: CycleRMW | CycleDL},                            // $14 TRB dp
	{Instruction: OraInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},                      // $15 ORA dp,X
	{Instruction: AslInst, Addressing: DirectPageIndexedXAddressing, Timing: 6, Cycles: CycleRMW | CycleDL},                    // $16 ASL dp,X
	{Instruction: OraInst, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6, Cycles: CycleM | CycleDL},          // $17 ORA [dp],Y
	{Instruction: ClcInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $18 CLC
	{Instruction: OraInst, Addressing: AbsoluteIndexedYAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $19 ORA abs,Y
	{Instruction: IncInst, Addressing: AccumulatorAddressing, Timing: 2},                                                       // $1A INC A
	{Instruction: TcsInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $1B TCS
	{Instruction: TrbInst, Addressing: AbsoluteAddressing, Timing: 6, Cycles: CycleRMW},                                        // $1C TRB abs
	{Instruction: OraInst, Addressing: AbsoluteIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $1D ORA abs,X
	{Instruction: AslInst, Addressing: AbsoluteIndexedXAddressing, Timing: 7, Cycles: CycleRMW},                                // $1E ASL abs,X
	{Instruction: OraInst, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5, Cycles: CycleM},                              // $1F ORA al,X

	// $20 - $2F
	{Instruction: JsrInst, Addressing: AbsoluteAddressing, Timing: 6},                                             // $20 JSR abs
	{Instruction: AndInst, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6, Cycles: CycleM | CycleDL}, // $21 AND (dp,X)
	{Instruction: JslInst, Addressing: AbsoluteLongAddressing, Timing: 8},                                         // $22 JSL al
	{Instruction: AndInst, Addressing: StackRelativeAddressing, Timing: 4, Cycles: CycleM},                        // $23 AND sr,S
	{Instruction: BitInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $24 BIT dp
	{Instruction: AndInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $25 AND dp
	{Instruction: RolInst, Addressing: DirectPageAddressing, Timing: 5, Cycles: CycleRMW | CycleDL},               // $26 ROL dp
	{Instruction: AndInst, Addressing: DirectPageIndirectLongAddressing, Timing: 6, Cycles: CycleM | CycleDL},     // $27 AND [dp]
	{Instruction: PlpInst, Addressing: ImpliedAddressing, Timing: 4},                                              // $28 PLP
	{Instruction: AndInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleM, WidthFlag: WidthM},         // $29 AND #imm
	{Instruction: RolInst, Addressing: AccumulatorAddressing, Timing: 2},                                          // $2A ROL A
	{Instruction: PldInst, Addressing: ImpliedAddressing, Timing: 5},                                              // $2B PLD
	{Instruction: BitInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                             // $2C BIT abs
	{Instruction: AndInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                             // $2D AND abs
	{Instruction: RolInst, Addressing: AbsoluteAddressing, Timing: 6, Cycles: CycleRMW},                           // $2E ROL abs
	{Instruction: AndInst, Addressing: AbsoluteLongAddressing, Timing: 5, Cycles: CycleM},                         // $2F AND al

	// $30 - $3F
	{Instruction: BmiInst, Addressing: RelativeAddressing, Timing: 2},                                                          // $30 BMI rel
	{Instruction: AndInst, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, Cycles: CycleM | CycleDL | CycleIndex}, // $31 AND (dp),Y
	{Instruction: AndInst, Addressing: DirectPageIndirectAddressing, Timing: 5, Cycles: CycleM | CycleDL},                      // $32 AND (dp)
	{Instruction: AndInst, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7, Cycles: CycleM},                     // $33 AND (sr,S),Y
	{Instruction: BitInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},                      // $34 BIT dp,X
	{Instruction: AndInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},                      // $35 AND dp,X
	{Instruction: RolInst, Addressing: DirectPageIndexedXAddressing, Timing: 6, Cycles: CycleRMW | CycleDL},                    // $36 ROL dp,X
	{Instruction: AndInst, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6, Cycles: CycleM | CycleDL},          // $37 AND [dp],Y
	{Instruction: SecInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $38 SEC
	{Instruction: AndInst, Addressing: AbsoluteIndexedYAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $39 AND abs,Y
	{Instruction: DecInst, Addressing: AccumulatorAddressing, Timing: 2},                                                       // $3A DEC A
	{Instruction: TscInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $3B TSC
	{Instruction: BitInst, Addressing: AbsoluteIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $3C BIT abs,X
	{Instruction: AndInst, Addressing: AbsoluteIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $3D AND abs,X
	{Instruction: RolInst, Addressing: AbsoluteIndexedXAddressing, Timing: 7, Cycles: CycleRMW},                                // $3E ROL abs,X
	{Instruction: AndInst, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5, Cycles: CycleM},                              // $3F AND al,X

	// $40 - $4F
	{Instruction: RtiInst, Addressing: ImpliedAddressing, Timing: 6, Cycles: CycleNative},                         // $40 RTI
	{Instruction: EorInst, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6, Cycles: CycleM | CycleDL}, // $41 EOR (dp,X)
	{Instruction: WdmInst, Addressing: ImmediateAddressing, Timing: 2},                                            // $42 WDM
	{Instruction: EorInst, Addressing: StackRelativeAddressing, Timing: 4, Cycles: CycleM},                        // $43 EOR sr,S
	{Instruction: MvpInst, Addressing: BlockMoveAddressing, Timing: 7},                                            // $44 MVP
	{Instruction: EorInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $45 EOR dp
	{Instruction: LsrInst, Addressing: DirectPageAddressing, Timing: 5, Cycles: CycleRMW | CycleDL},               // $46 LSR dp
	{Instruction: EorInst, Addressing: DirectPageIndirectLongAddressing, Timing: 6, Cycles: CycleM | CycleDL},     // $47 EOR [dp]
	{Instruction: PhaInst, Addressing: ImpliedAddressing, Timing: 3, Cycles: CycleM},                              // $48 PHA
	{Instruction: EorInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleM, WidthFlag: WidthM},         // $49 EOR #imm
	{Instruction: LsrInst, Addressing: AccumulatorAddressing, Timing: 2},                                          // $4A LSR A
	{Instruction: PhkInst, Addressing: ImpliedAddressing, Timing: 3},                                              // $4B PHK
	{Instruction: JmpInst, Addressing: AbsoluteAddressing, Timing: 3},                                             // $4C JMP abs
	{Instruction: EorInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                             // $4D EOR abs
	{Instruction: LsrInst, Addressing: AbsoluteAddressing, Timing: 6, Cycles: CycleRMW},                           // $4E LSR abs
	{Instruction: EorInst, Addressing: AbsoluteLongAddressing, Timing: 5, Cycles: CycleM},                         // $4F EOR al

	// $50 - $5F
	{Instruction: BvcInst, Addressing: RelativeAddressing, Timing: 2},                                                          // $50 BVC rel
	{Instruction: EorInst, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, Cycles: CycleM | CycleDL | CycleIndex}, // $51 EOR (dp),Y
	{Instruction: EorInst, Addressing: DirectPageIndirectAddressing, Timing: 5, Cycles: CycleM | CycleDL},                      // $52 EOR (dp)
	{Instruction: EorInst, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7, Cycles: CycleM},                     // $53 EOR (sr,S),Y
	{Instruction: MvnInst, Addressing: BlockMoveAddressing, Timing: 7},                                                         // $54 MVN
	{Instruction: EorInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},                      // $55 EOR dp,X
	{Instruction: LsrInst, Addressing: DirectPageIndexedXAddressing, Timing: 6, Cycles: CycleRMW | CycleDL},                    // $56 LSR dp,X
	{Instruction: EorInst, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6, Cycles: CycleM | CycleDL},          // $57 EOR [dp],Y
	{Instruction: CliInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $58 CLI
	{Instruction: EorInst, Addressing: AbsoluteIndexedYAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $59 EOR abs,Y
	{Instruction: PhyInst, Addressing: ImpliedAddressing, Timing: 3, Cycles: CycleX},                                           // $5A PHY
	{Instruction: TcdInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $5B TCD
	{Instruction: JmlInst, Addressing: AbsoluteLongAddressing, Timing: 4},                                                      // $5C JML al
	{Instruction: EorInst, Addressing: AbsoluteIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $5D EOR abs,X
	{Instruction: LsrInst, Addressing: AbsoluteIndexedXAddressing, Timing: 7, Cycles: CycleRMW},                                // $5E LSR abs,X
	{Instruction: EorInst, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5, Cycles: CycleM},                              // $5F EOR al,X

	// $60 - $6F
	{Instruction: RtsInst, Addressing: ImpliedAddressing, Timing: 6},                                              // $60 RTS
	{Instruction: AdcInst, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6, Cycles: CycleM | CycleDL}, // $61 ADC (dp,X)
	{Instruction: PerInst, Addressing: RelativeLongAddressing, Timing: 6},                                         // $62 PER rl
	{Instruction: AdcInst, Addressing: StackRelativeAddressing, Timing: 4, Cycles: CycleM},                        // $63 ADC sr,S
	{Instruction: StzInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $64 STZ dp
	{Instruction: AdcInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $65 ADC dp
	{Instruction: RorInst, Addressing: DirectPageAddressing, Timing: 5, Cycles: CycleRMW | CycleDL},               // $66 ROR dp
	{Instruction: AdcInst, Addressing: DirectPageIndirectLongAddressing, Timing: 6, Cycles: CycleM | CycleDL},     // $67 ADC [dp]
	{Instruction: PlaInst, Addressing: ImpliedAddressing, Timing: 4, Cycles: CycleM},                              // $68 PLA
	{Instruction: AdcInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleM, WidthFlag: WidthM},         // $69 ADC #imm
	{Instruction: RorInst, Addressing: AccumulatorAddressing, Timing: 2},                                          // $6A ROR A
	{Instruction: RtlInst, Addressing: ImpliedAddressing, Timing: 6},                                              // $6B RTL
	{Instruction: JmpInst, Addressing: AbsoluteIndirectAddressing, Timing: 5},                                     // $6C JMP (abs)
	{Instruction: AdcInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                             // $6D ADC abs
	{Instruction: RorInst, Addressing: AbsoluteAddressing, Timing: 6, Cycles: CycleRMW},                           // $6E ROR abs
	{Instruction: AdcInst, Addressing: AbsoluteLongAddressing, Timing: 5, Cycles: CycleM},                         // $6F ADC al

	// $70 - $7F
	{Instruction: BvsInst, Addressing: RelativeAddressing, Timing: 2},                                                          // $70 BVS rel
	{Instruction: AdcInst, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, Cycles: CycleM | CycleDL | CycleIndex}, // $71 ADC (dp),Y
	{Instruction: AdcInst, Addressing: DirectPageIndirectAddressing, Timing: 5, Cycles: CycleM | CycleDL},                      // $72 ADC (dp)
	{Instruction: AdcInst, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7, Cycles: CycleM},                     // $73 ADC (sr,S),Y
	{Instruction: StzInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},                      // $74 STZ dp,X
	{Instruction: AdcInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},                      // $75 ADC dp,X
	{Instruction: RorInst, Addressing: DirectPageIndexedXAddressing, Timing: 6, Cycles: CycleRMW | CycleDL},                    // $76 ROR dp,X
	{Instruction: AdcInst, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6, Cycles: CycleM | CycleDL},          // $77 ADC [dp],Y
	{Instruction: SeiInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $78 SEI
	{Instruction: AdcInst, Addressing: AbsoluteIndexedYAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $79 ADC abs,Y
	{Instruction: PlyInst, Addressing: ImpliedAddressing, Timing: 4, Cycles: CycleX},                                           // $7A PLY
	{Instruction: TdcInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $7B TDC
	{Instruction: JmpInst, Addressing: AbsoluteIndexedXIndirectAddressing, Timing: 6},                                          // $7C JMP (abs,X)
	{Instruction: AdcInst, Addressing: AbsoluteIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $7D ADC abs,X
	{Instruction: RorInst, Addressing: AbsoluteIndexedXAddressing, Timing: 7, Cycles: CycleRMW},                                // $7E ROR abs,X
	{Instruction: AdcInst, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5, Cycles: CycleM},                              // $7F ADC al,X

	// $80 - $8F
	{Instruction: BraInst, Addressing: RelativeAddressing, Timing: 2},                                             // $80 BRA rel
	{Instruction: StaInst, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6, Cycles: CycleM | CycleDL}, // $81 STA (dp,X)
	{Instruction: BrlInst, Addressing: RelativeLongAddressing, Timing: 4},                                         // $82 BRL rl
	{Instruction: StaInst, Addressing: StackRelativeAddressing, Timing: 4, Cycles: CycleM},                        // $83 STA sr,S
	{Instruction: StyInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleX | CycleDL},                 // $84 STY dp
	{Instruction: StaInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $85 STA dp
	{Instruction: StxInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleX | CycleDL},                 // $86 STX dp
	{Instruction: StaInst, Addressing: DirectPageIndirectLongAddressing, Timing: 6, Cycles: CycleM | CycleDL},     // $87 STA [dp]
	{Instruction: DeyInst, Addressing: ImpliedAddressing, Timing: 2},                                              // $88 DEY
	{Instruction: BitInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleM, WidthFlag: WidthM},         // $89 BIT #imm
	{Instruction: TxaInst, Addressing: ImpliedAddressing, Timing: 2},                                              // $8A TXA
	{Instruction: PhbInst, Addressing: ImpliedAddressing, Timing: 3},                                              // $8B PHB
	{Instruction: StyInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleX},                             // $8C STY abs
	{Instruction: StaInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                             // $8D STA abs
	{Instruction: StxInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleX},                             // $8E STX abs
	{Instruction: StaInst, Addressing: AbsoluteLongAddressing, Timing: 5, Cycles: CycleM},                         // $8F STA al

	// $90 - $9F
	{Instruction: BccInst, Addressing: RelativeAddressing, Timing: 2},                                                 // $90 BCC rel
	{Instruction: StaInst, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 6, Cycles: CycleM | CycleDL},     // $91 STA (dp),Y
	{Instruction: StaInst, Addressing: DirectPageIndirectAddressing, Timing: 5, Cycles: CycleM | CycleDL},             // $92 STA (dp)
	{Instruction: StaInst, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7, Cycles: CycleM},            // $93 STA (sr,S),Y
	{Instruction: StyInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleX | CycleDL},             // $94 STY dp,X
	{Instruction: StaInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},             // $95 STA dp,X
	{Instruction: StxInst, Addressing: DirectPageIndexedYAddressing, Timing: 4, Cycles: CycleX | CycleDL},             // $96 STX dp,Y
	{Instruction: StaInst, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6, Cycles: CycleM | CycleDL}, // $97 STA [dp],Y
	{Instruction: TyaInst, Addressing: ImpliedAddressing, Timing: 2},                                                  // $98 TYA
	{Instruction: StaInst, Addressing: AbsoluteIndexedYAddressing, Timing: 5, Cycles: CycleM},                         // $99 STA abs,Y
	{Instruction: TxsInst, Addressing: ImpliedAddressing, Timing: 2},                                                  // $9A TXS
	{Instruction: TxyInst, Addressing: ImpliedAddressing, Timing: 2},                                                  // $9B TXY
	{Instruction: StzInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                                 // $9C STZ abs
	{Instruction: StaInst, Addressing: AbsoluteIndexedXAddressing, Timing: 5, Cycles: CycleM},                         // $9D STA abs,X
	{Instruction: StzInst, Addressing: AbsoluteIndexedXAddressing, Timing: 5, Cycles: CycleM},                         // $9E STZ abs,X
	{Instruction: StaInst, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5, Cycles: CycleM},                     // $9F STA al,X

	// $A0 - $AF
	{Instruction: LdyInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleX, WidthFlag: WidthX},         // $A0 LDY #imm
	{Instruction: LdaInst, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6, Cycles: CycleM | CycleDL}, // $A1 LDA (dp,X)
	{Instruction: LdxInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleX, WidthFlag: WidthX},         // $A2 LDX #imm
	{Instruction: LdaInst, Addressing: StackRelativeAddressing, Timing: 4, Cycles: CycleM},                        // $A3 LDA sr,S
	{Instruction: LdyInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleX | CycleDL},                 // $A4 LDY dp
	{Instruction: LdaInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $A5 LDA dp
	{Instruction: LdxInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleX | CycleDL},                 // $A6 LDX dp
	{Instruction: LdaInst, Addressing: DirectPageIndirectLongAddressing, Timing: 6, Cycles: CycleM | CycleDL},     // $A7 LDA [dp]
	{Instruction: TayInst, Addressing: ImpliedAddressing, Timing: 2},                                              // $A8 TAY
	{Instruction: LdaInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleM, WidthFlag: WidthM},         // $A9 LDA #imm
	{Instruction: TaxInst, Addressing: ImpliedAddressing, Timing: 2},                                              // $AA TAX
	{Instruction: PlbInst, Addressing: ImpliedAddressing, Timing: 4},                                              // $AB PLB
	{Instruction: LdyInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleX},                             // $AC LDY abs
	{Instruction: LdaInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                             // $AD LDA abs
	{Instruction: LdxInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleX},                             // $AE LDX abs
	{Instruction: LdaInst, Addressing: AbsoluteLongAddressing, Timing: 5, Cycles: CycleM},                         // $AF LDA al

	// $B0 - $BF
	{Instruction: BcsInst, Addressing: RelativeAddressing, Timing: 2},                                                          // $B0 BCS rel
	{Instruction: LdaInst, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, Cycles: CycleM | CycleDL | CycleIndex}, // $B1 LDA (dp),Y
	{Instruction: LdaInst, Addressing: DirectPageIndirectAddressing, Timing: 5, Cycles: CycleM | CycleDL},                      // $B2 LDA (dp)
	{Instruction: LdaInst, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7, Cycles: CycleM},                     // $B3 LDA (sr,S),Y
	{Instruction: LdyInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleX | CycleDL},                      // $B4 LDY dp,X
	{Instruction: LdaInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},                      // $B5 LDA dp,X
	{Instruction: LdxInst, Addressing: DirectPageIndexedYAddressing, Timing: 4, Cycles: CycleX | CycleDL},                      // $B6 LDX dp,Y
	{Instruction: LdaInst, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6, Cycles: CycleM | CycleDL},          // $B7 LDA [dp],Y
	{Instruction: ClvInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $B8 CLV
	{Instruction: LdaInst, Addressing: AbsoluteIndexedYAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $B9 LDA abs,Y
	{Instruction: TsxInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $BA TSX
	{Instruction: TyxInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $BB TYX
	{Instruction: LdyInst, Addressing: AbsoluteIndexedXAddressing, Timing: 4, Cycles: CycleX | CycleIndex},                     // $BC LDY abs,X
	{Instruction: LdaInst, Addressing: AbsoluteIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $BD LDA abs,X
	{Instruction: LdxInst, Addressing: AbsoluteIndexedYAddressing, Timing: 4, Cycles: CycleX | CycleIndex},                     // $BE LDX abs,Y
	{Instruction: LdaInst, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5, Cycles: CycleM},                              // $BF LDA al,X

	// $C0 - $CF
	{Instruction: CpyInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleX, WidthFlag: WidthX},         // $C0 CPY #imm
	{Instruction: CmpInst, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6, Cycles: CycleM | CycleDL}, // $C1 CMP (dp,X)
	{Instruction: RepInst, Addressing: ImmediateAddressing, Timing: 3},                                            // $C2 REP #imm
	{Instruction: CmpInst, Addressing: StackRelativeAddressing, Timing: 4, Cycles: CycleM},                        // $C3 CMP sr,S
	{Instruction: CpyInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleX | CycleDL},                 // $C4 CPY dp
	{Instruction: CmpInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $C5 CMP dp
	{Instruction: DecInst, Addressing: DirectPageAddressing, Timing: 5, Cycles: CycleRMW | CycleDL},               // $C6 DEC dp
	{Instruction: CmpInst, Addressing: DirectPageIndirectLongAddressing, Timing: 6, Cycles: CycleM | CycleDL},     // $C7 CMP [dp]
	{Instruction: InyInst, Addressing: ImpliedAddressing, Timing: 2},                                              // $C8 INY
	{Instruction: CmpInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleM, WidthFlag: WidthM},         // $C9 CMP #imm
	{Instruction: DexInst, Addressing: ImpliedAddressing, Timing: 2},                                              // $CA DEX
	{Instruction: WaiInst, Addressing: ImpliedAddressing, Timing: 3},                                              // $CB WAI
	{Instruction: CpyInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleX},                             // $CC CPY abs
	{Instruction: CmpInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                             // $CD CMP abs
	{Instruction: DecInst, Addressing: AbsoluteAddressing, Timing: 6, Cycles: CycleRMW},                           // $CE DEC abs
	{Instruction: CmpInst, Addressing: AbsoluteLongAddressing, Timing: 5, Cycles: CycleM},                         // $CF CMP al

	// $D0 - $DF
	{Instruction: BneInst, Addressing: RelativeAddressing, Timing: 2},                                                          // $D0 BNE rel
	{Instruction: CmpInst, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, Cycles: CycleM | CycleDL | CycleIndex}, // $D1 CMP (dp),Y
	{Instruction: CmpInst, Addressing: DirectPageIndirectAddressing, Timing: 5, Cycles: CycleM | CycleDL},                      // $D2 CMP (dp)
	{Instruction: CmpInst, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7, Cycles: CycleM},                     // $D3 CMP (sr,S),Y
	{Instruction: PeiInst, Addressing: DirectPageIndirectAddressing, Timing: 6, Cycles: CycleDL},                               // $D4 PEI (dp)
	{Instruction: CmpInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},                      // $D5 CMP dp,X
	{Instruction: DecInst, Addressing: DirectPageIndexedXAddressing, Timing: 6, Cycles: CycleRMW | CycleDL},                    // $D6 DEC dp,X
	{Instruction: CmpInst, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6, Cycles: CycleM | CycleDL},          // $D7 CMP [dp],Y
	{Instruction: CldInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $D8 CLD
	{Instruction: CmpInst, Addressing: AbsoluteIndexedYAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $D9 CMP abs,Y
	{Instruction: PhxInst, Addressing: ImpliedAddressing, Timing: 3, Cycles: CycleX},                                           // $DA PHX
	{Instruction: StpInst, Addressing: ImpliedAddressing, Timing: 3},                                                           // $DB STP
	{Instruction: JmlInst, Addressing: AbsoluteIndirectLongAddressing, Timing: 6},                                              // $DC JML [abs]
	{Instruction: CmpInst, Addressing: AbsoluteIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $DD CMP abs,X
	{Instruction: DecInst, Addressing: AbsoluteIndexedXAddressing, Timing: 7, Cycles: CycleRMW},                                // $DE DEC abs,X
	{Instruction: CmpInst, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5, Cycles: CycleM},                              // $DF CMP al,X

	// $E0 - $EF
	{Instruction: CpxInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleX, WidthFlag: WidthX},         // $E0 CPX #imm
	{Instruction: SbcInst, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6, Cycles: CycleM | CycleDL}, // $E1 SBC (dp,X)
	{Instruction: SepInst, Addressing: ImmediateAddressing, Timing: 3},                                            // $E2 SEP #imm
	{Instruction: SbcInst, Addressing: StackRelativeAddressing, Timing: 4, Cycles: CycleM},                        // $E3 SBC sr,S
	{Instruction: CpxInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleX | CycleDL},                 // $E4 CPX dp
	{Instruction: SbcInst, Addressing: DirectPageAddressing, Timing: 3, Cycles: CycleM | CycleDL},                 // $E5 SBC dp
	{Instruction: IncInst, Addressing: DirectPageAddressing, Timing: 5, Cycles: CycleRMW | CycleDL},               // $E6 INC dp
	{Instruction: SbcInst, Addressing: DirectPageIndirectLongAddressing, Timing: 6, Cycles: CycleM | CycleDL},     // $E7 SBC [dp]
	{Instruction: InxInst, Addressing: ImpliedAddressing, Timing: 2},                                              // $E8 INX
	{Instruction: SbcInst, Addressing: ImmediateAddressing, Timing: 2, Cycles: CycleM, WidthFlag: WidthM},         // $E9 SBC #imm
	{Instruction: NopInst, Addressing: ImpliedAddressing, Timing: 2},                                              // $EA NOP
	{Instruction: XbaInst, Addressing: ImpliedAddressing, Timing: 3},                                              // $EB XBA
	{Instruction: CpxInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleX},                             // $EC CPX abs
	{Instruction: SbcInst, Addressing: AbsoluteAddressing, Timing: 4, Cycles: CycleM},                             // $ED SBC abs
	{Instruction: IncInst, Addressing: AbsoluteAddressing, Timing: 6, Cycles: CycleRMW},                           // $EE INC abs
	{Instruction: SbcInst, Addressing: AbsoluteLongAddressing, Timing: 5, Cycles: CycleM},                         // $EF SBC al

	// $F0 - $FF
	{Instruction: BeqInst, Addressing: RelativeAddressing, Timing: 2},                                                          // $F0 BEQ rel
	{Instruction: SbcInst, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, Cycles: CycleM | CycleDL | CycleIndex}, // $F1 SBC (dp),Y
	{Instruction: SbcInst, Addressing: DirectPageIndirectAddressing, Timing: 5, Cycles: CycleM | CycleDL},                      // $F2 SBC (dp)
	{Instruction: SbcInst, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7, Cycles: CycleM},                     // $F3 SBC (sr,S),Y
	{Instruction: PeaInst, Addressing: AbsoluteAddressing, Timing: 5},                                                          // $F4 PEA abs
	{Instruction: SbcInst, Addressing: DirectPageIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleDL},                      // $F5 SBC dp,X
	{Instruction: IncInst, Addressing: DirectPageIndexedXAddressing, Timing: 6, Cycles: CycleRMW | CycleDL},                    // $F6 INC dp,X
	{Instruction: SbcInst, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6, Cycles: CycleM | CycleDL},          // $F7 SBC [dp],Y
	{Instruction: SedInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $F8 SED
	{Instruction: SbcInst, Addressing: AbsoluteIndexedYAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $F9 SBC abs,Y
	{Instruction: PlxInst, Addressing: ImpliedAddressing, Timing: 4, Cycles: CycleX},                                           // $FA PLX
	{Instruction: XceInst, Addressing: ImpliedAddressing, Timing: 2},                                                           // $FB XCE
	{Instruction: JsrInst, Addressing: AbsoluteIndexedXIndirectAddressing, Timing: 8},                                          // $FC JSR (abs,X)
	{Instruction: SbcInst, Addressing: AbsoluteIndexedXAddressing, Timing: 4, Cycles: CycleM | CycleIndex},                     // $FD SBC abs,X
	{Instruction: IncInst, Addressing: AbsoluteIndexedXAddressing, Timing: 7, Cycles: CycleRMW},                                // $FE INC abs,X
	{Instruction: SbcInst, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5, Cycles: CycleM},                              // $FF SBC al,X
}
