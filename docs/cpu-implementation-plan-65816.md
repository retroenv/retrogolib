# CPU Implementation Plan: WDC 65C816 (65816)

## Context

The WDC 65C816 is the 16-bit successor to the 65C02, designed by Bill Mensch at Western Design
Center. It was the CPU of two major retro platforms:

- **Super Nintendo Entertainment System (SNES/Super Famicom)** -- The 16-bit console that
  defined an era of gaming (1990-1998), running the 65816 at 3.58 MHz (fast) / 2.68 MHz (slow)
  / 1.79 MHz (extra slow) depending on memory region.
- **Apple IIGS** -- Apple's 16-bit successor to the Apple II line (1986), running at 2.8 MHz.

The 65816 maintains full backward compatibility with the 65C02 in "emulation mode" while
providing a dramatically expanded architecture in "native mode": 16-bit registers, 24-bit
address space, new addressing modes, and new instructions.

---

## 1. Architecture Overview

### 1.1 Registers

| Register | Width | Description |
|----------|-------|-------------|
| C | 16-bit | Accumulator (splits into A low byte / B high byte in 8-bit mode) |
| X | 16-bit/8-bit | Index register X (width controlled by X flag) |
| Y | 16-bit/8-bit | Index register Y (width controlled by X flag) |
| SP | 16-bit | Stack pointer (full 64K range in native mode, page 1 in emulation) |
| DP | 16-bit | Direct Page register (replaces fixed zero page $00xx) |
| DB | 8-bit | Data Bank register (bits 16-23 of data addresses) |
| PB | 8-bit | Program Bank register (bits 16-23 of instruction fetches) |
| PC | 16-bit | Program counter (within current program bank) |
| P | 8-bit | Processor status: N V M X D I Z C |
| E | 1-bit | Emulation flag (hidden, toggled via XCE) |

**Effective address width:** PB:PC forms a 24-bit program address. DB is prepended to 16-bit
data addresses to form 24-bit data addresses. The full address space is 16 MB (2^24).

### 1.2 Processor Modes

**Emulation Mode (E=1):**
- Behaves like a 65C02 (8-bit A, 8-bit X/Y, 256-byte stack on page 1)
- M and X flags are forced to 1
- Direct Page register still works (unlike real zero page)
- DB and PB still function (can access full 24-bit space)
- Default mode after reset

**Native Mode (E=0):**
- Full 16-bit capabilities
- M flag controls accumulator width: M=1 -> 8-bit, M=0 -> 16-bit
- X flag controls index register width: X=1 -> 8-bit, X=0 -> 16-bit
- Stack pointer is full 16-bit (can be anywhere in bank 0)
- When X flag transitions 0->1, high bytes of X and Y are zeroed

### 1.3 Memory Map

The 65816 uses a 24-bit address space organized into 256 banks of 64K each:

- **Bank $00** -- Direct Page, Stack, interrupt vectors (mirrors 65C02 layout)
- **Bank $00-$7F** -- Available for ROM/RAM mapping (system-dependent)
- **Bank $80-$FF** -- Mirror or additional address space (system-dependent)
- Interrupt vectors are always at $00:FFE0-$00:FFFF (native) and $00:FFF0-$00:FFFF (emulation)

The actual memory map is system-specific (SNES LoROM vs HiROM vs ExHiROM, Apple IIGS slots).

---

## 2. Key Differences from 6502/65C02

| Feature | 6502 | 65C02 | 65816 |
|---------|------|-------|-------|
| Address bus | 16-bit | 16-bit | 24-bit |
| Accumulator | 8-bit | 8-bit | 8 or 16-bit (M flag) |
| Index registers | 8-bit | 8-bit | 8 or 16-bit (X flag) |
| Stack pointer | 8-bit (page 1) | 8-bit (page 1) | 16-bit (any location) |
| Zero/Direct page | Fixed $00xx | Fixed $00xx | Relocatable (DP register) |
| Bank registers | None | None | DB (data), PB (program) |
| Addressing modes | 13 | 14 (+ZP indirect) | ~38 (+24 new) |
| Instructions | ~56 | ~86 (+30) | ~114 (+28 new) |
| Opcode slots used | 151 | 212 | 256 (all filled) |
| Processor modes | 1 | 1 | 2 (emulation + native) |

