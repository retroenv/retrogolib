# System Implementation Plan: Commodore 64

## Current Status

- **Status:** PLANNED
- **Last Updated:** 2026-03-11
- **Dependencies:** m6502 package (complete)

## Context

The Commodore 64 (1982) is the best-selling single personal computer model of all time,
with an estimated 12.5-17 million units sold. It has one of the largest retro software
libraries (>10,000 commercial titles) and the most active retro computing community,
with new software, demos, and hardware released regularly.

### Hardware Overview

| Component | Details |
|-----------|---------|
| CPU | MOS 6510 @ 1.023 MHz (NTSC) / 0.985 MHz (PAL) |
| RAM | 64 KB |
| ROM | 20 KB (BASIC $A000, KERNAL $E000, Character $D000) |
| Video | VIC-II (MOS 6567/6569), 320x200 / 160x200, 16 colors |
| Audio | SID (MOS 6581/8580), 3 voices + filter |
| I/O | CIA x2 (MOS 6526), keyboard, joysticks, serial, user port |
| Address bus | 16 bits (64 KB) |

### The MOS 6510

The 6510 is a MOS 6502 with one addition: a **built-in 6-bit bidirectional I/O port**
mapped to addresses $0000 (data direction register) and $0001 (port register). This I/O
port controls the C64's bank switching, selecting which combination of RAM, ROM, and I/O
chips are visible in the memory map.

| Feature | 6502 | 6510 |
|---------|------|------|
| Package | 40-pin | 40-pin |
| Instruction set | Full | **Identical** (including undocumented opcodes) |
| Address bus | 16 bits | 16 bits |
| I/O port | None | **6-bit at $0000-$0001** |
| IRQ/NMI | Yes | Yes |

The instruction set is byte-for-byte identical to the NMOS 6502, including all undocumented
opcodes. The only hardware difference is the built-in I/O port.

---

## Part 1: CPU Variant -- MOS 6510

### 1.1 Approach: Extend m6502 Package with Variant

Since the instruction set is identical, the 6510 is a variant of the existing m6502 package.
The I/O port at $0000-$0001 is handled by the Memory implementation, not the CPU.

### 1.2 The I/O Port

The 6510's I/O port uses two memory-mapped registers:

**$0000 -- Data Direction Register (DDR):**
Each bit controls whether the corresponding port bit is input (0) or output (1).
Default after reset: $2F (bits 0-3,5 are outputs; bits 4 is input; bits 6-7 unused).

**$0001 -- Port Register:**
| Bit | Name | Direction | Description |
|-----|------|-----------|-------------|
| 0 | LORAM | Output | BASIC ROM at $A000-$BFFF (1=ROM, 0=RAM) |
| 1 | HIRAM | Output | KERNAL ROM at $E000-$FFFF (1=ROM, 0=RAM) |
| 2 | CHAREN | Output | Character ROM at $D000-$DFFF (1=I/O, 0=Char ROM) |
| 3 | Cassette | Output | Cassette motor (0=on, 1=off) |
| 4 | Cassette | Input | Cassette switch sense (1=button pressed) |
| 5 | Cassette | Output | Cassette write line |
| 6-7 | - | - | Not connected |

Default after reset: $37 (BASIC, KERNAL, and I/O all visible).

The lower 3 bits control bank switching. All 8 combinations produce different memory maps
(see Part 2).

### 1.3 Implementation

#### Phase 1: Add Variant Constant

**Files to modify:**
- `arch/cpu/m6502/option.go` -- Add `Variant6510` constant

```go
const (
    VariantNMOS6502 CPUVariant = iota
    VariantNES6502
    Variant6507
    Variant6510    // MOS 6510: 6502 with built-in 6-bit I/O port at $0000-$0001
    Variant65C02
)
```

**Note:** The `Variant6510` is placed before `Variant65C02` to preserve the
`>= Variant65C02` opcode table selection. The 6510 uses the NMOS 6502 opcode table
(including undocumented opcodes).

#### Phase 2: Architecture Registration

**Files to modify:**
- `arch/arch.go` -- Add `M6510 Architecture = "6510"` (optional)

This is optional because the 6510 is instruction-identical to the 6502. However, having
the constant allows tools to distinguish the variant.

