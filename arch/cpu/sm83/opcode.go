package sm83

// MaxOpcodeSize is the maximum size of an opcode and its operands in bytes.
const MaxOpcodeSize = 3

// SM83 instruction prefix byte
const (
	PrefixCB = 0xCB // Bit operations prefix
)

// Opcode is a CPU opcode that contains the instruction info and used addressing mode.
type Opcode struct {
	Instruction *Instruction
	Addressing  AddressingMode // Addressing mode
	Timing      byte           // Timing in M-cycles
	Size        byte           // Size of opcode in bytes
}

// Opcodes maps the first opcode byte to CPU instruction information.
// Timing is in M-cycles (machine cycles).
var Opcodes = [256]Opcode{
	{Instruction: NopInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},             // 0x00 NOP
	{Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 3, Size: 3},           // 0x01 LD BC,nn
	{Instruction: LdIndirect, Addressing: RegisterIndirectAddressing, Timing: 2, Size: 1}, // 0x02 LD (BC),A
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1},           // 0x03 INC BC
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},            // 0x04 INC B
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},            // 0x05 DEC B
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 2, Size: 2},            // 0x06 LD B,n
	{Instruction: RlcaInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},            // 0x07 RLCA
	{Instruction: LdAddrSP, Addressing: ExtendedAddressing, Timing: 5, Size: 3},           // 0x08 LD (nn),SP
	{Instruction: AddHL, Addressing: RegisterAddressing, Timing: 2, Size: 1},              // 0x09 ADD HL,BC
	{Instruction: LdIndirect, Addressing: RegisterIndirectAddressing, Timing: 2, Size: 1}, // 0x0A LD A,(BC)
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1},           // 0x0B DEC BC
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},            // 0x0C INC C
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},            // 0x0D DEC C
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 2, Size: 2},            // 0x0E LD C,n
	{Instruction: RrcaInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},            // 0x0F RRCA

	{Instruction: StopInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},            // 0x10 STOP
	{Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 3, Size: 3},           // 0x11 LD DE,nn
	{Instruction: LdIndirect, Addressing: RegisterIndirectAddressing, Timing: 2, Size: 1}, // 0x12 LD (DE),A
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1},           // 0x13 INC DE
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},            // 0x14 INC D
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},            // 0x15 DEC D
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 2, Size: 2},            // 0x16 LD D,n
	{Instruction: RlaInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},             // 0x17 RLA
	{Instruction: JrRel, Addressing: RelativeAddressing, Timing: 3, Size: 2},              // 0x18 JR e
	{Instruction: AddHL, Addressing: RegisterAddressing, Timing: 2, Size: 1},              // 0x19 ADD HL,DE
	{Instruction: LdIndirect, Addressing: RegisterIndirectAddressing, Timing: 2, Size: 1}, // 0x1A LD A,(DE)
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1},           // 0x1B DEC DE
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},            // 0x1C INC E
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},            // 0x1D DEC E
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 2, Size: 2},            // 0x1E LD E,n
	{Instruction: RraInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},             // 0x1F RRA

	{Instruction: JrCond, Addressing: RelativeAddressing, Timing: 2, Size: 2},   // 0x20 JR NZ,e (3 if taken)
	{Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 3, Size: 3}, // 0x21 LD HL,nn
	{Instruction: LdHLPlusA, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // 0x22 LD (HL+),A
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x23 INC HL
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x24 INC H
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x25 DEC H
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 2, Size: 2},  // 0x26 LD H,n
	{Instruction: DaaInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},   // 0x27 DAA
	{Instruction: JrCond, Addressing: RelativeAddressing, Timing: 2, Size: 2},   // 0x28 JR Z,e (3 if taken)
	{Instruction: AddHL, Addressing: RegisterAddressing, Timing: 2, Size: 1},    // 0x29 ADD HL,HL
	{Instruction: LdAHLPlus, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // 0x2A LD A,(HL+)
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x2B DEC HL
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x2C INC L
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x2D DEC L
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 2, Size: 2},  // 0x2E LD L,n
	{Instruction: CplInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},   // 0x2F CPL

	{Instruction: JrCond, Addressing: RelativeAddressing, Timing: 2, Size: 2},              // 0x30 JR NC,e (3 if taken)
	{Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 3, Size: 3},            // 0x31 LD SP,nn
	{Instruction: LdHLMinusA, Addressing: ImpliedAddressing, Timing: 2, Size: 1},           // 0x32 LD (HL-),A
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1},            // 0x33 INC SP
	{Instruction: IncIndirect, Addressing: RegisterIndirectAddressing, Timing: 3, Size: 1}, // 0x34 INC (HL)
	{Instruction: DecIndirect, Addressing: RegisterIndirectAddressing, Timing: 3, Size: 1}, // 0x35 DEC (HL)
	{Instruction: LdIndirectImm, Addressing: ImmediateAddressing, Timing: 3, Size: 2},      // 0x36 LD (HL),n
	{Instruction: ScfInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},              // 0x37 SCF
	{Instruction: JrCond, Addressing: RelativeAddressing, Timing: 2, Size: 2},              // 0x38 JR C,e (3 if taken)
	{Instruction: AddHL, Addressing: RegisterAddressing, Timing: 2, Size: 1},               // 0x39 ADD HL,SP
	{Instruction: LdAHLMinus, Addressing: ImpliedAddressing, Timing: 2, Size: 1},           // 0x3A LD A,(HL-)
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1},            // 0x3B DEC SP
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},             // 0x3C INC A
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},             // 0x3D DEC A
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 2, Size: 2},             // 0x3E LD A,n
	{Instruction: CcfInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},              // 0x3F CCF

	// 0x40-0x7F: LD r,r instructions (register to register)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x40 LD B,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x41 LD B,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x42 LD B,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x43 LD B,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x44 LD B,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x45 LD B,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x46 LD B,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x47 LD B,A
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x48 LD C,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x49 LD C,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x4A LD C,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x4B LD C,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x4C LD C,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x4D LD C,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x4E LD C,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x4F LD C,A

	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x50 LD D,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x51 LD D,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x52 LD D,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x53 LD D,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x54 LD D,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x55 LD D,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x56 LD D,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x57 LD D,A
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x58 LD E,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x59 LD E,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x5A LD E,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x5B LD E,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x5C LD E,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x5D LD E,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x5E LD E,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x5F LD E,A

	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x60 LD H,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x61 LD H,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x62 LD H,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x63 LD H,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x64 LD H,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x65 LD H,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x66 LD H,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x67 LD H,A
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x68 LD L,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x69 LD L,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x6A LD L,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x6B LD L,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x6C LD L,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x6D LD L,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x6E LD L,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x6F LD L,A

	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1},  // 0x70 LD (HL),B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1},  // 0x71 LD (HL),C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1},  // 0x72 LD (HL),D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1},  // 0x73 LD (HL),E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1},  // 0x74 LD (HL),H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1},  // 0x75 LD (HL),L
	{Instruction: HaltInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1}, // 0x76 HALT
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1},  // 0x77 LD (HL),A
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x78 LD A,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x79 LD A,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x7A LD A,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x7B LD A,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x7C LD A,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x7D LD A,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 2, Size: 1},  // 0x7E LD A,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 1, Size: 1},  // 0x7F LD A,A

	// 0x80-0xBF: ALU operations
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x80 ADD A,B
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x81 ADD A,C
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x82 ADD A,D
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x83 ADD A,E
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x84 ADD A,H
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x85 ADD A,L
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x86 ADD A,(HL)
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x87 ADD A,A
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x88 ADC A,B
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x89 ADC A,C
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x8A ADC A,D
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x8B ADC A,E
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x8C ADC A,H
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x8D ADC A,L
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x8E ADC A,(HL)
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x8F ADC A,A

	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x90 SUB B
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x91 SUB C
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x92 SUB D
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x93 SUB E
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x94 SUB H
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x95 SUB L
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x96 SUB (HL)
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x97 SUB A
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x98 SBC A,B
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x99 SBC A,C
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x9A SBC A,D
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x9B SBC A,E
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x9C SBC A,H
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x9D SBC A,L
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0x9E SBC A,(HL)
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0x9F SBC A,A

	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xA0 AND B
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xA1 AND C
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xA2 AND D
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xA3 AND E
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xA4 AND H
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xA5 AND L
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0xA6 AND (HL)
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xA7 AND A
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xA8 XOR B
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xA9 XOR C
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xAA XOR D
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xAB XOR E
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xAC XOR H
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xAD XOR L
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0xAE XOR (HL)
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xAF XOR A

	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xB0 OR B
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xB1 OR C
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xB2 OR D
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xB3 OR E
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xB4 OR H
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xB5 OR L
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0xB6 OR (HL)
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xB7 OR A
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xB8 CP B
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xB9 CP C
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xBA CP D
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xBB CP E
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xBC CP H
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xBD CP L
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 2, Size: 1}, // 0xBE CP (HL)
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 1, Size: 1}, // 0xBF CP A

	// 0xC0-0xFF: Conditional returns, jumps, calls, and immediate operations
	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 2, Size: 1},    // 0xC0 RET NZ (5 if taken)
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 3, Size: 1},  // 0xC1 POP BC
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 3, Size: 3},    // 0xC2 JP NZ,nn (4 if taken)
	{Instruction: JpAbs, Addressing: ExtendedAddressing, Timing: 4, Size: 3},     // 0xC3 JP nn
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 3, Size: 3},  // 0xC4 CALL NZ,nn (6 if taken)
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xC5 PUSH BC
	{Instruction: AddA, Addressing: ImmediateAddressing, Timing: 2, Size: 2},     // 0xC6 ADD A,n
	{Instruction: RstInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1},    // 0xC7 RST 00H
	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 2, Size: 1},    // 0xC8 RET Z (5 if taken)
	{Instruction: RetInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1},    // 0xC9 RET
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 3, Size: 3},    // 0xCA JP Z,nn (4 if taken)
	{}, // 0xCB - Prefix for bit operations
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 3, Size: 3}, // 0xCC CALL Z,nn (6 if taken)
	{Instruction: CallInst, Addressing: ExtendedAddressing, Timing: 6, Size: 3}, // 0xCD CALL nn
	{Instruction: AdcA, Addressing: ImmediateAddressing, Timing: 2, Size: 2},    // 0xCE ADC A,n
	{Instruction: RstInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1},   // 0xCF RST 08H

	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 2, Size: 1},   // 0xD0 RET NC (5 if taken)
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 3, Size: 1}, // 0xD1 POP DE
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 3, Size: 3},   // 0xD2 JP NC,nn (4 if taken)
	{}, // 0xD3 - ILLEGAL
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 3, Size: 3},  // 0xD4 CALL NC,nn (6 if taken)
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xD5 PUSH DE
	{Instruction: SubA, Addressing: ImmediateAddressing, Timing: 2, Size: 2},     // 0xD6 SUB n
	{Instruction: RstInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1},    // 0xD7 RST 10H
	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 2, Size: 1},    // 0xD8 RET C (5 if taken)
	{Instruction: RetiInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1},   // 0xD9 RETI
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 3, Size: 3},    // 0xDA JP C,nn (4 if taken)
	{}, // 0xDB - ILLEGAL
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 3, Size: 3}, // 0xDC CALL C,nn (6 if taken)
	{}, // 0xDD - ILLEGAL
	{Instruction: SbcA, Addressing: ImmediateAddressing, Timing: 2, Size: 2},  // 0xDE SBC A,n
	{Instruction: RstInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1}, // 0xDF RST 18H

	{Instruction: LdhNA, Addressing: ImmediateAddressing, Timing: 3, Size: 2},   // 0xE0 LDH (n),A
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 3, Size: 1}, // 0xE1 POP HL
	{Instruction: LdCA, Addressing: ImpliedAddressing, Timing: 2, Size: 1},      // 0xE2 LD (C),A
	{}, // 0xE3 - ILLEGAL
	{}, // 0xE4 - ILLEGAL
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 4, Size: 1},    // 0xE5 PUSH HL
	{Instruction: AndA, Addressing: ImmediateAddressing, Timing: 2, Size: 2},        // 0xE6 AND n
	{Instruction: RstInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1},       // 0xE7 RST 20H
	{Instruction: AddSPE, Addressing: ImmediateAddressing, Timing: 4, Size: 2},      // 0xE8 ADD SP,e
	{Instruction: JpHL, Addressing: RegisterIndirectAddressing, Timing: 1, Size: 1}, // 0xE9 JP (HL)
	{Instruction: LdAddrA, Addressing: ExtendedAddressing, Timing: 4, Size: 3},      // 0xEA LD (nn),A
	{}, // 0xEB - ILLEGAL
	{}, // 0xEC - ILLEGAL
	{}, // 0xED - ILLEGAL
	{Instruction: XorA, Addressing: ImmediateAddressing, Timing: 2, Size: 2},  // 0xEE XOR n
	{Instruction: RstInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1}, // 0xEF RST 28H

	{Instruction: LdhAN, Addressing: ImmediateAddressing, Timing: 3, Size: 2},   // 0xF0 LDH A,(n)
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 3, Size: 1}, // 0xF1 POP AF
	{Instruction: LdAC, Addressing: ImpliedAddressing, Timing: 2, Size: 1},      // 0xF2 LD A,(C)
	{Instruction: DiInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},    // 0xF3 DI
	{}, // 0xF4 - ILLEGAL
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 4, Size: 1},     // 0xF5 PUSH AF
	{Instruction: OrA, Addressing: ImmediateAddressing, Timing: 2, Size: 2},          // 0xF6 OR n
	{Instruction: RstInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1},        // 0xF7 RST 30H
	{Instruction: LdHLSPOffset, Addressing: ImmediateAddressing, Timing: 3, Size: 2}, // 0xF8 LD HL,SP+e
	{Instruction: LdSPHL, Addressing: RegisterAddressing, Timing: 2, Size: 1},        // 0xF9 LD SP,HL
	{Instruction: LdAAddr, Addressing: ExtendedAddressing, Timing: 4, Size: 3},       // 0xFA LD A,(nn)
	{Instruction: EiInst, Addressing: ImpliedAddressing, Timing: 1, Size: 1},         // 0xFB EI
	{}, // 0xFC - ILLEGAL
	{}, // 0xFD - ILLEGAL
	{Instruction: CpA, Addressing: ImmediateAddressing, Timing: 2, Size: 2},   // 0xFE CP n
	{Instruction: RstInst, Addressing: ImpliedAddressing, Timing: 4, Size: 1}, // 0xFF RST 38H
}
