# Work Branch Changelog

Tracks every file changed on the `work` branch compared to `main`. This document must be
kept up to date as changes are made. Features will be extracted from `work` to `main`
individually.

**Last Updated:** 2026-04-27

---

## Merge Plan to Main

Organize changes into logical, testable groups for merging to `main`. Each group should be
independent and maintain a working codebase.

### Group 1: Foundation & Test Infrastructure ✅
**Status:** Ready to merge
**Dependencies:** None

- `testdata/Makefile` - Test ROM download automation
- `testdata/.gitignore` - Git ignore rules for test data
- `Makefile` - Test target changes (`-short` flag, `test-integration` target)

**Why first:** Enables test infrastructure without changing CPU code. Tests can run
with existing test ROMs once downloaded.

---

### Group 2: Architecture Registration ✅
**Status:** Ready to merge
**Dependencies:** None

- `arch/arch.go` - Add `M65C02`, `M65816`, `M68000`, `SM83`, `X86` constants
- `arch/arch_test.go` - Update test expectations
- `arch/system.go` - Add `AppleIIGS`, `Atari2600`, `SNES` constants
- `arch/system_test.go` - Update test expectations

**Why second:** Pure constants, no behavioral changes. Safe to merge early.

---

### Group 3: 65C02 CPU Implementation ✅
**Status:** Ready to merge
**Dependencies:** Group 2 (architecture constants)

**New Files:**
- `arch/cpu/m6502/emulation_65c02.go` - 65C02 instruction handlers
- `arch/cpu/m6502/emulation_65c02_test.go` - 65C02 unit tests
- `arch/cpu/m6502/instruction_65c02.go` - 65C02 instruction definitions
- `arch/cpu/m6502/opcode_65c02.go` - Complete 65C02 opcode table
- `arch/cpu/m6502/emulation_6507_test.go` - 6507 variant tests

**Modified Files:**
- `arch/cpu/m6502/addressing.go` - New addressing modes
- `arch/cpu/m6502/categories.go` - Add 65C02 instructions
- `arch/cpu/m6502/cpu.go` - Add `branchTaken` field
- `arch/cpu/m6502/emulation.go` - 65C02 decimal mode fixes
- `arch/cpu/m6502/interrupt.go` - 65C02 D flag clearing
- `arch/cpu/m6502/instruction.go` - Extended instruction definitions
- `arch/cpu/m6502/memory.go` - New memory helpers
- `arch/cpu/m6502/opcode.go` - Reorganized NMOS table
- `arch/cpu/m6502/opcode_test.go` - 65C02 table tests
- `arch/cpu/m6502/option.go` - CPU variant support
- `arch/cpu/m6502/param.go` - New parameter readers
- `arch/cpu/m6502/step.go` - Opcode table selection

**Test Files (build tags):**
- `arch/cpu/m6502/dormann_test.go` - Klaus Dormann test suite
- `arch/cpu/m6502/singlestep_test.go` - SingleStepTests integration

**Why third:** Builds on architecture constants. 65C02 is backward-compatible with 6502.

---

### Group 4: Z80 Bus Interface & Undocumented Instructions ✅
**Status:** Ready to merge
**Dependencies:** None (internal Z80 refactor)

**New Files:**
- `arch/cpu/z80/emulation_dd_undoc.go` - Undocumented DD prefix handlers
- `arch/cpu/z80/emulation_fd_undoc.go` - Undocumented FD prefix handlers
- `arch/cpu/z80/emulation_helpers.go` - Shared helpers
- `arch/cpu/z80/emulation_index.go` - IX/IY shared logic
- `arch/cpu/z80/emulation_jump.go` - Jump instruction handlers
- `arch/cpu/z80/emulation_load.go` - Load instruction handlers
- `arch/cpu/z80/instruction_dd_undoc.go` - DD instruction metadata
- `arch/cpu/z80/instruction_fd_undoc.go` - FD instruction metadata
- `arch/cpu/z80/opcode_id.go` - Type-safe opcode constants

**Test Files:**
- `arch/cpu/z80/singlestep_test.go` - SingleStepTests integration
- `arch/cpu/z80/zexall_test.go` - ZEXDOC/ZEXALL validation

**Modified Files:**
- `arch/cpu/z80/cpu.go` - MEMPTR/Q registers, Bus interface
- `arch/cpu/z80/memory.go` - Bus interface
- `arch/cpu/z80/step.go` - MEMPTR tracking
- `arch/cpu/z80/emulation*.go` - Refactored handlers
- `arch/cpu/z80/instruction*.go` - Updated definitions
- `arch/cpu/z80/opcode.go` - Table reorganization

**Why now:** Major refactor but self-contained. Improves Z80 accuracy to 100%.

---

### Group 5: 65816 CPU (SNES/Apple IIGS) ✅
**Status:** Complete
**Dependencies:** Group 2 (architecture constants)

**Entirely New Package:** `arch/cpu/m65816/`
- All files new (22 files total)
- 16-bit successor to 65C02
- Emulation and native modes
- SNES and Apple IIGS support

**Test Files:**
- `arch/cpu/m65816/singlestep_test.go` - SingleStepTests integration

**Why here:** Builds on 65C02 concepts but independent implementation.

---

### Group 6: Motorola 68000 CPU (Genesis/Amiga/Atari ST) ✅
**Status:** Complete
**Dependencies:** Group 2 (architecture constants)

**Entirely New Package:** `arch/cpu/m68000/`
- All files new (27 files total)
- 32-bit CISC with 16-bit data bus
- Big-endian memory model
- 14 addressing modes
- Exception processing with 256 vectors

**Test Files:**
- `arch/cpu/m68000/singlestep_test.go` - SingleStepTests integration

