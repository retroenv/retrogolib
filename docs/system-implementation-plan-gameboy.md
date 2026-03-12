# System Implementation Plan: Game Boy / Game Boy Color

## Current Status

- **Status:** PLANNED
- **Last Updated:** 2026-03-11
- **Existing Work:** Z80 package provides `GameBoyMemory` with MBC1 banking and
  `WithSystemType(arch.GameBoy)` option (sets PC=$0100, SP=$FFFE). The `GameBoy` system
  constant is already registered in `arch/system.go`.
- **Dependencies:** Z80 package (complete), new LR35902 CPU package (required)

## Context

The Game Boy (1989) and Game Boy Color (1998) are among the most popular retro platforms,
with a combined library of ~2,500 games and an exceptionally active homebrew community.

### Hardware Overview

| Component | Game Boy (DMG) | Game Boy Color (CGB) |
|-----------|---------------|---------------------|
| CPU | Sharp LR35902 @ 4.19 MHz | Sharp LR35902 @ 4.19/8.39 MHz (double speed) |
| RAM | 8 KB | 32 KB (8 banks of 4 KB) |
| VRAM | 8 KB | 16 KB (2 banks of 8 KB) |
| ROM | 32 KB - 8 MB (cartridge) | 32 KB - 8 MB (cartridge) |
| Display | 160x144, 4 shades | 160x144, 32,768 colors |
| Audio | 4 channels | 4 channels |
| I/O | Serial link, joypad | Serial link, joypad, infrared |

### Why a New CPU Package is Required

The Sharp LR35902 is often described as "a modified Z80" but the differences are substantial
enough to require a new package rather than a Z80 variant:

| Feature | Z80 | LR35902 |
|---------|-----|---------|
| Clock speed | Variable | 4.19 MHz (fixed) |
| Address bus | 16-bit | 16-bit |
| Registers | A,B,C,D,E,H,L + shadow set | A,B,C,D,E,H,L (**no shadow set**) |
| Index registers | IX, IY | **None** |
| I/O instructions | IN/OUT (256 ports) | **None** (memory-mapped I/O only) |
| Interrupt modes | IM 0/1/2 | **Single mode** (fixed vector table) |
| CB prefix | Bit/shift ops on registers | Bit/shift ops + **SWAP** (rotate nibbles) |
| Unique instructions | None | **STOP** (enter low-power), **SWAP** |
| Flag bits | S, Z, H, P/V, N, C (+ X, Y undocumented) | Z, N, H, C (**only 4 flags**) |
| DAA behavior | Complex (both add/sub, all flags) | Simplified (fewer flag interactions) |
| Halt bug | Different behavior | **HALT bug** (PC not incremented when IME=0 and interrupt pending) |
| Prefix handling | CB, DD, ED, FD | **CB only** (no DD/ED/FD) |

**Key architectural differences preventing a variant approach:**

1. **Missing entire register sets** -- No IX, IY, shadow registers (AF', BC', DE', HL').
   These are not just unused; the opcodes that would access them are **repurposed** for
   different instructions.
2. **Missing entire instruction groups** -- No ED prefix (block transfers LDIR/LDDR/CPIR/
   CPDR, I/O instructions, interrupt mode setting). No DD/FD prefix at all.
3. **Different flag register** -- Only 4 flags vs 8. No Sign, no Parity/Overflow, no
   undocumented X/Y flags. This changes the behavior of every flag-affecting instruction.
4. **Repurposed opcodes** -- Z80 opcodes $08, $10, $22, $2A, $32, $3A, $D9, $E0, $E2,
   $E8, $EA, $F0, $F2, $F8 have **completely different meanings** on the LR35902.

Attempting to shoehorn these differences into the Z80 package would require conditionals
in virtually every instruction handler and would make the Z80 code harder to maintain.

---

## Part 1: CPU Package -- Sharp LR35902

### 1.1 Approach: New Package `arch/cpu/lr35902/`

A new package following the `cpu-architecture-guidelines.md` patterns.

### 1.2 Registers

| Register | Width | Description |
|----------|-------|-------------|
| A | 8-bit | Accumulator |
| F | 8-bit | Flags: Z (bit 7), N (bit 6), H (bit 5), C (bit 4); bits 3-0 always 0 |
| B, C | 8-bit | General purpose (BC pair for 16-bit ops) |
| D, E | 8-bit | General purpose (DE pair for 16-bit ops) |
| H, L | 8-bit | General purpose (HL pair, also memory pointer) |
| SP | 16-bit | Stack pointer |
| PC | 16-bit | Program counter |
| IME | 1-bit | Interrupt Master Enable (not part of any register) |

