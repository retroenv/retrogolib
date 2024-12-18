package m6502

// MaxOpcodeSize is the maximum size of an opcode and its operands in bytes.
const MaxOpcodeSize = 3

// Opcode is a CPU opcode that contains the instruction info and used addressing mode.
type Opcode struct {
	Instruction    *Instruction
	Addressing     AddressingMode // Addressing mode
	Timing         byte           // Timing in cycles
	PageCrossCycle bool           // Crossing page boundary takes an additional cycle
}

// OpcodeInfo contains the opcode and timing info for an instruction addressing mode.
type OpcodeInfo struct {
	Opcode byte // First byte of opcode
	Size   byte // Size of opcode in bytes
}

// Opcodes maps the first opcode byte to CPU instruction information.
// Reference https://www.masswerk.at/6502/6502_instruction_set.html
var Opcodes = [256]Opcode{
	{Instruction: Brk, Addressing: ImpliedAddressing, Timing: 7},   // 0x00
	{Instruction: Ora, Addressing: IndirectXAddressing, Timing: 6}, // 0x01
	{}, // 0x02
	{Instruction: Slo, Addressing: IndirectXAddressing, Timing: 8},          // 0x03
	{Instruction: NopUnofficial, Addressing: ZeroPageAddressing, Timing: 3}, // 0x04
	{Instruction: Ora, Addressing: ZeroPageAddressing, Timing: 3},           // 0x05
	{Instruction: Asl, Addressing: ZeroPageAddressing, Timing: 5},           // 0x06
	{Instruction: Slo, Addressing: ZeroPageAddressing, Timing: 5},           // 0x07
	{Instruction: Php, Addressing: ImpliedAddressing, Timing: 3},            // 0x08
	{Instruction: Ora, Addressing: ImmediateAddressing, Timing: 2},          // 0x09
	{Instruction: Asl, Addressing: AccumulatorAddressing, Timing: 2},        // 0x0a
	{}, // 0x0b
	{Instruction: NopUnofficial, Addressing: AbsoluteAddressing, Timing: 4},              // 0x0c
	{Instruction: Ora, Addressing: AbsoluteAddressing, Timing: 4},                        // 0x0d
	{Instruction: Asl, Addressing: AbsoluteAddressing, Timing: 6},                        // 0x0e
	{Instruction: Slo, Addressing: AbsoluteAddressing, Timing: 6},                        // 0x0f
	{Instruction: Bpl, Addressing: RelativeAddressing, Timing: 2},                        // 0x10
	{Instruction: Ora, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0x11
	{}, // 0x12
	{Instruction: Slo, Addressing: IndirectYAddressing, Timing: 8},                                 // 0x13
	{Instruction: NopUnofficial, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0x14
	{Instruction: Ora, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x15
	{Instruction: Asl, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x16
	{Instruction: Slo, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x17
	{Instruction: Clc, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x18
	{Instruction: Ora, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0x19
	{Instruction: NopUnofficial, Addressing: ImpliedAddressing, Timing: 2},                         // 0x1a
	{Instruction: Slo, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0x1b
	{Instruction: NopUnofficial, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x1c
	{Instruction: Ora, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0x1d
	{Instruction: Asl, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x1e
	{Instruction: Slo, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x1f
	{Instruction: Jsr, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x20
	{Instruction: And, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x21
	{}, // 0x22
	{Instruction: Rla, Addressing: IndirectXAddressing, Timing: 8},   // 0x23
	{Instruction: Bit, Addressing: ZeroPageAddressing, Timing: 3},    // 0x24
	{Instruction: And, Addressing: ZeroPageAddressing, Timing: 3},    // 0x25
	{Instruction: Rol, Addressing: ZeroPageAddressing, Timing: 5},    // 0x26
	{Instruction: Rla, Addressing: ZeroPageAddressing, Timing: 5},    // 0x27
	{Instruction: Plp, Addressing: ImpliedAddressing, Timing: 4},     // 0x28
	{Instruction: And, Addressing: ImmediateAddressing, Timing: 2},   // 0x29
	{Instruction: Rol, Addressing: AccumulatorAddressing, Timing: 2}, // 0x2a
	{}, // 0x2b
	{Instruction: Bit, Addressing: AbsoluteAddressing, Timing: 4},                        // 0x2c
	{Instruction: And, Addressing: AbsoluteAddressing, Timing: 4},                        // 0x2d
	{Instruction: Rol, Addressing: AbsoluteAddressing, Timing: 6},                        // 0x2e
	{Instruction: Rla, Addressing: AbsoluteAddressing, Timing: 6},                        // 0x2f
	{Instruction: Bmi, Addressing: RelativeAddressing, Timing: 2},                        // 0x30
	{Instruction: And, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0x31
	{}, // 0x32
	{Instruction: Rla, Addressing: IndirectYAddressing, Timing: 8},                                 // 0x33
	{Instruction: NopUnofficial, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0x34
	{Instruction: And, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x35
	{Instruction: Rol, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x36
	{Instruction: Rla, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x37
	{Instruction: Sec, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x38
	{Instruction: And, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0x39
	{Instruction: NopUnofficial, Addressing: ImpliedAddressing, Timing: 2},                         // 0x3a
	{Instruction: Rla, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0x3b
	{Instruction: NopUnofficial, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x3c
	{Instruction: And, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0x3d
	{Instruction: Rol, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x3e
	{Instruction: Rla, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x3f
	{Instruction: Rti, Addressing: ImpliedAddressing, Timing: 6},                                   // 0x40
	{Instruction: Eor, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x41
	{}, // 0x42
	{Instruction: Sre, Addressing: IndirectXAddressing, Timing: 8},          // 0x43
	{Instruction: NopUnofficial, Addressing: ZeroPageAddressing, Timing: 3}, // 0x44
	{Instruction: Eor, Addressing: ZeroPageAddressing, Timing: 3},           // 0x45
	{Instruction: Lsr, Addressing: ZeroPageAddressing, Timing: 5},           // 0x46
	{Instruction: Sre, Addressing: ZeroPageAddressing, Timing: 5},           // 0x47
	{Instruction: Pha, Addressing: ImpliedAddressing, Timing: 3},            // 0x48
	{Instruction: Eor, Addressing: ImmediateAddressing, Timing: 2},          // 0x49
	{Instruction: Lsr, Addressing: AccumulatorAddressing, Timing: 2},        // 0x4a
	{}, // 0x4b
	{Instruction: Jmp, Addressing: AbsoluteAddressing, Timing: 3},                        // 0x4c
	{Instruction: Eor, Addressing: AbsoluteAddressing, Timing: 4},                        // 0x4d
	{Instruction: Lsr, Addressing: AbsoluteAddressing, Timing: 6},                        // 0x4e
	{Instruction: Sre, Addressing: AbsoluteAddressing, Timing: 6},                        // 0x4f
	{Instruction: Bvc, Addressing: RelativeAddressing, Timing: 2},                        // 0x50
	{Instruction: Eor, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0x51
	{}, // 0x52
	{Instruction: Sre, Addressing: IndirectYAddressing, Timing: 8},                                 // 0x53
	{Instruction: NopUnofficial, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0x54
	{Instruction: Eor, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x55
	{Instruction: Lsr, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x56
	{Instruction: Sre, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x57
	{Instruction: Cli, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x58
	{Instruction: Eor, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0x59
	{Instruction: NopUnofficial, Addressing: ImpliedAddressing, Timing: 2},                         // 0x5a
	{Instruction: Sre, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0x5b
	{Instruction: NopUnofficial, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x5c
	{Instruction: Eor, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0x5d
	{Instruction: Lsr, Addressing: AbsoluteXAddressing, Timing: 7, PageCrossCycle: true},           // 0x5e
	{Instruction: Sre, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x5f
	{Instruction: Rts, Addressing: ImpliedAddressing, Timing: 6},                                   // 0x60
	{Instruction: Adc, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x61
	{}, // 0x62
	{Instruction: Rra, Addressing: IndirectXAddressing, Timing: 8},          // 0x63
	{Instruction: NopUnofficial, Addressing: ZeroPageAddressing, Timing: 3}, // 0x64
	{Instruction: Adc, Addressing: ZeroPageAddressing, Timing: 3},           // 0x65
	{Instruction: Ror, Addressing: ZeroPageAddressing, Timing: 5},           // 0x66
	{Instruction: Rra, Addressing: ZeroPageAddressing, Timing: 5},           // 0x67
	{Instruction: Pla, Addressing: ImpliedAddressing, Timing: 4},            // 0x68
	{Instruction: Adc, Addressing: ImmediateAddressing, Timing: 2},          // 0x69
	{Instruction: Ror, Addressing: AccumulatorAddressing, Timing: 2},        // 0x6a
	{}, // 0x6b
	{Instruction: Jmp, Addressing: IndirectAddressing, Timing: 5},                        // 0x6c
	{Instruction: Adc, Addressing: AbsoluteAddressing, Timing: 4},                        // 0x6d
	{Instruction: Ror, Addressing: AbsoluteAddressing, Timing: 6},                        // 0x6e
	{Instruction: Rra, Addressing: AbsoluteAddressing, Timing: 6},                        // 0x6f
	{Instruction: Bvs, Addressing: RelativeAddressing, Timing: 2},                        // 0x70
	{Instruction: Adc, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0x71
	{}, // 0x72
	{Instruction: Rra, Addressing: IndirectYAddressing, Timing: 8},                                 // 0x73
	{Instruction: NopUnofficial, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0x74
	{Instruction: Adc, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x75
	{Instruction: Ror, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x76
	{Instruction: Rra, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x77
	{Instruction: Sei, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x78
	{Instruction: Adc, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0x79
	{Instruction: NopUnofficial, Addressing: ImpliedAddressing, Timing: 2},                         // 0x7a
	{Instruction: Rra, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0x7b
	{Instruction: NopUnofficial, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x7c
	{Instruction: Adc, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0x7d
	{Instruction: Ror, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x7e
	{Instruction: Rra, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x7f
	{Instruction: NopUnofficial, Addressing: ImmediateAddressing, Timing: 2},                       // 0x80
	{Instruction: Sta, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x81
	{Instruction: NopUnofficial, Addressing: ImmediateAddressing, Timing: 2},                       // 0x82
	{Instruction: Sax, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x83
	{Instruction: Sty, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x84
	{Instruction: Sta, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x85
	{Instruction: Stx, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x86
	{Instruction: Sax, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x87
	{Instruction: Dey, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x88
	{Instruction: NopUnofficial, Addressing: ImmediateAddressing, Timing: 2},                       // 0x89
	{Instruction: Txa, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x8a
	{}, // 0x8b
	{Instruction: Sty, Addressing: AbsoluteAddressing, Timing: 4},  // 0x8c
	{Instruction: Sta, Addressing: AbsoluteAddressing, Timing: 4},  // 0x8d
	{Instruction: Stx, Addressing: AbsoluteAddressing, Timing: 4},  // 0x8e
	{Instruction: Sax, Addressing: AbsoluteAddressing, Timing: 4},  // 0x8f
	{Instruction: Bcc, Addressing: RelativeAddressing, Timing: 2},  // 0x90
	{Instruction: Sta, Addressing: IndirectYAddressing, Timing: 6}, // 0x91
	{}, // 0x92
	{}, // 0x93
	{Instruction: Sty, Addressing: ZeroPageXAddressing, Timing: 4}, // 0x94
	{Instruction: Sta, Addressing: ZeroPageXAddressing, Timing: 4}, // 0x95
	{Instruction: Stx, Addressing: ZeroPageYAddressing, Timing: 4}, // 0x96
	{Instruction: Sax, Addressing: ZeroPageYAddressing, Timing: 4}, // 0x97
	{Instruction: Tya, Addressing: ImpliedAddressing, Timing: 2},   // 0x98
	{Instruction: Sta, Addressing: AbsoluteYAddressing, Timing: 5}, // 0x99
	{Instruction: Txs, Addressing: ImpliedAddressing, Timing: 2},   // 0x9a
	{}, // 0x9b
	{}, // 0x9c
	{Instruction: Sta, Addressing: AbsoluteXAddressing, Timing: 5}, // 0x9d
	{}, // 0x9e
	{}, // 0x9f
	{Instruction: Ldy, Addressing: ImmediateAddressing, Timing: 2}, // 0xa0
	{Instruction: Lda, Addressing: IndirectXAddressing, Timing: 6}, // 0xa1
	{Instruction: Ldx, Addressing: ImmediateAddressing, Timing: 2}, // 0xa2
	{Instruction: Lax, Addressing: IndirectXAddressing, Timing: 6}, // 0xa3
	{Instruction: Ldy, Addressing: ZeroPageAddressing, Timing: 3},  // 0xa4
	{Instruction: Lda, Addressing: ZeroPageAddressing, Timing: 3},  // 0xa5
	{Instruction: Ldx, Addressing: ZeroPageAddressing, Timing: 3},  // 0xa6
	{Instruction: Lax, Addressing: ZeroPageAddressing, Timing: 3},  // 0xa7
	{Instruction: Tay, Addressing: ImpliedAddressing, Timing: 2},   // 0xa8
	{Instruction: Lda, Addressing: ImmediateAddressing, Timing: 2}, // 0xa9
	{Instruction: Tax, Addressing: ImpliedAddressing, Timing: 2},   // 0xaa
	{}, // 0xab
	{Instruction: Ldy, Addressing: AbsoluteAddressing, Timing: 4},                        // 0xac
	{Instruction: Lda, Addressing: AbsoluteAddressing, Timing: 4},                        // 0xad
	{Instruction: Ldx, Addressing: AbsoluteAddressing, Timing: 4},                        // 0xae
	{Instruction: Lax, Addressing: AbsoluteAddressing, Timing: 4},                        // 0xaf
	{Instruction: Bcs, Addressing: RelativeAddressing, Timing: 2},                        // 0xb0
	{Instruction: Lda, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0xb1
	{}, // 0xb2
	{Instruction: Lax, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0xb3
	{Instruction: Ldy, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0xb4
	{Instruction: Lda, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0xb5
	{Instruction: Ldx, Addressing: ZeroPageYAddressing, Timing: 4},                       // 0xb6
	{Instruction: Lax, Addressing: ZeroPageYAddressing, Timing: 4},                       // 0xb7
	{Instruction: Clv, Addressing: ImpliedAddressing, Timing: 2},                         // 0xb8
	{Instruction: Lda, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true}, // 0xb9
	{Instruction: Tsx, Addressing: ImpliedAddressing, Timing: 2},                         // 0xba
	{}, // 0xbb
	{Instruction: Ldy, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0xbc
	{Instruction: Lda, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0xbd
	{Instruction: Ldx, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true}, // 0xbe
	{Instruction: Lax, Addressing: AbsoluteYAddressing, Timing: 4},                       // 0xbf
	{Instruction: Cpy, Addressing: ImmediateAddressing, Timing: 2},                       // 0xc0
	{Instruction: Cmp, Addressing: IndirectXAddressing, Timing: 6},                       // 0xc1
	{Instruction: NopUnofficial, Addressing: ImmediateAddressing, Timing: 2},             // 0xc2
	{Instruction: Dcp, Addressing: IndirectXAddressing, Timing: 8},                       // 0xc3
	{Instruction: Cpy, Addressing: ZeroPageAddressing, Timing: 3},                        // 0xc4
	{Instruction: Cmp, Addressing: ZeroPageAddressing, Timing: 3},                        // 0xc5
	{Instruction: Dec, Addressing: ZeroPageAddressing, Timing: 5},                        // 0xc6
	{Instruction: Dcp, Addressing: ZeroPageAddressing, Timing: 5},                        // 0xc7
	{Instruction: Iny, Addressing: ImpliedAddressing, Timing: 2},                         // 0xc8
	{Instruction: Cmp, Addressing: ImmediateAddressing, Timing: 2},                       // 0xc9
	{Instruction: Dex, Addressing: ImpliedAddressing, Timing: 2},                         // 0xca
	{}, // 0xcb
	{Instruction: Cpy, Addressing: AbsoluteAddressing, Timing: 4},                        // 0xcc
	{Instruction: Cmp, Addressing: AbsoluteAddressing, Timing: 4},                        // 0xcd
	{Instruction: Dec, Addressing: AbsoluteAddressing, Timing: 6},                        // 0xce
	{Instruction: Dcp, Addressing: AbsoluteAddressing, Timing: 6},                        // 0xcf
	{Instruction: Bne, Addressing: RelativeAddressing, Timing: 2},                        // 0xd0
	{Instruction: Cmp, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true}, // 0xd1
	{}, // 0xd2
	{Instruction: Dcp, Addressing: IndirectYAddressing, Timing: 8},                                 // 0xd3
	{Instruction: NopUnofficial, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0xd4
	{Instruction: Cmp, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0xd5
	{Instruction: Dec, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0xd6
	{Instruction: Dcp, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0xd7
	{Instruction: Cld, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xd8
	{Instruction: Cmp, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0xd9
	{Instruction: NopUnofficial, Addressing: ImpliedAddressing, Timing: 2},                         // 0xda
	{Instruction: Dcp, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0xdb
	{Instruction: NopUnofficial, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0xdc
	{Instruction: Cmp, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0xdd
	{Instruction: Dec, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0xde
	{Instruction: Dcp, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0xdf
	{Instruction: Cpx, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xe0
	{Instruction: Sbc, Addressing: IndirectXAddressing, Timing: 6},                                 // 0xe1
	{Instruction: NopUnofficial, Addressing: ImmediateAddressing, Timing: 2},                       // 0xe2
	{Instruction: Isc, Addressing: IndirectXAddressing, Timing: 8},                                 // 0xe3
	{Instruction: Cpx, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xe4
	{Instruction: Sbc, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xe5
	{Instruction: Inc, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0xe6
	{Instruction: Isc, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0xe7
	{Instruction: Inx, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xe8
	{Instruction: Sbc, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xe9
	{Instruction: Nop, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xea
	{Instruction: SbcUnofficial, Addressing: ImmediateAddressing, Timing: 2},                       // 0xeb
	{Instruction: Cpx, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xec
	{Instruction: Sbc, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xed
	{Instruction: Inc, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0xee
	{Instruction: Isc, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0xef
	{Instruction: Beq, Addressing: RelativeAddressing, Timing: 2},                                  // 0xf0
	{Instruction: Sbc, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true},           // 0xf1
	{}, // 0xf2
	{Instruction: Isc, Addressing: IndirectYAddressing, Timing: 8},                                 // 0xf3
	{Instruction: NopUnofficial, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0xf4
	{Instruction: Sbc, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0xf5
	{Instruction: Inc, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0xf6
	{Instruction: Isc, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0xf7
	{Instruction: Sed, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xf8
	{Instruction: Sbc, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0xf9
	{Instruction: NopUnofficial, Addressing: ImpliedAddressing, Timing: 2},                         // 0xfa
	{Instruction: Isc, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0xfb
	{Instruction: NopUnofficial, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0xfc
	{Instruction: Sbc, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0xfd
	{Instruction: Inc, Addressing: AbsoluteXAddressing, Timing: 7, PageCrossCycle: true},           // 0xfe
	{Instruction: Isc, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0xff
}

// ReadsMemory returns whether the instruction accesses memory reading.
func (opcode Opcode) ReadsMemory(memoryReadInstructions map[string]struct{}) bool {
	switch opcode.Addressing {
	case ImmediateAddressing, ImpliedAddressing, RelativeAddressing:
		return false
	}

	_, ok := memoryReadInstructions[opcode.Instruction.Name]
	return ok
}

// WritesMemory returns whether the instruction accesses memory writing.
func (opcode Opcode) WritesMemory(memoryWriteInstructions map[string]struct{}) bool {
	switch opcode.Addressing {
	case ImmediateAddressing, ImpliedAddressing, RelativeAddressing:
		return false
	}

	_, ok := memoryWriteInstructions[opcode.Instruction.Name]
	return ok
}

// ReadWritesMemory returns whether the instruction accesses memory reading and writing.
func (opcode Opcode) ReadWritesMemory(memoryReadWriteInstructions map[string]struct{}) bool {
	switch opcode.Addressing {
	case ImmediateAddressing, ImpliedAddressing, RelativeAddressing:
		return false
	}

	_, ok := memoryReadWriteInstructions[opcode.Instruction.Name]
	return ok
}
