# System Implementation Plan: Commodore 64

## Current Status

- **Status:** System package planned; CPU variant available.
- **Last Updated:** 2026-10-02
- **Dependencies:** `cpu6502.WithVariant(cpu6502.Variant6510)` is available.
  C64 system registration, memory mapping, and device behavior remain unimplemented.

The variant is defined in [option.go](../arch/cpu/cpu6502/option.go).
[step.go](../arch/cpu/cpu6502/step.go) selects the NMOS opcode table, and
[instruction_registry.go](../arch/cpu/cpu6502/instruction_registry.go) selects
the NMOS instruction registry. These paths do not implement the 6510 I/O port.
There is no C64 entry in [system.go](../arch/system.go).

## Context

This plan covers C64 system definitions, memory mapping, and file loading.
VIC-II rendering, SID audio, CIA execution, and peripheral emulation need separate
implementation plans. Register definitions alone do not implement these devices.

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

The 6510 uses the NMOS 6502 instruction set and has a **6-bit bidirectional I/O port**
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

Reuse the NMOS instruction implementation. C64 bus sharing and the processor
port still require system behavior beyond opcode selection.

---

## Part 1: CPU Variant -- MOS 6510

### 1.1 Approach: Use the Existing cpu6502 Variant

Use `WithVariant(Variant6510)`. Implement the I/O port at $0000-$0001 in
the C64 implementation of `cpu6502.BasicMemory`.

### 1.2 The I/O Port

The 6510's I/O port uses two memory-mapped registers:

