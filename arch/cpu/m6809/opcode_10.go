package m6809

// OpcodesPage2 is the page 2 ($10 prefix) opcode table for the 6809.
// Only defined entries are populated; all others have nil Instruction.
var OpcodesPage2 = [256]Opcode{
	// $10 $21 - $10 $2F: Long branches
	0x21: {Instruction: LbrnInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $1021 LBRN
	0x22: {Instruction: LbhiInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $1022 LBHI
	0x23: {Instruction: LblsInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $1023 LBLS
	0x24: {Instruction: LbccInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $1024 LBCC
	0x25: {Instruction: LbcsInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $1025 LBCS
	0x26: {Instruction: LbneInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $1026 LBNE
	0x27: {Instruction: LbeqInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $1027 LBEQ
	0x28: {Instruction: LbvcInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $1028 LBVC
	0x29: {Instruction: LbvsInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $1029 LBVS
	0x2A: {Instruction: LbplInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $102A LBPL
	0x2B: {Instruction: LbmiInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $102B LBMI
	0x2C: {Instruction: LbgeInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $102C LBGE
	0x2D: {Instruction: LbltInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $102D LBLT
	0x2E: {Instruction: LbgtInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $102E LBGT
	0x2F: {Instruction: LbleInst, Addressing: RelativeLongAddressing, Timing: 5, Size: 4}, // $102F LBLE

	// $10 $3F: SWI2
	0x3F: {Instruction: Swi2Inst, Addressing: ImpliedAddressing, Timing: 20, Size: 2}, // $103F SWI2

	// $10 $83 - $10 $BC: CMPD, CMPY
	0x83: {Instruction: CmpdInst, Addressing: Immediate16Addressing, Timing: 5, Size: 4}, // $1083 CMPD imm16
	0x8C: {Instruction: CmpyInst, Addressing: Immediate16Addressing, Timing: 5, Size: 4}, // $108C CMPY imm16
	0x8E: {Instruction: LdyInst, Addressing: Immediate16Addressing, Timing: 4, Size: 4},  // $108E LDY imm16

	0x93: {Instruction: CmpdInst, Addressing: DirectAddressing, Timing: 7, Size: 3}, // $1093 CMPD direct
	0x9C: {Instruction: CmpyInst, Addressing: DirectAddressing, Timing: 7, Size: 3}, // $109C CMPY direct
	0x9E: {Instruction: LdyInst, Addressing: DirectAddressing, Timing: 6, Size: 3},  // $109E LDY direct
	0x9F: {Instruction: StyInst, Addressing: DirectAddressing, Timing: 6, Size: 3},  // $109F STY direct

	0xA3: {Instruction: CmpdInst, Addressing: IndexedAddressing, Timing: 7, Size: 3}, // $10A3 CMPD indexed
	0xAC: {Instruction: CmpyInst, Addressing: IndexedAddressing, Timing: 7, Size: 3}, // $10AC CMPY indexed
	0xAE: {Instruction: LdyInst, Addressing: IndexedAddressing, Timing: 6, Size: 3},  // $10AE LDY indexed
	0xAF: {Instruction: StyInst, Addressing: IndexedAddressing, Timing: 6, Size: 3},  // $10AF STY indexed

	0xB3: {Instruction: CmpdInst, Addressing: ExtendedAddressing, Timing: 8, Size: 4}, // $10B3 CMPD extended
	0xBC: {Instruction: CmpyInst, Addressing: ExtendedAddressing, Timing: 8, Size: 4}, // $10BC CMPY extended
	0xBE: {Instruction: LdyInst, Addressing: ExtendedAddressing, Timing: 7, Size: 4},  // $10BE LDY extended
	0xBF: {Instruction: StyInst, Addressing: ExtendedAddressing, Timing: 7, Size: 4},  // $10BF STY extended

	// $10 $CE - $10 $FF: LDS, STS
	0xCE: {Instruction: LdsInst, Addressing: Immediate16Addressing, Timing: 4, Size: 4}, // $10CE LDS imm16
	0xDE: {Instruction: LdsInst, Addressing: DirectAddressing, Timing: 6, Size: 3},      // $10DE LDS direct
	0xDF: {Instruction: StsInst, Addressing: DirectAddressing, Timing: 6, Size: 3},      // $10DF STS direct
	0xEE: {Instruction: LdsInst, Addressing: IndexedAddressing, Timing: 6, Size: 3},     // $10EE LDS indexed
	0xEF: {Instruction: StsInst, Addressing: IndexedAddressing, Timing: 6, Size: 3},     // $10EF STS indexed
	0xFE: {Instruction: LdsInst, Addressing: ExtendedAddressing, Timing: 7, Size: 4},    // $10FE LDS extended
	0xFF: {Instruction: StsInst, Addressing: ExtendedAddressing, Timing: 7, Size: 4},    // $10FF STS extended
}
