# Work Branch Changelog

Tracks the changes introduced by `work2` since its common ancestor with `main`.

**Last Updated:** 2026-08-29

---

## Current Branch State

- Comparison of current remote-tracking refs: `origin/main...work2`
  (`e7fb46b...af3c5fd`).
- Merge base: `e7fb46b`; this revision of `main` was merged into `work2` by
  `af3c5fd` on 2026-08-29. Earlier merges remain part of the branch history.
- Committed baseline delta: 119 files, with 118 added and 1 modified; 24,259
  insertions and 5 deletions.
- The source of truth for this document is the current `origin/main...work2`
  diff plus the branch history, not older merge-plan assumptions.

## Changes Already Absorbed From `main`

The following earlier `work2` areas are absent from the branch-side diff because
they were present at the latest merge base:

- Test infrastructure changes in `Makefile` and `testdata/`.
- Architecture registration updates in `arch/`.
- 6502/65C02 implementation, package rename, and integration-test work.
- Most Z80 bus, opcode, and integration-test work, apart from the three-file
  documentation and declaration-order cleanup listed below.
- x86 support in `arch/cpu/x86/`.
- CLI and configuration declaration-order cleanup.
- Logging safety, formatting, and handler-consistency changes under `log/`.

## Branch-Specific Changes

### CPU Emulators

- **WDC 65C816:** Adds the `arch/cpu/cpu65816/` package, including
  emulation/native modes, 24-bit memory support, interrupts, tracing, and unit,
  opcode, and SingleStepTests coverage.
- **Motorola 68000:** Adds the `arch/cpu/cpu68000/` package, including
  bus and memory abstractions, effective-address handling, interrupts, tracing,
  and unit, opcode, memory, and SingleStepTests coverage.
- **Motorola 6809:** Adds the `arch/cpu/cpu6809/` package, including
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

### Small Follow-Up Cleanups

- `arch/cpu/cpu68000/instruction.go` and `arch/cpu/sm83/option.go` order private
  function types before their exported dependents to satisfy lint rules.
- `arch/cpu/z80/categories.go`, `instruction.go`, and `opcode_test.go` clean up
  comments and test labels without changing opcode behavior.

### Documentation

- Updates `README.md` to list the newly added CPU packages.
- Adds active gap-closure plans for the Motorola 68000 and Z80.
- Adds the planned Commodore 64 system implementation document.
- Maintains this branch changelog.

## Files

| Status | Files | Purpose |
| --- | --- | --- |
| Modified | `README.md` | Lists the newly added CPU packages. |
| Added | `arch/cpu/cpu65816/`, `arch/cpu/cpu68000/`, `arch/cpu/cpu6809/`, `arch/cpu/sm83/` | New CPU emulator packages and tests relative to the merge base. |
| Added | `arch/system/atari2600/`, `arch/system/coco/`, `arch/system/vectrex/` | New system, cartridge, memory-map, and register definitions with tests. |
| Modified | `arch/cpu/cpu68000/instruction.go`, `arch/cpu/sm83/option.go`, `arch/cpu/z80/` | Non-behavioral declaration, comment, and label cleanup. |
| Added | `docs/cpu68000-gap-closure-plan.md`, `docs/system-implementation-plan-c64.md`, `docs/z80-gap-closure-plan.md`, `docs/work-branch-changelog.md` | Active plans and branch tracking. |

## Merge Summary

Most of the remaining branch delta is additive CPU and system-package work. The
only remaining changes to existing files are the README package list and the
small Z80 cleanups.

## Verification

- Inspected: `git status --short`, `git diff --stat origin/main...work2`,
  `git diff --name-status origin/main...work2`, and the substantive
  branch-side diffs.
- Passed: `go fmt ./...`.
- Passed: `make lint`.
- Passed: `make test`.
