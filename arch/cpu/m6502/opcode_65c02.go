// This file contains the 65C02 opcode table.

package m6502

// Nop65C02 is a NOP instruction used in the 65C02 opcode table for slots that
// were undocumented instructions on the NMOS 6502.
var Nop65C02 = &Instruction{
	Name: NopName,
	Addressing: map[AddressingMode]OpcodeInfo{
		ImpliedAddressing: {Opcode: 0xea, Size: 1},
	},
	NoParamFunc: nop,
}

// Opcodes65C02 maps the first opcode byte to CPU instruction information for the 65C02.
// Based on the NMOS 6502 table with undocumented opcodes replaced by NOPs and
// new 65C02 instructions/addressing modes added.
var Opcodes65C02 = [256]Opcode{
	{Instruction: Brk, Addressing: ImpliedAddressing, Timing: 7},                          // 0x00
	{Instruction: Ora65C02, Addressing: IndirectXAddressing, Timing: 6},                   // 0x01
	{Instruction: Nop65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x02 - NOP (was KIL)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x03 - NOP (was SLO)
	{Instruction: Tsb, Addressing: ZeroPageAddressing, Timing: 5},                         // 0x04 - TSB zp
	{Instruction: Ora65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0x05
	{Instruction: Asl, Addressing: ZeroPageAddressing, Timing: 5},                         // 0x06
	{Instruction: Rmb0, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x07 - RMB0 zp
	{Instruction: Php, Addressing: ImpliedAddressing, Timing: 3},                          // 0x08
	{Instruction: Ora65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x09
	{Instruction: Asl, Addressing: AccumulatorAddressing, Timing: 2},                      // 0x0a
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x0b - NOP (was ANC)
	{Instruction: Tsb, Addressing: AbsoluteAddressing, Timing: 6},                         // 0x0c - TSB abs
	{Instruction: Ora65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0x0d
	{Instruction: Asl, Addressing: AbsoluteAddressing, Timing: 6},                         // 0x0e
	{Instruction: Bbr0, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x0f - BBR0 zp,rel
	{Instruction: Bpl, Addressing: RelativeAddressing, Timing: 2},                         // 0x10
	{Instruction: Ora65C02, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0x11
	{Instruction: Ora65C02, Addressing: ZeroPageIndirectAddressing, Timing: 5},            // 0x12 - ORA (zp)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x13 - NOP (was SLO)
	{Instruction: Trb, Addressing: ZeroPageAddressing, Timing: 5},                         // 0x14 - TRB zp
	{Instruction: Ora65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0x15
	{Instruction: Asl, Addressing: ZeroPageXAddressing, Timing: 6},                        // 0x16
	{Instruction: Rmb1, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x17 - RMB1 zp
	{Instruction: Clc, Addressing: ImpliedAddressing, Timing: 2},                          // 0x18
	{Instruction: Ora65C02, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true}, // 0x19
	{Instruction: Inc65C02, Addressing: AccumulatorAddressing, Timing: 2},                 // 0x1a - INC A
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x1b - NOP (was SLO)
	{Instruction: Trb, Addressing: AbsoluteAddressing, Timing: 6},                         // 0x1c - TRB abs
	{Instruction: Ora65C02, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x1d
	{Instruction: Asl, Addressing: AbsoluteXAddressing, Timing: 7},                       // 0x1e
	{Instruction: Bbr1, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x1f - BBR1 zp,rel
	{Instruction: Jsr, Addressing: AbsoluteAddressing, Timing: 6},                         // 0x20
	{Instruction: And65C02, Addressing: IndirectXAddressing, Timing: 6},                   // 0x21
	{Instruction: Nop65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x22 - NOP (was KIL)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x23 - NOP (was RLA)
	{Instruction: Bit65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0x24
	{Instruction: And65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0x25
	{Instruction: Rol, Addressing: ZeroPageAddressing, Timing: 5},                         // 0x26
	{Instruction: Rmb2, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x27 - RMB2 zp
	{Instruction: Plp, Addressing: ImpliedAddressing, Timing: 4},                          // 0x28
	{Instruction: And65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x29
	{Instruction: Rol, Addressing: AccumulatorAddressing, Timing: 2},                      // 0x2a
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x2b - NOP (was ANC)
	{Instruction: Bit65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0x2c
	{Instruction: And65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0x2d
	{Instruction: Rol, Addressing: AbsoluteAddressing, Timing: 6},                         // 0x2e
	{Instruction: Bbr2, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x2f - BBR2 zp,rel
	{Instruction: Bmi, Addressing: RelativeAddressing, Timing: 2},                         // 0x30
	{Instruction: And65C02, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0x31
	{Instruction: And65C02, Addressing: ZeroPageIndirectAddressing, Timing: 5},            // 0x32 - AND (zp)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x33 - NOP (was RLA)
	{Instruction: Bit65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0x34 - BIT zp,X
	{Instruction: And65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0x35
	{Instruction: Rol, Addressing: ZeroPageXAddressing, Timing: 6},                        // 0x36
	{Instruction: Rmb3, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x37 - RMB3 zp
	{Instruction: Sec, Addressing: ImpliedAddressing, Timing: 2},                          // 0x38
	{Instruction: And65C02, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true}, // 0x39
	{Instruction: Dec65C02, Addressing: AccumulatorAddressing, Timing: 2},                 // 0x3a - DEC A
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x3b - NOP (was RLA)
	{Instruction: Bit65C02, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x3c - BIT abs,X
	{Instruction: And65C02, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x3d
	{Instruction: Rol, Addressing: AbsoluteXAddressing, Timing: 7},                       // 0x3e
	{Instruction: Bbr3, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x3f - BBR3 zp,rel
	{Instruction: Rti, Addressing: ImpliedAddressing, Timing: 6},                          // 0x40
	{Instruction: Eor65C02, Addressing: IndirectXAddressing, Timing: 6},                   // 0x41
	{Instruction: Nop65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x42 - NOP (was KIL)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x43 - NOP (was SRE)
	{Instruction: Nop65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0x44 - NOP zp
	{Instruction: Eor65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0x45
	{Instruction: Lsr, Addressing: ZeroPageAddressing, Timing: 5},                         // 0x46
	{Instruction: Rmb4, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x47 - RMB4 zp
	{Instruction: Pha, Addressing: ImpliedAddressing, Timing: 3},                          // 0x48
	{Instruction: Eor65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x49
	{Instruction: Lsr, Addressing: AccumulatorAddressing, Timing: 2},                      // 0x4a
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x4b - NOP (was ALR)
	{Instruction: Jmp65C02, Addressing: AbsoluteAddressing, Timing: 3},                    // 0x4c
	{Instruction: Eor65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0x4d
	{Instruction: Lsr, Addressing: AbsoluteAddressing, Timing: 6},                         // 0x4e
	{Instruction: Bbr4, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x4f - BBR4 zp,rel
	{Instruction: Bvc, Addressing: RelativeAddressing, Timing: 2},                         // 0x50
	{Instruction: Eor65C02, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0x51
	{Instruction: Eor65C02, Addressing: ZeroPageIndirectAddressing, Timing: 5},            // 0x52 - EOR (zp)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x53 - NOP (was SRE)
	{Instruction: Nop65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0x54 - NOP zp,X
	{Instruction: Eor65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0x55
	{Instruction: Lsr, Addressing: ZeroPageXAddressing, Timing: 6},                        // 0x56
	{Instruction: Rmb5, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x57 - RMB5 zp
	{Instruction: Cli, Addressing: ImpliedAddressing, Timing: 2},                          // 0x58
	{Instruction: Eor65C02, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true}, // 0x59
	{Instruction: Phy, Addressing: ImpliedAddressing, Timing: 3},                          // 0x5a - PHY
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x5b - NOP (was SRE)
	{Instruction: Nop65C02, Addressing: AbsoluteAddressing, Timing: 8},                    // 0x5c - NOP abs (8 cycles)
	{Instruction: Eor65C02, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x5d
	{Instruction: Lsr, Addressing: AbsoluteXAddressing, Timing: 7},                       // 0x5e
	{Instruction: Bbr5, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x5f - BBR5 zp,rel
	{Instruction: Rts, Addressing: ImpliedAddressing, Timing: 6},                          // 0x60
	{Instruction: Adc65C02, Addressing: IndirectXAddressing, Timing: 6},                   // 0x61
	{Instruction: Nop65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x62 - NOP (was KIL)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x63 - NOP (was RRA)
	{Instruction: Stz, Addressing: ZeroPageAddressing, Timing: 3},                         // 0x64 - STZ zp
	{Instruction: Adc65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0x65
	{Instruction: Ror, Addressing: ZeroPageAddressing, Timing: 5},                         // 0x66
	{Instruction: Rmb6, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x67 - RMB6 zp
	{Instruction: Pla, Addressing: ImpliedAddressing, Timing: 4},                          // 0x68
	{Instruction: Adc65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x69
	{Instruction: Ror, Addressing: AccumulatorAddressing, Timing: 2},                      // 0x6a
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x6b - NOP (was ARR)
	{Instruction: Jmp65C02, Addressing: IndirectAddressing, Timing: 6},                    // 0x6c - JMP (abs) - 65C02 fixes page bug, 6 cycles
	{Instruction: Adc65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0x6d
	{Instruction: Ror, Addressing: AbsoluteAddressing, Timing: 6},                         // 0x6e
	{Instruction: Bbr6, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x6f - BBR6 zp,rel
	{Instruction: Bvs, Addressing: RelativeAddressing, Timing: 2},                         // 0x70
	{Instruction: Adc65C02, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0x71
	{Instruction: Adc65C02, Addressing: ZeroPageIndirectAddressing, Timing: 5},            // 0x72 - ADC (zp)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x73 - NOP (was RRA)
	{Instruction: Stz, Addressing: ZeroPageXAddressing, Timing: 4},                        // 0x74 - STZ zp,X
	{Instruction: Adc65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0x75
	{Instruction: Ror, Addressing: ZeroPageXAddressing, Timing: 6},                        // 0x76
	{Instruction: Rmb7, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x77 - RMB7 zp
	{Instruction: Sei, Addressing: ImpliedAddressing, Timing: 2},                          // 0x78
	{Instruction: Adc65C02, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true}, // 0x79
	{Instruction: Ply, Addressing: ImpliedAddressing, Timing: 4},                          // 0x7a - PLY
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x7b - NOP (was RRA)
	{Instruction: Jmp65C02, Addressing: AbsoluteXIndirectAddressing, Timing: 6},           // 0x7c - JMP (abs,X)
	{Instruction: Adc65C02, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x7d
	{Instruction: Ror, Addressing: AbsoluteXAddressing, Timing: 7},                       // 0x7e
	{Instruction: Bbr7, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x7f - BBR7 zp,rel
	{Instruction: Bra, Addressing: RelativeAddressing, Timing: 3},                         // 0x80 - BRA
	{Instruction: Sta65C02, Addressing: IndirectXAddressing, Timing: 6},                   // 0x81
	{Instruction: Nop65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x82 - NOP (was NOP imm)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x83 - NOP (was SAX)
	{Instruction: Sty, Addressing: ZeroPageAddressing, Timing: 3},                         // 0x84
	{Instruction: Sta65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0x85
	{Instruction: Stx, Addressing: ZeroPageAddressing, Timing: 3},                         // 0x86
	{Instruction: Smb0, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x87 - SMB0 zp
	{Instruction: Dey, Addressing: ImpliedAddressing, Timing: 2},                          // 0x88
	{Instruction: Bit65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0x89 - BIT #imm
	{Instruction: Txa, Addressing: ImpliedAddressing, Timing: 2},                          // 0x8a
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x8b - NOP (was ANE)
	{Instruction: Sty, Addressing: AbsoluteAddressing, Timing: 4},                         // 0x8c
	{Instruction: Sta65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0x8d
	{Instruction: Stx, Addressing: AbsoluteAddressing, Timing: 4},                         // 0x8e
	{Instruction: Bbs0, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x8f - BBS0 zp,rel
	{Instruction: Bcc, Addressing: RelativeAddressing, Timing: 2},                         // 0x90
	{Instruction: Sta65C02, Addressing: IndirectYAddressing, Timing: 6},                   // 0x91
	{Instruction: Sta65C02, Addressing: ZeroPageIndirectAddressing, Timing: 5},            // 0x92 - STA (zp)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x93 - NOP (was SHA)
	{Instruction: Sty, Addressing: ZeroPageXAddressing, Timing: 4},                        // 0x94
	{Instruction: Sta65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0x95
	{Instruction: Stx, Addressing: ZeroPageYAddressing, Timing: 4},                        // 0x96
	{Instruction: Smb1, Addressing: ZeroPageAddressing, Timing: 5},                        // 0x97 - SMB1 zp
	{Instruction: Tya, Addressing: ImpliedAddressing, Timing: 2},                          // 0x98
	{Instruction: Sta65C02, Addressing: AbsoluteYAddressing, Timing: 5},                   // 0x99
	{Instruction: Txs, Addressing: ImpliedAddressing, Timing: 2},                          // 0x9a
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0x9b - NOP (was TAS)
	{Instruction: Stz, Addressing: AbsoluteAddressing, Timing: 4},                         // 0x9c - STZ abs
	{Instruction: Sta65C02, Addressing: AbsoluteXAddressing, Timing: 5},                   // 0x9d
	{Instruction: Stz, Addressing: AbsoluteXAddressing, Timing: 5},                        // 0x9e - STZ abs,X
	{Instruction: Bbs1, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0x9f - BBS1 zp,rel
	{Instruction: Ldy, Addressing: ImmediateAddressing, Timing: 2},                        // 0xa0
	{Instruction: Lda65C02, Addressing: IndirectXAddressing, Timing: 6},                   // 0xa1
	{Instruction: Ldx, Addressing: ImmediateAddressing, Timing: 2},                        // 0xa2
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xa3 - NOP (was LAX)
	{Instruction: Ldy, Addressing: ZeroPageAddressing, Timing: 3},                         // 0xa4
	{Instruction: Lda65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0xa5
	{Instruction: Ldx, Addressing: ZeroPageAddressing, Timing: 3},                         // 0xa6
	{Instruction: Smb2, Addressing: ZeroPageAddressing, Timing: 5},                        // 0xa7 - SMB2 zp
	{Instruction: Tay, Addressing: ImpliedAddressing, Timing: 2},                          // 0xa8
	{Instruction: Lda65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0xa9
	{Instruction: Tax, Addressing: ImpliedAddressing, Timing: 2},                          // 0xaa
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xab - NOP (was LXA)
	{Instruction: Ldy, Addressing: AbsoluteAddressing, Timing: 4},                         // 0xac
	{Instruction: Lda65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0xad
	{Instruction: Ldx, Addressing: AbsoluteAddressing, Timing: 4},                         // 0xae
	{Instruction: Bbs2, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0xaf - BBS2 zp,rel
	{Instruction: Bcs, Addressing: RelativeAddressing, Timing: 2},                         // 0xb0
	{Instruction: Lda65C02, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0xb1
	{Instruction: Lda65C02, Addressing: ZeroPageIndirectAddressing, Timing: 5},            // 0xb2 - LDA (zp)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xb3 - NOP (was LAX)
	{Instruction: Ldy, Addressing: ZeroPageXAddressing, Timing: 4},                        // 0xb4
	{Instruction: Lda65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0xb5
	{Instruction: Ldx, Addressing: ZeroPageYAddressing, Timing: 4},                        // 0xb6
	{Instruction: Smb3, Addressing: ZeroPageAddressing, Timing: 5},                        // 0xb7 - SMB3 zp
	{Instruction: Clv, Addressing: ImpliedAddressing, Timing: 2},                          // 0xb8
	{Instruction: Lda65C02, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true}, // 0xb9
	{Instruction: Tsx, Addressing: ImpliedAddressing, Timing: 2},                          // 0xba
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xbb - NOP (was LAS)
	{Instruction: Ldy, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},  // 0xbc
	{Instruction: Lda65C02, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0xbd
	{Instruction: Ldx, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},  // 0xbe
	{Instruction: Bbs3, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0xbf - BBS3 zp,rel
	{Instruction: Cpy, Addressing: ImmediateAddressing, Timing: 2},                        // 0xc0
	{Instruction: Cmp65C02, Addressing: IndirectXAddressing, Timing: 6},                   // 0xc1
	{Instruction: Nop65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0xc2 - NOP (was NOP imm)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xc3 - NOP (was DCP)
	{Instruction: Cpy, Addressing: ZeroPageAddressing, Timing: 3},                         // 0xc4
	{Instruction: Cmp65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0xc5
	{Instruction: Dec65C02, Addressing: ZeroPageAddressing, Timing: 5},                    // 0xc6
	{Instruction: Smb4, Addressing: ZeroPageAddressing, Timing: 5},                        // 0xc7 - SMB4 zp
	{Instruction: Iny, Addressing: ImpliedAddressing, Timing: 2},                          // 0xc8
	{Instruction: Cmp65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0xc9
	{Instruction: Dex, Addressing: ImpliedAddressing, Timing: 2},                          // 0xca
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xcb - NOP (was AXS)
	{Instruction: Cpy, Addressing: AbsoluteAddressing, Timing: 4},                         // 0xcc
	{Instruction: Cmp65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0xcd
	{Instruction: Dec65C02, Addressing: AbsoluteAddressing, Timing: 6},                    // 0xce
	{Instruction: Bbs4, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0xcf - BBS4 zp,rel
	{Instruction: Bne, Addressing: RelativeAddressing, Timing: 2},                         // 0xd0
	{Instruction: Cmp65C02, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0xd1
	{Instruction: Cmp65C02, Addressing: ZeroPageIndirectAddressing, Timing: 5},            // 0xd2 - CMP (zp)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xd3 - NOP (was DCP)
	{Instruction: Nop65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0xd4 - NOP zp,X
	{Instruction: Cmp65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0xd5
	{Instruction: Dec65C02, Addressing: ZeroPageXAddressing, Timing: 6},                   // 0xd6
	{Instruction: Smb5, Addressing: ZeroPageAddressing, Timing: 5},                        // 0xd7 - SMB5 zp
	{Instruction: Cld, Addressing: ImpliedAddressing, Timing: 2},                          // 0xd8
	{Instruction: Cmp65C02, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true}, // 0xd9
	{Instruction: Phx, Addressing: ImpliedAddressing, Timing: 3},                          // 0xda - PHX
	{Instruction: Nop65C02, Addressing: ImmediateAddressing, Timing: 4},                    // 0xdb - NOP imm (STP on WDC)
	{Instruction: Nop65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0xdc - NOP abs
	{Instruction: Cmp65C02, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0xdd
	{Instruction: Dec65C02, Addressing: AbsoluteXAddressing, Timing: 7},                   // 0xde
	{Instruction: Bbs5, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0xdf - BBS5 zp,rel
	{Instruction: Cpx, Addressing: ImmediateAddressing, Timing: 2},                        // 0xe0
	{Instruction: Sbc65C02, Addressing: IndirectXAddressing, Timing: 6},                   // 0xe1
	{Instruction: Nop65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0xe2 - NOP (was NOP imm)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xe3 - NOP (was ISC)
	{Instruction: Cpx, Addressing: ZeroPageAddressing, Timing: 3},                         // 0xe4
	{Instruction: Sbc65C02, Addressing: ZeroPageAddressing, Timing: 3},                    // 0xe5
	{Instruction: Inc65C02, Addressing: ZeroPageAddressing, Timing: 5},                    // 0xe6
	{Instruction: Smb6, Addressing: ZeroPageAddressing, Timing: 5},                        // 0xe7 - SMB6 zp
	{Instruction: Inx, Addressing: ImpliedAddressing, Timing: 2},                          // 0xe8
	{Instruction: Sbc65C02, Addressing: ImmediateAddressing, Timing: 2},                   // 0xe9
	{Instruction: Nop, Addressing: ImpliedAddressing, Timing: 2},                          // 0xea
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xeb - NOP (was SBC unofficial)
	{Instruction: Cpx, Addressing: AbsoluteAddressing, Timing: 4},                         // 0xec
	{Instruction: Sbc65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0xed
	{Instruction: Inc65C02, Addressing: AbsoluteAddressing, Timing: 6},                    // 0xee
	{Instruction: Bbs6, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0xef - BBS6 zp,rel
	{Instruction: Beq, Addressing: RelativeAddressing, Timing: 2},                         // 0xf0
	{Instruction: Sbc65C02, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0xf1
	{Instruction: Sbc65C02, Addressing: ZeroPageIndirectAddressing, Timing: 5},            // 0xf2 - SBC (zp)
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xf3 - NOP (was ISC)
	{Instruction: Nop65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0xf4 - NOP zp,X
	{Instruction: Sbc65C02, Addressing: ZeroPageXAddressing, Timing: 4},                   // 0xf5
	{Instruction: Inc65C02, Addressing: ZeroPageXAddressing, Timing: 6},                   // 0xf6
	{Instruction: Smb7, Addressing: ZeroPageAddressing, Timing: 5},                        // 0xf7 - SMB7 zp
	{Instruction: Sed, Addressing: ImpliedAddressing, Timing: 2},                          // 0xf8
	{Instruction: Sbc65C02, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true}, // 0xf9
	{Instruction: Plx, Addressing: ImpliedAddressing, Timing: 4},                          // 0xfa - PLX
	{Instruction: Nop65C02, Addressing: ImpliedAddressing, Timing: 1},                     // 0xfb - NOP (was ISC)
	{Instruction: Nop65C02, Addressing: AbsoluteAddressing, Timing: 4},                    // 0xfc - NOP abs
	{Instruction: Sbc65C02, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0xfd
	{Instruction: Inc65C02, Addressing: AbsoluteXAddressing, Timing: 7},                   // 0xfe
	{Instruction: Bbs7, Addressing: ZeroPageRelativeAddressing, Timing: 5},                // 0xff - BBS7 zp,rel
}