**Register pairs** for 16-bit operations: AF, BC, DE, HL, SP.

### 1.3 Flags

Only 4 flags, stored in the upper nibble of F:

| Bit | Flag | Name | Description |
|-----|------|------|-------------|
| 7 | Z | Zero | Set when result is zero |
| 6 | N | Subtract | Set when last operation was subtraction |
| 5 | H | Half-carry | Set on carry from bit 3 to bit 4 |
| 4 | C | Carry | Set on carry from bit 7 |
| 3-0 | - | Unused | Always 0 |

### 1.4 Instruction Set

The LR35902 has 245 base opcodes (11 unused slots) and 256 CB-prefixed opcodes.

**Instructions shared with Z80 (same opcode, same behavior):**
Most ALU, load, push/pop, call/return, and branch instructions are identical. However,
the flag behavior differs due to the reduced flag set (no S, P/V, X, Y flags).

**Instructions unique to LR35902 (not in Z80):**

| Opcode | Mnemonic | Description |
|--------|----------|-------------|
| $10 | STOP | Enter low-power mode (until button press or interrupt) |
| $08 | LD (nn),SP | Store SP at 16-bit address |
| $E0 | LDH (n),A | Store A at $FF00+n (high RAM / I/O ports) |
| $F0 | LDH A,(n) | Load A from $FF00+n |
| $E2 | LD (C),A | Store A at $FF00+C |
| $F2 | LD A,(C) | Load A from $FF00+C |
| $E8 | ADD SP,e | Add signed 8-bit to SP (flags: 0 0 H C) |
| $F8 | LD HL,SP+e | Load SP + signed 8-bit into HL (flags: 0 0 H C) |
| $EA | LD (nn),A | Store A at 16-bit address |
| $FA | LD A,(nn) | Load A from 16-bit address |
| $D9 | RETI | Return from interrupt and enable interrupts (IME=1) |
| $CB 3x | SWAP r | Swap upper and lower nibbles of register |
| $CB 36 | SWAP (HL) | Swap upper and lower nibbles of (HL) |

**Z80 opcodes NOT present (unused, should trigger illegal opcode):**

$D3, $DB, $DD, $E3, $E4, $EB, $EC, $ED, $F4, $FC, $FD (11 opcodes).

**Z80 opcodes with DIFFERENT behavior:**

| Opcode | Z80 | LR35902 |
|--------|-----|---------|
| $08 | EX AF,AF' | LD (nn),SP |
| $10 | DJNZ e | STOP |
| $22 | LD (nn),HL | LD (HL+),A (store A, increment HL) |
| $2A | LD HL,(nn) | LD A,(HL+) (load A, increment HL) |
| $32 | LD (nn),A | LD (HL-),A (store A, decrement HL) |
| $3A | LD A,(nn) | LD A,(HL-) (load A, decrement HL) |
| $D9 | EXX | RETI |
| $E0 | RET PO | LDH (n),A |
| $E2 | JP PO,nn | LD (C),A |
| $E8 | RET PE | ADD SP,e |
| $EA | JP PE,nn | LD (nn),A |
| $F0 | RET P | LDH A,(n) |
| $F2 | JP P,nn | LD A,(C) |
| $F8 | RET M | LD HL,SP+e |

### 1.5 Interrupt System

The LR35902 has a single interrupt mode with 5 fixed vectors:

| Priority | Vector | Bit | Source |
|----------|--------|-----|--------|
| Highest | $0040 | 0 | V-Blank |
| | $0048 | 1 | LCD STAT |
| | $0050 | 2 | Timer |
| | $0058 | 3 | Serial |
| Lowest | $0060 | 4 | Joypad |

Controlled by two I/O registers:
- **IE** ($FFFF) -- Interrupt Enable: which interrupts are enabled
- **IF** ($FF0F) -- Interrupt Flag: which interrupts are pending

Interrupt dispatch: when `IME=1` and `(IE & IF) != 0`, the highest-priority pending
interrupt is serviced: IME is cleared, the corresponding IF bit is cleared, PC is pushed
to the stack, and PC is set to the vector address. Total: 20 cycles.