---

## 3. Why a New Package (Not a Variant)

Unlike the 65C02 (which was implemented as a variant of m6502), the 65816 requires a **new
package `arch/cpu/m65816/`**. The delta is too large for a variant approach:

1. **24-bit address space** -- Memory interfaces change fundamentally from `uint16` to
   `uint32` (masked to 24 bits). Every memory read/write function signature differs.
2. **Variable-width registers** -- A, X, Y can be 8 or 16 bits depending on M/X flags.
   This affects instruction size, cycle counts, and all ALU operations.
3. **~24 new addressing modes** -- Stack-relative, long indirect, block move, etc.
   These are not minor additions but entire new categories of memory access.
4. **Direct Page register** -- Replaces the fixed zero page concept. All "zero page"
   addressing modes now use DP as a base, with potential bank-crossing behavior.
5. **Bank registers (DB, PB)** -- Every data access implicitly uses DB, every instruction
   fetch uses PB. This is pervasive, not a conditional tweak.
6. **Dual-mode operation** -- Emulation vs Native mode affects nearly every instruction's
   behavior, register widths, stack behavior, and interrupt handling.
7. **Same opcode, different sizes** -- The same opcode can be 2 or 3 bytes depending on
   M/X flags. This makes static analysis significantly more complex than 6502.

Attempting to shoehorn these changes into `m6502/` would require conditionals in virtually
every function, making the code harder to understand and maintain than a clean implementation.

---

## 4. New Addressing Modes (~24)

The 65816 inherits all 14 addressing modes from the 65C02 and adds approximately 24 new ones.
Many are "long" (24-bit) variants of existing modes.

### Stack-Relative Modes (2)

| Mode | Syntax | Description |
|------|--------|-------------|
| Stack Relative | `sr,S` | Offset from stack pointer |
| Stack Relative Indirect Indexed Y | `(sr,S),Y` | Indirect through stack, indexed by Y |

### Long Modes (4)

| Mode | Syntax | Description |
|------|--------|-------------|
| Absolute Long | `al` | 24-bit absolute address |
| Absolute Long Indexed X | `al,X` | 24-bit absolute + X |
| Direct Page Indirect Long | `[dp]` | 24-bit pointer at direct page location |
| Direct Page Indirect Long Indexed Y | `[dp],Y` | 24-bit pointer + Y |

### Block Move Modes (2)

| Mode | Syntax | Description |
|------|--------|-------------|
| Block Move (MVN) | `src,dst` | Source and destination banks for block move next |
| Block Move (MVP) | `src,dst` | Source and destination banks for block move previous |

### Relative Long Mode (1)

| Mode | Syntax | Description |
|------|--------|-------------|
| Program Counter Relative Long | `rl` | 16-bit signed offset (for BRL, PER) |

### Other New Modes

| Mode | Syntax | Description |
|------|--------|-------------|
| Absolute Indirect Long | `[abs]` | 24-bit pointer at absolute address (JML) |
| Absolute Indexed Indirect | `(abs,X)` | Indirect through abs+X in program bank (JSR/JMP) |

### Inherited 65C02 Modes (14)

Implied, Accumulator, Immediate, Zero Page (now Direct Page), Zero Page X, Zero Page Y,
Absolute, Absolute X, Absolute Y, Indirect, Indexed Indirect X, Indirect Indexed Y,
Relative, Zero Page Indirect.

**Note:** Immediate mode is variable-width: 1 byte when M=1/X=1, 2 bytes when M=0/X=0
(depending on which register the instruction targets).

---

## 5. New Instructions (~28 New Mnemonics)