**Why here:** Independent architecture, large package, test with SingleStepTests.

---

### Group 7: x86 Real Mode CPU (8086-80486) ✅
**Status:** Complete
**Dependencies:** None

**Entirely New Package:** `arch/cpu/x86/`
- All files new (14 files total)
- Real mode only (1MB address space)
- Segmented memory model
- ModR/M byte decoding
- 8086 base through 80486 extensions

**Why here:** Static analysis focus (not emulation). Independent of other CPUs.

---

### Group 8: SM83 CPU (Game Boy) ✅
**Status:** Complete
**Dependencies:** Group 2 (architecture constants)

**Entirely New Package:** `arch/cpu/sm83/`
- All files new (17 files total)
- Game Boy/Game Boy Color CPU
- Z80-derived but distinct
- 4 flags only (Z, N, H, C)
- Game Boy-specific instructions (STOP, SWAP, LDH)

**Test Files:**
- `arch/cpu/sm83/singlestep_test.go` - SingleStepTests integration

**Why here:** Independent CPU, well-tested with SingleStepTests.

---

### Group 9: Atari 2600 System Package
**Status:** Ready to merge
**Dependencies:** Group 2 (architecture constants), 6507 variant from Group 3

**New Package:** `arch/system/atari2600/`
- `atari2600.go` - Memory map constants
- `atari2600_test.go` - Memory map tests
- `cartridge/cartridge.go` - Bank switching logic
- `cartridge/cartridge_test.go` - Cartridge tests
- `register/tia.go` - TIA video/audio registers
- `register/riot.go` - RIOT I/O registers
- `register/register_test.go` - Register completeness

**Why here:** System-level package, depends on 6507 (6502 variant).

---

### Group 10: Documentation Updates
**Status:** Ongoing
**Dependencies:** None

**New Documentation:**
- `cpu-implementation-plan-65816.md`
- `cpu-implementation-plan-65c02-68000.md`
- `cpu-implementation-plan-sm83.md`
- `m68000-emulator-comparison.md`
- `m68000-gap-closure-plan.md`
- `supported-systems.md`
- `system-implementation-plan-atari2600.md`
- `system-implementation-plan-c64.md`
- `system-implementation-plan-gameboy.md`
- `z80-emulator-comparison.md`
- `z80-emulator-comparison-cross-language.md`
- `z80-gap-closure-plan.md`

**Why last:** Documentation can be merged anytime, but best to keep with related code.

---

## Merge Order Summary

1. ✅ **Group 1:** Test infrastructure (Makefile, testdata)
2. ✅ **Group 2:** Architecture constants (`arch/`)
3. ✅ **Group 3:** 65C02 CPU (extends 6502)
4. ✅ **Group 4:** Z80 Bus/MEMPTR (internal refactor)
5. ✅ **Group 5:** 65816 CPU (SNES/Apple IIGS)
6. ✅ **Group 6:** 68000 CPU (Genesis/Amiga/Atari ST)
7. ✅ **Group 7:** x86 Real Mode (8086-80486)
8. ✅ **Group 8:** SM83 CPU (Game Boy)
9. 🔄 **Group 9:** Atari 2600 System (next to merge)
10. 📝 **Group 10:** Documentation (merge anytime)

**Legend:** ✅ Complete | 🔄 Next Priority | 📝 Anytime

---

---

## Build & Configuration

### `Makefile`
- **What:** Changed `test` target to use `-short` flag; added `test-integration` target
- **Why:** Integration tests (SingleStepTests, ZEXDOC, ZEXALL) take minutes to run and
  should not block normal `make test`. They run via `make test-integration` separately

---

## Architecture & System Registration

### `arch/arch.go`
- **What:** Added architecture constants: `M65C02`, `M65816`, `M68000`, `SM83`
- **Why:** Registers the four new CPU architectures for validation and lookup

### `arch/arch_test.go`
- **What:** Updated test count expectation for new architectures
- **Why:** Keeps the architecture completeness test passing

### `arch/system.go`
- **What:** Added system constants: `AppleIIGS`, `Atari2600`, `SNES`
- **Why:** Registers SNES, Apple IIGS, and Atari 2600 as supported systems

### `arch/system_test.go`
- **What:** Updated test expectations for new systems
- **Why:** Keeps the system completeness test passing

---

## MOS 6502 / 65C02 (`arch/cpu/m6502/`)

### `addressing.go`
- **What:** Added 3 new addressing modes: `ZeroPageIndirectAddressing`,
  `AbsoluteXIndirectAddressing`, `ZeroPageRelativeAddressing`
- **Why:** Required by 65C02 instructions: `(zp)` indirect, `(abs,X)` for JMP, and
  `zp,rel` for BBR/BBS Rockwell bit-branch instructions

### `categories.go`
- **What:** Added 65C02 instructions to category sets (BRA to branching, STZ/TRB/TSB
  to memory write, etc.)
- **Why:** Static analysis category sets must include all variant instructions

### `cpu.go`
- **What:** Added `branchTaken` field to CPU struct
- **Why:** Tracks branch outcomes for cycle-accurate timing (branch taken adds +1 cycle)

### `interrupt.go`
- **What:** Added 6507 variant guards to `TriggerIrq()` and `TriggerNMI()`
- **Why:** The MOS 6507 has no IRQ or NMI pins; these must be no-ops for this variant

### `emulation_6507_test.go` (new)
- **What:** 6507 variant tests: 11 opcode execution cases, IRQ/NMI no-op verification,
  opcode table selection, interrupt behavior comparison against NMOS 6502
- **Why:** Validates that Variant6507 correctly uses the NMOS opcode table and rejects interrupts

