package z80

// This file contains patches for the most critical opcodes identified by our consistency tests
// This demonstrates the enhanced opcode structure without updating all 256 entries manually

func init() {
	// Patch the 15 problematic opcodes identified by our consistency tests
	// These were the RegisterIndirectAddressing opcodes that caused test failures

	// LD r,(HL) instructions
	Opcodes[0x46] = Opcode{LdReg8, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegB, RegNone} // LD B,(HL)
	Opcodes[0x4E] = Opcode{LdReg8, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegC, RegNone} // LD C,(HL)
	Opcodes[0x56] = Opcode{LdReg8, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegD, RegNone} // LD D,(HL)
	Opcodes[0x5E] = Opcode{LdReg8, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegE, RegNone} // LD E,(HL)
	Opcodes[0x66] = Opcode{LdReg8, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegH, RegNone} // LD H,(HL)
	Opcodes[0x6E] = Opcode{LdReg8, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegL, RegNone} // LD L,(HL)
	Opcodes[0x7E] = Opcode{LdReg8, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegA, RegNone} // LD A,(HL)

	// ALU operations with (HL)
	Opcodes[0x86] = Opcode{AddA, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegA, RegNone}    // ADD A,(HL)
	Opcodes[0x8E] = Opcode{AdcA, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegA, RegNone}    // ADC A,(HL)
	Opcodes[0x96] = Opcode{SubA, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegNone, RegNone} // SUB (HL)
	Opcodes[0x9E] = Opcode{SbcA, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegA, RegNone}    // SBC A,(HL)
	Opcodes[0xA6] = Opcode{AndA, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegNone, RegNone} // AND (HL)
	Opcodes[0xAE] = Opcode{XorA, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegNone, RegNone} // XOR (HL)
	Opcodes[0xB6] = Opcode{OrA, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegNone, RegNone}  // OR (HL)
	Opcodes[0xBE] = Opcode{CpA, RegisterIndirectAddressing, 7, 1, RegHLIndirect, RegNone, RegNone}  // CP (HL)

	// Sample of register-to-register LD instructions to show disambiguation
	Opcodes[0x40] = Opcode{LdReg8, RegisterAddressing, 4, 1, RegB, RegB, RegNone} // LD B,B
	Opcodes[0x41] = Opcode{LdReg8, RegisterAddressing, 4, 1, RegC, RegB, RegNone} // LD B,C
	Opcodes[0x42] = Opcode{LdReg8, RegisterAddressing, 4, 1, RegD, RegB, RegNone} // LD B,D
	Opcodes[0x43] = Opcode{LdReg8, RegisterAddressing, 4, 1, RegE, RegB, RegNone} // LD B,E
	Opcodes[0x47] = Opcode{LdReg8, RegisterAddressing, 4, 1, RegA, RegB, RegNone} // LD B,A
	Opcodes[0x7F] = Opcode{LdReg8, RegisterAddressing, 4, 1, RegA, RegA, RegNone} // LD A,A

	// Sample of other enhanced opcodes
	Opcodes[0x01] = Opcode{LdReg16, ImmediateAddressing, 10, 3, RegImm16, RegBC, RegNone} // LD BC,nn
	Opcodes[0x04] = Opcode{IncReg8, RegisterAddressing, 4, 1, RegNone, RegNone, RegB}     // INC B
	Opcodes[0x05] = Opcode{DecReg8, RegisterAddressing, 4, 1, RegNone, RegNone, RegB}     // DEC B
	Opcodes[0x06] = Opcode{LdImm8, ImmediateAddressing, 7, 2, RegImm8, RegB, RegNone}     // LD B,n
	Opcodes[0x09] = Opcode{AddHl, RegisterAddressing, 11, 1, RegBC, RegHL, RegNone}       // ADD HL,BC
	Opcodes[0xC7] = Opcode{Rst, ImpliedAddressing, 11, 1, RegNone, RegNone, RegRst00}     // RST 00H
	Opcodes[0xCF] = Opcode{Rst, ImpliedAddressing, 11, 1, RegNone, RegNone, RegRst08}     // RST 08H
}
