# System Implementation Plan: Atari 2600

## Current Status

- **Status:** PLANNED
- **Last Updated:** 2026-03-11
- **Dependencies:** m6502 package (complete)

## Context

The Atari 2600 (1977) is one of the most iconic retro gaming platforms, with a large library
of cartridge-based games and an active homebrew community. Its architecture is uniquely
constrained: the CPU has only 128 bytes of RAM, no frame buffer, and the programmer must
"race the beam" -- feeding data to the Television Interface Adapter (TIA) in real time as
each scanline is drawn.

### Hardware Overview

| Component | Details |
|-----------|---------|
| CPU | MOS 6507 (6502 with 13-bit address bus, no IRQ/NMI pins) |
| RAM | 128 bytes (RIOT chip, $0080-$00FF) |
| ROM | 2 KB - 64 KB (cartridge, bank-switched for >4 KB) |
| Video | TIA (Television Interface Adapter), 160x~192 resolution |
| Audio | TIA (2 channels) |
| I/O | RIOT (6532, timers + I/O ports for controllers) |
| Address bus | 13 bits (8 KB address space, mirrored) |
| Clock | 1.19 MHz (NTSC), 1.18 MHz (PAL) |

### Why the 2600 is Architecturally Unique

The Atari 2600 has no frame buffer. The TIA generates video output one scanline at a time,
and the CPU must update TIA registers between scanlines (during HBLANK) or mid-scanline
for effects. This "racing the beam" technique means:

1. **Cycle-exact CPU timing is critical** -- the CPU and TIA share a clock relationship
   (1 CPU cycle = 3 TIA color clocks). Software depends on exact cycle counts.
2. **No vertical blank interrupt** -- the CPU polls VSYNC/VBLANK and counts scanlines
   manually. There is no IRQ or NMI hardware.
3. **Bank switching is mapper-dependent** -- cartridges >4 KB use various banking schemes
   triggered by reading specific addresses (not writing).

---

## Part 1: CPU Variant -- MOS 6507

### 1.1 Approach: Extend m6502 Package with Variant

The 6507 is a pin-reduced 6502 in a 28-pin package. It executes the **identical instruction
set** as the NMOS 6502 -- same opcodes, same timing, same undocumented instructions. The
differences are purely in the external interface:

| Feature | 6502 | 6507 |
|---------|------|------|
| Package | 40-pin | 28-pin |
| Address bus | 16 bits (64 KB) | 13 bits (8 KB) |
| IRQ pin | Yes | **No** |
| NMI pin | Yes | **No** |
| RDY pin | Yes | **No** |
| Data bus | 8 bits | 8 bits |
| Instruction set | Full | **Identical** |

Since the instruction set is identical, the 6507 is a variant of the existing m6502 package,
not a new CPU package. The differences are:

1. **13-bit address masking:** All addresses are masked to $0000-$1FFF (8 KB). The full
   64 KB space mirrors this 8 KB region. This is handled by the Memory implementation,
   not the CPU itself.
2. **No IRQ/NMI:** The interrupt pins are not connected. `TriggerIRQ()` and `TriggerNMI()`
   should be no-ops or forbidden for this variant.
3. **No RDY pin:** The CPU cannot be halted externally. (Not currently modeled anyway.)

### 1.2 Implementation

#### Phase 1: Add Variant Constant

**Files to modify:**
- `arch/cpu/m6502/option.go` -- Add `Variant6507` constant
- `arch/arch.go` -- Add `M6507 Architecture = "6507"` (optional, since it's instruction-
  identical to 6502)

```go
// option.go - add to CPUVariant enum:
const (
    VariantNMOS6502 CPUVariant = iota
    VariantNES6502
    Variant6507    // MOS 6507: 6502 with 13-bit address bus, no IRQ/NMI
    Variant65C02
)
```

**Note:** `Variant6507` is placed before `Variant65C02` because the comparison
`c.opts.variant >= Variant65C02` is used to select the 65C02 opcode table. The 6507 must
use the NMOS 6502 table.

#### Phase 2: Disable Interrupts for 6507

**Files to modify:**
- `arch/cpu/m6502/cpu.go` -- Guard `TriggerIRQ()` and `TriggerNMI()` to be no-ops for 6507

```go
func (c *CPU) TriggerIRQ() {
    if c.opts.variant == Variant6507 {
        return // 6507 has no IRQ pin
    }
    c.triggerIrq = true
}

func (c *CPU) TriggerNMI() {
    if c.opts.variant == Variant6507 {
        return // 6507 has no NMI pin
    }
    c.triggerNmi = true
}
```

#### Phase 3: Testing

- Verify all existing 6502 tests pass with `Variant6507`
- Add unit test confirming `TriggerIRQ()`/`TriggerNMI()` are no-ops
- Test address mirroring at the Memory level (see Part 2)

### 1.3 Estimated Effort

