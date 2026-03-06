# Supported Systems

This document describes the CPU architectures and system-level support currently implemented in
retrogolib, along with notable gaps that would be valuable to add.

retrogolib is a library of building blocks — CPU emulators, cartridge formats, memory maps, and
register definitions. Full system emulation is implemented in separate repos that build on top of
this library.

---

## CPUs

### MOS Technology 6502 (`arch/cpu/m6502`)

Full cycle-accurate emulation of the original MOS 6502 processor.

- Complete instruction set with all addressing modes
- Accurate flag handling (N, V, Z, C, I, D, B)
- NMI and IRQ interrupt handling
- Unofficial/undocumented opcodes
- DMA stall cycle support
- Debugging and tracing

**65C02 variant** (`emulation_65c02.go`, `instruction_65c02.go`, `opcode_65c02.go`):
Extended instruction set of the CMOS successor (BRA, STZ, TRB, TSB, and additional addressing
modes for BIT, INC, DEC, JMP).

Target systems: NES, Apple II, Atari 8-bit, Commodore 64 (6510 variant), BBC Micro

---

### WDC 65C816 (`arch/cpu/m65816`)

Full emulation of the 16-bit successor to the 65C02.

- Dual-mode operation: Emulation mode (65C02-compatible) and Native mode
- 24-bit address space (16 MB)
- Mode-switchable register widths (M flag for accumulator, X flag for index registers)
- Direct Page register replacing the fixed zero page
- Data Bank and Program Bank registers for 24-bit addressing
- NMI and IRQ interrupt handling
- STP (stop) and WAI (wait) power-management instructions
- Debugging and tracing

Target systems: SNES/Super Famicom, Apple IIGS

---

### Zilog Z80 (`arch/cpu/z80`)

High-performance Z80 emulator with full instruction set support.

- Complete instruction set including undocumented opcodes
- All four prefix groups: base, CB, DD, FD (including undocumented DD/FD variants)
- Shadow register set (AF', BC', DE', HL')
- Index registers IX and IY with offset addressing
- Interrupt handling: NMI, maskable interrupts, modes 0/1/2
- Memory banking support
- Game Boy MBC1 memory controller (`memory_gameboy.go`)
- Cycle-accurate timing
- Debugging and tracing

Target systems: ZX Spectrum, MSX, ColecoVision, CP/M systems, Amstrad CPC, Game Boy (partial)

---

### Motorola 68000 (`arch/cpu/m68000`)

Full emulator for the 32-bit CISC processor with 16-bit external data bus.

- ~68 instruction mnemonics
- Hierarchical line-based opcode decoder (16 lines from top 4 bits)
- 14 addressing modes with effective address resolution
- 8 data registers (D0-D7, 32-bit) and 8 address registers (A0-A7)
- Dual stack pointers (USP/SSP) for user/supervisor modes
- 16-bit status register with CCR and system byte
- 256-vector exception model
- Big-endian byte order
- Cycle-accurate timing
- Debugging and tracing

Target systems: Sega Genesis/Mega Drive, Commodore Amiga, Atari ST, early Apple Macintosh,
Neo Geo

---

### CHIP-8 (`arch/cpu/chip8`)

Complete CHIP-8 virtual machine implementation.

- 16 general-purpose registers (V0-VF)
- 4 KB memory
- 64x32 monochrome display with XOR sprite drawing and collision detection
- 16-key hexadecimal keypad
- Sound and delay timers
- Built-in font set for hexadecimal digits (0-F)
- Comprehensive bounds checking

Target systems: COSMAC VIP, Telmac 1800, and derived hobby platforms

---

### Intel x86 (`arch/cpu/x86`)

Instruction set definitions for **static analysis** (disassemblers, assemblers).
This is not a runtime emulator.

- Opcode tables from 8086 through 80486
- Addressing mode metadata
- ModR/M byte parsing
- Instruction name and category lookup

Covered generations: 8086/8088, 80186, 80286, 80386, 80486

Target use: disassembly and assembler tooling for DOS-era PC software

---

## System Packages

### Nintendo Entertainment System (`arch/system/nes`)

System-level support for the NES.

- Memory map constants (code base, I/O register range, RAM end, name tables, palette size)
- iNES and NES 2.0 cartridge format parsing (`cartridge/`)
- Mirroring modes (horizontal, vertical, single-screen, four-screen)
- Code Data Log (CDL) format support for disassembler annotation
- Hardware register names and translations (`register/`)
- Parameter conversion and string representation utilities


---

## Missing Systems and CPUs

The following are notable gaps in the retrogolib foundation, grouped by priority based on
ecosystem size and tool-building demand.

### High-Value CPU Additions

| CPU | Used By | Notes |
|-----|---------|-------|
| Sharp LR35902 | Game Boy, Game Boy Color | Z80-derived with differences: no IX/IY, no shadow registers, added GB-specific opcodes (STOP, SWAP, RL/RR on (HL)), different flag behavior. The z80 package has `GameBoyMemory` (MBC1) but no GB-specific CPU. |
| Ricoh 2A03 / 2A07 | NES | 6502 without decimal mode (BCD), with built-in APU. The existing m6502 package is close but decimal mode behavior differs. |
| MOS 6510 | Commodore 64 | 6502 with a built-in 6-bit bidirectional I/O port at address $0000-$0001, used for bank switching. |
| HuC6280 | PC Engine / TurboGrafx-16 | 65C02 superset with additional instructions, 8 KB zero page, timer, and I/O. |
| SPC700 | SNES audio | The SNES audio coprocessor runs independently of the 65816. Required for full SNES emulation. |
| ARM7TDMI | Game Boy Advance | 32-bit ARM with Thumb mode. Very large target platform. |
| MIPS R3000A | PlayStation 1 | 32-bit MIPS I, little-endian. Large library of games. |

### High-Value System Package Additions

| System | CPU(s) | Notes |
|--------|--------|-------|
| Game Boy / Game Boy Color | Sharp LR35902 | The z80 package already provides `GameBoyMemory` with MBC1 banking. Missing: GB system package (memory map, I/O registers, cartridge format). |
| SNES / Super Famicom | 65816 + SPC700 | The 65816 CPU is already implemented. Missing: system package with memory map constants, cartridge format (SMC/SFC headers), and register definitions. |
| Sega Genesis / Mega Drive | 68000 + Z80 | Both CPUs are implemented. Missing: system package with the Genesis memory map, VDP registers, and ROM header parsing. |
| Commodore 64 | 6510 | Large software library, active retro scene. Needs 6510 CPU variant and C64 memory map/banking constants. |
| ZX Spectrum | Z80 | Z80 is already fully implemented. A system package (memory layout, ROM addresses, ULA I/O ports) would be a small addition with high value for the active Spectrum community. |
| Atari 2600 | 6507 (6502 variant) | 13-bit address bus (8 KB), no RAM on CPU side, TIA chip drives output. Unusual architecture. |
| MSX | Z80 | Z80 is implemented. MSX system package needs slot-based memory mapper and standard I/O port definitions. |
| Amstrad CPC | Z80 | Z80 is implemented. Needs Gate Array, CRTC, and memory banking constants. |
| Apple II | 6502 | 6502 is implemented. Needs soft switch addresses, slot memory map, and disk image format support. |

### x86 Runtime Emulation

The existing `arch/cpu/x86` package covers instruction definitions only. A runtime emulator for
8086/80286 would enable DOS-era tool development (debuggers, tracers, binary analysis). This is a
significantly larger effort than a typical retro CPU due to segmented memory and privilege levels.
