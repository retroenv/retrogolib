package x86

import "github.com/retroenv/retrogolib/set"

// MaxOpcodeSize is the maximum size of an opcode and its operands in bytes.
const MaxOpcodeSize = 6

// Opcode represents a CPU opcode with instruction information and addressing mode.
type Opcode struct {
	Instruction *Instruction   // pointer to instruction definition
	Addressing  AddressingMode // addressing mode used
	Timing      uint8          // execution time in cycles
	Size        uint8          // instruction size in bytes

	// Register disambiguation fields
	SrcRegister RegisterParam // source register
	DstRegister RegisterParam // destination register
	Register    RegisterParam // single register operand

	// ModR/M and displacement info
	HasModRM         bool  // instruction uses ModR/M byte
	HasDisplacement  bool  // instruction has displacement
	DisplacementSize uint8 // size of displacement (1 or 2 bytes)
}

// OpcodeInfo contains opcode and timing information for instruction variants.
type OpcodeInfo struct {
	Opcode   uint8 // primary opcode byte
	Size     uint8 // total instruction size in bytes
	Cycles   uint8 // execution cycles
	HasModRM bool  // uses ModR/M byte
}

// Opcodes maps the first opcode byte to CPU instruction information.
// Based on Intel 8086/8088 instruction set for DOS compatibility.
var Opcodes = [256]Opcode{
	// 0x00-0x0F: Basic arithmetic and data movement
	{Instruction: AddRMReg8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x00 ADD r/m8, r8
	{Instruction: AddRMReg16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x01 ADD r/m16, r16
	{Instruction: AddRegRM8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x02 ADD r8, r/m8
	{Instruction: AddRegRM16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x03 ADD r16, r/m16
	{Instruction: AddALImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},     // 0x04 ADD AL, imm8
	{Instruction: AddAXImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX},    // 0x05 ADD AX, imm16
	{Instruction: PushES, Addressing: ImpliedAddressing, Timing: 10, Size: 1, Register: RegES},         // 0x06 PUSH ES
	{Instruction: PopES, Addressing: ImpliedAddressing, Timing: 8, Size: 1, Register: RegES},           // 0x07 POP ES
	{Instruction: OrRMReg8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},   // 0x08 OR r/m8, r8
	{Instruction: OrRMReg16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x09 OR r/m16, r16
	{Instruction: OrRegRM8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},   // 0x0A OR r8, r/m8
	{Instruction: OrRegRM16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x0B OR r16, r/m16
	{Instruction: OrALImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},      // 0x0C OR AL, imm8
	{Instruction: OrAXImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX},     // 0x0D OR AX, imm16
	{Instruction: PushCS, Addressing: ImpliedAddressing, Timing: 10, Size: 1, Register: RegCS},         // 0x0E PUSH CS
	{Instruction: Undefined, Addressing: ImpliedAddressing, Timing: 1, Size: 1},                        // 0x0F (reserved)

	// 0x10-0x1F: More arithmetic operations
	{Instruction: AdcRMReg8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x10 ADC r/m8, r8
	{Instruction: AdcRMReg16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x11 ADC r/m16, r16
	{Instruction: AdcRegRM8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x12 ADC r8, r/m8
	{Instruction: AdcRegRM16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x13 ADC r16, r/m16
	{Instruction: AdcALImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},     // 0x14 ADC AL, imm8
	{Instruction: AdcAXImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX},    // 0x15 ADC AX, imm16
	{Instruction: PushSS, Addressing: ImpliedAddressing, Timing: 10, Size: 1, Register: RegSS},         // 0x16 PUSH SS
	{Instruction: PopSS, Addressing: ImpliedAddressing, Timing: 8, Size: 1, Register: RegSS},           // 0x17 POP SS
	{Instruction: SbbRMReg8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x18 SBB r/m8, r8
	{Instruction: SbbRMReg16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x19 SBB r/m16, r16
	{Instruction: SbbRegRM8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x1A SBB r8, r/m8
	{Instruction: SbbRegRM16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x1B SBB r16, r/m16
	{Instruction: SbbALImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},     // 0x1C SBB AL, imm8
	{Instruction: SbbAXImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX},    // 0x1D SBB AX, imm16
	{Instruction: PushDS, Addressing: ImpliedAddressing, Timing: 10, Size: 1, Register: RegDS},         // 0x1E PUSH DS
	{Instruction: PopDS, Addressing: ImpliedAddressing, Timing: 8, Size: 1, Register: RegDS},           // 0x1F POP DS

	// 0x20-0x2F: AND operations and segment prefixes
	{Instruction: AndRMReg8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x20 AND r/m8, r8
	{Instruction: AndRMReg16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x21 AND r/m16, r16
	{Instruction: AndRegRM8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x22 AND r8, r/m8
	{Instruction: AndRegRM16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x23 AND r16, r/m16
	{Instruction: AndALImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},     // 0x24 AND AL, imm8
	{Instruction: AndAXImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX},    // 0x25 AND AX, imm16
	{Instruction: SegES, Addressing: ImpliedAddressing, Timing: 2, Size: 1},                            // 0x26 ES: (segment prefix)
	{Instruction: Daa, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                              // 0x27 DAA
	{Instruction: SubRMReg8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x28 SUB r/m8, r8
	{Instruction: SubRMReg16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x29 SUB r/m16, r16
	{Instruction: SubRegRM8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x2A SUB r8, r/m8
	{Instruction: SubRegRM16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x2B SUB r16, r/m16
	{Instruction: SubALImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},     // 0x2C SUB AL, imm8
	{Instruction: SubAXImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX},    // 0x2D SUB AX, imm16
	{Instruction: SegCS, Addressing: ImpliedAddressing, Timing: 2, Size: 1},                            // 0x2E CS: (segment prefix)
	{Instruction: Das, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                              // 0x2F DAS

	// 0x30-0x3F: XOR operations and segment prefixes
	{Instruction: XorRMReg8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x30 XOR r/m8, r8
	{Instruction: XorRMReg16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x31 XOR r/m16, r16
	{Instruction: XorRegRM8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x32 XOR r8, r/m8
	{Instruction: XorRegRM16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x33 XOR r16, r/m16
	{Instruction: XorALImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},     // 0x34 XOR AL, imm8
	{Instruction: XorAXImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX},    // 0x35 XOR AX, imm16
	{Instruction: SegSS, Addressing: ImpliedAddressing, Timing: 2, Size: 1},                            // 0x36 SS: (segment prefix)
	{Instruction: Aaa, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                              // 0x37 AAA
	{Instruction: CmpRMReg8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x38 CMP r/m8, r8
	{Instruction: CmpRMReg16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x39 CMP r/m16, r16
	{Instruction: CmpRegRM8, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true},  // 0x3A CMP r8, r/m8
	{Instruction: CmpRegRM16, Addressing: ModRMRegisterAddressing, Timing: 3, Size: 2, HasModRM: true}, // 0x3B CMP r16, r/m16
	{Instruction: CmpALImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},     // 0x3C CMP AL, imm8
	{Instruction: CmpAXImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX},    // 0x3D CMP AX, imm16
	{Instruction: SegDS, Addressing: ImpliedAddressing, Timing: 2, Size: 1},                            // 0x3E DS: (segment prefix)
	{Instruction: Aas, Addressing: ImpliedAddressing, Timing: 4, Size: 1},                              // 0x3F AAS

	// 0x40-0x4F: INC/DEC register instructions
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegAX}, // 0x40 INC AX
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegCX}, // 0x41 INC CX
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegDX}, // 0x42 INC DX
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegBX}, // 0x43 INC BX
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegSP}, // 0x44 INC SP
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegBP}, // 0x45 INC BP
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegSI}, // 0x46 INC SI
	{Instruction: IncReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegDI}, // 0x47 INC DI
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegAX}, // 0x48 DEC AX
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegCX}, // 0x49 DEC CX
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegDX}, // 0x4A DEC DX
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegBX}, // 0x4B DEC BX
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegSP}, // 0x4C DEC SP
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegBP}, // 0x4D DEC BP
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegSI}, // 0x4E DEC SI
	{Instruction: DecReg16, Addressing: RegisterAddressing, Timing: 2, Size: 1, Register: RegDI}, // 0x4F DEC DI

	// 0x50-0x5F: PUSH/POP register instructions
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1, Register: RegAX}, // 0x50 PUSH AX
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1, Register: RegCX}, // 0x51 PUSH CX
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1, Register: RegDX}, // 0x52 PUSH DX
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1, Register: RegBX}, // 0x53 PUSH BX
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1, Register: RegSP}, // 0x54 PUSH SP
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1, Register: RegBP}, // 0x55 PUSH BP
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1, Register: RegSI}, // 0x56 PUSH SI
	{Instruction: PushReg16, Addressing: RegisterAddressing, Timing: 11, Size: 1, Register: RegDI}, // 0x57 PUSH DI
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegAX},   // 0x58 POP AX
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegCX},   // 0x59 POP CX
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegDX},   // 0x5A POP DX
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegBX},   // 0x5B POP BX
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegSP},   // 0x5C POP SP
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegBP},   // 0x5D POP BP
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegSI},   // 0x5E POP SI
	{Instruction: PopReg16, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegDI},   // 0x5F POP DI

	// 0x60-0x6F: More operations (8086 reserved, but we'll add some common ones)
	{}, {}, {}, {}, {}, {}, {}, {},
	{}, {}, {}, {}, {}, {}, {}, {},

	// 0x70-0x7F: Conditional jump instructions
	{Instruction: Jo, Addressing: RelativeAddressing, Timing: 16, Size: 2},   // 0x70 JO rel8
	{Instruction: Jno, Addressing: RelativeAddressing, Timing: 16, Size: 2},  // 0x71 JNO rel8
	{Instruction: Jb, Addressing: RelativeAddressing, Timing: 16, Size: 2},   // 0x72 JB/JNAE/JC rel8
	{Instruction: Jnb, Addressing: RelativeAddressing, Timing: 16, Size: 2},  // 0x73 JNB/JAE/JNC rel8
	{Instruction: Jz, Addressing: RelativeAddressing, Timing: 16, Size: 2},   // 0x74 JZ/JE rel8
	{Instruction: Jnz, Addressing: RelativeAddressing, Timing: 16, Size: 2},  // 0x75 JNZ/JNE rel8
	{Instruction: Jbe, Addressing: RelativeAddressing, Timing: 16, Size: 2},  // 0x76 JBE/JNA rel8
	{Instruction: Jnbe, Addressing: RelativeAddressing, Timing: 16, Size: 2}, // 0x77 JNBE/JA rel8
	{Instruction: Js, Addressing: RelativeAddressing, Timing: 16, Size: 2},   // 0x78 JS rel8
	{Instruction: Jns, Addressing: RelativeAddressing, Timing: 16, Size: 2},  // 0x79 JNS rel8
	{Instruction: Jp, Addressing: RelativeAddressing, Timing: 16, Size: 2},   // 0x7A JP/JPE rel8
	{Instruction: Jnp, Addressing: RelativeAddressing, Timing: 16, Size: 2},  // 0x7B JNP/JPO rel8
	{Instruction: Jl, Addressing: RelativeAddressing, Timing: 16, Size: 2},   // 0x7C JL/JNGE rel8
	{Instruction: Jnl, Addressing: RelativeAddressing, Timing: 16, Size: 2},  // 0x7D JNL/JGE rel8
	{Instruction: Jle, Addressing: RelativeAddressing, Timing: 16, Size: 2},  // 0x7E JLE/JNG rel8
	{Instruction: Jnle, Addressing: RelativeAddressing, Timing: 16, Size: 2}, // 0x7F JNLE/JG rel8

	// 0x80-0x8F: Group 1 ALU operations and MOV instructions
	{}, {}, {}, {}, {}, {}, {}, {}, // 0x80-0x87 (Group 1 ALU operations with ModR/M)
	{Instruction: MovRMReg8, Addressing: ModRMRegisterAddressing, Timing: 2, Size: 2, HasModRM: true},    // 0x88 MOV r/m8, r8
	{Instruction: MovRMReg16, Addressing: ModRMRegisterAddressing, Timing: 2, Size: 2, HasModRM: true},   // 0x89 MOV r/m16, r16
	{Instruction: MovRegRM8, Addressing: ModRMRegisterAddressing, Timing: 2, Size: 2, HasModRM: true},    // 0x8A MOV r8, r/m8
	{Instruction: MovRegRM16, Addressing: ModRMRegisterAddressing, Timing: 2, Size: 2, HasModRM: true},   // 0x8B MOV r16, r/m16
	{Instruction: MovMemImm16, Addressing: ModRMRegisterAddressing, Timing: 10, Size: 4, HasModRM: true}, // 0x8C MOV r/m16, Sreg
	{Instruction: Lea, Addressing: ModRMRegisterAddressing, Timing: 2, Size: 2, HasModRM: true},          // 0x8D LEA r16, m
	{Instruction: MovMemImm16, Addressing: ModRMRegisterAddressing, Timing: 2, Size: 2, HasModRM: true},  // 0x8E MOV Sreg, r/m16
	{}, // 0x8F Group 1A (POP r/m16)

	// 0x90-0x9F: NOP, XCHG, and other instructions
	{Instruction: Nop, Addressing: ImpliedAddressing, Timing: 3, Size: 1},                    // 0x90 NOP (XCHG AX, AX)
	{Instruction: Xchg, Addressing: RegisterAddressing, Timing: 3, Size: 1, Register: RegCX}, // 0x91 XCHG AX, CX
	{Instruction: Xchg, Addressing: RegisterAddressing, Timing: 3, Size: 1, Register: RegDX}, // 0x92 XCHG AX, DX
	{Instruction: Xchg, Addressing: RegisterAddressing, Timing: 3, Size: 1, Register: RegBX}, // 0x93 XCHG AX, BX
	{Instruction: Xchg, Addressing: RegisterAddressing, Timing: 3, Size: 1, Register: RegSP}, // 0x94 XCHG AX, SP
	{Instruction: Xchg, Addressing: RegisterAddressing, Timing: 3, Size: 1, Register: RegBP}, // 0x95 XCHG AX, BP
	{Instruction: Xchg, Addressing: RegisterAddressing, Timing: 3, Size: 1, Register: RegSI}, // 0x96 XCHG AX, SI
	{Instruction: Xchg, Addressing: RegisterAddressing, Timing: 3, Size: 1, Register: RegDI}, // 0x97 XCHG AX, DI
	{Instruction: Cbw, Addressing: ImpliedAddressing, Timing: 2, Size: 1},                    // 0x98 CBW
	{Instruction: Cwd, Addressing: ImpliedAddressing, Timing: 5, Size: 1},                    // 0x99 CWD
	{Instruction: CallFar, Addressing: SegmentOffsetAddressing, Timing: 28, Size: 5},         // 0x9A CALL ptr16:16
	{}, // 0x9B WAIT
	{}, // 0x9C PUSHF
	{}, // 0x9D POPF
	{}, // 0x9E SAHF
	{}, // 0x9F LAHF

	// 0xA0-0xAF: MOV direct, string operations
	{Instruction: MovRegImm8, Addressing: DirectAddressing, Timing: 10, Size: 3, Register: RegAL},  // 0xA0 MOV AL, moffs8
	{Instruction: MovRegImm16, Addressing: DirectAddressing, Timing: 10, Size: 3, Register: RegAX}, // 0xA1 MOV AX, moffs16
	{Instruction: MovMemImm8, Addressing: DirectAddressing, Timing: 10, Size: 3, Register: RegAL},  // 0xA2 MOV moffs8, AL
	{Instruction: MovMemImm16, Addressing: DirectAddressing, Timing: 10, Size: 3, Register: RegAX}, // 0xA3 MOV moffs16, AX
	{Instruction: Movsb, Addressing: StringAddressing, Timing: 18, Size: 1},                        // 0xA4 MOVSB
	{Instruction: Movsw, Addressing: StringAddressing, Timing: 18, Size: 1},                        // 0xA5 MOVSW
	{Instruction: Cmpsb, Addressing: StringAddressing, Timing: 22, Size: 1},                        // 0xA6 CMPSB
	{Instruction: Cmpsw, Addressing: StringAddressing, Timing: 22, Size: 1},                        // 0xA7 CMPSW
	{Instruction: Test, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},      // 0xA8 TEST AL, imm8
	{Instruction: Test, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX},      // 0xA9 TEST AX, imm16
	{Instruction: Stosb, Addressing: StringAddressing, Timing: 11, Size: 1},                        // 0xAA STOSB
	{Instruction: Stosw, Addressing: StringAddressing, Timing: 11, Size: 1},                        // 0xAB STOSW
	{Instruction: Lodsb, Addressing: StringAddressing, Timing: 12, Size: 1},                        // 0xAC LODSB
	{Instruction: Lodsw, Addressing: StringAddressing, Timing: 12, Size: 1},                        // 0xAD LODSW
	{Instruction: Scasb, Addressing: StringAddressing, Timing: 15, Size: 1},                        // 0xAE SCASB
	{Instruction: Scasw, Addressing: StringAddressing, Timing: 15, Size: 1},                        // 0xAF SCASW

	// 0xB0-0xBF: MOV immediate to register
	{Instruction: MovRegImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAL},  // 0xB0 MOV AL, imm8
	{Instruction: MovRegImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegCL},  // 0xB1 MOV CL, imm8
	{Instruction: MovRegImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegDL},  // 0xB2 MOV DL, imm8
	{Instruction: MovRegImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegBL},  // 0xB3 MOV BL, imm8
	{Instruction: MovRegImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegAH},  // 0xB4 MOV AH, imm8
	{Instruction: MovRegImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegCH},  // 0xB5 MOV CH, imm8
	{Instruction: MovRegImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegDH},  // 0xB6 MOV DH, imm8
	{Instruction: MovRegImm8, Addressing: ImmediateAddressing, Timing: 4, Size: 2, Register: RegBH},  // 0xB7 MOV BH, imm8
	{Instruction: MovRegImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegAX}, // 0xB8 MOV AX, imm16
	{Instruction: MovRegImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegCX}, // 0xB9 MOV CX, imm16
	{Instruction: MovRegImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegDX}, // 0xBA MOV DX, imm16
	{Instruction: MovRegImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegBX}, // 0xBB MOV BX, imm16
	{Instruction: MovRegImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegSP}, // 0xBC MOV SP, imm16
	{Instruction: MovRegImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegBP}, // 0xBD MOV BP, imm16
	{Instruction: MovRegImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegSI}, // 0xBE MOV SI, imm16
	{Instruction: MovRegImm16, Addressing: ImmediateAddressing, Timing: 4, Size: 3, Register: RegDI}, // 0xBF MOV DI, imm16

	// 0xC0-0xCF: Shifts, RET, LES/LDS, MOV immediate to memory, INT
	{}, {}, // 0xC0-0xC1 (Shift Group 2 - not in 8086)
	{Instruction: Ret, Addressing: ImmediateAddressing, Timing: 20, Size: 3}, // 0xC2 RET imm16
	{Instruction: Ret, Addressing: ImpliedAddressing, Timing: 16, Size: 1},   // 0xC3 RET
	{}, {}, // 0xC4-0xC5 LES, LDS
	{Instruction: MovMemImm8, Addressing: ModRMImmediateAddressing, Timing: 10, Size: 3, HasModRM: true},  // 0xC6 MOV r/m8, imm8
	{Instruction: MovMemImm16, Addressing: ModRMImmediateAddressing, Timing: 10, Size: 4, HasModRM: true}, // 0xC7 MOV r/m16, imm16
	{}, {}, // 0xC8-0xC9 ENTER, LEAVE (not in 8086)
	{Instruction: RetFar, Addressing: ImmediateAddressing, Timing: 25, Size: 3},                 // 0xCA RETF imm16
	{Instruction: RetFar, Addressing: ImpliedAddressing, Timing: 34, Size: 1},                   // 0xCB RETF
	{Instruction: Int, Addressing: ImmediateAddressing, Timing: 52, Size: 2, Register: RegImm8}, // 0xCC INT 3
	{Instruction: Int, Addressing: ImmediateAddressing, Timing: 51, Size: 2},                    // 0xCD INT imm8
	{Instruction: Into, Addressing: ImpliedAddressing, Timing: 53, Size: 1},                     // 0xCE INTO
	{Instruction: Iret, Addressing: ImpliedAddressing, Timing: 32, Size: 1},                     // 0xCF IRET

	// 0xD0-0xDF: Shift/rotate, AAM/AAD, XLAT
	{}, {}, {}, {}, // 0xD0-0xD3 (Shift Group 2)
	{}, {}, // 0xD4-0xD5 AAM, AAD
	{}, // 0xD6 (reserved)
	{Instruction: Xlat, Addressing: ImpliedAddressing, Timing: 11, Size: 1}, // 0xD7 XLAT
	{}, {}, {}, {}, {}, {}, {}, {}, // 0xD8-0xDF (FPU instructions)

	// 0xE0-0xEF: Loop, jump, I/O
	{}, {}, // 0xE0-0xE1 LOOPNZ, LOOPZ
	{}, // 0xE2 LOOP
	{}, // 0xE3 JCXZ
	{Instruction: In, Addressing: ImmediateAddressing, Timing: 10, Size: 2, Register: RegAL},  // 0xE4 IN AL, imm8
	{Instruction: In, Addressing: ImmediateAddressing, Timing: 10, Size: 2, Register: RegAX},  // 0xE5 IN AX, imm8
	{Instruction: Out, Addressing: ImmediateAddressing, Timing: 10, Size: 2, Register: RegAL}, // 0xE6 OUT imm8, AL
	{Instruction: Out, Addressing: ImmediateAddressing, Timing: 10, Size: 2, Register: RegAX}, // 0xE7 OUT imm8, AX
	{Instruction: Call, Addressing: RelativeAddressing, Timing: 19, Size: 3},                  // 0xE8 CALL rel16
	{Instruction: Jmp, Addressing: RelativeAddressing, Timing: 15, Size: 3},                   // 0xE9 JMP rel16
	{Instruction: JmpFar, Addressing: SegmentOffsetAddressing, Timing: 15, Size: 5},           // 0xEA JMP ptr16:16
	{Instruction: Jmp, Addressing: RelativeAddressing, Timing: 15, Size: 2},                   // 0xEB JMP rel8
	{Instruction: In, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegAL},    // 0xEC IN AL, DX
	{Instruction: In, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegAX},    // 0xED IN AX, DX
	{Instruction: Out, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegAL},   // 0xEE OUT DX, AL
	{Instruction: Out, Addressing: RegisterAddressing, Timing: 8, Size: 1, Register: RegAX},   // 0xEF OUT DX, AX

	// 0xF0-0xFF: Prefixes, TEST, control flags, INC/DEC memory, CALL/JMP indirect
	{}, // 0xF0 LOCK prefix
	{}, // 0xF1 (reserved)
	{Instruction: Repnz, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // 0xF2 REPNZ prefix
	{Instruction: Rep, Addressing: ImpliedAddressing, Timing: 2, Size: 1},   // 0xF3 REP/REPZ prefix
	{Instruction: Hlt, Addressing: ImpliedAddressing, Timing: 2, Size: 1},   // 0xF4 HLT
	{Instruction: Cmc, Addressing: ImpliedAddressing, Timing: 2, Size: 1},   // 0xF5 CMC
	{}, {}, // 0xF6-0xF7 (Group 3 - TEST, NOT, NEG, MUL, IMUL, DIV, IDIV)
	{Instruction: Clc, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // 0xF8 CLC
	{Instruction: Stc, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // 0xF9 STC
	{Instruction: Cli, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // 0xFA CLI
	{Instruction: Sti, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // 0xFB STI
	{Instruction: Cld, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // 0xFC CLD
	{Instruction: Std, Addressing: ImpliedAddressing, Timing: 2, Size: 1}, // 0xFD STD
	{}, {}, // 0xFE-0xFF (Group 4/5 - INC/DEC, CALL/JMP indirect)
}

// InstructionOpcodeMap provides reverse lookup from instruction to opcode.
var InstructionOpcodeMap = make(map[*Instruction][]uint8)

// RegisterOpcodeMap maps register-specific opcodes for faster lookup.
var RegisterOpcodeMap = make(map[RegisterParam]map[uint8]*Instruction)

// AddressingModeOpcodeMap maps addressing modes to their opcodes.
var AddressingModeOpcodeMap = make(map[AddressingMode]map[uint8]*Instruction)

// ValidOpcodes contains all valid opcode values.
var ValidOpcodes = set.New[uint8]()

// InitializeOpcodeMaps initializes the reverse lookup maps.
func InitializeOpcodeMaps() {
	for opcode, opcodeInfo := range Opcodes {
		if opcodeInfo.Instruction == nil {
			continue
		}

		opcodeValue := uint8(opcode)
		ValidOpcodes.Add(opcodeValue)

		// Build instruction to opcode map
		if InstructionOpcodeMap[opcodeInfo.Instruction] == nil {
			InstructionOpcodeMap[opcodeInfo.Instruction] = []uint8{}
		}
		InstructionOpcodeMap[opcodeInfo.Instruction] = append(
			InstructionOpcodeMap[opcodeInfo.Instruction], opcodeValue)

		// Build register to opcode map
		if opcodeInfo.Register != 0 {
			if RegisterOpcodeMap[opcodeInfo.Register] == nil {
				RegisterOpcodeMap[opcodeInfo.Register] = make(map[uint8]*Instruction)
			}
			RegisterOpcodeMap[opcodeInfo.Register][opcodeValue] = opcodeInfo.Instruction
		}

		// Build addressing mode to opcode map
		if AddressingModeOpcodeMap[opcodeInfo.Addressing] == nil {
			AddressingModeOpcodeMap[opcodeInfo.Addressing] = make(map[uint8]*Instruction)
		}
		AddressingModeOpcodeMap[opcodeInfo.Addressing][opcodeValue] = opcodeInfo.Instruction
	}
}

// GetOpcodeInfo returns the opcode information for a given opcode byte.
func GetOpcodeInfo(opcode uint8) (Opcode, bool) {
	if int(opcode) >= len(Opcodes) {
		return Opcode{}, false
	}

	opcodeInfo := Opcodes[opcode]
	if opcodeInfo.Instruction == nil {
		return Opcode{}, false
	}

	return opcodeInfo, true
}

// IsValidOpcode returns whether the given byte is a valid opcode.
func IsValidOpcode(opcode uint8) bool {
	return ValidOpcodes.Contains(opcode)
}

// GetInstructionOpcodes returns all opcode bytes for a given instruction.
func GetInstructionOpcodes(instruction *Instruction) []uint8 {
	if opcodes, exists := InstructionOpcodeMap[instruction]; exists {
		return opcodes
	}
	return nil
}

// GetRegisterOpcodes returns all opcodes that use a specific register.
func GetRegisterOpcodes(register RegisterParam) map[uint8]*Instruction {
	if opcodes, exists := RegisterOpcodeMap[register]; exists {
		return opcodes
	}
	return nil
}

// GetAddressingModeOpcodes returns all opcodes that use a specific addressing mode.
func GetAddressingModeOpcodes(mode AddressingMode) map[uint8]*Instruction {
	if opcodes, exists := AddressingModeOpcodeMap[mode]; exists {
		return opcodes
	}
	return nil
}