#### Phase 3: Testing

- Verify all existing 6502 tests pass with `Variant6510`
- The I/O port behavior is entirely in the Memory implementation (Part 2), so no
  CPU-level tests are needed beyond confirming the variant selects the correct opcode table

### 1.4 Estimated Effort

| Component | New LOC | Modified LOC |
|-----------|---------|-------------|
| Variant constant | ~5 | ~5 |
| Architecture constant | ~5 | ~5 |
| Tests | ~30 | ~0 |
| **Total** | **~40** | **~10** |

---

## Part 2: System Package -- Commodore 64

### 2.1 System Registration

**Files to modify:**
- `arch/system.go` -- Add `C64 System = "c64"`

### 2.2 Memory Map

The C64's memory map is controlled by the 6510 I/O port (bits 0-2 of $0001) and the
VIC-II/CIA chip select lines. The processor sees different devices at the same addresses
depending on the bank configuration.

**Default configuration ($0001 = $37, LORAM=1 HIRAM=1 CHAREN=1):**

| Address Range | Size | Contents |
|---------------|------|----------|
| $0000-$0001 | 2 B | 6510 I/O port (DDR + Port) |
| $0002-$00FF | 254 B | Zero page RAM |
| $0100-$01FF | 256 B | Stack (RAM) |
| $0200-$03FF | 512 B | Operating system work area |
| $0400-$07FF | 1 KB | Screen memory (default) |
| $0800-$9FFF | 38 KB | BASIC program area (RAM) |
| $A000-$BFFF | 8 KB | BASIC ROM |
| $C000-$CFFF | 4 KB | RAM |
| $D000-$D3FF | 1 KB | VIC-II registers |
| $D400-$D7FF | 1 KB | SID registers |
| $D800-$DBFF | 1 KB | Color RAM (4-bit, always visible) |
| $DC00-$DCFF | 256 B | CIA 1 registers |
| $DD00-$DDFF | 256 B | CIA 2 registers |
| $DE00-$DFFF | 512 B | I/O expansion area |
| $E000-$FFFF | 8 KB | KERNAL ROM |

**Bank switching configurations (bits 2-0 of $0001):**

| LORAM | HIRAM | CHAREN | $A000-$BFFF | $D000-$DFFF | $E000-$FFFF |
|-------|-------|--------|-------------|-------------|-------------|
| 1 | 1 | 1 | BASIC ROM | I/O chips | KERNAL ROM |
| 1 | 1 | 0 | BASIC ROM | Char ROM | KERNAL ROM |
| 1 | 0 | 1 | RAM | I/O chips | RAM |
| 1 | 0 | 0 | RAM | Char ROM | RAM |
| 0 | 1 | 1 | RAM | I/O chips | KERNAL ROM |
| 0 | 1 | 0 | RAM | Char ROM | KERNAL ROM |
| 0 | 0 | 1 | RAM | I/O chips | RAM |
| 0 | 0 | 0 | RAM | RAM | RAM |

**Important:** The VIC-II chip always sees RAM (not ROM) when reading memory for display.
The bank switching only affects the CPU's view. The VIC-II has its own 14-bit address space
(16 KB banks) selected by CIA 2 port A bits 0-1.

### 2.3 VIC-II Registers ($D000-$D3FF)

| Address | Name | Description |
|---------|------|-------------|
| $D000-$D00F | SP0X-SP7Y | Sprite 0-7 X/Y positions |
| $D010 | MSIGX | Sprite X position MSBs |
| $D011 | CR1 | Control register 1 (scroll Y, screen height, mode) |
| $D012 | RASTER | Raster counter (read) / Raster compare (write) |
| $D013-$D014 | LPX/LPY | Light pen X/Y |
| $D015 | SPENA | Sprite enable |
| $D016 | CR2 | Control register 2 (scroll X, screen width, multicolor) |
| $D017 | SPYEX | Sprite Y expansion |
| $D018 | VMCSB | Memory pointers (screen, character, bitmap base) |
| $D019 | IRQST | Interrupt status register |
| $D01A | IRQEN | Interrupt enable register |
| $D01B | SPDP | Sprite-data priority |
| $D01C | SPMC | Sprite multicolor |
| $D01D | SPXEX | Sprite X expansion |
| $D01E | SSCOL | Sprite-sprite collision |
| $D01F | SDCOL | Sprite-data collision |
| $D020 | BORDER | Border color |
| $D021 | BGCOL0 | Background color 0 |
| $D022-$D024 | BGCOL1-3 | Background colors 1-3 |
| $D025-$D026 | SPMCOL0-1 | Sprite multicolors |
| $D027-$D02E | SP0COL-SP7COL | Sprite 0-7 colors |