### Branch and Jump (3)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| BRL | $82 | Relative Long | Branch Always Long (16-bit offset) |
| JML | $5C/$DC | Abs Long / Indirect Long | Jump Long (sets PB) |
| JSL | $22 | Absolute Long | Jump to Subroutine Long (pushes PB:PC) |

### Return (1)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| RTL | $6B | Implied | Return from Subroutine Long (pulls PB:PC) |

### Block Move (2)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| MVN | $54 | Block Move | Move block Next (C+1 bytes, increment) |
| MVP | $44 | Block Move | Move block Previous (C+1 bytes, decrement) |

### Stack Operations (7)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| PEA | $F4 | Absolute | Push Effective Absolute address |
| PEI | $D4 | Direct Page Indirect | Push Effective Indirect address |
| PER | $62 | Relative Long | Push Effective Relative address |
| PHB | $8B | Implied | Push Data Bank register |
| PHD | $0B | Implied | Push Direct Page register (16-bit) |
| PHK | $4B | Implied | Push Program Bank register |
| PLB | $AB | Implied | Pull Data Bank register |
| PLD | $2B | Implied | Pull Direct Page register (16-bit) |

### Processor Status (2)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| REP | $C2 | Immediate | Reset Processor status bits (clear bits) |
| SEP | $E2 | Immediate | Set Processor status bits (set bits) |

### Register Transfer (6)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| TCD | $5B | Implied | Transfer C (16-bit accumulator) to Direct Page |
| TCS | $1B | Implied | Transfer C to Stack Pointer |
| TDC | $7B | Implied | Transfer Direct Page to C |
| TSC | $3B | Implied | Transfer Stack Pointer to C |
| TXY | $9B | Implied | Transfer X to Y |
| TYX | $BB | Implied | Transfer Y to X |

### Accumulator (1)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| XBA | $EB | Implied | Exchange B and A (swap high/low bytes of C) |

### Mode Switching (1)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| XCE | $FB | Implied | Exchange Carry and Emulation flags |

### System (3)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| COP | $02 | Immediate | Co-Processor enable (software interrupt) |
| STP | $DB | Implied | Stop the Processor (until reset) |
| WAI | $CB | Implied | Wait for Interrupt |

### WDM (1)

| Mnemonic | Opcode | Mode | Description |
|----------|--------|------|-------------|
| WDM | $42 | Immediate | Reserved for future expansion (2-byte NOP) |

---

## 6. Opcode Table

The 65816 fills all 256 opcode slots. Unlike the 6502/65C02, there are no unused opcodes.

**Key complexity:** The same opcode can have different instruction sizes depending on the
M and X processor flags:

- Instructions operating on the accumulator (LDA, STA, ADC, etc.) are 1 byte wider when M=0
- Instructions operating on index registers (LDX, LDY, CPX, CPY) are 1 byte wider when X=0
- This affects both immediate operands and cycle counts

The opcode table should use a single `[256]Opcode` array. The `Opcode` struct needs
additional fields or methods to handle variable sizing:

```go
type Opcode struct {
    Instruction *Instruction
    Addressing  AddressingMode
    Timing      byte // Base cycles (may vary with M/X flags and page crossing)
    BaseSize    byte // Size when M=1/X=1 (8-bit mode)
    // For instructions affected by M/X flag:
    // actual size = BaseSize + 1 when relevant flag is 0 (16-bit mode)
    WidthFlag   WidthFlag // None, M, or X -- which flag affects this instruction's width
}

type WidthFlag byte
const (
    WidthNone WidthFlag = iota // Fixed size
    WidthM                     // Size varies with M flag (accumulator operations)
    WidthX                     // Size varies with X flag (index operations)
)
```

---

## Current Status

- **Status:** IN_PROGRESS
- **Last Updated:** 2026-03-06
- **Summary:** Phases 1–5 complete; Phase 7 tests cover all major instruction groups. Two timing bugs fixed: BRA cycle count and branch page-crossing penalty in emulation mode. 84 tests passing.

## Completed Work

