# Work Branch Changelog

Tracks the committed changes introduced by `work2` since its common ancestor
with the current remote-tracking `main` branch.

**Last Updated:** 2026-09-01

---

## Current Branch State

- Comparison range: `origin/main...HEAD`, currently
  `e7fb46b...3751fde`.
- Merge base: `e7fb46b` (`readme: remove goreportcard`, 2026-08-29).
- Committed branch delta: 185 files, with 150 added and 35 modified; 27,234
  insertions and 1,283 deletions.
- `git diff --name-status --find-renames origin/main...HEAD` reports no deleted
  or renamed files in the current range.
- The range covers committed changes only; the working-tree edit that refreshes
  this document is intentionally not included in those statistics.

## Changes Already Absorbed From `main`

The merge base already contains earlier branch work that is therefore absent
from the current branch-side diff:

- The base CHIP-8, 6502/65C02, x86, and Z80 implementations and their initial
  architecture registration.
- Earlier shared test infrastructure and test data support; this branch only
  contributes the CPU targets and CHIP-8 additions described below.
- Earlier CLI, configuration, and logging cleanups.

## Branch-Specific Changes

### New CPU Emulators

- **WDC 65C816:** Adds `arch/cpu/cpu65816/` with emulation and native modes,
  24-bit addressing, interrupts, tracing, structured instruction metadata, and
  unit, opcode, effective-address, and SingleStepTests coverage.
- **Motorola 68000:** Adds `arch/cpu/cpu68000/` with bus and memory abstractions,
  effective-address handling, interrupts, tracing, structured instruction
  metadata, and unit, opcode, memory, and SingleStepTests coverage.
- **Motorola 6809:** Adds `arch/cpu/cpu6809/` with indexed addressing, page
  0/10/11 opcode tables, interrupts, tracing, structured instruction metadata,
  and unit and opcode coverage.
- **Sharp SM83:** Adds `arch/cpu/sm83/` with Game Boy CPU behavior, CB-prefixed
  instructions, interrupts, tracing, structured instruction metadata, and
  SingleStepTests coverage.

### Existing CPU Refactors and Correctness

- **Cross-architecture structure:** Aligns instruction names, instruction
  registries, option handling, constructor setup, and focused tests across the
  CHIP-8, 6502, 65C816, 68000, 6809, SM83, x86, and Z80 packages. Small internal
  declarations remain with their cohesive implementation instead of occupying
  single-purpose files.
- **Public APIs:** Replaces exposed option implementation details with private
  option state and typed functional options. Obsolete compatibility APIs are
  removed rather than retained as deprecated wrappers.
- **CHIP-8:** Makes COSMAC VIP behavior the default and adds explicit quirks for
  later interpreter variants. Corrects display-wait, key-release, shift, logic,
  jump, load/store, sprite, arithmetic-flag, and bounds behavior; expands unit
  tests and adds Timendus ROM conformance tests with screenshot comparison.
- **6502/65C02:** Reorganizes instruction metadata and unofficial operations,
  tightens constructor, memory, interrupt, stepping, and option behavior, and
  adds focused CPU, interrupt, option, opcode, and step tests.
- **x86 and Z80:** Separates substantial instruction-name and registry data from
  core instruction definitions. Z80 also adopts private option state, a typed
  pre-execution hook, nil-option handling, and registry and option tests.

### System Foundations

- **Atari 2600:** Adds system and memory-map definitions, TIA and RIOT register
  constants, cartridge loading, supported banking schemes, and tests under
  `arch/system/atari2600/`.
- **TRS-80 Color Computer:** Adds CoCo memory-map, hardware, interrupt, PIA, and
  SAM definitions with tests under `arch/system/coco/`.
- **Vectrex:** Adds Vectrex memory-map, hardware, interrupt, and VIA definitions
  with tests under `arch/system/vectrex/`.

### Test and Documentation Integration

- Extends the root integration-test target for the 65C816, 68000, and SM83
  suites.
- Adds a `testdata/Makefile` target for the Timendus CHIP-8 test suite and keeps
  the existing CPU test-data targets in the aggregate workflow.
- Updates `README.md` with the new CPU packages.
- Adds gap-closure plans for the Motorola 68000 and Z80, a Commodore 64 system
  implementation plan, and this branch changelog.

## Files

| Status | Count | Files | Purpose |
| --- | ---: | --- | --- |
| Modified | 1 | `Makefile` | Extends CPU integration-test coverage. |
| Modified | 1 | `README.md` | Lists the new CPU packages. |
| Added / Modified | 4 / 6 | `arch/cpu/chip8/` | Adds compatibility options, registry separation, correctness fixes, and Timendus ROM tests. |
| Added / Modified | 7 / 17 | `arch/cpu/cpu6502/` | Refactors package structure and behavior and expands focused tests. |
| Added | 30 | `arch/cpu/cpu65816/` | Adds the WDC 65C816 emulator and tests. |
| Added | 29 | `arch/cpu/cpu68000/` | Adds the Motorola 68000 emulator and tests. |
| Added | 27 | `arch/cpu/cpu6809/` | Adds the Motorola 6809 emulator and tests. |
| Added | 24 | `arch/cpu/sm83/` | Adds the Sharp SM83 emulator and tests. |
| Added / Modified | 2 / 2 | `arch/cpu/x86/` | Moves instruction names and registry data into cohesive files. |
| Added / Modified | 4 / 7 | `arch/cpu/z80/` | Aligns registry and option structure and adds focused tests. |
| Added | 19 | `arch/system/atari2600/`, `arch/system/coco/`, `arch/system/vectrex/` | Adds system foundations, register definitions, cartridge support, and tests. |
| Added | 4 | `docs/` | Adds implementation plans and branch tracking. |
| Modified | 1 | `testdata/Makefile` | Integrates the Timendus CHIP-8 ROM suite. |

The grouped counts above total the exact 150 added and 35 modified files
reported for `origin/main...HEAD`; no row represents a rename.

## Merge Summary

The current branch delta combines four additive CPU emulators and three system
foundations with a consistency and correctness pass over every pre-existing CPU
architecture. The largest changes to existing code are the CHIP-8 and 6502
behavioral refactors; x86 and Z80 changes primarily align instruction metadata
and option organization.

## Verification

- Refreshed `origin/main` before deriving the comparison.
- Derived the merge base, totals, and file classifications from
  `origin/main...HEAD`, using `--find-renames` for the Files table.
- Verified the final document with `git diff --check` and rechecked its stated
  range and statistics against the live committed diff.
