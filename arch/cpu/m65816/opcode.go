package m65816

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

// Opcode contains decoded instruction information for a single opcode byte.
type Opcode struct {
	Instruction    *Instruction
	Addressing     AddressingMode
	Timing         byte      // Base cycle count
	PageCrossCycle bool      // Extra cycle on page boundary crossing
	WidthFlag      WidthFlag // Which flag affects operand width
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

// IsBranching returns true if this opcode is a branching instruction.
func (op Opcode) IsBranching(branchInstructions set.Set[string]) bool {
	return branchInstructions.Contains(op.Instruction.Name)
}

// Opcodes is the complete 65816 opcode table.
// All 256 entries are defined (the 65816 has no unused opcodes).
// Timing values are base cycles assuming 8-bit M and X flags.
// Reference: WDC W65C816S Datasheet, Programming the 65816 (Eyes & Lichty).
var Opcodes = [256]Opcode{
	// $00 - $0F
	{Instruction: Brk, Addressing: ImmediateAddressing, Timing: 7},                    // $00 BRK
	{Instruction: Ora, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6},   // $01 ORA (dp,X)
	{Instruction: Cop, Addressing: ImmediateAddressing, Timing: 7},                    // $02 COP
	{Instruction: Ora, Addressing: StackRelativeAddressing, Timing: 4},                // $03 ORA sr,S
	{Instruction: Tsb, Addressing: DirectPageAddressing, Timing: 5},                   // $04 TSB dp
	{Instruction: Ora, Addressing: DirectPageAddressing, Timing: 3},                   // $05 ORA dp
	{Instruction: Asl, Addressing: DirectPageAddressing, Timing: 5},                   // $06 ASL dp
	{Instruction: Ora, Addressing: DirectPageIndirectLongAddressing, Timing: 6},       // $07 ORA [dp]
	{Instruction: Php, Addressing: ImpliedAddressing, Timing: 3},                      // $08 PHP
	{Instruction: Ora, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthM}, // $09 ORA #imm
	{Instruction: Asl, Addressing: AccumulatorAddressing, Timing: 2},                  // $0A ASL A
	{Instruction: Phd, Addressing: ImpliedAddressing, Timing: 4},                      // $0B PHD
	{Instruction: Tsb, Addressing: AbsoluteAddressing, Timing: 6},                     // $0C TSB abs
	{Instruction: Ora, Addressing: AbsoluteAddressing, Timing: 4},                     // $0D ORA abs
	{Instruction: Asl, Addressing: AbsoluteAddressing, Timing: 6},                     // $0E ASL abs
	{Instruction: Ora, Addressing: AbsoluteLongAddressing, Timing: 5},                 // $0F ORA al

	// $10 - $1F
	{Instruction: Bpl, Addressing: RelativeAddressing, Timing: 2, PageCrossCycle: true},                   // $10 BPL rel
	{Instruction: Ora, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, PageCrossCycle: true}, // $11 ORA (dp),Y
	{Instruction: Ora, Addressing: DirectPageIndirectAddressing, Timing: 5},                               // $12 ORA (dp)
	{Instruction: Ora, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7},                    // $13 ORA (sr,S),Y
	{Instruction: Trb, Addressing: DirectPageAddressing, Timing: 5},                                       // $14 TRB dp
	{Instruction: Ora, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $15 ORA dp,X
	{Instruction: Asl, Addressing: DirectPageIndexedXAddressing, Timing: 6},                               // $16 ASL dp,X
	{Instruction: Ora, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6},                   // $17 ORA [dp],Y
	{Instruction: Clc, Addressing: ImpliedAddressing, Timing: 2},                                          // $18 CLC
	{Instruction: Ora, Addressing: AbsoluteIndexedYAddressing, Timing: 4, PageCrossCycle: true},           // $19 ORA abs,Y
	{Instruction: Inc, Addressing: AccumulatorAddressing, Timing: 2},                                      // $1A INC A
	{Instruction: Tcs, Addressing: ImpliedAddressing, Timing: 2},                                          // $1B TCS
	{Instruction: Trb, Addressing: AbsoluteAddressing, Timing: 6},                                         // $1C TRB abs
	{Instruction: Ora, Addressing: AbsoluteIndexedXAddressing, Timing: 4, PageCrossCycle: true},           // $1D ORA abs,X
	{Instruction: Asl, Addressing: AbsoluteIndexedXAddressing, Timing: 7},                                 // $1E ASL abs,X
	{Instruction: Ora, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5},                             // $1F ORA al,X

	// $20 - $2F
	{Instruction: Jsr, Addressing: AbsoluteAddressing, Timing: 6},                     // $20 JSR abs
	{Instruction: And, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6},   // $21 AND (dp,X)
	{Instruction: Jsl, Addressing: AbsoluteLongAddressing, Timing: 8},                 // $22 JSL al
	{Instruction: And, Addressing: StackRelativeAddressing, Timing: 4},                // $23 AND sr,S
	{Instruction: Bit, Addressing: DirectPageAddressing, Timing: 3},                   // $24 BIT dp
	{Instruction: And, Addressing: DirectPageAddressing, Timing: 3},                   // $25 AND dp
	{Instruction: Rol, Addressing: DirectPageAddressing, Timing: 5},                   // $26 ROL dp
	{Instruction: And, Addressing: DirectPageIndirectLongAddressing, Timing: 6},       // $27 AND [dp]
	{Instruction: Plp, Addressing: ImpliedAddressing, Timing: 4},                      // $28 PLP
	{Instruction: And, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthM}, // $29 AND #imm
	{Instruction: Rol, Addressing: AccumulatorAddressing, Timing: 2},                  // $2A ROL A
	{Instruction: Pld, Addressing: ImpliedAddressing, Timing: 5},                      // $2B PLD
	{Instruction: Bit, Addressing: AbsoluteAddressing, Timing: 4},                     // $2C BIT abs
	{Instruction: And, Addressing: AbsoluteAddressing, Timing: 4},                     // $2D AND abs
	{Instruction: Rol, Addressing: AbsoluteAddressing, Timing: 6},                     // $2E ROL abs
	{Instruction: And, Addressing: AbsoluteLongAddressing, Timing: 5},                 // $2F AND al

	// $30 - $3F
	{Instruction: Bmi, Addressing: RelativeAddressing, Timing: 2, PageCrossCycle: true},                   // $30 BMI rel
	{Instruction: And, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, PageCrossCycle: true}, // $31 AND (dp),Y
	{Instruction: And, Addressing: DirectPageIndirectAddressing, Timing: 5},                               // $32 AND (dp)
	{Instruction: And, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7},                    // $33 AND (sr,S),Y
	{Instruction: Bit, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $34 BIT dp,X
	{Instruction: And, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $35 AND dp,X
	{Instruction: Rol, Addressing: DirectPageIndexedXAddressing, Timing: 6},                               // $36 ROL dp,X
	{Instruction: And, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6},                   // $37 AND [dp],Y
	{Instruction: Sec, Addressing: ImpliedAddressing, Timing: 2},                                          // $38 SEC
	{Instruction: And, Addressing: AbsoluteIndexedYAddressing, Timing: 4, PageCrossCycle: true},           // $39 AND abs,Y
	{Instruction: Dec, Addressing: AccumulatorAddressing, Timing: 2},                                      // $3A DEC A
	{Instruction: Tsc, Addressing: ImpliedAddressing, Timing: 2},                                          // $3B TSC
	{Instruction: Bit, Addressing: AbsoluteIndexedXAddressing, Timing: 4, PageCrossCycle: true},           // $3C BIT abs,X
	{Instruction: And, Addressing: AbsoluteIndexedXAddressing, Timing: 4, PageCrossCycle: true},           // $3D AND abs,X
	{Instruction: Rol, Addressing: AbsoluteIndexedXAddressing, Timing: 7},                                 // $3E ROL abs,X
	{Instruction: And, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5},                             // $3F AND al,X

	// $40 - $4F
	{Instruction: Rti, Addressing: ImpliedAddressing, Timing: 6},                      // $40 RTI
	{Instruction: Eor, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6},   // $41 EOR (dp,X)
	{Instruction: Wdm, Addressing: ImmediateAddressing, Timing: 2},                    // $42 WDM
	{Instruction: Eor, Addressing: StackRelativeAddressing, Timing: 4},                // $43 EOR sr,S
	{Instruction: Mvp, Addressing: BlockMoveAddressing, Timing: 7},                    // $44 MVP
	{Instruction: Eor, Addressing: DirectPageAddressing, Timing: 3},                   // $45 EOR dp
	{Instruction: Lsr, Addressing: DirectPageAddressing, Timing: 5},                   // $46 LSR dp
	{Instruction: Eor, Addressing: DirectPageIndirectLongAddressing, Timing: 6},       // $47 EOR [dp]
	{Instruction: Pha, Addressing: ImpliedAddressing, Timing: 3},                      // $48 PHA
	{Instruction: Eor, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthM}, // $49 EOR #imm
	{Instruction: Lsr, Addressing: AccumulatorAddressing, Timing: 2},                  // $4A LSR A
	{Instruction: Phk, Addressing: ImpliedAddressing, Timing: 3},                      // $4B PHK
	{Instruction: Jmp, Addressing: AbsoluteAddressing, Timing: 3},                     // $4C JMP abs
	{Instruction: Eor, Addressing: AbsoluteAddressing, Timing: 4},                     // $4D EOR abs
	{Instruction: Lsr, Addressing: AbsoluteAddressing, Timing: 6},                     // $4E LSR abs
	{Instruction: Eor, Addressing: AbsoluteLongAddressing, Timing: 5},                 // $4F EOR al

	// $50 - $5F
	{Instruction: Bvc, Addressing: RelativeAddressing, Timing: 2, PageCrossCycle: true},                   // $50 BVC rel
	{Instruction: Eor, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, PageCrossCycle: true}, // $51 EOR (dp),Y
	{Instruction: Eor, Addressing: DirectPageIndirectAddressing, Timing: 5},                               // $52 EOR (dp)
	{Instruction: Eor, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7},                    // $53 EOR (sr,S),Y
	{Instruction: Mvn, Addressing: BlockMoveAddressing, Timing: 7},                                        // $54 MVN
	{Instruction: Eor, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $55 EOR dp,X
	{Instruction: Lsr, Addressing: DirectPageIndexedXAddressing, Timing: 6},                               // $56 LSR dp,X
	{Instruction: Eor, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6},                   // $57 EOR [dp],Y
	{Instruction: Cli, Addressing: ImpliedAddressing, Timing: 2},                                          // $58 CLI
	{Instruction: Eor, Addressing: AbsoluteIndexedYAddressing, Timing: 4, PageCrossCycle: true},           // $59 EOR abs,Y
	{Instruction: Phy, Addressing: ImpliedAddressing, Timing: 3},                                          // $5A PHY
	{Instruction: Tcd, Addressing: ImpliedAddressing, Timing: 2},                                          // $5B TCD
	{Instruction: Jml, Addressing: AbsoluteLongAddressing, Timing: 4},                                     // $5C JML al
	{Instruction: Eor, Addressing: AbsoluteIndexedXAddressing, Timing: 4, PageCrossCycle: true},           // $5D EOR abs,X
	{Instruction: Lsr, Addressing: AbsoluteIndexedXAddressing, Timing: 7},                                 // $5E LSR abs,X
	{Instruction: Eor, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5},                             // $5F EOR al,X

	// $60 - $6F
	{Instruction: Rts, Addressing: ImpliedAddressing, Timing: 6},                      // $60 RTS
	{Instruction: Adc, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6},   // $61 ADC (dp,X)
	{Instruction: Per, Addressing: RelativeLongAddressing, Timing: 6},                 // $62 PER rl
	{Instruction: Adc, Addressing: StackRelativeAddressing, Timing: 4},                // $63 ADC sr,S
	{Instruction: Stz, Addressing: DirectPageAddressing, Timing: 3},                   // $64 STZ dp
	{Instruction: Adc, Addressing: DirectPageAddressing, Timing: 3},                   // $65 ADC dp
	{Instruction: Ror, Addressing: DirectPageAddressing, Timing: 5},                   // $66 ROR dp
	{Instruction: Adc, Addressing: DirectPageIndirectLongAddressing, Timing: 6},       // $67 ADC [dp]
	{Instruction: Pla, Addressing: ImpliedAddressing, Timing: 4},                      // $68 PLA
	{Instruction: Adc, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthM}, // $69 ADC #imm
	{Instruction: Ror, Addressing: AccumulatorAddressing, Timing: 2},                  // $6A ROR A
	{Instruction: Rtl, Addressing: ImpliedAddressing, Timing: 6},                      // $6B RTL
	{Instruction: Jmp, Addressing: AbsoluteIndirectAddressing, Timing: 5},             // $6C JMP (abs)
	{Instruction: Adc, Addressing: AbsoluteAddressing, Timing: 4},                     // $6D ADC abs
	{Instruction: Ror, Addressing: AbsoluteAddressing, Timing: 6},                     // $6E ROR abs
	{Instruction: Adc, Addressing: AbsoluteLongAddressing, Timing: 5},                 // $6F ADC al

	// $70 - $7F
	{Instruction: Bvs, Addressing: RelativeAddressing, Timing: 2, PageCrossCycle: true},                   // $70 BVS rel
	{Instruction: Adc, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, PageCrossCycle: true}, // $71 ADC (dp),Y
	{Instruction: Adc, Addressing: DirectPageIndirectAddressing, Timing: 5},                               // $72 ADC (dp)
	{Instruction: Adc, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7},                    // $73 ADC (sr,S),Y
	{Instruction: Stz, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $74 STZ dp,X
	{Instruction: Adc, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $75 ADC dp,X
	{Instruction: Ror, Addressing: DirectPageIndexedXAddressing, Timing: 6},                               // $76 ROR dp,X
	{Instruction: Adc, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6},                   // $77 ADC [dp],Y
	{Instruction: Sei, Addressing: ImpliedAddressing, Timing: 2},                                          // $78 SEI
	{Instruction: Adc, Addressing: AbsoluteIndexedYAddressing, Timing: 4, PageCrossCycle: true},           // $79 ADC abs,Y
	{Instruction: Ply, Addressing: ImpliedAddressing, Timing: 4},                                          // $7A PLY
	{Instruction: Tdc, Addressing: ImpliedAddressing, Timing: 2},                                          // $7B TDC
	{Instruction: Jmp, Addressing: AbsoluteIndexedXIndirectAddressing, Timing: 6},                         // $7C JMP (abs,X)
	{Instruction: Adc, Addressing: AbsoluteIndexedXAddressing, Timing: 4, PageCrossCycle: true},           // $7D ADC abs,X
	{Instruction: Ror, Addressing: AbsoluteIndexedXAddressing, Timing: 7},                                 // $7E ROR abs,X
	{Instruction: Adc, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5},                             // $7F ADC al,X

	// $80 - $8F
	{Instruction: Bra, Addressing: RelativeAddressing, Timing: 2, PageCrossCycle: true}, // $80 BRA rel
	{Instruction: Sta, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6},   // $81 STA (dp,X)
	{Instruction: Brl, Addressing: RelativeLongAddressing, Timing: 4},                 // $82 BRL rl
	{Instruction: Sta, Addressing: StackRelativeAddressing, Timing: 4},                // $83 STA sr,S
	{Instruction: Sty, Addressing: DirectPageAddressing, Timing: 3},                   // $84 STY dp
	{Instruction: Sta, Addressing: DirectPageAddressing, Timing: 3},                   // $85 STA dp
	{Instruction: Stx, Addressing: DirectPageAddressing, Timing: 3},                   // $86 STX dp
	{Instruction: Sta, Addressing: DirectPageIndirectLongAddressing, Timing: 6},       // $87 STA [dp]
	{Instruction: Dey, Addressing: ImpliedAddressing, Timing: 2},                      // $88 DEY
	{Instruction: Bit, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthM}, // $89 BIT #imm
	{Instruction: Txa, Addressing: ImpliedAddressing, Timing: 2},                      // $8A TXA
	{Instruction: Phb, Addressing: ImpliedAddressing, Timing: 3},                      // $8B PHB
	{Instruction: Sty, Addressing: AbsoluteAddressing, Timing: 4},                     // $8C STY abs
	{Instruction: Sta, Addressing: AbsoluteAddressing, Timing: 4},                     // $8D STA abs
	{Instruction: Stx, Addressing: AbsoluteAddressing, Timing: 4},                     // $8E STX abs
	{Instruction: Sta, Addressing: AbsoluteLongAddressing, Timing: 5},                 // $8F STA al

	// $90 - $9F
	{Instruction: Bcc, Addressing: RelativeAddressing, Timing: 2, PageCrossCycle: true}, // $90 BCC rel
	{Instruction: Sta, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 6},     // $91 STA (dp),Y
	{Instruction: Sta, Addressing: DirectPageIndirectAddressing, Timing: 5},             // $92 STA (dp)
	{Instruction: Sta, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7},  // $93 STA (sr,S),Y
	{Instruction: Sty, Addressing: DirectPageIndexedXAddressing, Timing: 4},             // $94 STY dp,X
	{Instruction: Sta, Addressing: DirectPageIndexedXAddressing, Timing: 4},             // $95 STA dp,X
	{Instruction: Stx, Addressing: DirectPageIndexedYAddressing, Timing: 4},             // $96 STX dp,Y
	{Instruction: Sta, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6}, // $97 STA [dp],Y
	{Instruction: Tya, Addressing: ImpliedAddressing, Timing: 2},                        // $98 TYA
	{Instruction: Sta, Addressing: AbsoluteIndexedYAddressing, Timing: 5},               // $99 STA abs,Y
	{Instruction: Txs, Addressing: ImpliedAddressing, Timing: 2},                        // $9A TXS
	{Instruction: Txy, Addressing: ImpliedAddressing, Timing: 2},                        // $9B TXY
	{Instruction: Stz, Addressing: AbsoluteAddressing, Timing: 4},                       // $9C STZ abs
	{Instruction: Sta, Addressing: AbsoluteIndexedXAddressing, Timing: 5},               // $9D STA abs,X
	{Instruction: Stz, Addressing: AbsoluteIndexedXAddressing, Timing: 5},               // $9E STZ abs,X
	{Instruction: Sta, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5},           // $9F STA al,X

	// $A0 - $AF
	{Instruction: Ldy, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthX}, // $A0 LDY #imm
	{Instruction: Lda, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6},   // $A1 LDA (dp,X)
	{Instruction: Ldx, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthX}, // $A2 LDX #imm
	{Instruction: Lda, Addressing: StackRelativeAddressing, Timing: 4},                // $A3 LDA sr,S
	{Instruction: Ldy, Addressing: DirectPageAddressing, Timing: 3},                   // $A4 LDY dp
	{Instruction: Lda, Addressing: DirectPageAddressing, Timing: 3},                   // $A5 LDA dp
	{Instruction: Ldx, Addressing: DirectPageAddressing, Timing: 3},                   // $A6 LDX dp
	{Instruction: Lda, Addressing: DirectPageIndirectLongAddressing, Timing: 6},       // $A7 LDA [dp]
	{Instruction: Tay, Addressing: ImpliedAddressing, Timing: 2},                      // $A8 TAY
	{Instruction: Lda, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthM}, // $A9 LDA #imm
	{Instruction: Tax, Addressing: ImpliedAddressing, Timing: 2},                      // $AA TAX
	{Instruction: Plb, Addressing: ImpliedAddressing, Timing: 4},                      // $AB PLB
	{Instruction: Ldy, Addressing: AbsoluteAddressing, Timing: 4},                     // $AC LDY abs
	{Instruction: Lda, Addressing: AbsoluteAddressing, Timing: 4},                     // $AD LDA abs
	{Instruction: Ldx, Addressing: AbsoluteAddressing, Timing: 4},                     // $AE LDX abs
	{Instruction: Lda, Addressing: AbsoluteLongAddressing, Timing: 5},                 // $AF LDA al

	// $B0 - $BF
	{Instruction: Bcs, Addressing: RelativeAddressing, Timing: 2, PageCrossCycle: true},                   // $B0 BCS rel
	{Instruction: Lda, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, PageCrossCycle: true}, // $B1 LDA (dp),Y
	{Instruction: Lda, Addressing: DirectPageIndirectAddressing, Timing: 5},                               // $B2 LDA (dp)
	{Instruction: Lda, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7},                    // $B3 LDA (sr,S),Y
	{Instruction: Ldy, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $B4 LDY dp,X
	{Instruction: Lda, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $B5 LDA dp,X
	{Instruction: Ldx, Addressing: DirectPageIndexedYAddressing, Timing: 4},                               // $B6 LDX dp,Y
	{Instruction: Lda, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6},                   // $B7 LDA [dp],Y
	{Instruction: Clv, Addressing: ImpliedAddressing, Timing: 2},                                          // $B8 CLV
	{Instruction: Lda, Addressing: AbsoluteIndexedYAddressing, Timing: 4, PageCrossCycle: true},           // $B9 LDA abs,Y
	{Instruction: Tsx, Addressing: ImpliedAddressing, Timing: 2},                                          // $BA TSX
	{Instruction: Tyx, Addressing: ImpliedAddressing, Timing: 2},                                          // $BB TYX
	{Instruction: Ldy, Addressing: AbsoluteIndexedXAddressing, Timing: 4, PageCrossCycle: true},           // $BC LDY abs,X
	{Instruction: Lda, Addressing: AbsoluteIndexedXAddressing, Timing: 4, PageCrossCycle: true},           // $BD LDA abs,X
	{Instruction: Ldx, Addressing: AbsoluteIndexedYAddressing, Timing: 4, PageCrossCycle: true},           // $BE LDX abs,Y
	{Instruction: Lda, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5},                             // $BF LDA al,X

	// $C0 - $CF
	{Instruction: Cpy, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthX}, // $C0 CPY #imm
	{Instruction: Cmp, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6},   // $C1 CMP (dp,X)
	{Instruction: Rep, Addressing: ImmediateAddressing, Timing: 3},                    // $C2 REP #imm
	{Instruction: Cmp, Addressing: StackRelativeAddressing, Timing: 4},                // $C3 CMP sr,S
	{Instruction: Cpy, Addressing: DirectPageAddressing, Timing: 3},                   // $C4 CPY dp
	{Instruction: Cmp, Addressing: DirectPageAddressing, Timing: 3},                   // $C5 CMP dp
	{Instruction: Dec, Addressing: DirectPageAddressing, Timing: 5},                   // $C6 DEC dp
	{Instruction: Cmp, Addressing: DirectPageIndirectLongAddressing, Timing: 6},       // $C7 CMP [dp]
	{Instruction: Iny, Addressing: ImpliedAddressing, Timing: 2},                      // $C8 INY
	{Instruction: Cmp, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthM}, // $C9 CMP #imm
	{Instruction: Dex, Addressing: ImpliedAddressing, Timing: 2},                      // $CA DEX
	{Instruction: Wai, Addressing: ImpliedAddressing, Timing: 3},                      // $CB WAI
	{Instruction: Cpy, Addressing: AbsoluteAddressing, Timing: 4},                     // $CC CPY abs
	{Instruction: Cmp, Addressing: AbsoluteAddressing, Timing: 4},                     // $CD CMP abs
	{Instruction: Dec, Addressing: AbsoluteAddressing, Timing: 6},                     // $CE DEC abs
	{Instruction: Cmp, Addressing: AbsoluteLongAddressing, Timing: 5},                 // $CF CMP al

	// $D0 - $DF
	{Instruction: Bne, Addressing: RelativeAddressing, Timing: 2, PageCrossCycle: true},                   // $D0 BNE rel
	{Instruction: Cmp, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, PageCrossCycle: true}, // $D1 CMP (dp),Y
	{Instruction: Cmp, Addressing: DirectPageIndirectAddressing, Timing: 5},                               // $D2 CMP (dp)
	{Instruction: Cmp, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7},                    // $D3 CMP (sr,S),Y
	{Instruction: Pei, Addressing: DirectPageIndirectAddressing, Timing: 6},                               // $D4 PEI (dp)
	{Instruction: Cmp, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $D5 CMP dp,X
	{Instruction: Dec, Addressing: DirectPageIndexedXAddressing, Timing: 6},                               // $D6 DEC dp,X
	{Instruction: Cmp, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6},                   // $D7 CMP [dp],Y
	{Instruction: Cld, Addressing: ImpliedAddressing, Timing: 2},                                          // $D8 CLD
	{Instruction: Cmp, Addressing: AbsoluteIndexedYAddressing, Timing: 4, PageCrossCycle: true},           // $D9 CMP abs,Y
	{Instruction: Phx, Addressing: ImpliedAddressing, Timing: 3},                                          // $DA PHX
	{Instruction: Stp, Addressing: ImpliedAddressing, Timing: 3},                                          // $DB STP
	{Instruction: Jml, Addressing: AbsoluteIndirectLongAddressing, Timing: 6},                             // $DC JML [abs]
	{Instruction: Cmp, Addressing: AbsoluteIndexedXAddressing, Timing: 4, PageCrossCycle: true},           // $DD CMP abs,X
	{Instruction: Dec, Addressing: AbsoluteIndexedXAddressing, Timing: 7},                                 // $DE DEC abs,X
	{Instruction: Cmp, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5},                             // $DF CMP al,X

	// $E0 - $EF
	{Instruction: Cpx, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthX}, // $E0 CPX #imm
	{Instruction: Sbc, Addressing: DirectPageIndexedXIndirectAddressing, Timing: 6},   // $E1 SBC (dp,X)
	{Instruction: Sep, Addressing: ImmediateAddressing, Timing: 3},                    // $E2 SEP #imm
	{Instruction: Sbc, Addressing: StackRelativeAddressing, Timing: 4},                // $E3 SBC sr,S
	{Instruction: Cpx, Addressing: DirectPageAddressing, Timing: 3},                   // $E4 CPX dp
	{Instruction: Sbc, Addressing: DirectPageAddressing, Timing: 3},                   // $E5 SBC dp
	{Instruction: Inc, Addressing: DirectPageAddressing, Timing: 5},                   // $E6 INC dp
	{Instruction: Sbc, Addressing: DirectPageIndirectLongAddressing, Timing: 6},       // $E7 SBC [dp]
	{Instruction: Inx, Addressing: ImpliedAddressing, Timing: 2},                      // $E8 INX
	{Instruction: Sbc, Addressing: ImmediateAddressing, Timing: 2, WidthFlag: WidthM}, // $E9 SBC #imm
	{Instruction: Nop, Addressing: ImpliedAddressing, Timing: 2},                      // $EA NOP
	{Instruction: Xba, Addressing: ImpliedAddressing, Timing: 3},                      // $EB XBA
	{Instruction: Cpx, Addressing: AbsoluteAddressing, Timing: 4},                     // $EC CPX abs
	{Instruction: Sbc, Addressing: AbsoluteAddressing, Timing: 4},                     // $ED SBC abs
	{Instruction: Inc, Addressing: AbsoluteAddressing, Timing: 6},                     // $EE INC abs
	{Instruction: Sbc, Addressing: AbsoluteLongAddressing, Timing: 5},                 // $EF SBC al

	// $F0 - $FF
	{Instruction: Beq, Addressing: RelativeAddressing, Timing: 2, PageCrossCycle: true},                   // $F0 BEQ rel
	{Instruction: Sbc, Addressing: DirectPageIndirectIndexedYAddressing, Timing: 5, PageCrossCycle: true}, // $F1 SBC (dp),Y
	{Instruction: Sbc, Addressing: DirectPageIndirectAddressing, Timing: 5},                               // $F2 SBC (dp)
	{Instruction: Sbc, Addressing: StackRelativeIndirectIndexedYAddressing, Timing: 7},                    // $F3 SBC (sr,S),Y
	{Instruction: Pea, Addressing: AbsoluteAddressing, Timing: 5},                                         // $F4 PEA abs
	{Instruction: Sbc, Addressing: DirectPageIndexedXAddressing, Timing: 4},                               // $F5 SBC dp,X
	{Instruction: Inc, Addressing: DirectPageIndexedXAddressing, Timing: 6},                               // $F6 INC dp,X
	{Instruction: Sbc, Addressing: DirectPageIndirectLongIndexedYAddressing, Timing: 6},                   // $F7 SBC [dp],Y
	{Instruction: Sed, Addressing: ImpliedAddressing, Timing: 2},                                          // $F8 SED
	{Instruction: Sbc, Addressing: AbsoluteIndexedYAddressing, Timing: 4, PageCrossCycle: true},           // $F9 SBC abs,Y
	{Instruction: Plx, Addressing: ImpliedAddressing, Timing: 4},                                          // $FA PLX
	{Instruction: Xce, Addressing: ImpliedAddressing, Timing: 2},                                          // $FB XCE
	{Instruction: Jsr, Addressing: AbsoluteIndexedXIndirectAddressing, Timing: 8},                         // $FC JSR (abs,X)
	{Instruction: Sbc, Addressing: AbsoluteIndexedXAddressing, Timing: 4, PageCrossCycle: true},           // $FD SBC abs,X
	{Instruction: Inc, Addressing: AbsoluteIndexedXAddressing, Timing: 7},                                 // $FE INC abs,X
	{Instruction: Sbc, Addressing: AbsoluteLongIndexedXAddressing, Timing: 5},                             // $FF SBC al,X
}