| Date | What | Notes |
|------|------|-------|
| 2026-03-06 | All Phases 1–4 implemented | All files created, 256-opcode table filled, all addressing mode param readers, all instruction handlers |
| 2026-03-06 | Fix `resolveEA` for `DirectPageX`/`DirectPageY` | Was always using 8-bit X/Y and truncating to `uint8`; now respects `IdxWidth()` and emulation-mode page-0 wrap |
| 2026-03-06 | Fix `mvn`/`mvp` off-by-one | Redundant "final byte" copy removed; loop alone correctly exits with C=0xFFFF |
| 2026-03-06 | Expand test coverage (Phase 7) | Added JSR/RTS, JSL/RTL, JMP, JML, PEA/PEI/PER, MVN 1+3-byte, BRK native, dp,X 8/16-bit, ADC/SBC 16-bit |
| 2026-03-06 | ADC/SBC decimal mode (BCD) — Phase 5 | `adcBCD8/16` and `sbcBCD8/16` helpers; `adc`/`sbc` dispatch on `Flags.D≠0`. V from binary intermediate; N/Z/C from BCD result. |
| 2026-03-06 | Expand test coverage round 2 (Phase 7) | Added BCD ADC/SBC (8+16-bit), MVP, RTI native, PHB/PLB, PHD/PLD, WAI+NMI dispatch |
| 2026-03-06 | Expand test coverage round 3 (Phase 7) | Added BRK/COP emulation mode (stack layout + vectors), CLC→XCE→REP→LDA mode-switch sequence, abs,X bank-boundary crossing. 81 tests total. |
| 2026-03-06 | Fix BRA timing + branch page-crossing (emulation mode compatibility) | BRA had Timing=3 but branch() always adds +1, giving 4 cycles (wrong). Fixed to Timing=2. Added page-crossing detection to paramReaderRelative; step.go now only applies branch page-cross penalty in emulation mode (E=1). 84 tests total. |

### Next Target: Remaining Phase 7 gaps (lower priority)

- **More cycle accuracy tests:** Verify page-crossing penalties for abs,Y and (dp),Y addressing modes; verify DP penalty cycle when DP_low≠0.
- **Emulation mode compatibility sweep:** Verify remaining 65C02 edge cases (ORA/AND/EOR, CMP, LDX/LDY in various addressing modes) in emulation mode.
- **SNES system layer (Phase 6):** LoROM/HiROM memory mapping, DMA, VBlank NMI timing.

---

## 7. Implementation Phases

### Phase 1: Architecture Registration + Static Analysis Foundation

**Files to modify:**
- `arch/arch.go` -- Add `M65816 Architecture = "65816"` constant
- `arch/system.go` -- Add `SNES System = "snes"` and `AppleIIGS System = "apple-iigs"`

**Files to create:**
- `arch/cpu/m65816/doc.go` -- Architecture overview, modes, usage examples
- `arch/cpu/m65816/addressing.go` -- ~38 addressing modes as typed constants
- `arch/cpu/m65816/instruction.go` -- ~114 instruction definitions (all 6502 + 65C02 + new)
- `arch/cpu/m65816/errors.go` -- Package-specific errors
- `arch/cpu/m65816/flag.go` -- Processor status flags including M, X, E

```go
// arch/arch.go - add:
M65816 Architecture = "65816"

// arch/system.go - add:
SNES     System = "snes"
AppleIIGS System = "apple-iigs"
```

