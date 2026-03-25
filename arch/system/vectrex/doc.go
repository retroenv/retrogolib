// Package vectrex provides constants and definitions for the GCE/Milton Bradley Vectrex system.
//
// The Vectrex uses a Motorola 6809 processor with a VIA (MC6522 Versatile Interface Adapter)
// for I/O, a DAC for vector display output, and an AY-3-8912 PSG for sound.
// The system has 1 KB of RAM and 8 KB of system ROM (executive/BIOS).
//
// Memory map (64 KB address space):
//
//	$0000-$7FFF  Cartridge ROM (32 KB window)
//	$C800-$CFFF  RAM (1 KB, mirrored)
//	$D000-$D7FF  VIA (MC6522) registers (mirrored)
//	$D800-$DFFF  VIA registers mirror
//	$E000-$FFFF  System ROM (8 KB executive/BIOS)
package vectrex
