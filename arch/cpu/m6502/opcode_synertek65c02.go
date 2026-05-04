// This file contains the Synertek 65C02 opcode table.
//
// The Synertek 65C02 implements the core 65C02 instruction set but does NOT
// include the Rockwell extensions (RMB/SMB/BBR/BBS). Those opcodes are treated
// as NOPs of the same byte length: 2-byte for RMBx/SMBx, 3-byte for BBRx/BBSx.

package m6502

// OpcodesSynertek65C02 is the opcode table for the Synertek 65C02 variant.
// Identical to Opcodes65C02 except that the 32 Rockwell bit-manipulation
// opcodes are replaced with size-matched NOPs.
var OpcodesSynertek65C02 = func() [256]Opcode {
	table := Opcodes65C02

	// RMB0-7: opcodes 0x07, 0x17, 0x27, 0x37, 0x47, 0x57, 0x67, 0x77 (2-byte NOPs)
	// SMB0-7: opcodes 0x87, 0x97, 0xa7, 0xb7, 0xc7, 0xd7, 0xe7, 0xf7 (2-byte NOPs)
	for _, op := range []uint8{0x07, 0x17, 0x27, 0x37, 0x47, 0x57, 0x67, 0x77,
		0x87, 0x97, 0xa7, 0xb7, 0xc7, 0xd7, 0xe7, 0xf7} {
		table[op] = Opcode{Instruction: Nop65C02Inst, Addressing: ZeroPageAddressing, Timing: 3}
	}

	// BBR0-7: opcodes 0x0f, 0x1f, 0x2f, 0x3f, 0x4f, 0x5f, 0x6f, 0x7f (3-byte NOPs)
	// BBS0-7: opcodes 0x8f, 0x9f, 0xaf, 0xbf, 0xcf, 0xdf, 0xef, 0xff (3-byte NOPs)
	for _, op := range []uint8{0x0f, 0x1f, 0x2f, 0x3f, 0x4f, 0x5f, 0x6f, 0x7f,
		0x8f, 0x9f, 0xaf, 0xbf, 0xcf, 0xdf, 0xef, 0xff} {
		table[op] = Opcode{Instruction: Nop65C02Inst, Addressing: ZeroPageRelativeAddressing, Timing: 3}
	}

	return table
}()
