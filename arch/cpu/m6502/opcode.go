package m6502

import "github.com/retroenv/retrogolib/set"

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
	{Instruction: BrkInst, Addressing: ImpliedAddressing, Timing: 7},                                   // 0x00
	{Instruction: OraInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x01
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x02 - KIL/JAM
	{Instruction: SloInst, Addressing: IndirectXAddressing, Timing: 8},                                 // 0x03
	{Instruction: NopUnofficialInst, Addressing: ZeroPageAddressing, Timing: 3},                        // 0x04
	{Instruction: OraInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x05
	{Instruction: AslInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0x06
	{Instruction: SloInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0x07
	{Instruction: PhpInst, Addressing: ImpliedAddressing, Timing: 3},                                   // 0x08
	{Instruction: OraInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0x09
	{Instruction: AslInst, Addressing: AccumulatorAddressing, Timing: 2},                               // 0x0a
	{Instruction: AncInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0x0b
	{Instruction: NopUnofficialInst, Addressing: AbsoluteAddressing, Timing: 4},                        // 0x0c
	{Instruction: OraInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0x0d
	{Instruction: AslInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x0e
	{Instruction: SloInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x0f
	{Instruction: BplInst, Addressing: RelativeAddressing, Timing: 2},                                  // 0x10
	{Instruction: OraInst, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true},           // 0x11
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x12 - KIL/JAM
	{Instruction: SloInst, Addressing: IndirectYAddressing, Timing: 8},                                 // 0x13
	{Instruction: NopUnofficialInst, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0x14
	{Instruction: OraInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x15
	{Instruction: AslInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x16
	{Instruction: SloInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x17
	{Instruction: ClcInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x18
	{Instruction: OraInst, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0x19
	{Instruction: NopUnofficialInst, Addressing: ImpliedAddressing, Timing: 2},                         // 0x1a
	{Instruction: SloInst, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0x1b
	{Instruction: NopUnofficialInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x1c
	{Instruction: OraInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0x1d
	{Instruction: AslInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x1e
	{Instruction: SloInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x1f
	{Instruction: JsrInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x20
	{Instruction: AndInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x21
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x22 - KIL/JAM
	{Instruction: RlaInst, Addressing: IndirectXAddressing, Timing: 8},                                 // 0x23
	{Instruction: BitInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x24
	{Instruction: AndInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x25
	{Instruction: RolInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0x26
	{Instruction: RlaInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0x27
	{Instruction: PlpInst, Addressing: ImpliedAddressing, Timing: 4},                                   // 0x28
	{Instruction: AndInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0x29
	{Instruction: RolInst, Addressing: AccumulatorAddressing, Timing: 2},                               // 0x2a
	{Instruction: AncUnofficialInst, Addressing: ImmediateAddressing, Timing: 2},                       // 0x2b (alternate)
	{Instruction: BitInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0x2c
	{Instruction: AndInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0x2d
	{Instruction: RolInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x2e
	{Instruction: RlaInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x2f
	{Instruction: BmiInst, Addressing: RelativeAddressing, Timing: 2},                                  // 0x30
	{Instruction: AndInst, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true},           // 0x31
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x32 - KIL/JAM
	{Instruction: RlaInst, Addressing: IndirectYAddressing, Timing: 8},                                 // 0x33
	{Instruction: NopUnofficialInst, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0x34
	{Instruction: AndInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x35
	{Instruction: RolInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x36
	{Instruction: RlaInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x37
	{Instruction: SecInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x38
	{Instruction: AndInst, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0x39
	{Instruction: NopUnofficialInst, Addressing: ImpliedAddressing, Timing: 2},                         // 0x3a
	{Instruction: RlaInst, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0x3b
	{Instruction: NopUnofficialInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x3c
	{Instruction: AndInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0x3d
	{Instruction: RolInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x3e
	{Instruction: RlaInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x3f
	{Instruction: RtiInst, Addressing: ImpliedAddressing, Timing: 6},                                   // 0x40
	{Instruction: EorInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x41
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x42 - KIL/JAM
	{Instruction: SreInst, Addressing: IndirectXAddressing, Timing: 8},                                 // 0x43
	{Instruction: NopUnofficialInst, Addressing: ZeroPageAddressing, Timing: 3},                        // 0x44
	{Instruction: EorInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x45
	{Instruction: LsrInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0x46
	{Instruction: SreInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0x47
	{Instruction: PhaInst, Addressing: ImpliedAddressing, Timing: 3},                                   // 0x48
	{Instruction: EorInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0x49
	{Instruction: LsrInst, Addressing: AccumulatorAddressing, Timing: 2},                               // 0x4a
	{Instruction: AlrInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0x4b
	{Instruction: JmpInst, Addressing: AbsoluteAddressing, Timing: 3},                                  // 0x4c
	{Instruction: EorInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0x4d
	{Instruction: LsrInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x4e
	{Instruction: SreInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x4f
	{Instruction: BvcInst, Addressing: RelativeAddressing, Timing: 2},                                  // 0x50
	{Instruction: EorInst, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true},           // 0x51
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x52 - KIL/JAM
	{Instruction: SreInst, Addressing: IndirectYAddressing, Timing: 8},                                 // 0x53
	{Instruction: NopUnofficialInst, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0x54
	{Instruction: EorInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x55
	{Instruction: LsrInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x56
	{Instruction: SreInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x57
	{Instruction: CliInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x58
	{Instruction: EorInst, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0x59
	{Instruction: NopUnofficialInst, Addressing: ImpliedAddressing, Timing: 2},                         // 0x5a
	{Instruction: SreInst, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0x5b
	{Instruction: NopUnofficialInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x5c
	{Instruction: EorInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0x5d
	{Instruction: LsrInst, Addressing: AbsoluteXAddressing, Timing: 7, PageCrossCycle: true},           // 0x5e
	{Instruction: SreInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x5f
	{Instruction: RtsInst, Addressing: ImpliedAddressing, Timing: 6},                                   // 0x60
	{Instruction: AdcInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x61
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x62 - KIL/JAM
	{Instruction: RraInst, Addressing: IndirectXAddressing, Timing: 8},                                 // 0x63
	{Instruction: NopUnofficialInst, Addressing: ZeroPageAddressing, Timing: 3},                        // 0x64
	{Instruction: AdcInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x65
	{Instruction: RorInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0x66
	{Instruction: RraInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0x67
	{Instruction: PlaInst, Addressing: ImpliedAddressing, Timing: 4},                                   // 0x68
	{Instruction: AdcInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0x69
	{Instruction: RorInst, Addressing: AccumulatorAddressing, Timing: 2},                               // 0x6a
	{Instruction: ArrInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0x6b
	{Instruction: JmpInst, Addressing: IndirectAddressing, Timing: 5},                                  // 0x6c
	{Instruction: AdcInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0x6d
	{Instruction: RorInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x6e
	{Instruction: RraInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0x6f
	{Instruction: BvsInst, Addressing: RelativeAddressing, Timing: 2},                                  // 0x70
	{Instruction: AdcInst, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true},           // 0x71
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x72 - KIL/JAM
	{Instruction: RraInst, Addressing: IndirectYAddressing, Timing: 8},                                 // 0x73
	{Instruction: NopUnofficialInst, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0x74
	{Instruction: AdcInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x75
	{Instruction: RorInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x76
	{Instruction: RraInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0x77
	{Instruction: SeiInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x78
	{Instruction: AdcInst, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0x79
	{Instruction: NopUnofficialInst, Addressing: ImpliedAddressing, Timing: 2},                         // 0x7a
	{Instruction: RraInst, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0x7b
	{Instruction: NopUnofficialInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0x7c
	{Instruction: AdcInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0x7d
	{Instruction: RorInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x7e
	{Instruction: RraInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0x7f
	{Instruction: NopUnofficialInst, Addressing: ImmediateAddressing, Timing: 2},                       // 0x80
	{Instruction: StaInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x81
	{Instruction: NopUnofficialInst, Addressing: ImmediateAddressing, Timing: 2},                       // 0x82
	{Instruction: SaxInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0x83
	{Instruction: StyInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x84
	{Instruction: StaInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x85
	{Instruction: StxInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x86
	{Instruction: SaxInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0x87
	{Instruction: DeyInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x88
	{Instruction: NopUnofficialInst, Addressing: ImmediateAddressing, Timing: 2},                       // 0x89
	{Instruction: TxaInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x8a
	{Instruction: AneInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0x8b
	{Instruction: StyInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0x8c
	{Instruction: StaInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0x8d
	{Instruction: StxInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0x8e
	{Instruction: SaxInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0x8f
	{Instruction: BccInst, Addressing: RelativeAddressing, Timing: 2},                                  // 0x90
	{Instruction: StaInst, Addressing: IndirectYAddressing, Timing: 6},                                 // 0x91
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x92 - KIL/JAM
	{Instruction: ShaInst, Addressing: IndirectYAddressing, Timing: 6},                                 // 0x93
	{Instruction: StyInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x94
	{Instruction: StaInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0x95
	{Instruction: StxInst, Addressing: ZeroPageYAddressing, Timing: 4},                                 // 0x96
	{Instruction: SaxInst, Addressing: ZeroPageYAddressing, Timing: 4},                                 // 0x97
	{Instruction: TyaInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x98
	{Instruction: StaInst, Addressing: AbsoluteYAddressing, Timing: 5},                                 // 0x99
	{Instruction: TxsInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0x9a
	{Instruction: TasInst, Addressing: AbsoluteYAddressing, Timing: 5},                                 // 0x9b
	{Instruction: ShyInst, Addressing: AbsoluteXAddressing, Timing: 5},                                 // 0x9c
	{Instruction: StaInst, Addressing: AbsoluteXAddressing, Timing: 5},                                 // 0x9d
	{Instruction: ShxInst, Addressing: AbsoluteYAddressing, Timing: 5},                                 // 0x9e
	{Instruction: ShaInst, Addressing: AbsoluteYAddressing, Timing: 5},                                 // 0x9f
	{Instruction: LdyInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xa0
	{Instruction: LdaInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0xa1
	{Instruction: LdxInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xa2
	{Instruction: LaxInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0xa3
	{Instruction: LdyInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xa4
	{Instruction: LdaInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xa5
	{Instruction: LdxInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xa6
	{Instruction: LaxInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xa7
	{Instruction: TayInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xa8
	{Instruction: LdaInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xa9
	{Instruction: TaxInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xaa
	{Instruction: LxaInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xab
	{Instruction: LdyInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xac
	{Instruction: LdaInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xad
	{Instruction: LdxInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xae
	{Instruction: LaxInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xaf
	{Instruction: BcsInst, Addressing: RelativeAddressing, Timing: 2},                                  // 0xb0
	{Instruction: LdaInst, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true},           // 0xb1
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xb2 - KIL/JAM
	{Instruction: LaxInst, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true},           // 0xb3
	{Instruction: LdyInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0xb4
	{Instruction: LdaInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0xb5
	{Instruction: LdxInst, Addressing: ZeroPageYAddressing, Timing: 4},                                 // 0xb6
	{Instruction: LaxInst, Addressing: ZeroPageYAddressing, Timing: 4},                                 // 0xb7
	{Instruction: ClvInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xb8
	{Instruction: LdaInst, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0xb9
	{Instruction: TsxInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xba
	{Instruction: LasInst, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0xbb
	{Instruction: LdyInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0xbc
	{Instruction: LdaInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0xbd
	{Instruction: LdxInst, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0xbe
	{Instruction: LaxInst, Addressing: AbsoluteYAddressing, Timing: 4},                                 // 0xbf
	{Instruction: CpyInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xc0
	{Instruction: CmpInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0xc1
	{Instruction: NopUnofficialInst, Addressing: ImmediateAddressing, Timing: 2},                       // 0xc2
	{Instruction: DcpInst, Addressing: IndirectXAddressing, Timing: 8},                                 // 0xc3
	{Instruction: CpyInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xc4
	{Instruction: CmpInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xc5
	{Instruction: DecInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0xc6
	{Instruction: DcpInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0xc7
	{Instruction: InyInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xc8
	{Instruction: CmpInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xc9
	{Instruction: DexInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xca
	{Instruction: AxsInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xcb
	{Instruction: CpyInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xcc
	{Instruction: CmpInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xcd
	{Instruction: DecInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0xce
	{Instruction: DcpInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0xcf
	{Instruction: BneInst, Addressing: RelativeAddressing, Timing: 2},                                  // 0xd0
	{Instruction: CmpInst, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true},           // 0xd1
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xd2 - KIL/JAM
	{Instruction: DcpInst, Addressing: IndirectYAddressing, Timing: 8},                                 // 0xd3
	{Instruction: NopUnofficialInst, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0xd4
	{Instruction: CmpInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0xd5
	{Instruction: DecInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0xd6
	{Instruction: DcpInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0xd7
	{Instruction: CldInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xd8
	{Instruction: CmpInst, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0xd9
	{Instruction: NopUnofficialInst, Addressing: ImpliedAddressing, Timing: 2},                         // 0xda
	{Instruction: DcpInst, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0xdb
	{Instruction: NopUnofficialInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0xdc
	{Instruction: CmpInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0xdd
	{Instruction: DecInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0xde
	{Instruction: DcpInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0xdf
	{Instruction: CpxInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xe0
	{Instruction: SbcInst, Addressing: IndirectXAddressing, Timing: 6},                                 // 0xe1
	{Instruction: NopUnofficialInst, Addressing: ImmediateAddressing, Timing: 2},                       // 0xe2
	{Instruction: IscInst, Addressing: IndirectXAddressing, Timing: 8},                                 // 0xe3
	{Instruction: CpxInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xe4
	{Instruction: SbcInst, Addressing: ZeroPageAddressing, Timing: 3},                                  // 0xe5
	{Instruction: IncInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0xe6
	{Instruction: IscInst, Addressing: ZeroPageAddressing, Timing: 5},                                  // 0xe7
	{Instruction: InxInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xe8
	{Instruction: SbcInst, Addressing: ImmediateAddressing, Timing: 2},                                 // 0xe9
	{Instruction: NopInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xea
	{Instruction: SbcUnofficialInst, Addressing: ImmediateAddressing, Timing: 2},                       // 0xeb
	{Instruction: CpxInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xec
	{Instruction: SbcInst, Addressing: AbsoluteAddressing, Timing: 4},                                  // 0xed
	{Instruction: IncInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0xee
	{Instruction: IscInst, Addressing: AbsoluteAddressing, Timing: 6},                                  // 0xef
	{Instruction: BeqInst, Addressing: RelativeAddressing, Timing: 2},                                  // 0xf0
	{Instruction: SbcInst, Addressing: IndirectYAddressing, Timing: 5, PageCrossCycle: true},           // 0xf1
	{Instruction: KilInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xf2 - KIL/JAM
	{Instruction: IscInst, Addressing: IndirectYAddressing, Timing: 8},                                 // 0xf3
	{Instruction: NopUnofficialInst, Addressing: ZeroPageXAddressing, Timing: 4},                       // 0xf4
	{Instruction: SbcInst, Addressing: ZeroPageXAddressing, Timing: 4},                                 // 0xf5
	{Instruction: IncInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0xf6
	{Instruction: IscInst, Addressing: ZeroPageXAddressing, Timing: 6},                                 // 0xf7
	{Instruction: SedInst, Addressing: ImpliedAddressing, Timing: 2},                                   // 0xf8
	{Instruction: SbcInst, Addressing: AbsoluteYAddressing, Timing: 4, PageCrossCycle: true},           // 0xf9
	{Instruction: NopUnofficialInst, Addressing: ImpliedAddressing, Timing: 2},                         // 0xfa
	{Instruction: IscInst, Addressing: AbsoluteYAddressing, Timing: 7},                                 // 0xfb
	{Instruction: NopUnofficialInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true}, // 0xfc
	{Instruction: SbcInst, Addressing: AbsoluteXAddressing, Timing: 4, PageCrossCycle: true},           // 0xfd
	{Instruction: IncInst, Addressing: AbsoluteXAddressing, Timing: 7, PageCrossCycle: true},           // 0xfe
	{Instruction: IscInst, Addressing: AbsoluteXAddressing, Timing: 7},                                 // 0xff
}

// ReadsMemory returns whether the instruction accesses memory reading.
func (opcode Opcode) ReadsMemory(memoryReadInstructions set.Set[string]) bool {
	switch opcode.Addressing {
	case ImmediateAddressing, ImpliedAddressing, RelativeAddressing:
		return false
	}

	return memoryReadInstructions.Contains(opcode.Instruction.Name)
}

// WritesMemory returns whether the instruction accesses memory writing.
func (opcode Opcode) WritesMemory(memoryWriteInstructions set.Set[string]) bool {
	switch opcode.Addressing {
	case ImmediateAddressing, ImpliedAddressing, RelativeAddressing:
		return false
	}

	return memoryWriteInstructions.Contains(opcode.Instruction.Name)
}

// ReadWritesMemory returns whether the instruction accesses memory reading and writing.
func (opcode Opcode) ReadWritesMemory(memoryReadWriteInstructions set.Set[string]) bool {
	switch opcode.Addressing {
	case ImmediateAddressing, ImpliedAddressing, RelativeAddressing:
		return false
	}

	return memoryReadWriteInstructions.Contains(opcode.Instruction.Name)
}
