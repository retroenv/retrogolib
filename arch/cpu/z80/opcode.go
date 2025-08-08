package z80

// MaxOpcodeSize is the maximum size of an opcode and its operands in bytes.
const MaxOpcodeSize = 4 // Z80 has some 4-byte instructions with prefixes

// Opcode is a CPU opcode that contains the instruction info and used addressing mode.
type Opcode struct {
	Instruction *Instruction
	Addressing  AddressingMode // Addressing mode
	Timing      uint8          // Timing in cycles
	Size        uint8          // Size in bytes
}

// OpcodeInfo contains the opcode and timing info for an instruction addressing mode.
type OpcodeInfo struct {
	Opcode uint8 // First byte of opcode (or main opcode for prefixed instructions)
	Size   uint8 // Size of opcode in bytes
	Cycles uint8 // Timing in cycles
}

// Opcodes maps the first opcode byte to CPU instruction information.
// Reference: Z80 CPU User Manual
var Opcodes = [256]Opcode{
	{Instruction: Nop, Addressing: ImpliedAddressing, Timing: 4, Size: 1}, // 0x00 - NOP
	{}, // 0x01 - LD BC,nn (will be handled by 16-bit load instructions)
	{}, // 0x02 - LD (BC),A
	{}, // 0x03 - INC BC
	{}, // 0x04 - INC B
	{}, // 0x05 - DEC B
	{}, // 0x06 - LD B,n
	{}, // 0x07 - RLCA
	{}, // 0x08 - EX AF,AF'
	{}, // 0x09 - ADD HL,BC
	{}, // 0x0A - LD A,(BC)
	{}, // 0x0B - DEC BC
	{}, // 0x0C - INC C
	{}, // 0x0D - DEC C
	{}, // 0x0E - LD C,n
	{}, // 0x0F - RRCA

	{},                                 // 0x10 - DJNZ e
	{},                                 // 0x11 - LD DE,nn
	{},                                 // 0x12 - LD (DE),A
	{},                                 // 0x13 - INC DE
	{},                                 // 0x14 - INC D
	{},                                 // 0x15 - DEC D
	{},                                 // 0x16 - LD D,n
	{},                                 // 0x17 - RLA
	{JrRel, RelativeAddressing, 12, 2}, // 0x18 - JR e
	{},                                 // 0x19 - ADD HL,DE
	{},                                 // 0x1A - LD A,(DE)
	{},                                 // 0x1B - DEC DE
	{},                                 // 0x1C - INC E
	{},                                 // 0x1D - DEC E
	{},                                 // 0x1E - LD E,n
	{},                                 // 0x1F - RRA

	{}, // 0x20 - JR NZ,e
	{}, // 0x21 - LD HL,nn
	{}, // 0x22 - LD (nn),HL
	{}, // 0x23 - INC HL
	{}, // 0x24 - INC H
	{}, // 0x25 - DEC H
	{}, // 0x26 - LD H,n
	{}, // 0x27 - DAA
	{}, // 0x28 - JR Z,e
	{}, // 0x29 - ADD HL,HL
	{}, // 0x2A - LD HL,(nn)
	{}, // 0x2B - DEC HL
	{}, // 0x2C - INC L
	{}, // 0x2D - DEC L
	{}, // 0x2E - LD L,n
	{}, // 0x2F - CPL

	{},                                  // 0x30 - JR NC,e
	{},                                  // 0x31 - LD SP,nn
	{},                                  // 0x32 - LD (nn),A
	{},                                  // 0x33 - INC SP
	{},                                  // 0x34 - INC (HL)
	{},                                  // 0x35 - DEC (HL)
	{},                                  // 0x36 - LD (HL),n
	{},                                  // 0x37 - SCF
	{},                                  // 0x38 - JR C,e
	{},                                  // 0x39 - ADD HL,SP
	{},                                  // 0x3A - LD A,(nn)
	{},                                  // 0x3B - DEC SP
	{IncReg8, RegisterAddressing, 4, 1}, // 0x3C - INC A
	{DecReg8, RegisterAddressing, 4, 1}, // 0x3D - DEC A
	{LdImm8, ImmediateAddressing, 7, 2}, // 0x3E - LD A,n
	{},                                  // 0x3F - CCF

	// 0x40-0x7F: LD r,r' instructions (8-bit register to register)
	{LdReg8, RegisterAddressing, 4, 1}, // 0x40 - LD B,B
	{LdReg8, RegisterAddressing, 4, 1}, // 0x41 - LD B,C
	{LdReg8, RegisterAddressing, 4, 1}, // 0x42 - LD B,D
	{LdReg8, RegisterAddressing, 4, 1}, // 0x43 - LD B,E
	{LdReg8, RegisterAddressing, 4, 1}, // 0x44 - LD B,H
	{LdReg8, RegisterAddressing, 4, 1}, // 0x45 - LD B,L
	{},                                 // 0x46 - LD B,(HL)
	{LdReg8, RegisterAddressing, 4, 1}, // 0x47 - LD B,A
	{LdReg8, RegisterAddressing, 4, 1}, // 0x48 - LD C,B
	{LdReg8, RegisterAddressing, 4, 1}, // 0x49 - LD C,C
	{LdReg8, RegisterAddressing, 4, 1}, // 0x4A - LD C,D
	{LdReg8, RegisterAddressing, 4, 1}, // 0x4B - LD C,E
	{LdReg8, RegisterAddressing, 4, 1}, // 0x4C - LD C,H
	{LdReg8, RegisterAddressing, 4, 1}, // 0x4D - LD C,L
	{},                                 // 0x4E - LD C,(HL)
	{LdReg8, RegisterAddressing, 4, 1}, // 0x4F - LD C,A

	{LdReg8, RegisterAddressing, 4, 1}, // 0x50 - LD D,B
	{LdReg8, RegisterAddressing, 4, 1}, // 0x51 - LD D,C
	{LdReg8, RegisterAddressing, 4, 1}, // 0x52 - LD D,D
	{LdReg8, RegisterAddressing, 4, 1}, // 0x53 - LD D,E
	{LdReg8, RegisterAddressing, 4, 1}, // 0x54 - LD D,H
	{LdReg8, RegisterAddressing, 4, 1}, // 0x55 - LD D,L
	{},                                 // 0x56 - LD D,(HL)
	{LdReg8, RegisterAddressing, 4, 1}, // 0x57 - LD D,A
	{LdReg8, RegisterAddressing, 4, 1}, // 0x58 - LD E,B
	{LdReg8, RegisterAddressing, 4, 1}, // 0x59 - LD E,C
	{LdReg8, RegisterAddressing, 4, 1}, // 0x5A - LD E,D
	{LdReg8, RegisterAddressing, 4, 1}, // 0x5B - LD E,E
	{LdReg8, RegisterAddressing, 4, 1}, // 0x5C - LD E,H
	{LdReg8, RegisterAddressing, 4, 1}, // 0x5D - LD E,L
	{},                                 // 0x5E - LD E,(HL)
	{LdReg8, RegisterAddressing, 4, 1}, // 0x5F - LD E,A

	{LdReg8, RegisterAddressing, 4, 1}, // 0x60 - LD H,B
	{LdReg8, RegisterAddressing, 4, 1}, // 0x61 - LD H,C
	{LdReg8, RegisterAddressing, 4, 1}, // 0x62 - LD H,D
	{LdReg8, RegisterAddressing, 4, 1}, // 0x63 - LD H,E
	{LdReg8, RegisterAddressing, 4, 1}, // 0x64 - LD H,H
	{LdReg8, RegisterAddressing, 4, 1}, // 0x65 - LD H,L
	{},                                 // 0x66 - LD H,(HL)
	{LdReg8, RegisterAddressing, 4, 1}, // 0x67 - LD H,A
	{LdReg8, RegisterAddressing, 4, 1}, // 0x68 - LD L,B
	{LdReg8, RegisterAddressing, 4, 1}, // 0x69 - LD L,C
	{LdReg8, RegisterAddressing, 4, 1}, // 0x6A - LD L,D
	{LdReg8, RegisterAddressing, 4, 1}, // 0x6B - LD L,E
	{LdReg8, RegisterAddressing, 4, 1}, // 0x6C - LD L,H
	{LdReg8, RegisterAddressing, 4, 1}, // 0x6D - LD L,L
	{},                                 // 0x6E - LD L,(HL)
	{LdReg8, RegisterAddressing, 4, 1}, // 0x6F - LD L,A

	{},                                 // 0x70 - LD (HL),B
	{},                                 // 0x71 - LD (HL),C
	{},                                 // 0x72 - LD (HL),D
	{},                                 // 0x73 - LD (HL),E
	{},                                 // 0x74 - LD (HL),H
	{},                                 // 0x75 - LD (HL),L
	{Halt, ImpliedAddressing, 4, 1},    // 0x76 - HALT
	{},                                 // 0x77 - LD (HL),A
	{LdReg8, RegisterAddressing, 4, 1}, // 0x78 - LD A,B
	{LdReg8, RegisterAddressing, 4, 1}, // 0x79 - LD A,C
	{LdReg8, RegisterAddressing, 4, 1}, // 0x7A - LD A,D
	{LdReg8, RegisterAddressing, 4, 1}, // 0x7B - LD A,E
	{LdReg8, RegisterAddressing, 4, 1}, // 0x7C - LD A,H
	{LdReg8, RegisterAddressing, 4, 1}, // 0x7D - LD A,L
	{},                                 // 0x7E - LD A,(HL)
	{LdReg8, RegisterAddressing, 4, 1}, // 0x7F - LD A,A

	// 0x80-0xBF: Arithmetic and logical operations
	{AddA, RegisterAddressing, 4, 1}, // 0x80 - ADD A,B
	{AddA, RegisterAddressing, 4, 1}, // 0x81 - ADD A,C
	{AddA, RegisterAddressing, 4, 1}, // 0x82 - ADD A,D
	{AddA, RegisterAddressing, 4, 1}, // 0x83 - ADD A,E
	{AddA, RegisterAddressing, 4, 1}, // 0x84 - ADD A,H
	{AddA, RegisterAddressing, 4, 1}, // 0x85 - ADD A,L
	{},                               // 0x86 - ADD A,(HL)
	{AddA, RegisterAddressing, 4, 1}, // 0x87 - ADD A,A
	{},                               // 0x88 - ADC A,B
	{},                               // 0x89 - ADC A,C
	{},                               // 0x8A - ADC A,D
	{},                               // 0x8B - ADC A,E
	{},                               // 0x8C - ADC A,H
	{},                               // 0x8D - ADC A,L
	{},                               // 0x8E - ADC A,(HL)
	{},                               // 0x8F - ADC A,A

	{SubA, RegisterAddressing, 4, 1}, // 0x90 - SUB B
	{SubA, RegisterAddressing, 4, 1}, // 0x91 - SUB C
	{SubA, RegisterAddressing, 4, 1}, // 0x92 - SUB D
	{SubA, RegisterAddressing, 4, 1}, // 0x93 - SUB E
	{SubA, RegisterAddressing, 4, 1}, // 0x94 - SUB H
	{SubA, RegisterAddressing, 4, 1}, // 0x95 - SUB L
	{},                               // 0x96 - SUB (HL)
	{SubA, RegisterAddressing, 4, 1}, // 0x97 - SUB A
	{},                               // 0x98 - SBC A,B
	{},                               // 0x99 - SBC A,C
	{},                               // 0x9A - SBC A,D
	{},                               // 0x9B - SBC A,E
	{},                               // 0x9C - SBC A,H
	{},                               // 0x9D - SBC A,L
	{},                               // 0x9E - SBC A,(HL)
	{},                               // 0x9F - SBC A,A

	{AndA, RegisterAddressing, 4, 1}, // 0xA0 - AND B
	{AndA, RegisterAddressing, 4, 1}, // 0xA1 - AND C
	{AndA, RegisterAddressing, 4, 1}, // 0xA2 - AND D
	{AndA, RegisterAddressing, 4, 1}, // 0xA3 - AND E
	{AndA, RegisterAddressing, 4, 1}, // 0xA4 - AND H
	{AndA, RegisterAddressing, 4, 1}, // 0xA5 - AND L
	{},                               // 0xA6 - AND (HL)
	{AndA, RegisterAddressing, 4, 1}, // 0xA7 - AND A
	{XorA, RegisterAddressing, 4, 1}, // 0xA8 - XOR B
	{XorA, RegisterAddressing, 4, 1}, // 0xA9 - XOR C
	{XorA, RegisterAddressing, 4, 1}, // 0xAA - XOR D
	{XorA, RegisterAddressing, 4, 1}, // 0xAB - XOR E
	{XorA, RegisterAddressing, 4, 1}, // 0xAC - XOR H
	{XorA, RegisterAddressing, 4, 1}, // 0xAD - XOR L
	{},                               // 0xAE - XOR (HL)
	{XorA, RegisterAddressing, 4, 1}, // 0xAF - XOR A

	{OrA, RegisterAddressing, 4, 1}, // 0xB0 - OR B
	{OrA, RegisterAddressing, 4, 1}, // 0xB1 - OR C
	{OrA, RegisterAddressing, 4, 1}, // 0xB2 - OR D
	{OrA, RegisterAddressing, 4, 1}, // 0xB3 - OR E
	{OrA, RegisterAddressing, 4, 1}, // 0xB4 - OR H
	{OrA, RegisterAddressing, 4, 1}, // 0xB5 - OR L
	{},                              // 0xB6 - OR (HL)
	{OrA, RegisterAddressing, 4, 1}, // 0xB7 - OR A
	{CpA, RegisterAddressing, 4, 1}, // 0xB8 - CP B
	{CpA, RegisterAddressing, 4, 1}, // 0xB9 - CP C
	{CpA, RegisterAddressing, 4, 1}, // 0xBA - CP D
	{CpA, RegisterAddressing, 4, 1}, // 0xBB - CP E
	{CpA, RegisterAddressing, 4, 1}, // 0xBC - CP H
	{CpA, RegisterAddressing, 4, 1}, // 0xBD - CP L
	{},                              // 0xBE - CP (HL)
	{CpA, RegisterAddressing, 4, 1}, // 0xBF - CP A

	{},                                 // 0xC0 - RET NZ
	{},                                 // 0xC1 - POP BC
	{},                                 // 0xC2 - JP NZ,nn
	{JpAbs, ExtendedAddressing, 10, 3}, // 0xC3 - JP nn
	{},                                 // 0xC4 - CALL NZ,nn
	{},                                 // 0xC5 - PUSH BC
	{AddA, ImmediateAddressing, 7, 2},  // 0xC6 - ADD A,n
	{},                                 // 0xC7 - RST 00H
	{},                                 // 0xC8 - RET Z
	{},                                 // 0xC9 - RET
	{},                                 // 0xCA - JP Z,nn
	{},                                 // 0xCB - Extended instructions prefix
	{},                                 // 0xCC - CALL Z,nn
	{},                                 // 0xCD - CALL nn
	{},                                 // 0xCE - ADC A,n
	{},                                 // 0xCF - RST 08H

	{},                                // 0xD0 - RET NC
	{},                                // 0xD1 - POP DE
	{},                                // 0xD2 - JP NC,nn
	{},                                // 0xD3 - OUT (n),A
	{},                                // 0xD4 - CALL NC,nn
	{},                                // 0xD5 - PUSH DE
	{SubA, ImmediateAddressing, 7, 2}, // 0xD6 - SUB n
	{},                                // 0xD7 - RST 10H
	{},                                // 0xD8 - RET C
	{},                                // 0xD9 - EXX
	{},                                // 0xDA - JP C,nn
	{},                                // 0xDB - IN A,(n)
	{},                                // 0xDC - CALL C,nn
	{},                                // 0xDD - IX instructions prefix
	{},                                // 0xDE - SBC A,n
	{},                                // 0xDF - RST 18H

	{},                                // 0xE0 - RET PO
	{},                                // 0xE1 - POP HL
	{},                                // 0xE2 - JP PO,nn
	{},                                // 0xE3 - EX (SP),HL
	{},                                // 0xE4 - CALL PO,nn
	{},                                // 0xE5 - PUSH HL
	{AndA, ImmediateAddressing, 7, 2}, // 0xE6 - AND n
	{},                                // 0xE7 - RST 20H
	{},                                // 0xE8 - RET PE
	{},                                // 0xE9 - JP (HL)
	{},                                // 0xEA - JP PE,nn
	{},                                // 0xEB - EX DE,HL
	{},                                // 0xEC - CALL PE,nn
	{},                                // 0xED - Extended instructions prefix
	{XorA, ImmediateAddressing, 7, 2}, // 0xEE - XOR n
	{},                                // 0xEF - RST 28H

	{},                               // 0xF0 - RET P
	{},                               // 0xF1 - POP AF
	{},                               // 0xF2 - JP P,nn
	{},                               // 0xF3 - DI
	{},                               // 0xF4 - CALL P,nn
	{},                               // 0xF5 - PUSH AF
	{OrA, ImmediateAddressing, 7, 2}, // 0xF6 - OR n
	{},                               // 0xF7 - RST 30H
	{},                               // 0xF8 - RET M
	{},                               // 0xF9 - LD SP,HL
	{},                               // 0xFA - JP M,nn
	{},                               // 0xFB - EI
	{},                               // 0xFC - CALL M,nn
	{},                               // 0xFD - IY instructions prefix
	{CpA, ImmediateAddressing, 7, 2}, // 0xFE - CP n
	{},                               // 0xFF - RST 38H
}
