package z80

// This file contains patches for the most critical opcodes identified by our consistency tests
// This demonstrates the enhanced opcode structure without updating all 256 entries manually

func init() {
	// Patch the 15 problematic opcodes identified by our consistency tests
	// These were the RegisterIndirectAddressing opcodes that caused test failures

	// LD r,(HL) instructions
	Opcodes[0x46] = Opcode{Instruction: LdReg8, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegB} // LD B,(HL)
	Opcodes[0x4E] = Opcode{Instruction: LdReg8, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegC} // LD C,(HL)
	Opcodes[0x56] = Opcode{Instruction: LdReg8, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegD} // LD D,(HL)
	Opcodes[0x5E] = Opcode{Instruction: LdReg8, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegE} // LD E,(HL)
	Opcodes[0x66] = Opcode{Instruction: LdReg8, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegH} // LD H,(HL)
	Opcodes[0x6E] = Opcode{Instruction: LdReg8, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegL} // LD L,(HL)
	Opcodes[0x7E] = Opcode{Instruction: LdReg8, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegA} // LD A,(HL)

	// ALU operations with (HL)
	Opcodes[0x86] = Opcode{Instruction: AddA, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegA} // ADD A,(HL)
	Opcodes[0x8E] = Opcode{Instruction: AdcA, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegA} // ADC A,(HL)
	Opcodes[0x96] = Opcode{Instruction: SubA, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect}                    // SUB (HL)
	Opcodes[0x9E] = Opcode{Instruction: SbcA, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect, DstRegister: RegA} // SBC A,(HL)
	Opcodes[0xA6] = Opcode{Instruction: AndA, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect}                    // AND (HL)
	Opcodes[0xAE] = Opcode{Instruction: XorA, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect}                    // XOR (HL)
	Opcodes[0xB6] = Opcode{Instruction: OrA, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect}                     // OR (HL)
	Opcodes[0xBE] = Opcode{Instruction: CpA, Addressing: RegisterIndirectAddressing, Timing: 7, Size: 1, SrcRegister: RegHLIndirect}                     // CP (HL)

	// Sample of register-to-register LD instructions to show disambiguation
	Opcodes[0x40] = Opcode{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1, SrcRegister: RegB, DstRegister: RegB} // LD B,B
	Opcodes[0x41] = Opcode{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1, SrcRegister: RegC, DstRegister: RegB} // LD B,C
	Opcodes[0x42] = Opcode{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1, SrcRegister: RegD, DstRegister: RegB} // LD B,D
	Opcodes[0x43] = Opcode{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1, SrcRegister: RegE, DstRegister: RegB} // LD B,E
	Opcodes[0x47] = Opcode{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1, SrcRegister: RegA, DstRegister: RegB} // LD B,A
	Opcodes[0x7F] = Opcode{Instruction: LdReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1, SrcRegister: RegA, DstRegister: RegA} // LD A,A

	// Sample of other enhanced opcodes
	Opcodes[0x01] = Opcode{Instruction: LdReg16, Addressing: ImmediateAddressing, Timing: 10, Size: 3, SrcRegister: RegImm16, DstRegister: RegBC} // LD BC,nn
	Opcodes[0x04] = Opcode{Instruction: IncReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1, Register: RegB}                              // INC B
	Opcodes[0x05] = Opcode{Instruction: DecReg8, Addressing: RegisterAddressing, Timing: 4, Size: 1, Register: RegB}                              // DEC B
	Opcodes[0x06] = Opcode{Instruction: LdImm8, Addressing: ImmediateAddressing, Timing: 7, Size: 2, SrcRegister: RegImm8, DstRegister: RegB}     // LD B,n
	Opcodes[0x09] = Opcode{Instruction: AddHl, Addressing: RegisterAddressing, Timing: 11, Size: 1, SrcRegister: RegBC, DstRegister: RegHL}       // ADD HL,BC
	Opcodes[0xC7] = Opcode{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1, Register: RegRst00}                              // RST 00H
	Opcodes[0xCF] = Opcode{Instruction: Rst, Addressing: ImpliedAddressing, Timing: 11, Size: 1, Register: RegRst08}                              // RST 08H
}