| Component | New LOC | Modified LOC |
|-----------|---------|-------------|
| Variant constant | ~5 | ~5 |
| Interrupt guards | ~10 | ~10 |
| Tests | ~50 | ~0 |
| **Total** | **~65** | **~15** |

---

## Part 2: System Package -- Atari 2600

### 2.1 Architecture Registration

**Files to modify:**
- `arch/system.go` -- Add `Atari2600 System = "atari-2600"`

### 2.2 Memory Map

The 2600's 13-bit address space (8 KB) is divided as follows:

| Address Range | Size | Component | Description |
|---------------|------|-----------|-------------|
| $0000-$002C | 45 bytes | TIA | Write registers (player, missile, ball, playfield) |
| $0030-$003D | 14 bytes | TIA | Read registers (collision, input, timing) |
| $0080-$00FF | 128 bytes | RIOT | RAM |
| $0280-$0297 | 24 bytes | RIOT | I/O ports and timer registers |
| $1000-$1FFF | 4 KB | Cartridge | ROM (bank-switchable for larger ROMs) |

The entire space mirrors within the 13-bit range. Notable mirrors:
- TIA registers mirror every 64 bytes in $0000-$003F range
- RIOT RAM mirrors at $0180-$01FF
- RIOT I/O mirrors at $0280-$02FF and $0380-$03FF

### 2.3 TIA Registers

**Write registers ($00-$2C):**

| Address | Name | Description |
|---------|------|-------------|
| $00 | VSYNC | Vertical sync set-clear |
| $01 | VBLANK | Vertical blank set-clear |
| $02 | WSYNC | Wait for leading edge of horizontal blank |
| $03 | RSYNC | Reset horizontal sync counter |
| $04 | NUSIZ0 | Number-size player-missile 0 |
| $05 | NUSIZ1 | Number-size player-missile 1 |
| $06 | COLUP0 | Color-luminance player 0 |
| $07 | COLUP1 | Color-luminance player 1 |
| $08 | COLUPF | Color-luminance playfield |
| $09 | COLUBK | Color-luminance background |
| $0A | CTRLPF | Control playfield ball size and collisions |
| $0B | REFP0 | Reflect player 0 |
| $0C | REFP1 | Reflect player 1 |
| $0D | PF0 | Playfield register byte 0 |
| $0E | PF1 | Playfield register byte 1 |
| $0F | PF2 | Playfield register byte 2 |
| $10 | RESP0 | Reset player 0 |
| $11 | RESP1 | Reset player 1 |
| $12 | RESM0 | Reset missile 0 |
| $13 | RESM1 | Reset missile 1 |
| $14 | RESBL | Reset ball |
| $15 | AUDC0 | Audio control 0 |
| $16 | AUDC1 | Audio control 1 |
| $17 | AUDF0 | Audio frequency 0 |
| $18 | AUDF1 | Audio frequency 1 |
| $19 | AUDV0 | Audio volume 0 |
| $1A | AUDV1 | Audio volume 1 |
| $1B | GRP0 | Graphics player 0 |
| $1C | GRP1 | Graphics player 1 |
| $1D | ENAM0 | Graphics enable missile 0 |
| $1E | ENAM1 | Graphics enable missile 1 |
| $1F | ENABL | Graphics enable ball |
| $20 | HMP0 | Horizontal motion player 0 |
| $21 | HMP1 | Horizontal motion player 1 |
| $22 | HMM0 | Horizontal motion missile 0 |
| $23 | HMM1 | Horizontal motion missile 1 |
| $24 | HMBL | Horizontal motion ball |
| $25 | VDELP0 | Vertical delay player 0 |
| $26 | VDELP1 | Vertical delay player 1 |
| $27 | VDELBL | Vertical delay ball |
| $28 | RESMP0 | Reset missile 0 to player 0 |
| $29 | RESMP1 | Reset missile 1 to player 1 |
| $2A | HMOVE | Apply horizontal motion |
| $2B | HMCLR | Clear horizontal motion registers |
| $2C | CXCLR | Clear collision latches |

**Read registers ($00-$0D, active bits only):**

| Address | Name | Description |
|---------|------|-------------|
| $00 | CXM0P | Collision M0-P1, M0-P0 (bits 7-6) |
| $01 | CXM1P | Collision M1-P0, M1-P1 (bits 7-6) |
| $02 | CXP0FB | Collision P0-PF, P0-BL (bits 7-6) |
| $03 | CXP1FB | Collision P1-PF, P1-BL (bits 7-6) |
| $04 | CXM0FB | Collision M0-PF, M0-BL (bits 7-6) |
| $05 | CXM1FB | Collision M1-PF, M1-BL (bits 7-6) |
| $06 | CXBLPF | Collision BL-PF (bit 7) |
| $07 | CXPPMM | Collision P0-P1, M0-M1 (bits 7-6) |
| $08 | INPT0 | Paddle 0 input (bit 7) |
| $09 | INPT1 | Paddle 1 input (bit 7) |
| $0A | INPT2 | Paddle 2 input (bit 7) |
| $0B | INPT3 | Paddle 3 input (bit 7) |
| $0C | INPT4 | Joystick 0 trigger (bit 7) |
| $0D | INPT5 | Joystick 1 trigger (bit 7) |