**$0000 -- Data Direction Register (DDR):**
Each bit controls whether the corresponding port bit is input (0) or output (1).
Do not use the usual software value $2F as the hardware reset value. At reset,
the port pins are inputs. Model external pull-ups separately from the output
latch, as in [VICE's processor-port implementation](https://github.com/VICE-Team/svn-mirror/blob/main/vice/src/c64/c64pla.c).

**$0001 -- Port Register:**

| Bit | Name | Direction | Description |
|-----|------|-----------|-------------|
| 0 | LORAM | Output | Selects memory mapping with HIRAM and CHAREN; see the table below |
| 1 | HIRAM | Output | Selects memory mapping with LORAM and CHAREN |
| 2 | CHAREN | Output | Selects I/O or character ROM when LORAM or HIRAM is high |
| 3 | Cassette | Output | Cassette write data |
| 4 | Cassette | Input | Cassette switch sense (0=button pressed) |
| 5 | Cassette | Output | Cassette motor (0=on, 1=off) |
| 6-7 | - | - | Not connected |

The directions above describe normal software use. DDR bits select the actual
directions. The [Commodore service manual](https://zimmers.net/anonftp/pub/cbm/schematics/computers/c64/manual-html/Page_11.html)
defines the cassette signals; VICE also models the active-low sense input.

The usual software value $37 selects BASIC, KERNAL, and I/O with the normal
DDR configuration. The effective pin levels select the memory map. With input
pins, the external signals and pull-ups matter. The eight combinations produce
seven distinct CPU read maps without a cartridge.

### 1.3 Remaining CPU Integration

- Add focused tests for 6510 opcode selection, decimal arithmetic, and interrupts.
- Test CPU access to the system port and memory map through `BasicMemory`.
- Consider an optional `CPU6510` architecture identifier only if tooling needs it.
  If added, update architecture validation and registration tests as well.
- Use the existing `WithCycleHook` API when device timing is implemented.
  Its read-cycle DMA support needs C64-specific validation before use for VIC-II
  bus sharing. See [cycle.go](../arch/cpu/cpu6502/cycle.go).

---

## Part 2: System Package -- Commodore 64

### 2.1 System Registration

**Files to modify:**

- `arch/system.go` -- Add `C64 System = "c64"` and include it in `allSupportedSystems`.
- `arch/system_test.go` -- Test parsing, validation, and supported-system listing.

### 2.2 Memory Map

The PLA selects CPU memory from the processor-port signals and the cartridge
GAME/EXROM lines. The table below assumes no cartridge: GAME=1 and EXROM=1.
Cartridge modes need additional maps.

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
| $D800-$DBFF | 1 KB | Color RAM (4-bit; CPU access requires I/O to be selected) |
| $DC00-$DCFF | 256 B | CIA 1 registers |
| $DD00-$DDFF | 256 B | CIA 2 registers |
| $DE00-$DFFF | 512 B | I/O expansion area |
| $E000-$FFFF | 8 KB | KERNAL ROM |

**CPU read maps from effective LORAM, HIRAM, and CHAREN pin levels:**

| LORAM | HIRAM | CHAREN | $A000-$BFFF | $D000-$DFFF | $E000-$FFFF |
|-------|-------|--------|-------------|-------------|-------------|
| 1 | 1 | 1 | BASIC ROM | I/O chips | KERNAL ROM |
| 1 | 1 | 0 | BASIC ROM | Char ROM | KERNAL ROM |
| 1 | 0 | 1 | RAM | I/O chips | RAM |
| 1 | 0 | 0 | RAM | Char ROM | RAM |
| 0 | 1 | 1 | RAM | I/O chips | KERNAL ROM |
| 0 | 1 | 0 | RAM | Char ROM | KERNAL ROM |
| 0 | 0 | 1 | RAM | RAM | RAM |
| 0 | 0 | 0 | RAM | RAM | RAM |

These maps agree with [VICE's CPU memory tables](https://github.com/VICE-Team/svn-mirror/blob/main/vice/src/c64/c64meminit.c).
In these modes, CPU writes under BASIC, KERNAL, or character ROM reach RAM.
When I/O is selected, writes in $D000-$DFFF reach the selected I/O device.

The VIC-II uses a separate memory view. CIA 2 port A bits 0-1 select a 16 KB
bank. With no cartridge, character ROM replaces RAM at physical addresses
$1000-$1FFF and $9000-$9FFF in this view. CPU bank switching does not remove
these ROM windows. See the [chips C64 memory implementation](https://github.com/floooh/chips/blob/master/systems/c64.h).

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

### 2.6 Program and Cartridge Formats

PRG stores a loadable program. CRT stores a cartridge image and its metadata.

**PRG format (simplest):**

- 2-byte load address header (little-endian) followed by raw data
- Load address is where the data should be placed in memory
- Used for programs loaded from disk

**CRT format (cartridge images):**

- CHIP packets containing ROM chip data
- Supports bank switching for large cartridges (Ocean, EasyFlash, etc.)

CRT numeric fields with more than one byte use big-endian order. Read the
header length to locate the first CHIP packet; $0040 is the minimum header
length. Follow the [VICE CRT specification](https://vice-emu.sourceforge.io/vice_17.html).

| Offset | Size | Description |
|--------|------|-------------|
| $0000 | 16 B | Signature: "C64 CARTRIDGE   " |
| $0010 | 4 B | Header length |
| $0014 | 2 B | CRT version |
| $0016 | 2 B | Hardware type (cartridge mapper) |
| $0018 | 1 B | EXROM line |
| $0019 | 1 B | GAME line |
| $001A | 1 B | Hardware revision/subtype in CRT version 1.1 and later |
| $001B | 5 B | Reserved |
| $0020 | 32 B | Cartridge name |
| Header length | Variable | CHIP packets |

**CHIP packet:**

| Offset | Size | Description |
|--------|------|-------------|
| $0000 | 4 B | Signature: "CHIP" |
| $0004 | 4 B | Total packet length |
| $0008 | 2 B | Chip type (ROM/RAM/Flash/EEPROM) |
| $000A | 2 B | Bank number |
| $000C | 2 B | Load address |
| $000E | 2 B | ROM size |
| $0010+ | var | ROM data |

Check header and packet bounds before reading data. Reject truncated files,
invalid packet lengths, and unsupported hardware types. CRT parsing alone
does not implement cartridge bank switching.

### 2.7 File Structure

```
arch/system/c64/
    doc.go              -- Package documentation
    c64.go              -- Memory map constants, bank switching table
    memory.go           -- CPU memory access, processor port, and VIC-II memory view
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

### Phase 1: CPU Integration (6510)

- Use the existing `Variant6510` and add the tests listed in Part 1.3.
- Decide whether tooling needs a separate architecture identifier.

### Phase 2: System Registration

- Add `C64` to the system constants and `allSupportedSystems`; test registration.
- Create `arch/system/c64/` package

### Phase 3: Memory Map and Bank Switching

- Define memory map address ranges and constants
- Implement all eight no-cartridge CPU configurations, including writes under ROM.
- Keep the VIC-II memory view separate and include its character-ROM windows.
- Define 6510 I/O port bit constants ($0000-$0001)

### Phase 4: Hardware Registers

- VIC-II register definitions ($D000-$D02E)
- SID register definitions ($D400-$D41C)
- CIA 1/2 register definitions ($DC00-$DD0F)
- Color RAM address ($D800-$DBFF)

### Phase 5: Program and Cartridge Support

- PRG file loading (2-byte header)
- CRT format parsing (header + CHIP packets)
- Hardware type (mapper) detection

### Phase 6: Testing

- Opcode execution tests with 6510 variant
- Bank switching tests for all eight no-cartridge configurations and supported cartridge modes
- Tests for writes under ROM, I/O visibility, and the separate VIC-II memory view
- I/O port DDR and port register behavior tests
- Cartridge format parsing tests
- Register address completeness tests
- Select and pin the Wolfgang Lorenz test suite when the required devices work.
  Record the selected tests and their system dependencies before claiming conformance.

---

## Part 4: Design Decisions

### 6510 as Variant vs Separate Package

**Decision: Variant within cpu6502**

- Reuse the existing NMOS implementation and `Variant6510`. Keep port and system
  behavior in the C64 package. A separate CPU package would duplicate instruction code.

### I/O Port in Memory vs CPU

**Decision: Implement the port in the C64 `BasicMemory` implementation.**

- Handle DDR, output latch, external inputs, and effective pin levels separately.
  Verify read and write behavior through CPU instructions.

### Bank Switching at Memory Level

**Decision: `BasicMemory.Read` and `BasicMemory.Write` select CPU memory.**

- Use effective processor-port signals and cartridge lines to select the map.
  Reads and writes can select different storage at the same address.
  Give VIC-II access its own memory view.

---

## Part 5: Estimated Effort

The estimates below cover system foundations and CPU integration tests. They are
planning estimates, not measured implementation costs. Device execution is outside
this estimate.

| Component | New LOC |
|-----------|---------|
| CPU integration tests (existing 6510 variant) | ~50 |
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

- **VICE:** [Source repository](https://github.com/VICE-Team/svn-mirror).
- **C64 Programmer's Reference Guide** (Commodore, 1982): Official hardware and software
  reference. Available online.
- **Mapping the Commodore 64** (Sheldon Leemon): Complete memory map reference with every
  address documented.
- **C64 Wiki** (c64-wiki.com): Community-maintained technical reference.
- **Wolfgang Lorenz CPU test suite**: Tests 6510 behavior including undocumented opcodes
  and decimal mode edge cases.