```go
// addressing.go
type AddressingMode int

const (
    NoAddressing                            AddressingMode = 0
    ImpliedAddressing                       AddressingMode = 1 << iota
    AccumulatorAddressing
    ImmediateAddressing                     // 1 or 2 bytes depending on M/X
    DirectPageAddressing                    // dp
    DirectPageIndexedXAddressing            // dp,X
    DirectPageIndexedYAddressing            // dp,Y
    DirectPageIndirectAddressing            // (dp)
    DirectPageIndexedIndirectXAddressing    // (dp,X)
    DirectPageIndirectIndexedYAddressing    // (dp),Y
    DirectPageIndirectLongAddressing        // [dp]
    DirectPageIndirectLongIndexedYAddressing // [dp],Y
    AbsoluteAddressing                      // abs
    AbsoluteIndexedXAddressing              // abs,X
    AbsoluteIndexedYAddressing              // abs,Y
    AbsoluteIndirectAddressing              // (abs)
    AbsoluteIndexedIndirectAddressing       // (abs,X)
    AbsoluteIndirectLongAddressing          // [abs]
    AbsoluteLongAddressing                  // al (24-bit)
    AbsoluteLongIndexedXAddressing          // al,X
    StackRelativeAddressing                 // sr,S
    StackRelativeIndirectIndexedYAddressing // (sr,S),Y
    RelativeAddressing                      // 8-bit offset
    RelativeLongAddressing                  // 16-bit offset
    BlockMoveAddressing                     // src,dst
)
```

```go
// flag.go
const (
    FlagCarry     = 0 // C
    FlagZero      = 1 // Z
    FlagIRQ       = 2 // I - IRQ disable
    FlagDecimal   = 3 // D
    FlagIndex     = 4 // X - Index register width (native mode)
    FlagMemory    = 5 // M - Accumulator width (native mode)
    FlagOverflow  = 6 // V
    FlagNegative  = 7 // N
)

// Emulation flag (E) is separate, not part of P register
```

### Phase 2: Opcode Table and Categories

**Files to create:**
- `arch/cpu/m65816/opcode.go` -- Full 256-entry opcode table with width flag metadata
- `arch/cpu/m65816/categories.go` -- Instruction category sets for static analysis

The opcode table must encode which instructions are affected by M vs X flag for proper
disassembly and static analysis. This is critical because the instruction stream cannot be
decoded without tracking processor state.

```go
var Opcodes = [256]Opcode{
    // $00: BRK #imm (2 bytes in native, 2 bytes in emulation)
    {Instruction: Brk, Addressing: ImmediateAddressing, Timing: 7, BaseSize: 2, WidthFlag: WidthNone},
    // $01: ORA (dp,X)
    {Instruction: Ora, Addressing: DirectPageIndexedIndirectXAddressing, Timing: 6, BaseSize: 2, WidthFlag: WidthM},
    // $02: COP #imm
    {Instruction: Cop, Addressing: ImmediateAddressing, Timing: 7, BaseSize: 2, WidthFlag: WidthNone},
    // ... all 256 entries ...
}
```

### Phase 3: CPU State and Memory Model (24-bit)

**Files to create:**
- `arch/cpu/m65816/cpu.go` -- CPU state with dual-mode registers
- `arch/cpu/m65816/option.go` -- Functional options
- `arch/cpu/m65816/memory.go` -- 24-bit memory interface

```go
type Memory interface {
    ReadByte(address uint32) uint8
    WriteByte(address uint32, value uint8)
    // ReadWord reads a 16-bit value in little-endian byte order.
    ReadWord(address uint32) uint16
    WriteWord(address uint32, value uint16)
    // ReadLong reads a 24-bit value in little-endian byte order.
    ReadLong(address uint32) uint32
    WriteLong(address uint32, value uint32)
}

type CPU struct {
    // 16-bit accumulator C (accessed as full 16-bit or split A/B)
    C uint16

    X  uint16 // Index register X
    Y  uint16 // Index register Y
    SP uint16 // Stack pointer
    DP uint16 // Direct Page register
    DB uint8  // Data Bank register
    PB uint8  // Program Bank register
    PC uint16 // Program counter (within bank)

    P uint8   // Processor status register
    E bool    // Emulation mode flag

    cycles uint64
    halted bool
    stopped bool

    memory Memory
    opts   Options
}

// FullPC returns the 24-bit program address (PB:PC).
func (c *CPU) FullPC() uint32 {
    return uint32(c.PB)<<16 | uint32(c.PC)
}

// A returns the low byte of the accumulator.
func (c *CPU) A() uint8 { return uint8(c.C) }

// B returns the high byte of the accumulator.
func (c *CPU) B() uint8 { return uint8(c.C >> 8) }

// AccWidth returns the current accumulator width (1 or 2 bytes).
func (c *CPU) AccWidth() int {
    if c.E || c.P&MaskMemory != 0 { return 1 }
    return 2
}

// IdxWidth returns the current index register width (1 or 2 bytes).
func (c *CPU) IdxWidth() int {
    if c.E || c.P&MaskIndex != 0 { return 1 }
    return 2
}
```

