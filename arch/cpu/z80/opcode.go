package z80

import "github.com/retroenv/retrogolib/set"

// MaxOpcodeSize is the maximum size of an opcode and its operands in bytes.
const MaxOpcodeSize = 4

// Z80 instruction prefix bytes
const (
	PrefixCB = 0xCB // Bit operations prefix
	PrefixDD = 0xDD // IX operations prefix
	PrefixED = 0xED // Extended operations prefix
	PrefixFD = 0xFD // IY operations prefix
)

// Opcode is a CPU opcode that contains the instruction info and used addressing mode.
type Opcode struct {
	Instruction *Instruction
	Addressing  AddressingMode // Addressing mode
	Timing      byte           // Timing in T-states
	Size        byte           // Size of opcode in bytes
}

// OpcodeInfo contains the opcode and timing info for an instruction addressing mode.
type OpcodeInfo struct {
	Prefix byte // Prefix byte (0x00 for none, 0xCB/0xED/0xDD/0xFD for prefixed)
	Opcode byte // Opcode byte (after prefix if applicable)
	Size   byte // Size of opcode in bytes
	Cycles byte // Timing in T-states
}

// Opcodes maps the first opcode byte to CPU instruction information.
// Reference: https://jnz.dk/z80/opref.html
var Opcodes = [256]Opcode{
	{Instruction: Nop, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                 // 0x00 NOP
	{Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 10, Size: 3},          // 0x01 LD BC,nn
	{Instruction: LdIndirect, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1}, // 0x02 LD (BC),A
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 6, Size: 1},           // 0x03 INC BC
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},            // 0x04 INC B
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},            // 0x05 DEC B
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 7, Size: 2},            // 0x06 LD B,n
	{Instruction: Rlca, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                // 0x07 RLCA
	{Instruction: ExAf, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                // 0x08 EX AF,AF'
	{Instruction: AddHl, Addressing: RegisterAddressing, Timing: 11, Size: 1},             // 0x09 ADD HL,BC
	{Instruction: LdIndirect, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1}, // 0x0A LD A,(BC)
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 6, Size: 1},           // 0x0B DEC BC
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},            // 0x0C INC C
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},            // 0x0D DEC C
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 7, Size: 2},            // 0x0E LD C,n
	{Instruction: Rrca, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                // 0x0F RRCA

	{Instruction: Djnz, Addressing: RelativeAddressing, Timing: 8, Size: 2},               // 0x10 DJNZ e (8 if not taken, 13 if taken)
	{Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 10, Size: 3},          // 0x11 LD DE,nn
	{Instruction: LdIndirect, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1}, // 0x12 LD (DE),A
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 6, Size: 1},           // 0x13 INC DE
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},            // 0x14 INC D
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},            // 0x15 DEC D
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 7, Size: 2},            // 0x16 LD D,n
	{Instruction: Rla, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                 // 0x17 RLA
	{Instruction: JrRel, Addressing: RelativeAddressing, Timing: 12, Size: 2},             // 0x18 JR e
	{Instruction: AddHl, Addressing: RegisterAddressing, Timing: 11, Size: 1},             // 0x19 ADD HL,DE
	{Instruction: LdIndirect, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1}, // 0x1A LD A,(DE)
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 6, Size: 1},           // 0x1B DEC DE
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},            // 0x1C INC E
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},            // 0x1D DEC E
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 7, Size: 2},            // 0x1E LD E,n
	{Instruction: Rra, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                 // 0x1F RRA

	{Instruction: JrCond, Addressing: RelativeAddressing, Timing: 7, Size: 2},      // 0x20 JR NZ,e (7 if not taken, 12 if taken)
	{Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 10, Size: 3},   // 0x21 LD HL,nn
	{Instruction: LdExtended, Addressing: ExtendedAddressing, Timing: 16, Size: 3}, // 0x22 LD (nn),HL
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 6, Size: 1},    // 0x23 INC HL
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},     // 0x24 INC H
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},     // 0x25 DEC H
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 7, Size: 2},     // 0x26 LD H,n
	{Instruction: Daa, Addressing: ImpliedAddressing, Timing: 4, Size: 1},          // 0x27 DAA
	{Instruction: JrCond, Addressing: RelativeAddressing, Timing: 7, Size: 2},      // 0x28 JR Z,e (7 if not taken, 12 if taken)
	{Instruction: AddHl, Addressing: RegisterAddressing, Timing: 11, Size: 1},      // 0x29 ADD HL,HL
	{Instruction: LdExtended, Addressing: ExtendedAddressing, Timing: 16, Size: 3}, // 0x2A LD HL,(nn)
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 6, Size: 1},    // 0x2B DEC HL
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},     // 0x2C INC L
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},     // 0x2D DEC L
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 7, Size: 2},     // 0x2E LD L,n
	{Instruction: Cpl, Addressing: ImpliedAddressing, Timing: 4, Size: 1},          // 0x2F CPL

	{Instruction: JrCond, Addressing: RelativeAddressing, Timing: 7, Size: 2},                 // 0x30 JR NC,e (7 if not taken, 12 if taken)
	{Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 10, Size: 3},              // 0x31 LD SP,nn
	{Instruction: LdExtended, Addressing: ExtendedAddressing, Timing: 13, Size: 3},            // 0x32 LD (nn),A
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 6, Size: 1},               // 0x33 INC SP
	{Instruction: IncIndirect, Addressing: RegisterIndirectAddressing, Timing: 11, Size: 1},   // 0x34 INC (HL)
	{Instruction: DecIndirect, Addressing: RegisterIndirectAddressing, Timing: 11, Size: 1},   // 0x35 DEC (HL)
	{Instruction: LdIndirectImm, Addressing: RegisterIndirectAddressing, Timing: 10, Size: 2}, // 0x36 LD (HL),n
	{Instruction: Scf, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                     // 0x37 SCF
	{Instruction: JrCond, Addressing: RelativeAddressing, Timing: 7, Size: 2},                 // 0x38 JR C,e (7 if not taken, 12 if taken)
	{Instruction: AddHl, Addressing: RegisterAddressing, Timing: 11, Size: 1},                 // 0x39 ADD HL,SP
	{Instruction: LdExtended, Addressing: ExtendedAddressing, Timing: 13, Size: 3},            // 0x3A LD A,(nn)
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 6, Size: 1},               // 0x3B DEC SP
	{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},                // 0x3C INC A
	{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1},                // 0x3D DEC A
	{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 7, Size: 2},                // 0x3E LD A,n
	{Instruction: Ccf, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                     // 0x3F CCF

	// 0x40-0x7F: LD r,r instructions (register to register)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x40 LD B,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x41 LD B,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x42 LD B,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x43 LD B,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x44 LD B,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x45 LD B,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x46 LD B,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x47 LD B,A
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x48 LD C,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x49 LD C,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x4A LD C,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x4B LD C,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x4C LD C,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x4D LD C,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x4E LD C,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x4F LD C,A

	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x50 LD D,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x51 LD D,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x52 LD D,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x53 LD D,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x54 LD D,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x55 LD D,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x56 LD D,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x57 LD D,A
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x58 LD E,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x59 LD E,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x5A LD E,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x5B LD E,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x5C LD E,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x5D LD E,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x5E LD E,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x5F LD E,A

	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x60 LD H,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x61 LD H,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x62 LD H,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x63 LD H,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x64 LD H,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x65 LD H,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x66 LD H,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x67 LD H,A
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x68 LD L,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x69 LD L,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x6A LD L,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x6B LD L,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x6C LD L,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x6D LD L,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x6E LD L,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x6F LD L,A

	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x70 LD (HL),B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x71 LD (HL),C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x72 LD (HL),D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x73 LD (HL),E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x74 LD (HL),H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x75 LD (HL),L
	{Instruction: Halt, Addressing: ImpliedAddressing, Timing: 4, Size: 1},    // 0x76 HALT
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x77 LD (HL),A
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x78 LD A,B
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x79 LD A,C
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x7A LD A,D
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x7B LD A,E
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x7C LD A,H
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x7D LD A,L
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x7E LD A,(HL)
	{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x7F LD A,A

	// 0x80-0xBF: ALU operations
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x80 ADD A,B
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x81 ADD A,C
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x82 ADD A,D
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x83 ADD A,E
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x84 ADD A,H
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x85 ADD A,L
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x86 ADD A,(HL)
	{Instruction: AddA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x87 ADD A,A
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x88 ADC A,B
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x89 ADC A,C
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x8A ADC A,D
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x8B ADC A,E
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x8C ADC A,H
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x8D ADC A,L
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x8E ADC A,(HL)
	{Instruction: AdcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x8F ADC A,A

	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x90 SUB B
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x91 SUB C
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x92 SUB D
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x93 SUB E
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x94 SUB H
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x95 SUB L
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x96 SUB (HL)
	{Instruction: SubA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x97 SUB A
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x98 SBC A,B
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x99 SBC A,C
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x9A SBC A,D
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x9B SBC A,E
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x9C SBC A,H
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x9D SBC A,L
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0x9E SBC A,(HL)
	{Instruction: SbcA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0x9F SBC A,A

	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xA0 AND B
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xA1 AND C
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xA2 AND D
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xA3 AND E
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xA4 AND H
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xA5 AND L
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0xA6 AND (HL)
	{Instruction: AndA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xA7 AND A
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xA8 XOR B
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xA9 XOR C
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xAA XOR D
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xAB XOR E
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xAC XOR H
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xAD XOR L
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0xAE XOR (HL)
	{Instruction: XorA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xAF XOR A

	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xB0 OR B
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xB1 OR C
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xB2 OR D
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xB3 OR E
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xB4 OR H
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xB5 OR L
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0xB6 OR (HL)
	{Instruction: OrA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xB7 OR A
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xB8 CP B
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xB9 CP C
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xBA CP D
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xBB CP E
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xBC CP H
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xBD CP L
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 7, Size: 1}, // 0xBE CP (HL)
	{Instruction: CpA, Addressing: RegisterAddressing, Timing: 4, Size: 1}, // 0xBF CP A

	// 0xC0-0xFF: Conditional returns, jumps, calls, and immediate operations
	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 5, Size: 1},     // 0xC0 RET NZ (5 if not taken, 11 if taken)
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 10, Size: 1},  // 0xC1 POP BC
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},    // 0xC2 JP NZ,nn
	{Instruction: JpAbs, Addressing: ExtendedAddressing, Timing: 10, Size: 3},     // 0xC3 JP nn
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},  // 0xC4 CALL NZ,nn (10 if not taken, 17 if taken)
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1}, // 0xC5 PUSH BC
	{Instruction: AddA, Addressing: ImmediateAddressing, Timing: 7, Size: 2},      // 0xC6 ADD A,n
	{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1},        // 0xC7 RST 00H
	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 5, Size: 1},     // 0xC8 RET Z (5 if not taken, 11 if taken)
	{Instruction: Ret, Addressing: ImpliedAddressing, Timing: 10, Size: 1},        // 0xC9 RET
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},    // 0xCA JP Z,nn
	{}, // 0xCB - Prefix for bit operations
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3}, // 0xCC CALL Z,nn (10 if not taken, 17 if taken)
	{Instruction: Call, Addressing: ExtendedAddressing, Timing: 17, Size: 3},     // 0xCD CALL nn
	{Instruction: AdcA, Addressing: ImmediateAddressing, Timing: 7, Size: 2},     // 0xCE ADC A,n
	{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1},       // 0xCF RST 08H

	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 5, Size: 1},     // 0xD0 RET NC (5 if not taken, 11 if taken)
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 10, Size: 1},  // 0xD1 POP DE
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},    // 0xD2 JP NC,nn
	{Instruction: OutPort, Addressing: PortAddressing, Timing: 11, Size: 2},       // 0xD3 OUT (n),A
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},  // 0xD4 CALL NC,nn (10 if not taken, 17 if taken)
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1}, // 0xD5 PUSH DE
	{Instruction: SubA, Addressing: ImmediateAddressing, Timing: 7, Size: 2},      // 0xD6 SUB n
	{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1},        // 0xD7 RST 10H
	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 5, Size: 1},     // 0xD8 RET C (5 if not taken, 11 if taken)
	{Instruction: Exx, Addressing: ImpliedAddressing, Timing: 4, Size: 1},         // 0xD9 EXX
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},    // 0xDA JP C,nn
	{Instruction: InPort, Addressing: PortAddressing, Timing: 11, Size: 2},        // 0xDB IN A,(n)
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},  // 0xDC CALL C,nn (10 if not taken, 17 if taken)
	{}, // 0xDD - Prefix for IX operations
	{Instruction: SbcA, Addressing: ImmediateAddressing, Timing: 7, Size: 2}, // 0xDE SBC A,n
	{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1},   // 0xDF RST 18H

	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 5, Size: 1},             // 0xE0 RET PO (5 if not taken, 11 if taken)
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 10, Size: 1},          // 0xE1 POP HL
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},            // 0xE2 JP PO,nn
	{Instruction: ExSp, Addressing: RegisterIndirectAddressing, Timing: 19, Size: 1},      // 0xE3 EX (SP),HL
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},          // 0xE4 CALL PO,nn (10 if not taken, 17 if taken)
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1},         // 0xE5 PUSH HL
	{Instruction: AndA, Addressing: ImmediateAddressing, Timing: 7, Size: 2},              // 0xE6 AND n
	{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1},                // 0xE7 RST 20H
	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 5, Size: 1},             // 0xE8 RET PE (5 if not taken, 11 if taken)
	{Instruction: JpIndirect, Addressing: RegisterIndirectAddressing, Timing: 4, Size: 1}, // 0xE9 JP (HL)
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},            // 0xEA JP PE,nn
	{Instruction: ExDeHl, Addressing: ImpliedAddressing, Timing: 4, Size: 1},              // 0xEB EX DE,HL
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},          // 0xEC CALL PE,nn (10 if not taken, 17 if taken)
	{}, // 0xED - Prefix for extended operations
	{Instruction: XorA, Addressing: ImmediateAddressing, Timing: 7, Size: 2}, // 0xEE XOR n
	{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1},   // 0xEF RST 28H

	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 5, Size: 1},     // 0xF0 RET P (5 if not taken, 11 if taken)
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 10, Size: 1},  // 0xF1 POP AF
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},    // 0xF2 JP P,nn
	{Instruction: Di, Addressing: ImpliedAddressing, Timing: 4, Size: 1},          // 0xF3 DI
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},  // 0xF4 CALL P,nn (10 if not taken, 17 if taken)
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1}, // 0xF5 PUSH AF
	{Instruction: OrA, Addressing: ImmediateAddressing, Timing: 7, Size: 2},       // 0xF6 OR n
	{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1},        // 0xF7 RST 30H
	{Instruction: RetCond, Addressing: ImpliedAddressing, Timing: 5, Size: 1},     // 0xF8 RET M (5 if not taken, 11 if taken)
	{Instruction: LdSp, Addressing: RegisterAddressing, Timing: 6, Size: 1},       // 0xF9 LD SP,HL
	{Instruction: JpCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},    // 0xFA JP M,nn
	{Instruction: Ei, Addressing: ImpliedAddressing, Timing: 4, Size: 1},          // 0xFB EI
	{Instruction: CallCond, Addressing: ExtendedAddressing, Timing: 10, Size: 3},  // 0xFC CALL M,nn (10 if not taken, 17 if taken)
	{}, // 0xFD - Prefix for IY operations
	{Instruction: CpA, Addressing: ImmediateAddressing, Timing: 7, Size: 2}, // 0xFE CP n
	{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1},  // 0xFF RST 38H
}

