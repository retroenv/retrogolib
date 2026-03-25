package m6809

// OpcodesPage3 is the page 3 ($11 prefix) opcode table for the 6809.
// Only defined entries are populated; all others have nil Instruction.
var OpcodesPage3 = [256]Opcode{
	// $11 $3F: SWI3
	0x3F: {Instruction: Swi3Inst, Addressing: ImpliedAddressing, Timing: 20, Size: 2}, // $113F SWI3

	// $11 $83 - $11 $BC: CMPU, CMPS
	0x83: {Instruction: CmpuInst, Addressing: Immediate16Addressing, Timing: 5, Size: 4}, // $1183 CMPU imm16
	0x8C: {Instruction: CmpsInst, Addressing: Immediate16Addressing, Timing: 5, Size: 4}, // $118C CMPS imm16

	0x93: {Instruction: CmpuInst, Addressing: DirectAddressing, Timing: 7, Size: 3}, // $1193 CMPU direct
	0x9C: {Instruction: CmpsInst, Addressing: DirectAddressing, Timing: 7, Size: 3}, // $119C CMPS direct

	0xA3: {Instruction: CmpuInst, Addressing: IndexedAddressing, Timing: 7, Size: 3}, // $11A3 CMPU indexed
	0xAC: {Instruction: CmpsInst, Addressing: IndexedAddressing, Timing: 7, Size: 3}, // $11AC CMPS indexed

	0xB3: {Instruction: CmpuInst, Addressing: ExtendedAddressing, Timing: 8, Size: 4}, // $11B3 CMPU extended
	0xBC: {Instruction: CmpsInst, Addressing: ExtendedAddressing, Timing: 8, Size: 4}, // $11BC CMPS extended
}
