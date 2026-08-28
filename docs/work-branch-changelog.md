# Work Branch Changelog

Tracks the changes introduced by `work2` since its common ancestor with `main`.

**Last Updated:** 2026-08-28

---

## Current Branch State

- Comparison of local refs: `main...work2` (`ea7c09d...feac27c`).
- Merge base: `316c334`; this revision of `main` was merged into `work2` by
  `f6f724a` on 2026-05-03. The earlier `6042731` and `85d794f` merges are also
  part of the branch history.
- Delta: 126 files, with 115 added and 11 modified; 24,493 insertions and 91
  deletions.
- The source of truth for this document is the current `main...work2` diff, not
  older merge-plan assumptions.

## Changes Already Absorbed From `main`

The following earlier `work2` areas are absent from the branch-side diff because
they were present at the latest merge base:

- Test infrastructure changes in `Makefile` and `testdata/`.
- Architecture registration updates in `arch/`.
- 6502/65C02 implementation and integration-test work, apart from the small
  housekeeping changes listed below.
- Most Z80 bus, opcode, and integration-test work, apart from the four-file
  documentation and declaration-order cleanup listed below.
- x86 support in `arch/cpu/x86/`.

## Branch-Specific Changes

### CPU Emulators

- **WDC 65C816:** Adds the `arch/cpu/m65816/` package, including
  emulation/native modes, 24-bit memory support, interrupts, tracing, and unit,
  opcode, and SingleStepTests coverage.
- **Motorola 68000:** Adds the `arch/cpu/m68000/` package, including
  bus and memory abstractions, effective-address handling, interrupts, tracing,
  and unit, opcode, memory, and SingleStepTests coverage.
- **Motorola 6809:** Adds the `arch/cpu/m6809/` package, including
  indexed addressing, page 0/10/11 opcode tables, interrupts, tracing, and unit
  and opcode coverage.
- **Sharp SM83:** Adds the `arch/cpu/sm83/` package for Game Boy and
  Game Boy Color CPU behavior, including CB-prefixed instructions, interrupts,
  tracing, and SingleStepTests coverage.

### System Foundations

- **Atari 2600:** Adds system and memory-map definitions, TIA and RIOT register
  constants, and cartridge loading with supported banking schemes under
  `arch/system/atari2600/`, with tests.
- **TRS-80 Color Computer:** Adds CoCo memory-map, hardware, interrupt, PIA, and
  SAM definitions under `arch/system/coco/`, with tests.
- **Vectrex:** Adds Vectrex memory-map, hardware, interrupt, and VIA definitions
  under `arch/system/vectrex/`, with tests.

### Logging

- Adds lazy `Int64Func`, `Float64Func`, `BoolFunc`, `DurationFunc`, and
  `StringerFunc` field constructors, including nil handling for lazy stringers.
- Makes `Err` and `Type` produce explicit `<nil>` values, replaces integer hex
  formatting with a fixed-width uppercase formatter, and tests the added and
  revised field behavior.
- Passes fields to `slog.Record.AddAttrs` directly instead of allocating and
  populating an intermediate `[]any` slice.

### Small Follow-Up Cleanups

- `arch/cpu/m6502/option.go` and `arch/cpu/m6502/singlestep_test.go` move a
  private hook declaration and normalize test log capitalization.
- `arch/cpu/z80/categories.go`, `instruction.go`, `opcode_test.go`, and
  `option.go` clean up comments and test labels and move a private hook
  declaration; they do not change opcode behavior.
- `cli/flags.go` and `config/config.go` move private type declarations closer to
  their owning types without changing their data shapes.

### Documentation

- Adds active gap-closure plans for the Motorola 68000 and Z80.
- Adds the planned Commodore 64 system implementation document.
- Maintains this branch changelog.

## Files

| Status | Files | Purpose |
| --- | --- | --- |
| Added | `arch/cpu/m65816/`, `arch/cpu/m68000/`, `arch/cpu/m6809/`, `arch/cpu/sm83/` | New CPU emulator packages and tests. |
| Added | `arch/system/atari2600/`, `arch/system/coco/`, `arch/system/vectrex/` | New system, cartridge, memory-map, and register definitions with tests. |
| Modified | `log/field.go`, `log/field_test.go`, `log/logger.go` | Expanded lazy fields, revised formatting and nil handling, and reduced record-building allocation. |
| Modified | `arch/cpu/m6502/`, `arch/cpu/z80/`, `cli/flags.go`, `config/config.go` | Non-behavioral declaration, comment, label, and capitalization cleanup. |
| Added | `docs/m68000-gap-closure-plan.md`, `docs/system-implementation-plan-c64.md`, `docs/z80-gap-closure-plan.md`, `docs/work-branch-changelog.md` | Active plans and branch tracking. |

## Merge Summary

Most of the remaining branch delta is additive CPU and system-package work. The
logging changes are the main behavioral modification to an existing package;
the m6502, Z80, CLI, and configuration differences are small cleanups.

## Verification

- Inspected: `git status --short`, `git diff --stat main...work2`,
  `git diff --name-status main...work2`, and the substantive branch-side diffs.
- Not run: Go tests; this update changes documentation only.