### 2.4 RIOT (6532) Registers

| Address | Name | Description |
|---------|------|-------------|
| $0280 | SWCHA | Port A: joystick directions (read/write) |
| $0281 | SWACNT | Port A DDR (data direction register) |
| $0282 | SWCHB | Port B: console switches (read) |
| $0283 | SWBCNT | Port B DDR |
| $0284 | INTIM | Timer output (read) |
| $0285 | INSTAT | Timer interrupt status (read) |
| $0294 | TIM1T | Set 1-clock interval timer (write) |
| $0295 | TIM8T | Set 8-clock interval timer (write) |
| $0296 | TIM64T | Set 64-clock interval timer (write) |
| $0297 | T1024T | Set 1024-clock interval timer (write) |

### 2.5 Cartridge Format

Atari 2600 ROMs are raw binary files (no header). The cartridge size determines the
banking scheme:

| Size | Banks | Scheme | Trigger Addresses |
|------|-------|--------|-------------------|
| 2 KB | 1 | None | N/A |
| 4 KB | 1 | None | N/A |
| 8 KB | 2 | F8 | $1FF8-$1FF9 |
| 12 KB | 3 | FA | $1FF8-$1FFA |
| 16 KB | 4 | F6 | $1FF6-$1FF9 |
| 32 KB | 8 | F4 | $1FF4-$1FFB |
| 64 KB | 16 | 3F (Tigervision) | Write bank to $003F |

Bank switching is triggered by **reading** (not writing) the trigger addresses.
This is a unique quirk of the 2600 -- the cartridge detects address bus activity
regardless of read/write.

Some cartridges also include extra RAM (128 bytes - 256 bytes) mapped into the
cartridge address space.

### 2.6 File Structure

```
arch/system/atari2600/
    doc.go              -- Package documentation
    atari2600.go        -- Memory map constants, address ranges
    register/
        tia.go          -- TIA register addresses and names
        riot.go         -- RIOT register addresses and names
    cartridge/
        cartridge.go    -- Cartridge struct, banking scheme detection
        format.go       -- Raw ROM loading, size-based scheme selection
```

---

## Part 3: Implementation Phases

### Phase 1: CPU Variant (6507)
- Add `Variant6507` to m6502 option.go
- Guard IRQ/NMI triggers
- Unit tests

### Phase 2: System Registration
- Add `Atari2600` system constant to `arch/system.go`
- Create `arch/system/atari2600/` package

### Phase 3: Memory Map and Registers
- Define TIA write/read register constants
- Define RIOT register constants
- Define memory map address ranges and mirrors

### Phase 4: Cartridge Support
- Raw ROM loading (no header)
- Size-based banking scheme detection
- F8/FA/F6/F4 bank switching via address read triggers

### Phase 5: Testing
- Opcode execution tests with 6507 variant
- Memory mirroring tests
- Cartridge loading and bank switching tests
- Register address completeness tests

---

## Part 4: Design Decisions

### 6507 as Variant vs Separate Package
**Decision: Variant within m6502**
- Rationale: The instruction set is byte-for-byte identical to the NMOS 6502. The only
  differences (13-bit address bus, no IRQ/NMI) are external interface constraints handled
  by the Memory implementation and interrupt guards. A separate package would duplicate
  the entire m6502 codebase for zero instruction-level benefit.

### Address Masking
**Decision: Memory implementation handles 13-bit masking, not the CPU**
- Rationale: The 6507 internally computes 16-bit addresses identically to the 6502.
  The address bus is simply truncated externally. The Memory implementation applies the
  `& 0x1FFF` mask, which is more flexible (allows the same CPU to run in either
  configuration) and matches how real hardware works.

### Bank Switching via Reads
**Decision: Memory implementation detects read-triggered bank switches**
- Rationale: The 2600's bank switching is triggered by the address bus, not by explicit
  write operations. The Memory.Read() implementation must detect accesses to trigger
  addresses and switch banks before returning data. This is the only correct approach.

---

## Part 5: Estimated Effort

| Component | New LOC |
|-----------|---------|
| CPU variant (6507 in m6502) | ~80 |
| System package (constants, memory map) | ~300 |
| TIA/RIOT register definitions | ~200 |
| Cartridge format support | ~250 |
| Tests | ~400 |
| **Total** | **~1,230** |

---

## Part 6: References

- **Stella** (C++): The gold standard Atari 2600 emulator. Mature, well-documented source.
  Repository: `stella-emu/stella`
- **Atari 2600 Programming Guide** (Nick Bensema): Comprehensive hardware reference.
- **TIA Hardware Notes** (Andrew Towers): Detailed TIA operation and timing.
- **Stella Programmer's Guide** (Steve Wright, 1979): Original Atari developer documentation.
