// Package atari2600 provides constants and definitions for the Atari 2600 (VCS) system.
//
// The Atari 2600 uses a MOS 6507 processor (a pin-reduced 6502 with a 13-bit address bus)
// combined with a TIA (Television Interface Adapter) for video/audio and a RIOT (6532) chip
// for RAM, timers, and I/O. The system has 128 bytes of RAM and no frame buffer -- the
// programmer must feed graphics data to the TIA in real time ("racing the beam").
//
// Memory map (13-bit, 8 KB address space with mirrors):
//
//	$0000-$002C  TIA write registers (player, missile, ball, playfield, audio)
//	$0030-$003D  TIA read registers (collisions, inputs, timing)
//	$0080-$00FF  RIOT RAM (128 bytes)
//	$0280-$0297  RIOT I/O and timer registers
//	$1000-$1FFF  Cartridge ROM (4 KB window, bank-switched for larger ROMs)
package atari2600