**HALT bug:** When HALT is executed with IME=0 and an interrupt is pending (IE & IF != 0),
the CPU resumes but **fails to increment PC** for the next instruction, causing the byte
after HALT to be read twice. This is a well-documented silicon bug.

### 1.6 CB Prefix Instructions

All 256 CB-prefixed opcodes are valid. They follow a regular pattern:

| Range | Instruction | Description |
|-------|-------------|-------------|
| $00-$07 | RLC r | Rotate left circular |
| $08-$0F | RRC r | Rotate right circular |
| $10-$17 | RL r | Rotate left through carry |
| $18-$1F | RR r | Rotate right through carry |
| $20-$27 | SLA r | Shift left arithmetic |
| $28-$2F | SRA r | Shift right arithmetic |
| $30-$37 | SWAP r | **Swap nibbles** (unique to LR35902) |
| $38-$3F | SRL r | Shift right logical |
| $40-$7F | BIT b,r | Test bit |
| $80-$BF | RES b,r | Reset bit |
| $C0-$FF | SET b,r | Set bit |

Register encoding (bits 2-0): B, C, D, E, H, L, (HL), A.

---

## Part 2: System Package -- Game Boy

### 2.1 Memory Map

| Address Range | Size | Description |
|---------------|------|-------------|
| $0000-$00FF | 256 B | Boot ROM / Interrupt vectors ($0000-$0067) |
| $0100-$014F | 80 B | Cartridge header |
| $0150-$3FFF | ~16 KB | ROM bank 0 (fixed) |
| $4000-$7FFF | 16 KB | ROM bank 1-N (switchable) |
| $8000-$9FFF | 8 KB | Video RAM (VRAM) |
| $A000-$BFFF | 8 KB | External RAM (cartridge, bank-switchable) |
| $C000-$CFFF | 4 KB | Work RAM bank 0 |
| $D000-$DFFF | 4 KB | Work RAM bank 1 (CGB: bank 1-7) |
| $E000-$FDFF | ~8 KB | Echo RAM (mirror of $C000-$DDFF) |
| $FE00-$FE9F | 160 B | OAM (Object Attribute Memory, 40 sprites) |
| $FEA0-$FEFF | 96 B | Unusable |
| $FF00-$FF7F | 128 B | I/O registers |
| $FF80-$FFFE | 127 B | High RAM (HRAM) |
| $FFFF | 1 B | Interrupt Enable register (IE) |

### 2.2 I/O Registers ($FF00-$FF7F)

**Joypad:**

| Address | Name | Description |
|---------|------|-------------|
| $FF00 | P1/JOYP | Joypad input (active low, select via bits 4-5) |

**Serial:**

| Address | Name | Description |
|---------|------|-------------|
| $FF01 | SB | Serial transfer data |
| $FF02 | SC | Serial transfer control |

**Timer:**

| Address | Name | Description |
|---------|------|-------------|
| $FF04 | DIV | Divider register (increments at 16384 Hz) |
| $FF05 | TIMA | Timer counter |
| $FF06 | TMA | Timer modulo (reload value) |
| $FF07 | TAC | Timer control (enable, frequency select) |

**Interrupt:**

| Address | Name | Description |
|---------|------|-------------|
| $FF0F | IF | Interrupt flag (pending interrupts) |

**Audio (Sound):**

| Address | Name | Description |
|---------|------|-------------|
| $FF10-$FF14 | NR10-NR14 | Channel 1: Pulse with sweep |
| $FF16-$FF19 | NR21-NR24 | Channel 2: Pulse |
| $FF1A-$FF1E | NR30-NR34 | Channel 3: Wave |
| $FF20-$FF23 | NR41-NR44 | Channel 4: Noise |
| $FF24-$FF26 | NR50-NR52 | Sound control |
| $FF30-$FF3F | Wave RAM | 16 bytes wave pattern |

**Video (LCD):**

| Address | Name | Description |
|---------|------|-------------|
| $FF40 | LCDC | LCD control |
| $FF41 | STAT | LCD status |
| $FF42 | SCY | Background scroll Y |
| $FF43 | SCX | Background scroll X |
| $FF44 | LY | LCD Y coordinate (current scanline, read-only) |
| $FF45 | LYC | LY compare |
| $FF46 | DMA | OAM DMA transfer start address |
| $FF47 | BGP | Background palette (DMG) |
| $FF48 | OBP0 | Object palette 0 (DMG) |
| $FF49 | OBP1 | Object palette 1 (DMG) |
| $FF4A | WY | Window Y position |
| $FF4B | WX | Window X position |