// EDOpcodes maps ED-prefixed opcodes to instruction information.
// ED prefix (0xED) provides extended Z80 instructions.
// Reference: https://jnz.dk/z80/opref.html
var EDOpcodes = [256]Opcode{
	// 0x40-0x4F: I/O and control instructions
	0x40: {Instruction: EdInBC, Addressing: ImpliedAddressing, Timing: 12, Size: 2},    // IN B,(C)
	0x41: {Instruction: EdOutCB, Addressing: ImpliedAddressing, Timing: 12, Size: 2},   // OUT (C),B
	0x42: {Instruction: EdSbcHlBc, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // SBC HL,BC
	0x43: {Instruction: EdLdNnBc, Addressing: ImpliedAddressing, Timing: 20, Size: 4},  // LD (nn),BC
	0x44: {Instruction: EdNeg, Addressing: ImpliedAddressing, Timing: 8, Size: 2},      // NEG
	0x45: {Instruction: EdRetn, Addressing: ImpliedAddressing, Timing: 14, Size: 2},    // RETN
	0x46: {Instruction: EdIm0, Addressing: ImpliedAddressing, Timing: 8, Size: 2},      // IM 0
	0x47: {Instruction: EdLdIA, Addressing: ImpliedAddressing, Timing: 9, Size: 2},     // LD I,A
	0x48: {Instruction: EdInCC, Addressing: ImpliedAddressing, Timing: 12, Size: 2},    // IN C,(C)
	0x49: {Instruction: EdOutCC, Addressing: ImpliedAddressing, Timing: 12, Size: 2},   // OUT (C),C
	0x4A: {Instruction: EdAdcHlBc, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // ADC HL,BC
	0x4B: {Instruction: EdLdBcNn, Addressing: ImpliedAddressing, Timing: 20, Size: 4},  // LD BC,(nn)
	0x4C: {Instruction: EdNeg, Addressing: ImpliedAddressing, Timing: 8, Size: 2},      // NEG (undocumented)
	0x4D: {Instruction: EdReti, Addressing: ImpliedAddressing, Timing: 14, Size: 2},    // RETI
	0x4E: {Instruction: edIm0Alias, Addressing: ImpliedAddressing, Timing: 8, Size: 2}, // IM 0 (undocumented alias)
	0x4F: {Instruction: EdLdRA, Addressing: ImpliedAddressing, Timing: 9, Size: 2},     // LD R,A

	// 0x50-0x5F
	0x50: {Instruction: EdInDC, Addressing: ImpliedAddressing, Timing: 12, Size: 2},      // IN D,(C)
	0x51: {Instruction: EdOutCD, Addressing: ImpliedAddressing, Timing: 12, Size: 2},     // OUT (C),D
	0x52: {Instruction: EdSbcHlDe, Addressing: ImpliedAddressing, Timing: 15, Size: 2},   // SBC HL,DE
	0x53: {Instruction: EdLdNnDe, Addressing: ImpliedAddressing, Timing: 20, Size: 4},    // LD (nn),DE
	0x54: {Instruction: EdNeg, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // NEG (undocumented)
	0x55: {Instruction: EdRetn, Addressing: ImpliedAddressing, Timing: 14, Size: 2},      // RETN (undocumented)
	0x56: {Instruction: EdIm1, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // IM 1
	0x57: {Instruction: EdLdAI, Addressing: ImpliedAddressing, Timing: 9, Size: 2},       // LD A,I
	0x58: {Instruction: EdInEC, Addressing: ImpliedAddressing, Timing: 12, Size: 2},      // IN E,(C)
	0x59: {Instruction: EdOutCE, Addressing: ImpliedAddressing, Timing: 12, Size: 2},     // OUT (C),E
	0x5A: {Instruction: EdAdcHlDe, Addressing: ImpliedAddressing, Timing: 15, Size: 2},   // ADC HL,DE
	0x5B: {Instruction: EdLdDeNn, Addressing: ImpliedAddressing, Timing: 20, Size: 4},    // LD DE,(nn)
	0x5C: {Instruction: EdNeg, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // NEG (undocumented)
	0x5D: {Instruction: edRetnAlias, Addressing: ImpliedAddressing, Timing: 14, Size: 2}, // RETN (undocumented alias)
	0x5E: {Instruction: EdIm2, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // IM 2
	0x5F: {Instruction: EdLdAR, Addressing: ImpliedAddressing, Timing: 9, Size: 2},       // LD A,R

	// 0x60-0x6F
	0x60: {Instruction: EdInHC, Addressing: ImpliedAddressing, Timing: 12, Size: 2},      // IN H,(C)
	0x61: {Instruction: EdOutCH, Addressing: ImpliedAddressing, Timing: 12, Size: 2},     // OUT (C),H
	0x62: {Instruction: EdSbcHlHl, Addressing: ImpliedAddressing, Timing: 15, Size: 2},   // SBC HL,HL
	0x63: {Instruction: EdLdNnHl, Addressing: ImpliedAddressing, Timing: 20, Size: 4},    // LD (nn),HL
	0x64: {Instruction: EdNeg, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // NEG (undocumented)
	0x65: {Instruction: EdRetn, Addressing: ImpliedAddressing, Timing: 14, Size: 2},      // RETN (undocumented)
	0x66: {Instruction: EdIm0, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // IM 0 (undocumented)
	0x67: {Instruction: EdRrd, Addressing: ImpliedAddressing, Timing: 18, Size: 2},       // RRD
	0x68: {Instruction: EdInLC, Addressing: ImpliedAddressing, Timing: 12, Size: 2},      // IN L,(C)
	0x69: {Instruction: EdOutCL, Addressing: ImpliedAddressing, Timing: 12, Size: 2},     // OUT (C),L
	0x6A: {Instruction: EdAdcHlHl, Addressing: ImpliedAddressing, Timing: 15, Size: 2},   // ADC HL,HL
	0x6B: {Instruction: EdLdHlNn, Addressing: ImpliedAddressing, Timing: 20, Size: 4},    // LD HL,(nn)
	0x6C: {Instruction: EdNeg, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // NEG (undocumented)
	0x6D: {Instruction: edRetnAlias, Addressing: ImpliedAddressing, Timing: 14, Size: 2}, // RETN (undocumented alias)
	0x6E: {Instruction: edIm0Alias, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // IM 0 (undocumented alias)
	0x6F: {Instruction: EdRld, Addressing: ImpliedAddressing, Timing: 18, Size: 2},       // RLD

	// 0x70-0x7F
	0x70: {Instruction: EdInFC, Addressing: ImpliedAddressing, Timing: 12, Size: 2},      // IN F,(C) (undocumented)
	0x71: {Instruction: EdOut0C, Addressing: ImpliedAddressing, Timing: 12, Size: 2},     // OUT (C),0 (undocumented)
	0x72: {Instruction: EdSbcHlSp, Addressing: ImpliedAddressing, Timing: 15, Size: 2},   // SBC HL,SP
	0x73: {Instruction: EdLdNnSp, Addressing: ImpliedAddressing, Timing: 20, Size: 4},    // LD (nn),SP
	0x74: {Instruction: EdNeg, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // NEG (undocumented)
	0x75: {Instruction: EdRetn, Addressing: ImpliedAddressing, Timing: 14, Size: 2},      // RETN (undocumented)
	0x76: {Instruction: EdIm1, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // IM 1 (undocumented)
	0x78: {Instruction: EdInAC, Addressing: ImpliedAddressing, Timing: 12, Size: 2},      // IN A,(C)
	0x79: {Instruction: EdOutCA, Addressing: ImpliedAddressing, Timing: 12, Size: 2},     // OUT (C),A
	0x7A: {Instruction: EdAdcHlSp, Addressing: ImpliedAddressing, Timing: 15, Size: 2},   // ADC HL,SP
	0x7B: {Instruction: EdLdSpNn, Addressing: ImpliedAddressing, Timing: 20, Size: 4},    // LD SP,(nn)
	0x7C: {Instruction: EdNeg, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // NEG (undocumented)
	0x7D: {Instruction: edRetnAlias, Addressing: ImpliedAddressing, Timing: 14, Size: 2}, // RETN (undocumented alias)
	0x7E: {Instruction: EdIm2, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // IM 2 (undocumented)
	0x77: {Instruction: EdNop, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // NOP (undocumented)
	0x7F: {Instruction: EdNop, Addressing: ImpliedAddressing, Timing: 8, Size: 2},        // NOP (undocumented)

	// 0xA0-0xBF: Block operations
	0xA0: {Instruction: EdLdi, Addressing: ImpliedAddressing, Timing: 16, Size: 2},  // LDI
	0xA1: {Instruction: EdCpi, Addressing: ImpliedAddressing, Timing: 16, Size: 2},  // CPI
	0xA2: {Instruction: EdIni, Addressing: ImpliedAddressing, Timing: 16, Size: 2},  // INI
	0xA3: {Instruction: EdOuti, Addressing: ImpliedAddressing, Timing: 16, Size: 2}, // OUTI
	0xA8: {Instruction: EdLdd, Addressing: ImpliedAddressing, Timing: 16, Size: 2},  // LDD
	0xA9: {Instruction: EdCpd, Addressing: ImpliedAddressing, Timing: 16, Size: 2},  // CPD
	0xAA: {Instruction: EdInd, Addressing: ImpliedAddressing, Timing: 16, Size: 2},  // IND
	0xAB: {Instruction: EdOutd, Addressing: ImpliedAddressing, Timing: 16, Size: 2}, // OUTD
	0xB0: {Instruction: EdLdir, Addressing: ImpliedAddressing, Timing: 21, Size: 2}, // LDIR
	0xB1: {Instruction: EdCpir, Addressing: ImpliedAddressing, Timing: 21, Size: 2}, // CPIR
	0xB2: {Instruction: EdInir, Addressing: ImpliedAddressing, Timing: 21, Size: 2}, // INIR
	0xB3: {Instruction: EdOtir, Addressing: ImpliedAddressing, Timing: 21, Size: 2}, // OTIR
	0xB8: {Instruction: EdLddr, Addressing: ImpliedAddressing, Timing: 21, Size: 2}, // LDDR
	0xB9: {Instruction: EdCpdr, Addressing: ImpliedAddressing, Timing: 21, Size: 2}, // CPDR
	0xBA: {Instruction: EdIndr, Addressing: ImpliedAddressing, Timing: 21, Size: 2}, // INDR
	0xBB: {Instruction: EdOtdr, Addressing: ImpliedAddressing, Timing: 21, Size: 2}, // OTDR
}

// CBOpcodes maps CB-prefixed opcodes to instruction information.
// CB prefix (0xCB) provides bit operations and rotate/shift instructions.
// Pattern: bits 0-2 = register (B,C,D,E,H,L,(HL),A), bits 3-5 = bit number (for BIT/RES/SET).
// Reference: https://jnz.dk/z80/opref.html
var CBOpcodes = [256]Opcode{
	// 0x00-0x07: RLC r
	0x00: {Instruction: CBRlc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x01: {Instruction: CBRlc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x02: {Instruction: CBRlc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x03: {Instruction: CBRlc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x04: {Instruction: CBRlc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x05: {Instruction: CBRlc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x06: {Instruction: CBRlc, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x07: {Instruction: CBRlc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0x08-0x0F: RRC r
	0x08: {Instruction: CBRrc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x09: {Instruction: CBRrc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x0A: {Instruction: CBRrc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x0B: {Instruction: CBRrc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x0C: {Instruction: CBRrc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x0D: {Instruction: CBRrc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x0E: {Instruction: CBRrc, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x0F: {Instruction: CBRrc, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0x10-0x17: RL r
	0x10: {Instruction: CBRl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x11: {Instruction: CBRl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x12: {Instruction: CBRl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x13: {Instruction: CBRl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x14: {Instruction: CBRl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x15: {Instruction: CBRl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x16: {Instruction: CBRl, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x17: {Instruction: CBRl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0x18-0x1F: RR r
	0x18: {Instruction: CBRr, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x19: {Instruction: CBRr, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x1A: {Instruction: CBRr, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x1B: {Instruction: CBRr, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x1C: {Instruction: CBRr, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x1D: {Instruction: CBRr, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x1E: {Instruction: CBRr, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x1F: {Instruction: CBRr, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0x20-0x27: SLA r
	0x20: {Instruction: CBSla, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x21: {Instruction: CBSla, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x22: {Instruction: CBSla, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x23: {Instruction: CBSla, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x24: {Instruction: CBSla, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x25: {Instruction: CBSla, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x26: {Instruction: CBSla, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x27: {Instruction: CBSla, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0x28-0x2F: SRA r
	0x28: {Instruction: CBSra, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x29: {Instruction: CBSra, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x2A: {Instruction: CBSra, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x2B: {Instruction: CBSra, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x2C: {Instruction: CBSra, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x2D: {Instruction: CBSra, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x2E: {Instruction: CBSra, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x2F: {Instruction: CBSra, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0x30-0x37: SLL r (unofficial)
	0x30: {Instruction: CBSll, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x31: {Instruction: CBSll, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x32: {Instruction: CBSll, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x33: {Instruction: CBSll, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x34: {Instruction: CBSll, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x35: {Instruction: CBSll, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x36: {Instruction: CBSll, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x37: {Instruction: CBSll, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0x38-0x3F: SRL r
	0x38: {Instruction: CBSrl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x39: {Instruction: CBSrl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x3A: {Instruction: CBSrl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x3B: {Instruction: CBSrl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x3C: {Instruction: CBSrl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x3D: {Instruction: CBSrl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x3E: {Instruction: CBSrl, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x3F: {Instruction: CBSrl, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0x40-0x7F: BIT b,r
	0x40: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x41: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x42: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x43: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x44: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x45: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x46: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 12, Size: 2},
	0x47: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x48: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x49: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x4A: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x4B: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x4C: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x4D: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x4E: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 12, Size: 2},
	0x4F: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x50: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x51: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x52: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x53: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x54: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x55: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x56: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 12, Size: 2},
	0x57: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x58: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x59: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x5A: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x5B: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x5C: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x5D: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x5E: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 12, Size: 2},
	0x5F: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x60: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x61: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x62: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x63: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x64: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x65: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x66: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 12, Size: 2},
	0x67: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x68: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x69: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x6A: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x6B: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x6C: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x6D: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x6E: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 12, Size: 2},
	0x6F: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x70: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x71: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x72: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x73: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x74: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x75: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x76: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 12, Size: 2},
	0x77: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x78: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x79: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x7A: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x7B: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x7C: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x7D: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x7E: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 12, Size: 2},
	0x7F: {Instruction: CBBit, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0x80-0xBF: RES b,r
	0x80: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x81: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x82: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x83: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x84: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x85: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x86: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x87: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x88: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x89: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x8A: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x8B: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x8C: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x8D: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x8E: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x8F: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x90: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x91: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x92: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x93: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x94: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x95: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x96: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x97: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x98: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x99: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x9A: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x9B: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x9C: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x9D: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0x9E: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0x9F: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xA0: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xA1: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xA2: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xA3: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xA4: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xA5: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xA6: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xA7: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xA8: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xA9: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xAA: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xAB: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xAC: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xAD: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xAE: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xAF: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xB0: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xB1: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xB2: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xB3: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xB4: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xB5: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xB6: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xB7: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xB8: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xB9: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xBA: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xBB: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xBC: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xBD: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xBE: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xBF: {Instruction: CBRes, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	// 0xC0-0xFF: SET b,r
	0xC0: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xC1: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xC2: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xC3: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xC4: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xC5: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xC6: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xC7: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xC8: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xC9: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xCA: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xCB: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xCC: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xCD: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xCE: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xCF: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xD0: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xD1: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xD2: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xD3: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xD4: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xD5: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xD6: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xD7: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xD8: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xD9: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xDA: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xDB: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xDC: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xDD: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xDE: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xDF: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xE0: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xE1: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xE2: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xE3: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xE4: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xE5: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xE6: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xE7: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xE8: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xE9: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xEA: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xEB: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xEC: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xED: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xEE: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xEF: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xF0: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xF1: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xF2: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xF3: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xF4: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xF5: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xF6: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xF7: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xF8: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xF9: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xFA: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xFB: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xFC: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xFD: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
	0xFE: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 15, Size: 2},
	0xFF: {Instruction: CBSet, Addressing: ImpliedAddressing, Timing: 8, Size: 2},
}

// DDOpcodes maps DD-prefixed opcodes to instruction information.
// DD prefix (0xDD) provides IX register operations.
// Reference: https://jnz.dk/z80/opref.html
var DDOpcodes = [256]Opcode{
	0x09: {Instruction: DdAddIXBc, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // ADD IX,BC
	0x19: {Instruction: DdAddIXDe, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // ADD IX,DE
	0x21: {Instruction: DdLdIXnn, Addressing: ImpliedAddressing, Timing: 14, Size: 4},  // LD IX,nn
	0x22: {Instruction: DdLdNnIX, Addressing: ImpliedAddressing, Timing: 20, Size: 4},  // LD (nn),IX
	0x23: {Instruction: DdIncIX, Addressing: ImpliedAddressing, Timing: 10, Size: 2},   // INC IX
	0x24: {Instruction: DdIncIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // INC IXH (undoc)
	0x25: {Instruction: DdDecIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // DEC IXH (undoc)
	0x26: {Instruction: DdLdIXHn, Addressing: ImpliedAddressing, Timing: 11, Size: 3},  // LD IXH,n (undoc)
	0x29: {Instruction: DdAddIXIX, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // ADD IX,IX
	0x2A: {Instruction: DdLdIXNn, Addressing: ImpliedAddressing, Timing: 20, Size: 4},  // LD IX,(nn)
	0x2B: {Instruction: DdDecIX, Addressing: ImpliedAddressing, Timing: 10, Size: 2},   // DEC IX
	0x2C: {Instruction: DdIncIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // INC IXL (undoc)
	0x2D: {Instruction: DdDecIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // DEC IXL (undoc)
	0x2E: {Instruction: DdLdIXLn, Addressing: ImpliedAddressing, Timing: 11, Size: 3},  // LD IXL,n (undoc)
	0x34: {Instruction: DdIncIXd, Addressing: ImpliedAddressing, Timing: 23, Size: 3},  // INC (IX+d)
	0x35: {Instruction: DdDecIXd, Addressing: ImpliedAddressing, Timing: 23, Size: 3},  // DEC (IX+d)
	0x36: {Instruction: DdLdIXdN, Addressing: ImpliedAddressing, Timing: 19, Size: 4},  // LD (IX+d),n
	0x39: {Instruction: DdAddIXSp, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // ADD IX,SP

	// LD r,IXH/IXL and LD IXH/IXL,r (undocumented)
	0x44: {Instruction: DdLdBIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD B,IXH
	0x45: {Instruction: DdLdBIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD B,IXL
	0x46: {Instruction: DdLdBIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD B,(IX+d)
	0x4C: {Instruction: DdLdCIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD C,IXH
	0x4D: {Instruction: DdLdCIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD C,IXL
	0x4E: {Instruction: DdLdCIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD C,(IX+d)
	0x54: {Instruction: DdLdDIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD D,IXH
	0x55: {Instruction: DdLdDIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD D,IXL
	0x56: {Instruction: DdLdDIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD D,(IX+d)
	0x5C: {Instruction: DdLdEIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD E,IXH
	0x5D: {Instruction: DdLdEIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD E,IXL
	0x5E: {Instruction: DdLdEIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD E,(IX+d)
	0x60: {Instruction: DdLdIXHB, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXH,B
	0x61: {Instruction: DdLdIXHC, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXH,C
	0x62: {Instruction: DdLdIXHD, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXH,D
	0x63: {Instruction: DdLdIXHE, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXH,E
	0x64: {Instruction: DdLdIXHIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2}, // LD IXH,IXH
	0x65: {Instruction: DdLdIXHIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2}, // LD IXH,IXL
	0x66: {Instruction: DdLdHIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD H,(IX+d)
	0x67: {Instruction: DdLdIXHA, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXH,A
	0x68: {Instruction: DdLdIXLB, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXL,B
	0x69: {Instruction: DdLdIXLC, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXL,C
	0x6A: {Instruction: DdLdIXLD, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXL,D
	0x6B: {Instruction: DdLdIXLE, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXL,E
	0x6C: {Instruction: DdLdIXLIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2}, // LD IXL,IXH
	0x6D: {Instruction: DdLdIXLIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2}, // LD IXL,IXL
	0x6E: {Instruction: DdLdLIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD L,(IX+d)
	0x6F: {Instruction: DdLdIXLA, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IXL,A
	0x70: {Instruction: DdLdIXdB, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IX+d),B
	0x71: {Instruction: DdLdIXdC, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IX+d),C
	0x72: {Instruction: DdLdIXdD, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IX+d),D
	0x73: {Instruction: DdLdIXdE, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IX+d),E
	0x74: {Instruction: DdLdIXdH, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IX+d),H
	0x75: {Instruction: DdLdIXdL, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IX+d),L
	0x77: {Instruction: DdLdIXdA, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IX+d),A
	0x7C: {Instruction: DdLdAIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD A,IXH
	0x7D: {Instruction: DdLdAIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD A,IXL
	0x7E: {Instruction: DdLdAIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD A,(IX+d)

	// Arithmetic with IXH/IXL (undocumented)
	0x84: {Instruction: DdAddAIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // ADD A,IXH
	0x85: {Instruction: DdAddAIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // ADD A,IXL
	0x86: {Instruction: DdAddAIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // ADD A,(IX+d)
	0x8C: {Instruction: DdAdcAIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // ADC A,IXH
	0x8D: {Instruction: DdAdcAIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // ADC A,IXL
	0x8E: {Instruction: DdAdcAIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // ADC A,(IX+d)
	0x94: {Instruction: DdSubIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // SUB IXH
	0x95: {Instruction: DdSubIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // SUB IXL
	0x96: {Instruction: DdSubAIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // SUB (IX+d)
	0x9C: {Instruction: DdSbcAIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // SBC A,IXH
	0x9D: {Instruction: DdSbcAIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // SBC A,IXL
	0x9E: {Instruction: DdSbcAIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // SBC A,(IX+d)
	0xA4: {Instruction: DdAndIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // AND IXH
	0xA5: {Instruction: DdAndIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // AND IXL
	0xA6: {Instruction: DdAndAIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // AND (IX+d)
	0xAC: {Instruction: DdXorIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // XOR IXH
	0xAD: {Instruction: DdXorIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // XOR IXL
	0xAE: {Instruction: DdXorAIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // XOR (IX+d)
	0xB4: {Instruction: DdOrIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},    // OR IXH
	0xB5: {Instruction: DdOrIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},    // OR IXL
	0xB6: {Instruction: DdOrAIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // OR (IX+d)
	0xBC: {Instruction: DdCpIXH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},    // CP IXH
	0xBD: {Instruction: DdCpIXL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},    // CP IXL
	0xBE: {Instruction: DdCpAIXd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // CP (IX+d)
	0xE1: {Instruction: DdPopIX, Addressing: ImpliedAddressing, Timing: 14, Size: 2},   // POP IX
	0xE3: {Instruction: DdExSpIX, Addressing: ImpliedAddressing, Timing: 23, Size: 2},  // EX (SP),IX
	0xE5: {Instruction: DdPushIX, Addressing: ImpliedAddressing, Timing: 15, Size: 2},  // PUSH IX
	0xE9: {Instruction: DdJpIX, Addressing: ImpliedAddressing, Timing: 8, Size: 2},     // JP (IX)
	0xF9: {Instruction: DdLdSpIX, Addressing: ImpliedAddressing, Timing: 10, Size: 2},  // LD SP,IX
}

// FDOpcodes maps FD-prefixed opcodes to instruction information.
// FD prefix (0xFD) provides IY register operations.
// Reference: https://jnz.dk/z80/opref.html
var FDOpcodes = [256]Opcode{
	0x09: {Instruction: FdAddIYBc, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // ADD IY,BC
	0x19: {Instruction: FdAddIYDe, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // ADD IY,DE
	0x21: {Instruction: FdLdIYnn, Addressing: ImpliedAddressing, Timing: 14, Size: 4},  // LD IY,nn
	0x22: {Instruction: FdLdNnIY, Addressing: ImpliedAddressing, Timing: 20, Size: 4},  // LD (nn),IY
	0x23: {Instruction: FdIncIY, Addressing: ImpliedAddressing, Timing: 10, Size: 2},   // INC IY
	0x24: {Instruction: FdIncIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // INC IYH (undoc)
	0x25: {Instruction: FdDecIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // DEC IYH (undoc)
	0x26: {Instruction: FdLdIYHn, Addressing: ImpliedAddressing, Timing: 11, Size: 3},  // LD IYH,n (undoc)
	0x29: {Instruction: FdAddIYIY, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // ADD IY,IY
	0x2A: {Instruction: FdLdIYNn, Addressing: ImpliedAddressing, Timing: 20, Size: 4},  // LD IY,(nn)
	0x2B: {Instruction: FdDecIY, Addressing: ImpliedAddressing, Timing: 10, Size: 2},   // DEC IY
	0x2C: {Instruction: FdIncIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // INC IYL (undoc)
	0x2D: {Instruction: FdDecIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // DEC IYL (undoc)
	0x2E: {Instruction: FdLdIYLn, Addressing: ImpliedAddressing, Timing: 11, Size: 3},  // LD IYL,n (undoc)
	0x34: {Instruction: FdIncIYd, Addressing: ImpliedAddressing, Timing: 23, Size: 3},  // INC (IY+d)
	0x35: {Instruction: FdDecIYd, Addressing: ImpliedAddressing, Timing: 23, Size: 3},  // DEC (IY+d)
	0x36: {Instruction: FdLdIYdN, Addressing: ImpliedAddressing, Timing: 19, Size: 4},  // LD (IY+d),n
	0x39: {Instruction: FdAddIYSp, Addressing: ImpliedAddressing, Timing: 15, Size: 2}, // ADD IY,SP

	// LD r,IYH/IYL and LD IYH/IYL,r (undocumented)
	0x44: {Instruction: FdLdBIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD B,IYH
	0x45: {Instruction: FdLdBIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD B,IYL
	0x46: {Instruction: FdLdBIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD B,(IY+d)
	0x4C: {Instruction: FdLdCIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD C,IYH
	0x4D: {Instruction: FdLdCIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD C,IYL
	0x4E: {Instruction: FdLdCIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD C,(IY+d)
	0x54: {Instruction: FdLdDIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD D,IYH
	0x55: {Instruction: FdLdDIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD D,IYL
	0x56: {Instruction: FdLdDIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD D,(IY+d)
	0x5C: {Instruction: FdLdEIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD E,IYH
	0x5D: {Instruction: FdLdEIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD E,IYL
	0x5E: {Instruction: FdLdEIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD E,(IY+d)
	0x60: {Instruction: FdLdIYHB, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYH,B
	0x61: {Instruction: FdLdIYHC, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYH,C
	0x62: {Instruction: FdLdIYHD, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYH,D
	0x63: {Instruction: FdLdIYHE, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYH,E
	0x64: {Instruction: FdLdIYHIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2}, // LD IYH,IYH
	0x65: {Instruction: FdLdIYHIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2}, // LD IYH,IYL
	0x66: {Instruction: FdLdHIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD H,(IY+d)
	0x67: {Instruction: FdLdIYHA, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYH,A
	0x68: {Instruction: FdLdIYLB, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYL,B
	0x69: {Instruction: FdLdIYLC, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYL,C
	0x6A: {Instruction: FdLdIYLD, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYL,D
	0x6B: {Instruction: FdLdIYLE, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYL,E
	0x6C: {Instruction: FdLdIYLIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2}, // LD IYL,IYH
	0x6D: {Instruction: FdLdIYLIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2}, // LD IYL,IYL
	0x6E: {Instruction: FdLdLIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD L,(IY+d)
	0x6F: {Instruction: FdLdIYLA, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD IYL,A
	0x70: {Instruction: FdLdIYdB, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IY+d),B
	0x71: {Instruction: FdLdIYdC, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IY+d),C
	0x72: {Instruction: FdLdIYdD, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IY+d),D
	0x73: {Instruction: FdLdIYdE, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IY+d),E
	0x74: {Instruction: FdLdIYdH, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IY+d),H
	0x75: {Instruction: FdLdIYdL, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IY+d),L
	0x77: {Instruction: FdLdIYdA, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD (IY+d),A
	0x7C: {Instruction: FdLdAIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD A,IYH
	0x7D: {Instruction: FdLdAIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // LD A,IYL
	0x7E: {Instruction: FdLdAIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // LD A,(IY+d)

	// Arithmetic with IYH/IYL (undocumented)
	0x84: {Instruction: FdAddAIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // ADD A,IYH
	0x85: {Instruction: FdAddAIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // ADD A,IYL
	0x86: {Instruction: FdAddAIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // ADD A,(IY+d)
	0x8C: {Instruction: FdAdcAIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // ADC A,IYH
	0x8D: {Instruction: FdAdcAIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // ADC A,IYL
	0x8E: {Instruction: FdAdcAIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // ADC A,(IY+d)
	0x94: {Instruction: FdSubIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // SUB IYH
	0x95: {Instruction: FdSubIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // SUB IYL
	0x96: {Instruction: FdSubAIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // SUB (IY+d)
	0x9C: {Instruction: FdSbcAIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // SBC A,IYH
	0x9D: {Instruction: FdSbcAIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},  // SBC A,IYL
	0x9E: {Instruction: FdSbcAIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // SBC A,(IY+d)
	0xA4: {Instruction: FdAndIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // AND IYH
	0xA5: {Instruction: FdAndIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // AND IYL
	0xA6: {Instruction: FdAndAIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // AND (IY+d)
	0xAC: {Instruction: FdXorIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // XOR IYH
	0xAD: {Instruction: FdXorIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},   // XOR IYL
	0xAE: {Instruction: FdXorAIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3}, // XOR (IY+d)
	0xB4: {Instruction: FdOrIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},    // OR IYH
	0xB5: {Instruction: FdOrIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},    // OR IYL
	0xB6: {Instruction: FdOrAIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // OR (IY+d)
	0xBC: {Instruction: FdCpIYH, Addressing: ImpliedAddressing, Timing: 8, Size: 2},    // CP IYH
	0xBD: {Instruction: FdCpIYL, Addressing: ImpliedAddressing, Timing: 8, Size: 2},    // CP IYL
	0xBE: {Instruction: FdCpAIYd, Addressing: ImpliedAddressing, Timing: 19, Size: 3},  // CP (IY+d)
	0xE1: {Instruction: FdPopIY, Addressing: ImpliedAddressing, Timing: 14, Size: 2},   // POP IY
	0xE3: {Instruction: FdExSpIY, Addressing: ImpliedAddressing, Timing: 23, Size: 2},  // EX (SP),IY
	0xE5: {Instruction: FdPushIY, Addressing: ImpliedAddressing, Timing: 15, Size: 2},  // PUSH IY
	0xE9: {Instruction: FdJpIY, Addressing: ImpliedAddressing, Timing: 8, Size: 2},     // JP (IY)
	0xF9: {Instruction: FdLdSpIY, Addressing: ImpliedAddressing, Timing: 10, Size: 2},  // LD SP,IY
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
