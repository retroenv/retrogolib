# CPU Implementation Plan: WDC 65C02 and Motorola 68000

## Context

The retrogolib library currently supports 4 CPU architectures (6502, Z80, CHIP-8, x86-stub)
and 6 systems (NES, ZX Spectrum, Game Boy, CHIP-8, DOS, Generic). Adding the 65C02 and 68000
expands support to major retro platforms: Apple IIe/IIc, Atari Lynx, TurboGrafx-16 (65C02)
and Sega Genesis/Mega Drive, Amiga, Atari ST, Macintosh (68000).

---

## Part 1: WDC 65C02

### 1.1 Approach: Extend m6502 Package

The 65C02 is a strict superset of the NMOS 6502 with ~30 new instructions, 1 new addressing
mode, and behavioral fixes. Rather than creating a separate package, extend `arch/cpu/m6502/`
with a variant option. This mirrors how rust-z80emu handles NMOS/CMOS/BM1 as flavours of
the same CPU.

### 1.2 New Instructions (27 total)

**Stack operations (2):**
- PHX ($DA, implied, 3 cycles) -- Push X
- PHY ($5A, implied, 3 cycles) -- Push Y
- PLX ($FA, implied, 4 cycles) -- Pull X
- PLY ($7A, implied, 4 cycles) -- Pull Y

**Accumulator operations (2):**
- INC A ($1A, accumulator, 2 cycles) -- Increment A
- DEC A ($3A, accumulator, 2 cycles) -- Decrement A

**Branch (1):**
- BRA ($80, relative, 3+1 cycles) -- Branch Always

**Memory operations (4):**
- STZ ($64/$74/$9C/$9E, 4 modes, 3-5 cycles) -- Store Zero
- TRB ($14/$1C, zp/abs, 5-6 cycles) -- Test and Reset Bits
- TSB ($04/$0C, zp/abs, 5-6 cycles) -- Test and Set Bits

**New addressing modes for existing instructions (8 instructions):**
- ORA, AND, EOR, ADC, STA, LDA, CMP, SBC gain zero page indirect `(zp)`:
  - $12, $32, $52, $72, $92, $B2, $D2, $F2 (5 cycles each)
- BIT gains immediate ($89), zp,X ($34), abs,X ($3C)
- JMP gains absolute indexed indirect `(abs,X)` ($7C)

**Rockwell/WDC extensions (32 opcodes, optional):**
- RMB0-7 ($07/$17/$27/$37/$47/$57/$67/$77) -- Reset Memory Bit
- SMB0-7 ($87/$97/$A7/$B7/$C7/$D7/$E7/$F7) -- Set Memory Bit
- BBR0-7 ($0F/$1F/$2F/$3F/$4F/$5F/$6F/$7F) -- Branch on Bit Reset
- BBS0-7 ($8F/$9F/$AF/$BF/$CF/$DF/$EF/$FF) -- Branch on Bit Set

**WDC-only (2, optional):**
- STP ($DB) -- Stop the Processor
- WAI ($CB) -- Wait for Interrupt

### 1.3 Behavioral Changes

These require conditional logic based on CPU variant:

1. **Decimal mode flag clearing**: On BRK/IRQ/NMI/RESET, D flag is cleared (NMOS doesn't)
2. **Decimal mode ADC/SBC**: Flags N, V, Z set correctly (+1 cycle penalty)
3. **JMP ($xxFF) page boundary bug fixed**: Reads correct high byte from next page
4. **RMW instruction behavior**: Read-read-write instead of NMOS read-write-write
   (affects STA timing with hardware registers)
5. **Undocumented opcodes**: All 105 NMOS undocumented opcodes become NOPs
   (varying 1-4 cycles, 1-3 bytes)

### 1.4 Implementation Plan

#### Phase 1: Variant Infrastructure
**Files to modify:**
- `arch/cpu/m6502/option.go` -- Add `WithVariant(variant CPUVariant)` option
- `arch/cpu/m6502/cpu.go` -- Add `variant CPUVariant` field to CPU struct
- `arch/arch.go` -- Add `M65C02 Architecture = "65c02"` constant

```go
// option.go
type CPUVariant int
const (
    VariantNMOS6502 CPUVariant = iota // Original NMOS 6502
    Variant65C02                       // WDC 65C02 (base)
    Variant65C02Rockwell               // 65C02 + Rockwell bit extensions
    Variant65C02WDC                    // 65C02 + Rockwell + STP/WAI
)

func WithVariant(v CPUVariant) Option {
    return func(o *Options) { o.variant = v }
}
```

#### Phase 2: New Addressing Mode
**Files to modify:**
- `arch/cpu/m6502/addressing.go` -- Add `ZeroPageIndirectAddressing`
- `arch/cpu/m6502/param.go` -- Add parameter reader for `(zp)` mode
- `arch/cpu/m6502/memory.go` -- Add `ReadZeroPageIndirect()` method

```go
// addressing.go - add:
ZeroPageIndirectAddressing  // 65C02: (zp) - zero page indirect without indexing

// New type:
type ZeroPageIndirect uint8
```

#### Phase 3: New Instructions
**Files to create:**
- `arch/cpu/m6502/instruction_65c02.go` -- 65C02 instruction definitions
- `arch/cpu/m6502/emulation_65c02.go` -- 65C02 instruction handlers

**Files to modify:**
- `arch/cpu/m6502/opcode.go` -- Add 65C02 opcode table variant
- `arch/cpu/m6502/categories.go` -- Add new instructions to category sets

```go
// opcode.go - add:
var Opcodes65C02 = [256]Opcode{...} // Full 256-entry table with 65C02 additions

// instruction_65c02.go - new instructions:
var Bra = &Instruction{Name: BraName, Addressing: ...}
var Phx = &Instruction{Name: PhxName, Addressing: ...}
// ... etc
```

#### Phase 4: Behavioral Fixes
**Files to modify:**
- `arch/cpu/m6502/emulation.go` -- Conditional ADC/SBC decimal mode, RMW behavior
- `arch/cpu/m6502/interrupt.go` -- Clear D flag on interrupt (conditional on variant)
- `arch/cpu/m6502/memory.go` -- Fix `ReadWordBug()` to be correct on 65C02
- `arch/cpu/m6502/step.go` -- Select opcode table based on variant

#### Phase 5: Testing
- Klaus Dormann 65C02 functional test ROM
- Unit tests for each new instruction
- Unit tests for behavioral differences (decimal mode, JMP bug fix, etc.)
- Opcode table completeness test (all 256 entries defined for 65C02)

### 1.5 Estimated Effort

| Component | New LOC | Modified LOC |
|-----------|---------|-------------|
| Variant infrastructure | ~40 | ~30 |
| New addressing mode | ~30 | ~20 |
| New instructions + handlers | ~400 | ~20 |
| Behavioral fixes | ~60 | ~80 |
| Tests | ~500 | ~50 |
| **Total** | **~1,030** | **~200** |

---

## Part 2: Motorola 68000

### 2.1 Approach: New Package

The 68000 is architecturally distinct from all existing CPUs: 32-bit registers, 16-bit data
bus, 24-bit address bus, big-endian, privilege modes, complex exception model. It requires a
new package at `arch/cpu/m68000/`.

### 2.2 Architecture Overview

**Registers:**
- D0-D7: 8 x 32-bit data registers (support byte/word/long operations)
- A0-A6: 7 x 32-bit address registers
- A7/USP: User stack pointer
- SSP: Supervisor stack pointer (swapped with A7 in supervisor mode)
- PC: 32-bit program counter (24-bit external address bus)
- SR: 16-bit status register
  - Low byte (CCR): C, V, Z, N, X (extend) flags
  - High byte (system): T (trace), S (supervisor), I2-I0 (interrupt mask)

**Addressing modes (14):**
1. Data Register Direct: Dn
2. Address Register Direct: An
3. Address Register Indirect: (An)
4. Address Register Indirect with Postincrement: (An)+
5. Address Register Indirect with Predecrement: -(An)
6. Address Register Indirect with Displacement: d16(An)
7. Address Register Indirect with Index: d8(An,Xn)
8. Absolute Short: (xxx).W
9. Absolute Long: (xxx).L
10. Program Counter with Displacement: d16(PC)
11. Program Counter with Index: d8(PC,Xn)
12. Immediate: #imm
13. Status Register (implicit): SR, CCR
14. Quick Immediate (3-bit or 8-bit embedded in opcode)

**Operand sizes:**
- Byte (.B), Word (.W, 16-bit), Long (.L, 32-bit)
- Size encoded in opcode bits (usually bits 7-6)

**Instruction encoding:**
- 16-bit operation word minimum
- 0-4 extension words for source/destination effective addresses
- Top 4 bits = "line" (instruction group), rest = operation + EA encoding
- Effective address: 6 bits = 3 mode + 3 register

### 2.3 Instruction Set (~68 mnemonics)

**Data movement (11):** MOVE, MOVEA, MOVEQ, MOVEM, MOVEP, EXG, LEA, PEA, LINK, UNLK, SWAP

**Arithmetic (14):** ADD, ADDA, ADDI, ADDQ, ADDX, SUB, SUBA, SUBI, SUBQ, SUBX, MULU, MULS, DIVU, DIVS, NEG, NEGX, CLR, EXT

**Logical (8):** AND, ANDI, OR, ORI, EOR, EORI, NOT, TST

**Shift/Rotate (8):** ASL, ASR, LSL, LSR, ROL, ROR, ROXL, ROXR

**Bit manipulation (4):** BTST, BSET, BCLR, BCHG

**BCD (3):** ABCD, SBCD, NBCD

**Comparison (3):** CMP, CMPA, CMPI, CMPM

**Branch/Jump (8):** Bcc (14 conditions), BRA, BSR, DBcc, Scc, JMP, JSR, NOP

**System control (10):** TRAP, TRAPV, CHK, RTE, RTS, RTR, STOP, RESET, MOVE to/from SR/USP, ANDI/ORI/EORI to SR/CCR

### 2.4 Implementation Plan

#### Phase 1: Static Analysis Foundation (Required Files)
**Files to create:**
- `arch/cpu/m68000/doc.go` -- Architecture overview and usage
- `arch/cpu/m68000/addressing.go` -- 14 addressing modes + operand size types
- `arch/cpu/m68000/instruction.go` -- ~68 instruction definitions
- `arch/cpu/m68000/opcode.go` -- Opcode decoding (NOT a flat 256-entry table)
- `arch/cpu/m68000/categories.go` -- Instruction category sets
- `arch/cpu/m68000/errors.go` -- Package-specific errors
- `arch/cpu/m68000/flag.go` -- CCR/SR flag definitions

**Key design difference from 6502/Z80:** The 68000 uses variable-length instructions
(2-10 bytes) with a hierarchical opcode structure. A flat [65536]Opcode table would be
wasteful. Instead, use a two-level decode:

```go
// Level 1: Decode by line (top 4 bits of opcode word)
type lineDecoder func(opcode uint16) (*Instruction, AddressingMode, OperandSize)

var lineDecoders = [16]lineDecoder{
    decodeLine0, // ORI, ANDI, SUBI, ADDI, EORI, CMPI, BTST, BSET, BCLR, BCHG, MOVEP
    decodeLine1, // MOVE.B
    decodeLine2, // MOVE.L, MOVEA.L
    decodeLine3, // MOVE.W, MOVEA.W
    decodeLine4, // Miscellaneous (LEA, PEA, CHK, SWAP, EXT, TRAP, LINK, UNLK, etc.)
    decodeLine5, // ADDQ, SUBQ, Scc, DBcc
    decodeLine6, // Bcc, BRA, BSR
    decodeLine7, // MOVEQ
    decodeLine8, // OR, DIV, SBCD
    decodeLine9, // SUB, SUBA, SUBX
    decodeLineA, // Unassigned (Line A emulator trap)
    decodeLineB, // CMP, CMPA, CMPM, EOR
    decodeLineC, // AND, MUL, ABCD, EXG
    decodeLineD, // ADD, ADDA, ADDX
    decodeLineE, // Shift/Rotate
    decodeLineF, // Unassigned (Line F emulator trap)
}
```

**Architecture registration:**
- `arch/arch.go` -- Add `M68000 Architecture = "m68000"`
- `arch/system.go` -- Add systems: `SegaGenesis`, `Amiga`, `AtariST`, `Macintosh68k`

#### Phase 2: Effective Address Decoder
**Files to create:**
- `arch/cpu/m68000/ea.go` -- Effective address decoding and resolution

The EA decoder is the core of the 68000. It interprets the 6-bit mode+register field
and any extension words to compute the operand address.

```go
type OperandSize int
const (
    SizeByte OperandSize = 1
    SizeWord OperandSize = 2
    SizeLong OperandSize = 4
)

type EffectiveAddress struct {
    Mode     AddressingMode
    Register uint8       // 0-7
    Size     OperandSize
    Address  uint32      // Resolved address (for memory modes)
    Value    uint32      // Immediate or register value
}

// Decode reads the EA from the instruction stream and resolves it.
func (c *CPU) decodeEA(mode, reg uint8, size OperandSize) (EffectiveAddress, error)

// ReadEA reads the value at an effective address.
func (c *CPU) readEA(ea EffectiveAddress) (uint32, error)

// WriteEA writes a value to an effective address.
func (c *CPU) writeEA(ea EffectiveAddress, value uint32) error
```

#### Phase 3: Memory and Bus Interface
**Files to create:**
- `arch/cpu/m68000/memory.go` -- Memory interface (big-endian, 24-bit, word-aligned)
- `arch/cpu/m68000/bus.go` -- Bus interface with interrupt acknowledge

```go
type Memory interface {
    ReadByte(address uint32) uint8
    WriteByte(address uint32, value uint8)
    ReadWord(address uint32) uint16  // Must be word-aligned
    WriteWord(address uint32, value uint16)
    ReadLong(address uint32) uint32
    WriteLong(address uint32, value uint32)
}

type Bus interface {
    Memory
    // IRQLevel returns the current interrupt priority level (0-7).
    IRQLevel() uint8
    // IRQAcknowledge is called when the CPU acknowledges an interrupt.
    // Returns the vector number for the interrupt.
    IRQAcknowledge(level uint8) uint32
    // Reset is called when the CPU executes the RESET instruction.
    OnReset()
}
```

Key difference from Z80: 32-bit addresses, big-endian byte order, word alignment
requirement (odd-address word/long access = Address Error exception).

#### Phase 4: CPU State and Emulation Core
**Files to create:**
- `arch/cpu/m68000/cpu.go` -- CPU state, registers, privilege mode
- `arch/cpu/m68000/option.go` -- Functional options
- `arch/cpu/m68000/step.go` -- Instruction fetch/decode/execute cycle
- `arch/cpu/m68000/interrupt.go` -- Exception and interrupt processing
- `arch/cpu/m68000/param.go` -- Parameter/operand reading

```go
type CPU struct {
    mu sync.RWMutex

    D [8]uint32  // Data registers D0-D7
    A [7]uint32  // Address registers A0-A6
    USP uint32   // User stack pointer
    SSP uint32   // Supervisor stack pointer
    PC  uint32   // Program counter

    // Status register
    Flags CCR    // Condition code register (C, V, Z, N, X)
    SR    uint16 // Full status register (includes system byte)

    cycles uint64
    halted bool
    stopped bool // STP instruction state

    bus Bus
    opts Options
}

// A7 returns the active stack pointer based on privilege mode.
func (c *CPU) A7() uint32 {
    if c.SR&FlagSupervisor != 0 {
        return c.SSP
    }
    return c.USP
}
```

#### Phase 5: Instruction Handlers
**Files to create:**
- `arch/cpu/m68000/emulation.go` -- Core ALU operations (ADD, SUB, AND, OR, etc.)
- `arch/cpu/m68000/emulation_move.go` -- Data movement (MOVE, MOVEM, LEA, etc.)
- `arch/cpu/m68000/emulation_branch.go` -- Branch and jump instructions
- `arch/cpu/m68000/emulation_shift.go` -- Shift and rotate instructions
- `arch/cpu/m68000/emulation_bit.go` -- Bit manipulation instructions
- `arch/cpu/m68000/emulation_system.go` -- Privileged and system instructions

Instruction handlers follow the same signature pattern as Z80/6502:
```go
func moveWord(c *CPU, params ...any) error { ... }
func addLong(c *CPU, params ...any) error { ... }
```

#### Phase 6: Exception Processing
**File:** `arch/cpu/m68000/interrupt.go`

The 68000 exception model is significantly more complex than Z80/6502:

```go
const (
    VectorReset          = 0   // Reset SSP and PC (vectors 0-1)
    VectorBusError       = 2
    VectorAddressError   = 3
    VectorIllegal        = 4
    VectorDivideByZero   = 5
    VectorCHK            = 6
    VectorTRAPV          = 7
    VectorPrivilege      = 8
    VectorTrace          = 9
    VectorLineA          = 10
    VectorLineF          = 11
    VectorAutoVector     = 25  // Vectors 25-31 for interrupt levels 1-7
    VectorTrap0          = 32  // Vectors 32-47 for TRAP #0-#15
    VectorUser           = 64  // Vectors 64-255 for user interrupts
)

func (c *CPU) processException(vector uint8) error {
    // 1. Save current SR
    // 2. Set supervisor mode (S bit)
    // 3. Clear trace bit (T)
    // 4. For interrupts: set interrupt mask
    // 5. Push PC and SR to supervisor stack
    // 6. For bus/address errors: push additional state
    // 7. Load PC from vector table (vector * 4)
}
```

#### Phase 7: Testing
- TomHarte/SingleStepTests JSON test vectors (generated from MAME/Musashi)
- Unit tests for each instruction with all operand sizes
- Effective address decoding tests
- Exception processing tests
- Privilege mode tests
- Alignment error tests

### 2.5 Estimated Effort

| Component | New LOC |
|-----------|---------|
| Static analysis files (doc, addressing, instruction, opcode, categories, errors, flag) | ~1,500 |
| Effective address decoder | ~500 |
| Memory/Bus interfaces | ~200 |
| CPU state + options | ~400 |
| Step/decode cycle | ~400 |
| Instruction handlers (6 files) | ~3,000 |
| Exception/interrupt processing | ~300 |
| Tests | ~2,000 |
| **Total** | **~8,300** |

### 2.6 Reference Implementations

For accuracy validation and edge case guidance:
- **Musashi** (C, MAME): The gold standard. Git: `mamedev/mame` at `src/devices/cpu/m68000/`
- **Moira** (C++, cycle-exact): Modern, faster than Musashi. Git: `dirkwhoffmann/Moira`
- **TomHarte/ProcessorTests**: JSON test vectors for bus-level validation

---

## Part 3: Architecture and System Registration

### New Architecture Constants
```go
// arch/arch.go - add:
M65C02 Architecture = "65c02"
M68000 Architecture = "m68000"
```

### New System Constants
```go
// arch/system.go - add:
AppleII      System = "apple-ii"       // Apple IIe/IIc (65C02)
AtariLynx    System = "atari-lynx"     // Atari Lynx (65C02)
TurboGrafx   System = "turbografx"     // TurboGrafx-16/PC Engine (65C02 variant)
SegaGenesis  System = "sega-genesis"   // Sega Genesis/Mega Drive (68000)
Amiga        System = "amiga"          // Commodore Amiga (68000)
AtariST      System = "atari-st"       // Atari ST (68000)
Macintosh68k System = "macintosh-68k"  // Original Macintosh (68000)
```

---

## Part 4: Implementation Order

### Recommended Sequence

**Step 1: 65C02 (2-3 weeks effort)**
- Smallest delta, extends proven codebase
- Unlocks Apple II, Atari Lynx, TurboGrafx-16
- Validates the variant pattern for future use

**Step 2: 68000 Static Analysis (1-2 weeks)**
- Instruction definitions, opcode decoding, categories
- No emulation yet -- useful for disassemblers/analyzers
- Follows existing x86 stub pattern

**Step 3: 68000 Emulation (4-6 weeks)**
- Full instruction execution
- Exception handling
- Test suite integration

### Quality Gates

Each step must pass before proceeding:

1. **65C02**: All existing 6502 tests still pass + Klaus Dormann 65C02 test ROM passes
2. **68000 Static**: Opcode decode roundtrip tests + instruction completeness tests
3. **68000 Emulation**: TomHarte SingleStep JSON test vectors pass

---

## Part 5: Key Design Decisions

### 65C02: Variant vs Separate Package
**Decision: Variant within m6502**
- Rationale: 95% code reuse, only ~30 new instructions, behavioral changes are small
  conditionals. A separate package would duplicate thousands of lines.

### 68000: Opcode Decode Strategy
**Decision: Line-based hierarchical decoder (NOT flat 65536-entry table)**
- Rationale: The 68000 has a 16-bit opcode word with structured encoding. A flat
  table would be 65536 entries, most empty. A two-level decoder (line -> specific)
  is both smaller and more maintainable.

### 68000: Memory Model
**Decision: 32-bit addresses in the Memory interface**
- Rationale: Although the 68000 only has 24 address lines, the registers and
  instruction set are 32-bit. The Memory interface should use uint32. The
  implementation can mask to 24 bits internally.

### 68000: Endianness
**Decision: Memory interface handles big-endian natively**
- Rationale: The 68000 is big-endian. Unlike the Z80/6502 where Memory.ReadWord()
  is little-endian, the 68000 Memory.ReadWord() returns big-endian word. This avoids
  byte-swapping on every memory access.