**CGB-Only Registers:**

| Address | Name | Description |
|---------|------|-------------|
| $FF4D | KEY1 | Speed switch (prepare double speed) |
| $FF4F | VBK | VRAM bank select |
| $FF51-$FF55 | HDMA1-5 | HDMA transfer |
| $FF68-$FF6B | BCPS-OCPD | Color palette specification/data |
| $FF70 | SVBK | WRAM bank select |

### 2.3 Cartridge Format

Game Boy ROMs have a 80-byte header at $0100-$014F:

| Offset | Size | Description |
|--------|------|-------------|
| $0100-$0103 | 4 B | Entry point (usually NOP + JP $0150) |
| $0104-$0133 | 48 B | Nintendo logo (validated by boot ROM) |
| $0134-$0143 | 16 B | Title (uppercase ASCII) |
| $013F-$0142 | 4 B | Manufacturer code (CGB) |
| $0143 | 1 B | CGB flag ($80=CGB compatible, $C0=CGB only) |
| $0144-$0145 | 2 B | New licensee code |
| $0146 | 1 B | SGB flag |
| $0147 | 1 B | Cartridge type (MBC type + features) |
| $0148 | 1 B | ROM size (32KB << value) |
| $0149 | 1 B | RAM size |
| $014A | 1 B | Destination (Japan/International) |
| $014B | 1 B | Old licensee code |
| $014C | 1 B | ROM version |
| $014D | 1 B | Header checksum |
| $014E-$014F | 2 B | Global checksum |

**Memory Bank Controllers (MBC):**

| Type Byte | MBC | Features |
|-----------|-----|----------|
| $00 | None | ROM only (32 KB max) |
| $01-$03 | MBC1 | Up to 2 MB ROM, 32 KB RAM |
| $05-$06 | MBC2 | Up to 256 KB ROM, 512x4 bits RAM |
| $0F-$13 | MBC3 | Up to 2 MB ROM, 32 KB RAM, RTC |
| $19-$1E | MBC5 | Up to 8 MB ROM, 128 KB RAM |

The existing `GameBoyMemory` in the Z80 package implements MBC1. This should be moved or
extended in the new system package.

### 2.4 File Structure

```
arch/cpu/lr35902/
    doc.go              -- Package documentation
    addressing.go       -- Addressing modes
    instruction.go      -- Instruction definitions
    opcode.go           -- 256-entry base opcode table
    opcode_cb.go        -- 256-entry CB-prefix opcode table
    categories.go       -- Instruction category sets
    errors.go           -- Package-specific errors
    flag.go             -- Flag definitions (Z, N, H, C only)
    cpu.go              -- CPU state and registers
    option.go           -- Functional options
    memory.go           -- Memory interface
    step.go             -- Fetch/decode/execute cycle
    param.go            -- Operand reading
    emulation.go        -- ALU instruction handlers
    emulation_load.go   -- Load/store instruction handlers
    emulation_branch.go -- Branch/call/return handlers
    emulation_cb.go     -- CB-prefix instruction handlers (shifts, bits, swap)
    interrupt.go        -- Interrupt dispatch (5 vectors)

arch/system/gameboy/
    doc.go              -- Package documentation
    gameboy.go          -- Memory map constants, system configuration
    register/
        io.go           -- I/O register addresses and names ($FF00-$FF7F)
        lcd.go          -- LCD register constants
        audio.go        -- Audio register constants
    cartridge/
        cartridge.go    -- Cartridge struct, header parsing
        header.go       -- Header format definitions
        mbc.go          -- MBC type detection
```

---

## Part 3: Migration of Existing Game Boy Code

The Z80 package currently contains Game Boy-specific code that should be evaluated:

| Current Location | Action |
|-----------------|--------|
| `z80/memory_gameboy.go` (`GameBoyMemory`) | Keep as-is (MBC1 memory for Z80-based use) or deprecate in favor of system package |
| `z80/option.go` (`WithSystemType(arch.GameBoy)`) | Keep for backward compatibility; new LR35902 package handles GB natively |
| `arch/system.go` (`GameBoy` constant) | Already registered; no change needed |

The `GameBoyMemory` type in the Z80 package can remain for backward compatibility. The new
`arch/system/gameboy/` package will provide a more complete implementation with all MBC
types and proper I/O register handling.

---

## Part 4: Implementation Phases