### 2.4 SID Registers ($D400-$D7FF)

| Address | Name | Description |
|---------|------|-------------|
| $D400-$D406 | V1FREQ-V1AD/SR | Voice 1: frequency, pulse width, control, envelope |
| $D407-$D40D | V2FREQ-V2AD/SR | Voice 2: frequency, pulse width, control, envelope |
| $D40E-$D414 | V3FREQ-V3AD/SR | Voice 3: frequency, pulse width, control, envelope |
| $D415-$D418 | FCLO-RESON | Filter: cutoff, resonance, mode, volume |
| $D419 | POTX | Paddle X (read) |
| $D41A | POTY | Paddle Y (read) |
| $D41B | RANDOM | Voice 3 oscillator output (read) |
| $D41C | ENV3 | Voice 3 envelope output (read) |

### 2.5 CIA Registers ($DC00-$DDFF)

Two identical CIA 6526 chips provide timers, I/O, and interrupt control.

**CIA 1 ($DC00-$DC0F) -- Keyboard, joystick, IRQ:**

| Address | Name | Description |
|---------|------|-------------|
| $DC00 | PRA | Port A: keyboard column / joystick 2 |
| $DC01 | PRB | Port B: keyboard row / joystick 1 |
| $DC02-$DC03 | DDRA/DDRB | Data direction registers |
| $DC04-$DC05 | TALO/TAHI | Timer A (16-bit) |
| $DC06-$DC07 | TBLO/TBHI | Timer B (16-bit) |
| $DC08-$DC0B | TOD | Time-of-Day clock (10ths, sec, min, hours BCD) |
| $DC0C | SDR | Serial data register |
| $DC0D | ICR | Interrupt control register |
| $DC0E | CRA | Control register A |
| $DC0F | CRB | Control register B |

**CIA 2 ($DD00-$DD0F) -- Serial bus, VIC bank, NMI:**
Same register layout as CIA 1, but:
- Port A bits 0-1: VIC-II bank select (inverted: %00=bank 3, %11=bank 0)
- Port A bits 2-7: Serial bus (IEC) and RS-232
- Interrupts trigger NMI instead of IRQ

### 2.6 Cartridge Format

C64 cartridge ROMs come in two common formats:

**PRG format (simplest):**
- 2-byte load address header (little-endian) followed by raw data
- Load address is where the data should be placed in memory
- Used for programs loaded from disk

**CRT format (cartridge images):**
- CHIP packets containing ROM chip data
- Supports bank switching for large cartridges (Ocean, EasyFlash, etc.)

| Offset | Size | Description |
|--------|------|-------------|
| $0000 | 16 B | Signature: "C64 CARTRIDGE   " |
| $0010 | 4 B | Header length |
| $0014 | 2 B | CRT version |
| $0016 | 2 B | Hardware type (cartridge mapper) |
| $0018 | 1 B | EXROM line |
| $0019 | 1 B | GAME line |
| $001A | 6 B | Reserved |
| $0020 | 32 B | Cartridge name |
| $0040+ | var | CHIP packets |

**CHIP packet:**

| Offset | Size | Description |
|--------|------|-------------|
| $0000 | 4 B | Signature: "CHIP" |
| $0004 | 4 B | Total packet length |
| $0008 | 2 B | Chip type (ROM/RAM/Flash) |
| $000A | 2 B | Bank number |
| $000C | 2 B | Load address |
| $000E | 2 B | ROM size |
| $0010+ | var | ROM data |

### 2.7 File Structure

