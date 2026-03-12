// Package sm83 provides a Sharp SM83 (LR35902) CPU emulator for Game Boy systems.
//
// The SM83 is the CPU used in the Nintendo Game Boy and Game Boy Color. While derived
// from the Zilog Z80, it is a distinct processor that removes shadow registers, IX/IY
// index registers, I/O instructions, and DD/ED/FD opcode prefixes. It repurposes
// several Z80 opcodes for Game Boy-specific operations like LDH and STOP.
//
// Key differences from Z80:
//   - No shadow registers (AF', BC', DE', HL')
//   - No IX/IY index registers
//   - No I/O port instructions (IN/OUT)
//   - Only CB prefix (no DD, ED, FD prefixes)
//   - 4 flags only: Z (bit 7), N (bit 6), H (bit 5), C (bit 4)
//   - 4 condition codes: NZ, Z, NC, C (no PO, PE, P, M)
//   - SM83-unique instructions: STOP, SWAP, LDH, LD (HL+/-), ADD SP,e, LD HL,SP+e
//   - 11 illegal opcodes: $D3, $DB, $DD, $E3, $E4, $EB, $EC, $ED, $F4, $FC, $FD
//   - 5 fixed interrupt vectors ($0040, $0048, $0050, $0058, $0060)
//   - HALT bug: if HALT executed with IME=0 and pending interrupt, PC fails to increment
package sm83
