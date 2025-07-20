// Package gameboy provides constants, memory mappings, and utilities for
// Game Boy (DMG) system emulation.
//
// The Game Boy uses a modified Z80 CPU (often called the SM83 or GBZ80) that
// lacks some Z80 features but adds Game Boy-specific functionality. This package
// provides the memory map, I/O registers, and system constants needed for
// Game Boy emulation.
//
// Key features:
//   - Complete memory map definitions for Game Boy DMG
//   - I/O register constants and bit masks
//   - Cartridge header structure and MBC support
//   - Audio, video, and timer register definitions
//   - Interrupt vector and control definitions
//
// Memory Layout:
//
//	0x0000-0x3FFF: ROM Bank 0 (16KB)
//	0x4000-0x7FFF: ROM Bank 1-N (16KB, switchable)
//	0x8000-0x9FFF: Video RAM (8KB)
//	0xA000-0xBFFF: External RAM (8KB, switchable)
//	0xC000-0xDFFF: Work RAM (8KB)
//	0xE000-0xFDFF: Echo RAM (mirror of 0xC000-0xDDFF)
//	0xFE00-0xFE9F: Object Attribute Memory (OAM, 160 bytes)
//	0xFEA0-0xFEFF: Not usable
//	0xFF00-0xFF7F: I/O Registers
//	0xFF80-0xFFFE: High RAM (127 bytes)
//	0xFFFF: Interrupt Enable Register
package gameboy