```
arch/system/c64/
    doc.go              -- Package documentation
    c64.go              -- Memory map constants, bank switching table
    register/
        vic.go          -- VIC-II register addresses and names
        sid.go          -- SID register addresses and names
        cia.go          -- CIA register addresses and names
        io_port.go      -- 6510 I/O port bit definitions
    cartridge/
        prg.go          -- PRG file format (2-byte header + data)
        crt.go          -- CRT cartridge format parsing
```

---

## Part 3: Implementation Phases

### Phase 1: CPU Variant (6510)
- Add `Variant6510` to m6502 option.go
- Add `M6510` architecture constant (optional)
- Unit tests confirming correct opcode table selection

### Phase 2: System Registration
- Add `C64` system constant to `arch/system.go`
- Create `arch/system/c64/` package

### Phase 3: Memory Map and Bank Switching
- Define memory map address ranges and constants
- Document all 8 bank switching configurations
- Define 6510 I/O port bit constants ($0000-$0001)

### Phase 4: Hardware Registers
- VIC-II register definitions ($D000-$D02E)
- SID register definitions ($D400-$D41C)
- CIA 1/2 register definitions ($DC00-$DD0F)
- Color RAM address ($D800-$DBFF)

### Phase 5: Cartridge Support
- PRG file loading (2-byte header)
- CRT format parsing (header + CHIP packets)
- Hardware type (mapper) detection

### Phase 6: Testing
- Opcode execution tests with 6510 variant
- Bank switching configuration tests (all 8 modes)
- I/O port DDR and port register behavior tests
- Cartridge format parsing tests
- Register address completeness tests
- Lorenz test suite (C64-specific CPU test suite, tests undocumented opcodes and
  interrupt timing in the C64 environment)

---

## Part 4: Design Decisions

### 6510 as Variant vs Separate Package
**Decision: Variant within m6502**
- Rationale: The instruction set (including all undocumented opcodes) is byte-for-byte
  identical to the NMOS 6502. The only addition is the I/O port at $0000-$0001, which
  is properly handled by the Memory implementation. Creating a separate package would
  duplicate the entire 6502 codebase for a single hardware feature that lives outside
  the CPU instruction pipeline.

### I/O Port in Memory vs CPU
**Decision: I/O port handled by Memory implementation, not CPU**
- Rationale: The 6510 I/O port appears at memory addresses $0000-$0001 and is accessed
  via normal LDA/STA instructions -- there are no special I/O instructions. The Memory
  implementation intercepts reads/writes to these addresses and manages the DDR/port
  register state, including bank switching side effects. This keeps the CPU variant
  minimal (just a constant) while the system complexity lives in the system package.

### Bank Switching at Memory Level
**Decision: Memory.Read()/Memory.Write() handle bank visibility**
- Rationale: The C64's bank switching is controlled by bits 0-2 of the 6510 port register
  ($0001). When these bits change, the Memory implementation changes which backing
  store (RAM, ROM, or I/O) is visible at each address range. This is transparent to
  the CPU, which simply reads and writes addresses. The VIC-II's separate memory view
  (always sees RAM) is also handled at the Memory level.

---

## Part 5: Estimated Effort

| Component | New LOC |
|-----------|---------|
| CPU variant (6510 in m6502) | ~50 |
| System package (constants, memory map, bank switching) | ~400 |
| VIC-II register definitions | ~150 |
| SID register definitions | ~100 |
| CIA register definitions | ~150 |
| 6510 I/O port definitions | ~50 |
| PRG format support | ~80 |
| CRT cartridge format support | ~250 |
| Tests | ~500 |
| **Total** | **~1,730** |

---

## Part 6: References

- **VICE** (C): The Versatile Commodore Emulator. The gold standard C64 emulator with
  exceptional accuracy. Repository: `VICE-Team/svn-mirror`
- **C64 Programmer's Reference Guide** (Commodore, 1982): Official hardware and software
  reference. Available online.
- **Mapping the Commodore 64** (Sheldon Leemon): Complete memory map reference with every
  address documented.
- **C64 Wiki** (c64-wiki.com): Community-maintained technical reference.
- **Lorenz test suite**: CPU test suite specifically designed for the C64 environment,
  testing undocumented opcodes and interrupt timing.
- **Wolfgang Lorenz CPU test suite**: Tests 6510 behavior including undocumented opcodes
  and decimal mode edge cases.