### Phase 4: Instruction Execution

**Files to create:**
- `arch/cpu/m65816/step.go` -- Fetch/decode/execute cycle with mode awareness
- `arch/cpu/m65816/param.go` -- Parameter reading for all addressing modes
- `arch/cpu/m65816/emulation.go` -- Core ALU operations (ADC, SBC, AND, ORA, etc.)
- `arch/cpu/m65816/emulation_move.go` -- Data movement (LDA, STA, MVN, MVP, transfers)
- `arch/cpu/m65816/emulation_branch.go` -- Branch instructions (BRA, BRL, BCC, etc.)
- `arch/cpu/m65816/emulation_stack.go` -- Stack operations (PHA, PEA, PEI, PER, PHB, etc.)
- `arch/cpu/m65816/emulation_system.go` -- System instructions (REP, SEP, XCE, COP, etc.)
- `arch/cpu/m65816/interrupt.go` -- Interrupt/exception handling (different vectors for E/N)

Key complexity in Step:

```go
func (c *CPU) Step() error {
    opcode := c.memory.ReadByte(c.FullPC())
    c.PC++

    info := Opcodes[opcode]
    if info.Instruction == nil {
        return ErrInvalidOpcode
    }

    // Determine actual instruction size based on M/X flags
    size := int(info.BaseSize)
    switch info.WidthFlag {
    case WidthM:
        if c.AccWidth() == 2 { size++ }
    case WidthX:
        if c.IdxWidth() == 2 { size++ }
    }

    // Read operand bytes and execute
    // ...
}
```

**Interrupt vectors (native mode):**
- COP: $00:FFE4
- BRK: $00:FFE6
- ABORT: $00:FFE8
- NMI: $00:FFEA
- IRQ: $00:FFEE

**Interrupt vectors (emulation mode):**
- COP: $00:FFF4
- ABORT: $00:FFF8
- NMI: $00:FFFA
- RESET: $00:FFFC (emulation only)
- IRQ/BRK: $00:FFFE

### Phase 5: Emulation Core Refinement

- Cycle-accurate timing (add page-crossing penalties, mode-dependent cycles)
- Proper wrapping behavior (direct page wrapping in emulation mode)
- Bank boundary behavior (wrapping within bank vs crossing banks)
- WAI/STP instruction behavior
- Edge cases: REP/SEP immediate interaction with E flag

### Phase 6: SNES System Support

**Files to create (future, outside m65816 package):**
- SNES memory map (LoROM, HiROM, ExHiROM mapping)
- SNES DMA and HDMA support
- SNES interrupt timing (NMI on VBlank)
- ROM header parsing

This phase is out of scope for the CPU package itself but is noted as the primary
motivation for 65816 support.

### Phase 7: Testing

- **Opcode table completeness test** -- All 256 entries filled, no nil instructions
- **Instruction encoding roundtrip tests** -- Assemble/disassemble match
- **Per-instruction unit tests** -- Each instruction with both 8-bit and 16-bit modes
- **Mode switching tests** -- XCE, REP, SEP behavior and side effects
- **Addressing mode tests** -- Each mode with various DP, DB values
- **Block move tests** -- MVN/MVP with various bank configurations
- **Emulation mode compatibility** -- Verify 65C02 behavior in emulation mode
- **bsnes/higan accuracy tests** -- Compare against known-accurate SNES emulator behavior
- **Tom Harte processor tests** -- If available for 65816

---

## 8. File Structure

Following `cpu-architecture-guidelines.md`:

```
arch/cpu/m65816/
    doc.go              -- Package documentation
    addressing.go       -- ~38 addressing modes
    instruction.go      -- ~114 instruction definitions
    opcode.go           -- 256-entry opcode table with width metadata
    categories.go       -- Instruction category sets
    errors.go           -- Package-specific errors
    flag.go             -- P register flags + E flag
    cpu.go              -- CPU state, registers, mode queries
    option.go           -- Functional options
    memory.go           -- 24-bit memory interface
    step.go             -- Fetch/decode/execute cycle
    param.go            -- Operand reading for all addressing modes
    emulation.go        -- Core ALU instruction handlers
    emulation_move.go   -- Data movement handlers
    emulation_branch.go -- Branch/jump handlers
    emulation_stack.go  -- Stack operation handlers
    emulation_system.go -- System/mode instruction handlers
    interrupt.go        -- Interrupt/exception processing
```

---

## 9. Estimated Effort

| Component | New LOC |
|-----------|---------|
| Architecture/system registration | ~20 |
| Static analysis files (doc, addressing, instruction, opcode, categories, errors, flag) | ~2,500 |
| CPU state + options + memory interface | ~400 |
| Step/decode cycle + parameter reading | ~600 |
| Instruction handlers (5 emulation files) | ~3,500 |
| Interrupt/exception processing | ~300 |
| Tests | ~2,500 |
| **Total** | **~9,800** |

---

## 10. Reference Implementations

For accuracy validation and edge case guidance:

- **bsnes/higan** (C++, cycle-exact): The gold standard for SNES accuracy.
  Byuu's 65816 core is the most accurate known implementation.
- **Mesen-S** (C++): Highly accurate SNES emulator with good debugger support.
- **65816 Programming Manual** (WDC official): Definitive instruction reference
  with cycle counts and flag effects. Available from westerndesigncenter.com.
- **Eyes & Lichty, "Programming the 65816"** (book): Comprehensive reference
  covering all instructions, addressing modes, and edge cases.
- **TomHarte/ProcessorTests**: May include 65816 test vectors.

---

## 11. Design Decisions

### New Package vs Variant

**Decision: New package `arch/cpu/m65816/`**

Rationale: The 65816's 24-bit address space, variable-width registers, ~24 new addressing
modes, bank registers, and dual-mode operation represent a fundamental architectural change.
The code reuse with m6502 would be minimal (mainly instruction names and basic ALU logic),
while the conditionals required would make the code unmaintainable. A clean package allows
proper modeling of 24-bit memory, bank semantics, and mode-dependent behavior.

### Opcode Table Strategy

**Decision: Single [256]Opcode table with width metadata**

Rationale: Unlike the 68000 (which needs hierarchical decoding of 16-bit opcodes), the
65816 retains the 6502's single-byte opcode structure. All 256 slots are filled. The
variable instruction size is handled by metadata on each opcode entry indicating whether
the M or X flag affects its width, rather than separate tables per mode.

### Memory Interface

**Decision: uint32 addresses in the Memory interface**

Rationale: Although the 65816 has 24 address lines, using uint32 avoids a custom 24-bit
type and aligns with Go's native integer sizes. Implementations mask to 24 bits internally.
This matches the approach used for the 68000 package (which also uses uint32 for its
24-bit address bus).

### Emulation Mode Implementation

**Decision: Single CPU struct with mode-dependent behavior (not two separate implementations)**

Rationale: Emulation mode and native mode share the same opcode table and most instruction
logic. The differences are in register widths, stack behavior, and interrupt vectors.
Conditionals on the E flag at key points (register access, stack operations, interrupts)
are cleaner than duplicating the entire execution engine.

### Accumulator Model

**Decision: Single uint16 field (C) with accessors for A/B bytes**

Rationale: The 65816 treats the accumulator as a 16-bit register C that can be accessed
as two 8-bit halves (A = low, B = high). Using a single uint16 with byte accessors maps
directly to the hardware model and simplifies 16-bit operations. The XBA instruction
(exchange B and A) becomes a simple byte swap.