### `dormann_test.go` (new)
- **What:** Integration test runner for Klaus Dormann 6502/65C02 functional test ROMs
- **Why:** Industry-standard CPU validation suite, runs the full test ROM to completion
  and verifies the success trap address is reached

### `emulation.go`
- **What:** Extended ADC/SBC with correct 65C02 decimal mode flag behavior; added
  `hasAccumulatorParam` helper; refactored interrupt flag clearing for 65C02 variant
- **Why:** 65C02 fixes decimal mode N/V/Z flags (NMOS gets them wrong); BRK/IRQ/NMI
  clear D flag on 65C02 (NMOS doesn't)

### `emulation_65c02.go` (new)
- **What:** Instruction handlers for all 65C02-specific instructions
- **Why:** Implements BRA, STZ, TRB, TSB, PHX/PHY/PLX/PLY, INC A, DEC A, BIT immediate,
  BIT zp,X / abs,X, JMP (abs,X), and Rockwell extensions (RMB, SMB, BBR, BBS)

### `emulation_65c02_test.go` (new)
- **What:** Unit tests for all 65C02-specific instructions and behavioral changes
- **Why:** Validates correctness of each new instruction and decimal mode flag fixes

### `errors_test.go`
- **What:** Updated error test for new variant-related errors
- **Why:** Keeps error handling tests current

### `instruction.go`
- **What:** Added instruction name constants for new 65C02 instructions; extended existing
  instructions (BIT, INC, DEC, JMP) with new addressing mode entries
- **Why:** Instruction definitions must include all valid addressing modes for each variant

### `instruction_65c02.go` (new)
- **What:** Full 65C02 instruction definitions with addressing mode maps
- **Why:** Defines all new 65C02 instructions (BRA, STZ, TRB, TSB, PHX, PHY, PLX, PLY,
  WAI, STP, RMB0-7, SMB0-7, BBR0-7, BBS0-7) and extended versions of existing instructions
  (BIT, INC, DEC, JMP with new addressing modes)

### `interrupt.go`
- **What:** Added D flag clearing on interrupt for 65C02 variant
- **Why:** 65C02 behavioral fix: BRK/IRQ/NMI clear the decimal flag (NMOS 6502 doesn't)

### `memory.go`
- **What:** Added `ReadZeroPageIndirect`, `ReadAbsoluteXIndirect` memory helpers
- **Why:** New addressing modes for 65C02 need memory access functions that handle the
  indirect reads correctly (zero page indirect wraps within page 0)

### `opcode.go`
- **What:** Reorganized NMOS opcode table entries (formatting, consistency)
- **Why:** Cleanup pass for consistency before adding 65C02 table

### `opcode_65c02.go` (new)
- **What:** Full 256-entry opcode table for the 65C02 variant
- **Why:** The 65C02 replaces all undocumented NMOS opcodes with NOPs and adds new
  instructions in those slots; requires a separate complete table

### `opcode_test.go`
- **What:** Added opcode table completeness test for 65C02 table
- **Why:** Validates all 256 entries are defined in the 65C02 table

### `option.go`
- **What:** Added `CPUVariant` type with `VariantNMOS6502`, `VariantNES6502`,
  `Variant6507`, `Variant6510`, `Variant65C02` constants; added `WithVariant()` option function
- **Why:** Enables selecting CPU variant at construction time; the NES variant (2A03)
  disables decimal mode, the 65C02 variant uses the extended opcode table, 6507/6510 for
  Atari 2600 and Commodore 64 respectively

### `param.go`
- **What:** Added parameter readers for `ZeroPageIndirectAddressing`,
  `AbsoluteXIndirectAddressing`, `ZeroPageRelativeAddressing`
- **Why:** Each new addressing mode needs a function to read its operand bytes and
  compute the effective address

### `singlestep_test.go` (new)
- **What:** Integration test runner for SingleStepTests/65x02 JSON test vectors
- **Why:** Validates every opcode against the SingleStepTests reference (generated from
  known-accurate emulators). Tests initial state, final state, and memory side effects

### `step.go`
- **What:** Opcode table selection based on variant (`>= Variant65C02` uses `Opcodes65C02`);
  added branch cycle penalty tracking
- **Why:** The step function must use the correct opcode table for the active variant and
  accurately account for branch timing

---

## WDC 65C816 (`arch/cpu/m65816/`) -- All New

Entire package is new. Implements the 16-bit successor to the 65C02.

### `doc.go`
- Architecture overview, modes (emulation/native), usage examples

### `addressing.go`
- ~24 addressing modes as typed constants (stack-relative, long, block move, etc.)

### `instruction.go`
- ~114 instruction definitions (all 6502 + 65C02 + 28 new 65816 instructions)

### `opcode.go`
- Full 256-entry opcode table with `WidthFlag` metadata (WidthM/WidthX/WidthNone) for
  variable-size instructions based on M and X processor flags

### `opcode_test.go`
- Opcode table completeness test (all 256 entries filled)

### `categories.go`
- Instruction category sets for static analysis

### `errors.go`
- Package-specific errors

### `flag.go`
- Processor status flags including M (memory/accumulator width), X (index width),
  E (emulation mode). Flags struct with width query methods (`AccWidth`, `IdxWidth`)

### `cpu.go`
- CPU state: 16-bit accumulator C (split A/B access), X, Y, SP, DP (direct page),
  DB (data bank), PB (program bank), PC, P status, E emulation flag
- `FullPC()` returns 24-bit PB:PC address
- Mode-dependent register width helpers

### `option.go`
- Functional options (tracing, pre-execution hooks)

### `memory.go`
- 24-bit memory interface using `uint32` addresses (masked to 24 bits)
- `ReadByte`, `WriteByte`, `ReadWord`, `WriteWord`, `ReadLong`, `WriteLong`

### `step.go`
- Fetch/decode/execute cycle with M/X flag-dependent instruction sizing
- Branch cycle penalty (emulation mode only)

### `param.go`
- Operand readers for all ~24 addressing modes, handling bank wrapping, direct page
  offset, and variable-width immediates

### `emulation.go`
- Core ALU: ADC, SBC (including BCD mode for 8-bit and 16-bit), AND, ORA, EOR, CMP,
  CPX, CPY, BIT, INC, DEC, ASL, LSR, ROL, ROR, TSB, TRB

### `emulation_branch.go`
- Branch instructions: BCC, BCS, BEQ, BNE, BMI, BPL, BVC, BVS, BRA, BRL
- Jump: JMP, JML, JSR, JSL, RTS, RTL

### `emulation_move.go`
- Data movement: LDA, LDX, LDY, STA, STX, STY, STZ
- Register transfers: TAX, TAY, TXA, TYA, TSX, TXS, TCD, TCS, TDC, TSC, TXY, TYX
- Block moves: MVN, MVP (with cycle-budget cap for step-based execution)
- XBA (exchange A and B bytes of accumulator)

### `emulation_stack.go`
- Stack operations: PHA, PHX, PHY, PHP, PLA, PLX, PLY, PLP
- New 65816 stack: PHB, PHD, PHK, PLB, PLD, PEA, PEI, PER

### `emulation_system.go`
- Mode control: REP, SEP, XCE (exchange carry and emulation flags)
- System: NOP, BRK, COP, STP, WAI, WDM

### `interrupt.go`
- Dual interrupt vectors (native mode at $00:FFE0-$00:FFEF, emulation at $00:FFF0-$00:FFFF)
- NMI and IRQ dispatch with mode-dependent stack frame format

### `cpu_test.go`
- Unit tests for CPU state, mode switching, register width queries

### `emulation_test.go`
- Comprehensive unit tests: JSR/RTS, JSL/RTL, JMP, JML, PEA/PEI/PER, MVN/MVP, BRK,
  ADC/SBC (binary and BCD in 8-bit and 16-bit modes), mode switch sequences,
  bank boundary crossing, WAI+NMI dispatch

### `singlestep_test.go`
- Integration test runner for SingleStepTests/65816 (512 test files, 512,000+ test cases)

---

## Motorola 68000 (`arch/cpu/m68000/`) -- All New

Entire package is new. Implements the 32-bit CISC processor with 16-bit data bus.

### `doc.go`
- Architecture overview, register set, addressing modes

### `addressing.go`
- 14 addressing modes + operand size types (Byte/Word/Long)

### `instruction.go`
- ~75 instruction definitions with addressing mode maps

### `opcode.go`
- Line-based hierarchical decoder: 16 line decoders (top 4 bits of 16-bit opcode word),
  each handling a logical instruction group
- Effective address mode/register extraction from 6-bit EA field

### `opcode_test.go`
- Opcode decode tests for all 16 lines, verifying instruction identification

### `categories.go`
- Instruction category sets for static analysis

### `errors.go`
- Package-specific errors

### `flag.go`
- CCR flags (C, V, Z, N, X) and SR system byte (T, S, I2-I0)
- `updateFlags` helper with per-flag control (set/clear/unchanged/calculated)

### `cpu.go`
- CPU state: D0-D7 (8x32-bit data), A0-A6 (7x32-bit address), USP/SSP (dual stack
  pointers), PC, SR. `A7()` returns active stack pointer based on privilege mode

### `option.go`
- Functional options (tracing, pre-execution hooks)

### `memory.go`
- Big-endian memory interface using `uint32` addresses
- `ReadByte`, `WriteByte`, `ReadWord`, `WriteWord`, `ReadLong`, `WriteLong`
- `BasicMemory` implementation (16 MB flat)

### `ea.go`
- Effective address decoder: interprets 6-bit mode+register field, reads extension words,
  computes operand addresses for all 14 modes

### `ea_test.go`
- Tests for each addressing mode's EA resolution

### `step.go`
- Fetch/decode/execute cycle: reads 16-bit opcode word, dispatches via line decoder

### `emulation.go`
- Core ALU: ADD, ADDA, ADDI, ADDQ, ADDX, SUB, SUBA, SUBI, SUBQ, SUBX, MULU, MULS,
  DIVU, DIVS, NEG, NEGX, CLR, EXT, AND, ANDI, OR, ORI, EOR, EORI, NOT, TST, CMP,
  CMPA, CMPI, CMPM, ABCD, SBCD, NBCD

### `emulation_move.go`
- Data movement: MOVE, MOVEA, MOVEQ, MOVEM, MOVEP, EXG, LEA, PEA, LINK, UNLK, SWAP,
  Scc, TAS

### `emulation_branch.go`
- Branch: Bcc (14 conditions), BRA, BSR, DBcc, JMP, JSR, RTS, RTR, NOP

### `emulation_shift.go`
- Shift/rotate: ASL, ASR, LSL, LSR, ROL, ROR, ROXL, ROXR (register and memory forms)

### `emulation_bit.go`
- Bit manipulation: BTST, BSET, BCLR, BCHG (register and immediate bit number)

### `emulation_system.go`
- System: TRAP, TRAPV, CHK, RTE, STOP, RESET, MOVE to/from SR/USP,
  ANDI/ORI/EORI to SR/CCR, illegal instruction trap

### `interrupt.go`
- 256-vector exception model: reset, bus/address error, illegal instruction, divide by
  zero, CHK, TRAPV, privilege violation, trace, Line A/F traps, auto-vectored interrupts,
  TRAP #0-15

### `memory_test.go`
- Tests for big-endian memory operations

### `cpu_test.go`
- Tests for CPU state, privilege mode, stack pointer switching

### `emulation_test.go`
- Comprehensive instruction tests with all three operand sizes

### `singlestep_test.go`
- Integration test runner for SingleStepTests/680x0 JSON test vectors

---

## Zilog Z80 (`arch/cpu/z80/`)

### `cpu.go`
- **What:** Added `MEMPTR` (WZ) register, `q` register for SCF/CCF X/Y flag tracking,
  `lastWasLdAIR` flag for LD A,{I|R} interrupt bug. Replaced `memory Memory` field with
  `bus Bus`. Added `NewWithBus()` constructor. Added `Bus()` accessor. Changed all
  internal `c.memory.*` calls to `c.bus.*`
- **Why:** MEMPTR/Q are undocumented internal registers needed for 100% zexall pass rate.
  Bus interface consolidates Memory + IOHandler + interrupt acknowledge into one interface

### `cpu_test.go`
- **What:** Updated tests for Bus interface, added MEMPTR/Q state verification
- **Why:** Tests must validate new register tracking and Bus-based construction

### `memory.go`
- **What:** Added `Bus` interface (extends Memory with `ReadPort`, `WritePort`, `IRQData`,
  `OnRETI`), added `legacyBusAdapter` for backward compatibility
- **Why:** Unified bus interface enables full 16-bit I/O port addressing, interrupt data bus
  emulation (IM 0/IM 2), and RETI notification for daisy chains

### `interrupt.go`
- **What:** Changed internal memory/port calls to use `c.bus.*`
- **Why:** Part of Memory-to-Bus migration

### `option.go`
- **What:** Minor cleanup
- **Why:** Consistency with Bus interface changes

### `step.go`
- **What:** Major refactor: added MEMPTR updates throughout instruction execution, added Q
  register tracking (capture flags after each instruction, reset for non-flag ops), added
  `lastWasLdAIR` handling, updated all memory/port calls to use Bus
- **Why:** Correct undocumented flag behavior requires MEMPTR/Q tracking at every instruction
  boundary. Bus migration is pervasive

### `emulation.go`
- **What:** Major refactor: extracted load instructions to `emulation_load.go`, jump
  instructions to `emulation_jump.go`, shared index helpers to `emulation_index.go`,
  utility helpers to `emulation_helpers.go`. Added MEMPTR updates to all memory/jump/IO
  instructions. Changed all `c.memory.*` to `c.bus.*`
- **Why:** File was too large (~1,800 lines). MEMPTR must be updated on every memory access,
  jump, call, return, I/O, and block instruction for correct undocumented flag behavior

### `emulation_cb.go`
- **What:** Updated memory calls to use `c.bus.*`, added MEMPTR-based X/Y flags for
  BIT n,(HL)
- **Why:** BIT n,(HL) takes X/Y flags from MEMPTR high byte, not the tested value

### `emulation_dd.go`
- **What:** Extracted shared IX/IY instruction logic to `emulation_index.go`, extracted
  undocumented DD instructions to `emulation_dd_undoc.go`, updated memory calls
- **Why:** Reduces code duplication between DD and FD files. Separates undocumented
  instruction handling for clarity

### `emulation_dd_undoc.go` (new)
- **What:** Undocumented DD-prefix instruction handlers (IXH/IXL arithmetic, DD+CB
  indexed operations with register copy)
- **Why:** Clean separation of undocumented behavior from standard instructions

### `emulation_ed.go`
- **What:** Added MEMPTR updates to all ED-prefix instructions (block transfers, I/O
  block ops, LD (nn) pairs). Added `lastWasLdAIR` tracking for LD A,I and LD A,R.
  Updated port operations to use 16-bit Bus addresses (B<<8|C)
- **Why:** ED instructions are the primary users of MEMPTR (LDI/LDD/CPI/CPD/INI/IND
  all update MEMPTR). LD A,{I|R} must set the flag for the NMOS interrupt bug. Port
  operations must pass the full 16-bit address per Z80 hardware behavior

### `emulation_fd.go`
- **What:** Same changes as emulation_dd.go but for IY (FD prefix). Extracted shared
  logic to `emulation_index.go`, undocumented instructions to `emulation_fd_undoc.go`
- **Why:** Same reasons as DD refactor

### `emulation_fd_undoc.go` (new)
- **What:** Undocumented FD-prefix instruction handlers (IYH/IYL arithmetic)
- **Why:** Clean separation of undocumented behavior

### `emulation_helpers.go` (new)
- **What:** Shared arithmetic and flag helpers extracted from emulation.go
- **Why:** Reduces file size, avoids duplication between base and prefix instruction files

### `emulation_index.go` (new)
- **What:** Shared IX/IY instruction implementations (LD, ADD, ADC, SBC, INC, DEC, etc.
  on index registers and indexed memory)
- **Why:** DD and FD prefix instructions are identical except for which index register
  they use. Shared implementations eliminate ~700 lines of duplication

### `emulation_jump.go` (new)
- **What:** Jump, call, return, and RST instruction handlers extracted from emulation.go.
  All include MEMPTR updates
- **Why:** Jump/call/return are a logical group. MEMPTR is set to the target address on
  every jump, call, and return instruction

### `emulation_load.go` (new)
- **What:** 16-bit load instructions (LD rr,nn / LD rr,(nn) / LD (nn),rr) and PUSH/POP
  handlers extracted from emulation.go. All include MEMPTR updates
- **Why:** 16-bit loads are a logical group. LD (nn),rr and LD rr,(nn) update MEMPTR to
  the address + 1

### `emulation_test.go`
- **What:** Updated tests for Bus interface, added MEMPTR verification
- **Why:** Tests must validate MEMPTR tracking

### `instruction.go`
- **What:** Added instruction definitions for newly separated instruction handlers,
  updated instruction function references
- **Why:** Instruction table entries must point to the correct handler functions after
  refactoring

### `instruction_dd.go`
- **What:** Added undocumented DD instruction table entries
- **Why:** Register undocumented DD-prefix opcodes in the opcode table

### `instruction_dd_undoc.go` (new)
- **What:** Instruction definitions for undocumented DD-prefix opcodes
- **Why:** Undocumented IXH/IXL register operations need instruction metadata for
  disassembly support

### `instruction_ed.go`
- **What:** Added new ED instruction entries for extended operations
- **Why:** Complete ED-prefix instruction coverage

### `instruction_fd.go`
- **What:** Added undocumented FD instruction table entries
- **Why:** Same as instruction_dd.go but for IY

### `instruction_fd_undoc.go` (new)
- **What:** Instruction definitions for undocumented FD-prefix opcodes
- **Why:** Same as instruction_dd_undoc.go but for IY

### `opcode.go`
- **What:** Reorganized opcode table entries, updated handler references for refactored
  instruction functions, added entries for newly defined instructions
- **Why:** Opcode tables must reflect the new handler function locations after the
  emulation.go split

### `opcode_test.go`
- **What:** Updated opcode table completeness tests
- **Why:** Tests must pass with reorganized tables

### `param.go`
- **What:** Updated parameter reading to use `c.bus.*` instead of `c.memory.*`
- **Why:** Part of Memory-to-Bus migration

### `unofficial.go`
- **What:** Updated I/O port operations to use `c.bus.ReadPort`/`c.bus.WritePort` with
  16-bit addresses; general cleanup
- **Why:** I/O operations must use the Bus interface for correct 16-bit port addressing

### `unofficial_test.go`
- **What:** Updated tests for Bus interface
- **Why:** Tests must work with new port operation signatures

### `singlestep_test.go` (new)
- **What:** Integration test runner for SingleStepTests/z80 JSON test vectors
- **Why:** Bus-level validation of every Z80 opcode against the Ares Z80 core reference

### `zexall_test.go` (new)
- **What:** Integration test runner for ZEXDOC (67 tests) and ZEXALL (67 tests) exercisers
- **Why:** ZEXALL validates all undocumented flag behavior including MEMPTR and Q register
  effects. 100% pass rate (67/67 for both) confirms correctness

---

## Documentation (`docs/`)

### `cpu-implementation-plan-65816.md` (new)
- 65816 CPU implementation plan. Status: COMPLETE for CPU emulation scope

### `cpu-implementation-plan-65c02-68000.md` (new)
- 65C02 and 68000 implementation plan. Status: COMPLETE for both CPUs

### `m68000-emulator-comparison.md` (new)
- Comparison of retrogolib 68000 vs Musashi, Moira, MAME, Cyclone, m68k-rs, r68k

### `m68000-gap-closure-plan.md` (new)
- 68000 accuracy gap closure: address errors, exception frames, bus errors, timing

### `supported-systems.md` (new)
- Overview of all implemented CPUs and systems with gap analysis

### `system-implementation-plan-atari2600.md` (new)
- Atari 2600 implementation plan: 6507 variant + TIA/RIOT system package

### `system-implementation-plan-c64.md` (new)
- Commodore 64 implementation plan: 6510 variant + VIC-II/SID/CIA system package

### `cpu-implementation-plan-sm83.md` (new)
- SM83 CPU implementation plan. Status: COMPLETE for CPU emulation scope

### `system-implementation-plan-gameboy.md` (new)
- Game Boy implementation plan: LR35902 CPU package + GB system package

### `z80-emulator-comparison.md` (new)
- Comparison of retrogolib Z80 vs Go-based Z80 emulators

### `z80-emulator-comparison-cross-language.md` (new)
- Comparison of retrogolib Z80 vs C/C++/Rust Z80 emulators

### `z80-gap-closure-plan.md` (new)
- Z80 gap closure: Bus interface, IM 0/2, RETI, LD A,{I|R} bug, ED mirrors

---

## Sharp SM83 (`arch/cpu/sm83/`) -- All New

Entire package is new. Implements the SM83 (LR35902) CPU used in Game Boy / Game Boy Color.
Architecturally distinct from Z80 — removes shadow registers, IX/IY, I/O instructions, DD/ED/FD
prefixes, and repurposes 14+ opcodes for Game Boy-specific operations.

### `doc.go`
- Package documentation with architecture overview and key Z80 differences

### `addressing.go`
- 7 addressing modes (no PortAddressing), SM83-specific RegisterParam constants
  (RegHLPlus, RegHLMinus, RegHighMem, RegCIndirect, RegSPOffset), string representations

### `instruction.go`
- Instruction struct, OpcodeInfo struct, 44 instruction name constants, ~60 instruction
  variable definitions covering all SM83 opcodes including SM83-unique instructions
  (STOP, SWAP, LDH, LD HL+/-, ADD SP,e, LD HL,SP+e, LD (nn),SP, RETI at 0xD9)

### `opcode.go`
- Opcode struct, 256-entry base opcode table with M-cycle timing, 11 illegal opcode
  slots ($D3, $DB, $DD, $E3, $E4, $EB, $EC, $ED, $F4, $FC, $FD)

### `opcode_cb.go`
- 256-entry CB-prefix opcode table, SWAP at 0x30-0x37 (replaces Z80's SLL)

### `categories.go`
- Instruction category sets: branching, non-returning, memory read/write

### `errors.go`
- Package-specific errors including ErrIllegalOpcode

### `flag.go`
- 4 flags only: Z (bit 7), N (bit 6), H (bit 5), C (bit 4), lower nibble always 0

### `cpu.go`
- CPU state: A,F,B,C,D,E,H,L (8-bit), SP,PC (16-bit), IME, imeDelay, haltBug
- Register pair accessors, stack operations, register value get/set by 3-bit encoding

### `option.go`
- Functional options with Game Boy defaults (PC=0x0100, SP=0xFFFE)

### `memory.go`
- Memory interface (Read/Write/ReadWord/WriteWord), BasicMemory flat 64KB implementation

### `step.go`
- Fetch/decode/execute with CB prefix handling, HALT bug, delayed EI semantics,
  TraceStep support, jump instruction detection

### `param.go`
- Parameter reading for all 7 addressing modes

### `interrupt.go`
- 5 fixed vectors (VBlank $0040, LCD STAT $0048, Timer $0050, Serial $0058, Joypad $0060),
  IME, IE ($FFFF) / IF ($FF0F) register interaction, HALT wake-up and HALT bug activation

### `emulation.go`
- ALU: ADD, ADC, SUB, SBC, AND, OR, XOR, CP, INC, DEC (8-bit and 16-bit), DAA, CPL,
  CCF, SCF, ADD HL,rr, ADD SP,e, NOP, HALT, STOP, DI, EI, rotate accumulator (RLCA/RRCA/RLA/RRA)

### `emulation_load.go`
- LD variants: register, immediate, indirect (BC/DE), HL+/-, LDH ($FF00+n), LD (C),
  LD (nn),A / LD A,(nn), LD (nn),SP, LD SP,HL, LD HL,SP+e, LD (HL),n, PUSH, POP

### `emulation_branch.go`
- JP, JR (absolute/relative, conditional/unconditional), CALL, RET (conditional/unconditional),
  RETI, RST, condition checking helper

### `emulation_cb.go`
- RLC, RRC, RL, RR, SLA, SRA, SWAP, SRL, BIT, RES, SET with (HL) indirect support

### `singlestep_test.go` (new)
- Integration test runner for SingleStepTests/sm83 JSON test vectors

---

## Intel x86 Real Mode (`arch/cpu/x86/`) -- All New

Entire package is new. Implements Intel 8086/8088 through 80486 in real mode for static analysis and tooling development.

### `addressing.go`
- ModR/M byte decoder with 16-bit addressing modes (DIB, SI, DI, BX, BP, SP, direct)
- Segmented addressing: segment:offset to 20-bit physical address calculation
- Effective address computation for all ModR/M combinations

### `categories.go`
- Instruction category sets: branching, non-returning, memory read/write
- Covers all instruction types from 8086 base through 80486 extensions

### `cpu.go`
- CPU state: AX, BX, CX, DX, SP, BP, SI, DI (16-bit registers)
- Segment registers: CS, DS, ES, SS
- Flags register with all x86 flags (CF, PF, AF, ZF, SF, TF, IF, DF, OF)
- 20-bit physical address space (1MB real mode)
- Thread-safe register access with mutex protection

### `cpu_test.go`
- Unit tests for register access, segment calculations, flag operations

### `doc.go`
- Architecture overview, real mode memory model, usage examples
- Documents supported CPU generations (8086 through 80486)

### `errors.go`
- Package-specific errors for invalid addressing modes, bounds violations

### `flag.go`
- 16-bit flags register with individual bit constants
- Helper methods for flag access and manipulation

### `instruction.go`
- ~200 instruction definitions with opcode mappings
- Covers 8086 base, 186/286 extensions, 386 protected mode instructions (in real mode)
- Includes 0x0F prefix opcodes for 386+/486 instructions

### `instruction_test.go`
- Comprehensive instruction decode tests
- ModR/M byte interpretation validation
- Instruction size and timing verification

### `instructions.go`
- Instruction execution handlers for all opcodes
- Data movement: MOV, XCHG, XLAT, LEA, LDS, LES, PUSH, POP, PUSHA, POPA, PUSHF, POPF
- Arithmetic: ADD, ADC, SUB, SBC, MUL, IMUL, DIV, IDIV, INC, DEC, NEG, CMP
- Logical: AND, OR, XOR, NOT, TEST
- Shift/rotate: SHL, SHR, SAL, SAR, ROL, ROR, RCL, RCR, SHLD, SHRD
- Control flow: JMP, CALL, RET, IRET, INT, INTO, BOUND, LOOP, LOOPE, LOOPNE, JCXZ
- String: MOVSB/W, CMPSB/W, SCASB/W, LODSB/W, STOSB/W (with REP/REPE/REPNE)
- Bit operations: BSF, BSR, BT, BTC, BTR, BTS
- Data movement extensions: MOVZX, MOVSX, XADD, CMPXCHG, BSWAP
- System: SMSW, LMSW, CLC, STC, CMC, CLD, STD, CLI, STI, CLTS, LOCK

### `memory.go`
- Segmented memory model with segment:offset to physical address translation
- 1MB addressable space (20-bit addresses)
- Read/Write methods for byte, word, dword
- Memory region validation and bounds checking
- Hex dump utility for debugging

### `memory_test.go`
- Segment:offset calculation tests
- Bounds checking validation
- Mirror behavior verification

### `opcode.go`
- 256-entry base opcode table with instruction metadata
- 0x0F extended opcode table for 386+ instructions
- Opcode classification by type and CPU generation

### `opcode_id.go` (new)
- `OpcodeID` type with 256+ constants for type-safe opcode references
- Compile-time safety for opcode lookups

---

## Z80 Undocumented Instructions & Bus Interface (`arch/cpu/z80/`)

### `emulation_dd_undoc.go` (new)
- **What:** Undocumented DD-prefix instruction handlers (IXH/IXL arithmetic, DD+CB indexed operations)
- **Why:** Clean separation of undocumented behavior from standard instructions

### `emulation_fd_undoc.go` (new)
- **What:** Undocumented FD-prefix instruction handlers (IYH/IYL arithmetic)
- **Why:** Clean separation of undocumented behavior

### `emulation_helpers.go` (new)
- **What:** Shared arithmetic and flag helpers extracted from emulation.go
- **Why:** Reduces file size, avoids duplication between base and prefix instruction files

### `emulation_index.go` (new)
- **What:** Shared IX/IY instruction implementations (LD, ADD, ADC, SBC, INC, DEC, etc.)
- **Why:** DD and FD prefix instructions are identical except for which index register they use

### `emulation_jump.go` (new)
- **What:** Jump, call, return, and RST instruction handlers with MEMPTR updates
- **Why:** Jump/call/return are a logical group; MEMPTR is set to the target address

### `emulation_load.go` (new)
- **What:** 16-bit load instructions (LD rr,nn / LD rr,(nn) / LD (nn),rr) and PUSH/POP handlers
- **Why:** 16-bit loads are a logical group; MEMPTR is updated on memory accesses

### `instruction_dd_undoc.go` (new)
- **What:** Instruction definitions for undocumented DD-prefix opcodes
- **Why:** Undocumented IXH/IXL register operations need instruction metadata for disassembly

### `instruction_fd_undoc.go` (new)
- **What:** Instruction definitions for undocumented FD-prefix opcodes
- **Why:** Same as instruction_dd_undoc.go but for IY

### `opcode_id.go` (new)
- **What:** Added `OpcodeID` type with 256 constants for type-safe opcode references
- **Why:** Provides compile-time safety and IDE autocomplete for opcode lookups

### `singlestep_test.go` (new)
- **What:** Integration test runner for SingleStepTests/z80 JSON test vectors
- **Why:** Bus-level validation of every Z80 opcode against the Ares Z80 core reference

### `zexall_test.go` (new)
- **What:** Integration test runner for ZEXDOC (67 tests) and ZEXALL (67 tests) exercisers
- **Why:** ZEXALL validates all undocumented flag behavior including MEMPTR and Q register effects

---

## Atari 2600 System (`arch/system/atari2600/`) -- All New

### `atari2600.go` (new)
- **What:** Memory map constants: TIA ranges, RAM (128 bytes), RIOT I/O, ROM window, cartridge sizes, reset vector
- **Why:** Defines Atari 2600 (VCS) memory layout for system emulation

### `atari2600_test.go` (new)
- **What:** Tests for address constants, cartridge sizes, reset vector, 13-bit address masking, TIA mirroring
- **Why:** Validates memory map correctness and mirroring behavior

### `cartridge/cartridge.go` (new)
- **What:** BankingScheme type (None, F8, FA, F6, F4, 3F), Load function, size-based scheme detection
- **Why:** Handles various cartridge bank switching schemes used in Atari 2600 games

### `cartridge/cartridge_test.go` (new)
- **What:** 14 tests for scheme detection, ROM loading, bank offsets, trigger address mapping
- **Why:** Validates cartridge detection and bank switching behavior

### `doc.go` (new)
- **What:** Package documentation with memory map overview
- **Why:** Documents Atari 2600 system architecture

### `register/register_test.go` (new)
- **What:** Register completeness tests (TIA write: 45, TIA read: 14, RIOT: 10), address range validation
- **Why:** Ensures all TIA and RIOT registers are properly defined

### `register/riot.go` (new)
- **What:** 10 RIOT registers ($0280-$0297), timer interval constants, console switch bits (SWCHB), joystick direction bits (SWCHA)
- **Why:** Defines RIOT I/O chip registers for Atari 2600

### `register/tia.go` (new)
- **What:** All 45 TIA write registers ($00-$2C) and 14 TIA read registers ($00-$0D)
- **Why:** Defines TIA video/audio registers for Atari 2600

---

## Additional Documentation

### `supported-systems.md` (new)
- Overview of all implemented CPUs and systems with gap analysis

### `system-implementation-plan-c64.md` (new)
- Commodore 64 implementation plan: 6510 variant + VIC-II/SID/CIA system package

### `system-implementation-plan-gameboy.md` (new)
- Game Boy implementation plan: LR35902 CPU package + GB system package

---

## Test Infrastructure

### `testdata/Makefile` (new)
- **What:** Build automation for downloading and preparing test ROMs (Dormann, SingleStepTests)
- **Why:** Standardizes test data preparation across CPU packages

### `testdata/.gitignore` (new)
- **What:** Git ignore rules for test data directories
- **Why:** Prevents large test ROMs from being committed to repository

---

## Architecture Registration Updates

### `arch/arch.go`
- **What:** Added architecture constants: `M65C02`, `M65816`, `M68000`, `SM83`, `X86`
- **Why:** Registers five new CPU architectures for validation and lookup

### `arch/arch_test.go`
- **What:** Updated test count expectation for new architectures
- **Why:** Keeps the architecture completeness test passing

### `arch/system.go`
- **What:** Added system constants: `AppleIIGS`, `Atari2600`, `SNES`
- **Why:** Registers SNES, Apple IIGS, and Atari 2600 as supported systems

### `arch/system_test.go`
- **What:** Updated test expectations for new systems
- **Why:** Keeps the system completeness test passing

### `chip8/`
- **What:** Added `opcode_id.go` with OpcodeID type
- **Why:** Consistency across CPU packages for type-safe opcode references