### Phase 1: LR35902 Static Analysis Foundation
- `arch/arch.go` -- Add `LR35902 Architecture = "lr35902"`
- Create `arch/cpu/lr35902/` with doc.go, addressing.go, instruction.go, opcode.go,
  opcode_cb.go, categories.go, errors.go, flag.go
- Full 256+256 opcode tables with instruction metadata

### Phase 2: LR35902 CPU Emulation Core
- cpu.go, option.go, memory.go, step.go, param.go
- Interrupt handling with 5 fixed vectors and IME
- HALT instruction with HALT bug emulation

### Phase 3: LR35902 Instruction Handlers
- emulation.go -- ALU (ADD, ADC, SUB, SBC, AND, OR, XOR, CP, INC, DEC, DAA, CPL, CCF, SCF)
- emulation_load.go -- LD variants, LDH, LD (HL+/-), PUSH, POP
- emulation_branch.go -- JP, JR, CALL, RET, RETI, RST
- emulation_cb.go -- All CB-prefix ops including SWAP

### Phase 4: System Package
- Create `arch/system/gameboy/` with memory map constants
- I/O register definitions (joypad, timer, serial, LCD, audio)
- Cartridge header parsing and MBC type detection

### Phase 5: Testing
- Opcode table completeness test (245 base + 256 CB = 501 valid opcodes)
- Per-instruction unit tests with flag verification
- HALT bug test
- Interrupt dispatch tests (priority, IME, IE/IF interaction)
- Cartridge header parsing tests
- Blargg's cpu_instrs test ROM (11 sub-tests, the standard GB CPU validation suite)
- Blargg's instr_timing test ROM (cycle accuracy)

---

## Part 5: Design Decisions

### New Package vs Z80 Variant
**Decision: New package `arch/cpu/lr35902/`**
- Rationale: The LR35902 removes entire register sets (IX, IY, shadow registers), removes
  3 of 4 prefix groups (DD, ED, FD), removes I/O instructions, has only 4 flags instead
  of 8, and repurposes 14+ opcodes for completely different instructions. The remaining
  shared logic (basic ALU, common loads) is not enough to justify the complexity of
  conditionalizing the Z80 package. A clean implementation is clearer and more maintainable.

### Flag Storage
**Decision: Packed byte (upper nibble of F register)**
- Rationale: With only 4 flags that must always occupy the upper nibble of F (lower nibble
  always 0), a packed byte is natural and efficient. The Z80's individual flag fields
  approach was justified by 8 flags including undocumented bits; the LR35902's simpler
  flag model doesn't need that complexity.

### Code Reuse from Z80
**Decision: No direct code sharing, but same architectural patterns**
- Rationale: While the LR35902 and Z80 share some instruction mnemonics, the
  implementations differ in flag behavior, register sets, and instruction variants. Copying
  and adapting is acceptable for the ALU core, but sharing code via imports would create
  tight coupling between architecturally distinct CPUs. Follow the same patterns
  (opcode tables, step loop, param readers) without creating dependencies.

---

## Part 6: Estimated Effort

| Component | New LOC |
|-----------|---------|
| Architecture registration | ~10 |
| Static analysis files (doc, addressing, instruction, opcode, categories, errors, flag) | ~1,800 |
| CPU state + options + memory interface | ~300 |
| Step/decode cycle + parameter reading | ~400 |
| Instruction handlers (4 emulation files) | ~2,000 |
| Interrupt handling + HALT bug | ~200 |
| System package (constants, registers) | ~400 |
| Cartridge format support | ~300 |
| Tests | ~1,500 |
| **Total** | **~6,910** |

---

## Part 7: References

- **Blargg's test ROMs**: cpu_instrs (11 tests) and instr_timing are the standard CPU
  validation suite. Available at: github.com/retrio/gb-test-roms
- **Pan Docs**: The comprehensive Game Boy technical reference.
  Available at: gbdev.io/pandocs
- **RGBDS**: The Game Boy assembler/linker toolchain, useful for test ROM generation.
- **SameBoy** (C): Highly accurate Game Boy emulator by Lior Halphon. Reference for
  edge cases and timing accuracy.
- **Gambatte** (C++): Well-regarded emulator focused on timing accuracy.
- **Game Boy CPU Manual** (various community authors): Instruction set reference with
  flag effects and cycle counts.
- **TCAGBD** (The Cycle-Accurate Game Boy Docs): Detailed timing documentation.
